package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

// loginFunc is an execution in the login process to perform checks and build the login response
// the returned boolean indicates whether the execution is successful or not
type loginFunc func(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool

// checks whether the user profile exists and sets it in tne response
func (us *UseCasesUserImpl) userProfileCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	userProfile, err := us.Query.GetUserProfileByPhoneNumber(ctx, *credentials.PhoneNumber, credentials.Flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := "failed to get user profile by phone number"
		code := exceptions.ProfileNotFound.Code()
		response.SetResponseCode(code, message)

		return false
	}

	response.SetUserProfile(userProfile)

	return true
}

// checks whether a user is active
func (us *UseCasesUserImpl) checkUserIsActive(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
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
func (us *UseCasesUserImpl) clientProfileCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		// If the user is only a caregiver
		// don't proceed with client logic
		if user.UserType == enums.CaregiverUser &&
			response.GetCaregiverProfile() != nil {
			return true
		}

		clientProfile, err := us.Query.GetClientProfileByUserID(ctx, *user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := "failed to get client profile"
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		cccIdentifier, err := us.Query.GetClientCCCIdentifier(ctx, *clientProfile.ID)
		if err != nil {
			// Do not lock the user(client) out of the
			// app if their identifier is not found.
			//
			// A workflow for this should be developed.
			helpers.ReportErrorToSentry(err)
		}
		// because the error above is ignored
		if cccIdentifier != nil {
			clientProfile.CCCNumber = cccIdentifier.IdentifierValue
		}

		if clientProfile.CHVUserID != nil {
			CHVProfile, err := us.Query.GetUserProfileByUserID(ctx, *clientProfile.CHVUserID)
			if err != nil {
				helpers.ReportErrorToSentry(err)

				message := "failed to get CHV profile"
				code := exceptions.UserNotFound.Code()
				response.SetResponseCode(code, message)

				return false
			}

			clientProfile.CHVUserName = CHVProfile.Name
		}

		response.SetClientProfile(clientProfile)

		return true
	default:
		return true
	}
}

// Checks whether a user has an active PIN reset request
func (us *UseCasesUserImpl) pinResetRequestCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	var exists bool

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		// If the user is only a caregiver
		// don't proceed with client logic
		if response.GetUserProfile().UserType == enums.CaregiverUser &&
			response.GetCaregiverProfile() != nil {
			return true
		}

		client := response.GetClientProfile()

		status, err := us.Query.CheckIfClientHasUnresolvedServiceRequests(ctx, *client.ID, enums.ServiceRequestTypePinReset.String())
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := "failed to check if client has unresolved pin reset request"
			code := exceptions.Internal.Code()
			response.SetResponseCode(code, message)

			return false
		}

		exists = status

	case feedlib.FlavourPro:
		staff := response.GetStaffProfile()

		status, err := us.Query.CheckIfStaffHasUnresolvedServiceRequests(ctx, *staff.ID, enums.ServiceRequestTypeStaffPinReset.String())
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := "failed to check if staff has unresolved pin reset request"
			code := exceptions.Internal.Code()
			response.SetResponseCode(code, message)

			return false
		}

		exists = status

	}

	if exists {
		message := exceptions.ClientHasUnresolvedPinResetRequestErr().Error()
		code := exceptions.ClientHasUnresolvedPinResetRequestError.Code()
		response.SetResponseCode(code, message)

		return false
	}

	return true
}

// checks whether the staff profile exists and sets it in tne response
func (us *UseCasesUserImpl) staffProfileCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourPro:
		staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, *user.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)

			message := exceptions.StaffProfileNotFoundErr(err).Error()
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		response.SetStaffProfile(staffProfile)
		return true

	default:
		return true
	}
}

// Checks whether a user as an exponential back-off that prevents them from singing in
func (us *UseCasesUserImpl) loginTimeoutCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	currentTime := time.Now()
	timeOutOccurred := currentTime.Before(*user.NextAllowedLogin)

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

func (us *UseCasesUserImpl) checkPIN(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	userPIN, err := us.Query.GetUserPINByUserID(ctx, *user.ID, credentials.Flavour)
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

	matched := us.ExternalExt.ComparePIN(
		*credentials.PIN,
		userPIN.Salt,
		userPIN.HashedPIN,
		nil,
	)
	if !matched {
		// increment number of failed logins
		user.FailedLoginCount++

		// increase the retry time exponentially
		nextAttempt := utils.NextAllowedLoginTime(user.FailedLoginCount)
		user.NextAllowedLogin = &nextAttempt

		userUpdates := map[string]interface{}{
			"failed_login_count": user.FailedLoginCount,
			"next_allowed_login": user.NextAllowedLogin,
			"last_failed_login":  time.Now(),
		}

		err := us.Update.UpdateUser(ctx, user, userUpdates)
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
		"last_successful_login": time.Now(),
		"failed_login_count":    0,
		"next_allowed_login":    time.Now(),
		"last_failed_login":     nil,
	}

	err = us.Update.UpdateUser(ctx, user, userUpdates)
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

func (us *UseCasesUserImpl) addAuthCredentials(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	customToken, err := us.ExternalExt.CreateFirebaseCustomToken(ctx, *user.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.InternalErr(err).Error()
		code := exceptions.Internal.Code()
		response.SetResponseCode(code, message)

		return false
	}

	userTokens, err := us.ExternalExt.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.InternalErr(err).Error()
		code := exceptions.Internal.Code()
		response.SetResponseCode(code, message)

		return false
	}

	creds := domain.AuthCredentials{
		RefreshToken: userTokens.RefreshToken,
		IDToken:      userTokens.IDToken,
		ExpiresIn:    userTokens.ExpiresIn,
	}

	response.SetAuthCredentials(creds)

	return true
}

func (us *UseCasesUserImpl) addGetStreamToken(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	var user *stream.User
	userProfile := response.GetUserProfile()

	userAge := utils.CalculateAge(*userProfile.DateOfBirth)

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		// If the user is only a caregiver
		// don't proceed with client logic
		if userProfile.UserType == enums.CaregiverUser &&
			response.GetCaregiverProfile() != nil {
			return true
		}

		client := response.GetClientProfile()

		user = &stream.User{
			ID:   *client.ID,
			Name: userProfile.Name,
			ExtraData: map[string]interface{}{
				"userID":        userProfile.ID,
				"userType":      enums.ClientUser.String(),
				"username":      userProfile.Username,
				"ageUpperBound": userAge,
				"ageLowerBound": userAge,
				"clientTypes":   client.ClientTypes,
				"gender":        userProfile.Gender,
			},
		}

	case feedlib.FlavourPro:
		staff := response.GetStaffProfile()
		user = &stream.User{
			ID:   *staff.ID,
			Name: userProfile.Name,
			ExtraData: map[string]interface{}{
				"userID":   userProfile.ID,
				"userType": enums.StaffUser.String(),
				"username": userProfile.Username,
			},
		}

	}

	_, err := us.GetStream.CreateGetStreamUser(ctx, user)
	if err != nil {
		helpers.ReportErrorToSentry(err)
	}

	token, err := us.GetStream.CreateGetStreamUserToken(ctx, user.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
	}

	response.SetStreamToken(token)

	return true
}

func (us *UseCasesUserImpl) addRolesPermissions(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	roles, err := us.Authority.GetUserRoles(ctx, *user.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.GetUserRolesErr(err).Error()
		code := exceptions.GetUserRolesError.Code()
		response.SetResponseCode(code, message)

		return false
	}
	response.SetRoles(roles)

	permissions, err := us.Authority.GetUserPermissions(ctx, *user.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)

		message := exceptions.GetUserPermissionsErr(err).Error()
		code := exceptions.GetUserPermissionsError.Code()
		response.SetResponseCode(code, message)

		return false
	}
	response.SetPermissions(permissions)

	return true
}

// checks whether a caregiver profile exists and assigns
// - There should be a caregiver profile
func (us *UseCasesUserImpl) caregiverProfileCheck(ctx context.Context, credentials *dto.LoginInput, response domain.ILoginResponse) bool {
	user := response.GetUserProfile()

	switch credentials.Flavour {
	case feedlib.FlavourConsumer:
		// 1. Check if caregiver profile exists
		caregiver, err := us.Query.GetCaregiverByUserID(ctx, *user.ID)
		if err != nil {
			// User is a client without caregiver profile proceed i.e client that isn't a caregiver
			if errors.Is(err, gorm.ErrRecordNotFound) && user.UserType == enums.ClientUser {
				return true
			}

			helpers.ReportErrorToSentry(err)

			// TODO: Caregiver error message
			message := exceptions.ProfileNotFoundErr(err).Error()
			code := exceptions.ProfileNotFound.Code()
			response.SetResponseCode(code, message)

			return false
		}

		profile := &domain.CaregiverProfile{
			ID:              caregiver.ID,
			UserID:          caregiver.UserID,
			User:            *user,
			CaregiverNumber: caregiver.CaregiverNumber,
		}

		response.SetCaregiverProfile(profile)
	}

	return true
}
