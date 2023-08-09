package user

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// loginFunc is an execution in the login process to perform checks and build the login response
// the returned boolean indicates whether the execution is successful or not
type loginFunc func(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool

// checks whether the user profile exists and sets it in tne response
func (us *UseCasesUserImpl) userProfileCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user, err := us.Query.GetUserProfileByUsername(ctx, credentials.Username)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := "failed to get user profile by username"
		code := exceptions.ProfileNotFound.Code()
		response.SetResponseCode(code, message)

		return false
	}

	profile := &dto.User{
		ID:                     *user.ID,
		Name:                   user.Name,
		Username:               user.Username,
		Active:                 user.Active,
		NextAllowedLogin:       *user.NextAllowedLogin,
		FailedLoginCount:       user.FailedLoginCount,
		PinChangeRequired:      user.PinChangeRequired,
		HasSetPin:              user.HasSetPin,
		HasSetSecurityQuestion: user.HasSetSecurityQuestion,
		IsPhoneVerified:        user.IsPhoneVerified,
		TermsAccepted:          user.TermsAccepted,
		Suspended:              user.Suspended,
		FailedSecurityCount:    user.FailedSecurityCount,
		PinUpdateRequired:      user.PinUpdateRequired,
		HasSetNickname:         user.HasSetNickname,
		CurrentProgramID:       user.CurrentProgramID,
	}
	response.SetUserProfile(profile)

	return true
}

// checks whether a user is active
func (us *UseCasesUserImpl) checkUserIsActive(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()
	if !user.Active {

		message := "user profile is not active"
		code := exceptions.InactiveUser.Code()
		response.SetResponseCode(code, message)

		return false
	}

	return true
}

// checks whether the client profile exists and sets it in tne response
func (us *UseCasesUserImpl) clientProfileCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		exists, err := us.Query.CheckClientExists(ctx, user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := "failed to get client profile"
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		if exists {
			response.SetIsClient()
		}

		return true

	default:
		return true
	}
}

// consumerProfilesCheck ensures that a user logging in to consumer app has the required profiles
// i.e a client or a caregiver profile exists for the user
// if no profile exists the user types required to log in don't exist
func (us *UseCasesUserImpl) consumerProfilesCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		if !response.GetIsClient() && !response.GetIsCaregiver() {

			message := "failed to get client or caregiver profile"
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		return true

	default:
		return true
	}
}

// Checks whether a user has an active PIN reset request
func (us *UseCasesUserImpl) pinResetRequestCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	userProfile := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		clientProfile, err := us.Query.GetClientProfile(ctx, userProfile.ID, userProfile.CurrentProgramID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false
		}

		serviceRequest, err := us.Query.GetClientServiceRequests(ctx, enums.ServiceRequestTypePinReset.String(), enums.ServiceRequestStatusPending.String(), *clientProfile.ID, *clientProfile.DefaultFacility.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false
		}

		if len(serviceRequest) > 0 {
			message := exceptions.PINResetServiceRequestFoundErr(err).Error()
			code := exceptions.PINResetServiceRequest.Code()
			response.SetResponseCode(code, message)

			return false
		}

		return true

	default:
		return true
	}
}

// checks whether the staff profile exists and sets it in tne response
func (us *UseCasesUserImpl) staffProfileCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourPro:
		exists, err := us.Query.CheckStaffExists(ctx, user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := exceptions.StaffProfileNotFoundErr(err).Error()
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		if !exists {
			message := exceptions.StaffProfileNotFoundErr(err).Error()
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)
		}

		return true

	default:
		return true
	}
}

// Checks whether a user as an exponential back-off that prevents them from singing in
func (us *UseCasesUserImpl) loginTimeoutCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()

	currentTime := time.Now().UTC()
	timeOutOccurred := currentTime.Before(user.NextAllowedLogin)

	if timeOutOccurred {
		loginRetryTime := user.NextAllowedLogin.Sub(currentTime).Seconds()
		response.SetRetryTime(loginRetryTime)

		message := fmt.Sprintf("please try again after %v seconds", loginRetryTime)
		code := exceptions.RetryLoginError.Code()

		response.SetResponseCode(code, message)

		return false
	}

	return true
}

func (us *UseCasesUserImpl) checkPIN(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()

	userPIN, err := us.Query.GetUserPINByUserID(ctx, user.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.PinNotFoundError(err).Error()
		code := exceptions.PINNotFound.Code()
		response.SetResponseCode(code, message)

		return false
	}

	// If pin `ValidTo` field is in the past (expired). This means the user has to change their pin
	currentTime := time.Now()
	expired := currentTime.After(userPIN.ValidTo)
	if expired {
		message := exceptions.ExpiredPinErr().Error()
		code := exceptions.ExpiredPinError.Code()
		response.SetResponseCode(code, message)

		return false
	}

	matched := utils.ComparePIN(
		credentials.PIN,
		userPIN.Salt,
		userPIN.HashedPIN,
		nil,
	)
	if !matched {
		// increment number of failed logins
		user.FailedLoginCount++

		// increase the retry time exponentially
		nextAttempt := utils.NextAllowedLoginTime(user.FailedLoginCount)
		user.NextAllowedLogin = nextAttempt

		userUpdates := map[string]interface{}{
			"failed_login_count": user.FailedLoginCount,
			"next_allowed_login": user.NextAllowedLogin,
			"last_failed_login":  time.Now().UTC(),
		}

		err := us.Update.UpdateUser(ctx, &domain.User{ID: &user.ID}, userUpdates)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			message := exceptions.InternalErr(err).Error()
			code := exceptions.Internal.Code()
			response.SetResponseCode(code, message)

			return false
		}

		response.SetFailedLoginCount(user.FailedLoginCount)
		response.SetUserProfile(user)

		message := exceptions.PinMismatchError().Error()
		code := exceptions.PINMismatch.Code()
		response.SetResponseCode(code, message)

		return false
	}

	userUpdates := map[string]interface{}{
		"last_successful_login": time.Now().UTC(),
		"failed_login_count":    0,
		"next_allowed_login":    time.Now().UTC(),
		"last_failed_login":     nil,
	}

	err = us.Update.UpdateUser(ctx, &domain.User{ID: &user.ID}, userUpdates)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.InternalErr(err).Error()
		code := exceptions.Internal.Code()
		response.SetResponseCode(code, message)

		return false
	}

	response.SetUserProfile(user)

	return true
}

func (us *UseCasesUserImpl) addRolesPermissions(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {

	return true
}

// checks whether a caregiver profile exists and assigns
// - There should be a caregiver profile
func (us *UseCasesUserImpl) caregiverProfileCheck(ctx context.Context, credentials *dto.LoginInput, response dto.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		exists, err := us.Query.CheckCaregiverExists(ctx, user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)

			// TODO: Caregiver error message
			message := exceptions.ProfileNotFoundErr(err).Error()
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		if exists {
			response.SetIsCaregiver()
		}

		return true

	}

	return true
}
