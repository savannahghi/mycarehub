package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ILogin ...
type ILogin interface {
	// ...
	// when successful: return the user object
	// when not successful: nil user, error **code**, error
	// error codes should be standardized (enum)
	// the second param: intended for the clients (mobile, web) to understand
	// the third param: a technical error that can be handled in Go e.g logged
	// TODO: After verifying PIN, check PIN valid to
	//	if in future; allow login
	// 	if in past; require change
	//	require change: communicate to mobile/web client via error code (second return value)
	//  ONLY create access token/cookie etc AFTER all checks pass
	// TODO: error codes (second param) need to be a controlled list (enum) that is
	// 	synchronized between the frontend clients, Go code and GraphQL schema.
	//	it needs to be discussed by mobile + backend devs together.
	// TODO Only allow active users to log in
	// TODO For successful logins, reset last failed login and failed login count; set last successful login
	// TODO For failed logins:
	//	increment failed login count
	//	update last failed login timestamp
	//	set next allowed login timestamp
	//	use the failed login count (post increment) as the exponent to calculate the duration/interval
	//		to add in order to get the next allowed login timestamp
	//	the base (for the exponential backoff calculation) is a setting (env + default)
	//	default this base to 4...but override to 3 for a start in env
	// TODO: Only users who have accepted terms can login
	// TODO: Update metrics e.g login count, failed login count, successful login count etc
	Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error)
}

// IUserForget models the behavior needed to comply with privacy laws e.g GDPR
// and "forget me"
type IUserForget interface {

	// Forget inactivates the user record AND hashes all identifiable information
	// After this: the user should not be available on user lists or able to log in
	// After this: it should not be possible to re-identify the user
	// This is irreversible and the UX should ensure confirmation
	// Validate: A user can only forget themselves
	// Validate: PIN is correct
	//Forget(userID string, pin string, flavour string) (bool, error)
}

// IRequestDataExport allows a user to request data known about them
// Mostly for legal compliance.
// The first impl. will simply create a task (for manual follow up) and acknowledge
type IRequestDataExport interface {
	//RequestDataExport(userID string, pin string, flavour string) (bool, error)
}

// ISetUserPIN ...
type ISetUserPIN interface {
	// SetUserPIN sets a user's PIN.
	// It can be used to set a PIN for the first time.
	// It can be used to change the PIN.
	// It can also be used to change a PIN e.g on first login after invite or
	// after expiry.
	// TODO: auditable
	// TODO: Consult CLIENT_PIN_VALIDITY_DAYS and PRO_PIN_VALIDITY DAYS env/setting to set expiry
	// TODO: flavour is an enum...same enum used in profile e.g Client, Pro
	// TODO: ensure that old PINs are not re-used
	//	this presumes that we keep a record of **hashed** PINs per user
	// TODO Each time a PIN is set, recalculate valid to / valid from and update the
	//	cached IsActive value as appropriate i.e latest PIN active, others inactive
	//
	// PINs should not be re-used (compare hashed PINs)
	// TODO: the user pin table has validity and each new PIN that is set should be a new
	// entry in the table; and also invalidate past PINs.
	// it means that the same table can be used to check for PIN reuse.
	// TODO: all PINs are hashed
	SetUserPIN(ctx context.Context, input *dto.PinInput) (bool, error)
}

// IResetPIN ...
type IResetPIN interface {
	// ResetPIN can be used by admins or healthcare workers to generate and send
	// a new PIN for a client or other user.
	// The new PIN is generated automatically and set to expire immediately so
	// that a PIN change is forced on next login.
	// TODO: Notify user after PIN reset
	//ResetPIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
}

// IVerifyPIN is used e.g to check the PIN when accessing sensitive content
type IVerifyPIN interface {
	//VerifyPIN(userID string, flavour string, pin string) (bool, error)
}

// IReviewTerms ...
type IReviewTerms interface {

	// ReviewTerms can be used to accept or review terms
	//ReviewTerms(userID string, accepted bool, flavour string) (bool, error)
}

// IAnonymizedIdentifier ...
type IAnonymizedIdentifier interface {
	// GetAnonymizedUserIdentifier is used to get an opaque (but **stable**) user
	// identifier for events, analytics etc
	//GetAnonymizedUserIdentifier(userID string, flavour string) (string, error)
}

// IAddPushToken ...
type IAddPushToken interface {
	//AddPushtoken(userID string, flavour string) (bool, error)
}

// IRemovePushtoken ...
type IRemovePushtoken interface {
	//RemovePushToken(userID string, flavour string) (bool, error)
}

// IUpdateLanguagePreferences ...
type IUpdateLanguagePreferences interface {
	//UpdateLanguagePreferences(userID string, language string) (bool, error)
}

// IUserInvite ...
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
	//Invite(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
}

// UseCasesUser group all business logic usecases related to user
type UseCasesUser interface {
	IUserInvite
	IUserForget
	ISetUserPIN
	ILogin
	IRequestDataExport
	IReviewTerms
	IAnonymizedIdentifier
	IAddPushToken
	IRemovePushtoken
	IUpdateLanguagePreferences
	IResetPIN
}

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Create        infrastructure.Create
	Query         infrastructure.Query
	Delete        infrastructure.Delete
	OnboardingExt extension.OnboardingLibraryExtension
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	onboardingExt extension.OnboardingLibraryExtension,
) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Create:        create,
		Query:         query,
		Delete:        delete,
		OnboardingExt: onboardingExt,
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
		return nil, "", err
	}

	pinData, err := us.Query.GetUserPINByUserID(ctx, *profile.ID)
	if err != nil {
		return nil, "", err
	}

	// If pin `ValidTo` field is in the past (expired), throw an error. This means the user has to
	// change their pin on the next login
	currentTime := time.Now()
	expired := utils.CheckPINExpiry(currentTime, pinData)
	if expired {
		return nil, "", exceptions.ExpiredPinError()
	}

	matched := us.OnboardingExt.ComparePIN(pin, pinData.Salt, pinData.HashedPIN, nil)
	if !matched {
		return nil, "", exceptions.PinMismatchError(err)
	}

	customToken, err := us.OnboardingExt.CreateFirebaseCustomToken(ctx, *profile.ID)
	if err != nil {
		return nil, "", err
	}

	userTokens, err := us.OnboardingExt.AuthenticateCustomFirebaseToken(customToken)
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

// ResetPIN resets user PIN
func (us *UseCasesUserImpl) ResetPIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return true, nil
}

// VerifyPIN verifies user PIN
func (us *UseCasesUserImpl) VerifyPIN(userID string, flavour string, pin string) (bool, error) {
	return true, nil
}

// ReviewTerms is used to accept or review terms
func (us *UseCasesUserImpl) ReviewTerms(userID string, accepted bool, flavour string) (bool, error) {
	return true, nil
}

// GetAnonymizedUserIdentifier is used to get an opaque (but **stable**) user
//
// identifier for events, analytics etc
func (us *UseCasesUserImpl) GetAnonymizedUserIdentifier(userID string, flavour string) (string, error) {
	return "", nil
}

// AddPushtoken adds push token to a user
func (us *UseCasesUserImpl) AddPushtoken(userID string, flavour string) (bool, error) {
	return true, nil
}

// RemovePushToken removes/retires user push token
func (us *UseCasesUserImpl) RemovePushToken(userID string, flavour string) (bool, error) {
	return true, nil
}

// UpdateLanguagePreferences updates user language preferences
func (us *UseCasesUserImpl) UpdateLanguagePreferences(userID string, language string) (bool, error) {
	return true, nil
}

// Invite sends an invite to a  user (client/staff)
// The invite contains: link to app/play store, temporary PIN that **MUST** be changed on first login
//
// The default invite channel is SMS
func (us *UseCasesUserImpl) Invite(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return false, nil
}

// SetUserPIN sets a user's PIN.
func (us *UseCasesUserImpl) SetUserPIN(ctx context.Context, input *dto.PinInput) (bool, error) {
	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return false, fmt.Errorf("failed to get a user profile by phonenumber: %v", err)
	}

	err = utils.ValidatePIN(input.PIN)
	if err != nil {
		return false, fmt.Errorf("invalid PIN provided: %v", err)
	}

	salt, encryptedPIN := us.OnboardingExt.EncryptPIN(input.PIN, nil)

	isMatch := us.OnboardingExt.ComparePIN(input.ConfirmedPin, salt, encryptedPIN, nil)
	if !isMatch {
		return false, fmt.Errorf("the provided PINs do not match")
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
		return false, err
	}

	return true, nil
}

// Forget ...
func (us *UseCasesUserImpl) Forget(userID string, pin string, flavour string) (bool, error) {
	return true, nil
}

// RequestDataExport ...
func (us *UseCasesUserImpl) RequestDataExport(userID string, pin string, flavour string) (bool, error) {
	return true, nil
}
