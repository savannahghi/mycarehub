package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// FacilityUsecaseMock mocks the implementation of facility usecase methods
type FacilityUsecaseMock struct {
	MockRetrieveFacilityFn             func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	MockRetrieveFacilityByIdentifierFn func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
	MockGetFacilitiesFn                func(ctx context.Context) ([]*domain.Facility, error)
	MockListProgramFacilitiesFn        func(ctx context.Context, programID *string, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockDeleteFacilityFn               func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	FetchFacilitiesFn                  func(ctx context.Context) ([]*domain.Facility, error)
	MockInactivateFacilityFn           func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockReactivateFacilityFn           func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockAddFacilityContactFn           func(ctx context.Context, facilityID string, contact string) (bool, error)
	MockListFacilitiesFn               func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockSyncFacilitiesFn               func(ctx context.Context) error
	MockCreateFacilitiesFn             func(ctx context.Context, facilities []*dto.FacilityInput) ([]*domain.Facility, error)
	MockPublishFacilitiesToCMSFn       func(ctx context.Context, facilities []*domain.Facility) error
	MockAddFacilityToProgramFn         func(ctx context.Context, facilityIDs []string, programID string) (bool, error)
	MockGetNearbyFacilitiesFn          func(ctx context.Context, locationInput *dto.LocationInput, serviceIDs []string, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error)
	MockGetServicesFn                  func(ctx context.Context, pagination *dto.PaginationsInput) (*dto.FacilityServiceOutputPage, error)
	MockSearchFacilitiesByServiceFn    func(ctx context.Context, locationInput *dto.LocationInput, serviceName string, pagination *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockBookServiceFn                  func(ctx context.Context, facilityID string, serviceIDs []string, serviceBookingTime time.Time) (*dto.BookingOutput, error)
	MockListBookingsFn                 func(ctx context.Context, clientID string, bookingState enums.BookingState, pagination dto.PaginationsInput) (*dto.BookingPage, error)
	MockVerifyBookingCodeFn            func(ctx context.Context, bookingID string, code string, programID string) (bool, error)
}

// NewFacilityUsecaseMock initializes a new instance of `GormMock` then mocking the case of success.
func NewFacilityUsecaseMock() *FacilityUsecaseMock {
	ID := uuid.New().String()
	name := "test facility"
	country := "Kenya"
	phone := interserviceclient.TestUserPhoneNumber
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()

	facilityInput := &domain.Facility{
		ID:                 &ID,
		Name:               name,
		Phone:              phone,
		Active:             true,
		Country:            country,
		County:             country,
		Address:            "1234",
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
		Identifiers: []*domain.FacilityIdentifier{
			{
				ID:     ID,
				Active: true,
				Type:   enums.FacilityIdentifierTypeHealthCRM,
				Value:  gofakeit.UUID(),
			},
		},
		WorkStationDetails: domain.WorkStationDetails{},
		Coordinates: &domain.Coordinates{
			Lat: -1.33456786787,
			Lng: 36.56789034552,
		},
		Services: []domain.FacilityService{
			{
				ID:          ID,
				Name:        name,
				Description: description,
				Identifiers: []domain.ServiceIdentifier{
					{
						ID:              ID,
						IdentifierType:  "CIEL",
						IdentifierValue: ID,
						ServiceID:       ID,
					},
				},
			},
		},
		BusinessHours: []domain.BusinessHours{
			{
				ID:          ID,
				Day:         "MONDAY",
				OpeningTime: "12:12",
				ClosingTime: "18:00",
				FacilityID:  ID,
			},
		},
	}

	var facilitiesList []*domain.Facility
	facilitiesList = append(facilitiesList, facilityInput)

	nextPage := 3
	previousPage := 1
	facilitiesPage := &domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:        1,
			CurrentPage:  2,
			Count:        3,
			TotalPages:   3,
			NextPage:     &nextPage,
			PreviousPage: &previousPage,
		},
		Facilities: []*domain.Facility{
			{
				ID:          &ID,
				Name:        name,
				Phone:       phone,
				Active:      true,
				Country:     country,
				Description: description,
			},
		},
	}

	return &FacilityUsecaseMock{

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockRetrieveFacilityByIdentifierFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockListProgramFacilitiesFn: func(ctx context.Context, programID *string, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},

		FetchFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockAddFacilityContactFn: func(ctx context.Context, facilityID string, contact string) (bool, error) {
			return true, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},
		MockSyncFacilitiesFn: func(ctx context.Context) error {
			return nil
		},
		MockCreateFacilitiesFn: func(ctx context.Context, facilities []*dto.FacilityInput) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockPublishFacilitiesToCMSFn: func(ctx context.Context, facilities []*domain.Facility) error {
			return nil
		},
		MockAddFacilityToProgramFn: func(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
			return true, nil
		},
		MockGetNearbyFacilitiesFn: func(ctx context.Context, locationInput *dto.LocationInput, serviceIDs []string, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
			return &domain.FacilityPage{
				Pagination: domain.Pagination{
					Limit:       50,
					CurrentPage: 1,
					Count:       0,
					TotalPages:  100,
				},
				Facilities: facilitiesList,
			}, nil
		},
		MockGetServicesFn: func(ctx context.Context, pagination *dto.PaginationsInput) (*dto.FacilityServiceOutputPage, error) {
			return &dto.FacilityServiceOutputPage{
				Results: []domain.FacilityService{
					{
						ID:          ID,
						Name:        name,
						Description: description,
						Identifiers: []domain.ServiceIdentifier{
							{
								ID:              ID,
								IdentifierType:  "CIEL",
								IdentifierValue: ID,
								ServiceID:       ID,
							},
						},
					},
				},
				Pagination: domain.Pagination{
					Limit:       10,
					CurrentPage: 30,
					Count:       45,
					TotalPages:  230,
				},
			}, nil
		},
		MockSearchFacilitiesByServiceFn: func(ctx context.Context, locationInput *dto.LocationInput, serviceName string, pagination *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return &domain.FacilityPage{
				Pagination: domain.Pagination{
					Limit:       50,
					CurrentPage: 1,
					Count:       0,
					TotalPages:  100,
				},
				Facilities: facilitiesList,
			}, nil
		},
		MockBookServiceFn: func(ctx context.Context, facilityID string, serviceIDs []string, serviceBookingTime time.Time) (*dto.BookingOutput, error) {
			return &dto.BookingOutput{
				ID: ID,
				Services: []domain.FacilityService{
					{
						ID:   gofakeit.UUID(),
						Name: gofakeit.BeerName(),
					},
				},
				Date:           time.Now(),
				Facility:       domain.Facility{},
				Client:         domain.ClientProfile{},
				OrganisationID: FHIROrganisationID,
				ProgramID:      ID,
			}, nil
		},
		MockVerifyBookingCodeFn: func(ctx context.Context, bookingID, code, programID string) (bool, error) {
			return true, nil
		},
		MockListBookingsFn: func(ctx context.Context, clientID string, bookingState enums.BookingState, pagination dto.PaginationsInput) (*dto.BookingPage, error) {
			return &dto.BookingPage{
				Results: []dto.BookingOutput{
					{
						ID: ID,
						Services: []domain.FacilityService{
							{
								ID:   gofakeit.UUID(),
								Name: gofakeit.BeerName(),
							},
						},
						Date:           time.Now(),
						Facility:       domain.Facility{},
						Client:         domain.ClientProfile{},
						OrganisationID: FHIROrganisationID,
						ProgramID:      ID,
					},
				},
				Pagination: domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			}, nil
		},
	}
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *FacilityUsecaseMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByIdentifier mocks the implementation of `gorm's` RetrieveFacilityByIdentifier method.
func (f *FacilityUsecaseMock) RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityByIdentifierFn(ctx, identifier, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (f *FacilityUsecaseMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.MockGetFacilitiesFn(ctx)
}

// ListProgramFacilities mocks the implementation of  ListProgramFacilities method.
func (f *FacilityUsecaseMock) ListProgramFacilities(
	ctx context.Context,
	programID *string,
	searchTerm *string,
	filterInput []*dto.FiltersInput,
	paginationsInput *dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return f.MockListProgramFacilitiesFn(ctx, programID, searchTerm, filterInput, paginationsInput)
}

// DeleteFacility mocks the implementation of deleting a facility by ID
func (f *FacilityUsecaseMock) DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockDeleteFacilityFn(ctx, identifier)
}

// FetchFacilities mocks the implementation of fetch all facilities
func (f *FacilityUsecaseMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.FetchFacilitiesFn(ctx)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (f *FacilityUsecaseMock) InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockInactivateFacilityFn(ctx, identifier)
}

// ReactivateFacility mocks the implementation of reactivating the active status of a particular facility
func (f *FacilityUsecaseMock) ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockReactivateFacilityFn(ctx, identifier)
}

// AddFacilityContact mock the implementation of the AddFacilityContact method
func (f *FacilityUsecaseMock) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	return f.MockAddFacilityContactFn(ctx, facilityID, contact)
}

// ListFacilities mock the implementation of the ListFacilities method
func (f *FacilityUsecaseMock) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	return f.MockListFacilitiesFn(ctx, searchTerm, filterInput, paginationsInput)
}

// SyncFacilities mock the implementation of the SyncFacilities method
func (f *FacilityUsecaseMock) SyncFacilities(ctx context.Context) error {
	return f.MockSyncFacilitiesFn(ctx)
}

// CreateFacilities Mocks the implementation of CreateFacilities method
func (f *FacilityUsecaseMock) CreateFacilities(ctx context.Context, facilities []*dto.FacilityInput) ([]*domain.Facility, error) {
	return f.MockCreateFacilitiesFn(ctx, facilities)
}

// PublishFacilitiesToCMS mock the implementation of the PublishFacilitiesToCMS method
func (f *FacilityUsecaseMock) PublishFacilitiesToCMS(ctx context.Context, facilities []*domain.Facility) error {
	return f.MockPublishFacilitiesToCMSFn(ctx, facilities)
}

// AddFacilityToProgram mock the implementation of the CmdAddFacilityToProgram method
func (f *FacilityUsecaseMock) AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
	return f.MockAddFacilityToProgramFn(ctx, facilityIDs, programID)
}

// GetNearbyFacilities mocks the implementation of getting nearby facilities
func (f *FacilityUsecaseMock) GetNearbyFacilities(ctx context.Context, locationInput *dto.LocationInput, serviceIDs []string, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return f.MockGetNearbyFacilitiesFn(ctx, locationInput, serviceIDs, paginationInput)
}

// GetServices mocks the implementation of getting available services from health cem
func (f *FacilityUsecaseMock) GetServices(ctx context.Context, pagination *dto.PaginationsInput) (*dto.FacilityServiceOutputPage, error) {
	return f.MockGetServicesFn(ctx, pagination)
}

// SearchFacilitiesByService mocks the implementation of searching facilities by a service name
func (f *FacilityUsecaseMock) SearchFacilitiesByService(ctx context.Context, locationInput *dto.LocationInput, serviceName string, pagination *dto.PaginationsInput) (*domain.FacilityPage, error) {
	return f.MockSearchFacilitiesByServiceFn(ctx, locationInput, serviceName, pagination)
}

// BookService is used to mock the booking of a service
func (f *FacilityUsecaseMock) BookService(ctx context.Context, facilityID string, serviceIDs []string, serviceBookingTime time.Time) (*dto.BookingOutput, error) {
	return f.MockBookServiceFn(ctx, facilityID, serviceIDs, serviceBookingTime)
}

// VerifyBookingCode mocks the implementation of verifying booking code
func (f *FacilityUsecaseMock) VerifyBookingCode(ctx context.Context, bookingID string, code string, programID string) (bool, error) {
	return f.MockVerifyBookingCodeFn(ctx, bookingID, code, programID)
}

// ListBookings mocks the implementation that lists client bookings
func (f *FacilityUsecaseMock) ListBookings(ctx context.Context, clientID string, bookingState enums.BookingState, pagination dto.PaginationsInput) (*dto.BookingPage, error) {
	return f.MockListBookingsFn(ctx, clientID, bookingState, pagination)
}
