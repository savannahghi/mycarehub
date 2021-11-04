package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	onboardingExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/serverutils"
)

const (
	engagementService     = "engagement"
	testHTTPClientTimeout = 180
)

var (
	srv       *http.Server
	baseURL   string
	serverErr error
)

var (
	testInteractor     *interactor.Interactor
	testIntractorError error

	testClient      *domain.ClientUserProfile
	testClientError error

	hasSetClientPin      bool
	hasSetClientPinError error
)

var (
	testPhone         = interserviceclient.TestUserPhoneNumber
	testPin           = "1234"
	testClientFlavour = feedlib.FlavourConsumer

	userInput = &dto.UserInput{
		FirstName:   gofakeit.FirstName(),
		MiddleName:  gofakeit.BeerAlcohol(),
		LastName:    gofakeit.LastName(),
		Username:    gofakeit.Username(),
		DisplayName: gofakeit.BeerHop(),
		UserType:    enums.ClientUser,
		Gender:      enumutils.GenderMale,
		Contacts: []*dto.ContactInput{
			{
				Type:    enums.PhoneContact,
				Contact: testPhone,
				Active:  true,
				OptedIn: true,
			},
		},
		Languages: enumutils.AllLanguage,
		Flavour:   testClientFlavour,
	}
	clientInput = &dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}

	pinInput = &dto.PinInput{
		PhoneNumber:  testPhone,
		PIN:          testPin,
		ConfirmedPin: testPin,
		Flavour:      testClientFlavour,
	}
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")

	initialEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "testing")

	ctx := context.Background()

	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	)
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
		os.Exit(1)
	}

	pg, pgError := gorm.NewPGInstance()
	if pgError != nil {
		log.Printf("can't instantiate test repository: %v", pgError)
		os.Exit(1)
	}

	testInteractor, testIntractorError = initializeTestService(ctx)
	if testIntractorError != nil {
		log.Printf("Error initializing test service: %v", testIntractorError)
		os.Exit(1)
	}

	testClient, testClientError = registerTestClient(ctx, userInput, clientInput)
	if testClientError != nil {
		log.Printf("Error registering test client: %v", testClientError)
		os.Exit(1)
	}

	hasSetClientPin, hasSetClientPinError = setTestClientPin(ctx, pinInput)
	if hasSetClientPinError != nil {
		log.Printf("Error setting test client pin: %v", hasSetClientPinError)
		os.Exit(1)
	}
	if !hasSetClientPin {
		log.Printf("Error setting test client pin")
		os.Exit(1)
	}

	// run tests
	log.Printf("Running tests ...")
	code := m.Run()

	// teardown
	// pgError = pg.DB.Unscoped().Where("user_id", testClient.User.ID).Delete(&gorm.PINData{}).Error
	// if pgError != nil {
	// 	log.Printf("Error deleting test client pin data: %v", pgError)
	// 	os.Exit(1)
	// }
	// for _, c := range testClient.User.Contacts {
	// 	err := pg.DB.Unscoped().Where("contact_id", c.ID).Delete(&gorm.Contact{}).Error
	// 	if err != nil {
	// 		log.Printf("Error deleting user contact: %v", err)
	// 		os.Exit(1)
	// 	}
	// }
	// pgError = pg.DB.Unscoped().Where("id", testClient.Client.ID).Delete(&gorm.ClientProfile{}).Error
	// if pgError != nil {
	// 	log.Printf("Error deleting test client: %v", pgError)
	// 	os.Exit(1)
	// }
	// pgError = pg.DB.Unscoped().Where("user_id", testClient.User.ID).Delete(&gorm.User{}).Error
	// if pgError != nil {
	// 	log.Printf("Error deleting test client user: %v", pgError)
	// 	os.Exit(1)
	// }
	pg.DB.Migrator().DropTable(&gorm.Contact{})
	pg.DB.Migrator().DropTable(&gorm.PINData{})
	pg.DB.Migrator().DropTable(&gorm.ClientProfile{})
	pg.DB.Migrator().DropTable(&gorm.User{})
	pg.DB.Migrator().DropTable(&gorm.Facility{})
	// restore envs
	os.Setenv(initialEnv, "ENVIRONMENT")

	log.Printf("finished running tests")

	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()

	os.Exit(code)
}

func initializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fc := &firebasetools.FirebaseClient{}
	_, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate test repository: %v", err)
	}

	onboardingExt := onboardingExtension.NewOnboardingLibImpl()

	db := postgres.NewMyCareHubDb(pg, pg, pg)

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db)

	// Initialize client usecase
	clientUseCase := client.NewUseCasesClientImpl(db, db, db)

	userUsecase := user.NewUseCasesUserImpl(db, db, db, onboardingExt)

	i := interactor.NewMyCareHubInteractor(facilityUseCase, clientUseCase, userUsecase)
	return i, nil
}

func registerTestClient(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
	i, err := initializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test service: %v", err)
	}

	client, err := i.ClientUseCase.RegisterClient(ctx, userInput, clientInput)
	if err != nil {
		return nil, fmt.Errorf("failed to register test client: %v", err)
	}

	return client, nil
}

func setTestClientPin(ctx context.Context, pinInput *dto.PinInput) (bool, error) {
	i, err := initializeTestService(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to initialize test service: %v", err)
	}

	ok, err := i.UserUsecase.SetUserPIN(ctx, pinInput)
	if err != nil {
		return false, fmt.Errorf("failed to set test user pin: %v", err)
	}
	if !ok {
		return false, fmt.Errorf("failed to set test user pin: %v", err)
	}

	return ok, nil

}

func composeValidClientUserPayload(t *testing.T, phoneNumber string) (*dto.LoginInput, error) {
	pin := testPin
	flavour := testClientFlavour

	return &dto.LoginInput{
		PhoneNumber: &phoneNumber,
		PIN:         &pin,
		Flavour:     flavour,
	}, nil
}

func loginByPhone(t *testing.T, phoneNumber string) (*domain.AuthCredentials, error) {
	validPayload, err := composeValidClientUserPayload(t, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to compose a valid payload: %v", err)
	}

	bs, err := json.Marshal(validPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)

	url := fmt.Sprintf("%s/login_by_phone", baseURL)

	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)

	}
	if r == nil {
		return nil, fmt.Errorf("nil request")
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)

	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP error: %v", err)
	}

	var loginResponse domain.AuthCredentials
	err = json.Unmarshal(data, &loginResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall response: %v", err)
	}
	return &loginResponse, nil
}
