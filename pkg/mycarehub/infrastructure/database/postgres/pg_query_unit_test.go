package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
)

func TestOnboardingDb_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

	id := ksuid.New().String()

	invalidID := uuid.New().String()

	type args struct {
		ctx    context.Context
		id     *string
		active bool
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid ID passed",
			args: args{
				ctx:    ctx,
				id:     &id,
				active: true,
			},
			wantErr: false,
		},
		{
			name: "sad case - no ID passed",
			args: args{
				ctx:    ctx,
				active: false,
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid ID",
			args: args{
				ctx:    ctx,
				id:     &invalidID,
				active: false,
			},
			wantErr: true,
		},

		{
			name: "sad case - nil ID",
			args: args{
				ctx:    ctx,
				id:     nil,
				active: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case - no ID passed" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - invalid ID" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - nil ID" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			got, err := d.RetrieveFacility(ctx, tt.args.id, tt.args.active)

			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facility to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facility not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_GetFacilities(t *testing.T) {
	ctx := context.Background()

	id := uuid.New().String()
	name := gofakeit.Name()
	code := "KN001"
	county := "Kanairo"
	description := gofakeit.HipsterSentence(15)

	facility := &domain.Facility{
		ID:          &id,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	var facilityData []*domain.Facility
	facilityData = append(facilityData, facility)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Facility
		wantErr bool
	}{
		{
			name:    "happy case - valid payload",
			args:    args{ctx: ctx},
			want:    facilityData,
			wantErr: false,
		},
		{
			name:    "sad case - facility want data not given",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case - facility want data not given" {
				fakeGorm.MockGetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}

			got, err := d.GetFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_RetrieveByFacilityMFLCode(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

	name := gofakeit.Name()
	code := "KN001"
	county := "Kanairo"
	description := gofakeit.HipsterSentence(15)

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	// Setup, create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflCode := facility.Code

	invalidMFLCode := ksuid.New().String()

	type args struct {
		ctx      context.Context
		MFLCode  string
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      ctx,
				MFLCode:  mflCode,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				MFLCode:  invalidMFLCode,
				isActive: true,
			},
			wantErr: true,
		},

		{
			name: "Sad case - nil mfl code",
			args: args{
				ctx:      ctx,
				MFLCode:  "",
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}

			if tt.name == "Sad case - nil mfl code" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}

			got, err := d.RetrieveByFacilityMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RetrieveByFacilityMFLCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}
