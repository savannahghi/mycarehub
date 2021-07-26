package tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

func TestCreateSMSData_acceptance(t *testing.T) {
	client := http.DefaultClient

	validLinkID := uuid.New().String()
	text := "Test Covers"
	to := "3601"
	id := "60119"
	from := "+254705385894"
	date := "2021-05-17T13:20:04.490Z"

	// valid payload
	validSMSData := &dto.AfricasTalkingMessage{
		LinkID: validLinkID,
		Text:   text,
		To:     to,
		ID:     id,
		Date:   date,
		From:   from,
	}

	validPayload := composeSMSMessageDataPayload(t, validSMSData)

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
			name: "Success: receive sms message data with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/incoming_ait_messages", baseURL),
				httpMethod: http.MethodPost,
				body:       validPayload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "Failure: receive sms with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/incoming_ait_messages", baseURL),
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

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "onboarding") {
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
