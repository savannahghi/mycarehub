package user_test

import (
	"context"
	"fmt"
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
	"github.com/tj/assert"
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

func TestUseCasesUserImpl_SetUserPIN_Integration_test(t *testing.T) {
	// clean up before running any test
	TearDown(t)

	ctx := context.Background()

	phone := gofakeit.Phone()
	PIN := "1234"
	ConfirmedPin := "1234"
	flavour := feedlib.FlavourConsumer

	validPINInput := createPINInputData(phone, PIN, ConfirmedPin, flavour)

	invalidPINLength := createPINInputData(phone, "123", "123", flavour)

	pinNotDigits := createPINInputData(phone, "page", ConfirmedPin, flavour)

	u := InitializeTestService(ctx)

	AddressesInput := &dto.AddressesInput{
		Type:       enums.AddressesTypePostal,
		Text:       gofakeit.BeerAlcohol(),
		Country:    enums.CountryTypeKenya,
		PostalCode: "10400",
		County:     enums.CountyTypeNairobi,
		Active:     true,
	}

	userInput := &dto.UserInput{
		Username:    gofakeit.Username(),
		DisplayName: gofakeit.BeerHop(),
		FirstName:   gofakeit.FirstName(),
		MiddleName:  gofakeit.BeerAlcohol(),
		LastName:    gofakeit.LastName(),
		UserType:    enums.ClientUser,
		Gender:      enumutils.GenderMale,
		Contacts:    []*dto.ContactInput{{Type: enums.PhoneContact, Contact: phone, Active: true, OptedIn: true}},
		Languages:   enumutils.AllLanguage,
		Flavour:     feedlib.FlavourConsumer,
		Address:     []*dto.AddressesInput{AddressesInput},
	}
	clientProfileInput := &dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}

	testClient, err := u.ClientUseCase.RegisterClient(ctx, userInput, clientProfileInput)
	if err != nil {
		fmt.Printf("unable to register client: %v\n", err)
		return
	}
	assert.NotNil(t, testClient)

	ext := mycareHubExtension.NewOnboardingLibImpl()
	salt, encodedPIN := ext.EncryptPIN(PIN, nil)

	isMatch := ext.ComparePIN(PIN, salt, encodedPIN, nil)
	if !isMatch {
		fmt.Printf("pin mismatch")
		return
	}

	type args struct {
		ctx   context.Context
		input *dto.PinInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: validPINInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid PIN length",
			args: args{
				ctx:   ctx,
				input: invalidPINLength,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin not digits",
			args: args{
				ctx:   ctx,
				input: pinNotDigits,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := u.UserUsecase.SetUserPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}

	//clean up
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
