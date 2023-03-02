package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// AuthorityRole defines user roles
type AuthorityRole struct {
	AuthorityRoleID string                `json:"authorityRoleID"`
	Name            string                `json:"name"`
	Active          bool                  `json:"active"`
	OrganisationID  string                `json:"organisationID"`
	ProgramID       string                `json:"programID"`
	Permissions     []AuthorityPermission `json:"permissions"`
}

// AuthorityPermission defines user permissions
type AuthorityPermission struct {
	PermissionID string               `json:"permissionID"`
	Active       bool                 `json:"active"`
	Name         enums.PermissionType `json:"name"`
	Description  string               `json:"description"`
	Category     string               `json:"category"`
	Scope        string               `json:"scope"`
}
