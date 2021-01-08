package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// func composeInValidPinPayload(t *testing.T) *domain.SetPINRequest {
// 	return &domain.SetPINRequest{
// 		PhoneNumber: "",
// 		PIN:         "1234",
// 	}

// }

func composeValidChangePinPayload(t *testing.T, otp string) *domain.ChangePINRequest {
	return &domain.ChangePINRequest{
		PhoneNumber: base.TestUserPhoneNumberWithPin,
		PIN:         "1234",
		OTP:         otp,
	}
}

func composeInValidChangePinPayload(t *testing.T, otp string) *domain.ChangePINRequest {
	return &domain.ChangePINRequest{
		PhoneNumber: "",
		PIN:         "1234",
		OTP:         otp,
	}

}

// func composeValidPinPayload(t *testing.T) *domain.SetPINRequest {
// 	return &domain.SetPINRequest{
// 		PhoneNumber: base.TestUserPhoneNumber,
// 		PIN:         "1234",
// 	}
// }

// func composeUnregisteredPhone(t *testing.T) *domain.SetPINRequest {
// 	return &domain.SetPINRequest{
// 		PhoneNumber: base.TestUserPhoneNumber,
// 		PIN:         "1234",
// 	}
// }

func composeInValidPinResetPayload(t *testing.T) *resources.PhoneNumberPayload {
	emptyString := ""
	return &resources.PhoneNumberPayload{
		PhoneNumber: &emptyString,
	}

}

func composeValidPinResetPayload(t *testing.T) *resources.PhoneNumberPayload {
	validNumber := base.TestUserPhoneNumberWithPin
	return &resources.PhoneNumberPayload{
		PhoneNumber: &validNumber,
	}
}

func TestResetPin(t *testing.T) {
	client := http.DefaultClient
	// create a user and their profile
	phoneNumber := base.TestUserPhoneNumberWithPin
	_, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		// return
	}

	// valid change pin payload
	otpResp, err := generateTestOTP(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	validPayload := composeValidChangePinPayload(t, otpResp.OTP)
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// invalid change payload
	secondOtpResp, err := generateTestOTP(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to generate a second test OTP: %v", err)
		return
	}

	badPayload := composeInValidChangePinPayload(t, secondOtpResp.OTP)
	bs2, err := json.Marshal(badPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(bs2)

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
			name: "failure: change pin with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: change pin with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: change pin with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/reset_pin", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
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
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestRequestPINReset(t *testing.T) {
	client := http.DefaultClient
	// create a user and their profile
	phoneNumber := base.TestUserPhoneNumberWithPin
	_, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		// return
	}
	// valid change pin payload
	validPayload := composeValidPinResetPayload(t)
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// invalid change payload
	badPayload := composeInValidPinResetPayload(t)
	bs2, err := json.Marshal(badPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	invalidPayload := bytes.NewBuffer(bs2)

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
			name: "failure: pin reset request with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: pin reset request with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "success: pin reset request with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/request_pin_reset", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
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
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestUpdateUserPIN(t *testing.T) {
	// create a user and thier profile
	phoneNumber := base.TestUserPhoneNumber
	_, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}
	ctx := base.GetAuthenticatedContext(t)

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphqlMutation := `
	mutation updateUserPIN($phone:String!, $pin:String!){
		updateUserPIN(phone:$phone, pin:$pin){
		  profileID
		  pinNumber
		}
	  }
	`
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
			name: "success: update user pin with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phone": base.TestUserPhoneNumber,
						"pin":   "1234",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		// {
		// 	name: "failure: update pin for unregistred user",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"phone": base.TestUserPhoneNumber,
		// 				"pin":   "1234",
		// 			},
		// 		},
		// 	},
		// 	wantStatus: 465,
		// 	wantErr:    true,
		// },
		// {
		// 	name: "failure: update pin with bogus payload",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"phone": "*",
		// 				"pin":   "*",
		// 			},
		// 		},
		// 	},
		// 	wantStatus: 452,
		// 	wantErr:    true,
		// },
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
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
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
			// if tt.wantStatus != resp.StatusCode {
			// 	t.Errorf("Bad status response returned")
			// 	return
			// }
		})
	}
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, phoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}
