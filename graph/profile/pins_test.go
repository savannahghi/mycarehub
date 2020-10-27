package profile

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestPinResetFunc(t *testing.T) {
	ctx := context.Background()
	srv := NewService()
	requestPinReset := RequestPinResetFunc(ctx, srv)

	goodResetDataJSONBytes, err := json.Marshal(&PinRecovery{MSISDN: "+254778990088"})
	assert.Nil(t, err)
	assert.NotNil(t, goodResetDataJSONBytes)
	goodResetDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	goodResetDataRequest.Body = ioutil.NopCloser(bytes.NewReader(goodResetDataJSONBytes))

	emptyResetDataJSONBytes, err := json.Marshal(&PinRecovery{MSISDN: ""})
	assert.Nil(t, err)
	assert.NotNil(t, emptyResetDataJSONBytes)
	emptyResetDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	emptyResetDataRequest.Body = ioutil.NopCloser(bytes.NewReader(emptyResetDataJSONBytes))

	incorrectResetDataJSONBytes, err := json.Marshal(&PinRecovery{MSISDN: "this is not a number"})
	assert.Nil(t, err)
	assert.NotNil(t, incorrectResetDataJSONBytes)
	incorrectResetDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	incorrectResetDataRequest.Body = ioutil.NopCloser(bytes.NewReader(incorrectResetDataJSONBytes))

	nonExistentResetDataJSONBytes, err := json.Marshal(&PinRecovery{MSISDN: "+254712789456"})
	assert.Nil(t, err)
	assert.NotNil(t, nonExistentResetDataJSONBytes)
	nonExistentResetDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	nonExistentResetDataRequest.Body = ioutil.NopCloser(bytes.NewReader(nonExistentResetDataJSONBytes))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "good pin reset request",
			args: args{
				w: httptest.NewRecorder(),
				r: goodResetDataRequest,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "empty pin reset request",
			args: args{
				w: httptest.NewRecorder(),
				r: emptyResetDataRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "incorrect pin reset request",
			args: args{
				w: httptest.NewRecorder(),
				r: incorrectResetDataRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "non existent user pin reset request",
			args: args{
				w: httptest.NewRecorder(),
				r: nonExistentResetDataRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPinReset(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, rec.Code, tt.wantStatusCode)
		})
	}
}

func TestUpdatePinHandler(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	loginFunc := UpdatePinHandler(ctx, service)
	otp, _ := service.otpService.GenerateAndSendOTP("+254778990088")

	goodUpdatePinPayloadJSONBytes, err := json.Marshal(&PinRecovery{
		MSISDN: "+254778990088",
		PIN:    "1234",
		OTP:    otp,
	})

	assert.Nil(t, err)
	assert.NotNil(t, goodUpdatePinPayloadJSONBytes)

	goodUpdatePinPayloadReq := httptest.NewRequest(http.MethodGet, "/", nil)
	goodUpdatePinPayloadReq.Body = ioutil.NopCloser(bytes.NewReader(goodUpdatePinPayloadJSONBytes))

	invalidOTPJSONBytes, err := json.Marshal(&PinRecovery{
		MSISDN: "+254778990088",
		PIN:    "1234",
		OTP:    otp,
	})

	assert.Nil(t, err)
	assert.NotNil(t, invalidOTPJSONBytes)

	invalidOTPloadReq := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidOTPloadReq.Body = ioutil.NopCloser(bytes.NewReader(invalidOTPJSONBytes))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "correct login credentials",
			args: args{
				w: httptest.NewRecorder(),
				r: goodUpdatePinPayloadReq,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid otp",
			args: args{
				w: httptest.NewRecorder(),
				r: invalidOTPloadReq,
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginFunc(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, rec.Code, tt.wantStatusCode)
		})
	}
}
