package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	utilsExt "github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
)

// ILogin is an interface that contans login related methods
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error)
	InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
}

// ISetUserPIN is an interface that contains all the user use cases for pins
type ISetUserPIN interface {
	SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error)
}

// IVerifyPIN is used e.g to check the PIN when accessing sensitive content
type IVerifyPIN interface {
	VerifyPIN(ctx context.Context, userID string, pin string) (bool, error)
}

// ISetNickName is used change and or set user nickname
type ISetNickName interface {
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
}

// IRequestPinReset defines a method signature that is used to request a pin reset
type IRequestPinReset interface {
	RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
	ISetUserPIN
	IVerifyPIN
	ISetNickName
	IRequestPinReset
}

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	Delete      infrastructure.Delete
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
	OTP         otp.UsecaseOTP
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
	otp otp.UsecaseOTP,
) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
		Update:      update,
		ExternalExt: externalExt,
		OTP:         otp,
	}
}

// VerifyPIN checks whether a pin is valid. If a pin is invalid, it will prompt
// the user to change their pin
func (us *UseCasesUserImpl) VerifyPIN(ctx context.Context, userID string, pin string) (bool, error) {
	pinData, err := us.Query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		return false, exceptions.PinNotFoundError(err)
	}

	userProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return false, exceptions.UserNotFoundError(err)
	}

	// If pin `ValidTo` field is in the past (expired), throw an error. This means the user has to
	// change their pin on the next login
	currentTime := time.Now()
	expired := currentTime.After(pinData.ValidTo)
	if expired {
		return false, exceptions.ExpiredPinErr(fmt.Errorf("the provided pin has expired"))
	}

	matched := us.ExternalExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		failedLoginAttempts := userProfile.FailedLoginCount + 1
		err := us.Update.UpdateUserFailedLoginCount(ctx, userID, failedLoginAttempts)
		if err != nil {
			return false, exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user failed login count"))
		}

		err = us.Update.UpdateUserLastFailedLoginTime(ctx, userID)
		if err != nil {
			return false, exceptions.LoginTimeUpdateErr(fmt.Errorf("failed to update user last failed login time"))
		}

		nextAllowedLoginTime := utilsExt.NextAllowedLoginTime(failedLoginAttempts)
		err = us.Update.UpdateUserNextAllowedLoginTime(ctx, userID, nextAllowedLoginTime)
		if err != nil {
			return false, exceptions.NexAllowedLOginTimeErr(fmt.Errorf("failed to update user next allowed login time"))
		}

		return false, exceptions.PinMismatchError(err)
	}

	// In the event of a successful login, reset the failed login count to 0
	if userProfile.FailedLoginCount > 0 {
		err := us.Update.UpdateUserFailedLoginCount(ctx, userID, 0)
		if err != nil {
			return false, exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user failed login count"))
		}
	}

	return true, nil
}

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return nil, int(errorcodeutil.InvalidPhoneNumberFormat), exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return nil, int(errorcodeutil.InvalidFlavour), exceptions.InvalidFlavourDefinedError()
	}

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		return nil, int(errorcodeutil.ProfileNotFound), exceptions.UserNotFoundError(err)
	}

	if !userProfile.TermsAccepted {
		return nil, int(errorcodeutil.Internal), fmt.Errorf("user has not accepted the terms and conditions")
	}

	if !userProfile.Active {
		return nil, int(errorcodeutil.Internal), fmt.Errorf("user is not active")
	}

	// If the next allowed login time is after the current time, don't log in the user
	// The user has to retry after some time. We check whether time out (the current time being greater than
	// the next allowed login time) has happened. If not, the user will have to wait before trying to log in.
	currentTime := time.Now()
	timeOutOccured := currentTime.Before(*userProfile.NextAllowedLogin)
	if timeOutOccured {
		return nil, int(errorcodeutil.Internal), fmt.Errorf("please try again after a while")
	}

	_, err = us.VerifyPIN(ctx, *userProfile.ID, pin)
	if err != nil {
		return nil, int(errorcodeutil.PINMismatch), err
	}

	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, *userProfile.ID)
	if err != nil {
		return nil, int(errorcodeutil.Internal), err
	}

	userTokens, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, int(errorcodeutil.Internal), err
	}

	err = us.Update.UpdateUserLastSuccessfulLoginTime(ctx, *userProfile.ID)
	if err != nil {
		return nil, int(errorcodeutil.Internal), fmt.Errorf("failed to update user last successful login time")
	}

	clientProfile, err := us.Query.GetClientProfileByUserID(ctx, *userProfile.ID)
	if err != nil {
		return nil, int(errorcodeutil.Internal), fmt.Errorf("failed to return the client profile")
	}

	clientProfile.User = userProfile
	loginResponse := &domain.LoginResponse{
		Client: clientProfile,
		AuthCredentials: domain.AuthCredentials{
			RefreshToken: userTokens.RefreshToken,
			IDToken:      userTokens.IDToken,
			ExpiresIn:    userTokens.ExpiresIn,
		},
		Code:    int(errorcodeutil.OK),
		Message: "Success",
	}

	return loginResponse, int(errorcodeutil.OK), nil
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
		return false, exceptions.GeneratePinErr(fmt.Errorf("failed to generate temporary pin: %v", err))
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
		return false, exceptions.GetInviteLinkErr(err)
	}

	message := helpers.CreateInviteMessage(userProfile, inviteLink, tempPin)

	err = us.ExternalExt.SendInviteSMS(ctx, []string{*phone}, message)
	if err != nil {
		return false, exceptions.SendSMSErr(fmt.Errorf("failed to send invite SMS: %v", err))
	}

	return true, nil
}

// SetUserPIN is used to set the user's PIN
func (us *UseCasesUserImpl) SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error) {

	if err := input.Validate(); err != nil {
		return false, exceptions.EmptyInputErr(fmt.Errorf("empty value passed in input: %v", err))
	}
	userProfile, err := us.Query.GetUserProfileByUserID(ctx, *input.UserID)
	if err != nil {
		return false, exceptions.UserNotFoundError(fmt.Errorf("failed to get a user profile by phonenumber: %v", err))
	}

	err = utils.ValidatePIN(*input.PIN)
	if err != nil {
		return false, exceptions.ValidatePINDigitsErr(err)
	}

	salt, encryptedPIN := us.ExternalExt.EncryptPIN(*input.PIN, nil)

	isMatch := us.ExternalExt.ComparePIN(*input.ConfirmPIN, salt, encryptedPIN, nil)
	if !isMatch {
		return false, exceptions.PinMismatchError(fmt.Errorf("the provided PINs do not match"))
	}
	// TODO: Make this an env variable
	expiryDate := time.Now().AddDate(0, 0, 7)

	pinDataPayload := &domain.UserPIN{
		UserID:    *userProfile.ID,
		HashedPIN: encryptedPIN,
		ValidFrom: time.Now(),
		ValidTo:   expiryDate,
		Flavour:   input.Flavour,
		IsValid:   true,
		Salt:      salt,
	}

	_, err = us.Create.SavePin(ctx, pinDataPayload)
	if err != nil {
		return false, exceptions.SaveUserPinError(fmt.Errorf("failed to save user pin: %v", err))
	}

	return true, nil
}

// SetNickName is used to set the user's nickname
func (us *UseCasesUserImpl) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	ok, err := us.Update.SetNickName(ctx, userID, nickname)
	if err != nil {
		return false, exceptions.FailedToUpdateItemErr(fmt.Errorf("failed to set user nickname %v", err))
	}
	return ok, err
}

// RequestPINReset sends an OTP to the phone number that is provided. It begins the workflow of resetting a pin
func (us *UseCasesUserImpl) RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return "", exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return "", exceptions.InvalidFlavourDefinedError()
	}

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		return "", exceptions.UserNotFoundError(err)
	}

	exists, err := us.Query.CheckUserHasPin(ctx, *userProfile.ID, flavour)
	if !exists {
		return "", exceptions.ExistingPINError(err)
	}

	code, err := us.OTP.GenerateAndSendOTP(ctx, *phone, flavour)
	if err != nil {
		return "", fmt.Errorf("failed to generate and send OTP")
	}

	otpDataPayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Hour * 1),
		Channel:     "SMS",
		Flavour:     flavour,
		PhoneNumber: phoneNumber,
		OTP:         code,
	}

	err = us.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		return "", fmt.Errorf("failed to save otp")
	}

	return code, nil
}
