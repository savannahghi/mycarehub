package facility_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	healthCRMMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
	serviceRequestMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	"github.com/savannahghi/scalarutils"
)

func TestUseCaseFacilityImpl_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()

	type args struct {
		ctx      context.Context
		id       *string
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility",
			args: args{
				ctx:      ctx,
				id:       &ID,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case - no id",
			args: args{
				ctx:      ctx,
				isActive: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to retrieve facility by id",
			args: args{
				ctx:      ctx,
				id:       &ID,
				isActive: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get facility business hours",
			args: args{
				ctx:      ctx,
				id:       &ID,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFacility := mock.NewFacilityUsecaseMock()

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case - no id" {
				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
				}
			}
			if tt.name == "Sad case - unable to retrieve facility by id" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case - unable to get facility business hours" {
				fakeHealthCRM.MockGetCRMFacilityByIDFn = func(ctx context.Context, id string) (*domain.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}

			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_RetrieveFacilityByIdentifier_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
		isActive   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retrieve facility by identifier",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
				isActive: true,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to retrieve facility by identifier",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
				isActive: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing input",
			args: args{
				ctx:      ctx,
				isActive: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case: unable to retrieve facility by identifier" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by Identifier")
				}
			}

			got, err := f.RetrieveFacilityByIdentifier(tt.args.ctx, &tt.args.identifier, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUnit_ListProgramFacilities(t *testing.T) {
	ctx := context.Background()
	searchTerm := "term"
	programID := gofakeit.UUID()

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}

	type args struct {
		ctx              context.Context
		programID        *string
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Happy case: list facilities for a different program",
			args: args{
				ctx:        ctx,
				programID:  &programID,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Happy case- empty search input",
			args: args{
				ctx:              ctx,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case- failed to get logged in user",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- failed to get user profile",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- failed to list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case- failed to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case- failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case- failed to list facilities" {
				fakeDB.MockListProgramFacilitiesFn = func(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to list facilities")
				}
			}

			got, err := f.ListProgramFacilities(tt.args.ctx, tt.args.programID, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.ListProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_Inactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockInactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.InactivateFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.Inactivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUseCaseFacilityImpl_Reactivate_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - very bad mflCode",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad Case - empty mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - invalid mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - very bad mflCode" {
				fakeDB.MockReactivateFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.ReactivateFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestUseCaseFacilityImpl_DeleteFacility(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		identifier dto.FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete facility",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable delete facility",
			args: args{
				ctx: ctx,
				identifier: dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "30290320932",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Happy Case - Successfully delete facility" {
				fakeFacility.MockDeleteFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad Case - unable delete facility" {
				fakeFacility.MockDeleteFacilityFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := f.DeleteFacility(tt.args.ctx, &tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseFacilityImpl_ListFacilities(t *testing.T) {
	ctx := context.Background()
	searchTerm := "term"

	paginationInput := dto.PaginationsInput{
		Limit:       1,
		CurrentPage: 1,
	}

	type args struct {
		ctx              context.Context
		searchTerm       *string
		filterInput      []*dto.FiltersInput
		paginationsInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Happy case- empty search input",
			args: args{
				ctx:              ctx,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case- failed to list facilities",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility business hours",
			args: args{
				ctx:        ctx,
				searchTerm: &searchTerm,
				filterInput: []*dto.FiltersInput{
					{
						DataType: enums.FilterSortDataTypeName,
						Value:    "value",
					},
				},
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case- failed to list facilities" {
				fakeDB.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to list facilities")
				}
			}
			if tt.name == "Sad case: unable to get facility business hours" {
				fakeHealthCRM.MockGetCRMFacilityByIDFn = func(ctx context.Context, id string) (*domain.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}

			got, err := f.ListFacilities(tt.args.ctx, tt.args.searchTerm, tt.args.filterInput, tt.args.paginationsInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.ListFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_SyncFacilities(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeExt := extensionMock.NewFakeExtension()
	fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
	fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: ctx,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to notify pubsub",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get facilities without FHIROrganisationID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - unable to notify pubsub" {
				fakePubsub.MockNotifyCreateOrganizationFn = func(ctx context.Context, facility *domain.Facility) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - unable to get facilities without FHIROrganisationID" {
				fakeDB.MockGetFacilitiesWithoutFHIRIDFn = func(ctx context.Context) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			err := f.SyncFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.SyncFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_AddFacilityContact(t *testing.T) {

	type args struct {
		ctx        context.Context
		facilityID string
		contact    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: success adding facility contact",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    interserviceclient.TestUserPhoneNumber,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: fail to normalize phone number",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    "072897",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to update facility",
			args: args{
				ctx:        context.Background(),
				facilityID: gofakeit.UUID(),
				contact:    interserviceclient.TestUserPhoneNumber,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "sad case: fail to update facility" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update facility")
				}
			}

			got, err := f.AddFacilityContact(tt.args.ctx, tt.args.facilityID, tt.args.contact)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.AddFacilityContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.AddFacilityContact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseFacilityImpl_CreateFacilities(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilities []*dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "392893828",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
							Lng: "2.4565",
						},
						BusinessHours: []dto.BusinessHoursInput{
							{
								Day:         "MONDAY",
								OpeningTime: "08:00",
								ClosingTime: "17:00",
							},
						},
						Services: []dto.FacilityServiceInput{
							{
								Name:        "Weight",
								Description: "Weight(Kgs)",
								Identifiers: []dto.ServiceIdentifierInput{
									{
										IdentifierType:  "CIEL",
										IdentifierValue: "1234",
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: no facilities input",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to parse float coordinates (lat)",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "392893828",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to parse float coordinates (lng)",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "392893828",
						},
						Coordinates: dto.CoordinatesInput{
							Lng: "2.4565",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create health crm facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "392893828",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
							Lng: "2.4567",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create facilities in mch",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "09090908",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
							Lng: "2.4567",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to notify create facility in clinical service",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "09090908",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
							Lng: "2.4567",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to notify create facility in cms",
			args: args{
				ctx: context.Background(),
				facilities: []*dto.FacilityInput{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: dto.FacilityIdentifierInput{
							Type:  enums.FacilityIdentifierTypeMFLCode,
							Value: "09090908",
						},
						Coordinates: dto.CoordinatesInput{
							Lat: "34.4565",
							Lng: "2.4567",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case: failed to create health crm facilities" {
				fakeHealthCRM.MockCreateFacilityFn = func(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("unable to eat facility")
				}
			}
			if tt.name == "Sad case: failed to create facilities in mch" {
				fakeDB.MockCreateFacilitiesFn = func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to notify create facility in clinical service" {
				fakePubsub.MockNotifyCreateOrganizationFn = func(ctx context.Context, facility *domain.Facility) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to notify create facility in cms" {
				fakePubsub.MockNotifyCreateCMSFacilityFn = func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := f.CreateFacilities(tt.args.ctx, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CreateFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a a value to be returned, got: %v", got)
			}
		})
	}
}

func TestUseCaseFacilityImpl_PublishFacilitiesToCMS(t *testing.T) {
	id := gofakeit.UUID()
	id2 := gofakeit.UUID()
	type args struct {
		ctx        context.Context
		facilities []*domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: publish facilities to cms",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						ID:   &id,
						Name: gofakeit.BS(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed publish facilities to cms",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						ID:   &id2,
						Name: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()

			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case: failed publish facilities to cms" {
				fakePubsub.MockNotifyCreateCMSFacilityFn = func(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if err := f.PublishFacilitiesToCMS(tt.args.ctx, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.PublishFacilitiesToCMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCaseFacilityImpl_CmdAddFacilityToProgram(t *testing.T) {

	type args struct {
		ctx         context.Context
		facilityIDs []string
		programID   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: add facilities to program",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: fail to add facility to program",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to notify cms",
			args: args{
				ctx:         context.Background(),
				facilityIDs: []string{gofakeit.UUID()},
				programID:   gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		fakeDB := pgMock.NewPostgresMock()
		fakePubsub := pubsubMock.NewPubsubServiceMock()
		fakeExt := extensionMock.NewFakeExtension()
		fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
		fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

		f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

		if tt.name == "sad case: fail to add facility to program" {
			fakeDB.MockAddFacilityToProgramFn = func(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error) {
				return nil, fmt.Errorf("failed to add facility to program")
			}
		}
		if tt.name == "sad case: failed to notify cms" {
			fakePubsub.MockNotifyCMSAddFacilityToProgramFn = func(ctx context.Context, payload *dto.CMSLinkFacilityToProgramPayload) error {
				return fmt.Errorf("failed to notify cms")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.AddFacilityToProgram(tt.args.ctx, tt.args.facilityIDs, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CmdAddFacilityToProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseFacilityImpl.CmdAddFacilityToProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseFacilityImpl_GetNearbyFacilities(t *testing.T) {
	latitude := -1.2979512335313856
	longitude := 36.78882506563385
	type args struct {
		ctx             context.Context
		locationInput   *dto.LocationInput
		paginationInput dto.PaginationsInput
		serviceIDs      []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get nearby facilities successfully",
			args: args{
				ctx: context.Background(),
				locationInput: &dto.LocationInput{
					Lat: latitude,
					Lng: longitude,
				},
				paginationInput: dto.PaginationsInput{
					Limit:       20,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get facilities",
			args: args{
				ctx: context.Background(),
				locationInput: &dto.LocationInput{
					Lat: latitude,
					Lng: longitude,
				},
				paginationInput: dto.PaginationsInput{
					Limit:       20,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad Case - Fail to get facilities" {
				fakeHealthCRM.MockGetFacilitiesFn = func(ctx context.Context, location *dto.LocationInput, serviceIDs []string, searchParameter string, pagination *domain.Pagination) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}

			_, err := f.GetNearbyFacilities(tt.args.ctx, tt.args.locationInput, tt.args.serviceIDs, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetNearbyFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_GetServices(t *testing.T) {
	type args struct {
		ctx        context.Context
		pagination *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get services",
			args: args{
				ctx: context.Background(),
				pagination: &dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get services",
			args: args{
				ctx: context.Background(),
				pagination: &dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case: unable to get services" {
				fakeHealthCRM.MockGetServicesFn = func(ctx context.Context, pagination *domain.Pagination) (*domain.FacilityServicePage, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := f.GetServices(tt.args.ctx, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_SearchFacilitiesByService(t *testing.T) {
	type args struct {
		ctx           context.Context
		locationInput *dto.LocationInput
		serviceName   string
		pagination    *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FacilityPage
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully get facilities by the service name",
			args: args{
				serviceName: "Prep",
				pagination: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Empty service name",
			args: args{
				serviceName: "",
				pagination: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get facilities by the service name",
			args: args{
				serviceName: "Prep",
				pagination: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad Case: Fail to get facilities by the service name" {
				fakeHealthCRM.MockGetFacilitiesFn = func(ctx context.Context, location *dto.LocationInput, serviceIDs []string, searchParameter string, pagination *domain.Pagination) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}

			got, err := f.SearchFacilitiesByService(tt.args.ctx, tt.args.locationInput, tt.args.serviceName, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.SearchFacilitiesByService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("UseCaseFacilityImpl.SearchFacilitiesByService() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestUseCaseFacilityImpl_BookService(t *testing.T) {
	newTime := "2023-10-29T14:30:00.000Z"
	type args struct {
		ctx                context.Context
		facilityID         string
		serviceIDs         []string
		serviceBookingTime scalarutils.DateTime
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: book successfully",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to book a service",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get a facility from health crm",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create service request",
			args: args{
				ctx:                context.Background(),
				facilityID:         gofakeit.UUID(),
				serviceIDs:         []string{gofakeit.UUID()},
				serviceBookingTime: scalarutils.DateTime(newTime),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeHealthCRM := healthCRMMock.NewHealthServiceMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub, fakeExt, fakeHealthCRM, fakeServiceRequest)

			if tt.name == "Sad case: unable to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to get client profile" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to book a service" {
				fakeDB.MockCreateBookingFn = func(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to get a facility from health crm" {
				fakeHealthCRM.MockGetCRMFacilityByIDFn = func(ctx context.Context, id string) (*domain.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to create service request" {
				fakeServiceRequest.MockCreateServiceRequestFn = func(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
					return false, fmt.Errorf("error")
				}
			}

			_, err := f.BookService(tt.args.ctx, tt.args.facilityID, tt.args.serviceIDs, &tt.args.serviceBookingTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.BookService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
