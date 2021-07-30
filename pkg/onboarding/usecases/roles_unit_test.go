package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/profileutils"
)

func TestRoleUseCaseImpl_CreateRole(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	input := dto.RoleInput{
		Name: "Agents",
	}

	allPerms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		t.Error("error did not get all permissions")
		return
	}

	perms := []profileutils.Permission{}
	for _, perm := range allPerms {
		if perm.Scope == "role.edit" {
			perm.Allowed = true
		}
		perms = append(perms, perm)
	}
	expectedOutput := &dto.RoleOutput{
		Scopes:      []string{"role.edit"},
		Permissions: perms,
	}

	type args struct {
		ctx   context.Context
		input dto.RoleInput
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.RoleOutput
		wantErr bool
	}{
		{
			name: "sad: unable to get logged in user",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to check if user has permissions",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: user do not have required permissions",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to get user's profile",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create role in database",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy:created role",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    expectedOutput,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad: unable to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "sad: unable to check if user has permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("error unable to check if user has required permissions")
				}
			}

			if tt.name == "sad: user do not have required permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad: unable to get user's profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to get user profile")
				}
			}

			if tt.name == "sad: unable to create role in database" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, profileID string, role dto.RoleInput) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error un able to create role in db")
				}
			}

			if tt.name == "happy:created role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, profileID string, role dto.RoleInput) (*profileutils.Role, error) {
					return &profileutils.Role{
						Scopes: []string{"role.edit"},
					}, nil
				}
			}

			got, err := i.Role.CreateRole(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleUseCaseImpl.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleUseCaseImpl.CreateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleUseCaseImpl_GetAllRoles(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	allPerms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		t.Errorf("failed to get all permissions")
		return
	}
	rolePerms := []profileutils.Permission{}
	for _, perm := range allPerms {
		if perm.Scope == "role.create" {
			perm.Allowed = true
		}
		rolePerms = append(rolePerms, perm)
	}

	expectedOutput := []*dto.RoleOutput{
		{
			ID:          "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
			Scopes:      []string{"role.create"},
			Permissions: rolePerms,
		},
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.RoleOutput
		wantErr bool
	}{
		{
			name:    "sad: did not get logged in user",
			args:    args{ctx: ctx},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad: unable to check if user has permission",
			args:    args{ctx: ctx},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad: user do not have required permission",
			args:    args{ctx: ctx},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad: did not get roles from database",
			args:    args{ctx: ctx},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "happy: got roles",
			args:    args{ctx: ctx},
			want:    expectedOutput,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad: did not get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("error, did not get logged in user")
				}
			}

			if tt.name == "sad: unable to check if user has permission" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("error unable to check is user has permission")
				}
			}

			if tt.name == "sad: user do not have required permission" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "sad: did not get roles from database" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetAllRolesFn = func(ctx context.Context) (*[]profileutils.Role, error) {
					return nil, fmt.Errorf("error, did not get roles from database")
				}
			}

			if tt.name == "happy: got roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetAllRolesFn = func(ctx context.Context) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{
							ID:     "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
							Scopes: []string{"role.create"},
						},
					}, nil
				}
			}
			got, err := i.Role.GetAllRoles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleUseCaseImpl.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleUseCaseImpl.GetAllRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleUseCaseImpl_GetAllPermissions(t *testing.T) {

	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	allPerms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		t.Error("error did not get all permissions")
		return
	}

	output := []*profileutils.Permission{}
	for _, perm := range allPerms {
		output = append(output, &perm)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*profileutils.Permission
		wantErr bool
	}{
		{
			name:    "happy got all permissions",
			args:    args{ctx: ctx},
			want:    output,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Role.GetAllPermissions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleUseCaseImpl.GetAllPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleUseCaseImpl.GetAllPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleUseCaseImpl_AddPermissionsToRole(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	input := dto.RolePermissionInput{
		RoleID: "123",
		Scopes: []string{"role.create"},
	}

	allPerms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		t.Error("error did not get all permissions")
		return
	}

	perms := []profileutils.Permission{}
	for _, perm := range allPerms {
		if perm.Scope == "role.create" {
			perm.Allowed = true
		}
		perms = append(perms, perm)
	}

	expectedOutput := dto.RoleOutput{
		ID:          "123",
		Scopes:      []string{"role.create"},
		Permissions: perms,
	}

	type args struct {
		ctx   context.Context
		input dto.RolePermissionInput
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.RoleOutput
		wantErr bool
	}{
		{
			name: "sad unable to get logged in user",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad unable to check if user has permissions",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad user do not have required permission",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad unable to get role by id",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad unable to get user profile",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad unable to update role details",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy_added_permission_to_roles",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    &expectedOutput,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad unable to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "sad unable to check if user has permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("unable to check permissions")
				}
			}

			if tt.name == "sad user do not have required permission" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad unable to get role by id" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error unable to get role to edit")
				}
			}

			if tt.name == "sad unable to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}

			if tt.name == "sad unable to update role details" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.UpdateRoleDetailsFn = func(ctx context.Context, profileID string, role profileutils.Role) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error unable to update role")
				}
			}

			if tt.name == "happy_added_permission_to_roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "123"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.UpdateRoleDetailsFn = func(ctx context.Context, profileID string, role profileutils.Role) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "123",
						Scopes: []string{"role.create"},
					}, nil
				}
			}

			got, err := i.Role.AddPermissionsToRole(tt.args.ctx, tt.args.input)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"RoleUseCaseImpl.AddPermissionsToRole() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleUseCaseImpl.AddPermissionsToRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleUseCaseImpl_GetRole(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx context.Context
		ID  string
	}

	var input = args{ctx: ctx, ID: "123"}

	allPerms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		t.Error("error did not get all permissions")
		return
	}

	perms := []profileutils.Permission{}
	for _, perm := range allPerms {
		if perm.Scope == "agent.register" {
			perm.Allowed = true
		}
		perms = append(perms, perm)
	}
	var expectedOutput = dto.RoleOutput{
		ID:          "123",
		Scopes:      []string{"agent.register"},
		Permissions: perms,
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.RoleOutput
		wantErr bool
	}{
		{
			name:    "sad unable to get logged in user",
			args:    input,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad unable to check user permissions",
			args:    input,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad user do not have required permission",
			args:    input,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sad unable to get role",
			args:    input,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "happy got role",
			args:    input,
			want:    &expectedOutput,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad unable to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("error unable to  get logged in user")
				}
			}
			if tt.name == "sad unable to check user permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "1234"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("error unable to check permission")
				}
			}

			if tt.name == "sad user do not have required permission" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "1234"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad unable to get role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "1234"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error unable to get role by id")
				}
			}
			if tt.name == "happy got role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: "1234"}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "123",
						Scopes: []string{"agent.register"},
					}, nil
				}
			}

			got, err := i.Role.GetRole(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoleUseCaseImpl.GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleUseCaseImpl.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
