package gorm_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
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
		{
			name: "Sad Case - unacceptable MFL code",
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
	searchParameter := "Nairobi"
	type args struct {
		ctx             context.Context
		searchParameter *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Happy Case - List all facilities",
			args:    args{ctx: ctx, searchParameter: &searchParameter},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchFacility(tt.args.ctx, tt.args.searchParameter)
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
	invalid := "invalid"

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
			wantErr: false,
		},
		{
			name: "Sad case: invalid security question id",
			args: args{
				ctx:                ctx,
				securityQuestionID: &invalid,
			},
			wantErr: true,
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
	OTPnotFound := &dto.VerifyOTPInput{
		PhoneNumber: testPhone,
		OTP:         "5555",
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
			name: "Sad case - invalid OTP",
			args: args{
				ctx:     ctx,
				payload: OTPnotFound,
			},
			want:    false,
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

func TestPGInstance_GetStaffProfileByUserID(t *testing.T) {
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
			name: "Happy Case - Successfully get staff profile",
			args: args{
				ctx:    ctx,
				userID: userIDtoAssignStaff,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get staff profile",
			args: args{
				ctx:    ctx,
				userID: invalidID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetStaffProfileByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected a staff profile to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a staff profile not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_SearchStaffProfileByStaffNumber(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx         context.Context
		staffNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				staffNumber: staffNumber,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchStaffProfile(tt.args.ctx, tt.args.staffNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchStaffProfileByStaffNumber() error = %v, wantErr %v", err, tt.wantErr)
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
		UserID:         userID2,
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
				userID: userID2,
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
		{
			name: "invalid: empty user ID and contact type",
			args: args{
				ctx:         ctx,
				userID:      nil,
				contactType: "",
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
		UserID:         userID2,
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

func TestPGInstance_CheckIfUserHasViewedContent(t *testing.T) {
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
			name: "Happy case: check that user has viewed content",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: contentID,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: invalid user ID",
			args: args{
				ctx:       ctx,
				userID:    "userID",
				contentID: contentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: invalid content ID",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: -101,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIfUserHasViewedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfUserHasViewedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfUserHasViewedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserProfileByUserID(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()

	type args struct {
		ctx    context.Context
		userID *string
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
				userID: &userID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: &UUID,
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty user ID",
			args: args{
				ctx:    ctx,
				userID: nil,
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

	sharedAt := time.Now().Add(time.Hour * -25)
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
		SharedAt:              &sharedAt,
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
		SharedAt:              &sharedAt,
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
		UserID:         userID2,
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
				userID:    userID2,
				contentID: bookmarkInput.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid user id",
			args: args{
				ctx:       ctx,
				userID:    "userID",
				contentID: bookmarkInput.ContentID,
			},
			want:    false,
			wantErr: true,
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
			_, err := testingDB.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetUserProfileByPhoneNumber(t *testing.T) {
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
			name: "happy case: get user profile by phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: testPhone,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserProfileByPhoneNumber(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
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
			name: "happy case: get user pin by user id",
			args: args{
				ctx:     context.Background(),
				userID:  userID,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid-flavour",
			args: args{
				ctx:     context.Background(),
				userID:  userID,
				flavour: "Invalid-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserPINByUserID(tt.args.ctx, tt.args.userID, tt.args.flavour)
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

func TestPGInstance_GetSecurityQuestionResponse(t *testing.T) {
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
			name: "happy case: get security question response by id",
			args: args{
				ctx:        context.Background(),
				questionID: securityQuestionID,
				userID:     userID,
			},
			wantErr: true, // TODO @maxwellgithinji investigate why this is passing locally but not on the CI
		},
		{
			name: "sad case: invalid question id",
			args: args{
				ctx:        context.Background(),
				questionID: "securityQuestionID",
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid user ID",
			args: args{
				ctx:        context.Background(),
				questionID: "securityQuestionID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetSecurityQuestionResponse(tt.args.ctx, tt.args.questionID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
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
		{
			name: "Sad case: invalid flavour",
			args: args{
				ctx:     context.Background(),
				userID:  userID,
				flavour: "invalid",
			},
			want:    false,
			wantErr: true,
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

func TestPGInstance_GetClientCaregiver(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
	}

	caregiver := &gorm.Caregiver{}

	err := testingDB.DB.Where("id = ?", testCaregiverID).First(&caregiver).Error
	if err != nil {
		t.Errorf("failed to get caregiver: %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    *gorm.Caregiver
		wantErr bool
	}{
		{
			name: "happy case: get client caregiver",
			args: args{
				ctx:         context.Background(),
				caregiverID: testCaregiverID,
			},

			want:    caregiver,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientCaregiver(tt.args.ctx, tt.args.caregiverID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PGInstance.GetClientCaregiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientByClientID(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Client
		wantErr bool
	}{
		{
			name: "happy case: get client by client id",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientProfileByClientID(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientProfileByClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetClientsPendingServiceRequestsCount(t *testing.T) {
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
			name: "Happy case",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Happy case - No request type",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.GetClientsPendingServiceRequestsCount(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsPendingServiceRequestsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestPGInstance_GetCurrentTerms(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	termsID := gofakeit.Number(50, 9999999)
	now := time.Now()
	future := time.Now().AddDate(0, 0, 2)
	termsOfServiceInput := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &termsText,
		Flavour:   feedlib.FlavourPro,
		ValidFrom: &now,
		ValidTo:   &future,
		Active:    true,
	}

	err = pg.DB.Create(termsOfServiceInput).Error
	if err != nil {
		t.Errorf("Create terms of service failed: %v", err)
	}

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.TermsOfService
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetCurrentTerms(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}

	// TearDown
	if err = pg.DB.Where("id", termsOfServiceInput.TermsID).Unscoped().Delete(&gorm.TermsOfService{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_GetServiceRequests(t *testing.T) {
	var requesttype = "RED_FLAG"
	var requeststatus = "PENDING"

	type args struct {
		ctx           context.Context
		requestType   *string
		requestStatus *string
		facilityID    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.ClientServiceRequest
		wantErr bool
	}{
		{
			name: "happy case: get service requests by type",
			args: args{
				ctx:         context.Background(),
				requestType: &requesttype,
				facilityID:  facilityID,
			},
			wantErr: false,
		},
		{
			name: "happy case: get service requests by status",
			args: args{
				ctx:           context.Background(),
				requestStatus: &requeststatus,
				facilityID:    facilityID,
			},
			wantErr: false,
		},
		{
			name: "happy case: get service requests by type and status",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: &requeststatus,
				facilityID:    facilityID,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_CheckUserRole(t *testing.T) {
	ctx := context.Background()

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
			name: "happy case: user has role",
			args: args{
				ctx:    ctx,
				userID: userWithRolesID,
				role:   systemAdminRole,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case: user does not have role",
			args: args{
				ctx:    ctx,
				userID: userID2,
				role:   systemAdminRole,
			},
			wantErr: false,
			want:    false,
		},

		{
			name: "sad case: invalid user passed",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				role:   systemAdminRole,
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "sad case: invalid role passed",
			args: args{
				ctx:    ctx,
				userID: userID,
				role:   uuid.New().String(),
			},
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CheckUserRole(tt.args.ctx, tt.args.userID, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckUserRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckUserRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CheckUserPermission(t *testing.T) {
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
			name: "happy case: user has permission",
			args: args{
				ctx:        context.Background(),
				userID:     userWithRolesID,
				permission: string(enums.PermissionTypeCanEditOwnRole),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case: user does not have permission",
			args: args{
				ctx:        context.Background(),
				userID:     userID2,
				permission: string(enums.PermissionTypeCanEditOwnRole),
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "sad case: invalid user passed",
			args: args{
				ctx:        context.Background(),
				userID:     uuid.New().String(),
				permission: string(enums.PermissionTypeCanEditOwnRole),
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "sad case: invalid permission passed",
			args: args{
				ctx:        context.Background(),
				userID:     userID,
				permission: uuid.New().String(),
			},
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckUserPermission(tt.args.ctx, tt.args.userID, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckUserPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckUserPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserRoles(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*gorm.AuthorityRole
	}{
		{
			name: "happy case: user has roles",
			args: args{
				ctx:    context.Background(),
				userID: userWithRolesID,
			},
			wantErr: false,
		},
		{
			name: "happy case: user does not have roles",
			args: args{
				ctx:    context.Background(),
				userID: userID2,
			},
			wantErr: false,
		},
		{
			name: "sad case: user does not exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserRoles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetUserRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserPermissions(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.AuthorityPermission
		wantErr bool
	}{
		{
			name: "happy case: user has permissions",
			args: args{
				ctx:    context.Background(),
				userID: userWithRolesID,
			},
			wantErr: false,
		},
		{
			name: "happy case: user does not have permissions",
			args: args{
				ctx:    context.Background(),
				userID: userID2,
			},
			wantErr: false,
		},
		{
			name: "sad case: user does not exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserPermissions(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetUserRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetCommunityByID(t *testing.T) {
	ctx := context.Background()

	communityInput := &gorm.Community{
		ID:             communityID,
		Name:           "test",
		Description:    "test",
		Active:         true,
		MinimumAge:     19,
		MaximumAge:     30,
		Gender:         []string{"MALE"},
		ClientTypes:    []string{"PMTCT"},
		InviteOnly:     true,
		Discoverable:   true,
		OrganisationID: uuid.New().String(),
	}

	community, err := testingDB.CreateCommunity(ctx, communityInput)
	if err != nil {
		t.Errorf("unable to create community: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		communityID string
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
				communityID: community.ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				communityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetCommunityByID(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCommunityByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetCommunityByID() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CheckIdentifierExists(t *testing.T) {
	ctx := context.Background()

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
			name: "Sad case: identifier doesn't exist",
			args: args{
				ctx:             ctx,
				identifierType:  "CCC",
				identifierValue: "no-exist",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Happy case: identifier exists",
			args: args{
				ctx:             ctx,
				identifierType:  "CCC",
				identifierValue: "123456",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIdentifierExists(tt.args.ctx, tt.args.identifierType, tt.args.identifierValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIdentifierExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIdentifierExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientsByParams(t *testing.T) {

	unknown := "a908710b-d3d2-4ca7-bf99-5eda47356bd8"
	syncTime := time.Now()

	type args struct {
		ctx          context.Context
		params       gorm.Client
		lastSyncTime *time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Client
		wantErr bool
	}{
		{
			name: "happy case: retrieve non existent client",
			args: args{
				ctx: context.Background(),
				params: gorm.Client{
					ID: &unknown,
				},
				lastSyncTime: nil,
			},
			want:    []*gorm.Client{},
			wantErr: false,
		},
		{
			name: "happy case: retrieve non existent client",
			args: args{
				ctx: context.Background(),
				params: gorm.Client{
					ID: &unknown,
				},
				lastSyncTime: &syncTime,
			},
			want:    []*gorm.Client{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientsByParams(tt.args.ctx, tt.args.params, tt.args.lastSyncTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsByParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PGInstance.GetClientsByParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientCCCIdentifier(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Identifier
		wantErr bool
	}{
		{
			name: "happy case: retrieve ccc identifier",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
			},
			wantErr: false,
		},
		{
			name: "sad case: client has no identifier",
			args: args{
				ctx:      context.Background(),
				clientID: clientID2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientCCCIdentifier(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientCCCIdentifier() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPGInstance_GetServiceRequestsForKenyaEMR(t *testing.T) {
	ctx := context.Background()
	currentTime := time.Now()

	requestID := uuid.New().String()
	serviceRequest := &gorm.ClientServiceRequest{
		ID:             &requestID,
		Active:         true,
		RequestType:    "RED_FLAG",
		Request:        "SAD",
		Status:         "PENDING",
		InProgressAt:   &currentTime,
		ResolvedAt:     &currentTime,
		ClientID:       clientID,
		InProgressByID: &staffID,
		OrganisationID: uuid.New().String(),
		ResolvedByID:   &staffID,
		FacilityID:     facilityID,
		Meta:           `{}`,
	}

	err := testingDB.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		t.Errorf("an error occurred %v", err)
		return
	}

	type args struct {
		ctx          context.Context
		facilityID   string
		lastSyncTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.ClientServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:          ctx,
				facilityID:   facilityID,
				lastSyncTime: time.Now().AddDate(0, 0, 10),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:          ctx,
				facilityID:   "123Q4",
				lastSyncTime: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetServiceRequestsForKenyaEMR(tt.args.ctx, tt.args.facilityID, tt.args.lastSyncTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetServiceRequestsForKenyaEMR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetServiceRequestsForKenyaEMR() error = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetServiceRequestsForKenyaEMR() error = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolsQuestions(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		toolType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:      ctx,
				toolType: enums.ScreeningToolTypeTB.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetScreeningToolQuestions(tt.args.ctx, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetScreeningToolQuestions() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CheckFacilityExistsByMFLCode(t *testing.T) {
	ctx := context.Background()

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
				ctx:     ctx,
				MFLCode: mflCode,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				MFLCode: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckFacilityExistsByMFLCode(tt.args.ctx, tt.args.MFLCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckFacilityExistsByMFLCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckFacilityExistsByMFLCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientsInAFacility(t *testing.T) {
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
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientsInAFacility(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsInAFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPGInstance_GetRecentHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx          context.Context
		lastSyncTime time.Time
		clientID     string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.ClientHealthDiaryEntry
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:          ctx,
				lastSyncTime: time.Now().AddDate(0, 0, 10),
				clientID:     clientID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				clientID: "qwert123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetRecentHealthDiaryEntries(tt.args.ctx, tt.args.lastSyncTime, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetRecentHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected client health diary to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected client health diary not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolQuestionByQuestionID(t *testing.T) {
	ctx := context.Background()
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
			name: "happy case",
			args: args{
				ctx:        ctx,
				questionID: screeningToolsQuestionID,
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid uuid",
			args: args{
				ctx:        ctx,
				questionID: "123Q4",
			},
			wantErr: true,
		},
		{
			name: "sad case: question not found",
			args: args{
				ctx:        ctx,
				questionID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetScreeningToolQuestionByQuestionID(tt.args.ctx, tt.args.questionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolQuestionByQuestionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetScreeningToolQuestions() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_ListAppointments(t *testing.T) {

	type args struct {
		ctx        context.Context
		params     *gorm.Appointment
		filters    []*firebasetools.FilterParam
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Appointment
		wantErr bool
	}{
		{
			name: "happy case: list all appointments",
			args: args{
				ctx:        context.Background(),
				params:     nil,
				filters:    nil,
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: list paginated appointments",
			args: args{
				ctx:     context.Background(),
				params:  nil,
				filters: nil,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: list filtered appointments date",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "date",
						FieldType:           enumutils.FieldTypeTimestamp,
						ComparisonOperation: enumutils.OperationLessThanOrEqualTo,
						FieldValue:          time.Now().Format(time.RFC3339),
					},
				},
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: list filtered appointments time",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "created",
						FieldType:           enumutils.FieldTypeTimestamp,
						ComparisonOperation: enumutils.OperationLessThanOrEqualTo,
						FieldValue:          time.Now().Format(time.RFC3339),
					},
				},
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: list filtered appointments active",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "active",
						FieldType:           enumutils.FieldTypeBoolean,
						ComparisonOperation: enumutils.OperationLessThanOrEqualTo,
						FieldValue:          false,
					},
				},
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "sad case: list filtered appointments invalid filter",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "active",
						FieldType:           "INVALID",
						ComparisonOperation: enumutils.OperationEqual,
						FieldValue:          false,
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
		{
			name: "sad case: list filtered appointments invalid operation",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "active",
						FieldType:           enumutils.FieldTypeBoolean,
						ComparisonOperation: "INVALID",
						FieldValue:          false,
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
		{
			name: "sad case: list filtered appointments invalid boolean filter",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "active",
						FieldType:           enumutils.FieldTypeBoolean,
						ComparisonOperation: enumutils.OperationEqual,
						FieldValue:          "INVALID",
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
		{
			name: "sad case: list filtered appointments invalid string filter",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "status",
						FieldType:           enumutils.FieldTypeString,
						ComparisonOperation: enumutils.OperationEqual,
						FieldValue:          1234,
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
		{
			name: "sad case: list filtered appointments invalid integer filter",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "status",
						FieldType:           enumutils.FieldTypeInteger,
						ComparisonOperation: enumutils.OperationEqual,
						FieldValue:          "INVALID",
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
		{
			name: "sad case: list filtered appointments invalid timestamp filter",
			args: args{
				ctx:    context.Background(),
				params: nil,
				filters: []*firebasetools.FilterParam{
					{
						FieldName:           "date",
						FieldType:           enumutils.FieldTypeTimestamp,
						ComparisonOperation: enumutils.OperationEqual,
						FieldValue:          123456,
					},
				},
				pagination: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.ListAppointments(tt.args.ctx, tt.args.params, tt.args.filters, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListAppointments() error = %v, wantErr %v", err, tt.wantErr)
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
		})
	}
}

func TestPGInstance_ListNotifications(t *testing.T) {

	type args struct {
		ctx        context.Context
		params     *gorm.Notification
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Notification
		wantErr bool
	}{
		{
			name: "happy case: list user notifications",
			args: args{
				ctx: context.Background(),
				params: &gorm.Notification{
					UserID:  &userID,
					Flavour: feedlib.FlavourConsumer,
				},
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: list facility notifications",
			args: args{
				ctx: context.Background(),
				params: &gorm.Notification{
					UserID:     &userIDtoAssignStaff,
					FacilityID: &facilityID,
					Flavour:    feedlib.FlavourPro,
				},
				pagination: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: list user paginated notifications",
			args: args{
				ctx: context.Background(),
				params: &gorm.Notification{
					UserID:  &userID,
					Flavour: feedlib.FlavourConsumer,
				},
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.ListNotifications(tt.args.ctx, tt.args.params, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected notifications to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected notifications not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetClientProfileByCCCNumber(t *testing.T) {
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
				CCCNumber: "123456",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to get client profile by CCC number",
			args: args{
				ctx:       ctx,
				CCCNumber: "3232873827382",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientProfileByCCCNumber(tt.args.ctx, tt.args.CCCNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientProfileByCCCNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetClientProfileByCCCNumber() Expected a response but got = %v", got)
			}
		})
	}
}

func TestPGInstance_CheckIfClientHasUnresolvedServiceRequests(t *testing.T) {
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
			name: "happy case: check pending pin reset",
			args: args{
				ctx:                context.Background(),
				clientID:           clientUnresolvedRequestID,
				serviceRequestType: string(enums.ServiceRequestTypePinReset),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: check pending pin reset, no pending pin reset request",
			args: args{
				ctx:                context.Background(),
				clientID:           clientID,
				serviceRequestType: string(enums.ServiceRequestTypePinReset),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "sad case: invalid client id",
			args: args{
				ctx:                context.Background(),
				clientID:           "123Q4",
				serviceRequestType: string(enums.ServiceRequestTypePinReset),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CheckIfClientHasUnresolvedServiceRequests(tt.args.ctx, tt.args.clientID, tt.args.serviceRequestType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfClientHasUnresolvedServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfClientHasUnresolvedServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetAllRoles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get all roles",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetAllRoles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetAllRoles() Expected a response but got = %v", got)
			}
		})
	}
}

func TestPGInstance_SearchClientProfilesByCCCNumber(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx             context.Context
		searchParameter string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Client
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profiles by CCC number",
			args: args{
				ctx:             ctx,
				searchParameter: "123456",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchClientProfile(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchClientProfile() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPGInstance_GetUserProfileByStaffID(t *testing.T) {
	type args struct {
		ctx     context.Context
		staffID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user profile by staff ID",
			args: args{
				ctx:     context.Background(),
				staffID: staffID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to get user profile by staff ID, invalid uuid",
			args: args{
				ctx:     context.Background(),
				staffID: "123456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserProfileByStaffID(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserProfileByStaffID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetUserProfileByStaffID() Expected a response but got = %v", got)
			}
		})
	}
}

func TestPGInstance_GetHealthDiaryEntryByID(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                context.Context
		healthDiaryEntryID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.ClientHealthDiaryEntry
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: clientsHealthDiaryEntryID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: "entryID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetHealthDiaryEntryByID(tt.args.ctx, tt.args.healthDiaryEntryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetHealthDiaryEntryByID() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPGInstance_CheckIfUsernameExists(t *testing.T) {
	ctx := context.Background()

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
			name: "Happy case",
			args: args{
				ctx:      ctx,
				username: "test user",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIfUsernameExists(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfUsernameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfUsernameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetStaffProfileByStaffID(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		staffID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				staffID: staffID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				staffID: "staffID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetStaffProfileByStaffID(tt.args.ctx, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffProfileByStaffID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected staff not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetStaffPendingServiceRequestsCount(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ServiceRequestsCount
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
			name: "Happy case",
			args: args{
				ctx:        ctx,
				facilityID: "facilityID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetStaffPendingServiceRequestsCount(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffPendingServiceRequestsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected staff not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetStaffServiceRequests(t *testing.T) {
	ctx := context.Background()

	var requesttype = "STAFF_PIN_RESET"
	var requestStatus = "PENDING"

	type args struct {
		ctx           context.Context
		requestType   *string
		requestStatus *string
		facilityID    string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.StaffServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:           ctx,
				requestType:   &requesttype,
				requestStatus: &requestStatus,
				facilityID:    facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid facility ID",
			args: args{
				ctx:           ctx,
				requestType:   &requesttype,
				requestStatus: &requestStatus,
				facilityID:    "facilityID",
			},
			wantErr: true,
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
			got, err := testingDB.GetStaffServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected staff not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetServiceRequestByID(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.ClientServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				serviceRequestID: serviceRequestID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:              ctx,
				serviceRequestID: "serviceRequestID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetServiceRequestByID(tt.args.ctx, tt.args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetServiceRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected service request to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected service request not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetAppointmentServiceRequests(t *testing.T) {
	type args struct {
		ctx          context.Context
		lastSyncTime time.Time
		facilityID   string
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
				facilityID:   facilityID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetAppointmentServiceRequests(tt.args.ctx, tt.args.lastSyncTime, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetAppointmentServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected service requests not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetFacilitiesWithoutFHIRID(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Facility
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
			got, err := testingDB.GetFacilitiesWithoutFHIRID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetFacilitiesWithoutFHIRID() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPGInstance_GetClientServiceRequests(t *testing.T) {
	type args struct {
		ctx         context.Context
		requestType string
		status      string
		clientID    string
		FacilityID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         context.Background(),
				requestType: string(enums.ServiceRequestTypeScreeningToolsRedFlag),
				status:      string(enums.ServiceRequestStatusResolved),
				clientID:    clientID,
				FacilityID:  facilityID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.status, tt.args.clientID, tt.args.FacilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected service requests not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetActiveScreeningToolResponses(t *testing.T) {
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
			name: "Happy case",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: clientID is invalid",
			args: args{
				ctx:      context.Background(),
				clientID: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetActiveScreeningToolResponses(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetActiveScreeningToolResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected service requests not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_CheckAppointmentExistsByExternalID(t *testing.T) {
	type args struct {
		ctx        context.Context
		externalID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: existing appointment",
			args: args{
				ctx:        context.Background(),
				externalID: externalAppointmentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case: non existent appointment",
			args: args{
				ctx:        context.Background(),
				externalID: "non-existent",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckAppointmentExistsByExternalID(tt.args.ctx, tt.args.externalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckAppointmentExistsByExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckAppointmentExistsByExternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetAnsweredScreeningToolQuestions(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		toolType   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.ScreeningToolsResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				toolType:   "TB_ASSESSMENT",
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        ctx,
				toolType:   "",
				facilityID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid tool type",
			args: args{
				ctx:        ctx,
				toolType:   "INVALID-TOOL-TYPE",
				facilityID: "facilityID",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid facilityID",
			args: args{
				ctx:        ctx,
				toolType:   "INVALID-TOOL-TYPE",
				facilityID: "facilityID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.GetAnsweredScreeningToolQuestions(tt.args.ctx, tt.args.facilityID, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetAnsweredScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetSharedHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.ClientHealthDiaryEntry
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				clientID:   clientID,
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid facility",
			args: args{
				ctx:        ctx,
				facilityID: gofakeit.HipsterSentence(44),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility",
			args: args{
				ctx:        ctx,
				facilityID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid client ID",
			args: args{
				ctx:        ctx,
				clientID:   gofakeit.HipsterSentence(45),
				facilityID: facilityID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetSharedHealthDiaryEntries(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSharedHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected shared health diary entries to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected shared health diary entries not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetClientScreeningToolResponsesByToolType(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
		toolType string
		active   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
				toolType: string(enums.ScreeningToolTypeGBV),
				active:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientScreeningToolResponsesByToolType(tt.args.ctx, tt.args.clientID, tt.args.toolType, tt.args.active)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientScreeningToolResponsesByToolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetAppointment(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *gorm.Appointment
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Appointment
		wantErr bool
	}{
		{
			name: "Happy case: retrieve an appointment",
			args: args{
				ctx: context.Background(),
				params: &gorm.Appointment{
					ClientID: clientID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid param to retrieve an appointment",
			args: args{
				ctx: context.Background(),
				params: &gorm.Appointment{
					ClientID: "-",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetAppointment(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetAppointment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected appointment to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected appointment not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_GetClientScreeningToolServiceRequestByToolType(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
		toolType string
		status   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
				toolType: enums.ScreeningToolTypeGBV.String(),
				status:   string(enums.ServiceRequestStatusPending),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientScreeningToolServiceRequestByToolType(tt.args.ctx, tt.args.clientID, tt.args.toolType, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientScreeningToolServiceRequestByToolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_CheckIfStaffHasUnresolvedServiceRequests(t *testing.T) {
	type args struct {
		ctx                context.Context
		staffID            string
		serviceRequestType string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: check pending pin reset",
			args: args{
				ctx:                context.Background(),
				staffID:            staffUnresolvedRequestID,
				serviceRequestType: string(enums.ServiceRequestTypeStaffPinReset),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: check pending pin reset, no pending pin reset request",
			args: args{
				ctx:                context.Background(),
				staffID:            staffID,
				serviceRequestType: string(enums.ServiceRequestTypeStaffPinReset),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "sad case: invalid client id",
			args: args{
				ctx:                context.Background(),
				staffID:            "123Q4",
				serviceRequestType: string(enums.ServiceRequestTypeStaffPinReset),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIfStaffHasUnresolvedServiceRequests(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestType)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfStaffHasUnresolvedServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfStaffHasUnresolvedServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetFacilityStaffs(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility staff",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid facility id",
			args: args{
				ctx:        context.Background(),
				facilityID: "-",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetFacilityStaffs(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetFacilityStaffs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetUserSurveyForms(t *testing.T) {
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
				userID: userWithRolesID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: gofakeit.BeerAlcohol(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserSurveyForms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserSurveyForms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetNotification(t *testing.T) {

	type args struct {
		ctx            context.Context
		notificationID string
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Notification
		wantErr bool
	}{
		{
			name: "happy case: retrieve notification",
			args: args{
				ctx:            context.Background(),
				notificationID: notificationID,
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid notification id",
			args: args{
				ctx:            context.Background(),
				notificationID: "non-existent-id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetNotification(tt.args.ctx, tt.args.notificationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetClientsByFilterParams(t *testing.T) {
	type args struct {
		ctx          context.Context
		facilityID   string
		filterParams *dto.ClientFilterParamsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility clients",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve facility clients by client type and age range, gender",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve facility clients by client type",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve facility clients by age range",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve facility clients by gender",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					Gender: []enumutils.Gender{enumutils.GenderMale, enumutils.GenderFemale, enumutils.GenderOther},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: retrieve facility clients by client type and age range, gender, invalid facility id",
			args: args{
				ctx:        context.Background(),
				facilityID: "facilityID",
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetClientsByFilterParams(tt.args.ctx, tt.args.facilityID, tt.args.filterParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsByFilterParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}
