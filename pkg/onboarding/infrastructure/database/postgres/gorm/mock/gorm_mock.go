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
	SetUserPINFn                func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	GetUserPINByUserIDFn        func(ctx context.Context, userID string) (*gorm.PINData, error)
	GetUserProfileByUserIDFn    func(ctx context.Context, userID string, flavour string) (*gorm.User, error)

	//Updates
	UpdateUserLastSuccessfulLoginFn func(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error
	UpdateUserLastFailedLoginFn     func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error
	UpdateUserFailedLoginCountFn    func(ctx context.Context, userID string, failedLoginCount string, flavour string) error
	UpdateUserNextAllowedLoginFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error
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

		SetUserPINFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},

		GetUserProfileByUserIDFn: func(ctx context.Context, userID, flavour string) (*gorm.User, error) {
			id := uuid.New().String()
			usercontact := &gorm.Contact{
				ContactID: &id,
				Type:      "test",
				Contact:   "test",
				Active:    true,
				OptedIn:   true,
			}
			time := time.Now()
			return &gorm.User{
				Base:                gorm.Base{},
				UserID:              &id,
				Username:            "test",
				DisplayName:         "test",
				FirstName:           "test",
				MiddleName:          "test",
				LastName:            "test",
				Flavour:             "test",
				UserType:            "test",
				Gender:              "test",
				Active:              false,
				Contacts:            []gorm.Contact{*usercontact},
				Languages:           []string{"en"},
				PushTokens:          []string{"push-token"},
				LastSuccessfulLogin: &time,
				LastFailedLogin:     &time,
				FailedLoginCount:    "test",
				NextAllowedLogin:    &time,
				TermsAccepted:       false,
				AcceptedTermsID:     "test",
			}, nil
		},

		GetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return &gorm.PINData{
				UserID:    userID,
				HashedPIN: "mbzcbvhbxchjbvhdbvhhjdfskgbfhas832y38hjsdnfkjbh73y73y72",
				ValidFrom: time.Now(),
				ValidTo:   time.Now(),
				Flavour:   "CONSUMER",
				IsValid:   true,
				Salt:      "test-salt",
			}, nil
		},

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

//SetUserPIN mocks the implementation of SetUserPIN method
func (gm *GormMock) SetUserPIN(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.SetUserPINFn(ctx, pinData)
}

// GetUserPINByUserID ...
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.GetUserPINByUserIDFn(ctx, userID)
}

// GetUserProfileByUserID gets user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*gorm.User, error) {
	return gm.GetUserProfileByUserIDFn(ctx, userID, flavour)
}

//UpdateUserLastSuccessfulLogin ...
func (gm *GormMock) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	return gm.UpdateUserLastSuccessfulLoginFn(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (gm *GormMock) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	return gm.UpdateUserLastFailedLoginFn(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (gm *GormMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	return gm.UpdateUserFailedLoginCountFn(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (gm *GormMock) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	return gm.UpdateUserNextAllowedLoginFn(ctx, userID, nextAllowedLoginTime, flavour)
}
