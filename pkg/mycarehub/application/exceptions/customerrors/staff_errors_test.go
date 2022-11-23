package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestStaffErrors(t *testing.T) {
	err := StaffProfileNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = StaffHasUnresolvedPinResetRequestErr()
	assert.NotNil(t, err)

}
