package twilio_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/kevinburke/twilio-go"
	serviceTwilio "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio"
	twilioClientMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio/mock"
)

func TestTwilioServiceImpl_SendSMSViaTwilio(t *testing.T) {
	type args struct {
		ctx         context.Context
		phonenumber string
		message     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: send sms via twilio",
			args: args{
				ctx:         context.Background(),
				phonenumber: "+254700000000",
				message:     "Hello World",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to send sms via twilio",
			args: args{
				ctx:         context.Background(),
				phonenumber: "+254700000000",
				message:     "Hello World",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeTwilio := twilioClientMock.NewTwilioClientMock()
			tw := serviceTwilio.NewServiceTwilio(fakeTwilio)

			if tt.name == "Sad case: unable to send sms via twilio" {
				fakeTwilio.MockSendMessageFn = func(from, to, body string, mediaURLs []*url.URL) (*twilio.Message, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if err := tw.SendSMSViaTwilio(tt.args.ctx, tt.args.phonenumber, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.SendSMSViaTwilio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
