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
				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed query and retrieve facility by MFLCode")
				}

				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					return &domain.Facility{
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      facility.Active,
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
			}

			got, err := f.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCaseFacilityImpl_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	f := testFakeInfrastructureInteractor

	ID := uuid.New()

	type args struct {
		ctx      context.Context
		id       *uuid.UUID
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
				id:       &ID,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				id:       &ID,
				isActive: true,
			},
			wantErr: true,
		},

		{
			name: "Sad case - no id",
			args: args{
				ctx:      ctx,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					ID := uuid.New()
					return &domain.Facility{
						ID:          ID,
						Name:        "facility.Name",
						Code:        "facility.Code",
						Active:      true,
						County:      "facility.County",
						Description: "facility.Description",
					}, nil
				}

				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
					ID := uuid.New()
					return &domain.Facility{
						ID:          ID,
						Name:        "facility.Name",
						Code:        "facility.Code",
						Active:      true,
						County:      "facility.County",
						Description: "facility.Description",
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
				}
			}

			if tt.name == "Sad case - no id" {
				fakeQuery.RetrieveFacilityFn = func(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
				}
			}
			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCaseFacilityImpl_RetrieveFacilityByMFLCode_Unittest(t *testing.T) {
	ctx := context.Background()

	f := testFakeInfrastructureInteractor

	MFLCode := ksuid.New().String()

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
				MFLCode:  MFLCode,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				MFLCode:  MFLCode,
				isActive: true,
			},
			wantErr: true,
		},

		{
			name: "Sad case#1",
			args: args{
				ctx:      ctx,
				MFLCode:  MFLCode,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					ID := uuid.New()
					return &domain.Facility{
						ID:          ID,
						Name:        "facility.Name",
						Code:        "facility.Code",
						Active:      true,
						County:      "facility.County",
						Description: "facility.Description",
					}, nil
				}

				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
					ID := uuid.New()
					return &domain.Facility{
						ID:          ID,
						Name:        "facility.Name",
						Code:        "facility.Code",
						Active:      true,
						County:      "facility.County",
						Description: "facility.Description",
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					ID := uuid.New()
					return &domain.Facility{
						ID:          ID,
						Name:        "facility.Name",
						Code:        "facility.Code",
						Active:      true,
						County:      "facility.County",
						Description: "facility.Description",
					}, nil
				}

				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFLCode")
				}
			}

			if tt.name == "Sad case#1" {
				fakeQuery.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFLCode")
				}
			}
			got, err := f.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
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
