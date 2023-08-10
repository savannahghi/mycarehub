package mail_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail"
	mailMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail/mock"
)

func TestMailgunServiceImpl_SendFeedback(t *testing.T) {
	type args struct {
		ctx             context.Context
		subject         string
		feedbackMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: send feedback",
			args: args{
				ctx:             context.Background(),
				subject:         "Feedback",
				feedbackMessage: "Hello World",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to send feedback",
			args: args{
				ctx:             context.Background(),
				subject:         "Feedback",
				feedbackMessage: "Hello World",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeMail := mailMock.NewMailGunClientMock()
			mg := mail.NewServiceMail(fakeMail)

			if tt.name == "Sad case: unable to send feedback" {
				fakeMail.MockSendFn = func(ctx context.Context, m *mailgun.Message) (string, string, error) {
					return "", "", errors.New("error")
				}
			}

			got, err := mg.SendFeedback(tt.args.ctx, tt.args.subject, tt.args.feedbackMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("MailgunServiceImpl.SendFeedback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("MailgunServiceImpl.SendFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}
