package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// OrgIDForProgramExistErr returns an error message when an organization id exist for a program that is being created
func OrgIDForProgramExistErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "program with the selected organization exist",
		Code:    int(exceptions.OrgIDForProgramExistError),
		Detail:  "please select another organization for this program",
	}
}

// CreateProgramErr returns an error message when a program is not saved successfully
func CreateProgramErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to create program",
		Code:    int(exceptions.CreateProgramError),
		Detail:  "check your input or try again",
	}
}
