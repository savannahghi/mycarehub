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

	// try adding the first cover again. This should add the cover because the cover already exists
	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer1", PayerSladeCode: 1, MemberName: "name1", MemberNumber: "mem1"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer2", PayerSladeCode: 2, MemberName: "name2", MemberNumber: "mem2"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.Covers))

	// try adding the second cover again. This should add the cover because the cover already exists
	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer2", PayerSladeCode: 2, MemberName: "name2", MemberNumber: "mem2"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer1", PayerSladeCode: 2, MemberName: "name11", MemberNumber: "mem22"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer1", PayerSladeCode: 2, MemberName: "name111", MemberNumber: "mem222"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 4, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer2", PayerSladeCode: 1, MemberName: "name2", MemberNumber: "mem2"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 5, len(pr.Covers))

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
