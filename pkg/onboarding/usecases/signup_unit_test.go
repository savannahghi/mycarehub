package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestSignUpUseCasesImpl_RetirePushToken(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_successfully_retire_pushtoken",
			args: args{
				ctx:   ctx,
				token: "VAL1IDT0K3N",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_fail_to_retire_pushtoken",
			args: args{
				ctx:   ctx,
				token: "VAL1IDT0K3N",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_fail_to_retire_pushtoken_invalid_length",
			args: args{
				ctx:   ctx,
				token: "*",
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_successfully_retire_pushtoken" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:         "f4f39af7--91bd-42b3af315a4e",
						PushTokens: []string{"token12", "token23", "token34"},
					}, nil
				}

				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_fail_to_retire_pushtoken" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:         "f4f39af7--91bd-42b3af315a4e",
						PushTokens: []string{"token12", "token23", "token34"},
					}, nil
				}

				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
					return fmt.Errorf("failed to retire push token")
				}
			}

			if tt.name == "invalid:_fail_to_retire_pushtoken_invalid_length" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:         "f4f39af7--91bd-42b3af315a4e",
						PushTokens: []string{"token12", "token23", "token34"},
					}, nil
				}

				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
					return fmt.Errorf("failed to retire push token")
				}
			}

			got, err := i.Signup.RetirePushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RetirePushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.RetirePushToken() = %v, want %v", got, tt.want)
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
