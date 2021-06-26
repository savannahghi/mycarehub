package otp_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/base"
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

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
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

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make http request")
				}
			}

			if tt.name == "invalid:make_http_request_returns_unexpected_status_code" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
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

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
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

func TestServiceOTPImpl_GenerateAndSendOTP(t *testing.T) {
	ctx := context.Background()
	p := otp.NewOTPService(otpClient, baseExt)

	validRespPayload := `"234234"`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `"otp":"234234"`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.OtpResponse
		wantErr bool
	}{
		{
			name: "valid:_successfully_generate_and_send_otp",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_generate_and_send_otp" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			resp, err := p.GenerateAndSendOTP(tt.args.ctx, tt.args.phone)

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

				if resp.OTP != tt.want.OTP {
					t.Errorf("expected %v, got %v", tt.want.OTP, resp.OTP)
					return
				}
			}

		})
	}
}

func TestServiceOTPImpl_SendRetryOTP(t *testing.T) {
	ctx := context.Background()
	p := otp.NewOTPService(otpClient, baseExt)

	validRespPayload := `"123123"`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `"otp":"123123"`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx       context.Context
		msisdn    string
		retryStep int
	}
	tests := []struct {
		name    string
		args    args
		want    *base.OtpResponse
		wantErr bool
	}{
		{
			name: "valid:_successfully_send_retry_otp",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want: &base.OtpResponse{
				OTP: "123123",
			},
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_normalize_msisdn",
			args: args{
				ctx:       ctx,
				msisdn:    "+asc719ASD678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_send_retry_otp" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}

				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_normalize_msisdn" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize msisdn")
				}
			}

			resp, err := p.SendRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep)

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

				if resp.OTP != tt.want.OTP {
					t.Errorf("expected %v, got %v", tt.want.OTP, resp.OTP)
					return
				}
			}
		})
	}
}

func TestServiceOTPImpl_VerifyEmailOTP(t *testing.T) {
	ctx := context.Background()
	p := otp.NewOTPService(otpClient, baseExt)

	validRespPayload := `{"IsVerified":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `{""}`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))
	type args struct {
		ctx   context.Context
		email string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_successfully_verify_email_otp",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_verify_email_otp" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeotpExt.MakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			resp, err := p.VerifyEmailOTP(tt.args.ctx, tt.args.email, tt.args.otp)

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

				if resp != tt.want {
					t.Errorf("expected %v, got %v", tt.want, resp)
					return
				}
			}
		})
	}
}
