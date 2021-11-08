package gorm_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func createInactiveTestFacility() *gorm.Facility {
	ID := uuid.New().String()
	name := ksuid.New().String()
	code := uuid.New().String()
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      strconv.FormatBool(false),
		County:      county,
		Description: description,
	}

	return facility
}

func TestPGInstance_InactivateFacility(t *testing.T) {

	ctx := context.Background()

	testID := uuid.New().String()
	testFacility := &gorm.Facility{
		FacilityID:  &testID,
		Name:        ksuid.New().String(),
		Code:        ksuid.New().String(),
		Active:      strconv.FormatBool(true),
		County:      enums.CountyTypeNairobi,
		Description: "test description",
	}

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

func TestPGInstance_UpdateFacility(t *testing.T) {
	ctx := context.Background()

	testFacility := createTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	updateFacilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        ksuid.New().String(),
		Active:      false,
		County:      enums.CountyTypeBaringo,
		Description: "test description",
	}

	type args struct {
		ctx           context.Context
		id            *string
		facilityInput *dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: successful update",
			args: args{
				id:            facility.FacilityID,
				facilityInput: updateFacilityInput,
			},
			wantErr: false,
			want:    true,
		},

		{
			name: "valid: update active only",
			args: args{
				id: facility.FacilityID,
				facilityInput: &dto.FacilityInput{
					Active: false,
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: update county only",
			args: args{
				id: facility.FacilityID,
				facilityInput: &dto.FacilityInput{
					County: enums.CountyTypeKakamega,
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: update description only",
			args: args{
				id: facility.FacilityID,
				facilityInput: &dto.FacilityInput{
					Description: "test description2",
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: update code only",
			args: args{
				id: facility.FacilityID,
				facilityInput: &dto.FacilityInput{
					Code: ksuid.New().String(),
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: update name only",
			args: args{
				id: facility.FacilityID,
				facilityInput: &dto.FacilityInput{
					Name: ksuid.New().String(),
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: facility is nil",
			args: args{
				id:            facility.FacilityID,
				facilityInput: nil,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateFacility(tt.args.ctx, tt.args.id, tt.args.facilityInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ReactivateFacility(t *testing.T) {

	ctx := context.Background()

	testFacility := createInactiveTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility: %v", err)
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
}
