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
	SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
	SaveOTP(ctx context.Context, otpInput *domain.OTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error)
	CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	CreateCommunity(ctx context.Context, community *domain.Community) (*domain.Community, error)
	GetOrCreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error
	GetOrCreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	CreateAppointment(ctx context.Context, appointment domain.Appointment) error
	CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	SaveNotification(ctx context.Context, payload *domain.Notification) error
	CreateUserSurveys(ctx context.Context, userSurvey []*dto.UserSurveyInput) error
	CreateMetric(ctx context.Context, payload *domain.Metric) error
	RegisterStaff(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
	RegisterExistingUserAsStaff(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
	SaveFeedback(ctx context.Context, payload *domain.FeedbackResponse) error
	RegisterClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	RegisterExistingUserAsClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	RegisterCaregiver(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error)
	CreateCaregiver(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error)
	CreateScreeningTool(ctx context.Context, input *domain.ScreeningTool) error
	CreateScreeningToolResponse(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error)
	AddCaregiverToClient(ctx context.Context, clientCaregiver *domain.CaregiverClient) error
	RegisterExistingUserAsCaregiver(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error)
	CreateOrganisation(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error)
	AddFacilityToProgram(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error)
	CreateProgram(ctx context.Context, input *dto.ProgramInput) (*domain.Program, error)
	CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error)
	CreateSecurityQuestions(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error)
	CreateTermsOfService(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error)
	CreateOauthClientJWT(ctx context.Context, jwt *domain.OauthClientJWT) error
	CreateOauthClient(ctx context.Context, client *domain.OauthClient) error
	CreateOrUpdateSession(ctx context.Context, session *domain.Session) error
	CreateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode) error
	CreateAccessToken(ctx context.Context, token *domain.AccessToken) error
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	DeleteStaffProfile(ctx context.Context, staffID string) error
	DeleteCommunity(ctx context.Context, communityID string) error
	RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error
	RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error
	DeleteOrganisation(ctx context.Context, organisation *domain.Organisation) error
	DeleteAccessToken(ctx context.Context, signature string) error
	DeleteRefreshToken(ctx context.Context, signature string) error
	DeleteClientProfile(ctx context.Context, clientID string, userID *string) error
}

// Query contains all query methods
type Query interface {
	GetCaregiverByUserID(ctx context.Context, userID string) (*domain.Caregiver, error)
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
	ListProgramFacilities(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error)
	GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, optedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	CheckStaffExists(ctx context.Context, userID string) (bool, error)
	CheckClientExists(ctx context.Context, userID string) (bool, error)
	CheckCaregiverExists(ctx context.Context, userID string) (bool, error)
	GetOrganisation(ctx context.Context, id string) (*domain.Organisation, error)
	GetClientProfile(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error)
	GetStaffProfile(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error)
	CheckUserHasPin(ctx context.Context, userID string) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	FindContacts(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error)
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
	GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error)
	GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error)
	GetPendingServiceRequestsCount(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCountResponse, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string, flavour feedlib.Flavour, pagination *domain.Pagination) ([]*domain.ServiceRequest, *domain.Pagination, error)
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)
	GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error)
	CheckIdentifierExists(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error)
	CheckFacilityExistsByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error)
	GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	GetClientIdentifiers(ctx context.Context, clientID string) ([]*domain.Identifier, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	ListAppointments(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	ListNotifications(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error)
	ListAvailableNotificationTypes(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error)
	SearchStaffProfile(ctx context.Context, searchParameter string, programID *string) ([]*domain.StaffProfile, error)
	GetProgramClientProfileByIdentifier(ctx context.Context, programID, identifierType, value string) (*domain.ClientProfile, error)
	GetClientProfilesByIdentifier(ctx context.Context, identifierType, value string) ([]*domain.ClientProfile, error)
	SearchClientProfile(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
	CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	GetStaffProfileByStaffID(ctx context.Context, staffID string) (*domain.StaffProfile, error)
	GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error)
	GetClientServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error)
	GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error)
	GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]domain.AppointmentServiceRequests, error)
	GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error)
	CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error)
	GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error)
	GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error)
	GetAppointment(ctx context.Context, params domain.Appointment) (*domain.Appointment, error)
	GetFacilityStaffs(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error)
	CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	GetNotification(ctx context.Context, notificationID string) (*domain.Notification, error)
	GetClientsByFilterParams(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error)
	SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	GetScreeningToolByID(ctx context.Context, screeningToolID string) (*domain.ScreeningTool, error)
	GetAvailableScreeningTools(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error)
	GetAllScreeningTools(ctx context.Context, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error)
	GetScreeningToolResponsesWithin24Hours(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error)
	GetScreeningToolResponsesWithPendingServiceRequests(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error)
	GetFacilityRespondedScreeningTools(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error)
	ListSurveyRespondents(ctx context.Context, params *domain.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error)
	GetScreeningToolRespondents(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error)
	GetScreeningToolResponseByID(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error)
	GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error)
	GetSurveysWithServiceRequests(ctx context.Context, facilityID, programID string) ([]*dto.SurveysWithServiceRequest, error)
	GetStaffFacilities(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	GetClientFacilities(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error)
	SearchPlatformCaregivers(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error)
	GetCaregiverManagedClients(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error)
	ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error)
	CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error)
	CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error)
	ListOrganisations(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error)
	GetStaffUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error)
	GetClientUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error)
	GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error)
	GetProgramByID(ctx context.Context, programID string) (*domain.Program, error)
	GetCaregiverProfileByUserID(ctx context.Context, userID string, organisationID string) (*domain.CaregiverProfile, error)
	GetCaregiversClient(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error)
	SearchPrograms(ctx context.Context, searchParameter string, organisationID string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error)
	GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*domain.CaregiverProfile, error)
	ListPrograms(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error)
	CheckIfSuperUserExists(ctx context.Context) (bool, error)
	SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error)
	ListCommunities(ctx context.Context, programID string, organisationID string) ([]*domain.Community, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	GetStaffServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error)
	GetClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error)
	GetOauthClient(ctx context.Context, id string) (*domain.OauthClient, error)
	GetValidClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error)
	GetAuthorizationCode(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	GetAccessToken(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error)
	GetRefreshToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error)
	CheckIfClientHasPendingSurveyServiceRequest(ctx context.Context, clientID string, projectID int, formID string) (bool, error)
	GetUserProfileByPushToken(ctx context.Context, pushToken string) (*domain.User, error)
	CheckStaffExistsInProgram(ctx context.Context, userID, programID string) (bool, error)
	CheckIfFacilityExistsInProgram(ctx context.Context, programID, facilityID string) (bool, error)
	CheckIfClientExistsInProgram(ctx context.Context, userID, programID string) (bool, error)
	GetUserClientProfiles(ctx context.Context, userID string) ([]*domain.ClientProfile, error)
	GetUserStaffProfiles(ctx context.Context, userID string) ([]*domain.StaffProfile, error)
	ListBookings(ctx context.Context, clientID string, bookingState enums.BookingState, pagination *domain.Pagination) ([]*domain.Booking, *domain.Pagination, error)
}

// Update represents all the update action interfaces
type Update interface {
	InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	UpdateFacility(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
	UpdateClient(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error)
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error
	UpdateAppointment(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error)
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
	UpdateStaff(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error
	AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) error
	AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) error
	UpdateCaregiverClient(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error
	UpdateCaregiver(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error
	UpdateClientIdentifier(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error
	UpdateUserContact(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error
	UpdateProgram(ctx context.Context, program *domain.Program, updateData map[string]interface{}) error
	UpdateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode, updateData map[string]interface{}) error
	UpdateAccessToken(ctx context.Context, token *domain.AccessToken, updateData map[string]interface{}) error
	UpdateRefreshToken(ctx context.Context, token *domain.RefreshToken, updateData map[string]interface{}) error
	UpdateBooking(ctx context.Context, booking *domain.Booking, updateData map[string]interface{}) error
}
