package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

func TestUseCasesUserImpl_SetUserPIN_Unittest(t *testing.T) {
	ctx := context.Background()

	f := testFakeInfrastructureInteractor

	validPINInput := &dto.PINInput{
		PIN:          "1234",
		ConfirmedPin: "1234",
		Flavour:      feedlib.FlavourConsumer,
	}

	invalidPINInput := &dto.PINInput{
		PIN:          "",
		ConfirmedPin: "1234",
		Flavour:      "CONSUMER",
	}

	type args struct {
		ctx   context.Context
		input *dto.PINInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: validPINInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: invalidPINInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeCreate.SetUserPINFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Happy case" {
				fakePIN.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "encryptedPIN"
				}
				fakeCreate.SetUserPINFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return true, nil
				}
			}
			_, err := f.SetUserPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
