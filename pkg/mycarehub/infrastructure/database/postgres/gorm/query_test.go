package gorm_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func createTestFacility() *gorm.Facility {
	ID := uuid.New().String()
	name := ksuid.New().String()
	code := uuid.New().String()
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      strconv.FormatBool(true),
		County:      county,
		Description: description,
	}

	return facility
}

func TestPGInstance_RetrieveFacility(t *testing.T) {
	ctx := context.Background()
	fakeID := "1234"

	testFacility := createTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

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
			name: "Happy Case - Successfully retrieve created facility",
			args: args{
				ctx:      ctx,
				id:       facility.FacilityID,
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to retrieve facility",
			args: args{
				ctx:      ctx,
				id:       &fakeID,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestPGInstance_RetrieveFacilityByMFLCode(t *testing.T) {
	ctx := context.Background()
	testFacility := createTestFacility()
	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

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
			name: "Happy Case - Successfully fetch facility by MFL code",
			args: args{
				ctx:      ctx,
				MFLCode:  facility.Code,
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to fetch facility by MFL code",
			args: args{
				ctx:      ctx,
				MFLCode:  "invalid code",
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RetrieveFacilityByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestPGInstance_ListFacilities(t *testing.T) {
	ctx := context.Background()

	d := testingDB

	code := ksuid.New().String()
	code2 := ksuid.New().String()

	facilityInput := &gorm.Facility{
		Name:        ksuid.New().String(),
		Code:        code,
		Active:      strconv.FormatBool(true),
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	facilityInput2 := &gorm.Facility{
		Name:        ksuid.New().String(),
		Code:        code2,
		Active:      strconv.FormatBool(true),
		County:      enums.CountyTypeBaringo,
		Description: "This is just for mocking",
	}

	noSearchTerm := ""
	searchTerm := "ro"

	noFilterInput := []*domain.FiltersParam{}

	formatBool := strconv.FormatBool(true)

	filterInput := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterSortDataTypeActive.String(),
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterSortDataTypeCounty.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterEmptyName := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    "",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterSortDataTypeActive.String(),
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterSortDataTypeCounty.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}
	filterEmptyMFLCode := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    "",
		},
		{
			Name:     enums.FilterSortDataTypeActive.String(),
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterSortDataTypeCounty.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidBool := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterSortDataTypeActive.String(),
			DataType: enums.FilterSortDataTypeActive,
			Value:    "invalid",
		},
		{
			Name:     enums.FilterSortDataTypeCounty.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidCounty := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterSortDataTypeActive.String(),
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterSortDataTypeCounty.String(),
			DataType: enums.FilterSortDataTypeCounty,
			Value:    "Kanairo",
		},
	}

	noSortValues := domain.SortParam{
		// Field:     enums.FilterSortDataTypeCreatedAt,
		// Direction: enums.SortDataTypeAsc,
	}

	sortParams := domain.SortParam{
		Field:     enums.FilterSortDataTypeCreatedAt,
		Direction: enums.SortDataTypeAsc,
	}
	invalidSortInput := domain.SortParam{
		Field:     "invalid",
		Direction: enums.SortDataTypeAsc,
	}

	paginationInput := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       1,
			CurrentPage: 1,
			Sort:        &noSortValues,
		},
	}
	paginationInputWithSort := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       1,
			CurrentPage: 1,
			Sort:        &sortParams,
		},
	}

	paginationInputInvalidSort := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       1,
			CurrentPage: 1,
			Sort:        &invalidSortInput,
		},
	}

	// Setup
	// create a facility
	_, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}
	// Create another Facility
	_, err = d.GetOrCreateFacility(ctx, facilityInput2)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*domain.FiltersParam
		PaginationsInput domain.FacilityPage
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
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},

		{
			name: "valid: with valid filters",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      filterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},

		{
			name: "valid: with valid searchterm",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with no sort values",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with  sort values",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInputWithSort,
			},
			wantErr: false,
		},

		{
			name: "valid: with valid searchterm and filter",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				PaginationsInput: paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm and sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid sort and filter",
			args: args{
				ctx:              ctx,
				filterInput:      filterInput,
				searchTerm:       &noSearchTerm,
				PaginationsInput: paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm, filter and sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				PaginationsInput: paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "invalid: with invalid sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				PaginationsInput: paginationInputInvalidSort,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty name passed",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyName,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty MFL code",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyMFLCode,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid bool",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidBool,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},

		{
			name: "invalid: invalid county",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidCounty,
				PaginationsInput: paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, &tt.args.PaginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}

}
