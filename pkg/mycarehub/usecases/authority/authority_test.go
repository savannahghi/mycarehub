package authority

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
)

func TestUsecaseAuthorityImpl_CheckUserRole(t *testing.T) {
	type fields struct {
		Query infrastructure.Query
	}
	type args struct {
		ctx    context.Context
		userID string
		role   enums.UserRoleType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successfully check if a user has a role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   enums.UserRoleTypeContentManagement,
			},
			wantErr: false,
		},
		{
			name: "sad case: missing role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to check if use has role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   enums.UserRoleTypeClientManagement,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get logged in user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   enums.UserRoleTypeContentManagement,
			},
			wantErr: true,
		},
		{
			name: "sad case: user not authorized",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   enums.UserRoleTypeContentManagement,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := NewUsecaseAuthority(fakeDB, fakeExtension)

			if tt.name == "sad case: failed to check if use has role" {
				fakeDB.MockCheckUserRoleFn = func(ctx context.Context, userID string, role string) (bool, error) {
					return false, fmt.Errorf("failed to check if use has role")
				}
			}

			if tt.name == "sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "sad case: user not authorized" {
				fakeDB.MockCheckUserRoleFn = func(ctx context.Context, userID string, role string) (bool, error) {
					return false, nil
				}
			}

			err := u.CheckUserRole(tt.args.ctx, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.CheckUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseAuthorityImpl_CheckUserPermission(t *testing.T) {
	type args struct {
		ctx        context.Context
		permission enums.PermissionType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successfully check if a user has a permission",
			args: args{
				ctx:        context.Background(),
				permission: enums.PermissionTypeCanEditOwnRole,
			},

			wantErr: false,
		},
		{
			name: "sad case: missing permission",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to check if use has permission",
			args: args{
				ctx:        context.Background(),
				permission: enums.PermissionTypeCanEditOwnRole,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get logged in user",
			args: args{
				ctx:        context.Background(),
				permission: enums.PermissionTypeCanEditOwnRole,
			},
			wantErr: true,
		},
		{
			name: "sad case: user not authorized",
			args: args{
				ctx:        context.Background(),
				permission: enums.PermissionTypeCanEditOwnRole,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := NewUsecaseAuthority(fakeDB, fakeExtension)

			if tt.name == "sad case: failed to check if use has permission" {
				fakeDB.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, fmt.Errorf("failed to check if use has permission")
				}
			}

			if tt.name == "sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "sad case: user not authorized" {
				fakeDB.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, nil
				}
			}

			err := u.CheckUserPermission(tt.args.ctx, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.CheckUserPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
