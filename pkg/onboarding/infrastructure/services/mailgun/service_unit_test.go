package mailgun_test

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
)

var fakeISCExt extMock.ISCClientExtension
var mailgunClient extension.ISCClientExtension = &fakeISCExt

func TestServiceMailgunImpl_SendMail(t *testing.T) {
	mailgun := mailgun.NewServiceMailgunImpl(mailgunClient)

	type args struct {
		email   string
		message string
		subject string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid:successfully_send_email",
			args: args{
				email:   "johndoe@gmail.com",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid:use_an_invalid_email",
			args: args{
				email:   "12345",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid:error_while_sending_request",
			args: args{
				email:   "johndoe",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_send_email" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       nil,
					}, nil
				}
			}

			if tt.name == "invalid:use_an_invalid_email" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: 400,
						Body:       nil,
					}, fmt.Errorf("an error occured! Invalid email address")
				}
			}

			if tt.name == "invalid:error_while_sending_request" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occured!")
				}
			}
			err := mailgun.SendMail(tt.args.email, tt.args.message, tt.args.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceMailgunImpl.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}
