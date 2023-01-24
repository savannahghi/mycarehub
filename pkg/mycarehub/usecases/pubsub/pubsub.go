package pubsub

import (
	"net/http"

	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
)

// IReceivePubsub represents the interface that receives pubsub messages
type IReceivePubsub interface {
	ReceivePubSubPushMessages(w http.ResponseWriter, r *http.Request)
}

// UseCasePubSub groups all the logic related to the PubSub service
type UseCasePubSub interface {
	IReceivePubsub
}

// ServicePubSubImpl represents PubSub implementation object
type ServicePubSubImpl struct {
	Pubsub pubsubmessaging.ServicePubsub
}

// NewUseCasePubSub is the controller for the PubSub usecases
func NewUseCasePubSub(
	pubsub pubsubmessaging.ServicePubsub,
) *ServicePubSubImpl {
	return &ServicePubSubImpl{
		Pubsub: pubsub,
	}
}

// ReceivePubSubPushMessages receives and processes a pubsub message
func (p *ServicePubSubImpl) ReceivePubSubPushMessages(w http.ResponseWriter, r *http.Request) {
	p.Pubsub.ReceivePubSubPushMessages(w, r)
}
