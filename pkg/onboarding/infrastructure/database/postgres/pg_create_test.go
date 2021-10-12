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

func TestOnboardingDb_CreateFacility(t *testing.T) {
	ctx := context.Background()
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"
	type args struct {
		ctx      context.Context
		facility *dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm)
			got, err := d.CreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.CreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.CreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return facility, nil
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.CreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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
