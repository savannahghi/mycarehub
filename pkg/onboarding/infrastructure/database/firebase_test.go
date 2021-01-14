package database_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, err
	}
	ext := extension.NewBaseExtensionImpl()
	otp := otp.NewOTPService(fr, ext)
	profile := usecases.NewProfileUseCase(fr, otp, ext)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes, ext)
	login := usecases.NewLoginUseCases(fr, profile, ext)
	survey := usecases.NewSurveyUseCases(fr, ext)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile, ext)
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
	exists, err := s.Signup.CheckPhoneExists(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}
	if !exists {
		otp, err := s.Otp.GenerateAndSendOTP(ctx, phone)
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
func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	originalENV := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")
	originaDEBUG := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("onboarding_ci_%v", time.Now().Unix()))
	originalROOT_S := os.Getenv("ROOT_COLLECTION_SUFFIX")
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
	purgeRecords()
	os.Setenv("ENVIRONMENT", originalENV)
	os.Setenv("DEBUG", originaDEBUG)
	os.Setenv("ROOT_COLLECTION_SUFFIX", originalROOT_S)

	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	os.Exit(code)
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

func TestCreateEmptyCustomerProfile(t *testing.T) {
	ctx := context.Background()
	firestoreDB, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Error(err)
		return
	}

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
	firestoreDB, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
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
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}
	profileID := uuid.New().String()
	_, err = fr.CreateEmptySupplierProfile(ctx, profileID)
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
	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}
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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

func TestRepository_GetSupplierProfileByProfileID(t *testing.T) {
	ctx := context.Background()

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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

	fr, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

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
