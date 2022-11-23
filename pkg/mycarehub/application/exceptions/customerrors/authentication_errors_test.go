package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestAuthenticationErrors(t *testing.T) {
	err := LoginCountUpdateErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = LoginTimeUpdateErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NexAllowedLoginTimeErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = RetryLoginErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = FailedSecurityCountExceededErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = SecurityQuestionResponseMismatchErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
}
