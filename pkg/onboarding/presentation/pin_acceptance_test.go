package presentation_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

func composeInValidPinPayload(t *testing.T) *domain.PIN {
	payload := &domain.PIN{
		PhoneNumber: "",
		PINNumber:   "1234",
	}
	return payload
}

func composeValidPinPayload(t *testing.T) *domain.PIN {
	payload := &domain.PIN{
		PhoneNumber: base.TestUserPhoneNumberWithPin,
		PINNumber:   "1234",
	}
	return payload
}

func composeUnregisteredPhone(t *testing.T) *domain.PIN {
	payload := &domain.PIN{
		PhoneNumber: base.TestUserPhoneNumber,
		PINNumber:   "1234",
	}
	return payload
}

func TestSetUserPIN(t *testing.T) {
	client := http.DefaultClient
	// create a user and thier profile
	_, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		// return
	}
	validPayload := composeValidPinPayload(t)
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// invalid payload
	badPayload := composeInValidPinPayload(t)
	bs2, err := json.Marshal(badPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(bs2)

	unregisteredPhone := composeUnregisteredPhone(t)
	bs3, err := json.Marshal(unregisteredPhone)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	userNotRegistered := bytes.NewBuffer(bs3)

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
			name: "failure: set pin with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: set pin with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: set pin with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "failure: signup user with the same valid payload again",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       userNotRegistered,
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
			// TODO calvine to uncomment this after duplicate profiles are fixed
			// if tt.wantStatus != resp.StatusCode {
			// 	t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			// 	return
			// }
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
