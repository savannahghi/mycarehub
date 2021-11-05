package client_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	mycareHubExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

func InitializeTestService(ctx context.Context) *interactor.Interactor {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil
	}

	onboardingExt := mycareHubExtension.NewOnboardingLibImpl()

	db := postgres.NewMyCareHubDb(pg, pg, pg)

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db)

	// Initialize client usecase
	clientUseCase := client.NewUseCasesClientImpl(db, db, db)

	userUsecase := user.NewUseCasesUserImpl(db, db, db, onboardingExt)

	i := interactor.NewMyCareHubInteractor(facilityUseCase, clientUseCase, userUsecase)
	return i

}

func TestUseCasesClientImpl_RegisterClient_Integration_test(t *testing.T) {
	//Cleanup befor running test
	TearDown(t)

	ctx := context.Background()

	u := InitializeTestService(ctx)

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

		{
			name: "Sad case - nil  input",
			args: args{
				ctx:         ctx,
				userInput:   nil,
				clientInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := u.ClientUseCase.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected client registration output not to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected client registration output not to be nil for %v", tt.name)
				return
			}
		})
	}

	//Cleanup after running tests
	TearDown(t)
}

func TearDown(t *testing.T) {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return
	}

	pg.DB.Migrator().DropTable(&gorm.Contact{})
	pg.DB.Migrator().DropTable(&gorm.PINData{})
	pg.DB.Migrator().DropTable(&gorm.User{})
	pg.DB.Migrator().DropTable(&gorm.ClientProfile{})
	pg.DB.Migrator().DropTable(&gorm.Facility{})
}
