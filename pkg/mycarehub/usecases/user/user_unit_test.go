package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

func TestUseCasesUserImpl_Login_Unittest(t *testing.T) {
	ctx := context.Background()

	phoneNumber := "+2547100000000"
	PIN := "1234"
	flavour := enums.CONSUMER

	type args struct {
		ctx         context.Context
		phoneNumber string
		pin         string
		flavour     enums.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get user PIN By User ID",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create firebase token",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to authenticate token",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case - no phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - unable to get user PIN By User ID" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get user PIN by user ID")
				}
			}

			if tt.name == "Sad case - pin mismatch" {
				fakeExtension.MockComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad Case - Fail to create firebase token" {
				fakeExtension.MockCreateFirebaseCustomTokenFn = func(ctx context.Context, uid string) (string, error) {
					return "", fmt.Errorf("failed to create custom token")
				}
			}

			if tt.name == "Sad Case - Fail to authenticate token" {
				fakeExtension.MockAuthenticateCustomFirebaseTokenFn = func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
					return nil, fmt.Errorf("failed to authenticate token")
				}
			}

			_, _, err := u.Login(tt.args.ctx, tt.args.phoneNumber, tt.args.pin, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
