package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
)

func TestMyCareHubDb_InactivateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	validMFLCode := ksuid.New().String()
	veryBadMFLCode := ksuid.New().String() + gofakeit.HipsterSentence(500)

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
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
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

	validMFLCode := ksuid.New().String()
	veryBadMFLCode := ksuid.New().String() + gofakeit.HipsterSentence(500)

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
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
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

func TestMyCareHubDb_UpdateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	validID := ksuid.New().String()
	longText := ksuid.New().String() + gofakeit.HipsterSentence(500)

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
			name: "valid: with valid values",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        gofakeit.BeerHop(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: nil ID",
			args: args{
				ctx: ctx,
				id:  nil,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        gofakeit.BeerHop(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid ID input",
			args: args{
				ctx: ctx,
				id:  &longText,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        gofakeit.BeerHop(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: short name",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        "te",
					Code:        gofakeit.BeerHop(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: long name",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        longText,
					Code:        gofakeit.BeerHop(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: long code",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        longText,
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "This is just for mocking",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: long description",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        ksuid.New().String(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: longText,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: short description",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name:        "test name",
					Code:        ksuid.New().String(),
					Active:      true,
					County:      enums.CountyTypeBaringo,
					Description: "sh",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "invalid: nil ID" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: invalid ID input" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: short name" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long name" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long code" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long description" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: short description" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}

			got, err := d.UpdateFacility(tt.args.ctx, tt.args.id, tt.args.facilityInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.UpdateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}
