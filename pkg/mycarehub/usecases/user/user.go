package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	utilsExt "github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/serverutils"
)

// ILogin is an interface that contans login related methods
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error)
	InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
}

// IRefreshToken contains the method refreshing a token
type IRefreshToken interface {
	RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error)
}

// ISetUserPIN is an interface that contains all the user use cases for pins
type ISetUserPIN interface {
	SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error)
}

// IVerifyLoginPIN is used to verify the user's pin when logging in
type IVerifyLoginPIN interface {
	VerifyLoginPIN(ctx context.Context, userID string, pin string) (bool, int, error)
}

// IVerifyPIN is used e.g to check the PIN when accessing sensitive content
type IVerifyPIN interface {
	VerifyPIN(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error)
}

// ISetNickName is used change and or set user nickname
type ISetNickName interface {
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
}

// IRequestPinReset defines a method signature that is used to request a pin reset
type IRequestPinReset interface {
	RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
}

// ICompleteOnboardingTour defines a method that is used to complete the onboarding tour
type ICompleteOnboardingTour interface {
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
}

// IResetPIN is an interface that contains all the user use cases for pin resets
type IResetPIN interface {
	ResetPIN(ctx context.Context, input dto.UserResetPinInput) (bool, error)
}

// ICreateClientCaregiver is an interface that contains all the create client caregiver use cases
type ICreateClientCaregiver interface {
	CreateOrUpdateClientCaregiver(ctx context.Context, clientCaregiver *dto.CaregiverInput) (bool, error)
}

// IGetClientCaregiver is an interface that contains all the query client caregiver use cases
type IGetClientCaregiver interface {
	GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
	ISetUserPIN
	IVerifyLoginPIN
	ISetNickName
	IRequestPinReset
	ICompleteOnboardingTour
	IResetPIN
	IRefreshToken
	IVerifyPIN
	ICreateClientCaregiver
	IGetClientCaregiver
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

// VerifyLoginPIN checks whether a pin is valid. If a pin is invalid, it will prompt
// the user to change their pin
func (us *UseCasesUserImpl) VerifyLoginPIN(ctx context.Context, userID string, pin string) (bool, int, error) {
	pinData, err := us.Query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, int(exceptions.PINNotFound), exceptions.PinNotFoundError(err)
	}

	userProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, int(exceptions.UserNotFound), exceptions.UserNotFoundError(err)
	}

	// If pin `ValidTo` field is in the past (expired), throw an error. This means the user has to
	// change their pin on the next login
	currentTime := time.Now()
	expired := currentTime.After(pinData.ValidTo)
	if expired {
		return false, int(exceptions.ExpiredPinError), exceptions.ExpiredPinErr(fmt.Errorf("the provided pin has expired"))
	}

	matched := us.ExternalExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		failedLoginAttempts := userProfile.FailedLoginCount + 1
		err := us.Update.UpdateUserFailedLoginCount(ctx, userID, failedLoginAttempts)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, int(exceptions.LoginCountUpdateError), exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user failed login count"))
		}

		err = us.Update.UpdateUserLastFailedLoginTime(ctx, userID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, int(exceptions.LoginTimeUpdateError), exceptions.LoginTimeUpdateErr(fmt.Errorf("failed to update user last failed login time"))
		}

		nextAllowedLoginTime := utilsExt.NextAllowedLoginTime(failedLoginAttempts)
		err = us.Update.UpdateUserNextAllowedLoginTime(ctx, userID, nextAllowedLoginTime)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, int(exceptions.NexAllowedLOginTimeError), exceptions.NexAllowedLOginTimeErr(fmt.Errorf("failed to update user next allowed login time"))
		}

		return false, int(exceptions.PINMismatch), exceptions.PinMismatchError(err)
	}

	// In the event of a successful login, reset the failed login count to 0
	if userProfile.FailedLoginCount > 0 {
		err := us.Update.UpdateUserFailedLoginCount(ctx, userID, 0)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, int(exceptions.LoginCountUpdateError), exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user failed login count"))
		}
	}

	return true, int(exceptions.OK), nil
}

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, int(exceptions.InvalidPhoneNumberFormat), exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return nil, int(exceptions.InvalidFlavour), exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, int(exceptions.ProfileNotFound), exceptions.ProfileNotFoundErr(err)
	}

	if !userProfile.Active {
		return nil, int(exceptions.Internal), fmt.Errorf("user is not active")
	}

	// If the next allowed login time is after the current time, don't log in the user
	// The user has to retry after some time. We check whether time out (the current time being greater than
	// the next allowed login time) has happened. If not, the user will have to wait before trying to log in.
	currentTime := time.Now()
	timeOutOccured := currentTime.Before(*userProfile.NextAllowedLogin)
	if timeOutOccured {
		return nil, int(exceptions.Internal), fmt.Errorf("please try again after a while")
	}

	_, statusCode, err := us.VerifyLoginPIN(ctx, *userProfile.ID, pin)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, statusCode, err
	}

	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, int(exceptions.Internal), err
	}

	userTokens, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, int(exceptions.Internal), err
	}

	err = us.Update.UpdateUserLastSuccessfulLoginTime(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, int(exceptions.Internal), fmt.Errorf("failed to update user last successful login time")
	}

	return us.ReturnLoginResponse(ctx, flavour, userProfile, userTokens)
}

// ReturnLoginResponse returns either a client's or staff's response depending on the specified flavour
func (us *UseCasesUserImpl) ReturnLoginResponse(ctx context.Context, flavour feedlib.Flavour, userProfile *domain.User, userTokens *firebasetools.FirebaseUserTokens) (*domain.LoginResponse, int, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		clientProfile, err := us.Query.GetClientProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, int(exceptions.ProfileNotFound), exceptions.ClientProfileNotFoundErr(err)
		}

		clientProfile.User = userProfile
		loginResponse := &domain.LoginResponse{
			Client: clientProfile,
			AuthCredentials: domain.AuthCredentials{
				RefreshToken: userTokens.RefreshToken,
				IDToken:      userTokens.IDToken,
				ExpiresIn:    userTokens.ExpiresIn,
			},
			Code:    int(exceptions.OK),
			Message: "Success",
		}

		return loginResponse, int(exceptions.OK), nil

	case feedlib.FlavourPro:
		staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, int(exceptions.ProfileNotFound), exceptions.StaffProfileNotFoundErr(err)
		}

		staffProfile.User = userProfile
		loginResponse := &domain.LoginResponse{
			Staff: staffProfile,
			AuthCredentials: domain.AuthCredentials{
				RefreshToken: userTokens.RefreshToken,
				IDToken:      userTokens.IDToken,
				ExpiresIn:    userTokens.ExpiresIn,
			},
			Code:    int(exceptions.OK),
			Message: "Success",
		}

		return loginResponse, int(exceptions.OK), nil

	default:
		return nil, 0, fmt.Errorf("an error occurred while logging in user with phone number")
	}
}

// InviteUser is used to invite a user to the application. The invite link that is sent to the
// user will open the app if installed OR goes to the store if not installed.
func (us *UseCasesUserImpl) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return false, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotFoundError(err)
	}

	tempPin, err := us.ExternalExt.GenerateTempPIN(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GeneratePinErr(fmt.Errorf("failed to generate temporary pin: %v", err))
	}

	pinExpiryDays := serverutils.MustGetEnvVar("INVITE_PIN_EXPIRY_DAYS")

	pinExpiryDaysInt, err := strconv.Atoi(pinExpiryDays)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InternalErr(fmt.Errorf("failed to convert invite pin expiry days to int"))
	}

	pinExpiryDate := time.Now().AddDate(0, 0, pinExpiryDaysInt)

	salt, encryptedTempPin := us.ExternalExt.EncryptPIN(tempPin, nil)
	pinPayload := &domain.UserPIN{
		UserID:    userID,
		HashedPIN: encryptedTempPin,
		Salt:      salt,
		ValidFrom: time.Now(),
		ValidTo:   pinExpiryDate,
		Flavour:   flavour,
		IsValid:   true,
	}

	_, err = us.Update.InvalidatePIN(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InvalidatePinErr(err)
	}

	_, err = us.Create.SaveTemporaryUserPin(ctx, pinPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.SaveUserPinError(err)
	}

	inviteLink, err := helpers.GetInviteLink(flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetInviteLinkErr(err)
	}

	message := helpers.CreateInviteMessage(userProfile, inviteLink, tempPin)

	err = us.ExternalExt.SendInviteSMS(ctx, *phone, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.SendSMSErr(fmt.Errorf("failed to send invite SMS: %v", err))
	}

	return true, nil
}

// SetUserPIN is used to set the user's PIN
func (us *UseCasesUserImpl) SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error) {

	if err := input.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(fmt.Errorf("empty value passed in input: %v", err))
	}
	userProfile, err := us.Query.GetUserProfileByUserID(ctx, *input.UserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotFoundError(fmt.Errorf("failed to get a user profile by phonenumber: %v", err))
	}

	err = utils.ValidatePIN(*input.PIN)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ValidatePINDigitsErr(err)
	}

	salt, encryptedPIN := us.ExternalExt.EncryptPIN(*input.PIN, nil)

	isMatch := us.ExternalExt.ComparePIN(*input.ConfirmPIN, salt, encryptedPIN, nil)
	if !isMatch {
		return false, exceptions.PinMismatchError(fmt.Errorf("the provided PINs do not match"))
	}

	expiryDate, err := helpers.GetPinExpiryDate()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InternalErr(err)
	}

	pinDataPayload := &domain.UserPIN{
		UserID:    *userProfile.ID,
		HashedPIN: encryptedPIN,
		ValidFrom: time.Now(),
		ValidTo:   *expiryDate,
		Flavour:   input.Flavour,
		IsValid:   true,
		Salt:      salt,
	}

	_, err = us.Update.InvalidatePIN(ctx, *input.UserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InvalidatePinErr(err)
	}

	_, err = us.Create.SavePin(ctx, pinDataPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.SaveUserPinError(fmt.Errorf("failed to save user pin: %v", err))
	}

	return true, nil
}

// SetNickName is used to set the user's nickname
func (us *UseCasesUserImpl) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	ok, err := us.Update.SetNickName(ctx, userID, nickname)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.FailedToUpdateItemErr(fmt.Errorf("failed to set user nickname %v", err))
	}
	return ok, err
}

// RequestPINReset sends an OTP to the phone number that is provided. It begins the workflow of resetting a pin
func (us *UseCasesUserImpl) RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return "", exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.UserNotFoundError(err)
	}

	exists, err := us.Query.CheckUserHasPin(ctx, *userProfile.ID, flavour)
	if !exists {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.ExistingPINError(err)
	}

	code, err := us.OTP.GenerateAndSendOTP(ctx, *phone, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		return "", fmt.Errorf("failed to save otp")
	}

	return code, nil
}

// CompleteOnboardingTour is used to complete the onboarding tour for first time users. When a new user is
// set up, their field `pinChangeRequired` is set to true, this will inform the front end to redirect the new user
// through the process of setting a new pin, accepting terms and setting security questions. After all this is done,
// the field will be set to false. It will enable the user to be directed to the login page when they log in again.
func (us *UseCasesUserImpl) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return us.Update.UpdateUserPinChangeRequiredStatus(ctx, userID, flavour)
}

// ResetPIN resets the user's PIN when they start the reset pin process. this is a user driven request
// ensure phone/flavor is verified
// ensure the OTP for the phone is valid
// ensure the security questions were answered correctly
// ensure to invlidate the old PIN
// save new pin to db and ensure it is not duplicate for the same user
// return true if the pin was reset successfully
func (us *UseCasesUserImpl) ResetPIN(ctx context.Context, input dto.UserResetPinInput) (bool, error) {

	if err := input.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InputValidationErr(fmt.Errorf("failed to validate PIN reset Input: %v", err))
	}

	ok := input.Flavour.IsValid()
	if !ok {
		return false, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	phone, err := converterandformatter.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.NormalizeMSISDNError(err)
	}

	_, err = strconv.ParseInt(input.PIN, 10, 64)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.PINErr(fmt.Errorf("PIN should be a number: %v", err))
	}

	if len([]byte(input.PIN)) != 4 {
		return false, exceptions.PINErr(fmt.Errorf("PIN length be 4 digits: %v", err))
	}

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ContactNotFoundErr(err)

	}

	ok, err = us.Query.VerifyOTP(ctx, &dto.VerifyOTPInput{
		PhoneNumber: *phone,
		OTP:         input.OTP,
		Flavour:     input.Flavour,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotFoundError(fmt.Errorf("failed to verify otp: %v", err))
	}
	if !ok {
		return false, exceptions.InternalErr(fmt.Errorf("failed to verify otp: %v", err))
	}

	userResponse, err := us.Query.GetUserSecurityQuestionsResponses(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InternalErr(fmt.Errorf("failed to get user security question responses: %v", err))
	}

	for _, response := range userResponse {
		if !response.IsCorrect {
			return false, fmt.Errorf("user security question response is not correct")
		}
	}

	salt, encryptedPin := us.ExternalExt.EncryptPIN(input.PIN, nil)
	expiryDate, err := helpers.GetPinExpiryDate()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InternalErr(fmt.Errorf("failed to get pin expiry date: %v", err))
	}

	pinPayload := &domain.UserPIN{
		UserID:    *userProfile.ID,
		HashedPIN: encryptedPin,
		Salt:      salt,
		ValidFrom: time.Now(),
		ValidTo:   *expiryDate,
		Flavour:   input.Flavour,
		IsValid:   true,
	}

	ok, err = us.Update.InvalidatePIN(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InvalidatePinErr(err)
	}
	if !ok {
		return false, exceptions.InvalidatePinErr(err)
	}

	ok, err = us.Create.SavePin(ctx, pinPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ResetPinErr(err)
	}
	if !ok {
		return false, exceptions.ResetPinErr(err)
	}

	return true, nil
}

// RefreshToken takes a user ID and creates a custom Firebase refresh token. It then tries to fetch
// an ID token and returns auth credentials if successful
func (us *UseCasesUserImpl) RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error) {
	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	tokenResponse, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.AuthCredentials{
		RefreshToken: tokenResponse.RefreshToken,
		IDToken:      tokenResponse.IDToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}, nil
}

// VerifyPIN is used to verify the user's PIN when they are acessing e.g. sensitive information
// such as their health diary
func (us *UseCasesUserImpl) VerifyPIN(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
	if userID == "" {
		return false, exceptions.UserNotFoundError(fmt.Errorf("user id is empty"))
	}
	if !flavour.IsValid() {
		return false, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}
	if pin == "" {
		return false, exceptions.PINErr(fmt.Errorf("pin is empty"))
	}
	pinData, err := us.Query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.PinNotFoundError(err)
	}

	// If pin data does not match, this means the user cant access the data
	matched := us.ExternalExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		return false, exceptions.PinMismatchError(err)
	}
	return true, nil
}

// CreateOrUpdateClientCaregiver creates a client caregiver
func (us *UseCasesUserImpl) CreateOrUpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
	if err := caregiverInput.Validate(); err != nil {
		return false, exceptions.InputValidationErr(fmt.Errorf("failed to validate client caregiver input: %v", err))
	}
	var phone = &caregiverInput.PhoneNumber
	var err error

	if caregiverInput.PhoneNumber != "" {
		phone, err = converterandformatter.NormalizeMSISDN(caregiverInput.PhoneNumber)
		if err != nil {
			return false, exceptions.NormalizeMSISDNError(err)
		}
	}

	if !caregiverInput.CaregiverType.IsValid() {
		return false, exceptions.InputValidationErr(fmt.Errorf("caregiver type is not valid"))
	}

	caregiverInput.PhoneNumber = *phone

	client, err := us.Query.GetClientByClientID(ctx, caregiverInput.ClientID)
	if err != nil {
		return false, exceptions.ClientProfileNotFoundErr(err)
	}

	if client.CaregiverID != nil {
		err := us.Update.UpdateClientCaregiver(ctx, caregiverInput)
		if err != nil {
			return false, exceptions.UpdateClientCaregiverErr(err)
		}
	} else {

		err = us.Create.CreateClientCaregiver(ctx, caregiverInput)
		if err != nil {
			return false, exceptions.CreateClientCaregiverErr(err)
		}
	}
	return true, nil
}

// GetClientCaregiver returns a client's caregiver
func (us *UseCasesUserImpl) GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error) {
	if clientID == "" {
		return nil, exceptions.EmptyInputErr(fmt.Errorf("client id is empty"))
	}

	client, err := us.Query.GetClientByClientID(ctx, clientID)
	if err != nil {
		return nil, exceptions.ClientProfileNotFoundErr(err)
	}

	if client.CaregiverID == nil {
		return &domain.Caregiver{}, nil
	}

	caregiver, err := us.Query.GetClientCaregiver(ctx, *client.CaregiverID)
	if err != nil {
		return nil, err
	}
	return caregiver, nil
}
