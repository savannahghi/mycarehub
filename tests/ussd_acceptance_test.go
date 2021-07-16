package tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestGenerateUSSD_acceptance_tests(t *testing.T) {
	client := http.DefaultClient

	level := 0
	phoneNumber := interserviceclient.TestUserPhoneNumber
	sessionID := "123456"

	// valid payload
	payload := &dto.SessionDetails{
		SessionID:   sessionID,
		PhoneNumber: &phoneNumber,
		Level:       level,
		Text:        "",
	}

	payload2 := &dto.SessionDetails{
		SessionID:   "",
		PhoneNumber: &phoneNumber,
		Level:       level,
		Text:        "Test Text",
	}

	validPayload := composeUSSDPayload(t, payload)
	emptySessionIDPayload := composeUSSDPayload(t, payload2)

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
			name: "Success: receive ussd data with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Fail: receive ussd data with an empty sessionID",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       emptySessionIDPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Fail: receive ussd data with a nil payload",
			args: args{
				url:        fmt.Sprintf("%s/ait_ussd", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
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
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			if err != nil {
				t.Errorf("can't create a new request: %v", err)
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
