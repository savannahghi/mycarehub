package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestAuthorizationErrors(t *testing.T) {
	err := CheckUserRoleErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = UserNotAuthorizedErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = CheckUserPermissionErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = AssignRolesErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GetUserRolesErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GetUserPermissionsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = RevokeRolesErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = GetAllRolesErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
