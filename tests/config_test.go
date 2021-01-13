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
	"os"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation"

	"gitlab.slade360emr.com/go/base"
)

const (
	testHTTPClientTimeout = 180
)

const (
	testSladeCode            = "BRA-PRO-4190-4"
	testEDIPortalUsername    = "avenue-4190@healthcloud.co.ke"
	testEDIPortalPassword    = "test provider"
	testChargeMasterBranchID = "94294577-6b27-4091-9802-1ce0f2ce4153"
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
	fc := base.FirebaseClient{}
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
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, err
	}

	profile := usecases.NewProfileUseCase(fr)
	otp := otp.NewOTPService(fr)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes)
	login := usecases.NewLoginUseCases(fr, profile)
	survey := usecases.NewSurveyUseCases(fr)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier, otp)

	return &interactor.Interactor{
		Onboarding:   profile,
		Signup:       su,
		Otp:          otp,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
	}, nil
}

func composeInValidUserPayload(t *testing.T) *resources.SignUpInput {
	phone := base.TestUserPhoneNumber
	pin := "" // empty string
	flavour := base.FlavourPro
	payload := &resources.SignUpInput{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
	}
	return payload
}

func composeValidUserPayload(t *testing.T, phone string) (*resources.SignUpInput, error) {
	pin := "2030"
	flavour := base.FlavourPro
	otp, err := generateTestOTP(t, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to generate test OTP: %v", err)
	}
	return &resources.SignUpInput{
		PhoneNumber: &phone,
		PIN:         &pin,
		Flavour:     flavour,
		OTP:         &otp.OTP,
	}, nil
}

func CreateTestUserByPhone(t *testing.T, phone string) (*base.UserResponse, error) {
	client := http.DefaultClient
	validPayload, err := composeValidUserPayload(t, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to compose a valid payload")
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
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create user: expected status to be %v got %v ", http.StatusCreated, resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP error: %v", err)

	}

	var userResponse base.UserResponse
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall response: %v", err)
	}
	return &userResponse, nil
}

func RemoveTestUserByPhone(t *testing.T, phone string) (bool, error) {
	client := http.DefaultClient
	validPayload := &resources.PhoneNumberPayload{PhoneNumber: &phone}
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

func generateTestOTP(t *testing.T, phone string) (*base.OtpResponse, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	return s.Otp.GenerateAndSendOTP(ctx, phone)
}

func setUpLoggedInTestUserGraphHeaders(t *testing.T) map[string]string {
	// create a user and thier profile
	phoneNumber := base.TestUserPhoneNumber
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

	ctx := context.Background()
	srv, baseURL, serverErr = base.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	r := database.Repository{} // They are nil
	fsc, _ := initializeAcceptanceTestFirebaseClient(ctx)

	purgeRecords := func() {
		collections := []string{
			r.GetCustomerProfileCollectionName(),
			r.GetPINsCollectionName(),
			r.GetUserProfileCollectionName(),
			r.GetSupplierProfileCollectionName(),
			r.GetSurveyCollectionName(),
		}
		for _, collection := range collections {
			ref := fsc.Collection(collection)
			base.DeleteCollection(ctx, fsc, ref, 10)
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

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
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
