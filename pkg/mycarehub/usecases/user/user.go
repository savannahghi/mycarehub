package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	getStreamClient "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	utilsExt "github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/serverutils"
)

var (
	registerClientAPIEndpoint = serverutils.MustGetEnvVar("CLIENT_REGISTRATION_URL")
	registerStaffAPIEndpoint  = serverutils.MustGetEnvVar("STAFF_REGISTRATION_URL")
)

// ILogin is an interface that contans login related methods
type ILogin interface {
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error)
	InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
}

// IRefreshToken contains the method refreshing a token
type IRefreshToken interface {
	RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error)
	RefreshGetStreamToken(ctx context.Context, userID string) (*domain.GetStreamToken, error)
}

// ISetUserPIN is an interface that contains all the user use cases for pins
type ISetUserPIN interface {
	SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error)
}

// IVerifyLoginPIN is used to verify the user's pin when logging in
type IVerifyLoginPIN interface {
	VerifyLoginPIN(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error)
}

// IVerifyPIN is used e.g to check the PIN when accessing sensitive content
type IVerifyPIN interface {
	VerifyPIN(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error)
}

// ISetNickName is used change and or set user nickname
type ISetNickName interface {
	SetNickName(ctx context.Context, userID string, nickname string) (bool, error)
}

// IRequestPinReset defines a method signature that is used to request a pin reset
type IRequestPinReset interface {
	RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
}

// ICompleteOnboardingTour defines a method that is used to complete the onboarding tour
type ICompleteOnboardingTour interface {
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
}

// IPIN is an interface that contains all the user use cases for pins
type IPIN interface {
	ResetPIN(ctx context.Context, input dto.UserResetPinInput) (bool, error)
	GenerateTemporaryPin(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error)
}

// IClientCaregiver is an interface that contains all the client caregiver use cases
type IClientCaregiver interface {
	GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error)
	CreateOrUpdateClientCaregiver(ctx context.Context, clientCaregiver *dto.CaregiverInput) (bool, error)
}

// IRegisterUser interface defines a method signature that is used to register users
type IRegisterUser interface {
	RegisterClient(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error)
	RegisterKenyaEMRPatients(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.ClientRegistrationOutput, error)
	RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error)
}

// IClientMedicalHistory interface defines method signature for dealing with medical history
type IClientMedicalHistory interface {
	RegisteredFacilityPatients(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error)
}

// ISearchClientByCCCNumber interface contain the method used to get a client using his/her CCC number
type ISearchClientByCCCNumber interface {
	SearchClientsByCCCNumber(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error)
}

// ISearchStaffByStaffNumber interface contain the method used to get a client using his/her CCC number
type ISearchStaffByStaffNumber interface {
	SearchStaffByStaffNumber(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error)
}

// IConsent interface contains the method used to opt out a client
type IConsent interface {
	Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour, active bool) (bool, error)
}

// IUserProfile interface contains the methods to retrieve a user profile
type IUserProfile interface {
	GetUserProfile(ctx context.Context, userID string) (*domain.User, error)
}

// IClientProfile interface contains method signatures related to a client profile
type IClientProfile interface {
	AddClientFHIRID(ctx context.Context, input dto.ClientFHIRPayload) error
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
	ISetUserPIN
	IVerifyLoginPIN
	ISetNickName
	IRequestPinReset
	ICompleteOnboardingTour
	IPIN
	IRefreshToken
	IVerifyPIN
	IClientCaregiver
	IRegisterUser
	IClientMedicalHistory
	ISearchClientByCCCNumber
	ISearchStaffByStaffNumber
	IConsent
	IUserProfile
	IClientProfile
}

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	Delete      infrastructure.Delete
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
	OTP         otp.UsecaseOTP
	Authority   authority.UsecaseAuthority
	GetStream   getstream.ServiceGetStream
	Pubsub      pubsubmessaging.ServicePubsub
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
	otp otp.UsecaseOTP,
	authority authority.UsecaseAuthority,
	getstream getstream.ServiceGetStream,
	pubsub pubsubmessaging.ServicePubsub,
) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
		Update:      update,
		ExternalExt: externalExt,
		OTP:         otp,
		Authority:   authority,
		GetStream:   getstream,
		Pubsub:      pubsub,
	}
}

// GetUserProfile returns a user profile given the user ID
func (us *UseCasesUserImpl) GetUserProfile(ctx context.Context, userID string) (*domain.User, error) {
	return us.Query.GetUserProfileByUserID(ctx, userID)
}

// AddClientFHIRID updates the client profile with the patient fhir ID from clinical
func (us *UseCasesUserImpl) AddClientFHIRID(ctx context.Context, input dto.ClientFHIRPayload) error {
	client, err := us.Query.GetClientProfileByClientID(ctx, input.ClientID)
	if err != nil {
		return fmt.Errorf("error retrieving client profile: %v", err)
	}

	_, err = us.Update.UpdateClient(ctx, client, map[string]interface{}{"fhir_patient_id": input.FHIRID})
	if err != nil {
		return fmt.Errorf("error updating client profile: %v", err)
	}

	return nil
}

// VerifyLoginPIN checks whether a pin is valid. If a pin is invalid, it will prompt
// the user to change their pin
func (us *UseCasesUserImpl) VerifyLoginPIN(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error) {

	pinData, err := us.Query.GetUserPINByUserID(ctx, *userProfile.ID, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.PinNotFoundError(err)
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
		err := us.Update.UpdateUserFailedLoginCount(ctx, *userProfile.ID, failedLoginAttempts)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user failed login count"))
		}
		userProfile.FailedLoginCount++

		err = us.Update.UpdateUserLastFailedLoginTime(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.LoginTimeUpdateErr(fmt.Errorf("failed to update user last failed login time"))
		}

		nextAllowedLoginTime := utilsExt.NextAllowedLoginTime(failedLoginAttempts)
		err = us.Update.UpdateUserNextAllowedLoginTime(ctx, *userProfile.ID, nextAllowedLoginTime)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.NexAllowedLoginTimeErr(fmt.Errorf("failed to update user next allowed login time"))
		}

		return false, exceptions.PinMismatchError(fmt.Errorf("the provided pin is incorrect"))
	}

	// In the event of a successful pin match and user is not in exponential backoff, reset:
	// 1. failed login count to 0
	// 2. next allowed login time
	// 3. last successful login time
	// 4. last failed login time
	if userProfile.FailedLoginCount > 0 {
		err = us.Update.UpdateUserProfileAfterLoginSuccess(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.LoginCountUpdateErr(fmt.Errorf("failed to update user profile after login success"))
		}
	}
	return true, nil
}

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error) {
	var userProfile *domain.User

	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "Invalid phone number",
			Code:    int(exceptions.Internal),
		}, exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return &domain.LoginResponse{
			Message: "Invalid flavour",
			Code:    int(exceptions.InvalidFlavour),
		}, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err = us.Query.GetUserProfileByPhoneNumber(ctx, *phone, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "failed to get user profile by phone number",
			Code:    int(exceptions.ProfileNotFound),
		}, exceptions.ProfileNotFoundErr(err)
	}

	if !userProfile.Active {
		return &domain.LoginResponse{
			Message: "user profile is not active",
			Code:    int(exceptions.InactiveUser),
		}, fmt.Errorf("user is not active")
	}

	// If the next allowed login time is after the current time, don't log in the user
	// The user has to retry after some time. We check whether time out (the current time being greater than
	// the next allowed login time) has happened. If not, the user will have to wait before trying to log in.
	currentTime := time.Now()

	timeOutOccurred := currentTime.Before(*userProfile.NextAllowedLogin)

	if timeOutOccurred {
		loginRetryTime := userProfile.NextAllowedLogin.Sub(currentTime).Seconds()
		err := fmt.Errorf("please try again after %v seconds", loginRetryTime)
		return &domain.LoginResponse{
			Message:   err.Error(),
			RetryTime: loginRetryTime,
			Attempts:  userProfile.FailedLoginCount,
			Code:      int(exceptions.RetryLoginError),
		}, exceptions.RetryLoginErr(err)
	}

	_, err = us.VerifyLoginPIN(ctx, userProfile, pin, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message:  err.Error(),
			Code:     exceptions.GetErrorCode(err),
			Attempts: userProfile.FailedLoginCount,
		}, exceptions.GetError(err)
	}

	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "failed to create firebase custom token",
			Code:    int(exceptions.Internal),
		}, err
	}

	userTokens, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "failed to authenticate firebase custom token",
			Code:    int(exceptions.Internal),
		}, err
	}

	return us.ReturnLoginResponse(ctx, flavour, userProfile, userTokens)
}

// ReturnLoginResponse returns either a client's or staff's response depending on the specified flavour
func (us *UseCasesUserImpl) ReturnLoginResponse(ctx context.Context, flavour feedlib.Flavour, userProfile *domain.User, userTokens *firebasetools.FirebaseUserTokens) (*domain.LoginResponse, error) {
	// add user roles and permissions to the response
	roles, err := us.Authority.GetUserRoles(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "failed to get user roles",
			Code:    int(exceptions.GetUserRolesError),
		}, exceptions.GetUserRolesErr(err)
	}
	userProfile.Roles = roles

	permissions, err := us.Authority.GetUserPermissions(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return &domain.LoginResponse{
			Message: "failed to get user permissions",
			Code:    int(exceptions.GetUserPermissionsError),
		}, exceptions.GetUserPermissionsErr(err)
	}
	userProfile.Permissions = permissions

	switch flavour {
	case feedlib.FlavourConsumer:
		clientProfile, err := us.Query.GetClientProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return &domain.LoginResponse{
				Message: "failed to get client profile",
				Code:    int(exceptions.ProfileNotFound),
			}, exceptions.ClientProfileNotFoundErr(err)
		}
		clientCCCIdentifier, err := us.Query.GetClientCCCIdentifier(ctx, *clientProfile.ID)
		if err != nil {
			// An exception is not raised since we do not want to lock the user(client) out of the
			// app if their identifiers are not found.
			helpers.ReportErrorToSentry(err)
		}
		if clientCCCIdentifier != nil {
			clientProfile.CCCNumber = clientCCCIdentifier.IdentifierValue
		}

		// check if client has unresolved pin reset request
		ok, err := us.Query.CheckIfClientHasUnresolvedServiceRequests(ctx, *clientProfile.ID, string(enums.ServiceRequestTypePinReset))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return &domain.LoginResponse{
				Message: "failed to check if client has unresolved pin reset request",
				Code:    int(exceptions.Internal),
			}, exceptions.InternalErr(err)
		}
		if ok {
			err := fmt.Errorf("client has unresolved pin reset request")
			helpers.ReportErrorToSentry(err)
			return &domain.LoginResponse{
				Message: "client has unresolved pin reset request",
				Code:    int(exceptions.ClientHasUnresolvedPinResetRequestError),
			}, exceptions.ClientHasUnresolvedPinResetRequestErr(err)
		}
		if clientProfile.CHVUserID != "" {
			CHVProfile, err := us.Query.GetUserProfileByUserID(ctx, clientProfile.CHVUserID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return &domain.LoginResponse{
					Message: "failed to get CHV profile",
					Code:    int(exceptions.ProfileNotFound),
				}, exceptions.UserNotFoundError(err)
			}
			clientProfile.CHVUserName = CHVProfile.Name
		}
		// Create/update a client's getstream user
		getStreamUser := &getStreamClient.User{
			ID:   *clientProfile.ID,
			Role: "user",
			Name: userProfile.Name,
			ExtraData: map[string]interface{}{
				"userType": "CLIENT",
				"userID":   userProfile.ID,
				"nickName": userProfile.Username,
			},
		}

		_, err = us.GetStream.CreateGetStreamUser(ctx, getStreamUser)
		if err != nil {
			return &domain.LoginResponse{
				Message: "failed to create/update client's getstream user",
				Code:    int(exceptions.Internal),
			}, fmt.Errorf("failed to create client's getstream user account: %v", err)
		}

		getStreamToken, err := us.GetStream.CreateGetStreamUserToken(ctx, *clientProfile.ID)
		if err != nil {
			return &domain.LoginResponse{
				Message: "failed to create client's getstream user token",
				Code:    int(exceptions.Internal),
			}, fmt.Errorf("failed to generate getstream token: %v", err)
		}

		clientProfile.User = userProfile
		loginResponse := &domain.Response{
			Client: clientProfile,
			AuthCredentials: domain.AuthCredentials{
				RefreshToken: userTokens.RefreshToken,
				IDToken:      userTokens.IDToken,
				ExpiresIn:    userTokens.ExpiresIn,
			},
			GetStreamToken: getStreamToken,
		}

		return &domain.LoginResponse{
			Response: loginResponse,
			Message:  "Success",
			Code:     int(exceptions.OK),
		}, nil

	case feedlib.FlavourPro:
		staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return &domain.LoginResponse{
				Message: "failed to get staff profile",
				Code:    int(exceptions.ProfileNotFound),
			}, exceptions.StaffProfileNotFoundErr(err)
		}

		// Create/update a staff's getstream user
		getStreamUser := &getStreamClient.User{
			ID:   *staffProfile.ID,
			Role: "user",
			Name: userProfile.Name,
			ExtraData: map[string]interface{}{
				"userType": "STAFF",
				"userID":   userProfile.ID,
				"nickName": userProfile.Username,
			},
		}

		_, err = us.GetStream.CreateGetStreamUser(ctx, getStreamUser)
		if err != nil {
			return &domain.LoginResponse{
				Message: "failed to create/update staff's getstream user",
				Code:    int(exceptions.Internal),
			}, fmt.Errorf("failed to create staff's getstream user account: %v", err)
		}

		getStreamToken, err := us.GetStream.CreateGetStreamUserToken(ctx, *staffProfile.ID)
		if err != nil {
			return &domain.LoginResponse{
				Message: "failed to create staff's getstream user token",
				Code:    int(exceptions.Internal),
			}, fmt.Errorf("failed to generate getstream token: %v", err)
		}

		staffProfile.User = userProfile
		loginResponse := &domain.Response{
			Staff: staffProfile,
			AuthCredentials: domain.AuthCredentials{
				RefreshToken: userTokens.RefreshToken,
				IDToken:      userTokens.IDToken,
				ExpiresIn:    userTokens.ExpiresIn,
			},
			GetStreamToken: getStreamToken,
		}

		return &domain.LoginResponse{
			Response: loginResponse,
			Message:  "Success",
			Code:     int(exceptions.OK),
		}, nil

	default:
		return &domain.LoginResponse{
			Message: "an error occurred while logging in user with phone number",
			Code:    int(exceptions.Internal),
		}, fmt.Errorf("an error occurred while logging in user with phone number")
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

	inviteLink, err := helpers.GetInviteLink(flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetInviteLinkErr(err)
	}

	tempPin, err := us.GenerateTemporaryPin(ctx, userID, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetError(err)
	}

	message := helpers.CreateInviteMessage(userProfile, inviteLink, tempPin, flavour)

	err = us.ExternalExt.SendInviteSMS(ctx, *phone, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.SendSMSErr(fmt.Errorf("failed to send invite SMS: %v", err))
	}

	return true, nil
}

// GenerateTemporaryPin generates a temporary user pin and invalidates the previous user pins
func (us *UseCasesUserImpl) GenerateTemporaryPin(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
	tempPin, err := us.ExternalExt.GenerateTempPIN(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.GeneratePinErr(fmt.Errorf("failed to generate temporary pin: %v", err))
	}

	pinExpiryDays := serverutils.MustGetEnvVar("INVITE_PIN_EXPIRY_DAYS")

	pinExpiryDaysInt, err := strconv.Atoi(pinExpiryDays)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.InternalErr(fmt.Errorf("failed to convert invite pin expiry days to int"))
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

	_, err = us.Update.InvalidatePIN(ctx, userID, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.InvalidatePinErr(err)
	}

	_, err = us.Create.SaveTemporaryUserPin(ctx, pinPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.SaveUserPinError(err)
	}

	return tempPin, nil

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

	_, err = us.Update.InvalidatePIN(ctx, *input.UserID, input.Flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InvalidatePinErr(err)
	}

	_, err = us.Create.SavePin(ctx, pinDataPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.SaveUserPinError(fmt.Errorf("failed to save user pin: %v", err))
	}

	err = us.Update.UpdateUser(ctx, userProfile, map[string]interface{}{
		"pin_update_required": false,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(fmt.Errorf("failed to update user profile: %v", err))
	}

	return true, nil
}

// SetNickName is used to set the user's nickname. The nickname is also the username
func (us *UseCasesUserImpl) SetNickName(ctx context.Context, userID string, nickname string) (bool, error) {
	exists, err := us.Query.CheckIfUsernameExists(ctx, nickname)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	if exists {
		return false, exceptions.UserNameExistsErr(fmt.Errorf("username has already been taken"))
	}

	ok, err := us.Update.SetNickName(ctx, &userID, &nickname)
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

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone, flavour)
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
	return us.Update.CompleteOnboardingTour(ctx, userID, flavour)
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

	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *phone, input.Flavour)
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

	ok, err = us.Update.InvalidatePIN(ctx, *userProfile.ID, input.Flavour)
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

	// change pin update required to true
	err = us.Update.UpdateUserPinUpdateRequiredStatus(ctx, *userProfile.ID, input.Flavour, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.InternalErr(fmt.Errorf("failed to update user pin update required status: %v", err))
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
	pinData, err := us.Query.GetUserPINByUserID(ctx, userID, flavour)
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

	client, err := us.Query.GetClientProfileByClientID(ctx, caregiverInput.ClientID)
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

	client, err := us.Query.GetClientProfileByClientID(ctx, clientID)
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

// RegisterClient is used to register a client on our application. When a client is registered, their corresponding
// user profile, contacts and identifiers are created.
func (us *UseCasesUserImpl) RegisterClient(
	ctx context.Context,
	input *dto.ClientRegistrationInput,
) (*dto.ClientRegistrationOutput, error) {
	var registrationOutput *dto.ClientRegistrationOutput

	// TODO: Restore after aligning with frontend
	// check if logged in user can register client
	// err := us.Authority.CheckUserPermission(ctx, enums.PermissionTypeCanInviteClient)
	// if err != nil {
	// 	helpers.ReportErrorToSentry(err)
	// 	return nil, exceptions.UserNotAuthorizedErr(err)
	// }

	input.Gender = enumutils.Gender(strings.ToUpper(input.Gender.String()))
	resp, err := us.ExternalExt.MakeRequest(ctx, http.MethodPost, registerClientAPIEndpoint, input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	// Success is indicated with 2xx status codes
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if strings.Contains(string(dataResponse), "Identifier with this Identifier") {
			return nil, fmt.Errorf("a client with this identifier type and value already exists")
		} else if strings.Contains(string(dataResponse), "Contact with this Contact value and Flavour already exists") {
			return nil, fmt.Errorf("a contact with this value and flavour already exists")
		}
		return nil, fmt.Errorf("%v", string(dataResponse))
	}

	err = json.Unmarshal(dataResponse, &registrationOutput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if input.InviteClient {
		_, err := us.InviteUser(ctx, registrationOutput.UserID, input.PhoneNumber, feedlib.FlavourConsumer)
		if err != nil {
			return nil, fmt.Errorf("failed to invite client: %v", err)
		}
	}

	err = us.Pubsub.NotifyCreatePatient(ctx, registrationOutput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to publish to create patient topic: %v", err)
	}

	return registrationOutput, nil
}

// RefreshGetStreamToken update a getstream token as soon as a token exception occurs. The implementation
// is that frontend will call backend with the ID of the user as well as a valid session id or secret needed to authenticate them.
func (us *UseCasesUserImpl) RefreshGetStreamToken(ctx context.Context, userID string) (*domain.GetStreamToken, error) {
	_, err := us.GetStream.RevokeGetStreamUserToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke user token: %v", err)
	}

	token, err := us.GetStream.CreateGetStreamUserToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh getstream token: %v", err)
	}

	return &domain.GetStreamToken{
		Token: token,
	}, nil
}

// RegisterKenyaEMRPatients is the usecase for registering patients from KenyaEMR as clients
func (us *UseCasesUserImpl) RegisterKenyaEMRPatients(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.ClientRegistrationOutput, error) {
	clients := []*dto.ClientRegistrationOutput{}

	for _, patient := range input {
		MFLCode, err := strconv.Atoi(patient.MFLCode)
		if err != nil {

			return nil, err
		}

		exists, err := us.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
		if err != nil {
			return nil, fmt.Errorf("error checking for facility")
		}
		if !exists {
			return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", patient.MFLCode)
		}

		facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
		if err != nil {
			return nil, fmt.Errorf("error retrieving facility: %v", err)
		}

		exists, err = us.Query.CheckIdentifierExists(ctx, "CCC", patient.CCCNumber)
		if err != nil {
			return nil, fmt.Errorf("error checking for identifier")
		}
		if exists {
			return nil, fmt.Errorf("patient with that identifier exists: %v", patient.CCCNumber)
		}

		input := &dto.ClientRegistrationInput{
			Facility:       facility.Name,
			ClientType:     enums.ClientType(patient.ClientType),
			ClientName:     patient.Name,
			Gender:         enumutils.Gender(patient.Gender),
			DateOfBirth:    patient.DateOfBirth,
			PhoneNumber:    patient.PhoneNumber,
			EnrollmentDate: patient.EnrollmentDate,
			CCCNumber:      patient.CCCNumber,
			Counselled:     patient.Counselled,
			InviteClient:   true,
		}

		client, err := us.RegisterClient(ctx, input)
		if err != nil {
			return nil, err
		}

		cntct := domain.Contact{
			ContactType:  "PHONE",
			ContactValue: patient.NextOfKin.Contact,
		}
		contact, err := us.Create.GetOrCreateContact(ctx, &cntct)
		if err != nil {
			return nil, err
		}

		err = us.Create.GetOrCreateNextOfKin(ctx, &patient.NextOfKin, client.ID, *contact.ID)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// RegisteredFacilityPatients checks for newly registered clients at a facility
// from a given time i,e sync time. It is useful to fetch all patient information
// from Kenya EMR and sync it to mycarehub
func (us *UseCasesUserImpl) RegisteredFacilityPatients(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error) {
	exists, err := us.Query.CheckFacilityExistsByMFLCode(ctx, input.MFLCode)
	if err != nil {
		return nil, fmt.Errorf("error checking for facility")
	}
	if !exists {
		return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", input.MFLCode)
	}

	facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, input.MFLCode, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving facility: %v", err)
	}

	var clients []*domain.ClientProfile

	if input.SyncTime == nil {
		clients, err = us.Query.GetClientsByParams(ctx, gorm.Client{FacilityID: *facility.ID}, nil)
		if err != nil {
			return nil, err
		}
	} else {
		clients, err = us.Query.GetClientsByParams(ctx, gorm.Client{
			FacilityID: *facility.ID,
		}, input.SyncTime)
		if err != nil {
			return nil, err
		}
	}

	output := dto.PatientSyncResponse{
		MFLCode: facility.Code,
	}

	for _, client := range clients {
		identifier, err := us.Query.GetClientCCCIdentifier(ctx, *client.ID)
		if err != nil {
			return nil, err
		}

		output.Patients = append(output.Patients, identifier.IdentifierValue)
	}

	return &output, nil
}

// RegisterStaff is used to register a staff user on our application
func (us *UseCasesUserImpl) RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
	var registrationOutput *dto.StaffRegistrationOutput

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	input.Gender = enumutils.Gender(strings.ToUpper(input.Gender.String()))
	resp, err := us.ExternalExt.MakeRequest(ctx, http.MethodPost, registerStaffAPIEndpoint, input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	// Success is indicated with 2xx status codes
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if strings.Contains(string(dataResponse), "Contact with this Contact value and Flavour already exists") {
			return nil, fmt.Errorf("a contact with this value and flavour already exists")
		}
		return nil, fmt.Errorf("%v", string(dataResponse))
	}

	err = json.Unmarshal(dataResponse, &registrationOutput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if input.InviteStaff {
		_, err := us.InviteUser(ctx, registrationOutput.UserID, input.PhoneNumber, feedlib.FlavourPro)
		if err != nil {
			return nil, fmt.Errorf("failed to invite staff user: %v", err)
		}
	}

	return registrationOutput, nil
}

// SearchClientsByCCCNumber is used to search for a client using their CCC number.
func (us *UseCasesUserImpl) SearchClientsByCCCNumber(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
	if CCCNumber == "" {
		return nil, fmt.Errorf("ccc number must not be empty")
	}
	clientProfile, err := us.Query.SearchClientProfilesByCCCNumber(ctx, CCCNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get client profile: %v", err)
	}

	return clientProfile, nil
}

// SearchStaffByStaffNumber is used to search for a staff using their staff number.
// The method may also return a list of staffs at a given time depending on the value of staff number provided
func (us *UseCasesUserImpl) SearchStaffByStaffNumber(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
	if staffNumber == "" {
		return nil, fmt.Errorf("staff number must not be empty")
	}
	staffProfile, err := us.Query.SearchStaffProfileByStaffNumber(ctx, staffNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return staffProfile, nil
}

// Consent gives the client an option to choose to withdraw from the app by withdrawing their consent. It can also
// be used to opt back in to the app
func (us *UseCasesUserImpl) Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour, active bool) (bool, error) {
	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, phoneNumber, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ProfileNotFoundErr(err)
	}

	err = us.Update.UpdateUserActiveStatus(ctx, *userProfile.ID, flavour, active)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update user's active status: %v", err)
	}

	return true, nil
}
