package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/interserviceclient"
	helpers_mock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_DeleteFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	name := gofakeit.Name()
	code := gofakeit.Number(0, 100)
	phone := interserviceclient.TestUserPhoneNumber
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Phone:       phone,
		Active:      true,
		County:      county,
		Description: description,
	}

	// create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflcode := facility.Code

	veryBadMFLCode := 987668878900987654

	type args struct {
		ctx context.Context
		id  int
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
				ctx: ctx,
				id:  mflcode,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  6789,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Sad case - empty MFL Code",
			args: args{
				ctx: ctx,
				id:  0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad MFL Code",
			args: args{
				ctx: ctx,
				id:  veryBadMFLCode,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			fakeHelpers := helpers_mock.NewFakeHelper()

			if tt.name == "Happy case" {
				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mflCode int) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeHelpers.MockFakeReportErrorToSentryFn = func(err error) {}

				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mflCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case - empty MFL Code" {
				fakeHelpers.MockFakeReportErrorToSentryFn = func(err error) {}

				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mflCode int) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}

			if tt.name == "Sad case - very bad MFL Code" {
				fakeHelpers.MockFakeReportErrorToSentryFn = func(err error) {}

				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mflCode int) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}

			_, err := d.DeleteFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_DeleteClientProfile(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case",
			args: args{
				ctx:      ctx,
				clientID: "123456789",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				clientID: "123456789",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockDeleteClientProfileFn = func(ctx context.Context, clientID string) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}
			got, err := d.DeleteClientProfile(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.DeleteClientProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_DeleteUser(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: "123456789",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: "123456789",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockDeleteUserFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}
			got, err := d.DeleteUser(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				fakeGorm.MockDeleteStaffProfileFn = func(ctx context.Context, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}
			got, err := d.DeleteStaffProfile(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.DeleteStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.DeleteStaffProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
