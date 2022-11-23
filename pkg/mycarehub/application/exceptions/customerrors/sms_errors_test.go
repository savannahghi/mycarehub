package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestSMSErrors(t *testing.T) {
	err := SendSMSErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
