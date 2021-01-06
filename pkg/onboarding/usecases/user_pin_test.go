package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func generateTestOTP(t *testing.T) (string, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to initialize test service: %v", err)
	}
	return s.Otp.GenerateAndSendOTP(ctx, base.TestUserPhoneNumber)
}

func TestUserPinUseCaseImpl_SetUserPIN(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx   context.Context
		pin   string
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: valid pin setup - valid payload",
			args: args{
				ctx:   ctx,
				pin:   "1234",
				phone: base.TestUserPhoneNumber,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid payload",
			args: args{
				ctx:   ctx,
				pin:   "",
				phone: base.TestUserPhoneNumber,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid payload - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			authResponse, err := pin.UserPIN.SetUserPIN(tt.args.ctx, tt.args.pin, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.SetUserPIN() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && authResponse != false {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == false {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestUserPinUseCaseImpl_RequestPINReset(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid pin reset request",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid pin reset request - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid pin reset request - wrong phone number",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumberWithPin, // Not the same with primary
				// phone number (TestUserPhoneNumber)
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number format",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			otpResponse, err := pin.UserPIN.RequestPINReset(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.RequestPINReset() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && otpResponse != "" {
				t.Errorf("expected empty string OTP response but got %v, since the error %v occurred",
					otpResponse,
					err,
				)
				return
			}

			if !tt.wantErr && otpResponse == "" {
				t.Errorf("expected an otp response but got empty string, since no error occurred")
				return
			}
		})
	}
}

func TestUserPinUseCaseImpl_ChangeUserPIN(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx   context.Context
		phone string
		pin   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: valid pin setup - valid payload",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				pin:   "12356",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid payload- empty payload",
			args: args{
				ctx:   ctx,
				pin:   "",
				phone: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid payload - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy case: restore the previous pin",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				pin:   base.TestUserPin,
			},
			want:    true,
			wantErr: false,
		}, // Revert to original PIN to prevent Test user login breakages
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			authResponse, err := pin.UserPIN.ChangeUserPIN(tt.args.ctx, tt.args.phone, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.ChangeUserPIN() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && authResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestUserPinUseCaseImpl_ResetUserPIN(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	otp, err := generateTestOTP(t)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	secondOtp, err := generateTestOTP(t)
	if err != nil {
		t.Errorf("failed to generate a second test OTP: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		phone string
		PIN   string
		OTP   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: valid PIN setup - valid payload",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				PIN:   "12356",
				OTP:   otp,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid payload- empty payload",
			args: args{
				ctx:   ctx,
				PIN:   "",
				phone: "",
				OTP:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid payload - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
				PIN:   base.TestUserPin,
				OTP:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
				PIN:   base.TestUserPin,
				OTP:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy case: restore the previous pin",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				PIN:   base.TestUserPin,
				OTP:   secondOtp,
			},
			want:    true,
			wantErr: false,
		}, // Revert to original PIN to prevent Test user login breakages
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			authResponse, err := pin.UserPIN.ResetUserPIN(
				tt.args.ctx,
				tt.args.phone,
				tt.args.PIN,
				tt.args.OTP,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.ResetUserPIN() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && authResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}
