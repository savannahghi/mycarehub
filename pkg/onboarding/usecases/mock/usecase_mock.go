package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

// CreateMock is a mock of the create methods
type CreateMock struct {
	GetOrCreateFacilityFn func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	CollectMetricsFn      func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
	RegisterStaffUserFn   func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error)
	SetUserPINFn          func(ctx context.Context, input *domain.UserPIN) (bool, error)
	RegisterClientFn      func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error)
	AddIdentifierFn       func(ctx context.Context, clientID string, idType enums.IdentifierType, idValue string, isPrimary bool) (*domain.Identifier, error)
}

// NewCreateMock creates in itializes create type mocks
func NewCreateMock() *CreateMock {
	return &CreateMock{
		RegisterClientFn: func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()

			return &domain.ClientUserProfile{
				User: &domain.User{
					ID:                  &ID,
					FirstName:           "FirstName",
					LastName:            "Last Name",
					Username:            "User Name",
					MiddleName:          "Middle Name",
					DisplayName:         "Display Name",
					Gender:              enumutils.GenderMale,
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
				},
				Client: &domain.ClientProfile{
					ID:             &ID,
					UserID:         &ID,
					ClientType:     enums.ClientTypeOvc,
					HealthRecordID: &ID,
				},
			}, nil
		},

		GetOrCreateFacilityFn: func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &domain.Facility{
				ID:          &id,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},

		CollectMetricsFn: func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
			metricID := uuid.New().String()
			return &domain.Metric{
				MetricID:  &metricID,
				Type:      enums.EngagementMetrics,
				Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
				Timestamp: time.Now(),
				UID:       ksuid.New().String(),
			}, nil
		},

		SetUserPINFn: func(ctx context.Context, input *domain.UserPIN) (bool, error) {
			return true, nil
		},
		RegisterStaffUserFn: func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
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

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (f *CreateMock) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.GetOrCreateFacilityFn(ctx, facility)
}

// CollectMetrics mocks the implementation of `gorm's` CollectMetrics method.
func (f *CreateMock) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
	return f.CollectMetricsFn(ctx, metric)
}

//SetUserPIN mocks the implementation of SetUserPIN method
func (f *CreateMock) SetUserPIN(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return f.SetUserPINFn(ctx, pinData)
}

// RegisterStaffUser mocks the implementation of  RegisterStaffUser method.
func (f *CreateMock) RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	return f.RegisterStaffUserFn(ctx, user, staff)
}

// RegisterClient mocks the implementation of `gorm's` RegisterClient method
func (f *CreateMock) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	return f.RegisterClientFn(ctx, userInput, clientInput)
}

// AddIdentifier mocks the implementation of `gorm's` AddIdentifier method
func (f *CreateMock) AddIdentifier(ctx context.Context, clientID string, idType enums.IdentifierType, idValue string, isPrimary bool) (*domain.Identifier, error) {
	return f.AddIdentifierFn(ctx, clientID, idType, idValue, isPrimary)
}

// QueryMock is a mock of the query methods
type QueryMock struct {
	RetrieveFacilityFn           func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCodeFn  func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetFacilitiesFn              func(ctx context.Context) ([]*domain.Facility, error)
	GetUserPINByUserIDFn         func(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetUserProfileByUserIDFn     func(ctx context.Context, userID string, flavour string) (*domain.User, error)
	GetClientProfileByClientIDFn func(ctx context.Context, clientID string) (*domain.ClientProfile, error)
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
func (f *QueryMock) GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*domain.User, error) {
	return f.GetUserProfileByUserIDFn(ctx, userID, flavour)
}

// GetClientProfileByClientID defines a mock for fetching a client profile using the client's ID
func (f *QueryMock) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	return f.GetClientProfileByClientIDFn(ctx, clientID)
}

// UpdateMock ...
type UpdateMock struct {
	UpdateUserLastSuccessfulLoginFn func(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error
	UpdateUserLastFailedLoginFn     func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error
	UpdateUserFailedLoginCountFn    func(ctx context.Context, userID string, failedLoginCount string, flavour string) error
	UpdateUserNextAllowedLoginFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error
}

// NewUpdateMock initializes a new instance of `GormMock` then mocking the case of success.
func NewUpdateMock() *UpdateMock {
	return &UpdateMock{
		UpdateUserLastSuccessfulLoginFn: func(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
			return nil
		},

		UpdateUserLastFailedLoginFn: func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
			return nil
		},

		UpdateUserFailedLoginCountFn: func(ctx context.Context, userID, failedLoginCount, flavour string) error {
			return nil
		},

		UpdateUserNextAllowedLoginFn: func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
			return nil
		},
	}
}

//UpdateUserLastSuccessfulLogin ...
func (um *UpdateMock) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	return um.UpdateUserLastSuccessfulLoginFn(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (um *UpdateMock) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	return um.UpdateUserLastFailedLoginFn(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (um *UpdateMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	return um.UpdateUserFailedLoginCountFn(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (um *UpdateMock) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	return um.UpdateUserNextAllowedLoginFn(ctx, userID, nextAllowedLoginTime, flavour)
}

// DeleteMock ....
type DeleteMock struct{}

// NewDeleteMock initializes a new instance of `GormMock` then mocking the case of success.
func NewDeleteMock() *DeleteMock {
	return &DeleteMock{}
}
