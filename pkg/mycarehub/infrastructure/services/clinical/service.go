package clinical

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
)

const (
	deletePatientPath = "internal/delete-patient"
)

// IServiceClinical ...
type IServiceClinical interface {
	DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) error
}

// ServiceClinical represents clinical isc actions
type ServiceClinical struct {
	Client extension.ISCClientExtension
}

// NewServiceClinical returns new instance of ServiceClinical
func NewServiceClinical(
	client extension.ISCClientExtension,
) *ServiceClinical {
	return &ServiceClinical{
		Client: client,
	}
}

// DeleteFHIRPatientByPhone deletes a patient by phone number
func (c *ServiceClinical) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) error {
	payload := struct {
		PhoneNumber string `json:"phoneNumber"`
	}{
		PhoneNumber: phoneNumber,
	}

	resp, err := c.Client.MakeRequest(ctx, http.MethodDelete, deletePatientPath, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to make request: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete fhir patient, status:%v body: %v", resp.Status, string(data))
	}

	return nil
}
