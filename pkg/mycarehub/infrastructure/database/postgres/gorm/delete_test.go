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
	ctx := context.Background()
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
				ctx:         ctx,
				communityID: communityIDToDelete,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete community, not found",
			args: args{
				ctx:         ctx,
				communityID: uuid.New().String(),
			},
			wantErr: false, // skip error checking for this case
		},
		{
			name: "Sad Case - Unable delete community, invalid id",
			args: args{
				ctx:         ctx,
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
