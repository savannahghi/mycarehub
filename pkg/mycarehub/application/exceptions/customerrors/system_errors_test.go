package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestSystemErrors(t *testing.T) {
	err := InternalErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = InvalidFlavourDefinedErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = EncryptionErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
}
