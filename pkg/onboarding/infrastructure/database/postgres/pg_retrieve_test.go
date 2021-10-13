package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm/mock"
)

func TestOnboardingDb_RetrieveFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm)

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        "KN001",
		County:      "Kanairo",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := d.CreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	id := facility.ID

	invalidID := int64(-100)

	type args struct {
		ctx context.Context
		id  *int64
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
				ctx: ctx,
				id:  &id,
			},
			wantErr: false,
		},
		{
			name: "sad case - no ID passed",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid ID",
			args: args{
				ctx: ctx,
				id:  &invalidID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.RetrieveFacility(ctx, tt.args.id)

			if tt.name == "happy case - valid ID passed" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *int64) (*gorm.Facility, error) {
					return &gorm.Facility{
						FacilityID:  &facility.ID,
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      facility.Active,
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
			}

			if tt.name == "sad case - no ID passed" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *int64) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - invalid ID" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *int64) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

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
