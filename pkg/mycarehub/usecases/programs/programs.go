package programs

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICreatePrograms creates the programs
type ICreatePrograms interface {
	CreateProgram(ctx context.Context, input *dto.ProgramInput) (bool, error)
}

// IListPrograms listing programs
type IListPrograms interface {
	ListUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error)
	SetCurrentProgram(ctx context.Context, programID string) (bool, error)
}

// UsecasePrograms groups al the interfaces for the Programs usecase
type UsecasePrograms interface {
	ICreatePrograms
	IListPrograms
}

// UsecaseProgramsImpl represents the Programs implementation
type UsecaseProgramsImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
}

// NewUsecasePrograms is the controller function for the Programs usecase
func NewUsecasePrograms(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	ext extension.ExternalMethodsExtension,
) UsecasePrograms {
	return &UsecaseProgramsImpl{
		Query:       query,
		Create:      create,
		Update:      update,
		ExternalExt: ext,
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

// ListUserPrograms lists the programs a user is part of in an organisation
// The user can then select and select the program
func (u *UsecaseProgramsImpl) ListUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	_, err := u.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(err)
	}

	programs, err := u.Query.GetUserPrograms(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	return programs, nil
}

// SetCurrentProgram sets the program that the user has selected from their programs
func (u *UsecaseProgramsImpl) SetCurrentProgram(ctx context.Context, programID string) (bool, error) {
	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	user, err := u.Query.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotFoundError(err)
	}

	update := map[string]interface{}{
		"current_program_id": programID,
	}
	err = u.Update.UpdateUser(ctx, user, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotFoundError(err)
	}

	return true, nil
}
