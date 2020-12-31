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

func composeWrongUserPINPayload(t *testing.T) *domain.LoginPayload {
	// TODO: Harmonize LoginPayload SignUpPayload
	phone := base.TestUserPhoneNumberWithPin
	pin := "qwer"
	flavour := base.FlavourPro
	payload := &domain.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeWrongUserPhonePayload(t *testing.T) *domain.LoginPayload {
	phone := "+254700000000"
	pin := base.TestUserPin
	flavour := base.FlavourPro
	payload := &domain.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeInvalidUserPhonePayload(t *testing.T) *domain.LoginPayload {
	phone := "+254-not-a-number"
	pin := base.TestUserPin
	flavour := base.FlavourPro
	payload := &domain.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}
func TestLoginInByPhone(t *testing.T) {
	_, err := CreateTestUserByPhone(t)
	if err != nil {
		t.Errorf("failed to create a user by phone")
		return
	}

	client := http.DefaultClient
	validPayload := composeValidUserPayload(t)
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// invalid payload
	badPayload := composeInValidUserPayload(t)
	bs2, err := json.Marshal(badPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(bs2)

	wrongPINPayload := composeWrongUserPINPayload(t)
	wrongPINBs, err := json.Marshal(wrongPINPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	badPINpayload := bytes.NewBuffer(wrongPINBs)

	wrongPhonePayload := composeWrongUserPhonePayload(t)
	wrongPhoneBs, err := json.Marshal(wrongPhonePayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	badPhonepayload := bytes.NewBuffer(wrongPhoneBs)

	invalidPhonePayload := composeInvalidUserPhonePayload(t)
	invalidPhoneBs, err := json.Marshal(invalidPhonePayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	badInvalidPhonepayload := bytes.NewBuffer(invalidPhoneBs)

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
			name: "success: login user with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "failure: login user with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: login user with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: login user with a wrong PIN",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       badPINpayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: login user with a wrong primary phone number",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       badPhonepayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: login user with invalid phone number",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       badInvalidPhonepayload,
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
