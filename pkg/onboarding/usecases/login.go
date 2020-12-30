package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// LoginUseCases  represents all the business logic involved in logging in a user and managing their authorization credentials.
type LoginUseCases interface {
	LoginByPhone(ctx context.Context, phone string, PIN string, flavour base.Flavour) (*domain.AuthCredentialResponse, error)
	RefreshToken(token string) (*domain.AuthCredentialResponse, error)
}

// LoginUseCasesImpl represents the usecase implementation object
type LoginUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewLoginUseCases initializes a new sign up usecase
func NewLoginUseCases(r repository.OnboardingRepository) LoginUseCases {
	return &LoginUseCasesImpl{r}
}

// LoginByPhone returns credentials that are used to log a user in
// provided the phone number and pin supplied are correct
func (l *LoginUseCasesImpl) LoginByPhone(
	ctx context.Context,
	phone string,
	PIN string,
	flavour base.Flavour,
) (*domain.AuthCredentialResponse, error) {
	profile, err := l.onboardingRepository.
		GetUserProfileByPrimaryPhoneNumber(ctx, phone)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("%v", base.ProfileNotFound)
	}

	PINData, err := l.onboardingRepository.
		GetPINByProfileID(ctx, profile.ID)

	if err != nil {
		return nil, err
	}

	if PINData == nil {
		return nil, fmt.Errorf("%v", base.PINNotFound)
	}

	// TODO: Save the specific PIN salt and use it during the matching (calvin)
	// matched := utils.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	// if !matched {
	// 	return nil, fmt.Errorf("%v", base.PINMismatch)
	// }

	return l.onboardingRepository.GenerateAuthCredentials(ctx, phone)
}

// RefreshToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (l *LoginUseCasesImpl) RefreshToken(token string) (*domain.AuthCredentialResponse, error) {
	return l.onboardingRepository.ExchangeRefreshTokenForIDToken(token)
}
