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
	GetOrCreateFacilityFn     func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	CollectMetricsFn          func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
	GetOrCreateStaffUserserFn func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error)
	SavePinFn                 func(ctx context.Context, input *domain.UserPIN) (bool, error)
	RegisterClientFn          func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error)
	AddIdentifierFn           func(ctx context.Context, clientID string, idType enums.IdentifierType, idValue string, isPrimary bool) (*domain.Identifier, error)
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

		SavePinFn: func(ctx context.Context, input *domain.UserPIN) (bool, error) {
			return true, nil
		},
		GetOrCreateStaffUserserFn: func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
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

//SavePin mocks the implementation of SavePin method
func (f *CreateMock) SavePin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return f.SavePinFn(ctx, pinData)
}

// GetOrCreateStaffUser mocks the implementation of  GetOrCreateStaffUser method.
func (f *CreateMock) GetOrCreateStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	return f.GetOrCreateStaffUserserFn(ctx, user, staff)
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
