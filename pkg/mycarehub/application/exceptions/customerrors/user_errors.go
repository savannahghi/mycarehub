package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to get a user",
		Code:    int(exceptions.UserNotFound),
	}
}

// EmptyUserIDErr returns an error message when the user id is empty
func EmptyUserIDErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user id input is empty",
		Code:    int(exceptions.EmptyUserIDInputError),
	}
}

// ProfileNotFoundErr returns an error message when the profile is not found
func ProfileNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user profile not found",
		Code:    int(exceptions.ProfileNotFound),
	}
}

// NotOptedInErr returns an error message when the user is not opted in
func NotOptedInErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user not opted in",
		Code:    int(exceptions.NotOptedInError),
	}
}

// GetLoggedInUserUIDErr returns an error message when the logged in user uid check fails
func GetLoggedInUserUIDErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to get logged in user uid",
		Code:    int(exceptions.GetLoggedInUserUIDError),
	}
}

// UserNameExistsErr returns an error message when the item update fails
func UserNameExistsErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "username has already been taken",
		Code:    int(exceptions.NickNameExistsError),
	}
}

// UpdateProfileErr returns an error message when the user profile update fails
func UpdateProfileErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to update profile",
		Code:    int(exceptions.UpdateProfileError),
	}
}

// GetInviteLinkErr returns an error message when the invite link generation fails
func GetInviteLinkErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to generate invite link",
		Code:    int(exceptions.GetInviteLinkError),
	}
}
