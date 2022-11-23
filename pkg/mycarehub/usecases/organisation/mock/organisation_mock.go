package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// OrganisationUseCaseMock mocks the implementation of organisation usecase
type OrganisationUseCaseMock struct {
	MockCreateOrganisationFn func(input *dto.OrganisationInput) (bool, error)
	MockDeleteOrganisationFn func(ctx context.Context, organisationID string) error
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
