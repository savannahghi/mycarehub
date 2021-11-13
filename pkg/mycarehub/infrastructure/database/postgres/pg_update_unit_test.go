package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
)

func TestMyCareHubDb_InactivateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	validMFLCode := gofakeit.Number(0, 100)
	veryBadMFLCode := gofakeit.Number(10000, 10000000)

	type args struct {
		ctx     context.Context
		mflCode *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				mflCode: &validMFLCode,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx:     ctx,
				mflCode: &veryBadMFLCode,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - empty mflCode" {
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, mflCode *int) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, mflCode *int) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}

			got, err := d.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_ReactivateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	validMFLCode := gofakeit.Number(0, 100)
	veryBadMFLCode := gofakeit.Number(10000, 10000000)

	type args struct {
		ctx     context.Context
		mflCode *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				mflCode: &validMFLCode,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx:     ctx,
				mflCode: &veryBadMFLCode,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - empty mflCode" {
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, mflCode *int) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, mflCode *int) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}

			got, err := d.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_AcceptTerms(t *testing.T) {
	ctx := context.Background()

	userID := ksuid.New().String()
	termsID := gofakeit.Number(0, 100000)
	negativeTermsID := gofakeit.Number(-100000, -1)

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
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
				userID:  &userID,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - negative termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - userID and negative termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - negative termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - userID and negative termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}
