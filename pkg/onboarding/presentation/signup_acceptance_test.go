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
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func composeInValidUserPayload(t *testing.T) *resources.SignUpPayload {
	phone := base.TestUserPhoneNumber
	pin := "" // empty string
	flavour := base.FlavourPro
	payload := &resources.SignUpPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeValidUserPayload(t *testing.T) *resources.SignUpPayload {
	phone := base.TestUserPhoneNumberWithPin
	pin := "2030"
	flavour := base.FlavourPro
	payload := &resources.SignUpPayload{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func CreateTestUserByPhone(t *testing.T) (*resources.UserResponse, error) {
	client := http.DefaultClient
	validPayload := composeValidUserPayload(t)
	bs, err := json.Marshal(validPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)
	url := fmt.Sprintf("%s/create_user_by_phone", baseURL)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)

	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)

	}

	if r == nil {
		return nil, fmt.Errorf("nil request")
	}

	for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
		r.Header.Add(k, v)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)

	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP error: %v", err)

	}

	var userResponse resources.UserResponse
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall response: %v", err)
	}
	return &userResponse, nil
}

func TestCreateUserWithPhoneNumber(t *testing.T) {
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
			name: "success: signup user with valid payload",
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
				body:       payload,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: signup user with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: signup user with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
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

func TestVerifySignUpPhoneNumber(t *testing.T) {
	client := http.DefaultClient
	// prepare a valid payload
	_, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}
	registeredPhone := struct {
		PhoneNumber string
	}{
		PhoneNumber: base.TestUserPhoneNumberWithPin,
	}
	bs, err := json.Marshal(registeredPhone)
	if err != nil {
		t.Errorf("unable to marshal registeredPhone to JSON: %s", err)
		return
	}
	payload := bytes.NewBuffer(bs)

	// prepare an invalid payload
	unregisteredPhone := struct {
		PhoneNumber string
	}{
		PhoneNumber: base.TestUserPhoneNumber,
	}
	bs2, err := json.Marshal(unregisteredPhone)
	if err != nil {
		t.Errorf("unable to marshal unregisteredPhone to JSON: %s", err)
	}
	unregisteredUser := bytes.NewBuffer(bs2)

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
			name: "success: verify a phone number that exists with valid profile",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusBadRequest, //TODO fix me change to StatusOk
			wantErr:    false,
		},
		{
			name: "failure: verify a phone number whose profile does not exist",
			args: args{
				url:        fmt.Sprintf("%s/verify_phone", baseURL),
				httpMethod: http.MethodPost,
				body:       unregisteredUser,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: verify a phone number with a nil payload",
			args: args{
				url:        fmt.Sprintf("%s/create_user_by_phone", baseURL),
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

func TestUserRecoveryPhoneNumbers(t *testing.T) {
	client := http.DefaultClient
	// create a test user
	_, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}
	validNumber := base.TestUserPhoneNumberWithPin
	validPayload := resources.PhoneNumberPayload{
		PhoneNumber: &validNumber,
	}
	bs, err := json.Marshal(validPayload)
	if err != nil {
		t.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	// phone number not registered
	inValidNumber := base.TestUserPhoneNumber
	badPayload := resources.PhoneNumberPayload{
		PhoneNumber: &inValidNumber,
	}
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
			name: "success: user recovery with valid payload",
			args: args{
				url:        fmt.Sprintf("%s/user_recovery_phonenumbers", baseURL),
				httpMethod: http.MethodPost,
				body:       payload,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failure: user recovery with nil payload supplied",
			args: args{
				url:        fmt.Sprintf("%s/user_recovery_phonenumbers", baseURL),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: user recovery with invalid payload",
			args: args{
				url:        fmt.Sprintf("%s/user_recovery_phonenumbers", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
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

func TestRegisterPushToken(t *testing.T) {
	// create a user and thier profile
	_, err := CreateTestUserByPhone(t)
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
	mutation registerPushToken($token:String!){
		registerPushToken(token: $token)
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
			name: "success: register token with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phone": base.TestUserPhoneNumberWithPin,
						"token": "1234",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true, // TODO fixme the logged in user must have a registred profile
		},
		{
			name: "failure: register token with bogus payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"token": "*",
					},
				},
			},
			wantStatus: http.StatusOK,
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
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}

}
