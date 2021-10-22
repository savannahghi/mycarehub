package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
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
	Login(ctx context.Context, userID string, pin string, flavour string) (*domain.AuthCredentials, string, error)
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
	Forget(userID string, pin string, flavour string) (bool, error)
}

// IRequestDataExport allows a user to request data known about them
// Mostly for legal compliance.
// The first impl. will simply create a task (for manual follow up) and acknowledge
type IRequestDataExport interface {
	RequestDataExport(userID string, pin string, flavour string) (bool, error)
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
	SetUserPIN(ctx context.Context, input *dto.PINInput) (bool, error)
}

// ResetPIN ...
type ResetPIN interface {
	// ResetPIN can be used by admins or healthcare workers to generate and send
	// a new PIN for a client or other user.
	// The new PIN is generated automatically and set to expire immediately so
	// that a PIN change is forced on next login.
	// TODO: Notify user after PIN reset
	ResetPIN(userID string, flavour string) (bool, error)
}

// IVerifyPIN is used e.g to check the PIN when accessing sensitive content
type IVerifyPIN interface {
	VerifyPIN(userID string, flavour string, pin string) (bool, error)
}

// IReviewTerms ...
type IReviewTerms interface {

	// ReviewTerms can be used to accept or review terms
	ReviewTerms(userID string, accepted bool, flavour string) (bool, error)
}

// IAnonymizedIdentifier ...
type IAnonymizedIdentifier interface {
	// GetAnonymizedUserIdentifier is used to get an opaque (but **stable**) user
	// identifier for events, analytics etc
	GetAnonymizedUserIdentifier(userID string, flavour string) (string, error)
}

// IAddPushToken ...
type IAddPushToken interface {
	AddPushtoken(userID string, flavour string) (bool, error)
}

// IRemovePushtoken ...
type IRemovePushtoken interface {
	RemovePushToken(userID string, flavour string) (bool, error)
}

// IUpdateLanguagePreferences ...
type IUpdateLanguagePreferences interface {
	UpdateLanguagePreferences(userID string, language string) (bool, error)
}

// IUserInvite ...
type IUserInvite interface {

	// TODO: flavour is an enum; client or pro app
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
	Invite(userID string, flavour string) (bool, error)
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
}

// UseCasesUserImpl represents user implementation object
type UseCasesUserImpl struct {
	Infrastructure infrastructure.Interactor
}

// NewUseCasesUserImpl returns a new user service
func NewUseCasesUserImpl(infra infrastructure.Interactor) *UseCasesUserImpl {
	return &UseCasesUserImpl{
		Infrastructure: infra,
	}
}

// Login is used to login the user into the application
func (us *UseCasesUserImpl) Login(ctx context.Context, userID string, pin string, flavour string) (*domain.AuthCredentials, string, error) {
	// Get user profile by UserID
	userProfile, err := us.Infrastructure.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get user profile by userID: %v", err)
	}

	//Fetch PIN by UserID
	userPINData, err := us.Infrastructure.GetUserPINByUserID(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get user PIN by userID: %v", err)
	}

	//Compare PIN to check for validity
	isMatch := extension.ComparePIN(pin, userPINData.Salt, userPINData.HashedPIN, nil)
	// On any mis-match:
	if !isMatch {
		failedLoginCount, err := strconv.Atoi(userProfile.FailedLoginCount)
		if err != nil {
			return nil, "", err
		}
		trials := failedLoginCount + 1
		//Convert trials to string
		numberOfTrials := strconv.Itoa(trials)

		// Implement exponential back-off, record the number of trials and last successful login
		// 1. Record the user's number of failed login time
		if err := us.Infrastructure.UpdateUserFailedLoginCount(ctx, userID, numberOfTrials, flavour); err != nil {
			return nil, "unable to update number of user failed login counts", fmt.Errorf("unable to update number of user failed login counts: %v", err)
		}

		// 2. Update user's last failed login time
		lastFailedLoginTime := time.Now()
		if err := us.Infrastructure.UpdateUserLastFailedLogin(ctx, userID, lastFailedLoginTime, flavour); err != nil {
			return nil, "unable to update number of user last failed login time", fmt.Errorf("unable to update number of user last failed login time: %v", err)
		}

		// TODO: Implement exponential back-off

		return nil, "", fmt.Errorf("an error occurred")
	}

	// Update last successful login
	currentTime := time.Now()
	if err := us.Infrastructure.UpdateUserLastSuccessfulLogin(ctx, userID, currentTime, flavour); err != nil {
		return nil, "", fmt.Errorf("unable to update user last successful login: %v", err)
	}

	// Reset failed login count to zero after successful login
	if err := us.Infrastructure.UpdateUserFailedLoginCount(ctx, userID, "0", flavour); err != nil {
		return nil, "unable to reset the number of failed user's failed login attempts", fmt.Errorf("unable to reset the number of failed user's failed login attempts: %v", err)
	}

	customToken, err := firebasetools.CreateFirebaseCustomToken(ctx, *userProfile.ID)
	if err != nil {
		return nil, "", exceptions.CustomTokenError(err)
	}

	userTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, "", exceptions.AuthenticateTokenError(err)
	}

	//Generate Authentication credentials
	authCredentials := &domain.AuthCredentials{
		User:         userProfile,
		RefreshToken: userTokens.RefreshToken,
		IDToken:      userTokens.IDToken,
		ExpiresIn:    userTokens.ExpiresIn,
	}

	//Return authentication credentials, string and an error.
	return authCredentials, "login successful", nil
}

// ResetPIN resets user PIN
func (us *UseCasesUserImpl) ResetPIN(userID string, flavour string) (bool, error) {
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

// Invite is sends an invite to a  user (client/staff)
// The invite contains: link to app/play store, temporary PIN that **MUST** be changed on first login
//
// The default invite channel is SMS
func (us *UseCasesUserImpl) Invite(userID string, flavour string) (bool, error) {
	return false, nil
}

// SetUserPIN sets a user's PIN.
func (us *UseCasesUserImpl) SetUserPIN(ctx context.Context, input *dto.PINInput) (bool, error) {
	//Get user profile PIN

	err := utils.ValidatePIN(input.PIN)
	if err != nil {
		return false, fmt.Errorf("invalid PIN provided: %v", err)
	}
	salt, encryptedPIN := extension.EncryptPIN(input.PIN, nil)

	isMatch := extension.ComparePIN(input.ConfirmedPin, salt, encryptedPIN, nil)
	if !isMatch {
		return false, fmt.Errorf("the provided PINs do not match")
	}

	// Ensure that the PIN is only valid for the next 24 hours.
	validTo := utils.GetHourMinuteSecond(24, 0, 0)

	pinDataInput := &domain.UserPIN{
		UserID:    "2c301ec0-7ee2-4b65-8e9f-9b99756a1072",
		HashedPIN: encryptedPIN,
		ValidFrom: time.Now(),
		ValidTo:   validTo, // Consult for appropriate timings
		Flavour:   input.Flavour,
		IsValid:   isMatch,
		Salt:      salt,
	}

	return us.Infrastructure.SetUserPIN(ctx, pinDataInput)
}

// Forget ...
func (us *UseCasesUserImpl) Forget(userID string, pin string, flavour string) (bool, error) {
	return true, nil
}

// RequestDataExport ...
func (us *UseCasesUserImpl) RequestDataExport(userID string, pin string, flavour string) (bool, error) {
	return true, nil
}
