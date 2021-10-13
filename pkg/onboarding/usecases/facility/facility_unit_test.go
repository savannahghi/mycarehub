package facility_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
)

func TestUnit_CreateFacility(t *testing.T) {
	f := testFakeInfrastructureInteractor
	ctx := context.Background()
	name := "Kanairo One"
	code := ksuid.New().String()
	county := "Kanairo"
	description := "This is just for mocking"

	type args struct {
		ctx      context.Context
		facility dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
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
				facility: dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case - facility code not defined" {
				fakeCreate.CreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeCreate.CreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					return &domain.Facility{
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      facility.Active,
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
			}

			got, err := f.CreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
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
	// TODO: add teardown
}

func TestUnit_RetrieveFacility(t *testing.T) {
	f := testFakeInfrastructureInteractor

	ctx := context.Background()

	id := uuid.New()
	facility := &domain.Facility{
		ID:          id,
		Name:        "test-name",
		Code:        "c123",
		Active:      true,
		County:      "test",
		Description: "test description",
	}

	invalidID := uuid.New()

	type args struct {
		ctx context.Context
		id  *uuid.UUID
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
			want:    facility,
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
			if tt.name == "happy case - valid ID passed" {
				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
					return facility, nil
				}
			}
			if tt.name == "sad case - no ID passed" {
				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - invalid ID" {
				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id)

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
	// TODO: add teardown
}
