package organisation

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// CreateOrganisation interface holds the method for creating an organisation
type CreateOrganisation interface {
	CreateOrganisation(ctx context.Context, input dto.OrganisationInput) (bool, error)
}

// UseCaseOrganisationImpl implements the CreateOrganisation interface
type UseCaseOrganisationImpl struct {
	Create infrastructure.Create
}

// NewUseCaseOrganisationImpl creates a new instance of UseCaseOrganisationImpl
func NewUseCaseOrganisationImpl(
	create infrastructure.Create,
) *UseCaseOrganisationImpl {
	return &UseCaseOrganisationImpl{
		Create: create,
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
