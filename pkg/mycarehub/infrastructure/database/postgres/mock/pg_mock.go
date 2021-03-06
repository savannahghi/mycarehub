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
	MockCreateUserFn                                     func(ctx context.Context, user domain.User) (*domain.User, error)
	MockCreateClientFn                                   func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error)
	MockCreateIdentifierFn                               func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error)
	MockGetOrCreateFacilityFn                            func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	MockSearchFacilityFn                                 func(ctx context.Context, searchParameter *string) ([]*domain.Facility, error)
	MockRetrieveFacilityFn                               func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	ListFacilitiesFn                                     func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockDeleteFacilityFn                                 func(ctx context.Context, id int) (bool, error)
	MockRetrieveFacilityByMFLCodeFn                      func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
	MockGetUserProfileByPhoneNumberFn                    func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error)
	MockGetUserPINByUserIDFn                             func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error)
	MockInactivateFacilityFn                             func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn                             func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn                         func(ctx context.Context, userID string) (*domain.User, error)
	MockSaveTemporaryUserPinFn                           func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	MockGetCurrentTermsFn                                func(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error)
	MockAcceptTermsFn                                    func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                        func(ctx context.Context, pin *domain.UserPIN) (bool, error)
	MockSetNickNameFn                                    func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockGetSecurityQuestionsFn                           func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	MockSaveOTPFn                                        func(ctx context.Context, otpInput *domain.OTP) error
	MockGetSecurityQuestionByIDFn                        func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn                   func(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	MockGetSecurityQuestionResponseFn                    func(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                       func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                                      func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetClientProfileByUserIDFn                       func(ctx context.Context, userID string) (*domain.ClientProfile, error)
	MockGetStaffProfileByUserIDFn                        func(ctx context.Context, userID string) (*domain.StaffProfile, error)
	MockCheckUserHasPinFn                                func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGenerateRetryOTPFn                               func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	MockCompleteOnboardingTourFn                         func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                         func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	MockGetUserSecurityQuestionsResponsesFn              func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error)
	MockInvalidatePINFn                                  func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetContactByUserIDFn                             func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn        func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockFetchFacilitiesFn                                func(ctx context.Context) ([]*domain.Facility, error)
	MockCreateHealthDiaryEntryFn                         func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error
	MockCreateServiceRequestFn                           func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	MockCanRecordHeathDiaryFn                            func(ctx context.Context, userID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn                      func(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error)
	MockGetClientHealthDiaryEntriesFn                    func(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error)
	MockCreateClientCaregiverFn                          func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockGetClientCaregiverFn                             func(ctx context.Context, caregiverID string) (*domain.Caregiver, error)
	MockUpdateClientCaregiverFn                          func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockUpdateFacilityFn                                 func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error
	MockInProgressByFn                                   func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn                     func(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	MockGetPendingServiceRequestsCountFn                 func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error)
	MockGetServiceRequestsFn                             func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error)
	MockResolveServiceRequestFn                          func(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error
	MockCreateCommunityFn                                func(ctx context.Context, community *dto.CommunityInput) (*domain.Community, error)
	MockCheckUserRoleFn                                  func(ctx context.Context, userID string, role string) (bool, error)
	MockCheckUserPermissionFn                            func(ctx context.Context, userID string, permission string) (bool, error)
	MockAssignRolesFn                                    func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetUserRolesFn                                   func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	MockGetUserPermissionsFn                             func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
	MockCheckIfUsernameExistsFn                          func(ctx context.Context, username string) (bool, error)
	MockRevokeRolesFn                                    func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetCommunityByIDFn                               func(ctx context.Context, communityID string) (*domain.Community, error)
	MockCheckIdentifierExists                            func(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	MockCheckFacilityExistsByMFLCode                     func(ctx context.Context, MFLCode int) (bool, error)
	MockGetOrCreateNextOfKin                             func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error
	MockGetOrCreateContactFn                             func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
	MockGetClientsInAFacilityFn                          func(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	MockGetRecentHealthDiaryEntriesFn                    func(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                               func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	MockGetClientCCCIdentifier                           func(ctx context.Context, clientID string) (*domain.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn                  func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	MockCreateAppointment                                func(ctx context.Context, appointment domain.Appointment) error
	MockUpdateAppointmentFn                              func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error)
	MockGetScreeningToolsQuestionsFn                     func(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error)
	MockAnswerScreeningToolQuestionsFn                   func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error
	MockGetScreeningToolQuestionByQuestionIDFn           func(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error)
	MockSearchStaffProfileFn                             func(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error)
	MockUpdateHealthDiaryFn                              func(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error
	MockInvalidateScreeningToolResponseFn                func(ctx context.Context, clientID string, questionID string) error
	MockUpdateServiceRequestsFn                          func(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error)
	MockListAppointments                                 func(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	MockGetClientProfileByCCCNumberFn                    func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error)
	MockUpdateUserPinChangeRequiredStatusFn              func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockSearchClientProfileFn                            func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
	MockCheckIfClientHasUnresolvedServiceRequestsFn      func(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	MockUpdateUserSurveysFn                              func(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error
	MockGetAllRolesFn                                    func(ctx context.Context) ([]*domain.AuthorityRole, error)
	MockUpdateUserPinUpdateRequiredStatusFn              func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockGetHealthDiaryEntryByIDFn                        func(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error)
	MockUpdateClientFn                                   func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error)
	MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn func(ctx context.Context, userID string, failCount int) error
	MockGetFacilitiesWithoutFHIRIDFn                     func(ctx context.Context) ([]*domain.Facility, error)
	MockGetSharedHealthDiaryEntriesFn                    func(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetServiceRequestByIDFn                          func(ctx context.Context, id string) (*domain.ServiceRequest, error)
	MockUpdateUserFn                                     func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error
	MockGetStaffProfileByStaffIDFn                       func(ctx context.Context, staffID string) (*domain.StaffProfile, error)
	MockResolveStaffServiceRequestFn                     func(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error)
	MockCreateStaffServiceRequestFn                      func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error
	MockGetAppointmentServiceRequestsFn                  func(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error)
	MockGetClientAppointmentByIDFn                       func(ctx context.Context, appointmentID string) (*domain.Appointment, error)
	MockGetAssessmentResponsesFn                         func(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error)
	MockGetAppointmentByAppointmentUUIDFn                func(ctx context.Context, appointmentUUID string) (*domain.Appointment, error)
	MockGetClientServiceRequestsFn                       func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error)
	MockGetActiveScreeningToolResponsesFn                func(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error)
	MockGetAppointmentByClientIDFn                       func(ctx context.Context, clientID string) (*domain.Appointment, error)
	MockCheckAppointmentExistsByExternalIDFn             func(ctx context.Context, externalID string) (bool, error)
	MockGetUserSurveyFormsFn                             func(ctx context.Context, userID string) ([]*domain.UserSurvey, error)
	MockListNotificationsFn                              func(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error)
	MockListAvailableNotificationTypesFn                 func(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error)
	MockSaveNotificationFn                               func(ctx context.Context, payload *domain.Notification) error
	MockGetClientScreeningToolResponsesByToolTypeFn      func(ctx context.Context, clientID, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error)
	MockGetClientScreeningToolServiceRequestByToolTypeFn func(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error)
	MockGetAppointmentFn                                 func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error)
	MockGetFacilityStaffsFn                              func(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error)
	MockCheckIfStaffHasUnresolvedServiceRequestsFn       func(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	MockDeleteUserFn                                     func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error
	MockDeleteStaffProfileFn                             func(ctx context.Context, staffID string) error
	MockUpdateNotificationFn                             func(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error
	MockGetNotificationFn                                func(ctx context.Context, notificationID string) (*domain.Notification, error)
	MockGetClientsByFilterParamsFn                       func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error)
	MockCreateUserSurveyFn                               func(ctx context.Context, userSurvey []*dto.UserSurveyInput) error
	MockCreateMetricFn                                   func(ctx context.Context, payload *domain.Metric) error
	MockUpdateClientServiceRequestFn                     func(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error
	MockSaveFeedbackFn                                   func(ctx context.Context, feedback *domain.FeedbackResponse) error
	MockSearchClientServiceRequestsFn                    func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	MockSearchStaffServiceRequestsFn                     func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	MockRegisterClientFn                                 func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error)
	MockRegisterStaffFn                                  func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	ID := uuid.New().String()

	name := gofakeit.Name()
	code := gofakeit.Number(0, 100)
	county := "Nairobi"
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
		Flavour:      "CONSUMER",
	}

	facilityInput := &domain.Facility{
		ID:          &ID,
		Name:        name,
		Code:        code,
		Phone:       phone,
		Active:      true,
		County:      county,
		Description: description,
	}

	var facilitiesList []*domain.Facility
	facilitiesList = append(facilitiesList, facilityInput)
	nextPage := 3
	previousPage := 1
	facilitiesPage := &domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:        1,
			CurrentPage:  2,
			Count:        3,
			TotalPages:   3,
			NextPage:     &nextPage,
			PreviousPage: &previousPage,
		},
		Facilities: []domain.Facility{
			{
				ID:          &ID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			},
		},
	}

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
		FacilityID:              ID,
		FacilityName:            name,
		CHVUserID:               &ID,
		CHVUserName:             name,
		CaregiverID:             &ID,
		CCCNumber:               "123456789",
	}
	staff := &domain.StaffProfile{
		ID:                &ID,
		User:              userProfile,
		UserID:            uuid.New().String(),
		Active:            false,
		StaffNumber:       gofakeit.BeerAlcohol(),
		Facilities:        []domain.Facility{*facilityInput},
		DefaultFacilityID: gofakeit.BeerAlcohol(),
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
		Mood:                  "test",
		Note:                  "test",
		EntryType:             "test",
		ShareWithHealthWorker: false,
		SharedAt:              &currentTime,
		ClientID:              ID,
		CreatedAt:             time.Now(),
		PhoneNumber:           phone,
		ClientName:            name,
	}

	return &PostgresMock{
		MockCreateMetricFn: func(ctx context.Context, payload *domain.Metric) error {
			return nil
		},
		MockUpdateNotificationFn: func(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error {
			return nil
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
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockGetFacilityStaffsFn: func(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error) {
			return []*domain.StaffProfile{staff}, nil
		},
		MockSearchFacilityFn: func(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},
		ListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},
		MockDeleteFacilityFn: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},
		MockGetOrCreateContactFn: func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
			return contactData, nil
		},
		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
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
		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
			return &domain.UserPIN{
				UserID:    userID,
				ValidFrom: time.Now().Add(time.Hour * 10),
				ValidTo:   time.Now().Add(time.Hour * 20),
				Flavour:   flavour,
				IsValid:   false,
			}, nil
		},
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
			return userProfile, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockGetStaffProfileByStaffIDFn: func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
			return &domain.StaffProfile{
				ID:                &ID,
				User:              userProfile,
				UserID:            ID,
				Active:            false,
				StaffNumber:       "TEST-00",
				Facilities:        []domain.Facility{},
				DefaultFacilityID: ID,
			}, nil
		},
		MockGetCurrentTermsFn: func(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error) {
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
		MockGetUserProfileByUserIDFn: func(ctx context.Context, userID string) (*domain.User, error) {
			return &domain.User{
				ID:            &userID,
				Username:      gofakeit.Name(),
				Name:          gofakeit.Name(),
				Active:        true,
				TermsAccepted: true,
				UserType:      enums.ClientUser,
				Gender:        enumutils.GenderMale,
				Contacts: &domain.Contact{
					ID:           &userID,
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
					Active:       true,
					OptedIn:      true,
				},
			}, nil
		},
		MockSetNickNameFn: func(ctx context.Context, userID, nickname *string) (bool, error) {
			return true, nil
		},
		MockCreateIdentifierFn: func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
			return &domain.Identifier{
				ID:                  ID,
				IdentifierType:      "CCC",
				IdentifierValue:     "123456789",
				IdentifierUse:       "OFFICIAL",
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
		MockGetClientProfileByUserIDFn: func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockGetStaffProfileByUserIDFn: func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
			return staff, nil
		},
		MockUpdateUserSurveysFn: func(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
			return nil
		},
		MockCheckUserHasPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
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
		MockGetUserSecurityQuestionsResponsesFn: func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
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
		MockInvalidatePINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGetContactByUserIDFn: func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
			return &domain.Contact{
				ID:           userID,
				ContactType:  "PHONE",
				ContactValue: gofakeit.Phone(),
				Active:       true,
				OptedIn:      true,
			}, nil
		},
		MockUpdateIsCorrectSecurityQuestionResponseFn: func(ctx context.Context, userID string, isCorrect bool) (bool, error) {
			return true, nil
		},
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
			return nil
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
		MockGetClientCaregiverFn: func(ctx context.Context, caregiverID string) (*domain.Caregiver, error) {
			return &domain.Caregiver{
				ID:            "26b20a42-cbb8-4553-aedb-c539602d04fc",
				FirstName:     "test",
				LastName:      "test",
				PhoneNumber:   "+61412345678",
				CaregiverType: enums.CaregiverTypeFather,
			}, nil
		},
		MockUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
			return nil
		},
		MockCreateCommunityFn: func(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {
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
				InviteOnly: true,
			}, nil
		},
		MockGetClientProfileByClientIDFn: func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
			client := &domain.ClientProfile{
				ID:          &ID,
				User:        userProfile,
				CaregiverID: &ID,
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
		MockCheckUserRoleFn: func(ctx context.Context, userID string, role string) (bool, error) {
			return true, nil
		},
		MockCheckUserPermissionFn: func(ctx context.Context, userID string, permission string) (bool, error) {
			return true, nil
		},
		MockAssignRolesFn: func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
			return true, nil
		},
		MockGetUserRolesFn: func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
			return []*domain.AuthorityRole{
				{
					AuthorityRoleID: uuid.New().String(),
					Name:            enums.UserRoleTypeClientManagement,
				},
			}, nil
		},
		MockGetUserPermissionsFn: func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
			return []*domain.AuthorityPermission{
				{
					PermissionID: uuid.New().String(),
					Name:         enums.PermissionTypeCanManageClient,
				},
			}, nil
		},
		MockCheckIfUsernameExistsFn: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
		MockRevokeRolesFn: func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
			return true, nil
		},
		MockGetCommunityByIDFn: func(ctx context.Context, communityID string) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.New().String(),
				CID:         uuid.New().String(),
				Name:        gofakeit.Name(),
				Disabled:    false,
				Frozen:      false,
				MemberCount: 10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Description: description,
				AgeRange: &domain.AgeRange{
					LowerBound: 0,
					UpperBound: 0,
				},
				Gender:     []enumutils.Gender{},
				ClientType: []enums.ClientType{},
				InviteOnly: false,
				Members:    []domain.CommunityMember{},
				CreatedBy:  &domain.Member{},
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
		MockCheckFacilityExistsByMFLCode: func(ctx context.Context, MFLCode int) (bool, error) {
			return true, nil
		},
		MockGetClientCCCIdentifier: func(ctx context.Context, clientID string) (*domain.Identifier, error) {
			return &domain.Identifier{
				ID:                  uuid.New().String(),
				IdentifierType:      "CCC",
				IdentifierValue:     "123456",
				IdentifierUse:       "OFFICIAL",
				Description:         description,
				ValidFrom:           time.Now(),
				ValidTo:             time.Now(),
				IsPrimaryIdentifier: false,
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
		MockGetScreeningToolsQuestionsFn: func(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error) {
			return []*domain.ScreeningToolQuestion{
				{
					ID:       ID,
					Question: gofakeit.Sentence(1),
					ToolType: enums.ScreeningToolTypeTB,
					ResponseChoices: map[string]interface{}{
						"1": "yes",
						"2": "no",
					},
					ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
					ResponseType:     enums.ScreeningToolResponseTypeInteger,
					Sequence:         1,
					Meta:             nil,
					Active:           true,
				},
			}, nil
		},
		MockAnswerScreeningToolQuestionsFn: func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error {
			return nil
		},
		MockUpdateHealthDiaryFn: func(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
			return nil
		},
		MockGetScreeningToolQuestionByQuestionIDFn: func(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error) {
			return &domain.ScreeningToolQuestion{
				ID:       ID,
				Question: gofakeit.Sentence(1),
				ToolType: enums.ScreeningToolTypeGBV,
				ResponseChoices: map[string]interface{}{
					"0": "yes",
					"1": "no",
				},
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				Sequence:         1,
				Meta: map[string]interface{}{
					"category":             "Violence",
					"category_description": "Response from GBV tool",
					"helper_text":          "Emotional violence Assessment",
					"violence_type":        "EMOTIONAL",
					"violence_code":        "GBV-EV",
				},
				Active: true,
			}, nil
		},
		MockInvalidateScreeningToolResponseFn: func(ctx context.Context, clientID string, questionID string) error {
			return nil
		},
		MockResolveStaffServiceRequestFn: func(ctx context.Context, staffID, serviceRequestID *string, verificationStatus string) (bool, error) {
			return true, nil
		},
		MockUpdateFacilityFn: func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
			return nil
		},
		MockGetClientProfileByCCCNumberFn: func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
			return clientProfile, nil
		},
		MockCheckIfClientHasUnresolvedServiceRequestsFn: func(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
			return false, nil
		},
		MockUpdateUserPinChangeRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockGetAllRolesFn: func(ctx context.Context) ([]*domain.AuthorityRole, error) {
			return []*domain.AuthorityRole{
				{
					AuthorityRoleID: ID,
					Name:            enums.UserRoleTypeClientManagement,
					Active:          true,
				},
			}, nil
		},
		MockUpdateUserPinUpdateRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
			return nil
		},
		MockUpdateClientFn: func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
			return client, nil
		},
		MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn: func(ctx context.Context, userID string, failCount int) error {
			return nil
		},
		MockGetAssessmentResponsesFn: func(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
			return []*domain.ScreeningToolAssessmentResponse{
				{
					ClientName:   name,
					DateAnswered: time.Now(),
					ClientID:     ID,
				},
			}, nil
		},
		MockGetServiceRequestByIDFn: func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
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
		MockGetActiveScreeningToolResponsesFn: func(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error) {
			return []*domain.ScreeningToolQuestionResponse{
				{
					ID:         ID,
					QuestionID: ID,
					ClientID:   clientID,
					Answer:     "0",
					Active:     true,
				},
			}, nil
		},
		MockGetClientScreeningToolResponsesByToolTypeFn: func(ctx context.Context, clientID, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error) {
			return []*domain.ScreeningToolQuestionResponse{
				{
					ID:         ID,
					QuestionID: ID,
					ClientID:   clientID,
					Answer:     "0",
					Active:     true,
				},
			}, nil
		},
		MockGetUserSurveyFormsFn: func(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
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
		MockCreateClientFn: func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
			return clientProfile, nil
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

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (gm *PostgresMock) ListFacilities(
	ctx context.Context,
	searchTerm *string,
	filterInput []*dto.FiltersInput,
	paginationsInput *dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return gm.ListFacilitiesFn(ctx, searchTerm, filterInput, paginationsInput)
}

// SearchFacility mocks the implementation of `gorm's` GetFacilities method
func (gm *PostgresMock) SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
	return gm.MockSearchFacilityFn(ctx, searchParameter)
}

// DeleteFacility mocks the implementation of deleting a facility by ID
func (gm *PostgresMock) DeleteFacility(ctx context.Context, id int) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, id)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacilityByMFLCode method.
func (gm *PostgresMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (gm *PostgresMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber, flavour)
}

// GetUserPINByUserID mocks the get user pin by ID implementation
func (gm *PostgresMock) GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID, flavour)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *PostgresMock) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, mflCode)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (gm *PostgresMock) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockReactivateFacilityFn(ctx, mflCode)
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *PostgresMock) GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx, flavour)
}

// GetUserProfileByUserID mocks the implementation of fetching a user profile by userID
func (gm *PostgresMock) GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, userID)
}

// SaveTemporaryUserPin mocks the implementation of saving a temporary user pin
func (gm *PostgresMock) SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return gm.MockSaveTemporaryUserPinFn(ctx, pinData)
}

// AcceptTerms mocks the implementation of accept current terms of service
func (gm *PostgresMock) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	return gm.MockAcceptTermsFn(ctx, userID, termsID)
}

// SavePin mocks the implementation of saving a user pin
func (gm *PostgresMock) SavePin(ctx context.Context, pin *domain.UserPIN) (bool, error) {
	return gm.MockSavePinFn(ctx, pin)
}

// UpdateUserSurveys mocks the implementation of `gorm's` UpdateUserSurveys method.
func (gm *PostgresMock) UpdateUserSurveys(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
	return gm.MockUpdateUserSurveysFn(ctx, survey, updateData)
}

//GetSecurityQuestions mocks the implementation of getting all the security questions.
func (gm *PostgresMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionsFn(ctx, flavour)
}

// SaveOTP mocks the implementation for saving an OTP
func (gm *PostgresMock) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	return gm.MockSaveOTPFn(ctx, otpInput)
}

// SetNickName is used to mock the implementation offset or changing the user's nickname
func (gm *PostgresMock) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	return gm.MockSetNickNameFn(ctx, userID, nickname)
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

// GetClientProfileByUserID mocks the method for fetching a client profile using the user ID
func (gm *PostgresMock) GetClientProfileByUserID(ctx context.Context, userID string) (*domain.ClientProfile, error) {
	return gm.MockGetClientProfileByUserIDFn(ctx, userID)
}

// GetStaffProfileByUserID mocks the method for fetching a staff profile using the user ID
func (gm *PostgresMock) GetStaffProfileByUserID(ctx context.Context, userID string) (*domain.StaffProfile, error) {
	return gm.MockGetStaffProfileByUserIDFn(ctx, userID)
}

// SearchStaffProfile mocks the implementation of getting staff profile using their staff number.
func (gm *PostgresMock) SearchStaffProfile(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	return gm.MockSearchStaffProfileFn(ctx, searchParameter)
}

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *PostgresMock) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID, flavour)
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

// GetUserSecurityQuestionsResponses mocks the implementation of getting the user's responded security questions
func (gm *PostgresMock) GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
	return gm.MockGetUserSecurityQuestionsResponsesFn(ctx, userID)
}

// InvalidatePIN mocks the implementation of invalidating a user pin
func (gm *PostgresMock) InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockInvalidatePINFn(ctx, userID, flavour)
}

// GetContactByUserID mocks the implementation of fetching a contact by userID
func (gm *PostgresMock) GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
	return gm.MockGetContactByUserIDFn(ctx, userID, contactType)
}

// UpdateIsCorrectSecurityQuestionResponse updates the IsCorrectSecurityQuestion response
func (gm *PostgresMock) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	return gm.MockUpdateIsCorrectSecurityQuestionResponseFn(ctx, userID, isCorrectSecurityQuestionResponse)
}

//FetchFacilities mocks the implementation of fetching facility
func (gm *PostgresMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockFetchFacilitiesFn(ctx)
}

// CreateHealthDiaryEntry mocks the method for creating a health diary entry
func (gm *PostgresMock) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
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

// CheckUserRole mocks the implementation of checking if a user has a role
func (gm *PostgresMock) CheckUserRole(ctx context.Context, userID string, role string) (bool, error) {
	return gm.MockCheckUserRoleFn(ctx, userID, role)
}

// CheckUserPermission mocks the implementation of checking if a user has a permission
func (gm *PostgresMock) CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error) {
	return gm.MockCheckUserPermissionFn(ctx, userID, permission)
}

// AssignRoles mocks the implementation of assigning roles to a user
func (gm *PostgresMock) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return gm.MockAssignRolesFn(ctx, userID, roles)
}

// CreateCommunity mocks the implementation of creating a channel
func (gm *PostgresMock) CreateCommunity(ctx context.Context, community *dto.CommunityInput) (*domain.Community, error) {
	return gm.MockCreateCommunityFn(ctx, community)
}

// GetUserRoles mocks the implementation of getting all roles for a user
func (gm *PostgresMock) GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
	return gm.MockGetUserRolesFn(ctx, userID)
}

// GetUserPermissions mocks the implementation of getting all permissions for a user
func (gm *PostgresMock) GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
	return gm.MockGetUserPermissionsFn(ctx, userID)
}

// CheckIfUsernameExists mocks the implementation of checking whether a username exists
func (gm *PostgresMock) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	return gm.MockCheckIfUsernameExistsFn(ctx, username)
}

// RevokeRoles mocks the implementation of revoking roles from a user
func (gm *PostgresMock) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return gm.MockRevokeRolesFn(ctx, userID, roles)
}

// GetCommunityByID mocks the implementation of getting the community by ID
func (gm *PostgresMock) GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error) {
	return gm.MockGetCommunityByIDFn(ctx, communityID)
}

// CheckIdentifierExists mocks checking an identifier exists
func (gm *PostgresMock) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	return gm.MockCheckIdentifierExists(ctx, identifierType, identifierValue)
}

// CheckFacilityExistsByMFLCode mocks checking a facility by mfl codes
func (gm *PostgresMock) CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error) {
	return gm.MockCheckFacilityExistsByMFLCode(ctx, MFLCode)
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

// GetClientCCCIdentifier retrieves client's ccc number
func (gm *PostgresMock) GetClientCCCIdentifier(ctx context.Context, clientID string) (*domain.Identifier, error) {
	return gm.MockGetClientCCCIdentifier(ctx, clientID)
}

// GetServiceRequestsForKenyaEMR mocks the getting of red flag service requests for use by KenyaEMR
func (gm *PostgresMock) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
	return gm.MockGetServiceRequestsForKenyaEMRFn(ctx, payload)
}

// GetScreeningToolQuestions mocks the implementation of getting screening tools questions
func (gm *PostgresMock) GetScreeningToolQuestions(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error) {
	return gm.MockGetScreeningToolsQuestionsFn(ctx, toolType)
}

// AnswerScreeningToolQuestions mocks the implementation of answering screening tool questions
func (gm *PostgresMock) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error {
	return gm.MockAnswerScreeningToolQuestionsFn(ctx, screeningToolResponses)
}

// GetScreeningToolQuestionByQuestionID mocks the implementation of getting screening tool questions by question ID
func (gm *PostgresMock) GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error) {
	return gm.MockGetScreeningToolQuestionByQuestionIDFn(ctx, questionID)
}

// InvalidateScreeningToolResponse mocks the implementation of invalidating screening tool responses
func (gm *PostgresMock) InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error {
	return gm.MockInvalidateScreeningToolResponseFn(ctx, clientID, questionID)
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

// GetClientProfileByCCCNumber mocks the implementation of getting a client profile using the CCC number
func (gm *PostgresMock) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
	return gm.MockGetClientProfileByCCCNumberFn(ctx, CCCNumber)
}

// CheckIfClientHasUnresolvedServiceRequests mocks the implementation of checking if a client has an unresolved service request
func (gm *PostgresMock) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfClientHasUnresolvedServiceRequestsFn(ctx, clientID, serviceRequestType)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation of updating a user pin change required state
func (gm *PostgresMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour, status)
}

// GetAllRoles mocks the implementation of getting all roles
func (gm *PostgresMock) GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error) {
	return gm.MockGetAllRolesFn(ctx)
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

// GetServiceRequestByID mocks the implementation of getting a service request by ID
func (gm *PostgresMock) GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	return gm.MockGetServiceRequestByIDFn(ctx, serviceRequestID)
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
func (gm *PostgresMock) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
	return gm.MockGetUserSurveyFormsFn(ctx, userID)
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

// GetActiveScreeningToolResponses mocks the implementation of getting active screening tool responses
func (gm *PostgresMock) GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error) {
	return gm.MockGetActiveScreeningToolResponsesFn(ctx, clientID)
}

// GetAssessmentResponses mocks the implementation of getting answered screening tool questions
func (gm *PostgresMock) GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
	return gm.MockGetAssessmentResponsesFn(ctx, facilityID, toolType)
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

//ListAvailableNotificationTypes retrieves the distinct notification types available for a user
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

// GetClientScreeningToolResponsesByToolType mocks the implementation of getting client screening tool responses by tool type
func (gm *PostgresMock) GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID string, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error) {
	return gm.MockGetClientScreeningToolResponsesByToolTypeFn(ctx, clientID, toolType, active)
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
