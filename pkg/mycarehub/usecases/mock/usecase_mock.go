package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateMock is a mock of the create methods
type CreateMock struct {
	GetOrCreateFacilityFn func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	RegisterClientFn      func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error)
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
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (f *CreateMock) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.GetOrCreateFacilityFn(ctx, facility)
}

// RegisterClient mocks the implementation of `gorm's` RegisterClient method
func (f *CreateMock) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	return f.RegisterClientFn(ctx, userInput, clientInput)
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

// UpdateMock ...
type UpdateMock struct {
}

// NewUpdateMock initializes a new instance of `GormMock` then mocking the case of success.
func NewUpdateMock() *UpdateMock {
	return &UpdateMock{}
}

// DeleteMock ....
type DeleteMock struct{}

// NewDeleteMock initializes a new instance of `GormMock` then mocking the case of success.
func NewDeleteMock() *DeleteMock {
	return &DeleteMock{}
}
