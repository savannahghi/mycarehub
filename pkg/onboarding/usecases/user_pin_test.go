package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestUserPinUseCaseUnitTest_SetUserPIN(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx       context.Context
		pin       string
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_set_user_pin",
			args: args{
				ctx:       ctx,
				pin:       "1234",
				profileID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_non_digits_included",
			args: args{
				ctx:       ctx,
				pin:       "12b4",
				profileID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_wrong_number_of_digits",
			args: args{
				ctx:       ctx,
				pin:       "12",
				profileID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_save_pin",
			args: args{
				ctx:       ctx,
				pin:       "1234",
				profileID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_set_user_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}

			}

			if tt.name == "invalid:_unable_to_save_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("unable to save pin")
				}
			}

			isSet, err := i.UserPIN.SetUserPIN(tt.args.ctx, tt.args.pin, tt.args.profileID)

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

				if isSet != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, isSet)
					return
				}
			}

		})
	}
}

func TestUserPinUseCaseUnitTest_ResetUserPIN(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx   context.Context
		phone string
		PIN   string
		OTP   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_reset_user_pin",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_verify_otp_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_verify_otp_returns_false",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_cant_get_user_by_phone",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_check_has_pin_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_update_pin_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				PIN:   base.TestUserPin,
				OTP:   "588214",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_reset_user_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
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
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid:_update_pin_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
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
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("unable to update pin")
				}
			}

			if tt.name == "invalid:_verify_otp_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify otp")
				}
			}

			if tt.name == "invalid:_verify_otp_returns_false" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:_cant_get_user_by_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user by phone")
				}
			}

			if tt.name == "invalid:_check_has_pin_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get user pin")
				}
			}

			got, err := i.UserPIN.ResetUserPIN(tt.args.ctx, tt.args.phone, tt.args.PIN, tt.args.OTP)

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

				if got != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, got)
					return
				}
			}

		})
	}
}

func TestUserPinUseCaseImpl_ChangeUserPINUnitTest(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx   context.Context
		phone string
		pin   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_change_pin",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				pin:   base.TestUserPin,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_get_profile_by_phone_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_check_has_pin_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_update_pin_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_normalize_msisdn_fails",
			args: args{
				ctx:   ctx,
				phone: "+254721456789",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_change_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
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
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid:_normalize_msisdn_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize phone")
				}
			}

			if tt.name == "invalid:_get_profile_by_phone_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user by phone")
				}
			}

			if tt.name == "invalid:_check_has_pin_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get user pin")
				}
			}

			if tt.name == "invalid:_update_pin_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
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
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("unable to update pin")
				}
			}

			got, err := i.UserPIN.ChangeUserPIN(tt.args.ctx, tt.args.phone, tt.args.pin)

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

				if got != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, got)
					return
				}
			}

		})
	}
}

func TestUserPinUseCaseImpl_CheckHasPIN(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	type args struct {
		ctx       context.Context
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_check_has_pin",
			args: args{
				ctx:       ctx,
				profileID: "1234",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_get_pin_by_profile_returns_nil",
			args: args{
				ctx:       ctx,
				profileID: "1234",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_get_pin_by_profile_returns_error",
			args: args{
				ctx:       ctx,
				profileID: "1234",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_check_has_pin" {
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
			}

			if tt.name == "invalid:_get_pin_by_profile_returns_nil" {
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, nil
				}
			}

			if tt.name == "invalid:_get_pin_by_profile_returns_error" {
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get pin by profile")
				}
			}

			got, err := i.UserPIN.CheckHasPIN(tt.args.ctx, tt.args.profileID)

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

				if got != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, got)
					return
				}
			}

		})
	}
}

func TestUserPinUseCaseImpl_RequestPINReset(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.OtpResponse
		wantErr bool
	}{
		{
			name: "valid:_request_pin_reset",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: false,
			want:    &base.OtpResponse{OTP: "1234"},
		},
		{
			name: "invalid:_unable_to_normalize_phone",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_user_profile",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: true,
		},
		{
			name: "invalid:_check_has_pin_fails",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: true,
		},
		{
			name: "invalid:_generate_and_send_otp_fails",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_request_pin_reset" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
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
				fakeOtp.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return &base.OtpResponse{OTP: "1234"}, nil
				}
			}

			if tt.name == "invalid:_unable_to_normalize_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize phone")
				}

			}

			if tt.name == "invalid:_unable_to_get_user_profile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}

			if tt.name == "invalid:_check_has_pin_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get user pin")
				}
			}

			if tt.name == "invalid:_generate_and_send_otp_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
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
				fakeOtp.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return nil, fmt.Errorf("unable to generate and send otp")
				}
			}

			got, err := i.UserPIN.RequestPINReset(tt.args.ctx, tt.args.phone)

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

				if got.OTP != tt.want.OTP {
					t.Errorf("expected %v got %v  ", tt.want.OTP, got.OTP)
					return
				}
			}

		})
	}
}
