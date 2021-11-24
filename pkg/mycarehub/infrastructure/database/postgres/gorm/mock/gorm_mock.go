package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockGetOrCreateFacilityFn               func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn                  func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn         func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error)
	MockGetFacilitiesFn                     func(ctx context.Context) ([]gorm.Facility, error)
	MockDeleteFacilityFn                    func(ctx context.Context, mflCode int) (bool, error)
	MockListFacilitiesFn                    func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	MockGetUserProfileByPhoneNumberFn       func(ctx context.Context, phoneNumber string) (*gorm.User, error)
	MockGetUserPINByUserIDFn                func(ctx context.Context, userID string) (*gorm.PINData, error)
	MockInactivateFacilityFn                func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn                func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn            func(ctx context.Context, UserID string) (*gorm.User, error)
	MockSaveTemporaryUserPinFn              func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetCurrentTermsFn                   func(ctx context.Context) (*gorm.TermsOfService, error)
	MockAcceptTermsFn                       func(ctx context.Context, userID *string, termsID *int) (bool, error)
	MockSavePinFn                           func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockUpdateUserFailedLoginCountFn        func(ctx context.Context, userID string, failedLoginAttempts int) error
	MockUpdateUserLastFailedLoginTimeFn     func(ctx context.Context, userID string) error
	MockUpdateUserNextAllowedLoginTimeFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	MockUpdateUserLastSuccessfulLoginTimeFn func(ctx context.Context, userID string) error
	MockSetNickNameFn                       func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockGetSecurityQuestionsFn              func(ctx context.Context, flavour feedlib.Flavour) ([]*gorm.SecurityQuestion, error)
	MockSaveOTPFn                           func(ctx context.Context, otpInput *gorm.UserOTP) error
	MockGetSecurityQuestionByIDFn           func(ctx context.Context, securityQuestionID *string) (*gorm.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn      func(ctx context.Context, securityQuestionResponse *gorm.SecurityQuestionResponse) error
	MockGetSecurityQuestionResponseByIDFn   func(ctx context.Context, questionID string) (*gorm.SecurityQuestionResponse, error)
	MockCheckIfPhoneNumberExistsFn          func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	MockVerifyOTPFn                         func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	MockGetSntProfileByUserIDFn             func(ctx context.Context, userID string) (*gorm.Client, error)
	MockGetClientProfileByUserIDFn          func(ctx context.Context, userID string) (*gorm.Client, error)
	MockCheckUserHasPinFn                   func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockUpdateUserPinChangeRequiredStatusFn func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockGetOTPFn                            func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*gorm.UserOTP, error)
	MockGetUserSecurityQuestionsResponsesFn func(ctx context.Context, userID string) ([]*gorm.SecurityQuestionResponse, error)

	MockInvalidatePINFn                           func(ctx context.Context, userID string) (bool, error)
	MockGetContactByUserIDFn                      func(ctx context.Context, userID *string, contactType string) (*gorm.Contact, error)
	MockUpdateIsCorrectSecurityQuestionResponseFn func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	MockListContentCategoriesFn                   func(ctx context.Context) ([]*gorm.ContentItemCategory, error)
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

	facility := &gorm.Facility{
		FacilityID:  &UUID,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
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
		ID: &UUID,
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

	contentItemCategory := &gorm.ContentItemCategory{
		ID:   ID,
		Name: name,
		Icon: gorm.WagtailImage{
			ID:   ID,
			File: "url",
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

		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},

		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return pinData, nil
		},

		MockInactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockGetCurrentTermsFn: func(ctx context.Context) (*gorm.TermsOfService, error) {
			termsID := gofakeit.Number(1, 1000)
			validFrom := time.Now()
			testText := "test"

			validTo := time.Now().AddDate(0, 0, 80)
			terms := &gorm.TermsOfService{
				TermsID:   &termsID,
				Text:      &testText,
				ValidFrom: &validFrom,
				ValidTo:   &validTo,
			}
			return terms, nil
		},
		MockGetUserProfileByUserIDFn: func(ctx context.Context, UserID string) (*gorm.User, error) {
			return &gorm.User{
				UserID: &UserID,
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
		MockSaveSecurityQuestionResponseFn: func(ctx context.Context, securityQuestionResponse *gorm.SecurityQuestionResponse) error {
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
		MockCheckUserHasPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockUpdateUserPinChangeRequiredStatusFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
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
		MockInvalidatePINFn: func(ctx context.Context, userID string) (bool, error) {
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
		MockListContentCategoriesFn: func(ctx context.Context) ([]*gorm.ContentItemCategory, error) {
			return []*gorm.ContentItemCategory{contentItemCategory}, nil
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
func (gm *GormMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*gorm.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// GetUserPINByUserID mocks the implementation of retrieving a user pin by user ID
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
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
func (gm *GormMock) GetCurrentTerms(ctx context.Context) (*gorm.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}

// GetUserProfileByUserID mocks the implementation of retrieving a user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, UserID string) (*gorm.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, UserID)
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

// UpdateUserLastSuccessfulLoginTime mocks the implementation of updating a user's last successful login time
func (gm *GormMock) UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error {
	return gm.MockUpdateUserLastSuccessfulLoginTimeFn(ctx, userID)
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
func (gm *GormMock) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse *gorm.SecurityQuestionResponse) error {
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

// CheckUserHasPin mocks the method for checking if a user has a pin
func (gm *GormMock) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockCheckUserHasPinFn(ctx, userID, flavour)
}

// UpdateUserPinChangeRequiredStatus mocks the implementation for updating a user's pin change required state
func (gm *GormMock) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return gm.MockUpdateUserPinChangeRequiredStatusFn(ctx, userID, flavour)
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

//ListContentCategories mocks the implementation listing content categories
func (gm *GormMock) ListContentCategories(ctx context.Context) ([]*gorm.ContentItemCategory, error) {
	return gm.MockListContentCategoriesFn(ctx)
}
