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

	resp2, err2 := s.Onboarding.CheckPhoneExists(context.Background(), phone)
	assert.Nil(t, err2)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

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

func TestRemoveUserByPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := "+254718658120"

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

	// login to asser that the user has really been created. Should pass
	login1, err := s.Login.LoginByPhone(context.Background(), phone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.Equal(t, login1.Profile.ID, resp.Profile.ID)
	assert.Equal(t, login1.Profile.PrimaryPhone, resp.Profile.PrimaryPhone)

	// now remove the user
	err = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
	assert.Nil(t, err)

	// login again. It should fail since the user has been removed
	login2, err := s.Login.LoginByPhone(context.Background(), phone, pin, base.FlavourPro)
	assert.NotNil(t, err)
	assert.Contains(t, "7: failed to get a user profile", err.Error())
	assert.Nil(t, login2)

	// call remove user again. This should fail since the user has been removed already
	err = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
	assert.NotNil(t, err)
}
