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
func (ps ServicePubSubMessaging) NotifyCreatePatient(ctx context.Context, client *dto.PatientCreationOutput) error {
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

// NotifyCreateCMSClient publishes to the create cms user topic and the user will be created in the CMS system
func (ps ServicePubSubMessaging) NotifyCreateCMSClient(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
	return ps.newPublish(ctx, user, common.CreateCMSClientTopicName, MyCareHubServiceName)
}

// NotifyDeleteCMSClient publishes to the delete cms user topic and the user will be deleted in the CMS.
func (ps ServicePubSubMessaging) NotifyDeleteCMSClient(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
	return ps.newPublish(ctx, user, common.DeleteCMSClientTopicName, MyCareHubServiceName)
}

// NotifyDeleteCMSStaff publishes to the delete cms staff topic and the staff will be deleted in the CMS.
func (ps ServicePubSubMessaging) NotifyDeleteCMSStaff(ctx context.Context, staff *dto.DeleteCMSUserPayload) error {
	return ps.newPublish(ctx, staff, common.DeleteCMSStaffTopicName, MyCareHubServiceName)
}

// NotifyCreateCMSStaff publishes to the create cms staff topic and the staff will be created in the CMS system
func (ps ServicePubSubMessaging) NotifyCreateCMSStaff(ctx context.Context, user *dto.PubsubCreateCMSStaffPayload) error {
	return ps.newPublish(ctx, user, common.CreateCMSStaffTopicName, MyCareHubServiceName)
}

// NotifyCreateCMSProgram publishes to the create cms program topic and the program will be created in the CMS.
func (ps ServicePubSubMessaging) NotifyCreateCMSProgram(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
	return ps.newPublish(ctx, program, common.CreateCMSProgramTopicName, MyCareHubServiceName)
}

// NotifyCreateCMSOrganisation publishes to the create cms organisation topic and the organisation will be created in the CMS.
func (ps ServicePubSubMessaging) NotifyCreateCMSOrganisation(ctx context.Context, organisation *dto.CreateCMSOrganisationPayload) error {
	return ps.newPublish(ctx, organisation, common.CreateCMSOrganisationTopicName, MyCareHubServiceName)
}

// NotifyCreateCMSFacility publishes to the create cms facility topic and the facility will be created in the CMS.
func (ps ServicePubSubMessaging) NotifyCreateCMSFacility(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
	return ps.newPublish(ctx, facility, common.CreateCMSFacilityTopicName, MyCareHubServiceName)
}

// NotifyCMSAddFacilityToProgram publishes to the add facility to program topic and the facility will be added to the program in the CMS.
func (ps ServicePubSubMessaging) NotifyCMSAddFacilityToProgram(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error {
	return ps.newPublish(ctx, payload, common.CreateCMSProgramFacilityTopicName, MyCareHubServiceName)
}
