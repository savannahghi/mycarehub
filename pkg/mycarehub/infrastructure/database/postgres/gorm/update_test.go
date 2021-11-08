package gorm_test

import (
	"context"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestPGInstance_InactivateFacility(t *testing.T) {
	ctx := context.Background()

	testFacility := createTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *string
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
				mflCode: &facility.Code,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TearDown(t *testing.T) {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return
	}

	pg.DB.Migrator().DropTable(&gorm.Facility{})
}

func TestPGInstance_ReactivateFacility(t *testing.T) {

	ctx := context.Background()

	inactiveFacility := createInactiveTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, inactiveFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *string
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
				mflCode: &facility.Code,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}

	TearDown(t)
}
