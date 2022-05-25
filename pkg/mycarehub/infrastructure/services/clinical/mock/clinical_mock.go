package mock

import "context"

// ClinicalMock mocks the clinical service
type ClinicalMock struct {
	MockDeleteFHIRPatientByPhoneFn func(ctx context.Context, phoneNumber string) (bool, error)
}

// NewClinicalMock initializes the clinical mock service
func NewClinicalMock() *ClinicalMock {
	return &ClinicalMock{
		MockDeleteFHIRPatientByPhoneFn: func(ctx context.Context, phoneNumber string) (bool, error) {
			return true, nil
		},
	}
}

// DeleteFHIRPatientByPhone mocks the delete fhir patient by phone number
func (c *ClinicalMock) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return c.MockDeleteFHIRPatientByPhoneFn(ctx, phoneNumber)
}
