package user_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func TestUseCasesUserImpl_SetUserPIN_Integration(t *testing.T) {
	ctx := context.Background()

	m := testInfrastructureInteractor

	validPINInput := &dto.PINInput{
		PIN:          "1234",
		ConfirmedPin: "1234",
		Flavour:      feedlib.FlavourConsumer,
	}

	invalidPINInput := &dto.PINInput{
		PIN:          "12",
		ConfirmedPin: "1234",
		Flavour:      "CONSUMER",
	}

	invalidPINInput2 := &dto.PINInput{
		PIN:          "",
		ConfirmedPin: "",
		Flavour:      "CONSUMER",
	}

	//check for valid PIN
	err1 := utils.ValidatePIN(validPINInput.PIN)
	assert.Nil(t, err1)

	// check for invalid PIN
	err2 := utils.ValidatePIN(invalidPINInput.PIN)
	assert.NotNil(t, err2)

	// check for empty PIN
	err3 := utils.ValidatePIN(invalidPINInput2.PIN)
	assert.NotNil(t, err3)

	salt, encodedPIN := extension.EncryptPIN(validPINInput.PIN, nil)

	isMatch := extension.ComparePIN(validPINInput.PIN, salt, encodedPIN, nil)

	pinDataInput := &domain.UserPIN{
		UserID:    ksuid.New().String(),
		HashedPIN: encodedPIN,
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Flavour:   validPINInput.Flavour,
		IsValid:   isMatch,
		Salt:      salt,
	}

	isTrue, err := m.SetUserPIN(ctx, pinDataInput)
	assert.Nil(t, err)
	assert.NotNil(t, isTrue)
	assert.Equal(t, true, isTrue)

}

func TestUseCasesUserImpl_Login_Integration_Test(t *testing.T) {
	ctx := context.Background()

	m := testInfrastructureInteractor
	i := testInteractor

	flavour := feedlib.FlavourConsumer
	pin := "1234"

	facilityInput := dto.FacilityInput{
		Name:        "test",
		Code:        "c123",
		Active:      true,
		County:      "test",
		Description: "test description",
	}

	// Create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("Failed to create facility: %v", err)
	}

	userInput := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Flavour:     feedlib.FlavourConsumer,
	}

	staffInput := &dto.StaffProfileInput{
		StaffNumber:       "s123",
		DefaultFacilityID: facility.ID,
	}

	// Register user
	staffUserProfile, err := i.StaffUsecase.RegisterStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, err)
	assert.NotNil(t, staffUserProfile)

	// Set PIN
	salt, encodedPIN := extension.EncryptPIN(pin, nil)
	assert.NotNil(t, encodedPIN)
	assert.NotNil(t, salt)

	PINInput := &domain.UserPIN{
		UserID:    *staffUserProfile.User.ID,
		HashedPIN: encodedPIN,
		ValidFrom: time.Now(),
		ValidTo:   utils.GetHourMinuteSecond(24, 0, 0),
		Flavour:   flavour,
		IsValid:   true,
		Salt:      salt,
	}

	isSet, err := m.SetUserPIN(ctx, PINInput)
	assert.Nil(t, err)
	assert.Equal(t, true, isSet)

	// Valid userID
	userProfile, err := m.GetUserProfileByUserID(ctx, *staffUserProfile.User.ID, string(flavour))
	assert.Nil(t, err)
	assert.NotNil(t, userProfile)

	//Valid: Fetch PIN by UserID
	userPINData, err := m.GetUserPINByUserID(ctx, *staffUserProfile.User.ID)
	assert.Nil(t, err)
	assert.NotNil(t, userPINData)

	isMatch := extension.ComparePIN("1234", userPINData.Salt, userPINData.HashedPIN, nil)
	assert.Equal(t, true, isMatch)

	successTime := time.Now()
	err = m.UpdateUserLastSuccessfulLogin(ctx, *staffUserProfile.User.ID, successTime, string(flavour))
	assert.Nil(t, err)

	err = m.UpdateUserFailedLoginCount(ctx, *staffUserProfile.User.ID, "0", string(flavour))
	assert.Nil(t, err)

	customToken, err := firebasetools.CreateFirebaseCustomToken(ctx, *userProfile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, customToken)

	userTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, err)
	assert.NotNil(t, userTokens)

	//Login
	authCred, str, err := i.UserUsecase.Login(ctx, *staffUserProfile.User.ID, pin, flavour.String())
	assert.Nil(t, err)
	assert.NotEmpty(t, str)
	assert.NotNil(t, str)
	assert.NotNil(t, authCred)

	// Invalid
	invalidPIN1 := "4321"
	profile, err1 := m.GetUserProfileByUserID(ctx, *staffUserProfile.User.ID, string(flavour))
	assert.Nil(t, err1)
	assert.NotNil(t, profile)

	//Valid: Fetch PIN by UserID
	userPINData2, err2 := m.GetUserPINByUserID(ctx, *profile.ID)
	assert.Nil(t, err2)
	assert.NotNil(t, userPINData2)

	isMatch = extension.ComparePIN(invalidPIN1, userPINData.Salt, userPINData.HashedPIN, nil)
	assert.Equal(t, false, isMatch)

	err3 := m.UpdateUserFailedLoginCount(ctx, *profile.ID, "1", string(flavour))
	assert.Nil(t, err3)

	lastFailedLoginTime := time.Now()
	err4 := m.UpdateUserLastFailedLogin(ctx, *staffUserProfile.User.ID, lastFailedLoginTime, string(flavour))
	assert.Nil(t, err4)

	//Cannot Login
	authCred, str, err = i.UserUsecase.Login(ctx, *staffUserProfile.User.ID, invalidPIN1, flavour.String())
	assert.NotNil(t, err)
	assert.Empty(t, str)
	assert.Nil(t, authCred)

	invalidPIN2 := "4321"
	profile2, err5 := m.GetUserProfileByUserID(ctx, *staffUserProfile.User.ID, string(flavour))
	assert.Nil(t, err5)
	assert.NotNil(t, profile2)

	//Valid: Fetch PIN by UserID
	userPINData3, err6 := m.GetUserPINByUserID(ctx, *profile2.ID)
	assert.Nil(t, err6)
	assert.NotNil(t, userPINData3)

	isMatch2 := extension.ComparePIN(invalidPIN2, userPINData.Salt, userPINData.HashedPIN, nil)
	assert.Equal(t, false, isMatch)
	if !isMatch2 {
		failedLoginCount, err7 := strconv.Atoi(profile2.FailedLoginCount)
		assert.Nil(t, err7)
		assert.NotNil(t, failedLoginCount)
		trials := failedLoginCount + 1
		//Convert trials to string
		numberOfTrials := strconv.Itoa(trials)
		assert.NotNil(t, numberOfTrials)

		err8 := m.UpdateUserFailedLoginCount(ctx, *staffUserProfile.User.ID, numberOfTrials, string(flavour))
		assert.Nil(t, err8)

		lastFailedLoginTime := time.Now()
		err9 := m.UpdateUserLastFailedLogin(ctx, *staffUserProfile.User.ID, lastFailedLoginTime, string(flavour))
		assert.Nil(t, err9)

	}

	//Cannot Login
	authCred, str, err = i.UserUsecase.Login(ctx, *staffUserProfile.User.ID, invalidPIN2, flavour.String())
	assert.NotNil(t, err)
	assert.Empty(t, str)
	assert.Nil(t, authCred)

}
