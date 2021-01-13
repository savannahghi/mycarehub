package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// LoginUseCases represents all the business logic involved in logging in a user and managing their authorization credentials.
type LoginUseCases interface {
	LoginByPhone(
		ctx context.Context,
		phone string,
		PIN string,
		flavour base.Flavour,
	) (*base.UserResponse, error)
	RefreshToken(token string) (*base.AuthCredentialResponse, error)
	LoginAsAnonymous(ctx context.Context) (*base.AuthCredentialResponse, error)
	ResumeWithPin(ctx context.Context, pin string) (bool, error)
}

// LoginUseCasesImpl represents the usecase implementation object
type LoginUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	profile              ProfileUseCase
	baseExt              extension.BaseExtension
}

// NewLoginUseCases initializes a new sign up usecase
func NewLoginUseCases(r repository.OnboardingRepository, p ProfileUseCase, ext extension.BaseExtension) LoginUseCases {
	return &LoginUseCasesImpl{onboardingRepository: r, profile: p, baseExt: ext}
}

// LoginByPhone returns credentials that are used to log a user in
// provided the phone number and pin supplied are correct
func (l *LoginUseCasesImpl) LoginByPhone(
	ctx context.Context,
	phone string,
	PIN string,
	flavour base.Flavour,
) (*base.UserResponse, error) {
	phoneNumber, err := l.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	profile, err := l.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, *phoneNumber)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}
	// defensive programming in action
	if profile == nil {
		return nil, exceptions.ProfileNotFoundError(nil)
	}

	PINData, err := l.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		return nil, exceptions.PinNotFoundError(err)
	}
	// defensive programming in action
	if PINData == nil {
		return nil, exceptions.PinNotFoundError(nil)
	}

	matched := utils.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		return nil, exceptions.PinMismatchError(nil)

	}

	auth, err := l.onboardingRepository.GenerateAuthCredentials(ctx, *phoneNumber)
	if err != nil {
		return nil, err
	}

	customer, supplier, err := l.onboardingRepository.GetCustomerOrSupplierProfileByProfileID(
		ctx,
		flavour,
		profile.ID,
	)
	if err != nil {
		return nil, exceptions.RetrieveRecordError(err)
	}

	return &base.UserResponse{
		Profile:         profile,
		CustomerProfile: customer,
		SupplierProfile: supplier,
		Auth:            *auth,
	}, nil
}

// RefreshToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (l *LoginUseCasesImpl) RefreshToken(token string) (*base.AuthCredentialResponse, error) {
	return l.onboardingRepository.ExchangeRefreshTokenForIDToken(token)
}

// LoginAsAnonymous logs in a user as anonymous. This anonymous user will not have a userProfile since we don't have
// their phone number. All that we return is auth credentials and an error
func (l *LoginUseCasesImpl) LoginAsAnonymous(ctx context.Context) (*base.AuthCredentialResponse, error) {
	return l.onboardingRepository.GenerateAuthCredentialsForAnonymousUser(ctx)
}

// ResumeWithPin called by the frontend check whether the currentlt logged in user is the once trying to get
// access to app
func (l *LoginUseCasesImpl) ResumeWithPin(ctx context.Context, pin string) (bool, error) {
	profile, err := l.profile.UserProfile(ctx)
	if err != nil {
		return false, exceptions.ProfileNotFoundError(err)
	}

	// defensive programming in action
	if profile == nil {
		return false, exceptions.ProfileNotFoundError(nil)
	}

	PINData, err := l.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		return false, exceptions.PinNotFoundError(err)
	}
	// defensive programming in action
	if PINData == nil {
		return false, exceptions.PinNotFoundError(nil)
	}

	matched := utils.ComparePIN(pin, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		// if the pins don't match, return false and dont throw an error.
		return false, nil

	}
	return true, nil
}
