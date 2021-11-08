package facility_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
	"github.com/segmentio/ksuid"
)

func TestUnit_CreateFacility(t *testing.T) {
	ctx := context.Background()
	name := "Kanairo One"
	code := ksuid.New().String()
	county := enums.CountyTypeNairobi
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
			name: "sad case - facility code empty",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
					Name:        name,
					Code:        "",
					Active:      true,
					County:      enums.CountyTypeNairobi,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case - facility code empty" {
				fakeFacility.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			got, err := f.GetOrCreateFacility(tt.args.ctx, &tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected facility to be nil for %v", tt.name)
				return
			}
		})
	}
	// TODO: add teardown
}

func TestUseCaseFacilityImpl_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()

	type args struct {
		ctx      context.Context
		id       *string
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
			fakeFacility := mock.NewFacilityUsecaseMock()

			fakeDB := pgMock.NewPostgresMock()
			facility := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case - no id" {
				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
				}
			}
			got, err := facility.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
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
				MFLCode:  "",
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case" {

				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
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

func TestUnit_ListFacilities(t *testing.T) {
	ctx := context.Background()

	searchTerm := "term"

	filterValue := "value"

	filterInput := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    filterValue,
		},
	}

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case- empty search term",
			args: args{
				ctx:              ctx,
				searchTerm:       nil,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- nil filter input",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- nil pagination input",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case- empty search term" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}
			}

			if tt.name == "Sad case- nil filter input" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}
			}

			if tt.name == "Sad case- nil pagination input" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}
			}

			got, err := f.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCaseFacilityImpl_Inactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	validMFLCode := ksuid.New().String()
	veryBadMFLCode := ksuid.New().String() + gofakeit.HipsterSentence(500)
	emptyMFLCode := ""

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
				mflCode: &emptyMFLCode,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx:     ctx,
				mflCode: &validMFLCode,
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
			fakeDB := pgMock.NewPostgresMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.Inactivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUseCaseFacilityImpl_Reactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	validMFLCode := ksuid.New().String()
	veryBadMFLCode := ksuid.New().String() + gofakeit.HipsterSentence(500)
	emptyMFLCode := ""

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
				mflCode: &emptyMFLCode,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx:     ctx,
				mflCode: &validMFLCode,
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
			fakeDB := pgMock.NewPostgresMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, mflCode *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUnit_UpdateFacility(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()

	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB)

	validID := ksuid.New().String()
	longText := ksuid.New().String() + gofakeit.HipsterSentence(500)
	emptyText := ""

	updateFacilityInput := &dto.FacilityInput{
		Name:        "updated Name",
		Code:        "1234",
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
				ctx:           ctx,
				id:            &validID,
				facilityInput: updateFacilityInput,
			},
			wantErr: false,
			want:    true,
		},

		{
			name: "valid: update active only",
			args: args{
				ctx: ctx,
				id:  &validID,
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
				ctx: ctx,
				id:  &validID,
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
				ctx: ctx,
				id:  &validID,
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
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Code: "12345",
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: update name only",
			args: args{
				ctx: ctx,
				id:  &validID,
				facilityInput: &dto.FacilityInput{
					Name: "test name",
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: nil facility input",
			args: args{
				ctx:           ctx,
				id:            &validID,
				facilityInput: nil,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: long ID",
			args: args{
				ctx:           ctx,
				id:            &longText,
				facilityInput: updateFacilityInput,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: empty ID",
			args: args{
				ctx:           ctx,
				id:            &emptyText,
				facilityInput: updateFacilityInput,
			},
			wantErr: true,
			want:    false,
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
			if tt.name == "invalid: nil facility input" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long ID" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: empty ID" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}

			if tt.name == "invalid: nil ID" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: invalid ID input" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: short name" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long name" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long code" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: long description" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "invalid: short description" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
					return false, fmt.Errorf("failed to update facility")
				}
			}

			got, err := f.UpdateFacility(tt.args.ctx, tt.args.id, tt.args.facilityInput)
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
