package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestProgramErrors(t *testing.T) {
	err := OrgIDForProgramExistErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = CreateProgramErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
