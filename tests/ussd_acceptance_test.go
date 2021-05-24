package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestCreateUSSDData(t *testing.T) {
	client := http.DefaultClient

	invalidPhoneNumber := "2345"
	phoneNumber := base.TestUserPhoneNumber
	data := url.Values{}
	data.Set("phoneNumber", phoneNumber)
	data.Set("sessionId", "1234567")
	data.Set("input", "1")

	payload := composeUssdPayload(t, phoneNumber, "1234567", "1")
	invalid, err := json.Marshal(dto.EndSessionDetails{
		PhoneNumber: &phoneNumber,
		Input:       "1",
	})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(invalid)
	payloadWithInvalidPhone := composeUssdPayload(t, invalidPhoneNumber, "1234567", "1")

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
			name: "success: create USSD with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_end_note_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: create USSD with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_end_note_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: create USSD with invalid phone number",
			args: args{
				url:        fmt.Sprintf("%s/ait_end_note_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       payloadWithInvalidPhone,
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

			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
			dataResponse, err := ioutil.ReadAll(resp.Body)
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

func TestGenerateUSSD(t *testing.T) {
	client := http.DefaultClient

	invalidPhoneNumber := "2345"
	phoneNumber := base.TestUserPhoneNumber
	data := url.Values{}
	data.Set("phoneNumber", phoneNumber)
	data.Set("sessionId", "1234567")
	data.Set("text", "1")

	payload := strings.NewReader(data.Encode())
	payloadWithInvalidPhone := composeUssdPayload(t, invalidPhoneNumber, "1234567", "1")
	payloadWithInvalidSessionID := composeUssdPayload(t, phoneNumber, "", "1")
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
			name: "success: send USSD with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: send USSD with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: send USSD with invalid phone number",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       payloadWithInvalidPhone,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: send USSD with invalid sessionId",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       payloadWithInvalidSessionID,
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

			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
			dataResponse, err := ioutil.ReadAll(resp.Body)
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
