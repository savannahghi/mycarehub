package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_DeleteStaffProfile(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case",
			args: args{
				ctx:     ctx,
				staffID: "123456789",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				staffID: "123456789",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockDeleteStaffProfileFn = func(ctx context.Context, staffID string) error {
					return fmt.Errorf("an error occurred while deleting")
				}
			}
			err := d.DeleteStaffProfile(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_DeleteCommunity(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case",
			args: args{
				ctx:         ctx,
				communityID: "123456789",
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				communityID: "123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockDeleteCommunityFn = func(ctx context.Context, communityID string) error {
					return fmt.Errorf("an error occurred while deleting")
				}
			}
			if err := d.DeleteCommunity(tt.args.ctx, tt.args.communityID); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteCommunity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RemoveFacilitiesFromClientProfile(t *testing.T) {
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
			name: "Happy case: remove facilities from  client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to remove facilities from  client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to remove facilities from  client profile" {
				fakeGorm.MockRemoveFacilitiesFromClientProfileFn = func(ctx context.Context, clientID string, facilities []string) error {
					return fmt.Errorf("failed to remove facilities from client profile")
				}
			}

			if err := d.RemoveFacilitiesFromClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RemoveFacilitiesFromClientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_DeleteOrganisation(t *testing.T) {
	type args struct {
		ctx          context.Context
		organisation *domain.Organisation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: delete organisation",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					ID: uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to delete organisation",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					ID: uuid.NewString(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to delete organisation" {
				fakeGorm.MockDeleteOrganisationFn = func(ctx context.Context, organisation *gorm.Organisation) error {
					return fmt.Errorf("unable to delete organisation")
				}
			}

			if err := d.DeleteOrganisation(tt.args.ctx, tt.args.organisation); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_DeleteClientProfile(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	userID := gofakeit.UUID()

	type args struct {
		ctx      context.Context
		clientID string
		userID   *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: delete client profile",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
				userID:   &userID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to delete client profile",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
				userID:   &userID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to delete client profile" {
				fakeGorm.MockDeleteClientProfileFn = func(ctx context.Context, clientID string, userID *string) error {
					return fmt.Errorf("an error occurred while deleting")
				}
			}
			err := d.DeleteClientProfile(tt.args.ctx, tt.args.clientID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
