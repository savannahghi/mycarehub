package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// StaffProfileNotFoundErr returns an error message when the client profile is not found
func StaffProfileNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "staff profile not found",
		Code:    int(exceptions.ProfileNotFound),
	}
}

// StaffHasUnresolvedPinResetRequestErr returns an error message when the staff has an unresolved pin reset request
func StaffHasUnresolvedPinResetRequestErr() error {
	return &exceptions.CustomError{
		Err:     nil,
		Message: "staff has unresolved pin reset request",
		Code:    int(exceptions.StaffHasUnresolvedPinResetRequestError),
	}
}
