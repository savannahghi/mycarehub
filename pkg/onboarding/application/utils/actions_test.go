package utils

import (
	"reflect"
	"testing"

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

func TestGroupNested(t *testing.T) {
	type args struct {
		actions []domain.NavigationAction
	}
	expectedOutput := make(map[domain.NavigationGroup]domain.NavigationAction)
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

	expectedOutput[domain.HomeGroup] = navAction

	tests := []struct {
		name string
		args args
		want map[domain.NavigationGroup]domain.NavigationAction
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
		actions map[domain.NavigationGroup]domain.NavigationAction
	}
	actions := make(map[domain.NavigationGroup]domain.NavigationAction)
	navAction1 := domain.NavigationAction{
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
	navAction2 := domain.NavigationAction{
		Group: domain.AgentGroup,
		Title: "Agent",
	}
	navAction3 := domain.NavigationAction{
		Group: domain.PatientGroup,
		Title: "Patients",
	}
	navAction4 := domain.NavigationAction{
		Group: domain.PartnerGroup,
		Title: "Partner",
	}
	navAction5 := domain.NavigationAction{
		Group:          domain.RoleGroup,
		Title:          "Role",
		IsHighPriority: true,
	}
	navAction6 := domain.NavigationAction{
		Group:          domain.ConsumerGroup,
		Title:          "Consumers",
		IsHighPriority: true,
	}
	navAction7 := domain.NavigationAction{
		Group:          domain.EmployeeGroup,
		Title:          "Employee",
		IsHighPriority: true,
	}

	actions[domain.HomeGroup] = navAction1
	actions[domain.AgentGroup] = navAction2
	actions[domain.PatientGroup] = navAction3
	actions[domain.PartnerGroup] = navAction4
	actions[domain.RoleGroup] = navAction5
	actions[domain.ConsumerGroup] = navAction6
	actions[domain.EmployeeGroup] = navAction7

	tests := []struct {
		name          string
		args          args
		wantPrimary   int
		wantSecondary int
	}{
		{
			name: "happy: grouped into priorities",
			args: args{
				actions: actions,
			},
			wantPrimary:   4,
			wantSecondary: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrimary, gotSecondary := GroupPriority(tt.args.actions)
			if !reflect.DeepEqual(len(gotPrimary), tt.wantPrimary) {
				t.Errorf("GroupPriority() gotPrimary = %v, want %v", len(gotPrimary), tt.wantPrimary)
			}
			if !reflect.DeepEqual(len(gotSecondary), tt.wantSecondary) {
				t.Errorf("GroupPriority() gotSecondary = %v, want %v", len(gotSecondary), tt.wantSecondary)
			}
		})
	}
}
