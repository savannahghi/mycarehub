package engagement

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

const (
	// Feed ISC paths
	publishNudge         = "feed/%s/PRO/false/nudges/"
	publishItem          = "feed/%s/PRO/false/items/"
	resolveDefaultNudges = "feed/%s/%s/false/defaultnudges/%s/resolve/"

	// Communication ISC paths
	sendEmail = "internal/send_email"
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	PublishKYCNudge(
		uid string,
		payload base.Nudge,
	) (*http.Response, error)

	PublishKYCFeedItem(
		uid string,
		payload base.Item,
	) (*http.Response, error)

	ResolveDefaultNudgeByTitle(
		UID string,
		flavour base.Flavour,
		nudgeTitle string,
	) error

	SendMail(
		email string,
		message string,
		subject string,
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

// SendMail sends emails to communicate to our users
func (en *ServiceEngagementImpl) SendMail(
	email string,
	message string,
	subject string,
) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email address: %v", email)
	}

	body := map[string]interface{}{
		"to":      []string{email},
		"text":    message,
		"subject": subject,
	}

	resp, err := en.Engage.MakeRequest(
		http.MethodPost,
		sendEmail,
		body,
	)
	if err != nil {
		return fmt.Errorf("unable to send email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send email : %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	return nil
}
