package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
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
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
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
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
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
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
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

func TestSignUpUseCasesImpl_CreateUserByPhone(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	phoneNumber := "+254777886622"
	pin := "1234"
	otp := "678251"

	validSignUpInput := &dto.SignUpInput{
		PhoneNumber: &phoneNumber,
		PIN:         &pin,
		Flavour:     base.FlavourConsumer,
		OTP:         &otp,
	}

	invalidPhoneNumber := "+254"
	invalidPin := ""
	invalidOTP := ""

	invalidSignUpInput := &dto.SignUpInput{
		PhoneNumber: &invalidPhoneNumber,
		PIN:         &invalidPin,
		Flavour:     base.FlavourConsumer,
		OTP:         &invalidOTP,
	}
	phone := gofakeit.Phone()

	type args struct {
		ctx   context.Context
		input *dto.SignUpInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:successfully_create_user_by_phone",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_verifyOTP",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:use_invalid_input",
			args: args{
				ctx:   ctx,
				input: invalidSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_create_user_profile",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_generate_auth_credentials",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_set_userPin",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_create_empty_supplier_profile",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_create_empty_customer_profile",
			args: args{
				ctx:   ctx,
				input: validSignUpInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_create_user_by_phone" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						PrimaryPhone: &phone,
					}, nil
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

				// Mock SetUserPin begins here
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

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "password"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
				// Finished mocking SetUserPin

				fakeRepo.CreateEmptySupplierProfileFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return &base.Customer{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string,
					allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
						ID:            uuid.New().String(),
						AllowWhatsApp: true,
						AllowTextSMS:  true,
						AllowEmail:    true,
						AllowPush:     true,
					}, nil
				}
				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_verifyOTP" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:use_invalid_input" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:fail_to_check_ifPhoneNumberExists" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid:fail_to_create_user_profile" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("fail to create user profile")
				}
			}

			if tt.name == "invalid:fail_to_generate_auth_credentials" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
					}, nil
				}

				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *base.UserProfile) (*base.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("failed to generate auth credentials")
				}
			}

			if tt.name == "invalid:fail_to_set_userPin" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						PrimaryPhone: &phone,
					}, nil
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

				// Mock SetUserPin begins here
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

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "password"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("failed to save user pin")
				}
			}

			if tt.name == "invalid:fail_to_create_empty_supplier_profile" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						PrimaryPhone: &phone,
					}, nil
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

				// Mock SetUserPin begins here
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

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "password"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
				// Finished mocking SetUserPin

				fakeRepo.CreateEmptySupplierProfileFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to create empty supplier profile")
				}
			}

			if tt.name == "fail_to_create_empty_customer_profile" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						DisplayName: "John Doe",
						Email:       "johndoe@gmail.com",
						PhoneNumber: phoneNumber,
					}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						PrimaryPhone: &phone,
					}, nil
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

				// Mock SetUserPin begins here
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

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "password"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
				// Finished mocking SetUserPin

				fakeRepo.CreateEmptySupplierProfileFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return nil, fmt.Errorf("failed to create empty customer profile")
				}
			}

			_, err := i.Signup.CreateUserByPhone(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.CreateUserByPhone() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSignUpUseCasesImpl_VerifyPhoneNumber(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:successfully_verify_a_phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+254777886622",
			},
			wantErr: false,
		},
		{
			name: "invalid:_phone_number_is_empty",
			args: args{
				ctx:   ctx,
				phone: "+",
			},
			wantErr: true,
		},
		{
			name: "invalid:_user_phone_already_exists",
			args: args{
				ctx:   ctx,
				phone: "+254777886622",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_generate_and_send_otp",
			args: args{
				ctx:   ctx,
				phone: "+254777886622",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_check_if_phone_exists",
			args: args{
				ctx:   ctx,
				phone: "+254777886622",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_verify_a_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}

				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return &base.OtpResponse{OTP: "1234"}, nil
				}
			}

			if tt.name == "invalid:_phone_number_is_empty" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number")
				}
			}

			if tt.name == "invalid:_user_phone_already_exists" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid:_unable_to_check_if_phone_exists" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("unable to check if phone exists")
				}
			}

			if tt.name == "invalid:fail_to_generate_and_send_otp" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}

				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return nil, fmt.Errorf("failed to generate and send otp")
				}
			}

			_, err := i.Signup.VerifyPhoneNumber(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSignUpUseCasesImpl_RemoveUserByPhoneNumber(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		phone string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:successfully_RemoveUserByPhoneNumber",
			args: args{
				ctx:   ctx,
				phone: "+254799739102",
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_RemoveUserByPhoneNumber",
			args: args{
				ctx:   ctx,
				phone: "+254799739102",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_normalize_phone",
			args: args{
				ctx:   ctx,
				phone: "+",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_RemoveUserByPhoneNumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.PurgeUserByPhoneNumberFn = func(ctx context.Context, phone string) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_RemoveUserByPhoneNumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeRepo.PurgeUserByPhoneNumberFn = func(ctx context.Context, phone string) error {
					return fmt.Errorf("failed to purge user by phonenumber")
				}
			}

			if tt.name == "invalid:fail_to_normalize_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("failed to normalize phonenumber")
				}
			}

			err := i.Signup.RemoveUserByPhoneNumber(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RemoveUserByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSignUpUseCasesImpl_SetPhoneAsPrimary(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		phone string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:set_primary_phoneNumber",
			args: args{
				ctx:   ctx,
				phone: "+254795941530",
				otp:   "567291",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_normalize_phoneNumber",
			args: args{
				ctx:   ctx,
				phone: "+",
				otp:   "567291",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_primary_phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+25463728192",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:set_primary_phoneNumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				// Begin Mocking SetPrimaryPhoneNumber
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				// Begin Mocking UpdatePrimaryPhoneNumber
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254755889922"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "ABCDE",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}

				fakeRepo.UpdateSecondaryPhoneNumbersFn = func(ctx context.Context, id string, phoneNumbers []string) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_normalize_phoneNumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("failed to normalize phonenumber")
				}
			}

			if tt.name == "invalid:_unable_to_update_primary_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				// Begin Mocking SetPrimaryPhoneNumber
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				// Begin Mocking UpdatePrimaryPhoneNumber
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254755889922"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "ABCDE",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return fmt.Errorf("failed to update primary phone")
				}
			}

			got, err := i.Signup.SetPhoneAsPrimary(tt.args.ctx, tt.args.phone, tt.args.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.SetPhoneAsPrimary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.SetPhoneAsPrimary() = %v, want %v", got, tt.want)
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

func TestSignUpUseCasesImpl_RegisterPushToken(t *testing.T) {
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
			name: "valid:register_pushtoken",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:nil_token",
			args: args{
				ctx:   ctx,
				token: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_userProfile",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_loggedInUser",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:register_pushtoken" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:        "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						Suspended: false,
					}, nil
				}
				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
					return nil
				}
			}

			if tt.name == "invalid:nil_token" {
				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
					return fmt.Errorf("failed to register push token")
				}
			}

			if tt.name == "invalid:fail_to_get_userProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "invalid:fail_to_get_loggedInUser" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			got, err := i.Signup.RegisterPushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RegisterPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.RegisterPushToken() = %v, want %v", got, tt.want)
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

func TestSignUpUseCasesImpl_CompleteSignup(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	userFirstName := "John"
	userLastName := "Doe"

	type args struct {
		ctx     context.Context
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:successfully_complete_signup",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_get_userProfile",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_loggedInUser",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:missing_bioData",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:invalid_flavour",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_AddCustomerSupplierERPAccount",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_FetchDefaultCurrency",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_complete_signup" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uuid.New().String(),
							},
						},
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.EnsureTopicsExistFn = func(ctx context.Context, topicIDs []string) error {
					return nil
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_userProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "invalid:fail_to_get_loggedInUser" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "invalid:missing_bioData" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uuid.New().String(),
							},
						},
					}, nil
				}

			}

			if tt.name == "invalid:invalid_flavour" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("invalid flavour defined")
				}
			}

			if tt.name == "invalid:fail_to_AddCustomerSupplierERPAccount" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uuid.New().String(),
							},
						},
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakeEPRSvc.CreateERPCustomerFn = func(
					ctx context.Context,
					customerPayload dto.CustomerPayload,
					UID string,
				) (*base.Customer, error) {
					return nil, fmt.Errorf("failed to add customer supplier ERP account")
				}
			}

			if tt.name == "invalid:fail_to_FetchDefaultCurrency" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uuid.New().String(),
							},
						},
						UserBioData: base.BioData{
							FirstName: &userFirstName,
							LastName:  &userLastName,
						},
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					return nil, fmt.Errorf("failed to fetch default currency")
				}
			}

			got, err := i.Signup.CompleteSignup(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.CompleteSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.CompleteSignup() = %v, want %v", got, tt.want)
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

func TestSignUpUseCasesImpl_GetUserRecoveryPhoneNumbers(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.AccountRecoveryPhonesResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_GetUserRecoveryPhoneNumbers",
			args: args{
				ctx:   ctx,
				phone: "+254766228822",
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_normalize_phone",
			args: args{
				ctx:   ctx,
				phone: "+254766228822",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_userProfile",
			args: args{
				ctx:   ctx,
				phone: "+254766228822",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_GetUserRecoveryPhoneNumbers" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0744610111", "0794959697",
						},
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_normalize_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("failed to normalize phone")
				}
			}

			if tt.name == "invalid:fail_to_get_userProfile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}
			got, err := i.Signup.GetUserRecoveryPhoneNumbers(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.GetUserRecoveryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("returned a nil account recovery phone response")
					return
				}
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestSignUpUseCasesImpl_UpdateUserProfile(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	photoUploadID := "somePhotoUploadID"
	firstName := "John"
	lastName := "Doe"
	gender := base.GenderMale
	dateOfBirth := base.Date{
		Year:  1990,
		Month: 3,
		Day:   10,
	}
	phone := gofakeit.Phone()
	validInput := &dto.UserProfileInput{
		PhotoUploadID: &photoUploadID,
		DateOfBirth:   &dateOfBirth,
		Gender:        &gender,
		FirstName:     &firstName,
		LastName:      &lastName,
	}
	invalidInput := &dto.UserProfileInput{
		PhotoUploadID: nil,
		DateOfBirth:   nil,
		Gender:        nil,
		FirstName:     nil,
		LastName:      nil,
	}
	type args struct {
		ctx   context.Context
		input *dto.UserProfileInput
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid:successfully_update_userProfile",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "invalid:missing_biodata",
			args: args{
				ctx:   ctx,
				input: invalidInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_updatePhotoUploadID",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getUserProfile",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getLoggedInUser",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_update_userProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           uuid.New().String(),
						PrimaryPhone: &phone,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           uuid.New().String(),
						PrimaryPhone: &phone,
					}, nil
				}

				fakeRepo.UpdatePhotoUploadIDFn = func(ctx context.Context, id string, uploadID string) error {
					return nil
				}

				// Begin mocking UpdateBioData
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryPhone: &phone,
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}

			}

			if tt.name == "invalid:missing_biodata" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePhotoUploadIDFn = func(ctx context.Context, id string, uploadID string) error {
					return nil
				}

				// Begin mocking UpdateBioData
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return fmt.Errorf("failed to update biodata")
				}
			}

			if tt.name == "invalid:fail_to_updatePhotoUploadID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePhotoUploadIDFn = func(ctx context.Context, id string, uploadID string) error {
					return fmt.Errorf("failed to update the photo upload ID")
				}
			}

			if tt.name == "invalid:fail_to_getUserProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get a user profile")
				}
			}

			if tt.name == "invalid:fail_to_getLoggedInUser" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			got, err := i.Signup.UpdateUserProfile(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.UpdateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("returned a nil account recovery phone response")
					return
				}
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}
