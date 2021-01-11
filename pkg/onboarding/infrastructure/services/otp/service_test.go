package otp_test

import (
	"context"
	"fmt"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
)

func TestSendRetryOTP(t *testing.T) {
	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		t.Errorf("can't instantiate firebase repository in resolver: %w", err)
		return
	}

	otpservice := otp.NewOTPService(fr)
	type args struct {
		ctx       context.Context
		msisdn    string
		retryStep int
	}

	tests := []struct {
		name        string
		args        args
		want        bool
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Happy case - Send Retry OTP via WhatsApp",
			args: args{
				ctx:       context.Background(),
				msisdn:    base.TestUserPhoneNumberWithPin,
				retryStep: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case - Send Retry OTP via Twilio",
			args: args{
				ctx:       context.Background(),
				msisdn:    base.TestUserPhoneNumberWithPin,
				retryStep: 2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - Send Retry OTP via a nonexistent channel",
			args: args{
				ctx:       context.Background(),
				msisdn:    base.TestUserPhoneNumberWithPin,
				retryStep: 121,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Send Retry OTP using nonexistent credentials",
			args: args{
				ctx:       context.Background(),
				msisdn:    "",
				retryStep: 1,
			},
			want:        false,
			wantErr:     true,
			expectedErr: exceptions.NormalizeMSISDNError(fmt.Errorf("Send Retry OTP using nonexistent credentials")).Error(),
		},
		{
			name: "Sad case - Send Retry OTP using an invalid phonenumber",
			args: args{
				ctx:       context.Background(),
				msisdn:    "+254712kjf787",
				retryStep: 1,
			},
			want:        false,
			wantErr:     true,
			expectedErr: exceptions.NormalizeMSISDNError(fmt.Errorf("Send Retry OTP using an invalid phonenumber")).Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := otpservice.SendRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRetryOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
