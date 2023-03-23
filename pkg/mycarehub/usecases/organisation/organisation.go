package organisation

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
)

// CreateOrganisation interface holds the method for creating an organisation
type CreateOrganisation interface {
	CreateOrganisation(ctx context.Context, organisationInput dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error)
}

// DeleteOrganisation interface holds the method for deleting an organisation
type DeleteOrganisation interface {
	DeleteOrganisation(ctx context.Context, organisationID string) (bool, error)
}

// ListOrganisation interface holds the method for listing organisations
type ListOrganisation interface {
	ListOrganisations(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error)
	SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error)
	GetOrganisationByID(ctx context.Context, organisationID string) (*domain.Organisation, error)
}

// UseCaseOrganisation is the interface for the organisation use case
type UseCaseOrganisation interface {
	CreateOrganisation
	DeleteOrganisation
	ListOrganisation
}

// UseCaseOrganisationImpl implements the CreateOrganisation interface
type UseCaseOrganisationImpl struct {
	Create      infrastructure.Create
	Delete      infrastructure.Delete
	Query       infrastructure.Query
	ExternalExt extension.ExternalMethodsExtension
	Pubsub      pubsubmessaging.ServicePubsub
}

// NewUseCaseOrganisationImpl creates a new instance of UseCaseOrganisationImpl
func NewUseCaseOrganisationImpl(
	create infrastructure.Create,
	delete infrastructure.Delete,
	query infrastructure.Query,
	ext extension.ExternalMethodsExtension,
	pubsub pubsubmessaging.ServicePubsub,
) UseCaseOrganisation {
	return &UseCaseOrganisationImpl{
		Create:      create,
		Delete:      delete,
		Query:       query,
		ExternalExt: ext,
		Pubsub:      pubsub,
	}
}

// CreateOrganisation creates an organisation
func (u *UseCaseOrganisationImpl) CreateOrganisation(ctx context.Context, organisationInput dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error) {
	organisation := &domain.Organisation{
		Active:          true,
		Code:            organisationInput.Code,
		Name:            organisationInput.Name,
		Description:     organisationInput.Description,
		EmailAddress:    organisationInput.EmailAddress,
		PhoneNumber:     organisationInput.PhoneNumber,
		PostalAddress:   organisationInput.PostalAddress,
		PhysicalAddress: organisationInput.PhysicalAddress,
		DefaultCountry:  organisationInput.DefaultCountry,
	}

	programs := []*domain.Program{}
	for _, program := range programInput {
		facilities := []*domain.Facility{}
		for _, facilityID := range program.Facilities {
			facilities = append(facilities, &domain.Facility{ID: &facilityID})
		}
		programs = append(programs, &domain.Program{
			Active:      true,
			Name:        program.Name,
			Description: program.Description,
			Facilities:  facilities,
		})
	}

	org, err := u.Create.CreateOrganisation(ctx, organisation, programs)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.CreateOrganisationErr(err)
	}

	randomInt, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	cmsOrganisationPayload := &dto.CreateCMSOrganisationPayload{
		OrganisationID: org.ID,
		Name:           org.Name,
		Email:          org.EmailAddress,
		PhoneNumber:    org.PhoneNumber,
		Code:           int(randomInt.Int64()),
	}

	err = u.Pubsub.NotifyCreateCMSOrganisation(ctx, cmsOrganisationPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	for _, program := range org.Programs {
		err = u.Pubsub.NotifyCreateCMSProgram(ctx, &dto.CreateCMSProgramPayload{
			ProgramID:      program.ID,
			Name:           program.Name,
			OrganisationID: org.ID,
		})
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}
	}

	return true, nil
}

// DeleteOrganisation deletes an organisation
func (u *UseCaseOrganisationImpl) DeleteOrganisation(ctx context.Context, organisationID string) (bool, error) {
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	userProfile, err := u.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	_, err = u.Query.GetStaffProfile(ctx, loggedInUserID, userProfile.CurrentProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	exists, err := u.Query.CheckOrganisationExists(ctx, organisationID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	if !exists {
		return false, exceptions.NonExistentOrganizationErr(err)
	}

	organisation := &domain.Organisation{
		ID: organisationID,
	}

	err = u.Delete.DeleteOrganisation(ctx, organisation)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// ListOrganisations lists all organisations
func (u *UseCaseOrganisationImpl) ListOrganisations(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
	var page *domain.Pagination

	if paginationInput != nil {
		if err := paginationInput.Validate(); err != nil {
			return nil, fmt.Errorf("pagination input validation failed: %v", err)
		}

		page = &domain.Pagination{
			Limit:       paginationInput.Limit,
			CurrentPage: paginationInput.CurrentPage,
		}
	}

	organisations, pageInfo, err := u.Query.ListOrganisations(ctx, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	organisationOutputPage := &dto.OrganisationOutputPage{
		Pagination:    pageInfo,
		Organisations: organisations,
	}

	return organisationOutputPage, nil
}

// SearchOrganisation searches for an organisation with the given search parameter
func (u *UseCaseOrganisationImpl) SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
	organisation, err := u.Query.SearchOrganisation(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return organisation, nil
}

// GetOrganisationByID retrieves an organisation by its ID and returns it associated data
func (u *UseCaseOrganisationImpl) GetOrganisationByID(ctx context.Context, organisationID string) (*domain.Organisation, error) {
	organisation, err := u.Query.GetOrganisation(ctx, organisationID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return organisation, nil
}
