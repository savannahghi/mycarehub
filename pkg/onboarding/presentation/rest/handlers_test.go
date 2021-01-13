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

	profile := usecases.NewProfileUseCase(r)
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

func TestHandlersInterfacesImpl_VerifySignUpPhoneNumber(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize onboarding interactor: %v", err)
		return
	}

	h := rest.NewHandlersInterfaces(i)
	// payload 1
	phoneNumber := struct {
		PhoneNumber string
	}{
		PhoneNumber: base.TestUserPhoneNumber,
	}
	bs, err := json.Marshal(phoneNumber)
	if err != nil {
		t.Errorf("unable to marshal phoneNumber to JSON: %s", err)
		return
	}
	payload := bytes.NewBuffer(bs)

	// payload 2
	phoneNumber2 := struct {
		PhoneNumber string
	}{
		PhoneNumber: base.TestUserPhoneNumberWithPin,
	}
	bs, err = json.Marshal(phoneNumber2)
	if err != nil {
		t.Errorf("unable to marshal phoneNumber2 to JSON: %s", err)
		return
	}
	payload2 := bytes.NewBuffer(bs)

	// payload 3
	phoneNumber3 := struct {
		PhoneNumber string
	}{
		PhoneNumber: "0700100200",
	}
	bs, err = json.Marshal(phoneNumber3)
	if err != nil {
		t.Errorf("unable to marshal phoneNumber3 to JSON: %s", err)
		return
	}
	payload3 := bytes.NewBuffer(bs)

	// payload 4
	phoneNumber4 := struct {
		PhoneNumber string
	}{
		PhoneNumber: "0700600300",
	}
	bs, err = json.Marshal(phoneNumber4)
	if err != nil {
		t.Errorf("unable to marshal phoneNumber4 to JSON: %s", err)
		return
	}
	payload4 := bytes.NewBuffer(bs)

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
			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
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
