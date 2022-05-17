package clinical

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
)

var (
	deleteFHIRPatientEndpoint = "internal/delete-patient"
)

// ServiceClinical defines various clinical methods
type ServiceClinical interface {
	DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error)
}

// ServiceClinicalImpl represents clinical implementation objects
type ServiceClinicalImpl struct {
	ClinicalISCClient extension.ISCClientExtension
	ExternalExt       extension.ExternalMethodsExtension
}

// NewServiceClinical returns new instance of ServiceClinicalImpl
func NewServiceClinical(
	clinicalISCClient extension.ISCClientExtension,
	externalExt extension.ExternalMethodsExtension,
) ServiceClinical {
	return &ServiceClinicalImpl{
		ClinicalISCClient: clinicalISCClient,
		ExternalExt:       externalExt,
	}
}

// DeleteFHIRPatientByPhone deletes a patient by phone number
func (c *ServiceClinicalImpl) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	type PhoneNumberPayload struct {
		PhoneNumber string `json:"phoneNumber"`
	}
	resp, err := c.ClinicalISCClient.MakeRequest(ctx, http.MethodDelete, deleteFHIRPatientEndpoint, &PhoneNumberPayload{PhoneNumber: phoneNumber})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request: %w", err)
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("failed to delete fhir patient: %v", string(dataResponse))
	}
	return true, nil
}
