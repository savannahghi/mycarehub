package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestSecurityQuestionError(t *testing.T) {
	err := SecurityQuestionNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
}
