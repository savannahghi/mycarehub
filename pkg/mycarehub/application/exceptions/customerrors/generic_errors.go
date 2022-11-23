package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// EmptyInputErr returns an error message when an input is empty
func EmptyInputErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "input is empty",
		Code:    int(exceptions.EmptyInputError),
	}
}

// NotActiveErr returns an error message when a field is not active
func NotActiveErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "field active is false",
		Code:    int(exceptions.NotActiveError),
	}
}

// FailedToUpdateItemErr returns an error message when the item update fails
func FailedToUpdateItemErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to update item",
		Code:    int(exceptions.FailedToUpdateItemError),
	}
}

// ItemNotFoundErr returns an error message when the item is not found
func ItemNotFoundErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "item not found",
		Code:    int(exceptions.ItemNotFoundError),
	}
}

// InputValidationErr returns an error message when the input is invalid
func InputValidationErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "input validation failed",
		Code:    int(exceptions.InputValidationError),
	}
}

// FailedToSaveItemErr returns an error message when the item save fails
func FailedToSaveItemErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to save item",
		Code:    int(exceptions.FailedToSaveItemError),
	}
}
