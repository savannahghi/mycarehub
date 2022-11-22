package exceptions

// OrgIDForProgramExistErr returns an error message when an organization id exist for a program that is being created
func OrgIDForProgramExistErr(err error) error {
	return &CustomError{
		Err:     err,
		Code:    int(OrgIDForProgramExistError),
		Message: "program with the selected organization exist",
		Detail:  "please select another organization for this program",
	}
}

// CreateProgramErr returns an error message when a program is not saved successfully
func CreateProgramErr(err error) error {
	return &CustomError{
		Err:     err,
		Code:    int(CreateProgramError),
		Message: "unable to create program",
		Detail:  "check your input or try again",
	}
}
