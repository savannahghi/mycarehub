package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	"github.com/segmentio/ksuid"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockCreateUserFn                                     func(ctx context.Context, user *gorm.User) error
	MockCreateClientFn                                   func(ctx context.Context, client *gorm.Client, contactID, identifierID string) error
	MockCreateIdentifierFn                               func(ctx context.Context, identifier *gorm.Identifier) error
	MockGetOrCreateFacilityFn                            func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn                               func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn                      func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error)
	MockSearchFacilityFn                                 func(ctx context.Context, searchParameter *string) ([]gorm.Facility, error)
	MockDeleteFacilityFn                                 func(ctx context.Context, mflCode int) (bool, error)
	MockListFacilitiesFn                                 func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	MockGetUserProfileByPhoneNumberFn                    func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error)
	MockGetUserPINByUserIDFn                             func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error)
	MockInactivateFacilityFn                             func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn                             func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn                         func(ctx context.Context, userID *string) (*gorm.User, error)
	MockSaveTemporaryUserPinFn                           func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetCurrentTermsFn                                func(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error)
	MockAcceptTermsFn                                    func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                        func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockSetNickNameFn                                    func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockGetSecurityQuestionsFn                           func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error)
	MockSaveOTPFn                                        func(ctx context.Context, otpInput *gorm.UserOTP) error
	MockGetSecurityQuestionByIDFn                        func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn                   func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error
	MockGetSecurityQuestionResponseFn                    func(ctx context.Context, questionID string, userID string) (*gorm.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                       func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                                      func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetClientProfileByUserIDFn                       func(ctx context.Context, userID string) (*gorm.Client, error)
	MockGetStaffProfileByUserIDFn                        func(ctx context.Context, userID string) (*gorm.StaffProfile, error)
	MockCheckUserHasPinFn                                func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockCompleteOnboardingTourFn                         func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                         func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error)
	MockGetUserSecurityQuestionsResponsesFn              func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error)
	MockInvalidatePINFn                                  func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetContactByUserIDFn                             func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn        func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockCreateHealthDiaryEntryFn                         func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) error
	MockCreateServiceRequestFn                           func(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error
	MockCanRecordHeathDiaryFn                            func(ctx context.Context, clientID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn                      func(ctx context.Context, limit int) ([]*gorm.ClientHealthDiaryQuote, error)
	MockGetClientHealthDiaryEntriesFn                    func(ctx context.Context, params map[string]interface{}) ([]*gorm.ClientHealthDiaryEntry, error)
	MockCreateClientCaregiverFn                          func(ctx context.Context, clientID string, clientCaregiver *gorm.Caregiver) error
	MockGetClientCaregiverFn                             func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error)
	MockUpdateClientCaregiverFn                          func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockInProgressByFn                                   func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn                     func(ctx context.Context, clientID string) (*gorm.Client, error)
	MockGetServiceRequestsFn                             func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockGetPendingServiceRequestsCountFn                 func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockCheckUserRoleFn                                  func(ctx context.Context, userID string, role string) (bool, error)
	MockCheckUserPermissionFn                            func(ctx context.Context, userID string, permission string) (bool, error)
	MockAssignRolesFn                                    func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockCreateCommunityFn                                func(ctx context.Context, community *gorm.Community) (*gorm.Community, error)
	MockGetUserRolesFn                                   func(ctx context.Context, userID string) ([]*gorm.AuthorityRole, error)
	MockGetUserPermissionsFn                             func(ctx context.Context, userID string) ([]*gorm.AuthorityPermission, error)
	MockCheckIfUsernameExistsFn                          func(ctx context.Context, username string) (bool, error)
	MockRevokeRolesFn                                    func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetCommunityByIDFn                               func(ctx context.Context, communityID string) (*gorm.Community, error)
	MockCheckIdentifierExists                            func(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	MockCheckFacilityExistsByMFLCode                     func(ctx context.Context, MFLCode int) (bool, error)
	MockGetOrCreateNextOfKin                             func(ctx context.Context, person *gorm.RelatedPerson, clientID, contactID string) error
	MockGetOrCreateContact                               func(ctx context.Context, contact *gorm.Contact) (*gorm.Contact, error)
	MockGetClientsInAFacilityFn                          func(ctx context.Context, facilityID string) ([]*gorm.Client, error)
	MockGetRecentHealthDiaryEntriesFn                    func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                               func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error)
	MockGetClientCCCIdentifier                           func(ctx context.Context, clientID string) (*gorm.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn                  func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error)
	MockCreateAppointment                                func(ctx context.Context, appointment *gorm.Appointment) error
	MockListAppointments                                 func(ctx context.Context, params *gorm.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error)
	MockUpdateAppointmentFn                              func(ctx context.Context, appointment *gorm.Appointment, updateData map[string]interface{}) (*gorm.Appointment, error)
	MockGetScreeningToolsQuestionsFn                     func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error)
	MockAnswerScreeningToolQuestionsFn                   func(ctx context.Context, screeningToolResponses []*gorm.ScreeningToolsResponse) error
	MockGetScreeningToolQuestionByQuestionIDFn           func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error)
	MockInvalidateScreeningToolResponseFn                func(ctx context.Context, clientID string, questionID string) error
	MockUpdateServiceRequestsFn                          func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error)
	MockGetClientProfileByCCCNumberFn                    func(ctx context.Context, CCCNumber string) (*gorm.Client, error)
	MockSearchClientProfileFn                            func(ctx context.Context, searchParameter string) ([]*gorm.Client, error)
	MockSearchStaffProfileFn                             func(ctx context.Context, searchParameter string) ([]*gorm.StaffProfile, error)
	MockUpdateUserPinChangeRequiredStatusFn              func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockCheckIfClientHasUnresolvedServiceRequestsFn      func(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	MockGetAllRolesFn                                    func(ctx context.Context) ([]*gorm.AuthorityRole, error)
	MockUpdateHealthDiaryFn                              func(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error
	MockUpdateUserPinUpdateRequiredStatusFn              func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockUpdateClientFn                                   func(ctx context.Context, client *gorm.Client, updates map[string]interface{}) (*gorm.Client, error)
	MockGetUserProfileByStaffIDFn                        func(ctx context.Context, staffID string) (*gorm.User, error)
	MockGetHealthDiaryEntryByIDFn                        func(ctx context.Context, healthDiaryEntryID string) (*gorm.ClientHealthDiaryEntry, error)
	MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn func(ctx context.Context, userID string, failCount int) error
	MockGetServiceRequestByIDFn                          func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error)
	MockUpdateUserFn                                     func(ctx context.Context, user *gorm.User, updateData map[string]interface{}) error
	MockGetStaffProfileByStaffIDFn                       func(ctx context.Context, staffID string) (*gorm.StaffProfile, error)
	MockCreateStaffServiceRequestFn                      func(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error
	MockGetStaffPendingServiceRequestsCountFn            func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockGetStaffServiceRequestsFn                        func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.StaffServiceRequest, error)
	MockResolveStaffServiceRequestFn                     func(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error)
	MockGetAppointmentServiceRequestsFn                  func(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockUpdateFacilityFn                                 func(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error
	MockGetFacilitiesWithoutFHIRIDFn                     func(ctx context.Context) ([]*gorm.Facility, error)
	MockGetSharedHealthDiaryEntriesFn                    func(ctx context.Context, clientID string, facilityID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetClientServiceRequestsFn                       func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockGetActiveScreeningToolResponsesFn                func(ctx context.Context, clientID string) ([]*gorm.ScreeningToolsResponse, error)
	MockCheckAppointmentExistsByExternalIDFn             func(ctx context.Context, externalID string) (bool, error)
	MockGetUserSurveyFormsFn                             func(ctx context.Context, userID string) ([]*gorm.UserSurvey, error)
	MockGetAnsweredScreeningToolQuestionsFn              func(ctx context.Context, facilityID string, toolType string) ([]*gorm.ScreeningToolsResponse, error)
	MockCreateNotificationFn                             func(ctx context.Context, notification *gorm.Notification) error
	MockUpdateUserSurveysFn                              func(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error
	MockSearchClientServiceRequestsFn                    func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.ClientServiceRequest, error)
	MockSearchStaffServiceRequestsFn                     func(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*gorm.StaffServiceRequest, error)
	MockListNotificationsFn                              func(ctx context.Context, params *gorm.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*gorm.Notification, *domain.Pagination, error)
	MockListAvailableNotificationTypesFn                 func(ctx context.Context, params *gorm.Notification) ([]enums.NotificationType, error)
	MockGetClientScreeningToolResponsesByToolTypeFn      func(ctx context.Context, clientID, toolType string, active bool) ([]*gorm.ScreeningToolsResponse, error)
	MockGetClientScreeningToolServiceRequestByToolTypeFn func(ctx context.Context, clientID, toolType, status string) (*gorm.ClientServiceRequest, error)
	MockGetAppointmentFn                                 func(ctx context.Context, params *gorm.Appointment) (*gorm.Appointment, error)
	MockCheckIfStaffHasUnresolvedServiceRequestsFn       func(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	MockGetFacilityStaffsFn                              func(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error)
	MockDeleteUserFn                                     func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error
	MockDeleteStaffProfileFn                             func(ctx context.Context, staffID string) error
	MockSaveFeedbackFn                                   func(ctx context.Context, feedback *gorm.Feedback) error
	MockUpdateNotificationFn                             func(ctx context.Context, notification *gorm.Notification, updateData map[string]interface{}) error
	MockGetNotificationFn                                func(ctx context.Context, notificationID string) (*gorm.Notification, error)
	MockGetClientsByFilterParamsFn                       func(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*gorm.Client, error)
	MockCreateUserSurveyFn                               func(ctx context.Context, userSurvey []*gorm.UserSurvey) error
	MockCreateMetricFn                                   func(ctx context.Context, metric *gorm.Metric) error
	MockRegisterStaffFn                                  func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) error
	MockUpdateClientServiceRequestFn                     func(ctx context.Context, clientServiceRequest *gorm.ClientServiceRequest, updateData map[string]interface{}) error
	MockRegisterClientFn                                 func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) error
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
//
// This initialization initializes all the good cases of your mock tests. i.e all success cases should be defined here.
func NewGormMock() *GormMock {

	/*
		In this section, you find commonly shared success case structs for mock tests
	*/

	ID := gofakeit.Number(300, 400)
	UUID := ksuid.New().String()
	name := gofakeit.Name()
	code := 1234567890
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)
	phoneContact := gofakeit.Phone()
	acceptedTermsID := gofakeit.Number(1, 10)
	currentTime := time.Now()

	facility := &gorm.Facility{
		FacilityID:  &UUID,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Phone:       phoneContact,
		Description: description,
	}

	var facilities []gorm.Facility
	facilities = append(facilities, *facility)

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
				ID:          &UUID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			},
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
			FirstName:              gofakeit.Name(),
			MiddleName:             name,
			LastName:               gofakeit.Name(),
			UserType:               enums.HealthcareWorkerUser,
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
			Flavour:                feedlib.FlavourPro,
			Avatar:                 gofakeit.URL(),
			IsSuspended:            true,
			PinChangeRequired:      true,
			HasSetPin:              true,
			HasSetSecurityQuestion: true,
			IsPhoneVerified:        true,
			OrganisationID:         uuid.New().String(),
			Password:               gofakeit.Name(),
			IsSuperuser:            true,
			IsStaff:                true,
			Email:                  gofakeit.Email(),
			DateJoined:             gofakeit.BeerIbu(),
			Name:                   name,
			IsApproved:             true,
			ApprovalNotified:       true,
			Handle:                 "@test",
			DateOfBirth:            &currentTime,
		},
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &fhirID,
		HealthRecordID:          &UUID,
		TreatmentBuddy:          gofakeit.Name(),
		ClientCounselled:        true,
		OrganisationID:          uuid.New().String(),
		FacilityID:              uuid.New().String(),
		CHVUserID:               &UUID,
		UserID:                  &UUID,
		CaregiverID:             &UUID,
	}

	userProfile := &gorm.User{
		UserID:                 &UUID,
		Username:               gofakeit.Name(),
		FirstName:              gofakeit.Name(),
		MiddleName:             name,
		LastName:               gofakeit.Name(),
		UserType:               enums.HealthcareWorkerUser,
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
		Flavour:                feedlib.FlavourPro,
		Avatar:                 "test",
		IsSuspended:            true,
		PinChangeRequired:      true,
		HasSetPin:              true,
		HasSetSecurityQuestion: true,
		IsPhoneVerified:        true,
		OrganisationID:         uuid.New().String(),
		Password:               "test",
		IsSuperuser:            true,
		IsStaff:                true,
		Email:                  gofakeit.Email(),
		DateJoined:             gofakeit.BeerIbu(),
		Name:                   name,
		IsApproved:             true,
		ApprovalNotified:       true,
		Handle:                 "@test",
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
		Flavour:   feedlib.FlavourConsumer,
	}

	nowTime := time.Now()
	laterTime := nowTime.Add(time.Hour * 24)
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
		},
	}

	return &GormMock{
		MockCreateMetricFn: func(ctx context.Context, metric *gorm.Metric) error {
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
				Flavour:    "PRO",
				IsRead:     false,
				UserID:     &UUID,
				FacilityID: &UUID,
			}, nil
		},
		MockUpdateNotificationFn: func(ctx context.Context, notification *gorm.Notification, updateData map[string]interface{}) error {
			return nil
		},
		MockGetFacilityStaffsFn: func(ctx context.Context, facilityID string) ([]*gorm.StaffProfile, error) {
			return []*gorm.StaffProfile{
				staff,
			}, nil
		},
		MockCreateUserFn: func(ctx context.Context, user *gorm.User) error {
			return nil
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
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			return facility, nil
		},
		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {

			return facility, nil
		},
		MockGetStaffPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
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
		MockSearchFacilityFn: func(ctx context.Context, searchParameter *string) ([]gorm.Facility, error) {
			return facilities, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, mflCode int) (bool, error) {
			return true, nil
		},

		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
			return facility, nil
		},
		MockRegisterStaffFn: func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) error {
			return nil
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
		MockGetAnsweredScreeningToolQuestionsFn: func(ctx context.Context, facilityID string, toolType string) ([]*gorm.ScreeningToolsResponse, error) {
			return []*gorm.ScreeningToolsResponse{
				{
					ID:             fhirID,
					ClientID:       uuid.New().String(),
					QuestionID:     uuid.New().String(),
					Response:       uuid.New().String(),
					Active:         true,
					OrganisationID: uuid.New().String(),
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
		MockRegisterClientFn: func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) error {
			return nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},
		MockGetUserSurveyFormsFn: func(ctx context.Context, userID string) ([]*gorm.UserSurvey, error) {
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

		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},

		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error) {
			return pinData, nil
		},

		MockInactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
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
		MockGetCurrentTermsFn: func(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error) {
			termsID := gofakeit.Number(1, 1000)
			validFrom := time.Now()
			testText := "test"

			validTo := time.Now().AddDate(0, 0, 80)
			terms := &gorm.TermsOfService{
				Base:      gorm.Base{},
				TermsID:   &termsID,
				Text:      &testText,
				Flavour:   feedlib.FlavourPro,
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
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			}
			return []*gorm.SecurityQuestion{securityQuestion}, nil
		},
		MockSaveOTPFn: func(ctx context.Context, otpInput *gorm.UserOTP) error {
			return nil
		},
		MockSetNickNameFn: func(ctx context.Context, userID, nickname *string) (bool, error) {
			return true, nil
		},
		MockGetSecurityQuestionByIDFn: func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error) {
			return &gorm.SecurityQuestion{
				SecurityQuestionID: &UUID,
				QuestionStem:       "test",
				Description:        "test",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			}, nil
		},
		MockSaveSecurityQuestionResponseFn: func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error {
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
		MockGetClientProfileByUserIDFn: func(ctx context.Context, userID string) (*gorm.Client, error) {
			return clientProfile, nil
		},
		MockGetStaffProfileByUserIDFn: func(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockSearchStaffProfileFn: func(ctx context.Context, staffNumber string) ([]*gorm.StaffProfile, error) {
			return []*gorm.StaffProfile{staff}, nil
		},
		MockCheckUserHasPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
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
		MockGetUserSecurityQuestionsResponsesFn: func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error) {
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
		MockInvalidatePINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGetContactByUserIDFn: func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error) {
			return &gorm.Contact{
				ContactID:    &UUID,
				UserID:       userID,
				ContactType:  "PHONE",
				ContactValue: phoneContact,
				Active:       true,
				OptedIn:      true,
			}, nil
		},
		MockUpdateIsCorrectSecurityQuestionResponseFn: func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
			return true, nil
		},
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) error {
			return nil
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
		MockCreateClientCaregiverFn: func(ctx context.Context, clientID string, clientCaregiver *gorm.Caregiver) error {
			return nil
		},
		MockGetPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
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
		MockGetClientCaregiverFn: func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error) {
			ID := uuid.New().String()
			return &gorm.Caregiver{
				CaregiverID:   &ID,
				FirstName:     "test",
				LastName:      "test",
				PhoneNumber:   gofakeit.Phone(),
				CaregiverType: enums.CaregiverTypeFather,
				Active:        true,
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
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
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
				InviteOnly:     true,
				Discoverable:   true,
				OrganisationID: uuid.New().String(),
			}, nil
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
		MockGetUserRolesFn: func(ctx context.Context, userID string) ([]*gorm.AuthorityRole, error) {
			return []*gorm.AuthorityRole{
				{
					AuthorityRoleID: &UUID,
					Name:            enums.UserRoleTypeClientManagement.String(),
				},
			}, nil
		},
		MockGetUserPermissionsFn: func(ctx context.Context, userID string) ([]*gorm.AuthorityPermission, error) {
			return []*gorm.AuthorityPermission{
				{
					AuthorityPermissionID: &UUID,
					Name:                  enums.PermissionTypeCanCreateGroup.String(),
				},
			}, nil
		},
		MockCheckIfUsernameExistsFn: func(ctx context.Context, username string) (bool, error) {
			return true, nil
		},
		MockRevokeRolesFn: func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
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
				InviteOnly:     false,
				Discoverable:   false,
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
		MockCheckFacilityExistsByMFLCode: func(ctx context.Context, MFLCode int) (bool, error) {
			return true, nil
		},
		MockCheckIdentifierExists: func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
			return true, nil
		},
		MockGetClientsByParams: func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error) {
			return []*gorm.Client{clientProfile}, nil
		},
		MockGetClientCCCIdentifier: func(ctx context.Context, clientID string) (*gorm.Identifier, error) {
			return &gorm.Identifier{
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
					Flavour:    "PRO",
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
		MockGetScreeningToolsQuestionsFn: func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error) {
			return []gorm.ScreeningToolQuestion{
				{
					ID:               UUID,
					Question:         gofakeit.Sentence(1),
					ToolType:         enums.ScreeningToolTypeTB.String(),
					ResponseChoices:  `{"1": "Yes", "2": "No"}`,
					ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
					ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
					Sequence:         1,
					Active:           true,
					Meta:             `{"meta": "data"}`,
					OrganisationID:   uuid.New().String(),
				},
			}, nil
		},
		MockAnswerScreeningToolQuestionsFn: func(ctx context.Context, screeningToolResponses []*gorm.ScreeningToolsResponse) error {
			return nil
		},
		MockGetScreeningToolQuestionByQuestionIDFn: func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error) {
			return &gorm.ScreeningToolQuestion{
				ID:               UUID,
				Question:         gofakeit.Sentence(1),
				ToolType:         enums.ScreeningToolTypeGBV.String(),
				ResponseChoices:  `{"O": "Yes", "1": "No"}`,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice.String(),
				ResponseType:     enums.ScreeningToolResponseTypeInteger.String(),
				Sequence:         1,
				Active:           true,
				Meta:             `{"meta": "data"}`,
				OrganisationID:   uuid.New().String(),
			}, nil
		},
		MockInvalidateScreeningToolResponseFn: func(ctx context.Context, clientID string, questionID string) error {
			return nil
		},
		MockGetClientProfileByCCCNumberFn: func(ctx context.Context, CCCNumber string) (*gorm.Client, error) {
			return clientProfile, nil
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
		MockGetAllRolesFn: func(ctx context.Context) ([]*gorm.AuthorityRole, error) {
			return []*gorm.AuthorityRole{
				{
					AuthorityRoleID: &UUID,
					Name:            enums.UserRoleTypeClientManagement.String(),
					Active:          true,
				},
			}, nil
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
		MockGetServiceRequestByIDFn: func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
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
		MockGetActiveScreeningToolResponsesFn: func(ctx context.Context, clientID string) ([]*gorm.ScreeningToolsResponse, error) {
			return []*gorm.ScreeningToolsResponse{
				{
					Base:           gorm.Base{},
					ID:             UUID,
					ClientID:       uuid.New().String(),
					QuestionID:     "",
					Response:       "",
					Active:         true,
					OrganisationID: "",
				},
			}, nil
		},
		MockGetClientScreeningToolResponsesByToolTypeFn: func(ctx context.Context, clientID, toolType string, active bool) ([]*gorm.ScreeningToolsResponse, error) {
			return []*gorm.ScreeningToolsResponse{
				{
					Base:           gorm.Base{},
					ID:             UUID,
					ClientID:       uuid.New().String(),
					QuestionID:     "",
					Response:       "",
					Active:         true,
					OrganisationID: "",
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
		MockDeleteUserFn: func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
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
	}
}

// DeleteStaffProfile mocks the implementation of deleting a staff
func (gm *GormMock) DeleteStaffProfile(ctx context.Context, staffID string) error {
	return gm.MockDeleteStaffProfileFn(ctx, staffID)
}

// DeleteUser mocks the implementation of deleting a user
func (gm *GormMock) DeleteUser(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
	return gm.MockDeleteUserFn(ctx, userID, clientID, staffID, flavour)
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// UpdateUserSurveys mocks the implementation of `gorm's` UpdateUserSurveys method.
func (gm *GormMock) UpdateUserSurveys(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error {
	return gm.MockUpdateUserSurveysFn(ctx, survey, updateData)
}

// SearchFacility mocks the implementation of `gorm's` SearchFacility method.
func (gm *GormMock) SearchFacility(ctx context.Context, searchParameter *string) ([]gorm.Facility, error) {
	return gm.MockSearchFacilityFn(ctx, searchParameter)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode int) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, mflcode)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (gm *GormMock) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
	return gm.MockListFacilitiesFn(ctx, searchTerm, filter, pagination)
}

// GetUserProfileByPhoneNumber mocks the implementation of retrieving a user profile by phonenumber
func (gm *GormMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber, flavour)
}

// GetUserPINByUserID mocks the implementation of retrieving a user pin by user ID
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID, flavour)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *GormMock) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, mflCode)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (gm *GormMock) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockReactivateFacilityFn(ctx, mflCode)
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *GormMock) GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx, flavour)
}

// GetUserProfileByUserID mocks the implementation of retrieving a user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, userID *string) (*gorm.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, userID)
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

//GetSecurityQuestions mocks the implementation of getting all the security questions.
func (gm *GormMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionsFn(ctx, flavour)
}

// SaveOTP mocks the implementation for saving an OTP
func (gm *GormMock) SaveOTP(ctx context.Context, otpInput *gorm.UserOTP) error {
	return gm.MockSaveOTPFn(ctx, otpInput)
}

// SetNickName is used to mock the implementation ofsetting or changing the user's nickname
func (gm *GormMock) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	return gm.MockSetNickNameFn(ctx, userID, nickname)
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

// VerifyOTP mocks the implementation of verify otp
func (gm *GormMock) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return gm.MockVerifyOTPFn(ctx, payload)
}

// GetClientProfileByUserID mocks the method for fetching a client profile using the user ID
func (gm *GormMock) GetClientProfileByUserID(ctx context.Context, userID string) (*gorm.Client, error) {
	return gm.MockGetClientProfileByUserIDFn(ctx, userID)
}

// GetStaffProfileByUserID mocks the method for fetching a staff profile using the user ID
func (gm *GormMock) GetStaffProfileByUserID(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
	return gm.MockGetStaffProfileByUserIDFn(ctx, userID)
}

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *GormMock) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID, flavour)
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
func (gm *GormMock) GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error) {
	return gm.MockGetUserSecurityQuestionsResponsesFn(ctx, userID)
}

// InvalidatePIN mocks the implementation of invalidating the pin
func (gm *GormMock) InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockInvalidatePINFn(ctx, userID, flavour)
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
func (gm *GormMock) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) error {
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

// CreateClientCaregiver mocks the implementation of creating a caregiver
func (gm *GormMock) CreateClientCaregiver(ctx context.Context, clientID string, caregiver *gorm.Caregiver) error {
	return gm.MockCreateClientCaregiverFn(ctx, clientID, caregiver)
}

// GetClientCaregiver mocks the implementation of getting a caregiver
func (gm *GormMock) GetClientCaregiver(ctx context.Context, caregiverID string) (*gorm.Caregiver, error) {
	return gm.MockGetClientCaregiverFn(ctx, caregiverID)
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
func (gm *GormMock) GetClientsPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetPendingServiceRequestsCountFn(ctx, facilityID)
}

// GetServiceRequests mocks the implementation of getting service requests by type
func (gm *GormMock) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID)
}

// CheckUserRole mocks the implementation of checking if a user has a role
func (gm *GormMock) CheckUserRole(ctx context.Context, userID string, role string) (bool, error) {
	return gm.MockCheckUserRoleFn(ctx, userID, role)
}

// CheckUserPermission mocks the implementation of checking if a user has a permission
func (gm *GormMock) CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error) {
	return gm.MockCheckUserPermissionFn(ctx, userID, permission)
}

// AssignRoles mocks the implementation of assigning roles to a user
func (gm *GormMock) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return gm.MockAssignRolesFn(ctx, userID, roles)
}

// CreateCommunity mocks the implementation of creating a channel
func (gm *GormMock) CreateCommunity(ctx context.Context, community *gorm.Community) (*gorm.Community, error) {
	return gm.MockCreateCommunityFn(ctx, community)
}

// GetUserRoles mocks the implementation of getting a user's roles
func (gm *GormMock) GetUserRoles(ctx context.Context, userID string) ([]*gorm.AuthorityRole, error) {
	return gm.MockGetUserRolesFn(ctx, userID)
}

// GetUserPermissions mocks the implementation of getting a user's permissions
func (gm *GormMock) GetUserPermissions(ctx context.Context, userID string) ([]*gorm.AuthorityPermission, error) {
	return gm.MockGetUserPermissionsFn(ctx, userID)
}

// CheckIfUsernameExists mocks the implementation of checking whether a username exists
func (gm *GormMock) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	return gm.MockCheckIfUsernameExistsFn(ctx, username)
}

// RevokeRoles mocks the implementation of revoking roles from a user
func (gm *GormMock) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return gm.MockRevokeRolesFn(ctx, userID, roles)
}

// GetCommunityByID mocks the implementation of getting the community by ID
func (gm *GormMock) GetCommunityByID(ctx context.Context, communityID string) (*gorm.Community, error) {
	return gm.MockGetCommunityByIDFn(ctx, communityID)
}

// CheckIdentifierExists mocks checking of identifiers
func (gm *GormMock) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	return gm.MockCheckIdentifierExists(ctx, identifierType, identifierValue)
}

// CheckFacilityExistsByMFLCode mocks checking a facility by MFL Code
func (gm *GormMock) CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error) {
	return gm.MockCheckFacilityExistsByMFLCode(ctx, MFLCode)
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

// GetClientCCCIdentifier retrieves a client's ccc identifier
func (gm *GormMock) GetClientCCCIdentifier(ctx context.Context, clientID string) (*gorm.Identifier, error) {
	return gm.MockGetClientCCCIdentifier(ctx, clientID)
}

// GetServiceRequestsForKenyaEMR mocks the getting of service requests attached to a specific facility for use by KenyaEMR
func (gm *GormMock) GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestsForKenyaEMRFn(ctx, facilityID, lastSyncTime)
}

// GetScreeningToolQuestions mocks the implementation of getting screening tools questions
func (gm *GormMock) GetScreeningToolQuestions(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error) {
	return gm.MockGetScreeningToolsQuestionsFn(ctx, toolType)
}

// AnswerScreeningToolQuestions mocks the implementation of answering screening tool questions
func (gm *GormMock) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*gorm.ScreeningToolsResponse) error {
	return gm.MockAnswerScreeningToolQuestionsFn(ctx, screeningToolResponses)
}

// GetScreeningToolQuestionByQuestionID mocks the implementation of getting screening tool questions by question ID
func (gm *GormMock) GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error) {
	return gm.MockGetScreeningToolQuestionByQuestionIDFn(ctx, questionID)
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

// InvalidateScreeningToolResponse mocks the implementation of invalidating screening tool responses
func (gm *GormMock) InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error {
	return gm.MockInvalidateScreeningToolResponseFn(ctx, clientID, questionID)
}

// UpdateServiceRequests mocks the implementation of updating service requests from KenyaEMR to MyCareHub
func (gm *GormMock) UpdateServiceRequests(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error) {
	return gm.MockUpdateServiceRequestsFn(ctx, payload)
}

// GetClientProfileByCCCNumber mocks the implementation of retrieving a client profile by CCC number
func (gm *GormMock) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*gorm.Client, error) {
	return gm.MockGetClientProfileByCCCNumberFn(ctx, CCCNumber)
}

// CheckIfClientHasUnresolvedServiceRequests mocks the implementation of checking if a client has a pending service request
func (gm *GormMock) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return gm.MockCheckIfClientHasUnresolvedServiceRequestsFn(ctx, clientID, serviceRequestType)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation of updating a user pin change required state
func (gm *GormMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour, status)
}

// GetAllRoles mocks the implementation of getting all roles
func (gm *GormMock) GetAllRoles(ctx context.Context) ([]*gorm.AuthorityRole, error) {
	return gm.MockGetAllRolesFn(ctx)
}

// SearchClientProfile mocks the implementation of searching for client profiles.
func (gm *GormMock) SearchClientProfile(ctx context.Context, CCCNumber string) ([]*gorm.Client, error) {
	return gm.MockSearchClientProfileFn(ctx, CCCNumber)
}

// SearchStaffProfile mocks the implementation of getting staff profile using their staff number.
func (gm *GormMock) SearchStaffProfile(ctx context.Context, searchParameter string) ([]*gorm.StaffProfile, error) {
	return gm.MockSearchStaffProfileFn(ctx, searchParameter)
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

// GetServiceRequestByID mocks the implementation of getting a service request by ID
func (gm *GormMock) GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestByIDFn(ctx, serviceRequestID)
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
func (gm *GormMock) GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetStaffPendingServiceRequestsCountFn(ctx, facilityID)
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

// GetActiveScreeningToolResponses mocks the implementation of getting active screening tool responses
func (gm *GormMock) GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*gorm.ScreeningToolsResponse, error) {
	return gm.MockGetActiveScreeningToolResponsesFn(ctx, clientID)
}

// GetAnsweredScreeningToolQuestions mocks the implementation of getting answered screening tool questions
func (gm *GormMock) GetAnsweredScreeningToolQuestions(ctx context.Context, facilityID string, toolType string) ([]*gorm.ScreeningToolsResponse, error) {
	return gm.MockGetAnsweredScreeningToolQuestionsFn(ctx, facilityID, toolType)
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
func (gm *GormMock) GetUserSurveyForms(ctx context.Context, userID string) ([]*gorm.UserSurvey, error) {
	return gm.MockGetUserSurveyFormsFn(ctx, userID)
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

// GetClientScreeningToolResponsesByToolType mocks the implementation of getting client screening tool responses by tool type
func (gm *GormMock) GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID, toolType string, active bool) ([]*gorm.ScreeningToolsResponse, error) {
	return gm.MockGetClientScreeningToolResponsesByToolTypeFn(ctx, clientID, toolType, active)
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

//UpdateNotification updates a notification with the new data
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
func (gm *GormMock) RegisterStaff(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) error {
	return gm.MockRegisterStaffFn(ctx, contact, identifier, staffProfile)
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
func (gm *GormMock) RegisterClient(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) error {
	return gm.MockRegisterClientFn(ctx, contact, identifier, client)
}
