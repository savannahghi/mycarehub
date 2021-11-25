package infrastructure

import (
	"context"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
	SaveOTP(ctx context.Context, otpInput *domain.OTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse *dto.SecurityQuestionResponseInput) error
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id int) (bool, error)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error)
	GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, optedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*domain.ClientProfile, error)
	CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
	GetUserBookmarkedContent(ctx context.Context, userID string) ([]*domain.ContentItem, error)
}

// Update represents all the update action interfaces
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error
	UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error
	UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
	UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}
