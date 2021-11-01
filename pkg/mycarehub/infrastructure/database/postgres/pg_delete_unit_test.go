package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestOnboardingDb_DeleteFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
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

	type args struct {
		ctx context.Context
		id  string
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
				id:  "mflcode",
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Sad case - empty MFL Code",
			args: args{
				ctx: ctx,
				id:  "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Happy case" {
				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mfl_code string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mfl_code string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case - empty MFL Code" {
				fakeGorm.MockDeleteFacilityFn = func(ctx context.Context, mfl_code string) (bool, error) {
					return false, fmt.Errorf("an error occurred while deleting")
				}
			}

			_, err := d.DeleteFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
