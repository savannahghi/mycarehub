package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	user "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func createPINInputData(phone string, PIN string, ConfirmedPIN string, flavour feedlib.Flavour) *dto.PinInput {
	return &dto.PinInput{
		PhoneNumber:  phone,
		PIN:          PIN,
		ConfirmedPin: ConfirmedPIN,
		Flavour:      flavour,
	}
}

func TestUseCasesUserImpl_SetUserPIN_Unittest(t *testing.T) {
	ctx := context.Background()

	phone := "+2547100000000"
	PIN := "1234"
	ConfirmedPin := "1234"
	flavour := feedlib.FlavourConsumer

	validPINInput := createPINInputData(phone, PIN, ConfirmedPin, flavour)

	noPhoneInput := createPINInputData("", PIN, ConfirmedPin, flavour)

	invalidPINLength := createPINInputData(phone, "123", "123", flavour)

	pinNotDigits := createPINInputData(phone, "page", ConfirmedPin, flavour)

	pinMismatch := createPINInputData(phone, "4321", ConfirmedPin, flavour)

	type args struct {
		ctx   context.Context
		input *dto.PinInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
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
			name: "Sad case - no phone",
			args: args{
				ctx:   ctx,
				input: noPhoneInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid PIN length",
			args: args{
				ctx:   ctx,
				input: invalidPINLength,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin not digits",
			args: args{
				ctx:   ctx,
				input: pinNotDigits,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx:   ctx,
				input: pinMismatch,
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to save PIN",
			args: args{
				ctx:   ctx,
				input: validPINInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeOnboardingLibraryExtension()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case - no phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to set user pin")
				}
			}
			if tt.name == "Sad case - invalid PIN length" {
				fakeUser.MockSetUserPINFn = func(ctx context.Context, input *dto.PinInput) (bool, error) {
					return false, fmt.Errorf("failed to set user pin")
				}
			}
			if tt.name == "Sad case - pin not digits" {
				fakeUser.MockSetUserPINFn = func(ctx context.Context, input *dto.PinInput) (bool, error) {
					return false, fmt.Errorf("failed to set user pin")
				}
			}

			if tt.name == "Sad case - pin mismatch" {
				fakeExtension.MockComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad case - fail to save PIN" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("failed to save PIN")
				}
			}

			_, err := u.SetUserPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
