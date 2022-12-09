package organisation

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// CreateOrganisation interface holds the method for creating an organisation
type CreateOrganisation interface {
	CreateOrganisation(ctx context.Context, input dto.OrganisationInput) (bool, error)
}

// DeleteOrganisation interface holds the method for deleting an organisation
type DeleteOrganisation interface {
	DeleteOrganisation(ctx context.Context, organisationID string) (bool, error)
}

// ListOrganisation interface holds the method for listing organisations
type ListOrganisation interface {
	ListOrganisations(ctx context.Context) ([]*domain.Organisation, error)
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
}

// NewUseCaseOrganisationImpl creates a new instance of UseCaseOrganisationImpl
func NewUseCaseOrganisationImpl(
	create infrastructure.Create,
	delete infrastructure.Delete,
	query infrastructure.Query,
	ext extension.ExternalMethodsExtension,
) *UseCaseOrganisationImpl {
	return &UseCaseOrganisationImpl{
		Create:      create,
		Delete:      delete,
		Query:       query,
		ExternalExt: ext,
	}
}

// CreateOrganisation creates an organisation
func (u *UseCaseOrganisationImpl) CreateOrganisation(ctx context.Context, input dto.OrganisationInput) (bool, error) {
	organisation := &domain.Organisation{
		Active:           true,
		OrganisationCode: input.OrganisationCode,
		Name:             input.Name,
		Description:      input.Description,
		EmailAddress:     input.EmailAddress,
		PhoneNumber:      input.PhoneNumber,
		PostalAddress:    input.PostalAddress,
		PhysicalAddress:  input.PhysicalAddress,
		DefaultCountry:   input.DefaultCountry,
	}

	err := u.Create.CreateOrganisation(ctx, organisation)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.CreateOrganisationErr(err)
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

	_, err = u.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
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
func (u *UseCaseOrganisationImpl) ListOrganisations(ctx context.Context) ([]*domain.Organisation, error) {
	organisations, err := u.Query.ListOrganisations(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return organisations, nil
}
