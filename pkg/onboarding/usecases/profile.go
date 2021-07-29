package usecases

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization/permission"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/common"
	"github.com/savannahghi/profileutils"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/cenkalti/backoff"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/crm"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/segmentio/ksuid"
)

const (
	// EmailsAttribute is an attribute that represents
	// a user profile's email addresses
	EmailsAttribute = "emails"

	// PhoneNumbersAttribute is an attribute that represents
	// a user profile's phone numbers
	PhoneNumbersAttribute = "phonenumbers"

	// FCMTokensAttribute is an attribute that represents
	// a user profile's FCM push tokens
	FCMTokensAttribute = "tokens"
)

// VerifyEmailNudgeTitle is the title defined in the `engagement service`
// for the `VerifyEmail` nudge
const VerifyEmailNudgeTitle = "Add Primary Email Address"

var tracer = otel.Tracer("github.com/savannahghi/onboarding/pkg/onboarding/usecases")

// ProfileUseCase represents all the profile business logic
type ProfileUseCase interface {
	// profile related
	UserProfile(ctx context.Context) (*profileutils.UserProfile, error)
	GetProfileByID(ctx context.Context, id *string) (*profileutils.UserProfile, error)
	UpdateUserName(ctx context.Context, userName string) error
	UpdatePrimaryPhoneNumber(ctx context.Context, phoneNumber string, useContext bool) error
	UpdatePrimaryEmailAddress(ctx context.Context, emailAddress string) error
	UpdateSecondaryPhoneNumbers(ctx context.Context, phoneNumbers []string) error
	UpdateSecondaryEmailAddresses(ctx context.Context, emailAddresses []string) error
	UpdateVerifiedIdentifiers(ctx context.Context, identifiers []profileutils.VerifiedIdentifier) error
	UpdateVerifiedUIDS(ctx context.Context, uids []string) error
	UpdateSuspended(ctx context.Context, status bool, phoneNumber string, useContext bool) error
	UpdatePhotoUploadID(ctx context.Context, uploadID string) error
	UpdateCovers(ctx context.Context, covers []profileutils.Cover) error
	UpdatePushTokens(ctx context.Context, pushToken string, retire bool) error
	UpdatePermissions(ctx context.Context, perms []profileutils.PermissionType) error
	AddAdminPermsToUser(ctx context.Context, phone string) error
	RemoveAdminPermsToUser(ctx context.Context, phone string) error
	AddRoleToUser(ctx context.Context, phone string, role profileutils.RoleType) error
	RemoveRoleToUser(ctx context.Context, phone string) error
	UpdateBioData(ctx context.Context, data profileutils.BioData) error
	GetUserProfileByUID(
		ctx context.Context,
		UID string,
	) (*profileutils.UserProfile, error)

	// masks phone number.
	MaskPhoneNumbers(phones []string) []string
	// called to set the primary phone number of a specific profile.
	// useContext is used to mark under which scenario the method is been called.
	SetPrimaryPhoneNumber(
		ctx context.Context,
		phoneNumber string,
		otp string,
		useContext bool,
	) error
	SetPrimaryEmailAddress(
		ctx context.Context,
		emailAddress string,
		otp string,
	) error
	// checks whether a phone number has been registered by another user. Checks both primary and
	// secondary phone numbers. If the the phone number is foreign, it returns false
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)

	// check whether a email has been registered by another user. Checks both primary and
	// secondary emails. If the the phone number is foreign, it returns false
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// called to remove specific secondary phone numbers from the user's profile.'
	RetireSecondaryPhoneNumbers(ctx context.Context, toRemovePhoneNumbers []string) (bool, error)

	// called to remove specific secondary email addresses from the user's profile.
	RetireSecondaryEmailAddress(ctx context.Context, toRemoveEmails []string) (bool, error)

	GetUserProfileAttributes(
		ctx context.Context,
		UIDs []string,
		attribute string,
	) (map[string][]string, error)

	ConfirmedEmailAddresses(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ConfirmedPhoneNumbers(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ValidFCMTokens(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ProfileAttributes(
		ctx context.Context,
		UIDs []string,
		attribute string,
	) (map[string][]string, error)

	SetupAsExperimentParticipant(ctx context.Context, participate *bool) (bool, error)

	AddAddress(
		ctx context.Context,
		input dto.UserAddressInput,
		addressType enumutils.AddressType,
	) (*profileutils.Address, error)

	GetAddresses(ctx context.Context) (*domain.UserAddresses, error)

	GetUserCommunicationsSettings(ctx context.Context) (*profileutils.UserCommunicationsSetting, error)

	SetUserCommunicationsSettings(
		ctx context.Context,
		allowWhatsApp *bool,
		allowTextSms *bool,
		allowPush *bool,
		allowEmail *bool,
	) (*profileutils.UserCommunicationsSetting, error)

	GetNavActions(ctx context.Context, user profileutils.UserProfile) (*profileutils.NavigationActions, error)
	GenerateDefaultNavActions(ctx context.Context) (profileutils.NavigationActions, error)
	GenerateAgentNavActions(ctx context.Context) (profileutils.NavigationActions, error)
	GenerateEmployeeNavActions(ctx context.Context) (profileutils.NavigationActions, error)

	SaveFavoriteNavActions(ctx context.Context, title string) (bool, error)
	DeleteFavoriteNavActions(ctx context.Context, title string) (bool, error)
	RefreshNavigationActions(ctx context.Context) (*profileutils.NavigationActions, error)
	SwitchUserFlaggedFeatures(ctx context.Context, phoneNumber string) (*dto.OKResp, error)
}

// ProfileUseCaseImpl represents usecase implementation object
type ProfileUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	baseExt              extension.BaseExtension
	engagement           engagement.ServiceEngagement
	pubsub               pubsubmessaging.ServicePubSub
	crm                  crm.ServiceCrm
}

// NewProfileUseCase returns a new a onboarding usecase
func NewProfileUseCase(
	r repository.OnboardingRepository,
	ext extension.BaseExtension,
	eng engagement.ServiceEngagement,
	pubsub pubsubmessaging.ServicePubSub,
	crm crm.ServiceCrm,
) ProfileUseCase {
	return &ProfileUseCaseImpl{
		onboardingRepository: r,
		baseExt:              ext,
		engagement:           eng,
		pubsub:               pubsub,
		crm:                  crm,
	}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) UserProfile(ctx context.Context) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "UserProfile")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.UserProfileView)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !isAuthorized {
		return nil, fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return profile, nil
}

// GetProfileByID returns the profile identified by the indicated ID
func (p *ProfileUseCaseImpl) GetProfileByID(
	ctx context.Context,
	id *string,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetProfileByID")
	defer span.End()

	profile, err := p.onboardingRepository.GetUserProfileByID(ctx, *id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return profile, nil
}

// UpdateUserName updates the user username.
func (p *ProfileUseCaseImpl) UpdateUserName(ctx context.Context, userName string) error {
	ctx, span := tracer.Start(ctx, "UpdateUserName")
	defer span.End()

	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	return p.onboardingRepository.UpdateUserName(ctx, profile.ID, userName)
}

// UpdatePrimaryPhoneNumber updates the primary phone number of a specific user profile
// this should be called after a prior check of uniqueness is done
// We use `useContext` to determine
// which mode to fetch the user profile
func (p *ProfileUseCaseImpl) UpdatePrimaryPhoneNumber(
	ctx context.Context,
	phone string,
	useContext bool,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePrimaryPhoneNumber")
	defer span.End()

	var profile *profileutils.UserProfile

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	// fetch the user profile
	if useContext {
		user, err := p.baseExt.GetLoggedInUser(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return fmt.Errorf("can't get user: %w", err)
		}

		isAuthorized, err := authorization.IsAuthorized(user, permission.PrimaryPhoneUpdate)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err
		}

		if !isAuthorized {
			return fmt.Errorf("user not authorized to access this resource")
		}

		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err
		}

	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err
		}

	}

	previousPrimaryPhone := profile.PrimaryPhone
	secondaryPhones := profile.SecondaryPhoneNumbers
	if err := p.onboardingRepository.UpdatePrimaryPhoneNumber(ctx, profile.ID, phone); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	contact, err := p.crm.GetContactByPhone(ctx, *profile.PrimaryPhone)
	if err != nil {
		return fmt.Errorf("failed to get contact %s: %w", *profile.PrimaryPhone, err)
	}
	if contact == nil {
		return nil
	}

	contact.Properties.Phone = phone

	if err = p.pubsub.NotifyUpdateContact(ctx, *contact); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to crm.contact.update topic: %v", err)
	}

	// check if number to be set as primary exists in the list of secondary phones
	index, exists := utils.FindItem(secondaryPhones, *phoneNumber)
	if exists {
		// remove the phoneNumber from the Secondary Phones slice
		secondaryPhones = append(secondaryPhones[:index], secondaryPhones[index+1:]...)
	}

	secondaryPhones = append(secondaryPhones, *previousPrimaryPhone)

	if len(secondaryPhones) >= 1 {
		if err := p.onboardingRepository.UpdateSecondaryPhoneNumbers(ctx, profile.ID, secondaryPhones); err != nil {
			utils.RecordSpanError(span, err)
			return err
		}
	}

	return nil
}

// UpdatePrimaryEmailAddress updates primary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdatePrimaryEmailAddress(
	ctx context.Context,
	emailAddress string,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePrimaryEmailAddress")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.PrimaryEmailUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	err = p.onboardingRepository.UpdatePrimaryEmailAddress(ctx, profile.ID, emailAddress)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	// After updating the primary email, update it on the CRM too
	contact, err := p.crm.GetContactByPhone(ctx, *profile.PrimaryPhone)
	if err != nil {
		return fmt.Errorf("failed to get contact %s: %w", *profile.PrimaryPhone, err)
	}
	if contact == nil {
		return nil
	}

	contact.Properties.Email = emailAddress

	if err = p.pubsub.NotifyUpdateContact(ctx, *contact); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to crm.contact.update topic: %v", err)
	}

	previousPrimaryEmail := profile.PrimaryEmailAddress
	secondaryEmails := profile.SecondaryEmailAddresses

	if profile.PrimaryEmailAddress != nil {
		// Check if the email to be set as primary exists in the list
		// of secondary emails.
		index, exists := utils.FindItem(secondaryEmails, emailAddress)
		if exists {
			secondaryEmails = append(secondaryEmails[:index], secondaryEmails[index+1:]...)
		}

		secondaryEmails = append(secondaryEmails, *previousPrimaryEmail)

		if len(secondaryEmails) >= 1 {
			if err := p.onboardingRepository.UpdateSecondaryEmailAddresses(ctx, profile.ID, secondaryEmails); err != nil {
				return err
			}
		}
	}

	return nil

}

// UpdateSecondaryPhoneNumbers updates secondary phone numbers of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryPhoneNumbers(
	ctx context.Context,
	phoneNumbers []string,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSecondaryPhoneNumbers")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.SecondaryPhoneNumberUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	// assert that the phone numbers are unique
	uniquePhones := []string{}
	for _, phone := range phoneNumbers {
		exist, err := p.CheckPhoneExists(ctx, phone)
		if err != nil {
			utils.RecordSpanError(span, err)
			// this is a wrapped error. No need to wrap it again
			return err
		}

		if !exist {
			uniquePhones = append(uniquePhones, phone)
		}
	}

	if len(uniquePhones) >= 1 {
		profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			// this is a wrapped error. No need to wrap it again
			return err
		}

		return p.onboardingRepository.UpdateSecondaryPhoneNumbers(ctx, profile.ID, phoneNumbers)
	}

	// throw an error indicating the phone number(s) is/are already in the use
	return exceptions.CheckPhoneNumberExistError()
}

// UpdateSecondaryEmailAddresses updates secondary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryEmailAddresses(
	ctx context.Context,
	emailAddresses []string,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSecondaryEmailAddresses")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.SecondaryEmailAddressUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	uniqueEmails := []string{}
	for _, email := range emailAddresses {
		exist, err := p.CheckEmailExists(ctx, email)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err
		}

		if !exist {
			uniqueEmails = append(uniqueEmails, email)
		}
	}

	if len(uniqueEmails) >= 1 {
		profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err
		}

		if profile.PrimaryEmailAddress != nil {
			return p.onboardingRepository.UpdateSecondaryEmailAddresses(
				ctx,
				profile.ID,
				uniqueEmails,
			)
		}

		// internal error. primary email addresses must be present before addong secondary email
		// addresses.
		return exceptions.InternalServerError(
			fmt.Errorf(
				"primary email addresses must be present before adding secondary email addresses",
			),
		)

	}

	// throw an error indicating the email(s) is/are already in the use
	return exceptions.CheckEmailExistError()
}

// UpdateVerifiedUIDS updates the profile's verified uids
func (p *ProfileUseCaseImpl) UpdateVerifiedUIDS(ctx context.Context, uids []string) error {
	ctx, span := tracer.Start(ctx, "UpdateVerifiedUIDS")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.VerifiedUIDUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return p.onboardingRepository.UpdateVerifiedUIDS(ctx, profile.ID, uids)
}

// UpdateVerifiedIdentifiers updates the profile's verified identifiers
func (p *ProfileUseCaseImpl) UpdateVerifiedIdentifiers(
	ctx context.Context,
	identifiers []profileutils.VerifiedIdentifier,
) error {
	ctx, span := tracer.Start(ctx, "UpdateVerifiedIdentifiers")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.VerifiedIdentifiersUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	return p.onboardingRepository.UpdateVerifiedIdentifiers(ctx, profile.ID, identifiers)
}

// UpdateSuspended updates primary suspend attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdateSuspended(
	ctx context.Context,
	status bool,
	phone string,
	useContext bool,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSuspended")
	defer span.End()

	var profile *profileutils.UserProfile

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.VerifiedIdentifiersUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}
	// fetch the user profile
	if useContext {
		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			// this is a wrapped error. No need to wrap it again
			return err
		}
	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			// this is a wrapped error. No need to wrap it again
			return err
		}

	}
	return p.onboardingRepository.UpdateSuspended(ctx, profile.ID, status)
}

// UpdatePhotoUploadID updates photouploadid attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdatePhotoUploadID(ctx context.Context, uploadID string) error {
	ctx, span := tracer.Start(ctx, "UpdatePhotoUploadID")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.PhotoUploadIDUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	return p.onboardingRepository.UpdatePhotoUploadID(ctx, profile.ID, uploadID)

}

// UpdateCovers updates primary covers of a specific user profile
func (p *ProfileUseCaseImpl) UpdateCovers(ctx context.Context, covers []profileutils.Cover) error {
	ctx, span := tracer.Start(ctx, "UpdateCovers")
	defer span.End()

	if len(covers) == 0 {
		return fmt.Errorf("no covers to update found")
	}

	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	return p.onboardingRepository.UpdateCovers(
		ctx,
		profile.ID,
		utils.AddHashToCovers(covers),
	)
}

// UpdatePushTokens updates primary push tokens of a specific user profile.
func (p *ProfileUseCaseImpl) UpdatePushTokens(
	ctx context.Context,
	pushToken string,
	retire bool,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePushTokens")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.PushTokensUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	if retire {
		// check if supplied push token exist in the list of pushtokens
		index, exist := utils.FindItem(profile.PushTokens, pushToken)
		if exist {
			// remove it from the list of push tokens
			profile.PushTokens = append(profile.PushTokens[:index], profile.PushTokens[index+1:]...)
		}

		return p.onboardingRepository.UpdatePushTokens(ctx, profile.ID, profile.PushTokens)

	}
	newToken := []string{}
	newTokens := append(newToken, pushToken)

	return p.onboardingRepository.UpdatePushTokens(ctx, profile.ID, newTokens)
}

// UpdatePermissions updates the profiles permissions
func (p *ProfileUseCaseImpl) UpdatePermissions(
	ctx context.Context,
	perms []profileutils.PermissionType,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePermissions")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.PermissionsUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return p.onboardingRepository.UpdatePermissions(ctx, profile.ID, perms)
}

// AddAdminPermsToUser updates the profiles permissions
func (p *ProfileUseCaseImpl) AddAdminPermsToUser(ctx context.Context, phone string) error {
	ctx, span := tracer.Start(ctx, "AddAdminPermsToUser")
	defer span.End()

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	perms := profileutils.DefaultSuperAdminPermissions
	return p.onboardingRepository.UpdatePermissions(ctx, profile.ID, perms)
}

// RemoveAdminPermsToUser updates the profiles permissions by removing the admin permissions
// This also flips back userProfile field IsAdmin to false
func (p *ProfileUseCaseImpl) RemoveAdminPermsToUser(ctx context.Context, phone string) error {
	ctx, span := tracer.Start(ctx, "RemoveAdminPermsToUser")
	defer span.End()

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	permissions := profile.Permissions
	if len(permissions) >= 1 {
		permissions = nil
	}
	return p.onboardingRepository.UpdatePermissions(ctx, profile.ID, permissions)
}

// AddRoleToUser updates the profiles role and permissions
func (p *ProfileUseCaseImpl) AddRoleToUser(
	ctx context.Context,
	phone string,
	role profileutils.RoleType,
) error {
	ctx, span := tracer.Start(ctx, "AddRoleToUser")
	defer span.End()

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !role.IsValid() {
		return &errorcodeutil.CustomError{
			Message: fmt.Sprintf("Invalid role `%v` not available", role),
		}
	}
	return p.onboardingRepository.UpdateRole(ctx, profile.ID, role)
}

// RemoveRoleToUser updates the profiles role and permissions by setting roles to default
func (p *ProfileUseCaseImpl) RemoveRoleToUser(ctx context.Context, phone string) error {
	ctx, span := tracer.Start(ctx, "RemoveRoleToUser")
	defer span.End()

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	return p.onboardingRepository.UpdateRole(ctx, profile.ID, "")
}

// UpdateBioData updates primary biodata of a specific user profile
func (p *ProfileUseCaseImpl) UpdateBioData(ctx context.Context, data profileutils.BioData) error {
	ctx, span := tracer.Start(ctx, "UpdateBioData")
	defer span.End()

	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.BioDataUpdate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err = p.onboardingRepository.UpdateBioData(ctx, profile.ID, data); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	contact, err := p.crm.GetContactByPhone(ctx, *profile.PrimaryPhone)
	if err != nil {
		return fmt.Errorf("failed to get contact %s: %w", *profile.PrimaryPhone, err)
	}
	if contact == nil {
		return nil
	}

	if data.FirstName != nil {
		contact.Properties.FirstName = *data.FirstName
	}

	if data.LastName != nil {
		contact.Properties.LastName = *data.LastName
	}

	if data.DateOfBirth != nil {
		dob := data.DateOfBirth.AsTime()
		contact.Properties.DateOfBirth = dob
	}

	if data.Gender != "" {
		contact.Properties.Gender = data.Gender.String()
	}

	if err = p.pubsub.NotifyUpdateContact(ctx, *contact); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to crm.contact.update topic: %v", err)
	}

	return nil
}

// MaskPhoneNumbers masks phone number. the masked phone numbers will be in the form +254700***123
func (p *ProfileUseCaseImpl) MaskPhoneNumbers(phones []string) []string {
	masked := make([]string, 0, len(phones))
	for _, num := range phones {
		var b strings.Builder
		max := len(num)
		for i, p := range num {
			if i+1 == max-3 || i+1 == max-4 || i+1 == max-5 {
				fmt.Fprintf(&b, "*")
			} else {
				fmt.Fprint(&b, string(p))
			}
		}
		masked = append(masked, b.String())
	}
	return masked
}

// GetUserProfileByUID retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) GetUserProfileByUID(
	ctx context.Context,
	UID string,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByUID")
	defer span.End()

	return p.onboardingRepository.GetUserProfileByUID(ctx, UID, false)
}

// SetPrimaryPhoneNumber set the primary phone number of the user after verifying the otp code
func (p *ProfileUseCaseImpl) SetPrimaryPhoneNumber(
	ctx context.Context,
	phoneNumber string,
	otp string,
	useContext bool,
) error {
	ctx, span := tracer.Start(ctx, "SetPrimaryPhoneNumber")
	defer span.End()

	// verify otp code
	verified, err := p.engagement.VerifyOTP(
		ctx,
		phoneNumber,
		otp,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.VerifyOTPError(err)
	}

	if !verified {
		return exceptions.VerifyOTPError(nil)
	}

	// now set the primary phone number
	if err := p.UpdatePrimaryPhoneNumber(ctx, phoneNumber, useContext); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	return nil
}

// SetPrimaryEmailAddress set the primary email address of the user after verifying the otp code
func (p *ProfileUseCaseImpl) SetPrimaryEmailAddress(
	ctx context.Context,
	emailAddress string,
	otp string,
) error {
	ctx, span := tracer.Start(ctx, "SetPrimaryEmailAddress")
	defer span.End()

	UID, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.UserNotFoundError(err)
	}

	// verify otp code
	verified, err := p.engagement.VerifyEmailOTP(
		ctx,
		emailAddress,
		otp,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.VerifyOTPError(err)
	}
	if !verified {
		return exceptions.VerifyOTPError(nil)
	}

	if err := p.UpdatePrimaryEmailAddress(ctx, emailAddress); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	contact, err := p.crm.GetContactByPhone(ctx, *profile.PrimaryPhone)
	if err != nil {
		return fmt.Errorf("failed to get contact %s: %w", *profile.PrimaryPhone, err)
	}
	if contact == nil {
		return nil
	}
	contact.Properties.Email = *profile.PrimaryPhone

	if err = p.pubsub.NotifyUpdateContact(ctx, *contact); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to crm.contact.update topic: %v", err)
	}

	// The `VerifyEmail` nudge is by default created for both flavours, `PRO`
	// and `CONSUMER`, thus if a user adds and verifies their `Primary Email`
	// we need to `Resolve` the nudge for this user in both flavours
	// Resolve the nudge in `CONSUMER`
	go func() {
		// get details of the current trace span
		s := trace.SpanContextFromContext(ctx)
		// create a new context using the span configuration
		newctx := trace.ContextWithSpanContext(context.Background(), s)

		// releases resources if retry fails after a set duration
		newctx, cancel := context.WithTimeout(newctx, time.Duration(10*time.Minute))
		defer cancel()

		b := backoff.WithContext(backoff.NewExponentialBackOff(), newctx)
		cons := func() error {
			return p.engagement.ResolveDefaultNudgeByTitle(
				newctx,
				UID,
				feedlib.FlavourConsumer,
				VerifyEmailNudgeTitle,
			)
		}
		if err := backoff.Retry(
			cons,
			b,
		); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	go func() {
		// get details of the current trace span
		s := trace.SpanContextFromContext(ctx)
		// create a new context using the span configuration
		newctx := trace.ContextWithSpanContext(context.Background(), s)

		// releases resources if retry fails after a set duration
		newctx, cancel := context.WithTimeout(newctx, time.Duration(10*time.Minute))
		defer cancel()

		b := backoff.WithContext(backoff.NewExponentialBackOff(), newctx)
		pro := func() error {
			return p.engagement.ResolveDefaultNudgeByTitle(
				newctx,
				UID,
				feedlib.FlavourPro,
				VerifyEmailNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			b,
		); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	return nil
}

// CheckPhoneExists checks whether a phone number has been registered by another user.
// Checks both primary and secondary phone numbers.
func (p *ProfileUseCaseImpl) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckPhoneExists")
	defer span.End()

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.NormalizeMSISDNError(err)
	}
	exists, err := p.onboardingRepository.CheckIfPhoneNumberExists(ctx, *phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return exists, nil
}

// CheckEmailExists checks whether a email has been registered by another user.
// Checks both primary and secondary emails.
func (p *ProfileUseCaseImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckEmailExists")
	defer span.End()

	exists, err := p.onboardingRepository.CheckIfEmailExists(ctx, email)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.CheckEmailExistError()
	}
	return exists, nil
}

// RetireSecondaryPhoneNumbers overwrites an existing secondary phone number,
// if any, with the provided phone number.
func (p *ProfileUseCaseImpl) RetireSecondaryPhoneNumbers(
	ctx context.Context,
	phoneNumbers []string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "RetireSecondaryPhoneNumbers")
	defer span.End()

	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	secondaryPhoneNumbers := profile.SecondaryPhoneNumbers
	if len(secondaryPhoneNumbers) == 0 {
		return false, exceptions.SecondaryResourceHardResetError()
	}

	for _, phoneNo := range phoneNumbers {
		// Check if the passed number exists in the list of secondary numbers
		index, exists := utils.FindItem(secondaryPhoneNumbers, phoneNo)
		if exists {
			// remove the passed number from the list of secondary numbers
			secondaryPhoneNumbers = append(
				secondaryPhoneNumbers[:index],
				secondaryPhoneNumbers[index+1:]...)
		} else {
			return false, exceptions.RecordDoesNotExistError(fmt.Errorf("record does not exist"))
		}
	}

	if err := p.onboardingRepository.HardResetSecondaryPhoneNumbers(
		ctx,
		profile,
		secondaryPhoneNumbers,
	); err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	return true, nil
}

// RetireSecondaryEmailAddress removes specific secondary email addresses from the user's profile.
func (p *ProfileUseCaseImpl) RetireSecondaryEmailAddress(
	ctx context.Context,
	emailAddresses []string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "RetireSecondaryEmailAddress")
	defer span.End()

	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	secondaryEmails := profile.SecondaryEmailAddresses
	if len(secondaryEmails) == 0 {
		return false, exceptions.SecondaryResourceHardResetError()
	}

	for _, email := range emailAddresses {
		// Check if the passed email exists in the list of secondary emails
		index, exists := utils.FindItem(secondaryEmails, email)
		if exists {
			// remove the passed email from the list of secondary emails
			secondaryEmails = append(secondaryEmails[:index], secondaryEmails[index+1:]...)
		} else {
			return false, exceptions.RecordDoesNotExistError(fmt.Errorf("record does not exist"))
		}
	}

	if err := p.onboardingRepository.HardResetSecondaryEmailAddress(
		ctx,
		profile,
		secondaryEmails,
	); err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	return true, nil
}

// GetUserProfileAttributes takes a slice of UIDs and for each UID,
// it fetches the user profiles confirmed emails, phone numbers and
// FCM push tokens
func (p *ProfileUseCaseImpl) GetUserProfileAttributes(
	ctx context.Context,
	UIDs []string,
	attribute string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileAttributes")
	defer span.End()

	output := make(map[string][]string)
	values := []string{}

	for _, UID := range UIDs {
		profile, err := p.onboardingRepository.GetUserProfileByUID(
			ctx,
			UID,
			false,
		)
		if err != nil {
			utils.RecordSpanError(span, err)
			return output, err
		}

		switch attribute {
		case EmailsAttribute:
			primaryEmail := profile.PrimaryEmailAddress

			if primaryEmail == nil {
				//if not found just show an empty list
				output[UID] = []string{}
				continue
			}

			output[UID] = append(values, *primaryEmail)

			secondaryEmails := profile.SecondaryEmailAddresses
			if len(secondaryEmails) != 0 {
				for _, secondaryEmail := range secondaryEmails {
					output[UID] = append(
						values,
						*primaryEmail,
						secondaryEmail,
					)
				}
			}

		case PhoneNumbersAttribute:
			output[UID] = append(values, *profile.PrimaryPhone)
			secondaryPhones := profile.SecondaryPhoneNumbers
			if len(secondaryPhones) != 0 {
				for _, secondaryPhone := range secondaryPhones {
					output[UID] = append(
						values,
						*profile.PrimaryPhone,
						secondaryPhone,
					)
				}
			}

		case FCMTokensAttribute:
			if len(profile.PushTokens) == 0 {
				// We do not expect there to be a user profile without
				// FCM push tokens, but we can't be too sure
				// if not found just show an empty list
				output[UID] = []string{}
				continue
			}
			output[UID] = append(values, profile.PushTokens...)

		default:
			err := fmt.Errorf("failed to retrieve user profile attribute %s",
				attribute,
			)
			return nil, exceptions.RetrieveRecordError(err)
		}
	}

	return output, nil
}

// ConfirmedEmailAddresses returns verified email addresses for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ConfirmedEmailAddresses(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "ConfirmedEmailAddresses")
	defer span.End()

	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		EmailsAttribute,
	)
}

// ConfirmedPhoneNumbers returns verified phone numbers for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ConfirmedPhoneNumbers(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "ConfirmedPhoneNumbers")
	defer span.End()

	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		PhoneNumbersAttribute,
	)
}

// ValidFCMTokens returns valid FCM push tokens for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ValidFCMTokens(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "ValidFCMTokens")
	defer span.End()

	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		FCMTokensAttribute,
	)
}

// ProfileAttributes retrieves the user profiles confirmed emails,
// phone numbers and FCM push tokens
func (p *ProfileUseCaseImpl) ProfileAttributes(
	ctx context.Context,
	UIDs []string,
	attribute string,
) (map[string][]string, error) {
	ctx, span := tracer.Start(ctx, "ProfileAttributes")
	defer span.End()

	switch attribute {
	case EmailsAttribute:
		return p.ConfirmedEmailAddresses(
			ctx,
			UIDs,
		)

	case PhoneNumbersAttribute:
		return p.ConfirmedPhoneNumbers(
			ctx,
			UIDs,
		)

	case FCMTokensAttribute:
		return p.ValidFCMTokens(
			ctx,
			UIDs,
		)

	default:
		err := fmt.Errorf("failed to retrieve user profile attribute %s",
			attribute,
		)
		return nil, exceptions.RetrieveRecordError(err)
	}
}

// SetupAsExperimentParticipant sets up the logged-in user as an experiment participant.
// An experiment participant will be able to see unstable or otherwise flaged-feature in the UI of
// the app
func (p *ProfileUseCaseImpl) SetupAsExperimentParticipant(
	ctx context.Context,
	participate *bool,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "SetupAsExperimentParticipant")
	defer span.End()

	// fetch the user profile
	pr, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return false, err
	}

	if *participate {
		// add the user to the list of experiment participants
		return p.onboardingRepository.AddUserAsExperimentParticipant(ctx, pr)
	}

	// remove the user to the list of experiment participants
	return p.onboardingRepository.RemoveUserAsExperimentParticipant(ctx, pr)
}

// AddAddress adds a user's home or work address to their user's profile
func (p *ProfileUseCaseImpl) AddAddress(
	ctx context.Context,
	input dto.UserAddressInput,
	addressType enumutils.AddressType,
) (*profileutils.Address, error) {
	ctx, span := tracer.Start(ctx, "AddAddress")
	defer span.End()

	var address *profileutils.Address
	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	address = &profileutils.Address{
		// Longitude and latitude coordinates are stored with
		// 15 decimal digits right of the decimal points
		Latitude:         fmt.Sprintf("%.15f", input.Latitude),
		Longitude:        fmt.Sprintf("%.15f", input.Longitude),
		Locality:         input.Locality,
		Name:             input.Name,
		PlaceID:          input.PlaceID,
		FormattedAddress: input.FormattedAddress,
	}
	err = p.onboardingRepository.UpdateAddresses(
		ctx,
		profile.ID,
		*address,
		addressType,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return address, err
	}

	return address, nil
}

// GetAddresses returns a user's home and work addresses
func (p *ProfileUseCaseImpl) GetAddresses(
	ctx context.Context,
) (*domain.UserAddresses, error) {
	ctx, span := tracer.Start(ctx, "GetAddresses")
	defer span.End()

	profile, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	var thinHomeAddress domain.ThinAddress
	if profile.HomeAddress != nil {
		homeLatitude, err := strconv.ParseFloat(
			profile.HomeAddress.Latitude,
			64,
		)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		homeLongitude, err := strconv.ParseFloat(
			profile.HomeAddress.Longitude,
			64,
		)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		thinHomeAddress = domain.ThinAddress{
			Latitude:  homeLatitude,
			Longitude: homeLongitude,
		}
	}

	var thinWorkAddress domain.ThinAddress
	if profile.WorkAddress != nil {
		workLatitude, err := strconv.ParseFloat(
			profile.WorkAddress.Latitude,
			64,
		)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		workLongitude, err := strconv.ParseFloat(
			profile.WorkAddress.Longitude,
			64,
		)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		thinWorkAddress = domain.ThinAddress{
			Latitude:  workLatitude,
			Longitude: workLongitude,
		}
	}

	return &domain.UserAddresses{
		HomeAddress: thinHomeAddress,
		WorkAddress: thinWorkAddress,
	}, nil
}

// GetUserCommunicationsSettings  retrives the logged in user communications settings.
func (p *ProfileUseCaseImpl) GetUserCommunicationsSettings(
	ctx context.Context,
) (*profileutils.UserCommunicationsSetting, error) {
	ctx, span := tracer.Start(ctx, "GetUserCommunicationsSettings")
	defer span.End()

	pr, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return p.onboardingRepository.GetUserCommunicationsSettings(ctx, pr.ID)
}

// SetUserCommunicationsSettings sets the user communication settings
func (p *ProfileUseCaseImpl) SetUserCommunicationsSettings(
	ctx context.Context,
	allowWhatsApp *bool,
	allowTextSms *bool,
	allowPush *bool,
	allowEmail *bool,
) (*profileutils.UserCommunicationsSetting, error) {
	ctx, span := tracer.Start(ctx, "SetUserCommunicationsSettings")
	defer span.End()

	pr, err := p.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return p.onboardingRepository.SetUserCommunicationsSettings(
		ctx,
		pr.ID,
		allowWhatsApp,
		allowTextSms,
		allowPush,
		allowEmail,
	)

}

// GetNavActions Generates and returns the navigation actions for a user based on their role
func (p *ProfileUseCaseImpl) GetNavActions(
	ctx context.Context,
	user profileutils.UserProfile,
) (*profileutils.NavigationActions, error) {
	ctx, span := tracer.Start(ctx, "GetNavActions")
	defer span.End()

	var navActions profileutils.NavigationActions
	var err error

	switch user.Role {
	case profileutils.RoleTypeEmployee:
		navActions, err = p.GenerateEmployeeNavActions(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

	case profileutils.RoleTypeAgent:
		navActions, err = p.GenerateAgentNavActions(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

	default:
		navActions, err = p.GenerateDefaultNavActions(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}
	}

	// set favorite navigation actions in response
	primaryActions := []profileutils.NavAction{}
	for _, navAction := range navActions.Primary {
		if utils.CheckUserHasFavNavAction(&user, navAction.Title) {
			navAction.Favourite = true
		}
		primaryActions = append(primaryActions, navAction)
	}

	secondaryActions := []profileutils.NavAction{}
	for _, navAction := range navActions.Secondary {
		if utils.CheckUserHasFavNavAction(&user, navAction.Title) {
			navAction.Favourite = true
		}
		secondaryActions = append(secondaryActions, navAction)
	}

	navActions.Primary = primaryActions
	navActions.Secondary = secondaryActions
	return &navActions, nil
}

//GenerateEmployeeNavActions generates the navigation actions for a SIL employee
func (p *ProfileUseCaseImpl) GenerateEmployeeNavActions(
	ctx context.Context,
) (profileutils.NavigationActions, error) {

	actions := profileutils.NavigationActions{
		Primary: []profileutils.NavAction{
			{
				Title:      common.HomeNavActionTitle,
				OnTapRoute: common.HomeRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HomeNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HomeNavActionTitle,
					Description: common.HomeNavActionDescription,
					Thumbnail:   common.HomeNavActionURL,
				},
				Favourite: false,
			},
			{
				Title:      common.RequestsNavActionTitle,
				OnTapRoute: common.RequestsRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.RequestNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.RequestsNavActionTitle,
					Description: common.RequestsNavActionDescription,
					Thumbnail:   common.RequestNavActionURL,
				},
				Favourite: false,
			},
			{
				Title: common.PartnerNavActionTitle,
				// Not provided yet
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.PartnerNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.PartnerNavActionTitle,
					Description: common.PartnerNavActionDescription,
					Thumbnail:   common.PartnerNavActionURL,
				},
				Favourite: false,
			},
			{
				Title: common.ConsumerNavActionTitle,
				// Not provided yet
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.ConsumerNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.ConsumerNavActionTitle,
					Description: common.ConsumerNavActionDescription,
					Thumbnail:   common.ConsumerNavActionURL,
				},
				Favourite: false,
			},
		},
		Secondary: []profileutils.NavAction{
			{
				Title:      common.AgentNavActionTitle,
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.AgentNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.AgentNavActionTitle,
					Description: common.AgentNavActionDescription,
					Thumbnail:   common.AgentNavActionURL,
				},
				Favourite: false,
				Nested: []profileutils.NestedNavAction{
					{
						Title:      common.AgentRegistrationActionTitle,
						OnTapRoute: common.AgentRegistrationRoute,
					},
					{
						Title:      common.AgentIdentificationActionTitle,
						OnTapRoute: common.AgentIdentificationRoute,
					},
				},
			},
			{
				Title: common.PatientNavActionTitle,
				// Empty string for parent with nested actions
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.PatientNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.PatientNavActionTitle,
					Description: common.PatientNavActionDescription,
					Thumbnail:   common.PatientNavActionURL,
				},
				Favourite: false,
				Nested: []profileutils.NestedNavAction{
					{
						Title:      common.PatientRegistrationActionTitle,
						OnTapRoute: common.PatientRegistrationRoute,
					},
					{
						Title:      common.PatientIdentificationActionTitle,
						OnTapRoute: common.PatientIdentificationRoute,
					},
				},
			},
			{
				Title:      common.HelpNavActionTitle,
				OnTapRoute: common.GetHelpRouteRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HelpNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HelpNavActionTitle,
					Description: common.HelpNavActionDescription,
					Thumbnail:   common.HelpNavActionURL,
				},
				Favourite: false,
			},
		},
	}

	return actions, nil
}

//GenerateAgentNavActions generates the navigation actions for a SIL employee
func (p *ProfileUseCaseImpl) GenerateAgentNavActions(
	ctx context.Context,
) (profileutils.NavigationActions, error) {
	actions := profileutils.NavigationActions{
		Primary: []profileutils.NavAction{
			{
				Title:      common.HomeNavActionTitle,
				OnTapRoute: common.HomeRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HomeNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HomeNavActionTitle,
					Description: common.HomeNavActionDescription,
					Thumbnail:   common.HomeNavActionURL,
				},
				Favourite: false,
			},
			{
				Title:      common.PartnerNavActionTitle,
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.PartnerNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.PartnerNavActionTitle,
					Description: common.PartnerNavActionDescription,
					Thumbnail:   common.PartnerNavActionURL,
				},
				Favourite: false,
			},
			{
				Title:      common.ConsumerNavActionTitle,
				OnTapRoute: "",
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.ConsumerNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.ConsumerNavActionTitle,
					Description: common.ConsumerNavActionDescription,
					Thumbnail:   common.ConsumerNavActionURL,
				},
				Favourite: false,
			},
			{
				Title:      common.HelpNavActionTitle,
				OnTapRoute: common.GetHelpRouteRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HelpNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HelpNavActionTitle,
					Description: common.HelpNavActionDescription,
					Thumbnail:   common.HelpNavActionURL,
				},
				Favourite: false,
			},
		},
	}
	return actions, nil
}

//GenerateDefaultNavActions generates the navigation actions for a SIL employee
func (p *ProfileUseCaseImpl) GenerateDefaultNavActions(
	ctx context.Context,
) (profileutils.NavigationActions, error) {
	actions := profileutils.NavigationActions{
		Primary: []profileutils.NavAction{
			{
				Title:      common.HomeNavActionTitle,
				OnTapRoute: common.HomeRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HomeNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HomeNavActionTitle,
					Description: common.HomeNavActionDescription,
					Thumbnail:   common.HomeNavActionURL,
				},
				Favourite: false,
			},
			{
				Title:      common.HelpNavActionTitle,
				OnTapRoute: common.GetHelpRouteRoute,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         common.HelpNavActionURL,
					LinkType:    feedlib.LinkTypeSvgImage,
					Title:       common.HelpNavActionTitle,
					Description: common.HelpNavActionDescription,
					Thumbnail:   common.HelpNavActionURL,
				},
				Favourite: false,
			},
		},
	}

	return actions, nil
}

// SaveFavoriteNavActions  saves the users favorite navigation actions
func (p *ProfileUseCaseImpl) SaveFavoriteNavActions(
	ctx context.Context,
	title string,
) (bool, error) {
	userinfo, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		return false, exceptions.ProfileNotFoundError(err)
	}

	user, err := p.onboardingRepository.GetUserProfileByUID(ctx, userinfo.UID, false)
	if err != nil {
		return false, exceptions.ProfileNotFoundError(err)
	}

	favActions := user.FavNavActions
	// if user does not have such favorite action, add it.
	if !utils.CheckUserHasFavNavAction(user, title) {
		favActions = append(favActions, title)
	}

	if len(favActions) != len(user.FavNavActions)+1 {
		return false, exceptions.NavigationActionsError(
			fmt.Errorf("failed to add user favorite actions"),
		)
	}

	err = p.onboardingRepository.UpdateFavNavActions(ctx, user.ID, favActions)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteFavoriteNavActions  removes a booked marked navigation action from user profile
func (p *ProfileUseCaseImpl) DeleteFavoriteNavActions(
	ctx context.Context,
	title string,
) (bool, error) {
	userinfo, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		return false, exceptions.ProfileNotFoundError(err)
	}

	user, err := p.onboardingRepository.GetUserProfileByUID(ctx, userinfo.UID, false)
	if err != nil {
		return false, exceptions.ProfileNotFoundError(err)
	}
	var favActions []string
	for _, t := range user.FavNavActions {
		// retain the favorite action if it's not the one removed by user
		if t != title {
			favActions = append(favActions, t)
		}
	}

	if len(favActions) != len(user.FavNavActions)-1 {
		return false, exceptions.NavigationActionsError(
			fmt.Errorf("failed to remove user favorite actions"),
		)
	}

	err = p.onboardingRepository.UpdateFavNavActions(ctx, user.ID, favActions)
	if err != nil {
		return false, err
	}
	return true, nil
}

// RefreshNavigationActions gets user navigation actions only
func (p *ProfileUseCaseImpl) RefreshNavigationActions(
	ctx context.Context,
) (*profileutils.NavigationActions, error) {
	user, err := p.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}

	navAction, err := p.GetNavActions(ctx, *profile)
	if err != nil {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to get user navigation actions"),
		)
	}
	return navAction, nil
}

// SwitchUserFlaggedFeatures flips the user as opt-in or opt-out to flagged features
// once flipped the, frontend will receive an updated user profile when the person logs in again
func (p *ProfileUseCaseImpl) SwitchUserFlaggedFeatures(
	ctx context.Context,
	phoneNumber string,
) (*dto.OKResp, error) {
	ctx, span := tracer.Start(ctx, "SwitchUserFlaggedFeatures")
	defer span.End()

	profile, err := p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, phoneNumber, false)
	if err != nil {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to get user with provider phone number"),
		)
	}

	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		&auth.Token{
			UID: profile.VerifiedUIDS[0],
		},
	)

	canExperiment, err := p.onboardingRepository.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to get user with provider phone number"),
		)
	}

	// switch to the opposite
	v := !canExperiment

	_, err = p.SetupAsExperimentParticipant(authenticatedContext, &v)
	if err != nil {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to get user with provider phone number"),
		)
	}

	return &dto.OKResp{Status: "SUCCESS", Response: struct {
		SwithedFrom bool
		SwithedTo   bool
	}{SwithedFrom: canExperiment, SwithedTo: v}}, nil
}
