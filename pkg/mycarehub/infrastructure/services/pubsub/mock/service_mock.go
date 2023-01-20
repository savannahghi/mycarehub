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

	MockNotifyCreateVitalsFn       func(ctx context.Context, vitals *dto.PatientVitalSignOutput) error
	MockNotifyCreateAllergyFn      func(ctx context.Context, allergy *dto.PatientAllergyOutput) error
	MockNotifyCreateMedicationFn   func(ctx context.Context, medication *dto.PatientMedicationOutput) error
	MockNotifyCreateTestOrderFn    func(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error
	MockNotifyCreateTestResultFn   func(ctx context.Context, testResult *dto.PatientTestResultOutput) error
	MockNotifyCreateOrganizationFn func(ctx context.Context, facility *domain.Facility) error
	MockNotifyGetStreamEventFn     func(ctx context.Context, event *dto.GetStreamEvent) error
	MockNotifyCreateCMSUserFn      func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error
	MockNotifyDeleteCMSClientFn    func(ctx context.Context, user *dto.DeleteCMSUserPayload) error
	MockNotifyDeleteCMSStaffFn     func(ctx context.Context, user *dto.DeleteCMSUserPayload) error
	MockNotifyCreateCMSStaffFn     func(ctx context.Context, user *dto.PubsubCreateCMSStaffPayload) error
	MockNotifyCreateCMSProgramFn   func(ctx context.Context, program *dto.CreateCMSProgramPayload) error
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
		MockNotifyGetStreamEventFn: func(ctx context.Context, event *dto.GetStreamEvent) error {
			return nil
		},
		MockNotifyCreateCMSUserFn: func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
			return nil
		},
		MockNotifyDeleteCMSClientFn: func(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
			return nil
		},
		MockNotifyDeleteCMSStaffFn: func(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
			return nil
		},
		MockNotifyCreateCMSStaffFn: func(ctx context.Context, user *dto.PubsubCreateCMSStaffPayload) error {
			return nil
		},
		MockNotifyCreateCMSProgramFn: func(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
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

// NotifyGetStreamEvent mocks the implementation of publishing getstream events to a pubsub topic
func (m *FakeServicePubSub) NotifyGetStreamEvent(ctx context.Context, event *dto.GetStreamEvent) error {
	return m.MockNotifyGetStreamEventFn(ctx, event)
}

// NotifyCreateCMSClient mocks the implementation of publishing create cms user events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSClient(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
	return m.MockNotifyCreateCMSUserFn(ctx, user)
}

// NotifyDeleteCMSClient mocks the implementation of publishing delete cms user events to a pubsub topic
func (m *FakeServicePubSub) NotifyDeleteCMSClient(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
	return m.MockNotifyDeleteCMSClientFn(ctx, user)
}

// NotifyDeleteCMSStaff mocks the implementation of publishing delete cms staff events to a pubsub topic
func (m *FakeServicePubSub) NotifyDeleteCMSStaff(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
	return m.MockNotifyDeleteCMSStaffFn(ctx, user)
}

// NotifyCreateCMSStaff mocks the implementation of publishing create cms staff events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSStaff(ctx context.Context, user *dto.PubsubCreateCMSStaffPayload) error {
	return m.MockNotifyCreateCMSStaffFn(ctx, user)
}

// NotifyCreateCMSProgram mocks the implementation of publishing create cms program events to a pubsub topic
func (m *FakeServicePubSub) NotifyCreateCMSProgram(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
	return m.MockNotifyCreateCMSProgramFn(ctx, program)
}
