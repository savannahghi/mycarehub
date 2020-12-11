package graph_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func getOTPCode(msisdn string, s *profile.Service) (string, error) {

	//  set up ISC call to get an actual  OTP code from otp service
	body := map[string]interface{}{
		"msisdn": msisdn,
	}
	defaultOTP := ""

	resp, err := s.Otp.MakeRequest(http.MethodPost, profile.SendOTP, body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to generate and send otp: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return defaultOTP, fmt.Errorf("unable to generate and send otp, with status code %v", resp.StatusCode)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to convert response to byte: %v", err)
	}

	var codeResp map[string]string
	if err := json.Unmarshal(code, &codeResp); err != nil {
		return defaultOTP, fmt.Errorf("unable to convert response to map: %v", err)
	}

	otpCode := codeResp["otp"]

	return otpCode, nil
}

func TestMSISDNLogin(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	var s *profile.Service = profile.NewService()

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := s.SetUserPIN(ctx, base.TestUserPhoneNumberWithPin, base.TestUserPin)
	if !set {
		t.Errorf("can't set a test pin")
	}
	if err != nil {
		t.Errorf("can't set a test pin: %v", err)
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
			name: "happy case: a correct phone number and pin",
			args: args{
				PhoneNumber: base.TestUserPhoneNumberWithPin,
				Pin:         base.TestUserPin,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "edge case: invalid phone number and pin",
			args: args{
				PhoneNumber: "not a real phone number",
				Pin:         "not a pin",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    false,
		},
		{
			name: "sad case: a correct phone number with a wrong pin",
			args: args{
				PhoneNumber: base.TestUserPhoneNumberWithPin,
				Pin:         "wrong pin number",
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    false,
		},
		{
			name: "sad case: a non-existent phone number and pin",
			args: args{
				PhoneNumber: "+254780654321",
				Pin:         "0000",
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
				log.Printf("raw response: %s", string(data))
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

func TestRESTResetUserPIN(t *testing.T) {

	// Simulate similar running environment to staging otp service
	// get existing envars
	existingGoogleAppCredentials := base.MustGetEnvVar("GOOGLE_APPLICATION_CREDENTIALS")
	existingGcloudProject := base.MustGetEnvVar("GOOGLE_CLOUD_PROJECT")
	existingGcloudProjectNo := base.MustGetEnvVar("GOOGLE_PROJECT_NUMBER")
	existingFirebaseWebApiKey := base.MustGetEnvVar("FIREBASE_WEB_API_KEY")
	existingRootCollectionSuffix := base.MustGetEnvVar("ROOT_COLLECTION_SUFFIX")
	existingEnvironment := base.MustGetEnvVar("ENVIRONMENT")

	// Staging environment envars
	stagingGoogleAppCredentials := base.MustGetEnvVar("GCLOUD_STAGING_SERVICE_ACCOUNT")
	stagingGcloudProject := base.MustGetEnvVar("STAGING_GOOGLE_CLOUD_PROJECT")
	stagingGcloudProjectNo := base.MustGetEnvVar("STAGING_GOOGLE_PROJECT_NUMBER")
	stagingFirebaseWebApiKey := base.MustGetEnvVar("STAGING_FIREBASE_WEB_API_KEY")
	stagingRootCollectionSuffix := base.MustGetEnvVar("STAGING_ROOT_COLLECTION_SUFFIX")
	stagingEnvironment := base.MustGetEnvVar("STAGING_ENVIRONMENT")

	// finally set envars to match staging environment
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", stagingGoogleAppCredentials)
	if err != nil {
		t.Errorf("unable to reset Google Cloud Project env var: %v", err)
		return
	}
	err = os.Setenv("GOOGLE_CLOUD_PROJECT", stagingGcloudProject)
	if err != nil {
		t.Errorf("unable to reset Google Cloud Project env var: %v", err)
		return
	}
	err = os.Setenv("GOOGLE_CLOUD_PROJECT_NUMBER", stagingGcloudProjectNo)
	if err != nil {
		t.Errorf("unable to reset Google Cloud Project env var: %v", err)
		return
	}
	err = os.Setenv("FIREBASE_WEB_API_KEY", stagingFirebaseWebApiKey)
	if err != nil {
		t.Errorf("unable to reset Firebase Web Api Key env var: %v", err)
		return
	}
	err = os.Setenv("ENVIRONMENT", stagingEnvironment)
	if err != nil {
		t.Errorf("unable to reset Environment env var: %v", err)
		return
	}
	err = os.Setenv("ROOT_COLLECTION_SUFFIX", stagingRootCollectionSuffix)
	if err != nil {
		t.Errorf("unable to reset Root Collection Suffix env var: %v", err)
		return
	}

	ctx := base.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	var s *profile.Service = profile.NewService()
	set, err := s.SetUserPIN(ctx, base.TestUserPhoneNumberWithPin, base.TestUserPin)
	if !set {
		t.Errorf("can't set a test pin")
	}
	if err != nil {
		t.Errorf("can't set a test pin: %v", err)
		return
	}

	msisdnLoginURL := fmt.Sprintf("%s/%s", baseURL, "reset_pin")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get request headers %v", err)
		return
	}

	type args struct {
		msisdn string
		pin    string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "happy case: a correct phone number and pin",
			args: args{
				msisdn: base.TestUserPhoneNumberWithPin,
				pin:    base.TestUserPin,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "edge case: invalid phone number and pin",
			args: args{
				msisdn: "not a real phone number",
				pin:    "not a pin",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestInput := map[string]interface{}{}
			requestInput["msisdn"] = tt.args.msisdn
			requestInput["pin_number"] = tt.args.pin

			otpCode, err := getOTPCode(tt.args.msisdn, s)
			if (err != nil) != tt.wantErr {
				t.Errorf("unable to get otp code from the otp service: %v, wantErr %v", err, tt.wantErr)
				return
			}
			requestInput["otp"] = otpCode

			body, err := mapToJSONReader(requestInput)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s, wantErr %v", err, tt.wantErr)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				msisdnLoginURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s, wantErr %v", err, tt.wantErr)
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

			if resp == nil {
				t.Errorf("nil response")
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s, wantErr %v", err, tt.wantErr)
				return
			}

			if data == nil {
				t.Errorf("nil response data")
				return
			}

			if err != nil {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantStatus != resp.StatusCode {
				log.Printf("raw response: %s", string(data))
				t.Errorf("statusCode = %v, wantStatus %v", resp.StatusCode, tt.wantStatus)
				return
			}

		})

	}

	// Remember to restore everything to how it was before the test started running
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", existingGoogleAppCredentials)
	if err != nil {
		t.Errorf("unable to restore Google Cloud Project env var: %v", err)
		return
	}
	err = os.Setenv("GOOGLE_CLOUD_PROJECT", existingGcloudProject)
	if err != nil {
		t.Errorf("unable to restore Google Cloud Project env var: %v", err)
		return
	}
	err = os.Setenv("GOOGLE_CLOUD_PROJECT_NUMBER", existingGcloudProjectNo)
	if err != nil {
		t.Errorf("unable to restore Google Cloud Project Number env var: %v", err)
		return
	}
	err = os.Setenv("FIREBASE_WEB_API_KEY", existingFirebaseWebApiKey)
	if err != nil {
		t.Errorf("unable to restore Firebase Web Api Key env var: %v", err)
		return
	}
	err = os.Setenv("ENVIRONMENT", existingEnvironment)
	if err != nil {
		t.Errorf("unable to restore Environment env var: %v", err)
		return
	}
	err = os.Setenv("ROOT_COLLECTION_SUFFIX", existingRootCollectionSuffix)
	if err != nil {
		t.Errorf("unable to restore Root Collection Suffix env var: %v", err)
		return
	}
}

func TestRequestPinReset(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	requestPinResetUrl := fmt.Sprintf("%s/%s", baseURL, "request_pin_reset")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get request headers %v", err)
		return
	}

	type args struct {
		msisdn string
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid case - PIN that is not registered",
			args: args{
				msisdn: base.TestUserPhoneNumberWithPin,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid phone number",
			args: args{
				msisdn: "011",
			},
			wantErr:    false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestInput := map[string]interface{}{}
			requestInput["msisdn"] = tt.args.msisdn

			body, err := mapToJSONReader(requestInput)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}

			request, err := http.NewRequest(
				http.MethodPost,
				requestPinResetUrl,
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
				log.Printf("raw response: %s", string(data))
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

func TestFindProvider(t *testing.T) {
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

	variables := map[string]interface{}{
		"filters": map[string]string{
			"search": "khan",
		},
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `
						query findProvider($pagination: PaginationInput,$filter:[BusinessPartnerFilterInput],$sort:[BusinessPartnerSortInput]) {
							findProvider(pagination:$pagination,filter:$filter,sort:$sort){
								edges {
									cursor
									node {
									id
									name
									sladeCode
									}
								}
								pageInfo {
									hasNextPage
									hasPreviousPage
									startCursor
									endCursor
								}
							}
						}`,
					"variables": variables,
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestFindBranch(t *testing.T) {
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

	variables := map[string]interface{}{
		"filters": map[string]string{
			"search": "khan",
		},
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": ` 
						query findBranch($pagination: PaginationInput,$filter:[BranchFilterInput],$sort:[BranchSortInput]) {
							findBranch(pagination:$pagination,filter:$filter,sort:$sort){
							edges {
								cursor
								node {
								id
								name
								organizationSladeCode
								branchSladeCode
								}
							}
							pageInfo {
								hasNextPage
								hasPreviousPage
								startCursor
								endCursor
							}
							}
						}`,
					"variables": variables,
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestVerifyEmailOTPMutation(t *testing.T) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()

	if err != nil {
		t.Errorf("failed to initialize firebase: %s", err)
		return
	}

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)

	if err != nil {
		t.Errorf("unable to initialize firestore client:%s", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	otpCode := rand.Int()
	validData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(otpCode),
		"isValid":           true,
		"message":           "Testing email OTP message",
		"timestamp":         time.Now(),
		"email":             "automated.test.user.bewell-app-ci@healthcloud.co.ke",
	}

	_, err = base.SaveDataToFirestore(firestoreClient,
		base.SuffixCollection(base.OTPCollectionName), validData)

	if err != nil {
		t.Errorf("unable to create an otp: %s", err)
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
			name: "Valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": `mutation verifyEmailOTP($email: String!, $otp: String!){
						verifyEmailOTP(email: $email, otp: $otp)
					}`,
					"variables": map[string]interface{}{
						"email": base.TestUserEmail,
						"otp":   strconv.Itoa(otpCode),
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request with a wrong otp",
			args: args{
				query: map[string]interface{}{
					"query": `mutation verifyEmailOTP($email: String!, $otp: String!){
						verifyEmailOTP(email: $email, otp: $otp)
					}`,
					"variables": map[string]interface{}{
						"email": base.TestUserEmail,
						"otp":   "1234",
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
				t.Errorf("unable to make request: %s", err)
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

func TestGraphQLVerifyMSISDNAndPIN(t *testing.T) {
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %v", err)
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query VerifyMSISDNandPIN($msisdn: String!, $pin: String!) {
						verifyMSISDNandPIN(msisdn: $msisdn, pin: $pin)
					  }`,
					"variables": map[string]interface{}{
						"msisdn": base.TestUserPhoneNumberWithPin,
						"pin":    base.TestUserPin,
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
					"query": `query VerifyMSISDNandPIN($msisdn: String!, $pin: String!) {
						verifyMSISDNandPIN(msisdn: $msisdn, pin: $pin)
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
					"query": `query VerifyMSISDNandPIN($msisdn: String!, $pin: String!) {
						verifyMSISDNandPIN(msisdn: $msisdn, pin: $pin)
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
		{
			name: "invalid pin",
			args: args{
				query: map[string]interface{}{
					"query": `query VerifyMSISDNandPIN($msisdn: String!, $pin: String!) {
						verifyMSISDNandPIN(msisdn: $msisdn, pin: $pin)
					  }`,
					"variables": map[string]interface{}{
						"msisdn": base.TestUserPhoneNumber,
						"pin":    "112",
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
					t.Errorf("error not expected %v", data["errors"])
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

func TestSetLanguagePreference(t *testing.T) {
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
					"query": `mutation SetLanguagePreference($language: Language!){
						setLanguagePreference(language: $language)
					  }`,
					"variables": map[string]interface{}{
						"language": "en",
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid: wrong laguage parameter",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SetLanguagePreference($language: Language!){
						setLanguagePreference(language: $language)
					  }`,
					"variables": map[string]interface{}{
						"language": "KGB",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
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
					t.Errorf("error not expected got error: %w", err)
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

func TestUserProfileQuery(t *testing.T) {
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `
					query userProfile(){
						userProfile(){
							id
							verifiedIdentifiers
							isApproved
							termsAccepted
							msisdns
							emails
							photoBase64
							photoContentType
							pushTokens
							covers{
							payerName
							payerSladeCode
							memberNumber
							memberName
							}
							isTester
							active
							dateOfBirth
							gender
							patientID
							name
							bio
							language
							practitionerApproved
							practitionerTermsOfServiceAccepted
							canExperiment
							askAgainToSetIsTester
							askAgainToSetCanExperiment
							VerifiedEmails{
							email
							verified
							}
							VerifiedPhones{
							msisdn
							verified
							}
							hasPin
							hasSupplierAccount
							hasCustomerAccount
							practitionerHasServices

						}
					}`,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request",
			args: args{
				query: map[string]interface{}{
					"query": `
					query userProfile(){
						userProfile(){
							id
							verifiedIdentifiers
							isApproved
							randomQuery
							testingQuery
							emails
							photoBase64
							photoContentType
							pushTokens
							covers{
							payerName
							payerSladeCode
							memberNumber
							memberName
							}
							isTester
							active
							dateOfBirth
							gender
							patientID
							name
							bio
							language
							practitionerApproved
							practitionerTermsOfServiceAccepted
							canExperiment
							askAgainToSetIsTester
							askAgainToSetCanExperiment
							VerifiedEmails{
							email
							verified
							}
							VerifiedPhones{
							msisdn
							verified
							}
							hasPin
							hasSupplierAccount
							hasCustomerAccount
							practitionerHasServices

						}
					}`,
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
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestGraphQLFindProfile(t *testing.T) {
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %v", err)
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query FindProfile{
						findProfile{
							id
						verifiedIdentifiers
						isApproved
						termsAccepted
						msisdns
						emails
						photoBase64
						photoContentType
						pushTokens
						covers{
								payerName
						  payerSladeCode
						  memberNumber
						  memberName
						}
						isTester
						active
						dateOfBirth
						gender
						patientID
						name
						bio
						language
						practitionerApproved
						practitionerTermsOfServiceAccepted
						canExperiment
						askAgainToSetIsTester
						askAgainToSetCanExperiment
						VerifiedEmails {
						  email
						  verified
						}
						VerifiedPhones {
						  msisdn
						  verified
						}
						hasPin
						hasSupplierAccount
						hasCustomerAccount
						practitionerHasServices
					  }
					}`,
				},
			},
			wantStatus: 200,
			wantErr:    false,
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
					t.Errorf("error not expected %v", data["errors"])
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

func TestGraphGQLqueryGetProfile(t *testing.T) {
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query GetProfile($uid: String!){
						getProfile(uid: $uid){
						  id
						  emails
						  active
						  msisdns
						  verifiedIdentifiers
						}
					  }`,
					"variables": map[string]interface{}{
						"uid": "00000000000000000000000001",
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "test get profile with empty uid",
			args: args{
				query: map[string]interface{}{
					"query": `query GetProfile($uid: String!){
						getProfile(uid: $uid){
						  id
						  emails
						  active
						  msisdns
						  verifiedIdentifiers
						}
					  }`,
					"variables": map[string]interface{}{
						"uid": "",
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphQlConfirmEmail(t *testing.T) {
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

	graphqlMutation := `
	mutation ConfirmEmail($emailInput: String!) {
		confirmEmail(email: $emailInput){
			id
			verifiedIdentifiers
			isApproved
			termsAccepted
			msisdns
			emails
			photoBase64
			photoContentType
			pushTokens
			covers{
				payerName
				payerSladeCode
				memberNumber
				memberName
			}
			isTester
			active
			dateOfBirth
			gender
			patientID
			name
			bio
			language
			practitionerApproved
			practitionerTermsOfServiceAccepted
			canExperiment
			askAgainToSetIsTester
			askAgainToSetCanExperiment
			VerifiedEmails {
						email
				verified
			}
			VerifiedPhones {
				verified
				msisdn
			}
			hasPin
			hasSupplierAccount
			hasCustomerAccount
			practitionerHasServices
		}
	}`

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
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"emailInput": fmt.Sprintf("test-%s@healthcloud.co.ke", uuid.New()),
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid email",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"emailInput": "not avalid email",
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

func TestGraphQLListTestersQuery(t *testing.T) {
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
			name: "Valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query listTesters{
						listTesters
					  }`,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query listTesters{
						invalidQuery
					  }`,
				},
			},
			wantStatus: 422,
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphQLUpdateUserProfile(t *testing.T) {
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

	graphqlMutation := `
	mutation updateUserProfile($userProfileInput: UserProfileInput!) {
		updateUserProfile(input: $userProfileInput){
			id
		verifiedIdentifiers
		isApproved
		termsAccepted
		msisdns
		emails
		photoBase64
		photoContentType
		pushTokens
		covers{
			payerName
			payerSladeCode
			memberNumber
			memberName
		}
		isTester
		active
		dateOfBirth
		gender
		patientID
		name
		bio
		language
		practitionerApproved
		practitionerTermsOfServiceAccepted
		canExperiment
		askAgainToSetIsTester
		askAgainToSetCanExperiment
		VerifiedEmails {
					email
			verified
		}
		VerifiedPhones {
			verified
			msisdn
		}
		hasPin
		hasSupplierAccount
		hasCustomerAccount
		practitionerHasServices
		}
	}`

	bs, err := ioutil.ReadFile("profile/testdata/photo.jpg")
	if err != nil {
		t.Errorf("unable to readfile: %v", err)
		return
	}
	photoBase64 := base64.StdEncoding.EncodeToString(bs)

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
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"userProfileInput": map[string]interface{}{
							"photoBase64":                photoBase64,
							"photoContentType":           base.ContentTypeJpg,
							"emails":                     []string{gofakeit.Email()},
							"canExperiment":              false,
							"askAgainToSetIsTester":      false,
							"askAgainToSetCanExperiment": false,
						},
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid email",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"userProfileInput": map[string]interface{}{
							"photoBase64":                photoBase64,
							"photoContentType":           base.ContentTypeJpg,
							"emails":                     []string{"not an email"},
							"canExperiment":              false,
							"askAgainToSetIsTester":      false,
							"askAgainToSetCanExperiment": false,
						},
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

func TestSupplierProfileQuery(t *testing.T) {
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

	grapqhQLQueryPayload := `
	query supplierProfile($uid:String!){
		supplierProfile(uid:$uid){
			userProfile {
				id
				verifiedIdentifiers
				isApproved
				termsAccepted
				msisdns
				emails
				photoBase64
				photoContentType
				pushTokens
				covers {
				  payerName
				  payerSladeCode
				  memberNumber
				  memberName
				}
				isTester
				active
				dateOfBirth
				gender
				patientID
				name
				bio
				language
				practitionerApproved
				practitionerTermsOfServiceAccepted
				canExperiment
				askAgainToSetIsTester
				askAgainToSetCanExperiment
				VerifiedEmails {
				  email
				  verified
				}
				VerifiedPhones {
				  msisdn
				  verified
				}
				hasPin
				hasSupplierAccount
				hasCustomerAccount
				practitionerHasServices
			  }
			  supplierId
			  payablesAccount {
				id
				name
				isActive
				number
				tag
				description
			  }
			  supplierKYC {
				accountType
				identificationDocType
				identificationDocNumber
				identificationDocPhotoBase64
				identificationDocPhotoContentType
				license
				cadre
				profession
				kraPIN
				kraPINDocPhoto
				businessNumber
				businessNumberDocPhotoBase64
				businessNumberDocPhotoContentType
			  }
			  active

		}
	}
	`

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Valid query request",
			args: args{
				query: map[string]interface{}{
					"query": grapqhQLQueryPayload,
					"variables": map[string]interface{}{
						"uid": "e59Ag9JaKNRzmrsqMREQrdNJl5m1",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request - has a non-existent uid",
			args: args{
				query: map[string]interface{}{
					"query": grapqhQLQueryPayload,
					"variables": map[string]interface{}{
						"uid": "not a valid uid",
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
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestGraphGQLmutationPractitionerSignUp(t *testing.T) {
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

	query := map[string]interface{}{}
	query["query"] = `mutation PractitionerSignUp($input: PractitionerSignupInput!){
		practitionerSignUp(input: $input)
	  }	`

	type args struct {
		license   string
		cadre     string
		specialty string
		emails    []string
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
				license:   "123456",
				cadre:     "DOCTOR",
				specialty: "PUBLIC_HEALTH",
				emails:    []string{"be.well@bewell.co.ke", "info@bewell.co.ke"},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid: wrong cadre",
			args: args{
				license:   "123456",
				cadre:     "JUST_PROFESSIONAL",
				specialty: "PUBLIC_HEALTH",
				emails:    []string{"be.well@bewell.co.ke", "info@bewell.co.ke"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: wrong speciality",
			args: args{
				license:   "123456",
				cadre:     "DOCTOR",
				specialty: "JUST_A_SPECIALITY",
				emails:    []string{"be.well@bewell.co.ke", "info@bewell.co.ke"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variables := map[string]interface{}{
				"input": map[string]interface{}{
					"license":   tt.args.license,
					"cadre":     tt.args.cadre,
					"specialty": tt.args.specialty,
					"emails":    tt.args.emails,
				},
			}
			query["variables"] = variables

			body, err := mapToJSONReader(query)
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphGQLmutationCompleteSignUp(t *testing.T) {
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
		query     map[string]interface{}
		variables map[string]interface{}
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
					"query": `mutation {
						completeSignup
					  }	`,
				},
				variables: map[string]interface{}{},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": `mutation CompleteSignup($random: String!) {
						completeSignup
					  }	`,
				},
				variables: map[string]interface{}{
					"random": "unknown parameters",
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.query["variables"] = tt.args.variables
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
					t.Errorf("error not expected got error: %w", err)
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

func TestAcceptTermsAndConditionsMutation(t *testing.T) {
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

	graphQLMutationPayload := `
	mutation acceptTermsAndConditions($accept:Boolean!){
		acceptTermsAndConditions(accept:$accept)
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"accept": true,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"accept": "invalid variable",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestSetUpSupplier(t *testing.T) {
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

	graphQlMutation := `
	mutation createSupplier($dataInput:SupplierAccountInput!){
		setUpSupplier(input:$dataInput) {
		  userProfile {
			id
		  }
		  underOrganization
		  isOrganizationVerified
		  sladeCode
		  parentOrganizationID

		}
	}`

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
					"query": graphQlMutation,
					"variables": map[string]interface{}{
						"dataInput": map[string]interface{}{
							"accountType":       "INDIVIDUAL",
							"underOrganization": false,
						},
					},
				},
			},
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name: "invalid mutation (invalid account type)",
			args: args{
				query: map[string]interface{}{
					"query": graphQlMutation,
					"variables": map[string]interface{}{
						"dataInput": map[string]interface{}{
							"accountType":       "NOT VALID ACCOUNT TYPE",
							"underOrganization": false,
						},
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
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
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
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

func TestGetRegisteredPractitionerQuery(t *testing.T) {
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

	graphQLQueryPayload := `query getKMPDURegisteredPractitioner($regno: String!){
		getKMPDURegisteredPractitioner(regno: $regno){
			name
			regno
			address
			qualifications
			speciality
			subspeciality
			licensetype
			active
		}
	  }`

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
			name: "Valid query request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQueryPayload,
					"variables": map[string]interface{}{
						"regno": "A0008",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request - Has no specified regno",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQueryPayload,
					"variables": map[string]interface{}{
						"regno": "not a valid regno",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid query request - Has a valid regno but a wrong variable specified as reg ",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQueryPayload,
					"variables": map[string]interface{}{
						"reg": "AOOO8",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
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
			log.Printf(string(dataResponse))
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphGQLmutationRemoveTester(t *testing.T) {
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
		query     map[string]interface{}
		variables map[string]interface{}
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
					"query": `mutation RemoveTester($email: String!){
								removeTester(email: $email)
					  }`,
				},
				variables: map[string]interface{}{
					"email": "be.well@bewell.co.ke",
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": `mutation RemoveTester($email: String!){
								removeTester(email: $email)
			  				}`,
				},
				variables: map[string]interface{}{
					"random": "unknown parameters",
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.query["variables"] = tt.args.variables
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphGQLmutationCreateSignUpMethod(t *testing.T) {
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
		query     map[string]interface{}
		variables map[string]interface{}
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
					"query": `mutation CreateSignUpMethod($method: SignUpMethod!){
						createSignUpMethod(signUpMethod: $method)
					  }`,
				},
				variables: map[string]interface{}{
					"method": "google",
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: wrong sign up method",
			args: args{
				query: map[string]interface{}{
					"query": `mutation CreateSignUpMethod($method: SignUpMethod!){
						createSignUpMethod(signUpMethod: $method)
					  }`,
				},
				variables: map[string]interface{}{
					"method": "some alien sign up method",
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.query["variables"] = tt.args.variables
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphGQLmutationGetSignUpMethod(t *testing.T) {
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
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query GetSignUpMethod($id: String!){
						getSignUpMethod(id: $id)
					  }`,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query request",
			args: args{
				query: map[string]interface{}{
					"query": `query GetSignUpMethod($id: String!){
						getSignUpMethod(id: $idd)
					  }`,
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := base.GetLoggedInUserUID(ctx)
			if err != nil {
				t.Errorf("failed to get uid of signed in user: %w", err)
				return
			}
			tt.args.query["variables"] = map[string]interface{}{
				"id": uid,
			}
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
					t.Errorf("error not expected got error: %w", err)
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

func TestGraphGQLmutationAddSupplier(t *testing.T) {
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
		query     map[string]interface{}
		variables map[string]interface{}
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
					"query": `mutation AddSuplier($na: String!, $partnerType: PartnerTypes!){
						addSupplier(name: $na, partnerType: $partnerType){
						  supplierId
						  active
						  sladeCode
						  parentOrganizationID
						}
					  }`,
				},
				variables: map[string]interface{}{
					"na":          "just a name",
					"partnerType": "PRACTITIONER",
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: wrong partner type",
			args: args{
				query: map[string]interface{}{
					"query": `mutation AddSuplier($na: String!, $partnerType: PartnerTypes!){
						addSupplier(name: $na, partnerType: $partnerType){
						  supplierId
						  active
						  sladeCode
						  parentOrganizationID
						}
					  }`,
				},
				variables: map[string]interface{}{
					"na":          "just a name",
					"partnerType": "alien partner type",
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.query["variables"] = tt.args.variables

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
					t.Errorf("error not expected got error: %w", err)
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
