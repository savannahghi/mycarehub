package gorm_test

import (
	"context"
	"reflect"
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
				ctx:      context.Background(),
				id:       &facilityID,
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to retrieve facility",
			args: args{
				ctx:      context.Background(),
				id:       &fakeID,
				isActive: false,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - nil facility",
			args: args{
				ctx:      context.Background(),
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

func TestPGInstance_RetrieveFacilityByIdentifier(t *testing.T) {

	type args struct {
		ctx        context.Context
		identifier *gorm.FacilityIdentifier
		isActive   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully fetch facility by MFL code",
			args: args{
				ctx: context.Background(),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: mflIdentifier,
				},
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to fetch facility by MFL code",
			args: args{
				ctx: context.Background(),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: "88888888",
				},
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RetrieveFacilityByIdentifier(tt.args.ctx, tt.args.identifier, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RetrieveFacilityByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
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
	emptySearchTerm := ""
	searchTerm := "ny"
	noSearchTerm := "this will not be found"
	type args struct {
		ctx        context.Context
		searchTerm *string
		filter     []*domain.FiltersParam
		pagination *domain.Pagination
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        context.Background(),
				searchTerm: nil,
				filter:     nil,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Happy case: search facility",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				filter:     nil,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Happy case: search facility, empty",
			args: args{
				ctx:        context.Background(),
				searchTerm: &emptySearchTerm,
				filter:     nil,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Happy case: search facility, no results",
			args: args{
				ctx:        context.Background(),
				searchTerm: &noSearchTerm,
				filter:     nil,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Happy case: filter",
			args: args{
				ctx:        context.Background(),
				searchTerm: nil,
				filter: []*domain.FiltersParam{
					{
						Name:     "country",
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Sad case: filter, invalid filter",
			args: args{
				ctx:        context.Background(),
				searchTerm: nil,
				filter: []*domain.FiltersParam{
					{
						Name:     "invalid",
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filter, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("PGInstance.ListFacilities() got = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestPGInstance_SearchFacility(t *testing.T) {
	searchParameter := "Kenya"
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
			name:    "Happy Case - search facility",
			args:    args{ctx: context.Background(), searchParameter: &searchParameter},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchFacility(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchFacility() error = %v, wantErr %v", err, tt.wantErr)
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
				ctx:     context.Background(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     context.Background(),
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
				ctx:                context.Background(),
				securityQuestionID: &securityQuestionID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid security question id",
			args: args{
				ctx:                context.Background(),
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
				ctx:       context.Background(),
				phone:     testPhone,
				isOptedIn: true,
				flavour:   testFlavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:       context.Background(),
				phone:     testPhone,
				isOptedIn: true,
				flavour:   "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - ungistered phone",
			args: args{
				ctx:       context.Background(),
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
				ctx:     context.Background(),
				payload: validOTPPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - invalid OTP",
			args: args{
				ctx:     context.Background(),
				payload: OTPnotFound,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:     context.Background(),
				payload: invalidOTPPayload2,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no otp",
			args: args{
				ctx:     context.Background(),
				payload: invalidOTPPayload3,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     context.Background(),
				payload: invalidOTPPayload4,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad inputs",
			args: args{
				ctx:     context.Background(),
				payload: invalidOTPPayload5,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad inputs",
			args: args{
				ctx:     context.Background(),
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

func TestPGInstance_GetClientProfile(t *testing.T) {

	invalidID := uuid.New().String()

	type args struct {
		ctx       context.Context
		userID    string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profile",
			args: args{
				ctx:       context.Background(),
				userID:    userID,
				programID: programID,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get client profile with empty programID",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:       context.Background(),
				userID:    invalidID,
				programID: "programID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientProfile(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetStaffProfile(t *testing.T) {

	invalidID := uuid.New().String()

	type args struct {
		ctx       context.Context
		userID    string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get staff profile",
			args: args{
				ctx:       context.Background(),
				userID:    userIDtoAssignStaff,
				programID: programID,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get staff profile with empty programID",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignStaff,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get staff profile",
			args: args{
				ctx:    context.Background(),
				userID: invalidID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetStaffProfile(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
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
				ctx:         context.Background(),
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
				ctx:         context.Background(),
				phoneNumber: testPhone,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get OTP, nonexisting phone",
			args: args{
				ctx:         context.Background(),
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

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	securityQuestionResponseInput := &gorm.SecurityQuestionResponse{
		UserID:     userID2,
		QuestionID: securityQuestionID,
		Response:   "1917",
		Timestamp:  time.Now(),
	}

	err = pg.DB.WithContext(addRequiredContext(context.Background(), t)).Create(securityQuestionResponseInput).Error
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
				ctx:    context.Background(),
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
				ctx:         context.Background(),
				userID:      &userID,
				contactType: "PHONE",
			},
		},
		{
			name: "invalid: contact type invalid",
			args: args{
				ctx:         context.Background(),
				userID:      &userID,
				contactType: "Invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid: user ID invalid",
			args: args{
				ctx:         context.Background(),
				userID:      &ID,
				contactType: "EMAIL",
			},
			wantErr: true,
		},
		{
			name: "invalid: empty user ID and contact type",
			args: args{
				ctx:         context.Background(),
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

func TestPGInstance_GetUserProfileByUserID(t *testing.T) {

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
				ctx:    context.Background(),
				userID: &userID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    context.Background(),
				userID: &UUID,
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty user ID",
			args: args{
				ctx:    context.Background(),
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

	type args struct {
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Get a random quote",
			args: args{
				ctx:   context.Background(),
				limit: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetClientHealthDiaryQuote(tt.args.ctx, tt.args.limit)
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
}

func TestPGInstance_CanRecordHeathDiary(t *testing.T) {

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
		ProgramID:             programID,
	}
	err = pg.DB.WithContext(addRequiredContext(context.Background(), t)).Create(&canShowHealthDiaryInput).Error
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
		ProgramID:             programID,
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
				ctx:      context.Background(),
				clientID: clientID2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case, cant record new entry before 24 hours",
			args: args{
				ctx:      context.Background(),
				clientID: clientID2,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Invalid: id not type UUID",
			args: args{
				ctx:      context.Background(),
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
				err = pg.DB.WithContext(addRequiredContext(context.Background(), t)).Create(&cannotShowHealthDiaryInput).Error
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
		{
			name: "happy case: get user pin by user id",
			args: args{
				ctx:    context.Background(),
				userID: "userID",
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
			name: "happy case: check user has pin",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: invalid user ID",
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckUserHasPin(tt.args.ctx, tt.args.userID)
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

func TestPGInstance_GetClientHealthDiaryEntries(t *testing.T) {
	params := map[string]interface{}{
		"client_id": clientID,
	}
	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get client health diary entries",
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.params)
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
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Happy case - No request type",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - unable to get clients service request count with invalid facility ID",
			args: args{
				ctx:        context.Background(),
				facilityID: "facilityID",
			},
			wantErr: true,
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
	type args struct {
		ctx context.Context
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
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetCurrentTerms(tt.args.ctx)
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

func TestPGInstance_GetCommunityByID(t *testing.T) {

	communityInput := &gorm.Community{
		ID:             communityID,
		Name:           "test",
		Description:    "test",
		Active:         true,
		MinimumAge:     19,
		MaximumAge:     30,
		Gender:         []string{"MALE"},
		ClientTypes:    []string{"PMTCT"},
		OrganisationID: orgID,
		ProgramID:      programID,
	}

	community, err := testingDB.CreateCommunity(addRequiredContext(context.Background(), t), communityInput)
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
				ctx:         context.Background(),
				communityID: community.ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         context.Background(),
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
				ctx:             context.Background(),
				identifierType:  "CCC",
				identifierValue: "no-exist",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Happy case: identifier exists",
			args: args{
				ctx:             context.Background(),
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
		OrganisationID: orgID,
		ResolvedByID:   &staffID,
		FacilityID:     facilityID,
		Meta:           `{}`,
		ProgramID:      programID,
	}

	err := testingDB.CreateServiceRequest(addRequiredContext(context.Background(), t), serviceRequest)
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
				ctx:          context.Background(),
				facilityID:   facilityID,
				lastSyncTime: time.Now().AddDate(0, 0, 10),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:          context.Background(),
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

func TestPGInstance_CheckFacilityExistsByIdentifier(t *testing.T) {

	type args struct {
		ctx        context.Context
		identifier *gorm.FacilityIdentifier
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check identifier exist",
			args: args{
				ctx: context.Background(),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: "98329839208309",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: record not found",
			args: args{
				ctx: context.Background(),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: "does not exist",
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckFacilityExistsByIdentifier(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckFacilityExistsByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckFacilityExistsByIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientsInAFacility(t *testing.T) {

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
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        context.Background(),
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
				ctx:          context.Background(),
				lastSyncTime: time.Now().AddDate(0, 0, 10),
				clientID:     clientID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      context.Background(),
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
		filters    []*firebasetools.FilterParam
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
				filters:    nil,
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
				filters:    nil,
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
				filters: nil,
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
			got, _, err := testingDB.ListNotifications(tt.args.ctx, tt.args.params, tt.args.filters, tt.args.pagination)
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

func TestPGInstance_ListAvailableNotificationTypes(t *testing.T) {
	type args struct {
		ctx    context.Context
		params *gorm.Notification
	}
	tests := []struct {
		name    string
		args    args
		want    []enums.NotificationType
		wantErr bool
	}{
		{
			name: "happy case: list user notifications",
			args: args{
				ctx: context.Background(),
				params: &gorm.Notification{
					UserID:  &userID2,
					Flavour: feedlib.FlavourConsumer,
				},
			},
			want:    []enums.NotificationType{enums.NotificationTypeAppointment, enums.NotificationTypeServiceRequest},
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
			},
			want:    []enums.NotificationType{enums.NotificationTypeAppointment, enums.NotificationTypeServiceRequest},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ListAvailableNotificationTypes(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListAvailableNotificationTypes() error = %v, wantErr %v", err, tt.wantErr)
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
				ctx:       context.Background(),
				CCCNumber: "123456",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to get client profile by CCC number",
			args: args{
				ctx:       context.Background(),
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

func TestPGInstance_SearchClientProfilesByCCCNumber(t *testing.T) {

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
				ctx:             context.Background(),
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
				ctx:                context.Background(),
				healthDiaryEntryID: clientsHealthDiaryEntryID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                context.Background(),
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
				ctx:      context.Background(),
				username: "a_test_user",
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
				ctx:     context.Background(),
				staffID: staffID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     context.Background(),
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
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Happy case",
			args: args{
				ctx:        context.Background(),
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
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: &requestStatus,
				facilityID:    facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid facility ID",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: &requestStatus,
				facilityID:    "facilityID",
			},
			wantErr: true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
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
				ctx:              context.Background(),
				serviceRequestID: serviceRequestID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:              context.Background(),
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

// func TestPGInstance_GetFacilitiesWithoutFHIRID(t *testing.T) {

// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []*gorm.Facility
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: context.Background(),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := testingDB.GetFacilitiesWithoutFHIRID(tt.args.ctx)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("PGInstance.GetFacilitiesWithoutFHIRID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected client to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected client not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// }

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

func TestPGInstance_GetSharedHealthDiaryEntry(t *testing.T) {

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
				ctx:        context.Background(),
				clientID:   clientID,
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid facility",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.HipsterSentence(44),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility",
			args: args{
				ctx:        context.Background(),
				facilityID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid client ID",
			args: args{
				ctx:        context.Background(),
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

	hasSubmitted := false

	type args struct {
		ctx    context.Context
		params map[string]interface{}
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
				params: map[string]interface{}{
					"user_id":       userID,
					"project_id":    projectID,
					"form_id":       formID,
					"has_submitted": hasSubmitted,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"user_id":       gofakeit.BeerAlcohol(),
					"project_id":    projectID,
					"form_id":       formID,
					"has_submitted": hasSubmitted,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserSurveyForms(tt.args.ctx, tt.args.params)
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

func TestPGInstance_SearchClientServiceRequests(t *testing.T) {

	type args struct {
		ctx             context.Context
		searchParameter string
		requestType     string
		facilityID      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.ClientServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: search client service requests",
			args: args{
				ctx:             context.Background(),
				searchParameter: "PENDING",
				requestType:     "RED_FLAG",
				facilityID:      facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid search client service requests",
			args: args{
				ctx:             context.Background(),
				searchParameter: "PENDING",
				requestType:     gofakeit.HipsterParagraph(1, 10, 200, ""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchClientServiceRequests(tt.args.ctx, tt.args.searchParameter, tt.args.requestType, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchClientServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_SearchStaffServiceRequests(t *testing.T) {

	type args struct {
		ctx             context.Context
		searchParameter string
		requestType     string
		facilityID      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.StaffServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: search staff service requests",
			args: args{
				ctx:             context.Background(),
				searchParameter: "PENDING",
				requestType:     "STAFF_PIN_RESET",
				facilityID:      facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid search staff service requests",
			args: args{
				ctx:             context.Background(),
				searchParameter: "PENDING",
				requestType:     gofakeit.HipsterParagraph(1, 10, 200, ""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SearchStaffServiceRequests(tt.args.ctx, tt.args.searchParameter, tt.args.requestType, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchStaffServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetQuestionnaireByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get questionnaire by id",
			args: args{
				ctx: context.Background(),
				id:  questionnaireID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid get questionnaire by id",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetQuestionnaireByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetQuestionnaireByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool by id",
			args: args{
				ctx: context.Background(),
				id:  screeningToolID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get screening tool by id",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetScreeningToolByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetQuestionsByQuestionnaireID(t *testing.T) {
	type args struct {
		ctx             context.Context
		questionnaireID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get questions by questionnaire id",
			args: args{
				ctx:             context.Background(),
				questionnaireID: questionnaireID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get questions by questionnaire id",
			args: args{
				ctx:             context.Background(),
				questionnaireID: "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetQuestionsByQuestionnaireID(tt.args.ctx, tt.args.questionnaireID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetQuestionsByQuestionnaireID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetQuestionInputChoicesByQuestionID(t *testing.T) {
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
			name: "Happy case: get question input choices by question id",
			args: args{
				ctx:        context.Background(),
				questionID: questionID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get question input choices by question id",
			args: args{
				ctx:        context.Background(),
				questionID: "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetQuestionInputChoicesByQuestionID(tt.args.ctx, tt.args.questionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetQuestionInputChoicesByQuestionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetAvailableScreeningTools(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
		programID  string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.ScreeningTool
		wantErr bool
	}{
		{
			name: "Happy case: get available screening tools",
			args: args{
				ctx:        context.Background(),
				clientID:   clientID,
				facilityID: facilityID,
				programID:  programID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid get available screening tools",
			args: args{
				ctx:        context.Background(),
				clientID:   gofakeit.HipsterParagraph(1, 10, 200, ""),
				facilityID: gofakeit.HipsterParagraph(1, 10, 200, ""),
				programID:  gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.GetAvailableScreeningTools(tt.args.ctx, tt.args.clientID, tt.args.facilityID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetAvailableScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetFacilityRespondedScreeningTools(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilityID string
		programID  string
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get facility responded screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				programID:  programID,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid get facility responded screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.HipsterParagraph(1, 10, 200, ""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.GetFacilityRespondedScreeningTools(tt.args.ctx, tt.args.facilityID, tt.args.programID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetFacilityRespondedScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_ListSurveyRespondents(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilityID string
		params     map[string]interface{}
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list survey respondents",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"project_id":    projectID,
					"has_submitted": true,
				},
				facilityID: facilityID,
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list survey respondents",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"project_id":    gofakeit.HipsterParagraph(1, 10, 200, ""),
					"has_submitted": true,
				},
				facilityID: facilityID,
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := testingDB.ListSurveyRespondents(tt.args.ctx, tt.args.params, tt.args.facilityID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListSurveyRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolServiceRequestOfRespondents(t *testing.T) {
	type args struct {
		ctx             context.Context
		facilityID      string
		programID       string
		screeningToolID string
		searchTerm      string
		pagination      *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool service request by respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      facilityID,
				programID:       programID,
				screeningToolID: screeningToolID,
				searchTerm:      "",
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: get screening tool service request by respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      "INVALID_ID",
				programID:       programID,
				screeningToolID: "INVALID_ID",
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.GetScreeningToolServiceRequestOfRespondents(tt.args.ctx, tt.args.facilityID, tt.args.programID, tt.args.screeningToolID, tt.args.searchTerm, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolServiceRequestOfRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolResponseByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool response by id",
			args: args{
				ctx: context.Background(),
				id:  screeningToolsResponseID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get screening tool response by id",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetScreeningToolResponseByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolResponseByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetScreeningToolQuestionResponsesByResponseID(t *testing.T) {
	type args struct {
		ctx        context.Context
		responseID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool question responses by response id",
			args: args{
				ctx:        context.Background(),
				responseID: screeningToolsResponseID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get screening tool question responses by response id",
			args: args{
				ctx:        context.Background(),
				responseID: "INVALID_ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetScreeningToolQuestionResponsesByResponseID(tt.args.ctx, tt.args.responseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetScreeningToolQuestionResponsesByResponseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestPGInstance_GetSurveysWithServiceRequests(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.UserSurvey
		wantErr bool
	}{
		{
			name: "Happy case: get surveys with service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: get surveys with service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.HipsterSentence(100),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetSurveysWithServiceRequests(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetSurveysWithServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PGInstance.GetSurveysWithServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientsSurveyServiceRequest(t *testing.T) {

	type args struct {
		ctx        context.Context
		facilityID string
		projectID  int
		formID     string
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get clients service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				projectID:  projectID,
				formID:     formID,
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get clients service requests by facility ID",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.HipsterSentence(100),
				projectID:  projectID,
				formID:     formID,
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := testingDB.GetClientsSurveyServiceRequest(tt.args.ctx, tt.args.facilityID, tt.args.projectID, tt.args.formID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsSurveyServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetStaffFacilities(t *testing.T) {
	invalidID := "invalid"
	type args struct {
		ctx           context.Context
		staffFacility gorm.StaffFacilities
		pagination    *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get staff facilities by staff ID",
			args: args{
				ctx: context.Background(),
				staffFacility: gorm.StaffFacilities{
					StaffID: &staffID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: get staff facilities by facility ID",
			args: args{
				ctx: context.Background(),
				staffFacility: gorm.StaffFacilities{
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: get staff facilities by facility ID and staff ID",
			args: args{
				ctx: context.Background(),
				staffFacility: gorm.StaffFacilities{
					StaffID:    &staffID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get staff facilities, invalid staff ID",
			args: args{
				ctx: context.Background(),
				staffFacility: gorm.StaffFacilities{
					StaffID:    &invalidID,
					FacilityID: &facilityID,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get staff facilities, invalid facility ID",
			args: args{
				ctx: context.Background(),
				staffFacility: gorm.StaffFacilities{
					StaffID:    &staffID,
					FacilityID: &invalidID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.GetStaffFacilities(tt.args.ctx, tt.args.staffFacility, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestPGInstance_GetClientFacilities(t *testing.T) {
	invalidID := "invalid"
	type args struct {
		ctx            context.Context
		clientFacility gorm.ClientFacilities
		pagination     *domain.Pagination
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			name: "happy case: get client facilities by client ID",
			args: args{
				ctx: context.Background(),
				clientFacility: gorm.ClientFacilities{
					ClientID: &clientID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: get client facilities by facility ID",
			args: args{
				ctx: context.Background(),
				clientFacility: gorm.ClientFacilities{
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: get client facilities by facility ID and client ID",
			args: args{
				ctx: context.Background(),
				clientFacility: gorm.ClientFacilities{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get client facilities, invalid client ID",
			args: args{
				ctx: context.Background(),
				clientFacility: gorm.ClientFacilities{
					ClientID:   &invalidID,
					FacilityID: &facilityID,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get client facilities, invalid facility ID",
			args: args{
				ctx: context.Background(),
				clientFacility: gorm.ClientFacilities{
					ClientID:   &clientID,
					FacilityID: &invalidID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _, err := testingDB.GetClientFacilities(tt.args.ctx, tt.args.clientFacility, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestPGInstance_GetNotificationsCount(t *testing.T) {

	invalidID := "invalid"
	type args struct {
		ctx          context.Context
		notification gorm.Notification
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "happy case: get staffs notifications count",
			args: args{
				ctx: context.Background(),
				notification: gorm.Notification{
					UserID:     &userToRegisterStaff,
					IsRead:     false,
					FacilityID: &facilityID,
					Flavour:    feedlib.FlavourPro,
				},
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "sad case: get notifications count",
			args: args{
				ctx: context.Background(),
				notification: gorm.Notification{
					UserID:     &invalidID,
					IsRead:     false,
					FacilityID: &facilityID,
					Flavour:    feedlib.FlavourPro,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetNotificationsCount(tt.args.ctx, tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetNotificationsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.GetNotificationsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetClientsSurveyCount(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Happy case: get clients survey count",
			args: args{
				ctx:    context.Background(),
				userID: clientID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get clients survey count",
			args: args{
				ctx:    context.Background(),
				userID: clientID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientsSurveyCount(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientsSurveyCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.GetClientsSurveyCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_FindContacts(t *testing.T) {
	type args struct {
		ctx          context.Context
		contactType  string
		contactValue string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "happy case: find contacts",
			args: args{
				ctx:          context.Background(),
				contactType:  "PHONE",
				contactValue: testPhone2,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "happy case: no contacts",
			args: args{
				ctx:          context.Background(),
				contactType:  "PHONE",
				contactValue: gofakeit.Phone(),
			},
			wantCount: 0,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.FindContacts(tt.args.ctx, tt.args.contactType, tt.args.contactValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.FindContacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("PGInstance.FindContacts() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}

func TestPGInstance_GetOrganisation(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: fetch an organisation",
			args: args{
				ctx: context.Background(),
				id:  orgID,
			},
			wantErr: false,
		},
		{
			name: "sad case: fetch using invalid id",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetOrganisation(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got %s", err)
			}
		})
	}
}

func TestPGInstance_SearchCaregiverUser(t *testing.T) {
	type args struct {
		ctx             context.Context
		searchParameter string
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Caregiver
		wantErr bool
	}{
		{
			name: "happy case: search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: testCaregiverNumber,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.SearchCaregiverUser(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchCaregiverUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetCaregiverByUserID(t *testing.T) {
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
			name: "sad case: invalid user id",
			args: args{
				ctx:    context.Background(),
				userID: "invalid",
			},
			wantErr: true,
		},
		{
			name: "sad case: missing caregiver",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignStaff,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := testingDB.GetCaregiverByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCaregiverByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestPGInstance_GetCaregiversClient(t *testing.T) {
	type args struct {
		ctx             context.Context
		caregiverClient gorm.CaregiverClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: get caregiver's client",
			args: args{
				ctx: context.Background(),
				caregiverClient: gorm.CaregiverClient{
					CaregiverID: testCaregiverID,
					ClientID:    clientID,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case: get caregiver's client, no caregiver id",
			args: args{
				ctx: context.Background(),
				caregiverClient: gorm.CaregiverClient{
					ClientID: clientID,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case: get caregiver's client, no client id",
			args: args{
				ctx: context.Background(),
				caregiverClient: gorm.CaregiverClient{
					CaregiverID: testCaregiverID,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case: get caregiver's client, no caregiver client params",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid ID",
			args: args{
				ctx: context.Background(),
				caregiverClient: gorm.CaregiverClient{
					CaregiverID: testCaregiverID,
					ClientID:    "invalid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetCaregiversClient(tt.args.ctx, tt.args.caregiverClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCaregiversClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value but got %v", got)
			}
		})
	}
}

func TestPGInstance_GetCaregiverManagedClients(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
		pagination  *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: get managed clients",
			args: args{
				ctx:         context.Background(),
				caregiverID: testCaregiverID,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to get managed clients, invalid caregiver id",
			args: args{
				ctx:         context.Background(),
				caregiverID: "invalid",
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := testingDB.GetCaregiverManagedClients(tt.args.ctx, tt.args.caregiverID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCaregiverManagedClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected value but got %v", got)
			}

		})
	}
}

func TestPGInstance_GetCaregiverProfileByCaregiverID(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get caregiver profile by caregiver id",
			args: args{
				ctx:         context.Background(),
				caregiverID: testCaregiverID,
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get caregiver profile by caregiver id",
			args: args{
				ctx:         context.Background(),
				caregiverID: "testCaregiverID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.GetCaregiverProfileByCaregiverID(tt.args.ctx, tt.args.caregiverID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCaregiverProfileByCaregiverID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_ListClientsCaregivers(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list clients caregivers",
			args: args{
				ctx:      context.Background(),
				clientID: clientID,
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to list clients caregivers",
			args: args{
				ctx:      context.Background(),
				clientID: "clientID",
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := testingDB.ListClientsCaregivers(tt.args.ctx, tt.args.clientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListClientsCaregivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_RetrieveFacilityIdentifierByFacilityID(t *testing.T) {
	invalidID := "invalid"
	nonExistentFacilityID := uuid.New().String()
	type args struct {
		ctx        context.Context
		facilityID *string
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
				facilityID: &facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        context.Background(),
				facilityID: &invalidID,
			},
			wantErr: true,
		},
		{
			name: "Sad case: record not found",
			args: args{
				ctx:        context.Background(),
				facilityID: &nonExistentFacilityID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RetrieveFacilityIdentifierByFacilityID(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RetrieveFacilityIdentifierByFacilityID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got %s", err)
			}
		})
	}
}

func TestPGInstance_CheckStaffExists(t *testing.T) {
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
			name: "happy case: staff exists",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignStaff,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: staff does not exist",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignClient,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckStaffExists(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckStaffExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckStaffExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CheckClientExists(t *testing.T) {
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
			name: "happy case: client exists",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: client does not exist",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignStaff,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckClientExists(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckClientExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckClientExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CheckCaregiverExists(t *testing.T) {
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
			name: "happy case: caregiver exists",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: caregiver does not exist",
			args: args{
				ctx:    context.Background(),
				userID: userIDtoAssignStaff,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckCaregiverExists(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckCaregiverExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckCaregiverExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CheckOrganisationExists(t *testing.T) {
	type args struct {
		ctx            context.Context
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: organization exists",
			args: args{
				ctx:            context.Background(),
				organisationID: orgID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: organization does not exists",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: invalid organisation ID",
			args: args{
				ctx:            context.Background(),
				organisationID: "invalid",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CheckOrganisationExists(tt.args.ctx, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckOrganisationExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckOrganisationExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CheckIfProgramNameExists(t *testing.T) {
	type args struct {
		ctx            context.Context
		organisationID string
		programName    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: organisation not associated with the program",
			args: args{
				ctx:            context.Background(),
				organisationID: orgID,
				programName:    gofakeit.Name(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: organization already associated with another program",
			args: args{
				ctx:            context.Background(),
				organisationID: orgID,
				programName:    programName,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: organization does not exists",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
				programName:    gofakeit.Name(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: invalid organisation ID",
			args: args{
				ctx:            context.Background(),
				organisationID: "invalid",
				programName:    gofakeit.Name(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CheckIfProgramNameExists(tt.args.ctx, tt.args.organisationID, tt.args.programName)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfProgramNameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfProgramNameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetUserProfileByUsername(t *testing.T) {

	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get user profile",
			args: args{
				ctx:      context.Background(),
				username: "a_test_user",
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid username",
			args: args{
				ctx:      context.Background(),
				username: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetUserProfileByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetUserProfileByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got %s", err)
			}
		})
	}
}

func TestPGInstance_GetStaffUserPrograms(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "happy case: retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "sad case: invalid user id",
			args: args{
				ctx:    context.Background(),
				userID: "randy",
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetStaffUserPrograms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetStaffUserPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("PGInstance.GetStaffUserPrograms() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}

func TestPGInstance_GetClientUserPrograms(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "happy case: retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "sad case: invalid user id",
			args: args{
				ctx:    context.Background(),
				userID: "randy",
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetClientUserPrograms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetClientUserPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("PGInstance.GetClientUserPrograms() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}

func TestPGInstance_GetProgramFacilities(t *testing.T) {
	type args struct {
		ctx       context.Context
		programID string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "happy case: retrieve program facilities",
			args: args{
				ctx:       context.Background(),
				programID: programID,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "sad case: invalid program ID",
			args: args{
				ctx:       context.Background(),
				programID: "test",
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetProgramFacilities(tt.args.ctx, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("PGInstance.GetClientUserPrograms() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}

func TestPGInstance_ListPrograms(t *testing.T) {
	type args struct {
		ctx            context.Context
		organisationID *string
		pagination     *domain.Pagination
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "Happy Case: list programs",
			args: args{
				ctx:            context.Background(),
				organisationID: &orgID,
				pagination:     &domain.Pagination{Limit: 1, CurrentPage: 1},
			},
			wantErr: false,
		},
		{
			name: "Happy Case: list programs without checking organisation",
			args: args{
				ctx:        context.Background(),
				pagination: &domain.Pagination{Limit: 1, CurrentPage: 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := testingDB.ListPrograms(tt.args.ctx, tt.args.organisationID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected programs not to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got1 == nil {
				t.Errorf("expected pagination not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_CheckIfSuperUserExists(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check if superuser exists, no superuser",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CheckIfSuperUserExists(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CheckIfSuperUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CheckIfSuperUserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_GetCaregiverProfileByUserID(t *testing.T) {
	type args struct {
		ctx            context.Context
		userID         string
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: retrieve caregiver profile by user ID",
			args: args{
				ctx:            context.Background(),
				userID:         userID,
				organisationID: orgID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: user does not exist",
			args: args{
				ctx:            context.Background(),
				userID:         uuid.NewString(),
				organisationID: orgID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.GetCaregiverProfileByUserID(tt.args.ctx, tt.args.userID, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetCaregiverProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetCaregiverProfileByUserID() got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_ListOrganisations(t *testing.T) {
	type args struct {
		ctx        context.Context
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*gorm.Organisation
		want1   *domain.Pagination
		wantErr bool
	}{
		{
			name: "Happy Case: list organisations",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 30,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := testingDB.ListOrganisations(tt.args.ctx, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListOrganisations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_SearchOrganisations(t *testing.T) {
	type args struct {
		ctx             context.Context
		searchParameter string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: search organisations",
			args: args{
				ctx:             context.Background(),
				searchParameter: "test org",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.SearchOrganisation(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SearchOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_ListCommunities(t *testing.T) {
	type args struct {
		ctx            context.Context
		programID      string
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: list communities",
			args: args{
				ctx:            context.Background(),
				programID:      programID,
				organisationID: orgID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: unable to list communities - invalid program id",
			args: args{
				ctx:            context.Background(),
				programID:      "programID",
				organisationID: orgID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case: unable to list communities - invalid org id",
			args: args{
				ctx:            context.Background(),
				programID:      programID,
				organisationID: "orgID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.ListCommunities(tt.args.ctx, tt.args.programID, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ListCommunities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
