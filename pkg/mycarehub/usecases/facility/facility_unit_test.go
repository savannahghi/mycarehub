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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
)

// func TestUnit_CreateFacility(t *testing.T) {
// 	ctx := context.Background()
// 	name := "Kanairo One"
// 	county := "Nairobi"
// 	description := "This is just for mocking"

// 	type args struct {
// 		ctx        context.Context
// 		facility   dto.FacilityInput
// 		identifier dto.FacilityIdentifierInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantNil bool
// 		wantErr bool
// 	}{
// 		// {
// 		// 	name: "happy case - valid payload",
// 		// 	args: args{
// 		// 		ctx: ctx,
// 		// 		facility: dto.FacilityInput{
// 		// 			Name:        name,
// 		// 			Active:      true,
// 		// 			County:      county,
// 		// 			Description: description,
// 		// 		},
// 		// 		identifier: dto.FacilityIdentifierInput{
// 		// 			Type:  enums.FacilityIdentifierTypeMFLCode,
// 		// 			Value: "30290320932",
// 		// 		},
// 		// 	},
// 		// 	wantErr: false,
// 		// },
// 		{
// 			name: "sad case - facility code empty",
// 			args: args{
// 				ctx: ctx,
// 				facility: dto.FacilityInput{
// 					Name:        name,
// 					Active:      true,
// 					County:      county,
// 					Description: description,
// 				},
// 				identifier: dto.FacilityIdentifierInput{
// 					Type:  enums.FacilityIdentifierTypeMFLCode,
// 					Value: "30290320932",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Happy case - Create facility",
// 			args: args{
// 				ctx: ctx,
// 				facility: dto.FacilityInput{
// 					Name:        name,
// 					Active:      true,
// 					County:      county,
// 					Description: description,
// 				},
// 				identifier: dto.FacilityIdentifierInput{
// 					Type:  enums.FacilityIdentifierTypeMFLCode,
// 					Value: "30290320932",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			fakeDB := pgMock.NewPostgresMock()
// 			fakeFacility := mock.NewFacilityUsecaseMock()
// 			fakePubsub := pubsubMock.NewPubsubServiceMock()

// 			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

// 			if tt.name == "Happy case - Create facility" {
// 				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("failed query and retrieve facility by MFLCode")
// 				}
// 			}

// 			if tt.name == "sad case - facility code empty" {
// 				fakeFacility.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("failed to create facility")
// 				}
// 			}

// 			got, err := f.GetOrCreateFacility(tt.args.ctx, &tt.args.facility, &tt.args.identifier)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCaseFacilityImpl.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facility to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// 	// TODO: add teardown
// }

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
			name: "Happy case",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFacility := mock.NewFacilityUsecaseMock()

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

			if tt.name == "Sad case - no id" {
				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility")
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

// func TestUseCaseFacilityImpl_RetrieveFacilityByIdentifier_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	type args struct {
// 		ctx        context.Context
// 		identifier dto.FacilityIdentifierInput
// 		isActive   bool
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				identifier: dto.FacilityIdentifierInput{
// 					Type:  enums.FacilityIdentifierTypeMFLCode,
// 					Value: "30290320932",
// 				},
// 				isActive: true,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 				identifier: dto.FacilityIdentifierInput{
// 					Type:  enums.FacilityIdentifierTypeMFLCode,
// 					Value: "30290320932",
// 				},
// 				isActive: false,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			fakeDB := pgMock.NewPostgresMock()
// 			fakeFacility := mock.NewFacilityUsecaseMock()
// 			fakePubsub := pubsubMock.NewPubsubServiceMock()
// 			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

// 			if tt.name == "Sad case" {

// 				fakeFacility.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
// 					return nil, fmt.Errorf("an error occurred while retrieving facility by MFLCode")
// 				}
// 			}

// 			got, err := f.RetrieveFacilityByIdentifier(tt.args.ctx, &tt.args.identifier, tt.args.isActive)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCaseFacilityImpl.RetrieveFacilityByIdentifier() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.wantErr && got != nil {
// 				t.Errorf("expected facilities to be nil for %v", tt.name)
// 				return
// 			}

// 			if !tt.wantErr && got == nil {
// 				t.Errorf("expected facilities not to be nil for %v", tt.name)
// 				return
// 			}
// 		})
// 	}
// }

func TestUnit_ListFacilities(t *testing.T) {
	ctx := context.Background()

	searchTerm := "term"

	filterValue := "value"

	filterInput := []*dto.FiltersInput{
		{
			DataType: enums.FilterSortDataTypeName,
			Value:    filterValue,
		},
	}

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
			name: "Happy case",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case- empty search term",
			args: args{
				ctx:              ctx,
				searchTerm:       nil,
				filterInput:      filterInput,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- nil filter input",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      nil,
				paginationsInput: &paginationInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case- nil pagination input",
			args: args{
				ctx:              ctx,
				searchTerm:       &searchTerm,
				filterInput:      filterInput,
				paginationsInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeFacility := mock.NewFacilityUsecaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

			if tt.name == "Sad case- empty search term" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}
			}

			if tt.name == "Sad case- nil filter input" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
				}
			}

			if tt.name == "Sad case- nil pagination input" {
				fakeFacility.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("failed to list facilities")
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

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

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

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

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

			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

			if tt.name == "Happy Case - Successfully delete facility" {
				fakeFacility.DeleteFacilityFn = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad Case - unable delete facility" {
				fakeFacility.DeleteFacilityFn = func(ctx context.Context, id int) (bool, error) {
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

func TestUseCaseFacilityImpl_FetchFacilities(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx             context.Context
		searchParameter *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully fetch facilities",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()

			fakePubsub := pubsubMock.NewPubsubServiceMock()
			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

			got, err := f.SearchFacility(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCaseFacilityImpl_SyncFacilities(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

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

func TestUseCaseFacilityImpl_UpdateFacility(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

	UUID := uuid.New().String()

	type args struct {
		ctx                context.Context
		updateFacilityData *dto.UpdateFacilityPayload
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
				updateFacilityData: &dto.UpdateFacilityPayload{
					FacilityID:         UUID,
					FHIROrganisationID: UUID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case",
			args: args{
				ctx: ctx,
				updateFacilityData: &dto.UpdateFacilityPayload{
					FacilityID:         UUID,
					FHIROrganisationID: UUID,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case" {
				fakeDB.MockUpdateFacilityFn = func(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			err := f.UpdateFacility(tt.args.ctx, tt.args.updateFacilityData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
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
			f := facility.NewFacilityUsecase(fakeDB, fakeDB, fakeDB, fakeDB, fakePubsub)

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
