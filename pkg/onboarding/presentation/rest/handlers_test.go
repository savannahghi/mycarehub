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

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	chargemasterMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	engagementMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	erpMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	mailgunMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	messagingMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	otpMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/rest"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	mockRepo "gitlab.slade360emr.com/go/profile/pkg/onboarding/repository/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

var fakeRepo mockRepo.FakeOnboardingRepository
var fakeOtp otpMock.FakeServiceOTP
var serverUrl = "http://localhost:5000"

func InitializeFakeOnboaridingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var otpSvc otp.ServiceOTP = &fakeOtp
	var erpSvc erp.ServiceERP = &erpMock.FakeServiceERP{}
	var chargemasterSvc chargemaster.ServiceChargeMaster = &chargemasterMock.FakeServiceChargeMaster{}
	var engagementSvc engagement.ServiceEngagement = &engagementMock.FakeServiceEngagement{}
	var mailgunSvc mailgun.ServiceMailgun = &mailgunMock.FakeServiceMailgun{}
	var messagingSvc messaging.ServiceMessaging = &messagingMock.FakeServiceMessaging{}

	profile := usecases.NewProfileUseCase(r, otpSvc)
	login := usecases.NewLoginUseCases(r, profile)
	survey := usecases.NewSurveyUseCases(r)
	supplier := usecases.NewSupplierUseCases(
		r, profile, erpSvc, chargemasterSvc, engagementSvc, mailgunSvc, messagingSvc,
	)
	userpin := usecases.NewUserPinUseCase(r, otpSvc, profile)
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, otpSvc)

	i, err := interactor.NewOnboardingInteractor(
		r, profile, su, otpSvc, supplier, login,
		survey, userpin, erpSvc, chargemasterSvc, engagementSvc, mailgunSvc, messagingSvc,
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

func composeSignupPayload(t *testing.T, phone, pin, otp string, flavour base.Flavour) *bytes.Buffer {
	payload := resources.SignUpInput{
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
func TestHandlersInterfacesImpl_VerifySignUpPhoneNumber(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// valid:_successfully_verifies_a_phone_number
	payload := composeValidPhonePayload(t, base.TestUserPhoneNumber)

	// payload 2
	payload2 := composeValidPhonePayload(t, base.TestUserPhoneNumberWithPin)

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
			if tt.name == "valid:_successfully_verifies_a_phone_number" {
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeOtp.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return &base.OtpResponse{OTP: "1234"}, nil
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
				fakeOtp.GenerateAndSendOTPFn = func(ctx context.Context, phone string) (*base.OtpResponse, error) {
					return nil, fmt.Errorf("unable generate and send otp")
				}
			}
			// Our handlers satisfy http.Handler, so we can call its ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.VerifySignUpPhoneNumber(ctx)
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
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)

	// payload
	pin := "2030"
	flavour := base.FlavourPro
	otp := "1234"
	phoneNumber := base.TestUserPhoneNumber
	payload := composeSignupPayload(t, phoneNumber, pin, otp, flavour)

	// payload 2
	pin2 := "1000"
	flavour2 := base.FlavourConsumer
	otp2 := "9000"
	phoneNumber2 := "+254720125456"
	payload2 := composeSignupPayload(t, phoneNumber2, pin2, otp2, flavour2)

	// payload 3
	pin3 := "2000"
	flavour3 := base.FlavourConsumer
	otp3 := "3000"
	phoneNumber3 := "+254721100200"
	payload3 := composeSignupPayload(t, phoneNumber3, pin3, otp3, flavour3)
	// payload 4
	pin4 := "1228"
	flavour4 := base.FlavourConsumer
	otp4 := "9652"
	phoneNumber4 := "+254721410698"
	payload4 := composeSignupPayload(t, phoneNumber4, pin4, otp4, flavour4)
	// payload 5
	pin5 := "0000"
	flavour5 := base.FlavourConsumer
	otp5 := "9520"
	phoneNumber5 := "+254721410589"
	payload5 := composeSignupPayload(t, phoneNumber5, pin5, otp5, flavour5) // payload6
	// payload6
	pin6 := "0000"
	flavour6 := base.FlavourConsumer
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
			name: "invalid_check_phone_exists_returns_error",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload4,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_check_phone_exists_returns_true",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", serverUrl),
				httpMethod: http.MethodPost,
				body:       payload5,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid_create_user_via_their_phone_number_fails",
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
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*resources.CreatedUserResponse, error) {
					return &resources.CreatedUserResponse{
						UID:         "1106f10f-bea6-4fa3-bdba-16b1e39bd318",
						DisplayName: "kalulu juha",
						Email:       "juha@gmail.com",
						PhoneNumber: "0756232452",
						PhotoURL:    "",
						ProviderID:  "google.com",
					}, nil
				}
				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID:           "125",
								LoginProvider: "Phone",
							},
						},
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GenerateAuthCredentialsFn = func(ctx context.Context, phone string) (*base.AuthCredentialResponse, error) {
					return &base.AuthCredentialResponse{
						UID: "5550",
						// IDToken:      "555",
						RefreshToken: "55550",
					}, nil
				}
				// SetUserPINFn =
				fakeRepo.CreateEmptySupplierProfileFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "5550",
						// ProfileID:  "555",
						SupplierID: "5555",
					}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return &base.Customer{
						ID: "0000",
						// ProfileID:  "1230",
						CustomerID: "22222",
					}, nil
				}
				// should return a profile with an ID
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "123",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
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
			}

			// mock VerifyOTP to return an error
			if tt.name == "invalid_unable_to_verify_otp" {
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify otp")
				}
			}

			// mock VerifyOTP to return an false
			if tt.name == "invalid_verify_otp_returns_false" {
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			// mock CheckPhoneExists to return an error
			if tt.name == "invalid_check_phone_exists_returns_error" {
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("unable to check phone")
				}
			}

			// mock CheckPhoneExists to returns true for a number that exists
			// also mocking CheckIfPhoneNumberExists is necessary to reach `CheckPhoneExists`
			if tt.name == "invalid_check_phone_exists_returns_true" {
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return true, nil
				}
			}

			// mock `GetOrCreatePhoneNumberUser` to return an error
			if tt.name == "invalid_check_phone_exists_returns_true" {
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*resources.CreatedUserResponse, error) {
					return nil, fmt.Errorf("unable to create user")
				}
			}

			// Our handlers satisfy http.Handler, so we can call its ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.CreateUserWithPhoneNumber(ctx)
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

func TestHandlersInterfacesImpl_UserRecoveryPhoneNumbers(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload 1
	payload := composeValidPhonePayload(t, base.TestUserPhoneNumber)

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
			name: "valid:_unable_to_get_profile",
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
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}
			}

			// we set GetUserProfileByPhoneNumber to return an error
			if tt.name == "valid:_unable_to_get_profile" {
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to retreive profile")
				}
			}

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			svr := h.UserRecoveryPhoneNumbers(ctx)
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
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload sucessfully_request_pin_reset
	payload := composeValidPhonePayload(t, base.TestUserPhoneNumber)
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
			name: "valid:sucessfully_request_pin_reset",
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
			if tt.name == "valid:sucessfully_request_pin_reset" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
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

			if tt.name == "invalid:_inable_to_get_primary_phone" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to fetch profile")
				}
			}

			if tt.name == "invalid:check_has_pin_failed" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return &base.UserProfile{
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
			svr := h.RequestPINReset(ctx)
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
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
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
			name: "valid:sucessfully_reset_pin",
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
			if tt.name == "valid:sucessfully_reset_pin" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			// we set `UpdatePIN` to return an error
			if tt.name == "invalid:unable_to_update_pin" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
					}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, profileID string) (*domain.PIN, error) {
					return &domain.PIN{ID: "123", ProfileID: "456"}, nil
				}
				fakeOtp.VerifyOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return false, fmt.Errorf("unable to update pin")
				}
			}
			svr := h.ResetPin(ctx)
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
