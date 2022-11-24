package organisation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation"
)

func TestUseCaseOrganisationImpl_CreateOrganisation(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.OrganisationInput
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
				input: dto.OrganisationInput{
					OrganisationCode: uuid.New().String(),
					Name:             "name",
					Description:      "description",
					EmailAddress:     "email_address",
					PhoneNumber:      "phone_number",
					PostalAddress:    "postal_address",
					PhysicalAddress:  "physical_address",
					DefaultCountry:   "default_country",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create organisation",
			args: args{
				ctx: context.Background(),
				input: dto.OrganisationInput{
					OrganisationCode: uuid.New().String(),
					Name:             "name",
					Description:      "description",
					EmailAddress:     "email_address",
					PhoneNumber:      "phone_number",
					PostalAddress:    "postal_address",
					PhysicalAddress:  "physical_address",
					DefaultCountry:   "default_country",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()

			if tt.name == "sad case: unable to create organisation" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation) error {
					return fmt.Errorf("unable to create organisation")
				}
			}

			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension)
			_, err := o.CreateOrganisation(tt.args.ctx, tt.args.input)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()

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
			if tt.name == "sad case: unable to get staff profile" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("unable to get staff profile")
				}
			}
			if tt.name == "sad case: unable to check if organisation exists" {
				fakeDB.MockCheckOrganisationExistsFn = func(ctx context.Context, organisationID string) (bool, error) {
					return false, fmt.Errorf("unable to check if the organisation exists")
				}
			}

			o := organisation.NewUseCaseOrganisationImpl(fakeDB, fakeDB, fakeDB, fakeExtension)
			_, err := o.DeleteOrganisation(tt.args.ctx, tt.args.organisationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOrganisationImpl.DeleteOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
