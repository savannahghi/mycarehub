package mock

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// FakeServicePubSub ...
type FakeServicePubSub struct {
	MockPublishToPubsubFn func(
		ctx context.Context,
		topicID string,
		serviceName string,
		payload []byte,
	) error

	MockReceivePubSubPushMessagesFn func(
		w http.ResponseWriter,
		r *http.Request,
	)

	MockNotifyCreatePatientFn func(ctx context.Context, client *dto.PatientCreationOutput) error

	MockNotifyCreateVitalsFn            func(ctx context.Context, vitals *dto.PatientVitalSignOutput) error
	MockNotifyCreateAllergyFn           func(ctx context.Context, allergy *dto.PatientAllergyOutput) error
	MockNotifyCreateMedicationFn        func(ctx context.Context, medication *dto.PatientMedicationOutput) error
	MockNotifyCreateTestOrderFn         func(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error
	MockNotifyCreateTestResultFn        func(ctx context.Context, testResult *dto.PatientTestResultOutput) error
	MockNotifyCreateOrganizationFn      func(ctx context.Context, facility *domain.Facility) error
	MockNotifyCreateCMSUserFn           func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error
	MockNotifyDeleteCMSClientFn         func(ctx context.Context, user *dto.DeleteCMSUserPayload) error
	MockNotifyCreateCMSProgramFn        func(ctx context.Context, program *dto.CreateCMSProgramPayload) error
	MockNotifyCreateCMSOrganisationFn   func(ctx context.Context, program *dto.CreateCMSOrganisationPayload) error
	MockNotifyCreateCMSFacilityFn       func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error
	MockNotifyCMSAddFacilityToProgramFn func(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error
	MockNotifyCreateClinicalTenantFn    func(ctx context.Context, tenant *dto.ClinicalTenantPayload) error
	MockNotifyRegisterMatrixUserFn      func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error
	MockNotifyCreateCMSClientFn         func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error
}

// NewPubsubServiceMock mocks the pubsub service implementation
func NewPubsubServiceMock() *FakeServicePubSub {
	return &FakeServicePubSub{
		MockPublishToPubsubFn: func(ctx context.Context, topicID string, serviceName string, payload []byte) error {
			return nil
		},
		MockReceivePubSubPushMessagesFn: func(w http.ResponseWriter, r *http.Request) {
			resp := map[string]string{"Status": "Success"}
			returnedResponse, _ := json.Marshal(resp)
			_, _ = w.Write(returnedResponse)
		},
		MockNotifyCreatePatientFn: func(ctx context.Context, client *dto.PatientCreationOutput) error {
			return nil
		},
		MockNotifyCreateVitalsFn: func(ctx context.Context, vitals *dto.PatientVitalSignOutput) error {
			return nil
		},
		MockNotifyCreateAllergyFn: func(ctx context.Context, allergy *dto.PatientAllergyOutput) error {
			return nil
		},
		MockNotifyCreateMedicationFn: func(ctx context.Context, medication *dto.PatientMedicationOutput) error {
			return nil
		},
		MockNotifyCreateTestOrderFn: func(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error {
			return nil
		},
		MockNotifyCreateTestResultFn: func(ctx context.Context, testResult *dto.PatientTestResultOutput) error {
			return nil
		},
		MockNotifyCreateOrganizationFn: func(ctx context.Context, facility *domain.Facility) error {
			return nil
		},
		MockNotifyCreateCMSUserFn: func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
			return nil
		},
		MockNotifyDeleteCMSClientFn: func(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
			return nil
		},

		MockNotifyCreateCMSProgramFn: func(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
			return nil
		},
		MockNotifyCreateCMSOrganisationFn: func(ctx context.Context, program *dto.CreateCMSOrganisationPayload) error {
			return nil
		},
		MockNotifyCreateCMSFacilityFn: func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
			return nil
		},
		MockNotifyCMSAddFacilityToProgramFn: func(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error {
			return nil
		},
		MockNotifyCreateClinicalTenantFn: func(ctx context.Context, tenant *dto.ClinicalTenantPayload) error {
			return nil
		},
		MockNotifyRegisterMatrixUserFn: func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
			return nil
		},
		MockNotifyCreateCMSClientFn: func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
			return nil
		},
	}
}

// PublishToPubsub publishes a message to a specified topic
func (m *FakeServicePubSub) PublishToPubsub(
	ctx context.Context,
	topicID string,
	serviceName string,
	payload []byte,
) error {
	return m.MockPublishToPubsubFn(ctx, topicID, serviceName, payload)
}

// NotifyCreatePatient publishes to the create patient topic
func (m *FakeServicePubSub) NotifyCreatePatient(ctx context.Context, client *dto.PatientCreationOutput) error {
	return m.MockNotifyCreatePatientFn(ctx, client)
}

// ReceivePubSubPushMessages receives and processes a pubsub message
func (m *FakeServicePubSub) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	m.MockReceivePubSubPushMessagesFn(w, r)
}

// NotifyCreateVitals publishes to the create vitals topic
func (m *FakeServicePubSub) NotifyCreateVitals(ctx context.Context, vitals *dto.PatientVitalSignOutput) error {
	return m.MockNotifyCreateVitalsFn(ctx, vitals)
}

// NotifyCreateAllergy publishes to the create allergy topic
func (m *FakeServicePubSub) NotifyCreateAllergy(ctx context.Context, allergy *dto.PatientAllergyOutput) error {
	return m.MockNotifyCreateAllergyFn(ctx, allergy)
}

// NotifyCreateMedication publishes to the create medication topic
func (m *FakeServicePubSub) NotifyCreateMedication(ctx context.Context, medication *dto.PatientMedicationOutput) error {
	return m.MockNotifyCreateMedicationFn(ctx, medication)
}

// NotifyCreateTestOrder publishes to the create test order topic
func (m *FakeServicePubSub) NotifyCreateTestOrder(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error {
	return m.MockNotifyCreateTestOrderFn(ctx, testOrder)
}

// NotifyCreateTestResult publishes to the create test result topic
func (m *FakeServicePubSub) NotifyCreateTestResult(ctx context.Context, testResult *dto.PatientTestResultOutput) error {
	return m.MockNotifyCreateTestResultFn(ctx, testResult)
}

// NotifyCreateOrganization publishes to the create organization create topic
func (m *FakeServicePubSub) NotifyCreateOrganization(ctx context.Context, facility *domain.Facility) error {
	return m.MockNotifyCreateOrganizationFn(ctx, facility)
}

// NotifyCreateCMSClient mocks the implementation of publishing create cms user events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSClient(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
	return m.MockNotifyCreateCMSClientFn(ctx, user)
}

// NotifyDeleteCMSClient mocks the implementation of publishing delete cms user events to a pubsub topic
func (m *FakeServicePubSub) NotifyDeleteCMSClient(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
	return m.MockNotifyDeleteCMSClientFn(ctx, user)
}

// NotifyCreateCMSProgram mocks the implementation of publishing create cms program events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSProgram(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
	return m.MockNotifyCreateCMSProgramFn(ctx, program)
}

// NotifyCreateCMSOrganisation mocks the implementation of publishing create cms organisation events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSOrganisation(ctx context.Context, organisation *dto.CreateCMSOrganisationPayload) error {
	return m.MockNotifyCreateCMSOrganisationFn(ctx, organisation)
}

// NotifyCreateCMSFacility mocks the implementation of publishing create cms facility events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSFacility(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
	return m.MockNotifyCreateCMSFacilityFn(ctx, facility)
}

// NotifyCMSAddFacilityToProgram mocks the implementation of publishing add facility to program events to a pubsub topic
func (m *FakeServicePubSub) NotifyCMSAddFacilityToProgram(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error {
	return m.MockNotifyCMSAddFacilityToProgramFn(ctx, payload)
}

// NotifyCreateClinicalTenant mocks the implementation of creating a clinical service tenant
func (m *FakeServicePubSub) NotifyCreateClinicalTenant(ctx context.Context, tenant *dto.ClinicalTenantPayload) error {
	return m.MockNotifyCreateClinicalTenantFn(ctx, tenant)
}

// NotifyRegisterMatrixUser mocks the implementation of registering a matrix user
func (m *FakeServicePubSub) NotifyRegisterMatrixUser(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
	return m.MockNotifyRegisterMatrixUserFn(ctx, payload)
}
