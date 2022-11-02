package sms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms"
	mockSMS "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms/mock"
	"github.com/savannahghi/silcomms"
)

func TestSILCommsClient_SendSMS(t *testing.T) {
	type args struct {
		ctx        context.Context
		message    string
		recipients []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: send sms",
			args: args{
				ctx:        context.Background(),
				message:    "Hello WOrld",
				recipients: []string{interserviceclient.TestUserPhoneNumber},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to send sms",
			args: args{
				ctx:        context.Background(),
				message:    "Hello WOrld",
				recipients: []string{interserviceclient.TestUserPhoneNumber},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockSMS.NewSMSClientMock()
			sc := sms.NewServiceSMS(fakeClient)

			if tt.name == "Sad case: unable to send sms" {
				fakeClient.MockSendBulkSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := sc.SendSMS(tt.args.ctx, tt.args.message, tt.args.recipients)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsClient.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
