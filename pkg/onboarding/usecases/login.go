package usecases

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	libExtension "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	libUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	"github.com/savannahghi/profileutils"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/onboarding-service/pkg/onboarding/usecases")

// LoginUseCases represents all the business logic involved in logging in a user and managing their
// authorization credentials.
type LoginUseCases interface {
	libUsecase.LoginUseCases
}

// LoginUseCasesImpl represents the usecase implementation object
type LoginUseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
	profile        libUsecase.ProfileUseCase
	baseExt        libExtension.BaseExtension
	pinExt         libExtension.PINExtension
	library        libUsecase.LoginUseCases
}

// NewLoginUseCases initializes a new sign up usecase
func NewLoginUseCases(
	infrastructure infrastructure.Infrastructure,
	p libUsecase.ProfileUseCase,
	ext libExtension.BaseExtension,
	pin libExtension.PINExtension,
) LoginUseCases {

	opensourceinfra := libInfra.NewInfrastructureInteractor()
	lib := libUsecase.NewLoginUseCases(opensourceinfra, p, ext, pin)

	return &LoginUseCasesImpl{
		infrastructure: infrastructure,
		profile:        p,
		baseExt:        ext,
		pinExt:         pin,
		library:        lib,
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
	return l.library.LoginByPhone(
		ctx,
		phone,
		PIN,
		flavour,
	)
}

// RefreshToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (l *LoginUseCasesImpl) RefreshToken(
	ctx context.Context,
	token string,
) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "LoginByPhone")
	defer span.End()
	return l.library.RefreshToken(
		ctx,
		token,
	)
}

// LoginAsAnonymous logs in a user as anonymous. This anonymous user will not have a userProfile
// since we don't have
// their phone number. All that we return is auth credentials and an error
func (l *LoginUseCasesImpl) LoginAsAnonymous(
	ctx context.Context,
) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "LoginByPhone")
	defer span.End()
	return l.library.LoginAsAnonymous(
		ctx,
	)
}

// ResumeWithPin called by the frontend check whether the currently logged in user is the one trying
// to get
// access to app
func (l *LoginUseCasesImpl) ResumeWithPin(
	ctx context.Context,
	pin string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "LoginByPhone")
	defer span.End()
	return l.library.ResumeWithPin(
		ctx,
		pin,
	)
}
