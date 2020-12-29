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
	LoginByPhone(ctx context.Context, phone, pin, flavour base.Flavour) (*domain.UserResponse, error)
	RefreshToken(ctx context.Context, token string) (*domain.AuthCredentialResponse, error)
}

// LoginUseCasesImpl represents the usecase implementation object
type LoginUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewLoginUseCases initializes a new sign up usecase
func NewLoginUseCases(r repository.OnboardingRepository) *LoginUseCasesImpl {
	return &LoginUseCasesImpl{r}
}

// LoginByPhone returns credentials that are used to log a user in
// provided the phone number and pin supplied are correct
func (o *LoginUseCasesImpl) LoginByPhone(
	ctx context.Context,
	phone string,
	PIN string,
	flavour base.Flavour,
) (*domain.AuthCredentialResponse, error) {
	profile, err := o.onboardingRepository.
		GetUserProfileByPrimaryPhoneNumber(ctx, phone)

	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("%v", base.ProfileNotFound)
	}

	PINData, err := o.onboardingRepository.
		GetPINByProfileID(ctx, profile.ID)

	if err != nil {
		return nil, err
	}

	if PINData == nil {
		return nil, fmt.Errorf("%v", "base.PINNotFound")
	}

	// TODO: Save the specific PIN salt and use it during the matching (calvin)
	// matched := utils.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	// if !matched {
	// 	return nil, fmt.Errorf("%v", base.PINMismatch)
	// }

	return o.onboardingRepository.GenerateAuthCredentials(ctx, phone)
}
