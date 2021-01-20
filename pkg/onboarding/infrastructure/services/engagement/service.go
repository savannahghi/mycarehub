package engagement

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

const engagementService = "engagement"

const (
	// engagement ISC paths
	publishNudge = "feed/%s/PRO/false/nudges/"
	publishItem  = "feed/%s/PRO/false/items/"
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error)
}

// ServiceEngagementImpl ... represents engagement usecases
type ServiceEngagementImpl struct {
	Engage *base.InterServiceClient
}

// NewServiceEngagementImpl ...
func NewServiceEngagementImpl() ServiceEngagement {

	client := utils.NewInterServiceClient(engagementService)

	return &ServiceEngagementImpl{Engage: client}
}

// PublishKYCNudge ...
func (en *ServiceEngagementImpl) PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishNudge, uid), payload)
}

// PublishKYCFeedItem ...
func (en *ServiceEngagementImpl) PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishItem, uid), payload)
}
