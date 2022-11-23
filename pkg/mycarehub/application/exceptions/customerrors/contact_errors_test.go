package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestContactErrors(t *testing.T) {
	err := InvalidContactTypeErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NoContactsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ContactNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NormalizeMSISDNError(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
