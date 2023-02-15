package programs

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
)

// ICreatePrograms creates the programs
type ICreatePrograms interface {
	CreateProgram(ctx context.Context, input *dto.ProgramInput) (bool, error)
}

// IListPrograms listing programs
type IListPrograms interface {
	ListUserPrograms(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error)
	SetCurrentProgram(ctx context.Context, programID string) (bool, error)
	GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error)
	ListPrograms(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error)
	SearchPrograms(ctx context.Context, searchParameter string) ([]*domain.Program, error)
	GetProgramByID(ctx context.Context, programID string) (*domain.Program, error)
	CmdListPrograms(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error)
}

// IUpdatePrograms updates programs
type IUpdatePrograms interface {
	SetStaffProgram(ctx context.Context, programID string) (*domain.StaffResponse, error)
	SetClientProgram(ctx context.Context, programID string) (*domain.ClientResponse, error)
}

// UsecasePrograms groups al the interfaces for the Programs usecase
type UsecasePrograms interface {
	ICreatePrograms
	IListPrograms
	IUpdatePrograms
}

// UsecaseProgramsImpl represents the Programs implementation
type UsecaseProgramsImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
	GetStream   getstream.ServiceGetStream
	Pubsub      pubsubmessaging.ServicePubsub
}

// NewUsecasePrograms is the controller function for the Programs usecase
func NewUsecasePrograms(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	ext extension.ExternalMethodsExtension,
	getStream getstream.ServiceGetStream,
	pubsub pubsubmessaging.ServicePubsub,
) UsecasePrograms {
	return &UsecaseProgramsImpl{
		Query:       query,
		Create:      create,
		Update:      update,
		ExternalExt: ext,
		GetStream:   getStream,
		Pubsub:      pubsub,
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

	program, err := u.Create.CreateProgram(ctx, input)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, exceptions.CreateProgramErr(err)
	}

	_, err = u.Create.AddFacilityToProgram(ctx, program.ID, input.Facilities)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to add facilities to program:%w", err))
		return false, fmt.Errorf("failed to add facilities to program:%w", err)
	}

	cmsProgramPayload := &dto.CreateCMSProgramPayload{
		ProgramID:      program.ID,
		Name:           program.Name,
		OrganisationID: program.Organisation.ID,
	}

	err = u.Pubsub.NotifyCreateCMSProgram(ctx, cmsProgramPayload)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return false, err
	}

	return true, nil
}

// ListUserPrograms lists the programs a user is part of in an organisation
func (u *UsecaseProgramsImpl) ListUserPrograms(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
	_, err := u.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(err)
	}

	programOutput := &dto.ProgramOutput{}

	switch flavour {
	case feedlib.FlavourPro:
		programs, err := u.Query.GetStaffUserPrograms(ctx, userID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.GetLoggedInUserUIDErr(err)
		}
		programOutput = &dto.ProgramOutput{
			Count:    len(programs),
			Programs: programs,
		}
	case feedlib.FlavourConsumer:
		programs, err := u.Query.GetClientUserPrograms(ctx, userID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.GetLoggedInUserUIDErr(err)
		}

		programOutput = &dto.ProgramOutput{
			Count:    len(programs),
			Programs: programs,
		}
	}
	return programOutput, nil
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

// GetProgramFacilities gets the facilities that belong the program
func (u *UsecaseProgramsImpl) GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error) {
	return u.Query.GetProgramFacilities(ctx, programID)
}

// SetStaffProgram sets the program that the staff user has selected from their programs
func (u *UsecaseProgramsImpl) SetStaffProgram(ctx context.Context, programID string) (*domain.StaffResponse, error) {
	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	programStaffProfile, err := u.Query.GetStaffProfile(ctx, uid, programID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	program, err := u.Query.GetProgramByID(ctx, programStaffProfile.ProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	update := map[string]interface{}{
		"current_program_id":      programID,
		"current_organisation_id": program.Organisation.ID,
	}

	err = u.Update.UpdateUser(ctx, programStaffProfile.User, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	communityToken, err := u.GetStream.CreateGetStreamUserToken(ctx, *programStaffProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.StaffResponse{
		StaffProfile:   *programStaffProfile,
		Roles:          []*domain.AuthorityRole{},
		Permissions:    []*domain.AuthorityPermission{},
		CommunityToken: communityToken,
	}, nil
}

// SetClientProgram sets the program that the client user has selected from their programs
func (u *UsecaseProgramsImpl) SetClientProgram(ctx context.Context, programID string) (*domain.ClientResponse, error) {
	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	programClientProfile, err := u.Query.GetClientProfile(ctx, uid, programID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(err)
	}

	program, err := u.Query.GetProgramByID(ctx, programClientProfile.ProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	update := map[string]interface{}{
		"current_program_id":      programID,
		"current_organisation_id": program.Organisation.ID,
	}

	err = u.Update.UpdateUser(ctx, programClientProfile.User, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	communityToken, err := u.GetStream.CreateGetStreamUserToken(ctx, *programClientProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.ClientResponse{
		ClientProfile:  programClientProfile,
		Roles:          []*domain.AuthorityRole{},
		Permissions:    []*domain.AuthorityPermission{},
		CommunityToken: communityToken,
	}, nil
}

// ListPrograms is responsible for returning a list of paginated facilities
func (u *UsecaseProgramsImpl) ListPrograms(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
	var page *domain.Pagination

	if paginationsInput != nil {
		if err := paginationsInput.Validate(); err != nil {
			return nil, fmt.Errorf("pagination input validation failed: %v", err)
		}

		page = &domain.Pagination{
			Limit:       paginationsInput.Limit,
			CurrentPage: paginationsInput.CurrentPage,
		}
	}

	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	userProfile, err := u.Query.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	programs, pageInfo, err := u.Query.ListPrograms(ctx, &userProfile.CurrentOrganizationID, page)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return nil, err
	}

	return &domain.ProgramPage{
		Pagination: *pageInfo,
		Programs:   programs,
	}, nil
}

// SearchPrograms is used to search for programs from the organisation of the currently logged in user
func (u *UsecaseProgramsImpl) SearchPrograms(ctx context.Context, searchParameter string) ([]*domain.Program, error) {
	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	user, err := u.Query.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(err)
	}

	programs, err := u.Query.SearchPrograms(ctx, searchParameter, user.CurrentOrganizationID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return programs, nil
}

// GetProgramByID retrieves a program from the database using the provided program id
func (u *UsecaseProgramsImpl) GetProgramByID(ctx context.Context, programID string) (*domain.Program, error) {
	program, err := u.Query.GetProgramByID(ctx, programID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return program, nil
}

// CmdListPrograms is responsible for returning a list of paginated programs to be used in cmd
func (u *UsecaseProgramsImpl) CmdListPrograms(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
	var page *domain.Pagination

	if paginationsInput != nil {
		if err := paginationsInput.Validate(); err != nil {
			return nil, fmt.Errorf("pagination input validation failed: %v", err)
		}

		page = &domain.Pagination{
			Limit:       paginationsInput.Limit,
			CurrentPage: paginationsInput.CurrentPage,
		}
	}

	programs, pageInfo, err := u.Query.ListPrograms(ctx, &organisationID, page)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return nil, err
	}

	return &domain.ProgramPage{
		Pagination: *pageInfo,
		Programs:   programs,
	}, nil
}
