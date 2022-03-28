package authority

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICheckUserRole contains methods to check if a user has a given role
type ICheckUserRole interface {
	CheckUserRole(ctx context.Context, role enums.UserRoleType) error
}

// ICheckUserPermission contains methods to check if a user has a given permission
type ICheckUserPermission interface {
	CheckUserPermission(ctx context.Context, permission enums.PermissionType) error
}

// IManageRoles contains methods to manage user roles
type IManageRoles interface {
	AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	AssignOrRevokeRoles(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error)
}

// IGetRoles contains methods that get the roles
type IGetRoles interface {
	GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error)
}

// IGetPermissions contains methods that get the permissions
type IGetPermissions interface {
	GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
}

// UsecaseAuthority groups al the interfaces for the Authority usecase
type UsecaseAuthority interface {
	ICheckUserRole
	ICheckUserPermission
	IManageRoles
	IGetRoles
	IGetPermissions
}

// UsecaseAuthorityImpl represents the Authority implementation
type UsecaseAuthorityImpl struct {
	Query       infrastructure.Query
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
}

// NewUsecaseAuthority is the controller function for the Authority usecase
func NewUsecaseAuthority(
	query infrastructure.Query,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
) *UsecaseAuthorityImpl {
	return &UsecaseAuthorityImpl{
		Query:       query,
		Update:      update,
		ExternalExt: externalExt,
	}
}

// CheckUserRole checks if the user had the specified role
func (u *UsecaseAuthorityImpl) CheckUserRole(ctx context.Context, role enums.UserRoleType) error {
	if role == "" {
		err := fmt.Errorf("role must not be empty")
		helpers.ReportErrorToSentry(err)
		return exceptions.EmptyInputErr(err)
	}

	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return exceptions.GetLoggedInUserUIDErr(err)
	}

	ok, err := u.Query.CheckUserRole(ctx, loggedInUserID, role.String())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return exceptions.CheckUserRoleErr(err)
	}
	if !ok {
		err := fmt.Errorf("user is not authorized perform action, missing role: %s", role)
		helpers.ReportErrorToSentry(err)
		return exceptions.UserNotAuthorizedErr(err)
	}
	return nil
}

// CheckUserPermission checks if the user had the specified permission
func (u *UsecaseAuthorityImpl) CheckUserPermission(ctx context.Context, permission enums.PermissionType) error {
	if permission == "" {
		err := fmt.Errorf("permission must not be empty")
		helpers.ReportErrorToSentry(err)
		return exceptions.EmptyInputErr(err)
	}

	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return exceptions.GetLoggedInUserUIDErr(err)
	}

	ok, err := u.Query.CheckUserPermission(ctx, loggedInUserID, permission.String())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return exceptions.CheckUserPermissionErr(err)
	}
	if !ok {
		err := fmt.Errorf("user is not authorized perform action, missing permission: %s", permission)
		helpers.ReportErrorToSentry(err)
		return exceptions.UserNotAuthorizedErr(err)
	}
	return nil
}

// AssignRoles assigns the specified roles to the user
func (u *UsecaseAuthorityImpl) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	if userID == "" {
		err := fmt.Errorf("userID must not be empty")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)
	}
	if len(roles) == 0 {
		err := fmt.Errorf("roles must not be empty")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)
	}
	// check if user can assign role
	err := u.CheckUserPermission(ctx, enums.PermissionTypeCanEditUserRole)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotAuthorizedErr(err)
	}

	ok, err := u.Update.AssignRoles(ctx, userID, roles)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.AssignRolesErr(err)
	}
	return ok, nil
}

// GetUserRoles returns the roles of the user
func (u *UsecaseAuthorityImpl) GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
	if userID == "" {
		err := fmt.Errorf("userID must not be empty")
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.EmptyInputErr(err)
	}

	roles, err := u.Query.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetUserRolesErr(err)
	}
	return roles, nil
}

// GetUserPermissions returns the permissions of the user
func (u *UsecaseAuthorityImpl) GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
	if userID == "" {
		err := fmt.Errorf("userID must not be empty")
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.EmptyInputErr(err)
	}

	permissions, err := u.Query.GetUserPermissions(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetUserPermissionsErr(err)
	}
	return permissions, nil
}

// RevokeRoles revokes the specified roles from the user
func (u *UsecaseAuthorityImpl) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	if userID == "" {
		err := fmt.Errorf("userID must not be empty")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)
	}
	if len(roles) == 0 {
		err := fmt.Errorf("roles must not be empty")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)
	}
	// check if user can revoke role
	err := u.CheckUserPermission(ctx, enums.PermissionTypeCanEditUserRole)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotAuthorizedErr(err)
	}

	ok, err := u.Update.RevokeRoles(ctx, userID, roles)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.RevokeRolesErr(err)
	}
	return ok, nil
}

// GetAllRoles returns all roles
func (u *UsecaseAuthorityImpl) GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error) {
	roles, err := u.Query.GetAllRoles(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetAllRolesErr(err)
	}
	return roles, nil
}

// AssignOrRevokeRoles assigns the specified roles to the user
func (u *UsecaseAuthorityImpl) AssignOrRevokeRoles(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error) {
	if userID == "" {
		err := fmt.Errorf("userID must not be empty")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.EmptyInputErr(err)
	}

	newRoles := []enums.UserRoleType{}
	for _, role := range roles {
		if role != nil {
			newRoles = append(newRoles, *role)
		}
	}

	err := u.CheckUserPermission(ctx, enums.PermissionTypeCanEditUserRole)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UserNotAuthorizedErr(err)
	}

	currentRoles, err := u.Query.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetUserRolesErr(err)
	}

	currentRoleList := []enums.UserRoleType{}
	for _, role := range currentRoles {
		currentRoleList = append(currentRoleList, role.Name)
	}
	_, err = u.Update.RevokeRoles(ctx, userID, currentRoleList)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.RevokeRolesErr(err)
	}

	_, err = u.Update.AssignRoles(ctx, userID, newRoles)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.AssignRolesErr(err)
	}
	return true, nil
}
