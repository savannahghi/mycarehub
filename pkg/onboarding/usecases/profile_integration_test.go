package usecases_test

import (
	"context"
	"log"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
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

	// try to login with secondaryPhone. This should fail because secondaryPhone != primaryPhone
	login2, err := s.Login.LoginByPhone(context.Background(), secondaryPhone, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login2)

	// add a secondary phone number to the user
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))

	// login to add assert the secondary phone number has been added
	login3, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login3)
	assert.Equal(t, 1, len(login3.Profile.SecondaryPhoneNumbers))

	// send otp to the secondary phone number we intend to make primary
	otpResp, err := s.Otp.GenerateAndSendOTP(context.Background(), secondaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otpResp)

	// set the old secondary phone number as the new primary phone number
	setResp, err := s.Signup.SetPhoneAsPrimary(context.Background(), secondaryPhone, otpResp.OTP)
	assert.Nil(t, err)
	assert.NotNil(t, setResp)

	// login with the old primary phone number. This should fail
	login4, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login4)

	// login with the new primary phone number. This should not fail. Assert that the primary phone number
	// is the new one and the secondary phone slice contains the old primary phone number.
	login5, err := s.Login.LoginByPhone(context.Background(), secondaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login5)
	assert.Equal(t, secondaryPhone, *login5.Profile.PrimaryPhone)
	assert.Contains(t, login5.Profile.SecondaryPhoneNumbers, secondaryPhone)

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

	// add the first secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))

	// try adding secondaryPhone1 again. this should fail because secondaryPhone1 already exists
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	assert.NotNil(t, err)

	// add the second secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryPhoneNumbers))

	// try adding secondaryPhone2 again. this should fail because secondaryPhone2 already exists
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	assert.NotNil(t, err)

	// add the third secondary phone number
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.SecondaryPhoneNumbers))

	// try to login with each secondary phone number. This should fail
	login2, err := s.Login.LoginByPhone(context.Background(), secondaryPhone1, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login2)

	login3, err := s.Login.LoginByPhone(context.Background(), secondaryPhone2, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login3)

	login4, err := s.Login.LoginByPhone(context.Background(), secondaryPhone3, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login4)
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

	// try adding a secondary email address. This should fail because the profile does not have a primary email first
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	assert.NotNil(t, err)

	// add the profile's primary email address. This is necessary. primary email must first exist before adding secondary emails
	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryEmailAddresses))

	// try adding secondaryemail1 again since secondaryemail1 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	assert.NotNil(t, err)

	// now add secondaryemail2
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail2})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryEmailAddresses))

	// try adding secondaryemail2 again since secondaryemail1 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail2})
	assert.NotNil(t, err)

	// now add secondaryemail3
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail3})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.SecondaryEmailAddresses))

	// try adding secondaryemail3 again since secondaryemail3 is already in use
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail3})
	assert.NotNil(t, err)

}

func TestUpdateUserProfileCovers(t *testing.T) {
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

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer1", PayerSladeCode: 1, MemberName: "name1", MemberNumber: "mem1"}})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer2", PayerSladeCode: 2, MemberName: "name2", MemberNumber: "mem2"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.Covers))
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

func TestRetireSecondaryPhoneNumbers(t *testing.T) {
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

	// try to retire a secondary phone number. This should fail since we have not added a secondary phone number yet
	rm1, err := s.Onboarding.RetireSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	assert.NotNil(t, err)
	assert.NotNil(t, rm1)
	assert.Equal(t, false, rm1)

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryPhoneNumbers))

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.SecondaryPhoneNumbers))

	// remove secondaryPhone3 and assert the length of new secondary phone numbers slice
	rm2, err := s.Onboarding.RetireSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	assert.Nil(t, err)
	assert.NotNil(t, rm2)
	assert.Equal(t, true, rm2)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryPhoneNumbers))
	assert.Equal(t, false, base.StringSliceContains(pr.SecondaryPhoneNumbers, secondaryPhone3))

	// remove secondaryPhone2 and assert the length of new secondary phone numbers slice
	rm3, err := s.Onboarding.RetireSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	assert.Nil(t, err)
	assert.NotNil(t, rm3)
	assert.Equal(t, true, rm3)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))
	assert.Equal(t, false, base.StringSliceContains(pr.SecondaryPhoneNumbers, secondaryPhone2))
}

func TestRetireSecondaryEmailAddress(t *testing.T) {
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

	// try adding a secondary email. this should fail because a primary email has not been added prior
	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	assert.NotNil(t, err)

	// now add the profile's primary email address. This is necessary. primary email must first exist before adding secondary emails
	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	// try to retire a secondary email. This should fail since we have not added a secondary email address yet
	rm1, err := s.Onboarding.RetireSecondaryEmailAddress(authenticatedContext, []string{secondaryemail1})
	assert.NotNil(t, err)
	assert.NotNil(t, rm1)
	assert.Equal(t, false, rm1)

	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail1})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryEmailAddresses))

	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail2})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryEmailAddresses))

	err = s.Onboarding.UpdateSecondaryEmailAddresses(authenticatedContext, []string{secondaryemail3})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.SecondaryEmailAddresses))

	// remove secondaryemail3 and assert the length of new secondary phone numbers slice
	rm2, err := s.Onboarding.RetireSecondaryEmailAddress(authenticatedContext, []string{secondaryemail3})
	assert.Nil(t, err)
	assert.NotNil(t, rm2)
	assert.Equal(t, true, rm2)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryEmailAddresses))
	assert.Equal(t, false, base.StringSliceContains(pr.SecondaryEmailAddresses, secondaryemail3))

	// remove secondaryemail2 and assert the length of new secondary phone numbers slice
	rm3, err := s.Onboarding.RetireSecondaryEmailAddress(authenticatedContext, []string{secondaryemail2})
	assert.Nil(t, err)
	assert.NotNil(t, rm3)
	assert.Equal(t, true, rm3)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryEmailAddresses))
	assert.Equal(t, false, base.StringSliceContains(pr.SecondaryPhoneNumbers, secondaryemail2))
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

func TestProfileUseCaseImpl_GetUserProfileByUID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
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
		ctx context.Context
		UID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: get a user profile given their UID",
			args: args{
				ctx: ctx,
				UID: auth.UID,
			},
			wantErr: false,
		},
		{
			name: "failure: fail get a user profile given a bad UID",
			args: args{
				ctx: ctx,
				UID: "not-a-valid-uid",
			},
			wantErr: true,
		},
		{
			name: "failure: fail get a user profile given an empty UID",
			args: args{
				ctx: ctx,
				UID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := s.Onboarding.GetUserProfileByUID(tt.args.ctx, tt.args.UID)
			if tt.wantErr && profile != nil {
				t.Errorf("expected nil but got %v, since the error %v occurred",
					profile,
					err,
				)
				return
			}

			if !tt.wantErr && profile == nil {
				t.Errorf("expected a profile but got nil, since no error occurred")
				return
			}

		})
	}
}

func TestProfileUseCaseImpl_UserProfile(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid: user profile retrieved",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context supplied",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.UserProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_GetProfileByID(t *testing.T) {

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	profile, err := s.Onboarding.UserProfile(ctx)
	if err != nil {
		t.Errorf("could not retrieve user profile")
		return
	}

	type args struct {
		ctx context.Context
		id  *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: user profile retrieved",
			args: args{
				ctx: ctx,
				id:  &profile.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.GetProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateBioData(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
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
	bioData := base.BioData{
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
	}

	var gender base.Gender = "female"
	updateGender := base.BioData{
		Gender: gender,
	}

	updateDOB := base.BioData{
		DateOfBirth: &dateOfBirth,
	}

	updateFirstName := base.BioData{
		FirstName: &firstName,
	}

	updateLastName := base.BioData{
		LastName: &lastName,
	}

	type args struct {
		ctx  context.Context
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update biodata",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the firstname",
			args: args{
				ctx:  ctx,
				data: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the lastname",
			args: args{
				ctx:  ctx,
				data: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the date of birth",
			args: args{
				ctx:  ctx,
				data: updateDOB,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the gender",
			args: args{
				ctx:  ctx,
				data: updateGender,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unauthenticated context",
			args: args{
				ctx:  context.Background(),
				data: bioData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdateBioData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePhotoUploadID(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		t.Errorf("could not get the logged in user")
		return
	}

	profile, err := s.Onboarding.GetUserProfileByUID(ctx, uid)
	if err != nil {
		t.Errorf("could not retrieve user profile")
		return
	}

	uploadID := "some-photo-upload-id"
	log.Printf("THE UPLOAD ID IS %v", profile.PhotoUploadID)

	type args struct {
		ctx      context.Context
		uploadID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update the PhotoUploadID",
			args: args{
				ctx:      ctx,
				uploadID: uploadID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an unauthenticated context",
			args: args{
				ctx:      context.Background(),
				uploadID: uploadID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdatePhotoUploadID(tt.args.ctx, tt.args.uploadID); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdatePhotoUploadID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
