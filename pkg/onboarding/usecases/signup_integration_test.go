package usecases_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPhoneNumber(t *testing.T) {

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := interserviceclient.TestUserPhoneNumber
	validPIN := "1234"

	invalidPhoneNumber := "+25471865"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// try to verify with invalidPhoneNumber. this should fail
	testAppID := uuid.New().String()
	resp, err := s.Signup.VerifyPhoneNumber(context.Background(), invalidPhoneNumber, &testAppID)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// verify with validPhoneNumber
	resp, err = s.Signup.VerifyPhoneNumber(context.Background(), validPhoneNumber, &testAppID)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// register the phone number then try to verify it
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp1, err := s.Signup.CreateUserByPhone(
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
	resp, err = s.Signup.VerifyPhoneNumber(context.Background(), validPhoneNumber, &testAppID)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

}

func TestCreateUserWithPhoneNumber_Consumer(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := interserviceclient.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
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
	assert.NotNil(t, resp.CommunicationSettings)
	assert.Equal(t, true, resp.CommunicationSettings.AllowEmail)
	assert.Equal(t, true, resp.CommunicationSettings.AllowPush)
	assert.Equal(t, true, resp.CommunicationSettings.AllowTextSMS)
	assert.Equal(t, true, resp.CommunicationSettings.AllowWhatsApp)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestCreateUserWithPhoneNumber_Pro(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := interserviceclient.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)

	otp, err := generateTestOTP(t, phone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourPro,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CommunicationSettings)
	assert.Equal(t, true, resp.CommunicationSettings.AllowEmail)
	assert.Equal(t, true, resp.CommunicationSettings.AllowPush)
	assert.Equal(t, true, resp.CommunicationSettings.AllowTextSMS)
	assert.Equal(t, true, resp.CommunicationSettings.AllowWhatsApp)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
}

func TestCreateUserByPhone(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	invalidPhoneNumber := "+25471865"
	invalidPinLength1 := "123"
	invalidPinLength2 := "12345678"
	invalidOtp := "00001234"

	validPhoneNumber := interserviceclient.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := feedlib.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// assert that no user has been created with the same phone number
	// we use GetUserRecoveryPhoneNumbers which searches for users using phone numbers.
	// this should fail since no user has been created yet
	resp1, err := s.Signup.GetUserRecoveryPhoneNumbers(context.Background(), validPhoneNumber)
	assert.NotNil(t, err)
	assert.Nil(t, resp1)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// try registration with invalid phone number. This should fail
	resp2, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &invalidPhoneNumber,
			PIN:         &validPIN,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp2)

	// try registration with invalidPinLength1. this should fail
	resp3, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &invalidPinLength1,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp3)

	// try registration with invalidPinLength2. this should fail
	resp4, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &invalidPinLength2,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp4)

	// try registration with invalidOtp. this should fail
	resp5, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &invalidPinLength2,
			Flavour:     validFlavourConsumer,
			OTP:         &invalidOtp,
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp5)

	// this should pass
	resp6, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &validPIN,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp6)
	assert.NotNil(t, resp6.Profile)
	assert.Equal(t, validPhoneNumber, *resp6.Profile.PrimaryPhone)
	assert.NotNil(t, resp6.Profile.UserName)
	assert.NotNil(t, resp6.CommunicationSettings)
	assert.Equal(t, true, resp6.CommunicationSettings.AllowEmail)
	assert.Equal(t, true, resp6.CommunicationSettings.AllowPush)
	assert.Equal(t, true, resp6.CommunicationSettings.AllowTextSMS)
	assert.Equal(t, true, resp6.CommunicationSettings.AllowWhatsApp)

	// try registration again with the same phone number. This should fail
	// this should pass
	resp7, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &validPIN,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, resp7)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)
}

func TestUpdateUserProfile(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := interserviceclient.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := feedlib.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &validPIN,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.Equal(t, validPhoneNumber, *resp.Profile.PrimaryPhone)
	assert.NotNil(t, resp.Profile.UserName)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	dateOfBirth1 := scalarutils.Date{
		Day:   12,
		Year:  1998,
		Month: 2,
	}
	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}

	firstName1 := "makmende1"
	lastName1 := "Omera1"
	firstName2 := "makmende2"
	lastName2 := "Omera2"
	uploadID1 := "photo-url1"
	uploadID2 := "photo-url2"

	justUploadID := &dto.UserProfileInput{
		PhotoUploadID: &uploadID1,
	}

	justDOB := &dto.UserProfileInput{
		DateOfBirth: &dateOfBirth1,
	}

	justFirstName := &dto.UserProfileInput{
		FirstName: &firstName1,
	}

	justLastName := &dto.UserProfileInput{
		LastName: &lastName1,
	}

	completeUserDetails := &dto.UserProfileInput{
		PhotoUploadID: &uploadID2,
		DateOfBirth:   &dateOfBirth2,
		FirstName:     &firstName2,
		LastName:      &lastName2,
	}

	// update just the  photo upload id
	resp1, err := s.Signup.UpdateUserProfile(authenticatedContext, justUploadID)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.Equal(t, *justUploadID.PhotoUploadID, resp1.PhotoUploadID)

	// update just the date of birth
	resp2, err := s.Signup.UpdateUserProfile(authenticatedContext, justDOB)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, *justDOB.DateOfBirth, *resp2.UserBioData.DateOfBirth)

	// update just the firstname
	resp3, err := s.Signup.UpdateUserProfile(authenticatedContext, justFirstName)
	assert.Nil(t, err)
	assert.NotNil(t, resp3)
	assert.Equal(t, firstName1, *resp3.UserBioData.FirstName)

	// update just the lastname
	resp4, err := s.Signup.UpdateUserProfile(authenticatedContext, justLastName)
	assert.Nil(t, err)
	assert.NotNil(t, resp4)
	assert.Equal(t, lastName1, *resp4.UserBioData.LastName)

	// update with the entire update input
	resp5, err := s.Signup.UpdateUserProfile(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)
	assert.NotNil(t, resp5)
	assert.Equal(t, firstName2, *resp5.UserBioData.FirstName)
	assert.Equal(t, lastName2, *resp5.UserBioData.LastName)
	assert.Equal(t, *completeUserDetails.DateOfBirth, *resp5.UserBioData.DateOfBirth)
	assert.Equal(t, *completeUserDetails.PhotoUploadID, resp5.PhotoUploadID)
	assert.NotEqual(t, firstName1, *resp5.UserBioData.FirstName)
	assert.NotEqual(t, lastName1, *resp5.UserBioData.LastName)
	assert.NotEqual(t, *justDOB.DateOfBirth, *resp5.UserBioData.DateOfBirth)
	assert.NotEqual(t, *justUploadID.PhotoUploadID, resp5.PhotoUploadID)

	// try update with an invalid context
	resp6, err := s.Signup.UpdateUserProfile(context.Background(), completeUserDetails)
	assert.NotNil(t, err)
	assert.Nil(t, resp6)
}

func TestRegisterPushToken(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := interserviceclient.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// use wrong context. this should fail
	respUpt, err := s.Signup.RegisterPushToken(context.Background(), "token1")
	assert.NotNil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, false, respUpt)

	respUpt, err = s.Signup.RegisterPushToken(authenticatedContext, "token1")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	respUpt, err = s.Signup.RegisterPushToken(authenticatedContext, "token2")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	// invalid token length
	respUpt, err = s.Signup.RegisterPushToken(authenticatedContext, "tok")
	assert.NotNil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, false, respUpt)
}

func TestRetirePushToken(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := interserviceclient.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// use wrong context. this should fail
	respUpt, err := s.Signup.RegisterPushToken(context.Background(), "token1")
	assert.NotNil(t, err)
	assert.NotNil(t, respUpt)

	respUpt, err = s.Signup.RegisterPushToken(authenticatedContext, "token1")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	respUpt, err = s.Signup.RegisterPushToken(authenticatedContext, "token2")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	// retire token1
	respUpt, err = s.Signup.RetirePushToken(authenticatedContext, "token1")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	// retire token2
	respUpt, err = s.Signup.RetirePushToken(authenticatedContext, "token2")
	assert.Nil(t, err)
	assert.NotNil(t, respUpt)
	assert.Equal(t, true, respUpt)

	// fetch the profile and assert the push token slice length
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.PushTokens))
}

func TestGetUserRecoveryPhoneNumbers(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := interserviceclient.TestUserPhoneNumber
	pin := "1234"
	invalidPhoneNumber := "+25471865"
	nonExistentPhone := interserviceclient.TestUserPhoneNumberWithPin

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), nonExistentPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)

	// fetch user recovery phone numbers
	resp1, err := s.Signup.GetUserRecoveryPhoneNumbers(context.Background(), primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.Equal(t, 1, len(resp1.MaskedPhoneNumbers))
	assert.Equal(t, 1, len(resp1.UnMaskedPhoneNumbers))
	assert.Equal(t, primaryPhone, resp1.UnMaskedPhoneNumbers[0])

	// fetch user recovery phone numbers using invalidPhoneNumber. should fail
	resp2, err := s.Signup.GetUserRecoveryPhoneNumbers(context.Background(), invalidPhoneNumber)
	assert.NotNil(t, err)
	assert.Nil(t, resp2)

	// fetch user recovery phone numbers using nonExistentPhone. should fail
	resp3, err := s.Signup.GetUserRecoveryPhoneNumbers(context.Background(), nonExistentPhone)
	assert.NotNil(t, err)
	assert.Nil(t, resp3)

}

func TestRemoveUserByPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	phone := "+254718658120"
	pin := "1234"

	otp, err := generateTestOTP(t, phone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&dto.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourPro,
			OTP:         &otp.OTP,
		},
	)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)

	// login to asser that the user has really been created. Should pass
	login1, err := s.Login.LoginByPhone(context.Background(), phone, pin, feedlib.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.Equal(t, login1.Profile.ID, resp.Profile.ID)
	assert.Equal(t, login1.Profile.PrimaryPhone, resp.Profile.PrimaryPhone)

	// now remove the user
	err = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
	assert.Nil(t, err)

	// login again. It should fail since the user has been removed
	login2, err := s.Login.LoginByPhone(context.Background(), phone, pin, feedlib.FlavourPro)
	assert.NotNil(t, err)
	assert.Contains(t, "7: failed to get a user profile", err.Error())
	assert.Nil(t, login2)

	// call remove user again. This should fail since the user has been removed already
	err = s.Signup.RemoveUserByPhoneNumber(context.Background(), phone)
	assert.NotNil(t, err)
}
