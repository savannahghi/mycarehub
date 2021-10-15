package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

// CreateMock is a mock of the create methods
type CreateMock struct {
	GetOrCreateFacilityFn func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	CollectMetricsFn      func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
	RegisterStaffUserFn   func(ctx context.Context, user dto.UserInput, profile dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error)
}

// NewCreateMock initializes a new instance of `GormMock` then mocking the case of success.
func NewCreateMock() *CreateMock {
	return &CreateMock{
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
				Type:      domain.EngagementMetrics,
				Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
				Timestamp: time.Now(),
				UID:       ksuid.New().String(),
			}, nil
		},
		RegisterStaffUserFn: func(ctx context.Context, user dto.UserInput, profile dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error) {
			userID := uuid.New().String()
			staffID := uuid.New().String()
			contactID := uuid.New().String()
			testTime := time.Now()
			facilityID := uuid.New().String()
			addressesID := uuid.New().String()

			testText := "testtext"

			userOutput := &domain.User{
				ID:          &userID,
				Username:    "user",
				DisplayName: "alias",
				FirstName:   "firstname",
				MiddleName:  &testText,
				LastName:    "lastname",
				UserType:    "doctor", //TODO: enum
				Gender:      "female", // TODO: enum
				Contacts: []*domain.Contact{
					{
						ID:      &contactID,
						Type:    "email",          //TODO: enum
						Contact: "user@email.com", //TODO: validate
						Active:  true,
						OptedIn: true,
					},
				},
				Languages:           []string{"en", "ksw"}, // TODO: slice of enums
				PushTokens:          []string{string(ksuid.New().String())},
				LastSuccessfulLogin: &testTime,
				LastFailedLogin:     &testTime,
				FailedLoginCount:    0,
				NextAllowedLogin:    &testTime,
				TermsAccepted:       true,
				AcceptedTermsID:     ksuid.New().String(), //TODO: add terms relation in db
			}

			staffProfileOutput := &domain.StaffProfile{
				ID:          &staffID,
				UserID:      &userID,
				StaffNumber: "st1010101",
				Facilities: []*domain.Facility{
					{
						ID:          &facilityID,
						Name:        "test-name",
						Code:        "c0032",
						Active:      true,
						County:      "Nakuru",
						Description: "This is just for mocking",
					},
				},
				DefaultFacilityID: &facilityID,
				Roles:             []domain.RoleType{domain.RoleTypePractitioner}, //TODO: enum
				Addresses: []*domain.UserAddress{
					{
						ID:         &addressesID,
						Type:       "postal", //TODO: enum
						Text:       "1123 Nairobi",
						Country:    "Kenya", //TODO: enum
						PostalCode: "10100",
						County:     "Nakuru", //TODO: counties belong to a country
						Active:     true,
					},
				},
			}

			return &domain.StaffUserProfileOutput{
				User:         *userOutput,
				StaffProfile: *staffProfileOutput,
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

// RegisterStaffUser mocks the implementation of `gorm's` RegisterStaffUser method
func (f *CreateMock) RegisterStaffUser(ctx context.Context, user dto.UserInput, profile dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error) {
	return f.RegisterStaffUserFn(ctx, user, profile)
}

// QueryMock is a mock of the query methods
type QueryMock struct {
	RetrieveFacilityFn          func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCodeFn func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetFacilitiesFn             func(ctx context.Context) ([]*domain.Facility, error)
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
