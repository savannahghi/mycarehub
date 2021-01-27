package mailgun

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

const mailgunService = "mailgun"

const (
	sendEmail = "internal/send_email"
)

// ServiceMailgun represents mailgun usecases
type ServiceMailgun interface {
	SendMail(email string, message string, subject string) error
	FetchClient() *base.InterServiceClient
}

// ServiceMailgunImpl ... represents mailgun usecases
type ServiceMailgunImpl struct {
	Mailgun *base.InterServiceClient
}

// NewServiceMailgunImpl ...
func NewServiceMailgunImpl(mailgun extension.ISCClientExtension) ServiceMailgun {
	client := utils.NewInterServiceClient(mailgunService)
	return &ServiceMailgunImpl{Mailgun: client}
}

// FetchClient ...
func (mg *ServiceMailgunImpl) FetchClient() *base.InterServiceClient {
	return mg.Mailgun
}

// SendMail ...
func (mg *ServiceMailgunImpl) SendMail(email string, message string, subject string) error {

	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email address: %v", email)
	}

	body := map[string]interface{}{
		"to":      []string{email},
		"text":    message,
		"subject": subject,
	}

	resp, err := mg.FetchClient().MakeRequest(http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send KYC email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send KYC email : %w, with status code %v", err, resp.StatusCode)
	}

	return nil
}
