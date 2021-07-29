package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

const (
	// CoverLinkingStatusStarted ...
	CoverLinkingStatusStarted = "coverlinking started"
)

// SignUpUseCases represents all the business logic involved in setting up a user
type SignUpUseCases interface {
	// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
	VerifyPhoneNumber(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error)

	// creates an account for the user, setting the provided phone number as the PRIMARY PHONE
	// NUMBER
	CreateUserByPhone(ctx context.Context, input *dto.SignUpInput) (*profileutils.UserResponse, error)

	// updates the user profile of the currently logged in user
	UpdateUserProfile(
		ctx context.Context,
		input *dto.UserProfileInput,
	) (*profileutils.UserProfile, error)

	// adds a new push token in the users profile if the push token does not exist
	RegisterPushToken(ctx context.Context, token string) (bool, error)

	// called to create a customer account in the ERP. This API is only valid for `BEWELL CONSUMER`
	// it should be the last call after updating the users bio data. Its should not return an error
	// when it fails due to unreachable errors, rather it should retry
	CompleteSignup(ctx context.Context, flavour feedlib.Flavour) (bool, error)

	// removes a push token from the users profile
	RetirePushToken(ctx context.Context, token string) (bool, error)

	// fetches the phone numbers of a user for the purposes of recoverying an account.
	// the returned phone numbers should be masked
	GetUserRecoveryPhoneNumbers(
		ctx context.Context,
		phoneNumber string,
	) (*dto.AccountRecoveryPhonesResponse, error)

	// called to set the provided phone number as the PRIMARY PHONE NUMBER in the user profile of
	// the user
	// where the phone number is associated with.
	SetPhoneAsPrimary(ctx context.Context, phone, otp string) (bool, error)

	RemoveUserByPhoneNumber(ctx context.Context, phone string) error
}

// SignUpUseCasesImpl represents usecase implementation object
type SignUpUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	profileUsecase       ProfileUseCase
	pinUsecase           UserPINUseCases
	supplierUsecase      SupplierUseCases
	baseExt              extension.BaseExtension
	engagement           engagement.ServiceEngagement
	pubsub               pubsubmessaging.ServicePubSub
	edi                  edi.ServiceEdi
}

// NewSignUpUseCases returns a new a onboarding usecase
func NewSignUpUseCases(
	r repository.OnboardingRepository,
	profile ProfileUseCase,
	pin UserPINUseCases,
	supplier SupplierUseCases,
	ext extension.BaseExtension,
	eng engagement.ServiceEngagement,
	pubsub pubsubmessaging.ServicePubSub,
	edi edi.ServiceEdi,
) SignUpUseCases {
	return &SignUpUseCasesImpl{
		onboardingRepository: r,
		profileUsecase:       profile,
		pinUsecase:           pin,
		supplierUsecase:      supplier,
		baseExt:              ext,
		engagement:           eng,
		pubsub:               pubsub,
		edi:                  edi,
	}
}

// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
func (s *SignUpUseCasesImpl) VerifyPhoneNumber(
	ctx context.Context,
	phone string,
	appID *string,
) (*profileutils.OtpResponse, error) {
	ctx, span := tracer.Start(ctx, "VerifyPhoneNumber")
	defer span.End()

	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}
	// check if phone number exists
	exists, err := s.profileUsecase.CheckPhoneExists(ctx, *phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	// if phone exists return early
	if exists {
		return nil, exceptions.CheckPhoneNumberExistError()
	}
	// generate and send otp to the phone number
	otpResp, err := s.engagement.GenerateAndSendOTP(ctx, *phoneNumber, appID)

	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.GenerateAndSendOTPError(err)
	}
	// return the generated otp
	return otpResp, nil
}

// CreateUserByPhone creates an account for the user, setting the provided phone number as the
// PRIMARY PHONE NUMBER
func (s *SignUpUseCasesImpl) CreateUserByPhone(
	ctx context.Context,
	input *dto.SignUpInput,
) (*profileutils.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "CreateUserByPhone")
	defer span.End()

	userData, err := utils.ValidateSignUpInput(input)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	verified, err := s.engagement.VerifyOTP(
		ctx,
		*userData.PhoneNumber,
		*userData.OTP,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.VerifyOTPError(err)
	}

	if !verified {
		return nil, exceptions.VerifyOTPError(nil)
	}

	// get or create user via their phone number
	user, err := s.onboardingRepository.GetOrCreatePhoneNumberUser(ctx, *userData.PhoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// create a user profile
	profile, err := s.onboardingRepository.CreateUserProfile(
		ctx,
		*userData.PhoneNumber,
		user.UID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	// generate auth credentials
	auth, err := s.onboardingRepository.GenerateAuthCredentials(
		ctx,
		*userData.PhoneNumber,
		profile,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	// save the user pin
	_, err = s.pinUsecase.SetUserPIN(
		ctx,
		*userData.PIN,
		profile.ID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	var supplier *profileutils.Supplier
	var customer *profileutils.Customer
	supplier, err = s.onboardingRepository.CreateEmptySupplierProfile(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	customer, err = s.onboardingRepository.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	// set the user default communications settings
	defaultCommunicationSetting := true
	comms, err := s.onboardingRepository.SetUserCommunicationsSettings(
		ctx,
		profile.ID,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	contact := domain.CRMContact{
		Properties: domain.ContactProperties{
			Phone:                 *profile.PrimaryPhone,
			FirstChannelOfContact: domain.ChannelOfContactApp,
			BeWellEnrolled:        domain.GeneralOptionTypeYes,
			BeWellAware:           domain.GeneralOptionTypeYes,
		},
	}

	if err = s.pubsub.NotifyCreateContact(ctx, contact); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to crm.contact.create topic: %v", err)
	}

	navActions, err := s.profileUsecase.GetNavActions(ctx, *profile)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return &profileutils.UserResponse{
		Profile:               profile,
		SupplierProfile:       supplier,
		CustomerProfile:       customer,
		CommunicationSettings: comms,
		Auth:                  *auth,
		NavActions:            navActions,
	}, nil
}

// UpdateUserProfile  updates the user profile of the currently logged in user
func (s *SignUpUseCasesImpl) UpdateUserProfile(
	ctx context.Context,
	input *dto.UserProfileInput,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "UpdateUserProfile")
	defer span.End()

	// get the old user profile
	pr, err := s.profileUsecase.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if input.PhotoUploadID != nil {
		if err := s.profileUsecase.UpdatePhotoUploadID(ctx, *input.PhotoUploadID); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}
	}

	if err := s.profileUsecase.UpdateBioData(ctx, profileutils.BioData{
		FirstName: func(n *string) *string {
			if n != nil {
				return n
			}
			return pr.UserBioData.FirstName
		}(input.FirstName),
		LastName: func(n *string) *string {
			if n != nil {
				return n
			}
			return pr.UserBioData.LastName
		}(input.LastName),
		DateOfBirth: func(n *scalarutils.Date) *scalarutils.Date {
			if n != nil {
				return n
			}
			return pr.UserBioData.DateOfBirth
		}(input.DateOfBirth),
		Gender: func(n *enumutils.Gender) enumutils.Gender {
			if n != nil {
				return *n
			}
			return pr.UserBioData.Gender
		}(input.Gender),
	}); err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return s.profileUsecase.UserProfile(ctx)
}

// RegisterPushToken adds a new push token in the users profile if the push token does not exist
func (s *SignUpUseCasesImpl) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	ctx, span := tracer.Start(ctx, "RegisterPushToken")
	defer span.End()

	if len(token) < 5 {
		return false, exceptions.InValidPushTokenLengthError()
	}
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, false); err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// CompleteSignup called to create a customer account in the ERP. This API is only valid for `BEWELL
// CONSUMER`
func (s *SignUpUseCasesImpl) CompleteSignup(
	ctx context.Context,
	flavour feedlib.Flavour,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "CompleteSignup")
	defer span.End()

	if flavour != feedlib.FlavourConsumer {
		return false, exceptions.InvalidFlavourDefinedError()
	}

	profile, err := s.profileUsecase.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, err
	}
	if len(profile.PushTokens) > 0 {
		logrus.Printf("This piece of code was called")
		coverLinkingDetails := dto.LinkCoverPubSubMessage{
			PhoneNumber: *profile.PrimaryPhone,
			UID:         uid,
			PushToken:   profile.PushTokens,
		}

		logrus.Printf("Publishing to covers.link topic")
		if err := s.pubsub.NotifyCoverLinking(ctx, coverLinkingDetails); err != nil {
			utils.RecordSpanError(span, err)
			log.Printf("failed to publish to covers.link topic: %v", err)
		}

		currentTime := time.Now()
		coverLinkingEvent := &dto.CoverLinkingEvent{
			ID:                    uuid.NewString(),
			CoverLinkingEventTime: &currentTime,
			CoverStatus:           CoverLinkingStatusStarted,
			PhoneNumber:           *profile.PrimaryPhone,
		}

		if _, err := s.onboardingRepository.SaveCoverAutolinkingEvents(ctx, coverLinkingEvent); err != nil {
			utils.RecordSpanError(span, err)
			log.Printf("failed to save coverlinking `started` event: %v", err)
		}

	}

	if profile.UserBioData.FirstName == nil || profile.UserBioData.LastName == nil {
		return false, exceptions.CompleteSignUpError(nil)
	}
	fullName := fmt.Sprintf("%v %v",
		*profile.UserBioData.FirstName,
		*profile.UserBioData.LastName,
	)

	err = s.supplierUsecase.CreateCustomerAccount(
		ctx,
		fullName,
		profileutils.PartnerTypeConsumer,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		logrus.Printf("failed to create customer account with error: %v", err)
	}

	return true, nil
}

// RetirePushToken removes a push token from the users profile
func (s *SignUpUseCasesImpl) RetirePushToken(ctx context.Context, token string) (bool, error) {
	ctx, span := tracer.Start(ctx, "RetirePushToken")
	defer span.End()

	if len(token) < 5 {
		return false, exceptions.InValidPushTokenLengthError()
	}
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, true); err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}

// GetUserRecoveryPhoneNumbers fetches the phone numbers of a user for the purposes of recoverying
// an account.
func (s *SignUpUseCasesImpl) GetUserRecoveryPhoneNumbers(
	ctx context.Context,
	phone string,
) (*dto.AccountRecoveryPhonesResponse, error) {
	ctx, span := tracer.Start(ctx, "GetUserRecoveryPhoneNumbers")
	defer span.End()

	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	pr, err := s.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	// cherrypick the phone numbers and mask them
	phones := func(p *profileutils.UserProfile) []string {
		phs := []string{}
		phs = append(phs, *p.PrimaryPhone)
		phs = append(phs, p.SecondaryPhoneNumbers...)
		return phs

	}(pr)
	masked := s.profileUsecase.MaskPhoneNumbers(phones)
	return &dto.AccountRecoveryPhonesResponse{
		MaskedPhoneNumbers:   masked,
		UnMaskedPhoneNumbers: phones,
	}, nil
}

// SetPhoneAsPrimary called to set the provided phone number as the PRIMARY PHONE NUMBER in the user
// profile of the user
// where the phone number is associated with.
func (s *SignUpUseCasesImpl) SetPhoneAsPrimary(
	ctx context.Context,
	phone, otp string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "SetPhoneAsPrimary")
	defer span.End()

	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.NormalizeMSISDNError(err)
	}

	err = s.profileUsecase.SetPrimaryPhoneNumber(ctx, *phoneNumber, otp, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// RemoveUserByPhoneNumber removes the record of a user using the provided phone number. This method
// will ONLY be called
// in testing environment.
func (s *SignUpUseCasesImpl) RemoveUserByPhoneNumber(ctx context.Context, phone string) error {
	ctx, span := tracer.Start(ctx, "RemoveUserByPhoneNumber")
	defer span.End()

	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}
	return s.onboardingRepository.PurgeUserByPhoneNumber(ctx, *phoneNumber)
}
