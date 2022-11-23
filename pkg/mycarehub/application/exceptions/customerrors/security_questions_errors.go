package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// SecurityQuestionNotFoundErr returns an error message when the security question is not found
func SecurityQuestionNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "security question not found",
		Code:    int(exceptions.SecurityQuestionNotFoundError),
	}
}
