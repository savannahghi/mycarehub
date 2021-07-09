package ussd_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestImpl_HandleResponseFromUSSDGateway_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	sessionID := uuid.New().String()
	unregisteredPhoneNumber := "0723456756"
	registeredPhoneNumber := "+254700100200"

	unregisteredValidPayload := &dto.SessionDetails{
		SessionID:   sessionID,
		PhoneNumber: &unregisteredPhoneNumber,
	}

	registeredValidPayload := &dto.SessionDetails{
		SessionID:   sessionID,
		PhoneNumber: &registeredPhoneNumber,
	}

	invalidPayload := &dto.SessionDetails{
		SessionID:   "",
		PhoneNumber: &registeredPhoneNumber,
	}

	type args struct {
		ctx     context.Context
		payload *dto.SessionDetails
	}
	tests := []struct {
		name     string
		args     args
		response string
	}{
		{
			name: "Happy case ):_Success case_Unregistered_user",
			args: args{
				ctx:     ctx,
				payload: unregisteredValidPayload,
			},
			response: "CON Welcome to Be.Well\r\n" +
				"1. Register\r\n" +
				"2. Opt Out\r\n",
		},
		{
			name: "Happy case ):_Success case_Registered_user",
			args: args{
				ctx:     ctx,
				payload: registeredValidPayload,
			},
			response: "CON Welcome to Be.Well\r\n" +
				"1. Register\r\n" +
				"2. Opt Out\r\n",
		},
		{
			name: "SAD case ):Fail case_invalid_sessionID",
			args: args{
				ctx:     ctx,
				payload: invalidPayload,
			},
			response: "END Something went wrong. Please try again.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case ):_Success case_Unregistered_user" {
				fakeRepo.HandleResponseFromUSSDGatewayFn = func(context context.Context, input *dto.SessionDetails) string {
					return "CON Welcome to Be.Well\r\n" +
						"1. Register\r\n" +
						"2. Opt Out\r\n"
				}
			}

			if tt.name == "Happy case ):_Success case_Registered_user" {
				fakeRepo.HandleResponseFromUSSDGatewayFn = func(context context.Context, input *dto.SessionDetails) string {
					return "CON Welcome to Be.Well\r\n" +
						"1. Opt out from marketing messages\r\n" +
						"2. Change PIN"
				}
			}

			if tt.name == "SAD case ):Fail case_invalid_sessionID" {
				fakeRepo.HandleResponseFromUSSDGatewayFn = func(context context.Context, input *dto.SessionDetails) string {
					return "END Something went wrong. Please try again."
				}
			}

			if gotresp := u.AITUSSD.HandleResponseFromUSSDGateway(tt.args.ctx, tt.args.payload); gotresp != tt.response {
				t.Errorf("Impl.HandleResponseFromUSSDGateway() = %v, want %v", gotresp, tt.response)
			}
		})
	}
}
