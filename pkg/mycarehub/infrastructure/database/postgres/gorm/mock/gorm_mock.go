package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockCreateUserFn                                          func(ctx context.Context, user *gorm.User) error
	MockCreateClientFn                                        func(ctx context.Context, client *gorm.Client, contactID, identifierID string) error
	MockCreateIdentifierFn                                    func(ctx context.Context, identifier *gorm.Identifier) error
	MockRetrieveFacilityFn                                    func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByIdentifierFn                        func(ctx context.Context, identifier *gorm.FacilityIdentifier, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityIdentifiersByFacilityIDFn             func(ctx context.Context, facilityID *string) ([]*gorm.FacilityIdentifier, error)
	MockListFacilitiesFn                                      func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error)
	MockDeleteFacilityFn                                      func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error)
	MockListProgramFacilitiesFn                               func(ctx context.Context, programID, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error)
	MockGetUserProfileByUsernameFn                            func(ctx context.Context, username string) (*gorm.User, error)
	MockGetUserProfileByPhoneNumberFn                         func(ctx context.Context, phoneNumber string) (*gorm.User, error)
	MockGetUserPINByUserIDFn                                  func(ctx context.Context, userID string) (*gorm.PINData, error)
	MockInactivateFacilityFn                                  func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error)
	MockReactivateFacilityFn                                  func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error)
	MockGetUserProfileByUserIDFn                              func(ctx context.Context, userID *string) (*gorm.User, error)
	MockSaveTemporaryUserPinFn                                func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetCurrentTermsFn                                     func(ctx context.Context) (*gorm.TermsOfService, error)
	MockAcceptTermsFn                                         func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                             func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetSecurityQuestionsFn                                func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error)
	MockSaveOTPFn                                             func(ctx context.Context, otpInput *gorm.UserOTP) error
	MockGetSecurityQuestionByIDFn                             func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn                        func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error
	MockGetSecurityQuestionResponseFn                         func(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                            func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                                           func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockCheckStaffExistsFn                                    func(ctx context.Context, userID string) (bool, error)
	MockCheckClientExistsFn                                   func(ctx context.Context, userID string) (bool, error)
	MockCheckCaregiverExistsFn                                func(ctx context.Context, userID string) (bool, error)
	MockGetClientProfileFn                                    func(ctx context.Context, userID string, programID string) (*gorm.Client, error)
	MockGetCaregiverByUserIDFn                                func(ctx context.Context, userID string) (*gorm.Caregiver, error)
	MockGetStaffProfileFn                                     func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error)
	MockGetClientsSurveyServiceRequestFn                      func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error)
	MockCheckUserHasPinFn                                     func(ctx context.Context, userID string) (bool, error)
	MockCompleteOnboardingTourFn                              func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                              func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error)
	MockGetUserSecurityQuestionsResponsesFn                   func(ctx context.Context, userID, flavour string) ([]*gorm.SecurityQuestionResponse, error)
	MockInvalidatePINFn                                       func(ctx context.Context, userID string) (bool, error)
	MockGetContactByUserIDFn                                  func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn             func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockCreateHealthDiaryEntryFn                              func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) (*gorm.ClientHealthDiaryEntry, error)
	MockCreateServiceRequestFn                                func(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error
	MockCanRecordHeathDiaryFn                                 func(ctx context.Context, clientID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn                           func(ctx context.Context, limit int) ([]*gorm.ClientHealthDiaryQuote, error)
	MockGetClientHealthDiaryEntriesFn                         func(ctx context.Context, params map[string]interface{}) ([]*gorm.ClientHealthDiaryEntry, error)
	MockUpdateClientCaregiverFn                               func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockInProgressByFn                                        func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn                          func(ctx context.Context, clientID string) (*gorm.Client, error)
	MockGetServiceRequestsFn                                  func(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string) ([]*gorm.ClientServiceRequest, error)
	MockGetClientPendingServiceRequestsCountFn                func(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error)
	MockCreateCommunityFn                                     func(ctx context.Context, community *gorm.Community) (*gorm.Community, error)
	MockCheckIfUsernameExistsFn                               func(ctx context.Context, username string) (bool, error)
	MockGetCommunityByIDFn                                    func(ctx context.Context, communityID string) (*gorm.Community, error)
	MockCheckIdentifierExists                                 func(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	MockCheckFacilityExistsByIdentifier                       func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error)
	MockGetOrCreateNextOfKin                                  func(ctx context.Context, person *gorm.RelatedPerson, clientID, contactID string) error
	MockGetOrCreateContact                                    func(ctx context.Context, contact *gorm.Contact) (*gorm.Contact, error)
	MockGetClientsInAFacilityFn                               func(ctx context.Context, facilityID string) ([]*gorm.Client, error)
	MockGetRecentHealthDiaryEntriesFn                         func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                                    func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error)
	MockGetClientIdentifiers                                  func(ctx context.Context, clientID string) ([]*gorm.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn                       func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error)
	MockCreateAppointment                                     func(ctx context.Context, appointment *gorm.Appointment) error
	MockListAppointments                                      func(ctx context.Context, params *gorm.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error)
	MockUpdateAppointmentFn                                   func(ctx context.Context, appointment *gorm.Appointment, updateData map[string]interface{}) (*gorm.Appointment, error)
	MockUpdateServiceRequestsFn                               func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error)
	MockGetProgramClientProfileByIdentifierFn                 func(ctx context.Context, programID, identifierType, value string) (*gorm.Client, error)
	MockGetClientProfilesByIdentifierFn                       func(ctx context.Context, identifierType, value string) ([]*gorm.Client, error)
	MockSearchClientProfileFn                                 func(ctx context.Context, searchParameter string) ([]*gorm.Client, error)
	MockSearchStaffProfileFn                                  func(ctx context.Context, searchParameter string, programID *string) ([]*gorm.StaffProfile, error)
	MockUpdateUserPinChangeRequiredStatusFn                   func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockCheckIfClientHasUnresolvedServiceRequestsFn           func(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	MockUpdateHealthDiaryFn                                   func(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error
	MockUpdateUserPinUpdateRequiredStatusFn                   func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockUpdateClientFn                                        func(ctx context.Context, client *gorm.Client, updates map[string]interface{}) (*gorm.Client, error)
	MockGetUserProfileByStaffIDFn                             func(ctx context.Context, staffID string) (*gorm.User, error)
	MockGetHealthDiaryEntryByIDFn                             func(ctx context.Context, healthDiaryEntryID string) (*gorm.ClientHealthDiaryEntry, error)
	MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn      func(ctx context.Context, userID string, failCount int) error
	MockGetClientServiceRequestByIDFn                         func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error)
	MockUpdateUserFn                                          func(ctx context.Context, user *gorm.User, updateData map[string]interface{}) error
	MockGetStaffProfileByStaffIDFn                            func(ctx context.Context, staffID string) (*gorm.StaffProfile, error)
	MockCreateStaffServiceRequestFn                           func(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error
	MockGetStaffPendingServiceRequestsCountFn                 func(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCount, error)
	MockGetStaffServiceRequestsFn                             func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.StaffServiceRequest, error)
	MockResolveStaffServiceRequestFn                          func(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error)
	MockGetAppointmentServiceRequestsFn                       func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockUpdateFacilityFn                                      func(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error
	MockGetFacilitiesWithoutFHIRIDFn                          func(ctx context.Context) ([]*gorm.Facility, error)
	MockGetSharedHealthDiaryEntriesFn                         func(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetClientServiceRequestsFn                            func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockCheckAppointmentExistsByExternalIDFn                  func(ctx context.Context, externalID string) (bool, error)
	MockGetUserSurveyFormsFn                                  func(ctx context.Context, params map[string]interface{}) ([]*gorm.UserSurvey, error)
	MockCreateNotificationFn                                  func(ctx context.Context, notification *gorm.Notification) error
	MockUpdateUserSurveysFn                                   func(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error
	MockSearchClientServiceRequestsFn                         func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockSearchStaffServiceRequestsFn                          func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error)
	MockListNotificationsFn                                   func(ctx context.Context, params *gorm.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Notification, *domain.Pagination, error)
	MockListAvailableNotificationTypesFn                      func(ctx context.Context, params *gorm.Notification) ([]enums.NotificationType, error)
	MockGetClientScreeningToolServiceRequestByToolTypeFn      func(ctx context.Context, clientID, toolType, status string) (*gorm.ClientServiceRequest, error)
	MockGetAppointmentFn                                      func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error)
	MockCheckIfStaffHasUnresolvedServiceRequestsFn            func(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	MockGetFacilityStaffsFn                                   func(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error)
	MockDeleteClientProfileFn                                 func(ctx context.Context, clientID string, userID *string) error
	MockDeleteStaffProfileFn                                  func(ctx context.Context, staffID string) error
	MockSaveFeedbackFn                                        func(ctx context.Context, feedback *gorm.Feedback) error
	MockUpdateNotificationFn                                  func(ctx context.Context, notification *gorm.Notification, updateData map[string]interface{}) error
	MockGetNotificationFn                                     func(ctx context.Context, notificationID string) (*gorm.Notification, error)
	MockGetClientsByFilterParamsFn                            func(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*gorm.Client, error)
	MockCreateUserSurveyFn                                    func(ctx context.Context, userSurvey []*gorm.UserSurvey) error
	MockCreateMetricFn                                        func(ctx context.Context, metric *gorm.Metric) error
	MockFindContactsFn                                        func(ctx context.Context, contactType string, contactValue string) ([]*gorm.Contact, error)
	MockRegisterStaffFn                                       func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) (*gorm.StaffProfile, error)
	MockUpdateClientServiceRequestFn                          func(ctx context.Context, clientServiceRequest *gorm.ClientServiceRequest, updateData map[string]interface{}) error
	MockRegisterClientFn                                      func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error)
	MockDeleteCommunityFn                                     func(ctx context.Context, communityID string) error
	MockCreateQuestionnaireFn                                 func(ctx context.Context, input *gorm.Questionnaire) error
	MockCreateScreeningToolFn                                 func(ctx context.Context, input *gorm.ScreeningTool) error
	MockCreateQuestionFn                                      func(ctx context.Context, input *gorm.Question) error
	MockCreateQuestionChoiceFn                                func(ctx context.Context, input *gorm.QuestionInputChoice) error
	MockGetScreeningToolByIDFn                                func(ctx context.Context, toolID string) (*gorm.ScreeningTool, error)
	MockGetQuestionnaireByIDFn                                func(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error)
	MockGetQuestionsByQuestionnaireIDFn                       func(ctx context.Context, questionnaireID string) ([]*gorm.Question, error)
	MockGetQuestionInputChoicesByQuestionIDFn                 func(ctx context.Context, questionID string) ([]*gorm.QuestionInputChoice, error)
	MockCreateScreeningToolResponseFn                         func(ctx context.Context, screeningToolResponse *gorm.ScreeningToolResponse, screeningToolQuestionResponses []*gorm.ScreeningToolQuestionResponse) (*string, error)
	MockGetAvailableScreeningToolsFn                          func(ctx context.Context, clientID string, screeningTool gorm.ScreeningTool, screeningToolIDs []string) ([]*gorm.ScreeningTool, error)
	MockGetScreeningToolResponsesWithin24HoursFn              func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error)
	MockGetScreeningToolResponsesWithPendingServiceRequestsFn func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error)
	MockGetFacilityRespondedScreeningToolsFn                  func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*gorm.ScreeningTool, *domain.Pagination, error)
	MockListSurveyRespondentsFn                               func(ctx context.Context, params *gorm.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*gorm.UserSurvey, *domain.Pagination, error)
	MockGetScreeningToolServiceRequestOfRespondentsFn         func(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error)
	MockGetScreeningToolResponseByIDFn                        func(ctx context.Context, id string) (*gorm.ScreeningToolResponse, error)
	MockGetScreeningToolQuestionResponsesByResponseIDFn       func(ctx context.Context, responseID string) ([]*gorm.ScreeningToolQuestionResponse, error)
	MockGetSurveysWithServiceRequestsFn                       func(ctx context.Context, facilityID, programID string) ([]*gorm.UserSurvey, error)
	MockGetStaffFacilitiesFn                                  func(ctx context.Context, staffFacility gorm.StaffFacilities, pagination *domain.Pagination) ([]*gorm.StaffFacilities, *domain.Pagination, error)
	MockGetClientFacilitiesFn                                 func(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error)
	MockUpdateStaffFn                                         func(ctx context.Context, staff *gorm.StaffProfile, updates map[string]interface{}) (*gorm.StaffProfile, error)
	MockAddFacilitiesToStaffProfileFn                         func(ctx context.Context, staffID string, facilities []string) error
	MockAddFacilitiesToClientProfileFn                        func(ctx context.Context, clientID string, facilities []string) error
	MockGetNotificationsCountFn                               func(ctx context.Context, notification gorm.Notification) (int, error)
	MockRegisterCaregiverFn                                   func(ctx context.Context, user *gorm.User, contact *gorm.Contact, caregiver *gorm.Caregiver) error
	MockCreateCaregiverFn                                     func(ctx context.Context, caregiver *gorm.Caregiver) error
	MockGetClientsSurveyCountFn                               func(ctx context.Context, userID string) (int, error)
	MockSearchCaregiverUserFn                                 func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error)
	MockSearchPlatformCaregiversFn                            func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error)
	MockRemoveFacilitiesFromClientProfileFn                   func(ctx context.Context, clientID string, facilities []string) error
	MockAddCaregiverToClientFn                                func(ctx context.Context, clientCaregiver *gorm.CaregiverClient) error
	MockRemoveFacilitiesFromStaffProfileFn                    func(ctx context.Context, staffID string, facilities []string) error
	MockGetCaregiverManagedClientsFn                          func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error)
	MockGetCaregiversClientFn                                 func(ctx context.Context, caregiverClient gorm.CaregiverClient) ([]*gorm.CaregiverClient, error)
	MockListClientsCaregiversFn                               func(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error)
	MockGetCaregiverProfileByCaregiverIDFn                    func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error)
	MockUpdateCaregiverClientFn                               func(ctx context.Context, caregiverClient *gorm.CaregiverClient, updates map[string]interface{}) error
	MockDeleteOrganisationFn                                  func(ctx context.Context, organisation *gorm.Organisation) error
	MockGetOrganisationFn                                     func(ctx context.Context, id string) (*gorm.Organisation, error)
	MockCreateOrganisationFn                                  func(ctx context.Context, organization *gorm.Organisation) (*gorm.Organisation, error)
	MockGetStaffUserProgramsFn                                func(ctx context.Context, userID string) ([]*gorm.Program, error)
	MockGetClientUserProgramsFn                               func(ctx context.Context, userID string) ([]*gorm.Program, error)
	MockCreateProgramFn                                       func(ctx context.Context, program *gorm.Program) (*gorm.Program, error)
	MockCheckOrganisationExistsFn                             func(ctx context.Context, organisationID string) (bool, error)
	MockCheckIfProgramNameExistsFn                            func(ctx context.Context, organisationID string, programName string) (bool, error)
	MockAddFacilityToProgramFn                                func(ctx context.Context, programID string, facilityID []string) error
	MockRegisterExistingUserAsClientFn                        func(ctx context.Context, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error)
	MockRegisterExistingUserAsStaffFn                         func(ctx context.Context, identifier *gorm.Identifier, staff *gorm.StaffProfile) (*gorm.StaffProfile, error)
	MockListOrganisationsFn                                   func(ctx context.Context, pagination *domain.Pagination) ([]*gorm.Organisation, *domain.Pagination, error)
	MockGetProgramFacilitiesFn                                func(ctx context.Context, programID string) ([]*gorm.ProgramFacility, error)
	MockGetProgramByIDFn                                      func(ctx context.Context, programID string) (*gorm.Program, error)
	MockListProgramsFn                                        func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error)
	MockCheckIfSuperUserExistsFn                              func(ctx context.Context) (bool, error)
	MockRegisterExistingUserAsCaregiverFn                     func(ctx context.Context, caregiver *gorm.Caregiver) (*gorm.Caregiver, error)
	MockUpdateClientIdentifierFn                              func(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error
	MockGetCaregiverProfileByUserIDFn                         func(ctx context.Context, userID string, organisationID string) (*gorm.Caregiver, error)
	MockUpdateCaregiverFn                                     func(ctx context.Context, caregiver *gorm.Caregiver, updates map[string]interface{}) error
	MockSearchProgramsFn                                      func(ctx context.Context, searchParameter string, organisationID string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error)
	MockUpdateUserContactFn                                   func(ctx context.Context, userContact *gorm.Contact, updates map[string]interface{}) error
	MockSearchOrganisationsFn                                 func(ctx context.Context, searchParameter string) ([]*gorm.Organisation, error)
	MockCreateFacilitiesFn                                    func(ctx context.Context, facilities []*gorm.Facility) ([]*gorm.Facility, error)
	MockListCommunitiesFn                                     func(ctx context.Context, programID string, organisationID string) ([]*gorm.Community, error)
	MockCreateSecurityQuestionsFn                             func(ctx context.Context, securityQuestions []*gorm.SecurityQuestion) ([]*gorm.SecurityQuestion, error)
	MockCreateTermsOfServiceFn                                func(ctx context.Context, termsOfService *gorm.TermsOfService) (*gorm.TermsOfService, error)
	MockCheckPhoneExistsFn                                    func(ctx context.Context, phone string) (bool, error)
	MockUpdateProgramFn                                       func(ctx context.Context, program *gorm.Program, updateData map[string]interface{}) error
	MockGetStaffServiceRequestByIDFn                          func(ctx context.Context, serviceRequestID string) (*gorm.StaffServiceRequest, error)
	MockCreateOauthClientJWT                                  func(ctx context.Context, jwt *gorm.OauthClientJWT) error
	MockCreateOauthClient                                     func(ctx context.Context, client *gorm.OauthClient) error
	MockGetClientJWT                                          func(ctx context.Context, jti string) (*gorm.OauthClientJWT, error)
	MockGetOauthClient                                        func(ctx context.Context, id string) (*gorm.OauthClient, error)
	MockGetValidClientJWT                                     func(ctx context.Context, jti string) (*gorm.OauthClientJWT, error)
	MockCreateOrUpdateSessionFn                               func(ctx context.Context, session *gorm.Session) error
	MockCreateAuthorizationCodeFn                             func(ctx context.Context, code *gorm.AuthorizationCode) error
	MockGetAuthorizationCodeFn                                func(ctx context.Context, code string) (*gorm.AuthorizationCode, error)
	MockUpdateAuthorizationCodeFn                             func(ctx context.Context, code *gorm.AuthorizationCode, updateData map[string]interface{}) error
	MockCreateAccessTokenFn                                   func(ctx context.Context, token *gorm.AccessToken) error
	MockCreateRefreshTokenFn                                  func(ctx context.Context, token *gorm.RefreshToken) error
	MockDeleteAccessTokenFn                                   func(ctx context.Context, signature string) error
	MockDeleteRefreshTokenFn                                  func(ctx context.Context, signature string) error
	MockGetAccessTokenFn                                      func(ctx context.Context, token gorm.AccessToken) (*gorm.AccessToken, error)
	MockGetRefreshTokenFn                                     func(ctx context.Context, token gorm.RefreshToken) (*gorm.RefreshToken, error)
	MockUpdateAccessTokenFn                                   func(ctx context.Context, code *gorm.AccessToken, updateData map[string]interface{}) error
	MockUpdateRefreshTokenFn                                  func(ctx context.Context, code *gorm.RefreshToken, updateData map[string]interface{}) error
	MockCheckIfClientHasPendingSurveyServiceRequestFn         func(ctx context.Context, clientID string, projectID int, formID string) (bool, error)
	MockGetUserProfileByPushTokenFn                           func(ctx context.Context, pushToken string) (*gorm.User, error)
	MockCheckStaffExistsInProgramFn                           func(ctx context.Context, userID, programID string) (bool, error)
	MockCheckIfFacilityExistsInProgramFn                      func(ctx context.Context, programID, facilityID string) (bool, error)
	MockGetStaffIdentifiersFn                                 func(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error)
	MockCheckIfClientExistsInProgramFn                        func(ctx context.Context, userID, programID string) (bool, error)
	MockGetUserClientProfilesFn                               func(ctx context.Context, userID string) ([]*gorm.Client, error)
	MockGetUserStaffProfilesFn                                func(ctx context.Context, userID string) ([]*gorm.StaffProfile, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
//
// This initialization initializes all the good cases of your mock tests. i.e all success cases should be defined here.
func NewGormMock() *GormMock {

	/*
		In this section, you find commonly shared success case structs for mock tests
	*/

	ID := gofakeit.Number(300, 400)
	UUID := uuid.New().String()
	name := gofakeit.Name()
	country := "Kenya"
	description := gofakeit.HipsterSentence(15)
	phoneContact := gofakeit.Phone()
	acceptedTermsID := gofakeit.Number(1, 10)
	currentTime := time.Now()

	facility := &gorm.Facility{
		FacilityID:  &UUID,
		Name:        name,
		Active:      true,
		Country:     country,
		Phone:       phoneContact,
		Description: description,
	}

	var facilities []gorm.Facility
	facilities = append(facilities, *facility)

	facilitiesPage := []*gorm.Facility{
		{
			FacilityID:  &UUID,
			Name:        name,
			Active:      true,
			Country:     country,
			Description: description,
		},
	}

	fhirID := uuid.New().String()
	clientProfile := &gorm.Client{
		ID:          &UUID,
		Active:      true,
		ClientTypes: []string{"PMTCT"},
		User: gorm.User{
			UserID:                 &UUID,
			Username:               gofakeit.Name(),
			Gender:                 enumutils.GenderMale,
			Active:                 true,
			Contacts:               gorm.Contact{},
			PushTokens:             []string{},
			LastSuccessfulLogin:    &currentTime,
			LastFailedLogin:        &currentTime,
			FailedLoginCount:       3,
			NextAllowedLogin:       &currentTime,
			TermsAccepted:          true,
			AcceptedTermsID:        &acceptedTermsID,
			Avatar:                 gofakeit.URL(),
			IsSuspended:            true,
			PinChangeRequired:      true,
			HasSetPin:              true,
			HasSetSecurityQuestion: true,
			IsPhoneVerified:        true,
			CurrentOrganisationID:  uuid.New().String(),
			IsSuperuser:            true,
			Name:                   name,
			DateOfBirth:            &currentTime,
		},
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &fhirID,
		HealthRecordID:          &UUID,
		ClientCounselled:        true,
		OrganisationID:          uuid.New().String(),
		FacilityID:              uuid.New().String(),
		UserID:                  &UUID,
	}

	userProfile := &gorm.User{
		UserID:                 &UUID,
		Username:               gofakeit.Name(),
		Gender:                 enumutils.GenderMale,
		Active:                 true,
		Contacts:               gorm.Contact{},
		PushTokens:             []string{gofakeit.HipsterSentence(50)},
		LastSuccessfulLogin:    &currentTime,
		LastFailedLogin:        &currentTime,
		FailedLoginCount:       3,
		NextAllowedLogin:       &currentTime,
		TermsAccepted:          true,
		AcceptedTermsID:        &acceptedTermsID,
		Avatar:                 "test",
		IsSuspended:            true,
		PinChangeRequired:      true,
		HasSetPin:              true,
		HasSetSecurityQuestion: true,
		IsPhoneVerified:        true,
		CurrentOrganisationID:  uuid.New().String(),
		IsSuperuser:            true,
		Name:                   name,
		DateOfBirth:            &currentTime,
	}

	staff := &gorm.StaffProfile{
		ID:                &UUID,
		UserProfile:       *userProfile,
		UserID:            uuid.New().String(),
		Active:            true,
		StaffNumber:       gofakeit.BeerAlcohol(),
		Facilities:        []gorm.Facility{*facility},
		DefaultFacilityID: gofakeit.BeerAlcohol(),
		OrganisationID:    gofakeit.BeerAlcohol(),
	}

	pinData := &gorm.PINData{
		PINDataID: &ID,
		UserID:    gofakeit.UUID(),
		HashedPIN: uuid.New().String(),
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
	}

	nowTime := time.Now()
	laterTime := nowTime.Add(time.Hour * 24)
	caregiverID := gofakeit.UUID()
	serviceRequests := []*gorm.ClientServiceRequest{
		{
			ID:             &UUID,
			ClientID:       uuid.New().String(),
			Active:         true,
			RequestType:    enums.ServiceRequestTypeRedFlag.String(),
			Status:         enums.ServiceRequestStatusPending.String(),
			InProgressAt:   &nowTime,
			InProgressByID: &UUID,
			ResolvedAt:     &laterTime,
			ResolvedByID:   &UUID,
			Meta:           `{"formID": "test", "projectID": 1, "submitterID": 1, "surveyName": "test"}`,
			CaregiverID:    &caregiverID,
		},
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

	caregiversClient := gorm.CaregiverClient{
		CaregiverID:        UUID,
		ClientID:           UUID,
		Active:             true,
		RelationshipType:   enums.CaregiverTypeFather,
		CaregiverConsent:   enums.ConsentStateAccepted,
		CaregiverConsentAt: &currentTime,
		ClientConsent:      enums.ConsentStateAccepted,
		ClientConsentAt:    &currentTime,
		OrganisationID:     UUID,
		AssignedBy:         UUID,
		ProgramID:          UUID,
	}

	program := gorm.Program{
		ID:             UUID,
		Active:         true,
		Name:           name,
		Description:    gofakeit.Sentence(10),
		OrganisationID: UUID,
	}

	return &GormMock{
		MockCreateMetricFn: func(ctx context.Context, metric *gorm.Metric) error {
			return nil
		},
		MockRegisterCaregiverFn: func(ctx context.Context, user *gorm.User, contact *gorm.Contact, caregiver *gorm.Caregiver) error {
			return nil
		},
		MockCreateNotificationFn: func(ctx context.Context, notification *gorm.Notification) error {
			return nil
		},
		MockGetNotificationFn: func(ctx context.Context, notificationID string) (*gorm.Notification, error) {
			return &gorm.Notification{
				Title:      "A notification",
				Body:       "This is what it's about",
				Type:       "TELECONSULT",
				IsRead:     false,
				UserID:     &UUID,
				FacilityID: &UUID,
			}, nil
		},
		MockUpdateNotificationFn: func(ctx context.Context, notification *gorm.Notification, updateData map[string]interface{}) error {
			return nil
		},
		MockCreateCaregiverFn: func(ctx context.Context, caregiver *gorm.Caregiver) error {
			return nil
		},
		MockGetFacilityStaffsFn: func(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error) {
			return []*gorm.StaffProfile{
				staff,
			}, nil
		},
		MockGetCaregiverByUserIDFn: func(ctx context.Context, userID string) (*gorm.Caregiver, error) {
			return &gorm.Caregiver{
				ID:              gofakeit.UUID(),
				Active:          true,
				CaregiverNumber: gofakeit.SSN(),
				UserID:          userID,
			}, nil
		},
		MockCreateUserFn: func(ctx context.Context, user *gorm.User) error {
			return nil
		},
		MockSearchProgramsFn: func(ctx context.Context, searchParameter string, organisationID string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
			return []*gorm.Program{
					{
						ID:             UUID,
						Name:           name,
						OrganisationID: UUID,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				}, nil
		},
		MockSaveFeedbackFn: func(ctx context.Context, feedback *gorm.Feedback) error {
			return nil
		},
		MockCreateClientFn: func(ctx context.Context, client *gorm.Client, contactID, identifierID string) error {
			return nil
		},
		MockCreateIdentifierFn: func(ctx context.Context, identifier *gorm.Identifier) error {
			return nil
		},
		MockGetFacilityRespondedScreeningToolsFn: func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*gorm.ScreeningTool, *domain.Pagination, error) {
			return []*gorm.ScreeningTool{
					{
						ID:              UUID,
						Active:          true,
						QuestionnaireID: UUID,
						Threshold:       1,
						ClientTypes:     []string{enums.ClientTypeHighRisk.String()},
						Genders:         []string{enumutils.GenderMale.String()},
						MinimumAge:      18,
						MaximumAge:      25,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				}, nil
		},
		MockListSurveyRespondentsFn: func(ctx context.Context, params *gorm.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*gorm.UserSurvey, *domain.Pagination, error) {
			return []*gorm.UserSurvey{
					{
						Base: gorm.Base{
							UpdatedAt: time.Now(),
						},

						ID:             UUID,
						Active:         true,
						Link:           "https://www.google.com",
						Title:          "Test",
						Description:    description,
						HasSubmitted:   true,
						FormID:         "1",
						ProjectID:      ID,
						LinkID:         ID,
						Token:          "",
						SubmittedAt:    &time.Time{},
						UserID:         UUID,
						OrganisationID: UUID,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				}, nil
		},
		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {

			return facility, nil
		},
		MockCheckStaffExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockCheckClientExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockCheckCaregiverExistsFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetClientsSurveyCountFn: func(ctx context.Context, userID string) (int, error) {
			return 1, nil
		},
		MockGetProgramByIDFn: func(ctx context.Context, programID string) (*gorm.Program, error) {
			return &gorm.Program{
				ID:             programID,
				Active:         true,
				Name:           "Test",
				OrganisationID: UUID,
			}, nil
		},
		MockGetStaffPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCount, error) {
			return &domain.ServiceRequestsCount{
				Total: 20,
				RequestsTypeCount: []*domain.RequestTypeCount{
					{
						RequestType: enums.ServiceRequestTypeStaffPinReset,
						Total:       10,
					},
				},
			}, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
			return facilitiesPage, pagination, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
			return true, nil
		},
		MockGetClientsSurveyServiceRequestFn: func(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
			return serviceRequests, &domain.Pagination{
				Limit:       10,
				CurrentPage: 1,
			}, nil
		},
		MockRetrieveFacilityByIdentifierFn: func(ctx context.Context, identifier *gorm.FacilityIdentifier, isActive bool) (*gorm.Facility, error) {
			return facility, nil
		},
		MockRetrieveFacilityIdentifiersByFacilityIDFn: func(ctx context.Context, facilityID *string) ([]*gorm.FacilityIdentifier, error) {
			return []*gorm.FacilityIdentifier{
				{
					ID:         UUID,
					Active:     true,
					Type:       "MFLCode",
					Value:      "21332433",
					FacilityID: UUID,
				},
			}, nil
		},
		MockRegisterStaffFn: func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockDeleteOrganisationFn: func(ctx context.Context, organisation *gorm.Organisation) error {
			return nil
		},
		MockGetAvailableScreeningToolsFn: func(ctx context.Context, clientID string, screeningTool gorm.ScreeningTool, screeningToolIDs []string) ([]*gorm.ScreeningTool, error) {
			return []*gorm.ScreeningTool{
				{
					OrganisationID:  uuid.New().String(),
					ID:              UUID,
					Active:          true,
					QuestionnaireID: uuid.New().String(),
					Threshold:       4,
					ClientTypes:     pq.StringArray{"PMTCT"},
					Genders:         pq.StringArray{"MALE"},
					MinimumAge:      14,
					MaximumAge:      24,
				},
			}, nil
		},
		MockGetScreeningToolResponsesWithin24HoursFn: func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
			return []*gorm.ScreeningToolResponse{
				{
					ID:              UUID,
					Active:          true,
					ScreeningToolID: UUID,
					FacilityID:      UUID,
					ClientID:        uuid.New().String(),
					AggregateScore:  3,
				},
			}, nil
		},
		MockGetScreeningToolResponsesWithPendingServiceRequestsFn: func(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
			return []*gorm.ScreeningToolResponse{
				{
					ID:              UUID,
					Active:          true,
					ScreeningToolID: UUID,
					FacilityID:      UUID,
					ClientID:        uuid.New().String(),
					AggregateScore:  3,
				},
			}, nil
		},
		MockGetSharedHealthDiaryEntriesFn: func(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error) {
			return []*gorm.ClientHealthDiaryEntry{
				{
					ClientHealthDiaryEntryID: &UUID,
					Active:                   true,
					Mood:                     "Bad",
					Note:                     "Note",
					EntryType:                "EntryType",
					ShareWithHealthWorker:    true,
					SharedAt:                 &currentTime,
					ClientID:                 UUID,
					OrganisationID:           UUID,
				},
			}, nil
		},
		MockGetAppointmentFn: func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error) {
			date := time.Now().Add(time.Duration(100))

			return &gorm.Appointment{
				ID:             gofakeit.UUID(),
				OrganisationID: gofakeit.UUID(),
				Active:         true,
				ExternalID:     strconv.Itoa(gofakeit.Number(0, 1000)),
				ClientID:       gofakeit.UUID(),
				FacilityID:     gofakeit.UUID(),
				Reason:         "Knocked up",
				Date:           date,
			}, nil
		},
		MockUpdateFacilityFn: func(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error {
			return nil
		},
		MockRegisterClientFn: func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
			return clientProfile, nil
		},
		MockListProgramFacilitiesFn: func(ctx context.Context, programID, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
			return facilitiesPage, pagination, nil
		},
		MockRegisterExistingUserAsCaregiverFn: func(ctx context.Context, caregiver *gorm.Caregiver) (*gorm.Caregiver, error) {
			return caregiver, nil
		},
		MockGetUserSurveyFormsFn: func(ctx context.Context, params map[string]interface{}) ([]*gorm.UserSurvey, error) {
			return []*gorm.UserSurvey{
				{
					Base:           gorm.Base{},
					ID:             fhirID,
					Active:         false,
					Link:           uuid.New().String(),
					Title:          "Title",
					Description:    description,
					HasSubmitted:   false,
					OrganisationID: uuid.New().String(),
					UserID:         uuid.New().String(),
				},
			}, nil
		},
		MockGetHealthDiaryEntryByIDFn: func(ctx context.Context, healthDiaryEntryID string) (*gorm.ClientHealthDiaryEntry, error) {
			return &gorm.ClientHealthDiaryEntry{
				ClientHealthDiaryEntryID: new(string),
				Active:                   false,
				Mood:                     "",
				Note:                     "",
				EntryType:                "",
				ShareWithHealthWorker:    false,
				SharedAt:                 &currentTime,
				ClientID:                 "",
				OrganisationID:           "",
			}, nil
		},
		MockGetUserProfileByUsernameFn: func(ctx context.Context, username string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},
		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return pinData, nil
		},

		MockInactivateFacilityFn: func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
			return true, nil
		},
		MockSearchCaregiverUserFn: func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
			return []*gorm.Caregiver{
				{
					ID:              UUID,
					Active:          true,
					CaregiverNumber: "CG001",
					UserID:          UUID,
				},
			}, nil
		},
		MockSearchPlatformCaregiversFn: func(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
			return []*gorm.Caregiver{
				{
					ID:              UUID,
					Active:          true,
					CaregiverNumber: "CG001",
					UserID:          UUID,
				},
			}, nil
		},
		MockGetCaregiverProfileByCaregiverIDFn: func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error) {
			return &gorm.Caregiver{
				ID:              UUID,
				Active:          true,
				CaregiverNumber: "CG001",
				UserID:          UUID,
				UserProfile:     *userProfile,
			}, nil
		},
		MockUpdateProgramFn: func(ctx context.Context, program *gorm.Program, updateData map[string]interface{}) error {
			return nil
		},
		MockListClientsCaregiversFn: func(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
			now := time.Now()
			return []*gorm.CaregiverClient{
					{
						CaregiverID:        uuid.New().String(),
						ClientID:           UUID,
						Active:             true,
						RelationshipType:   enums.CaregiverTypeFather,
						CaregiverConsent:   enums.ConsentStateAccepted,
						CaregiverConsentAt: &now,
						ClientConsent:      enums.ConsentStateAccepted,
						ClientConsentAt:    &now,
						OrganisationID:     UUID,
						AssignedBy:         UUID,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 2,
					TotalPages:  20,
				}, nil
		},
		MockUpdateCaregiverClientFn: func(ctx context.Context, caregiverClient *gorm.CaregiverClient, updates map[string]interface{}) error {
			return nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
			return true, nil
		},
		MockGetStaffServiceRequestsFn: func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
			UUID := uuid.New().String()
			rt := time.Now()
			serviceRequest := &gorm.StaffServiceRequest{
				ID:             &UUID,
				Active:         true,
				RequestType:    "test",
				Request:        "test",
				Status:         "test",
				ResolvedAt:     &rt,
				StaffID:        "test",
				OrganisationID: "test",
				ResolvedByID:   &UUID,
				Meta:           `{"key":"value"}`,
			}
			return []*gorm.StaffServiceRequest{serviceRequest}, nil
		},
		MockGetCurrentTermsFn: func(ctx context.Context) (*gorm.TermsOfService, error) {
			termsID := gofakeit.Number(1, 1000)
			validFrom := time.Now()
			testText := "test"

			validTo := time.Now().AddDate(0, 0, 80)
			terms := &gorm.TermsOfService{
				Base:      gorm.Base{},
				TermsID:   &termsID,
				Text:      &testText,
				ValidFrom: &validFrom,
				ValidTo:   &validTo,
				Active:    false,
			}
			return terms, nil
		},
		MockUpdateUserFn: func(ctx context.Context, user *gorm.User, updateData map[string]interface{}) error {
			return nil
		},
		MockGetUserProfileByUserIDFn: func(ctx context.Context, userID *string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
				Name:   "test",
			}, nil
		},
		MockSaveTemporaryUserPinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
		MockGetFacilitiesWithoutFHIRIDFn: func(ctx context.Context) ([]*gorm.Facility, error) {
			return []*gorm.Facility{facility}, nil
		},
		MockAcceptTermsFn: func(ctx context.Context, userID *string, termsID *int) (bool, error) {
			return true, nil
		},
		MockCreateStaffServiceRequestFn: func(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error {
			return nil
		},
		MockGetStaffProfileByStaffIDFn: func(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
			return &gorm.StaffProfile{
				ID: &UUID,
				UserProfile: gorm.User{
					UserID: &UUID,
				},
				UserID:            UUID,
				Active:            true,
				StaffNumber:       "TEST-001",
				Facilities:        facilities,
				DefaultFacilityID: UUID,
				OrganisationID:    UUID,
			}, nil
		},
		MockSavePinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
		MockRegisterExistingUserAsClientFn: func(ctx context.Context, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
			return client, nil
		},
		MockUpdateServiceRequestsFn: func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error) {
			return true, nil
		},
		MockResolveStaffServiceRequestFn: func(ctx context.Context, staffID, serviceRequestID *string, verificationStatus string) (bool, error) {
			return true, nil
		},
		MockUpdateHealthDiaryFn: func(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
			return nil
		},
		MockGetSecurityQuestionsFn: func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
			sq := ksuid.New().String()
			securityQuestion := &gorm.SecurityQuestion{
				SecurityQuestionID: &sq,
				QuestionStem:       "test",
				Description:        "test",
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			}
			return []*gorm.SecurityQuestion{securityQuestion}, nil
		},
		MockSaveOTPFn: func(ctx context.Context, otpInput *gorm.UserOTP) error {
			return nil
		},
		MockGetSecurityQuestionByIDFn: func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error) {
			return &gorm.SecurityQuestion{
				SecurityQuestionID: &UUID,
				QuestionStem:       "test",
				Description:        "test",
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			}, nil
		},
		MockSaveSecurityQuestionResponseFn: func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error {
			return nil
		},
		MockAddCaregiverToClientFn: func(ctx context.Context, clientCaregiver *gorm.CaregiverClient) error {
			return nil
		},
		MockGetSecurityQuestionResponseFn: func(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error) {
			return &gorm.SecurityQuestionResponse{
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
		MockGetClientProfileFn: func(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
			return clientProfile, nil
		},
		MockGetStaffProfileFn: func(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockListCommunitiesFn: func(ctx context.Context, programID, organisationID string) ([]*gorm.Community, error) {
			return []*gorm.Community{
				{
					ID:     uuid.NewString(),
					RoomID: uuid.NewString(),
				},
			}, nil
		},
		MockSearchStaffProfileFn: func(ctx context.Context, staffNumber string, programID *string) ([]*gorm.StaffProfile, error) {
			return []*gorm.StaffProfile{staff}, nil
		},
		MockCheckUserHasPinFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockCompleteOnboardingTourFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGetOTPFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error) {
			return &gorm.UserOTP{
				OTP: "1234",
			}, nil
		},
		MockGetUserSecurityQuestionsResponsesFn: func(ctx context.Context, userID, flavour string) ([]*gorm.SecurityQuestionResponse, error) {
			return []*gorm.SecurityQuestionResponse{
				{
					ResponseID: "1234",
					QuestionID: "1234",
					Active:     true,
					Response:   "Yes",
					IsCorrect:  true,
				},
			}, nil
		},
		MockInvalidatePINFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetContactByUserIDFn: func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error) {
			return &gorm.Contact{
				ID:      UUID,
				UserID:  userID,
				Type:    "PHONE",
				Value:   phoneContact,
				Active:  true,
				OptedIn: true,
			}, nil
		},
		MockFindContactsFn: func(ctx context.Context, contactType, contactValue string) ([]*gorm.Contact, error) {
			return []*gorm.Contact{
				{
					ID:      UUID,
					UserID:  &UUID,
					Type:    "PHONE",
					Value:   phoneContact,
					Active:  true,
					OptedIn: true,
				},
			}, nil
		},
		MockUpdateIsCorrectSecurityQuestionResponseFn: func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
			return true, nil
		},
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) (*gorm.ClientHealthDiaryEntry, error) {
			return &gorm.ClientHealthDiaryEntry{
				ClientHealthDiaryEntryID: &UUID,
				Active:                   true,
				Mood:                     "VERY_SAD",
				Note:                     gofakeit.BS(),
				EntryType:                "HOME_PAGE_HEALTH_DIARY_ENTRY",
				ShareWithHealthWorker:    false,
				SharedAt:                 &nowTime,
				ClientID:                 uuid.NewString(),
				OrganisationID:           uuid.NewString(),
			}, nil
		},
		MockCreateServiceRequestFn: func(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error {
			return nil
		},
		MockCanRecordHeathDiaryFn: func(ctx context.Context, clientID string) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryQuoteFn: func(ctx context.Context, limit int) ([]*gorm.ClientHealthDiaryQuote, error) {
			return []*gorm.ClientHealthDiaryQuote{
				{
					Quote:  "Quote",
					Author: "Author",
				},
			}, nil
		},
		MockGetClientHealthDiaryEntriesFn: func(ctx context.Context, params map[string]interface{}) ([]*gorm.ClientHealthDiaryEntry, error) {
			return []*gorm.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockGetClientPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error) {
			return &domain.ServiceRequestsCount{
				Total: 0,
				RequestsTypeCount: []*domain.RequestTypeCount{
					{
						RequestType: enums.ServiceRequestTypeRedFlag,
						Total:       0,
					},
				},
			}, nil
		},
		MockUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
			return nil
		},
		MockInProgressByFn: func(ctx context.Context, requestID, staffID string) (bool, error) {
			return true, nil
		},
		MockGetClientProfileByClientIDFn: func(ctx context.Context, clientID string) (*gorm.Client, error) {
			return clientProfile, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string) ([]*gorm.ClientServiceRequest, error) {
			return serviceRequests, nil
		},
		MockCreateCommunityFn: func(ctx context.Context, community *gorm.Community) (*gorm.Community, error) {
			return &gorm.Community{
				Base:           gorm.Base{},
				ID:             UUID,
				Name:           name,
				Description:    description,
				Active:         false,
				Gender:         []string{"test"},
				ClientTypes:    []string{"test"},
				OrganisationID: uuid.New().String(),
			}, nil
		},
		MockCheckIfUsernameExistsFn: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
		MockGetCommunityByIDFn: func(ctx context.Context, communityID string) (*gorm.Community, error) {
			return &gorm.Community{
				ID:             uuid.New().String(),
				Name:           name,
				Description:    description,
				Active:         false,
				MinimumAge:     0,
				MaximumAge:     0,
				Gender:         []string{"MALE"},
				ClientTypes:    []string{"PMTCT"},
				OrganisationID: uuid.New().String(),
			}, nil
		},
		MockGetClientsInAFacilityFn: func(ctx context.Context, facilityID string) ([]*gorm.Client, error) {
			return []*gorm.Client{
				clientProfile,
			}, nil
		},
		MockGetRecentHealthDiaryEntriesFn: func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
			return []*gorm.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockCheckFacilityExistsByIdentifier: func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
			return true, nil
		},
		MockCheckIdentifierExists: func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
			return true, nil
		},
		MockGetClientsByParams: func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error) {
			return []*gorm.Client{clientProfile}, nil
		},
		MockGetClientIdentifiers: func(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
			return []*gorm.Identifier{
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
		MockGetServiceRequestsForKenyaEMRFn: func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facility := uuid.New().String()
			requestID := uuid.New().String()
			serviceReq := &gorm.ClientServiceRequest{
				ID:             &requestID,
				Active:         true,
				RequestType:    "TYPE",
				Request:        "REQUEST",
				Status:         "PENDING",
				InProgressAt:   &currentTime,
				ResolvedAt:     &currentTime,
				ClientID:       uuid.New().String(),
				InProgressByID: &staffID,
				OrganisationID: "",
				ResolvedByID:   &staffID,
				FacilityID:     facility,
				Meta:           `{"key":"value"}`,
			}
			return []*gorm.ClientServiceRequest{serviceReq}, nil
		},
		MockCreateAppointment: func(ctx context.Context, appointment *gorm.Appointment) error {
			return nil
		},
		MockListAppointments: func(ctx context.Context, params *gorm.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error) {
			date := time.Now().Add(time.Duration(100))
			return []*gorm.Appointment{
				{
					ID:             gofakeit.UUID(),
					OrganisationID: gofakeit.UUID(),
					Active:         true,
					ExternalID:     strconv.Itoa(gofakeit.Number(0, 1000)),
					ClientID:       gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
					Reason:         "Knocked up",
					Date:           date,
				},
			}, &domain.Pagination{Limit: 10, CurrentPage: 1}, nil
		},
		MockListNotificationsFn: func(ctx context.Context, params *gorm.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Notification, *domain.Pagination, error) {
			return []*gorm.Notification{
				{
					Title:      "A notification",
					Body:       "This is what it's about",
					Type:       "TELECONSULT",
					IsRead:     false,
					UserID:     &UUID,
					FacilityID: &UUID,
				},
			}, &domain.Pagination{Limit: 10, CurrentPage: 1}, nil
		},
		MockListAvailableNotificationTypesFn: func(ctx context.Context, params *gorm.Notification) ([]enums.NotificationType, error) {
			return []enums.NotificationType{enums.NotificationTypeAppointment}, nil
		},
		MockUpdateAppointmentFn: func(ctx context.Context, appointment *gorm.Appointment, updateData map[string]interface{}) (*gorm.Appointment, error) {
			return appointment, nil
		},
		MockGetStaffUserProgramsFn: func(ctx context.Context, userID string) ([]*gorm.Program, error) {
			return []*gorm.Program{
				{
					ID:             UUID,
					Active:         true,
					Name:           "Fake Program",
					OrganisationID: UUID,
				},
			}, nil
		},
		MockGetClientUserProgramsFn: func(ctx context.Context, userID string) ([]*gorm.Program, error) {
			return []*gorm.Program{
				{
					ID:             UUID,
					Active:         true,
					Name:           "Fake Program",
					OrganisationID: UUID,
				},
			}, nil
		},
		MockGetProgramClientProfileByIdentifierFn: func(ctx context.Context, programID, identifierType, value string) (*gorm.Client, error) {
			return clientProfile, nil
		},

		MockGetClientProfilesByIdentifierFn: func(ctx context.Context, identifierType, value string) ([]*gorm.Client, error) {
			return []*gorm.Client{clientProfile}, nil
		},
		MockSearchClientProfileFn: func(ctx context.Context, searchParameter string) ([]*gorm.Client, error) {
			return []*gorm.Client{clientProfile}, nil
		},
		MockCheckIfClientHasUnresolvedServiceRequestsFn: func(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
			return true, nil
		},
		MockUpdateUserPinChangeRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockUpdateUserPinUpdateRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockUpdateClientFn: func(ctx context.Context, client *gorm.Client, updates map[string]interface{}) (*gorm.Client, error) {
			return clientProfile, nil
		},
		MockGetUserProfileByStaffIDFn: func(ctx context.Context, staffID string) (*gorm.User, error) {
			return userProfile, nil
		},
		MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn: func(ctx context.Context, userID string, failCount int) error {
			return nil
		},
		MockGetClientServiceRequestByIDFn: func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facility := uuid.New().String()
			requestID := uuid.New().String()
			serviceReq := &gorm.ClientServiceRequest{
				ID:             &requestID,
				Active:         true,
				RequestType:    "TYPE",
				Request:        "REQUEST",
				Status:         "PENDING",
				InProgressAt:   &currentTime,
				ResolvedAt:     &currentTime,
				ClientID:       uuid.New().String(),
				InProgressByID: &staffID,
				OrganisationID: "",
				ResolvedByID:   &staffID,
				FacilityID:     facility,
				Meta:           `{"meta": "data"}`,
			}
			return serviceReq, nil
		},
		MockGetAppointmentServiceRequestsFn: func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
			meta := map[string]interface{}{
				"appointmentID":  uuid.New().String(),
				"rescheduleTime": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			}

			bs, err := json.Marshal(meta)
			if err != nil {
				return nil, err
			}
			return []*gorm.ClientServiceRequest{
				{
					ID:             &UUID,
					Active:         true,
					RequestType:    "TYPE",
					Request:        "REQUEST",
					Status:         "PENDING",
					InProgressAt:   nil,
					ResolvedAt:     nil,
					ClientID:       uuid.New().String(),
					InProgressByID: &UUID,
					OrganisationID: "",
					ResolvedByID:   &UUID,
					FacilityID:     uuid.New().String(),
					Meta:           string(bs),
				},
			}, nil
		},
		MockCheckAppointmentExistsByExternalIDFn: func(ctx context.Context, externalID string) (bool, error) {
			return true, nil
		},
		MockSearchClientServiceRequestsFn: func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
			return []*gorm.ClientServiceRequest{
				{
					ID:             &UUID,
					Active:         true,
					RequestType:    "RED_FLAG",
					Request:        "REQUEST",
					Status:         "PENDING",
					InProgressAt:   nil,
					ResolvedAt:     nil,
					ClientID:       uuid.New().String(),
					InProgressByID: &UUID,
					OrganisationID: "",
					ResolvedByID:   &UUID,
					FacilityID:     uuid.New().String(),
					Meta:           `{"meta": "data"}`,
				},
			}, nil
		},
		MockSearchStaffServiceRequestsFn: func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
			return []*gorm.StaffServiceRequest{
				{
					ID:             &UUID,
					Active:         true,
					RequestType:    "RED_FLAG",
					Request:        "REQUEST",
					Status:         "PENDING",
					ResolvedAt:     nil,
					StaffID:        uuid.New().String(),
					OrganisationID: "",
					ResolvedByID:   &UUID,
					Meta:           `{"meta": "data"}`,
				},
			}, nil
		},
		MockGetClientServiceRequestsFn: func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*gorm.ClientServiceRequest, error) {
			return []*gorm.ClientServiceRequest{
				{
					ID:             &UUID,
					Active:         true,
					RequestType:    enums.ServiceRequestTypeRedFlag.String(),
					Request:        "REQUEST",
					Status:         string(enums.ServiceRequestStatusResolved),
					InProgressAt:   nil,
					ResolvedAt:     nil,
					ClientID:       uuid.New().String(),
					InProgressByID: &UUID,
					OrganisationID: "",
					ResolvedByID:   &UUID,
					FacilityID:     uuid.New().String(),
					Meta:           fmt.Sprintf(`{"question_id":"%s"}`, "screening_tool_question_id"),
				},
			}, nil
		},
		MockGetClientScreeningToolServiceRequestByToolTypeFn: func(ctx context.Context, clientID, toolType, status string) (*gorm.ClientServiceRequest, error) {
			return &gorm.ClientServiceRequest{
				ID:             &UUID,
				Active:         true,
				RequestType:    enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
				Request:        "REQUEST",
				Status:         string(enums.ServiceRequestStatusPending),
				InProgressAt:   nil,
				ResolvedAt:     nil,
				ClientID:       uuid.New().String(),
				InProgressByID: &UUID,
				OrganisationID: "",
				ResolvedByID:   &UUID,
				FacilityID:     uuid.New().String(),
				Meta:           fmt.Sprintf(`{"question_id":"%s"}`, "screening_tool_question_id"),
			}, nil
		},
		MockCheckIfStaffHasUnresolvedServiceRequestsFn: func(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
			return false, nil
		},
		MockDeleteStaffProfileFn: func(ctx context.Context, staffID string) error {
			return nil
		},
		MockDeleteClientProfileFn: func(ctx context.Context, clientID string, userID *string) error {
			return nil
		},
		MockGetClientsByFilterParamsFn: func(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*gorm.Client, error) {
			return []*gorm.Client{
				clientProfile,
			}, nil
		},
		MockCreateUserSurveyFn: func(ctx context.Context, userSurvey []*gorm.UserSurvey) error {
			return nil
		},
		MockUpdateUserSurveysFn: func(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error {
			return nil
		},
		MockUpdateClientServiceRequestFn: func(ctx context.Context, clientServiceRequest *gorm.ClientServiceRequest, updateData map[string]interface{}) error {
			return nil
		},
		MockDeleteCommunityFn: func(ctx context.Context, communityID string) error {
			return nil
		},
		MockCreateQuestionnaireFn: func(ctx context.Context, input *gorm.Questionnaire) error {
			return nil
		},
		MockCreateScreeningToolFn: func(ctx context.Context, input *gorm.ScreeningTool) error {
			return nil
		},
		MockCreateQuestionFn: func(ctx context.Context, input *gorm.Question) error {
			return nil
		},
		MockCreateQuestionChoiceFn: func(ctx context.Context, input *gorm.QuestionInputChoice) error {
			return nil
		},
		MockGetScreeningToolByIDFn: func(ctx context.Context, toolID string) (*gorm.ScreeningTool, error) {
			return &gorm.ScreeningTool{
				ID:              UUID,
				Active:          true,
				QuestionnaireID: UUID,
				Threshold:       1,
				ClientTypes:     []string{enums.ClientTypeHighRisk.String()},
				Genders:         []string{enumutils.GenderMale.String()},
				MinimumAge:      18,
				MaximumAge:      25,
			}, nil
		},
		MockGetQuestionnaireByIDFn: func(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error) {
			return &gorm.Questionnaire{
				ID:          UUID,
				Active:      true,
				Name:        name,
				Description: description,
			}, nil

		},
		MockGetQuestionsByQuestionnaireIDFn: func(ctx context.Context, questionnaireID string) ([]*gorm.Question, error) {
			return []*gorm.Question{
				{
					ID:                UUID,
					Active:            true,
					QuestionnaireID:   UUID,
					Text:              gofakeit.BS(),
					QuestionType:      string(enums.QuestionTypeOpenEnded),
					ResponseValueType: string(enums.QuestionResponseValueTypeString),
					SelectMultiple:    false,
					Required:          true,
					Sequence:          0,
				},
			}, nil
		},
		MockGetQuestionInputChoicesByQuestionIDFn: func(ctx context.Context, questionID string) ([]*gorm.QuestionInputChoice, error) {
			return []*gorm.QuestionInputChoice{
				{
					ID:         UUID,
					Active:     false,
					QuestionID: UUID,
					Choice:     "0",
					Value:      "yes",
					Score:      0,
				},
			}, nil

		},
		MockCreateScreeningToolResponseFn: func(ctx context.Context, screeningToolResponse *gorm.ScreeningToolResponse, screeningToolQuestionResponses []*gorm.ScreeningToolQuestionResponse) (*string, error) {
			return &UUID, nil
		},
		MockGetSurveysWithServiceRequestsFn: func(ctx context.Context, facilityID, programID string) ([]*gorm.UserSurvey, error) {
			return []*gorm.UserSurvey{
				{
					Base: gorm.Base{
						UpdatedAt: time.Now(),
					},

					ID:             UUID,
					Active:         true,
					Link:           "https://www.google.com",
					Title:          "Test",
					Description:    description,
					HasSubmitted:   true,
					FormID:         "1",
					ProjectID:      ID,
					LinkID:         ID,
					Token:          "",
					SubmittedAt:    &time.Time{},
					UserID:         UUID,
					OrganisationID: UUID,
				},
			}, nil
		},
		MockGetScreeningToolServiceRequestOfRespondentsFn: func(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
			nextPage := 2
			return []*gorm.ClientServiceRequest{
					{
						ID:             &UUID,
						Active:         true,
						RequestType:    enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
						Request:        "REQUEST",
						Status:         string(enums.ServiceRequestStatusPending),
						InProgressAt:   nil,
						ResolvedAt:     nil,
						ClientID:       uuid.New().String(),
						InProgressByID: &UUID,
						OrganisationID: "",
						ResolvedByID:   &UUID,
						FacilityID:     uuid.New().String(),
						Meta:           fmt.Sprintf(`{"response_id":"%s"}`, "screening_tool_response_id"),
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
		MockGetScreeningToolResponseByIDFn: func(ctx context.Context, id string) (*gorm.ScreeningToolResponse, error) {
			return &gorm.ScreeningToolResponse{
				ID:              UUID,
				Active:          true,
				ScreeningToolID: UUID,
				FacilityID:      id,
				ClientID:        uuid.New().String(),
				AggregateScore:  3,
			}, nil
		},
		MockGetScreeningToolQuestionResponsesByResponseIDFn: func(ctx context.Context, responseID string) ([]*gorm.ScreeningToolQuestionResponse, error) {
			return []*gorm.ScreeningToolQuestionResponse{
				{
					ID:                      UUID,
					Active:                  true,
					ScreeningToolResponseID: UUID,
					QuestionID:              UUID,
					Response:                responseID,
					Score:                   3,
				},
			}, nil
		},
		MockGetStaffFacilitiesFn: func(ctx context.Context, staffFacility gorm.StaffFacilities, pagination *domain.Pagination) ([]*gorm.StaffFacilities, *domain.Pagination, error) {
			return []*gorm.StaffFacilities{
					{
						ID:         ID,
						StaffID:    &UUID,
						FacilityID: &UUID,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 2,
				}, nil
		},
		MockGetClientFacilitiesFn: func(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error) {
			return []*gorm.ClientFacilities{
					{
						ID:         ID,
						ClientID:   &UUID,
						FacilityID: &UUID,
					},
				}, &domain.Pagination{
					Limit:       10,
					CurrentPage: 2,
				}, nil
		},
		MockUpdateStaffFn: func(ctx context.Context, staff *gorm.StaffProfile, updates map[string]interface{}) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockAddFacilityToProgramFn: func(ctx context.Context, programID string, facilityID []string) error {
			return nil
		},
		MockAddFacilitiesToStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) error {
			return nil
		},
		MockUpdateClientIdentifierFn: func(ctx context.Context, clientID, identifierType, identifierValue, programID string) error {
			return nil
		},
		MockAddFacilitiesToClientProfileFn: func(ctx context.Context, clientID string, facilities []string) error {
			return nil
		},
		MockGetNotificationsCountFn: func(ctx context.Context, notification gorm.Notification) (int, error) {
			return 1, nil
		},
		MockRegisterExistingUserAsStaffFn: func(ctx context.Context, identifier *gorm.Identifier, staff *gorm.StaffProfile) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockRemoveFacilitiesFromClientProfileFn: func(ctx context.Context, clientID string, facilities []string) error {
			return nil
		},
		MockUpdateUserContactFn: func(ctx context.Context, userContact *gorm.Contact, updates map[string]interface{}) error {
			return nil
		},
		MockRemoveFacilitiesFromStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) error {
			return nil
		},
		MockGetCaregiverManagedClientsFn: func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
			return []*gorm.CaregiverClient{&caregiversClient}, paginationOutput, nil
		},
		MockListOrganisationsFn: func(ctx context.Context, pagination *domain.Pagination) ([]*gorm.Organisation, *domain.Pagination, error) {
			return []*gorm.Organisation{
				{
					ID:              &UUID,
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
			}, paginationOutput, nil
		},
		MockGetCaregiversClientFn: func(ctx context.Context, caregiverClient gorm.CaregiverClient) ([]*gorm.CaregiverClient, error) {
			return []*gorm.CaregiverClient{
				{
					CaregiverID:        uuid.NewString(),
					ClientID:           uuid.NewString(),
					Active:             true,
					RelationshipType:   enums.CaregiverTypeHealthCareProfessional,
					CaregiverConsent:   enums.ConsentStateAccepted,
					CaregiverConsentAt: &nowTime,
					ClientConsent:      enums.ConsentStateAccepted,
					ClientConsentAt:    &nowTime,
					AssignedBy:         uuid.NewString(),
				},
			}, nil
		},
		MockGetOrganisationFn: func(ctx context.Context, id string) (*gorm.Organisation, error) {
			return &gorm.Organisation{
				ID:             new(string),
				Active:         true,
				Code:           gofakeit.SSN(),
				Name:           gofakeit.Company(),
				Description:    description,
				EmailAddress:   gofakeit.Email(),
				PhoneNumber:    gofakeit.Phone(),
				DefaultCountry: gofakeit.Country(),
			}, nil
		},
		MockCreateOrganisationFn: func(ctx context.Context, organization *gorm.Organisation) (*gorm.Organisation, error) {
			return &gorm.Organisation{
				ID:             &UUID,
				Active:         true,
				Code:           gofakeit.SSN(),
				Name:           gofakeit.Company(),
				Description:    description,
				EmailAddress:   gofakeit.Email(),
				PhoneNumber:    gofakeit.Phone(),
				DefaultCountry: gofakeit.Country(),
			}, nil
		},
		MockCreateProgramFn: func(ctx context.Context, program *gorm.Program) (*gorm.Program, error) {
			return program, nil
		},
		MockSearchOrganisationsFn: func(ctx context.Context, searchParameter string) ([]*gorm.Organisation, error) {
			return []*gorm.Organisation{
				{
					ID:             &UUID,
					Active:         true,
					Code:           gofakeit.SSN(),
					Name:           gofakeit.Company(),
					Description:    description,
					EmailAddress:   gofakeit.Email(),
					PhoneNumber:    gofakeit.Phone(),
					DefaultCountry: gofakeit.Country(),
				},
			}, nil
		},
		MockCheckOrganisationExistsFn: func(ctx context.Context, organisationID string) (bool, error) {
			return true, nil
		},
		MockCheckIfProgramNameExistsFn: func(ctx context.Context, organisationID string, programName string) (bool, error) {
			return false, nil
		},
		MockGetProgramFacilitiesFn: func(ctx context.Context, programID string) ([]*gorm.ProgramFacility, error) {
			return []*gorm.ProgramFacility{
				{
					ID:         ID,
					ProgramID:  programID,
					FacilityID: UUID,
				},
			}, nil
		},
		MockListProgramsFn: func(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
			return []*gorm.Program{&program}, pagination, nil
		},
		MockCheckIfSuperUserExistsFn: func(ctx context.Context) (bool, error) {
			return true, nil
		},
		MockGetCaregiverProfileByUserIDFn: func(ctx context.Context, userID string, organisationID string) (*gorm.Caregiver, error) {
			return &gorm.Caregiver{
				ID:              UUID,
				Active:          false,
				CaregiverNumber: "c11",
				OrganisationID:  UUID,
				UserID:          userID,
				UserProfile:     *userProfile,
				CurrentClient:   &UUID,
				CurrentFacility: &UUID,
			}, nil
		},
		MockUpdateCaregiverFn: func(ctx context.Context, caregiver *gorm.Caregiver, updates map[string]interface{}) error {
			return nil
		},
		MockCreateFacilitiesFn: func(ctx context.Context, facilities []*gorm.Facility) ([]*gorm.Facility, error) {
			return []*gorm.Facility{facility}, nil
		},
		MockCreateSecurityQuestionsFn: func(ctx context.Context, securityQuestions []*gorm.SecurityQuestion) ([]*gorm.SecurityQuestion, error) {
			return []*gorm.SecurityQuestion{{
				SecurityQuestionID: &UUID,
				QuestionStem:       gofakeit.Question(),
				Description:        description,
				ResponseType:       enums.SecurityQuestionResponseTypeText,
				Flavour:            feedlib.FlavourPro,
				Active:             true,
				Sequence:           new(int),
			}}, nil
		},
		MockCreateTermsOfServiceFn: func(ctx context.Context, termsOfService *gorm.TermsOfService) (*gorm.TermsOfService, error) {
			return &gorm.TermsOfService{
				TermsID:   &ID,
				Text:      &name,
				ValidFrom: &time.Time{},
				ValidTo:   &time.Time{},
				Active:    true,
			}, nil
		},
		MockCheckPhoneExistsFn: func(ctx context.Context, phone string) (bool, error) {
			return false, nil
		},
		MockGetStaffServiceRequestByIDFn: func(ctx context.Context, serviceRequestID string) (*gorm.StaffServiceRequest, error) {
			return &gorm.StaffServiceRequest{
				ID:                &UUID,
				Active:            true,
				RequestType:       enums.ServiceRequestTypeRedFlag.String(),
				Request:           gofakeit.BS(),
				Status:            string(enums.ServiceRequestStatusPending),
				ResolvedAt:        &nowTime,
				Meta:              "",
				StaffID:           UUID,
				OrganisationID:    UUID,
				ResolvedByID:      &UUID,
				DefaultFacilityID: &UUID,
				ProgramID:         "",
			}, nil
		},
		MockCreateOauthClientJWT: func(ctx context.Context, jwt *gorm.OauthClientJWT) error {
			return nil
		},
		MockCreateOauthClient: func(ctx context.Context, client *gorm.OauthClient) error {
			return nil
		},
		MockGetClientJWT: func(ctx context.Context, jti string) (*gorm.OauthClientJWT, error) {
			return &gorm.OauthClientJWT{}, nil
		},
		MockGetOauthClient: func(ctx context.Context, id string) (*gorm.OauthClient, error) {
			return &gorm.OauthClient{}, nil
		},
		MockGetValidClientJWT: func(ctx context.Context, jti string) (*gorm.OauthClientJWT, error) {
			return &gorm.OauthClientJWT{}, nil
		},
		MockCreateOrUpdateSessionFn: func(ctx context.Context, session *gorm.Session) error {
			return nil
		},
		MockCreateAuthorizationCodeFn: func(ctx context.Context, code *gorm.AuthorizationCode) error {
			return nil
		},
		MockGetAuthorizationCodeFn: func(ctx context.Context, code string) (*gorm.AuthorizationCode, error) {
			return &gorm.AuthorizationCode{ID: gofakeit.UUID()}, nil
		},
		MockUpdateAuthorizationCodeFn: func(ctx context.Context, code *gorm.AuthorizationCode, updateData map[string]interface{}) error {
			return nil
		},
		MockCreateAccessTokenFn: func(ctx context.Context, token *gorm.AccessToken) error {
			return nil
		},
		MockCreateRefreshTokenFn: func(ctx context.Context, token *gorm.RefreshToken) error {
			return nil
		},
		MockDeleteAccessTokenFn: func(ctx context.Context, signature string) error {
			return nil
		},
		MockDeleteRefreshTokenFn: func(ctx context.Context, signature string) error {
			return nil
		},
		MockGetAccessTokenFn: func(ctx context.Context, token gorm.AccessToken) (*gorm.AccessToken, error) {
			return &gorm.AccessToken{}, nil
		},
		MockGetRefreshTokenFn: func(ctx context.Context, token gorm.RefreshToken) (*gorm.RefreshToken, error) {
			return &gorm.RefreshToken{}, nil
		},
		MockUpdateAccessTokenFn: func(ctx context.Context, code *gorm.AccessToken, updateData map[string]interface{}) error {
			return nil
		},
		MockUpdateRefreshTokenFn: func(ctx context.Context, code *gorm.RefreshToken, updateData map[string]interface{}) error {
			return nil
		},
		MockCheckIfClientHasPendingSurveyServiceRequestFn: func(ctx context.Context, clientID string, projectID int, formID string) (bool, error) {
			return false, nil
		},
		MockGetUserProfileByPushTokenFn: func(ctx context.Context, pushToken string) (*gorm.User, error) {
			return userProfile, nil
		},
		MockCheckStaffExistsInProgramFn: func(ctx context.Context, userID, programID string) (bool, error) {
			return false, nil
		},
		MockCheckIfFacilityExistsInProgramFn: func(ctx context.Context, programID, facilityID string) (bool, error) {
			return true, nil
		},
		MockGetStaffIdentifiersFn: func(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error) {
			return []*gorm.Identifier{
				{
					ID:                  uuid.New().String(),
					Type:                enums.UserIdentifierTypeNationalID.String(),
					Value:               "123456",
					Use:                 "OFFICIAL",
					Description:         description,
					ValidFrom:           time.Now(),
					ValidTo:             time.Now(),
					IsPrimaryIdentifier: false,
				},
			}, nil
		},
		MockCheckIfClientExistsInProgramFn: func(ctx context.Context, userID, programID string) (bool, error) {
			return false, nil
		},
		MockGetUserClientProfilesFn: func(ctx context.Context, userID string) ([]*gorm.Client, error) {
			return []*gorm.Client{
				clientProfile,
			}, nil
		},
		MockGetUserStaffProfilesFn: func(ctx context.Context, userID string) ([]*gorm.StaffProfile, error) {
			return []*gorm.StaffProfile{
				staff,
			}, nil
		},
	}
}

// DeleteStaffProfile mocks the implementation of deleting a staff
func (gm *GormMock) DeleteStaffProfile(ctx context.Context, staffID string) error {
	return gm.MockDeleteStaffProfileFn(ctx, staffID)
}

// GetOrganisation retrieves an organisation using the provided id
func (gm *GormMock) GetOrganisation(ctx context.Context, id string) (*gorm.Organisation, error) {
	return gm.MockGetOrganisationFn(ctx, id)
}

// DeleteClientProfile mocks the implementation of deleting a client user
func (gm *GormMock) DeleteClientProfile(ctx context.Context, clientID string, userID *string) error {
	return gm.MockDeleteClientProfileFn(ctx, clientID, userID)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByIdentifier mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByIdentifier(ctx context.Context, identifier *gorm.FacilityIdentifier, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityByIdentifierFn(ctx, identifier, isActive)
}

// RetrieveFacilityIdentifiersByFacilityID mocks the implementation of getting facility identifier by facility id
func (gm *GormMock) RetrieveFacilityIdentifiersByFacilityID(ctx context.Context, facilityID *string) ([]*gorm.FacilityIdentifier, error) {
	return gm.MockRetrieveFacilityIdentifiersByFacilityIDFn(ctx, facilityID)
}

// CheckStaffExists checks if there is a staff profile that exists for a user
func (gm *GormMock) CheckStaffExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckStaffExistsFn(ctx, userID)
}

// CheckClientExists checks if there is a client profile that exists for a user
func (gm *GormMock) CheckClientExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckClientExistsFn(ctx, userID)
}

// CheckCaregiverExists checks if there is a caregiver profile that exists for a user
func (gm *GormMock) CheckCaregiverExists(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckCaregiverExistsFn(ctx, userID)
}

// UpdateUserSurveys mocks the implementation of `gorm's` UpdateUserSurveys method.
func (gm *GormMock) UpdateUserSurveys(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error {
	return gm.MockUpdateUserSurveysFn(ctx, survey, updateData)
}

// ListFacilities mocks the implementation of `gorm's` ListFacilities method.
func (gm *GormMock) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
	return gm.MockListFacilitiesFn(ctx, searchTerm, filter, pagination)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, identifier)
}

// ListProgramFacilities mocks the implementation of  ListProgramFacilities method.
func (gm *GormMock) ListProgramFacilities(ctx context.Context, programID, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Facility, *domain.Pagination, error) {
	return gm.MockListProgramFacilitiesFn(ctx, programID, searchTerm, filter, pagination)
}

// GetUserProfileByUsername retrieves a user using their username
func (gm *GormMock) GetUserProfileByUsername(ctx context.Context, username string) (*gorm.User, error) {
	return gm.MockGetUserProfileByUsernameFn(ctx, username)
}

// GetUserProfileByPhoneNumber mocks the implementation of retrieving a user profile by phonenumber
func (gm *GormMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*gorm.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// GetUserPINByUserID mocks the implementation of retrieving a user pin by user ID
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *GormMock) InactivateFacility(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, identifier)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (gm *GormMock) ReactivateFacility(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
	return gm.MockReactivateFacilityFn(ctx, identifier)
}

// GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *GormMock) GetCurrentTerms(ctx context.Context) (*gorm.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}

// GetUserProfileByUserID mocks the implementation of retrieving a user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, userID *string) (*gorm.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, userID)
}

// GetCaregiverByUserID returns the caregiver record of the provided user ID
func (gm *GormMock) GetCaregiverByUserID(ctx context.Context, userID string) (*gorm.Caregiver, error) {
	return gm.MockGetCaregiverByUserIDFn(ctx, userID)
}

// SaveTemporaryUserPin mocks the implementation of saving a temporary user pin
func (gm *GormMock) SaveTemporaryUserPin(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.MockSaveTemporaryUserPinFn(ctx, pinData)
}

// AcceptTerms mocks the implementation of accept current terms of service
func (gm *GormMock) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	return gm.MockAcceptTermsFn(ctx, userID, termsID)
}

// SavePin mocks the implementation of saving the pin to the database
func (gm *GormMock) SavePin(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.MockSavePinFn(ctx, pinData)
}

// GetSecurityQuestions mocks the implementation of getting all the security questions.
func (gm *GormMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionsFn(ctx, flavour)
}

// SaveOTP mocks the implementation for saving an OTP
func (gm *GormMock) SaveOTP(ctx context.Context, otpInput *gorm.UserOTP) error {
	return gm.MockSaveOTPFn(ctx, otpInput)
}

// GetSecurityQuestionByID mocks the implementation of getting a security question by ID
func (gm *GormMock) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionByIDFn(ctx, securityQuestionID)
}

// SaveSecurityQuestionResponse mocks the implementation of saving a security question response
func (gm *GormMock) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error {
	return gm.MockSaveSecurityQuestionResponseFn(ctx, securityQuestionResponse)
}

// GetSecurityQuestionResponse mocks the get security question implementation
func (gm *GormMock) GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error) {
	return gm.MockGetSecurityQuestionResponseFn(ctx, questionID, userID)
}

// CheckIfPhoneNumberExists mock the implementation of checking the existence of phone number
func (gm *GormMock) CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckIfPhoneNumberExistsFn(ctx, phone, isOptedIn, flavour)
}

// GetStaffUserPrograms retrieves all programs associated with a staff user
func (gm *GormMock) GetStaffUserPrograms(ctx context.Context, userID string) ([]*gorm.Program, error) {
	return gm.MockGetStaffUserProgramsFn(ctx, userID)
}

// GetClientUserPrograms retrieves all programs associated with a client user
func (gm *GormMock) GetClientUserPrograms(ctx context.Context, userID string) ([]*gorm.Program, error) {
	return gm.MockGetClientUserProgramsFn(ctx, userID)
}

// VerifyOTP mocks the implementation of verify otp
func (gm *GormMock) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return gm.MockVerifyOTPFn(ctx, payload)
}

// GetClientProfile mocks the method for fetching a client profile using the user ID
func (gm *GormMock) GetClientProfile(ctx context.Context, userID string, programID string) (*gorm.Client, error) {
	return gm.MockGetClientProfileFn(ctx, userID, programID)
}

// GetStaffProfile mocks the method for fetching a staff profile using the user ID
func (gm *GormMock) GetStaffProfile(ctx context.Context, userID string, programID string) (*gorm.StaffProfile, error) {
	return gm.MockGetStaffProfileFn(ctx, userID, programID)
}

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *GormMock) CheckUserHasPin(ctx context.Context, userID string) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID)
}

// CompleteOnboardingTour mocks the implementation for updating a user's pin change required state
func (gm *GormMock) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCompleteOnboardingTourFn(ctx, userID, flavour)
}

// GetOTP fetches the OTP for the given phone number
func (gm *GormMock) GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error) {
	return gm.MockGetOTPFn(ctx, phoneNumber, flavour)
}

// GetUserSecurityQuestionsResponses mocks the implementation of getting the user's responded security questions
func (gm *GormMock) GetUserSecurityQuestionsResponses(ctx context.Context, userID, flavour string) ([]*gorm.SecurityQuestionResponse, error) {
	return gm.MockGetUserSecurityQuestionsResponsesFn(ctx, userID, flavour)
}

// InvalidatePIN mocks the implementation of invalidating the pin
func (gm *GormMock) InvalidatePIN(ctx context.Context, userID string) (bool, error) {
	return gm.MockInvalidatePINFn(ctx, userID)
}

// GetContactByUserID mocks the implementation of retrieving a contact by user ID
func (gm *GormMock) GetContactByUserID(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error) {
	return gm.MockGetContactByUserIDFn(ctx, userID, contactType)
}

// UpdateIsCorrectSecurityQuestionResponse updates the is_correct security question response
func (gm *GormMock) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	return gm.MockUpdateIsCorrectSecurityQuestionResponseFn(ctx, userID, isCorrectSecurityQuestionResponse)
}

// CreateHealthDiaryEntry mocks the method for creating a health diary entry
func (gm *GormMock) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) (*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockCreateHealthDiaryEntryFn(ctx, healthDiaryInput)
}

// CreateServiceRequest mocks creating a service request method
func (gm *GormMock) CreateServiceRequest(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error {
	return gm.MockCreateServiceRequestFn(ctx, serviceRequestInput)
}

// CanRecordHeathDiary mocks the implementation of checking if a user can record a health diary
func (gm *GormMock) CanRecordHeathDiary(ctx context.Context, userID string) (bool, error) {
	return gm.MockCanRecordHeathDiaryFn(ctx, userID)
}

// GetClientHealthDiaryQuote mocks the implementation of getting a client's health diary quote
func (gm *GormMock) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*gorm.ClientHealthDiaryQuote, error) {
	return gm.MockGetClientHealthDiaryQuoteFn(ctx, limit)
}

// GetClientHealthDiaryEntries mocks the implementation of getting all health diary entries that belong to a specific user
func (gm *GormMock) GetClientHealthDiaryEntries(ctx context.Context, params map[string]interface{}) ([]*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockGetClientHealthDiaryEntriesFn(ctx, params)
}

// UpdateClientCaregiver mocks the implementation of updating a caregiver
func (gm *GormMock) UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	return gm.MockUpdateClientCaregiverFn(ctx, caregiverInput)
}

// SetInProgressBy mocks the implementation of the `SetInProgressBy` update method
func (gm *GormMock) SetInProgressBy(ctx context.Context, requestID, staffID string) (bool, error) {
	return gm.MockInProgressByFn(ctx, requestID, staffID)
}

// GetClientProfileByClientID mocks the implementation of getting a client by client ID
func (gm *GormMock) GetClientProfileByClientID(ctx context.Context, clientID string) (*gorm.Client, error) {
	return gm.MockGetClientProfileByClientIDFn(ctx, clientID)
}

// GetClientsPendingServiceRequestsCount mocks the implementation of getting the service requests count
func (gm *GormMock) GetClientsPendingServiceRequestsCount(ctx context.Context, facilityID string, programID *string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetClientPendingServiceRequestsCountFn(ctx, facilityID, programID)
}

// GetServiceRequests mocks the implementation of getting service requests by type
func (gm *GormMock) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID, programID)
}

// CreateCommunity mocks the implementation of creating a channel
func (gm *GormMock) CreateCommunity(ctx context.Context, community *gorm.Community) (*gorm.Community, error) {
	return gm.MockCreateCommunityFn(ctx, community)
}

// CheckIfUsernameExists mocks the implementation of checking whether a username exists
func (gm *GormMock) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	return gm.MockCheckIfUsernameExistsFn(ctx, username)
}

// GetCommunityByID mocks the implementation of getting the community by ID
func (gm *GormMock) GetCommunityByID(ctx context.Context, communityID string) (*gorm.Community, error) {
	return gm.MockGetCommunityByIDFn(ctx, communityID)
}

// CheckIdentifierExists mocks checking of identifiers
func (gm *GormMock) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	return gm.MockCheckIdentifierExists(ctx, identifierType, identifierValue)
}

// CheckFacilityExistsByIdentifier mocks checking a facility by MFL Code
func (gm *GormMock) CheckFacilityExistsByIdentifier(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
	return gm.MockCheckFacilityExistsByIdentifier(ctx, identifier)
}

// GetOrCreateNextOfKin mocks creating a related person
func (gm *GormMock) GetOrCreateNextOfKin(ctx context.Context, person *gorm.RelatedPerson, clientID, contactID string) error {
	return gm.MockGetOrCreateNextOfKin(ctx, person, clientID, contactID)
}

// GetOrCreateContact mocks creating a contact
func (gm *GormMock) GetOrCreateContact(ctx context.Context, contact *gorm.Contact) (*gorm.Contact, error) {
	return gm.MockGetOrCreateContact(ctx, contact)
}

// GetClientsInAFacility mocks getting clients that belong to a certain facility
func (gm *GormMock) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*gorm.Client, error) {
	return gm.MockGetClientsInAFacilityFn(ctx, facilityID)
}

// GetRecentHealthDiaryEntries mocks getting the most recent health diary entries
func (gm *GormMock) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockGetRecentHealthDiaryEntriesFn(ctx, lastSyncTime, clientID)
}

// GetClientsByParams retrieves clients using the parameters provided
func (gm *GormMock) GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error) {
	return gm.MockGetClientsByParams(ctx, params, lastSyncTime)
}

// GetClientIdentifiers retrieves a client's ccc identifier
func (gm *GormMock) GetClientIdentifiers(ctx context.Context, clientID string) ([]*gorm.Identifier, error) {
	return gm.MockGetClientIdentifiers(ctx, clientID)
}

// GetServiceRequestsForKenyaEMR mocks the getting of service requests attached to a specific facility for use by KenyaEMR
func (gm *GormMock) GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestsForKenyaEMRFn(ctx, facilityID, lastSyncTime)
}

// CreateAppointment creates an appointment in the database
func (gm *GormMock) CreateAppointment(ctx context.Context, appointment *gorm.Appointment) error {
	return gm.MockCreateAppointment(ctx, appointment)
}

// ListAppointments Retrieves appointments using the provided parameters and filters
func (gm *GormMock) ListAppointments(ctx context.Context, params *gorm.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error) {
	return gm.MockListAppointments(ctx, params, filters, pagination)
}

// UpdateAppointment updates the details of an appointment requires the ID or appointment_uuid to be provided
func (gm *GormMock) UpdateAppointment(ctx context.Context, appointment *gorm.Appointment, updateData map[string]interface{}) (*gorm.Appointment, error) {
	return gm.MockUpdateAppointmentFn(ctx, appointment, updateData)
}

// UpdateServiceRequests mocks the implementation of updating service requests from KenyaEMR to MyCareHub
func (gm *GormMock) UpdateServiceRequests(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error) {
	return gm.MockUpdateServiceRequestsFn(ctx, payload)
}

// GetProgramClientProfileByIdentifier mocks the implementation of retrieving a client profile by CCC number
func (gm *GormMock) GetProgramClientProfileByIdentifier(ctx context.Context, programID, identifierType, value string) (*gorm.Client, error) {
	return gm.MockGetProgramClientProfileByIdentifierFn(ctx, programID, identifierType, value)
}

// GetClientProfilesByIdentifier mocks the implementation of retrieving a client profiles by identifier
func (gm *GormMock) GetClientProfilesByIdentifier(ctx context.Context, identifierType, value string) ([]*gorm.Client, error) {
	return gm.MockGetClientProfilesByIdentifierFn(ctx, identifierType, value)
}

// CheckIfClientHasUnresolvedServiceRequests mocks the implementation of checking if a client has a pending service request
func (gm *GormMock) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfClientHasUnresolvedServiceRequestsFn(ctx, clientID, serviceRequestType)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation of updating a user pin change required state
func (gm *GormMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour, status)
}

// SearchClientProfile mocks the implementation of searching for client profiles.
func (gm *GormMock) SearchClientProfile(ctx context.Context, searchTerm string) ([]*gorm.Client, error) {
	return gm.MockSearchClientProfileFn(ctx, searchTerm)
}

// SearchStaffProfile mocks the implementation of getting staff profile using their staff number.
func (gm *GormMock) SearchStaffProfile(ctx context.Context, searchParameter string, programID *string) ([]*gorm.StaffProfile, error) {
	return gm.MockSearchStaffProfileFn(ctx, searchParameter, programID)
}

// UpdateUserPinUpdateRequiredStatus mocks updating a user `pin update required status`
func (gm *GormMock) UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinUpdateRequiredStatusFn(ctx, userID, flavour, status)
}

// UpdateClient updates details for a particular client
func (gm *GormMock) UpdateClient(ctx context.Context, client *gorm.Client, updates map[string]interface{}) (*gorm.Client, error) {
	return gm.MockUpdateClientFn(ctx, client, updates)
}

// GetUserProfileByStaffID mocks the implementation of getting a user profile by staff ID
func (gm *GormMock) GetUserProfileByStaffID(ctx context.Context, staffID string) (*gorm.User, error) {
	return gm.MockGetUserProfileByStaffIDFn(ctx, staffID)
}

// UpdateHealthDiary mocks the implementation of updating the share status of a health diary entry when the client opts for the sharing
func (gm *GormMock) UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
	return gm.MockUpdateHealthDiaryFn(ctx, clientHealthDiaryEntry, updateData)
}

// GetHealthDiaryEntryByID mocks the implementation of getting health diary entry bu a given ID
func (gm *GormMock) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockGetHealthDiaryEntryByIDFn(ctx, healthDiaryEntryID)
}

// UpdateFailedSecurityQuestionsAnsweringAttempts mocks the implementation of resetting failed security attempts
func (gm *GormMock) UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error {
	return gm.MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn(ctx, userID, failCount)
}

// GetClientServiceRequestByID mocks the implementation of getting a service request by ID
func (gm *GormMock) GetClientServiceRequestByID(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
	return gm.MockGetClientServiceRequestByIDFn(ctx, serviceRequestID)
}

// FindContacts retrieves all the contacts that match the given contact type and value.
// Contacts can be shared by users thus the same contact can have multiple records stored
func (gm *GormMock) FindContacts(ctx context.Context, contactType, contactValue string) ([]*gorm.Contact, error) {
	return gm.MockFindContactsFn(ctx, contactType, contactValue)
}

// UpdateUser mocks the implementation of updating a user profile
func (gm *GormMock) UpdateUser(ctx context.Context, user *gorm.User, updateData map[string]interface{}) error {
	return gm.MockUpdateUserFn(ctx, user, updateData)
}

// GetStaffProfileByStaffID mocks the implementation getting staff profile by staff ID
func (gm *GormMock) GetStaffProfileByStaffID(ctx context.Context, staffID string) (*gorm.StaffProfile, error) {
	return gm.MockGetStaffProfileByStaffIDFn(ctx, staffID)
}

// CreateStaffServiceRequest mocks the implementation creating a staff's service request
func (gm *GormMock) CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error {
	return gm.MockCreateStaffServiceRequestFn(ctx, serviceRequestInput)
}

// GetStaffPendingServiceRequestsCount mocks the implementation getting staffs pin reset requests
func (gm *GormMock) GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetStaffPendingServiceRequestsCountFn(ctx, facilityID, programID)
}

// GetStaffServiceRequests mocks the implementation of getting staffs requests
func (gm *GormMock) GetStaffServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
	return gm.MockGetStaffServiceRequestsFn(ctx, requestType, requestStatus, facilityID)
}

// ResolveStaffServiceRequest mocks the implementation resolving staff service requests
func (gm *GormMock) ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error) {
	return gm.MockResolveStaffServiceRequestFn(ctx, staffID, serviceRequestID, verificationStatus)
}

// GetAppointmentServiceRequests mocks the implementation of getting appointments service requests
func (gm *GormMock) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetAppointmentServiceRequestsFn(ctx, lastSyncTime, facilityID)
}

// UpdateFacility mocks the implementation of updating a facility
func (gm *GormMock) UpdateFacility(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error {
	return gm.MockUpdateFacilityFn(ctx, facility, updateData)
}

// GetFacilitiesWithoutFHIRID mocks the implementation of getting a facility without FHIR Organisation
func (gm *GormMock) GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*gorm.Facility, error) {
	return gm.MockGetFacilitiesWithoutFHIRIDFn(ctx)
}

// CheckAppointmentExistsByExternalID checks if an appointment with the external id exists
func (gm *GormMock) CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	return gm.MockCheckAppointmentExistsByExternalIDFn(ctx, externalID)
}

// GetClientServiceRequests mocks the implementation of getting system generated client service requests
func (gm *GormMock) GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetClientServiceRequestsFn(ctx, requestType, status, clientID, facilityID)
}

// CreateUser creates a new user
func (gm *GormMock) CreateUser(ctx context.Context, user *gorm.User) error {
	return gm.MockCreateUserFn(ctx, user)
}

// CreateClient creates a new client
func (gm *GormMock) CreateClient(ctx context.Context, client *gorm.Client, contactID, identifierID string) error {
	return gm.MockCreateClientFn(ctx, client, contactID, identifierID)
}

// GetUserSurveyForms mocks the implementation of getting user survey forms
func (gm *GormMock) GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*gorm.UserSurvey, error) {
	return gm.MockGetUserSurveyFormsFn(ctx, params)
}

// CreateIdentifier creates a new identifier
func (gm *GormMock) CreateIdentifier(ctx context.Context, identifier *gorm.Identifier) error {
	return gm.MockCreateIdentifierFn(ctx, identifier)
}

// CreateNotification saves a new notification to the database
func (gm *GormMock) CreateNotification(ctx context.Context, notification *gorm.Notification) error {
	return gm.MockCreateNotificationFn(ctx, notification)
}

// ListNotifications Retrieves notifications using the provided parameters and filters
func (gm *GormMock) ListNotifications(ctx context.Context, params *gorm.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Notification, *domain.Pagination, error) {
	return gm.MockListNotificationsFn(ctx, params, filters, pagination)
}

// ListAvailableNotificationTypes retrieves the distinct notification types available for a user
func (gm *GormMock) ListAvailableNotificationTypes(ctx context.Context, params *gorm.Notification) ([]enums.NotificationType, error) {
	return gm.MockListAvailableNotificationTypesFn(ctx, params)
}

// GetSharedHealthDiaryEntries mocks the implementation of getting the most recently shared health diary entires by the client to a health care worker
func (gm *GormMock) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockGetSharedHealthDiaryEntriesFn(ctx, clientID, facilityID)
}

// GetClientScreeningToolServiceRequestByToolType mocks the implementation of getting client screening tool service request
func (gm *GormMock) GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*gorm.ClientServiceRequest, error) {
	return gm.MockGetClientScreeningToolServiceRequestByToolTypeFn(ctx, clientID, toolType, status)
}

// GetAppointment returns an appointment by provided params
func (gm *GormMock) GetAppointment(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error) {
	return gm.MockGetAppointmentFn(ctx, params)
}

// CheckIfStaffHasUnresolvedServiceRequests mocks the implementation of checking if a staff has unresolved service requests
func (gm *GormMock) CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfStaffHasUnresolvedServiceRequestsFn(ctx, staffID, serviceRequestType)
}

// GetFacilityStaffs returns a list of staff at a particular facility
func (gm *GormMock) GetFacilityStaffs(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error) {
	return gm.MockGetFacilityStaffsFn(ctx, facilityID)
}

// UpdateNotification updates a notification with the new data
func (gm *GormMock) UpdateNotification(ctx context.Context, notification *gorm.Notification, updateData map[string]interface{}) error {
	return gm.MockUpdateNotificationFn(ctx, notification, updateData)
}

// GetNotification retrieve a notification using the provided ID
func (gm *GormMock) GetNotification(ctx context.Context, notificationID string) (*gorm.Notification, error) {
	return gm.MockGetNotificationFn(ctx, notificationID)
}

// GetClientsByFilterParams returns a list of clients based on the provided filter params
func (gm *GormMock) GetClientsByFilterParams(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*gorm.Client, error) {
	return gm.MockGetClientsByFilterParamsFn(ctx, facilityID, filterParams)
}

// CreateUserSurveys creates a new user survey
func (gm *GormMock) CreateUserSurveys(ctx context.Context, survey []*gorm.UserSurvey) error {
	return gm.MockCreateUserSurveyFn(ctx, survey)
}

// CreateMetric saves a metric to the database
func (gm *GormMock) CreateMetric(ctx context.Context, metric *gorm.Metric) error {
	return gm.MockCreateMetricFn(ctx, metric)
}

// UpdateClientServiceRequest updates a client service request
func (gm *GormMock) UpdateClientServiceRequest(ctx context.Context, serviceRequest *gorm.ClientServiceRequest, updateData map[string]interface{}) error {
	return gm.MockUpdateClientServiceRequestFn(ctx, serviceRequest, updateData)
}

// SaveFeedback mocks the implementation of saving feedback into the database
func (gm *GormMock) SaveFeedback(ctx context.Context, feedback *gorm.Feedback) error {
	return gm.MockSaveFeedbackFn(ctx, feedback)
}

// RegisterStaff mocks the implementation of registering a staff
func (gm *GormMock) RegisterStaff(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) (*gorm.StaffProfile, error) {
	return gm.MockRegisterStaffFn(ctx, user, contact, identifier, staffProfile)
}

// SearchClientServiceRequests mocks the implementation of searching client service requests
func (gm *GormMock) SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockSearchClientServiceRequestsFn(ctx, searchParameter, requestType, facilityID)
}

// SearchStaffServiceRequests mocks the implementation of searching client service requests
func (gm *GormMock) SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error) {
	return gm.MockSearchStaffServiceRequestsFn(ctx, searchParameter, requestType, facilityID)
}

// RegisterClient mocks the implementation of registering a client
func (gm *GormMock) RegisterClient(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
	return gm.MockRegisterClientFn(ctx, user, contact, identifier, client)
}

// DeleteCommunity deletes the specified community from the database
func (gm *GormMock) DeleteCommunity(ctx context.Context, communityID string) error {
	return gm.MockDeleteCommunityFn(ctx, communityID)
}

// CreateQuestionnaire mocks the implementation of creating a questionnaire
func (gm *GormMock) CreateQuestionnaire(ctx context.Context, questionnaire *gorm.Questionnaire) error {
	return gm.MockCreateQuestionnaireFn(ctx, questionnaire)
}

// CreateScreeningTool mocks the implementation of creating a screening tool
func (gm *GormMock) CreateScreeningTool(ctx context.Context, screeningTool *gorm.ScreeningTool) error {
	return gm.MockCreateScreeningToolFn(ctx, screeningTool)
}

// CreateQuestion mocks the implementation of creating a question
func (gm *GormMock) CreateQuestion(ctx context.Context, question *gorm.Question) error {
	return gm.MockCreateQuestionFn(ctx, question)
}

// CreateQuestionChoice mocks the implementation of creating a question input choice
func (gm *GormMock) CreateQuestionChoice(ctx context.Context, questionChoice *gorm.QuestionInputChoice) error {
	return gm.MockCreateQuestionChoiceFn(ctx, questionChoice)
}

// GetScreeningToolByID mocks the implementation of getting a screening tool by ID
func (gm *GormMock) GetScreeningToolByID(ctx context.Context, screeningToolID string) (*gorm.ScreeningTool, error) {
	return gm.MockGetScreeningToolByIDFn(ctx, screeningToolID)
}

// GetQuestionnaireByID mocks the implementation of getting a questionnaire by ID
func (gm *GormMock) GetQuestionnaireByID(ctx context.Context, questionnaireID string) (*gorm.Questionnaire, error) {
	return gm.MockGetQuestionnaireByIDFn(ctx, questionnaireID)
}

// GetQuestionsByQuestionnaireID mocks the implementation of getting questions by questionnaire ID
func (gm *GormMock) GetQuestionsByQuestionnaireID(ctx context.Context, questionnaireID string) ([]*gorm.Question, error) {
	return gm.MockGetQuestionsByQuestionnaireIDFn(ctx, questionnaireID)
}

// GetQuestionInputChoicesByQuestionID mocks the implementation of getting question input choices by question ID
func (gm *GormMock) GetQuestionInputChoicesByQuestionID(ctx context.Context, questionID string) ([]*gorm.QuestionInputChoice, error) {
	return gm.MockGetQuestionInputChoicesByQuestionIDFn(ctx, questionID)
}

// CreateScreeningToolResponse mocks the implementation of creating a screening tool response
func (gm *GormMock) CreateScreeningToolResponse(ctx context.Context, screeningToolResponse *gorm.ScreeningToolResponse, screeningToolQuestionResponses []*gorm.ScreeningToolQuestionResponse) (*string, error) {
	return gm.MockCreateScreeningToolResponseFn(ctx, screeningToolResponse, screeningToolQuestionResponses)
}

// GetAvailableScreeningTools mocks the implementation of getting available screening tools
func (gm *GormMock) GetAvailableScreeningTools(ctx context.Context, clientID string, screeningTool gorm.ScreeningTool, screeningToolIDs []string) ([]*gorm.ScreeningTool, error) {
	return gm.MockGetAvailableScreeningToolsFn(ctx, clientID, screeningTool, screeningToolIDs)
}

// GetScreeningToolResponsesWithin24Hours mocks the implementation of GetScreeningToolResponsesWithin24Hours method
func (gm *GormMock) GetScreeningToolResponsesWithin24Hours(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponsesWithin24HoursFn(ctx, clientID, programID)
}

// GetScreeningToolResponsesWithPendingServiceRequests mocks the implementation of GetScreeningToolResponsesWithPendingServiceRequests method
func (gm *GormMock) GetScreeningToolResponsesWithPendingServiceRequests(ctx context.Context, clientID, programID string) ([]*gorm.ScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponsesWithPendingServiceRequestsFn(ctx, clientID, programID)
}

// GetFacilityRespondedScreeningTools mocks the response returned by getting facility responded screening tools
func (gm *GormMock) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*gorm.ScreeningTool, *domain.Pagination, error) {
	return gm.MockGetFacilityRespondedScreeningToolsFn(ctx, facilityID, programID, pagination)
}

// ListSurveyRespondents mocks the implementation of listing survey respondents
func (gm *GormMock) ListSurveyRespondents(ctx context.Context, params *gorm.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*gorm.UserSurvey, *domain.Pagination, error) {
	return gm.MockListSurveyRespondentsFn(ctx, params, facilityID, pagination)
}

// GetScreeningToolServiceRequestOfRespondents mocks the implementation of getting screening tool service requests by respondents
func (gm *GormMock) GetScreeningToolServiceRequestOfRespondents(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
	return gm.MockGetScreeningToolServiceRequestOfRespondentsFn(ctx, facilityID, programID, screeningToolID, searchTerm, pagination)
}

// GetScreeningToolResponseByID mocks the implementation of getting a screening tool response by ID
func (gm *GormMock) GetScreeningToolResponseByID(ctx context.Context, id string) (*gorm.ScreeningToolResponse, error) {
	return gm.MockGetScreeningToolResponseByIDFn(ctx, id)
}

// GetScreeningToolQuestionResponsesByResponseID mocks the implementation of getting screening tool question responses by response ID
func (gm *GormMock) GetScreeningToolQuestionResponsesByResponseID(ctx context.Context, responseID string) ([]*gorm.ScreeningToolQuestionResponse, error) {
	return gm.MockGetScreeningToolQuestionResponsesByResponseIDFn(ctx, responseID)
}

// GetSurveysWithServiceRequests mocks the implementation of getting surveys with service requests
func (gm *GormMock) GetSurveysWithServiceRequests(ctx context.Context, facilityID, programID string) ([]*gorm.UserSurvey, error) {
	return gm.MockGetSurveysWithServiceRequestsFn(ctx, facilityID, programID)
}

// GetClientsSurveyServiceRequest mocks the implementation of getting clients with survey service request
func (gm *GormMock) GetClientsSurveyServiceRequest(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*gorm.ClientServiceRequest, *domain.Pagination, error) {
	return gm.MockGetClientsSurveyServiceRequestFn(ctx, facilityID, projectID, formID, pagination)
}

// GetStaffFacilities mocks the implementation of getting a list of staff facilities
func (gm *GormMock) GetStaffFacilities(ctx context.Context, staffFacility gorm.StaffFacilities, pagination *domain.Pagination) ([]*gorm.StaffFacilities, *domain.Pagination, error) {
	return gm.MockGetStaffFacilitiesFn(ctx, staffFacility, pagination)
}

// GetClientFacilities mocks the implementation of getting a list of client facilities
func (gm *GormMock) GetClientFacilities(ctx context.Context, clientFacility gorm.ClientFacilities, pagination *domain.Pagination) ([]*gorm.ClientFacilities, *domain.Pagination, error) {
	return gm.MockGetClientFacilitiesFn(ctx, clientFacility, pagination)
}

// GetClientsSurveyCount mocks the implementation of getting clients survey count
func (gm *GormMock) GetClientsSurveyCount(ctx context.Context, userID string) (int, error) {
	return gm.MockGetClientsSurveyCountFn(ctx, userID)
}

// UpdateStaff mock the implementation of updating a staff profile
func (gm *GormMock) UpdateStaff(ctx context.Context, staff *gorm.StaffProfile, updates map[string]interface{}) (*gorm.StaffProfile, error) {
	return gm.MockUpdateStaffFn(ctx, staff, updates)
}

// AddFacilitiesToStaffProfile mocks the implementation of adding facilities to a staff profile
func (gm *GormMock) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return gm.MockAddFacilitiesToStaffProfileFn(ctx, staffID, facilities)
}

// AddFacilitiesToClientProfile mocks the implementation of adding facilities to a client profile
func (gm *GormMock) AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return gm.MockAddFacilitiesToClientProfileFn(ctx, clientID, facilities)
}

// GetNotificationsCount mocks the implementation of getting notification count
func (gm *GormMock) GetNotificationsCount(ctx context.Context, notification gorm.Notification) (int, error) {
	return gm.MockGetNotificationsCountFn(ctx, notification)
}

// RegisterCaregiver registers a new caregiver
func (gm *GormMock) RegisterCaregiver(ctx context.Context, user *gorm.User, contact *gorm.Contact, caregiver *gorm.Caregiver) error {
	return gm.MockRegisterCaregiverFn(ctx, user, contact, caregiver)
}

// CreateCaregiver creates a caregiver record linked to a user
func (gm *GormMock) CreateCaregiver(ctx context.Context, caregiver *gorm.Caregiver) error {
	return gm.MockCreateCaregiverFn(ctx, caregiver)
}

// SearchCaregiverUser mocks the searching of caregiver user
func (gm *GormMock) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
	return gm.MockSearchCaregiverUserFn(ctx, searchParameter)
}

// SearchPlatformCaregivers mocks the searching of caregiver user
func (gm *GormMock) SearchPlatformCaregivers(ctx context.Context, searchParameter string) ([]*gorm.Caregiver, error) {
	return gm.MockSearchPlatformCaregiversFn(ctx, searchParameter)
}

// RemoveFacilitiesFromClientProfile mocks the implementation of removing facilities from a client profile
func (gm *GormMock) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return gm.MockRemoveFacilitiesFromClientProfileFn(ctx, clientID, facilities)
}

// AddCaregiverToClient mocks the implementation of adding a caregiver to a client
func (gm *GormMock) AddCaregiverToClient(ctx context.Context, clientCaregiver *gorm.CaregiverClient) error {
	return gm.MockAddCaregiverToClientFn(ctx, clientCaregiver)
}

// RemoveFacilitiesFromStaffProfile mocks the implementation of removing facilities from a staff profile
func (gm *GormMock) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return gm.MockRemoveFacilitiesFromStaffProfileFn(ctx, staffID, facilities)
}

// GetCaregiverManagedClients mocks the implementation of getting caregiver's managed clients
func (gm *GormMock) GetCaregiverManagedClients(ctx context.Context, userID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
	return gm.MockGetCaregiverManagedClientsFn(ctx, userID, pagination)
}

// GetCaregiversClient mocks the implementation of getting a record of client - caregiver association
func (gm *GormMock) GetCaregiversClient(ctx context.Context, caregiverClient gorm.CaregiverClient) ([]*gorm.CaregiverClient, error) {
	return gm.MockGetCaregiversClientFn(ctx, caregiverClient)
}

// ListClientsCaregivers mocks the implementation of listing clients caregivers
func (gm *GormMock) ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*gorm.CaregiverClient, *domain.Pagination, error) {
	return gm.MockListClientsCaregiversFn(ctx, clientID, pagination)
}

// GetCaregiverProfileByCaregiverID mocks the implementation of getting a caregiver profile by caregiver ID
func (gm *GormMock) GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*gorm.Caregiver, error) {
	return gm.MockGetCaregiverProfileByCaregiverIDFn(ctx, caregiverID)
}

// UpdateCaregiverClient mocks the accepting of caregiver consent
func (gm *GormMock) UpdateCaregiverClient(ctx context.Context, caregiverClient *gorm.CaregiverClient, updates map[string]interface{}) error {
	return gm.MockUpdateCaregiverClientFn(ctx, caregiverClient, updates)
}

// CreateOrganisation mocks the implementation of creating an organisation
func (gm *GormMock) CreateOrganisation(ctx context.Context, organisation *gorm.Organisation) (*gorm.Organisation, error) {
	return gm.MockCreateOrganisationFn(ctx, organisation)
}

// CreateProgram mocks the implementation of creating a program
func (gm *GormMock) CreateProgram(ctx context.Context, program *gorm.Program) (*gorm.Program, error) {
	return gm.MockCreateProgramFn(ctx, program)
}

// CheckOrganisationExists mocks the implementation checking if the an organisation exists
func (gm *GormMock) CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error) {
	return gm.MockCheckOrganisationExistsFn(ctx, organisationID)
}

// CheckIfProgramNameExists mocks the implementation checking if an organisation is associated with a program
func (gm *GormMock) CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error) {
	return gm.MockCheckIfProgramNameExistsFn(ctx, organisationID, programName)
}

// DeleteOrganisation mocks the implementation of deleting an organisation
func (gm *GormMock) DeleteOrganisation(ctx context.Context, organisation *gorm.Organisation) error {
	return gm.MockDeleteOrganisationFn(ctx, organisation)
}

// AddFacilityToProgram mocks the implementation of adding a facility to a program
func (gm *GormMock) AddFacilityToProgram(ctx context.Context, programID string, facilityID []string) error {
	return gm.MockAddFacilityToProgramFn(ctx, programID, facilityID)
}

// RegisterExistingUserAsClient mocks the implementation of registering an existing user as a client
func (gm *GormMock) RegisterExistingUserAsClient(ctx context.Context, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
	return gm.MockRegisterExistingUserAsClientFn(ctx, identifier, client)
}

// RegisterExistingUserAsStaff mocks the implementation of registering an existing user as staff
func (gm *GormMock) RegisterExistingUserAsStaff(ctx context.Context, identifier *gorm.Identifier, staff *gorm.StaffProfile) (*gorm.StaffProfile, error) {
	return gm.MockRegisterExistingUserAsStaffFn(ctx, identifier, staff)
}

// ListOrganisations mocks the implementation of listing organisations
func (gm *GormMock) ListOrganisations(ctx context.Context, pagination *domain.Pagination) ([]*gorm.Organisation, *domain.Pagination, error) {
	return gm.MockListOrganisationsFn(ctx, pagination)
}

// GetProgramFacilities mocks the implementation of listing program facilities
func (gm *GormMock) GetProgramFacilities(ctx context.Context, programID string) ([]*gorm.ProgramFacility, error) {
	return gm.MockGetProgramFacilitiesFn(ctx, programID)
}

// GetProgramByID mocks the implementation of getting a program by ID
func (gm *GormMock) GetProgramByID(ctx context.Context, programID string) (*gorm.Program, error) {
	return gm.MockGetProgramByIDFn(ctx, programID)
}

// ListPrograms mocks the implementation of getting programs
func (gm *GormMock) ListPrograms(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
	return gm.MockListProgramsFn(ctx, organisationID, pagination)
}

// CheckIfSuperUserExists mocks the implementation of checking if a superuser exists
func (gm *GormMock) CheckIfSuperUserExists(ctx context.Context) (bool, error) {
	return gm.MockCheckIfSuperUserExistsFn(ctx)
}

// GetCaregiverProfileByUserID mocks the implementation of getting a caregiver profile
func (gm *GormMock) GetCaregiverProfileByUserID(ctx context.Context, userID string, organisationID string) (*gorm.Caregiver, error) {
	return gm.MockGetCaregiverProfileByUserIDFn(ctx, userID, organisationID)
}

// UpdateCaregiver mocks the implementation of updating a caregiver
func (gm *GormMock) UpdateCaregiver(ctx context.Context, caregiver *gorm.Caregiver, updates map[string]interface{}) error {
	return gm.MockUpdateCaregiverFn(ctx, caregiver, updates)
}

// RegisterExistingUserAsCaregiver mocks the implementation of registering an existing user as a caregiver
func (gm *GormMock) RegisterExistingUserAsCaregiver(ctx context.Context, caregiver *gorm.Caregiver) (*gorm.Caregiver, error) {
	return gm.MockRegisterExistingUserAsCaregiverFn(ctx, caregiver)
}

// UpdateUserContact mocks the implementation of updating a user contact
func (gm *GormMock) UpdateUserContact(ctx context.Context, userContact *gorm.Contact, updates map[string]interface{}) error {
	return gm.MockUpdateUserContactFn(ctx, userContact, updates)
}

// UpdateClientIdentifier mocks the implementation of updating a client identifier
func (gm *GormMock) UpdateClientIdentifier(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error {
	return gm.MockUpdateClientIdentifierFn(ctx, clientID, identifierType, identifierValue, programID)
}

// SearchOrganisation mocks the implementation of searching for organisations
func (gm *GormMock) SearchOrganisation(ctx context.Context, searchParameter string) ([]*gorm.Organisation, error) {
	return gm.MockSearchOrganisationsFn(ctx, searchParameter)
}

// SearchPrograms mocks the implementation of searching programs
func (gm *GormMock) SearchPrograms(ctx context.Context, searchParameter string, organisationID string, pagination *domain.Pagination) ([]*gorm.Program, *domain.Pagination, error) {
	return gm.MockSearchProgramsFn(ctx, searchParameter, organisationID, pagination)
}

// CreateFacilities Mocks the implementation of CreateFacilities method
func (gm *GormMock) CreateFacilities(ctx context.Context, facilities []*gorm.Facility) ([]*gorm.Facility, error) {
	return gm.MockCreateFacilitiesFn(ctx, facilities)
}

// ListCommunities mocks the implementation of listing communities
func (gm *GormMock) ListCommunities(ctx context.Context, programID string, organisationID string) ([]*gorm.Community, error) {
	return gm.MockListCommunitiesFn(ctx, programID, organisationID)
}

// CreateSecurityQuestions mocks the implementation of CreateSecurityQuestions method
func (gm *GormMock) CreateSecurityQuestions(ctx context.Context, securityQuestions []*gorm.SecurityQuestion) ([]*gorm.SecurityQuestion, error) {
	return gm.MockCreateSecurityQuestionsFn(ctx, securityQuestions)
}

// CreateTermsOfService mocks the implementation of CreateTermsOfService method
func (gm *GormMock) CreateTermsOfService(ctx context.Context, termsOfService *gorm.TermsOfService) (*gorm.TermsOfService, error) {
	return gm.MockCreateTermsOfServiceFn(ctx, termsOfService)
}

// CheckPhoneExists mocks the implementation of CheckPhoneExists method
func (gm *GormMock) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	return gm.MockCheckPhoneExistsFn(ctx, phone)
}

// UpdateProgram updates the program details
func (gm *GormMock) UpdateProgram(ctx context.Context, program *gorm.Program, updateData map[string]interface{}) error {
	return gm.MockUpdateProgramFn(ctx, program, updateData)
}

// GetStaffServiceRequestByID mocks the implementation of GetStaffServiceRequestByID method
func (gm *GormMock) GetStaffServiceRequestByID(ctx context.Context, serviceRequestID string) (*gorm.StaffServiceRequest, error) {
	return gm.MockGetStaffServiceRequestByIDFn(ctx, serviceRequestID)
}

// GetClientJWT retrieves a JWT by unique JTI
func (gm *GormMock) GetClientJWT(ctx context.Context, jti string) (*gorm.OauthClientJWT, error) {
	return gm.MockGetClientJWT(ctx, jti)
}

// GetOauthClient retrieves a client by ID
func (gm *GormMock) GetOauthClient(ctx context.Context, id string) (*gorm.OauthClient, error) {
	return gm.MockGetOauthClient(ctx, id)
}

// GetValidClientJWT retrieves a JWT that is still valid i.e not expired
func (gm *GormMock) GetValidClientJWT(ctx context.Context, jti string) (*gorm.OauthClientJWT, error) {
	return gm.MockGetValidClientJWT(ctx, jti)
}

// CreateOauthClientJWT creates a new oauth jwt client
func (gm *GormMock) CreateOauthClientJWT(ctx context.Context, jwt *gorm.OauthClientJWT) error {
	return gm.MockCreateOauthClientJWT(ctx, jwt)
}

// CreateOauthClient creates a new oauth client
func (gm *GormMock) CreateOauthClient(ctx context.Context, client *gorm.OauthClient) error {
	return gm.MockCreateOauthClient(ctx, client)
}

// CreateOrUpdateSession creates a new session or updates an existing session
func (gm *GormMock) CreateOrUpdateSession(ctx context.Context, session *gorm.Session) error {
	return gm.MockCreateOrUpdateSessionFn(ctx, session)
}

// CreateAuthorizationCode creates a new authorization code.
func (gm *GormMock) CreateAuthorizationCode(ctx context.Context, code *gorm.AuthorizationCode) error {
	return gm.MockCreateAuthorizationCodeFn(ctx, code)
}

// GetAuthorizationCode retrieves an authorization code using the code
func (gm *GormMock) GetAuthorizationCode(ctx context.Context, code string) (*gorm.AuthorizationCode, error) {
	return gm.MockGetAuthorizationCodeFn(ctx, code)
}

// UpdateAuthorizationCode updates the details of a given code
func (gm *GormMock) UpdateAuthorizationCode(ctx context.Context, code *gorm.AuthorizationCode, updateData map[string]interface{}) error {
	return gm.MockUpdateAuthorizationCodeFn(ctx, code, updateData)
}

// CreateAccessToken creates a new access token.
func (gm *GormMock) CreateAccessToken(ctx context.Context, token *gorm.AccessToken) error {
	return gm.MockCreateAccessTokenFn(ctx, token)
}

// CreateRefreshToken creates a new refresh token.
func (gm *GormMock) CreateRefreshToken(ctx context.Context, token *gorm.RefreshToken) error {
	return gm.MockCreateRefreshTokenFn(ctx, token)
}

// DeleteAccessToken retrieves an access token using the signature
func (gm *GormMock) DeleteAccessToken(ctx context.Context, signature string) error {
	return gm.MockDeleteAccessTokenFn(ctx, signature)
}

// DeleteRefreshToken retrieves a refresh token using the signature
func (gm *GormMock) DeleteRefreshToken(ctx context.Context, signature string) error {
	return gm.MockDeleteRefreshTokenFn(ctx, signature)
}

// GetAccessToken retrieves an access token using the signature
func (gm *GormMock) GetAccessToken(ctx context.Context, token gorm.AccessToken) (*gorm.AccessToken, error) {
	return gm.MockGetAccessTokenFn(ctx, token)
}

// GetRefreshToken retrieves a refresh token using the signature
func (gm *GormMock) GetRefreshToken(ctx context.Context, token gorm.RefreshToken) (*gorm.RefreshToken, error) {
	return gm.MockGetRefreshTokenFn(ctx, token)
}

// UpdateAccessToken updates the details of a given access token
func (gm *GormMock) UpdateAccessToken(ctx context.Context, code *gorm.AccessToken, updateData map[string]interface{}) error {
	return gm.MockUpdateAccessTokenFn(ctx, code, updateData)
}

// UpdateRefreshToken updates the details of a given refresh token
func (gm *GormMock) UpdateRefreshToken(ctx context.Context, code *gorm.RefreshToken, updateData map[string]interface{}) error {
	return gm.MockUpdateRefreshTokenFn(ctx, code, updateData)
}

// CheckIfClientHasPendingSurveyServiceRequest mocks the implementation of CheckIfClientHasPendingSurveyServiceRequest method
func (gm *GormMock) CheckIfClientHasPendingSurveyServiceRequest(ctx context.Context, clientID string, projectID int, formID string) (bool, error) {
	return gm.MockCheckIfClientHasPendingSurveyServiceRequestFn(ctx, clientID, projectID, formID)
}

// GetUserProfileByPushToken mocks the implementation of get user profile by push token
func (gm *GormMock) GetUserProfileByPushToken(ctx context.Context, pushToken string) (*gorm.User, error) {
	return gm.MockGetUserProfileByPushTokenFn(ctx, pushToken)
}

// CheckIfFacilityExistsInProgram mocks the implementation of CheckIfFacilityExistsInProgram method
func (gm *GormMock) CheckIfFacilityExistsInProgram(ctx context.Context, programID, facilityID string) (bool, error) {
	return gm.MockCheckIfFacilityExistsInProgramFn(ctx, programID, facilityID)
}

// CheckIfFacilityExistsInProgram mocks the implementation of CheckIfFacilityExistsInProgram method
func (gm *GormMock) GetStaffIdentifiers(ctx context.Context, staffID string, identifierType *string) ([]*gorm.Identifier, error) {
	return gm.MockGetStaffIdentifiersFn(ctx, staffID, identifierType)
}

// CheckStaffExistsInProgram mocks the implementation of CheckStaffExistsInProgram method
func (gm *GormMock) CheckStaffExistsInProgram(ctx context.Context, userID, programID string) (bool, error) {
	return gm.MockCheckStaffExistsInProgramFn(ctx, userID, programID)
}

// CheckStaffExistsInProgram mocks the implementation of CheckStaffExistsInProgram method
func (gm *GormMock) CheckIfClientExistsInProgram(ctx context.Context, userID, programID string) (bool, error) {
	return gm.MockCheckIfClientExistsInProgramFn(ctx, userID, programID)
}

// GetUserClientProfiles mocks the implementation of GetUserClientProfiles method
func (gm *GormMock) GetUserClientProfiles(ctx context.Context, userID string) ([]*gorm.Client, error) {
	return gm.MockGetUserClientProfilesFn(ctx, userID)
}

// GetUserStaffProfiles mocks the implementation of GetUserStaffProfiles method
func (gm *GormMock) GetUserStaffProfiles(ctx context.Context, userID string) ([]*gorm.StaffProfile, error) {
	return gm.MockGetUserStaffProfilesFn(ctx, userID)
}
