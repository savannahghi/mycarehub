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
	SavePinFn             func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
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
			county := enums.CountyTypeNairobi
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
		SavePinFn: func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
			return true, nil
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

// SavePin mocks the save pin implementation
func (f *CreateMock) SavePin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return f.SavePinFn(ctx, pinData)
}

// QueryMock is a mock of the query methods
type QueryMock struct {
	RetrieveFacilityFn            func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCodeFn   func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetFacilitiesFn               func(ctx context.Context) ([]*domain.Facility, error)
	GetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*domain.User, error)
	ListFacilitiesFn              func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error)
	GetUserPINByUserIDFn          func(ctx context.Context, userID string) (*domain.UserPIN, error)
}

// NewQueryMock initializes a new instance of `GormMock` then mocking the case of success.
func NewQueryMock() *QueryMock {
	return &QueryMock{
		GetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*domain.User, error) {
			id := uuid.New().String()
			return &domain.User{
				ID: &id,
			}, nil
		},

		GetUserPINByUserIDFn: func(ctx context.Context, userID string) (*domain.UserPIN, error) {
			return &domain.UserPIN{}, nil
		},

		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := enums.CountyTypeNairobi
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
			county := enums.CountyTypeNairobi
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
			county := enums.CountyTypeNairobi
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
		ListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := enums.CountyTypeNairobi
			description := "test description"
			nextPage := 1
			previousPage := 1
			return &domain.FacilityPage{
				Pagination: domain.Pagination{
					Limit:        1,
					CurrentPage:  1,
					Count:        1,
					TotalPages:   1,
					NextPage:     &nextPage,
					PreviousPage: &previousPage,
				},
				Facilities: []domain.Facility{
					{
						ID:          &facilityID,
						Name:        name,
						Code:        code,
						Active:      true,
						County:      county,
						Description: description,
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

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (f *QueryMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return f.GetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (f *QueryMock) ListFacilities(
	ctx context.Context,
	searchTerm *string,
	filterInput []*dto.FiltersInput,
	PaginationsInput dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return f.ListFacilitiesFn(ctx, searchTerm, filterInput, PaginationsInput)
}

// GetUserPINByUserID mocks the get user pin by ID implementation
func (f *QueryMock) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return f.GetUserPINByUserIDFn(ctx, userID)
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
