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
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	CreateClient(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error)
	CreateIdentifier(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error)
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
	SaveOTP(ctx context.Context, otpInput *domain.OTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error)
	CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	CreateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error)
	GetOrCreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error
	GetOrCreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	CreateAppointment(ctx context.Context, appointment domain.Appointment) error
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error
	CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	SaveNotification(ctx context.Context, payload *domain.Notification) error
	CreateUserSurveys(ctx context.Context, userSurvey []*dto.UserSurveyInput) error
	CreateMetric(ctx context.Context, payload *domain.Metric) error
	RegisterStaff(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
	SaveFeedback(ctx context.Context, payload *domain.FeedbackResponse) error
	RegisterClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	CreateScreeningTool(ctx context.Context, input *domain.ScreeningTool) error
	CreateScreeningToolResponse(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id int) (bool, error)
	DeleteStaffProfile(ctx context.Context, staffID string) error
	DeleteUser(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error
	DeleteCommunity(ctx context.Context, communityID string) error
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error)
	GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error)
	GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, optedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*domain.ClientProfile, error)
	GetStaffProfileByUserID(ctx context.Context, userID string) (*domain.StaffProfile, error)
	CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
	GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error)
	GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error)
	GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error)
	GetClientCaregiver(ctx context.Context, caregiverID string) (*domain.Caregiver, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error)
	CheckUserRole(ctx context.Context, userID string, role string) (bool, error)
	CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error)
	GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)
	GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error)
	CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error)
	GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error)
	GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	GetClientCCCIdentifier(ctx context.Context, clientID string) (*domain.Identifier, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	ListAppointments(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	ListNotifications(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error)
	ListAvailableNotificationTypes(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error)
	GetScreeningToolQuestions(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error)
	GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error)
	SearchStaffProfile(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error)
	GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error)
	SearchClientProfile(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
	CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error)
	GetStaffProfileByStaffID(ctx context.Context, staffID string) (*domain.StaffProfile, error)
	GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error)
	GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error)
	GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error)
	GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]domain.AppointmentServiceRequests, error)
	GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error)
	GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error)
	CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error)
	GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error)
	GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error)
	GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error)
	GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error)
	GetAppointment(ctx context.Context, params domain.Appointment) (*domain.Appointment, error)
	GetFacilityStaffs(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error)
	CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	GetNotification(ctx context.Context, notificationID string) (*domain.Notification, error)
	GetClientsByFilterParams(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error)
	SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	GetScreeningToolByID(ctx context.Context, screeningToolID string) (*domain.ScreeningTool, error)
	GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error)
	GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error)
	ListSurveyRespondents(ctx context.Context, projectID int, formID string, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error)
	GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error)
	GetScreeningToolResponseByID(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error)
	GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error)
	GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error)
}

// Update represents all the update action interfaces
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	UpdateFacility(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
	UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	UpdateClient(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error)
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error
	AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	UpdateAppointment(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error)
	InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error
	ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error)
	UpdateServiceRequests(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error)
	UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error
	UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error
	UpdateUser(ctx context.Context, user *domain.User, updateData map[string]interface{}) error
	CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error)
	UpdateNotification(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error
	UpdateUserSurveys(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error
	UpdateClientServiceRequest(ctx context.Context, serviceRequest *domain.ServiceRequest, updateData map[string]interface{}) error
	UpdateClientIdentifier(ctx context.Context, clientID string, identifierType string, identifierValue string) error
	UpdateUserContact(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error
	ActivateUser(ctx context.Context, userID string, flavour feedlib.Flavour) error
	DeActivateUser(ctx context.Context, userID string, flavour feedlib.Flavour) error
}
