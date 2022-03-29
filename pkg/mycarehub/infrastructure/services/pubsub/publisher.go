package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (ps *ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic, serviceName string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data received: %v", err)
	}
	return ps.PublishToPubsub(
		ctx,
		ps.AddPubSubNamespace(topic, serviceName),
		serviceName,
		payload,
	)
}

// NotifyCreatePatient publishes to the create patient topic
func (ps ServicePubSubMessaging) NotifyCreatePatient(ctx context.Context, client *dto.ClientRegistrationOutput) error {
	return ps.newPublish(ctx, client, common.CreatePatientTopic, ClinicalServiceName)
}

// NotifyCreateVitals publishes to the create vitals topic
func (ps ServicePubSubMessaging) NotifyCreateVitals(ctx context.Context, vitals *dto.PatientVitalSignOutput) error {
	return ps.newPublish(ctx, vitals, common.CreateVitalsTopicName, ClinicalServiceName)
}

// NotifyCreateAllergy publishes to the create allergy topic
func (ps ServicePubSubMessaging) NotifyCreateAllergy(ctx context.Context, allergy *dto.PatientAllergyOutput) error {
	return ps.newPublish(ctx, allergy, common.CreateAllergyTopicName, ClinicalServiceName)
}

// NotifyCreateMedication publishes to the create medication topic
func (ps ServicePubSubMessaging) NotifyCreateMedication(ctx context.Context, medication *dto.PatientMedicationOutput) error {
	return ps.newPublish(ctx, medication, common.CreateMedicationTopicName, ClinicalServiceName)
}

// NotifyCreateTestOrder publishes to the create test order topic
func (ps ServicePubSubMessaging) NotifyCreateTestOrder(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error {
	return ps.newPublish(ctx, testOrder, common.CreateTestOrderTopicName, ClinicalServiceName)
}

// NotifyCreateTestResult publishes to the create test result topic
func (ps ServicePubSubMessaging) NotifyCreateTestResult(ctx context.Context, testResult *dto.PatientTestResultOutput) error {
	return ps.newPublish(ctx, testResult, common.CreateTestResultTopicName, ClinicalServiceName)
}

// NotifyCreateOrganization publishes to the create organisation topic the facilities without a FHIR organisation ID
func (ps ServicePubSubMessaging) NotifyCreateOrganization(ctx context.Context, facility *domain.Facility) error {
	return ps.newPublish(ctx, facility, common.CreateOrganizationTopicName, ClinicalServiceName)
}
