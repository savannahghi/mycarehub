package utils

import (
	"context"
	"reflect"
	"testing"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
)

func TestCheckUserHasPermission(t *testing.T) {
	type args struct {
		roles      []profileutils.Role
		permission profileutils.Permission
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{
			name: "sad: user do not have permission, role deactivated",
			args: args{
				roles: []profileutils.Role{
					{Name: "Employee Role", Scopes: []string{"agent.view"}, Active: false},
				},
				permission: profileutils.CanViewAgent,
			},
			want: false,
		},

		{
			name: "sad: user do not have permission, no such scope",
			args: args{
				roles: []profileutils.Role{
					{Name: "Employee Role", Scopes: []string{"patient.create"}, Active: true},
				},
				permission: profileutils.CanViewAgent,
			},
			want: false,
		},
		{
			name: "happy: user has permission",
			args: args{
				roles: []profileutils.Role{
					{Name: "Employee Role", Scopes: []string{"agent.view"}, Active: true},
				},
				permission: profileutils.CanViewAgent,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUserHasPermission(tt.args.roles, tt.args.permission); got != tt.want {
				t.Errorf("CheckUserHasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserNavigationActions(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		user  profileutils.UserProfile
		roles []profileutils.Role
	}

	homeNavAction := domain.HomeNavAction
	homeNavAction.Favorite = true

	agentNavActions := domain.AgentNavActions
	agentNavActions.Nested = []interface{}{
		domain.AgentRegistrationNavAction,
		domain.AgentidentificationNavAction,
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.GroupedNavigationActions
		wantErr bool
	}{
		{
			name: "happy got user navigation actions",
			args: args{
				ctx: ctx,
				user: profileutils.UserProfile{
					FavNavActions: []string{"Home"},
				},
				roles: []profileutils.Role{
					{
						Scopes: []string{"agent.view", "agent.register", "agent.identify"},
						Active: true,
					},
				},
			},
			want: &dto.GroupedNavigationActions{
				Primary: []domain.NavigationAction{
					homeNavAction,
					domain.HelpNavAction,
				},
				Secondary: []domain.NavigationAction{
					agentNavActions,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserNavigationActions(tt.args.ctx, tt.args.user, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserNavigationActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserNavigationActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupNested(t *testing.T) {
	type args struct {
		actions []domain.NavigationAction
	}
	expectedOutput := []domain.NavigationAction{}

	navAction := domain.NavigationAction{
		Group: domain.HomeGroup,
		Title: "Home",
		Nested: []interface{}{
			domain.NavigationAction{
				Group:     domain.HomeGroup,
				Title:     "Child 1",
				HasParent: true,
			},
			domain.NavigationAction{
				Group:     domain.HomeGroup,
				Title:     "Child 2",
				HasParent: true,
			},
		},
	}

	expectedOutput = append(expectedOutput, navAction)

	tests := []struct {
		name string
		args args
		want []domain.NavigationAction
	}{
		{
			name: "happy grouped nested navigation actions",
			args: args{
				actions: []domain.NavigationAction{
					{
						Group: domain.HomeGroup,
						Title: "Home",
					},
					{
						Group:     domain.HomeGroup,
						Title:     "Child 1",
						HasParent: true,
					},
					{
						Group:     domain.HomeGroup,
						Title:     "Child 2",
						HasParent: true,
					},
				},
			},
			want: expectedOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupNested(tt.args.actions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupNested() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupPriority(t *testing.T) {
	type args struct {
		actions []domain.NavigationAction
	}
	actions := []domain.NavigationAction{}

	navAction1 := domain.NavigationAction{
		Group:          domain.HomeGroup,
		Title:          "Home",
		SequenceNumber: 1,
	}
	navAction2 := domain.NavigationAction{
		Group:          domain.AgentGroup,
		Title:          "Agent",
		SequenceNumber: 2,
		Nested: []interface{}{
			domain.NavigationAction{
				Group:     domain.AgentGroup,
				Title:     "Child 1",
				HasParent: true,
			},
			domain.NavigationAction{
				Group:     domain.AgentGroup,
				Title:     "Child 2",
				HasParent: true,
			},
		},
	}
	navAction3 := domain.NavigationAction{
		Group:          domain.PatientGroup,
		Title:          "Patients",
		SequenceNumber: 3,
	}
	navAction4 := domain.NavigationAction{
		Group:          domain.PartnerGroup,
		Title:          "Partner",
		SequenceNumber: 4,
	}
	navAction5 := domain.NavigationAction{
		Group:          domain.RoleGroup,
		Title:          "Role",
		SequenceNumber: 5,
	}
	navAction6 := domain.NavigationAction{
		Group:          domain.ConsumerGroup,
		Title:          "Consumers",
		SequenceNumber: 6,
	}
	navAction7 := domain.NavigationAction{
		Group:          domain.EmployeeGroup,
		Title:          "Employee",
		SequenceNumber: 7,
	}

	actions = append(actions, navAction1)
	actions = append(actions, navAction2)
	actions = append(actions, navAction3)
	actions = append(actions, navAction4)
	actions = append(actions, navAction5)
	actions = append(actions, navAction6)
	actions = append(actions, navAction7)

	tests := []struct {
		name          string
		args          args
		wantPrimary   []domain.NavigationAction
		wantSecondary []domain.NavigationAction
	}{
		{
			name: "happy: grouped into priorities",
			args: args{
				actions: actions,
			},
			wantPrimary: []domain.NavigationAction{
				navAction1,
				navAction3,
				navAction4,
				navAction5,
			},
			wantSecondary: []domain.NavigationAction{
				navAction2,
				navAction6,
				navAction7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimary, gotSecondary := GroupPriority(tt.args.actions)
			if !reflect.DeepEqual(gotPrimary, tt.wantPrimary) {
				t.Errorf("GroupPriority() gotPrimary = %v, want %v", gotPrimary, tt.wantPrimary)
			}
			if !reflect.DeepEqual(gotSecondary, tt.wantSecondary) {
				t.Errorf("GroupPriority() gotSecondary = %v, want %v", gotSecondary, tt.wantSecondary)
			}
		})
	}
}

func TestGetUserPermissions(t *testing.T) {
	type args struct {
		roles []profileutils.Role
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "happy got the right permissions",
			args: args{[]profileutils.Role{
				{Name: "Agent Role", Scopes: []string{"agent.register", "agent.view"}, Active: true},
				{Name: "Agent Role", Scopes: []string{"patient.create", "agent.view"}, Active: true},
				{Name: "Agent Role", Scopes: []string{"role.view", "role.create"}, Active: false},
			}},
			want: []string{"agent.register", "agent.view", "patient.create"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserPermissions(tt.args.roles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicateStrings(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "happy removed duplicates",
			args: args{strings: []string{"user", "tes", "user", "123", "another user"}},
			want: []string{"user", "tes", "123", "another user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicateStrings(tt.args.strings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicateStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
