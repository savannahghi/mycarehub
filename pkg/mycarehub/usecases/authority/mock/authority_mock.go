package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// AuthorityUseCaseMock mocks the implementation of usecase methods.
type AuthorityUseCaseMock struct {
	MockCheckUserRoleFn       func(ctx context.Context, role enums.UserRoleType) error
	MockCheckUserPermissionFn func(ctx context.Context, permission enums.PermissionType) error
	MockAssignRolesFn         func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetUserRolesFn        func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	MockGetUserPermissionsFn  func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
	MockRevokeRolesFn         func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetAllRolesFn         func(ctx context.Context) ([]*domain.AuthorityRole, error)
	MockAssignOrRevokeRolesFn func(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error)
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
		MockGetUserRolesFn: func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
			return []*domain.AuthorityRole{
				{
					AuthorityRoleID: uuid.New().String(),
					Name:            enums.UserRoleTypeClientManagement,
				},
			}, nil
		},
		MockGetUserPermissionsFn: func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
			return []*domain.AuthorityPermission{
				{
					PermissionID: uuid.New().String(),
					Name:         enums.PermissionTypeCanManageClient,
				},
			}, nil
		},
		MockRevokeRolesFn: func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
			return false, nil
		},
		MockGetAllRolesFn: func(ctx context.Context) ([]*domain.AuthorityRole, error) {
			return []*domain.AuthorityRole{
				{
					AuthorityRoleID: uuid.New().String(),
					Name:            enums.UserRoleTypeClientManagement,
				},
			}, nil
		},
		MockAssignOrRevokeRolesFn: func(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error) {
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

// GetUserRoles mocks the implementation of getting all roles for a user
func (f *AuthorityUseCaseMock) GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
	return f.MockGetUserRolesFn(ctx, userID)
}

// GetUserPermissions mocks the implementation of getting all permissions for a user
func (f *AuthorityUseCaseMock) GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
	return f.MockGetUserPermissionsFn(ctx, userID)
}

// RevokeRoles mocks the implementation for revoking the roles from the user
func (f *AuthorityUseCaseMock) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return f.MockRevokeRolesFn(ctx, userID, roles)
}

// GetAllRoles mocks the implementation for getting all roles
func (f *AuthorityUseCaseMock) GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error) {
	return f.MockGetAllRolesFn(ctx)
}

// AssignOrRevokeRoles mocks the implementation for assigning or revoking roles
func (f *AuthorityUseCaseMock) AssignOrRevokeRoles(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error) {
	return f.MockAssignOrRevokeRolesFn(ctx, userID, roles)
}
