package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestGetInviteLinkErr(t *testing.T) {
	err := UserNotFoundError(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = EmptyUserIDErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ProfileNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NotOptedInErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GetLoggedInUserUIDErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = UserNameExistsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = UpdateProfileErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GetInviteLinkErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
