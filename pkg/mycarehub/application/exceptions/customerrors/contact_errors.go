package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// InvalidContactTypeErr returns an error message when the contact type is invalid
func InvalidContactTypeErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "invalid contact type",
		Code:    int(exceptions.InvalidContactTypeError),
	}
}

// NoContactsErr returns an error message when there are no contacts
func NoContactsErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user has no contacts",
		Code:    int(exceptions.NoContactsError),
	}
}

// ContactNotFoundErr returns an error message when the contact is not found
func ContactNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "contact not found",
		Code:    int(exceptions.ContactNotFoundError),
	}
}

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to normalize the msisdn",
		Code:    int(exceptions.Internal),
	}
}
