package gorm_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_DeleteFacility(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()
	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := gofakeit.Name()
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()

	facility := &gorm.Facility{
		FacilityID:         &ID,
		Name:               name,
		Code:               code,
		Active:             true,
		County:             county,
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
	}

	facility, err := testingDB.GetOrCreateFacility(ctx, facility)
	if err != nil {
		t.Errorf("failed to create test facility: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		mflcode int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete facility",
			args: args{
				ctx:     ctx,
				mflcode: facility.Code,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to delete facility",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid facility",
			args: args{
				ctx:     ctx,
				mflcode: 789555,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.DeleteFacility(tt.args.ctx, tt.args.mflcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.DeleteFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_DeleteClientProfile(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete client profile",
			args: args{
				ctx:      ctx,
				clientID: clientID3,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.DeleteClientProfile(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.DeleteClientProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_DeleteUser(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete user",
			args: args{
				ctx:    ctx,
				userID: userIDToDelete,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to delete user",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.DeleteUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_DeleteStaffProfile(t *testing.T) {
	ctx := context.Background()

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
				ctx:     ctx,
				staffID: staffIDToDelete,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete staff profile",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.DeleteStaffProfile(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.DeleteStaffProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
