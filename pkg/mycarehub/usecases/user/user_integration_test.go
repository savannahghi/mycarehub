package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	mycareHubExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/tj/assert"
)

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

func TestUseCasesUserImpl_Login_Integration_test(t *testing.T) {
	//clean up before running test
	TearDown(t)

	ctx := context.Background()

	u := InitializeTestService(ctx)
	ext := mycareHubExtension.NewOnboardingLibImpl()

	phoneNumber := interserviceclient.TestUserPhoneNumber

	PIN := "1234"
	ConfirmPIN := "1234"
	flavour := feedlib.FlavourConsumer

	addresses := &dto.AddressesInput{
		Type:       enums.AddressesTypePostal,
		Text:       gofakeit.BeerAlcohol(),
		Country:    enums.CountryTypeKenya,
		PostalCode: gofakeit.BeerName(),
		County:     enums.CountyTypeNairobi,
		Active:     true,
	}

	contact := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: phoneNumber,
		Active:  true,
		OptedIn: true,
	}

	userInput := &dto.UserInput{
		Username:    gofakeit.Username(),
		DisplayName: gofakeit.BeerAlcohol(),
		FirstName:   gofakeit.FirstName(),
		MiddleName:  gofakeit.BeerAlcohol(),
		LastName:    gofakeit.LastName(),
		UserType:    enums.ClientUser,
		Gender:      enumutils.GenderFemale,
		Contacts:    []*dto.ContactInput{contact},
		Languages:   []enumutils.Language{enumutils.LanguageSw},
		Flavour:     feedlib.FlavourConsumer,
		Address:     []*dto.AddressesInput{addresses},
	}

	clientInput := &dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}

	clientProfile, err := u.ClientUseCase.RegisterClient(ctx, userInput, clientInput)
	if err != nil {
		fmt.Printf("an error occurred while registering client: %v", err)
		return
	}

	salt, encryptedPIN := ext.EncryptPIN(PIN, nil)

	pinData := &domain.UserPIN{
		UserID:    *clientProfile.User.ID,
		HashedPIN: encryptedPIN,
		ValidFrom: time.Now(),
		ValidTo:   time.Now().AddDate(0, 0, 7),
		Flavour:   clientProfile.User.Flavour,
		IsValid:   true,
		Salt:      salt,
	}

	validPINInput := createPINInputData(phoneNumber, PIN, ConfirmPIN, flavour)
	_, err = u.UserUsecase.SetUserPIN(ctx, validPINInput)
	if err != nil {
		fmt.Printf("an error occurred while saving PIN: %v", err)
		return
	}

	isMatch := ext.ComparePIN(ConfirmPIN, salt, encryptedPIN, nil)
	if !isMatch {
		fmt.Printf("the provided PINs do not match")
		return
	}

	currentTime := time.Now()
	isExpired := utils.CheckPINExpiry(currentTime, pinData)
	if !isExpired {
		fmt.Printf("the provided PINs has expired")
		return
	}

	customToken, customTokenErr := firebasetools.CreateFirebaseCustomToken(ctx, *clientProfile.User.ID)
	if customTokenErr != nil {
		fmt.Printf("an error occurred while creating firebase custom token: %v", customTokenErr)
		return
	}

	_, userTokensErr := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if userTokensErr != nil {
		fmt.Printf("an error occurred while authenticating firebase custom token: %v", customTokenErr)
		return
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
		pin         string
		flavour     feedlib.Flavour
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
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         "PIN",
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         "12345",
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin expired",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case - pin expired" {
				pg, err := gorm.NewPGInstance()
				if err != nil {
					return
				}

				pastDate := time.Now().AddDate(0, 0, -50)
				err = pg.DB.Model(&gorm.PINData{}).Where(&gorm.PINData{UserID: *clientProfile.User.ID}).Updates(&gorm.PINData{ValidTo: pastDate}).Error
				if err != nil {
					t.Errorf("an error occurred while getting user by ID: %v", err)
					return
				}
			}

			got, _, err := u.UserUsecase.Login(tt.args.ctx, tt.args.phoneNumber, tt.args.pin, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.Login() error = %v, wantErr %v", err, tt.wantErr)
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

	//clean up after running test
	TearDown(t)
}

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
