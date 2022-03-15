package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/scalarutils"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	//Get
	MockGetOrCreateFacilityFn                     func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	MockGetFacilitiesFn                           func(ctx context.Context) ([]*domain.Facility, error)
	MockRetrieveFacilityFn                        func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	ListFacilitiesFn                              func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockDeleteFacilityFn                          func(ctx context.Context, id int) (bool, error)
	MockRetrieveFacilityByMFLCodeFn               func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
	MockGetUserProfileByPhoneNumberFn             func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error)
	MockGetUserPINByUserIDFn                      func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error)
	MockInactivateFacilityFn                      func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn                      func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn                  func(ctx context.Context, userID string) (*domain.User, error)
	MockSaveTemporaryUserPinFn                    func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	MockGetCurrentTermsFn                         func(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error)
	MockAcceptTermsFn                             func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                                 func(ctx context.Context, pin *domain.UserPIN) (bool, error)
	MockUpdateUserFailedLoginCountFn              func(ctx context.Context, userID string, failedLoginAttempts int) error
	MockUpdateUserLastFailedLoginTimeFn           func(ctx context.Context, userID string) error
	MockUpdateUserNextAllowedLoginTimeFn          func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	MockSetNickNameFn                             func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockUpdateUserLastSuccessfulLoginTimeFn       func(ctx context.Context, userID string) error
	MockGetSecurityQuestionsFn                    func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	MockSaveOTPFn                                 func(ctx context.Context, otpInput *domain.OTP) error
	MockGetSecurityQuestionByIDFn                 func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn            func(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error
	MockGetSecurityQuestionResponseByIDFn         func(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn                func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                               func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetClientProfileByUserIDFn                func(ctx context.Context, userID string) (*domain.ClientProfile, error)
	MockGetStaffProfileByUserIDFn                 func(ctx context.Context, userID string) (*domain.StaffProfile, error)
	MockCheckUserHasPinFn                         func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGenerateRetryOTPFn                        func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	MockUpdateUserPinChangeRequiredStatusFn       func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                                  func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error)
	MockGetUserSecurityQuestionsResponsesFn       func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error)
	MockInvalidatePINFn                           func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetContactByUserIDFn                      func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockListContentCategoriesFn                   func(ctx context.Context) ([]*domain.ContentItemCategory, error)
	MockShareContentFn                            func(ctx context.Context, input dto.ShareContentInput) (bool, error)
	MockBookmarkContentFn                         func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnBookmarkContentFn                       func(ctx context.Context, userID string, contentID int) (bool, error)
	MockGetUserBookmarkedContentFn                func(ctx context.Context, userID string) ([]*domain.ContentItem, error)
	MockLikeContentFn                             func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCheckWhetherUserHasLikedContentFn         func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnlikeContentFn                           func(ctx context.Context, userID string, contentID int) (bool, error)
	MockFetchFacilitiesFn                         func(ctx context.Context) ([]*domain.Facility, error)
	MockViewContentFn                             func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCreateHealthDiaryEntryFn                  func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error
	MockCreateServiceRequestFn                    func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error
	MockCanRecordHeathDiaryFn                     func(ctx context.Context, userID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn               func(ctx context.Context) (*domain.ClientHealthDiaryQuote, error)
	MockCheckIfUserBookmarkedContentFn            func(ctx context.Context, userID string, contentID int) (bool, error)
	MockGetClientHealthDiaryEntriesFn             func(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetFAQContentFn                           func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error)
	MockCreateClientCaregiverFn                   func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockGetClientCaregiverFn                      func(ctx context.Context, caregiverID string) (*domain.Caregiver, error)
	MockUpdateClientCaregiverFn                   func(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	MockInProgressByFn                            func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetClientProfileByClientIDFn              func(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	MockGetServiceRequestsFn                      func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error)
	MockGetPendingServiceRequestsCountFn          func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockResolveServiceRequestFn                   func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	MockCreateCommunityFn                         func(ctx context.Context, community *dto.CommunityInput) (*domain.Community, error)
	MockCheckUserRoleFn                           func(ctx context.Context, userID string, role string) (bool, error)
	MockCheckUserPermissionFn                     func(ctx context.Context, userID string, permission string) (bool, error)
	MockAssignRolesFn                             func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetUserRolesFn                            func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error)
	MockGetUserPermissionsFn                      func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error)
	MockCheckIfUsernameExistsFn                   func(ctx context.Context, username string) (bool, error)
	MockRevokeRolesFn                             func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	MockGetCommunityByIDFn                        func(ctx context.Context, communityID string) (*domain.Community, error)
	MockCheckIdentifierExists                     func(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	MockCheckFacilityExistsByMFLCode              func(ctx context.Context, MFLCode int) (bool, error)
	MockCreateNextOfKin                           func(ctx context.Context, person *dto.NextOfKinPayload) error
	MockCreateContact                             func(ctx context.Context, contact *domain.Contact) error
	MockGetClientsInAFacilityFn                   func(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error)
	MockGetRecentHealthDiaryEntriesFn             func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetClientsByParams                        func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error)
	MockGetClientCCCIdentifier                    func(ctx context.Context, clientID string) (*domain.Identifier, error)
	MockGetServiceRequestsForKenyaEMRFn           func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	MockCreateAppointment                         func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error
	MockUpdateAppointment                         func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error
	MockGetScreeningToolsQuestionsFn              func(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error)
	MockAnswerScreeningToolQuestionsFn            func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error
	MockGetScreeningToolQuestionByQuestionIDFn    func(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error)
	MockInvalidateScreeningToolResponseFn         func(ctx context.Context, clientID string, questionID string) error
	MockUpdateServiceRequestsFn                   func(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error)
	MockListAppointments                          func(ctx context.Context, params *domain.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error)
	MockGetClientProfileByCCCNumberFn             func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error)
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
	}

	client := &domain.ClientProfile{
		ID: &ID,
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

	contentItemCategoryID := 1
	contentItemCategory := &domain.ContentItemCategory{
		ID:      contentItemCategoryID,
		Name:    name,
		IconURL: "test",
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

	return &PostgresMock{
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
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
		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
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
		MockUpdateServiceRequestsFn: func(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error) {
			return true, nil
		},
		MockSaveTemporaryUserPinFn: func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
			return true, nil
		},
		MockAcceptTermsFn: func(ctx context.Context, userID *string, termsID *int) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
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
		MockUpdateUserLastSuccessfulLoginTimeFn: func(ctx context.Context, userID string) error {
			return nil
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
		MockGetSecurityQuestionResponseByIDFn: func(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error) {
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
			return client, nil
		},
		MockGetStaffProfileByUserIDFn: func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
			return staff, nil
		},
		MockCheckUserHasPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGenerateRetryOTPFn: func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
			return "test-OTP", nil
		},
		MockUpdateUserPinChangeRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
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
		MockListContentCategoriesFn: func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
			return []*domain.ContentItemCategory{contentItemCategory}, nil
		},
		MockShareContentFn: func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
			return true, nil
		},
		MockBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockGetUserBookmarkedContentFn: func(ctx context.Context, userID string) ([]*domain.ContentItem, error) {
			return []*domain.ContentItem{
				{
					ID: int(uuid.New()[8]),
				},
			}, nil
		},
		MockLikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnlikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockFetchFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return []*domain.Facility{facilityInput}, nil
		},
		MockViewContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockCheckWhetherUserHasLikedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
			return nil
		},
		MockCreateServiceRequestFn: func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
			return nil
		},
		MockCanRecordHeathDiaryFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryQuoteFn: func(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
			return &domain.ClientHealthDiaryQuote{
				Quote:  "test",
				Author: "test",
			}, nil
		},
		MockCheckIfUserBookmarkedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryEntriesFn: func(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockGetFAQContentFn: func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
			return []*domain.FAQ{
				{
					ID:          &ID,
					Active:      true,
					Title:       gofakeit.Name(),
					Description: gofakeit.Name(),
					Body:        gofakeit.Name(),
				},
			}, nil
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
		MockGetPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
			return &domain.ServiceRequestsCount{
				Total: 5,
				RequestsTypeCount: []*domain.RequestTypeCount{
					{
						RequestType: enums.ServiceRequestTypeRedFlag,
						Total:       1,
					},
				},
			}, nil
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
				ID:   &ID,
				User: userProfile,
			}
			return client, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
			return serviceRequests, nil
		},
		MockResolveServiceRequestFn: func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
			return true, nil
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
					RoleID: uuid.New().String(),
					Name:   enums.UserRoleTypeClientManagement,
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
				client,
			}, nil
		},
		MockGetRecentHealthDiaryEntriesFn: func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
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
				ValidFrom:           time.Time{},
				ValidTo:             time.Time{},
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
				FacilityID:    &facilityID,
				ClientName:    &staffID,
				ClientContact: &contact,
			}
			return []*domain.ServiceRequest{serviceReq}, nil
		},
		MockCreateAppointment: func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error {
			return nil
		},
		MockUpdateAppointment: func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error {
			return nil
		},
		MockListAppointments: func(ctx context.Context, params *domain.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
			return []*domain.Appointment{{
				ID:       ID,
				Type:     "Dental",
				Status:   "COMPLETED",
				Reason:   "Bad tooth",
				Provider: "X",
				Date: scalarutils.Date{
					Year:  2023,
					Month: 1,
					Day:   1,
				},
				Start: time.Now(),
				End:   time.Now().Add(1 * time.Hour),
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
		MockGetScreeningToolQuestionByQuestionIDFn: func(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error) {
			return &domain.ScreeningToolQuestion{
				ID:       ID,
				Question: gofakeit.Sentence(1),
				ToolType: enums.ScreeningToolTypeTB,
				ResponseChoices: map[string]interface{}{
					"0": "yes",
					"1": "no",
				},
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "data"},
				Active:           true,
			}, nil
		},
		MockInvalidateScreeningToolResponseFn: func(ctx context.Context, clientID string, questionID string) error {
			return nil
		},
		MockGetClientProfileByCCCNumberFn: func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
			return client, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// CheckWhetherUserHasLikedContent mocks the implementation of `gorm's` CheckWhetherUserHasLikedContent method.
func (gm *PostgresMock) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockCheckWhetherUserHasLikedContentFn(ctx, userID, contentID)
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

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (gm *PostgresMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockGetFacilitiesFn(ctx)
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

// UpdateUserFailedLoginCount mocks the implementation of updating a user failed login count
func (gm *PostgresMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	return gm.MockUpdateUserFailedLoginCountFn(ctx, userID, failedLoginAttempts)
}

// UpdateUserLastFailedLoginTime mocks the implementation of updating a user's last failed login time
func (gm *PostgresMock) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	return gm.MockUpdateUserLastFailedLoginTimeFn(ctx, userID)
}

// UpdateUserNextAllowedLoginTime mocks the implementation of updating a user's next allowed login time
func (gm *PostgresMock) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	return gm.MockUpdateUserNextAllowedLoginTimeFn(ctx, userID, nextAllowedLoginTime)
}

// UpdateUserLastSuccessfulLoginTime mocks the implementation of updating a user's last successful login time
func (gm *PostgresMock) UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error {
	return gm.MockUpdateUserLastSuccessfulLoginTimeFn(ctx, userID)
}

//GetSecurityQuestions mocks the implementation of getting all the security questions.
func (gm *PostgresMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return gm.MockGetSecurityQuestionsFn(ctx, flavour)
}

// SaveOTP mocks the implementation for saving an OTP
func (gm *PostgresMock) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	return gm.MockSaveOTPFn(ctx, otpInput)
}

// SetNickName is used to mock the implementation ofsetting or changing the user's nickname
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

// GetSecurityQuestionResponseByID mocks the get security question implementation
func (gm *PostgresMock) GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error) {
	return gm.MockGetSecurityQuestionResponseByIDFn(ctx, questionID)
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

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *PostgresMock) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID, flavour)
}

// GenerateRetryOTP mock the implementtation of generating a retry OTP
func (gm *PostgresMock) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return gm.MockGenerateRetryOTPFn(ctx, payload)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation for updating a user's pin change required state
func (gm *PostgresMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour)
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

//ListContentCategories mocks the implementation listing content categories
func (gm *PostgresMock) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return gm.MockListContentCategoriesFn(ctx)
}

// ShareContent mock the implementation share content
func (gm *PostgresMock) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return gm.MockShareContentFn(ctx, input)
}

// BookmarkContent bookmarks a content
func (gm *PostgresMock) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockBookmarkContentFn(ctx, userID, contentID)
}

// UnBookmarkContent unbookmarks a content
func (gm *PostgresMock) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockUnBookmarkContentFn(ctx, userID, contentID)
}

// GetUserBookmarkedContent mocks the implementation of retrieving a user bookmarked content
func (gm *PostgresMock) GetUserBookmarkedContent(ctx context.Context, userID string) ([]*domain.ContentItem, error) {
	return gm.MockGetUserBookmarkedContentFn(ctx, userID)
}

//LikeContent mocks the implementation liking a feed content
func (gm *PostgresMock) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockLikeContentFn(ctx, userID, contentID)
}

//UnlikeContent mocks the implementation liking a feed content
func (gm *PostgresMock) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockUnlikeContentFn(ctx, userID, contentID)
}

//FetchFacilities mocks the implementation of fetching facility
func (gm *PostgresMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockFetchFacilitiesFn(ctx)
}

// ViewContent gets a content and updates the view count
func (gm *PostgresMock) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockViewContentFn(ctx, userID, contentID)
}

// CreateHealthDiaryEntry mocks the method for creating a health diary entry
func (gm *PostgresMock) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
	return gm.MockCreateHealthDiaryEntryFn(ctx, healthDiaryInput)
}

// CreateServiceRequest mocks creating a service request method
func (gm *PostgresMock) CreateServiceRequest(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
	return gm.MockCreateServiceRequestFn(ctx, serviceRequestInput)
}

// CanRecordHeathDiary mocks the implementation of checking if a user can record a health diary
func (gm *PostgresMock) CanRecordHeathDiary(ctx context.Context, userID string) (bool, error) {
	return gm.MockCanRecordHeathDiaryFn(ctx, userID)
}

// GetClientHealthDiaryQuote mocks the implementation of fetching client health diary quote
func (gm *PostgresMock) GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
	return gm.MockGetClientHealthDiaryQuoteFn(ctx)
}

// CheckIfUserBookmarkedContent mocks the implementation of checking if a user has bookmarked a content
func (gm *PostgresMock) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return gm.MockCheckIfUserBookmarkedContentFn(ctx, userID, contentID)
}

// GetClientHealthDiaryEntries mocks the implementation of getting all health diary entries that belong to a specific user
func (gm *PostgresMock) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetClientHealthDiaryEntriesFn(ctx, clientID)
}

// GetFAQContent mocks the implementation of getting FAQ content
func (gm *PostgresMock) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
	return gm.MockGetFAQContentFn(ctx, flavour, limit)
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
func (gm *PostgresMock) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return gm.MockGetPendingServiceRequestsCountFn(ctx, facilityID)
}

// GetServiceRequests mocks the implementation of getting all service requests for a client
func (gm *PostgresMock) GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
	return gm.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID)
}

// ResolveServiceRequest mocks the implementation of resolving a service request
func (gm *PostgresMock) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
	return gm.MockResolveServiceRequestFn(ctx, staffID, serviceRequestID)
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

// CreateNextOfKin mocks creating a next of kin
func (gm *PostgresMock) CreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload) error {
	return gm.MockCreateNextOfKin(ctx, person)
}

// CreateContact mocks creating a contact
func (gm *PostgresMock) CreateContact(ctx context.Context, contact *domain.Contact) error {
	return gm.MockCreateContact(ctx, contact)
}

// GetClientsInAFacility mocks getting all the clients in a facility
func (gm *PostgresMock) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
	return gm.MockGetClientsInAFacilityFn(ctx, facilityID)
}

// GetRecentHealthDiaryEntries mocks getting the most recent health diary entry
func (gm *PostgresMock) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return gm.MockGetRecentHealthDiaryEntriesFn(ctx, lastSyncTime, clientID)
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
func (gm *PostgresMock) ListAppointments(ctx context.Context, params *domain.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
	return gm.MockListAppointments(ctx, params, filter, pagination)
}

// CreateAppointment creates a new appointment
func (gm *PostgresMock) CreateAppointment(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error {
	return gm.MockCreateAppointment(ctx, appointment, appointmentUUID, clientID, staffID)
}

// UpdateAppointment updates an appointment
func (gm *PostgresMock) UpdateAppointment(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID, staffID string) error {
	return gm.MockUpdateAppointment(ctx, appointment, appointmentUUID, clientID, staffID)
}

// UpdateServiceRequests mocks the implementation of updating service requests
func (gm *PostgresMock) UpdateServiceRequests(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error) {
	return gm.MockUpdateServiceRequestsFn(ctx, payload)
}

// GetClientProfileByCCCNumber mocks the implementation of getting a client profile using the CCC number
func (gm *PostgresMock) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
	return gm.MockGetClientProfileByCCCNumberFn(ctx, CCCNumber)
}
