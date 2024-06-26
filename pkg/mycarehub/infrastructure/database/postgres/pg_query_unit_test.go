package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
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

func TestMyCareHubDb_ListFacilities(t *testing.T) {
	searchTerm := "ny"
	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to list facilities",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to retrieve facility identifier",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to list facilities" {
				fakeGorm.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to retrieve facility identifier" {
				fakeGorm.MockRetrieveFacilityIdentifiersByFacilityIDFn = func(ctx context.Context, facilityID *string) ([]*gorm.FacilityIdentifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, got1, err := d.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
			if !tt.wantErr && got1 == nil {
				t.Errorf("did nox expect error, got %v", got1)
			}
		})
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
		{
			name: "Sad Case - empty user id",
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

			if tt.name == "Sad Case - Fail to get user pin" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}
			if tt.name == "Sad Case - empty user id" {
				fakeGorm.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*gorm.PINData, error) {
					return nil, fmt.Errorf("failed to get user pin")
				}
			}
			if tt.name == "Sad Case - invalid-flavour" {
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
		ctx context.Context
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
		{
			name: "Sad case - nil context",
			args: args{
				ctx: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockGetCurrentTermsFn = func(ctx context.Context) (*gorm.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - nil context" {
				fakeGorm.MockGetCurrentTermsFn = func(ctx context.Context) (*gorm.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetCurrentTerms(tt.args.ctx)
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
		Username:    uuid.New().String(),
		PhoneNumber: gofakeit.Phone(),
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload1 := &dto.VerifyOTPInput{
		Username:    uuid.New().String(),
		PhoneNumber: gofakeit.Phone(),
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		Username:    uuid.New().String(),
		PhoneNumber: gofakeit.Phone(),
		OTP:         uuid.New().String(),
		Flavour:     "flavour",
	}
	invalidOTPPayload5 := &dto.VerifyOTPInput{
		Username: "otpInput.Username",
		OTP:      "otpInput.OTP",
		Flavour:  "flavour",
	}
	invalidOTPPayload6 := &dto.VerifyOTPInput{
		Username: gofakeit.HipsterParagraph(1, 10, 100, ""),
		OTP:      gofakeit.HipsterParagraph(1, 10, 100, ""),
		Flavour:  "gofakeit.HipsterParagraph(300, 10, 100)",
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

func TestMyCareHubDb_GetClientProfile(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully get client profile by user ID",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client identifier",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get facility by id",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
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
				fakeGorm.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by user ID")
				}
			}
			if tt.name == "Sad Case - Fail to get user profile" {
				fakeGorm.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by user ID")
				}
			}

			if tt.name == "Sad Case - Fail to get client identifier" {
				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeGorm.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by user ID")
				}
			}
			if tt.name == "Sad Case - Fail to get facility by id" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facility by id")
				}
			}

			got, err := d.GetClientProfile(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetStaffProfile(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully get staff profile by user ID",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get staff profile",
			args: args{
				ctx:       ctx,
				userID:    "1234",
				programID: uuid.New().String(),
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
		{
			name: "Sad Case - Unable to retrieve facility",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get staff identifiers",
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

			if tt.name == "Sad Case - Fail to get staff profile" {
				fakeGorm.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile" {
				fakeGorm.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeGorm.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}
			if tt.name == "Sad Case - Unable to retrieve facility" {
				fakeGorm.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
					return &gorm.StaffProfile{
						DefaultFacilityID: "1234",
					}, nil
				}
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facility by id")
				}
			}
			if tt.name == "Sad Case - Unable to get staff identifiers" {
				fakeGorm.MockGetStaffIdentifiersFn = func(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("error")
				}
			}

			got, err := d.GetStaffProfile(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchStaffProfile(t *testing.T) {
	ctx := context.Background()
	programID := uuid.NewString()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		searchTerm string
		programID  *string
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
				ctx:        ctx,
				searchTerm: uuid.New().String(),
				programID:  &programID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to get staff profile(s)",
			args: args{
				ctx:        ctx,
				searchTerm: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get user profile",
			args: args{
				ctx:        ctx,
				searchTerm: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to retrieve facility",
			args: args{
				ctx:        ctx,
				searchTerm: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - fail to get staff profile(s)" {
				fakeGorm.MockSearchStaffProfileFn = func(ctx context.Context, searchParameter string, programID *string) ([]*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			if tt.name == "Sad case - fail to get user profile" {
				fakeGorm.MockSearchStaffProfileFn = func(ctx context.Context, searchParameter string, programID *string) ([]*gorm.StaffProfile, error) {
					return []*gorm.StaffProfile{
						{
							DefaultFacilityID: uuid.New().String(),
							UserID:            uuid.New().String(),
						},
					}, nil
				}
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to retrieve facility" {
				fakeGorm.MockSearchStaffProfileFn = func(ctx context.Context, searchParameter string, programID *string) ([]*gorm.StaffProfile, error) {
					return []*gorm.StaffProfile{
						{
							DefaultFacilityID: uuid.New().String(),
							UserID:            uuid.New().String(),
						},
					}, nil
				}
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					UID := uuid.New().String()
					return &gorm.User{
						UserID: &UID,
					}, nil
				}
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchStaffProfile(tt.args.ctx, tt.args.searchTerm, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_CheckUserHasPin(t *testing.T) {
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
			name: "happy case - check user has pin",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing user ID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: could not check user has pin",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
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
				fakeGorm.MockCheckUserHasPinFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckUserHasPin(tt.args.ctx, tt.args.userID)
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
		ctx             context.Context
		userID, flavour string
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
				ctx:     ctx,
				userID:  gofakeit.UUID(),
				flavour: string(feedlib.FlavourPro),
			},
		},
		{
			name: "Happy case - no resposes",
			args: args{
				ctx:     ctx,
				userID:  gofakeit.UUID(),
				flavour: string(feedlib.FlavourPro),
			},
		},
		{
			name: "invalid: could not find security questions",
			args: args{
				ctx:     ctx,
				userID:  gofakeit.UUID(),
				flavour: string(feedlib.FlavourPro),
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
				fakeGorm.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*gorm.SecurityQuestionResponse, error) {
					return []*gorm.SecurityQuestionResponse{}, nil
				}
			}

			if tt.name == "invalid: could not find security questions" {
				fakeGorm.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*gorm.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get user security questions")
				}
			}

			got, err := d.GetUserSecurityQuestionsResponses(tt.args.ctx, tt.args.userID, tt.args.flavour)
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

func TestMyCareHubDb_GetClientHealthDiaryQuote(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully get client health diary quote",
			args: args{
				ctx:   ctx,
				limit: 1,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client health diary quote",
			args: args{
				ctx:   ctx,
				limit: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad Case - Fail to get client health diary quote" {
				fakeGorm.MockGetClientHealthDiaryQuoteFn = func(ctx context.Context, limit int) ([]*gorm.ClientHealthDiaryQuote, error) {
					return nil, fmt.Errorf("failed to get client health diary quote")
				}
			}
			got, err := d.GetClientHealthDiaryQuote(tt.args.ctx, tt.args.limit)
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

func TestMyCareHubDb_GetClientHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		clientID string
		moodType enums.Mood
		shared   bool
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
				moodType: enums.MoodSad,
				shared:   true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client health diary entries",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
				moodType: enums.MoodSad,
				shared:   true,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
				moodType: enums.MoodSad,
				shared:   false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client health diary entries" {
				fakeGorm.MockGetClientHealthDiaryEntriesFn = func(ctx context.Context, params map[string]interface{}) ([]*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - Fail to get client profile" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			got, err := d.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.clientID, &tt.args.moodType, &tt.args.shared)
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

func TestMyCareHubDb_GetPendingServiceRequestsCount(t *testing.T) {
	ctx := context.Background()

	facilityID := uuid.New().String()

	type args struct {
		ctx        context.Context
		facilityID string
		programID  string
	}
	tests := []struct {
		name    string
		args    args
		want    *int64
		wantErr bool
	}{
		{
			name: "Happy case: get pending service request count",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
				programID:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get pending service request count",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
				programID:  gofakeit.UUID(),
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
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to get pending service request count" {
				fakeGorm.MockGetClientPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty facility ID" {
				fakeGorm.MockGetClientPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - fail to get staff service requests count" {
				fakeGorm.MockGetStaffPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetPendingServiceRequestsCount(tt.args.ctx, tt.args.facilityID, tt.args.programID)
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
		{
			name: "Sad Case - Unable to retrieve facility",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get client identifiers",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get program by ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get organisation by ID",
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
			if tt.name == "Sad Case - Unable to retrieve facility" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad Case - Unable to get client identifiers" {
				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad Case - Unable to get program by ID" {
				fakeGorm.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*gorm.Program, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad Case - Unable to get organisation by ID" {
				fakeGorm.MockGetOrganisationFn = func(ctx context.Context, id string) (*gorm.Organisation, error) {
					return nil, fmt.Errorf("error")
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
	var requesttype = enums.ServiceRequestTypeRedFlag.String()
	facilityID := uuid.New().String()

	type args struct {
		ctx           context.Context
		requestType   *string
		requestStatus *string
		facilityID    string
		programID     string
		flavour       feedlib.Flavour
		pagination    *domain.Pagination
	}

	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get service requests - Consumer",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourConsumer,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get service requests - Pro",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourPro,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid Flavour",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       "invalid flavour",
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests - Consumer",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourConsumer,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile by client ID",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourConsumer,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by staff ID",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourConsumer,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff service requests",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourPro,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourPro,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by staff ID",
			args: args{
				ctx:           context.Background(),
				requestType:   &requesttype,
				requestStatus: new(string),
				facilityID:    facilityID,
				programID:     gofakeit.UUID(),
				flavour:       feedlib.FlavourPro,
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get service requests - Consumer" {
				fakeGorm.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get service requests by type")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by client ID" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by client ID")
				}
			}

			if tt.name == "Happy Case - Successfully get service requests - Consumer" {
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
						ResolvedByID:   &UUID,
						Meta:           "{}",
					},
				}
				fakeGorm.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
					return serviceRequests, &domain.Pagination{
						CurrentPage: 1,
						Limit:       10,
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to get staff service requests" {
				fakeGorm.MockGetStaffServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string, pagination *domain.Pagination) ([]*gorm.StaffServiceRequest, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get staff service request")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile" {
				fakeGorm.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by staff ID")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by staff ID" {
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by staff ID")
				}
			}

			got, _, err := d.GetServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID, tt.args.programID, tt.args.flavour, tt.args.pagination)
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
	id := gofakeit.UUID()

	type args struct {
		ctx          context.Context
		lastSyncTime time.Time
		client       *domain.ClientProfile
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
				lastSyncTime: time.Now(),
				client: &domain.ClientProfile{
					ID:     &id,
					UserID: id,
					User: &domain.User{
						Name: gofakeit.Name(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get recent entries",
			args: args{
				ctx:          ctx,
				lastSyncTime: time.Now(),
				client: &domain.ClientProfile{
					ID:     &id,
					UserID: id,
					User: &domain.User{
						Name: gofakeit.Name(),
					},
				},
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

			got, err := d.GetRecentHealthDiaryEntries(tt.args.ctx, tt.args.lastSyncTime, tt.args.client)
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

func TestMyCareHubDb_GetClientIdentifiers(t *testing.T) {
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
			fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
				return nil, fmt.Errorf("cannot get client identifier")
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := d.GetClientIdentifiers(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientIdentifiers() error = %v, wantErr %v", err, tt.wantErr)
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
		{
			name: "Sad case: Fail to retrieve facility by mfl code",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      123,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - Fail to get client profile by client ID",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      1234,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - Fail to get user profile by user ID",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      1234,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockGetServiceRequestsForKenyaEMRFn = func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: Fail to retrieve facility by mfl code" {
				fakeGorm.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility by mflcode")
				}
			}

			if tt.name == "Sad case - Fail to get client profile by client ID" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad case - Fail to get user profile by user ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
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

func TestMyCareHubDb_CheckIdentifierExists(t *testing.T) {
	type args struct {
		ctx             context.Context
		identifierType  enums.UserIdentifierType
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

func TestMyCareHubDb_CheckFacilityExistsByIdentifier(t *testing.T) {
	type args struct {
		ctx        context.Context
		identifier *dto.FacilityIdentifierInput
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
				ctx: context.Background(),
				identifier: &dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "490843984938948",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.CheckFacilityExistsByIdentifier(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckFacilityExistsByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckFacilityExistsByIdentifier() = %v, want %v", got, tt.want)
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

func TestMyCareHubDb_ListNotifications(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	id := gofakeit.UUID()

	type args struct {
		ctx        context.Context
		params     *domain.Notification
		filters    []*firebasetools.FilterParam
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Notification
		want1   *domain.Pagination
		wantErr bool
	}{
		{
			name: "happy case: success listing notifications",
			args: args{
				ctx: context.Background(),
				params: &domain.Notification{
					UserID:  &id,
					Flavour: feedlib.FlavourConsumer,
				},
				pagination: &domain.Pagination{
					Limit:       2,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error listing notifications",
			args: args{
				ctx: context.Background(),
				params: &domain.Notification{
					UserID:  &id,
					Flavour: feedlib.FlavourPro,
				},
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
			if tt.name == "sad case: error listing notifications" {
				fakeGorm.MockListNotificationsFn = func(ctx context.Context, params *gorm.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Notification, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("error listing notifications")
				}
			}

			got, got1, err := d.ListNotifications(tt.args.ctx, tt.args.params, tt.args.filters, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListNotifications() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_ListAvailableNotificationTypes(t *testing.T) {
	id := gofakeit.UUID()

	type args struct {
		ctx    context.Context
		params *domain.Notification
	}
	tests := []struct {
		name    string
		args    args
		want    []enums.NotificationType
		wantErr bool
	}{
		{
			name: "happy case: available notification types for a user",
			args: args{
				ctx: context.Background(),
				params: &domain.Notification{
					UserID: &id,
				},
			},
			want:    []enums.NotificationType{enums.NotificationTypeAppointment},
			wantErr: false,
		},
		{
			name: "sad case: fail to fetch available notification types",
			args: args{
				ctx: context.Background(),
				params: &domain.Notification{
					UserID: &id,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: fail to fetch available notification types" {
				fakeGorm.MockListAvailableNotificationTypesFn = func(ctx context.Context, params *gorm.Notification) ([]enums.NotificationType, error) {
					return []enums.NotificationType{}, fmt.Errorf("failed to fetch notification appointments")
				}
			}
			got, err := d.ListAvailableNotificationTypes(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListAvailableNotificationTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MyCareHubDb.ListAvailableNotificationTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetProgramClientProfileByIdentifier(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                              context.Context
		programID, identifierType, value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profile by CCC number",
			args: args{
				ctx:            ctx,
				programID:      gofakeit.UUID(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "123456",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:            ctx,
				programID:      gofakeit.UUID(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "123456",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client ccc identifier",
			args: args{
				ctx:            ctx,
				programID:      gofakeit.UUID(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "123456",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile by CCC number",
			args: args{
				ctx:            ctx,
				programID:      gofakeit.UUID(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "123456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client profile by CCC number" {
				fakeGorm.MockGetProgramClientProfileByIdentifierFn = func(ctx context.Context, programID string, identifierType string, value string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by CCC number")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Sad Case - Fail to get client ccc identifier" {
				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc identifier")
				}
			}

			got, err := d.GetProgramClientProfileByIdentifier(tt.args.ctx, tt.args.programID, tt.args.identifierType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetProgramClientProfileByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_SearchClientProfiles(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx             context.Context
		searchParameter string
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
				ctx:             ctx,
				searchParameter: "345678",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:             ctx,
				searchParameter: "111111",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client identifier",
			args: args{
				ctx:             ctx,
				searchParameter: "111111",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile",
			args: args{
				ctx:             ctx,
				searchParameter: "345",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to retrieve facility",
			args: args{
				ctx:             ctx,
				searchParameter: "345",
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
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					id := uuid.New().String()
					return &gorm.User{
						UserID: &id,
					}, nil
				}

				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return []*gorm.Identifier{{
						ID: clientID,
					},
					}, nil
				}

				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("failed to get client identifier")
				}
			}

			if tt.name == "Sad Case - Fail to retrieve facility" {
				fakeGorm.MockSearchClientProfileFn = func(ctx context.Context, searchParameter string) ([]*gorm.Client, error) {
					return []*gorm.Client{
						{
							FacilityID: uuid.New().String(),
						},
					}, nil
				}
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					id := uuid.New().String()
					return &gorm.User{
						UserID: &id,
					}, nil
				}

				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile" {
				fakeGorm.MockSearchClientProfileFn = func(ctx context.Context, searchParameter string) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("failed to search client profile")
				}
			}

			got, err := d.SearchClientProfile(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchClientProfile() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_GetClientServiceRequestByID(t *testing.T) {
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
				fakeGorm.MockGetClientServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service request by ID")
				}
			}

			got, err := d.GetClientServiceRequestByID(tt.args.ctx, tt.args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientServiceRequestByID() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "Happy case: get staff profile",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get staff identifier",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get staff profile",
			args: args{
				ctx:     ctx,
				staffID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to get staff identifier" {
				fakeGorm.MockGetStaffIdentifiersFn = func(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get staff identifier" {
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
			name: "Sad case:  failed to get facility",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: true,
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
			name: "Sad case:  failed to get client ccc identifier by client ID",
			args: args{
				ctx:          context.Background(),
				lastSyncTime: time.Now(),
				mflCode:      "1234567890",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case:  failed to get facility" {
				fakeGorm.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}

			if tt.name == "Sad case:  failed to get appointment service request" {
				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case:  invalid service request meta" {
				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					requestID := uuid.New().String()
					currentTime := time.Now()
					staffID := uuid.New().String()
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
							Meta:           `}{[]{`,
						},
					}, nil
				}
			}

			if tt.name == "Sad case:  failed to get appointment by id" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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

				fakeGorm.MockGetAppointmentFn = func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error) {
					return nil, fmt.Errorf("an error occurred")

				}
			}

			if tt.name == "Sad case:  failed to convert to suggested time" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    uuid.New().String(),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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

				fakeGorm.MockGetAppointmentFn = func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error) {
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

			if tt.name == "Sad case:  failed to get in progress by staff" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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
							ResolvedByID:   &UUID,
							InProgressByID: &UUID,
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
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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
							InProgressByID: nil,
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
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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

				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad case:  failed to get client ccc identifier by client ID" {
				now := time.Now()
				UUID := uuid.New().String()
				meta := map[string]interface{}{
					"appointmentID":     uuid.New().String(),
					"externalID":        uuid.New().String(),
					"appointmentReason": "reason",
					"provider":          "provider",
					"rescheduleTime":    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				}

				bs, err := json.Marshal(meta)
				if err != nil {
					t.Errorf("failed to marshal meta: %v", err)
				}

				fakeGorm.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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

				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc identifier")
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

func TestMyCareHubDb_GetClientServiceRequests(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx        context.Context
		toolType   string
		status     string
		clientID   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case:  get system generated client service requests",
			args: args{
				ctx:        context.Background(),
				toolType:   "SCREENING_TOOLS_RED_FLAG",
				status:     enums.ServiceRequestStatusPending.String(),
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.GetClientServiceRequests(tt.args.ctx, tt.args.toolType, tt.args.status, tt.args.clientID, tt.args.facilityID)
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

func TestMyCareHubDb_GetAppointment(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx    context.Context
		params domain.Appointment
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Appointment
		wantErr bool
	}{
		{
			name: "Happy case: retrieve an appointment",
			args: args{
				ctx: context.Background(),
				params: domain.Appointment{
					ClientID: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid param to retrieve an appointment",
			args: args{
				ctx: context.Background(),
				params: domain.Appointment{
					ClientID: "-",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: invalid param to retrieve an appointment" {
				fakeGorm.MockGetAppointmentFn = func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error) {
					return nil, fmt.Errorf("cannot retrieve appointment")
				}
			}

			got, err := d.GetAppointment(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAppointment() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_GetSharedHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
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
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid facility",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: gofakeit.HipsterSentence(44),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get client profile",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - invalid facility" {
				fakeGorm.MockGetSharedHealthDiaryEntriesFn = func(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get shared health diary entries")
				}
			}
			if tt.name == "Sad case - empty facility" {
				fakeGorm.MockGetSharedHealthDiaryEntriesFn = func(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get shared health diary entries")
				}
			}
			if tt.name == "Sad case - unable to get client profile" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			got, err := d.GetSharedHealthDiaryEntries(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSharedHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_GetClientScreeningToolServiceRequestByToolType(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx        context.Context
		clientID   string
		questionID string
		status     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case:  get client screening tool service request",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.New().String(),
				questionID: uuid.New().String(),
				status:     "PENDING",
			},
		},
		{
			name: "Sad case:  failed to get client screening tool service request",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.New().String(),
				questionID: uuid.New().String(),
				status:     "PENDING",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  failed to get client screening tool service request" {
				fakeGorm.MockGetClientScreeningToolServiceRequestByToolTypeFn = func(ctx context.Context, clientID string, toolType string, status string) (*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("failed to get client screening tool service request")
				}
			}
			got, err := d.GetClientScreeningToolServiceRequestByToolType(tt.args.ctx, tt.args.clientID, tt.args.questionID, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientScreeningToolServiceRequestByToolType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetUserSurveyForms(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case:  get user survey forms",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"user_id": uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case:  unable to get user survey forms",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"user_id": uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  unable to get user survey forms" {
				fakeGorm.MockGetUserSurveyFormsFn = func(ctx context.Context, params map[string]interface{}) ([]*gorm.UserSurvey, error) {
					return nil, fmt.Errorf("failed to get user survey forms")
				}
			}

			got, err := d.GetUserSurveyForms(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserSurveyForms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfStaffHasUnresolvedServiceRequests(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
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
			name: "Happy case:  check if staff has no unresolved service requests",
			args: args{
				ctx:                context.Background(),
				staffID:            uuid.New().String(),
				serviceRequestType: string(enums.ServiceRequestTypePinReset),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckIfStaffHasUnresolvedServiceRequests(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfStaffHasUnresolvedServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfStaffHasUnresolvedServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetFacilityStaffs(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility staff",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
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
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad case: invalid facility id" {
				fakeGorm.MockGetFacilityStaffsFn = func(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("cannot retrieve facilities")
				}
			}

			got, err := d.GetFacilityStaffs(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetFacilityStaffs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetNotification(t *testing.T) {

	type args struct {
		ctx            context.Context
		notificationID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Notification
		wantErr bool
	}{
		{
			name: "happy case: retrieve a facility",
			args: args{
				ctx:            context.Background(),
				notificationID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid facility id",
			args: args{
				ctx:            context.Background(),
				notificationID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: invalid facility id" {
				fakeGorm.MockGetNotificationFn = func(ctx context.Context, notificationID string) (*gorm.Notification, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}

			got, err := d.GetNotification(tt.args.ctx, tt.args.notificationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetClientsByFilterParams(t *testing.T) {
	facilityID := gofakeit.UUID()
	type args struct {
		ctx          context.Context
		facilityID   *string
		filterParams *dto.ClientFilterParamsInput
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy case: retrieve clients",
			args: args{
				ctx:          context.Background(),
				facilityID:   &facilityID,
				filterParams: &dto.ClientFilterParamsInput{},
			},
			wantErr: false,
		},
		{
			name: "Happy case: retrieve clients with filters",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
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
			name: "Sad case: failed to filter clients",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get user profile by user id",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to filter clients" {
				fakeGorm.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("cannot filter clients")
				}
			}

			if tt.name == "Sad case: failed to get user profile by user id" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("cannot get user profile")
				}
			}
			got, err := d.GetClientsByFilterParams(tt.args.ctx, tt.args.facilityID, tt.args.filterParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientsByFilterParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchStaffServiceRequests(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx             context.Context
		searchParameter string
		requestType     string
		facilityID      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: search staff service requests",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
				requestType:     "PIN_RESET",
				facilityID:      uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search staff service requests",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get staff profile by staff ID",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to search staff service requests" {
				fakeGorm.MockSearchStaffServiceRequestsFn = func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get staff profile by staff ID" {
				fakeGorm.MockSearchStaffServiceRequestsFn = func(ctx context.Context, searchParameter, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
					ID := uuid.New().String()
					return []*gorm.StaffServiceRequest{
						{
							ID:          &ID,
							Active:      true,
							RequestType: "TEST",
							Request:     "STAFF_PIN_RESET",
							Status:      "PENDING",
							ResolvedAt:  &time.Time{},
							StaffID:     ID,
						},
					}, nil
				}
				fakeGorm.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchStaffServiceRequests(tt.args.ctx, tt.args.searchParameter, tt.args.requestType, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchStaffServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchClientServiceRequests(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx             context.Context
		searchParameter string
		requestType     string
		facilityID      string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: search client service requests",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
				requestType:     "RED_FLAG",
				facilityID:      uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search client service requests",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
				requestType:     "RED_FLAG",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile by client id",
			args: args{
				ctx:             ctx,
				searchParameter: "PENDING",
				requestType:     "RED_FLAG",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to search client service requests" {
				fakeGorm.MockSearchClientServiceRequestsFn = func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get client profile by client id" {
				fakeGorm.MockSearchClientServiceRequestsFn = func(ctx context.Context, searchParameter, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
					ID := uuid.New().String()
					return []*gorm.ClientServiceRequest{
						{
							ID:          &ID,
							Active:      true,
							RequestType: "RED_FLAG",
							Request:     "test",
							Status:      "PENDING",
							ClientID:    ID,
						},
					}, nil
				}
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchClientServiceRequests(tt.args.ctx, tt.args.searchParameter, tt.args.requestType, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchClientServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_ReturnClientsServiceRequests(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	ID := uuid.New().String()
	meta := map[string]interface{}{
		"test": "test",
	}
	validjsonMeta, err := json.Marshal(meta)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}

	type args struct {
		ctx                   context.Context
		clientServiceRequests []*gorm.ClientServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: return clients service requests",
			args: args{
				ctx: ctx,
				clientServiceRequests: []*gorm.ClientServiceRequest{
					{
						ID:           &ID,
						Active:       true,
						RequestType:  "RED_FLAG",
						Request:      "test",
						Status:       "PENDING",
						Meta:         string(validjsonMeta),
						ResolvedByID: &ID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get client profile by client id",
			args: args{
				ctx: ctx,
				clientServiceRequests: []*gorm.ClientServiceRequest{
					{
						ID:          &ID,
						Active:      true,
						RequestType: "RED_FLAG",
						Request:     "test",
						Status:      "PENDING",
						Meta:        string(validjsonMeta),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get client profile by client id" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.ReturnClientsServiceRequests(tt.args.ctx, tt.args.clientServiceRequests)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ReturnClientsServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_ReturnStaffServiceRequests(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	ID := uuid.New().String()
	meta := map[string]interface{}{
		"test": "test",
	}
	validjsonMeta, err := json.Marshal(meta)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}

	invalidjsonMeta, err := json.Marshal("invalid json")
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}

	type args struct {
		ctx                  context.Context
		staffServiceRequests []*gorm.StaffServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: return staff service requests",
			args: args{
				ctx: ctx,
				staffServiceRequests: []*gorm.StaffServiceRequest{
					{
						ID:           &ID,
						Active:       true,
						RequestType:  "RED_FLAG",
						Request:      "test",
						Status:       "PENDING",
						Meta:         string(validjsonMeta),
						ResolvedByID: &ID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get staff profile by staff ID",
			args: args{
				ctx: ctx,
				staffServiceRequests: []*gorm.StaffServiceRequest{
					{
						ID:          &ID,
						Active:      true,
						RequestType: "RED_FLAG",
						Request:     "test",
						Status:      "PENDING",
						Meta:        string(invalidjsonMeta),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by staff ID",
			args: args{
				ctx: ctx,
				staffServiceRequests: []*gorm.StaffServiceRequest{
					{
						ID:           &ID,
						Active:       true,
						RequestType:  "RED_FLAG",
						Request:      "test",
						Status:       "PENDING",
						ResolvedByID: &ID,
						Meta:         string(validjsonMeta),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get staff profile by staff ID" {
				fakeGorm.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile by staff ID" {
				fakeGorm.MockGetStaffServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string, pagination *domain.Pagination) ([]*gorm.StaffServiceRequest, *domain.Pagination, error) {
					return []*gorm.StaffServiceRequest{
							{
								ID:           &ID,
								Active:       true,
								ResolvedByID: &ID,
							},
						}, &domain.Pagination{
							CurrentPage: 1,
							Limit:       10,
						}, fmt.Errorf("an error occurred")
				}
				fakeGorm.MockGetUserProfileByStaffIDFn = func(ctx context.Context, staffID string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.ReturnStaffServiceRequests(tt.args.ctx, tt.args.staffServiceRequests)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ReturnStaffServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetScreeningToolByID(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx    context.Context
		toolID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool by ID",
			args: args{
				ctx:    context.Background(),
				toolID: uuid.NewString(),
			},
		},
		{
			name: "Sad case: failed to get screenig tool by id",
			args: args{
				ctx:    context.Background(),
				toolID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get questionnaire by id",
			args: args{
				ctx:    context.Background(),
				toolID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get questions by questionnaire id",
			args: args{
				ctx:    context.Background(),
				toolID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get question input choice by question id",
			args: args{
				ctx:    context.Background(),
				toolID: uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to get screenig tool by id" {
				fakeGorm.MockGetScreeningToolByIDFn = func(ctx context.Context, toolID string) (*gorm.ScreeningTool, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get questionnaire by id" {
				fakeGorm.MockGetQuestionnaireByIDFn = func(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get questions by questionnaire id" {
				fakeGorm.MockGetQuestionsByQuestionnaireIDFn = func(ctx context.Context, questionnaireID string) ([]*gorm.Question, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get question input choice by question id" {
				fakeGorm.MockGetQuestionInputChoicesByQuestionIDFn = func(ctx context.Context, questionID string) ([]*gorm.QuestionInputChoice, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetScreeningToolByID(tt.args.ctx, tt.args.toolID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetFacilityRespondedScreeningTools(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case: return facility responded screening tools",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
		},
		{
			name: "Sad case: unable to get questionnaire by ID",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility responded screening tools",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get facility responded screening tools" {
				fakeGorm.MockGetFacilityRespondedScreeningToolsFn = func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*gorm.ScreeningTool, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get questionnaire by ID" {
				fakeGorm.MockGetQuestionnaireByIDFn = func(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, _, err := d.GetFacilityRespondedScreeningTools(tt.args.ctx, tt.args.facilityID, tt.args.programID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetFacilityRespondedScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_ListSurveyRespondents(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		pagination *domain.Pagination
		params     *domain.UserSurvey
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.SurveyRespondent
		want1   *domain.Pagination
		wantErr bool
	}{
		{
			name: "Happy case: return survey respondents",
			args: args{
				ctx: ctx,
				params: &domain.UserSurvey{
					HasSubmitted: true,
					FormID:       gofakeit.UUID(),
					ProjectID:    1,
					ProgramID:    gofakeit.UUID(),
				},
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get survey respondents",
			args: args{
				ctx: ctx,
				params: &domain.UserSurvey{
					HasSubmitted: true,
					FormID:       gofakeit.UUID(),
					ProjectID:    1,
					ProgramID:    gofakeit.UUID(),
				},
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile",
			args: args{
				ctx: ctx,
				params: &domain.UserSurvey{
					HasSubmitted: true,
					FormID:       gofakeit.UUID(),
					ProjectID:    1,
					ProgramID:    gofakeit.UUID(),
				},
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get survey respondents" {
				fakeGorm.MockListSurveyRespondentsFn = func(ctx context.Context, params *gorm.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*gorm.UserSurvey, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, _, err := d.ListSurveyRespondents(tt.args.ctx, tt.args.params, tt.args.facilityID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListSurveyRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetScreeningToolRespondents(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx             context.Context
		facilityID      string
		programID       string
		screeningToolID string
		searchTerm      string
		paginationInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: return respondents",
			args: args{
				ctx:             ctx,
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get screening tool service requests of respondent",
			args: args{
				ctx:             ctx,
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile by user id",
			args: args{
				ctx:             ctx,
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get screening tool service requests of respondent" {
				fakeGorm.MockGetScreeningToolServiceRequestOfRespondentsFn = func(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get client profile by user id" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, _, err := d.GetScreeningToolRespondents(tt.args.ctx, tt.args.facilityID, tt.args.programID, tt.args.screeningToolID, tt.args.searchTerm, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetScreeningToolResponseByID(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case: return screening tool response",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get screening tool response",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get screening tool question responses",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get screening tool",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get screening tool response" {
				fakeGorm.MockGetScreeningToolResponseByIDFn = func(ctx context.Context, id string) (*gorm.ScreeningToolResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get screening tool question responses" {
				fakeGorm.MockGetScreeningToolQuestionResponsesByResponseIDFn = func(ctx context.Context, id string) ([]*gorm.ScreeningToolQuestionResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get screening tool" {
				fakeGorm.MockGetQuestionsByQuestionnaireIDFn = func(ctx context.Context, questionnaireID string) ([]*gorm.Question, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetScreeningToolResponseByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolResponseByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetSurveysWithServiceRequests(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		facilityID string
		programID  string
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.SurveysWithServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case: return surveys with service requests",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				programID:  uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable return surveys with service requests",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				programID:  uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable return surveys with service requests" {
				fakeGorm.MockGetSurveysWithServiceRequestsFn = func(ctx context.Context, facilityID, programID string) ([]*gorm.UserSurvey, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetSurveysWithServiceRequests(tt.args.ctx, tt.args.facilityID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSurveysWithServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}

		})
	}
}

func TestMyCareHubDb_GetSurveyServiceRequestUser(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
			name: "Happy case: successfully get users with survey service request",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				projectID:  1,
				formID:     "test",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get clients survey service requests",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				projectID:  1,
				formID:     "test",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get clients profile by client ID",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				projectID:  1,
				formID:     "test",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to convert form ID to string",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				projectID:  1,
				formID:     "test",
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
			if tt.name == "Sad case: unable to get clients survey service requests" {
				fakeGorm.MockGetClientsSurveyServiceRequestFn = func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get clients profile by client ID" {
				fakeGorm.MockGetClientsSurveyServiceRequestFn = func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
					ID := uuid.New().String()
					return []*gorm.ClientServiceRequest{
							{
								ID:         &ID,
								Active:     true,
								FacilityID: uuid.New().String(),
								ClientID:   uuid.New().String(),
							},
						}, &domain.Pagination{
							Limit:       5,
							CurrentPage: 1,
						}, nil
				}

				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, _, err := d.GetSurveyServiceRequestUser(tt.args.ctx, tt.args.facilityID, tt.args.projectID, tt.args.formID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetSurveyServiceRequestUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetStaffFacilities(t *testing.T) {
	ctx := context.Background()
	staffID := uuid.New().String()
	facilityID := uuid.New().String()
	type args struct {
		ctx        context.Context
		input      dto.StaffFacilityInput
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get staff facilities",
			args: args{
				ctx: ctx,
				input: dto.StaffFacilityInput{
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
			name: "sad case: failed to get staff facilities",
			args: args{
				ctx: ctx,
				input: dto.StaffFacilityInput{
					StaffID:    &staffID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to retrieve facility",
			args: args{
				ctx: ctx,
				input: dto.StaffFacilityInput{
					StaffID:    &staffID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get staff notification count",
			args: args{
				ctx: ctx,
				input: dto.StaffFacilityInput{
					StaffID:    &staffID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get clients service request count",
			args: args{
				ctx: ctx,
				input: dto.StaffFacilityInput{
					StaffID:    &staffID,
					FacilityID: &facilityID,
				},
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
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "sad case: failed to get staff facilities" {
				fakeGorm.MockGetStaffFacilitiesFn = func(ctx context.Context, staffFacility gorm.StaffFacilities, pagination *domain.Pagination) ([]*gorm.StaffFacilities, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: failed to retrieve facility" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}
			if tt.name == "sad case: unable to get staff notification count" {
				fakeGorm.MockGetNotificationsCountFn = func(ctx context.Context, notification gorm.Notification) (int, error) {
					return 0, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get clients service request count" {
				fakeGorm.MockGetClientPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, _, err := d.GetStaffFacilities(tt.args.ctx, tt.args.input, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_GetClientFacilities(t *testing.T) {
	ctx := context.Background()
	clientID := uuid.New().String()
	facilityID := uuid.New().String()
	type args struct {
		ctx        context.Context
		input      dto.ClientFacilityInput
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get client facilities",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
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
			name: "sad case: failed to get client facilities",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to retrieve facility",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get client profile",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get notification count",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get surveys count",
			args: args{
				ctx: ctx,
				input: dto.ClientFacilityInput{
					ClientID:   &clientID,
					FacilityID: &facilityID,
				},
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

			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get client profile" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: failed to get client facilities" {
				fakeGorm.MockGetClientFacilitiesFn = func(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: failed to retrieve facility" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}
			if tt.name == "sad case: failed to get notification count" {
				fakeGorm.MockGetNotificationsCountFn = func(ctx context.Context, notification gorm.Notification) (int, error) {
					return 0, fmt.Errorf("failed to get notification count")
				}
			}
			if tt.name == "sad case: failed to get surveys count" {
				fakeGorm.MockGetClientsSurveyCountFn = func(ctx context.Context, userID string) (int, error) {
					return 0, fmt.Errorf("failed to get surveys count")
				}
			}

			got, _, err := d.GetClientFacilities(tt.args.ctx, tt.args.input, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_SearchCaregiverUser(t *testing.T) {
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
			name: "happy case: search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user ID",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile by user ID",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to search caregiver user" {
				fakeGorm.MockSearchCaregiverUserFn = func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get client profile by user ID" {
				fakeGorm.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchCaregiverUser(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchCaregiverUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_SearchPlatformCaregivers(t *testing.T) {
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
			name: "happy case: search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user ID",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile by user ID",
			args: args{
				ctx:             context.Background(),
				searchParameter: "CG001",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to search caregiver user" {
				fakeGorm.MockSearchPlatformCaregiversFn = func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user ID" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get client profile by user ID" {
				fakeGorm.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.SearchPlatformCaregivers(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchPlatformCaregivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_GetCaregiverManagedClients(t *testing.T) {
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
			name: "Happy Case: get caregiver's clients",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to get caregiver's clients",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get client profile by client id",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get client facilities",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get notification count",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get survey count",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
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
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case: failed to get caregiver's clients" {
				fakeGorm.MockGetCaregiverManagedClientsFn = func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get caregivers's client")
				}
			}

			if tt.name == "Sad Case: failed to get client profile by client id" {
				fakeGorm.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profile by client ID")
				}
			}

			if tt.name == "Sad Case: failed to get client facilities" {
				fakeGorm.MockGetClientFacilitiesFn = func(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get client facilities")
				}
			}

			if tt.name == "Sad Case: failed to get notification count" {
				fakeGorm.MockGetNotificationsCountFn = func(ctx context.Context, notification gorm.Notification) (int, error) {
					return 0, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case: failed to get survey count" {
				fakeGorm.MockGetClientsSurveyCountFn = func(ctx context.Context, userID string) (int, error) {
					return 0, fmt.Errorf("an error occurred")
				}
			}
			got, _, err := d.GetCaregiverManagedClients(tt.args.ctx, tt.args.caregiverID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCaregiverManagedClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_ListClientsCaregivers(t *testing.T) {
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
				clientID: "CL001",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
					Count:       30,
					TotalPages:  50,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to list clients caregivers",
			args: args{
				ctx:      context.Background(),
				clientID: "CL001",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
					Count:       30,
					TotalPages:  50,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get caregiver profile by caregiver ID",
			args: args{
				ctx:      context.Background(),
				clientID: "CL001",
				pagination: &domain.Pagination{
					Limit:       5,
					CurrentPage: 1,
					Count:       30,
					TotalPages:  50,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: unable to list clients caregivers" {
				fakeGorm.MockListClientsCaregiversFn = func(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get caregiver profile by caregiver ID" {
				fakeGorm.MockGetCaregiverProfileByCaregiverIDFn = func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, _, err := d.ListClientsCaregivers(tt.args.ctx, tt.args.clientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListClientsCaregivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_CheckOrganisationExists(t *testing.T) {
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
			name: "happy case: check if organisation exists",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.NewString(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: failed to check if organisation exists",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.NewString(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to check if organisation exists" {
				fakeGorm.MockCheckOrganisationExistsFn = func(ctx context.Context, organisationID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckOrganisationExists(tt.args.ctx, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckOrganisationExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckOrganisationExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfProgramNameExists(t *testing.T) {
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
			name: "happy case: check if organisation has a program",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.NewString(),
				programName:    gofakeit.Name(),
			},
			want:    false, // organization is unique for all programs
			wantErr: false,
		},
		{
			name: "sad case: failed to check if organisation has a program",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.NewString(),
				programName:    gofakeit.Name(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to check if organisation has a program" {
				fakeGorm.MockCheckIfProgramNameExistsFn = func(ctx context.Context, organisationID string, programName string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckIfProgramNameExists(tt.args.ctx, tt.args.organisationID, tt.args.programName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfProgramNameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfProgramNameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_ListOrganisations(t *testing.T) {
	type args struct {
		ctx        context.Context
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list organisations",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to list programs",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to list organisations",
			args: args{
				ctx: context.Background(),
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
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to list programs" {
				fakeGorm.MockListProgramsFn = func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: unable to list organisations" {
				fakeGorm.MockListOrganisationsFn = func(ctx context.Context, pagination *domain.Pagination) ([]*gorm.Organisation, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}

			_, _, err := d.ListOrganisations(tt.args.ctx, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListOrganisations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetStaffUserPrograms(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Program
		wantErr bool
	}{
		{
			name: "happy case: retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid failed to retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: invalid failed to retrieve user programs" {
				fakeGorm.MockGetStaffUserProgramsFn = func(ctx context.Context, userID string) ([]*gorm.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetStaffUserPrograms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffUserPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_GetClientUserPrograms(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Program
		wantErr bool
	}{
		{
			name: "happy case: retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid failed to retrieve user programs",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "sad case: invalid failed to retrieve user programs" {
				fakeGorm.MockGetClientUserProgramsFn = func(ctx context.Context, userID string) ([]*gorm.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetClientUserPrograms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientUserPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestMyCareHubDb_GetProgramFacilities(t *testing.T) {
	type args struct {
		ctx       context.Context
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get program facilities",
			args: args{
				ctx:       context.Background(),
				programID: uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get program facilities",
			args: args{
				ctx:       context.Background(),
				programID: uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get program facilities" {
				fakeGorm.MockGetProgramFacilitiesFn = func(ctx context.Context, programID string) ([]*gorm.ProgramFacility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.GetProgramFacilities(tt.args.ctx, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}
func TestMyCareHubDb_GetProgramByID(t *testing.T) {
	fakeGorm := gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx       context.Context
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: retrieve program by id",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to retrieve program by id",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: unable to retrieve program by id" {
				fakeGorm.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*gorm.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetProgramByID(tt.args.ctx, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetProgramByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_ListPrograms(t *testing.T) {
	orgID := gofakeit.UUID()
	type args struct {
		ctx            context.Context
		organisationID *string
		pagination     *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list programs",
			args: args{
				ctx:            context.Background(),
				organisationID: &orgID,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad caes: failed to get organisation by id",
			args: args{
				ctx:            context.Background(),
				organisationID: &orgID,
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad caes: failed to list programs",
			args: args{
				ctx: nil,
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

			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad caes: failed to get organisation by id" {
				fakeGorm.MockGetOrganisationFn = func(ctx context.Context, id string) (*gorm.Organisation, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad caes: failed to list programs" {
				fakeGorm.MockListProgramsFn = func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}

			got, got1, err := d.ListPrograms(tt.args.ctx, tt.args.organisationID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListPrograms() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_CheckIfSuperUserExists(t *testing.T) {
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
			name: "Happy case: check if super user exists",
			args: args{
				ctx: context.Background(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad caes: failed to check if superuser exists",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad caes: failed to check if superuser exists" {
				fakeGorm.MockCheckIfSuperUserExistsFn = func(ctx context.Context) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.CheckIfSuperUserExists(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfSuperUserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfSuperUserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetCaregiverProfileByUserID(t *testing.T) {
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
			name: "happy case: get caregiver profile by user id",
			args: args{
				ctx:            context.Background(),
				userID:         gofakeit.UUID(),
				organisationID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get caregiver profile by user id",
			args: args{
				ctx:            context.Background(),
				userID:         gofakeit.UUID(),
				organisationID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to check client exists",
			args: args{
				ctx:            context.Background(),
				userID:         gofakeit.UUID(),
				organisationID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get caregiver profile by user id" {
				fakeGorm.MockGetCaregiverProfileByUserIDFn = func(ctx context.Context, userID string, organisationID string) (*gorm.Caregiver, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to check client exists" {
				fakeGorm.MockCheckClientExistsFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetCaregiverProfileByUserID(tt.args.ctx, tt.args.userID, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCaregiverProfileByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetCaregiversClient(t *testing.T) {
	type args struct {
		ctx             context.Context
		caregiverClient domain.CaregiverClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get caregivers client",
			args: args{
				ctx:             context.Background(),
				caregiverClient: domain.CaregiverClient{ClientID: gofakeit.UUID(), CaregiverID: gofakeit.UUID()},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get caregivers clients",
			args: args{
				ctx:             context.Background(),
				caregiverClient: domain.CaregiverClient{ClientID: gofakeit.UUID(), CaregiverID: gofakeit.UUID()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to get caregivers clients" {
				fakeGorm.MockGetCaregiversClientFn = func(ctx context.Context, caregiverClient gorm.CaregiverClient) ([]*gorm.CaregiverClient, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetCaregiversClient(tt.args.ctx, tt.args.caregiverClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCaregiversClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}

func TestMyCareHubDb_GetCaregiverProfileByCaregiverID(t *testing.T) {
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
			name: "happy case: get caregiver profile by id",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to check client exists",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: unable to check client exists" {
				fakeGorm.MockCheckClientExistsFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.GetCaregiverProfileByCaregiverID(tt.args.ctx, tt.args.caregiverID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetCaregiverProfileByCaregiverID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}

func TestMyCareHubDb_ListProgramFacilities(t *testing.T) {
	searchTerm := "ny"
	programID := gofakeit.UUID()
	type args struct {
		ctx                   context.Context
		programID, searchTerm *string
		filterInput           []*dto.FiltersInput
		paginationsInput      *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				programID:  &programID,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to list facilities",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				programID:  &programID,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to retrieve facility identifier",
			args: args{
				ctx:        context.Background(),
				searchTerm: &searchTerm,
				programID:  &programID,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeCountry,
						Value:    "Kenya",
					},
				},
				paginationsInput: &domain.Pagination{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to list facilities" {
				fakeGorm.MockListProgramFacilitiesFn = func(ctx context.Context, programID, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to retrieve facility identifier" {
				fakeGorm.MockRetrieveFacilityIdentifiersByFacilityIDFn = func(ctx context.Context, facilityID *string) ([]*gorm.FacilityIdentifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, got1, err := d.ListProgramFacilities(tt.args.ctx, tt.args.programID, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
			if !tt.wantErr && got1 == nil {
				t.Errorf("did nox expect error, got %v", got1)
			}
		})
	}
}

func TestMyCareHubDb_SearchOrganisation(t *testing.T) {
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
			name: "Happy case: search organisation",
			args: args{
				ctx:             context.Background(),
				searchParameter: "te",
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to list programs",
			args: args{
				ctx:             context.Background(),
				searchParameter: "te",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to search organisation",
			args: args{
				ctx:             context.Background(),
				searchParameter: "te",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to list programs" {
				fakeGorm.MockListProgramsFn = func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to search organisation" {
				fakeGorm.MockSearchOrganisationsFn = func(ctx context.Context, searchParameter string) ([]*gorm.Organisation, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.SearchOrganisation(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_SearchPrograms(t *testing.T) {
	type args struct {
		ctx             context.Context
		searchParameter string
		organisationID  string
		pagination      *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: search programs",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
				organisationID:  uuid.NewString(),
				pagination: &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search programs",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
				organisationID:  uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organisation by id",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
				organisationID:  uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to search programs" {
				fakeGorm.MockSearchProgramsFn = func(ctx context.Context, searchParameter, organisationID string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get organisation by id" {
				fakeGorm.MockGetOrganisationFn = func(ctx context.Context, id string) (*gorm.Organisation, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, _, err := d.SearchPrograms(tt.args.ctx, tt.args.searchParameter, tt.args.organisationID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SearchPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_ListCommunities(t *testing.T) {
	type args struct {
		ctx            context.Context
		programID      string
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Happy case: list communities",
			args: args{
				ctx:            context.Background(),
				programID:      gofakeit.UUID(),
				organisationID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list communities",
			args: args{
				ctx:            context.Background(),
				programID:      gofakeit.UUID(),
				organisationID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to list communities" {
				fakeGorm.MockListCommunitiesFn = func(ctx context.Context, programID, organisationID string) ([]*gorm.Community, error) {
					return nil, errors.New("error")
				}
			}
			_, err := d.ListCommunities(tt.args.ctx, tt.args.programID, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListCommunities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckPhoneExists(t *testing.T) {
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check staff phone exists",
			args: args{
				ctx:   context.Background(),
				phone: gofakeit.Phone(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: failed to check staff phone exists",
			args: args{
				ctx:   context.Background(),
				phone: gofakeit.Phone(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to check staff phone exists" {
				fakeGorm.MockCheckPhoneExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, errors.New("an error occurred")
				}
			}
			got, err := d.CheckPhoneExists(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckPhoneExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckPhoneExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetAvailableScreeningTools(t *testing.T) {
	type args struct {
		ctx              context.Context
		clientID         string
		screeningTool    domain.ScreeningTool
		screeningToolIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get available screening tools",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				screeningTool: domain.ScreeningTool{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					Genders:     []enumutils.Gender{enumutils.GenderMale},
					AgeRange: domain.AgeRange{
						LowerBound: 20,
						UpperBound: 20,
					},
					ProgramID:      gofakeit.UUID(),
					OrganisationID: gofakeit.UUID(),
				},
				screeningToolIDs: []string{gofakeit.UUID()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get available screening tools",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				screeningTool: domain.ScreeningTool{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					Genders:     []enumutils.Gender{enumutils.GenderMale},
					AgeRange: domain.AgeRange{
						LowerBound: 20,
						UpperBound: 20,
					},
					ProgramID:      gofakeit.UUID(),
					OrganisationID: gofakeit.UUID(),
				},
				screeningToolIDs: []string{gofakeit.UUID()},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get questionnaire by id",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				screeningTool: domain.ScreeningTool{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					Genders:     []enumutils.Gender{enumutils.GenderMale},
					AgeRange: domain.AgeRange{
						LowerBound: 20,
						UpperBound: 20,
					},
					ProgramID:      gofakeit.UUID(),
					OrganisationID: gofakeit.UUID(),
				},
				screeningToolIDs: []string{gofakeit.UUID()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to get available screening tools" {
				fakeGorm.MockGetAvailableScreeningToolsFn = func(ctx context.Context, clientID string, screeningTool gorm.ScreeningTool, screeningToolIDs []string) ([]*gorm.ScreeningTool, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get questionnaire by id" {
				fakeGorm.MockGetAvailableScreeningToolsFn = func(ctx context.Context, clientID string, screeningTool gorm.ScreeningTool, screeningToolIDs []string) ([]*gorm.ScreeningTool, error) {
					return []*gorm.ScreeningTool{
						{
							OrganisationID:  uuid.New().String(),
							ID:              uuid.New().String(),
							Active:          true,
							QuestionnaireID: uuid.New().String(),
							Threshold:       10,
							MinimumAge:      0,
							MaximumAge:      0,
						},
					}, nil
				}
				fakeGorm.MockGetQuestionnaireByIDFn = func(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetAvailableScreeningTools(tt.args.ctx, tt.args.clientID, tt.args.screeningTool, tt.args.screeningToolIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetAvailableScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestMyCareHubDb_GetScreeningToolResponsesWithin24Hours(t *testing.T) {
	type args struct {
		ctx       context.Context
		clientID  string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get sscreening tool responses within 24 hours",
			args: args{
				ctx:       context.Background(),
				clientID:  gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get sscreening tool responses within 24 hours",
			args: args{
				ctx:       context.Background(),
				clientID:  gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to get sscreening tool responses within 24 hours" {
				fakeGorm.MockGetScreeningToolResponsesWithin24HoursFn = func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetScreeningToolResponsesWithin24Hours(tt.args.ctx, tt.args.clientID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolResponsesWithin24Hours() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetScreeningToolResponsesWithPendingServiceRequests(t *testing.T) {
	type args struct {
		ctx       context.Context
		clientID  string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get sscreening tool responses with penidng service requests",
			args: args{
				ctx:       context.Background(),
				clientID:  gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get sscreening tool responses with penidng service requests",
			args: args{
				ctx:       context.Background(),
				clientID:  gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to get sscreening tool responses with penidng service requests" {
				fakeGorm.MockGetScreeningToolResponsesWithPendingServiceRequestsFn = func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetScreeningToolResponsesWithPendingServiceRequests(tt.args.ctx, tt.args.clientID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetScreeningToolResponsesWithPendingServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetStaffServiceRequestByID(t *testing.T) {
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

			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case - Failed to get service request by ID" {
				fakeGorm.MockGetStaffServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.StaffServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service request by ID")
				}
			}

			_, err := d.GetStaffServiceRequestByID(tt.args.ctx, tt.args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetStaffServiceRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestMyCareHubDb_GetClientProfilesByIdentifier(t *testing.T) {
	type args struct {
		ctx            context.Context
		identifierType string
		value          string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy case: get client profiles by identifier",
			args: args{
				ctx:            context.Background(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "2332323",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profiles by identifier",
			args: args{
				ctx:            context.Background(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "2332323",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx:            context.Background(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "2332323",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client ccc identifier",
			args: args{
				ctx:            context.Background(),
				identifierType: string(enums.UserIdentifierTypeCCC),
				value:          "2332323",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to get client profiles by identifier" {
				fakeGorm.MockGetClientProfilesByIdentifierFn = func(ctx context.Context, identifierType string, value string) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("failed to get client profiles by identifier")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Sad Case - Fail to get client ccc identifier" {
				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc identifier")
				}
			}

			_, err := d.GetClientProfilesByIdentifier(tt.args.ctx, tt.args.identifierType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetClientProfilesByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetUserProfileByPushToken(t *testing.T) {
	type args struct {
		ctx       context.Context
		pushToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get user profile",
			args: args{
				ctx:       context.Background(),
				pushToken: gofakeit.HipsterSentence(50),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get user profile",
			args: args{
				ctx:       context.Background(),
				pushToken: gofakeit.HipsterSentence(50),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to get user profile" {
				fakeGorm.MockGetUserProfileByPushTokenFn = func(ctx context.Context, pushToken string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			_, err := d.GetUserProfileByPushToken(tt.args.ctx, tt.args.pushToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserProfileByPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_CheckStaffExistsInProgram(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: check if staff is registered in program",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: unable to check if staff is registered in program",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to check if staff is registered in program" {
				fakeGorm.MockCheckStaffExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckStaffExistsInProgram(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckStaffExistsInProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckStaffExistsInProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfFacilityExistsInProgram(t *testing.T) {
	type args struct {
		ctx        context.Context
		programID  string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: check if facility exists in program",
			args: args{
				ctx:        context.Background(),
				programID:  gofakeit.UUID(),
				facilityID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to check if facility exists in program",
			args: args{
				ctx:        context.Background(),
				programID:  gofakeit.UUID(),
				facilityID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to check if facility exists in program" {
				fakeGorm.MockCheckIfFacilityExistsInProgramFn = func(ctx context.Context, programID string, facilityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckIfFacilityExistsInProgram(tt.args.ctx, tt.args.programID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfFacilityExistsInProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfFacilityExistsInProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CheckIfClientExistsInProgram(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    string
		programID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: check if client exist in program",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "sad case: failed check if client exist in program",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed check if client exist in program" {
				fakeGorm.MockCheckIfClientExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CheckIfClientExistsInProgram(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckIfClientExistsInProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckIfClientExistsInProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_GetUserClientProfiles(t *testing.T) {
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
			name: "happy case: get user's client profiles",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - failed to get user's client profiles",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to retrieve facility",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get facility identifiers",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get client facilities",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - failed to get user profile",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - failed to get user's client profiles" {
				fakeGorm.MockGetUserClientProfilesFn = func(ctx context.Context, userID string) ([]*gorm.Client, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to retrieve facility" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to get facility identifiers" {
				fakeGorm.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to get client facilities" {
				fakeGorm.MockGetClientFacilitiesFn = func(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - failed to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetUserClientProfiles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserClientProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetUserStaffProfiles(t *testing.T) {
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
			name: "happy case: get user's staff profiles",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - failed to get user's staff profiles",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get user's staff profiles",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to retrieve facility",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get user profile",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get staff identifier",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - failed to get user's staff profiles" {
				fakeGorm.MockGetUserStaffProfilesFn = func(ctx context.Context, userID string) ([]*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to retrieve facility" {
				fakeGorm.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to get user profile" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - failed to get staff identifier" {
				fakeGorm.MockGetStaffIdentifiersFn = func(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.GetUserStaffProfiles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetUserStaffProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_ListBookings(t *testing.T) {
	type args struct {
		ctx          context.Context
		pagination   *domain.Pagination
		clientID     string
		bookingState enums.BookingState
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list bookings",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				bookingState: enums.UpcomingBooking,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list bookings",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				bookingState: enums.UpcomingBooking,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to list bookings" {
				fakeGorm.MockListBookingsFn = func(ctx context.Context, clientID string, bookingState enums.BookingState, pagination *domain.Pagination) ([]*gorm.Booking, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("error")
				}
			}

			_, _, err := d.ListBookings(tt.args.ctx, tt.args.clientID, tt.args.bookingState, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ListBookings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
