package facility_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func InitializeTestService(ctx context.Context) *interactor.Interactor {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil
	}

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)
	externalExt := externalExtension.NewExternalMethodsImpl()

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db)

	// Initialize user usecase
	userUsecase := user.NewUseCasesUserImpl(db, db, db, externalExt)

	i := interactor.NewMyCareHubInteractor(facilityUseCase, userUsecase)
	return i
}

func TestUseCaseFacilityImpl_CreateFacility(t *testing.T) {

	ctx := context.Background()
	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := "Nairobi"
	description := "This is just for testing"

	i := InitializeTestService(ctx)

	type args struct {
		ctx      context.Context
		facility *dto.FacilityInput
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
			name: "happy case - idempotent, should not error when saving again",
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
					Code:        0,
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
			got, err := i.FacilityUsecase.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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

	_, err := i.FacilityUsecase.DeleteFacility(ctx, code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_RetrieveFacility_Integration(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	Code := rand.Intn(1000000)

	facilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        Code,
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	ID := facility.ID
	active := facility.Active

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
				id:       ID,
				isActive: active,
			},
			wantErr: false,
		},
		{
			name: "Sad case - no id",
			args: args{
				ctx:      ctx,
				id:       nil,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.FacilityUsecase.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
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

	_, err = i.FacilityUsecase.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_DeleteFacility_Integrationtest(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	//Create facility
	facilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        rand.Intn(1000000),
		County:      "Nairobi",
		Active:      true,
		Description: "This is just for integration testing",
	}

	// create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	assert.Nil(t, err)
	assert.NotNil(t, facility)

	// retrieve the facility
	facility1, err := i.FacilityUsecase.RetrieveFacility(ctx, facility.ID, true)
	assert.Nil(t, err)
	assert.NotNil(t, facility1)

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  facility.Code,
			},
			wantErr: false,
		},
		{
			name: "Sad case - bad id",
			args: args{
				ctx: ctx,
				id:  8888,
			},
			wantErr: false,
		},
		{
			name: "Sad case - empty id",
			args: args{
				ctx: ctx,
				id:  0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := i.FacilityUsecase.DeleteFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_RetrieveFacilityByMFLCode_Integration(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	code := rand.Intn(1000000)
	facilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        code,
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflCode := facility.Code

	invalidMFLCode := 999999

	type args struct {
		ctx      context.Context
		MFLCode  int
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.FacilityUsecase.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
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

	_, err = i.FacilityUsecase.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_ListFacilities(t *testing.T) {
	ctx := context.Background()

	f := InitializeTestService(ctx)

	code := rand.Intn(1000000)
	code2 := rand.Intn(1000000)
	name := ksuid.New().String()

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for mocking",
	}

	facilityInput2 := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        code2,
		Active:      true,
		County:      "Baringo",
		Description: "This is just for mocking",
	}

	noSearchTerm := ""
	searchTerm := "ro"

	sortInput := dto.SortsInput{
		Field:     enums.FilterSortDataTypeCreatedAt,
		Direction: enums.SortDataTypeAsc,
	}

	invalidSortInput := dto.SortsInput{
		Field:     "invalid",
		Direction: enums.SortDataTypeAsc,
	}

	noFilterInput := []*dto.FiltersInput{}

	formatBool := strconv.FormatBool(true)

	filterInput := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    name,
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(code),
		},
		{
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterEmptyName := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    "",
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(code),
		},
		{
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}
	filterEmptyMFLCode := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    name,
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    "",
		},
		{
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidBool := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    name,
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(code),
		},
		{
			DataType: enums.FilterSortDataTypeActive,
			Value:    "invalid",
		},
		{
			DataType: enums.FilterSortDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidCounty := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    name,
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(code),
		},
		{
			DataType: enums.FilterSortDataTypeActive,
			Value:    formatBool,
		},
		{
			DataType: enums.FilterSortDataTypeCounty,
			Value:    "kanairo",
		},
	}

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}
	paginationInputWithSort := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
		Sort:        sortInput,
	}

	paginationInputInvalidSort := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
		Sort:        invalidSortInput,
	}
	paginationInputNoCurrentPage := dto.PaginationsInput{
		Limit: 1,
	}

	// Setup
	// create a facility
	facility, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}
	// Create another Facility
	facility2, err := f.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput2)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
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
		want    *domain.FacilityPage
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid filters",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      noFilterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				paginationsInput: &paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm and filter",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm and sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      noFilterInput,
				paginationsInput: &paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid sort and filter",
			args: args{
				ctx:              ctx,
				filterInput:      filterInput,
				searchTerm:       &noSearchTerm,
				paginationsInput: &paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "valid: with valid searchterm, filter and sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInputWithSort,
			},
			wantErr: false,
		},
		{
			name: "invalid: with invalid sort",
			args: args{
				ctx:              ctx,
				searchTerm:       &noSearchTerm,
				filterInput:      noFilterInput,
				paginationsInput: &paginationInputInvalidSort,
			},
			wantErr: true,
		},
		{
			name: "invalid: no params passed",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing current page",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInputNoCurrentPage,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty name passed",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyName,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty MFL code",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterEmptyMFLCode,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid bool",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidBool,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},

		{
			name: "invalid: invalid county",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInvalidCounty,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := f.FacilityUsecase.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// assert we get at least one value in the filter
			if tt.name == "valid: with valid filters" {
				assert.GreaterOrEqual(t, len(got.Facilities), 1)
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, got.Pagination.Limit, tt.args.paginationsInput.Limit)
				assert.Equal(t, got.Pagination.Limit, len(got.Facilities))
			}

		})
	}
	// Teardown
	_, err = f.FacilityUsecase.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
	_, err = f.FacilityUsecase.DeleteFacility(ctx, facility2.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_Inactivate_Integration_test(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	facilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        rand.Intn(1000000),
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for testing",
	}

	// Setup, create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &facility.Code,
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := i.FacilityUsecase.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.Inactivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
	_, err = i.FacilityUsecase.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestUseCaseFacilityImpl_Reactivate_Integration_test(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	facilityInput := &dto.FacilityInput{
		Name:        ksuid.New().String(),
		Code:        rand.Intn(1000000),
		Active:      false,
		County:      "Nairobi",
		Description: "This is just for testing",
	}

	// Setup, create a facility
	facility, err := i.FacilityUsecase.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &facility.Code,
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := i.FacilityUsecase.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
	_, err = i.FacilityUsecase.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}
