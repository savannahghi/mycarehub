package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm/mock"
)

func TestOnboardingDb_GetFacilities(t *testing.T) {
	ctx := context.Background()
	id := int64(1)
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"

	facility := &domain.Facility{
		ID:          id,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	facilityData := []*domain.Facility{}
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
				fakeGorm.GetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.GetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
					var facilities []gorm.Facility
					facilityID := int64(1)
					name := "Kanairo One"
					code := "KN001"
					county := "Kanairo"
					description := "This is just for mocking"
					facilities = append(facilities, gorm.Facility{
						FacilityID:  &facilityID,
						Name:        name,
						Code:        code,
						Active:      true,
						County:      county,
						Description: description,
					})
					return facilities, nil
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
