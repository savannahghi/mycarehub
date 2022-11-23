package customerrors

import (
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
)

// InternalErr returns an error message when the server fails
func InternalErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "internal error",
		Code:    int(exceptions.Internal),
	}
}

// InvalidFlavourDefinedErr is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedErr(err error) error {
	return &exceptions.CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: "invalid flavour defined",
		Code:    int(exceptions.InvalidFlavour),
	}
}

// EncryptionErr returns an error message when the encryption fails
func EncryptionErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "encryption failed",
		Code:    int(exceptions.EncryptionError),
	}
}
