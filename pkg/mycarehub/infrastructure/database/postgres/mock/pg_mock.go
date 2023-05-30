package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/scalarutils"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	MockCreateUserFn                                          func(ctx context.Context, user domain.User) (*domain.User, error)
	MockCreateClientFn                                        func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error)
	MockCreateIdentifierFn                                    func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error)
	MockListFacilitiesFn                                      func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	MockRetrieveFacilityFn                                    func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	MockListProgramFacilitiesFn                               func(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	MockDeleteFacilityFn                                      func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockRetrieveFacilityByIdentifierFn                        func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
	MockGetUserProfileByUsernameFn                            func(ctx context.Context, username string) (*domain.User, error)
	MockGetUserProfileByPhoneNumberFn                         func(ctx context.Context, phoneNumber string) (*domain.User, error)
	MockGetUserPINByUserIDFn                                  func(ctx context.Context, userID string) (*domain.UserPIN, error)
	MockInactivateFacilityFn                                  func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockReactivateFacilityFn                                  func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockGetUserProfileByUserIDFn                              func(ctx context.Context, userID string) (*domain.User, error)
	MockGetCaregiverByUserIDFn                                func(ctx context.Context, userID string) (*domain.Caregiver, error)
	MockSaveTemporaryUserPinFn                                func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	MockGetCurrentTermsFn                                     func(ctx context.Context) (*domain.TermsOfService, error)
	MockAcceptTermsFn                                         func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                             func(ctx context.Context, pin *domain.UserPIN) (bool, error)
	MockGetSecurityQuestionsFn                                func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	MockSaveOTPFn                                             func(ctx context.Context, otpInput *domain.OTP) error
	MockCheckStaffExistsFn                                    func(ctx context.Context, userID string) (bool, error)
	MockCheckClientExistsFn                                   func(ctx context.Context, userID string) (bool, error)
	MockCheckCaregiverExistsFn                                func(ctx context.Context, userID string) (bool, error)
	MockGetSecurityQuestionByIDFn                             func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn                        func(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	MockGetSecurityQuestionResponseFn                         func(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                            func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                                           func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetClientProfileFn                                    func(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error)
	MockGetStaffProfileFn                                     func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error)
	MockCheckUserHasPinFn                                     func(ctx context.Context, userID string) (bool, error)
	MockGenerateRetryOTPFn                                    func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	MockCompleteOnboardingTourFn                              func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                              func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	MockGetUserSecurityQuestionsResponsesFn                   func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error)
	MockInvalidatePINFn                                       func(ctx context.Context, userID string) (bool, error)
	MockGetContactByUserIDFn                                  func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	MockFindContactsFn                                        func(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn             func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockFetchFacilitiesFn                                     func(ctx context.Context) ([]*domain.Facility, error)
	MockCreateHealthDiaryEntryFn                              func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error)
	MockCreateServiceRequestFn                                func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	MockCanRecordHeathDiaryFn                                 func(ctx context.Context, userID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn                           func(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error)
	MockGetClientHealthDiaryEntriesFn                         func(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error)
	MockCreateClientCaregiverFn                               func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockGetClientCaregiverFn                                  func(ctx context.Context, caregiverID string) (*domain.Caregiver, error)
	MockUpdateClientCaregiverFn                               func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockUpdateFacilityFn                                      func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error
	MockInProgressByFn                                        func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn                          func(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	MockGetPendingServiceRequestsCountFn                      func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error)
	MockGetServiceRequestsFn                                  func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error)
	MockResolveServiceRequestFn                               func(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error
	MockCreateCommunityFn                                     func(ctx context.Context, community *domain.Community) (*domain.Community, error)
	MockCheckIfUsernameExistsFn                               func(ctx context.Context, username string) (bool, error)
	MockGetUsersWithSurveyServiceRequestFn                    func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error)
	MockGetCommunityByIDFn                                    func(ctx context.Context, communityID string) (*domain.Community, error)
	MockCheckIdentifierExists                                 func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error)
	MockCheckFacilityExistsByIdentifier                       func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockGetOrCreateNextOfKin                                  func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error
	MockGetOrCreateContactFn                                  func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	MockGetClientsInAFacilityFn                               func(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	MockGetRecentHealthDiaryEntriesFn                         func(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                                    func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	MockGetClientIdentifiers                                  func(ctx context.Context, clientID string) ([]*domain.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn                       func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	MockCreateAppointment                                     func(ctx context.Context, appointment domain.Appointment) error
	MockUpdateAppointmentFn                                   func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error)
	MockSearchStaffProfileFn                                  func(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error)
	MockUpdateHealthDiaryFn                                   func(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error
	MockUpdateServiceRequestsFn                               func(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error)
	MockListAppointments                                      func(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	MockGetProgramClientProfileByIdentifierFn                 func(ctx context.Context, programID, identifierType, value string) (*domain.ClientProfile, error)
	MockGetClientProfilesByIdentifierFn                       func(ctx context.Context, identifierType, value string) ([]*domain.ClientProfile, error)
	MockUpdateUserPinChangeRequiredStatusFn                   func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockSearchClientProfileFn                                 func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
	MockCheckIfClientHasUnresolvedServiceRequestsFn           func(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	MockUpdateUserSurveysFn                                   func(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error
	MockUpdateUserPinUpdateRequiredStatusFn                   func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockGetHealthDiaryEntryByIDFn                             func(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error)
	MockUpdateClientFn                                        func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error)
	MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn      func(ctx context.Context, userID string, failCount int) error
	MockGetFacilitiesWithoutFHIRIDFn                          func(ctx context.Context) ([]*domain.Facility, error)
	MockGetSharedHealthDiaryEntriesFn                         func(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetClientServiceRequestByIDFn                         func(ctx context.Context, id string) (*domain.ServiceRequest, error)
	MockUpdateUserFn                                          func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error
	MockGetStaffProfileByStaffIDFn                            func(ctx context.Context, staffID string) (*domain.StaffProfile, error)
	MockResolveStaffServiceRequestFn                          func(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error)
	MockCreateStaffServiceRequestFn                           func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	MockGetAppointmentServiceRequestsFn                       func(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error)
	MockGetClientAppointmentByIDFn                            func(ctx context.Context, appointmentID string) (*domain.Appointment, error)
	MockGetAppointmentByAppointmentUUIDFn                     func(ctx context.Context, appointmentUUID string) (*domain.Appointment, error)
	MockGetClientServiceRequestsFn                            func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error)
	MockGetAppointmentByClientIDFn                            func(ctx context.Context, clientID string) (*domain.Appointment, error)
	MockCheckAppointmentExistsByExternalIDFn                  func(ctx context.Context, externalID string) (bool, error)
	MockGetUserSurveyFormsFn                                  func(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error)
	MockListNotificationsFn                                   func(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error)
	MockListAvailableNotificationTypesFn                      func(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error)
	MockSaveNotificationFn                                    func(ctx context.Context, payload *domain.Notification) error
	MockGetClientScreeningToolServiceRequestByToolTypeFn      func(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error)
	MockGetAppointmentFn                                      func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error)
	MockGetFacilityStaffsFn                                   func(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error)
	MockCheckIfStaffHasUnresolvedServiceRequestsFn            func(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	MockDeleteUserFn                                          func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error
	MockDeleteStaffProfileFn                                  func(ctx context.Context, staffID string) error
	MockUpdateNotificationFn                                  func(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error
	MockGetNotificationFn                                     func(ctx context.Context, notificationID string) (*domain.Notification, error)
	MockGetClientsByFilterParamsFn                            func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error)
	MockCreateUserSurveyFn                                    func(ctx context.Context, userSurvey []*dto.UserSurveyInput) error
	MockCreateMetricFn                                        func(ctx context.Context, payload *domain.Metric) error
	MockUpdateClientServiceRequestFn                          func(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error
	MockSaveFeedbackFn                                        func(ctx context.Context, feedback *domain.FeedbackResponse) error
	MockSearchClientServiceRequestsFn                         func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	MockSearchStaffServiceRequestsFn                          func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	MockRegisterClientFn                                      func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	MockRegisterStaffFn                                       func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
	MockRegisterExistingUserAsStaffFn                         func(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
	MockDeleteCommunityFn                                     func(ctx context.Context, communityID string) error
	MockCreateScreeningToolFn                                 func(ctx context.Context, input *domain.ScreeningTool) error
	MockCreateScreeningToolResponseFn                         func(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error)
	MockGetScreeningToolByIDFn                                func(ctx context.Context, toolID string) (*domain.ScreeningTool, error)
	MockGetAvailableScreeningToolsFn                          func(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error)
	MockGetScreeningToolResponsesWithin24HoursFn              func(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error)
	MockGetScreeningToolResponsesWithPendingServiceRequestsFn func(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error)
	MockGetFacilityRespondedScreeningToolsFn                  func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error)
	MockListSurveyRespondentsFn                               func(ctx context.Context, params *domain.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error)
	MockGetScreeningToolRespondentsFn                         func(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error)
	MockGetScreeningToolResponseByIDFn                        func(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error)
	MockGetSurveysWithServiceRequestsFn                       func(ctx context.Context, facilityID, programID string) ([]*dto.SurveysWithServiceRequest, error)
	MockGetStaffFacilitiesFn                                  func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	MockGetClientFacilitiesFn                                 func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	MockUpdateStaffFn                                         func(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error
	MockAddFacilitiesToStaffProfileFn                         func(ctx context.Context, staffID string, facilities []string) error
	MockAddFacilitiesToClientProfileFn                        func(ctx context.Context, clientID string, facilities []string) error
	MockGetUserFacilitiesFn                                   func(ctx context.Context, user *domain.User, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error)
	MockRegisterCaregiverFn                                   func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error)
	MockCreateCaregiverFn                                     func(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error)
	MockSearchCaregiverUserFn                                 func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error)
	MockRemoveFacilitiesFromClientProfileFn                   func(ctx context.Context, clientID string, facilities []string) error
	MockAddCaregiverToClientFn                                func(ctx context.Context, clientCaregiver *domain.CaregiverClient) error
	MockRemoveFacilitiesFromStaffProfileFn                    func(ctx context.Context, staffID string, facilities []string) error
	MockGetOrganisationFn                                     func(ctx context.Context, id string) (*domain.Organisation, error)
	MockGetCaregiverManagedClientsFn                          func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error)
	MockListClientsCaregiversFn                               func(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error)
	MockUpdateCaregiverClientFn                               func(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error
	MockCreateProgramFn                                       func(ctx context.Context, program *dto.ProgramInput) (*domain.Program, error)
	MockCheckOrganisationExistsFn                             func(ctx context.Context, organisationID string) (bool, error)
	MockCheckIfProgramNameExistsFn                            func(ctx context.Context, organisationID string, programName string) (bool, error)
	MockDeleteOrganisationFn                                  func(ctx context.Context, organisation *domain.Organisation) error
	MockCreateOrganisationFn                                  func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error)
	MockAddFacilityToProgramFn                                func(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error)
	MockListOrganisationsFn                                   func(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error)
	MockGetStaffUserProgramsFn                                func(ctx context.Context, userID string) ([]*domain.Program, error)
	MockGetClientUserProgramsFn                               func(ctx context.Context, userID string) ([]*domain.Program, error)
	MockGetProgramFacilitiesFn                                func(ctx context.Context, programID string) ([]*domain.Facility, error)
	MockGetProgramByIDFn                                      func(ctx context.Context, programID string) (*domain.Program, error)
	MockRegisterExistingUserAsClientFn                        func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	MockGetCaregiverProfileByUserIDFn                         func(ctx context.Context, userID string, organisationID string) (*domain.CaregiverProfile, error)
	MockUpdateCaregiverFn                                     func(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error
	MockGetCaregiversClientFn                                 func(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error)
	MockGetCaregiverProfileByCaregiverIDFn                    func(ctx context.Context, caregiverID string) (*domain.CaregiverProfile, error)
	MockRegisterExistingUserAsCaregiverFn                     func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error)
	MockUpdateClientIdentifierFn                              func(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error
	MockUpdateUserContactFn                                   func(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error
	MockSearchProgramsFn                                      func(ctx context.Context, searchParameter string, organisationID string) ([]*domain.Program, error)
	MockListProgramsFn                                        func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error)
	MockCheckIfSuperUserExistsFn                              func(ctx context.Context) (bool, error)
	MockSearchOrganisationsFn                                 func(ctx context.Context, searchParameter string) ([]*domain.Organisation, error)
	MockCreateFacilitiesFn                                    func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error)
	MockListCommunitiesFn                                     func(ctx context.Context, programID string, organisationID string) ([]*domain.Community, error)
	MockCreateSecurityQuestionsFn                             func(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error)
	MockCreateTermsOfServiceFn                                func(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error)
	MockCheckPhoneExistsFn                                    func(ctx context.Context, phone string) (bool, error)
	MockUpdateProgramFn                                       func(ctx context.Context, program *domain.Program, updateData map[string]interface{}) error
	MockGetStaffServiceRequestByIDFn                          func(ctx context.Context, id string) (*domain.ServiceRequest, error)
	MockCreateOauthClientJWT                                  func(ctx context.Context, jwt *domain.OauthClientJWT) error
	MockCreateOauthClient                                     func(ctx context.Context, client *domain.OauthClient) error
	MockGetClientJWT                                          func(ctx context.Context, jti string) (*domain.OauthClientJWT, error)
	MockGetOauthClient                                        func(ctx context.Context, id string) (*domain.OauthClient, error)
	MockGetValidClientJWT                                     func(ctx context.Context, jti string) (*domain.OauthClientJWT, error)
	MockCreateOrUpdateSessionFn                               func(ctx context.Context, session *domain.Session) error
	MockCreateAuthorizationCodeFn                             func(ctx context.Context, code *domain.AuthorizationCode) error
	MockGetAuthorizationCodeFn                                func(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	MockUpdateAuthorizationCodeFn                             func(ctx context.Context, code *domain.AuthorizationCode, updateData map[string]interface{}) error
	MockCreateAccessTokenFn                                   func(ctx context.Context, token *domain.AccessToken) error
	MockCreateRefreshTokenFn                                  func(ctx context.Context, token *domain.RefreshToken) error
	MockDeleteAccessTokenFn                                   func(ctx context.Context, signature string) error
	MockDeleteRefreshTokenFn                                  func(ctx context.Context, signature string) error
	MockGetAccessTokenFn                                      func(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error)
	MockGetRefreshTokenFn                                     func(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error)
	MockUpdateAccessTokenFn                                   func(ctx context.Context, code *domain.AccessToken, updateData map[string]interface{}) error
	MockUpdateRefreshTokenFn                                  func(ctx context.Context, code *domain.RefreshToken, updateData map[string]interface{}) error
	MockCheckIfClientHasPendingSurveyServiceRequestFn         func(ctx context.Context, clientID string, projectID int, formID string) (bool, error)
	MockGetUserProfileByPushTokenFn                           func(ctx context.Context, pushToken string) (*domain.User, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	ID := "f3f8f8f8-f3f8-f3f8-f3f8-f3f8f8f8f8f8"
	screeningUUID := "f3f8f8f8-f3f8-f3f8-f3f8-f3f8f8f8f8f8"

	name := gofakeit.Name()
	country := "Kenya"
	phone := interserviceclient.TestUserPhoneNumber
	description := gofakeit.HipsterSentence(15)
	currentTime := time.Now()

	pastYear := time.Now().AddDate(-3, 0, 0)

	contactData := &domain.Contact{
		ID:           &ID,
		ContactType:  "PHONE",
		ContactValue: "+254711223344",
		Active:       true,
		OptedIn:      true,
		UserID:       &ID,
	}

	facilityInput := &domain.Facility{
		ID:                 &ID,
		Name:               name,
		Phone:              phone,
		Active:             true,
		Country:            country,
		Description:        description,
		FHIROrganisationID: ID,
		Identifier: domain.FacilityIdentifier{
			ID:     ID,
			Active: true,
			Type:   enums.FacilityIdentifierTypeMFLCode,
			Value:  "212121",
		},
	}

	var facilitiesList []*domain.Facility
	facilitiesList = append(facilitiesList, facilityInput)
	nextPage := 3

	userProfile := &domain.User{
		ID:                  &ID,
		Username:            gofakeit.Name(),
		Name:                gofakeit.Name(),
		Active:              true,
		TermsAccepted:       true,
		Gender:              enumutils.GenderMale,
		LastSuccessfulLogin: &currentTime,
		NextAllowedLogin:    &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    3,
		Contacts:            contactData,
		DateOfBirth:         &pastYear,
		CurrentProgramID:    "1234",
		PushTokens:          []string{"gofakeit.HipsterSentence(50)"},
	}

	clientProfile := &domain.ClientProfile{
		ID:                      &ID,
		User:                    userProfile,
		Active:                  false,
		ClientTypes:             []enums.ClientType{},
		UserID:                  ID,
		TreatmentEnrollmentDate: &time.Time{},
		FHIRPatientID:           &ID,
		HealthRecordID:          &ID,
		TreatmentBuddy:          "",
		ClientCounselled:        true,
		OrganisationID:          ID,
		DefaultFacility:         facilityInput,
		CHVUserID:               &ID,
		CHVUserName:             name,
		CaregiverID:             &ID,
		Facilities:              []*domain.Facility{facilityInput},
	}
	staff := &domain.StaffProfile{
		ID:              &ID,
		User:            userProfile,
		UserID:          uuid.New().String(),
		Active:          false,
		StaffNumber:     gofakeit.BeerAlcohol(),
		Facilities:      []*domain.Facility{facilityInput},
		DefaultFacility: facilityInput,
	}

	serviceRequests := []*domain.ServiceRequest{

		{
			ID:           ID,
			ClientID:     uuid.New().String(),
			RequestType:  enums.ServiceRequestTypeRedFlag.String(),
			Status:       enums.ServiceRequestStatusPending.String(),
			InProgressAt: &currentTime,
			InProgressBy: &ID,
			ResolvedAt:   &currentTime,
			ResolvedBy:   &ID,
		},
	}

	healthDiaryEntry := &domain.ClientHealthDiaryEntry{
		ID:                    &ID,
		Active:                false,
		Mood:                  "VERY_SAD",
		Note:                  "test",
		EntryType:             "test",
		ShareWithHealthWorker: false,
		SharedAt:              &currentTime,
		ClientID:              ID,
		CreatedAt:             time.Now(),
		PhoneNumber:           phone,
		ClientName:            name,
	}

	paginationOutput := &domain.Pagination{
		Limit:        10,
		CurrentPage:  1,
		Count:        1,
		TotalPages:   1,
		NextPage:     nil,
		PreviousPage: nil,
		Sort: &domain.SortParam{
			Field:     "id",
			Direction: enums.SortDataTypeDesc,
		},
	}

	caregiverProfile := &domain.CaregiverProfile{
		ID:              ID,
		User:            *userProfile,
		CaregiverNumber: gofakeit.SSN(),
		Consent: domain.ConsentStatus{
			ConsentStatus: enums.ConsentStateAccepted,
		},
	}

	caregiversClients := &domain.CaregiverClient{
		CaregiverID:        ID,
		ClientID:           ID,
		Active:             true,
		RelationshipType:   enums.CaregiverTypeFather,
		CaregiverConsent:   enums.ConsentStateAccepted,
		CaregiverConsentAt: &currentTime,
		ClientConsent:      enums.ConsentStateAccepted,
		ClientConsentAt:    &currentTime,
		OrganisationID:     ID,
		AssignedBy:         ID,
		ProgramID:          ID,
	}

	organisationPayload := domain.Organisation{
		ID:              ID,
		Active:          true,
		Code:            "323",
		Name:            name,
		Description:     description,
		EmailAddress:    "user@domain.com",
		PhoneNumber:     phone,
		PostalAddress:   gofakeit.BS(),
		PhysicalAddress: gofakeit.BS(),
		DefaultCountry:  country,
		Programs: []*domain.Program{
			{
				ID:          ID,
				Active:      true,
				Name:        name,
				Description: description,
				Facilities:  facilitiesList,
			},
		},
	}

	program := &domain.Program{
		ID:           ID,
		Active:       true,
		Name:         name,
		Organisation: organisationPayload,
	}

	return &PostgresMock{
		MockRegisterCaregiverFn: func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
			return &domain.CaregiverProfile{
				ID:              ID,
				User:            *userProfile,
				CaregiverNumber: gofakeit.SSN(),
			}, nil
		},
		MockCreateCaregiverFn: func(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error) {
			return &domain.Caregiver{
				ID:              gofakeit.UUID(),
				UserID:          gofakeit.UUID(),
				CaregiverNumber: gofakeit.SSN(),
				Active:          true,
			}, nil
		},
		MockCreateMetricFn: func(ctx context.Context, payload *domain.Metric) error {
			return nil
		},
		MockUpdateNotificationFn: func(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error {
			return nil
		},
		MockGetCaregiverByUserIDFn: func(ctx context.Context, userID string) (*domain.Caregiver, error) {
			return &domain.Caregiver{
				ID:              ID,
				UserID:          userID,
				CaregiverNumber: gofakeit.SSN(),
				Active:          false,
			}, nil
		},
		MockGetNotificationFn: func(ctx context.Context, notificationID string) (*domain.Notification, error) {
			return &domain.Notification{
				ID:     ID,
				Title:  "A notification",
				Body:   "The notification is about this",
				Type:   "Teleconsult",
				IsRead: false,
			}, nil
		},
		MockGetOrganisationFn: func(ctx context.Context, id string) (*domain.Organisation, error) {
			return &domain.Organisation{
				ID:             ID,
				Active:         true,
				Code:           gofakeit.SSN(),
				Name:           gofakeit.Company(),
				Description:    description,
				EmailAddress:   gofakeit.Email(),
				PhoneNumber:    phone,
				DefaultCountry: gofakeit.Country(),
			}, nil
		},
		MockCheckIdentifierExists: func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
			return false, nil
		},
		MockListCommunitiesFn: func(ctx context.Context, programID, organisationID string) ([]*domain.Community, error) {
			return []*domain.Community{
				{
					ID:          ID,
					RoomID:      ID,
					Name:        gofakeit.Name(),
					Description: gofakeit.BeerName(),
				},
			}, nil
		},
		MockCheckStaffExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockSearchOrganisationsFn: func(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
			return []*domain.Organisation{&organisationPayload}, nil
		},
		MockSearchProgramsFn: func(ctx context.Context, searchParameter, organisationID string) ([]*domain.Program, error) {
			return []*domain.Program{program}, nil
		},
		MockCheckClientExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockCheckCaregiverExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetFacilityStaffsFn: func(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error) {
			return []*domain.StaffProfile{staff}, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
			return facilitiesList, &domain.Pagination{
				Limit:       1,
				CurrentPage: 1,
			}, nil
		},
		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockListProgramFacilitiesFn: func(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
			return facilitiesList, &domain.Pagination{
				Limit:       1,
				CurrentPage: 1,
			}, nil

		},
		MockDeleteFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockRegisterExistingUserAsCaregiverFn: func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
			return &domain.CaregiverProfile{
				ID:              ID,
				User:            *userProfile,
				CaregiverNumber: gofakeit.SSN(),
			}, nil
		},
		MockGetOrCreateContactFn: func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
			return contactData, nil
		},
		MockUpdateClientIdentifierFn: func(ctx context.Context, clientID, identifierType, identifierValue, programID string) error {
			return nil
		},
		MockRetrieveFacilityByIdentifierFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockUpdateUserContactFn: func(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error {
			return nil
		},
		MockRegisterClientFn: func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockRegisterStaffFn: func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
			return staff, nil
		},
		MockGetAppointmentFn: func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
			return &domain.Appointment{
				ID:       ID,
				Reason:   "Bad tooth",
				Provider: "X",
				Date: scalarutils.Date{
					Year:  2023,
					Month: 1,
					Day:   1,
				},
			}, nil
		},
		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*domain.UserPIN, error) {
			return &domain.UserPIN{
				UserID:    userID,
				ValidFrom: time.Now().Add(time.Hour * 10),
				ValidTo:   time.Now().Add(time.Hour * 20),
				IsValid:   false,
			}, nil
		},
		MockGetUserFacilitiesFn: func(ctx context.Context, user *domain.User, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
			nextPage := 2
			previousPage := 0
			return []*domain.Facility{facilityInput}, &domain.Pagination{
				Limit:        10,
				CurrentPage:  1,
				Count:        20,
				TotalPages:   30,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
			}, nil
		},
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*domain.User, error) {
			return userProfile, nil
		},
		MockGetProgramByIDFn: func(ctx context.Context, programID string) (*domain.Program, error) {
			return &domain.Program{
				ID:     programID,
				Active: true,
				Name:   "Test",
				Organisation: domain.Organisation{
					ID: gofakeit.UUID(),
				},
			}, nil
		},
		MockGetUserProfileByUsernameFn: func(ctx context.Context, username string) (*domain.User, error) {
			return userProfile, nil
		},
		MockRegisterExistingUserAsClientFn: func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockUpdateProgramFn: func(ctx context.Context, program *domain.Program, updateData map[string]interface{}) error {
			return nil
		},
		MockGetSurveysWithServiceRequestsFn: func(ctx context.Context, facilityID, programID string) ([]*dto.SurveysWithServiceRequest, error) {
			return []*dto.SurveysWithServiceRequest{
				{
					Title:     "Test survey",
					ProjectID: 10,
					LinkID:    10,
					FormID:    ID,
				},
			}, nil
		},
		MockSearchCaregiverUserFn: func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
			return []*domain.CaregiverProfile{{
				ID:              ID,
				User:            *userProfile,
				CaregiverNumber: gofakeit.SSN(),
			}}, nil
		},
		MockUpdateCaregiverClientFn: func(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error {
			return nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockListClientsCaregiversFn: func(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error) {
			return &domain.ClientCaregivers{
					Caregivers: []*domain.CaregiverProfile{
						{
							ID:              ID,
							User:            *userProfile,
							CaregiverNumber: gofakeit.SSN(),
							Consent: domain.ConsentStatus{
								ConsentStatus: enums.ConsentStateAccepted,
							},
						},
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				}, nil
		},
		MockGetStaffProfileByStaffIDFn: func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
			return &domain.StaffProfile{
				ID:              &ID,
				User:            userProfile,
				UserID:          ID,
				Active:          false,
				StaffNumber:     "TEST-00",
				Facilities:      []*domain.Facility{},
				DefaultFacility: facilityInput,
			}, nil
		},
		MockGetCurrentTermsFn: func(ctx context.Context) (*domain.TermsOfService, error) {
			termsID := gofakeit.Number(1, 1000)
			testText := "test"
			terms := &domain.TermsOfService{
				TermsID: termsID,
				Text:    &testText,
				Flavour: feedlib.FlavourPro,
			}
			return terms, nil
		},
		MockUpdateUserFn: func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
			return nil
		},
		MockAddFacilityToProgramFn: func(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error) {
			return []*domain.Facility{facilityInput}, nil
		},
		MockGetUserProfileByUserIDFn: func(ctx context.Context, userID string) (*domain.User, error) {
			return &domain.User{
				ID:            &userID,
				Username:      gofakeit.Name(),
				Name:          gofakeit.Name(),
				Active:        true,
				TermsAccepted: true,
				Gender:        enumutils.GenderMale,
				Contacts: &domain.Contact{
					ID:           &userID,
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
					Active:       true,
					OptedIn:      true,
				},
				CurrentOrganizationID: uuid.NewString(),
				PushTokens:            []string{"AWERDNDFKLSNJDNFNASDJFNANFKJNADSFNADSKJNSJNFSJKDN"},
			}, nil
		},
		MockFindContactsFn: func(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
			return []*domain.Contact{
				{
					ID:             &ID,
					ContactType:    "PHONE",
					ContactValue:   gofakeit.Phone(),
					Active:         true,
					OptedIn:        true,
					OrganisationID: ID,
				},
			}, nil
		},
		MockListSurveyRespondentsFn: func(ctx context.Context, params *domain.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error) {
			return []*domain.SurveyRespondent{
					{
						ID:          ID,
						Name:        name,
						SubmittedAt: time.Time{},
						ProjectID:   1,
						SubmitterID: 10,
						FormID:      uuid.New().String(),
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				}, nil
		},
		MockCreateIdentifierFn: func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
			return &domain.Identifier{
				ID:                  ID,
				Type:                "CCC",
				Value:               "123456789",
				Use:                 "OFFICIAL",
				Description:         "CCC Number, Primary Identifier",
				ValidFrom:           time.Now(),
				ValidTo:             time.Now(),
				IsPrimaryIdentifier: true,
			}, nil
		},
		MockSaveNotificationFn: func(ctx context.Context, payload *domain.Notification) error {
			return nil
		},
		MockCreateStaffServiceRequestFn: func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
			return nil
		},
		MockUpdateServiceRequestsFn: func(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error) {
			return true, nil
		},
		MockGetSharedHealthDiaryEntriesFn: func(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{healthDiaryEntry}, nil
		},
		MockSaveTemporaryUserPinFn: func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
			return true, nil
		},
		MockGetHealthDiaryEntryByIDFn: func(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
			return &domain.ClientHealthDiaryEntry{
				ID:                    &ID,
				Active:                false,
				Mood:                  "VERY_SAD",
				Note:                  "",
				EntryType:             "",
				ShareWithHealthWorker: false,
				SharedAt:              &currentTime,
				ClientID:              ID,
				CreatedAt:             time.Now(),
				PhoneNumber:           phone,
				ClientName:            name,
			}, nil
		},
		MockAcceptTermsFn: func(ctx context.Context, userID *string, termsID *int) (bool, error) {
			return true, nil
		},
		MockGetPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
			return &domain.ServiceRequestsCountResponse{
				ClientsServiceRequestCount: &domain.ServiceRequestsCount{
					Total: 0,
					RequestsTypeCount: []*domain.RequestTypeCount{
						{
							RequestType: "test",
							Total:       0,
						},
					},
				},
				StaffServiceRequestCount: &domain.ServiceRequestsCount{
					Total: 0,
					RequestsTypeCount: []*domain.RequestTypeCount{
						{
							RequestType: "test",
							Total:       0,
						},
					},
				},
			}, nil
		},
		MockSavePinFn: func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
			return true, nil
		},
		MockGetSecurityQuestionsFn: func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
			securityQuestion := &domain.SecurityQuestion{
				QuestionStem: "test",
				Description:  "test",
				Flavour:      feedlib.FlavourConsumer,
				Active:       true,
				ResponseType: enums.SecurityQuestionResponseTypeNumber,
			}
			return []*domain.SecurityQuestion{securityQuestion}, nil
		},
		MockSaveOTPFn: func(ctx context.Context, otpInput *domain.OTP) error {
			return nil
		},
		MockGetSecurityQuestionByIDFn: func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
			return &domain.SecurityQuestion{
				QuestionStem: "test",
				Description:  "test",
				Flavour:      feedlib.FlavourConsumer,
				Active:       true,
				ResponseType: enums.SecurityQuestionResponseTypeNumber,
			}, nil
		},
		MockSaveSecurityQuestionResponseFn: func(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error {
			return nil
		},
		MockGetSecurityQuestionResponseFn: func(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error) {
			return &domain.SecurityQuestionResponse{
				ResponseID: "1234",
				QuestionID: "1234",
				Active:     true,
				Response:   "Yes",
			}, nil
		},
		MockCheckIfPhoneNumberExistsFn: func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockVerifyOTPFn: func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
			return true, nil
		},
		MockGetClientProfileFn: func(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockGetStaffProfileFn: func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
			return staff, nil
		},
		MockUpdateUserSurveysFn: func(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
			return nil
		},
		MockCheckUserHasPinFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGenerateRetryOTPFn: func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
			return "test-OTP", nil
		},
		MockCompleteOnboardingTourFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGetOTPFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error) {
			return &domain.OTP{
				OTP: "1234",
			}, nil
		},
		MockGetUserSecurityQuestionsResponsesFn: func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
			return []*domain.SecurityQuestionResponse{
				{
					ResponseID: "1234",
					QuestionID: "1234",
					Active:     true,
					Response:   "Yes",
				},
				{
					ResponseID: "1234",
					QuestionID: "1234",
					Active:     true,
					Response:   "Yes",
				},
				{
					ResponseID: "1234",
					QuestionID: "1234",
					Active:     true,
					Response:   "Yes",
				},
			}, nil
		},
		MockInvalidatePINFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetContactByUserIDFn: func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
			return &domain.Contact{
				ID:           userID,
				ContactType:  "PHONE",
				ContactValue: "+254700000000",
				Active:       true,
				OptedIn:      true,
			}, nil
		},
		MockUpdateIsCorrectSecurityQuestionResponseFn: func(ctx context.Context, userID string, isCorrect bool) (bool, error) {
			return true, nil
		},
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error) {
			return healthDiaryEntry, nil
		},
		MockCreateServiceRequestFn: func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
			return nil
		},
		MockCanRecordHeathDiaryFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryQuoteFn: func(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
			return []*domain.ClientHealthDiaryQuote{
				{
					Quote:  "test",
					Author: "test",
				},
			}, nil
		},
		MockListNotificationsFn: func(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error) {
			return []*domain.Notification{
				{
					ID:     ID,
					Title:  "A notification",
					Body:   "The notification is about this",
					Type:   "Teleconsult",
					IsRead: false,
				},
			}, &domain.Pagination{}, nil
		},
		MockListAvailableNotificationTypesFn: func(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error) {
			return []enums.NotificationType{enums.NotificationTypeAppointment}, nil
		},
		MockSearchClientProfileFn: func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
			return []*domain.ClientProfile{clientProfile}, nil
		},
		MockGetClientHealthDiaryEntriesFn: func(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{healthDiaryEntry}, nil
		},
		MockCreateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
			return nil
		},
		MockInProgressByFn: func(ctx context.Context, requestID, staffID string) (bool, error) {
			return true, nil
		},
		MockUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
			return nil
		},
		MockCreateCommunityFn: func(ctx context.Context, community *domain.Community) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.New().String(),
				Name:        name,
				Description: description,
				AgeRange: &domain.AgeRange{
					LowerBound: 10,
					UpperBound: 20,
				},
				Gender: []enumutils.Gender{
					enumutils.AllGender[0],
					enumutils.AllGender[1],
				},
				ClientType: []enums.ClientType{
					enums.AllClientType[0],
					enums.AllClientType[1],
				},
			}, nil
		},
		MockGetClientProfileByClientIDFn: func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
			client := &domain.ClientProfile{
				ID:          &ID,
				User:        userProfile,
				CaregiverID: &ID,
				DefaultFacility: &domain.Facility{
					ID:   &ID,
					Name: name,
				},
			}
			return client, nil
		},
		MockSearchStaffProfileFn: func(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
			return []*domain.StaffProfile{staff}, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
			return serviceRequests, nil
		},
		MockResolveServiceRequestFn: func(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error {
			return nil
		},
		MockCheckIfUsernameExistsFn: func(ctx context.Context, username string) (bool, error) {
			return false, nil
		},
		MockGetCommunityByIDFn: func(ctx context.Context, communityID string) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.New().String(),
				Name:        gofakeit.Name(),
				Description: description,
				AgeRange: &domain.AgeRange{
					LowerBound: 0,
					UpperBound: 0,
				},
				Gender:     []enumutils.Gender{},
				ClientType: []enums.ClientType{},
			}, nil
		},
		MockGetClientsInAFacilityFn: func(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
			return []*domain.ClientProfile{
				clientProfile,
			}, nil
		},
		MockGetFacilitiesWithoutFHIRIDFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return []*domain.Facility{facilityInput}, nil
		},
		MockGetRecentHealthDiaryEntriesFn: func(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockCheckFacilityExistsByIdentifier: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockGetClientIdentifiers: func(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
			return []*domain.Identifier{
				{
					ID:                  uuid.New().String(),
					Type:                "CCC",
					Value:               "123456",
					Use:                 "OFFICIAL",
					Description:         description,
					ValidFrom:           time.Now(),
					ValidTo:             time.Now(),
					IsPrimaryIdentifier: false,
				},
			}, nil
		},
		MockGetServiceRequestsForKenyaEMRFn: func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facilityID := uuid.New().String()
			contact := "123454323"
			serviceReq := &domain.ServiceRequest{
				ID:            ID,
				RequestType:   "SERVICE_REQUEST",
				Request:       "SERVICE_REQUEST",
				Status:        "PENDING",
				ClientID:      ID,
				InProgressAt:  &currentTime,
				InProgressBy:  &staffID,
				ResolvedAt:    &currentTime,
				ResolvedBy:    &staffID,
				FacilityID:    facilityID,
				ClientName:    &staffID,
				ClientContact: &contact,
			}
			return []*domain.ServiceRequest{serviceReq}, nil
		},
		MockCreateAppointment: func(ctx context.Context, appointment domain.Appointment) error {
			return nil
		},
		MockUpdateAppointmentFn: func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
			appointmentDate, _ := utils.ConvertTimeToScalarDate(time.Now())

			return &domain.Appointment{
				ID:         gofakeit.UUID(),
				ClientID:   gofakeit.UUID(),
				ExternalID: gofakeit.UUID(),
				Date:       appointmentDate,
			}, nil
		},
		MockListAppointments: func(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
			return []*domain.Appointment{{
				ID:       ID,
				Reason:   "Bad tooth",
				Provider: "X",
				Date: scalarutils.Date{
					Year:  2023,
					Month: 1,
					Day:   1,
				},
			}}, &domain.Pagination{}, nil
		},
		MockUpdateHealthDiaryFn: func(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
			return nil
		},
		MockAddCaregiverToClientFn: func(ctx context.Context, clientCaregiver *domain.CaregiverClient) error {
			return nil
		},
		MockResolveStaffServiceRequestFn: func(ctx context.Context, staffID, serviceRequestID *string, verificationStatus string) (bool, error) {
			return true, nil
		},
		MockUpdateFacilityFn: func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
			return nil
		},
		MockGetProgramClientProfileByIdentifierFn: func(ctx context.Context, programID, identifierType, value string) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},

		MockGetClientProfilesByIdentifierFn: func(ctx context.Context, identifierType, value string) ([]*domain.ClientProfile, error) {
			return []*domain.ClientProfile{clientProfile}, nil
		},

		MockCheckIfClientHasUnresolvedServiceRequestsFn: func(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
			return false, nil
		},
		MockUpdateUserPinChangeRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockUpdateUserPinUpdateRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockDeleteOrganisationFn: func(ctx context.Context, organisation *domain.Organisation) error {
			return nil
		},
		MockUpdateClientFn: func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
			return client, nil
		},
		MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn: func(ctx context.Context, userID string, failCount int) error {
			return nil
		},
		MockGetClientServiceRequestByIDFn: func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facilityID := uuid.New().String()
			serviceReq := &domain.ServiceRequest{
				ID:           ID,
				RequestType:  "SERVICE_REQUEST",
				Request:      "SERVICE_REQUEST",
				Status:       "PENDING",
				ClientID:     ID,
				InProgressAt: &currentTime,
				InProgressBy: &staffID,
				ResolvedAt:   &currentTime,
				ResolvedBy:   &staffID,
				FacilityID:   facilityID,
			}
			return serviceReq, nil
		},
		MockGetAppointmentServiceRequestsFn: func(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error) {
			return []domain.AppointmentServiceRequests{
				{
					ID:         ID,
					ExternalID: ID,
					Reason:     "reason",
					Date:       scalarutils.Date{Year: 2020, Month: 1, Day: 1},

					InProgressAt:  &currentTime,
					InProgressBy:  &name,
					ResolvedAt:    &currentTime,
					ResolvedBy:    &name,
					ClientName:    &name,
					ClientContact: &phone,
					CCCNumber:     "1234567890",
					MFLCODE:       "1234567890",
				},
			}, nil
		},
		MockGetAppointmentByClientIDFn: func(ctx context.Context, clientID string) (*domain.Appointment, error) {
			return &domain.Appointment{
				ID:         ID,
				ExternalID: ID,
				Reason:     "reason",
				Date:       scalarutils.Date{},
				ClientID:   ID,
				FacilityID: ID,
				Provider:   "provider",
			}, nil

		},
		MockCheckAppointmentExistsByExternalIDFn: func(ctx context.Context, externalID string) (bool, error) {
			return true, nil
		},
		MockSaveFeedbackFn: func(ctx context.Context, feedback *domain.FeedbackResponse) error {
			return nil
		},
		MockGetClientServiceRequestsFn: func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error) {
			return []*domain.ServiceRequest{
				{
					ID:           ID,
					RequestType:  enums.ServiceRequestTypeRedFlag.String(),
					Request:      "SERVICE_REQUEST",
					Status:       "PENDING",
					ClientID:     ID,
					InProgressAt: &currentTime,
					InProgressBy: &name,
					ResolvedAt:   &currentTime,
					ResolvedBy:   &name,
					FacilityID:   uuid.New().String(),
				},
			}, nil
		},
		MockGetUserSurveyFormsFn: func(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
			return []*domain.UserSurvey{
				{
					ID:           uuid.New().String(),
					Active:       false,
					Link:         uuid.New().String(),
					Title:        "SurveyTitle",
					Description:  description,
					HasSubmitted: false,
					UserID:       uuid.New().String(),
				},
			}, nil
		},
		MockGetClientScreeningToolServiceRequestByToolTypeFn: func(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error) {
			return &domain.ServiceRequest{
				ID:           ID,
				RequestType:  enums.ServiceRequestTypeRedFlag.String(),
				Request:      "SERVICE_REQUEST",
				Status:       "PENDING",
				ClientID:     ID,
				InProgressAt: &currentTime,
				InProgressBy: &name,
				ResolvedAt:   &currentTime,
				ResolvedBy:   &name,
				FacilityID:   uuid.New().String(),
			}, nil
		},
		MockCheckIfStaffHasUnresolvedServiceRequestsFn: func(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
			return false, nil
		},
		MockDeleteStaffProfileFn: func(ctx context.Context, staffID string) error {
			return nil
		},
		MockCreateUserFn: func(ctx context.Context, user domain.User) (*domain.User, error) {
			return userProfile, nil
		},
		MockDeleteUserFn: func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
			return nil
		},
		MockGetClientsByFilterParamsFn: func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
			return []*domain.ClientProfile{
				clientProfile,
			}, nil
		},
		MockCreateUserSurveyFn: func(ctx context.Context, userSurvey []*dto.UserSurveyInput) error {
			return nil
		},
		MockSearchClientServiceRequestsFn: func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
			UUID := uuid.New().String()
			return []*domain.ServiceRequest{
				{
					ID:           uuid.NewString(),
					Active:       true,
					RequestType:  "TYPE",
					Request:      "REQUEST",
					Status:       "PENDING",
					InProgressAt: nil,
					ResolvedAt:   nil,
					ClientID:     uuid.New().String(),
					InProgressBy: &UUID,
					ResolvedBy:   &UUID,
					FacilityID:   uuid.New().String(),
					Meta: map[string]interface{}{
						"meta": "meta",
					},
				},
			}, nil
		},
		MockSearchStaffServiceRequestsFn: func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
			UUID := uuid.New().String()
			return []*domain.ServiceRequest{
				{
					ID:           UUID,
					Active:       true,
					RequestType:  "TYPE",
					Request:      "REQUEST",
					Status:       "PENDING",
					ResolvedAt:   nil,
					StaffID:      uuid.New().String(),
					InProgressBy: &UUID,
					ResolvedBy:   &UUID,
					Meta: map[string]interface{}{
						"meta": "meta",
					},
				},
			}, nil
		},
		MockUpdateClientServiceRequestFn: func(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error {
			return nil
		},
		MockGetAvailableScreeningToolsFn: func(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error) {
			return []*domain.ScreeningTool{
				{
					ID:              ID,
					Active:          true,
					QuestionnaireID: ID,
					Threshold:       4,
					ClientTypes:     []enums.ClientType{"PMTCT"},
					Genders:         []enumutils.Gender{"MALE"},
					AgeRange: domain.AgeRange{
						LowerBound: 14,
						UpperBound: 20,
					},
					Questionnaire: domain.Questionnaire{},
				},
			}, nil
		},
		MockGetScreeningToolResponsesWithin24HoursFn: func(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
			return []*domain.QuestionnaireScreeningToolResponse{
				{
					ID:              ID,
					Active:          true,
					ScreeningToolID: ID,
					FacilityID:      ID,
					ClientID:        ID,
					DateOfResponse:  time.Now(),
					AggregateScore:  3,
					QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{
						{
							ID:                      ID,
							Active:                  true,
							ScreeningToolResponseID: ID,
							QuestionID:              ID,
							QuestionText:            gofakeit.Sentence(1),
							Response:                "0",
							NormalizedResponse: map[string]interface{}{
								"0": gofakeit.Sentence(1),
							},
							Score: 2,
						},
					},
				},
			}, nil
		},
		MockGetScreeningToolResponsesWithPendingServiceRequestsFn: func(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
			return []*domain.QuestionnaireScreeningToolResponse{
				{
					ID:              ID,
					Active:          true,
					ScreeningToolID: ID,
					FacilityID:      ID,
					ClientID:        ID,
					DateOfResponse:  time.Now(),
					AggregateScore:  3,
					QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{
						{
							ID:                      ID,
							Active:                  true,
							ScreeningToolResponseID: ID,
							QuestionID:              ID,
							QuestionText:            gofakeit.Sentence(1),
							Response:                "0",
							NormalizedResponse: map[string]interface{}{
								"0": gofakeit.Sentence(1),
							},
							Score: 2,
						},
					},
				},
			}, nil
		},
		MockCreateClientFn: func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockDeleteCommunityFn: func(ctx context.Context, communityID string) error {
			return nil
		},
		MockCreateScreeningToolFn: func(ctx context.Context, input *domain.ScreeningTool) error {
			return nil
		},
		MockCreateScreeningToolResponseFn: func(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error) {
			return &ID, nil
		},
		MockGetScreeningToolByIDFn: func(ctx context.Context, toolID string) (*domain.ScreeningTool, error) {
			return &domain.ScreeningTool{
				ID:              screeningUUID,
				Active:          true,
				QuestionnaireID: screeningUUID,
				Threshold:       0,
				ClientTypes:     []enums.ClientType{enums.ClientTypeDreams},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: domain.AgeRange{
					LowerBound: 21,
					UpperBound: 30,
				},
				Questionnaire: domain.Questionnaire{
					ID:          screeningUUID,
					Active:      true,
					Name:        name,
					Description: description,
					Questions: []domain.Question{
						{
							ID:                screeningUUID,
							Active:            false,
							QuestionnaireID:   screeningUUID,
							Text:              gofakeit.Sentence(10),
							QuestionType:      enums.QuestionTypeCloseEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          0,
							Choices: []domain.QuestionInputChoice{
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: screeningUUID,
									Choice:     "0",
									Value:      gofakeit.Sentence(10),
									Score:      0,
								},
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: screeningUUID,
									Choice:     "1",
									Value:      gofakeit.Sentence(10),
									Score:      0,
								},
							},
						},
					},
				},
			}, nil
		},
		MockGetFacilityRespondedScreeningToolsFn: func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error) {
			return []*domain.ScreeningTool{
					{
						ID:              screeningUUID,
						Active:          true,
						QuestionnaireID: screeningUUID,
						Questionnaire: domain.Questionnaire{
							ID:          screeningUUID,
							Active:      true,
							Name:        name,
							Description: description,
						},
					},
				}, &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				}, nil
		},
		MockGetUsersWithSurveyServiceRequestFn: func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error) {
			return []*domain.SurveyServiceRequestUser{
					{
						Name:        name,
						FormID:      "test",
						ProjectID:   1,
						SubmitterID: 1,
						SurveyName:  "test",
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
					TotalPages:  20,
				}, nil
		},
		MockGetScreeningToolRespondentsFn: func(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error) {
			return []*domain.ScreeningToolRespondent{
					{
						ClientID:                ID,
						ScreeningToolResponseID: screeningToolID,
						ServiceRequestID:        ID,
						Name:                    name,
						PhoneNumber:             phone,
						ServiceRequest:          gofakeit.Sentence(10),
					},
				},
				&domain.Pagination{
					Limit:        1,
					CurrentPage:  1,
					Count:        2,
					TotalPages:   2,
					NextPage:     &nextPage,
					PreviousPage: nil,
				},
				nil
		},
		MockGetScreeningToolResponseByIDFn: func(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
			return &domain.QuestionnaireScreeningToolResponse{
				ID:              ID,
				Active:          true,
				ScreeningToolID: ID,
				FacilityID:      ID,
				ClientID:        ID,
				DateOfResponse:  time.Now(),
				AggregateScore:  3,
				QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{
					{
						ID:                      ID,
						Active:                  true,
						ScreeningToolResponseID: ID,
						QuestionID:              ID,
						QuestionText:            gofakeit.Sentence(1),
						Response:                "0",
						NormalizedResponse: map[string]interface{}{
							"0": gofakeit.Sentence(1),
						},
						Score: 2,
					},
				},
			}, nil
		},
		MockGetStaffFacilitiesFn: func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
			return []*domain.Facility{
					{
						ID:                 &ID,
						Name:               name,
						Phone:              phone,
						Active:             true,
						Country:            country,
						Description:        description,
						FHIROrganisationID: ID,
					},
				}, &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				}, nil
		},
		MockGetClientFacilitiesFn: func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
			return []*domain.Facility{
					{
						ID:                 &ID,
						Name:               name,
						Phone:              phone,
						Active:             true,
						Country:            country,
						Description:        description,
						FHIROrganisationID: ID,
					},
				}, &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				}, nil
		},
		MockUpdateStaffFn: func(ctx context.Context, st *domain.StaffProfile, updates map[string]interface{}) error {
			return nil
		},
		MockAddFacilitiesToStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) error {
			return nil
		},
		MockListOrganisationsFn: func(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error) {
			return []*domain.Organisation{
					{
						ID:              ID,
						Active:          true,
						Code:            "",
						Name:            "Test Organisation",
						Description:     description,
						EmailAddress:    gofakeit.Email(),
						PhoneNumber:     interserviceclient.TestUserPhoneNumber,
						PostalAddress:   gofakeit.BeerAlcohol(),
						PhysicalAddress: gofakeit.BeerAlcohol(),
						DefaultCountry:  gofakeit.Country(),
					},
				}, &domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				}, nil
		},
		MockAddFacilitiesToClientProfileFn: func(ctx context.Context, clientID string, facilities []string) error {
			return nil
		},
		MockRemoveFacilitiesFromClientProfileFn: func(ctx context.Context, clientID string, facilities []string) error {
			return nil
		},
		MockRemoveFacilitiesFromStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) error {
			return nil
		},
		MockRegisterExistingUserAsStaffFn: func(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
			return staff, nil
		},
		MockGetCaregiverManagedClientsFn: func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error) {
			return []*domain.ManagedClient{
				{
					ClientProfile:    clientProfile,
					CaregiverConsent: enums.ConsentStateAccepted,
					ClientConsent:    enums.ConsentStateAccepted,
				},
			}, paginationOutput, nil
		},
		MockCreateProgramFn: func(ctx context.Context, program *dto.ProgramInput) (*domain.Program, error) {
			return &domain.Program{
				ID:     ID,
				Active: true,
				Name:   name,
				Organisation: domain.Organisation{
					ID: ID,
				},
			}, nil
		},
		MockCheckOrganisationExistsFn: func(ctx context.Context, organisationID string) (bool, error) {
			return true, nil
		},
		MockCheckIfProgramNameExistsFn: func(ctx context.Context, organisationID string, programName string) (bool, error) {
			return false, nil
		},
		MockCreateOrganisationFn: func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
			return &organisationPayload, nil
		},
		MockGetStaffUserProgramsFn: func(ctx context.Context, userID string) ([]*domain.Program, error) {
			return []*domain.Program{
				{
					ID:     ID,
					Active: true,
					Name:   name,
					Organisation: domain.Organisation{
						ID:              ID,
						Active:          true,
						Code:            "2121",
						Name:            name,
						Description:     description,
						EmailAddress:    "user@email.com",
						PhoneNumber:     phone,
						PostalAddress:   "322 er",
						PhysicalAddress: "323 er",
						DefaultCountry:  country,
					},
				},
			}, nil
		},
		MockGetClientUserProgramsFn: func(ctx context.Context, userID string) ([]*domain.Program, error) {
			return []*domain.Program{
				{
					ID:     ID,
					Active: true,
					Name:   name,
					Organisation: domain.Organisation{
						ID:              ID,
						Active:          true,
						Code:            "2121",
						Name:            name,
						Description:     description,
						EmailAddress:    "user@email.com",
						PhoneNumber:     phone,
						PostalAddress:   "322 er",
						PhysicalAddress: "323 er",
						DefaultCountry:  country,
					},
				},
			}, nil
		},
		MockGetProgramFacilitiesFn: func(ctx context.Context, programID string) ([]*domain.Facility, error) {
			return []*domain.Facility{
				{
					ID:                 &ID,
					Name:               name,
					Phone:              phone,
					Active:             true,
					Country:            country,
					Description:        description,
					FHIROrganisationID: ID,
					Identifier: domain.FacilityIdentifier{
						ID:     ID,
						Active: true,
						Type:   enums.FacilityIdentifierTypeMFLCode,
						Value:  "12345",
					},
					WorkStationDetails: domain.WorkStationDetails{},
				},
			}, nil
		},
		MockGetCaregiverProfileByUserIDFn: func(ctx context.Context, userID string, organisationID string) (*domain.CaregiverProfile, error) {
			return caregiverProfile, nil
		},
		MockUpdateCaregiverFn: func(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error {
			return nil
		},
		MockGetCaregiversClientFn: func(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
			return []*domain.CaregiverClient{
				caregiversClients,
			}, nil
		},
		MockGetCaregiverProfileByCaregiverIDFn: func(ctx context.Context, caregiverID string) (*domain.CaregiverProfile, error) {
			return caregiverProfile, nil
		},
		MockListProgramsFn: func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error) {
			return []*domain.Program{program}, pagination, nil
		},
		MockCheckIfSuperUserExistsFn: func(ctx context.Context) (bool, error) {
			return false, nil
		},
		MockCreateFacilitiesFn: func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockCreateSecurityQuestionsFn: func(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error) {
			return []*domain.SecurityQuestion{
				{
					SecurityQuestionID: gofakeit.UUID(),
					QuestionStem:       gofakeit.Question(),
					Description:        description,
					ResponseType:       enums.SecurityQuestionResponseTypeText,
					Flavour:            feedlib.FlavourPro,
					Active:             true,
					Sequence:           1,
				},
			}, nil
		},
		MockCreateTermsOfServiceFn: func(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error) {
			return &domain.TermsOfService{
				TermsID:   1,
				Text:      &name,
				ValidFrom: time.Now(),
				ValidTo:   time.Now(),
			}, nil
		},
		MockCheckPhoneExistsFn: func(ctx context.Context, phone string) (bool, error) {
			return false, nil
		},
		MockGetStaffServiceRequestByIDFn: func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facilityID := uuid.New().String()
			serviceReq := &domain.ServiceRequest{
				ID:           ID,
				RequestType:  "SERVICE_REQUEST",
				Request:      "SERVICE_REQUEST",
				Status:       "PENDING",
				ClientID:     ID,
				InProgressAt: &currentTime,
				InProgressBy: &staffID,
				ResolvedAt:   &currentTime,
				ResolvedBy:   &staffID,
				FacilityID:   facilityID,
			}
			return serviceReq, nil
		},
		MockCreateOauthClient: func(ctx context.Context, client *domain.OauthClient) error {
			return nil
		},
		MockCreateOauthClientJWT: func(ctx context.Context, jwt *domain.OauthClientJWT) error {
			return nil
		},
		MockGetClientJWT: func(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
			return &domain.OauthClientJWT{}, nil
		},
		MockGetOauthClient: func(ctx context.Context, id string) (*domain.OauthClient, error) {
			return &domain.OauthClient{
				Active: true,
				Grants: []string{"internal", "internal_refresh_token"},
			}, nil
		},
		MockGetValidClientJWT: func(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
			return &domain.OauthClientJWT{}, nil
		},
		MockCreateOrUpdateSessionFn: func(ctx context.Context, session *domain.Session) error {
			return nil
		},
		MockCreateAuthorizationCodeFn: func(ctx context.Context, code *domain.AuthorizationCode) error {
			return nil
		},
		MockGetAuthorizationCodeFn: func(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
			return &domain.AuthorizationCode{Active: true}, nil
		},
		MockUpdateAuthorizationCodeFn: func(ctx context.Context, code *domain.AuthorizationCode, updateData map[string]interface{}) error {
			return nil
		},
		MockCreateAccessTokenFn: func(ctx context.Context, token *domain.AccessToken) error {
			return nil
		},
		MockCreateRefreshTokenFn: func(ctx context.Context, token *domain.RefreshToken) error {
			return nil
		},
		MockDeleteAccessTokenFn: func(ctx context.Context, signature string) error {
			return nil
		},
		MockDeleteRefreshTokenFn: func(ctx context.Context, signature string) error {
			return nil
		},
		MockGetAccessTokenFn: func(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error) {
			return &domain.AccessToken{Active: true}, nil
		},
		MockGetRefreshTokenFn: func(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
			return &domain.RefreshToken{Active: true}, nil
		},
		MockUpdateAccessTokenFn: func(ctx context.Context, code *domain.AccessToken, updateData map[string]interface{}) error {
			return nil
		},
		MockUpdateRefreshTokenFn: func(ctx context.Context, code *domain.RefreshToken, updateData map[string]interface{}) error {
			return nil
		},
		MockCheckIfClientHasPendingSurveyServiceRequestFn: func(ctx context.Context, clientID string, projectID int, formID string) (bool, error) {
			return false, nil
		},
		MockGetUserProfileByPushTokenFn: func(ctx context.Context, pushToken string) (*domain.User, error) {
			return userProfile, nil
		},
	}
}

// DeleteStaffProfile mocks the implementation of deleting a staff
func (gm *PostgresMock) DeleteStaffProfile(ctx context.Context, staffID string) error {
	return gm.MockDeleteStaffProfileFn(ctx, staffID)
}

// DeleteUser mocks the implementation of deleting a user
func (gm *PostgresMock) DeleteUser(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
	return gm.MockDeleteUserFn(ctx, userID, clientID, staffID, flavour)
}

// CheckStaffExists checks if there is a staff profile that exists for a user
func (gm *PostgresMock) CheckStaffExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckStaffExistsFn(ctx, userID)
}

// CheckClientExists checks if there is a client profile that exists for a user
func (gm *PostgresMock) CheckClientExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckClientExistsFn(ctx, userID)
}

// CheckCaregiverExists checks if there is a caregiver profile that exists for a user
func (gm *PostgresMock) CheckCaregiverExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckCaregiverExistsFn(ctx, userID)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// ListProgramFacilities mocks the implementation of  ListProgramFacilities method.
func (gm *PostgresMock) ListProgramFacilities(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	return gm.MockListProgramFacilitiesFn(ctx, programID, searchTerm, filterInput, paginationsInput)
}

// ListFacilities mocks the implementation of `gorm's` GetFacilities method
func (gm *PostgresMock) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	return gm.MockListFacilitiesFn(ctx, searchTerm, filterInput, paginationsInput)
}

// DeleteFacility mocks the implementation of deleting a facility by ID
func (gm *PostgresMock) DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, identifier)
}

// RetrieveFacilityByIdentifier mocks the implementation of `gorm's` RetrieveFacilityByIdentifier method.
func (gm *PostgresMock) RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityByIdentifierFn(ctx, identifier, isActive)
}

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (gm *PostgresMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// GetUserProfileByUsername retrieves a user using their username
func (gm *PostgresMock) GetUserProfileByUsername(ctx context.Context, username string) (*domain.User, error) {
	return gm.MockGetUserProfileByUsernameFn(ctx, username)
}

// GetUserPINByUserID mocks the get user pin by ID implementation
func (gm *PostgresMock) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *PostgresMock) InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, identifier)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (gm *PostgresMock) ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return gm.MockReactivateFacilityFn(ctx, identifier)
}

// GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *PostgresMock) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}

// GetUserProfileByUserID mocks the implementation of fetching a user profile by userID
func (gm *PostgresMock) GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, userID)
}

// GetStaffUserPrograms retrieves all programs associated with a staff user
func (gm *PostgresMock) GetStaffUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	return gm.MockGetStaffUserProgramsFn(ctx, userID)
}

// GetClientUserPrograms retrieves all programs associated with a client user
func (gm *PostgresMock) GetClientUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	return gm.MockGetClientUserProgramsFn(ctx, userID)
}

// SaveTemporaryUserPin mocks the implementation of saving a temporary user pin
func (gm *PostgresMock) SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return gm.MockSaveTemporaryUserPinFn(ctx, pinData)
}

// AcceptTerms mocks the implementation of accept current terms of service
func (gm *PostgresMock) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	return gm.MockAcceptTermsFn(ctx, userID, termsID)
}

// GetCaregiverByUserID returns the caregiver record of the provided user ID
func (gm *PostgresMock) GetCaregiverByUserID(ctx context.Context, userID string) (*domain.Caregiver, error) {
	return gm.MockGetCaregiverByUserIDFn(ctx, userID)
}

// SavePin mocks the implementation of saving a user pin
func (gm *PostgresMock) SavePin(ctx context.Context, pin *domain.UserPIN) (bool, error) {
	return gm.MockSavePinFn(ctx, pin)
}

// UpdateUserSurveys mocks the implementation of `gorm's` UpdateUserSurveys method.
func (gm *PostgresMock) UpdateUserSurveys(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
	return gm.MockUpdateUserSurveysFn(ctx, survey, updateData)
}

// GetSecurityQuestions mocks the implementation of getting all the security questions.
func (gm *PostgresMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionsFn(ctx, flavour)
}

// SaveOTP mocks the implementation for saving an OTP
func (gm *PostgresMock) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	return gm.MockSaveOTPFn(ctx, otpInput)
}

// GetSecurityQuestionByID mocks the implementation of getting a security question by ID
func (gm *PostgresMock) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionByIDFn(ctx, securityQuestionID)
}

// SaveSecurityQuestionResponse saves the response of a security question
func (gm *PostgresMock) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error {
	return gm.MockSaveSecurityQuestionResponseFn(ctx, securityQuestionResponse)
}

// GetSecurityQuestionResponse mocks the get security question implementation
func (gm *PostgresMock) GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error) {
	return gm.MockGetSecurityQuestionResponseFn(ctx, questionID, userID)
}

// CheckIfPhoneNumberExists mock the implementation of checking the existence of phone number
func (gm *PostgresMock) CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckIfPhoneNumberExistsFn(ctx, phone, isOptedIn, flavour)
}

// VerifyOTP mocks the implementation of verify otp
func (gm *PostgresMock) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return gm.MockVerifyOTPFn(ctx, payload)
}

// GetClientProfile mocks the method for fetching a client profile using the user ID
func (gm *PostgresMock) GetClientProfile(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error) {
	return gm.MockGetClientProfileFn(ctx, userID, programID)
}

// FindContacts retrieves all the contacts that match the given contact type and value.
// Contacts can be shared by users thus the same contact can have multiple records stored
func (gm *PostgresMock) FindContacts(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
	return gm.MockFindContactsFn(ctx, contactType, contactValue)
}

// GetStaffProfile mocks the method for fetching a staff profile using the user ID
func (gm *PostgresMock) GetStaffProfile(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
	return gm.MockGetStaffProfileFn(ctx, userID, programID)
}

// SearchStaffProfile mocks the implementation of getting staff profile using their staff number.
func (gm *PostgresMock) SearchStaffProfile(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	return gm.MockSearchStaffProfileFn(ctx, searchParameter)
}

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *PostgresMock) CheckUserHasPin(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID)
}

// GenerateRetryOTP mock the implementtation of generating a retry OTP
func (gm *PostgresMock) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return gm.MockGenerateRetryOTPFn(ctx, payload)
}

// CompleteOnboardingTour mocks the implementation for updating a user's pin change required state
func (gm *PostgresMock) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCompleteOnboardingTourFn(ctx, userID, flavour)
}

// GetOTP mocks the implementation of fetching an OTP
func (gm *PostgresMock) GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error) {
	return gm.MockGetOTPFn(ctx, phoneNumber, flavour)
}

// GetOrganisation retrieves an organisation using the provided id
func (gm *PostgresMock) GetOrganisation(ctx context.Context, id string) (*domain.Organisation, error) {
	return gm.MockGetOrganisationFn(ctx, id)
}

// GetUserSecurityQuestionsResponses mocks the implementation of getting the user's responded security questions
func (gm *PostgresMock) GetUserSecurityQuestionsResponses(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
	return gm.MockGetUserSecurityQuestionsResponsesFn(ctx, userID, flavour)
}

// InvalidatePIN mocks the implementation of invalidating a user pin
func (gm *PostgresMock) InvalidatePIN(ctx context.Context, userID string) (bool, error) {
	return gm.MockInvalidatePINFn(ctx, userID)
}

// GetContactByUserID mocks the implementation of fetching a contact by userID
func (gm *PostgresMock) GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
	return gm.MockGetContactByUserIDFn(ctx, userID, contactType)
}

// UpdateIsCorrectSecurityQuestionResponse updates the IsCorrectSecurityQuestion response
func (gm *PostgresMock) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	return gm.MockUpdateIsCorrectSecurityQuestionResponseFn(ctx, userID, isCorrectSecurityQuestionResponse)
}

// FetchFacilities mocks the implementation of fetching facility
func (gm *PostgresMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockFetchFacilitiesFn(ctx)
}

// CreateHealthDiaryEntry mocks the method for creating a health diary entry
func (gm *PostgresMock) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error) {
	return gm.MockCreateHealthDiaryEntryFn(ctx, healthDiaryInput)
}

// CreateServiceRequest mocks creating a service request method
func (gm *PostgresMock) CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	return gm.MockCreateServiceRequestFn(ctx, serviceRequestInput)
}

// CanRecordHeathDiary mocks the implementation of checking if a user can record a health diary
func (gm *PostgresMock) CanRecordHeathDiary(ctx context.Context, userID string) (bool, error) {
	return gm.MockCanRecordHeathDiaryFn(ctx, userID)
}

// GetClientHealthDiaryQuote mocks the implementation of fetching client health diary quote
func (gm *PostgresMock) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	return gm.MockGetClientHealthDiaryQuoteFn(ctx, limit)
}

// GetClientHealthDiaryEntries mocks the implementation of getting all health diary entries that belong to a specific user
func (gm *PostgresMock) GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetClientHealthDiaryEntriesFn(ctx, clientID, moodType, shared)
}

// CreateClientCaregiver mocks the implementation of creating a caregiver
func (gm *PostgresMock) CreateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	return gm.MockCreateClientCaregiverFn(ctx, caregiverInput)
}

// GetClientCaregiver mocks the implementation of getting all caregivers for a client
func (gm *PostgresMock) GetClientCaregiver(ctx context.Context, caregiverID string) (*domain.Caregiver, error) {
	return gm.MockGetClientCaregiverFn(ctx, caregiverID)
}

// UpdateClientCaregiver mocks the implementation of updating a caregiver
func (gm *PostgresMock) UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	return gm.MockUpdateClientCaregiverFn(ctx, caregiverInput)
}

// SetInProgressBy mocks the implementation of the `SetInProgressBy` update method
func (gm *PostgresMock) SetInProgressBy(ctx context.Context, requestID, staffID string) (bool, error) {
	return gm.MockInProgressByFn(ctx, requestID, staffID)
}

// GetClientProfileByClientID mocks the implementation of getting a client by client id
func (gm *PostgresMock) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	return gm.MockGetClientProfileByClientIDFn(ctx, clientID)
}

// GetPendingServiceRequestsCount mocks the implementation of getting the service requests count
func (gm *PostgresMock) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
	return gm.MockGetPendingServiceRequestsCountFn(ctx, facilityID)
}

// GetServiceRequests mocks the implementation of getting all service requests for a client
func (gm *PostgresMock) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	return gm.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID, flavour)
}

// ResolveServiceRequest mocks the implementation of resolving a service request
func (gm *PostgresMock) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error {
	return gm.MockResolveServiceRequestFn(ctx, staffID, serviceRequestID, status, action, comment)
}

// CreateCommunity mocks the implementation of creating a channel
func (gm *PostgresMock) CreateCommunity(ctx context.Context, community *domain.Community) (*domain.Community, error) {
	return gm.MockCreateCommunityFn(ctx, community)
}

// CheckIfUsernameExists mocks the implementation of checking whether a username exists
func (gm *PostgresMock) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	return gm.MockCheckIfUsernameExistsFn(ctx, username)
}

// GetCommunityByID mocks the implementation of getting the community by ID
func (gm *PostgresMock) GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error) {
	return gm.MockGetCommunityByIDFn(ctx, communityID)
}

// CheckIdentifierExists mocks checking an identifier exists
func (gm *PostgresMock) CheckIdentifierExists(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
	return gm.MockCheckIdentifierExists(ctx, identifierType, identifierValue)
}

// CheckFacilityExistsByIdentifier mocks checking a facility by mfl codes
func (gm *PostgresMock) CheckFacilityExistsByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return gm.MockCheckFacilityExistsByIdentifier(ctx, identifier)
}

// GetOrCreateNextOfKin mocks creating a next of kin
func (gm *PostgresMock) GetOrCreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
	return gm.MockGetOrCreateNextOfKin(ctx, person, clientID, contactID)
}

// GetOrCreateContact mocks creating a contact
func (gm *PostgresMock) GetOrCreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
	return gm.MockGetOrCreateContactFn(ctx, contact)
}

// GetClientsInAFacility mocks getting all the clients in a facility
func (gm *PostgresMock) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
	return gm.MockGetClientsInAFacilityFn(ctx, facilityID)
}

// GetRecentHealthDiaryEntries mocks getting the most recent health diary entry
func (gm *PostgresMock) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetRecentHealthDiaryEntriesFn(ctx, lastSyncTime, client)
}

// GetClientsByParams retrieves client profiles matching the provided parameters
func (gm *PostgresMock) GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
	return gm.MockGetClientsByParams(ctx, params, lastSyncTime)
}

// GetClientIdentifiers retrieves client's ccc number
func (gm *PostgresMock) GetClientIdentifiers(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
	return gm.MockGetClientIdentifiers(ctx, clientID)
}

// GetServiceRequestsForKenyaEMR mocks the getting of red flag service requests for use by KenyaEMR
func (gm *PostgresMock) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
	return gm.MockGetServiceRequestsForKenyaEMRFn(ctx, payload)
}

// ListAppointments lists appointments based on provided criteria
func (gm *PostgresMock) ListAppointments(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
	return gm.MockListAppointments(ctx, params, filters, pagination)
}

// CreateAppointment creates a new appointment
func (gm *PostgresMock) CreateAppointment(ctx context.Context, appointment domain.Appointment) error {
	return gm.MockCreateAppointment(ctx, appointment)
}

// UpdateAppointment updates an appointment
func (gm *PostgresMock) UpdateAppointment(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
	return gm.MockUpdateAppointmentFn(ctx, appointment, updateData)
}

// UpdateServiceRequests mocks the implementation of updating service requests
func (gm *PostgresMock) UpdateServiceRequests(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error) {
	return gm.MockUpdateServiceRequestsFn(ctx, payload)
}

// GetProgramClientProfileByIdentifier mocks the implementation of getting a client profile using the CCC number
func (gm *PostgresMock) GetProgramClientProfileByIdentifier(ctx context.Context, programID, identifierType, value string) (*domain.ClientProfile, error) {
	return gm.MockGetProgramClientProfileByIdentifierFn(ctx, programID, identifierType, value)
}

// GetClientProfilesByIdentifier mocks the implementation of getting client profiles using identifiers
func (gm *PostgresMock) GetClientProfilesByIdentifier(ctx context.Context, identifierType, value string) ([]*domain.ClientProfile, error) {
	return gm.MockGetClientProfilesByIdentifierFn(ctx, identifierType, value)
}

// CheckIfClientHasUnresolvedServiceRequests mocks the implementation of checking if a client has an unresolved service request
func (gm *PostgresMock) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfClientHasUnresolvedServiceRequestsFn(ctx, clientID, serviceRequestType)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation of updating a user pin change required state
func (gm *PostgresMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour, status)
}

// SearchClientProfile mocks the implementation of searching for client profiles.
// It returns clients profiles whose parts of the CCC number matches
func (gm *PostgresMock) SearchClientProfile(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
	return gm.MockSearchClientProfileFn(ctx, CCCNumber)
}

// UpdateClient updates the client details for a particular client
func (gm *PostgresMock) UpdateClient(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
	return gm.MockUpdateClientFn(ctx, client, updates)
}

// UpdateUserPinUpdateRequiredStatus mocks updating a user `pin update required status`
func (gm *PostgresMock) UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinUpdateRequiredStatusFn(ctx, userID, flavour, status)
}

// UpdateHealthDiary mocks the implementation of updating the share status a health diary entry when the client opts for the sharing
func (gm *PostgresMock) UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
	return gm.MockUpdateHealthDiaryFn(ctx, clientHealthDiaryEntry, updateData)
}

// GetHealthDiaryEntryByID mocks the implementation of getting health diary entry bu a given ID
func (gm *PostgresMock) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetHealthDiaryEntryByIDFn(ctx, healthDiaryEntryID)
}

// UpdateFailedSecurityQuestionsAnsweringAttempts mocks the implementation of resetting failed security attempts
func (gm *PostgresMock) UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error {
	return gm.MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn(ctx, userID, failCount)
}

// GetClientServiceRequestByID mocks the implementation of getting a service request by ID
func (gm *PostgresMock) GetClientServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	return gm.MockGetClientServiceRequestByIDFn(ctx, serviceRequestID)
}

// UpdateUser mocks the implementation of updating a user profile
func (gm *PostgresMock) UpdateUser(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
	return gm.MockUpdateUserFn(ctx, user, updateData)
}

// GetStaffProfileByStaffID mocks the implementation getting staff profile by staff ID
func (gm *PostgresMock) GetStaffProfileByStaffID(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
	return gm.MockGetStaffProfileByStaffIDFn(ctx, staffID)
}

// CreateStaffServiceRequest mocks the implementation creating a staff's service request
func (gm *PostgresMock) CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	return gm.MockCreateStaffServiceRequestFn(ctx, serviceRequestInput)
}

// GetUserSurveyForms mocks the implementation of getting user survey forms
func (gm *PostgresMock) GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
	return gm.MockGetUserSurveyFormsFn(ctx, params)
}

// ResolveStaffServiceRequest mocks the implementation resolving staff service requests
func (gm *PostgresMock) ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error) {
	return gm.MockResolveStaffServiceRequestFn(ctx, staffID, serviceRequestID, verificationStatus)
}

// GetAppointmentServiceRequests mocks the implementation of getting appointments and service requests
func (gm *PostgresMock) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error) {
	return gm.MockGetAppointmentServiceRequestsFn(ctx, lastSyncTime, mflCode)
}

// UpdateFacility mocks the implementation of updating a facility
func (gm *PostgresMock) UpdateFacility(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
	return gm.MockUpdateFacilityFn(ctx, facility, updateData)
}

// GetFacilitiesWithoutFHIRID mocks the implementation of getting a facility without FHIR Organisation
func (gm *PostgresMock) GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockGetFacilitiesWithoutFHIRIDFn(ctx)
}

// GetAppointmentByClientID mocks the implementation of getting a client appointment by id
func (gm *PostgresMock) GetAppointmentByClientID(ctx context.Context, clientID string) (*domain.Appointment, error) {
	return gm.MockGetAppointmentByClientIDFn(ctx, clientID)
}

// CheckAppointmentExistsByExternalID checks if an appointment with the external id exists
func (gm *PostgresMock) CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	return gm.MockCheckAppointmentExistsByExternalIDFn(ctx, externalID)
}

// GetClientServiceRequests mocks the implementation of getting system generated client service requests
func (gm *PostgresMock) GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error) {
	return gm.MockGetClientServiceRequestsFn(ctx, requestType, status, clientID, facilityID)
}

// CreateUser creates a new user
func (gm *PostgresMock) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	return gm.MockCreateUserFn(ctx, user)
}

// CreateClient creates a new client
func (gm *PostgresMock) CreateClient(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
	return gm.MockCreateClientFn(ctx, client, contactID, identifierID)
}

// CreateIdentifier creates a new identifier
func (gm *PostgresMock) CreateIdentifier(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
	return gm.MockCreateIdentifierFn(ctx, identifier)
}

// ListNotifications lists notifications based on the provided parameters
func (gm *PostgresMock) ListNotifications(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error) {
	return gm.MockListNotificationsFn(ctx, params, filters, pagination)
}

// ListAvailableNotificationTypes retrieves the distinct notification types available for a user
func (gm *PostgresMock) ListAvailableNotificationTypes(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error) {
	return gm.MockListAvailableNotificationTypesFn(ctx, params)
}

// GetAppointment fetches an appointment by the external ID
func (gm *PostgresMock) GetAppointment(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
	return gm.MockGetAppointmentFn(ctx, params)
}

// SaveNotification saves the notifications to the database
func (gm *PostgresMock) SaveNotification(ctx context.Context, payload *domain.Notification) error {
	return gm.MockSaveNotificationFn(ctx, payload)
}

// GetSharedHealthDiaryEntries mocks the implementation of getting the most recently shared health diary entires by the client to a health care worker
func (gm *PostgresMock) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetSharedHealthDiaryEntriesFn(ctx, clientID, facilityID)
}

// GetClientScreeningToolServiceRequestByToolType mocks the implementation of getting client screening tool service request by question ID
func (gm *PostgresMock) GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID string, questionID string, status string) (*domain.ServiceRequest, error) {
	return gm.MockGetClientScreeningToolServiceRequestByToolTypeFn(ctx, clientID, questionID, status)
}

// CheckIfStaffHasUnresolvedServiceRequests mocks the implementation of checking if a staff has unresolved service requests
func (gm *PostgresMock) CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfStaffHasUnresolvedServiceRequestsFn(ctx, staffID, serviceRequestType)
}

// GetFacilityStaffs returns a list of staff at a particular facility
func (gm *PostgresMock) GetFacilityStaffs(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error) {
	return gm.MockGetFacilityStaffsFn(ctx, facilityID)
}

// UpdateNotification updates the notification with the provided notification details
func (gm *PostgresMock) UpdateNotification(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error {
	return gm.MockUpdateNotificationFn(ctx, notification, updateData)
}

// GetNotification retrieve a notification using the provided ID
func (gm *PostgresMock) GetNotification(ctx context.Context, notificationID string) (*domain.Notification, error) {
	return gm.MockGetNotificationFn(ctx, notificationID)
}

// GetClientsByFilterParams mocks the implementation of getting clients by filter params
func (gm *PostgresMock) GetClientsByFilterParams(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
	return gm.MockGetClientsByFilterParamsFn(ctx, facilityID, filterParams)
}

// CreateUserSurveys mocks the implementation of creating a user survey
func (gm *PostgresMock) CreateUserSurveys(ctx context.Context, survey []*dto.UserSurveyInput) error {
	return gm.MockCreateUserSurveyFn(ctx, survey)
}

// CreateMetric saves a metric to the database
func (gm *PostgresMock) CreateMetric(ctx context.Context, payload *domain.Metric) error {
	return gm.MockCreateMetricFn(ctx, payload)
}

// UpdateClientServiceRequest updates a service request
func (gm *PostgresMock) UpdateClientServiceRequest(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error {
	return gm.MockUpdateClientServiceRequestFn(ctx, clientServiceRequest, updateData)
}

// RegisterStaff mocks the implementation of registering a staff
func (gm *PostgresMock) RegisterStaff(ctx context.Context, staff *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
	return gm.MockRegisterStaffFn(ctx, staff)
}

// SaveFeedback mocks the implementation of saving feedback into the database
func (gm *PostgresMock) SaveFeedback(ctx context.Context, feedback *domain.FeedbackResponse) error {
	return gm.MockSaveFeedbackFn(ctx, feedback)
}

// SearchClientServiceRequests mocks the implementation of searching client service requests
func (gm *PostgresMock) SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	return gm.MockSearchClientServiceRequestsFn(ctx, searchParameter, requestType, facilityID)
}

// SearchStaffServiceRequests mocks the implementation of searching client service requests
func (gm *PostgresMock) SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	return gm.MockSearchStaffServiceRequestsFn(ctx, searchParameter, requestType, facilityID)
}

// RegisterClient mocks the implementation of registering a client
func (gm *PostgresMock) RegisterClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
	return gm.MockRegisterClientFn(ctx, payload)
}

// DeleteCommunity deletes the specified community from the database
func (gm *PostgresMock) DeleteCommunity(ctx context.Context, communityID string) error {
	return gm.MockDeleteCommunityFn(ctx, communityID)
}

// CreateScreeningTool mocks the implementation of creating a screening tool
func (gm *PostgresMock) CreateScreeningTool(ctx context.Context, payload *domain.ScreeningTool) error {
	return gm.MockCreateScreeningToolFn(ctx, payload)
}

// CreateScreeningToolResponse mocks the implementation of creating a screening tool response
func (gm *PostgresMock) CreateScreeningToolResponse(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error) {
	return gm.MockCreateScreeningToolResponseFn(ctx, input)
}

// GetScreeningToolByID mocks the implementation of getting a screening tool by ID
func (gm *PostgresMock) GetScreeningToolByID(ctx context.Context, toolID string) (*domain.ScreeningTool, error) {
	return gm.MockGetScreeningToolByIDFn(ctx, toolID)
}

// GetAvailableScreeningTools mocks the implementation of getting available screening tools
func (gm *PostgresMock) GetAvailableScreeningTools(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error) {
	return gm.MockGetAvailableScreeningToolsFn(ctx, clientID, screeningTool, screeningToolIDs)
}

// GetScreeningToolResponsesWithin24Hours mocks the implementation of GetScreeningToolResponsesWithin24Hours method
func (gm *PostgresMock) GetScreeningToolResponsesWithin24Hours(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponsesWithin24HoursFn(ctx, clientID, programID)
}

// GetScreeningToolResponsesWithPendingServiceRequests mocks the implementation of GetScreeningToolResponsesWithPendingServiceRequests method
func (gm *PostgresMock) GetScreeningToolResponsesWithPendingServiceRequests(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponsesWithPendingServiceRequestsFn(ctx, clientID, programID)
}

// GetFacilityRespondedScreeningTools mocks the implementation of getting responded screening tools
func (gm *PostgresMock) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error) {
	return gm.MockGetFacilityRespondedScreeningToolsFn(ctx, facilityID, programID, pagination)
}

// ListSurveyRespondents mocks the implementation of listing survey respondents
func (gm *PostgresMock) ListSurveyRespondents(ctx context.Context, params *domain.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error) {
	return gm.MockListSurveyRespondentsFn(ctx, params, facilityID, pagination)
}

// GetScreeningToolRespondents mocks the implementation of getting screening tool respondents
func (gm *PostgresMock) GetScreeningToolRespondents(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error) {
	return gm.MockGetScreeningToolRespondentsFn(ctx, facilityID, programID, screeningToolID, searchTerm, paginationInput)
}

// GetScreeningToolResponseByID mocks the implementation of getting a screening tool response by ID
func (gm *PostgresMock) GetScreeningToolResponseByID(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponseByIDFn(ctx, id)
}

// GetSurveysWithServiceRequests mocks the implementation of getting surveys with service requests
func (gm *PostgresMock) GetSurveysWithServiceRequests(ctx context.Context, facilityID, programID string) ([]*dto.SurveysWithServiceRequest, error) {
	return gm.MockGetSurveysWithServiceRequestsFn(ctx, facilityID, programID)
}

// GetSurveyServiceRequestUser mocks the implementation of getting users with survey service request
func (gm *PostgresMock) GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error) {
	return gm.MockGetUsersWithSurveyServiceRequestFn(ctx, facilityID, projectID, formID, pagination)
}

// GetStaffFacilities mocks the implementation of getting staff facilities
func (gm *PostgresMock) GetStaffFacilities(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	return gm.MockGetStaffFacilitiesFn(ctx, input, pagination)
}

// GetClientFacilities mocks the implementation of getting client facilities
func (gm *PostgresMock) GetClientFacilities(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	return gm.MockGetClientFacilitiesFn(ctx, input, pagination)
}

// UpdateStaff mocks the implementation of updating the staff profile
func (gm *PostgresMock) UpdateStaff(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error {
	return gm.MockUpdateStaffFn(ctx, staff, updates)
}

// AddFacilitiesToStaffProfile mocks the implementation of adding facilities to a staff profile
func (gm *PostgresMock) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return gm.MockAddFacilitiesToStaffProfileFn(ctx, staffID, facilities)
}

// AddFacilitiesToClientProfile mocks the implementation of adding facilities to a client profile
func (gm *PostgresMock) AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return gm.MockAddFacilitiesToClientProfileFn(ctx, clientID, facilities)
}

// GetUserFacilities mocks the implementation of getting user facilities
func (gm *PostgresMock) GetUserFacilities(ctx context.Context, user *domain.User, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	return gm.MockGetUserFacilitiesFn(ctx, user, pagination)
}

// RegisterCaregiver registers a new caregiver on the platform
func (gm *PostgresMock) RegisterCaregiver(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
	return gm.MockRegisterCaregiverFn(ctx, input)
}

// CreateCaregiver creates a caregiver record using the provided input
func (gm *PostgresMock) CreateCaregiver(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error) {
	return gm.MockCreateCaregiverFn(ctx, caregiver)
}

// SearchCaregiverUser mocks the implementation of searching caregivers
func (gm *PostgresMock) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	return gm.MockSearchCaregiverUserFn(ctx, searchParameter)
}

// RemoveFacilitiesFromClientProfile mocks the implementation of removing facilities from a client profile
func (gm *PostgresMock) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return gm.MockRemoveFacilitiesFromClientProfileFn(ctx, clientID, facilities)
}

// AddCaregiverToClient mocks the implementation of adding a caregiver to a client
func (gm *PostgresMock) AddCaregiverToClient(ctx context.Context, clientCaregiver *domain.CaregiverClient) error {
	return gm.MockAddCaregiverToClientFn(ctx, clientCaregiver)
}

// RemoveFacilitiesFromStaffProfile mocks the implementation of removing facilities from a staff profile
func (gm *PostgresMock) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return gm.MockRemoveFacilitiesFromStaffProfileFn(ctx, staffID, facilities)
}

// GetCaregiverManagedClients mocks the implementation of getting caregiver's managed clients
func (gm *PostgresMock) GetCaregiverManagedClients(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error) {
	return gm.MockGetCaregiverManagedClientsFn(ctx, userID, pagination)
}

// ListClientsCaregivers mocks the implementation of listing clients caregivers
func (gm *PostgresMock) ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error) {
	return gm.MockListClientsCaregiversFn(ctx, clientID, pagination)
}

// UpdateCaregiverClient mocks the action of updating a caregiver client details for either client or caregiver.
func (gm *PostgresMock) UpdateCaregiverClient(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error {
	return gm.MockUpdateCaregiverClientFn(ctx, caregiverClient, updateData)
}

// CreateProgram mocks the implementation of creating a program
func (gm *PostgresMock) CreateProgram(ctx context.Context, program *dto.ProgramInput) (*domain.Program, error) {
	return gm.MockCreateProgramFn(ctx, program)
}

// CheckOrganisationExists mocks the implementation checking if the an organisation exists
func (gm *PostgresMock) CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error) {
	return gm.MockCheckOrganisationExistsFn(ctx, organisationID)
}

// CheckIfProgramNameExists mocks the implementation checking if an organisation is associated with a program
func (gm *PostgresMock) CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error) {
	return gm.MockCheckIfProgramNameExistsFn(ctx, organisationID, programName)
}

// CreateOrganisation mocks the implementation of creating an organisation
func (gm *PostgresMock) CreateOrganisation(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
	return gm.MockCreateOrganisationFn(ctx, organisation, programs)
}

// DeleteOrganisation mocks the implementation of deleting an organisation
func (gm *PostgresMock) DeleteOrganisation(ctx context.Context, organisation *domain.Organisation) error {
	return gm.MockDeleteOrganisationFn(ctx, organisation)
}

// RegisterExistingUserAsClient mocks the implementation of registering an existing user as a client
func (gm *PostgresMock) RegisterExistingUserAsClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
	return gm.MockRegisterExistingUserAsClientFn(ctx, payload)
}

// AddFacilityToProgram mocks the implementation of adding a facility to a program
func (gm *PostgresMock) AddFacilityToProgram(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error) {
	return gm.MockAddFacilityToProgramFn(ctx, programID, facilityIDs)
}

// ListOrganisations mocks the implementation of listing organisations
func (gm *PostgresMock) ListOrganisations(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error) {
	return gm.MockListOrganisationsFn(ctx, pagination)
}

// GetProgramFacilities mocks the implementation of getting program facilities
func (gm *PostgresMock) GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error) {
	return gm.MockGetProgramFacilitiesFn(ctx, programID)
}

// GetProgramByID mocks the implementation of getting a program by ID
func (gm *PostgresMock) GetProgramByID(ctx context.Context, programID string) (*domain.Program, error) {
	return gm.MockGetProgramByIDFn(ctx, programID)
}

// GetCaregiverProfileByUserID mocks the implementation of getting a caregiver profile
func (gm *PostgresMock) GetCaregiverProfileByUserID(ctx context.Context, organisationID string, clientID string) (*domain.CaregiverProfile, error) {
	return gm.MockGetCaregiverProfileByUserIDFn(ctx, organisationID, clientID)
}

// UpdateCaregiver mocks the implementation of updating a caregiver
func (gm *PostgresMock) UpdateCaregiver(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error {
	return gm.MockUpdateCaregiverFn(ctx, caregiver, updates)
}

// GetCaregiversClient mocks the implementation of getting caregivers clients
func (gm *PostgresMock) GetCaregiversClient(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
	return gm.MockGetCaregiversClientFn(ctx, caregiverClient)
}

// RegisterExistingUserAsStaff mocks the implementation of registering an existing user as staff
func (gm *PostgresMock) RegisterExistingUserAsStaff(ctx context.Context, input *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
	return gm.MockRegisterExistingUserAsStaffFn(ctx, input)
}

// GetCaregiverProfileByCaregiverID mocks the implementation of getting a caregiver profile by caregiver id
func (gm *PostgresMock) GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*domain.CaregiverProfile, error) {
	return gm.MockGetCaregiverProfileByCaregiverIDFn(ctx, caregiverID)
}

// RegisterExistingUserAsCaregiver mocks the implementation of registering an existing user as a caregiver
func (gm *PostgresMock) RegisterExistingUserAsCaregiver(ctx context.Context, payload *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
	return gm.MockRegisterExistingUserAsCaregiverFn(ctx, payload)
}

// UpdateClientIdentifier mocks the implementation of updating a client's identifiers
func (gm *PostgresMock) UpdateClientIdentifier(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error {
	return gm.MockUpdateClientIdentifierFn(ctx, clientID, identifierType, identifierValue, programID)
}

// UpdateUserContact mocks the implementation of updating a user's contact
func (gm *PostgresMock) UpdateUserContact(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error {
	return gm.MockUpdateUserContactFn(ctx, contact, updateData)
}

// ListPrograms mocks the implementation of getting programs
func (gm *PostgresMock) ListPrograms(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error) {
	return gm.MockListProgramsFn(ctx, organisationID, pagination)
}

// CheckIfSuperUserExists mocks the implementation of checking if a superuser exists
func (gm *PostgresMock) CheckIfSuperUserExists(ctx context.Context) (bool, error) {
	return gm.MockCheckIfSuperUserExistsFn(ctx)
}

// SearchOrganisation mocks the implementation of searching organisations
func (gm *PostgresMock) SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
	return gm.MockSearchOrganisationsFn(ctx, searchParameter)
}

// SearchPrograms mocks the implementation of searching programs
func (gm *PostgresMock) SearchPrograms(ctx context.Context, searchQuery string, organisationID string) ([]*domain.Program, error) {
	return gm.MockSearchProgramsFn(ctx, searchQuery, organisationID)
}

// CreateFacilities Mocks the implementation of CreateFacilities method
func (gm *PostgresMock) CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
	return gm.MockCreateFacilitiesFn(ctx, facilities)
}

// ListCommunities mocks the implementation of listing communities
func (gm *PostgresMock) ListCommunities(ctx context.Context, programID string, organisationID string) ([]*domain.Community, error) {
	return gm.MockListCommunitiesFn(ctx, programID, organisationID)
}

// CreateSecurityQuestions mocks the implementation of CreateSecurityQuestions method
func (gm *PostgresMock) CreateSecurityQuestions(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error) {
	return gm.MockCreateSecurityQuestionsFn(ctx, securityQuestions)
}

// CreateTermsOfService mocks the implementation of CreateTermsOfService method
func (gm *PostgresMock) CreateTermsOfService(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error) {
	return gm.MockCreateTermsOfServiceFn(ctx, termsOfService)
}

// CheckPhoneExists mocks the implementation of CheckPhoneExists method
func (gm *PostgresMock) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	return gm.MockCheckPhoneExistsFn(ctx, phone)
}

// UpdateProgram mocks the implementation of updating a program
func (gm *PostgresMock) UpdateProgram(ctx context.Context, program *domain.Program, updateData map[string]interface{}) error {
	return gm.MockUpdateProgramFn(ctx, program, updateData)
}

// GetStaffServiceRequestByID mocks the implementation of getting a service request by ID
func (gm *PostgresMock) GetStaffServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	return gm.MockGetStaffServiceRequestByIDFn(ctx, serviceRequestID)
}

// CreateOauthClientJWT creates a new oauth jwt client
func (gm *PostgresMock) CreateOauthClientJWT(ctx context.Context, jwt *domain.OauthClientJWT) error {
	return gm.MockCreateOauthClientJWT(ctx, jwt)
}

// CreateOauthClient creates a new oauth client
func (gm *PostgresMock) CreateOauthClient(ctx context.Context, client *domain.OauthClient) error {
	return gm.MockCreateOauthClient(ctx, client)
}

// GetClientJWT retrieves a JWT by unique JTI
func (gm *PostgresMock) GetClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
	return gm.MockGetClientJWT(ctx, jti)
}

// GetOauthClient retrieves a client by ID
func (gm *PostgresMock) GetOauthClient(ctx context.Context, id string) (*domain.OauthClient, error) {
	return gm.MockGetOauthClient(ctx, id)
}

// GetValidClientJWT retrieves a JWT that is still valid i.e not expired
func (gm *PostgresMock) GetValidClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
	return gm.MockGetValidClientJWT(ctx, jti)
}

// CreateOrUpdateSession creates a new session or updates an existing session
func (gm *PostgresMock) CreateOrUpdateSession(ctx context.Context, session *domain.Session) error {
	return gm.MockCreateOrUpdateSessionFn(ctx, session)
}

// CreateAuthorizationCode creates a new authorization code.
func (gm *PostgresMock) CreateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode) error {
	return gm.MockCreateAuthorizationCodeFn(ctx, code)
}

// GetAuthorizationCode retrieves an authorization code using the code
func (gm *PostgresMock) GetAuthorizationCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	return gm.MockGetAuthorizationCodeFn(ctx, code)
}

// UpdateAuthorizationCode updates the details of a given code
func (gm *PostgresMock) UpdateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode, updateData map[string]interface{}) error {
	return gm.MockUpdateAuthorizationCodeFn(ctx, code, updateData)
}

// CreateAccessToken creates a new access token.
func (gm *PostgresMock) CreateAccessToken(ctx context.Context, token *domain.AccessToken) error {
	return gm.MockCreateAccessTokenFn(ctx, token)
}

// CreateRefreshToken creates a new refresh token.
func (gm *PostgresMock) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return gm.MockCreateRefreshTokenFn(ctx, token)
}

// DeleteAccessToken retrieves an access token using the signature
func (gm *PostgresMock) DeleteAccessToken(ctx context.Context, signature string) error {
	return gm.MockDeleteAccessTokenFn(ctx, signature)
}

// DeleteRefreshToken retrieves a refresh token using the signature
func (gm *PostgresMock) DeleteRefreshToken(ctx context.Context, signature string) error {
	return gm.MockDeleteRefreshTokenFn(ctx, signature)
}

// GetAccessToken retrieves an access token using the signature
func (gm *PostgresMock) GetAccessToken(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error) {
	return gm.MockGetAccessTokenFn(ctx, token)
}

// GetRefreshToken retrieves a refresh token using the signature
func (gm *PostgresMock) GetRefreshToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
	return gm.MockGetRefreshTokenFn(ctx, token)
}

// UpdateAccessToken updates the details of a given access token
func (gm *PostgresMock) UpdateAccessToken(ctx context.Context, code *domain.AccessToken, updateData map[string]interface{}) error {
	return gm.MockUpdateAccessTokenFn(ctx, code, updateData)
}

// UpdateRefreshToken updates the details of a given refresh token
func (gm *PostgresMock) UpdateRefreshToken(ctx context.Context, code *domain.RefreshToken, updateData map[string]interface{}) error {
	return gm.MockUpdateRefreshTokenFn(ctx, code, updateData)
}

// CheckIfClientHasPendingSurveyServiceRequest mocks the implementation of CheckIfClientHasPendingSurveyServiceRequest method
func (gm *PostgresMock) CheckIfClientHasPendingSurveyServiceRequest(ctx context.Context, clientID string, projectID int, formID string) (bool, error) {
	return gm.MockCheckIfClientHasPendingSurveyServiceRequestFn(ctx, clientID, projectID, formID)
}

// GetUserProfileByPushToken mocks the retrieving of user profile by push token
func (gm *PostgresMock) GetUserProfileByPushToken(ctx context.Context, pushToken string) (*domain.User, error) {
	return gm.MockGetUserProfileByPushTokenFn(ctx, pushToken)
}
