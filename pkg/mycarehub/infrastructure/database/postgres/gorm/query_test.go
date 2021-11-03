package gorm_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func createTestFacility() *gorm.Facility {
	ID := uuid.New().String()
	name := gofakeit.Name()
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

func TestPGInstance_GetUserPINByUserID(t *testing.T) {
	ctx := context.Background()
	User := &gorm.User{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Username:    gofakeit.Username(),
		MiddleName:  gofakeit.Name(),
		DisplayName: gofakeit.BeerAlcohol(),
		Gender:      enumutils.GenderMale,
	}
	Client := &gorm.ClientProfile{
		ClientType: enums.ClientTypeOvc,
	}

	user, err := testingDB.RegisterClient(ctx, User, Client)
	if err != nil {
		t.Errorf("failed to create test user")
		return
	}

	id := uuid.New().String()
	pinPayload := &gorm.PINData{
		PINDataID: &id,
		UserID:    *user.User.UserID,
		HashedPIN: "1234",
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
		Salt:      "1234",
	}

	_, err = testingDB.SavePin(ctx, pinPayload)
	if err != nil {
		t.Errorf("failed to save test pin")
		return
	}

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
				userID: *user.User.UserID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get pin",
			args: args{
				ctx:    ctx,
				userID: "1234amklca",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserPINByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserPINByUserID() error = %v, wantErr %v", err, tt.wantErr)
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
		Name:        "Kanairo One",
		Code:        code,
		Active:      strconv.FormatBool(true),
		County:      enums.CountyTypeNairobi,
		Description: "This is just for mocking",
	}

	facilityInput2 := &gorm.Facility{
		Name:        "Baringo 2",
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
			Name:     enums.FilterDataTypeName.String(),
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterDataTypeMFLCode.String(),
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterDataTypeActive.String(),
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterDataTypeCounty.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterEmptyName := []*domain.FiltersParam{
		{
			Name:     enums.FilterDataTypeName.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    "",
		},
		{
			Name:     enums.FilterDataTypeMFLCode.String(),
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterDataTypeActive.String(),
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterDataTypeCounty.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}
	filterEmptyMFLCode := []*domain.FiltersParam{
		{
			Name:     enums.FilterDataTypeName.String(),
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterDataTypeMFLCode.String(),
			DataType: enums.FilterDataTypeMFLCode,
			Value:    "",
		},
		{
			Name:     enums.FilterDataTypeActive.String(),
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterDataTypeCounty.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidBool := []*domain.FiltersParam{
		{
			Name:     enums.FilterDataTypeName.String(),
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterDataTypeMFLCode.String(),
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterDataTypeActive.String(),
			DataType: enums.FilterDataTypeActive,
			Value:    "invalid",
		},
		{
			Name:     enums.FilterDataTypeCounty.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    enums.CountyTypeNairobi.String(),
		},
	}

	filterInvalidCounty := []*domain.FiltersParam{
		{
			Name:     enums.FilterDataTypeName.String(),
			DataType: enums.FilterDataTypeName,
			Value:    "Kanairo One",
		},
		{
			Name:     enums.FilterDataTypeMFLCode.String(),
			DataType: enums.FilterDataTypeMFLCode,
			Value:    code,
		},
		{
			Name:     enums.FilterDataTypeActive.String(),
			DataType: enums.FilterDataTypeActive,
			Value:    formatBool,
		},
		{
			Name:     enums.FilterDataTypeCounty.String(),
			DataType: enums.FilterDataTypeCounty,
			Value:    "Kanairo",
		},
	}

	paginationInput := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       1,
			CurrentPage: 1,
		},
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

func TestPGInstance_GetFacilities(t *testing.T) {
	ctx := context.Background()
	testFacility := createTestFacility()

	_, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Happy Case - List all facilities",
			args:    args{ctx: ctx},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
