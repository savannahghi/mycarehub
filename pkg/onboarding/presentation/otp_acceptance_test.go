package presentation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func composeValidWARetryOTPPayload(t *testing.T) *domain.SendRetryOTPPayload {
	phoneNumber := base.TestUserPhoneNumber
	step := 1
	return &domain.SendRetryOTPPayload{
		Phone:     &phoneNumber,
		RetryStep: &step,
	}
}

func composeValidMessageRetryOTPPayload(t *testing.T) *domain.SendRetryOTPPayload {
	phoneNumber := base.TestUserPhoneNumber
	step := 2
	return &domain.SendRetryOTPPayload{
		Phone:     &phoneNumber,
		RetryStep: &step,
	}
}

func composeInvalidRetryOTPPayload(t *testing.T) *domain.SendRetryOTPPayload {
	phoneNumber := base.TestUserPhoneNumber
	step := 3 // There is no retry step > 2
	return &domain.SendRetryOTPPayload{
		Phone:     &phoneNumber,
		RetryStep: &step,
	}
}

func TestSendRetryOTP(t *testing.T) {
	client := http.DefaultClient
	validWAPayload := composeValidWARetryOTPPayload(t)
	bs, err := json.Marshal(validWAPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	validWAReqPayload := bytes.NewBuffer(bs)

	validMessagePayload := composeValidMessageRetryOTPPayload(t)
	validMessagePayloadBs, err := json.Marshal(validMessagePayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	validMessageReqPayload := bytes.NewBuffer(validMessagePayloadBs)

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
		{
			name: "success: send retry OTP via Twilio messages",
			args: args{
				url:        fmt.Sprintf("%s/send_retry_otp", baseURL),
				httpMethod: http.MethodPost,
				body:       validMessageReqPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
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

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
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
}
