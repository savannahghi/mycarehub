package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// LoginCountUpdateErr returns an error message when the login count update fails
func LoginCountUpdateErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to update login count",
		Code:    int(exceptions.LoginCountUpdateError),
	}
}

// LoginTimeUpdateErr returns an error message when the login time update fails
func LoginTimeUpdateErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to update login count",
		Code:    int(exceptions.LoginTimeUpdateError),
	}
}

// NexAllowedLoginTimeErr returns an error message when the login time update fails
func NexAllowedLoginTimeErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "login time not allowed",
		Code:    int(exceptions.NexAllowedLoginTimeError),
	}
}

// RetryLoginErr returns an error message when the user is not authorized to perform the action
func RetryLoginErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "retry login failed due to exponential backoff",
		Code:    int(exceptions.RetryLoginError),
	}
}

// FailedSecurityCountExceededErr returns an error message when the user is not authorized to verify the security question response
func FailedSecurityCountExceededErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "you have reached the maximum number of attempts",
		Code:    int(exceptions.FailedSecurityCountExceededError),
	}
}

// SecurityQuestionResponseMismatchErr returns an error message when the security question response does not match
func SecurityQuestionResponseMismatchErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "security question response does not match",
		Code:    int(exceptions.SecurityQuestionResponseMismatchError),
	}
}
