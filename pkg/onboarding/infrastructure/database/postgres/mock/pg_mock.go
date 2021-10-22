package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	//Get
	GetOrCreateFacilityFn    func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	GetFacilitiesFn          func(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityFn       func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	SetUserPINFn             func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
	GetUserPINByUserIDFn     func(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetUserProfileByUserIDFn func(ctx context.Context, userID string, flavour string) (*domain.User, error)
	RegisterStaffUserFn      func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error)

	//Updates
	UpdateUserLastSuccessfulLoginFn func(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error
	UpdateUserLastFailedLoginFn     func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error
	UpdateUserFailedLoginCountFn    func(ctx context.Context, userID string, failedLoginCount string, flavour string) error
	UpdateUserNextAllowedLoginFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	return &PostgresMock{
		GetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
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
		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return []*domain.Facility{
				{
					ID:          &id,
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			}, nil
		},
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
		SetUserPINFn: func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
			return true, nil
		},

		GetUserProfileByUserIDFn: func(ctx context.Context, userID, flavour string) (*domain.User, error) {
			id := uuid.New().String()
			contact := &domain.Contact{
				ID:      &id,
				Type:    enums.PhoneContact,
				Contact: "test",
				Active:  true,
				OptedIn: true,
			}
			time := time.Now()
			return &domain.User{
				ID:          &id,
				Username:    "test",
				DisplayName: "test",
				FirstName:   "test",
				MiddleName:  "test",
				LastName:    "test",
				UserType:    enums.HealthcareWorkerUser,
				Gender:      enumutils.GenderMale,
				Active:      false,
				Contacts:    []*domain.Contact{contact},
				Languages:   []enumutils.Language{enumutils.LanguageEn},
				// PushTokens:          []string{"push-token"},
				LastSuccessfulLogin: &time,
				LastFailedLogin:     &time,
				FailedLoginCount:    "test",
				NextAllowedLogin:    &time,
				TermsAccepted:       false,
				AcceptedTermsID:     "test",
				Flavour:             feedlib.FlavourPro,
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
				},
			}, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

//SetUserPIN mocks the implementation of SetUserPIN method
func (gm *PostgresMock) SetUserPIN(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return gm.SetUserPINFn(ctx, pinData)
}

// GetUserPINByUserID ...
func (gm *PostgresMock) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return gm.GetUserPINByUserIDFn(ctx, userID)
}

// GetUserProfileByUserID gets user profile by user ID
func (gm *PostgresMock) GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*domain.User, error) {
	return gm.GetUserProfileByUserIDFn(ctx, userID, flavour)
}

//UpdateUserLastSuccessfulLogin ...
func (gm *PostgresMock) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	return gm.UpdateUserLastSuccessfulLoginFn(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (gm *PostgresMock) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	return gm.UpdateUserLastFailedLoginFn(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (gm *PostgresMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	return gm.UpdateUserFailedLoginCountFn(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (gm *PostgresMock) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	return gm.UpdateUserNextAllowedLoginFn(ctx, userID, nextAllowedLoginTime, flavour)
}

// RegisterStaffUser mocks the implementation of `gorm's` RegisterStaffUser method.
func (gm *PostgresMock) RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	return gm.RegisterStaffUserFn(ctx, user, staff)
}
