package infrastructure

import (
	"context"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
	SaveOTP(ctx context.Context, otpInput *domain.OTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error
	CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	CreateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error)
	CreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error
	CreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	CreateAppointment(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error
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
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error)
	GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, optedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*domain.ClientProfile, error)
	GetStaffProfileByUserID(ctx context.Context, userID string) (*domain.StaffProfile, error)
	CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error)
	CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
	GetUserBookmarkedContent(ctx context.Context, userID string) ([]*domain.ContentItem, error)
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
	GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error)
	CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error)
	GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error)
	GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	GetClientCaregiver(ctx context.Context, caregiverID string) (*domain.Caregiver, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error)
	CheckUserRole(ctx context.Context, userID string, role string) (bool, error)
	CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error)
	GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)
	GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error)
	CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error)
	GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	GetClientCCCIdentifier(ctx context.Context, clientID string) (*domain.Identifier, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	ListAppointments(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	GetScreeningToolQuestions(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error)
	GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error)
	SearchStaffProfileByStaffNumber(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error)
	GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error)
	SearchClientProfilesByCCCNumber(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error)
	CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error)
	GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error)
}

// Update represents all the update action interfaces
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error
	UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error
	UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	UpdateUserProfileAfterLoginSuccess(ctx context.Context, userID string) error
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	LikeContent(ctx context.Context, userID string, contentID int) (bool, error)
	UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error)
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
	ViewContent(ctx context.Context, userID string, contentID int) (bool, error)
	UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	UpdateAppointment(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error
	InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error
	UpdateServiceRequests(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error)
	UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateUserActiveStatus(ctx context.Context, userID string, flavour feedlib.Flavour, active bool) error
	UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateHealthDiary(ctx context.Context, updateInput *gorm.ClientHealthDiaryEntry) (bool, error)
}
