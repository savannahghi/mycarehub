package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func composeInvalidUserPINPayload(t *testing.T) *resources.LoginPayload {
	phone := base.TestUserPhoneNumberWithPin
	pin := "" // empty pin
	flavour := base.FlavourPro
	payload := &resources.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeWrongUserPINPayload(t *testing.T) *resources.LoginPayload {
	phone := base.TestUserPhoneNumberWithPin
	pin := "qwer"
	flavour := base.FlavourPro
	payload := &resources.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeWrongUserPhonePayload(t *testing.T) *resources.LoginPayload {
	phone := "+254700000000"
	pin := base.TestUserPin
	flavour := base.FlavourPro
	payload := &resources.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeInvalidUserPhonePayload(t *testing.T) *resources.LoginPayload {
	phone := "+254-not-a-number"
	pin := base.TestUserPin
	flavour := base.FlavourPro
	payload := &resources.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeWrongFlavourPayload(t *testing.T) *resources.LoginPayload {
	phone := base.TestUserPhoneNumberWithPin
	pin := base.TestUserPin
	payload := &resources.LoginPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     "bad-flavour-supplied",
	}
	return payload
}

func TestLoginInByPhone(t *testing.T) {
	phoneNumber := base.TestUserPhoneNumber
	user, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to create a user by phone %v", err)
		return
	}
	if user == nil {
		t.Errorf("nil user found")
		return
	}

	client := http.DefaultClient
	validPayload, err := composeValidUserPayload(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to compose a valid payload")
		return
	}
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// invalid payload
	badPayload := composeInvalidUserPINPayload(t)
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

	emptyData := &resources.LoginPayload{}
	emptyBs, err := json.Marshal(emptyData)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	emptyPayload := bytes.NewBuffer(emptyBs)

	invalidFlavourPayload := composeWrongFlavourPayload(t)
	invalidFlavourBs, err := json.Marshal(invalidFlavourPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	badFlavourPayload := bytes.NewBuffer(invalidFlavourBs)

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
			wantErr:    false,
		},
		{
			name: "failure: login user with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       emptyPayload,
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
		{
			name: "failure: login user with invalid flavour",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       badFlavourPayload,
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

			for k, v := range base.GetDefaultHeaders(t, baseURL, "onboarding") {
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

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			// TODO ! uncomment/ remove after error message format has been standerdized
			// TODO! assert some data
			// if tt.wantErr {
			// 	errMsg, ok := data["error"]
			// 	if !ok {
			// 		t.Errorf("Request error: %s", errMsg)
			// 		return
			// 	}
			// }

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

		})
	}
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestLoginAsAnonymous(t *testing.T) {
	client := http.DefaultClient

	p1, err := json.Marshal(&resources.LoginPayload{
		Flavour: base.FlavourConsumer,
	})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
	}
	validPayload := bytes.NewBuffer(p1)

	p2, err := json.Marshal(&resources.LoginPayload{
		Flavour: base.FlavourPro,
	})
	if err != nil {
		t.Errorf("unable to marshal payload to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(p2)

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
			name: "valid : correct flavour",
			args: args{
				url:        fmt.Sprintf("%s/login_anonymous", baseURL),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "valid : incorrect flavour",
			args: args{
				url:        fmt.Sprintf("%s/login_anonymous", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
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

			for k, v := range base.GetDefaultHeaders(t, baseURL, "onboarding") {
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

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("Request error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

		})
	}

}

func TestRefreshToken(t *testing.T) {
	client := http.DefaultClient
	//todo(dexter) restore this after 11th release
	// validToken := "AOvuKvSiBjrtQ6WRdTbRUFeGm4q6KbKg1kdwACot-zZFSqAwZtePlLKTT4U5Ew7C6UFcQsu6HQPAKD-1Hr_jTrtUtwTJ2mrqTBEW0oxtWImbB7fnPtNnl3mSBMpnVewbj14w_quNw_AkvBaQKu2vIR5tjATqYaPHCRMM1d-W7GMQUneKlJNz-JQ"
	// validPayload := &resources.RefreshTokenPayload{
	// 	RefreshToken: &validToken,
	// }
	// bs, err := json.Marshal(validPayload)
	// if err != nil {
	// 	t.Errorf("unable to marshal test item to JSON: %s", err)
	// }
	// payload := bytes.NewBuffer(bs)

	inValidToken := "some-token"
	inValidPayload := &resources.RefreshTokenPayload{
		RefreshToken: &inValidToken,
	}
	badBs, err := json.Marshal(inValidPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	badPayload := bytes.NewBuffer(badBs)

	emptyData := &resources.LoginPayload{}
	emptyBs, err := json.Marshal(emptyData)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	emptyPayload := bytes.NewBuffer(emptyBs)

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
		//todo(dexter) restore this after 11th release
		// {
		// 	name: "success: refresh a token",
		// 	args: args{
		// 		url:        fmt.Sprintf("%s/refresh_token", baseURL),
		// 		httpMethod: http.MethodPost,
		// 		body:       payload,
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "failure: refresh token with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/refresh_token", baseURL),
				httpMethod: http.MethodPost,
				body:       emptyPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: refresh token with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/refresh_token", baseURL),
				httpMethod: http.MethodPost,
				body:       badPayload,
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

			for k, v := range base.GetDefaultHeaders(t, baseURL, "onboarding") {
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

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			// TODO ! uncomment/ remove after error message format has been standerdized
			// TODO! assert some data
			// if tt.wantErr {
			// 	errMsg, ok := data["error"]
			// 	if !ok {
			// 		t.Errorf("Request error: %s", errMsg)
			// 		return
			// 	}
			// }

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}
		})
	}
}
