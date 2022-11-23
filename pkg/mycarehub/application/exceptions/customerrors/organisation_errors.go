package customerrors

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

// CreateOrganisationErr returns an error message when the organisation creation fails
func CreateOrganisationErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to create an organization",
		Code:    int(exceptions.FailToCreateOrganisationError),
		Detail: "The organization could not be created. Please try again with a different name, organisation code, email or phone number. " +
			"If the problem persists, please contact support.",
	}
}

// NonExistentOrganizationErr returns an error if the organization does not exist
func NonExistentOrganizationErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "organisation in the input is does not exist",
		Code:    int(exceptions.NonExistentOrganizationError),
		Detail:  "please select another organization",
	}
}
