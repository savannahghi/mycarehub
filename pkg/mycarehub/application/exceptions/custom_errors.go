package exceptions

import "github.com/savannahghi/errorcodeutil"

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINMismatchErrMsg,
		Code:    int(errorcodeutil.PINMismatch),
	}
}

// ExpiredPinError displays an error when the pin
// is expired
func ExpiredPinError() error {
	return &errorcodeutil.CustomError{
		Message: ExpiredPinErrMsg,
		Code:    10,
	}
}
