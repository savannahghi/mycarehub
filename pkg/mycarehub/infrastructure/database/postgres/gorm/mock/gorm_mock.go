package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockGetOrCreateFacilityFn                       func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn                          func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn                 func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error)
	MockGetFacilitiesFn                             func(ctx context.Context) ([]gorm.Facility, error)
	MockDeleteFacilityFn                            func(ctx context.Context, mflCode int) (bool, error)
	MockListFacilitiesFn                            func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	MockGetUserProfileByPhoneNumberFn               func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.User, error)
	MockGetUserPINByUserIDFn                        func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.PINData, error)
	MockInactivateFacilityFn                        func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn                        func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn                    func(ctx context.Context, userID *string) (*gorm.User, error)
	MockSaveTemporaryUserPinFn                      func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetCurrentTermsFn                           func(ctx context.Context, flavour feedlib.Flavour) (*gorm.TermsOfService, error)
	MockAcceptTermsFn                               func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                   func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockUpdateUserFailedLoginCountFn                func(ctx context.Context, userID string, failedLoginAttempts int) error
	MockUpdateUserLastFailedLoginTimeFn             func(ctx context.Context, userID string) error
	MockUpdateUserNextAllowedLoginTimeFn            func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	MockUpdateUserProfileAfterLoginSuccessFn        func(ctx context.Context, userID string) error
	MockSetNickNameFn                               func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockGetSecurityQuestionsFn                      func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error)
	MockSaveOTPFn                                   func(ctx context.Context, otpInput *gorm.UserOTP) error
	MockGetSecurityQuestionByIDFn                   func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn              func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error
	MockGetSecurityQuestionResponseByIDFn           func(ctx context.Context, questionID string) (*gorm.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                  func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                                 func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetClientProfileByUserIDFn                  func(ctx context.Context, userID string) (*gorm.Client, error)
	MockGetStaffProfileByUserIDFn                   func(ctx context.Context, userID string) (*gorm.StaffProfile, error)
	MockCheckUserHasPinFn                           func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockCompleteOnboardingTourFn                    func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                    func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error)
	MockGetUserSecurityQuestionsResponsesFn         func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error)
	MockInvalidatePINFn                             func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetContactByUserIDFn                        func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn   func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockListContentCategoriesFn                     func(ctx context.Context) ([]*domain.ContentItemCategory, error)
	MockShareContentFn                              func(ctx context.Context, input dto.ShareContentInput) (bool, error)
	MockBookmarkContentFn                           func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnBookmarkContentFn                         func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCheckWhetherUserHasLikedContentFn           func(ctx context.Context, userID string, contentID int) (bool, error)
	MockGetUserBookmarkedContentFn                  func(ctx context.Context, userID string) ([]*gorm.ContentItem, error)
	MockLikeContentFn                               func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnlikeContentFn                             func(ctx context.Context, userID string, contentID int) (bool, error)
	MockViewContentFn                               func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCreateHealthDiaryEntryFn                    func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) error
	MockCreateServiceRequestFn                      func(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error
	MockCanRecordHeathDiaryFn                       func(ctx context.Context, clientID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn                 func(ctx context.Context) (*gorm.ClientHealthDiaryQuote, error)
	MockCheckIfUserBookmarkedContentFn              func(ctx context.Context, userID string, contentID int) (bool, error)
	MockGetClientHealthDiaryEntriesFn               func(ctx context.Context, clientID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetFAQContentFn                             func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*gorm.FAQ, error)
	MockCreateClientCaregiverFn                     func(ctx context.Context, clientID string, clientCaregiver *gorm.Caregiver) error
	MockGetClientCaregiverFn                        func(ctx context.Context, caregiverID string) (*gorm.Caregiver, error)
	MockUpdateClientCaregiverFn                     func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockInProgressByFn                              func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn                func(ctx context.Context, clientID string) (*gorm.Client, error)
	MockGetServiceRequestsFn                        func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*gorm.ClientServiceRequest, error)
	MockGetPendingServiceRequestsCountFn            func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockResolveServiceRequestFn                     func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	MockCheckUserRoleFn                             func(ctx context.Context, userID string, role string) (bool, error)
	MockCheckUserPermissionFn                       func(ctx context.Context, userID string, permission string) (bool, error)
	MockAssignRolesFn                               func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockCreateCommunityFn                           func(ctx context.Context, community *gorm.Community) (*gorm.Community, error)
	MockGetUserRolesFn                              func(ctx context.Context, userID string) ([]*gorm.AuthorityRole, error)
	MockGetUserPermissionsFn                        func(ctx context.Context, userID string) ([]*gorm.AuthorityPermission, error)
	MockCheckIfUsernameExistsFn                     func(ctx context.Context, username string) (bool, error)
	MockRevokeRolesFn                               func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetCommunityByIDFn                          func(ctx context.Context, communityID string) (*gorm.Community, error)
	MockCheckIdentifierExists                       func(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	MockCheckFacilityExistsByMFLCode                func(ctx context.Context, MFLCode int) (bool, error)
	MockCreateRelatedPerson                         func(ctx context.Context, person *gorm.RelatedPerson) error
	MockCreateContact                               func(ctx context.Context, contact *gorm.Contact) error
	MockGetClientsInAFacilityFn                     func(ctx context.Context, facilityID string) ([]*gorm.Client, error)
	MockGetRecentHealthDiaryEntriesFn               func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*gorm.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                          func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*gorm.Client, error)
	MockGetClientCCCIdentifier                      func(ctx context.Context, clientID string) (*gorm.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn             func(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*gorm.ClientServiceRequest, error)
	MockCreateAppointment                           func(ctx context.Context, appointment *gorm.Appointment) error
	MockListAppointments                            func(ctx context.Context, params *gorm.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error)
	MockUpdateAppointment                           func(ctx context.Context, payload *gorm.Appointment) error
	MockGetScreeningToolsQuestionsFn                func(ctx context.Context, toolType string) ([]gorm.ScreeningToolQuestion, error)
	MockAnswerScreeningToolQuestionsFn              func(ctx context.Context, screeningToolResponses []*gorm.ScreeningToolsResponse) error
	MockGetScreeningToolQuestionByQuestionIDFn      func(ctx context.Context, questionID string) (*gorm.ScreeningToolQuestion, error)
	MockInvalidateScreeningToolResponseFn           func(ctx context.Context, clientID string, questionID string) error
	MockUpdateServiceRequestsFn                     func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error)
	MockGetClientProfileByCCCNumberFn               func(ctx context.Context, CCCNumber string) (*gorm.Client, error)
	MockSearchClientProfilesByCCCNumberFn           func(ctx context.Context, CCCNumber string) ([]*gorm.Client, error)
	MockSearchStaffProfileByStaffNumberFn           func(ctx context.Context, staffNumber string) ([]*gorm.StaffProfile, error)
	MockUpdateUserPinChangeRequiredStatusFn         func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	MockCheckIfClientHasUnresolvedServiceRequestsFn func(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	MockGetAllRolesFn                               func(ctx context.Context) ([]*gorm.AuthorityRole, error)
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
	code := gofakeit.Number(0, 100)
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

	client := &gorm.Client{
		ID:         &UUID,
		Active:     true,
		ClientType: "",
		UserProfile: gorm.User{
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
		FHIRPatientID:           uuid.New().String(),
		HealthRecordID:          &UUID,
		TreatmentBuddy:          gofakeit.Name(),
		ClientCounselled:        true,
		OrganisationID:          uuid.New().String(),
		FacilityID:              uuid.New().String(),
		CHVUserID:               uuid.New().String(),
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

	contentItemCategory := &domain.ContentItemCategory{
		ID:      ID,
		Name:    name,
		IconURL: "https://test-icon-url/test.png",
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
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			return facility, nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {

			return facility, nil
		},
		MockGetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			return facilities, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, mflCode int) (bool, error) {
			return true, nil
		},

		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
			return facility, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
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
		MockGetUserProfileByUserIDFn: func(ctx context.Context, userID *string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},
		MockSaveTemporaryUserPinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
		MockAcceptTermsFn: func(ctx context.Context, userID *string, termsID *int) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
		MockUpdateServiceRequestsFn: func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error) {
			return true, nil
		},
		MockUpdateUserFailedLoginCountFn: func(ctx context.Context, userID string, failedLoginAttempts int) error {
			return nil
		},
		MockUpdateUserLastFailedLoginTimeFn: func(ctx context.Context, userID string) error {
			return nil
		},
		MockUpdateUserNextAllowedLoginTimeFn: func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
			return nil
		},
		MockUpdateUserProfileAfterLoginSuccessFn: func(ctx context.Context, userID string) error {
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
		MockGetSecurityQuestionResponseByIDFn: func(ctx context.Context, questionID string) (*gorm.SecurityQuestionResponse, error) {
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
			return client, nil
		},
		MockGetStaffProfileByUserIDFn: func(ctx context.Context, userID string) (*gorm.StaffProfile, error) {
			return staff, nil
		},
		MockSearchStaffProfileByStaffNumberFn: func(ctx context.Context, staffNumber string) ([]*gorm.StaffProfile, error) {
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
		MockListContentCategoriesFn: func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
			return []*domain.ContentItemCategory{contentItemCategory}, nil
		},
		MockShareContentFn: func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
			return true, nil
		}, MockGetUserBookmarkedContentFn: func(ctx context.Context, userID string) ([]*gorm.ContentItem, error) {
			return []*gorm.ContentItem{
				{
					PagePtrID: int(uuid.New()[9]),
				},
			}, nil
		},
		MockBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockLikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnlikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockViewContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
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
		MockGetClientHealthDiaryQuoteFn: func(ctx context.Context) (*gorm.ClientHealthDiaryQuote, error) {
			return &gorm.ClientHealthDiaryQuote{
				Quote:  "test",
				Author: "test",
			}, nil
		},
		MockCheckWhetherUserHasLikedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockCheckIfUserBookmarkedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryEntriesFn: func(ctx context.Context, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
			return []*gorm.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockGetFAQContentFn: func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*gorm.FAQ, error) {
			ID := uuid.New().String()
			return []*gorm.FAQ{
				{
					FAQID:       &ID,
					Active:      true,
					Title:       gofakeit.Name(),
					Description: gofakeit.Name(),
					Body:        gofakeit.Name(),
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
			return client, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*gorm.ClientServiceRequest, error) {
			return serviceRequests, nil
		},
		MockResolveServiceRequestFn: func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
			return true, nil
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
				client,
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
			return []*gorm.Client{client}, nil
		},
		MockGetClientCCCIdentifier: func(ctx context.Context, clientID string) (*gorm.Identifier, error) {
			return &gorm.Identifier{
				ID:                  uuid.New().String(),
				IdentifierType:      "CCC",
				IdentifierValue:     "123456",
				IdentifierUse:       "OFFICIAL",
				Description:         description,
				ValidFrom:           time.Time{},
				ValidTo:             time.Time{},
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
			}
			return []*gorm.ClientServiceRequest{serviceReq}, nil
		},
		MockCreateAppointment: func(ctx context.Context, appointment *gorm.Appointment) error {
			return nil
		},
		MockListAppointments: func(ctx context.Context, params *gorm.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error) {
			date := time.Now().Add(time.Duration(100))
			return []*gorm.Appointment{
				{
					ID:              gofakeit.UUID(),
					OrganisationID:  gofakeit.UUID(),
					Active:          true,
					AppointmentUUID: gofakeit.UUID(),
					AppointmentType: "Dental",
					Status:          enums.AppointmentStatusCompleted.String(),
					ClientID:        gofakeit.UUID(),
					FacilityID:      gofakeit.UUID(),
					Reason:          "Knocked up",
					Date:            date,
					StartTime:       gorm.CustomTime{Time: time.Now()},
					EndTime:         gorm.CustomTime{Time: time.Now().Add(30 * time.Minute)},
				},
			}, &domain.Pagination{Limit: 10, CurrentPage: 1}, nil
		},
		MockUpdateAppointment: func(ctx context.Context, payload *gorm.Appointment) error {
			return nil
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
				ToolType:         enums.ScreeningToolTypeTB.String(),
				ResponseChoices:  `{"1": "Yes", "2": "No"}`,
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
			return client, nil
		},
		MockSearchClientProfilesByCCCNumberFn: func(ctx context.Context, CCCNumber string) ([]*gorm.Client, error) {
			return []*gorm.Client{client}, nil
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
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// CheckWhetherUserHasLikedContent mocks the implementation of `gorm's` CheckWhetherUserHasLikedContent method.
func (gm *GormMock) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {

	return gm.MockCheckWhetherUserHasLikedContentFn(ctx, userID, contentID)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.MockGetFacilitiesFn(ctx)
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

// UpdateUserFailedLoginCount mocks the implementation of updating a user failed login count
func (gm *GormMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	return gm.MockUpdateUserFailedLoginCountFn(ctx, userID, failedLoginAttempts)
}

// UpdateUserLastFailedLoginTime mocks the implementation of updating a user's last failed login time
func (gm *GormMock) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	return gm.MockUpdateUserLastFailedLoginTimeFn(ctx, userID)
}

// UpdateUserNextAllowedLoginTime mocks the implementation of updating a user's next allowed login time
func (gm *GormMock) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	return gm.MockUpdateUserNextAllowedLoginTimeFn(ctx, userID, nextAllowedLoginTime)
}

// UpdateUserProfileAfterLoginSuccess mocks the implementation of updating a user's last successful login time
func (gm *GormMock) UpdateUserProfileAfterLoginSuccess(ctx context.Context, userID string) error {
	return gm.MockUpdateUserProfileAfterLoginSuccessFn(ctx, userID)
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

// GetSecurityQuestionResponseByID mocks the get security question implementation
func (gm *GormMock) GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*gorm.SecurityQuestionResponse, error) {
	return gm.MockGetSecurityQuestionResponseByIDFn(ctx, questionID)
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

//ListContentCategories mocks the implementation listing content categories
func (gm *GormMock) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return gm.MockListContentCategoriesFn(ctx)
}

// ShareContent mocks the implementation of sharing the content
func (gm *GormMock) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return gm.MockShareContentFn(ctx, input)
}

// BookmarkContent bookmarks a content
func (gm *GormMock) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockBookmarkContentFn(ctx, userID, contentID)
}

// UnBookmarkContent unbookmarks a content
func (gm *GormMock) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockUnBookmarkContentFn(ctx, userID, contentID)
}

// GetUserBookmarkedContent mocks the implementation of retrieving a user bookmarked content
func (gm *GormMock) GetUserBookmarkedContent(ctx context.Context, userID string) ([]*gorm.ContentItem, error) {
	return gm.MockGetUserBookmarkedContentFn(ctx, userID)
}

//LikeContent mocks the implementation liking a feed content
func (gm *GormMock) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockLikeContentFn(ctx, userID, contentID)
}

//UnlikeContent mocks the implementation liking a feed content
func (gm *GormMock) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockUnlikeContentFn(ctx, userID, contentID)
}

// ViewContent gets a content and updates the view count
func (gm *GormMock) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockViewContentFn(ctx, userID, contentID)
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
func (gm *GormMock) GetClientHealthDiaryQuote(ctx context.Context) (*gorm.ClientHealthDiaryQuote, error) {
	return gm.MockGetClientHealthDiaryQuoteFn(ctx)
}

// CheckIfUserBookmarkedContent mocks the implementation of checking if a user bookmarked a content
func (gm *GormMock) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockCheckIfUserBookmarkedContentFn(ctx, userID, contentID)
}

// GetClientHealthDiaryEntries mocks the implementation of getting all health diary entries that belong to a specific user
func (gm *GormMock) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*gorm.ClientHealthDiaryEntry, error) {
	return gm.MockGetClientHealthDiaryEntriesFn(ctx, clientID)
}

// GetFAQContent mocks the implementation of getting FAQ content
func (gm *GormMock) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*gorm.FAQ, error) {
	return gm.MockGetFAQContentFn(ctx, flavour, limit)
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

// GetPendingServiceRequestsCount mocks the implementation of getting the service requests count
func (gm *GormMock) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetPendingServiceRequestsCountFn(ctx, facilityID)
}

// GetServiceRequests mocks the implementation of getting service requests by type
func (gm *GormMock) GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*gorm.ClientServiceRequest, error) {
	return gm.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID)
}

// ResolveServiceRequest mocks the implementation of resolving a service request
func (gm *GormMock) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
	return gm.MockResolveServiceRequestFn(ctx, staffID, serviceRequestID)
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

// CreateRelatedPerson mocks creating a related person
func (gm *GormMock) CreateRelatedPerson(ctx context.Context, person *gorm.RelatedPerson) error {
	return gm.MockCreateRelatedPerson(ctx, person)
}

// CreateContact mocks creating a contact
func (gm *GormMock) CreateContact(ctx context.Context, contact *gorm.Contact) error {
	return gm.MockCreateContact(ctx, contact)
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

// GetServiceRequestsForKenyaEMR mocks the getting of service requests attched to a specific facility for use by KenyaEMR
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
func (gm *GormMock) ListAppointments(ctx context.Context, params *gorm.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*gorm.Appointment, *domain.Pagination, error) {
	return gm.MockListAppointments(ctx, params, filter, pagination)
}

// UpdateAppointment updates the details of an appointment requires the ID or appointment_uuid to be provided
func (gm *GormMock) UpdateAppointment(ctx context.Context, payload *gorm.Appointment) error {
	return gm.MockUpdateAppointment(ctx, payload)
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

// SearchClientProfilesByCCCNumber mocks the implementation of searching for client profiles.
// It returns clients profiles whose parts of the CCC number matches
func (gm *GormMock) SearchClientProfilesByCCCNumber(ctx context.Context, CCCNumber string) ([]*gorm.Client, error) {
	return gm.MockSearchClientProfilesByCCCNumberFn(ctx, CCCNumber)
}

// SearchStaffProfileByStaffNumber mocks the implementation of getting staff profile using their staff number.
func (gm *GormMock) SearchStaffProfileByStaffNumber(ctx context.Context, staffNumber string) ([]*gorm.StaffProfile, error) {
	return gm.MockSearchStaffProfileByStaffNumberFn(ctx, staffNumber)
}
