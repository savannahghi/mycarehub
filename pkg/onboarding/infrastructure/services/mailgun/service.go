package mailgun

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

const (
	sendEmail = "internal/send_email"
)

// ServiceMailgun represents mailgun usecases
type ServiceMailgun interface {
	SendMail(email string, message string, subject string) error
}

// ServiceMailgunImpl ... represents mailgun usecases
type ServiceMailgunImpl struct {
	Mailgun extension.ISCClientExtension
}

// NewServiceMailgunImpl ...
func NewServiceMailgunImpl(mailgun extension.ISCClientExtension) ServiceMailgun {
	return &ServiceMailgunImpl{Mailgun: mailgun}
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

	resp, err := mg.Mailgun.MakeRequest(http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send KYC email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send KYC email : %w, with status code %v", err, resp.StatusCode)
	}

	return nil
}
