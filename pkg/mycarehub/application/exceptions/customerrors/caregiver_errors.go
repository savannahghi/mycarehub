package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// UpdateClientCaregiverErr returns an error message when client caregiver update fails
func UpdateClientCaregiverErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to update caregiver",
		Code:    int(exceptions.UpdateClientCaregiverError),
	}
}

// CreateClientCaregiverErr returns an error message when the client caregiver creation fails
func CreateClientCaregiverErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to create caregiver",
		Code:    int(exceptions.CreateClientCaregiverError),
	}
}
