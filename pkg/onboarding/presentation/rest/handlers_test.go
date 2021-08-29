package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/profileutils"

	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	engagementMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement/mock"

	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	pubsubmessagingMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/rest"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	mockRepo "github.com/savannahghi/onboarding/pkg/onboarding/repository/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	adminSrv "github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin"
)

var fakeRepo mockRepo.FakeOnboardingRepository
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var serverUrl = "http://localhost:5000"
var fakePubSub pubsubmessagingMock.FakeServicePubSub

// InitializeFakeOnboardingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboardingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt
	var ps pubsubmessaging.ServicePubSub = &fakePubSub

	profile := usecases.NewProfileUseCase(r, ext, engagementSvc, ps)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	userpin := usecases.NewUserPinUseCase(r, profile, ext, pinExt, engagementSvc)
	su := usecases.NewSignUpUseCases(r, profile, userpin, ext, engagementSvc, ps)
	nhif := usecases.NewNHIFUseCases(r, profile, ext, engagementSvc)
	role := usecases.NewRoleUseCases(r, ext)

	adminSrv := adminSrv.NewService(ext)

	i, err := interactor.NewOnboardingInteractor(
		profile, su, login,
		survey, userpin,
		engagementSvc, nhif, ps, adminSrv,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}

func composeValidPhonePayload(t *testing.T, phone string) *bytes.Buffer {
	phoneNumber := struct {
		PhoneNumber string
	}{
		PhoneNumber: phone,
	}
	bs, err := json.Marshal(phoneNumber)
	if err != nil {
		t.Errorf("unable to marshal phoneNumber to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeValidRolePayload(t *testing.T, phone string, role profileutils.RoleType) *bytes.Buffer {
	payload := &dto.RolePayload{
		PhoneNumber: &phone,
		Role:        &role,
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal token string to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeSignupPayload(
	t *testing.T,
	phone, pin, otp string,
	flavour feedlib.Flavour,
) *bytes.Buffer {
	payload := dto.SignUpInput{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
		OTP:         &otp,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeChangePinPayload(t *testing.T, phone, pin, otp string) *bytes.Buffer {
	payload := domain.ChangePINRequest{
		PhoneNumber: phone,
		PIN:         pin,
		OTP:         otp,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeRefreshTokenPayload(t *testing.T, token *string) *bytes.Buffer {
	refreshToken := &dto.RefreshTokenPayload{RefreshToken: token}
	bs, err := json.Marshal(refreshToken)
	if err != nil {
		t.Errorf("unable to marshal token string to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeUIDPayload(t *testing.T, uid *string) *bytes.Buffer {
	uidPayload := &dto.UIDPayload{UID: uid}
	bs, err := json.Marshal(uidPayload)
	if err != nil {
		t.Errorf("unable to marshal token string to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

// TODO: restore
// func composePushTokenPayload(t *testing.T, UID, token string) *bytes.Buffer {
// 	payload := &dto.PushTokenPayload{
// 		PushToken: token,
// 		UID:       UID,
// 	}
// 	bs, err := json.Marshal(payload)
// 	if err != nil {
// 		t.Errorf("unable to marshal token string to JSON: %s", err)
// 	}
// 	return bytes.NewBuffer(bs)
// }

func composeLoginPayload(t *testing.T, phone, pin string, flavour feedlib.Flavour) *bytes.Buffer {
	payload := dto.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func composeSendRetryOTPPayload(t *testing.T, phone string, retryStep int) *bytes.Buffer {
	payload := dto.SendRetryOTPPayload{
		Phone:     &phone,
		RetryStep: &retryStep,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return nil
	}
	return bytes.NewBuffer(bs)
}

func composeCoversUpdatePayload(t *testing.T, payload *dto.UpdateCoversPayload) *bytes.Buffer {

	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return nil
	}
	return bytes.NewBuffer(bs)
}

func composeSetPrimaryPhoneNumberPayload(t *testing.T, phone, otp string) *bytes.Buffer {
	payload := dto.SetPrimaryPhoneNumberPayload{
		PhoneNumber: &phone,
		OTP:         &otp,
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return nil
	}
	return bytes.NewBuffer(bs)
}

func TestHandlersInterfacesImpl_VerifySignUpPhoneNumber(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// valid:_successfully_verifies_a_phone_number
	payload := composeValidPhonePayload(t, interserviceclient.TestUserPhoneNumber)

	// payload 2
	payload2 := composeValidPhonePayload(t, interserviceclient.TestUserPhoneNumberWithPin)

	// payload 3
	payload3 := composeValidPhonePayload(t, "0700100200")

	// payload 4
	payload4 := composeValidPhonePayload(t, "0700600300")

	// payload 5
	payload5 := composeValidPhonePayload(t, "*")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid:_phone_number_is_empty",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload5,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid:_successfully_verifies_a_phone_number",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_user_phone_already_exists",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_unable_to_check_if_phone_exists",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload3,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_sending_and_generation_of_OTP_fails",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload4,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
				phone := "+254721123123"
				return &phone, nil
			}
			if tt.name == "invalid:_phone_number_is_empty" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number")
				}
			}
			if tt.name == "valid:_successfully_verifies_a_phone_number" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return &profileutils.OtpResponse{OTP: "1234"}, nil
				}
			}
			// we mock `CheckPhoneExists` to return true
			// we dont need to mock `GenerateAndSendOTP` because we won't get there
			if tt.name == "invalid:_user_phone_already_exists" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}
			}
			// we mock `CheckPhoneExists` to return error,
			if tt.name == "invalid:_unable_to_check_if_phone_exists" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("unable to check if phone exists")
				}
			}
			// we mock `GenerateAndSendOTP` to return error,
			if tt.name == "invalid:_sending_and_generation_of_OTP_fails" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("unable generate and send otp")
				}
			}
			// Our handlers satisfy http.Handler, so we can call its ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.VerifySignUpPhoneNumber()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
			if !tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
				if !tt.wantErr {
					_, ok := data["error"]
					if ok {
						t.Errorf("error not expected")
						return
					}
				}
			}

		})

	}
}

func TestHandlersInterfacesImpl_CreateUserWithPhoneNumber(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	// payload
	pin := "2030"
	flavour := feedlib.FlavourPro
	otp := "1234"
	phoneNumber := interserviceclient.TestUserPhoneNumber
	payload := composeSignupPayload(t, phoneNumber, pin, otp, flavour)

	// payload 2
	pin2 := "1000"
	flavour2 := feedlib.FlavourConsumer
	otp2 := "9000"
	phoneNumber2 := "+254720125456"
	payload2 := composeSignupPayload(t, phoneNumber2, pin2, otp2, flavour2)

	// payload 3
	pin3 := "2000"
	flavour3 := feedlib.FlavourConsumer
	otp3 := "3000"
	phoneNumber3 := "+254721100200"
	payload3 := composeSignupPayload(t, phoneNumber3, pin3, otp3, flavour3)

	// payload6
	pin6 := "0000"
	flavour6 := feedlib.FlavourConsumer
	otp6 := "9520"
	phoneNumber6 := "+254721410589"
	payload6 := composeSignupPayload(t, phoneNumber6, pin6, otp6, flavour6)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_create_user_by_phone",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "invalid_unable_to_verify_otp",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_verify_otp_returns_false",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload3,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_get_or_create_phone_returns_error",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload6,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			if tt.name == "valid:_successfully_create_user_by_phone" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{
						UID:         "1106f10f-bea6-4fa3-bdba-16b1e39bd318",
						DisplayName: "kalulu juha",
						Email:       "juha@gmail.com",
						PhoneNumber: "0756232452",
						PhotoURL:    "",
						ProviderID:  "google.com",
					}, nil
				}
				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *profileutils.UserProfile) (*profileutils.AuthCredentialResponse, error) {
					return &profileutils.AuthCredentialResponse{
						UID:          "5550",
						RefreshToken: "55550",
					}, nil
				}
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}

				fakeRepo.CreateEmptySupplierProfileFn = func(ctx context.Context, profileID string) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{
						ID:         "5550",
						SupplierID: "5555",
					}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{
						ID:         "0000",
						CustomerID: "22222",
					}, nil
				}
				// should return a profile with an ID
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				// should return true
				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string,
					allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:            "111",
						ProfileID:     "profile-id",
						AllowWhatsApp: true,
						AllowEmail:    true,
						AllowTextSMS:  true,
						AllowPush:     true,
					}, nil
				}

				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					roles := []profileutils.Role{
						{
							ID: uuid.NewString(),
						},
					}
					return &roles, nil
				}

				fakeRepo.GetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:            "111",
						ProfileID:     "profile-id",
						AllowWhatsApp: true,
						AllowEmail:    true,
						AllowTextSMS:  true,
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

			// mock VerifyOTP to return an error
			if tt.name == "invalid_unable_to_verify_otp" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify otp")
				}
			}

			// mock VerifyOTP to return an false
			if tt.name == "invalid_verify_otp_returns_false" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			// mock `GetOrCreatePhoneNumberUser` to return an error
			if tt.name == "invalid_get_or_create_phone_returns_error" {
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return nil, fmt.Errorf("unable to create user")
				}
			}

			// Our handlers satisfy http.Handler, so we can call its ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.CreateUserWithPhoneNumber()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
			if !tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}
		})
	}
}

func TestHandlersInterfacesImpl_UserRecoveryPhoneNumbers(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload 1
	payload := composeValidPhonePayload(t, interserviceclient.TestUserPhoneNumber)

	// payload 2
	payload2 := composeValidPhonePayload(t, "0710100595")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_get_a_recovery_phone",
			args: args{
				url:        fmt.Sprintf("%s/user_recovery_phonenumbers", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_unable_to_get_profile",
			args: args{
				url:        fmt.Sprintf("%s/user_recovery_phonenumbers", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			if tt.name == "valid:_successfully_get_a_recovery_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}
			}

			// we set GetUserProfileByPhoneNumber to return an error
			if tt.name == "invalid:_unable_to_get_profile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to retrieve profile")
				}
			}

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.UserRecoveryPhoneNumbers()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
			if !tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
				if !tt.wantErr {
					_, ok := data["error"]
					if ok {
						t.Errorf("error not expected")
						return
					}
				}
			}
		})
	}
}

func TestHandlersInterfacesImpl_RequestPINReset(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload successfully_request_pin_reset
	payload := composeValidPhonePayload(t, interserviceclient.TestUserPhoneNumber)
	// _phone_number_invalid
	payload1 := composeValidPhonePayload(t, "")
	//invalid:_inable_to_get_primary_phone
	payload2 := composeValidPhonePayload(t, "0725123456")
	//invalid:check_has_pin_failed
	payload3 := composeValidPhonePayload(t, "0700100400")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:successfully_request_pin_reset",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_phone_number_invalid",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_inable_to_get_primary_phone",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:check_has_pin_failed",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload3,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:otp_generation_fails",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload3,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
				phone := "+254721123123"
				return &phone, nil
			}
			if tt.name == "invalid:_phone_number_invalid" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("invalid phone number")
				}
			}

			if tt.name == "valid:successfully_request_pin_reset" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return &profileutils.OtpResponse{OTP: "1234"}, nil
				}
			}

			if tt.name == "invalid:otp_generation_fails" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("unable to generate otp")
				}
			}

			if tt.name == "invalid:_inable_to_get_primary_phone" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to fetch profile")
				}
			}

			if tt.name == "invalid:check_has_pin_failed" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to retrieve pin")
				}
			}

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.RequestPINReset()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
			if !tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
				if !tt.wantErr {
					_, ok := data["error"]
					if ok {
						t.Errorf("error not expected")
						return
					}
				}
			}

		})
	}
}

func TestHandlersInterfacesImpl_ResetPin(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload
	phone := "0712456784"
	pin := "1897"
	otp := "000548"
	payload := composeChangePinPayload(t, phone, pin, otp)
	// payload2
	phone1 := "0710472196"
	pin1 := "02222"
	otp1 := "0002358"
	payload1 := composeChangePinPayload(t, phone1, pin1, otp1)
	// payload3
	phone2 := ""
	pin2 := ""
	otp2 := "6666"
	payload2 := composeChangePinPayload(t, phone2, pin2, otp2)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid:empty_payload",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "valid:successfully_reset_pin",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "invalid:unable_to_update_pin",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			if tt.name == "valid:successfully_reset_pin" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			// we set `UpdatePIN` to return an error
			if tt.name == "invalid:unable_to_update_pin" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("unable to update pin")
				}
			}
			svr := h.ResetPin()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

func TestHandlersInterfacesImpl_RefreshToken(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	token := "10c17f3b-a9a9-431c-ad0a-94c684eccd85"
	payload := composeRefreshTokenPayload(t, &token)

	token1 := "b5c52b32-7dd5-4dd5-9ddb-44cac9701d6c"
	payload1 := composeRefreshTokenPayload(t, &token1)

	token2 := "*"
	payload2 := composeRefreshTokenPayload(t, &token2)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_refresh_token",
			args: args{
				url:        fmt.Sprintf("%s/refresh_token", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_refresh_token_fails",
			args: args{
				url:        fmt.Sprintf("%s/refresh_token", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_refresh_token_with_invalid_payload",
			args: args{
				url:        fmt.Sprintf("%s/refresh_token", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			if tt.name == "valid:_successfully_refresh_token" {
				fakeRepo.ExchangeRefreshTokenForIDTokenFn = func(ctx context.Context, token string) (*profileutils.AuthCredentialResponse, error) {
					return &profileutils.AuthCredentialResponse{
						UID:          "5550",
						RefreshToken: "55550",
					}, nil
				}
			}

			if tt.name == "invalid:_refresh_token_fails" {
				fakeRepo.ExchangeRefreshTokenForIDTokenFn = func(ctx context.Context, token string) (*profileutils.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("unable to refresh token")
				}
			}

			svr := h.RefreshToken()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

func TestHandlersInterfacesImpl_GetUserProfileByUID(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	uid := "db963177-21b2-489f-83e6-3521bf5db516"
	payload := composeUIDPayload(t, &uid)

	uid1 := "584799be-97c5-4aa4-8b0f-094990bd55b1"
	payload1 := composeUIDPayload(t, &uid1)

	uid2 := "*"
	payload2 := composeUIDPayload(t, &uid2)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_get_profile_by_uid",
			args: args{
				url:        fmt.Sprintf("%s/user_profile", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_unable_to_get_profile_by_uid",
			args: args{
				url:        fmt.Sprintf("%s/user_profile", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_get_profile_by_uid_with_invalid_payload",
			args: args{
				url:        fmt.Sprintf("%s/user_profile", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			if tt.name == "valid:_successfully_get_profile_by_uid" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_by_uid" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			svr := h.GetUserProfileByUID()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func composePhoneOrEmailPayload(t *testing.T, payload *dto.RetrieveUserProfileInput) *bytes.Buffer {
	emailPayload := &dto.RetrieveUserProfileInput{
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}
	bs, err := json.Marshal(emailPayload)
	if err != nil {
		t.Errorf("unable to marshal email string to JSON: %s", err)
	}
	return bytes.NewBuffer(bs)
}

func TestHandlersInterfacesImpl_GetUserProfileByPhoneOrEmail(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	email := "test@kathurima.com"

	validPayload := &dto.RetrieveUserProfileInput{
		Email:       &email,
		PhoneNumber: nil,
	}
	payload := composePhoneOrEmailPayload(t, validPayload)

	invalidPayload := &dto.RetrieveUserProfileInput{
		Email:       nil,
		PhoneNumber: nil,
	}

	payload2 := composePhoneOrEmailPayload(t, invalidPayload)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:get_profile_by_email",
			args: args{
				url:        fmt.Sprintf("%s/get_profile_by_email", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},

		{
			name: "invalid:_unable_to_get_profile_by_email",
			args: args{
				url:        fmt.Sprintf("%s/get_profile_by_email", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			if tt.name == "valid:get_profile_by_email" {
				fakeRepo.UpdateUserProfileEmailFn = func(ctx context.Context, phone, email string) error {
					return nil
				}
				fakeRepo.GetUserProfileByPhoneOrEmailFn = func(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						PrimaryEmailAddress: &email,
					}, nil
				}
			}
			if tt.name == "invalid:_unable_to_get_profile_by_email" {
				fakeRepo.GetUserProfileByPhoneOrEmailFn = func(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to retrieve user profile")
				}
			}

			svr := h.GetUserProfileByPhoneOrEmail()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_SendOTP(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	payload := composeValidPhonePayload(t, interserviceclient.TestUserPhoneNumber)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_send_otp",
			args: args{
				url:        fmt.Sprintf("%s/send_otp", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_unable_to_send_otp",
			args: args{
				url:        fmt.Sprintf("%s/send_otp", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			if tt.name == "valid:_successfully_send_otp" {
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return &profileutils.OtpResponse{OTP: "1234"}, nil
				}
			}

			if tt.name == "invalid:_unable_to_send_otp" {
				fakeEngagementSvs.GenerateAndSendOTPFn = func(ctx context.Context, phone string, appID *string) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("unable to send otp")
				}
			}

			svr := h.SendOTP()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_LoginByPhone(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload
	phone := "0712456784"
	pin := "1897"
	flavour := feedlib.FlavourPro
	payload := composeLoginPayload(t, phone, pin, flavour)

	// payload1 : invalid:_get_userprofile_by_primary_phone_fails
	phone1 := "0708598520"
	pin1 := "1800"
	flavour1 := feedlib.FlavourConsumer
	payload1 := composeLoginPayload(t, phone1, pin1, flavour1)

	// payload2 : invalid:_get_pinbyprofileid_fails
	phone2 := "0708590000"
	pin2 := "1000"
	flavour2 := feedlib.FlavourConsumer
	payload2 := composeLoginPayload(t, phone2, pin2, flavour2)

	// payload4 invalid:_pin_mismatch
	phone4 := "0702960230"
	pin4 := "1023"
	flavour4 := feedlib.FlavourConsumer
	payload4 := composeLoginPayload(t, phone4, pin4, flavour4)

	// payload5 invalid:_generate_auth_credentials_fails
	phone5 := "0705222888"
	pin5 := "1093"
	flavour5 := feedlib.FlavourConsumer
	payload5 := composeLoginPayload(t, phone5, pin5, flavour5)

	// payload7 invalid:_invalid_flavour_used
	phone7 := "0712456784"
	pin7 := "1897"
	payload7 := composeLoginPayload(t, phone7, pin7, "invalidFlavour")
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_login_by_phone",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_get_userprofile_by_primary_phone_fails",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_get_pinbyprofileid_fails",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_pin_mismatch",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload4,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_generate_auth_credentials_fails",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload5,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_invalid_flavour_used",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload7,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// we mock the required methods for a valid case
			if tt.name == "valid:_successfully_login_by_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}
				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *profileutils.UserProfile) (*profileutils.AuthCredentialResponse, error) {
					return &profileutils.AuthCredentialResponse{
						UID: "5550",
						// IDToken:      "555",
						RefreshToken: "55550",
					}, nil
				}
				fakeRepo.GetCustomerOrSupplierProfileByProfileIDFn = func(ctx context.Context, flavour feedlib.Flavour, profileID string) (*profileutils.Customer, *profileutils.Supplier, error) {
					return &profileutils.Customer{
							ID: "5550",
						}, &profileutils.Supplier{
							ID: "5550",
						}, nil
				}
				fakeRepo.GetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:            "111",
						ProfileID:     "profile-id",
						AllowWhatsApp: true,
						AllowEmail:    true,
						AllowTextSMS:  true,
						AllowPush:     true,
					}, nil
				}
			}

			if tt.name == "invalid:_get_userprofile_by_primary_phone_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}

			if tt.name == "invalid:_get_pinbyprofileid_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return nil, fmt.Errorf("unable to get pin by profileID")
				}

			}

			if tt.name == "invalid:_pin_mismatch" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
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

			if tt.name == "invalid:_generate_auth_credentials_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return true
				}
				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string, profile *profileutils.UserProfile) (*profileutils.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("unable to generate auth credentials")
				}
			}

			if tt.name == "invalid:_invalid_flavour_used" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("invalid flavour defined")
				}
			}

			svr := h.LoginByPhone()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

			if !tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
				if !tt.wantErr {
					_, ok := data["error"]
					if ok {
						t.Errorf("error not expected")
						return
					}
				}
			}
		})
	}
}

func TestHandlersInterfacesImpl_SendRetryOTP(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	// valid payload
	validPayload := composeSendRetryOTPPayload(t, interserviceclient.TestUserPhoneNumber, 1)

	invalidPayload := composeSendRetryOTPPayload(t, "", 2)
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_send_retry_otp",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_unable_to_send_otp",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_unable_to_send_otp_due_to_missing_msisdn",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create a new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:_successfully_send_retry_otp" {
				fakeEngagementSvs.SendRetryOTPFn = func(ctx context.Context, msisdn string, retryStep int, appID *string) (*profileutils.OtpResponse, error) {
					return &profileutils.OtpResponse{
						OTP: "123456",
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_send_otp" {
				fakeEngagementSvs.SendRetryOTPFn = func(ctx context.Context, msisdn string, retryStep int, appID *string) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("unable to send OTP")
				}
			}

			if tt.name == "invalid:_unable_to_send_otp_due_to_missing_msisdn" {
				fakeEngagementSvs.SendRetryOTPFn = func(ctx context.Context, msisdn string, retryStep int, appID *string) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("unable to send OTP")
				}
			}

			svr := h.SendRetryOTP()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

func TestHandlersInterfacesImpl_LoginAnonymous(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	validPayload := composeLoginPayload(t, "", "", feedlib.FlavourConsumer)
	invalidPayload := composeLoginPayload(t, "", "", " ")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_login_as_anonymous",
			args: args{
				url:        fmt.Sprintf("%s/login_anonymous", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_invalid_flavour_defined",
			args: args{
				url:        fmt.Sprintf("%s/login_anonymous", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_missing_flavour",
			args: args{
				url:        fmt.Sprintf("%s/login_anonymous", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if tt.name == "valid:_successfully_login_as_anonymous" {
				fakeRepo.GenerateAuthCredentialsForAnonymousUserFn = func(ctx context.Context) (*profileutils.AuthCredentialResponse, error) {
					return &profileutils.AuthCredentialResponse{
						UID:          "6660",
						RefreshToken: "6660",
					}, nil
				}
			}

			if tt.name == "invalid:_invalid_flavour_defined" {
				fakeRepo.GenerateAuthCredentialsForAnonymousUserFn = func(ctx context.Context) (*profileutils.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("an invalid `flavour` defined")
				}
			}

			if tt.name == "invalid:_missing_flavour" {
				fakeRepo.GenerateAuthCredentialsForAnonymousUserFn = func(ctx context.Context) (*profileutils.AuthCredentialResponse, error) {
					return nil, fmt.Errorf("expected `flavour` to be defined")
				}
			}

			response := httptest.NewRecorder()

			svr := h.LoginAnonymous()

			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

func TestHandlersInterfacesImpl_UpdateCovers(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	invalidUID := " "
	uid := "5cf354a2-1d3e-400d-8716-7e2aead29f2c"
	payerName := "Payer Name"
	memberName := "Member Name"
	memberNumber := "5678"
	payerSladeCode := 1234
	beneficiaryID := 15689
	effectivePolicyNumber := "14582"

	validFromString := "2021-01-01T00:00:00+03:00"
	validFrom, err := time.Parse(time.RFC3339, validFromString)
	if err != nil {
		t.Errorf("failed parse date string: %v", err)
		return
	}

	validToString := "2022-01-01T00:00:00+03:00"
	validTo, err := time.Parse(time.RFC3339, validToString)
	if err != nil {
		t.Errorf("failed parse date string: %v", err)
		return
	}

	updateCoversPayloadValid := &dto.UpdateCoversPayload{
		UID:                   &uid,
		PayerName:             &payerName,
		PayerSladeCode:        &payerSladeCode,
		MemberName:            &memberName,
		MemberNumber:          &memberNumber,
		BeneficiaryID:         &beneficiaryID,
		EffectivePolicyNumber: &effectivePolicyNumber,
		ValidFrom:             &validFrom,
		ValidTo:               &validTo,
	}

	updateCoversPayloadInValid := &dto.UpdateCoversPayload{
		UID:                   &invalidUID,
		PayerName:             &payerName,
		PayerSladeCode:        &payerSladeCode,
		MemberName:            &memberName,
		MemberNumber:          &memberNumber,
		BeneficiaryID:         &beneficiaryID,
		EffectivePolicyNumber: &effectivePolicyNumber,
		ValidFrom:             &validFrom,
		ValidTo:               &validTo,
	}

	validPayload := composeCoversUpdatePayload(t, updateCoversPayloadValid)
	inValidPayload := composeCoversUpdatePayload(t, updateCoversPayloadInValid)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_Successfully_update_covers",
			args: args{
				url:        fmt.Sprintf("%s/update_covers", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "invalid:_update_covers_fails",
			args: args{
				url:        fmt.Sprintf("%s/update_covers", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_get_user_profile_by_UID_fails",
			args: args{
				url:        fmt.Sprintf("%s/update_covers", serverUrl),
				httpMethod: http.MethodPost,
				body:       inValidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:_Successfully_update_covers" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "f4f39af7",
					}, nil
				}

				fakeRepo.UpdateCoversFn = func(ctx context.Context, id string, covers []profileutils.Cover) error {
					return nil
				}
			}

			if tt.name == "invalid:_get_user_profile_by_UID_fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user UID")
				}
			}

			if tt.name == "invalid:_update_covers_fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: "f4f39af7",
					}, nil
				}

				fakeRepo.UpdateCoversFn = func(ctx context.Context, id string, covers []profileutils.Cover) error {
					return fmt.Errorf("unable to update covers")
				}
			}

			svr := h.UpdateCovers()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

func TestHandlersInterfacesImpl_RemoveUserByPhoneNumber(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	primaryPhone := "+254711445566"
	validPayload := composeValidPhonePayload(t, primaryPhone)
	validPayload1 := composeValidPhonePayload(t, "+254777882200")
	validPayload2 := composeValidPhonePayload(t, "+")

	invalidPayload := composeValidPhonePayload(t, " ")
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_Successfully_remove_user_by_phone",
			args: args{
				url:        fmt.Sprintf("%s/remove_user", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_unable_to_remove_user_by_phone",
			args: args{
				url:        fmt.Sprintf("%s/remove_user", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		{
			name: "invalid:_empty_phonenumber",
			args: args{
				url:        fmt.Sprintf("%s/remove_user", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_unable_to_check_if_phone_exists",
			args: args{
				url:        fmt.Sprintf("%s/remove_user", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:_Successfully_remove_user_by_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeRepo.PurgeUserByPhoneNumberFn = func(ctx context.Context, phone string) error {
					return nil
				}

			}

			if tt.name == "invalid:_unable_to_remove_user_by_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeRepo.PurgeUserByPhoneNumberFn = func(ctx context.Context, phone string) error {
					return fmt.Errorf("unable to purge user by phonenumber")
				}

			}

			if tt.name == "invalid:_empty_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number")
				}
			}

			if tt.name == "invalid:_unable_to_check_if_phone_exists" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "0788554422"
					return &phone, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("the phone does not exist")
				}

			}

			svr := h.RemoveUserByPhoneNumber()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_SetPrimaryPhoneNumber(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	primaryPhone := "+254701567839"
	otp := "890087"
	validPayload := composeSetPrimaryPhoneNumberPayload(t, primaryPhone, otp)

	primaryPhone1 := "+254765738293"
	otp1 := "345678"
	payload1 := composeSetPrimaryPhoneNumberPayload(t, primaryPhone1, otp1)

	primaryPhone2 := " "
	otp2 := " "
	payload2 := composeSetPrimaryPhoneNumberPayload(t, primaryPhone2, otp2)
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_successfully_set_a_primary_phonenumber",
			args: args{
				url:        fmt.Sprintf("%s/set_primary_phonenumber", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_fail_to_set_a_primary_phonenumber",
			args: args{
				url:        fmt.Sprintf("%s/set_primary_phonenumber", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_phonenumber_and_otp_missing",
			args: args{
				url:        fmt.Sprintf("%s/set_primary_phonenumber", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if tt.name == "valid:_successfully_set_a_primary_phonenumber" {

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254799774466"
					return &phone, nil
				}

				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254799774466"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
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

			if tt.name == "invalid:_fail_to_set_a_primary_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254799774466"
					return &phone, nil
				}

				fakeEngagementSvs.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254799774466"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return fmt.Errorf("failed to set a primary phone number")
				}
			}

			if tt.name == "invalid:_phonenumber_and_otp_missing" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number provided")
				}
			}

			response := httptest.NewRecorder()

			svr := h.SetPrimaryPhoneNumber()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
}

// TODO: restore
// func TestHandlersInterfacesImpl_RegisterPushToken(t *testing.T) {

// 	i, err := InitializeFakeOnboardingInteractor()
// 	if err != nil {
// 		t.Errorf("failed to initialize onboarding interactor: %v", err)
// 		return
// 	}

// 	h := rest.NewHandlersInterfaces(i)
// 	uid := "5cf354a2-1d3e-400d-8716-7e2aead29f2c"
// 	token := "10c17f3b-a9a9-431c-ad0a-94c684eccd85"
// 	payload := composePushTokenPayload(t, token, uid)

// 	token1 := ""
// 	uid1 := "5cf354a2-1d3e-400d-8716-7e2aead29f2c"
// 	invalidPayload := composePushTokenPayload(t, token1, uid1)

// 	type args struct {
// 		url        string
// 		httpMethod string
// 		body       io.Reader
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		want       http.HandlerFunc
// 		wantStatus int
// 		wantErr    bool
// 	}{
// 		{
// 			name: "valid:_successfully_register_push_token",
// 			args: args{
// 				url:        fmt.Sprintf("%s/register_push_token", serverUrl),
// 				httpMethod: http.MethodPost,
// 				body:       payload,
// 			},
// 			wantStatus: http.StatusOK,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "invalid:_unsuccessfully_register_push_token",
// 			args: args{
// 				url:        fmt.Sprintf("%s/register_push_token", serverUrl),
// 				httpMethod: http.MethodPost,
// 				body:       invalidPayload,
// 			},
// 			wantStatus: http.StatusBadRequest,
// 			wantErr:    true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
// 			if err != nil {
// 				t.Errorf("can't create new request: %v", err)
// 				return
// 			}

// 			if tt.name == "valid:_successfully_register_push_token" {
// 				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
// 					return "400d-8716-7e2aead29f2c", nil
// 				}

// 				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
// 					return &profileutils.UserProfile{
// 						ID:         "f4f39af7--91bd-42b3af315a4e",
// 						PushTokens: []string{"token12", "token23", "token34"},
// 					}, nil
// 				}

// 				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
// 					return nil
// 				}
// 			}
// 			if tt.name == "invalid:_unsuccessfully_register_push_token" {
// 				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
// 					return "400d-8716-7e2aead29f2c", nil
// 				}

// 				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
// 					return &profileutils.UserProfile{
// 						ID:         "f4f39af7--91bd-42b3af315a4e",
// 						PushTokens: []string{"token12", "token23", "token34"},
// 					}, nil
// 				}

// 				fakeRepo.UpdatePushTokensFn = func(ctx context.Context, id string, pushToken []string) error {
// 					return fmt.Errorf("failed to register push tokens")
// 				}
// 			}
// 			response := httptest.NewRecorder()

// 			svr := h.RegisterPushToken()
// 			svr.ServeHTTP(response, req)

// 			if tt.wantStatus != response.Code {
// 				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
// 				return
// 			}

// 			dataResponse, err := ioutil.ReadAll(response.Body)
// 			if err != nil {
// 				t.Errorf("can't read response body: %v", err)
// 				return
// 			}
// 			if dataResponse == nil {
// 				t.Errorf("nil response body data")
// 				return
// 			}
// 		})
// 	}
// }

func TestHandlersInterfacesImpl_AddAdminPermsToUser(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	primaryPhone := "+254711445566"
	validPayload := composeValidPhonePayload(t, primaryPhone)
	validPayload1 := composeValidPhonePayload(t, "+254777882200")
	validPayload2 := composeValidPhonePayload(t, "+")

	invalidPayload := composeValidPhonePayload(t, " ")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_Successfully_update_user_permissions",
			args: args{
				url:        fmt.Sprintf("%s/update_user_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_update_user_permissions",
			args: args{
				url:        fmt.Sprintf("%s/update_user_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		{
			name: "invalid:_empty_phonenumber",
			args: args{
				url:        fmt.Sprintf("%s/update_user_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_unable_to_check_if_phone_exists",
			args: args{
				url:        fmt.Sprintf("%s/update_user_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:_Successfully_update_user_permissions" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []profileutils.PermissionType) error {
					return nil
				}

			}

			if tt.name == "invalid:_update_user_permissions" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []profileutils.PermissionType) error {
					return fmt.Errorf("unable to update user permissions")
				}

			}

			if tt.name == "invalid:_empty_phonenumber" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number")
				}
			}

			if tt.name == "invalid:_unable_to_check_if_phone_exists" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("the phone does not exist")
				}

			}

			svr := h.AddAdminPermsToUser()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_RemoveAdminPermsToUser(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	primaryPhone := "+254711445566"
	validPayload := composeValidPhonePayload(t, primaryPhone)
	validPayload1 := composeValidPhonePayload(t, "+254777882200")
	validPayload2 := composeValidPhonePayload(t, "+")

	invalidPayload := composeValidPhonePayload(t, " ")

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:_Successfully_remove_admin_permissions",
			args: args{
				url:        fmt.Sprintf("%s/remove_admin_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:_remove_admin_permissions",
			args: args{
				url:        fmt.Sprintf("%s/remove_admin_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},

		{
			name: "invalid:_empty_phonenumber",
			args: args{
				url:        fmt.Sprintf("%s/remove_admin_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:_unable_to_check_if_phone_exists",
			args: args{
				url:        fmt.Sprintf("%s/remove_admin_permissions", serverUrl),
				httpMethod: http.MethodPost,
				body:       validPayload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:_Successfully_remove_admin_permissions" {

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []profileutils.PermissionType) error {
					return nil
				}

			}

			if tt.name == "invalid:_remove_admin_permissions" {

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []profileutils.PermissionType) error {
					return fmt.Errorf("unable to update user permissions")
				}

			}

			if tt.name == "invalid:_empty_phonenumber" {

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number")
				}
			}

			if tt.name == "invalid:_unable_to_check_if_phone_exists" {

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}

				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("the phone does not exist")
				}

			}

			svr := h.RemoveAdminPermsToUser()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_AddRoleToUser(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	validPhone := "+254711445566"
	invalidPhone := "+254777882200"
	validRole := profileutils.RoleTypeEmployee
	var invalidRole profileutils.RoleType = "STANGER"
	payload := composeValidRolePayload(t, validPhone, validRole)
	payload1 := composeValidRolePayload(t, invalidPhone, validRole)
	payload2 := composeValidRolePayload(t, validPhone, invalidRole)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:successful_added_user_role",
			args: args{
				url:        fmt.Sprintf("%s/add_user_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:failed_to_find_user",
			args: args{
				url:        fmt.Sprintf("%s/add_user_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid:failed_invalid_role",
			args: args{
				url:        fmt.Sprintf("%s/add_user_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload2,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:successful_added_user_role" {

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}
				fakeRepo.UpdateRoleFn = func(ctx context.Context, id string, role profileutils.RoleType) error {
					return nil
				}
			}

			if tt.name == "invalid:failed_to_find_user" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("Invalid phone number provided")
				}
			}

			if tt.name == "invalid:failed_invalid_role" {

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}
				fakeRepo.UpdateRoleFn = func(ctx context.Context, id string, role profileutils.RoleType) error {
					return fmt.Errorf("Invalid role provided")
				}
			}

			serverResponse := h.AddRoleToUser()
			serverResponse.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlersInterfacesImpl_RemoveRoleToUser(t *testing.T) {

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	validPhone := "+254711445566"
	payload := composeValidPhonePayload(t, validPhone)
	payload1 := composeValidPhonePayload(t, validPhone)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid:successful_removed_user_role",
			args: args{
				url:        fmt.Sprintf("%s/remove_user_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid:failed_to_find_user",
			args: args{
				url:        fmt.Sprintf("%s/remove_user_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload1,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			response := httptest.NewRecorder()

			if tt.name == "valid:successful_removed_user_role" {

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}
				fakeRepo.UpdateRoleFn = func(ctx context.Context, id string, role profileutils.RoleType) error {
					return nil
				}
			}

			if tt.name == "invalid:failed_to_find_user" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("Invalid phone number provided")
				}
			}

			serverResponse := h.RemoveRoleToUser()
			serverResponse.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}

			dataResponse, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response body data")
				return
			}

		})
	}
}

func TestHandlers_PollServices(t *testing.T) {
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}
	h := rest.NewHandlersInterfaces(i)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "return handler func",
			args: args{
				url:        fmt.Sprintf("%s/poll_services", serverUrl),
				httpMethod: http.MethodGet,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}
			// We create a ResponseRecorder to record the response.
			response := httptest.NewRecorder()

			/*
				Some Mockery
			*/

			// call its ServeHTTP method and pass in our Request and ResponseRecorder.
			svr := h.PollServices()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}
		})
	}
}

func TestHandlers_CheckPermission(t *testing.T) {
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}
	h := rest.NewHandlersInterfaces(i)

	uid := uuid.NewString()
	permission := profileutils.Permission{
		Scope:       "test.run",
		Description: "Running tests",
	}

	emptyUIDPayload, err := json.Marshal(dto.CheckPermissionPayload{Permission: &permission})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	emptyPermissionPayload, err := json.Marshal(dto.CheckPermissionPayload{UID: &uid})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	validPayload, err := json.Marshal(dto.CheckPermissionPayload{UID: &uid, Permission: &permission})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "fail: empty uid in payload",
			args: args{
				url:        fmt.Sprintf("%s/check_permission", serverUrl),
				httpMethod: http.MethodGet,
				body:       bytes.NewBuffer(emptyUIDPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: empty permission in payload",
			args: args{
				url:        fmt.Sprintf("%s/check_permission", serverUrl),
				httpMethod: http.MethodGet,
				body:       bytes.NewBuffer(emptyPermissionPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: usecase error checking permission",
			args: args{
				url:        fmt.Sprintf("%s/check_permission", serverUrl),
				httpMethod: http.MethodGet,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "success: user is not authorized",
			args: args{
				url:        fmt.Sprintf("%s/check_permission", serverUrl),
				httpMethod: http.MethodGet,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
		{
			name: "success: user is authorized",
			args: args{
				url:        fmt.Sprintf("%s/check_permission", serverUrl),
				httpMethod: http.MethodGet,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: usecase error checking permission" {
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("cannot check for permission")
				}
			}

			if tt.name == "success: user is not authorized" {
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "success: user is authorized" {
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
			}
			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			// We create a ResponseRecorder to record the response.
			response := httptest.NewRecorder()

			// call its ServeHTTP method and pass in our Request and ResponseRecorder.
			svr := h.CheckHasPermission()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}
		})
	}
}

func TestHandlers_CreateRole(t *testing.T) {
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}
	h := rest.NewHandlersInterfaces(i)

	invalidPayload, err := json.Marshal(dto.RoleInput{Description: "Test Role"})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	validPayload, err := json.Marshal(dto.RoleInput{Name: "Test Role", Description: "Test Role"})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "fail: missing name in request payload",
			args: args{
				url:        fmt.Sprintf("%s/create_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: create unauthorized role error",
			args: args{
				url:        fmt.Sprintf("%s/create_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "success: create role",
			args: args{
				url:        fmt.Sprintf("%s/create_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fail: create unauthorized role error" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, profileID string, role dto.RoleInput) (*profileutils.Role, error) {
					return nil, fmt.Errorf("duplicate role")
				}
			}

			if tt.name == "success: create role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, profileID string, role dto.RoleInput) (*profileutils.Role, error) {
					return &profileutils.Role{
						Scopes: []string{"role.edit"},
					}, nil
				}
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}
			// We create a ResponseRecorder to record the response.
			response := httptest.NewRecorder()

			// call its ServeHTTP method and pass in our Request and ResponseRecorder.
			svr := h.CreateRole()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}
		})
	}
}

func TestHandlers_AssignRole(t *testing.T) {
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}
	h := rest.NewHandlersInterfaces(i)

	emptyRoleIDPayload, err := json.Marshal(dto.AssignRolePayload{UserID: "123456", RoleID: ""})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	emptyUserIDPayload, err := json.Marshal(dto.AssignRolePayload{RoleID: "12345", UserID: ""})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	validPayload, err := json.Marshal(dto.AssignRolePayload{UserID: "123456", RoleID: "123456"})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "fail: missing role ID in request payload",
			args: args{
				url:        fmt.Sprintf("%s/assign_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(emptyRoleIDPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: missing user ID in request payload",
			args: args{
				url:        fmt.Sprintf("%s/assign_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(emptyUserIDPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: assign role error",
			args: args{
				url:        fmt.Sprintf("%s/assign_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "success: assign role",
			args: args{
				url:        fmt.Sprintf("%s/assign_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fail: assign role error" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: ""}, nil
				}

				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}

				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "",
						Scopes: []string{profileutils.CanAssignRole.Scope},
					}, nil
				}

				fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: ""}, nil
				}

				fakeRepo.UpdateUserRoleIDsFn = func(ctx context.Context, id string, roleIDs []string) error {
					return fmt.Errorf("cannot update uids")
				}
			}

			if tt.name == "success: assign role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: ""}, nil
				}

				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}

				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "",
						Scopes: []string{profileutils.CanAssignRole.Scope},
					}, nil
				}

				fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: ""}, nil
				}

				fakeRepo.UpdateUserRoleIDsFn = func(ctx context.Context, id string, roleIDs []string) error {
					return nil
				}
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}
			// We create a ResponseRecorder to record the response.
			response := httptest.NewRecorder()

			// call its ServeHTTP method and pass in our Request and ResponseRecorder.
			svr := h.AssignRole()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}
		})
	}
}

func TestHandlers_RemoveRoleByName(t *testing.T) {
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}
	h := rest.NewHandlersInterfaces(i)

	invalidPayload, err := json.Marshal(dto.DeleteRolePayload{Name: ""})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	validPayload, err := json.Marshal(dto.DeleteRolePayload{Name: "Test Role"})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "fail: missing name in request payload",
			args: args{
				url:        fmt.Sprintf("%s/remove_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "fail: find role by name error",
			args: args{
				url:        fmt.Sprintf("%s/remove_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "fail: delete role error",
			args: args{
				url:        fmt.Sprintf("%s/remove_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "success: remove role",
			args: args{
				url:        fmt.Sprintf("%s/remove_role", serverUrl),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: find role by name error" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByNameFn = func(ctx context.Context, roleName string) (*profileutils.Role, error) {
					return nil, fmt.Errorf("cannot find role")
				}

			}

			if tt.name == "fail: delete role error" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByNameFn = func(ctx context.Context, roleName string) (*profileutils.Role, error) {
					return &profileutils.Role{ID: "", Scopes: []string{profileutils.CanAssignRole.Scope}}, nil
				}

				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{ID: uuid.NewString(), Name: "Happy Test Role", Scopes: []string{profileutils.CanAssignRole.Scope}}, nil
				}

				fakeRepo.DeleteRoleFn = func(ctx context.Context, roleID string) (bool, error) {
					return false, fmt.Errorf("cannot delete role")
				}
			}

			if tt.name == "success: remove role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByNameFn = func(ctx context.Context, roleName string) (*profileutils.Role, error) {
					return &profileutils.Role{ID: "", Scopes: []string{profileutils.CanAssignRole.Scope}}, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{ID: uuid.NewString(), Name: "Happy Test Role", Scopes: []string{profileutils.CanAssignRole.Scope}}, nil
				}
				fakeRepo.DeleteRoleFn = func(ctx context.Context, roleID string) (bool, error) {
					return true, nil
				}
			}

			// Create a request to pass to our handler.
			req, err := http.NewRequest(tt.args.httpMethod, tt.args.url, tt.args.body)
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}
			// We create a ResponseRecorder to record the response.
			response := httptest.NewRecorder()

			// call its ServeHTTP method and pass in our Request and ResponseRecorder.
			svr := h.RemoveRoleByName()
			svr.ServeHTTP(response, req)

			if tt.wantStatus != response.Code {
				t.Errorf("expected status %d, got %d", tt.wantStatus, response.Code)
				return
			}
		})
	}
}
