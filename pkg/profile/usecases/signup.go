package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/profile/domain"
)

// SignUPUseCases represents all the business logic involved in setting up a user
type SignUPUseCases interface {
	VerifyPhone(ctx context.Context, phone string) (string, error)
	LoginByPhone(ctx context.Context, phone, pin, flavour string) (*domain.UserResponse, error)
	CreateUserByPhone(ctx context.Context, phoneNumber, pin, otp string) (*domain.UserResponse, error)
	SetPhoneAsPrimary(ctx context.Context, phone string) (bool, error)
	GetUserRecoveryPhoneNumbers(ctx context.Context, phoneNumber string) ([]string, error)
	RegisterPushToken(ctx context.Context, token string) (bool, error)
	UpdatePushToken(ctx context.Context, token string) (bool, error)
	RetirePushToken(ctx context.Context, token string) (bool, error)
	RefreshToken(ctx context.Context, token string) (*domain.AuthCredentialResponse, error)
	RecordPostVisitSurvey(ctx context.Context, input domain.PostVisitSurveyInput) (bool, error)
	CompleteSignup(ctx context.Context, flavour string) (bool, error)
}
