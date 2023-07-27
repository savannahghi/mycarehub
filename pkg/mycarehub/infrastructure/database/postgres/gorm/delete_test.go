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
				ctx:     context.Background(),
				staffID: staffIDToDelete,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete staff profile",
			args: args{
				ctx:     context.Background(),
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
				ctx:         context.Background(),
				communityID: communityIDToDelete,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete community, not found",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
			},
			wantErr: false, // skip error checking for this case
		},
		{
			name: "Sad Case - Unable delete community, invalid id",
			args: args{
				ctx:         context.Background(),
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

func TestPGInstance_DeleteAccessToken(t *testing.T) {

	type args struct {
		ctx       context.Context
		signature string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete token",
			args: args{
				ctx:       context.Background(),
				signature: "5Ueg0S3v3ZaoiFLgVD-ysjskmOgDs44koLcUY93rolI",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.DeleteAccessToken(tt.args.ctx, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_DeleteRefreshToken(t *testing.T) {

	type args struct {
		ctx       context.Context
		signature string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete token",
			args: args{
				ctx:       context.Background(),
				signature: "RYWZKMrji0MqV82zhnjpAaP4hP-L3kMNXQMsu4pw3gU",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.DeleteRefreshToken(tt.args.ctx, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_DeleteClientProfile(t *testing.T) {

	type args struct {
		ctx      context.Context
		clientID string
		userID   *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete client with one user profile",
			args: args{
				ctx:      context.Background(),
				clientID: testOPtOutClient,
				userID:   &testOPtOutClient,
			},
			wantErr: false,
		},
		{
			name: "happy case: delete client who is a caregiver",
			args: args{
				ctx:      context.Background(),
				clientID: testOPtOutClientCaregiver,
			},
			wantErr: false,
		},
		{
			name: "happy case: delete client who has a staff profile",
			args: args{
				ctx:      context.Background(),
				clientID: testOPtOutClientStaff,
			},
			wantErr: false,
		},
		{
			name: "happy case: delete client who has a staff profile 2",
			args: args{
				ctx:      context.Background(),
				clientID: testOptOutStaffClient,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid client id",
			args: args{
				ctx:      context.Background(),
				clientID: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.DeleteClientProfile(tt.args.ctx, tt.args.clientID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteClientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
