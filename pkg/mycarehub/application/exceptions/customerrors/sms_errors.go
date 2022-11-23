package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// SendSMSErr returns an error message when the SMS sending fails
func SendSMSErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to send SMS",
		Code:    int(exceptions.SendSMSError),
	}
}
