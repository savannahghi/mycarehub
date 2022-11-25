package gorm_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestPGInstance_DeleteStaffProfile(t *testing.T) {

	type args struct {
		ctx     context.Context
		staffID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete staff profile",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				staffID: staffIDToDelete,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete staff profile",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				staffID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testingDB.DeleteStaffProfile(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_DeleteCommunity(t *testing.T) {
	type args struct {
		ctx         context.Context
		communityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete community",
			args: args{
				ctx:         addRequiredContext(context.Background(), t),
				communityID: communityIDToDelete,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete community, not found",
			args: args{
				ctx:         addRequiredContext(context.Background(), t),
				communityID: uuid.New().String(),
			},
			wantErr: false, // skip error checking for this case
		},
		{
			name: "Sad Case - Unable delete community, invalid id",
			args: args{
				ctx:         addRequiredContext(context.Background(), t),
				communityID: "invalid id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.DeleteCommunity(tt.args.ctx, tt.args.communityID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteCommunity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_RemoveFacilitiesFromClientProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: remove facilities from client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   clientID,
				facilities: []string{facilityToRemoveFromUserProfile},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed remove facilities from client profile, invalid client ID",
			args: args{
				ctx:        context.Background(),
				clientID:   "Invalid",
				facilities: []string{facilityToRemoveFromUserProfile},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.RemoveFacilitiesFromClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RemoveFacilitiesFromClientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_RemoveFacilitiesFromStaffProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		staffID    string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: remove facilities from staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    staffID,
				facilities: []string{facilityToRemoveFromUserProfile},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed remove facilities from staff profile, invalid staff ID",
			args: args{
				ctx:        context.Background(),
				staffID:    "Invalid",
				facilities: []string{facilityToRemoveFromUserProfile},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.RemoveFacilitiesFromStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RemoveFacilitiesFromStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_DeleteOrganisation(t *testing.T) {
	invalidOrgID := "invalid"

	type args struct {
		ctx          context.Context
		organisation *gorm.Organisation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete organisation",
			args: args{
				ctx:          addRequiredContext(context.Background(), t),
				organisation: &gorm.Organisation{ID: &organisationIDToDelete},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to delete organisation",
			args: args{
				ctx:          addRequiredContext(context.Background(), t),
				organisation: &gorm.Organisation{ID: &invalidOrgID},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.DeleteOrganisation(tt.args.ctx, tt.args.organisation); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
