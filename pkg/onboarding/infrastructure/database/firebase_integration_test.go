package database_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"

	"fmt"

	"os"
	"reflect"

	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

const otpService = "otp"

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	envOriginalValue := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")
	debugEnvValue := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	collectionEnvValue := os.Getenv("ROOT_COLLECTION_SUFFIX")
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("onboarding_ci_%v", time.Now().Unix()))
	ctx := context.Background()
	r := database.Repository{} // They are nil
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}

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

	// try clean up first
	purgeRecords()

	// do clean up
	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	// restore environment varibles to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	ext := extension.NewBaseExtensionImpl()
	pinExt := extension.NewPINExtensionImpl()
	otpClient := utils.NewInterServiceClient(otpService)
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)
	otp := otp.NewOTPService(otpClient, ext)
	profile := usecases.NewProfileUseCase(fr, otp, ext)
	erp := erp.NewERPService()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	engage := engagement.NewServiceEngagementImpl()
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes, ext)
	login := usecases.NewLoginUseCases(fr, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(fr, ext)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile, ext, pinExt)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier, otp, ext)

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

func generateTestOTP(t *testing.T, phone string) (*base.OtpResponse, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	return s.Otp.GenerateAndSendOTP(ctx, phone)
}

// CreateTestUserByPhone creates a user that is to be used in
// running of our test cases.
// If the test user already exists then they are logged in
// to get their auth credentials
func CreateOrLoginTestUserByPhone(t *testing.T) (*auth.Token, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service")
	}
	phone := base.TestUserPhoneNumber
	flavour := base.FlavourConsumer
	pin := base.TestUserPin
	exists, err := s.Onboarding.CheckPhoneExists(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}
	if !exists {
		otp, err := generateTestOTP(t, phone)
		log.Println("The otp is:", otp)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test OTP: %v", err)
		}

		u, err := s.Signup.CreateUserByPhone(
			ctx,
			&resources.SignUpInput{
				PhoneNumber: &phone,
				PIN:         &pin,
				Flavour:     flavour,
				OTP:         &otp.OTP,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create a test user: %v", err)
		}
		if u == nil {
			return nil, fmt.Errorf("nil test user response")
		}
		authCred := &auth.Token{
			UID: u.Auth.UID,
		} // We add the test user UID to the expected auth.Token
		return authCred, nil
	}
	logInCreds, err := s.Login.LoginByPhone(
		ctx,
		phone,
		base.TestUserPin,
		flavour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to log in test user: %v", err)
	}
	authCred := &auth.Token{
		UID: logInCreds.Auth.UID,
	}
	return authCred, nil

}

// TestAuthenticatedContext returns a logged in context, useful for test purposes
func GetTestAuthenticatedContext(t *testing.T) (context.Context, *auth.Token, error) {
	ctx := context.Background()
	auth, err := CreateOrLoginTestUserByPhone(t)
	if err != nil {
		return nil, nil, err
	}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		auth,
	)
	return authenticatedContext, auth, nil
}

func InitializeTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up tests: %s", err)
	}
	return fsc, fbc
}

func TestRemoveKYCProcessingRequest(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	assert.Nil(t, err)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), base.TestUserPhoneNumber)

	ctx, auth, err := GetTestAuthenticatedContext(t)
	assert.Nil(t, err)
	assert.NotNil(t, auth)

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	input1 := domain.OrganizationNutrition{
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "someKRAPIN",
		KRAPINUploadID:                     "KRAPINUploadID",
		SupportingDocumentsUploadID:        []string{"SupportingDocumentsUploadID", "Support"},
		CertificateOfIncorporation:         "CertificateOfIncorporation",
		CertificateOfInCorporationUploadID: "CertificateOfInCorporationUploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeMilitary,
				IdentificationDocNumber:         "IdentificationDocNumber",
				IdentificationDocNumberUploadID: "IdentificationDocNumberUploadID",
			},
		},
		OrganizationCertificate: "OrganizationCertificate",
		RegistrationNumber:      "RegistrationNumber",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
	}

	kycJSON, err := json.Marshal(input1)
	assert.Nil(t, err)

	var kycAsMap map[string]interface{}
	err = json.Unmarshal(kycJSON, &kycAsMap)
	assert.Nil(t, err)

	// get the user profile
	profile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, profile)

	// fetch the supplier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, sup)

	//call remove kyc process request. this should fail since the user has not added a kyc yet
	err = fr.RemoveKYCProcessingRequest(ctx, sup.ID)
	assert.NotNil(t, err)

	sup.SupplierKYC = kycAsMap

	// now add the kyc processing request
	req1 := &domain.KYCRequest{
		ID:                  uuid.New().String(),
		ReqPartnerType:      sup.PartnerType,
		ReqOrganizationType: domain.OrganizationType(sup.AccountType),
		ReqRaw:              sup.SupplierKYC,
		Processed:           false,
		SupplierRecord:      sup,
		Status:              domain.KYCProcessStatusPending,
	}
	err = fr.StageKYCProcessingRequest(ctx, req1)
	assert.Nil(t, err)

	// call remove kypc processing request again. this should pass now since there is and existing processing request added
	err = fr.RemoveKYCProcessingRequest(ctx, sup.ID)
	assert.Nil(t, err)

}

func TestPurgeUserByPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	assert.Nil(t, err)
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), base.TestUserPhoneNumber)
	ctx, auth, err := GetTestAuthenticatedContext(t)
	assert.Nil(t, err)
	assert.NotNil(t, auth)
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)
	profile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, base.TestUserPhoneNumber, *profile.PrimaryPhone)

	// fetch the same profile but now using the primary phone number
	profile, err = fr.GetUserProfileByPrimaryPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, base.TestUserPhoneNumber, *profile.PrimaryPhone)

	// purge the record. this should not fail
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)

	// try purging the record again. this should fail since not user profile will be found with the phone number
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.NotNil(t, err)

	// create an invalid user profile
	fakeUID := uuid.New().String()
	invalidpr1, err := fr.CreateUserProfile(context.Background(), base.TestUserPhoneNumber, fakeUID)
	assert.Nil(t, err)
	assert.NotNil(t, invalidpr1)

	// fetch the pins related to invalidpr1. this should fail since no pin has been associated with invalidpr1
	pin, err := fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, pin)

	// fetch the customer profile related to invalidpr1. this should fail since no customer profile has been associated with invalidpr
	cpr, err := fr.GetCustomerProfileByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, cpr)

	// fetch the supplier profile related to invalidpr1. this should fail since no supplier profile has been associated with invalidpr
	spr, err := fr.GetSupplierProfileByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, spr)

	// call PurgeUserByPhoneNumber using the phone number associated with invalidpr1. this should fail since it does not have
	// an associated pin
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server error! unable to perform operation")

	// now set a  pin. this should not fail
	userpin := "1234"
	pset, err := s.UserPIN.SetUserPIN(ctx, userpin, base.TestUserPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, pset)
	assert.Equal(t, true, pset)

	// retrieve the pin and assert it matches the one set
	pin, err = fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pin)
	var pinExt extension.PINExtensionImpl
	matched := pinExt.ComparePIN(userpin, pin.Salt, pin.PINNumber, nil)
	assert.Equal(t, true, matched)

	// now remove. this should pass even though customer/supplier profile don't exist. What must be removed is the pins
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)

	// assert the pin has been removed
	pin, err = fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, pin)

}

func TestCreateEmptyCustomerProfile(t *testing.T) {
	ctx := context.Background()
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	firestoreDB := database.NewFirebaseRepository(firestoreExtension, fbc)

	tests := []struct {
		name      string
		profileID string
		wantErr   bool
	}{
		{
			name:      "valid case",
			profileID: uuid.New().String(),
			wantErr:   false,
		},
		{
			name:      "invalid case",
			profileID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := firestoreDB.CreateEmptyCustomerProfile(ctx, tt.profileID)
			if tt.wantErr && err != nil {
				t.Errorf("error expected but returned no error")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("error was not expected but got error: %v", err)
				return
			}

			if !tt.wantErr && customer == nil {
				t.Errorf("returned a nil customer")
				return
			}
		})
	}

}

func TestGetCustomerProfileByProfileID(t *testing.T) {
	ctx := context.Background()
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	firestoreDB := database.NewFirebaseRepository(firestoreExtension, fbc)
	tests := []struct {
		name      string
		profileID string
		wantErr   bool
	}{
		{
			name:      "valid case",
			profileID: uuid.New().String(),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerTest, err := firestoreDB.CreateEmptyCustomerProfile(ctx, tt.profileID)
			if err != nil {
				t.Errorf("failed to create a test Empty Customer profile err: %v", err)
				return
			}
			if customerTest.ProfileID == nil {
				t.Errorf("nil customer profile ID")
				return
			}
			customerProfile, err := firestoreDB.GetCustomerProfileByProfileID(ctx, tt.profileID)
			if err != nil && !tt.wantErr {
				t.Errorf("error not expected but got error: %v", err)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("error expected but got no error")
				return
			}
			if !tt.wantErr && customerProfile == nil {
				t.Errorf("nil customer profile")
				return
			}

			if !tt.wantErr {
				if customerTest.ProfileID == nil {
					t.Errorf("nil customer profile ID")
					return
				}

				if customerTest.ID == "" {
					t.Errorf("nil customer ID")
					return
				}
			}
		})
	}
}

func TestRepository_GetCustomerOrSupplierProfileByProfileID(t *testing.T) {
	ctx := context.Background()
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()
	_, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	_, err = fr.CreateEmptyCustomerProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty customer: %v", err)
	}
	type args struct {
		ctx       context.Context
		flavour   base.Flavour
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: get the customer profile",
			args: args{
				ctx:       ctx,
				flavour:   base.FlavourConsumer,
				profileID: profileID,
			},
			wantErr: false,
		},
		{
			name: "success: get the supplier profile",
			args: args{
				ctx:       ctx,
				flavour:   base.FlavourPro,
				profileID: profileID,
			},
			wantErr: false,
		},
		{
			name: "failure: bad flavour given",
			args: args{
				ctx:       ctx,
				flavour:   "not-a-flavour-bana",
				profileID: profileID,
			},
			wantErr: true,
		},
		{
			name: "failure: profile ID that does not exist",
			args: args{
				ctx:       ctx,
				flavour:   base.FlavourPro,
				profileID: "not-a-real-profile-ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, supplier, err := fr.GetCustomerOrSupplierProfileByProfileID(
				tt.args.ctx,
				tt.args.flavour,
				tt.args.profileID,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCustomerOrSupplierProfileByProfileID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if base.IsDebug() {
				log.Printf("Customer....%v", customer)
				log.Printf("Supplier....%v", supplier)
			}
		})
	}
}

func TestRepository_GetCustomerProfileByID(t *testing.T) {
	ctx := context.Background()
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()

	customer, err := fr.CreateEmptyCustomerProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty customer: %v", err)
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: get a customer profile by ID",
			args: args{
				ctx: ctx,
				id:  customer.ID,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to get a customer profile",
			args: args{
				ctx: ctx,
				id:  "not-a-customer-ID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerProfile, err := fr.GetCustomerProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetCustomerProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if base.IsDebug() {
				log.Printf("Customer....%v", customerProfile)
			}
		})
	}
}

func TestRepository_ExchangeRefreshTokenForIDToken(t *testing.T) {
	ctx, token, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GenerateAuthCredentials(ctx, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("failed to generate auth credentials: %v", err)
		return
	}

	type args struct {
		refreshToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *auth.Token
		wantErr bool
	}{
		{
			name: "valid firebase refresh token",
			args: args{
				refreshToken: user.RefreshToken,
			},
			want:    token,
			wantErr: false,
		},
		{
			name: "invalid firebase refresh token",
			args: args{
				refreshToken: "",
			},
			wantErr: true,
		},
		{
			name: "invalid firebase refresh token",
			args: args{
				refreshToken: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.ExchangeRefreshTokenForIDToken(tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.ExchangeRefreshTokenForIDToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// obtain auth token details from the id token string
				auth, err := base.ValidateBearerToken(ctx, *got.IDToken)
				if err != nil {
					t.Errorf("invalid token: %w", err)
					return
				}
				if auth.UID != tt.want.UID {
					t.Errorf("Repository.ExchangeRefreshTokenForIDToken() = %v, want %v", got.UID, tt.want.UID)
				}
			}
		})
	}
}

func TestRepository_GetUserProfileByPhoneNumber(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Get a user by valid phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Get a user by an invalid phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetUserProfileByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfileByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("returned a nil user")
				return
			}
		})
	}
}

func TestRepository_GetUserProfileByPrimaryPhoneNumber(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid : primary phone number in context",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "invalid : non-existent wrong phone number format",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254712qwe234",
			},
			wantErr: true,
		},
		{
			name: "invalid : non existent phone number",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254712098765",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetUserProfileByPrimaryPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfileByPrimaryPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("returned a nil user")
				return
			}
		})
	}
}

func TestRepository_GetSupplierProfileByProfileID(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()

	sup, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	type args struct {
		ctx       context.Context
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Supplier
		wantErr bool
	}{
		{
			name: "Happy Case - Get Supplier Profile By Valid profile ID",
			args: args{
				ctx:       ctx,
				profileID: profileID,
			},
			want:    sup,
			wantErr: false,
		},
		{
			name: "Sad Case - Get Supplier Profile By a non-existent profile ID",
			args: args{
				ctx:       ctx,
				profileID: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetSupplierProfileByProfileID(tt.args.ctx, tt.args.profileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetSupplierProfileByProfileID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetSupplierProfileByProfileID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetSupplierProfileByID(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()
	supplier, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Supplier
		wantErr bool
	}{
		{
			name: "Happy Case - Get Supplier by a valid ID",
			args: args{
				ctx: ctx,
				id:  supplier.ID,
			},
			want:    supplier,
			wantErr: false,
		},
		{
			name: "Sad Case - Get Supplier using a non-existent ID",
			args: args{
				ctx: ctx,
				id:  "randomID",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Get Supplier using an empty ID",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetSupplierProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetSupplierProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetSupplierProfileByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetUserProfileByUID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Get user profile by a valid UID",
			args: args{
				ctx: ctx,
				uid: auth.UID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Get user profile by a non-existent UID",
			args: args{
				ctx: context.Background(),
				uid: "random",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Get user profile using an empty UID",
			args: args{
				ctx: ctx,
				uid: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetUserProfileByUID(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfileByUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("returned a nil user")
				return
			}
		})
	}
}

func TestRepository_GetUserProfileByID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "Happy Case - Get user profile using a valid ID",
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    user,
			wantErr: false,
		},
		{
			name: "Sad Case - Get user profile using an invalid ID",
			args: args{
				ctx: ctx,
				id:  "invalid",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Get user profile using an empty ID",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetUserProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetUserProfileByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_CheckIfPhoneNumberExists(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Check for a valid number that does not exist",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254721524371",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Happy Case - Check for a number that exists",
			args: args{
				ctx:         ctx,
				phoneNumber: *user.PrimaryPhone,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.CheckIfPhoneNumberExists(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CheckIfPhoneNumberExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.CheckIfPhoneNumberExists() = %v, want %v", got, tt.want)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected, got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_CheckIfUsernameExists(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Check for a nonexistent username",
			args: args{
				ctx:      ctx,
				userName: "Jatelo Jakom",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Happy Case - Check for an existing username",
			args: args{
				ctx:      ctx,
				userName: *user.UserName,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.CheckIfUsernameExists(tt.args.ctx, tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CheckIfUsernameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.CheckIfUsernameExists() = %v, want %v", got, tt.want)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected, got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_GetPINByProfileID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	pin, err := fr.GetPINByProfileID(ctx, user.ID)
	if err != nil {
		t.Errorf("failed to get pin")
		return
	}

	type args struct {
		ctx       context.Context
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PIN
		wantErr bool
	}{
		{
			name: "Happy Case - Get pin using a valid profileID",
			args: args{
				ctx:       ctx,
				profileID: pin.ProfileID,
			},
			want:    pin,
			wantErr: false,
		},
		{
			name: "Sad Case - Get pin using an invalid profileID",
			args: args{
				ctx:       ctx,
				profileID: "invalidID",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Get pin using an empty profileID",
			args: args{
				ctx:       ctx,
				profileID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.GetPINByProfileID(tt.args.ctx, tt.args.profileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetPINByProfileID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected, got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetPINByProfileID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_SavePIN(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	validPin := base.TestUserPin

	var pin extension.PINExtensionImpl
	salt, encryptedPin := pin.EncryptPIN(validPin, nil)

	validSavePinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: user.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}

	type args struct {
		ctx context.Context
		pin *domain.PIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: save pin with valid pin payload",
			args: args{
				ctx: ctx,
				pin: validSavePinPayload,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "sad case: save pin with pin no payload",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.SavePIN(tt.args.ctx, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.SavePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.SavePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdatePIN(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	validPin := base.TestUserPin

	var pin extension.PINExtensionImpl
	salt, encryptedPin := pin.EncryptPIN(validPin, nil)

	validSavePinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: user.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}

	type args struct {
		ctx context.Context
		id  string
		pin *domain.PIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: update pin with valid pin payload",
			args: args{
				ctx: ctx,
				id:  user.ID,
				pin: validSavePinPayload,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "sad case: update pin with invalid payload",
			args: args{
				ctx: ctx,
				id:  "", // empty user profile
				pin: validSavePinPayload,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.UpdatePIN(tt.args.ctx, tt.args.id, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdatePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.UpdatePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_ActivateSupplierProfile(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()

	supplier, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	// Expected supplier after activation should be active
	// want == Updated supplier
	supplier.Active = true

	type args struct {
		ctx       context.Context
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Supplier
		wantErr bool
	}{
		{
			name: "Happy Case - Activate Supplier By Valid profile ID",
			args: args{
				ctx:       ctx,
				profileID: profileID,
			},
			want:    supplier,
			wantErr: false,
		},
		{
			name: "Sad Case - Activate Supplier By a non-existent profile ID",
			args: args{
				ctx:       ctx,
				profileID: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.ActivateSupplierProfile(tt.args.ctx, tt.args.profileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.ActivateSupplierProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.ActivateSupplierProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_AddPartnerType(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	testRiderName := "Test Rider"
	rider := base.PartnerTypeRider
	testPractitionerName := "Test Practitioner"
	practitioner := base.PartnerTypePractitioner
	testProviderName := "Test Provider"
	provider := base.PartnerTypeProvider

	profileID := uuid.New().String()

	supplier, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	type args struct {
		ctx         context.Context
		profileID   string
		name        *string
		partnerType *base.PartnerType
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Add a valid rider partner type",
			args: args{
				ctx:         ctx,
				profileID:   *supplier.ProfileID,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Add a valid practitioner partner type",
			args: args{
				ctx:         ctx,
				profileID:   *supplier.ProfileID,
				name:        &testPractitionerName,
				partnerType: &practitioner,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Add a valid provider partner type",
			args: args{
				ctx:         ctx,
				profileID:   *supplier.ProfileID,
				name:        &testProviderName,
				partnerType: &provider,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Use an invalid ID",
			args: args{
				ctx:         ctx,
				profileID:   "invalidid",
				name:        &testProviderName,
				partnerType: &provider,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.AddPartnerType(tt.args.ctx, tt.args.profileID, tt.args.name, tt.args.partnerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.AddPartnerType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected, got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}

			if got != tt.want {
				t.Errorf("Repository.AddPartnerType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_RecordPostVisitSurvey(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	type args struct {
		ctx   context.Context
		input resources.PostVisitSurveyInput
		UID   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record a post visit survey",
			args: args{
				ctx: ctx,
				input: resources.PostVisitSurveyInput{
					LikelyToRecommend: 10,
					Criticism:         "Nothing at all. Good job.",
					Suggestions:       "Can't think of anything.",
				},
				UID: auth.UID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid input",
			args: args{
				ctx: ctx,
				input: resources.PostVisitSurveyInput{
					LikelyToRecommend: 100,
					Criticism:         "Nothing at all. Good job.",
					Suggestions:       "Can't think of anything.",
				},
				UID: auth.UID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.RecordPostVisitSurvey(tt.args.ctx, tt.args.input, tt.args.UID); (err != nil) != tt.wantErr {
				t.Errorf("Repository.RecordPostVisitSurvey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateSuspended(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx    context.Context
		id     string
		status bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update the suspend status",
			args: args{
				ctx:    ctx,
				id:     user.ID,
				status: true,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully update the suspend status",
			args: args{
				ctx:    ctx,
				id:     user.ID,
				status: false,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an invalid id",
			args: args{
				ctx:    ctx,
				id:     "invalid id",
				status: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateSuspended(tt.args.ctx, tt.args.id, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateSuspended() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateVerifiedUIDS(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	user, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx  context.Context
		id   string
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update profile UIDs",
			args: args{
				ctx:  ctx,
				id:   user.ID,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid ID",
			args: args{
				ctx:  ctx,
				id:   "invalidid",
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateVerifiedUIDS(tt.args.ctx, tt.args.id, tt.args.uids); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateVerifiedUIDS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateVerifiedIdentifiers(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	presentIdentifiers := []base.VerifiedIdentifier{
		{
			UID:           "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
			LoginProvider: "Facebook",
		},
	}

	type args struct {
		ctx         context.Context
		id          string
		identifiers []base.VerifiedIdentifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update the user's verified identifiers",
			args: args{
				ctx: ctx,
				id:  userProfile.ID,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           auth.UID,
						LoginProvider: "Facebook",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Use a different UID",
			args: args{
				ctx:         ctx,
				id:          userProfile.ID,
				identifiers: presentIdentifiers,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Adding a new identifier",
			args: args{
				ctx: ctx,
				id:  userProfile.ID,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "5d46d3bd-a482-4787-9b87-3c94510c8b53",
						LoginProvider: "Google",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an invalid id",
			args: args{
				ctx: ctx,
				id:  "invalidid",
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           auth.UID,
						LoginProvider: "Facebook",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateVerifiedIdentifiers(tt.args.ctx, tt.args.id, tt.args.identifiers); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateVerifiedIdentifiers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateCovers(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	newCover := []base.Cover{
		{
			PayerName:      "Payer 6",
			PayerSladeCode: 27,
			MemberName:     "Jakom",
			MemberNumber:   "12345",
		},
	}

	type args struct {
		ctx    context.Context
		id     string
		covers []base.Cover
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad Case - Using an invalid ID",
			args: args{
				ctx: ctx,
				id:  "invalidID",
				covers: []base.Cover{
					{
						PayerName:      "payer1",
						PayerSladeCode: 1,
						MemberName:     "name1",
						MemberNumber:   "mem1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Add a valid new cover",
			args: args{
				ctx:    ctx,
				id:     userProfile.ID,
				covers: newCover,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateCovers(tt.args.ctx, tt.args.id, tt.args.covers); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateCovers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateSecondaryEmailAddresses(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	type args struct {
		ctx            context.Context
		id             string
		emailAddresses []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update Profile Secondary Email",
			args: args{
				ctx:            ctx,
				id:             userProfile.ID,
				emailAddresses: []string{"jatelo@gmail.com", "nyaras@gmail.com"},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Update Profile Secondary Email using an invalid ID",
			args: args{
				ctx:            ctx,
				id:             "invalid id",
				emailAddresses: []string{"jatelo@gmail.com", "nyaras@gmail.com"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateSecondaryEmailAddresses(tt.args.ctx, tt.args.id, tt.args.emailAddresses); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateSecondaryEmailAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_UpdateSupplierProfile(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	profileID := uuid.New().String()

	supplier, err := fr.CreateEmptySupplierProfile(ctx, profileID)
	if err != nil {
		t.Errorf("failed to create an empty supplier: %v", err)
	}

	validPayload := &base.Supplier{
		ID:        supplier.ID,
		ProfileID: supplier.ProfileID,
		Active:    true,
	}
	newprofileID := uuid.New().String()
	invalidPayload := &base.Supplier{
		ID:        uuid.New().String(),
		ProfileID: &newprofileID,
		Active:    true,
	}

	type args struct {
		ctx  context.Context
		data *base.Supplier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update Supplier Profile Supplier By Valid payload",
			args: args{
				ctx:  ctx,
				data: validPayload,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Update Supplier Profile By invalid payload",
			args: args{
				ctx:  ctx,
				data: invalidPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.UpdateSupplierProfile(tt.args.ctx, *tt.args.data.ProfileID, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateSupplierProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_FetchKYCProcessingRequests(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	reqPartnerType := base.PartnerTypeCoach
	organizationTypeLimitedCompany := domain.OrganizationTypeLimitedCompany
	id := uuid.New().String()
	kycReq := &domain.KYCRequest{
		ID:                  id,
		ReqPartnerType:      reqPartnerType,
		ReqOrganizationType: organizationTypeLimitedCompany,
	}

	err := fr.StageKYCProcessingRequest(ctx, kycReq)
	if err != nil {
		t.Errorf("failed to stage kyc: %v", err)
		return
	}

	kycRequests := []*domain.KYCRequest{}
	kycRequests = append(kycRequests, kycReq)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.KYCRequest
		wantErr bool
	}{
		{
			name: "Happy Case - Fetch KYC Processing Requests",
			args: args{
				ctx: ctx,
			},
			want:    kycRequests,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.FetchKYCProcessingRequests(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FetchKYCProcessingRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.FetchKYCProcessingRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdatePrimaryEmailAddress(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	newPrimaryEmail := "johndoe@gmail.com"

	type args struct {
		ctx          context.Context
		id           string
		emailAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update using a valid email",
			args: args{
				ctx:          ctx,
				id:           userProfile.ID,
				emailAddress: newPrimaryEmail,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to get logged in user",
			args: args{
				ctx:          ctx,
				id:           "invalidid",
				emailAddress: newPrimaryEmail,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdatePrimaryEmailAddress(tt.args.ctx, tt.args.id, tt.args.emailAddress); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdatePrimaryEmailAddress() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_UpdatePrimaryPhoneNumber(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	newPrimaryPhoneNumber := "+254711111111"
	type args struct {
		ctx         context.Context
		id          string
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update using a valid email",
			args: args{
				ctx:         ctx,
				id:          userProfile.ID,
				phoneNumber: newPrimaryPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to get logged in user",
			args: args{
				ctx:         ctx,
				id:          "invalidid",
				phoneNumber: newPrimaryPhoneNumber,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdatePrimaryPhoneNumber(tt.args.ctx, tt.args.id, tt.args.phoneNumber); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdatePrimaryPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_UpdateSecondaryPhoneNumbers(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	newSecondaryPhoneNumbers := []string{"+254744556677", "+254700998877"}

	type args struct {
		ctx          context.Context
		id           string
		phoneNumbers []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update secondary phonenumbers",
			args: args{
				ctx:          ctx,
				id:           userProfile.ID,
				phoneNumbers: newSecondaryPhoneNumbers,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Update secondary phonenumbers using an invalid ID",
			args: args{
				ctx:          ctx,
				id:           "invalid id",
				phoneNumbers: newSecondaryPhoneNumbers,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateSecondaryPhoneNumbers(tt.args.ctx, tt.args.id, tt.args.phoneNumbers); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateSecondaryPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_UpdateBioData(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	firstName := "Jatelo"
	lastName := "Mzima"
	dateOfBirth := base.Date{
		Year:  2000,
		Month: 12,
		Day:   17,
	}
	var gender base.Gender = "male"

	updateAllData := base.BioData{
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
		Gender:      gender,
	}

	updateFirstName := base.BioData{
		FirstName: &firstName,
	}
	updateLastName := base.BioData{
		LastName: &lastName,
	}
	updateDateOfBirth := base.BioData{
		DateOfBirth: &dateOfBirth,
	}
	updateGender := base.BioData{
		Gender: gender,
	}
	type args struct {
		ctx  context.Context
		id   string
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update all biodata",
			args: args{
				ctx:  ctx,
				id:   userProfile.ID,
				data: updateAllData,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Update firstname only",
			args: args{
				ctx:  ctx,
				id:   userProfile.ID,
				data: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Update lastname only",
			args: args{
				ctx:  ctx,
				id:   userProfile.ID,
				data: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Update date of birth only",
			args: args{
				ctx:  ctx,
				id:   userProfile.ID,
				data: updateDateOfBirth,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Update gender only",
			args: args{
				ctx:  ctx,
				id:   userProfile.ID,
				data: updateGender,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an invalid ID",
			args: args{
				ctx:  ctx,
				id:   "invalid id",
				data: updateAllData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateBioData(tt.args.ctx, tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
		})
	}
}

func TestRepository_FetchKYCProcessingRequestByID(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	reqPartnerType := base.PartnerTypeCoach
	organizationTypeLimitedCompany := domain.OrganizationTypeLimitedCompany
	id := uuid.New().String()
	kycReq := &domain.KYCRequest{
		ID:                  id,
		ReqPartnerType:      reqPartnerType,
		ReqOrganizationType: organizationTypeLimitedCompany,
	}

	err := fr.StageKYCProcessingRequest(ctx, kycReq)
	if err != nil {
		t.Errorf("failed to stage kyc: %v", err)
		return
	}

	kycRequests := []*domain.KYCRequest{}
	kycRequests = append(kycRequests, kycReq)

	kycRequest := kycRequests[0]

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.KYCRequest
		wantErr bool
	}{
		{
			name: "Happy Case - Fetch KYC Processing Requests by ID",
			args: args{
				ctx: ctx,
				id:  kycRequest.ID,
			},
			want:    kycRequest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fr.FetchKYCProcessingRequestByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FetchKYCProcessingRequestByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected, got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected, got %v", err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.FetchKYCProcessingRequestByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdateKYCProcessingRequest(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	reqPartnerType := base.PartnerTypeCoach
	organizationTypeLimitedCompany := domain.OrganizationTypeLimitedCompany
	id := uuid.New().String()
	kycReq := &domain.KYCRequest{
		ID:                  id,
		ReqPartnerType:      reqPartnerType,
		ReqOrganizationType: organizationTypeLimitedCompany,
	}

	err := fr.StageKYCProcessingRequest(ctx, kycReq)
	if err != nil {
		t.Errorf("failed to stage kyc: %v", err)
		return
	}

	kycRequests := []*domain.KYCRequest{}
	kycRequests = append(kycRequests, kycReq)

	kycRequest := kycRequests[0]

	kycStatus := domain.KYCProcessStatusApproved

	updatedKYCReq := &domain.KYCRequest{
		ID:     kycRequest.ID,
		Status: kycStatus,
	}

	updatedKYCRequests := []*domain.KYCRequest{}
	updatedKYCRequests = append(updatedKYCRequests, updatedKYCReq)

	updatedKYCRequest := updatedKYCRequests[0]

	type args struct {
		ctx        context.Context
		kycRequest *domain.KYCRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Update KYC Processing Requests",
			args: args{
				ctx:        ctx,
				kycRequest: updatedKYCRequest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fr.UpdateKYCProcessingRequest(tt.args.ctx, tt.args.kycRequest); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateKYCProcessingRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		if tt.wantErr {
			if err == nil {
				t.Errorf("error expected, got %v", err)
				return
			}
		}

		if !tt.wantErr {
			if err != nil {
				t.Errorf("error not expected, got %v", err)
				return
			}
		}
	}
}

func TestRepository_GenerateAuthCredentialsForAnonymousUser(t *testing.T) {
	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	anonymousPhoneNumber := "+254700000000"

	user, err := fr.GetOrCreatePhoneNumberUser(ctx, anonymousPhoneNumber)
	if err != nil {
		t.Errorf("failed to create a user")
		return
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, user.UID)
	if err != nil {
		t.Errorf("failed to create a custom auth token for the user")
		return
	}

	_, err = base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		t.Errorf("failed to fetch an ID token")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully generate auth credentials for anonymous user",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := fr.GenerateAuthCredentialsForAnonymousUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GenerateAuthCredentialsForAnonymousUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if *authResponse.CustomToken == "" {
					t.Errorf("nil custom token")
					return
				}

				if *authResponse.IDToken == "" {
					t.Errorf("nil ID token")
					return
				}

				if authResponse.RefreshToken == "" {
					t.Errorf("nil refresh token")
					return
				}

				if authResponse.UID == "" {
					t.Errorf("returned a nil user")
					return
				}

				if !authResponse.IsAnonymous {
					t.Errorf("the user should be anonymous")
					return
				}
			}
		})
	}
}

func TestRepository_GenerateAuthCredentials(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to create a custom auth token for the user")
		return
	}

	userToken, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		t.Errorf("failed to fetch an ID token")
		return
	}

	validCredentials := &base.AuthCredentialResponse{
		CustomToken:  &customToken,
		IDToken:      &userToken.IDToken,
		ExpiresIn:    userToken.ExpiresIn,
		RefreshToken: userToken.RefreshToken,
		UID:          auth.UID,
		IsAnonymous:  false,
		IsAdmin:      false,
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully generate valid auth credentials",
			args: args{
				ctx:   ctx,
				phone: *userProfile.PrimaryPhone,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an invalid phonenumber",
			args: args{
				ctx:   ctx,
				phone: "invalidphone",
			},
			want:    validCredentials,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := fr.GenerateAuthCredentials(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GenerateAuthCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if *authResponse.CustomToken == "" {
					t.Errorf("nil custom token")
					return
				}

				if *authResponse.IDToken == "" {
					t.Errorf("nil ID token")
					return
				}

				if authResponse.RefreshToken == "" {
					t.Errorf("nil refresh token")
					return
				}

				if authResponse.UID == "" {
					t.Errorf("returned a nil user")
					return
				}

			}
		})
	}
}

func TestRepository_FetchAdminUsers(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}

	permissions := base.DefaultAdminPermissions

	err = fr.UpdatePermissions(ctx, userProfile.ID, permissions)
	if err != nil {
		t.Errorf("failed to update user permissions: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*base.UserProfile
		wantErr bool
	}{
		{
			name: "Happy Case - Fetch admin users",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminResponse, err := fr.FetchAdminUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FetchAdminUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(adminResponse) == 0 {
					t.Errorf("nil admin response")
					return
				}

			}
		})
	}
}
