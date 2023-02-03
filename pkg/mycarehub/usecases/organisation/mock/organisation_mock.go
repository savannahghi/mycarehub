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
	MockCreateOrganisationFn  func(ctx context.Context, input dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error)
	MockDeleteOrganisationFn  func(ctx context.Context, organisationID string) (bool, error)
	MockListOrganisationsFn   func(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error)
	MockSearchOrganisationFn  func(ctx context.Context, searchParameter string) ([]*domain.Organisation, error)
	MOckGetOrganisationByIDFn func(ctx context.Context, organisationID string) (*domain.Organisation, error)
}

// NewOrganisationUseCaseMock initializes a new instance mock of the organisation usecase
func NewOrganisationUseCaseMock() *OrganisationUseCaseMock {
	org := &domain.Organisation{
		ID:              uuid.New().String(),
		Active:          true,
		Code:            "123",
		Name:            "Test Organisation",
		Description:     "Test Organisation",
		EmailAddress:    gofakeit.Email(),
		PhoneNumber:     interserviceclient.TestUserPhoneNumber,
		PostalAddress:   "1234 - Moi Avenue",
		PhysicalAddress: gofakeit.Address().Address,
		DefaultCountry:  "KE",
	}

	return &OrganisationUseCaseMock{
		MockCreateOrganisationFn: func(ctx context.Context, input dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error) {
			return true, nil
		},
		MockDeleteOrganisationFn: func(ctx context.Context, organisationID string) (bool, error) {
			return true, nil
		},
		MockListOrganisationsFn: func(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
			return &dto.OrganisationOutputPage{
				Organisations: []*domain.Organisation{
					{
						ID: "123",
					},
				},
			}, nil
		},
		MockSearchOrganisationFn: func(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
			return []*domain.Organisation{
				org,
			}, nil
		},
		MOckGetOrganisationByIDFn: func(ctx context.Context, organisationID string) (*domain.Organisation, error) {
			return org, nil
		},
	}
}

// CreateOrganisation mocks the create organisation method
func (m *OrganisationUseCaseMock) CreateOrganisation(ctx context.Context, organisationInput dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error) {
	return m.MockCreateOrganisationFn(ctx, organisationInput, programInput)
}

// DeleteOrganisation mocks the delete organisation method
func (m *OrganisationUseCaseMock) DeleteOrganisation(ctx context.Context, organisationID string) (bool, error) {
	return m.MockDeleteOrganisationFn(ctx, organisationID)
}

// ListOrganisations mocks the list organisations method
func (m *OrganisationUseCaseMock) ListOrganisations(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
	return m.MockListOrganisationsFn(ctx, paginationInput)
}

// SearchOrganisation mocks the search organisation method
func (m *OrganisationUseCaseMock) SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
	return m.MockSearchOrganisationFn(ctx, searchParameter)
}

// GetOrganisationByID mocks the get organisation by id method
func (m *OrganisationUseCaseMock) GetOrganisationByID(ctx context.Context, organisationID string) (*domain.Organisation, error) {
	return m.MOckGetOrganisationByIDFn(ctx, organisationID)
}
