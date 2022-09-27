package user

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	getStreamClient "github.com/GetStream/stream-chat-go/v5"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
	"github.com/lib/pq"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

// ILogin is an interface that contans login related methods
type ILogin interface {
	Login(ctx context.Context, input *dto.LoginInput) (*domain.LoginResponse, bool)
	InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error)
}

// IRefreshToken contains the method refreshing a token
type IRefreshToken interface {
	RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error)
	RefreshGetStreamToken(ctx context.Context, userID string) (*domain.GetStreamToken, error)
	RegisterPushToken(ctx context.Context, token string) (bool, error)
}

// ISetUserPIN is an interface that contains all the user use cases for pins
type ISetUserPIN interface {
	SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error)
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
	RegisterCaregiver(ctx context.Context, input dto.CaregiverInput) (*domain.CaregiverProfile, error)
	RegisterClientAsCaregiver(ctx context.Context, clientID string, caregiverNumber string) (*domain.CaregiverProfile, error)
	TransferClientToFacility(ctx context.Context, clientID *string, facilityID *string) (bool, error)
	AssignCaregiver(ctx context.Context, input dto.ClientCaregiverInput) (bool, error)
	ListClientsCaregivers(ctx context.Context, clientID string, pagination *dto.PaginationsInput) (*dto.CaregiverProfileOutputPage, error)
	ConsentToAClientCaregiver(ctx context.Context, clientID string, caregiverID string, consent bool) (bool, error)
	ConsentToManagingClient(ctx context.Context, caregiverID string, clientID string, consent bool) (bool, error)
}

// ICaregiversClients is an interface that contains all the caregiver clients use cases
type ICaregiversClients interface {
	GetCaregiverManagedClients(ctx context.Context, caregiverID string, input dto.PaginationsInput) (*dto.ManagedClientOutputPage, error)
}

// IRegisterUser interface defines a method signature that is used to register users
type IRegisterUser interface {
	RegisterClient(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error)
	RegisterKenyaEMRPatients(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error)
	RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error)
}

// IClientMedicalHistory interface defines method signature for dealing with medical history
type IClientMedicalHistory interface {
	RegisteredFacilityPatients(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error)
}

// ISearchClientUser interface contain the method used to retrieve client(s) from the database
type ISearchClientUser interface {
	SearchClientUser(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
}

// ISearchStaffUser interface contain the method used to retrieve staff(s) from the database
type ISearchStaffUser interface {
	SearchStaffUser(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error)
}

// ISearchCaregiverUser interface contain the method used to search for caregiver(s) from the database
type ISearchCaregiverUser interface {
	SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error)
}

// IConsent interface contains the method used to opt out a client
type IConsent interface {
	Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error)
}

// IUserProfile interface contains the methods to retrieve a user profile
type IUserProfile interface {
	GetUserProfile(ctx context.Context, userID string) (*domain.User, error)
	GetClientProfileByCCCNumber(ctx context.Context, cccNumber string) (*domain.ClientProfile, error)
}

// IClientProfile interface contains method signatures related to a client profile
type IClientProfile interface {
	AddClientFHIRID(ctx context.Context, input dto.ClientFHIRPayload) error
	AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error)
}

// IDeleteUser interface define the method signature that is used to delete user
type IDeleteUser interface {
	DeleteUser(ctx context.Context, payload *dto.PhoneInput) (bool, error)
}

// IUserFacility interface represents the user facility usecases
type IUserFacility interface {
	// SetDefaultFacility enables a client or a staff user to set their default facility from
	// a list of their assigned facilities
	SetStaffDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error)
	SetClientDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error)
	AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error)
	GetUserLinkedFacilities(ctx context.Context, userID string, paginationInput dto.PaginationsInput) (*dto.FacilityOutputPage, error)
	RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error)
	RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	ILogin
	ISetUserPIN
	ISetNickName
	IRequestPinReset
	ICompleteOnboardingTour
	IPIN
	IRefreshToken
	IVerifyPIN
	IClientCaregiver
	IRegisterUser
	IClientMedicalHistory
	ISearchClientUser
	ISearchStaffUser
	IConsent
	IUserProfile
	IClientProfile
	IDeleteUser
	IUserFacility
	ISearchCaregiverUser
	ICaregiversClients
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
	Clinical    clinical.IServiceClinical
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
	clinical clinical.IServiceClinical,
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
		Clinical:    clinical,
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

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, input *dto.LoginInput) (*domain.LoginResponse, bool) {
	response := domain.NewLoginResponse()

	steps := []loginFunc{
		us.userProfileCheck,
		us.checkUserIsActive,
		us.caregiverProfileCheck,
		us.clientProfileCheck,
		us.staffProfileCheck,
		us.pinResetRequestCheck,
		us.loginTimeoutCheck,
		us.checkPIN,
		us.addAuthCredentials,
		us.addRolesPermissions,
		us.addGetStreamToken,
	}

	for _, step := range steps {
		next := step(ctx, input, response)
		if !next {
			response.ClearProfiles()
			return response, false
		}
	}

	message := "login successful"
	code := exceptions.OK.Code()
	response.SetResponseCode(code, message)

	return response, true
}

// InviteUser is used to invite a user to the application. The invite link that is sent to the
// user will open the app if installed OR goes to the store if not installed.
func (us *UseCasesUserImpl) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
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
	if reinvite {
		err = us.ExternalExt.SendSMSViaTwilio(ctx, *phone, message)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.SendSMSErr(fmt.Errorf("failed to send invite SMS: %w", err))
		}
	} else {
		err = us.ExternalExt.SendInviteSMS(ctx, *phone, message)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.SendSMSErr(fmt.Errorf("failed to send invite SMS: %w", err))
		}
	}

	err = us.Update.UpdateUser(ctx, userProfile, map[string]interface{}{"pin_change_required": true})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update user: %w", err)
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
		return false, exceptions.PinMismatchError()
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

	err = us.Update.UpdateUser(ctx, &domain.User{ID: &userID}, map[string]interface{}{
		"has_set_nickname": true,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(fmt.Errorf("failed to update user profile: %v", err))
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

	// Check if the pin has expired
	// If pin `ValidTo` field is in the past (expired), throw an error. This means the user has to
	// change their pin on the next login
	currentTime := time.Now()
	expired := currentTime.After(pinData.ValidTo)
	if expired {
		return false, exceptions.ExpiredPinErr()
	}

	// If pin data does not match, this means the user cant access the data
	matched := us.ExternalExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		return false, exceptions.PinMismatchError()
	}

	return true, nil
}

// RegisterClient is used to register a client on our application. When a client is registered, their corresponding
// user profile, contacts and identifiers are created.
func (us *UseCasesUserImpl) RegisterClient(
	ctx context.Context,
	input *dto.ClientRegistrationInput,
) (*dto.ClientRegistrationOutput, error) {
	identifierExists, err := us.Query.CheckIdentifierExists(ctx, "CCC", input.CCCNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	if identifierExists {
		return nil, fmt.Errorf("an identifier with this CCC number %v already exists", input.CCCNumber)
	}

	normalized, err := converterandformatter.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		return nil, err
	}

	phoneExists, err := us.Query.CheckIfPhoneNumberExists(ctx, *normalized, false, feedlib.FlavourConsumer)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	if phoneExists {
		return nil, fmt.Errorf("a user registered with this phone number %v already exists", *normalized)
	}

	username := fmt.Sprintf("%v-%v", input.ClientName, input.CCCNumber)
	dob := input.DateOfBirth.AsTime()
	usr := &domain.User{
		Username:    username,
		Name:        input.ClientName,
		Gender:      enumutils.Gender(strings.ToUpper(input.Gender.String())),
		DateOfBirth: &dob,
		UserType:    enums.ClientUser,
		Flavour:     feedlib.FlavourConsumer,
		Active:      true,
	}

	phone := &domain.Contact{
		ContactType:  "PHONE",
		ContactValue: *normalized,
		Active:       true,
		OptedIn:      false,
		Flavour:      feedlib.FlavourConsumer,
	}

	ccc := domain.Identifier{
		IdentifierType:      "CCC",
		IdentifierValue:     input.CCCNumber,
		IdentifierUse:       "OFFICIAL",
		Description:         "CCC Number, Primary Identifier",
		IsPrimaryIdentifier: true,
		Active:              true,
	}

	MFLCode, err := strconv.Atoi(input.Facility)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	exists, err := us.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("facility with MFLCode %d does not exist", MFLCode)
	}

	facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var clientTypes []enums.ClientType
	clientTypes = append(clientTypes, input.ClientTypes...)
	clientEnrollmentDate := input.EnrollmentDate.AsTime()
	client := &domain.ClientProfile{
		ClientTypes:             clientTypes,
		TreatmentEnrollmentDate: &clientEnrollmentDate,
		FacilityID:              *facility.ID,
		ClientCounselled:        input.Counselled,
		Active:                  true,
	}

	registrationPayload := &domain.ClientRegistrationPayload{
		UserProfile:      *usr,
		Phone:            *phone,
		ClientIdentifier: ccc,
		Client:           *client,
	}

	registeredClient, err := us.Create.RegisterClient(ctx, registrationPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	payload := &dto.PatientCreationOutput{
		ID:     *registeredClient.ID,
		UserID: registeredClient.UserID,
	}
	err = us.Pubsub.NotifyCreatePatient(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to publish to create patient topic: %v", err)
	}

	handle := fmt.Sprintf("@%v", registeredClient.User.Username)
	cmsUserPayload := &dto.PubsubCreateCMSClientPayload{
		UserID:      registeredClient.UserID,
		Name:        registeredClient.User.Name,
		Gender:      registeredClient.User.Gender,
		UserType:    registeredClient.User.UserType,
		PhoneNumber: *normalized,
		Handle:      handle,
		Flavour:     registeredClient.User.Flavour,
		DateOfBirth: scalarutils.Date{
			Year:  registeredClient.User.DateOfBirth.Year(),
			Month: int(registeredClient.User.DateOfBirth.Month()),
			Day:   registeredClient.User.DateOfBirth.Day(),
		},
		ClientID:    *registeredClient.ID,
		ClientTypes: clientTypes,
		EnrollmentDate: scalarutils.Date{
			Year:  registeredClient.TreatmentEnrollmentDate.Year(),
			Month: int(registeredClient.TreatmentEnrollmentDate.Month()),
			Day:   registeredClient.TreatmentEnrollmentDate.Day(),
		},
		FacilityID:     registeredClient.FacilityID,
		FacilityName:   facility.Name,
		OrganisationID: registeredClient.OrganisationID,
	}

	err = us.Pubsub.NotifyCreateCMSClient(ctx, cmsUserPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to publish to create cms user topic: %v", err)
	}

	if input.InviteClient {
		_, err := us.InviteUser(ctx, registeredClient.UserID, input.PhoneNumber, feedlib.FlavourConsumer, false)
		if err != nil {
			return nil, fmt.Errorf("failed to invite client: %w", err)
		}
	}

	return &dto.ClientRegistrationOutput{
		ID:                *registeredClient.ID,
		Active:            registeredClient.Active,
		ClientTypes:       registeredClient.ClientTypes,
		EnrollmentDate:    registeredClient.TreatmentEnrollmentDate,
		TreatmentBuddy:    registeredClient.TreatmentBuddy,
		Counselled:        registeredClient.ClientCounselled,
		UserID:            registeredClient.UserID,
		CurrentFacilityID: registeredClient.FacilityID,
		Organisation:      registeredClient.OrganisationID,
	}, nil
}

// RegisterCaregiver is used to register a caregiver
func (us *UseCasesUserImpl) RegisterCaregiver(ctx context.Context, input dto.CaregiverInput) (*domain.CaregiverProfile, error) {
	normalized, err := converterandformatter.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to normalize phone number: %w", err)
	}

	phoneExists, err := us.Query.CheckIfPhoneNumberExists(ctx, *normalized, false, feedlib.FlavourConsumer)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check if phone number exists: %w", err)
	}
	if phoneExists {
		return nil, fmt.Errorf("phone number %v already exists", normalized)
	}

	username := fmt.Sprintf("%v-%v", input.Name, normalized)
	dob := input.DateOfBirth.AsTime()
	user := &domain.User{
		Username:    username,
		Name:        input.Name,
		Gender:      enumutils.Gender(strings.ToUpper(input.Gender.String())),
		DateOfBirth: &dob,
		UserType:    enums.CaregiverUser,
		Flavour:     feedlib.FlavourConsumer,
		Active:      true,
	}

	contact := &domain.Contact{
		ContactType:  "PHONE",
		ContactValue: *normalized,
		Active:       true,
		OptedIn:      false,
		Flavour:      feedlib.FlavourConsumer,
	}

	caregiver := &domain.Caregiver{
		CaregiverNumber: input.CaregiverNumber,
		Active:          true,
	}

	payload := &domain.CaregiverRegistration{
		User:      user,
		Contact:   contact,
		Caregiver: caregiver,
	}

	profile, err := us.Create.RegisterCaregiver(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	if input.SendInvite {
		_, err := us.InviteUser(ctx, *profile.User.ID, input.PhoneNumber, feedlib.FlavourConsumer, false)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to invite caregiver: %w", err)
		}
	}

	if len(input.AssignedClients) > 0 {
		for _, client := range input.AssignedClients {
			_, err = us.AssignCaregiver(ctx, dto.ClientCaregiverInput{
				ClientID:      client.ClientID,
				CaregiverID:   profile.ID,
				CaregiverType: client.CaregiverType,
			})
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, fmt.Errorf("failed to assign client to caregiver: %w", err)
			}
		}
	}

	return profile, nil
}

// RegisterClientAsCaregiver adds a caregiver profile to a client
func (us *UseCasesUserImpl) RegisterClientAsCaregiver(ctx context.Context, clientID string, caregiverNumber string) (*domain.CaregiverProfile, error) {
	client, err := us.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}

	// create caregiver
	caregiver, err := us.Create.CreateCaregiver(ctx, domain.Caregiver{
		UserID:          client.UserID,
		CaregiverNumber: caregiverNumber,
		Active:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create caregiver: %w", err)
	}

	return &domain.CaregiverProfile{
		ID:              caregiver.ID,
		User:            *client.User,
		CaregiverNumber: caregiver.CaregiverNumber,
	}, nil
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

func (us *UseCasesUserImpl) createClient(ctx context.Context, patient dto.PatientRegistrationPayload, facility domain.Facility) (*domain.ClientProfile, error) {
	// Adding ccc number makes it unique
	username := fmt.Sprintf("%s-%s", patient.Name, patient.CCCNumber)
	dob := patient.DateOfBirth.AsTime()
	usr := domain.User{
		Username:    username,
		Name:        patient.Name,
		Gender:      enumutils.Gender(strings.ToUpper(patient.Gender)),
		DateOfBirth: &dob,
		UserType:    enums.ClientUser,
		Flavour:     feedlib.FlavourConsumer,
	}
	user, err := us.Create.CreateUser(ctx, usr)
	if err != nil {
		return nil, err
	}

	normalized, err := converterandformatter.NormalizeMSISDN(patient.PhoneNumber)
	if err != nil {
		return nil, err
	}
	phone := domain.Contact{
		ContactType:  "PHONE",
		ContactValue: *normalized,
		Flavour:      feedlib.FlavourConsumer,
		UserID:       user.ID,
		OptedIn:      false,
	}
	contact, err := us.Create.GetOrCreateContact(ctx, &phone)
	if err != nil {
		return nil, err
	}

	ccc := domain.Identifier{
		IdentifierType:      "CCC",
		IdentifierValue:     patient.CCCNumber,
		IdentifierUse:       "OFFICIAL",
		Description:         "CCC Number, Primary Identifier",
		IsPrimaryIdentifier: true,
	}
	identifier, err := us.Create.CreateIdentifier(ctx, ccc)
	if err != nil {
		return nil, err
	}

	var clientList []enums.ClientType
	clientList = append(clientList, patient.ClientType)
	enrollment := patient.EnrollmentDate.AsTime()
	newClient := domain.ClientProfile{
		UserID:                  *user.ID,
		FacilityID:              *facility.ID,
		ClientCounselled:        patient.Counselled,
		ClientTypes:             clientList,
		TreatmentEnrollmentDate: &enrollment,
	}
	client, err := us.Create.CreateClient(ctx, newClient, *contact.ID, identifier.ID)
	if err != nil {
		return nil, err
	}

	payload := &dto.PatientCreationOutput{
		ID:     *client.ID,
		UserID: *user.ID,
	}
	err = us.Pubsub.NotifyCreatePatient(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to publish to create patient topic: %v", err)
		return client, nil
	}

	return client, nil
}

// RegisterKenyaEMRPatients is the usecase for registering patients from KenyaEMR as clients
func (us *UseCasesUserImpl) RegisterKenyaEMRPatients(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error) {
	patients := []*dto.PatientRegistrationPayload{}

	var errs error
	for _, patient := range input {
		MFLCode, err := strconv.Atoi(patient.MFLCode)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		exists, err := us.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("error checking for facility")
		}
		if !exists {

			return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", patient.MFLCode)
		}

		facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("error retrieving facility: %v", err)
		}

		// ---- Actual Client/Patient Registration begins here ----
		exists, err = us.Query.CheckIdentifierExists(ctx, "CCC", patient.CCCNumber)
		if err != nil {
			// accumulate errors rather than failing early for each client/patient
			errs = multierror.Append(errs, fmt.Errorf("error checking existing ccc number:%s, error:%w", patient.CCCNumber, err))
			helpers.ReportErrorToSentry(errs)
			continue
		}

		var client *domain.ClientProfile
		if exists {
			patients = append(patients, patient)
			continue
		} else {
			client, err = us.createClient(ctx, *patient, *facility)
			if err != nil {
				// accumulate errors rather than failing early for each client/patient
				errs = multierror.Append(errs, fmt.Errorf("error creating kenya emr client:%w", err))
				helpers.ReportErrorToSentry(errs)
				continue
			}
		}

		phone := domain.Contact{
			ContactType:  "PHONE",
			ContactValue: patient.NextOfKin.Contact,
			OptedIn:      false,
			Flavour:      feedlib.FlavourConsumer,
		}
		contact, err := us.Create.GetOrCreateContact(ctx, &phone)
		if err != nil {
			// accumulate errors rather than failing early for each client/patient
			errs = multierror.Append(errs, fmt.Errorf("error creating client next of kin contact:%w", err))
			helpers.ReportErrorToSentry(errs)
			continue
		}

		err = us.Create.GetOrCreateNextOfKin(ctx, &patient.NextOfKin, *client.ID, *contact.ID)
		if err != nil {
			// accumulate errors rather than failing early for each client/patient
			errs = multierror.Append(errs, fmt.Errorf("error creating client next of kin:%w", err))
			helpers.ReportErrorToSentry(errs)
			continue
		}

		patients = append(patients, patient)
	}

	return patients, errs
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

	var errs error
	facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, input.MFLCode, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving facility: %v", err)
	}

	var clients []*domain.ClientProfile

	if input.SyncTime == nil {
		clients, err = us.Query.GetClientsByParams(ctx, gorm.Client{FacilityID: *facility.ID}, nil)
		if err != nil {
			// accumulate errors rather than failing early for each client/patient
			errs = multierror.Append(errs, fmt.Errorf("error fetching client:%s", err))
			helpers.ReportErrorToSentry(errs)
		}
	} else {
		clients, err = us.Query.GetClientsByParams(ctx, gorm.Client{
			FacilityID: *facility.ID,
		}, input.SyncTime)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("error fetching client:%s", err))
			helpers.ReportErrorToSentry(errs)
		}
	}

	output := dto.PatientSyncResponse{
		MFLCode:  facility.Code,
		Patients: []string{},
	}

	for _, client := range clients {
		identifier, err := us.Query.GetClientCCCIdentifier(ctx, *client.ID)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to find client identifiers:%s", err))
			helpers.ReportErrorToSentry(errs)
			continue
		}

		output.Patients = append(output.Patients, identifier.IdentifierValue)
	}

	return &output, nil
}

// RegisterStaff is used to register a staff user on our application
func (us *UseCasesUserImpl) RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
	identifierExists, err := us.Query.CheckIdentifierExists(ctx, "NATIONAL_ID", input.IDNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check the existence of the identifier: %w", err)
	}
	if identifierExists {
		return nil, fmt.Errorf("identifier %v of identifier already exists", input.IDNumber)
	}

	normalized, err := converterandformatter.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to normalize phone number: %w", err)
	}

	phoneExists, err := us.Query.CheckIfPhoneNumberExists(ctx, *normalized, false, feedlib.FlavourPro)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check if phone number exists: %w", err)
	}
	if phoneExists {
		return nil, fmt.Errorf("phone number %v already exists", normalized)
	}

	username := fmt.Sprintf("%v-%v", input.StaffName, input.StaffName)
	dob := input.DateOfBirth.AsTime()
	user := &domain.User{
		Username:    username,
		Name:        input.StaffName,
		Gender:      enumutils.Gender(strings.ToUpper(input.Gender.String())),
		DateOfBirth: &dob,
		UserType:    enums.StaffUser,
		Flavour:     feedlib.FlavourPro,
		Active:      true,
	}

	contactData := &domain.Contact{
		ContactType:  "PHONE",
		ContactValue: *normalized,
		Active:       true,
		OptedIn:      false,
		Flavour:      feedlib.FlavourPro,
	}

	identifierData := &domain.Identifier{
		IdentifierType:      "NATIONAL_ID",
		IdentifierValue:     input.IDNumber,
		IdentifierUse:       "OFFICIAL",
		Description:         "NATIONAL ID, Official Identifier",
		IsPrimaryIdentifier: true,
		Active:              true,
	}

	MFLCode, err := strconv.Atoi(input.Facility)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	exists, err := us.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("facility with MFLCode %d does not exist", MFLCode)
	}

	facility, err := us.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	staffData := &domain.StaffProfile{
		Active:              true,
		StaffNumber:         input.StaffNumber,
		DefaultFacilityID:   *facility.ID,
		DefaultFacilityName: facility.Name,
	}

	staffRegistrationPayload := &domain.StaffRegistrationPayload{
		UserProfile:     *user,
		Phone:           *contactData,
		StaffIdentifier: *identifierData,
		Staff:           *staffData,
	}

	staff, err := us.Create.RegisterStaff(ctx, staffRegistrationPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to register staff: %w", err)
	}

	// UpdateRoles is used to update the roles of a user
	var staffRoles []enums.UserRoleType
	staffRoles = append(staffRoles, enums.UserRoleType(input.StaffRoles))
	_, err = us.Update.AssignRoles(ctx, staff.UserID, staffRoles)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to assign roles: %w", err)
	}

	handle := fmt.Sprintf("@%v", username)
	cmsStaffPayload := &dto.PubsubCreateCMSStaffPayload{
		UserID:      staff.UserID,
		Name:        staff.User.Name,
		Gender:      staff.User.Gender,
		UserType:    staff.User.UserType,
		PhoneNumber: *normalized,
		Handle:      handle,
		Flavour:     staff.User.Flavour,
		DateOfBirth: scalarutils.Date{
			Year:  staff.User.DateOfBirth.Year(),
			Month: int(staff.User.DateOfBirth.Month()),
			Day:   staff.User.DateOfBirth.Day(),
		},
		StaffNumber:    staff.StaffNumber,
		StaffID:        *staff.ID,
		FacilityID:     staff.DefaultFacilityID,
		FacilityName:   facility.Name,
		OrganisationID: staff.OrganisationID,
	}

	err = us.Pubsub.NotifyCreateCMSStaff(ctx, cmsStaffPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to publish staff creation event to the staff creation topic: %v", err)
	}

	if input.InviteStaff {
		_, err := us.InviteUser(ctx, staff.UserID, input.PhoneNumber, feedlib.FlavourPro, false)
		if err != nil {
			return nil, fmt.Errorf("failed to invite staff user: %v", err)
		}
	}

	return &dto.StaffRegistrationOutput{
		ID:              *staff.ID,
		Active:          staff.Active,
		StaffNumber:     input.StaffNumber,
		UserID:          staff.UserID,
		DefaultFacility: staff.DefaultFacilityID,
	}, nil
}

// SearchClientUser is used to search for a client member(s) using either of their phonenumber, username or CCC number.
func (us *UseCasesUserImpl) SearchClientUser(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
	if searchParameter == "" {
		return nil, fmt.Errorf("search parameter cannot be empty")
	}
	clientProfile, err := us.Query.SearchClientProfile(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get client profile: %w", err)
	}

	return clientProfile, nil
}

// SearchStaffUser is used to search for staff member(s) using either their phonenumber, username
// or staff number. It does this by matching of the strings based on comparison with the search Parameter
func (us *UseCasesUserImpl) SearchStaffUser(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	if searchParameter == "" {
		return nil, fmt.Errorf("search parameter cannot be empty")
	}
	staffProfile, err := us.Query.SearchStaffProfile(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return staffProfile, nil
}

// Consent gives the client an option to choose to withdraw from the app by withdrawing their consent.
func (us *UseCasesUserImpl) Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	_, err := us.DeleteUser(ctx, &dto.PhoneInput{
		PhoneNumber: phoneNumber,
		Flavour:     flavour,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to opt-out from the platform: %w", err)
	}

	return true, nil
}

// RegisterPushToken adds a new push token in the users profile
func (us *UseCasesUserImpl) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	if len(token) < 5 {
		return false, fmt.Errorf("invalid push token length")
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	var tokens pq.StringArray
	tokens = append(tokens, token)
	err = us.Update.UpdateUser(ctx, &domain.User{
		ID: &loggedInUserID,
	}, map[string]interface{}{
		"push_tokens": tokens,
	})
	if err != nil {
		return false, fmt.Errorf("failed to update user push token")
	}

	return true, nil
}

// GetClientProfileByCCCNumber is used to get a client profile by their CCC number
func (us *UseCasesUserImpl) GetClientProfileByCCCNumber(ctx context.Context, cccNumber string) (*domain.ClientProfile, error) {
	clientProfile, err := us.Query.GetClientProfileByCCCNumber(ctx, cccNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ProfileNotFoundErr(err)
	}

	return clientProfile, nil
}

// DeleteUser method is used to search for a user with a given phone number and flavour and deleted them.
// If the flavour is CONSUMER, their respective client profile as well as their user's profile.
// If flavour is PRO, their respective staff profile as well as their user's profile.
func (us *UseCasesUserImpl) DeleteUser(ctx context.Context, payload *dto.PhoneInput) (bool, error) {
	user, err := us.Query.GetUserProfileByPhoneNumber(ctx, payload.PhoneNumber, payload.Flavour)
	if err != nil {
		return false, fmt.Errorf("failed to get a user profile: %w", err)
	}

	switch payload.Flavour {
	case feedlib.FlavourConsumer:
		client, err := us.Query.GetClientProfileByUserID(ctx, *user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to get a client profile: %w", err)
		}

		go func() {
			timeoutContext, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*10))
			defer cancel()

			backOff := backoff.WithContext(backoff.NewExponentialBackOff(), timeoutContext)
			deletePatientProfile := func() error {
				err = us.Clinical.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
				if err != nil {
					helpers.ReportErrorToSentry(err)
					return fmt.Errorf("error deleting FHIR patient profile: %w", err)
				}
				return nil
			}
			if err := backoff.Retry(
				deletePatientProfile,
				backOff,
			); err != nil {
				helpers.ReportErrorToSentry(err)
				return
			}
		}()

		err = us.DeleteStreamUser(ctx, *client.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error deleting stream user: %w", err)
		}

		deleteCMSClientPayload := &dto.DeleteCMSUserPayload{
			UserID: client.UserID,
		}

		err = us.Pubsub.NotifyDeleteCMSClient(ctx, deleteCMSClientPayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("error notifying delete cms client: %v", err)
		}

		err = us.Delete.DeleteUser(ctx, *user.ID, client.ID, nil, feedlib.FlavourConsumer)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error deleting user profile: %w", err)
		}

	case feedlib.FlavourPro:
		staff, err := us.Query.GetStaffProfileByUserID(ctx, *user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error retrieving staff profile: %v", err)
		}

		err = us.DeleteStreamUser(ctx, *staff.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error deleting stream user: %v", err)
		}

		deleteCMSStaffPayload := &dto.DeleteCMSUserPayload{
			UserID: staff.UserID,
		}

		err = us.Pubsub.NotifyDeleteCMSStaff(ctx, deleteCMSStaffPayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("error notifying delete cms staff: %v", err)
		}

		err = us.Delete.DeleteUser(ctx, *user.ID, nil, staff.ID, feedlib.FlavourPro)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error deleting user profile: %v", err)
		}
	}

	return true, nil
}

// DeleteStreamUser is a helper method is used to delete a user from getstream using their ID
func (us *UseCasesUserImpl) DeleteStreamUser(ctx context.Context, id string) error {
	_, err := us.GetStream.DeleteUsers(
		ctx,
		[]string{id}, getStreamClient.DeleteUserOptions{
			User:     getStreamClient.HardDelete,
			Messages: getStreamClient.HardDelete,
		},
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("error deleting stream user: %w", err)
	}
	return nil
}

// TransferClientToFacility moves a client to a new facility
// A staff member should search for a client by their id and then transfer them to a facility
// The client profile is updated with the new facility
// The dependencies that relate to facility are updated with the current facility information
// The dependencies include:
// - All pending service requests (they should be updated to the new facility)
func (us *UseCasesUserImpl) TransferClientToFacility(ctx context.Context, clientID *string, facilityID *string) (bool, error) {
	var currentClientFacilityID string

	if clientID == nil || facilityID == nil {
		err := fmt.Errorf("clientID or facilityID is nil")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)

	}
	clientProfile, err := us.Query.GetClientProfileByClientID(ctx, *clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	currentClientFacilityID = clientProfile.FacilityID

	_, err = us.Update.UpdateClient(
		ctx,
		&domain.ClientProfile{ID: clientID},
		map[string]interface{}{"current_facility_id": facilityID},
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(err)
	}

	serviceRequests, err := us.Query.GetClientServiceRequests(ctx, "", enums.ServiceRequestStatusPending.String(), *clientID, currentClientFacilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	for _, serviceRequest := range serviceRequests {
		err = us.Update.UpdateClientServiceRequest(
			ctx,
			&domain.ServiceRequest{ID: serviceRequest.ID, Status: enums.ServiceRequestStatusPending.String()},
			map[string]interface{}{"facility_id": facilityID},
		)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.InternalErr(err)
		}
	}
	return true, nil
}

// SetStaffDefaultFacility enables a staff to set the default facility
func (us *UseCasesUserImpl) SetStaffDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error) {
	staff, err := us.Query.GetStaffProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	facilities, _, err := us.Query.GetStaffFacilities(ctx, dto.StaffFacilityInput{StaffID: staff.ID, FacilityID: &facilityID}, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get staff facilities %w", err)
	}

	if len(facilities) != 1 {
		return false, fmt.Errorf("staff user does not have  facility ID %s", facilityID)
	}

	update := map[string]interface{}{
		"default_facility_id": facilityID,
	}
	err = us.Update.UpdateStaff(ctx, staff, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	return true, nil
}

// SetClientDefaultFacility enables a client to set the default facility
func (us *UseCasesUserImpl) SetClientDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error) {

	client, err := us.Query.GetClientProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	facilities, _, err := us.Query.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: client.ID, FacilityID: &facilityID}, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get client facilities %w", err)
	}

	if len(facilities) != 1 {
		return false, fmt.Errorf("client user does not have  facility ID %s", facilityID)
	}

	update := map[string]interface{}{
		"current_facility_id": facilityID,
	}
	_, err = us.Update.UpdateClient(ctx, client, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// AddFacilitiesToClientProfile updates the client facility list
func (us *UseCasesUserImpl) AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error) {
	if clientID == "" {
		err := fmt.Errorf("client ID cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	if len(facilities) < 1 {
		err := fmt.Errorf("facilities cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	err := us.Update.AddFacilitiesToClientProfile(ctx, clientID, facilities)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update client facilities: %w", err)
	}
	return true, nil
}

// AddFacilitiesToStaffProfile updates the staff facility list
func (us *UseCasesUserImpl) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	if staffID == "" {
		err := fmt.Errorf("staff ID cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	if len(facilities) < 1 {
		err := fmt.Errorf("facilities cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	err := us.Update.AddFacilitiesToStaffProfile(ctx, staffID, facilities)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update staff facilities: %w", err)
	}
	return true, nil
}

// GetUserLinkedFacilities returns all the facilities that are linked to a user
func (us *UseCasesUserImpl) GetUserLinkedFacilities(ctx context.Context, userID string, paginationInput dto.PaginationsInput) (*dto.FacilityOutputPage, error) {
	if err := paginationInput.Validate(); err != nil {
		return nil, err
	}

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	userProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	switch userProfile.UserType {
	case enums.ClientUser:
		clientProfile, err := us.Query.GetClientProfileByUserID(ctx, userID)
		if err != nil {
			return nil, exceptions.ClientProfileNotFoundErr(err)
		}

		facilities, pageInfo, err := us.Query.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: clientProfile.ID}, page)
		if err != nil {
			return nil, fmt.Errorf("failed to get client facilities")
		}

		return &dto.FacilityOutputPage{
			Pagination: pageInfo,
			Facilities: facilities,
		}, nil

	case enums.StaffUser:
		staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, userID)
		if err != nil {
			return nil, exceptions.ClientProfileNotFoundErr(err)
		}

		facilities, pageInfo, err := us.Query.GetStaffFacilities(ctx, dto.StaffFacilityInput{StaffID: staffProfile.ID}, page)
		if err != nil {
			return nil, fmt.Errorf("failed to get client facilities")
		}

		return &dto.FacilityOutputPage{
			Pagination: pageInfo,
			Facilities: facilities,
		}, nil

	default:
		return nil, fmt.Errorf("the user has an invalid user type")
	}
}

// SearchCaregiverUser is used to search for a caregiver user
func (us *UseCasesUserImpl) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	caregiverProfile, err := us.Query.SearchCaregiverUser(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return caregiverProfile, nil
}

// RemoveFacilitiesFromClientProfile updates the client facility list to remove assigned facilities except the default facility
func (us *UseCasesUserImpl) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error) {
	if clientID == "" {
		err := fmt.Errorf("client ID cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	if len(facilities) < 1 {
		err := fmt.Errorf("facilities cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	client, err := us.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get client profile %w", err)
	}

	for _, facilityID := range facilities {
		if client.FacilityID == facilityID {
			return false, fmt.Errorf("cannot delete default facility ID: %s, please select another facility", facilityID)
		}
	}
	err = us.Delete.RemoveFacilitiesFromClientProfile(ctx, clientID, facilities)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to remove client facilities: %w", err)
	}
	return true, nil
}

// AssignCaregiver is used to assign a caregiver to a client
func (us *UseCasesUserImpl) AssignCaregiver(ctx context.Context, input dto.ClientCaregiverInput) (bool, error) {
	if input.CaregiverID == "" {
		return false, fmt.Errorf("caregiver ID is required")
	}

	userProfile, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, userProfile)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	caregiver := &domain.CaregiverClient{
		CaregiverID:      input.CaregiverID,
		ClientID:         input.ClientID,
		RelationshipType: input.CaregiverType,
		AssignedBy:       *staffProfile.ID,
	}

	err = us.Create.AddCaregiverToClient(ctx, caregiver)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to add caregiver to client: %w", err)
	}

	return true, nil
}

// RemoveFacilitiesFromStaffProfile updates the staff facility list to remove assigned facilities except the default facility
func (us *UseCasesUserImpl) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	if staffID == "" {
		err := fmt.Errorf("staff ID cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	if len(facilities) < 1 {
		err := fmt.Errorf("facilities cannot be empty")
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	staff, err := us.Query.GetStaffProfileByStaffID(ctx, staffID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get staff profile %w", err)
	}

	for _, facilityID := range facilities {
		if staff.DefaultFacilityID == facilityID {
			return false, fmt.Errorf("cannot delete default facility ID: %s, please select another facility", facilityID)
		}
	}

	err = us.Delete.RemoveFacilitiesFromStaffProfile(ctx, staffID, facilities)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update staff facilities: %w", err)
	}

	return true, nil
}

// GetCaregiverManagedClients lists clients who are managed by the caregivers
// The clients should have given their consent to be managed by the caregivers
func (us *UseCasesUserImpl) GetCaregiverManagedClients(ctx context.Context, caregiverID string, input dto.PaginationsInput) (*dto.ManagedClientOutputPage, error) {

	err := input.Validate()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("invalid pagination input: %w", err)
	}

	page := &domain.Pagination{
		Limit:       input.Limit,
		CurrentPage: input.CurrentPage,
	}

	managedClients, pageInfo, err := us.Query.GetCaregiverManagedClients(ctx, caregiverID, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get caregiver clients: %w", err)
	}

	return &dto.ManagedClientOutputPage{
		Pagination:     pageInfo,
		ManagedClients: managedClients,
	}, nil
}

// ListClientsCaregivers returns a list of caregivers for a client
func (us *UseCasesUserImpl) ListClientsCaregivers(ctx context.Context, clientID string, pagination *dto.PaginationsInput) (*dto.CaregiverProfileOutputPage, error) {
	if err := pagination.Validate(); err != nil {
		return nil, err
	}

	page := &domain.Pagination{
		Limit:       pagination.Limit,
		CurrentPage: pagination.CurrentPage,
	}

	caregivers, pageInfo, err := us.Query.ListClientsCaregivers(ctx, clientID, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to list client caregivers: %w", err)
	}

	return &dto.CaregiverProfileOutputPage{
		Pagination:       pageInfo,
		ClientCaregivers: caregivers,
	}, nil
}

// ConsentToAClientCaregiver is used to mark whether the client has acknowledged to having a certain caregiver assigned to them
func (us *UseCasesUserImpl) ConsentToAClientCaregiver(ctx context.Context, clientID string, caregiverID string, consent bool) (bool, error) {
	caregiverClient := &domain.CaregiverClient{
		ClientID:    clientID,
		CaregiverID: caregiverID,
	}

	updateData := map[string]interface{}{
		"client_consent":    consent,
		"client_consent_at": time.Now(),
	}

	if err := us.Update.UpdateCaregiverClient(ctx, caregiverClient, updateData); err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update client consent: %w", err)
	}

	return true, nil
}

// ConsentToManagingClient is used to update caregiver as having consented to offer their service to a caregiver
func (us *UseCasesUserImpl) ConsentToManagingClient(ctx context.Context, caregiverID string, clientID string, consent bool) (bool, error) {
	caregiverClient := &domain.CaregiverClient{
		CaregiverID: caregiverID,
		ClientID:    clientID,
	}

	updateData := map[string]interface{}{
		"caregiver_consent":    consent,
		"caregiver_consent_at": time.Now(),
	}

	if err := us.Update.UpdateCaregiverClient(ctx, caregiverClient, updateData); err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}
