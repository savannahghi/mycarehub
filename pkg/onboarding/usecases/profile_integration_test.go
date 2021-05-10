package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

func TestUpdateUserProfileUserName(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.Profile.UserName)

	// login and assert whether the profile matches the one created earlier
	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.Profile.UserName)
	assert.Equal(t, *login1.Profile.UserName, *resp.Profile.UserName)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	err = s.Onboarding.UpdateUserName(authenticatedContext, "makmende1")
	assert.Nil(t, err)

	pr1, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr1)
	assert.NotNil(t, pr1.UserName)
	assert.NotEqual(t, *login1.Profile.UserName, *pr1.UserName)
	assert.NotEqual(t, *resp.Profile.UserName, *pr1.UserName)

	// update the profile with the same userName. It should fail since the userName has already been taken.
	err = s.Onboarding.UpdateUserName(authenticatedContext, "makmende1")
	assert.NotNil(t, err)

	// update with a new unique user name
	err = s.Onboarding.UpdateUserName(authenticatedContext, "makmende2")
	assert.Nil(t, err)

	pr2, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr2)
	assert.NotNil(t, pr2.UserName)
	assert.NotEqual(t, *login1.Profile.UserName, *pr2.UserName)
	assert.NotEqual(t, *resp.Profile.UserName, *pr2.UserName)
	assert.NotEqual(t, *pr1.UserName, *pr2.UserName)

}

func TestSetPhoneAsPrimary(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	secondaryPhone := base.TestUserPhoneNumberWithPin
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), secondaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("failed to create a user by phone")
		return
	}

	if resp == nil {
		t.Error("nil user response returned")
		return
	}

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	if err != nil {
		t.Errorf("an error occured while logging in by phone")
		return
	}

	if login1 == nil {
		t.Errorf("nil response returned")
		return
	}

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// try to login with secondaryPhone. This should fail because secondaryPhone != primaryPhone
	login2, err := s.Login.LoginByPhone(context.Background(), secondaryPhone, pin, base.FlavourConsumer)
	if err == nil {
		t.Errorf("expected an error :%v", err)
		return
	}

	if login2 != nil {
		t.Errorf("the response was not expected")
		return
	}

	// add a secondary phone number to the user
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone})
	if err != nil {
		t.Errorf("failed to add a secondary number to the user")
		return
	}

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user")
		return
	}

	if pr == nil {
		t.Errorf("nil response returned")
		return
	}
	// check if the length of secondary number == 1
	if len(pr.SecondaryPhoneNumbers) != 1 {
		t.Errorf("expected the value to be equal to 1")
		return
	}

	// login to add assert the secondary phone number has been added
	login3, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	if err != nil {
		t.Errorf("expected an error :%v", err)
		return
	}

	if login3 == nil {
		t.Errorf("the response was not expected")
		return
	}

	// check if the length of secondary number == 1
	if len(login3.Profile.SecondaryPhoneNumbers) != 1 {
		t.Errorf("expected the value to be equal to 1")
		return
	}

	// send otp to the secondary phone number we intend to make primary
	otpResp, err := s.Engagement.GenerateAndSendOTP(context.Background(), secondaryPhone)
	if err != nil {
		t.Errorf("unable to send generate and send otp :%v", err)
		return
	}

	if otpResp == nil {
		t.Errorf("unexpected response")
		return
	}

	// set the old secondary phone number as the new primary phone number
	setResp, err := s.Signup.SetPhoneAsPrimary(context.Background(), secondaryPhone, otpResp.OTP)
	if err != nil {
		t.Errorf("failed to set phone as primary: %v", err)
		return
	}

	if setResp == false {
		t.Errorf("unexpected response")
		return
	}

	// login with the old primary phone number. This should fail
	login4, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	if err == nil {
		t.Errorf("unexpected error occured! :%v", err)
		return
	}

	if login4 != nil {
		t.Errorf("unexpected error occured! Expected this to fail")
		return
	}

	// login with the new primary phone number. This should not fail. Assert that the primary phone number
	// is the new one and the secondary phone slice contains the old primary phone number.
	login5, err := s.Login.LoginByPhone(context.Background(), secondaryPhone, pin, base.FlavourConsumer)
	if err != nil {
		t.Errorf("failed to login by phone :%v", err)
		return
	}

	if login5 == nil {
		t.Errorf("the response was not expected")
		return
	}

	if secondaryPhone != *login5.Profile.PrimaryPhone {
		t.Errorf("expected %v and %v to be equal", secondaryPhone, *login5.Profile.PrimaryPhone)
		return
	}

	_, exist := utils.FindItem(login5.Profile.SecondaryPhoneNumbers, secondaryPhone)
	if exist {
		t.Errorf("the secondary phonenumber slice %v, does not contain %v",
			login5.Profile.SecondaryPhoneNumbers,
			secondaryPhone,
		)
		return
	}

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), secondaryPhone)
}

func TestAddSecondaryPhoneNumbers(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	secondaryPhone1 := base.TestUserPhoneNumberWithPin
	secondaryPhone2 := "+25712345690"
	secondaryPhone3 := "+25710375600"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("failed to create a user by phone")
		return
	}

	if resp == nil {
		t.Error("nil user response returned")
		return
	}

	if resp.Profile == nil {
		t.Error("nil profile response returned")
		return
	}

	if resp.CustomerProfile == nil {
		t.Error("nil customer profile response returned")
		return
	}

	if resp.SupplierProfile == nil {
		t.Error("nil supplier profile response returned")
		return
	}

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	if err != nil {
		t.Errorf("an error occured while logging in by phone :%v", err)
		return
	}

	if login1 == nil {
		t.Errorf("nil response returned")
		return
	}

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// add the first secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	if err != nil {
		t.Errorf("failed to add secondary phonenumber :%v", err)
		return
	}

	userProfile, err := s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}

	// check if the length of secondary number == 1
	if len(userProfile.SecondaryPhoneNumbers) != 1 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryPhoneNumbers),
		)
		return
	}

	// try adding secondaryPhone1 again. this should fail because secondaryPhone1 already exists
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

	// add the second secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	if err != nil {
		t.Errorf("failed to add secondary phonenumber :%v", err)
		return
	}

	userProfile, err = s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}

	// check if the length of secondary number == 2
	if len(userProfile.SecondaryPhoneNumbers) != 2 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryPhoneNumbers),
		)
		return
	}

	// try adding secondaryPhone2 again. this should fail because secondaryPhone2 already exists
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

	// add the third secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	if err != nil {
		t.Errorf("failed to add secondary phonenumber :%v", err)
		return
	}

	userProfile, err = s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}

	// check if the length of secondary number == 3
	if len(userProfile.SecondaryPhoneNumbers) != 3 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryPhoneNumbers),
		)
		return
	}

	// try adding secondaryPhone3 again. this should fail because secondaryPhone3 already exists
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

	// try to login with each secondary phone number. This should fail
	login2, err := s.Login.LoginByPhone(context.Background(), secondaryPhone1, pin, base.FlavourConsumer)
	if err == nil {
		t.Errorf("an error %v was expected ", err)
		return
	}

	if login2 != nil {
		t.Errorf("an unexpected error occured :%v", err)
	}

	login3, err := s.Login.LoginByPhone(context.Background(), secondaryPhone2, pin, base.FlavourConsumer)
	if err == nil {
		t.Errorf("an error %v was expected ", err)
		return
	}

	if login3 != nil {
		t.Errorf("an unexpected error occured :%v", err)
	}

	login4, err := s.Login.LoginByPhone(context.Background(), secondaryPhone3, pin, base.FlavourConsumer)
	if err == nil {
		t.Errorf("an error %v was expected ", err)
		return
	}

	if login4 != nil {
		t.Errorf("an unexpected error occured :%v", err)
	}
}

func TestAddSecondaryEmailAddress(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	primaryEmail := "primary@example.com"
	secondaryemail1 := "user1@gmail.com"
	secondaryemail2 := "user2@gmail.com"
	secondaryemail3 := "user3@gmail.com"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("failed to create a user by phone")
		return
	}

	if resp == nil {
		t.Error("nil user response returned")
		return
	}

	if resp.Profile == nil {
		t.Error("nil profile response returned")
		return
	}

	if resp.CustomerProfile == nil {
		t.Error("nil customer profile response returned")
		return
	}

	if resp.SupplierProfile == nil {
		t.Error("nil supplier profile response returned")
		return
	}

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	if err != nil {
		t.Errorf("an error occured while logging in by phone :%v", err)
		return
	}

	if login1 == nil {
		t.Errorf("nil response returned")
		return
	}

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// try adding a secondary email address. This should fail because the profile does not have a primary email first
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	if err == nil {
		t.Errorf("expected an error: %v", err)
		return
	}

	// add the profile's primary email address. This is necessary. primary email must first exist before adding secondary emails
	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	if err != nil {
		t.Errorf("failed to add a primary email: %v", err)
		return
	}

	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	if err != nil {
		t.Errorf("failed to add secondary email: %v", err)
		return
	}

	userProfile, err := s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}
	// check if the length of secondary email == 1
	if len(userProfile.SecondaryEmailAddresses) != 1 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryEmailAddresses),
		)
		return
	}

	// try adding secondaryemail1 again since secondaryemail1 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

	// now add secondaryemail2
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail2})
	if err != nil {
		t.Errorf("failed to add secondary email: %v", err)
		return
	}

	userProfile, err = s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}
	// check if the length of secondary email == 2
	if len(userProfile.SecondaryEmailAddresses) != 2 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryEmailAddresses),
		)
		return
	}

	// try adding secondaryemail2 again since secondaryemail1 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail2})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

	// now add secondaryemail3
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail3})
	if err != nil {
		t.Errorf("failed to add secondary email: %v", err)
		return
	}

	userProfile, err = s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("failed to retrieve the profile of the logged in user :%v", err)
		return
	}

	if userProfile == nil {
		t.Errorf("nil response returned")
		return
	}
	// check if the length of secondary email == 3
	if len(userProfile.SecondaryEmailAddresses) != 3 {
		t.Errorf("expected the value to be equal to %v",
			len(userProfile.SecondaryEmailAddresses),
		)
		return
	}
	// try adding secondaryemail3 again since secondaryemail3 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail3})
	if err == nil {
		t.Errorf("an error %v was expected", err)
		return
	}

}

func TestUpdateUserProfilePushTokens(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
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

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	err = s.Onboarding.UpdatePushTokens(context.Background(), "token1", false)
	assert.NotNil(t, err)

	err = s.Onboarding.UpdatePushTokens(authenticatedContext, "token1", false)
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))

	err = s.Onboarding.UpdatePushTokens(authenticatedContext, "token2", false)
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.PushTokens))

	err = s.Onboarding.UpdatePushTokens(authenticatedContext, "token3", false)
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.PushTokens))

	// remove the token and assert new length
	err = s.Onboarding.UpdatePushTokens(context.Background(), "token2", true)
	assert.NotNil(t, err)

	err = s.Onboarding.UpdatePushTokens(authenticatedContext, "token2", true)
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.PushTokens))

	err = s.Onboarding.UpdatePushTokens(authenticatedContext, "token1", true)
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.PushTokens))
}

func TestCheckPhoneExists(t *testing.T) {
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

func TestGetUserProfileByUID(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
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

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the user profile using UID
	pr, err := s.Onboarding.GetUserProfileByUID(authenticatedContext, login1.Auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, login1.Profile.ID, pr.ID)
	assert.Equal(t, login1.Profile.UserName, pr.UserName)

	// now fetch using an authenticated context. should not fail
	pr2, err := s.Onboarding.GetUserProfileByUID(context.Background(), login1.Auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, pr2)
	assert.Equal(t, login1.Profile.ID, pr2.ID)
	assert.Equal(t, login1.Profile.UserName, pr2.UserName)
}

func TestUserProfile(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
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

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the user profile using authenticated context
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, login1.Profile.ID, pr.ID)
	assert.Equal(t, login1.Profile.UserName, pr.UserName)

	// now fetch using an unauthenticated context. should fail
	pr2, err := s.Onboarding.UserProfile(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, pr2)

}

func TestGetProfileByID(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	pin := "1234"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
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

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the user profile using ID
	pr, err := s.Onboarding.GetProfileByID(authenticatedContext, &login1.Profile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, login1.Profile.ID, pr.ID)
	assert.Equal(t, login1.Profile.UserName, pr.UserName)

	// now fetch using an authenticated context. should not fail
	pr2, err := s.Onboarding.GetProfileByID(context.Background(), &login1.Profile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pr2)
	assert.Equal(t, login1.Profile.ID, pr2.ID)
	assert.Equal(t, login1.Profile.UserName, pr2.UserName)

}

func TestUpdateBioData(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := base.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
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
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	dateOfBirth1 := base.Date{
		Day:   12,
		Year:  1998,
		Month: 2,
	}
	dateOfBirth2 := base.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}

	firstName1 := "makmende1"
	lastName1 := "Omera1"
	firstName2 := "makmende2"
	lastName2 := "Omera2"

	justDOB := base.BioData{
		DateOfBirth: &dateOfBirth1,
	}

	justFirstName := base.BioData{
		FirstName: &firstName1,
	}

	justLastName := base.BioData{
		LastName: &lastName1,
	}

	completeUserDetails := base.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update just the date of birth
	err = s.Onboarding.UpdateBioData(authenticatedContext, justDOB)
	assert.Nil(t, err)

	// fetch and assert dob has been updated
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, justDOB.DateOfBirth, pr.UserBioData.DateOfBirth)

	// update just the firstname
	err = s.Onboarding.UpdateBioData(authenticatedContext, justFirstName)
	assert.Nil(t, err)

	// fetch and assert firstname has been updated
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, justFirstName.FirstName, pr.UserBioData.FirstName)

	// update just the lastname
	err = s.Onboarding.UpdateBioData(authenticatedContext, justLastName)
	assert.Nil(t, err)

	// fetch and assert firstname has been updated
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, justLastName.LastName, pr.UserBioData.LastName)

	// update with the entire update input
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// fetch and assert dob, lastname & firstname have been updated
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, completeUserDetails.DateOfBirth, pr.UserBioData.DateOfBirth)
	assert.Equal(t, completeUserDetails.LastName, pr.UserBioData.LastName)
	assert.Equal(t, completeUserDetails.FirstName, pr.UserBioData.FirstName)

	assert.NotEqual(t, justDOB.DateOfBirth, pr.UserBioData.DateOfBirth)
	assert.NotEqual(t, justFirstName.FirstName, pr.UserBioData.LastName)
	assert.NotEqual(t, justLastName.LastName, pr.UserBioData.FirstName)

	// try update with an invalid context
	err = s.Onboarding.UpdateBioData(context.Background(), completeUserDetails)
	assert.NotNil(t, err)

}

func TestUpdatePhotoUploadID(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := base.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
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
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	uploadID1 := "photo-url1"
	uploadID2 := "photo-url2"

	err = s.Onboarding.UpdatePhotoUploadID(authenticatedContext, uploadID1)
	assert.Nil(t, err)

	// fetch and assert firstname has been updated
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, uploadID1, pr.PhotoUploadID)
	assert.NotEqual(t, resp.Profile.PhotoUploadID, pr.PhotoUploadID)

	err = s.Onboarding.UpdatePhotoUploadID(authenticatedContext, uploadID2)
	assert.Nil(t, err)

	// fetch and assert firstname has been updated again
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, uploadID2, pr.PhotoUploadID)
	assert.NotEqual(t, resp.Profile.PhotoUploadID, pr.PhotoUploadID)
	assert.NotEqual(t, uploadID1, pr.PhotoUploadID)
}

func TestUpdateSuspended(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := base.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
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
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert suspended
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, false, pr.Suspended)

	// now suspend the profile
	err = s.Onboarding.UpdateSuspended(authenticatedContext, true, *pr.PrimaryPhone, true)
	assert.Nil(t, err)

	// fetch the profile. this should fail because the profile has been suspended
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestUpdatePermissions(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := base.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
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
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []base.PermissionType{base.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	// use unauthenticated context. should fail
	err = s.Onboarding.UpdatePermissions(context.Background(), perms)
	assert.NotNil(t, err)

	pr, err = s.Onboarding.UserProfile(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, pr)
}

func TestSetupAsExperimentParticipant(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := "1234"

	validFlavourConsumer := base.FlavourConsumer

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), validPhoneNumber)

	// send otp to the phone number to initiate registration process
	otp, err := generateTestOTP(t, validPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, otp)

	// this should pass
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
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
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)
	// check that the currently created user can not experiment on new features
	assert.Equal(t, false, resp.Auth.CanExperiment)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	s, _ = InitializeTestService(authenticatedContext)

	// now add the user as an experiment participant
	input := true
	status, err := s.Onboarding.SetupAsExperimentParticipant(authenticatedContext, &input)
	assert.Nil(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, true, status)

	// try to add the user as an experiment participant. This should return the the same respones since th method internally is idempotent
	status, err = s.Onboarding.SetupAsExperimentParticipant(authenticatedContext, &input)
	assert.Nil(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, true, status)

	// login the user and assert they can experiment on new features
	login1, err := s.Login.LoginByPhone(context.Background(), validPhoneNumber, validPIN, validFlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.Equal(t, true, login1.Auth.CanExperiment)

	// now remove the user as an experiment participant
	input2 := false
	status, err = s.Onboarding.SetupAsExperimentParticipant(authenticatedContext, &input2)
	assert.Nil(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, true, status)

	// try removing the user as an experiment participant.This should return the the same respones since th method internally is idempotent
	status, err = s.Onboarding.SetupAsExperimentParticipant(authenticatedContext, &input2)
	assert.Nil(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, true, status)

	// login the user and assert they can not experiment on new features
	login2, err := s.Login.LoginByPhone(context.Background(), validPhoneNumber, validPIN, validFlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.Equal(t, false, login2.Auth.CanExperiment)
}

func TestMaskPhoneNumbers(t *testing.T) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		phones []string
	}

	tests := []struct {
		name string
		arg  args
		want []string
	}{
		{
			name: "valid case",
			arg: args{
				phones: []string{"+254789874267"},
			},
			want: []string{"+254789***267"},
		},
		{
			name: "valid case < 10 digits",
			arg: args{
				phones: []string{"+2547898742"},
			},
			want: []string{"+2547***742"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maskedPhone := s.Onboarding.MaskPhoneNumbers(tt.arg.phones)
			if len(maskedPhone) != len(tt.want) {
				t.Errorf("returned masked phone number not the expected one, wanted: %v got: %v", tt.want, maskedPhone)
				return
			}

			for i, number := range maskedPhone {
				if tt.want[i] != number {
					t.Errorf("wanted: %v, got: %v", tt.want[i], number)
					return
				}
			}
		})
	}
}

func TestAddAddress(t *testing.T) {
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

	addr := resources.UserAddressInput{
		Latitude:  1.2,
		Longitude: -34.001,
	}

	type args struct {
		ctx         context.Context
		input       resources.UserAddressInput
		addressType base.AddressType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:) add home address",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeHome,
			},
			wantErr: false,
		},
		{
			name: "happy:) add work address",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: false,
		},
		{
			name: "sad:( failed to get logged in user",
			args: args{
				ctx:         context.Background(),
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Onboarding.AddAddress(tt.args.ctx, tt.args.input, tt.args.addressType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.AddAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAddresses(t *testing.T) {
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

	addr := resources.UserAddressInput{
		Latitude:  1.2,
		Longitude: -34.001,
	}

	_, err = s.Onboarding.AddAddress(ctx, addr, base.AddressTypeWork)
	if err != nil {
		t.Errorf("unable to add test address")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "happy:) get addresses",
			args:    args{ctx: ctx},
			wantErr: false,
		},
		{
			name:    "sad:( failed to get logged in user",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Onboarding.GetAddresses(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestIntegrationGetAddresses(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup profile usecase")
	}

	validPhoneNumber := base.TestUserPhoneNumber
	validPIN := base.TestUserPin
	validFlavourConsumer := base.FlavourConsumer

	_ = s.Signup.RemoveUserByPhoneNumber(
		context.Background(),
		validPhoneNumber,
	)

	otp, err := generateTestOTP(t, validPhoneNumber)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}

	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &validPhoneNumber,
			PIN:         &validPIN,
			Flavour:     validFlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}
	if resp.Profile.HomeAddress != nil {
		t.Errorf("did not expect an address")
		return
	}
	if resp.Profile.WorkAddress != nil {
		t.Errorf("did not expect an address")
		return
	}

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: resp.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)

	lat := -1.2
	long := 34.56

	addr, err := s.Onboarding.AddAddress(
		authenticatedContext,
		resources.UserAddressInput{
			Latitude:  lat,
			Longitude: long,
		},
		base.AddressTypeHome,
	)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}
	if addr == nil {
		t.Errorf("expected an address")
		return
	}

	addrLat := addr.Latitude
	addrLong := addr.Longitude

	if addrLat != fmt.Sprintf("%.15f", lat) {
		t.Errorf("got a wrong address Latitude")
		return
	}
	if addrLong != fmt.Sprintf("%.15f", long) {
		t.Errorf("got a wrong address Longitude")
		return
	}

	profile, err := s.Onboarding.UserProfile(authenticatedContext)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}
	if profile == nil {
		t.Errorf("expected a user profile")
		return
	}

	if profile.HomeAddress == nil {
		t.Errorf("we expected an address")
		return
	}

	err = s.Signup.RemoveUserByPhoneNumber(
		authenticatedContext,
		validPhoneNumber,
	)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}
}

func TestProfileUseCaseImpl_UpdateCovers(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	p, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	memberNumber := uuid.New().String()
	cover := base.Cover{
		PayerName:      *utils.GetRandomName(),
		PayerSladeCode: 123,
		MemberNumber:   memberNumber,
		MemberName:     *utils.GetRandomName(),
	}

	type args struct {
		ctx    context.Context
		covers []base.Cover
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case:) update covers",
			args: args{
				ctx:    ctx,
				covers: []base.Cover{cover},
			},
			wantErr: false,
		},
		{
			name: "happy case:) update the same covers",
			args: args{
				ctx:    ctx,
				covers: []base.Cover{cover},
			},
			wantErr: false,
		},
		{
			name: "sad case:( unauthenticated context",
			args: args{
				ctx:    context.Background(),
				covers: []base.Cover{cover},
			},
			wantErr: true,
		},
		{
			name: "sad case:( update the nil covers",
			args: args{
				ctx:    ctx,
				covers: []base.Cover{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Onboarding.UpdateCovers(
				tt.args.ctx,
				tt.args.covers,
			); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdateCovers() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			profile, err := p.Onboarding.UserProfile(ctx)
			if err != nil {
				t.Errorf("unable to get user profile")
				return
			}

			covers := profile.Covers
			if len(covers) > 1 {
				t.Errorf("expected just one cover")
				return
			}
		})
	}
}

func TestRetireSecondaryPhoneNumbers(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	p, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		phoneNumbers []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad :( unable to get the user profile",
			args: args{
				phoneNumbers: []string{uuid.New().String()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad :( profile with no secondary phonenumbers",
			args: args{
				phoneNumbers: []string{uuid.New().String()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad :( adding an already existent phone number",
			args: args{
				phoneNumbers: []string{base.TestUserPhoneNumber},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy :) retire secondary phone numbers",
			args: args{
				phoneNumbers: []string{"+254700000003", "+254700000001"},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad :( unable to get the user profile" {
				got, err := p.Onboarding.RetireSecondaryPhoneNumbers(context.Background(), tt.args.phoneNumbers)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "sad :( profile with no secondary phonenumbers" {
				got, err := p.Onboarding.RetireSecondaryPhoneNumbers(ctx, tt.args.phoneNumbers)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "sad :( adding an already existent phone number" {
				err := p.Onboarding.UpdateSecondaryPhoneNumbers(ctx, []string{"+254700000001"})
				if err != nil {
					t.Errorf("unable to add secondary phone numbers: %v", err)
					return
				}

				got, err := p.Onboarding.RetireSecondaryPhoneNumbers(ctx, tt.args.phoneNumbers)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() = %v, want %v", got, tt.want)
				}
				profile, err := p.Onboarding.UserProfile(ctx)
				if err != nil {
					t.Errorf("unable to get user profile")
					return
				}
				if len(profile.SecondaryPhoneNumbers) > 1 {
					t.Errorf("expected 1 secondary phone numbers")
					return
				}
			}

			if tt.name == "happy :) retire secondary phone numbers" {
				err := p.Onboarding.UpdateSecondaryPhoneNumbers(ctx, []string{"+254700000003"})
				if err != nil {
					t.Errorf("unable to add secondary phone numbers: %v", err)
					return
				}

				got, err := p.Onboarding.RetireSecondaryPhoneNumbers(ctx, tt.args.phoneNumbers)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryPhoneNumbers() = %v, want %v", got, tt.want)
				}
				profile, err := p.Onboarding.UserProfile(ctx)
				if err != nil {
					t.Errorf("unable to get user profile")
					return
				}

				if len(profile.SecondaryPhoneNumbers) > 0 {
					t.Errorf("expected 0 secondary phone numbers but got: %v", len(profile.SecondaryPhoneNumbers))
					return
				}
			}
		})
	}
}

func TestRetireSecondaryEmailAddress(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	p, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	testEmail := "randommail@gmail.com"
	type args struct {
		emailAddresses []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad :( unable to get the user profile",
			args: args{
				emailAddresses: []string{base.GenerateRandomEmail()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad :( profile with no secondary email addresses",
			args: args{
				emailAddresses: []string{base.GenerateRandomEmail()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad :( adding an already existent email addresses",
			args: args{
				emailAddresses: []string{base.TestUserEmail},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy :) retire secondary email addresses",
			args: args{
				emailAddresses: []string{testEmail},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad :( unable to get the user profile" {
				got, err := p.Onboarding.RetireSecondaryEmailAddress(context.Background(), tt.args.emailAddresses)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "sad :( profile with no secondary email addresses" {
				got, err := p.Onboarding.RetireSecondaryEmailAddress(ctx, tt.args.emailAddresses)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "sad :( adding an already existent email addresses" {
				err := p.Onboarding.UpdatePrimaryEmailAddress(ctx, base.GenerateRandomEmail())
				if err != nil {
					t.Errorf("unable to set primary email address: %v", err)
					return
				}

				got, err := p.Onboarding.RetireSecondaryEmailAddress(ctx, tt.args.emailAddresses)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() = %v, want %v", got, tt.want)
				}
				profile, err := p.Onboarding.UserProfile(ctx)
				if err != nil {
					t.Errorf("unable to get user profile")
					return
				}
				if len(profile.SecondaryEmailAddresses) > 1 {
					t.Errorf("expected 1 secondary email addresses")
					return
				}
			}

			if tt.name == "happy :) retire secondary email addresses" {
				profile, err := p.Onboarding.UserProfile(ctx)
				if err != nil {
					t.Errorf("unable to get user profile")
					return
				}

				err = p.Onboarding.UpdatePrimaryEmailAddress(ctx, base.TestUserEmail)
				if err != nil {
					t.Errorf("unable to set primary email address: %v", err)
					return
				}

				time.Sleep(2 * time.Second)

				err = p.Onboarding.UpdateSecondaryEmailAddresses(ctx, []string{testEmail})
				if err != nil {
					t.Errorf("unable to set secondary email address: %v", err)
					return
				}

				got, err := p.Onboarding.RetireSecondaryEmailAddress(ctx, tt.args.emailAddresses)
				if (err != nil) != tt.wantErr {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ProfileUseCaseImpl.RetireSecondaryEmailAddress() = %v, want %v", got, tt.want)
				}
				if len(profile.SecondaryEmailAddresses) > 0 {
					t.Errorf("expected 0 secondary email addresses but got: %v", len(profile.SecondaryEmailAddresses))
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_RemoveAdminPermsToUser(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	p, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	phoneNumber := base.TestUserPhoneNumber
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup profile usecase")
	}

	_ = s.Signup.RemoveUserByPhoneNumber(
		context.Background(),
		phoneNumber,
	)
	phoneNumberWithNoUserProfile := "+2547898742"
	otp, err := generateTestOTP(t, phoneNumber)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	_, err = p.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &phoneNumber,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("failed to create a user by phone")
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
			name: "happy case:) remove admin permissions ",
			args: args{
				ctx:   ctx,
				phone: phoneNumber,
			},
			wantErr: false,
		},
		{
			name: "sade case:) remove admin permissions",
			args: args{
				ctx:   ctx,
				phone: phoneNumberWithNoUserProfile,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Onboarding.RemoveAdminPermsToUser(tt.args.ctx, tt.args.phone); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.RemoveAdminPermsToUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
