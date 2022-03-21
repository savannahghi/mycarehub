package mock

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
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

	MockNotifyCreatePatientFn func(ctx context.Context, client *dto.ClientRegistrationOutput) error
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
		MockNotifyCreatePatientFn: func(ctx context.Context, client *dto.ClientRegistrationOutput) error {
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
func (m *FakeServicePubSub) NotifyCreatePatient(ctx context.Context, client *dto.ClientRegistrationOutput) error {
	return m.MockNotifyCreatePatientFn(ctx, client)
}

// ReceivePubSubPushMessages receives and processes a pubsub message
func (m *FakeServicePubSub) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	m.MockReceivePubSubPushMessagesFn(w, r)
}
