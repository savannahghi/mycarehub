package gorm_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestPGInstance_RetrieveFacility(t *testing.T) {
	ctx := context.Background()
	fakeID := "1234"

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
				id:       &facilityID,
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
		{
			name: "Sad Case - nil facility",
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
				MFLCode:  mflCode,
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

	noSearchTerm := ""
	searchTerm := "ro"

	noFilterInput := []*domain.FiltersParam{}

	formatBool := strconv.FormatBool(true)

	filterInput := []*domain.FiltersParam{
		{
			Name:     enums.FilterSortDataTypeName.String(),
			DataType: enums.FilterSortDataTypeName,
			Value:    "Nairobi",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(mflCode),
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
			Value:    strconv.Itoa(mflCode),
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
			Value:    "Nairobi",
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
			Value:    "Nairobi",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(mflCode),
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
			Value:    "Nairobi",
		},
		{
			Name:     enums.FilterSortDataTypeMFLCode.String(),
			DataType: enums.FilterSortDataTypeMFLCode,
			Value:    strconv.Itoa(mflCode),
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

func TestPGInstance_GetFacilities(t *testing.T) {
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
}

func TestPGInstance_GetSecurityQuestionByID(t *testing.T) {
	ctx := context.Background()

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
				securityQuestionID: &securityQuestionID,
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
}

func TestPGInstance_CheckIfPhoneNumberExists(t *testing.T) {
	ctx := context.Background()

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
				phone:     testPhone,
				isOptedIn: true,
				flavour:   testFlavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:       ctx,
				phone:     testPhone,
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
				flavour:   testFlavour,
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
}

func TestPGInstance_VerifyOTP(t *testing.T) {
	ctx := context.Background()

	flavour := feedlib.FlavourConsumer

	validOTPPayload := &dto.VerifyOTPInput{
		PhoneNumber: testPhone,
		OTP:         testOTP,
		Flavour:     flavour,
	}
	invalidOTPPayload2 := &dto.VerifyOTPInput{
		PhoneNumber: "",
		OTP:         testOTP,
		Flavour:     flavour,
	}
	invalidOTPPayload3 := &dto.VerifyOTPInput{
		PhoneNumber: testPhone,
		OTP:         "",
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		PhoneNumber: testPhone,
		OTP:         testOTP,
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
}

func TestPGInstance_GetClientProfileByUserID(t *testing.T) {
	ctx := context.Background()

	invalidID := uuid.New().String()

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
				userID: userID,
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

func TestPGInstance_GetOTP(t *testing.T) {
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
			name: "Happy Case - Successfully get OTP",
			args: args{
				ctx:         ctx,
				phoneNumber: testPhone,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get OTP, nonexisting phone",
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

			got, err := testingDB.GetOTP(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}

}

func TestPGInstance_GetUserSecurityQuestionsResponses(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	securityQuestionResponseInput := &gorm.SecurityQuestionResponse{
		UserID:         userID,
		QuestionID:     securityQuestionID,
		Response:       "1917",
		Timestamp:      time.Now(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(securityQuestionResponseInput).Error
	if err != nil {
		t.Errorf("Create securityQuestionResponse failed: %v", err)
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
			name: "Happy Case - Successfully get security questions",

			args: args{
				ctx:    ctx,
				userID: userID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetUserSecurityQuestionsResponses(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserSecurityQuestionsResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetContactByUserID(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()

	type args struct {
		ctx         context.Context
		userID      *string
		contactType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get PHONE contact",
			args: args{
				ctx:         ctx,
				userID:      &userID,
				contactType: "PHONE",
			},
		},
		{
			name: "invalid: contact type invalid",
			args: args{
				ctx:         ctx,
				userID:      &userID,
				contactType: "Invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid: user ID invalid",
			args: args{
				ctx:         ctx,
				userID:      &ID,
				contactType: "EMAIL",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetContactByUserID(tt.args.ctx, tt.args.userID, tt.args.contactType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetContactByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_CheckWhetherUserHasLikedContent(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	contentLike := &gorm.ContentLike{
		Active:         true,
		ContentID:      contentID,
		UserID:         userID,
		OrganisationID: orgID,
	}
	err = pg.DB.Create(contentLike).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

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
				userID:    contentLike.UserID,
				contentID: contentLike.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - bad user ID",
			args: args{
				ctx:       ctx,
				userID:    gofakeit.HipsterSentence(100),
				contentID: contentLike.ContentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid contentID",
			args: args{
				ctx:       ctx,
				userID:    contentLike.UserID,
				contentID: -5,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CheckWhetherUserHasLikedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckWhetherUserHasLikedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserProfileByUserID(t *testing.T) {
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
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty user ID",
			args: args{
				ctx:    ctx,
				userID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetClientHealthDiaryQuote(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}
	quoteInput := &gorm.ClientHealthDiaryQuote{
		Author:         gofakeit.FirstName() + " " + gofakeit.LastName(),
		Quote:          gofakeit.Sentence(10),
		Active:         true,
		OrganisationID: orgID,
	}

	err = pg.DB.Create(&quoteInput).Error
	if err != nil {
		t.Errorf("failed to create quote: %v", err)
	}

	quoteInput2 := &gorm.ClientHealthDiaryQuote{
		Author:         gofakeit.FirstName() + " " + gofakeit.LastName(),
		Quote:          gofakeit.Sentence(10),
		Active:         true,
		OrganisationID: orgID,
	}

	err = pg.DB.Create(&quoteInput2).Error
	if err != nil {
		t.Errorf("failed to create quote: %v", err)
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
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetClientHealthDiaryQuote(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientHealthDiaryQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
	// tear down
	if err = pg.DB.Where("id", quoteInput.ClientHealthDiaryQuoteID).Unscoped().Delete(&gorm.ClientHealthDiaryQuote{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", quoteInput2.ClientHealthDiaryQuoteID).Unscoped().Delete(&gorm.ClientHealthDiaryQuote{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_CanRecordHeathDiary(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	canShowHealthDiaryInput := gorm.ClientHealthDiaryEntry{
		Base: gorm.Base{
			CreatedAt: time.Now().Add(time.Hour * -25),
			UpdatedAt: time.Now().Add(time.Hour * -25),
		},
		Active:                false,
		Mood:                  "HAPPY",
		Note:                  "test note",
		EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY",
		ShareWithHealthWorker: false,
		SharedAt:              time.Now().Add(time.Hour * -25),
		ClientID:              clientID2,
		OrganisationID:        orgID,
	}
	err = pg.DB.Create(&canShowHealthDiaryInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	cannotShowHealthDiaryInput := &gorm.ClientHealthDiaryEntry{
		Base: gorm.Base{
			CreatedAt: time.Now().Add(time.Hour * -20),
			UpdatedAt: time.Now().Add(time.Hour * -20),
		},
		Active:                false,
		Mood:                  "HAPPY",
		Note:                  "test note",
		EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY",
		ShareWithHealthWorker: false,
		SharedAt:              time.Now().Add(time.Hour * -20),
		ClientID:              clientID2,
		OrganisationID:        orgID,
	}

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case, can record after 24 hours",
			args: args{
				ctx:      ctx,
				clientID: clientID2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case, cant record new entry before 24 hours",
			args: args{
				ctx:      ctx,
				clientID: clientID2,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid: id not type UUID",
			args: args{
				ctx:      ctx,
				clientID: "invalid",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case, cant record new entry before 24 hours" {
				if err = pg.DB.Where("client_id", clientID2).Unscoped().Delete(&gorm.ClientHealthDiaryEntry{}).Error; err != nil {
					t.Errorf("failed to delete record = %v", err)
				}
				err = pg.DB.Create(&cannotShowHealthDiaryInput).Error
				if err != nil {
					t.Errorf("failed to create user: %v", err)
				}
			}
			got, err := testingDB.CanRecordHeathDiary(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CanRecordHeathDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CanRecordHeathDiary() = %v, want %v", got, tt.want)
			}
		})
	}
	// TearDown
	if err = pg.DB.Where("client_id", cannotShowHealthDiaryInput.ClientID).Unscoped().Delete(&gorm.ClientHealthDiaryEntry{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_CheckIfUserBookmarkedContent(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	bookmarkInput := &gorm.ContentBookmark{
		Active:         true,
		ContentID:      contentID,
		UserID:         userID,
		OrganisationID: orgID,
	}
	err = pg.DB.Create(bookmarkInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

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
			name: "happy case: get user bookmarked content by user id and content id",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: bookmarkInput.ContentID,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIfUserBookmarkedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfUserBookmarkedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ListContentCategories(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get all content categories",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetUserProfileByPhoneNumber(t *testing.T) {
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
			name: "happy case: get user profile by phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: testPhone,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserProfileByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserProfileByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetUserPINByUserID(t *testing.T) {
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
			name: "happy case: get user pin by user id",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantErr: false,
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
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetSecurityQuestionResponseByID(t *testing.T) {
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
			name: "happy case: get security question response by id",
			args: args{
				ctx:        context.Background(),
				questionID: securityQuestionID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetSecurityQuestionResponseByID(tt.args.ctx, tt.args.questionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSecurityQuestionResponseByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_CheckUserHasPin(t *testing.T) {
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
			name: "happy case: check user has pin",
			args: args{
				ctx:     context.Background(),
				userID:  userID,
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckUserHasPin(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckUserHasPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckUserHasPin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserBookmarkedContent(t *testing.T) {
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
			name: "happy case: get user bookmarked content",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserBookmarkedContent(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetClientHealthDiaryEntries(t *testing.T) {
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
			name: "happy case: get client health diary entries",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetFAQContent(t *testing.T) {
	ctx := context.Background()
	limit := 10
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}
	// create FAQ
	faqInput := &gorm.FAQ{
		Active:         true,
		Title:          gofakeit.Name(),
		Description:    gofakeit.Name(),
		Body:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(faqInput).Error
	if err != nil {
		t.Errorf("failed to create faq: %v", err)
	}

	type args struct {
		ctx     context.Context
		limit   *int
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get faq content",
			args: args{
				ctx:     ctx,
				limit:   &limit,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetFAQContent(tt.args.ctx, tt.args.flavour, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetFAQContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
	// TearDown
	if err = pg.DB.Where("id", faqInput.FAQID).Unscoped().Delete(&gorm.FAQ{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
