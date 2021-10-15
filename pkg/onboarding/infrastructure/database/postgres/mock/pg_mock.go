package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	GetOrCreateFacilityFn func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	GetFacilitiesFn       func(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityFn    func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RegisterStaffUserFn   func(ctx context.Context, user dto.UserInput, profile dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error)
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
				Roles:             []domain.RoleType{domain.RoleTypePractitioner},
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
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

// RegisterStaffUser mocks the implementation of `gorm's` RegisterStaffUser method.
func (gm *PostgresMock) RegisterStaffUser(ctx context.Context, user dto.UserInput, profile dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error) {
	return gm.RegisterStaffUserFn(ctx, user, profile)
}
