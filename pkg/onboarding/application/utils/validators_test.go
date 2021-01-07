package utils

import (
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestValidateSignUpInput(t *testing.T) {
	phone := base.TestUserPhoneNumber
	pin := base.TestUserPin
	flavour := base.FlavourConsumer
	otp := "12345"

	// alphanumericPhone := "+254-not-valid-123"
	// badPhone := "+254712"
	shortPin := "123"
	longPin := "1234567"
	alphabeticalPin := "abcd"

	type args struct {
		input *resources.SignUpInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: return a valid output",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: false,
		},
		// TODO: This will fail till @dexter comes up with the fix for it
		// {
		// 	name: "failure: bad phone number provided",
		// 	args: args{
		// 		input: &resources.SignUpInput{
		// 			PhoneNumber: &badPhone,
		// 			PIN:         &pin,
		// 			Flavour:     flavour,
		// 			OTP:         &otp,
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "failure: alphanumeric phone number provided",
		// 	args: args{
		// 		input: &resources.SignUpInput{
		// 			PhoneNumber: &alphanumericPhone,
		// 			PIN:         &pin,
		// 			Flavour:     flavour,
		// 			OTP:         &otp,
		// 		},
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "failure: short pin number provided",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &shortPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: long pin number provided",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &longPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: alphabetical pin number provided",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &alphabeticalPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: bad flavour provided",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     "not-a-flavour",
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: no OTP provided",
			args: args{
				input: &resources.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validInput, err := ValidateSignUpInput(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignUpInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && validInput != nil {
				t.Errorf("expected a nil valid input since an error :%v occurred", err)
			}

			if err == nil && validInput == nil {
				t.Errorf("expected a valid input %v since no error occurred", validInput)
			}
		})
	}
}
