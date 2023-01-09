package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// AuthorityRole defines user roles
type AuthorityRole struct {
	AuthorityRoleID string             `json:"authorityRoleID"`
	Name            enums.UserRoleType `json:"name"`
	Active          bool               `json:"active"`
}

// AuthorityPermission defines user permissions
type AuthorityPermission struct {
	PermissionID string               `json:"permissionID"`
	Name         enums.PermissionType `json:"name"`
	Active       bool                 `json:"active"`
}
