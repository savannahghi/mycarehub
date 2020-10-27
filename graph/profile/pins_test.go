package profile

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMsisdn(t *testing.T) {
	goodData := &PinRecovery{
		MSISDN: "+254712789456",
	}
	goodDataJSONBytes, err := json.Marshal(goodData)
	assert.Nil(t, err)
	assert.NotNil(t, goodDataJSONBytes)

	validRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validRequest.Body = ioutil.NopCloser(bytes.NewReader(goodDataJSONBytes))

	emptyDataRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	emptyDataRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *PinRecovery
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want: &PinRecovery{
				MSISDN: "+254712789456",
			},
			wantErr: false,
		},
		{
			name: "invalid data",
			args: args{
				w: httptest.NewRecorder(),
				r: emptyDataRequest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateMsisdn(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMsisdn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateMsisdn() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
