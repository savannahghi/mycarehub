package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/mock"
)

func TestUseCasesClientImpl_RegisterClient(t *testing.T) {
	ctx := context.Background()
	d := testFakeInfrastructureInteractor

	userPayload := &dto.UserInput{
		FirstName:   "FirstName",
		LastName:    "Last Name",
		UserName:    "User Name",
		MiddleName:  "Middle Name",
		DisplayName: "Display Name",
		Gender:      enumutils.GenderMale,
	}

	clientPayload := &dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}
	type args struct {
		ctx         context.Context
		userInput   *dto.UserInput
		clientInput *dto.ClientProfileInput
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// {
		// 	name: "Happy Case",
		// 	args: args{
		// 		ctx:         ctx,
		// 		userInput:   userPayload,
		// 		clientInput: clientPayload,
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "Sad Case: Fail to register user",
			args: args{
				ctx:         ctx,
				userInput:   userPayload,
				clientInput: clientPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewCreateMock()

			if tt.name == "Sad Case: Fail to register user" {
				fakeCreate.RegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register a client")
				}
			}

			got, err := d.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", nil)
			}
		})
	}
}
