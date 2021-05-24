package usecases_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestUSSDUseCaseImpl_CreateUSSDData(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	phone := "+254721026491"
	invalidPhone := ""
	validInput := dto.EndSessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Input:       "1",
	}

	invalidInput := dto.EndSessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &invalidPhone,
		Input:       "1",
	}
	aplhaNumbericphone := "+254-not-valid-123"
	alphanumericPhoneInput := dto.EndSessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &aplhaNumbericphone,
		Input:       "1",
	}

	type args struct {
		ctx   context.Context
		input dto.EndSessionDetails
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully add USSD Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "sad:( fail to add USSD Details",
			args: args{
				ctx:   ctx,
				input: invalidInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( invalid phone number",
			args: args{
				ctx:   ctx,
				input: alphanumericPhoneInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( invalid phone number",
			args: args{
				ctx:   ctx,
				input: invalidInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) successfully add USSD Details" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}
				fakeRepo.AddIncomingUSSDDataFn = func(ctx context.Context, input *dto.EndSessionDetails) error {
					return nil
				}

			}
			if tt.name == "sad:( fail to add USSD Details" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}
				fakeRepo.AddIncomingUSSDDataFn = func(ctx context.Context, input *dto.EndSessionDetails) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad:( invalid phone number" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			err := i.AITUSSD.CreateUSSDData(tt.args.ctx, &tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("USSDUseCaseImpl.AddUSSDDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
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

func TestUSSDUseCaseImpl_GenerateUSSD(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	phone := "+254721026491"
	invalidPhone := ""
	validInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "",
	}

	validInputWithEmptyText := dto.SessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &invalidPhone,
		Text:        "1",
	}
	invalidInput := dto.SessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &invalidPhone,
		Text:        "6",
	}

	type args struct {
		ctx   context.Context
		input dto.SessionDetails
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy:) text input is empty ",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want: "CON",
		},
		{
			name: "happy :) text with input",
			args: args{
				ctx:   ctx,
				input: validInputWithEmptyText,
			},
			want: "END",
		},
		{
			name: "happy :) text with invalid input",
			args: args{
				ctx:   ctx,
				input: invalidInput,
			},
			want: "CON Invalid choice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := i.AITUSSD.GenerateUSSD(&tt.args.input)
			if strings.Contains(resp, tt.want) != true {
				t.Errorf("expected %v to be in  %v  ", tt.want, resp)
				return
			}
		})
	}
}
