package graph_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestMSISDNLogin(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	msisdnLoginURL := fmt.Sprintf("%s/%s", baseURL, "msisdn_login")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get request headers %v", err)
		return
	}

	type args struct {
		PhoneNumber string
		Pin         string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "correct login credentials",
			args: args{
				PhoneNumber: base.TestUserPhoneNumberWithPin,
				Pin:         base.TestUserPin,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "totally incorrect login credentials",
			args: args{
				PhoneNumber: "not a real phone number",
				Pin:         "not a pin",
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "invalid login credentials format",
			args: args{
				PhoneNumber: "not a real phone number",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "wrong pin credentials",
			args: args{
				PhoneNumber: base.TestUserPhoneNumberWithPin,
				Pin:         "wrong pin number",
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
		{
			name: "non-existent login credentials",
			args: args{
				PhoneNumber: "+254780654321",
				Pin:         "0000",
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "no pin login credentials",
			args: args{
				PhoneNumber: "+254711223344",
				Pin:         "has no pin",
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestInput := map[string]interface{}{}
			requestInput["phonenumber"] = tt.args.PhoneNumber
			requestInput["pin"] = tt.args.Pin

			body, err := mapToJSONReader(requestInput)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				msisdnLoginURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("statusCode = %v, wantStatus %v", resp.StatusCode, tt.wantStatus)
				return
			}

		})
	}
}
