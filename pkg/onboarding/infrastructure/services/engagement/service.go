package engagement

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

const (
	// engagement ISC paths
	publishNudge         = "feed/%s/PRO/false/nudges/"
	publishItem          = "feed/%s/PRO/false/items/"
	resolveDefaultNudges = "feed/%s/%s/false/defaultnudges/%s/resolve/"
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error)
	ResolveDefaultNudgeByTitle(
		UID string,
		flavour base.Flavour,
		nudgeTitle string,
	) error
}

// ServiceEngagementImpl represents engagement usecases
type ServiceEngagementImpl struct {
	Engage extension.ISCClientExtension
}

// NewServiceEngagementImpl returns new instance of ServiceEngagementImpl
func NewServiceEngagementImpl(eng extension.ISCClientExtension) ServiceEngagement {
	return &ServiceEngagementImpl{Engage: eng}
}

// PublishKYCNudge calls the `engagement service` to publish
// a KYC nudge
func (en *ServiceEngagementImpl) PublishKYCNudge(
	uid string,
	payload base.Nudge,
) (*http.Response, error) {
	return en.Engage.MakeRequest(
		http.MethodPost,
		fmt.Sprintf(publishNudge, uid),
		payload,
	)
}

// PublishKYCFeedItem calls the `engagement service` to publish
// a KYC feed item
func (en *ServiceEngagementImpl) PublishKYCFeedItem(
	uid string,
	payload base.Item,
) (*http.Response, error) {
	return en.Engage.MakeRequest(
		http.MethodPost,
		fmt.Sprintf(publishItem, uid),
		payload,
	)
}

// ResolveDefaultNudgeByTitle calls the `engagement service`
// to resolve any default nudge by its `Title`
func (en *ServiceEngagementImpl) ResolveDefaultNudgeByTitle(
	UID string,
	flavour base.Flavour,
	nudgeTitle string,
) error {
	resp, err := en.Engage.MakeRequest(
		http.MethodPatch,
		fmt.Sprintf(
			resolveDefaultNudges,
			UID,
			flavour,
			nudgeTitle,
		),
		nil,
	)

	if err != nil {
		return exceptions.ResolveNudgeErr(
			err,
			flavour,
			nudgeTitle,
			nil,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return exceptions.ResolveNudgeErr(
			fmt.Errorf("unexpected status code %v", resp.StatusCode),
			flavour,
			nudgeTitle,
			&resp.StatusCode,
		)
	}

	return nil
}
