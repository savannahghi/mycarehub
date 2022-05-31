package enums

import (
	"fmt"
	"io"
	"strconv"
)

// UserRoleType is a list of all the user roles.
type UserRoleType string

const (
	// UserRoleTypeSystemAdministrator defines a system admin role
	UserRoleTypeSystemAdministrator UserRoleType = "SYSTEM_ADMINISTRATOR"

	// UserRoleTypeCommunityManagement defines a community manager role
	UserRoleTypeCommunityManagement UserRoleType = "COMMUNITY_MANAGEMENT"

	// UserRoleTypeContentManagement defines a content manager role
	UserRoleTypeContentManagement UserRoleType = "CONTENT_MANAGEMENT"

	// UserRoleTypeClientManagement defines a client manager role
	UserRoleTypeClientManagement UserRoleType = "CLIENT_MANAGEMENT"
)

// AllRoles is a set of a  valid and known user roles.
var AllRoles = []UserRoleType{
	UserRoleTypeSystemAdministrator,
	UserRoleTypeCommunityManagement,
	UserRoleTypeContentManagement,
	UserRoleTypeClientManagement,
}

// IsValid returns true if a role is valid
func (m UserRoleType) IsValid() bool {
	switch m {
	case UserRoleTypeSystemAdministrator,
		UserRoleTypeCommunityManagement,
		UserRoleTypeContentManagement,
		UserRoleTypeClientManagement:
		return true
	}
	return false
}

// String converts role type to string
func (m UserRoleType) String() string {
	return string(m)
}

// Name is the human readable role name
func (m UserRoleType) Name() string {
	switch m {
	case UserRoleTypeSystemAdministrator:
		return "System Administrator"
	case UserRoleTypeCommunityManagement:
		return "Community Management"
	case UserRoleTypeContentManagement:
		return "Content Management"
	case UserRoleTypeClientManagement:
		return "Client Management"
	default:
		return m.String()
	}
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *UserRoleType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = UserRoleType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid UserRoleType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m UserRoleType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
