package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
)

// ILogin ...
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error)
	InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
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
func (us *UseCasesUserImpl) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error) {
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

// InviteUser is used to invite a user to the application. The invite link that is sent to the
// user will open the app if installed OR goes to the store if not installed.
func (us *UseCasesUserImpl) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return false, exceptions.InvalidFlavourDefinedError()
	}

	userProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return false, exceptions.UserNotFoundError(err)
	}

	tempPin, err := us.ExternalExt.GenerateTempPIN(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to generate temporary pin: %v", err)
	}

	salt, encryptedTempPin := us.ExternalExt.EncryptPIN(tempPin, nil)
	pinPayload := &domain.UserPIN{
		UserID:    userID,
		HashedPIN: encryptedTempPin,
		Salt:      salt,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		Flavour:   flavour,
		IsValid:   true,
	}

	_, err = us.Create.SaveTemporaryUserPin(ctx, pinPayload)
	if err != nil {
		return false, exceptions.SaveUserPinError(err)
	}

	inviteLink, err := helpers.GetInviteLink(flavour)
	if err != nil {
		return false, fmt.Errorf("failed to get invite link: %v", err)
	}

	message := helpers.CreateInviteMessage(userProfile, inviteLink, tempPin)

	err = us.ExternalExt.SendInviteSMS(ctx, []string{*phone}, message)
	if err != nil {
		return false, fmt.Errorf("failed to send SMS: %v", err)
	}

	return true, nil
}
