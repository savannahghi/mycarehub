package exceptions

// NonExistentOrganizationErr returns an error if the organization does not exist
func NonExistentOrganizationErr(err error) error {
	return &CustomError{
		Err:     err,
		Code:    int(NonExistentOrganizationError),
		Message: "organisation in the input is does not exist",
		Detail:  "please select another organization",
	}
}
