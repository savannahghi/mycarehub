package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestPinErrors(t *testing.T) {
	err := SaveUserPinError(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = InvalidResetPinPayloadErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = PinNotFoundError(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = PinMismatchError()
	assert.NotNil(t, err)

	err = InvalidatePinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ResetPinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = PINExpiredErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = PINErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GenerateTempPINErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ExpiredPinErr()
	assert.NotNil(t, err)

	err = GeneratePinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ValidatePINDigitsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ExistingPINError(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
