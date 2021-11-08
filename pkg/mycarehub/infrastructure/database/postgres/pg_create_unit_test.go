package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	invalidFacility := &dto.FacilityInput{
		Name:        name,
		Active:      true,
		County:      county,
		Description: description,
	}

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
				ctx:      ctx,
				facility: facility,
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx:      ctx,
				facility: invalidFacility,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get ot create facility",
			args: args{
				ctx:      ctx,
				facility: invalidFacility,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm)
			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - nil facility input" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "Sad Case - Fail to get ot create facility" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get or create facility")
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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
