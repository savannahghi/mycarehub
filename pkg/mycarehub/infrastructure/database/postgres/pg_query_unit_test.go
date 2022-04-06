package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
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
	code := gofakeit.Number(0, 100)
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()

	facility := &domain.Facility{
		ID:                 &id,
		Name:               name,
		Code:               code,
		Active:             true,
		County:             county,
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
	}

	var facilityData []*domain.Facility
	facilityData = append(facilityData, facility)

	invalidFacility := &domain.Facility{
		ID:          &id,
		Name:        name,
		Active:      true,
		County:      county,
		Description: description,
	}

	var invalidFacilityData []*domain.Facility
	invalidFacilityData = append(invalidFacilityData, invalidFacility)

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
			name:    "sad case - invalid payload",
			args:    args{ctx: ctx},
			want:    invalidFacilityData,
			wantErr: true,
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
			if tt.name == "sad case - invalid payload" {
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
	code := gofakeit.Number(0, 100)
	phone := "+254711223344"
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)

	facilityInput := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Phone:       phone,
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

	invalidMFLCode := 23456

	type args struct {
		ctx      context.Context
		MFLCode  int
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
				MFLCode:  0,
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}

			if tt.name == "Sad case - nil mfl code" {
				fakeGorm.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
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

	code := gofakeit.Number(0, 100)
	code2 := gofakeit.Number(200, 300)

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        code,
		Phone:       "+254711223344",
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for mocking",
	}

	facilityInput2 := &dto.FacilityInput{
		Name:        "Baringo 2",
		Code:        code2,
		Phone:       "+254711223355",
		Active:      true,
		County:      "Baringo",
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
			Value:    "Kanairo One",
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
	_, err = d.DeleteFacility(ctx, facility.Code)
	if err != nil {
		t.Errorf("unable to delete facility")
		return
	}
	_, err = d.DeleteFacility(ctx, facility2.Code)
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
		flavour     feedlib.Flavour
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
				fakeGorm.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phonenumber")
				}
			}

			if tt.name == "Sad Case - Missing phonenumber" {
				fakeGorm.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error) {
					return nil, fmt.Errorf("phone number should be provided")
				}
			}

			got, err := d.GetUserProfileByPhoneNumber(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
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
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user pin by user ID",
			args: args{
				ctx:     ctx,
				userID:  "1234456",
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user pin",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - empty user id",
			args: args{
				ctx:     ctx,
				userID:  "",
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid-flavour",
			args: args{
				ctx:     ctx,
				userID:  "",
				flavour: "invalid-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get user pin" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}
			if tt.name == "Sad Case - empty user id" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}
			if tt.name == "Sad Case - invalid-flavour" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}

			got, err := d.GetUserPINByUserID(tt.args.ctx, tt.args.userID, tt.args.flavour)
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
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.TermsOfService
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				flavour: "invalid-flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - nil context",
			args: args{
				ctx:     nil,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockGetCurrentTermsFn = func(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - nil context" {
				fakeGorm.MockGetCurrentTermsFn = func(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
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

func TestOnboardingDb_GetUserProfileByUserID(t *testing.T) {
	ctx := context.Background()

	validUserID := ksuid.New().String()
	emptyUserID := ""

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user profile by user ID",
			args: args{
				ctx:    ctx,
				userID: validUserID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Passed in empty User ID",
			args: args{
				ctx:    ctx,
				userID: emptyUserID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by user ID",
			args: args{
				ctx:    ctx,
				userID: validUserID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Happy Case - Successfully get user profile by user ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return &gorm.User{
						UserID: &validUserID,
					}, nil
				}
			}
			if tt.name == "Sad Case - Passed in empty User ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("user ID should be provided")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by user ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
				}
			}

			got, err := d.GetUserProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetSecurityQuestions(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.SecurityQuestion
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavor",
			args: args{
				ctx:     ctx,
				flavour: "invalid-flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - nil flavor",
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

			if tt.name == "Sad case" {
				fakeGorm.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - invalid flavor" {
				fakeGorm.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - nil flavor" {
				fakeGorm.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetSecurityQuestions(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_GetSecurityQuestionByID(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()
	type args struct {
		ctx                context.Context
		securityQuestionID *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get security questions",
			args: args{
				ctx:                ctx,
				securityQuestionID: &ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - Fail to get security question",
			args: args{
				ctx:                ctx,
				securityQuestionID: &ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case - Fail to get security question" {
				fakeGorm.MockGetSecurityQuestionByIDFn = func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error) {
					return nil, fmt.Errorf("failed to get security question")
				}
			}

			got, err := d.GetSecurityQuestionByID(tt.args.ctx, tt.args.securityQuestionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSecurityQuestionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetSecurityQuestionResponse(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		questionID string
		userID     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Happy Case - Successfully get security question response",
			args:    args{ctx: ctx, questionID: "12345", userID: uuid.NewString()},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing question ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get security question response",
			args: args{
				ctx:        ctx,
				questionID: "12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Missing question ID" {
				fakeGorm.MockGetSecurityQuestionResponseFn = func(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get security question response")
				}
			}

			if tt.name == "Sad Case - Fail to get security question response" {
				fakeGorm.MockGetSecurityQuestionResponseFn = func(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get security question response")
				}
			}

			got, err := d.GetSecurityQuestionResponse(tt.args.ctx, tt.args.questionID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_VerifyOTP(t *testing.T) {
	ctx := context.Background()

	flavour := feedlib.FlavourConsumer

	validOTPPayload := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload1 := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload2 := &dto.VerifyOTPInput{
		PhoneNumber: "",
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload3 := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         "",
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         uuid.New().String(),
		Flavour:     "flavour",
	}
	invalidOTPPayload5 := &dto.VerifyOTPInput{
		PhoneNumber: "otpInput.PhoneNumber",
		OTP:         "otpInput.OTP",
		Flavour:     "flavour",
	}
	invalidOTPPayload6 := &dto.VerifyOTPInput{
		PhoneNumber: gofakeit.HipsterParagraph(1, 10, 100, ""),
		OTP:         gofakeit.HipsterParagraph(1, 10, 100, ""),
		Flavour:     "gofakeit.HipsterParagraph(300, 10, 100)",
	}

	type args struct {
		ctx     context.Context
		payload *dto.VerifyOTPInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				payload: validOTPPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no user ID",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload2,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no otp",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload3,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload4,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload5,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - extreme bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload6,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case - no user ID" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no phone" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no otp" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad flavour" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad inputs" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - extreme bad inputs" {
				fakeGorm.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.VerifyOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.VerifyOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfPhoneNumberExists(t *testing.T) {
	ctx := context.Background()

	phone := ksuid.New().String()
	veryBadPhone := gofakeit.HipsterSentence(200)

	type args struct {
		ctx       context.Context
		phone     string
		isOptedIn bool
		flavour   feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				phone:     phone,
				isOptedIn: true,
				flavour:   feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				phone:     phone,
				isOptedIn: false,
				flavour:   feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:       ctx,
				phone:     "",
				isOptedIn: false,
				flavour:   feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no phone no flavour",
			args: args{
				ctx:       ctx,
				phone:     "",
				isOptedIn: false,
				flavour:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no phone and invalid flavour",
			args: args{
				ctx:       ctx,
				phone:     "",
				isOptedIn: false,
				flavour:   "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no isOptedIn",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no isOptedIn and bad flavour",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: "feedlib.FlavourConsumer",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no isOptedIn and empty flavour",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no isOptedIn and empty flavour and phone",
			args: args{
				ctx:     ctx,
				phone:   "",
				flavour: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: "feedlib.FlavourConsumer",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and empty flavour",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedIn",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: true,
				flavour:   feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone, isOptedIn and invalid flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: true,
				flavour:   "feedlib.FlavourConsumer",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone, isOptedIn and empty flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: true,
				flavour:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedOut",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedOut and good flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: false,
				flavour:   feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedOut and bad flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: false,
				flavour:   "feedlib.FlavourConsumer",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedOut and empty flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: false,
				flavour:   "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and isOptedOut and very bad flavour",
			args: args{
				ctx:       ctx,
				phone:     veryBadPhone,
				isOptedIn: false,
				flavour:   feedlib.Flavour(gofakeit.HipsterParagraph(2, 40, 400, " ")),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case - no phone" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no phone no flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no phone and invalid flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no isOptedIn" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no isOptedIn and bad flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no isOptedIn and empty flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case - very bad phone" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and invalid flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and empty flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedIn" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone, isOptedIn and invalid flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone, isOptedIn and empty flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedOut" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedOut and good flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedOut and bad flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedOut and empty flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and isOptedOut and very bad flavour" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no isOptedIn and empty flavour and phone" {
				fakeGorm.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckIfPhoneNumberExists(tt.args.ctx, tt.args.phone, tt.args.isOptedIn, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfPhoneExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfPhoneExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetClientProfileByUserID(t *testing.T) {
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
			name: "Happy Case - Successfully get client profile by user ID",
			args: args{
				ctx:    ctx,
				userID: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:    ctx,
				userID: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing user ID",
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

			if tt.name == "Sad Case - Fail to get client profile" {
				fakeGorm.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by user ID")
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeGorm.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by user ID")
				}
			}

			got, err := d.GetClientProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetStaffProfileByUserID(t *testing.T) {
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
			name: "Happy Case - Successfully get staff profile by user ID",
			args: args{
				ctx:    ctx,
				userID: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get staff profile",
			args: args{
				ctx:    ctx,
				userID: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing user ID",
			args: args{
				ctx:    ctx,
				userID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get staff profile" {
				fakeGorm.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile" {
				fakeGorm.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeGorm.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			got, err := d.GetStaffProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchStaffProfileByStaffNumber(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx         context.Context
		staffNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to get staff profile(s)",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get user profile",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - fail to get staff profile(s)" {
				fakeGorm.MockSearchStaffProfileByStaffNumberFn = func(ctx context.Context, staffNumber string) ([]*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			if tt.name == "Sad case - fail to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchStaffProfileByStaffNumber(tt.args.ctx, tt.args.staffNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchStaffProfileByStaffNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff profiles to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected staff profiles not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_ListContentCategories(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ContentItemCategory
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
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

			if tt.name == "Sad case" {
				fakeGorm.MockListContentCategoriesFn = func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected content to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected content not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckUserHasPin(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing user ID",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: could not check user has pin",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "invalid: could not check user has pin" {
				fakeGorm.MockCheckUserHasPinFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckUserHasPin(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckUserHasPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckUserHasPin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     feedlib.FlavourConsumer,
			},
		},
		{
			name: "invalid:  no phone number",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid: no flavour",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
			},
			wantErr: true,
		},
		{
			name: "Invalid: invalid flavor",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     "invalid-flavour",
			},
			wantErr: true,
		},
		{
			name: "Invalid: failed to get otp",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Invalid: failed to get otp" {
				fakeGorm.MockGetOTPFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetOTP(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetUserSecurityQuestionsResponses(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.SecurityQuestionResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: gofakeit.UUID(),
			},
		},
		{
			name: "Happy case - no resposes",
			args: args{
				ctx:    ctx,
				userID: gofakeit.UUID(),
			},
		},
		{
			name: "invalid: could not find security questions",
			args: args{
				ctx:    ctx,
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "invalid: no userID",
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

			if tt.name == "Happy case - no resposes" {
				fakeGorm.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error) {
					return []*gorm.SecurityQuestionResponse{}, nil
				}
			}

			if tt.name == "invalid: could not find security questions" {
				fakeGorm.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get user security questions")
				}
			}

			got, err := d.GetUserSecurityQuestionsResponses(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserSecurityQuestionsResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetContactByUserID(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()
	type args struct {
		ctx         context.Context
		userID      *string
		contactType string
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				userID:      &ID,
				contactType: "EMAIL",
			},
		},
		{
			name: "invalid: no userID",
			args: args{
				ctx:         ctx,
				contactType: "EMAIL",
			},
			wantErr: true,
		},
		{
			name: "invalid: no contactType",
			args: args{
				ctx:         ctx,
				contactType: "",
			},
			wantErr: true,
		},
		{
			name: "Valid userID and invalid contactType",
			args: args{
				ctx:         ctx,
				userID:      &ID,
				contactType: "WHATSAPP",
			},
			wantErr: true,
		},
		{
			name: "Invalid contactType",
			args: args{
				ctx:         ctx,
				contactType: "WHATSAPP",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.GetContactByUserID(tt.args.ctx, tt.args.userID, tt.args.contactType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetContactByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckWhetherUserHasLikedContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: gofakeit.Number(1, 1001),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: -5,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty userID" {
				fakeGorm.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty contentID" {
				fakeGorm.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - invalid contentID" {
				fakeGorm.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckWhetherUserHasLikedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckWhetherUserHasLikedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckWhetherUserHasLikedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetUserBookmarkedContent(t *testing.T) {
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
			name: "Happy Case - Successfully get user bookmarked content",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get content",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get content" {
				fakeGorm.MockGetUserBookmarkedContentFn = func(ctx context.Context, userID string) ([]*gorm.ContentItem, error) {
					return nil, fmt.Errorf("failed to get bookmarked content")
				}
			}

			got, err := d.GetUserBookmarkedContent(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientHealthDiaryQuote(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client health diary quote",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client health diary quote",
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
			if tt.name == "Sad Case - Fail to get client health diary quote" {
				fakeGorm.MockGetClientHealthDiaryQuoteFn = func(ctx context.Context) (*gorm.ClientHealthDiaryQuote, error) {
					return nil, fmt.Errorf("failed to get client health diary quote")
				}
			}
			got, err := d.GetClientHealthDiaryQuote(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientHealthDiaryQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_CanRecordHeathDiary(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			got, err := d.CanRecordHeathDiary(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CanRecordHeathDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CanRecordHeathDiary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfUserBookmarkedContent(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       context.Background(),
				userID:    uuid.New().String(),
				contentID: 1,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockCheckIfUserBookmarkedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.CheckIfUserBookmarkedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfUserBookmarkedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetClientHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client health diary entries",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client health diary entries",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client health diary entries" {
				fakeGorm.MockGetClientHealthDiaryEntriesFn = func(ctx context.Context, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get client health diary entries")
				}
			}
			if tt.name == "Sad Case - Fail to get client profile" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}
			if tt.name == "Sad Case - Fail to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetFAQContent(t *testing.T) {
	ctx := context.Background()
	limit := 10
	type args struct {
		ctx     context.Context
		limit   *int
		flavour feedlib.Flavour
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get FAQ content",
			args: args{
				ctx:     ctx,
				limit:   &limit,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get FAQ content",
			args: args{
				ctx:     ctx,
				limit:   &limit,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad Case - Fail to get FAQ content" {
				fakeGorm.MockGetFAQContentFn = func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*gorm.FAQ, error) {
					return nil, fmt.Errorf("failed to get FAQ content")
				}
			}

			got, err := d.GetFAQContent(tt.args.ctx, tt.args.flavour, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetFAQContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientCaregiver(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Caregiver
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client caregiver",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},

			want: &domain.Caregiver{
				ID:            uuid.New().String(),
				FirstName:     "John",
				LastName:      "Doe",
				PhoneNumber:   "+1234567890",
				CaregiverType: enums.CaregiverTypeFather,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client caregiver",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client caregiver" {
				fakeGorm.MockGetClientCaregiverFn = func(ctx context.Context, clientID string) (*gorm.Caregiver, error) {
					return nil, fmt.Errorf("failed to get client caregiver")
				}
			}

			got, err := d.GetClientCaregiver(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetPendingServiceRequestsCount(t *testing.T) {
	ctx := context.Background()

	facilityID := uuid.New().String()

	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *int64
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get staff service requests count",
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

			if tt.name == "Sad case" {
				fakeGorm.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty facility ID" {
				fakeGorm.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - fail to get staff service requests count" {
				fakeGorm.MockGetStaffPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetPendingServiceRequestsCount(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetPendingServiceRequestsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetPendingServiceRequestsCount() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetPendingServiceRequestsCount() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientByClientID(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client by client ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client by client ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client by client ID" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client by client ID")
				}
			}
			got, err := d.GetClientProfileByClientID(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientProfileByClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetServiceRequests(t *testing.T) {
	var requeststatus = enums.ServiceRequestStatusPending.String()
	var requesttype = enums.ServiceRequestTypeRedFlag.String()
	facilityID := uuid.New().String()
	type args struct {
		ctx           context.Context
		requestType   *string
		requestStatus *string
		facilityID    string
		flavour       feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get service requests by type",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				flavour:       feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get service requests by status",
			args: args{
				ctx:           context.Background(),
				requestStatus: &requeststatus,
				facilityID:    facilityID,
				flavour:       feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get service requests by status and type",
			args: args{
				ctx:           context.Background(),
				requestStatus: &requeststatus,
				requestType:   &requesttype,
				facilityID:    facilityID,
				flavour:       feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get unresolved service requests",
			args: args{
				ctx:           context.Background(),
				requestStatus: &requeststatus,
				requestType:   &requesttype,
				facilityID:    facilityID,
				flavour:       feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: "",
				flavour:    feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile by staff ID",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				flavour:    feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by staff ID",
			args: args{
				ctx:        context.Background(),
				facilityID: "facilityID",
				flavour:    feedlib.FlavourPro,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad Case - Fail to get service requests" {
				fakeGorm.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service requests by type")
				}
			}

			if tt.name == "Happy Case - Successfully get unresolved service requests" {
				UUID := uuid.New().String()
				nowTime := time.Now()

				serviceRequests := []*gorm.ClientServiceRequest{
					{
						ID:             &UUID,
						ClientID:       uuid.New().String(),
						Active:         true,
						RequestType:    enums.ServiceRequestTypeRedFlag.String(),
						Status:         enums.ServiceRequestStatusPending.String(),
						InProgressAt:   &nowTime,
						InProgressByID: &UUID,
						ResolvedAt:     nil,
						ResolvedByID:   nil,
					},
				}
				fakeGorm.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					return serviceRequests, nil
				}
			}

			if tt.name == "Sad Case - Fail to get staff service requests" {
				fakeGorm.MockGetStaffServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred while getting staff service requests")
				}
			}
			if tt.name == "Sad Case - Fail to get staff profile by staff ID" {
				fakeGorm.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			if tt.name == "Sad Case - Fail to get user profile by staff ID" {
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckUserRole(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		role   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: user has a role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   string(enums.UserRoleTypeCommunityManagement),
			},
			wantErr: false,
			want:    true,
		},

		{
			name: "sad case: user has no role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
				role:   string(enums.UserRoleTypeCommunityManagement),
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: user has no role" {
				fakeGorm.MockCheckUserRoleFn = func(ctx context.Context, userID string, role string) (bool, error) {
					return false, fmt.Errorf("failed to get user role")
				}
			}
			got, err := d.CheckUserRole(tt.args.ctx, tt.args.userID, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckUserRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckUserPermission(t *testing.T) {
	type args struct {
		ctx        context.Context
		userID     string
		permission string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: user has a permission",
			args: args{
				ctx:        context.Background(),
				userID:     uuid.New().String(),
				permission: string(enums.PermissionTypeCanEditOwnRole),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case: user has no permission",
			args: args{
				ctx:        context.Background(),
				userID:     uuid.New().String(),
				permission: string(enums.PermissionTypeCanEditOwnRole),
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: user has no permission" {
				fakeGorm.MockCheckUserPermissionFn = func(ctx context.Context, userID string, permission string) (bool, error) {
					return false, fmt.Errorf("failed to get user permission")
				}
			}

			got, err := d.CheckUserPermission(tt.args.ctx, tt.args.userID, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckUserPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckUserPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetUserRoles(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*domain.AuthorityRole
	}{
		{
			name: "happy case: user has a role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: user has no role",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false, //should not error if user has no roles
		},
		{
			name: "sad case: failed to get user roles",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get user roles" {
				fakeGorm.MockGetUserRolesFn = func(ctx context.Context, userID string) ([]*gorm.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get user roles")
				}
			}

			got, err := d.GetUserRoles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("MyCareHubDb.GetUserRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetUserPermissions(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name string

		args    args
		want    []*domain.AuthorityPermission
		wantErr bool
	}{
		{
			name: "happy case: user has a permission",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: user has no permission",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false, //should not error if user has no permissions
		},
		{
			name: "sad case: failed to get user permissions",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get user permissions" {
				fakeGorm.MockGetUserPermissionsFn = func(ctx context.Context, userID string) ([]*gorm.AuthorityPermission, error) {
					return nil, fmt.Errorf("failed to get user permissions")
				}
			}

			got, err := d.GetUserPermissions(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("MyCareHubDb.GetUserRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfUserNameExists(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully check if username exists",
			args: args{
				ctx:      context.Background(),
				username: gofakeit.Username(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Username already exists",
			args: args{
				ctx:      context.Background(),
				username: gofakeit.Username(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing username",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Username already exists" {
				fakeGorm.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}

			got, err := d.CheckIfUsernameExists(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfUserNameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfUserNameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetCommunityByID(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		communityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Community
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				communityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				communityID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.GetCommunityByID(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCommunityByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected community to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientsInAFacility(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get clients in a facility",
			args: args{
				ctx:        ctx,
				facilityID: "1223445",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get clients in a facility",
			args: args{
				ctx:        ctx,
				facilityID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get clients in a facility" {
				fakeGorm.MockGetClientsInAFacilityFn = func(ctx context.Context, facilityID string) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get clients in a facility")
				}
			}

			got, err := d.GetClientsInAFacility(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientsInAFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetRecentHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx          context.Context
		lastSyncTime time.Time
		clientID     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get recent health diary entries",
			args: args{
				ctx:          ctx,
				lastSyncTime: time.Time{},
				clientID:     "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get recent entries",
			args: args{
				ctx:          ctx,
				lastSyncTime: time.Time{},
				clientID:     "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get recent entries" {
				fakeGorm.MockGetRecentHealthDiaryEntriesFn = func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get recent health diary entries")
				}
			}

			got, err := d.GetRecentHealthDiaryEntries(tt.args.ctx, tt.args.lastSyncTime, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetRecentHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientsByParams(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx          context.Context
		params       gorm.Client
		lastSyncTime *time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientProfile
		wantErr bool
	}{
		{
			name: "happy case: retrieving clients",
			args: args{
				ctx:    context.Background(),
				params: gorm.Client{},
			},
			wantErr: false,
		},
		{
			name: "sad case: error retrieving clients",
			args: args{
				ctx:    context.Background(),
				params: gorm.Client{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: error retrieving clients" {
				fakeGorm.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("cannot retrieve clients")
				}
			}

			got, err := d.GetClientsByParams(tt.args.ctx, tt.args.params, tt.args.lastSyncTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientsByParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected community to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientCCCIdentifier(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Identifier
		wantErr bool
	}{
		{
			name: "happy case: retrieve client ccc identifier",
			args: args{
				ctx:      context.Background(),
				clientID: "",
			},
			wantErr: false,
		},
		{
			name: "sad case: error retrieving client ccc identifier",
			args: args{
				ctx:      context.Background(),
				clientID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "sad case: error retrieving client ccc identifier" {
			fakeGorm.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*gorm.Identifier, error) {
				return nil, fmt.Errorf("cannot get client identifier")
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := d.GetClientCCCIdentifier(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientCCCIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected community to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetServiceRequestsForKenyaEMR(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	currentTime := time.Now()

	type args struct {
		ctx     context.Context
		payload *dto.ServiceRequestPayload
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      1234,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      123,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockGetServiceRequestsForKenyaEMRFn = func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetServiceRequestsForKenyaEMR(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetServiceRequestsForKenyaEMR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequestsForKenyaEMR = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequestsForKenyaEMR = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetScreeningToolsQuestions(t *testing.T) {
	type args struct {
		ctx      context.Context
		toolType string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolQuestion
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      context.Background(),
				toolType: string(enums.ScreeningToolTypeTB),
			},
			wantErr: false,
		},
		{
			name: "failed to convert response choices to map",
			args: args{
				ctx:      context.Background(),
				toolType: string(enums.ScreeningToolTypeTB),
			},
			wantErr: true,
		},
		{
			name: "failed to convert metadata to map",
			args: args{
				ctx:      context.Background(),
				toolType: string(enums.ScreeningToolTypeTB),
			},
			wantErr: true,
		},
		{
			name: "failed to get screeningtool questions",
			args: args{
				ctx:      context.Background(),
				toolType: string(enums.ScreeningToolTypeTB),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "failed to convert response choices to map" {
				fakeGorm.MockGetScreeningToolsQuestionsFn = func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error) {
					return []gorm.ScreeningToolQuestion{
						{
							ID:               uuid.New().String(),
							Question:         gofakeit.Sentence(1),
							ToolType:         enums.ScreeningToolTypeTB.String(),
							ResponseChoices:  `{"1": "Yes", "2": "No"`,
							ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
							ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
							Sequence:         1,
							Active:           true,
							OrganisationID:   uuid.New().String(),
						},
					}, nil
				}
			}

			if tt.name == "failed to convert metadata to map" {
				fakeGorm.MockGetScreeningToolsQuestionsFn = func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error) {
					return []gorm.ScreeningToolQuestion{
						{
							ID:               uuid.New().String(),
							Question:         gofakeit.Sentence(1),
							ToolType:         enums.ScreeningToolTypeTB.String(),
							ResponseChoices:  `{"1": "Yes", "2": "No"}`,
							ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
							ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
							Sequence:         1,
							Active:           true,
							Meta:             `{"meta": "data"`,
							OrganisationID:   uuid.New().String(),
						},
					}, nil
				}
			}

			if tt.name == "failed to get screeningtool questions" {
				fakeGorm.MockGetScreeningToolsQuestionsFn = func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetScreeningToolQuestions(tt.args.ctx, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckIdentifierExists(t *testing.T) {
	type args struct {
		ctx             context.Context
		identifierType  string
		identifierValue string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:             context.Background(),
				identifierType:  "ccc",
				identifierValue: "12345",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.CheckIdentifierExists(tt.args.ctx, tt.args.identifierType, tt.args.identifierValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIdentifierExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIdentifierExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckFacilityExistsByMFLCode(t *testing.T) {
	type args struct {
		ctx     context.Context
		MFLCode int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     context.Background(),
				MFLCode: 12345,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.CheckFacilityExistsByMFLCode(tt.args.ctx, tt.args.MFLCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckFacilityExistsByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckFacilityExistsByMFLCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_ListAppointments(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		params     *domain.Appointment
		filter     []*firebasetools.FilterParam
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Appointment
		want1   *domain.Pagination
		wantErr bool
	}{
		{
			name: "happy case: success listing appointments",
			args: args{
				ctx: context.Background(),
				params: &domain.Appointment{
					ClientID: gofakeit.UUID(),
				},
				filter: nil,
				pagination: &domain.Pagination{
					Limit:       2,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error listing appointments",
			args: args{
				ctx: context.Background(),
				params: &domain.Appointment{
					ClientID: gofakeit.UUID(),
				},
				filter: nil,
				pagination: &domain.Pagination{
					Limit:       2,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error listing appointments" {
				fakeGorm.MockListAppointments = func(ctx context.Context, params *gorm.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("error listing appointments")
				}
			}

			got, got1, err := d.ListAppointments(tt.args.ctx, tt.args.params, tt.args.filter, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListAppointments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected appointments to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected appointments not to be nil for %v", tt.name)
				return
			}

			if tt.wantErr && got1 != nil {
				t.Errorf("expected page info to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got1 == nil {
				t.Errorf("expected page info not to be nil for %v", tt.name)
				return
			}

		})
	}
}

func TestMyCareHubDb_GetScreeningToolQuestionByQuestionID(t *testing.T) {
	type args struct {
		ctx        context.Context
		questionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        context.Background(),
				questionID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "failed to convert response choices to map",
			args: args{
				ctx:        context.Background(),
				questionID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "failed to convert metadata to map",
			args: args{
				ctx:        context.Background(),
				questionID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "failed to get screeningtool questions",
			args: args{
				ctx:        context.Background(),
				questionID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "failed to convert response choices to map" {
				fakeGorm.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error) {
					return &gorm.ScreeningToolQuestion{
						ID:               uuid.New().String(),
						Question:         gofakeit.Sentence(1),
						ToolType:         enums.ScreeningToolTypeTB.String(),
						ResponseChoices:  `{"1": "Yes", "2": "No"`,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
						ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
						Sequence:         1,
						Active:           true,
						OrganisationID:   uuid.New().String(),
					}, nil
				}
			}

			if tt.name == "failed to convert metadata to map" {
				fakeGorm.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error) {
					return &gorm.ScreeningToolQuestion{
						ID:               uuid.New().String(),
						Question:         gofakeit.Sentence(1),
						ToolType:         enums.ScreeningToolTypeTB.String(),
						ResponseChoices:  `{"1": "Yes", "2": "No"}`,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
						ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
						Sequence:         1,
						Active:           true,
						Meta:             `{"meta": "data"`,
						OrganisationID:   uuid.New().String(),
					}, nil
				}
			}

			if tt.name == "failed to get screeningtool questions" {
				fakeGorm.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetScreeningToolQuestionByQuestionID(tt.args.ctx, tt.args.questionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolQuestionByQuestionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetClientProfileByCCCNumber(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		CCCNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profile by CCC number",
			args: args{
				ctx:       ctx,
				CCCNumber: "345678",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile by CCC number",
			args: args{
				ctx:       ctx,
				CCCNumber: "111111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client profile by CCC number" {
				fakeGorm.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by CCC number")
				}
			}

			got, err := d.GetClientProfileByCCCNumber(tt.args.ctx, tt.args.CCCNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientProfileByCCCNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckIfClientHasUnresolvedServiceRequests(t *testing.T) {
	type args struct {
		ctx                context.Context
		clientID           string
		serviceRequestType string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                context.Background(),
				clientID:           uuid.New().String(),
				serviceRequestType: string(enums.ServiceRequestTypePinReset),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.CheckIfClientHasUnresolvedServiceRequests(tt.args.ctx, tt.args.clientID, tt.args.serviceRequestType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfClientHasUnresolvedServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfClientHasUnresolvedServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetAllRoles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get roles",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to get roles" {
				fakeGorm.MockGetAllRolesFn = func(ctx context.Context) ([]*gorm.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get roles")
				}
			}

			got, err := d.GetAllRoles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchClientProfilesByCCCNumber(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx       context.Context
		CCCNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profiles by CCC number",
			args: args{
				ctx:       ctx,
				CCCNumber: "345678",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:       ctx,
				CCCNumber: "111111",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client identifier",
			args: args{
				ctx:       ctx,
				CCCNumber: "111111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - Fail to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - Fail to get client identifier" {
				fakeGorm.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*gorm.Identifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.SearchClientProfilesByCCCNumber(tt.args.ctx, tt.args.CCCNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchClientProfilesByCCCNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected client profiles to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected client profiles not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetHealthDiaryEntryByID(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx                context.Context
		healthDiaryEntryID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientHealthDiaryEntry
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockGetHealthDiaryEntryByIDFn = func(ctx context.Context, healthDiaryEntryID string) (*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetHealthDiaryEntryByID(tt.args.ctx, tt.args.healthDiaryEntryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetHealthDiaryEntryByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected healthdiary to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected health diary not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetServiceRequestByID(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              context.Background(),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case - Failed to get service request by ID",
			args: args{
				ctx:              context.Background(),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - Failed to get service request by ID" {
				fakeGorm.MockGetServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service request by ID")
				}
			}

			got, err := d.GetServiceRequestByID(tt.args.ctx, tt.args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetServiceRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetStaffProfileByStaffID(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx     context.Context
		staffID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetStaffProfileByStaffID(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffProfileByStaffID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetUserProfileByStaffID(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx     context.Context
		staffID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetUserProfileByStaffID(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserProfileByStaffID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetAppointmentServiceRequests(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx          context.Context
		lastSyncTime time.Time
		mflCode      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: false,
		},
		{
			name: "Sad case:  failed to get appointment service request",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  invalid service request meta",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get appointment by id",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to convert to suggested time",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to convert to suggested date",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get in progress by staff",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get resolved by staff",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get client profile by client ID",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get user profile by user ID",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to get client ccc identifier by client ID",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  failed to facility by facility by id",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case:  failed to get appointment service request" {
				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case:  invalid service request meta" {
				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					requestID := uuid.New().String()
					currentTime := time.Now()
					staffID := uuid.New().String()
					facilityID := uuid.New().String()
					return []*gorm.ClientServiceRequest{
						{
							ID:             &requestID,
							Active:         true,
							RequestType:    enums.ServiceRequestTypeAppointments.String(),
							Request:        "REQUEST",
							Status:         "PENDING",
							InProgressAt:   &currentTime,
							ResolvedAt:     &currentTime,
							ClientID:       uuid.New().String(),
							InProgressByID: &staffID,
							OrganisationID: "",
							ResolvedByID:   &staffID,
							FacilityID:     facilityID,
							Meta:           ``,
						},
					}, nil
				}
			}

			if tt.name == "Sad case:  failed to get appointment by id" {
				fakeGorm.MockGetAppointmentByIDFn = func(ctx context.Context, appointmentID string) (*gorm.Appointment, error) {
					return nil, fmt.Errorf("an error occurred")

				}
			}

			if tt.name == "Sad case:  failed to convert to suggested time" {
				fakeGorm.MockGetAppointmentByIDFn = func(ctx context.Context, appointmentID string) (*gorm.Appointment, error) {
					date := time.Now().Add(time.Duration(100))
					return &gorm.Appointment{
						ID:             gofakeit.UUID(),
						OrganisationID: gofakeit.UUID(),
						Active:         true,
						ExternalID:     gofakeit.UUID(),
						ClientID:       gofakeit.UUID(),
						FacilityID:     gofakeit.UUID(),
						Reason:         "Knocked up",
						Date:           date,
					}, nil
				}
			}

			if tt.name == "Sad case:  failed to convert to suggested date" {
				fakeGorm.MockGetAppointmentByIDFn = func(ctx context.Context, appointmentID string) (*gorm.Appointment, error) {
					return &gorm.Appointment{
						ID:             gofakeit.UUID(),
						OrganisationID: gofakeit.UUID(),
						Active:         true,
						ExternalID:     gofakeit.UUID(),
						ClientID:       gofakeit.UUID(),
						FacilityID:     gofakeit.UUID(),
						Reason:         "Knocked up",
						Date:           time.Now(),
					}, nil
				}
			}

			if tt.name == "Sad case:  failed to get in progress by staff" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"id":                uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					return []*gorm.ClientServiceRequest{
						{
							ID:             &UUID,
							Active:         true,
							RequestType:    enums.ServiceRequestTypeAppointments.String(),
							Request:        gofakeit.Sentence(1),
							Status:         enums.ServiceRequestStatusPending.String(),
							InProgressAt:   &now,
							ResolvedAt:     nil,
							ClientID:       gofakeit.Name(),
							InProgressByID: &UUID,
							ResolvedByID:   nil,
							FacilityID:     gofakeit.Name(),
							Meta:           string(bs),
						},
					}, nil
				}
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Sad case:  failed to get resolved by staff" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"id":                uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					return []*gorm.ClientServiceRequest{
						{
							ID:             &UUID,
							Active:         true,
							RequestType:    enums.ServiceRequestTypeAppointments.String(),
							Request:        gofakeit.Sentence(1),
							Status:         enums.ServiceRequestStatusPending.String(),
							InProgressAt:   &now,
							ResolvedAt:     &now,
							ClientID:       gofakeit.Name(),
							InProgressByID: &UUID,
							ResolvedByID:   &UUID,
							FacilityID:     gofakeit.Name(),
							Meta:           string(bs),
						},
					}, nil
				}
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Sad case:  failed to get client profile by client ID" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad case:  failed to get client ccc identifier by client ID" {
				fakeGorm.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client ccc identifier")
				}
			}

			if tt.name == "Sad case:  failed to facility by facility by id" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facility by id")
				}
			}
			got, err := d.GetAppointmentServiceRequests(tt.args.ctx, tt.args.lastSyncTime, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAppointmentServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetAssessmentResponses(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		facilityID string
		toolType   string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolAssesmentResponse
		wantErr bool
	}{
		{
			name: "Happy case:  get assessment responses",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				toolType:   "GBV",
			},
			wantErr: false,
		},
		{
			name: "Sad case:  unable to get assessment responses",
			args: args{
				ctx:        ctx,
				facilityID: "facilityID",
				toolType:   "GBV",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  unable to get client profile",
			args: args{
				ctx:        ctx,
				facilityID: "facilityID",
				toolType:   "GBV",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  unable to get user profile",
			args: args{
				ctx:        ctx,
				facilityID: "facilityID",
				toolType:   "GBV",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  unable to get assessment responses" {
				fakeGorm.MockGetAnsweredScreeningToolQuestionsFn = func(ctx context.Context, facilityID, toolType string) ([]*gorm.ScreeningToolsResponse, error) {
					return nil, fmt.Errorf("failed to get assessment responses")
				}
			}
			if tt.name == "Sad case:  unable to get client profile" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad case:  unable to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}
			got, err := d.GetAssessmentResponses(tt.args.ctx, tt.args.facilityID, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAssessmentResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected response to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetFacilitiesWithoutFHIRID(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case:  failed to get facilities without fhir id",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  failed to get facilities without fhir id" {
				fakeGorm.MockGetFacilitiesWithoutFHIRIDFn = func(ctx context.Context) ([]*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities without fhir id")

				}
			}
			got, err := d.GetFacilitiesWithoutFHIRID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetFacilitiesWithoutFHIRID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected client to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected client not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetAppointmentByClientID(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx           context.Context
		appointmentID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case:  get client appointment by id",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case:  failed to get client appointment by id",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  failed to get client appointment by id" {
				fakeGorm.MockGetAppointmentByIDFn = func(ctx context.Context, appointmentID string) (*gorm.Appointment, error) {
					return nil, fmt.Errorf("failed to get client appointment by id")
				}
			}
			got, err := d.GetAppointmentByClientID(tt.args.ctx, tt.args.appointmentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAppointmentByClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetAppointmentByExternalID(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		externalID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case:  get appointment by appointment UUID",
			args: args{
				ctx:        context.Background(),
				externalID: uuid.New().String(),
			},
		},
		{
			name: "Sad case:  failed to get appointment by appointment UUID",
			args: args{
				ctx:        context.Background(),
				externalID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  failed to get appointment by appointment UUID" {
				fakeGorm.MockGetAppointmentByExternalIDFn = func(ctx context.Context, externalID string) (*gorm.Appointment, error) {
					return nil, fmt.Errorf("failed to get appointment by appointment UUID")
				}
			}

			got, err := d.GetAppointmentByExternalID(tt.args.ctx, tt.args.externalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAppointmentByExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetClientServiceRequests(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx      context.Context
		toolType string
		status   string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case:  get system generated client service requests",
			args: args{
				ctx:      context.Background(),
				toolType: "SCREENING_TOOLS",
				status:   enums.ServiceRequestStatusPending.String(),
				clientID: uuid.New().String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.GetClientServiceRequests(tt.args.ctx, tt.args.toolType, tt.args.status, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetActiveScreeningToolResponses(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolQuestionResponse
		wantErr bool
	}{
		{
			name: "Happy case:  get active screening tool responses",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.GetActiveScreeningToolResponses(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetActiveScreeningToolResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}
