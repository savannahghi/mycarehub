package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// SignUpUseCases represents all the business logic involved in setting up a user
type SignUpUseCases interface {

	// checks whether a phone number has been registred by another user. Checks both primary and
	// secondary phone numbers. If the the phone number is foreign, it send an OTP to that phone number
	// Implemented for unauthenicated REST API
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)

	// creates an account for the user, setting the provided phone number as the PRIMARY PHONE NUMBER
	// Implemented for unauthenicated REST API
	CreateUserByPhone(ctx context.Context, phoneNumber, pin string, flavour base.Flavour) (*domain.UserResponse, error)

	// updates the user profile of the currently logged in user
	// Implemented for unauthenicated GRAPHQL API
	UpdateUserProfile(ctx context.Context, input *domain.UserProfileInput) (*base.UserProfile, error)

	// adds a new push token in the users profile if the push token does not exist
	RegisterPushToken(ctx context.Context, token string) (bool, error)

	// called to create a customer account in the ERP. This API is only valid for `BEWELL CONSUMER`
	// it should be the last call after updating the users bio data. Its should not return an error
	// when it fails due to unreachable errors, rather it should retry
	CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error)

	// removes a push token from the users profile
	RetirePushToken(ctx context.Context, token string) (bool, error)

	// fetches the phone numbers of a user for the purposes of recoverying an account.
	// the returned phone numbers should be masked
	// Implemented for unauthenicated REST API
	GetUserRecoveryPhoneNumbers(ctx context.Context, phoneNumber string) (*domain.AccountRecoveryPhonesResponse, error)

	// called to set the provided phone number as the PRIMARY PHONE NUMBER in the user profile of the user
	// where the phone number is associated with.
	SetPhoneAsPrimary(ctx context.Context, phone string) (bool, error)
}

// SignUpUseCasesImpl represents usecase implementation object
type SignUpUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	profileUsecase       ProfileUseCase
}

// NewSignUpUseCases returns a new a onboarding usecase
func NewSignUpUseCases(r repository.OnboardingRepository, ob ProfileUseCase) SignUpUseCases {
	return &SignUpUseCasesImpl{
		onboardingRepository: r,
		profileUsecase:       ob}
}

// CheckPhoneExists checks whether a phone number has been registred by another user.
// Checks both primary and secondary phone numbers.
func (s *SignUpUseCasesImpl) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return false, fmt.Errorf("failed to  normalize the phone number: %v", err)
	}

	v, err := s.onboardingRepository.CheckIfPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("failed to check the phone number: %v", err)
	}

	return v, nil
}

// CreateUserByPhone creates an account for the user, setting the provided phone number as the PRIMARY PHONE NUMBER
func (s *SignUpUseCasesImpl) CreateUserByPhone(ctx context.Context, phoneNumber, pin string, flavour base.Flavour) (*domain.UserResponse, error) {

	// check if phone number is registered to another user
	exists, err := s.CheckPhoneExists(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if exists {
		return nil, fmt.Errorf("%v", base.PhoneNumberInUse)
	}

	// get or create user via thier phone number
	user, err := base.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase user: %w", err)
	}

	auth, err := s.onboardingRepository.GenerateAuthCredentials(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	profile, err := s.onboardingRepository.CreateUserProfile(ctx, phoneNumber, user.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to create userProfile: %w", err)
	}

	var supplier *domain.Supplier
	var customer *domain.Customer

	if flavour == base.FlavourPro {
		supplier, err = s.onboardingRepository.CreateEmptySupplierProfile(ctx, profile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create supplierProfile: %w", err)
		}
	}

	if flavour == base.FlavourConsumer {
		customer, err = s.onboardingRepository.CreateEmptyCustomerProfile(ctx, profile.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create customerProfile: %w", err)
		}
	}

	//TODO(calvine) : add method to encrpyt and save pin

	return &domain.UserResponse{
		Profile:         profile,
		SupplierProfile: supplier,
		CustomerProfile: customer,
		Auth:            *auth,
	}, nil
}

// UpdateUserProfile  updates the user profile of the currently logged in user
func (s *SignUpUseCasesImpl) UpdateUserProfile(ctx context.Context, input *domain.UserProfileInput) (*base.UserProfile, error) {

	if err := s.profileUsecase.UpdatePhotoUploadID(ctx, input.PhotoUploadID); err != nil {
		return nil, err
	}
	if err := s.profileUsecase.UpdateBioData(ctx, base.BioData{
		FirstName:   *input.FirstName,
		LastName:    *input.LastName,
		DateOfBirth: input.DateOfBirth,
		Gender:      *input.Gender,
	}); err != nil {
		return nil, err
	}
	return s.profileUsecase.UserProfile(ctx)
}

// RegisterPushToken adds a new push token in the users profile if the push token does not exist
func (s *SignUpUseCasesImpl) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, false); err != nil {
		return false, err
	}
	return true, nil
}

// CompleteSignup called to create a customer account in the ERP. This API is only valid for `BEWELL CONSUMER`
func (s *SignUpUseCasesImpl) CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error) {

	if flavour == base.FlavourConsumer {
		// TODO : create ERP customer account (wellnesspoints and receivablesAccount)
		return true, nil
	}
	return false, nil
}

// RetirePushToken removes a push token from the users profile
func (s *SignUpUseCasesImpl) RetirePushToken(ctx context.Context, token string) (bool, error) {
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, true); err != nil {
		return false, err
	}
	return true, nil
}

// GetUserRecoveryPhoneNumbers fetches the phone numbers of a user for the purposes of recoverying an account.
func (s *SignUpUseCasesImpl) GetUserRecoveryPhoneNumbers(ctx context.Context, phone string) (*domain.AccountRecoveryPhonesResponse, error) {
	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to  normalize the phone number: %v", err)
	}

	pr, _, err := s.onboardingRepository.GetUserProfileByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	// cherrypick the phone numbers and mask them
	phones := func(p *base.UserProfile) []string {
		phs := []string{}
		phs = append(phs, p.PrimaryPhone)
		phs = append(phs, p.SecondaryPhoneNumbers...)
		return phs

	}(pr)

	masked := s.profileUsecase.MaskPhoneNumbers(phones)

	return &domain.AccountRecoveryPhonesResponse{
		MaskedPhoneNumbers:   masked,
		UnMaskedPhoneNumbers: phones,
	}, nil
}

// SetPhoneAsPrimary called to set the provided phone number as the PRIMARY PHONE NUMBER in the user profile of the user
// where the phone number is associated with.
func (s *SignUpUseCasesImpl) SetPhoneAsPrimary(ctx context.Context, phone string) (bool, error) {
	if err := s.profileUsecase.UpdatePrimaryPhoneNumber(ctx, phone, false); err != nil {
		return false, err
	}
	return true, nil
}
