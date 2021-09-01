package utils_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
	"github.com/stretchr/testify/assert"
)

func TestMatchAndReturn(t *testing.T) {
	tests := []struct {
		old  bool
		new  bool
		want bool
	}{
		{old: false, new: true, want: true},
		{old: true, new: false, want: false},
		{old: true, new: true, want: true},
		{old: false, new: false, want: false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			resp := utils.MatchAndReturn(tt.old, tt.new)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestFindNumber(t *testing.T) {
	type args struct {
		slice []string
		value string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "happy case - Find existing number",
			args: args{
				slice: []string{interserviceclient.TestUserPhoneNumber, "+254700998877"},
				value: interserviceclient.TestUserPhoneNumber,
			},
			// This is the index
			want:  0,
			want1: true,
		},
		{
			name: "sad case - non existent number",
			args: args{
				slice: []string{interserviceclient.TestUserPhoneNumber, "+254700998877"},
				value: "invalid",
			},
			// This is the index
			want:  -1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := utils.FindItem(tt.args.slice, tt.args.value)
			if got != tt.want {
				t.Errorf("FindNumber() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindNumber() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUniquePermissionsArray(t *testing.T) {
	duplicated := []profileutils.PermissionType{}
	duplicated = append(duplicated, profileutils.DefaultAdminPermissions...)
	duplicated = append(duplicated, profileutils.DefaultAdminPermissions...)
	duplicated = append(duplicated, profileutils.DefaultAdminPermissions...)

	duplicatedMixed := []profileutils.PermissionType{}
	duplicatedMixed = append(duplicatedMixed, profileutils.DefaultAdminPermissions...)
	duplicatedMixed = append(duplicatedMixed, profileutils.DefaultAgentPermissions...)
	duplicatedMixed = append(duplicatedMixed, profileutils.DefaultAdminPermissions...)
	duplicatedMixed = append(duplicatedMixed, profileutils.DefaultAgentPermissions...)
	mixed := []profileutils.PermissionType{}
	mixed = append(mixed, profileutils.DefaultAdminPermissions...)
	mixed = append(mixed, profileutils.DefaultAgentPermissions...)

	type args struct {
		arr []profileutils.PermissionType
	}
	tests := []struct {
		name string
		args args
		want []profileutils.PermissionType
	}{
		{
			name: "success:return unique array of permissions",
			args: args{
				arr: duplicatedMixed,
			},
			want: mixed,
		},
		{
			name: "success:return unique array of permissions",
			args: args{
				arr: duplicated,
			},
			want: profileutils.DefaultAdminPermissions,
		},
		{
			name: "success:return same unique array",
			args: args{
				arr: profileutils.DefaultAdminPermissions,
			},
			want: profileutils.DefaultAdminPermissions,
		},
		{
			name: "success:empty array of permissions",
			args: args{
				arr: []profileutils.PermissionType{},
			},
			want: []profileutils.PermissionType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.UniquePermissionsArray(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniquePermissionsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ServiceHealthEndPoint(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid_case",
			args:    "https://admin-staging.healthcloud.co.ke/graphql",
			want:    "https://admin-staging.healthcloud.co.ke/health",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ServiceHealthEndPoint(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("serviceHealthEndPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("serviceHealthEndPoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewActionsMapper(t *testing.T) {
	ctx := context.Background()

	nestedOne := domain.NavigationAction{
		Title:      "Nested One",
		OnTapRoute: "some/route",
		Icon:       "http://one.asset.com",
		Favorite:   false,
	}

	nestedTwo := domain.NavigationAction{
		Title:      "Nested Two",
		OnTapRoute: "some/route",
		Icon:       "http://two.asset.com",
		Favorite:   false,
	}

	new := &dto.GroupedNavigationActions{
		Primary: []domain.NavigationAction{
			{
				Title:      "Home",
				OnTapRoute: "some/route",
				Icon:       "http://home.asset.com",
				Favorite:   false,
			},
			{
				Title:      "Help",
				OnTapRoute: "some/route",
				Icon:       "http://help.asset.com",
				Favorite:   false,
			},
		},
		Secondary: []domain.NavigationAction{
			{
				Title:      "Secondary One",
				OnTapRoute: "some/route",
				Icon:       "http://one.asset.com",
				Favorite:   true,
				Nested:     []interface{}{nestedOne, nestedTwo},
			},
			{
				Title:      "Secondary Two",
				OnTapRoute: "some/route",
				Icon:       "http://two.asset.com",
				Favorite:   false,
				Nested:     []interface{}{nestedOne, nestedTwo},
			},
		},
	}

	old := &profileutils.NavigationActions{
		Primary:   []profileutils.NavAction{},
		Secondary: []profileutils.NavAction{},
	}

	type args struct {
		ctx     context.Context
		grouped *dto.GroupedNavigationActions
	}
	tests := []struct {
		name string
		args args
		want *profileutils.NavigationActions
	}{
		{
			name: "success: map empty new actions to old actions",
			args: args{
				ctx: ctx,
				grouped: &dto.GroupedNavigationActions{
					Primary:   []domain.NavigationAction{},
					Secondary: []domain.NavigationAction{},
				},
			},
			want: &profileutils.NavigationActions{},
		},
		{
			name: "success: map new actions to old actions",
			args: args{
				ctx:     ctx,
				grouped: new,
			},
			want: old,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NewActionsMapper(tt.args.ctx, tt.args.grouped)
			if got == nil {
				t.Errorf("NewActionsMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}
