package organisation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation"
)

func TestUseCaseOrganisationImpl_CreateOrganisation(t *testing.T) {
	type args struct {
		ctx               context.Context
		organisationInput dto.OrganisationInput
		programInput      []*dto.ProgramInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create organisation",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: create organisation, no program",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create organisation",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to publish to pubsub",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to publish program",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: organisation name exists",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: existing email address",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: organisation phone number exists",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: organisation code exists",
			args: args{
				ctx: context.Background(),
				organisationInput: dto.OrganisationInput{
					Code:            uuid.New().String(),
					Name:            "name",
					Description:     "description",
					EmailAddress:    "email_address",
					PhoneNumber:     "phone_number",
					PostalAddress:   "postal_address",
					PhysicalAddress: "physical_address",
					DefaultCountry:  "default_country",
				},
				programInput: []*dto.ProgramInput{
					{
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakePubsub)

			if tt.name == "sad case: unable to create organisation" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
					return nil, fmt.Errorf("unable to create organisation")
				}
			}
			if tt.name == "sad case: unable to publish to pubsub" {
				fakePubsub.MockNotifyCreateCMSOrganisationFn = func(ctx context.Context, program *dto.CreateCMSOrganisationPayload) error {
					return fmt.Errorf("unable to publish to pubsub")
				}
			}
			if tt.name == "sad case: failed to publish program" {
				fakePubsub.MockNotifyCreateCMSProgramFn = func(ctx context.Context, program *dto.CreateCMSProgramPayload) error {
					return fmt.Errorf("unable to publish to pubsub")
				}
			}
			if tt.name == "sad case: organisation name exists" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
					return nil, fmt.Errorf("common_organisation_name_key")
				}
			}
			if tt.name == "sad case: existing email address" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
					return nil, fmt.Errorf("common_organisation_email_address_key")
				}
			}
			if tt.name == "sad case: organisation phone number exists" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
					return nil, fmt.Errorf("common_organisation_phone_number_key")
				}
			}
			if tt.name == "sad case: organisation code exists" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
					return nil, fmt.Errorf("common_organisation_org_code_key")
				}
			}

			_, err := o.CreateOrganisation(tt.args.ctx, tt.args.organisationInput, tt.args.programInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseOrganisationImpl_DeleteOrganisation(t *testing.T) {
	type args struct {
		ctx            context.Context
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete organisation",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to delete organisation",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get logged in user",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get staff profile",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to check if organisation exists",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to user profile by logged in user id",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakePubsub)

			if tt.name == "sad case: unable to delete organisation" {
				fakeDB.MockDeleteOrganisationFn = func(ctx context.Context, organisation *domain.Organisation) error {
					return fmt.Errorf("unable to delete organisation")
				}
			}
			if tt.name == "sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}
			if tt.name == "Sad case - fail to user profile by logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("unable to get staff profile")
				}
			}
			if tt.name == "sad case: unable to check if organisation exists" {
				fakeDB.MockCheckOrganisationExistsFn = func(ctx context.Context, organisationID string) (bool, error) {
					return false, fmt.Errorf("unable to check if the organisation exists")
				}
			}

			_, err := o.DeleteOrganisation(tt.args.ctx, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOrganisationImpl.DeleteOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCaseOrganisationImpl_ListOrganisations(t *testing.T) {
	type args struct {
		ctx             context.Context
		paginationInput *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list organisations",
			args: args{
				ctx: context.Background(),
				paginationInput: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to list organisations",
			args: args{
				ctx: context.Background(),
				paginationInput: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 2,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakePubsub)

			if tt.name == "sad case: unable to list organisations" {
				fakeDB.MockListOrganisationsFn = func(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("unable to list organisations")
				}
			}
			_, err := o.ListOrganisations(tt.args.ctx, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOrganisationImpl.ListOrganisations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseOrganisationImpl_SearchOrganisation(t *testing.T) {
	type fields struct {
		Create      infrastructure.Create
		Delete      infrastructure.Delete
		Query       infrastructure.Query
		ExternalExt extension.ExternalMethodsExtension
		Pubsub      pubsubmessaging.ServicePubsub
	}
	type args struct {
		ctx             context.Context
		searchParameter string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case: search organisation",
			args: args{
				ctx:             context.Background(),
				searchParameter: "test",
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to search organisation",
			args: args{
				ctx:             context.Background(),
				searchParameter: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakePubsub)

			if tt.name == "sad case: unable to search organisation" {
				fakeDB.MockSearchOrganisationsFn = func(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
					return nil, fmt.Errorf("unable to search organisation")
				}
			}
			_, err := o.SearchOrganisation(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOrganisationImpl.SearchOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseOrganisationImpl_GetOrganisationByID(t *testing.T) {
	type args struct {
		ctx            context.Context
		organisationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get organisation by id",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get organisation by id",
			args: args{
				ctx:            context.Background(),
				organisationID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakePubsub)

			if tt.name == "sad case: unable to get organisation by id" {
				fakeDB.MockGetOrganisationFn = func(ctx context.Context, id string) (*domain.Organisation, error) {
					return nil, fmt.Errorf("unable to get organisation by id")
				}
			}
			_, err := o.GetOrganisationByID(tt.args.ctx, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOrganisationImpl.GetOrganisationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
