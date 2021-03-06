package authority

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	notificationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
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
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

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
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

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

func TestUsecaseAuthorityImpl_AssignRoles(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: successfully assign roles to a user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing user id",
			args: args{
				ctx:    context.Background(),
				userID: "",
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: missing roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: user not authorized",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to assign roles to a user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

			if tt.name == "sad case: user not authorized" {
				fakeDB.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: failed to assign roles to a user" {
				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to assign roles to a user")
				}
			}

			got, err := u.AssignRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.AssignRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseAuthorityImpl.AssignRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseAuthorityImpl_GetUserRoles(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.AuthorityRole
		wantErr bool
	}{
		{
			name: "happy case: successfully get user roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: missing user id",
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get user roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

			if tt.name == "sad case: failed to get user roles" {
				fakeDB.MockGetUserRolesFn = func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get user roles")
				}
			}

			got, err := u.GetUserRoles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.GetUserRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UsecaseAuthorityImpl.GetUserRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseAuthorityImpl_GetUserPermissions(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.AuthorityPermission
		wantErr bool
	}{
		{
			name: "happy case: successfully get user permissions",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: missing user id",
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get user permissions",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

			if tt.name == "sad case: failed to get user permissions" {
				fakeDB.MockGetUserPermissionsFn = func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
					return nil, fmt.Errorf("failed to get user permissions")
				}
			}

			got, err := u.GetUserPermissions(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.GetUserPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UsecaseAuthorityImpl.GetUserPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseAuthorityImpl_RevokeRoles(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: successfully revoke roles from a user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing user id",
			args: args{
				ctx:    context.Background(),
				userID: "",
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: missing roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: user not authorized",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to revoke roles from a user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)
			if tt.name == "sad case: user not authorized" {
				fakeDB.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: failed to revoke roles from a user" {
				fakeDB.MockRevokeRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to revoke roles from a user")
				}
			}

			got, err := u.RevokeRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.RevokeRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseAuthorityImpl.RevokeRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseAuthorityImpl_AssignOrRevokeRoles(t *testing.T) {

	adminRole := enums.UserRoleTypeSystemAdministrator
	type args struct {
		ctx    context.Context
		userID string
		roles  []*enums.UserRoleType
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: successfully assign or revoke roles from a user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing user id",
			args: args{
				ctx:    context.Background(),
				userID: "",
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to get current roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to revoke current roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: user not authorized",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to assign roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to get user profile",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to notify user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				roles:  []*enums.UserRoleType{&adminRole},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

			if tt.name == "sad case: failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad case: failed to notify user" {
				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
				}
			}

			if tt.name == "sad case: failed to get current roles" {
				fakeDB.MockGetUserRolesFn = func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get current roles")
				}
			}

			if tt.name == "sad case: failed to revoke current roles" {
				fakeDB.MockRevokeRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to revoke current roles")
				}
			}
			if tt.name == "sad case: user not authorized" {
				fakeDB.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: failed to assign roles" {
				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to assign roles")
				}
			}

			got, err := u.AssignOrRevokeRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.AssignOrRevokeRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseAuthorityImpl.AssignOrRevokeRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseAuthorityImpl_ComposeAndSendNotification(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

	type args struct {
		newRoles         []enums.UserRoleType
		ctx              context.Context
		user             *domain.User
		notificationType enums.NotificationType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case: compose and send notification",
			args: args{
				newRoles: []enums.UserRoleType{enums.UserRoleType(enums.NotificationTypeRoleAssignment)},
				ctx:      context.Background(),
				user: &domain.User{
					Username: "test",
				},
				notificationType: enums.NotificationTypeRoleAssignment,
			},
		},
		{
			name: "sad case: unable to compose and send notification",
			args: args{
				newRoles: []enums.UserRoleType{enums.UserRoleType(enums.NotificationTypeRoleAssignment)},
				ctx:      context.Background(),
				user: &domain.User{
					Username: "test",
				},
				notificationType: enums.NotificationTypeRoleRevocation,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: unable to compose and send notification" {
				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
				}
			}
			u.ComposeAndSendNotification(tt.args.ctx, tt.args.newRoles, tt.args.notificationType, tt.args.user)
		})
	}
}

func TestUsecaseAuthorityImpl_GetAllRoles(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	u := NewUsecaseAuthority(fakeDB, fakeDB, fakeExtension, fakeNotification)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get all roles",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to get all roles" {
				fakeDB.MockGetAllRolesFn = func(ctx context.Context) ([]*domain.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get all roles")
				}
			}
			_, err := u.GetAllRoles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseAuthorityImpl.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
