package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
)

const (
	inviteSMSMessage = "You have been invited to My Afya Hub. Download the app on %v. Your single use pin is %v"
	inviteLink       = "https://bwl.mobi/dl"
)

// ILogin ...
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error)
}

type IUserInvite interface {

	// TODO: send invite link via e.g SMS
	//    the invite deep link: opens the app if installed OR goes to the store if not installed
	//    a first time PIN is set and sent to the user
	//    this PIN must be changed on first use
	//    this PIN can be used only once
	//	  **encode** first use PIN and user ID into invite link
	//	  i.e not a generic invite link
	// TODO: generate first time PIN, must change, link to user
	// TODO: set the PIN valid to to the current moment so that the user is forced to change upon login
	// TODO determine communication channel for invite (e.g SMS) from settings
	Invite(ctx context.Context, userID string) (bool, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
	IUserInvite
}

// IUserInvite ...

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	Delete      infrastructure.Delete
	ExternalExt extension.ExternalMethodsExtension
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	externalExt extension.ExternalMethodsExtension,
) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
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

// Invite sends an invite to a  user (client/staff)
// The invite contains: link to app/play store, temporary PIN that **MUST** be changed on first login
//
// The default invite channel is SMS
func (us *UseCasesUserImpl) Invite(ctx context.Context, userID string) (bool, error) {

	user, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to fetch a user profile: %v", err)
	}

	fmt.Println(">>>>user", user)

	// TODO: set Flavour with the user type got from the user profile
	// if client => CONSUMER, if staff => PRO

	// var flavour feedlib.Flavour

	// if user.UserType == "Client" {
	// 	flavour = feedlib.FlavourConsumer
	// }
	// if user.UserType == "Staff" {
	// 	flavour = feedlib.FlavourPro
	// }

	pin, err := us.ExternalExt.GenerateTempPIN(ctx)
	if err != nil {
		return false, err
	}
	salt, encryptedPin := us.ExternalExt.EncryptPIN(pin, nil)
	// Set the pin to be valid for a week
	expiryDate := time.Now().AddDate(0, 0, 7)

	pinPayload := &domain.UserPIN{
		UserID:    userID,
		HashedPIN: encryptedPin,
		Salt:      salt,
		ValidFrom: time.Now(),
		ValidTo:   expiryDate,
		// TODO: add this field to the db
		// Flavour:   flavour,
		IsValid: true,
	}

	_, err = us.Create.SavePin(ctx, pinPayload)
	if err != nil {
		return false, err
	}

	// TODO: check user type to send invite to
	// if user type == client , get client_id by user_id fk in clients table, then get client contacts [] by client_id, {get contact_type by contact_id}
	// if user type == staff, get staff_id by user_id fk in staff table, then get staff contacts [] by staff_id, {get contact_type by contact_id}

	// var phoneNumber string
	// for _, contact := range user.Contacts {
	// 	if contact.Type == enums.PhoneContact.String() {
	// 		phoneNumber = contact.Contact
	// 	}
	// }

	// return false if user does not have a phone contact

	// message := fmt.Sprintf(inviteSMSMessage, inviteLink, pin)
	// err = us.ExternalExt.SendSMS(ctx, []string{phoneNumber}, message)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to send SMS: %v", err)
	// }

	return true, nil
}
