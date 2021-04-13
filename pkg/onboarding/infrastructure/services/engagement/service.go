package engagement

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
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
	SendAlertToSupplier(
		supplierName string,
		partnerType string,
		accountType string,
		subjectTitle string,
		emailBody string,
		emailAddress string,
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

//SendAlertToSupplier send email to admin notifying them of them of new
// KYC Request.
func (en *ServiceEngagementImpl) SendAlertToSupplier(
	supplierName string,
	partnerType string,
	accountType string,
	subjectTitle string,
	emailBody string,
	emailAddress string,
) error {
	var writer bytes.Buffer
	t := template.Must(template.New("profile").Parse(utils.AcknowledgementKYCEmail))
	_ = t.Execute(&writer, struct {
		SupplierName string
		PartnerType  string
		AccountType  string
		EmailBody    string
		EmailAddress string
	}{
		SupplierName: supplierName,
		PartnerType:  partnerType,
		AccountType:  accountType,
		EmailBody:    emailBody,
		EmailAddress: emailAddress,
	})

	body := map[string]interface{}{
		"to":      []string{emailAddress},
		"text":    writer.String(),
		"subject": subjectTitle,
	}

	resp, err := en.Engage.MakeRequest(http.MethodPost, sendEmail, body)

	if err != nil {
		return fmt.Errorf("unable to send Alert to admin email: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send Alert to admin email: %w", err)
	}

	return nil
}
