package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/chargemaster"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"

	crmExt "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/crm"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	erp "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

const (
	testHTTPClientTimeout = 180
)

const (
	testChargeMasterBranchID = "94294577-6b27-4091-9802-1ce0f2ce4153"
	engagementService        = "engagement"
	ediService               = "edi"
)

/// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func initializeAcceptanceTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}
	return fsc, fbc
}

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	var repo repository.OnboardingRepository

	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		fsc, fbc := initializeAcceptanceTestFirebaseClient(ctx)
		firestoreExtension := fb.NewFirestoreClientExtension(fsc)
		repo = fb.NewFirebaseRepository(firestoreExtension, fbc)
	}

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}
	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	ext := extension.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})

	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, ext)
	ediClient := utils.NewInterServiceClient(ediService, ext)

	erp := erp.NewAccounting()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(context.Background(), hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr)
	crmExt := crmExt.NewCrmService(hubspotUsecases)
	engage := engagement.NewServiceEngagementImpl(engagementClient, ext)
	edi := edi.NewEdiService(ediClient, repo)
	ps, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		ext,
		erp,
		crmExt,
		edi,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}
	mes := messaging.NewServiceMessagingImpl(ext)
	pinExt := extension.NewPINExtensionImpl()
	profile := usecases.NewProfileUseCase(repo, ext, engage, ps, crmExt)

	supplier := usecases.NewSupplierUseCases(repo, profile, erp, chrg, engage, mes, ext, ps)
	login := usecases.NewLoginUseCases(repo, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(repo, ext)
	userpin := usecases.NewUserPinUseCase(repo, profile, ext, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, ext, engage, ps, edi)
	nhif := usecases.NewNHIFUseCases(repo, profile, ext, engage)
	sms := usecases.NewSMSUsecase(repo, ext)

	return &interactor.Interactor{
		Onboarding:   profile,
		Signup:       su,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
		NHIF:         nhif,
		PubSub:       ps,
		SMS:          sms,
	}, nil
}

func composeInValidUserPayload(t *testing.T) *dto.SignUpInput {
	phone := interserviceclient.TestUserPhoneNumber
	pin := "" // empty string
	flavour := feedlib.FlavourPro
	payload := &dto.SignUpInput{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeValidUserPayload(t *testing.T, phone string) (*dto.SignUpInput, error) {
	pin := "2030"
	flavour := feedlib.FlavourPro
	otp, err := generateTestOTP(t, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to generate test OTP: %v", err)
	}
	return &dto.SignUpInput{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
		OTP:         &otp.OTP,
	}, nil
}

func composeSMSMessageDataPayload(t *testing.T, payload *dto.AfricasTalkingMessage) *strings.Reader {
	data := url.Values{}
	data.Set("date", payload.Date)
	data.Set("from", payload.From)
	data.Set("id", payload.ID)
	data.Set("linkId", payload.LinkID)
	data.Set("text", payload.Text)
	data.Set("to", payload.To)

	smspayload := strings.NewReader(data.Encode())
	return smspayload
}

func composeUSSDPayload(t *testing.T, payload *dto.SessionDetails) *strings.Reader {
	data := url.Values{}
	data.Set("sessionId", payload.SessionID)
	data.Set("phoneNumber", *payload.PhoneNumber)
	data.Set("text", payload.Text)

	smspayload := strings.NewReader(data.Encode())
	return smspayload
}

func CreateTestUserByPhone(t *testing.T, phone string) (*profileutils.UserResponse, error) {
	client := http.DefaultClient
	validPayload, err := composeValidUserPayload(t, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to compose a valid payload: %v", err)
	}
	bs, err := json.Marshal(validPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)
	url := fmt.Sprintf("%s/create_user_by_phone", baseURL)
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

	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %v", err)

	}
	// if resp.StatusCode != http.StatusCreated {
	// 	return nil, fmt.Errorf("failed to create user: expected status to be %v got %v ", http.StatusCreated, resp.StatusCode)
	// }
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP error: %v", err)
	}

	var userResponse profileutils.UserResponse
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall response: %v", err)
	}
	return &userResponse, nil
}

func TestCreateTestUserByPhone(t *testing.T) {
	userResponse, err := CreateTestUserByPhone(t, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("failed to create test user")
		return
	}
	if userResponse == nil {
		t.Errorf("got a nil user response")
		return
	}
}

func RemoveTestUserByPhone(t *testing.T, phone string) (bool, error) {
	client := http.DefaultClient
	validPayload := &dto.PhoneNumberPayload{PhoneNumber: &phone}
	bs, err := json.Marshal(validPayload)
	if err != nil {
		return false, fmt.Errorf("unable to marshal test item to JSON: %s", err)
	}
	payload := bytes.NewBuffer(bs)
	url := fmt.Sprintf("%s/remove_user", baseURL)
	r, err := http.NewRequest(
		http.MethodPost,
		url,
		payload,
	)

	if err != nil {
		return false, fmt.Errorf("can't create new request: %v", err)

	}

	if r == nil {
		return false, fmt.Errorf("nil request")
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return false, fmt.Errorf("HTTP error: %v", err)

	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to remove user: expected status to be %v got %v ", http.StatusCreated, resp.StatusCode)
	}
	return true, nil
}

func TestRemoveTestUserByPhone(t *testing.T) {
	phone := interserviceclient.TestUserPhoneNumber
	userResponse, err := CreateTestUserByPhone(t, phone)
	if err != nil {
		t.Errorf("failed to create test user")
		return
	}
	if userResponse == nil {
		t.Errorf("got a nil user response")
		return
	}

	removed, err := RemoveTestUserByPhone(t, phone)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}
	if !removed {
		t.Errorf("user was not removed")
		return
	}
}

func generateTestOTP(t *testing.T, phone string) (*profileutils.OtpResponse, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	testAppID := uuid.New().String()
	return s.Engagement.GenerateAndSendOTP(ctx, phone, &testAppID)
}

func setPrimaryEmailAddress(ctx context.Context, t *testing.T, emailAddress string) error {
	s, err := InitializeTestService(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize test service: %v", err)
	}

	return s.Onboarding.UpdatePrimaryEmailAddress(ctx, emailAddress)
}

func updateBioData(ctx context.Context, t *testing.T, data profileutils.BioData) error {
	s, err := InitializeTestService(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize test service: %v", err)
	}

	return s.Onboarding.UpdateBioData(ctx, data)
}

func addPartnerType(ctx context.Context, t *testing.T, name *string, partnerType profileutils.PartnerType) (bool, error) {
	s, err := InitializeTestService(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to initialize test service: %v", err)
	}

	return s.Supplier.AddPartnerType(ctx, name, &partnerType)
}

func setUpSupplier(ctx context.Context, t *testing.T, accountType profileutils.AccountType) (*profileutils.Supplier, error) {
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	return s.Supplier.SetUpSupplier(ctx, accountType)
}

func setUpLoggedInTestUserGraphHeaders(t *testing.T) map[string]string {
	// create a user and their profile
	phoneNumber := interserviceclient.TestUserPhoneNumber
	resp, err := CreateTestUserByPhone(t, phoneNumber)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return nil
	}

	if resp.Profile.ID == "" {
		t.Errorf(" user profile id should not be empty")
		return nil
	}

	if len(resp.Profile.VerifiedUIDS) == 0 {
		t.Errorf(" user profile VerifiedUIDS should not be empty")
		return nil
	}

	logrus.Infof("profile from create user : %v", resp.Profile)

	logrus.Infof("uid from create user : %v", resp.Auth.UID)

	return getGraphHeaders(*resp.Auth.IDToken)
}

// func setRoleForUserWithPhone(phoneNumber string, role profileutils.RoleType, headers map[string]string) error {
// 	url := fmt.Sprintf("%s/roles/add_user_role", baseURL)

// 	payload := dto.RolePayload{
// 		PhoneNumber: &phoneNumber,
// 		Role:        &role,
// 	}
// 	body, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("unable to marshal payload to JSON: %s", err)
// 	}

// 	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return fmt.Errorf("unable to compose request: %s", err)
// 	}
// 	if request == nil {
// 		return fmt.Errorf("nil request")
// 	}

// 	for header, value := range headers {
// 		request.Header.Add(header, value)
// 	}

// 	response, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		return fmt.Errorf("request error: %s", err)
// 	}

// 	if response.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to set role for user : expected status to be %v, but got %v ", http.StatusOK, response.StatusCode)
// 	}
// 	return nil
// }

func getGraphHeaders(idToken string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", idToken),
	}
}

func TestMain(m *testing.M) {
	// setup
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "onboarding_testing")
	os.Setenv("SAVANNAH_ADMIN_EMAIL", "test@bewell.co.ke")
	os.Setenv("REPOSITORY", "firebase")

	ctx := context.Background()
	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	fsc, _ := initializeAcceptanceTestFirebaseClient(ctx)

	purgeRecords := func() {
		if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
			r := fb.Repository{}
			collections := []string{
				r.GetCustomerProfileCollectionName(),
				r.GetPINsCollectionName(),
				r.GetUserProfileCollectionName(),
				r.GetSupplierProfileCollectionName(),
				r.GetSurveyCollectionName(),
				r.GetCommunicationsSettingsCollectionName(),
				r.GetCustomerProfileCollectionName(),
				r.GetExperimentParticipantCollectionName(),
				r.GetKCYProcessCollectionName(),
				r.GetMarketingDataCollectionName(),
				r.GetNHIFDetailsCollectionName(),
				r.GetProfileNudgesCollectionName(),
				r.GetSMSCollectionName(),
				r.GetUSSDDataCollectionName(),
			}
			for _, collection := range collections {
				ref := fsc.Collection(collection)
				firebasetools.DeleteCollection(ctx, fsc, ref, 10)
			}
		}

	}

	purgeRecords()

	// run the tests
	log.Printf("about to run tests")
	code := m.Run()
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

func TestRouter(t *testing.T) {
	ctx := context.Background()
	router, err := presentation.Router(ctx)
	if err != nil {
		t.Errorf("can't initialize router: %v", err)
		return
	}

	if router == nil {
		t.Errorf("nil router")
		return
	}
}

func TestHealthStatusCheck(t *testing.T) {
	client := http.DefaultClient

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful health check",
			args: args{
				url: fmt.Sprintf(
					"%s/health",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range interserviceclient.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}
