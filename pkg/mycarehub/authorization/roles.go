package authorization

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// DefaultRole defines various roles that are assigned to different users when they are created in the platform
type DefaultRole string

const (
	DefaultRoleAdmin     DefaultRole = "Default Admin"
	DefaultRoleClient    DefaultRole = "Default Client"
	DefaultRoleCaregiver DefaultRole = "Default Caregiver"
)

// IsValid checks if a string of type DefaultRole is of the valid type
func (p DefaultRole) IsValid() bool {
	switch p {
	case DefaultRoleAdmin,
		DefaultRoleClient,
		DefaultRoleCaregiver:
		return true
	}
	return false
}

// String converts DefaultRole type to type string
func (p DefaultRole) String() string {
	return string(p)
}

// UnmarshalGQL converts the supplied value to a DefaultRole type.
func (p *DefaultRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("value must be of type string")
	}

	*p = DefaultRole(str)
	if !p.IsValid() {
		return fmt.Errorf("%s is not a valid DefaultRole", str)
	}
	return nil
}

// MarshalGQL writes the DefaultRole type to the supplied
func (p DefaultRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(p.String()))
}

var (
	// DefaultAdminRole defines the default role assigned to an admin
	DefaultAdminRole = domain.AuthorityRole{
		Name:        DefaultRoleAdmin.String(),
		Permissions: AllPermissions(context.Background()),
	}

	// DefaultClientRole defines the default role assigned to a client
	DefaultClientRole = domain.AuthorityRole{
		Name:        DefaultRoleClient.String(),
		Permissions: DefaultClientPermissions(context.Background()),
	}
	// DefaultCaregiverRole defines the default role assigned to a caregiver
	DefaultCaregiverRole = domain.AuthorityRole{
		Name:        DefaultRoleCaregiver.String(),
		Permissions: DefaultCaregiverPermissions(context.Background()),
	}
)
