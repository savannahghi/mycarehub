package graph_test

import (
	"context"
	"fmt"
	"io"
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

	goodReqInput := map[string]interface{}{}
	goodReqInput["phonenumber"] = base.TestUserPhoneNumberWithPin
	goodReqInput["pin"] = base.TestUserPin

	validLoginReq, err := mapToJSONReader(goodReqInput)
	if err != nil {
		t.Errorf("unable to get goodReqInput JSON io Reader: %s", err)
		return
	}

	incorrectReqInput := map[string]interface{}{}
	incorrectReqInput["phonenumber"] = "not a real phone number"
	incorrectReqInput["pin"] = "not a pin"

	invalidLoginReq, err := mapToJSONReader(incorrectReqInput)
	if err != nil {
		t.Errorf("unable to get invalidLoginReq JSON io Reader: %s", err)
		return
	}

	invalidReqFormatInput := map[string]interface{}{}
	invalidReqFormatInput["invalidKey"] = "not a real phone number"

	invalidLoginFormatReq, err := mapToJSONReader(invalidReqFormatInput)
	if err != nil {
		t.Errorf("unable to get invalidLoginFormatReq JSON io Reader: %s", err)
		return
	}

	badReqInput := map[string]interface{}{}
	badReqInput["phonenumber"] = base.TestUserPhoneNumberWithPin
	badReqInput["pin"] = "wrong pin number"

	badLoginReq, err := mapToJSONReader(badReqInput)
	if err != nil {
		t.Errorf("unable to get badLoginReq JSON io Reader: %s", err)
		return
	}

	nonExistentCredReqInput := map[string]interface{}{}
	nonExistentCredReqInput["phonenumber"] = "+254780654321"
	nonExistentCredReqInput["pin"] = "0000"

	nonExistentCredLoginReq, err := mapToJSONReader(nonExistentCredReqInput)
	if err != nil {
		t.Errorf("unable to get nonExistentCredLoginReq JSON io Reader: %s", err)
		return
	}

	noPinCredReqInput := map[string]interface{}{}
	noPinCredReqInput["phonenumber"] = "+254711223344"
	noPinCredReqInput["pin"] = "has no pin"

	noPinLoginReq, err := mapToJSONReader(noPinCredReqInput)
	if err != nil {
		t.Errorf("unable to get noPinLoginReq JSON io Reader: %s", err)
		return
	}

	type args struct {
		body io.Reader
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
				body: validLoginReq,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "totally incorrect login credentials",
			args: args{
				body: invalidLoginReq,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "invalid login credentials format",
			args: args{
				body: invalidLoginFormatReq,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "wrong pin credentials",
			args: args{
				body: badLoginReq,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
		{
			name: "non-existent login credentials",
			args: args{
				body: nonExistentCredLoginReq,
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "no pin login credentials",
			args: args{
				body: noPinLoginReq,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				http.MethodPost,
				msisdnLoginURL,
				tt.args.body,
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
