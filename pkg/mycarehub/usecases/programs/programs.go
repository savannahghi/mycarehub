package programs

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICreatePrograms creates the programs
type ICreatePrograms interface {
	CreateProgram(ctx context.Context, input *dto.ProgramInput) (bool, error)
}

// UsecasePrograms groups al the interfaces for the Programs usecase
type UsecasePrograms interface {
	ICreatePrograms
}

// UsecaseProgramsImpl represents the Programs implementation
type UsecaseProgramsImpl struct {
	Query  infrastructure.Query
	Create infrastructure.Create
}

// NewUsecasePrograms is the controller function for the Programs usecase
func NewUsecasePrograms(
	query infrastructure.Query,
	create infrastructure.Create,
) UsecasePrograms {
	return &UsecaseProgramsImpl{
		Query:  query,
		Create: create,
	}
}

// CreateProgram adds a new record of programs
// the program name should be unique for each program in a given organization
func (u *UsecaseProgramsImpl) CreateProgram(ctx context.Context, input *dto.ProgramInput) (bool, error) {
	if err := input.Validate(); err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.InputValidationErr(err)
	}

	exists, err := u.Query.CheckOrganisationExists(ctx, input.OrganisationID)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.InternalErr(err)
	}
	if !exists {
		err := fmt.Errorf("organisation with ID %s does not exist", input.OrganisationID)
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.NonExistentOrganizationErr(err)
	}

	exists, err = u.Query.CheckIfProgramNameExists(ctx, input.OrganisationID, input.Name)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.InternalErr(err)
	}
	if exists {
		err := fmt.Errorf("a program with organisation ID %s and name %s already exists", input.OrganisationID, input.Name)
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.OrgIDForProgramExistErr(err)
	}

	err = u.Create.CreateProgram(ctx, input)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.CreateProgramErr(err)
	}

	return true, nil
}
