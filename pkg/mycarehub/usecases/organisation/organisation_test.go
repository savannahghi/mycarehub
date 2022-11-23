package organisation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
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

			if tt.name == "sad case: unable to create organisation" {
				fakeDB.MockCreateOrganisationFn = func(ctx context.Context, organisation *domain.Organisation) error {
					return fmt.Errorf("unable to create organisation")
				}
			}

			o := organisation.NewUseCaseOrganisationImpl(fakeDB)
			_, err := o.CreateOrganisation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrganisation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
