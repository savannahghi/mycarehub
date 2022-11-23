package mock

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"

// OrganisationUseCaseMock mocks the implementation of organisation usecase
type OrganisationUseCaseMock struct {
	MockCreateOrganisationFn func(input *dto.OrganisationInput) (bool, error)
}

// NewOrganisationUseCaseMock initializes a new instance mock of the organisation usecase
func NewOrganisationUseCaseMock() *OrganisationUseCaseMock {
	return &OrganisationUseCaseMock{
		MockCreateOrganisationFn: func(input *dto.OrganisationInput) (bool, error) {
			return true, nil
		},
	}
}

// CreateOrganisation mocks the create organisation method
func (m *OrganisationUseCaseMock) CreateOrganisation(input *dto.OrganisationInput) (bool, error) {
	return m.MockCreateOrganisationFn(input)
}
