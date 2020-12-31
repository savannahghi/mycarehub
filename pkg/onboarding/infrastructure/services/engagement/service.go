package engagement

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const engagementService = "engagement"

const (
	// engagement ISC paths
	publishNudge = "feed/%s/PRO/false/nudges/"
	publishItem  = "feed/%s/PRO/false/items/"
)

// Service  ...
type Service interface {
	PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error)
}

// ServiceImpl represents engagement usecases
type ServiceImpl struct {
	Engage *base.InterServiceClient
}

// NewEngagementService ...
func NewEngagementService(r repository.OnboardingRepository) Service {

	client := services.NewInterServiceClient(engagementService)
	return &ServiceImpl{Engage: client}
}

// PublishKYCNudge ...
func (en *ServiceImpl) PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishNudge, uid), payload)
}

// PublishKYCFeedItem ...
func (en *ServiceImpl) PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishItem, uid), payload)
}
