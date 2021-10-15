package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

// GormMock struct implements mocks of `gorm's`internal methods.
//
// This mock struct should be separate from our own internal methods.
type GormMock struct {
	GetOrCreateFacilityFn       func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	RetrieveFacilityFn          func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	RetrieveFacilityByMFLCodeFn func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error)
	GetFacilitiesFn             func(ctx context.Context) ([]gorm.Facility, error)
	DeleteFacilityFn            func(ctx context.Context, mfl_code string) (bool, error)
	CollectMetricsFn            func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error)
	RegisterStaffUserFn         func(ctx context.Context, user gorm.User, profile gorm.StaffProfile) (*gorm.StaffUserProfileTable, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
func NewGormMock() *GormMock {
	return &GormMock{
		GetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &id,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},

		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			var facilities []gorm.Facility
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			facilities = append(facilities, gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			})
			return facilities, nil
		},

		DeleteFacilityFn: func(ctx context.Context, mfl_code string) (bool, error) {
			return true, nil
		},

		CollectMetricsFn: func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
			now := time.Now()
			metricID := uuid.New().String()
			return &gorm.Metric{
				MetricID:  &metricID,
				Type:      domain.EngagementMetrics,
				Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
				Timestamp: now,
				UID:       ksuid.New().String(),
			}, nil
		},

		RetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},
		RegisterStaffUserFn: func(ctx context.Context, user gorm.User, profile gorm.StaffProfile) (*gorm.StaffUserProfileTable, error) {
			userID := uuid.New().String()
			staffID := uuid.New().String()
			contactID := uuid.New().String()
			testTime := time.Now()
			facilityID := uuid.New().String()
			addressesID := uuid.New().String()

			testText := "testtext"

			userOutput := &gorm.User{
				UserID:      &userID,
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

			staffProfileOutput := &gorm.StaffProfile{
				StaffProfileID: &staffID,
				UserID:         &userID,
				StaffNumber:    "st1010101",
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
				Roles:             []domain.RoleType{domain.RoleTypePractitioner, domain.RoleTypeModerator}, //TODO: enum
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

			return &gorm.StaffUserProfileTable{
				User:         *userOutput,
				StaffProfile: *staffProfileOutput,
			}, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.GetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode string) (bool, error) {
	return gm.DeleteFacilityFn(ctx, mflcode)
}

// CollectMetrics mocks the implementation of  CollectMetrics method.
func (gm *GormMock) CollectMetrics(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
	return gm.CollectMetricsFn(ctx, metrics)
}

// RegisterStaffUser mocks the implementation of  RegisterStaffUser method.
func (gm *GormMock) RegisterStaffUser(ctx context.Context, user gorm.User, profile gorm.StaffProfile) (*gorm.StaffUserProfileTable, error) {
	return gm.RegisterStaffUserFn(ctx, user, profile)
}
