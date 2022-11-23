package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// CheckUserRoleErr returns an error message when the user role check fails
func CheckUserRoleErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to check user role",
		Code:    int(exceptions.CheckUserRoleError),
	}
}

// UserNotAuthorizedErr returns an error message when the user is not authorized to perform the action
func UserNotAuthorizedErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user not authorized",
		Code:    int(exceptions.UserNotAuthorizedError),
	}
}

// CheckUserPermissionErr returns an error message when the user permission check fails
func CheckUserPermissionErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to check user permission",
		Code:    int(exceptions.CheckUserPermissionError),
	}
}

// AssignRolesErr returns an error message when the user role assignment fails
func AssignRolesErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to assign roles",
		Code:    int(exceptions.AssignRolesError),
	}
}

// GetUserRolesErr returns an error message when the user role retrieval fails
func GetUserRolesErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to get user roles",
		Code:    int(exceptions.GetUserRolesError),
	}
}

// GetUserPermissionsErr returns an error message when the user permission retrieval fails
func GetUserPermissionsErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to get user permissions",
		Code:    int(exceptions.GetUserPermissionsError),
	}
}

// RevokeRolesErr returns an error message when the user roles' revocation fails
func RevokeRolesErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to revoke roles",
		Code:    int(exceptions.RevokeRolesError),
	}
}

// GetAllRolesErr returns an error message when the user role retrieval fails
func GetAllRolesErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to query all roles",
		Code:    int(exceptions.GetAllRolesError),
	}
}
