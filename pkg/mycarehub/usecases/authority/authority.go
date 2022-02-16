package authority

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
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

// UsecaseAuthority groups al the interfaces for the Authority usecase
type UsecaseAuthority interface {
	ICheckUserRole
	ICheckUserPermission
}

// UsecaseAuthorityImpl represents the Authority implementation
type UsecaseAuthorityImpl struct {
	Query       infrastructure.Query
	ExternalExt extension.ExternalMethodsExtension
}

// NewUsecaseAuthority is the controller function for the Authority usecase
func NewUsecaseAuthority(
	query infrastructure.Query,
	externalExt extension.ExternalMethodsExtension,
) *UsecaseAuthorityImpl {
	return &UsecaseAuthorityImpl{
		Query:       query,
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
