package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

func TestLoginInClientByPhone(t *testing.T) {
	client := http.DefaultClient
	validPayload, err := composeValidClientUserPayload(t, testPhone)
	if err != nil {
		t.Errorf("Error composing client user payload: %v", err)
	}

	payload, err := convertToBytesBuffer(validPayload)
	if err != nil {
		t.Errorf("Error converting to bytes buffer: %v", err)
	}

	wrongPin := "11111"
	wrongPinPayloadInput := dto.LoginPayload{
		PhoneNumber: &testPhone,
		PIN:         &wrongPin,
		Flavour:     testClientFlavour,
	}

	wrongPinPayload, err := convertToBytesBuffer(wrongPinPayloadInput)
	if err != nil {
		t.Errorf("Error converting to bytes buffer: %v", err)
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
			name: "invalid: wrong pin passed",
			args: args{
				url:        fmt.Sprintf("%s/login_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       wrongPinPayload,
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

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "myafyahub") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}
			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			fmt.Println("got data: ", data)
			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
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

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}
