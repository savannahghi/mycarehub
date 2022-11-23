package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// ClientProfileNotFoundErr returns an error message when the client profile is not found
func ClientProfileNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "client profile not found",
		Code:    int(exceptions.ProfileNotFound),
	}
}

// ClientCCCIdentifierNotFoundErr returns an error message when the client profile is not found
func ClientCCCIdentifierNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "client ccc identifier not found",
		Code:    int(exceptions.CCCIdentifierNotFoundError),
	}
}

// ClientHasUnresolvedPinResetRequestErr returns an error message when the client has an unresolved pin reset request
func ClientHasUnresolvedPinResetRequestErr() error {
	return &exceptions.CustomError{
		Err:     nil,
		Message: "client has unresolved pin reset request",
		Code:    int(exceptions.ClientHasUnresolvedPinResetRequestError),
	}
}
