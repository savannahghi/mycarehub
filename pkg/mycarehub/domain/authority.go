package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// AuthorityRole defines user roles
type AuthorityRole struct {
	RoleID string             `json:"authorityRoleID"`
	Name   enums.UserRoleType `json:"name"`
}

// AuthorityPermission defines user permissions
type AuthorityPermission struct {
	PermissionID string               `json:"authorityPermissionID"`
	Name         enums.PermissionType `json:"name"`
}
