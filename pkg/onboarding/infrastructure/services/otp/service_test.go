package otp_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
)

var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakeotpExt extMock.ISCClientExtension
var otpClient extension.ISCClientExtension = &fakeotpExt
var baseExt extension.BaseExtension = &fakeBaseExt

func TestServiceOTPImpl_VerifyOTP(t *testing.T) {
	ctx := context.Background()
	p := otp.NewOTPService(otpClient, baseExt)

	validRespPayload := `{"IsVerified":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `{""}`
	respReader1 := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx   context.Context
		phone string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:verify_phone_with_valid_phone",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: false,
		},
		{
			name: "invalid:normalize_phone_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:make_http_request_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:make_http_request_returns_unexpected_status_code",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:unmarshalling_of_respose_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:verify_phone_with_valid_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:normalize_phone_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize phone")
				}
			}

			if tt.name == "invalid:make_http_request_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make http request")
				}
			}

			if tt.name == "invalid:make_http_request_returns_unexpected_status_code" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 400,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:unmarshalling_of_respose_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader1,
					}, nil
				}
			}

			resp, err := p.VerifyOTP(tt.args.ctx, tt.args.phone, tt.args.otp)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if resp != true {
					t.Errorf("response should be true got %v", resp)
					return
				}
			}

		})
	}
}
