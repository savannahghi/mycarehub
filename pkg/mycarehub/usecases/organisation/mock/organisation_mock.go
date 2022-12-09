package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// OrganisationUseCaseMock mocks the implementation of organisation usecase
type OrganisationUseCaseMock struct {
	MockCreateOrganisationFn  func(input *dto.OrganisationInput) (bool, error)
	MockDeleteOrganisationFn  func(ctx context.Context, organisationID string) error
	MockListOrganisationsFn   func(ctx context.Context) ([]*domain.Organisation, error)
	MockSearchOrganisationsFn func(ctx context.Context, searchParameter string, country string) ([]*domain.Organisation, error)
}

// NewOrganisationUseCaseMock initializes a new instance mock of the organisation usecase
func NewOrganisationUseCaseMock() *OrganisationUseCaseMock {
	org := &domain.Organisation{
		ID:               uuid.New().String(),
		Active:           true,
		OrganisationCode: "12345",
		Name:             "Test Organisation",
		Description:      "Test Organisation",
		EmailAddress:     gofakeit.Email(),
		PhoneNumber:      interserviceclient.TestUserPhoneNumber,
		PostalAddress:    "Test Address",
		PhysicalAddress:  "Test Address",
		DefaultCountry:   "KE",
	}

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
		MockSearchOrganisationsFn: func(ctx context.Context, searchParameter string, country string) ([]*domain.Organisation, error) {
			return []*domain.Organisation{
				org,
			}, nil
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

// SearchOrganisations mocks the search organisation method
func (m *OrganisationUseCaseMock) SearchOrganisations(ctx context.Context, searchParameter string, country string) ([]*domain.Organisation, error) {
	return m.MockSearchOrganisationsFn(ctx, searchParameter, country)
}
