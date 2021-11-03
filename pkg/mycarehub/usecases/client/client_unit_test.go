package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client/mock"
)

func TestUseCasesClientImpl_RegisterClient_Unittest(t *testing.T) {
	ctx := context.Background()

	addresses := &dto.AddressesInput{
		Type:       enums.AddressesTypePostal,
		Text:       gofakeit.BeerAlcohol(),
		Country:    enums.CountryTypeKenya,
		PostalCode: gofakeit.BeerName(),
		County:     enums.CountyTypeNairobi,
		Active:     true,
	}

	userInput := &dto.UserInput{
		Username:    gofakeit.Username(),
		DisplayName: gofakeit.BeerAlcohol(),
		FirstName:   gofakeit.FirstName(),
		MiddleName:  gofakeit.BeerAlcohol(),
		LastName:    gofakeit.LastName(),
		UserType:    enums.ClientUser,
		Gender:      enumutils.GenderFemale,
		Contacts:    []*dto.ContactInput{},
		Languages:   []enumutils.Language{enumutils.LanguageSw},
		Flavour:     feedlib.FlavourConsumer,
		Address:     []*dto.AddressesInput{addresses},
	}

	clientInput := &dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}

	type args struct {
		ctx         context.Context
		userInput   *dto.UserInput
		clientInput *dto.ClientProfileInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: clientInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case - nil user input",
			args: args{
				ctx:         ctx,
				userInput:   nil,
				clientInput: clientInput,
			},
			wantErr: true,
		},

		{
			name: "Sad case - nil client input",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeClient := mock.NewClientUsecaseMock()

			c := NewUseCasesClientImpl(fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case - nil user input" {
				fakeClient.MockRegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register client")
				}
			}

			if tt.name == "Sad case - nil client input" {
				fakeClient.MockRegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register client")
				}
			}

			if tt.name == "Sad case - nil context" {
				fakeClient.MockRegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register client")
				}
			}

			got, err := c.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected client registration output not to be nil for %v", tt.name)
				return
			}
		})
	}
}
