package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestCheckPhoneExists_ExistingPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	phone := base.TestUserPhoneNumber

	// remove user then signup user with the phone number then run phone number check
	// ignore the error since it is of no consequence to us
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := base.TestUserPin
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)

	resp2, err2 := s.Signup.CheckPhoneExists(context.Background(), phone)
	assert.Nil(t, err2)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestCreateUserWithPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	type testArgs struct {
		name    string
		phone   string
		pin     string
		flavour base.Flavour
		wantErr bool
		otp     string
	}

	testbed := []testArgs{
		{
			name:    "valid : should consumer create user",
			phone:   phone,
			pin:     "1234",
			flavour: base.FlavourConsumer,
			wantErr: false,
			otp:     otp.OTP,
		},
		{
			name:    "invalid : wrong phone number format",
			phone:   "+254718",
			pin:     "1234",
			flavour: base.FlavourConsumer,
			wantErr: true,
			otp:     otp.OTP,
		},
		{
			name:    "invalid : wrong flavour",
			phone:   phone,
			pin:     "1234",
			flavour: "wrong-flavour",
			wantErr: true,
			otp:     otp.OTP,
		},
		{
			name:    "invalid : pin length",
			phone:   phone,
			pin:     "1234",
			flavour: base.FlavourConsumer,
			wantErr: true,
			otp:     otp.OTP,
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
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Profile)
				assert.NotNil(t, resp.CustomerProfile)
				assert.NotNil(t, resp.SupplierProfile)
			}
		})
	}

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestCreateUserWithPhoneNumber_Consumer(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	pin := "1234"

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestCreateUserWithPhoneNumber_Pro(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}

	pin := "1234"

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestRegisterPushToken_SignUpUseCasesImpl_(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	type arg struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    arg
		want    bool
		wantErr bool
	}{
		{
			name: "valid: push token should be registered",
			args: arg{
				ctx:   ctx,
				token: "123456789",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: short push token should not be registered",
			args: arg{
				ctx:   ctx,
				token: "124",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: unauthenticated context provided",
			args: arg{
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

func TestSignUpUseCasesImpl_VerifyPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	phone := base.TestUserPhoneNumber

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
			name: "Sad Case - Using an existing phonenumber",
			args: args{
				ctx:   ctx,
				phone: phone,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Using an invalid phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+254719XYZ245",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Using an invalid phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+2546",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Using an invalid phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+254-not-vlid-123",
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Using a valid phonenumber",
			args: args{
				ctx:   ctx,
				phone: "+254733445566",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Signup.VerifyPhoneNumber(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestSignUpUseCasesImpl_CompleteSignup(t *testing.T) {
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
		ctx     context.Context
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully complete the signup process",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Try signup using an unauthenticated context",
			args: args{
				ctx:     context.Background(),
				flavour: base.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Using an invalid flavour",
			args: args{
				ctx:     ctx,
				flavour: base.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Signup.CompleteSignup(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.CompleteSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.CompleteSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignUpUseCasesImpl_RetirePushToken(t *testing.T) {
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
			name: "Sad Case - short push token (invalid) should not be retired ",
			args: args{
				ctx:   ctx,
				token: "1234",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unauthenticated context provided",
			args: args{
				ctx:   context.Background(),
				token: "1234",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy Case - retire a valid push token",
			args: args{
				ctx:   ctx,
				token: "12345678910",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Signup.RetirePushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpUseCasesImpl.RetirePushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SignUpUseCasesImpl.RetirePushToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
