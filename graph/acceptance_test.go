package graph_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

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

func TestSendRetryOTP(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	sendRetryOTP := fmt.Sprintf("%s/%s", baseURL, "send_retry_otp")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get request headers %v", err)
		return
	}

	type args struct {
		Msisdn    string
		RetryStep int
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "send retry OTP via whatsapp",
			args: args{
				Msisdn:    base.TestUserPhoneNumberWithPin,
				RetryStep: 1,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "send retry OTP via twilio",
			args: args{
				Msisdn:    base.TestUserPhoneNumberWithPin,
				RetryStep: 2,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "send retry OTP via a non-existent channel",
			args: args{
				Msisdn:    base.TestUserPhoneNumberWithPin,
				RetryStep: 300,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "send retry OTP using invalid credentials",
			args: args{
				Msisdn:    "+254795941530",
				RetryStep: 300,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestInput := map[string]interface{}{}
			requestInput["msisdn"] = tt.args.Msisdn
			requestInput["retryStep"] = tt.args.RetryStep

			body, err := mapToJSONReader(requestInput)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				sendRetryOTP,
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

func TestRequestPinRest(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	requestPinRestUrl := fmt.Sprintf("%s/%s", baseURL, "request_pin_reset")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get request headers %v", err)
		return
	}

	type args struct {
		msisdn    string
		PINNumber string
		otp       string
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid case",
			args: args{
				msisdn:    base.TestUserPhoneNumberWithPin,
				PINNumber: base.TestUserPin,
				otp:       "1234",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid phone number",
			args: args{
				msisdn:    "011",
				PINNumber: base.TestUserPin,
				otp:       "1234",
			},
			wantErr:    false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestInput := map[string]interface{}{}
			requestInput["msisdn"] = tt.args.msisdn
			requestInput["PINNumber"] = tt.args.PINNumber
			requestInput["otp"] = tt.args.otp

			body, err := mapToJSONReader(requestInput)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}

			request, err := http.NewRequest(
				http.MethodPost,
				requestPinRestUrl,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}
			if request == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				request.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(request)
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

func TestCreateUserByPhone(t *testing.T) {
	client := http.DefaultClient
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	createUserURL := fmt.Sprintf("%s/%s", baseURL, "create_user")
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name       string
		args       args
		want       http.HandlerFunc
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful create user",
			args: args{
				phoneNumber: base.TestUserPhoneNumber,
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "unsuccessful create user",
			args: args{
				phoneNumber: "072222222222568",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{}
			payload["msisdn"] = tt.args.phoneNumber

			body, err := mapToJSONReader(payload)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				createUserURL,
				body,
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

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
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
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestVerifySignUpPhoneNumber(t *testing.T) {
	client := http.DefaultClient
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	headers := base.GetDefaultHeaders(t, baseURL, "profile")

	VerifyPhoneURL := fmt.Sprintf("%s/%s", baseURL, "verify_phone")
	type args struct {
		phoneNumber string
	}

	tests := []struct {
		name       string
		args       args
		want       map[string]interface{}
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful verification of an existing user",
			args: args{
				phoneNumber: base.TestUserPhoneNumber,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
			want: map[string]interface{}{
				"isNewUser": false,
				"OTP":       "",
			},
		},
		{
			name: "successful verification of a nonexisting user",
			args: args{
				phoneNumber: "0722222222",
			},
			wantStatus: http.StatusOK,
			wantErr:    true, // Returns an error with status 401 due to an external isc call to otp service
			want: map[string]interface{}{
				"isNewUser": true,
				"OTP":       "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]interface{}{}
			payload["phoneNumber"] = tt.args.phoneNumber

			body, err := mapToJSONReader(payload)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				VerifyPhoneURL,
				body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
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

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestSetUserPin(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SetUserPin($msisdn: String!, $pin: String!){
						setUserPin(msisdn: $msisdn, pin: $pin)
					}`,
					"variables": map[string]interface{}{
						"msisdn": base.TestUserPhoneNumber,
						"pin":    "1234",
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid msisdn",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SetUserPin($msisdn: String!, $pin: String!){
						setUserPin(msisdn: $msisdn, pin: $pin)
					}`,
					"variables": map[string]interface{}{
						"msisdn": "+",
						"pin":    "1234",
					},
				},
			},
			wantStatus: 200,
			wantErr:    true,
		},
		{
			name: "invalid msisdn with string",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SetUserPin($msisdn: String!, $pin: String!){
						setUserPin(msisdn: $msisdn, pin: $pin)
					}`,
					"variables": map[string]interface{}{
						"msisdn": "qwer",
						"pin":    "1234",
					},
				},
			},
			wantStatus: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
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
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}
