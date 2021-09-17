package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPhoneNumber(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeTestService(ctx)
	fmt.Printf("43: THE ERROR IS: %v\n", err)
	if err != nil {
		t.Errorf("unable to initialize test service")
	}

	validPhoneNumber := interserviceclient.TestUserPhoneNumber
	validPIN := "1234"

	invalidPhoneNumber := "+25471865"

	// clean up
	_ = i.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// try to verify with invalidPhoneNumber. this should fail
	resp, err := i.VerifyPhoneNumber(context.Background(), invalidPhoneNumber, nil)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// verify with validPhoneNumber
	resp, err = i.VerifyPhoneNumber(context.Background(), validPhoneNumber, nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	// clean up
	_ = i.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// register the phone number then try to verify it
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp1, err := i.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &validPIN,
			Flavour:     feedlib.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp1.Profile)
	assert.Equal(t, validPhoneNumber, *resp1.Profile.PrimaryPhone)

	// now try to verify with the already registered phone number
	resp, err = i.VerifyPhoneNumber(context.Background(), validPhoneNumber, nil)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func TestCreateUserWithPhoneNumber_Consumer(t *testing.T) {
	i, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := interserviceclient.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = i.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := i.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)
	assert.NotNil(t, resp.CommunicationSettings)
	assert.Equal(t, true, resp.CommunicationSettings.AllowEmail)
	assert.Equal(t, true, resp.CommunicationSettings.AllowPush)
	assert.Equal(t, true, resp.CommunicationSettings.AllowTextSMS)
	assert.Equal(t, true, resp.CommunicationSettings.AllowWhatsApp)

	// clean up
	_ = i.RemoveUserByPhoneNumber(context.Background(), phone)
}
