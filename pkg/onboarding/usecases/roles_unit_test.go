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
			name: "sad_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_unable_to_get_userprofile_by_id",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_unable_to_get_userprofile",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_unable_to_get_user_roles",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_unable_to_get_role",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_user_do_not_have_required_permissions",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad_unable_to_create_role_in_database",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name: "happy_created_role",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: input,
		// 	},
		// 	want:    &dto.RoleOutput{},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "sad_unable_to_get_userprofile_by_id" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile by UID")
				}
			}

			if tt.name == "sad_unable_to_get_userprofile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to get user profile")
				}
			}

			if tt.name == "sad_unable_to_get_user_roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						Roles: []string{"123"},
					}, nil
				}

				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return nil, fmt.Errorf("error unable to get user roles")
				}

			}

			if tt.name == "sad_user_do_not_have_required_permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						Roles: []string{"123"},
					}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.edit"}},
					}, nil
				}

			}
			if tt.name == "sad_unable_to_create_role_in_database" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						Roles: []string{"123"},
					}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.create"}},
					}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, role profileutils.Role) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error un able to create role in db")
				}

			}
			if tt.name == "happy_created_role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						Roles: []string{"123"},
					}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.create"}},
					}, nil
				}
				fakeRepo.CreateRoleFn = func(ctx context.Context, role profileutils.Role) (*profileutils.Role, error) {
					return &profileutils.Role{}, nil
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

func TestRoleUseCaseImpl_AddPermissionsToRole(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()

	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	type args struct {
		ctx   context.Context
		input dto.RolePermissionInput
	}

	input := dto.RolePermissionInput{
		RoleID: "123",
		Scope:  "role.edit",
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		{
			name: "unable_to_get_userprofile_in_roles",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		{
			name: "unable_to_get_user_roles",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		{
			name: "sad_user_do_not_have_required_permissions",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		{
			name: "unable_to_get_role_to_edit",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		{
			name: "sad_unable_to_add_permission_to_role",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
		// {
		// 	name: "happy_added_permission_to_roles",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: input,
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "unable_to_get_userprofile_in_roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to get user profile")
				}
			}

			if tt.name == "unable_to_get_user_roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return nil, fmt.Errorf("Error unable to get role by ides")
				}
			}

			if tt.name == "sad_user_do_not_have_required_permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{Roles: []string{"123"}}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"employee.create"}},
					}, nil
				}
			}

			if tt.name == "unable_to_get_role_to_edit" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{Roles: []string{"123"}}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.edit"}},
					}, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return nil, fmt.Errorf("error unable to get role to edit")
				}
			}

			if tt.name == "sad_unable_to_add_permission_to_role" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{Roles: []string{"123"}}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.edit"}},
					}, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "123",
						Scopes: []string{"role.edit"},
					}, nil
				}
				fakeRepo.UpdateRoleDetailsFn = func(ctx context.Context, role profileutils.Role) error {
					return fmt.Errorf("error unable to update role")
				}
			}

			if tt.name == "happy_added_permission_to_roles" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{Roles: []string{"123"}}, nil
				}
				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					return &[]profileutils.Role{
						{ID: "123", Scopes: []string{"role.edit"}},
					}, nil
				}
				fakeRepo.GetRoleByIDFn = func(ctx context.Context, roleID string) (*profileutils.Role, error) {
					return &profileutils.Role{
						ID:     "123",
						Scopes: []string{"role.create"},
					}, nil
				}
				fakeRepo.UpdateRoleDetailsFn = func(ctx context.Context, role profileutils.Role) error {
					return nil
				}
			}

			if _, err := i.Role.AddPermissionsToRole(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf(
					"RoleUseCaseImpl.AddPermissionsToRole() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
