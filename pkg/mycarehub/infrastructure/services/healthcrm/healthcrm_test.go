package healthcrm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/healthcrm"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	healthCRMSvc "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm"
	mockHealthCRM "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm/mock"
)

func TestHealthCRMImpl_CreateFacility(t *testing.T) {
	uuid := gofakeit.UUID()

	type args struct {
		ctx      context.Context
		facility []*domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create facility",
			args: args{
				ctx: context.Background(),
				facility: []*domain.Facility{
					{
						ID:                 &uuid,
						Name:               gofakeit.BeerName(),
						Phone:              interserviceclient.TestUserPhoneNumber,
						Active:             true,
						Country:            gofakeit.CountryAbr(),
						County:             gofakeit.Country(),
						Address:            gofakeit.Address().Address,
						Description:        gofakeit.BeerMalt(),
						FHIROrganisationID: gofakeit.UUID(),
						Identifiers: []*domain.FacilityIdentifier{
							{
								ID:     uuid,
								Active: true,
								Type:   enums.FacilityIdentifierTypeMFLCode,
								Value:  "11094",
							},
						},
						WorkStationDetails: domain.WorkStationDetails{},
						Coordinates:        &domain.Coordinates{},
						BusinessHours: []domain.BusinessHours{
							{
								ID:          gofakeit.UUID(),
								Day:         "MONDAY",
								OpeningTime: "12:12",
								ClosingTime: "12:12",
								FacilityID:  gofakeit.UUID(),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create facility",
			args: args{
				ctx: context.Background(),
				facility: []*domain.Facility{
					{
						ID:                 &uuid,
						Name:               gofakeit.BeerName(),
						Phone:              interserviceclient.TestUserPhoneNumber,
						Active:             true,
						Country:            gofakeit.CountryAbr(),
						County:             gofakeit.Country(),
						Address:            gofakeit.Address().Address,
						Description:        gofakeit.BeerMalt(),
						FHIROrganisationID: gofakeit.UUID(),
						Identifiers: []*domain.FacilityIdentifier{
							{
								ID:     uuid,
								Active: true,
								Type:   enums.FacilityIdentifierTypeMFLCode,
								Value:  "11094",
							},
						},
						WorkStationDetails: domain.WorkStationDetails{},
						Coordinates:        &domain.Coordinates{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeHealthCRM := mockHealthCRM.NewHealthCRMClientMock()
			h := healthCRMSvc.NewHealthCRMService(fakeHealthCRM)

			if tt.name == "Sad case: unable to create facility" {
				fakeHealthCRM.MockCreateFacilityFn = func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := h.CreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCRMImpl.CreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestHealthCRMImpl_GetServices(t *testing.T) {
	type args struct {
		ctx        context.Context
		pagination *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get services offered in a facility",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 20,
					Count:       100,
					TotalPages:  0,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get services offered in a facility",
			args: args{
				ctx: context.Background(),
				pagination: &domain.Pagination{
					Limit:       1,
					CurrentPage: 20,
					Count:       100,
					TotalPages:  0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeHealthCRM := mockHealthCRM.NewHealthCRMClientMock()
			h := healthCRMSvc.NewHealthCRMService(fakeHealthCRM)

			if tt.name == "Sad case: unable to get services offered in a facility" {
				fakeHealthCRM.MockGetServicesFn = func(ctx context.Context, pagination *healthcrm.Pagination) (*healthcrm.FacilityServicePage, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := h.GetServices(tt.args.ctx, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCRMImpl.GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestHealthCRMImpl_GetCRMFacilityByID(t *testing.T) {
	id := gofakeit.UUID()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "Happy case: get facility by id",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get facility by id",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeHealthCRM := mockHealthCRM.NewHealthCRMClientMock()
			h := healthCRMSvc.NewHealthCRMService(fakeHealthCRM)

			if tt.name == "Sad case: unable to get facility by id" {
				fakeHealthCRM.MockGetFacilityByIDFn = func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := h.GetCRMFacilityByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCRMImpl.GetCRMFacilityByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestHealthCRMImpl_GetFacilities(t *testing.T) {
	latitude := -1.2979512335313856
	longitude := 36.78882506563385
	type args struct {
		ctx             context.Context
		location        *dto.LocationInput
		serviceIDs      []string
		searchParameter string
		pagination      *domain.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Facility
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get facilities",
			args: args{
				ctx: context.Background(),
				location: &dto.LocationInput{
					Lat: latitude,
					Lng: longitude,
				},
				serviceIDs: []string{},
				pagination: &domain.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get facilities",
			args: args{
				ctx: context.Background(),
				location: &dto.LocationInput{
					Lat: latitude,
					Lng: longitude,
				},
				serviceIDs: []string{},
				pagination: &domain.Pagination{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeHealthCRM := mockHealthCRM.NewHealthCRMClientMock()
			h := healthCRMSvc.NewHealthCRMService(fakeHealthCRM)

			if tt.name == "Sad Case: Fail to get facilities" {
				fakeHealthCRM.MockGetFacilitiesFn = func(ctx context.Context, location *healthcrm.Coordinates, serviceIDs []string, searchParameter string, pagination *healthcrm.Pagination) (*healthcrm.FacilityPage, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}

			got, err := h.GetFacilities(tt.args.ctx, tt.args.location, tt.args.serviceIDs, tt.args.searchParameter, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCRMImpl.GetFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("HealthCRMImpl.GetFacilities() = expected a response but got %v", got)
				}
			}
		})
	}
}

func TestHealthCRMImpl_CheckIfServiceExists(t *testing.T) {
	type args struct {
		ctx        context.Context
		serviceIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check if service exists",
			args: args{
				ctx:        context.Background(),
				serviceIDs: []string{gofakeit.UUID()},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeHealthCRM := mockHealthCRM.NewHealthCRMClientMock()
			h := healthCRMSvc.NewHealthCRMService(fakeHealthCRM)

			got, err := h.CheckIfServiceExists(tt.args.ctx, tt.args.serviceIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCRMImpl.CheckIfServiceExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HealthCRMImpl.CheckIfServiceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
