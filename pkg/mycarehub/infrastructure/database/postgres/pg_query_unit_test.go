package postgres

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
)

func TestMyCareHubDb_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	id := ksuid.New().String()

	invalidID := uuid.New().String()

	type args struct {
		ctx    context.Context
		id     *string
		active bool
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
				ctx:    ctx,
				id:     &id,
				active: true,
			},
			wantErr: false,
		},
		{
			name: "sad case - no ID passed",
			args: args{
				ctx:    ctx,
				active: false,
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid ID",
			args: args{
				ctx:    ctx,
				id:     &invalidID,
				active: false,
			},
			wantErr: true,
		},

		{
			name: "sad case - nil ID",
			args: args{
				ctx:    ctx,
				id:     nil,
				active: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case - no ID passed" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - invalid ID" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - nil ID" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			got, err := d.RetrieveFacility(ctx, tt.args.id, tt.args.active)

			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_GetFacilities(t *testing.T) {
	ctx := context.Background()

	id := uuid.New().String()
	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &domain.Facility{
		ID:          &id,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	var facilityData []*domain.Facility
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
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case - facility want data not given" {
				fakeGorm.MockGetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}

			got, err := d.GetFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetFacilities() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_RetrieveFacilityByMFLCode(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	// Setup, create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflCode := facility.Code

	invalidMFLCode := ksuid.New().String()

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

		{
			name: "Sad case - nil mfl code",
			args: args{
				ctx:      ctx,
				MFLCode:  "",
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}

			if tt.name == "Sad case - nil mfl code" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}

			got, err := d.RetrieveFacilityByMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RetrieveFacilityByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
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

func TestOnboardingDb_ListFacilities(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	code := ksuid.New().String()
	code2 := ksuid.New().String()

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        code,
		Active:      true,
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	facilityInput2 := &dto.FacilityInput{
		Name:        "Baringo 2",
		Code:        code2,
		Active:      true,
		County:      enums.CountyTypeBaringo,
		Description: "This is just for mocking",
	}

	noSearchTerm := ""
	searchTerm := "ro"

	noFilterInput := []*dto.FiltersInput{}

	formatBool := strconv.FormatBool(true)

	filterInput := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
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
			Value:    code,
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
			Value:    "Kanairo One",
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
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
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
			Value:    "Kanairo One",
		},
		{
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    code,
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
	paginationInputNoCurrentPage := dto.PaginationsInput{
		Limit: 1,
	}

	// Setup
	// create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}
	// Create another Facility
	facility2, err := d.GetOrCreateFacility(ctx, facilityInput2)
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
			name: "Sad case",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
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
			if tt.name == "Sad case" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return &gorm.Facility{
						FacilityID:  facility.FacilityID,
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      facility.Active,
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("current page not provided")
				}
			}

			if tt.name == "invalid: missing current page" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}
			if tt.name == "invalid: missing current page" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}
			if tt.name == "invalid: empty name passed" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}
			if tt.name == "invalid: empty MFL code" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}
			if tt.name == "invalid: invalid bool" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}
			if tt.name == "invalid: invalid county" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}

			}

			got, err := d.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
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
	// Teardown
	_, err = d.DeleteFacility(ctx, string(facility.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
	_, err = d.DeleteFacility(ctx, string(facility2.Code))
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
}

func TestOnboardingDb_GetUserProfileByPhoneNumber(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case- Successfully get a user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get user profile by phonenumber" {
				fakeGorm.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phonenumber")
				}
			}

			if tt.name == "Sad Case - Missing phonenumber" {
				fakeGorm.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
					return nil, fmt.Errorf("phone number should be provided")
				}
			}

			got, err := d.GetUserProfileByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetUserProfileByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestOnboardingDb_GetUserPINByUserID(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user pin by user ID",
			args: args{
				ctx:    ctx,
				userID: "1234456",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user pin",
			args: args{
				ctx:    ctx,
				userID: "12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get user pin" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}

			got, err := d.GetUserPINByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetUserPINByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetCurrentTerms(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		flavour enums.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				flavour: enums.PRO,
			},
			wantErr: false,
		},
		{
			name: "Sad case - empty flavour",
			args: args{
				ctx:     ctx,
				flavour: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case - empty flavour" {
				fakeGorm.MockGetCurrentTermsFn = func(ctx context.Context, flavour enums.Flavour) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetCurrentTerms(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
