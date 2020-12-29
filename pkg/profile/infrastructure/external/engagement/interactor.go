package engagement

import (
	"log"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/profile/infrastructure/external"
)

// specific endpoint paths for ISC
const (
	// feed isc paths
	createFeedItem = "internal/createFeedItem"
	feedSrv        = "feed"
)

// Interactor defines the interactions between fulfillment and the
// engagement service.
type Interactor interface {
}

// NewEngagementInteractor initializes a notifications interactor
func NewEngagementInteractor() Interactor {
	client := external.NewInterServiceClient(feedSrv)
	interactor := &InteractorImpl{
		client: client,
	}
	interactor.checkPreconditions()
	return interactor
}

// InteractorImpl holds an initialized inter-service client
type InteractorImpl struct {
	client *base.InterServiceClient
}

func (i InteractorImpl) checkPreconditions() {
	if i.client == nil {
		log.Panicf("nil inter-service client in engagement interactor")
	}
}
