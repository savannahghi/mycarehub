package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// OrganisationUseCaseMock mocks the implementation of organisation usecase
type OrganisationUseCaseMock struct {
	MockCreateOrganisationFn func(input *dto.OrganisationInput) (bool, error)
	MockDeleteOrganisationFn func(ctx context.Context, organisationID string) error
	MockListOrganisationsFn  func(ctx context.Context) ([]*domain.Organisation, error)
}

// NewOrganisationUseCaseMock initializes a new instance mock of the organisation usecase
func NewOrganisationUseCaseMock() *OrganisationUseCaseMock {
	return &OrganisationUseCaseMock{
		MockCreateOrganisationFn: func(input *dto.OrganisationInput) (bool, error) {
			return true, nil
		},
		MockDeleteOrganisationFn: func(ctx context.Context, organisationID string) error {
			return nil
		},
		MockListOrganisationsFn: func(ctx context.Context) ([]*domain.Organisation, error) {
			return []*domain.Organisation{}, nil
		},
	}
}

// CreateOrganisation mocks the create organisation method
func (m *OrganisationUseCaseMock) CreateOrganisation(input *dto.OrganisationInput) (bool, error) {
	return m.MockCreateOrganisationFn(input)
}

// DeleteOrganisation mocks the delete organisation method
func (m *OrganisationUseCaseMock) DeleteOrganisation(ctx context.Context, organisationID string) error {
	return m.MockDeleteOrganisationFn(ctx, organisationID)
}

// ListOrganisations mocks the list organisations method
func (m *OrganisationUseCaseMock) ListOrganisations(ctx context.Context) ([]*domain.Organisation, error) {
	return m.MockListOrganisationsFn(ctx)
}
