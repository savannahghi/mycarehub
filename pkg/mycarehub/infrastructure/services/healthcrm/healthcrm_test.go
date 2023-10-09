package healthcrm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/healthcrm"
	"github.com/savannahghi/interserviceclient"
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
						Identifier:         []*domain.FacilityIdentifier{},
						WorkStationDetails: domain.WorkStationDetails{},
						Coordinates:        &domain.Coordinates{},
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
						Identifier:         []*domain.FacilityIdentifier{},
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
