package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestCheckPhoneExists(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	type testArgs struct {
		name        string
		phone       string
		wantErr     bool
		expectedErr string
	}

	testbed := []testArgs{
		{
			name:    "valid : phone number does not exist",
			phone:   base.TestUserPhoneNumber,
			wantErr: false,
		},

		{
			name:    "valid : phone number already exists",
			phone:   "+254718123098", // use a different number since tear down has not happened yet
			wantErr: true,
		},

		{
			name:        "invalid : wrong phone number format",
			phone:       "71812308",
			wantErr:     true,
			expectedErr: "2: failed to verify OTP",
		},
	}

	for _, tt := range testbed {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				otp, err := generateTestOTP(t, tt.phone)
				if err != nil {
					t.Errorf("failed to generate test OTP: %v", err)
					return
				}
				pin := base.TestUserPin
				// signup user with the phone number then run phone number check
				resp, err := s.Signup.CreateUserByPhone(
					context.Background(),
					&resources.SignUpInput{
						PhoneNumber: &tt.phone,
						PIN:         &pin,
						Flavour:     base.FlavourConsumer,
						OTP:         &otp.OTP,
					},
				)
				if tt.expectedErr == "" {
					assert.Nil(t, err)
					assert.NotNil(t, resp)
				}

				if tt.expectedErr != "" {
					assert.NotNil(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)

					resp2, err2 := s.Signup.CheckPhoneExists(context.Background(), tt.phone)
					assert.Nil(t, err2)
					assert.NotNil(t, resp2)
					assert.Equal(t, false, resp2)
				}

			}

			//todo:(dexter) restore this
			// if !tt.wantErr {

			// 	resp, err := signup.CheckPhoneExists(context.Background(), tt.phone)
			// 	assert.Nil(t, err)
			// 	assert.Equal(t, true, resp)

			// 	// signup user with the phone number then run another phone number check
			// 	resp2, err := signup.CreateUserByPhone(context.Background(), tt.phone, "1234", base.FlavourConsumer)
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp2)

			// 	// now check the phone number that has been used above exists.
			// 	resp3, err := signup.CheckPhoneExists(context.Background(), tt.phone)
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp3)
			// 	assert.Equal(t, true, resp3)
			// }

		})
	}

}

func TestCreateUserWithPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := base.TestUserPhoneNumber
	otp, err := generateTestOTP(t, phone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	type testArgs struct {
		name        string
		phone       string
		pin         string
		flavour     base.Flavour
		wantErr     bool
		expectedErr string
		otp         string
	}

	testbed := []testArgs{
		{
			name:        "valid : should consumer create user",
			phone:       phone,
			pin:         "1234",
			flavour:     base.FlavourConsumer,
			wantErr:     false,
			expectedErr: "",
			otp:         otp.OTP,
		},
	}

	for _, tt := range testbed {

		t.Run(tt.name, func(t *testing.T) {

			resp, err := s.Signup.CreateUserByPhone(
				context.Background(),
				&resources.SignUpInput{
					PhoneNumber: &tt.phone,
					PIN:         &tt.pin,
					Flavour:     tt.flavour,
					OTP:         &tt.otp,
				},
			)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
			//todo:(dexter) restore this
			// if !tt.wantErr {
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp)
			// 	assert.NotNil(t, resp.Profile)
			// 	assert.NotNil(t, resp.CustomerProfile)
			// 	assert.NotNil(t, resp.SupplierProfile)
			// }
		})

	}
}

func TestRegisterPushToken(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: push token should be registered",
			args: args{
				ctx:   ctx,
				token: "123456789",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: short push token should not be registered",
			args: args{
				ctx:   ctx,
				token: "124",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: unauthenticated context provided",
			args: args{
				ctx:   context.Background(),
				token: "123456789",
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Signup.RegisterPushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RegisterPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.RegisterPushToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateUserProfile(t *testing.T) {
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

	dateOfBirth := base.Date{
		Day:   12,
		Year:  2000,
		Month: 2,
	}
	firstName := "Jatelo"
	lastName := "Omera"
	uploadID := "12345"

	userdetails := &resources.UserProfileInput{
		PhotoUploadID: &uploadID,
		DateOfBirth:   &dateOfBirth,
		FirstName:     &firstName,
		LastName:      &lastName,
	}

	updateDOB := &resources.UserProfileInput{
		DateOfBirth: &dateOfBirth,
	}

	updateFirstName := &resources.UserProfileInput{
		FirstName: &firstName,
	}

	updateLastName := &resources.UserProfileInput{
		LastName: &lastName,
	}

	type args struct {
		ctx   context.Context
		input *resources.UserProfileInput
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Successfully update a user profile",
			args: args{
				ctx:   ctx,
				input: userdetails,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update a user firstname",
			args: args{
				ctx:   ctx,
				input: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update a user lastname",
			args: args{
				ctx:   ctx,
				input: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update a user date of birth",
			args: args{
				ctx:   ctx,
				input: updateDOB,
			},
			wantErr: false,
		},
		{
			name: "Sad case - User profile not found",
			args: args{
				ctx:   context.Background(),
				input: userdetails,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Signup.UpdateUserProfile(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RegisterPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantErr {
				t.Errorf("profile was not updated")
				return
			}
		})
	}
}

func TestSignUpUseCasesImpl_GetUserRecoveryPhoneNumbers(t *testing.T) {

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
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
			name: "valid: recovery phone numbers returned",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "invalid: invalid phone input supplied",
			args: args{
				ctx:   ctx,
				phone: "not a valid phone number",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.Signup.GetUserRecoveryPhoneNumbers(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.GetUserRecoveryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil AccountRecoveryPhonesResponse returned")
				return
			}
		})
	}
}
