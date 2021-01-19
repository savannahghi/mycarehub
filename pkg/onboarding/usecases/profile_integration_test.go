package usecases_test

import (
	"context"
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
