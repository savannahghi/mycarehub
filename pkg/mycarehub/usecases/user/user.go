package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
)

// ILogin ...
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour enums.Flavour) (*domain.AuthCredentials, string, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
}

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	Delete      infrastructure.Delete
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
		Update:      update,
		ExternalExt: externalExt,
	}
}

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, phoneNumber string, pin string, flavour enums.Flavour) (*domain.AuthCredentials, string, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return nil, "", exceptions.NormalizeMSISDNError(err)
	}

	profile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		return nil, "", exceptions.UserNotFoundError(err)
	}

	pinData, err := us.Query.GetUserPINByUserID(ctx, *profile.ID)
	if err != nil {
		return nil, "", exceptions.PinNotFoundError(err)
	}

	// If pin `ValidTo` field is in the past (expired), throw an error. This means the user has to
	// change their pin on the next login
	currentTime := time.Now()
	expired := currentTime.After(pinData.ValidTo)
	if expired {
		return nil, "", fmt.Errorf("the provided pin has expired")
	}

	matched := us.ExternalExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		return nil, "", exceptions.PinMismatchError(err)
	}

	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, *profile.ID)
	if err != nil {
		return nil, "", err
	}

	userTokens, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, "", err
	}

	authCredentials := &domain.AuthCredentials{
		User:         profile,
		RefreshToken: userTokens.RefreshToken,
		IDToken:      userTokens.IDToken,
		ExpiresIn:    userTokens.ExpiresIn,
	}

	return authCredentials, "", nil
}
