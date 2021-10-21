package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
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
