package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func composeValidWARetryOTPPayload(t *testing.T) *dto.SendRetryOTPPayload {
	phoneNumber := interserviceclient.TestUserPhoneNumber
	step := 1
	return &dto.SendRetryOTPPayload{
		Phone:     &phoneNumber,
		RetryStep: &step,
	}
}

// func composeValidMessageRetryOTPPayload(t *testing.T) *dto.SendRetryOTPPayload {
// 	phoneNumber := interserviceclient.TestUserPhoneNumber
// 	step := 2
// 	return &dto.SendRetryOTPPayload{
// 		Phone:     &phoneNumber,
// 		RetryStep: &step,
// 	}
// }

func composeInvalidRetryOTPPayload(t *testing.T) *dto.SendRetryOTPPayload {
	phoneNumber := interserviceclient.TestUserPhoneNumber
	step := 3 // There is no retry step > 2
	return &dto.SendRetryOTPPayload{
		Phone:     &phoneNumber,
		RetryStep: &step,
	}
}

func TestSendRetryOTP(t *testing.T) {
	client := http.DefaultClient
	phoneNumber := interserviceclient.TestUserPhoneNumber
	user, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to create a user by phone %v", err)
		return
	}
	if user == nil {
		t.Errorf("nil user found")
		return
	}
	validWAPayload := composeValidWARetryOTPPayload(t)
	bs, err := json.Marshal(validWAPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	validWAReqPayload := bytes.NewBuffer(bs)

	// validMessagePayload := composeValidMessageRetryOTPPayload(t)
	// validMessagePayloadBs, err := json.Marshal(validMessagePayload)
	// if err != nil {
	// 	t.Errorf("unable to marshal test item to JSON: %s", err)
	// }
	// validMessageReqPayload := bytes.NewBuffer(validMessagePayloadBs)

	badPayload := composeInvalidRetryOTPPayload(t)
	bs2, err := json.Marshal(badPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(bs2)

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
			name: "success: send retry OTP via WhatsApp(WA)",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       validWAReqPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		// {
		// 	name: "success: send retry OTP via Twilio messages",
		// 	args: args{
		// 		url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
		// 		httpMethod: http.MethodPost,
		// 		body:       validMessageReqPayload,
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    true,
		// },
		{
			name: "failure: send retry OTP with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: send retry OTP with wrong retry step",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
				return
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if data == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestSendOTP(t *testing.T) {
	client := http.DefaultClient

	phoneNumber := interserviceclient.TestUserPhoneNumber
	validPayload := &dto.PhoneNumberPayload{
		PhoneNumber: &phoneNumber,
	}
	bs1, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	validReqPayload := bytes.NewBuffer(bs1)

	invalidPhone := "+2547569"
	invalidPayload := &dto.PhoneNumberPayload{
		PhoneNumber: &invalidPhone,
	}
	bs2, err := json.Marshal(invalidPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidReqPayload := bytes.NewBuffer(bs2)

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
			name: "valid: should send sms",
			args: args{
				url:        fmt.Sprintf("%s/send_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       validReqPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: wrong phone number provided",
			args: args{
				url:        fmt.Sprintf("%s/send_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidReqPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if tt.wantStatus != resp.StatusCode {
				dump, _ := httputil.DumpResponse(resp, true)
				t.Errorf("expected status %d, got %d with response %v", tt.wantStatus, resp.StatusCode, string(dump))
				return
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}
			if data == nil {
				t.Errorf("nil response body data")
				return
			}
		})
	}

}
