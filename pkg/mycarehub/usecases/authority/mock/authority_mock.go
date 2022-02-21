package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// AuthorityUseCaseMock mocks the implementation of usecase methods.
type AuthorityUseCaseMock struct {
	MockCheckUserRoleFn       func(ctx context.Context, role enums.UserRoleType) error
	MockCheckUserPermissionFn func(ctx context.Context, permission enums.PermissionType) error
	MockAssignRolesFn         func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
}

// NewAuthorityUseCaseMock creates in initializes create type mocks
func NewAuthorityUseCaseMock() *AuthorityUseCaseMock {

	return &AuthorityUseCaseMock{

		MockCheckUserRoleFn: func(ctx context.Context, role enums.UserRoleType) error {
			return nil
		},
		MockCheckUserPermissionFn: func(ctx context.Context, permission enums.PermissionType) error {
			return nil
		},
		MockAssignRolesFn: func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
			return false, nil
		},
	}
}

// CheckUserRole mocks the implementation for checking the user role
func (f *AuthorityUseCaseMock) CheckUserRole(ctx context.Context, role enums.UserRoleType) error {
	return f.MockCheckUserRoleFn(ctx, role)
}

// CheckUserPermission mocks the implementation for checking the user permission
func (f *AuthorityUseCaseMock) CheckUserPermission(ctx context.Context, permission enums.PermissionType) error {
	return f.MockCheckUserPermissionFn(ctx, permission)
}

// AssignRoles mocks the implementation for assigning the roles to the user
func (f *AuthorityUseCaseMock) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return f.MockAssignRolesFn(ctx, userID, roles)
}
