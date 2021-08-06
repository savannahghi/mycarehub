package usecases

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
)

// LoginUseCases represents all the business logic involved in logging in a user and managing their
// authorization credentials.
type LoginUseCases interface {
	LoginByPhone(
		ctx context.Context,
		phone string,
		PIN string,
		flavour feedlib.Flavour,
	) (*profileutils.UserResponse, error)
	RefreshToken(ctx context.Context, token string) (*profileutils.AuthCredentialResponse, error)
	LoginAsAnonymous(ctx context.Context) (*profileutils.AuthCredentialResponse, error)
	ResumeWithPin(ctx context.Context, pin string) (bool, error)
}

// LoginUseCasesImpl represents the usecase implementation object
type LoginUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	profile              ProfileUseCase
	baseExt              extension.BaseExtension
	pinExt               extension.PINExtension
}

// NewLoginUseCases initializes a new sign up usecase
func NewLoginUseCases(
	r repository.OnboardingRepository, p ProfileUseCase,
	ext extension.BaseExtension, pin extension.PINExtension) LoginUseCases {
	return &LoginUseCasesImpl{
		onboardingRepository: r,
		profile:              p,
		baseExt:              ext,
		pinExt:               pin,
	}
}

// LoginByPhone returns credentials that are used to log a user in
// provided the phone number and pin supplied are correct
func (l *LoginUseCasesImpl) LoginByPhone(
	ctx context.Context,
	phone string,
	PIN string,
	flavour feedlib.Flavour,
) (*profileutils.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "LoginByPhone")
	defer span.End()

	phoneNumber, err := l.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	profile, err := l.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		// the error is wrapped already. No need to wrap it again
		return nil, err
	}

	PINData, err := l.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		// the error is wrapped already. No need to wrap it again
		return nil, err
	}

	matched := l.pinExt.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		return nil, exceptions.PinMismatchError(fmt.Errorf("wrong PIN credentials supplied"))

	}

	auth, err := l.onboardingRepository.GenerateAuthCredentials(ctx, *phoneNumber, profile)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Check whether the PIN is temporary i.e OTP
	// Update the auth response
	if PINData.IsOTP {
		auth.ChangePIN = true
	}

	customer, supplier, err := l.onboardingRepository.GetCustomerOrSupplierProfileByProfileID(
		ctx,
		flavour,
		profile.ID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.RetrieveRecordError(err)
	}

	// fetch the user's communication settings
	comms, err := l.onboardingRepository.GetUserCommunicationsSettings(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// get navigation actions
	roles, err := l.onboardingRepository.GetRolesByIDs(ctx, profile.Roles)
	if err != nil {
		return nil, err
	}

	navActions, err := utils.GetUserNavigationActions(ctx, *profile, *roles)
	if err != nil {
		return nil, err
	}

	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return &profileutils.UserResponse{
		Profile:               profile,
		CustomerProfile:       customer,
		SupplierProfile:       supplier,
		Auth:                  *auth,
		CommunicationSettings: comms,
		NavActions:            utils.NewActionsMapper(ctx, navActions),
	}, nil
}

// RefreshToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (l *LoginUseCasesImpl) RefreshToken(ctx context.Context, token string) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "RefreshToken")
	defer span.End()

	return l.onboardingRepository.ExchangeRefreshTokenForIDToken(ctx, token)
}

// LoginAsAnonymous logs in a user as anonymous. This anonymous user will not have a userProfile
// since we don't have
// their phone number. All that we return is auth credentials and an error
func (l *LoginUseCasesImpl) LoginAsAnonymous(
	ctx context.Context,
) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "LoginAsAnonymous")
	defer span.End()

	return l.onboardingRepository.GenerateAuthCredentialsForAnonymousUser(ctx)
}

// ResumeWithPin called by the frontend check whether the currently logged in user is the one trying
// to get
// access to app
func (l *LoginUseCasesImpl) ResumeWithPin(ctx context.Context, pin string) (bool, error) {
	ctx, span := tracer.Start(ctx, "ResumeWithPin")
	defer span.End()

	profile, err := l.profile.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// the error is wrapped already. No need to wrap it again
		return false, err
	}
	if profile == nil {
		return false, exceptions.ProfileNotFoundError(err)
	}
	PINData, err := l.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.PinNotFoundError(err)
	}
	if PINData == nil {
		return false, exceptions.PinNotFoundError(nil)
	}
	matched := l.pinExt.ComparePIN(pin, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		// if the pins don't match, return false and dont throw an error.
		return false, nil

	}
	return true, nil
}
