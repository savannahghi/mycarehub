package mock

import "context"

// ClinicalMock mocks the clinical service implementations
type ClinicalMock struct {
	MockDeleteFHIRPatientByPhoneFn func(ctx context.Context, phoneNumber string) error
}

// NewClinicalServiceMock initializes the clinical mock service
func NewClinicalServiceMock() *ClinicalMock {
	return &ClinicalMock{
		MockDeleteFHIRPatientByPhoneFn: func(ctx context.Context, phoneNumber string) error {
			return nil
		},
	}
}

// DeleteFHIRPatientByPhone deletes a patient by phone number
func (c *ClinicalMock) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) error {
	return c.MockDeleteFHIRPatientByPhoneFn(ctx, phoneNumber)
}
