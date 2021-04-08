package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestProfileUseCaseImpl_ResumeWIthPin(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx context.Context
		pin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "valid:_login_with_pin",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid:_unable_to_get_profile",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_userprofile_returns_nil",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_unable_to_get_pin_by_profile_id",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_pin_data_returns_nil",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_pin_mismatch",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			// if the pins don't match, return false and dont throw an error.
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_login_with_pin" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}
			}

			if tt.name == "invalid:_unable_to_get_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to log in")
				}
			}

			if tt.name == "invalid:_userprofile_returns_nil" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, nil
				}

			}

			if tt.name == "invalid:_unable_to_get_pin_by_profile_id" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get pin by profile id")
				}
			}

			if tt.name == "invalid:_pin_data_returns_nil" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, nil
				}
			}

			if tt.name == "invalid:_pin_mismatch" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			isLogin, err := i.Login.ResumeWithPin(
				tt.args.ctx,
				tt.args.pin,
			)

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

				if tt.want != isLogin {
					t.Errorf("expected %v got %v", tt.want, isLogin)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_login_by_phone",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_normalize_phone",
			args: args{
				ctx:     ctx,
				phone:   "+21",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getUserProfile",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getPin",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_generateAuthCredentials",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getCustomerOrSupplierProfile",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_comparePin",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: base.FlavourConsumer,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_login_by_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}

				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *base.UserProfile) (*base.AuthCredentialResponse, error) {
					customToken := uuid.New().String()
					idToken := uuid.New().String()
					refreshToken := uuid.New().String()
					return &base.AuthCredentialResponse{
						CustomToken:  &customToken,
						IDToken:      &idToken,
						RefreshToken: refreshToken,
					}, nil
				}

				fakeRepo.GetCustomerOrSupplierProfileByProfileIDFn = func(ctx context.Context, flavour base.Flavour, profileID string) (*base.Customer, *base.Supplier, error) {
					return &base.Customer{ID: "5550"}, &base.Supplier{ID: "5550"}, nil
				}

				fakeRepo.GetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{ID: "111", ProfileID: "profile-id", AllowWhatsApp: true, AllowEmail: true, AllowTextSMS: true, AllowPush: true}, nil
				}
			}

			if tt.name == "invalid:fail_to_normalize_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("failed to normalize phone")
				}
			}

			if tt.name == "invalid:fail_to_getUserProfile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to user profile by phone number")
				}
			}

			if tt.name == "invalid:fail_to_getPin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("failed to get pin")
				}
			}

			if tt.name == "invalid:fail_to_generateAuthCredentials" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}

				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *base.UserProfile) (*base.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("failed to generate auth credentials")
				}
			}

			if tt.name == "invalid:fail_to_getCustomerOrSupplierProfile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}

				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *base.UserProfile) (*base.AuthCredentialResponse, error) {
					customToken := uuid.New().String()
					idToken := uuid.New().String()
					refreshToken := uuid.New().String()
					return &base.AuthCredentialResponse{
						CustomToken:  &customToken,
						IDToken:      &idToken,
						RefreshToken: refreshToken,
					}, nil
				}

				fakeRepo.GetCustomerOrSupplierProfileByProfileIDFn = func(ctx context.Context, flavour base.Flavour, profileID string) (*base.Customer, *base.Supplier, error) {
					return nil, nil, fmt.Errorf("failed to get customer or supplier profile by ID")
				}
			}

			if tt.name == "invalid:fail_to_comparePin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			got, err := i.Login.LoginByPhone(tt.args.ctx, tt.args.phone, tt.args.PIN, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.LoginByPhone() error = %v, wantErr %v", err, tt.wantErr)
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

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}

}

func TestProfileUseCaseImpl_RefreshToken(t *testing.T) {
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_refreshToken",
			args: args{
				token: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "invalid:invalid_refreshtoken",
			args: args{
				token: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_refreshToken" {
				fakeRepo.ExchangeRefreshTokenForIDTokenFn = func(token string) (*base.AuthCredentialResponse, error) {
					customToken := uuid.New().String()
					idToken := uuid.New().String()
					refreshToken := uuid.New().String()
					return &base.AuthCredentialResponse{
						CustomToken:  &customToken,
						IDToken:      &idToken,
						RefreshToken: refreshToken,
					}, nil
				}
			}

			if tt.name == "invalid:invalid_refreshtoken" {
				fakeRepo.ExchangeRefreshTokenForIDTokenFn = func(token string) (*base.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("invalid refresh token")
				}
			}
			got, err := i.Login.RefreshToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
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

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_LoginAsAnonymous(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_LoginAsAnonymous",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_generateAuthCredentials",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_LoginAsAnonymous" {
				fakeRepo.GenerateAuthCredentialsForAnonymousUserFn = func(ctx context.Context) (*base.AuthCredentialResponse, error) {
					customToken := uuid.New().String()
					idToken := uuid.New().String()
					refreshToken := uuid.New().String()
					return &base.AuthCredentialResponse{
						CustomToken:  &customToken,
						IDToken:      &idToken,
						RefreshToken: refreshToken,
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_generateAuthCredentials" {
				fakeRepo.GenerateAuthCredentialsForAnonymousUserFn = func(ctx context.Context) (*base.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("failed to generate auth credentials")
				}
			}

			got, err := i.Login.LoginAsAnonymous(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.LoginAsAnonymous() error = %v, wantErr %v", err, tt.wantErr)
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

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}
}
