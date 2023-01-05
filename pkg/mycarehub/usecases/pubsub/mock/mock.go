package mock

import (
	"net/http"
)

// PubSubUseCaseMock mocks the PubSubs usecase methods
type PubSubUseCaseMock struct {
	MockReceivePubSubPushMessagesFn func(w http.ResponseWriter, r *http.Request)
}

// NewServicePubSubMock initializes a new PubSub mock instance
func NewServicePubSubMock() *PubSubUseCaseMock {
	return &PubSubUseCaseMock{
		MockReceivePubSubPushMessagesFn: func(w http.ResponseWriter, r *http.Request) {

		},
	}
}

// NotifyUser mocks the implementation of sending a fcm PubSub to a user
func (p *PubSubUseCaseMock) ReceivePubSubPushMessages(w http.ResponseWriter, r *http.Request) {
	p.MockReceivePubSubPushMessagesFn(w, r)
}
