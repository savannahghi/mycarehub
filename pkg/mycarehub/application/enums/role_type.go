package enums

import (
	"fmt"
	"io"
	"strconv"
)

// RolesType is a list of all the role types.
type RolesType string

// contacts type constants
const (
	RolesTypeCanRegisterStaff   RolesType = "CAN_REGISTER_STAFF"
	RolesTypeCanInviteStaff     RolesType = "CAN_INVITE_STAFF"
	RolesTypeCanSuspendStaff    RolesType = "CAN_SUSPEND_STAFF"
	RolesTypeCanActivateStaff   RolesType = "CAN_ACTIVATE_STAFF"
	RolesTypeCanDeleteStaff     RolesType = "CAN_DELETE_STAFF"
	RolesTypeCanInactivateStaff RolesType = "CAN_INACTIVATE_STAFF"

	RolesTypeCanRegisterClient   RolesType = "CAN_REGISTER_CLIENT"
	RolesTypeCanInviteClient     RolesType = "CAN_INVITE_CLIENT"
	RolesTypeCanSuspendClient    RolesType = "CAN_SUSPEND_CLIENT"
	RolesTypeCanActivateClient   RolesType = "CAN_ACTIVATE_CLIENT"
	RolesTypeCanDeleteClient     RolesType = "CAN_DELETE_CLIENT"
	RolesTypeCanInactivateClient RolesType = "CAN_INACTIVATE_CLIENT"
)

// AllRoles is a set of a  valid and known role types.
var AllRoles = []RolesType{
	RolesTypeCanRegisterStaff,
	RolesTypeCanInviteStaff,
	RolesTypeCanSuspendStaff,
	RolesTypeCanActivateStaff,
	RolesTypeCanDeleteStaff,
	RolesTypeCanInactivateStaff,

	RolesTypeCanRegisterClient,
	RolesTypeCanInviteClient,
	RolesTypeCanSuspendClient,
	RolesTypeCanActivateClient,
	RolesTypeCanDeleteClient,
	RolesTypeCanInactivateClient,
}

// IsValid returns true if a role is valid
func (m RolesType) IsValid() bool {
	switch m {
	case RolesTypeCanRegisterStaff,
		RolesTypeCanInviteStaff,
		RolesTypeCanSuspendStaff,
		RolesTypeCanActivateStaff,
		RolesTypeCanDeleteStaff,
		RolesTypeCanInactivateStaff,

		RolesTypeCanRegisterClient,
		RolesTypeCanInviteClient,
		RolesTypeCanSuspendClient,
		RolesTypeCanActivateClient,
		RolesTypeCanDeleteClient,
		RolesTypeCanInactivateClient:
		return true
	}
	return false
}

// String converts roles to string.
func (m RolesType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a role type.
func (m *RolesType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = RolesType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid RolesType", str)
	}
	return nil
}

// MarshalGQL writes the role type to the supplied writer
func (m RolesType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
