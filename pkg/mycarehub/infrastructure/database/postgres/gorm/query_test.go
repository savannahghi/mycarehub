package gorm_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
)

func createTestFacility() *gorm.Facility {
	ID := uuid.New().String()
	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := gofakeit.Name()
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      true,
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
		MFLCode  int
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
				MFLCode:  102921083,
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

	code := rand.Intn(1000000)
	code2 := rand.Intn(1000000)

	facilityInput := &gorm.Facility{
		Name:        ksuid.New().String(),
		Code:        code,
		Active:      true,
		County:      "Nairobi",
		Description: "This is just for mocking",
	}

	facilityInput2 := &gorm.Facility{
		Name:        ksuid.New().String(),
		Code:        code2,
		Active:      true,
		County:      "Baringo",
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
			Value:    strconv.Itoa(code),
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
			Value:    strconv.Itoa(code),
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
			Value:    strconv.Itoa(code),
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
			Value:    strconv.Itoa(code),
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
		Field:     enums.FilterSortDataTypeCreatedAt,
		Direction: enums.SortDataTypeAsc,
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
	// Teardown
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	if err = pg.DB.Where("mfl_code", code).Unscoped().Delete(&gorm.Facility{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}

	if err = pg.DB.Where("mfl_code", code2).Unscoped().Delete(&gorm.Facility{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
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

func TestPGInstance_GetSecurityQuestions(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	sequence := 2
	securityQuestionID := ksuid.New().String()

	securityQuestionInput := &gorm.SecurityQuestion{
		SecurityQuestionID: &securityQuestionID,
		QuestionStem:       "test",
		Description:        "desc description",
		ResponseType:       enums.SecurityQuestionResponseTypeDate,
		Flavour:            feedlib.FlavourConsumer,
		Active:             true,
		Sequence:           &sequence,
		OrganisationID:     serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}
	err = pg.DB.Create(securityQuestionInput).Error
	if err != nil {
		t.Errorf("failed to security questions: %v", err)
	}

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
				flavour: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetSecurityQuestions(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
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

	//TearDown
	if err = pg.DB.Where("id", securityQuestionInput.SecurityQuestionID).Unscoped().Delete(&gorm.SecurityQuestion{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_GetSecurityQuestionByID(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	sequence := 2
	securityQuestionInput := &gorm.SecurityQuestion{
		QuestionStem:   "test",
		Description:    "desc description",
		ResponseType:   enums.SecurityQuestionResponseTypeDate,
		Flavour:        feedlib.FlavourConsumer,
		Active:         true,
		Sequence:       &sequence,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}

	err = pg.DB.Create(securityQuestionInput).Error
	if err != nil {
		t.Errorf("failed to security questions: %v", err)
	}

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
			name: "Happy case",
			args: args{
				ctx:                ctx,
				securityQuestionID: securityQuestionInput.SecurityQuestionID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetSecurityQuestionByID(tt.args.ctx, tt.args.securityQuestionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSecurityQuestionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}

	// TearDown
	if err = pg.DB.Where("id", securityQuestionInput.SecurityQuestionID).Unscoped().Delete(&gorm.SecurityQuestion{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}

}

func TestPGInstance_CheckIfPhoneNumberExists(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	userID := uuid.New().String()
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		UserID:              &userID,
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         false,
		OrganisationID:      serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contactID := uuid.New().String()
	contact := &gorm.Contact{
		ContactID:      &contactID,
		ContactType:    "SMS",
		ContactValue:   "+254710000001",
		Active:         true,
		OptedIn:        true,
		UserID:         userInput.UserID,
		Flavour:        userInput.Flavour,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}

	err = pg.DB.Create(&contact).Error
	if err != nil {
		t.Errorf("failed to create contact: %v", err)
	}

	type args struct {
		ctx       context.Context
		phone     string
		isOptedIn bool
		flavour   feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				phone:     contact.ContactValue,
				isOptedIn: true,
				flavour:   contact.Flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:       ctx,
				phone:     contact.ContactValue,
				isOptedIn: true,
				flavour:   "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - ungistered phone",
			args: args{
				ctx:       ctx,
				phone:     "+254711223344",
				isOptedIn: true,
				flavour:   contact.Flavour,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.CheckIfPhoneNumberExists(tt.args.ctx, tt.args.phone, tt.args.isOptedIn, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfPhoneNumberExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", contact.ContactID).Unscoped().Delete(&gorm.Contact{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_VerifyOTP(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	userID := uuid.New().String()
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		UserID:              &userID,
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	otpID := gofakeit.Number(1, 10000)
	generatedAt := time.Now()
	validUntil := time.Now().AddDate(0, 0, 2)

	ext := extension.NewExternalMethodsImpl()

	otp, err := ext.GenerateOTP(ctx)
	if err != nil {
		t.Errorf("unable to generate OTP")
	}

	otpInput := &gorm.UserOTP{
		OTPID:       otpID,
		UserID:      *userInput.UserID,
		Valid:       true,
		GeneratedAt: generatedAt,
		ValidUntil:  validUntil,
		Channel:     "SMS",
		Flavour:     userInput.Flavour,
		PhoneNumber: "+254710000111",
		OTP:         otp,
	}

	err = pg.DB.Create(&otpInput).Error
	if err != nil {
		t.Errorf("failed to create otp: %v", err)
	}

	validOTPPayload := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         otpInput.OTP,
		Flavour:     flavour,
	}
	invalidOTPPayload2 := &dto.VerifyOTPInput{
		PhoneNumber: "",
		OTP:         otpInput.OTP,
		Flavour:     flavour,
	}
	invalidOTPPayload3 := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         "",
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         otpInput.OTP,
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
			name: "Sad case - very bad inputs",
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
			got, err := testingDB.VerifyOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.VerifyOTP() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", otpInput.OTPID).Unscoped().Delete(&gorm.UserOTP{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_GetClientProfileByUserID(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	testFacility := createTestFacility()
	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             feedlib.FlavourConsumer,
		Avatar:              "",
		IsSuspended:         false,
		OrganisationID:      serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
		return
	}

	invalidID := uuid.New().String()
	time := time.Now()
	client := &gorm.Client{
		UserID:                  userInput.UserID,
		Active:                  true,
		ClientType:              "OVC",
		TreatmentEnrollmentDate: &time,
		FHIRPatientID:           uuid.New().String(),
		FacilityID:              *facility.FacilityID,
		OrganisationID:          serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
		CHVUserID:               *userInput.UserID,
	}

	err = pg.DB.Create(client).Error
	if err != nil {
		t.Errorf("failed to create client: %v", err)
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
			name: "Happy Case - Successfully get client profile",
			args: args{
				ctx:    ctx,
				userID: *userInput.UserID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:    ctx,
				userID: invalidID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}
