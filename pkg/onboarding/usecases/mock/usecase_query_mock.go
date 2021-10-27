package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// QueryMock is a mock of the query methods
type QueryMock struct {
	RetrieveFacilityFn             func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCodeFn    func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetFacilitiesFn                func(ctx context.Context) ([]*domain.Facility, error)
	GetUserPINByUserIDFn           func(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetUserProfileByUserIDFn       func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.User, error)
	GetClientProfileByClientIDFn   func(ctx context.Context, clientID string) (*domain.ClientProfile, error)
	GetStaffProfileByStaffIDFn     func(ctx context.Context, staffProfileID string) (*domain.StaffUserProfile, error)
	GetStaffProfileByStaffNumberFn func(ctx context.Context, staffNumber string) (*domain.StaffUserProfile, error)
	GetStaffProfileFn              func(ctx context.Context, staffNumber string) (*domain.StaffProfile, error)
}

// NewQueryMock initializes a new instance of `GormMock` then mocking the case of success.
func NewQueryMock() *QueryMock {
	return &QueryMock{

		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := "test-county"
			description := "test description"
			return &domain.Facility{
				ID:          &facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},

		RetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := "test-county"
			description := "test description"
			return &domain.Facility{
				ID:          &facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},

		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := "test-county"
			description := "test description"
			return []*domain.Facility{
				{
					ID:          &facilityID,
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			}, nil
		},

		GetUserPINByUserIDFn: func(ctx context.Context, userID string) (*domain.UserPIN, error) {
			return &domain.UserPIN{
				UserID:    userID,
				HashedPIN: "mbzcbvhbxchjbvhdbvhhjdfskgbfhas832y38hjsdnfkjbh73y73y72",
				ValidFrom: time.Now(),
				ValidTo:   time.Now(),
				Flavour:   "CONSUMER",
				IsValid:   true,
				Salt:      "test-salt",
			}, nil
		},
		GetStaffProfileByStaffIDFn: func(ctx context.Context, staffProfileID string) (*domain.StaffUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()
			return &domain.StaffUserProfile{
				User: &domain.User{
					ID:                  &ID,
					Username:            "test",
					DisplayName:         "test",
					FirstName:           "test",
					MiddleName:          "test",
					LastName:            "test",
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					FailedLoginCount:    "0",
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
				},
				Staff: &domain.StaffProfile{
					ID:                &ID,
					UserID:            &ID,
					StaffNumber:       "s123",
					DefaultFacilityID: &ID,
					Addresses: []*domain.Addresses{
						{
							ID:         ID,
							Type:       enums.AddressesTypePhysical,
							Text:       "test",
							Country:    enums.CountryTypeKenya,
							PostalCode: "test code",
							County:     enums.CountyTypeBaringo,
							Active:     true,
						},
					},
				},
			}, nil
		},
		GetStaffProfileByStaffNumberFn: func(ctx context.Context, staffNumber string) (*domain.StaffUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()
			return &domain.StaffUserProfile{
				User: &domain.User{
					ID:                  &ID,
					Username:            "test",
					DisplayName:         "test",
					FirstName:           "test",
					MiddleName:          "test",
					LastName:            "test",
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					FailedLoginCount:    "0",
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
				},
				Staff: &domain.StaffProfile{
					ID:                &ID,
					UserID:            &ID,
					StaffNumber:       "s123",
					DefaultFacilityID: &ID,
					Addresses: []*domain.Addresses{
						{
							ID:         ID,
							Type:       enums.AddressesTypePhysical,
							Text:       "test",
							Country:    enums.CountryTypeKenya,
							PostalCode: "test code",
							County:     enums.CountyTypeBaringo,
							Active:     true,
						},
					},
				},
			}, nil
		},
	}
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *QueryMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return f.RetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacilityByMFLCode method.
func (f *QueryMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return f.RetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (f *QueryMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.GetFacilitiesFn(ctx)
}

// GetUserPINByUserID ...
func (f *QueryMock) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return f.GetUserPINByUserIDFn(ctx, userID)
}

// GetUserProfileByUserID gets user profile by user ID
func (f *QueryMock) GetUserProfileByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.User, error) {
	return f.GetUserProfileByUserIDFn(ctx, userID, flavour)
}

// GetClientProfileByClientID defines a mock for fetching a client profile using the client's ID
func (f *QueryMock) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	return f.GetClientProfileByClientIDFn(ctx, clientID)
}

// GetStaffProfileByStaffID mocks the  GetStaffProfileByStaffID method.
func (f *QueryMock) GetStaffProfileByStaffID(ctx context.Context, staffProfileID string) (*domain.StaffUserProfile, error) {
	return f.GetStaffProfileByStaffIDFn(ctx, staffProfileID)
}

// GetStaffProfileByStaffNumber mocks the  GetStaffProfileByStaffNumber method.
func (f *QueryMock) GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*domain.StaffUserProfile, error) {
	return f.GetStaffProfileByStaffNumberFn(ctx, staffNumber)
}
