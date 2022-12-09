package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

func createSendOTPPayload(phonenumber string, flavour feedlib.Flavour) []byte {
	payload := &dto.SendOTPInput{
		PhoneNumber: phonenumber,
		Flavour:     flavour,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return marshalled
}

func createLoginPayload(flavour feedlib.Flavour, username, organisationID, phonenumber, pin string) []byte {
	payload := &dto.LoginInput{
		OrganisationID: organisationID,
		Username:       username,
		PhoneNumber:    phonenumber,
		PIN:            pin,
		Flavour:        flavour,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return marshalled
}

func TestMyCareHubHandlersInterfacesImpl_SendOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, "invalid flavour")
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)
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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/send_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/send_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Close = true

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

			defer resp.Body.Close()

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_SendRetryOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, "invalid flavour")
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/send_retry_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/send_retry_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Close = true

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

			defer resp.Body.Close()

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_RequestPINReset(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, "invalid flavour")
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)
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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/request_pin_reset",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/request_pin_reset",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidPayload := createLoginPayload("invalid flavour", gofakeit.Username(), gofakeit.UUID(), phoneNumber, "1234")
	invalidPayload1 := createLoginPayload(feedlib.FlavourConsumer, gofakeit.Username(), gofakeit.UUID(), "", "1234")
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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/login_by_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/login_by_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_VerifyPhone(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, "invalid flavour")
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_VerifyOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, "invalid flavour")
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

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
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_GetUserRespondedSecurityQuestions(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

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
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/get_user_responded_security_questions",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_RefreshToken(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := &dto.RefreshTokenPayload{
		UserID: nil,
	}
	marshalled, err := json.Marshal(invalidPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Missing user ID",
			args: args{
				url: fmt.Sprintf(
					"%s/refresh_token",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalled),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_RefreshGetStreamToken(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := &dto.RefreshTokenPayload{
		UserID: nil,
	}

	marshalled, err := json.Marshal(invalidPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Missing user ID",
			args: args{
				url: fmt.Sprintf(
					"%s/refresh_getstream_token",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalled),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_RegisterKenyaEMRPatients(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	emptyPayload := dto.PatientsPayload{
		Patients: []*dto.PatientRegistrationPayload{},
	}
	marshalledEmptyPayload, err := json.Marshal(emptyPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad case - empty payload",
			args: args{
				url: fmt.Sprintf(
					"%s/kenya-emr/register_patient",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalledEmptyPayload),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_CreateOrUpdateKenyaEMRAppointment(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	missingMFLCode := dto.FacilityAppointmentsPayload{
		MFLCode:      "",
		Appointments: []dto.AppointmentPayload{},
	}
	marshalledMissingMFLCode, err := json.Marshal(missingMFLCode)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	missingAppointment := dto.FacilityAppointmentsPayload{
		MFLCode:      "1234",
		Appointments: []dto.AppointmentPayload{},
	}
	marshalledMissingAppointment, err := json.Marshal(missingAppointment)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad case - invalid payload missing MFL code",
			args: args{
				url: fmt.Sprintf(
					"%s/kenya-emr/appointments",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalledMissingMFLCode),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad case - invalid payload Missing Appointment",
			args: args{
				url: fmt.Sprintf(
					"%s/kenya-emr/appointments",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalledMissingAppointment),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_AddFacilityFHIRID(t *testing.T) {
	emptyFacilityID := dto.UpdateFacilityPayload{
		FacilityID:         "",
		FHIROrganisationID: uuid.New().String(),
	}
	marshalledemptyFacilityID, err := json.Marshal(emptyFacilityID)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	noFHIROrgIDPayload := dto.UpdateFacilityPayload{
		FacilityID:         uuid.New().String(),
		FHIROrganisationID: "",
	}

	emptyFHIROrgIDPayload, err := json.Marshal(noFHIROrgIDPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/internal/facilities", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalledemptyFacilityID),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - No FHIROrganisationID in payload",
			args: args{
				url:        fmt.Sprintf("%s/internal/facilities", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(emptyFHIROrgIDPayload),
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
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "mycarehub") {
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_GetClientHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidparams := url.Values{}
	invalidparams.Add("invalid", "invalid")
	invalidparams.Add("invalid", "invalid")

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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/health_diary?%s", baseURL, invalidparams.Encode()),
				httpMethod: http.MethodGet,
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_RegisteredFacilityPatients(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidparams := url.Values{}
	invalidparams.Add("invalid", "invalid")
	invalidparams.Add("invalid", "invalid")

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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/patients?%s", baseURL, invalidparams.Encode()),
				httpMethod: http.MethodGet,
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_GetServiceRequests(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidparams := url.Values{}
	invalidparams.Add("invalid", "invalid")
	invalidparams.Add("invalid", "invalid")

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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/service_request?%s", baseURL, invalidparams.Encode()),
				httpMethod: http.MethodGet,
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_CreatePinResetServiceRequest(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	missingPayload := &dto.PinResetServiceRequestPayload{}
	invalidPayload, err := json.Marshal(missingPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	invalidPayload1, err := json.Marshal(&dto.PinResetServiceRequestPayload{CCCNumber: "12", PhoneNumber: "+254"})
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/service-requests", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Fail to create service request",
			args: args{
				url:        fmt.Sprintf("%s/service-requests", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_AddPatientsRecords(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	missingPayload := &dto.PatientsRecordsPayload{}
	invalidPayload, err := json.Marshal(missingPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	missingMFLCode := &dto.PatientsRecordsPayload{
		Records: []dto.PatientRecordPayload{
			{CCCNumber: "1234"},
		},
	}
	invalidPayload1, err := json.Marshal(missingMFLCode)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	missingRecords := &dto.PatientsRecordsPayload{
		MFLCode: "1234",
	}
	invalidPayload2, err := json.Marshal(missingRecords)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/observations", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing MFL Code",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/observations", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing Records",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/observations", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload2),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_GetUserProfile(t *testing.T) {
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
			name: "Sad Case -  invalid user id",
			args: args{
				url:        fmt.Sprintf("%s/internal/user-profile/ae701bee-aa65-41dc-8953-1c668d0aa14d", baseURL),
				httpMethod: http.MethodGet,
			},
			wantStatus: http.StatusInternalServerError,
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
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "mycarehub") {
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_AddClientFHIRID(t *testing.T) {
	invalidPayload, err := json.Marshal(&dto.ClientFHIRPayload{})
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	invalidPayload1, err := json.Marshal(&dto.ClientFHIRPayload{ClientID: gofakeit.UUID()})
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	invalidPayload2, err := json.Marshal(&dto.ClientFHIRPayload{FHIRID: gofakeit.UUID()})
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/internal/add-fhir-id", baseURL),
				httpMethod: http.MethodPatch,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - missing client ID in payload",
			args: args{
				url:        fmt.Sprintf("%s/internal/add-fhir-id", baseURL),
				httpMethod: http.MethodPatch,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - missing fhir id in payload",
			args: args{
				url:        fmt.Sprintf("%s/internal/add-fhir-id", baseURL),
				httpMethod: http.MethodPatch,
				body:       bytes.NewBuffer(invalidPayload2),
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
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "mycarehub") {
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_GetAppointmentServiceRequests(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	var (
		emptyPayload   = url.Values{}
		missingMFLCode = url.Values{
			"lastSyncTime": {"2006-01-02T15:04:05Z"},
		}
		missingLastSyncTime = url.Values{
			"MFLCODE": {"212121212121"},
		}
		zeroMFLCode = url.Values{
			"MFLCODE":      {"0"},
			"lastSyncTime": {"2006-01-02T15:04:05Z"},
		}
	)

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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/appointment-service-request?%s", baseURL, emptyPayload.Encode()),
				httpMethod: http.MethodGet,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing MFL Code",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/appointment-service-request?%s", baseURL, missingMFLCode.Encode()),
				httpMethod: http.MethodGet,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - missing last sync time",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/appointment-service-request?%s", baseURL, missingLastSyncTime.Encode()),
				httpMethod: http.MethodGet,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - zero MFL Code",
			args: args{
				url:        fmt.Sprintf("%s/kenya-emr/appointment-service-request?%s", baseURL, zeroMFLCode.Encode()),
				httpMethod: http.MethodGet,
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_DeleteUser(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidInput := &dto.PhoneInput{}
	invalidPayload, err := json.Marshal(invalidInput)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	invalidFlavourInput := &dto.PhoneInput{
		Flavour: "invalid",
	}
	invalidFlavourPayload, err := json.Marshal(invalidFlavourInput)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/delete-user", baseURL),
				httpMethod: http.MethodDelete,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Invalid flavour",
			args: args{
				url:        fmt.Sprintf("%s/delete-user", baseURL),
				httpMethod: http.MethodDelete,
				body:       bytes.NewBuffer(invalidFlavourPayload),
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_FetchContactOrganisations(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
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
			name: "Sad Case - Empty payload",
			args: args{
				url:        fmt.Sprintf("%s/contact_organisations", baseURL),
				httpMethod: http.MethodGet,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Invalide phone number",
			args: args{
				url:        fmt.Sprintf("%s/contact_organisations?phoneNumber=%s", baseURL, gofakeit.Phone()),
				httpMethod: http.MethodGet,
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

			dataResponse, err := io.ReadAll(resp.Body)
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

func TestMyCareHubHandlersInterfacesImpl_Organisations(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
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
			name: "Happy Case - list organisations",
			args: args{
				url:        fmt.Sprintf("%s/organisations", baseURL),
				httpMethod: http.MethodGet,
				body:       nil,
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

			dataResponse, err := io.ReadAll(resp.Body)
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
