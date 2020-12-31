package usecases_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("DEBUG", "true")
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("onboarding_ci_%v", time.Now().Unix()))
	ctx := context.Background()
	r := database.Repository{} // They are nil
	fsc, _ := initializeTestFirebaseClient(ctx)

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

	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	os.Exit(code)
}

func initializeTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
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
	otp := usecases.NewOTPUseCasesImpl(fr)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage)
	login := usecases.NewLoginUseCases(fr)
	survey := usecases.NewSurveyUseCases(fr)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier)

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

func createTestPhoneUser(ctx context.Context, flavour base.Flavour) (*domain.UserResponse, error) {
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service")
	}
	u, err := s.Signup.CreateUserByPhone(
		ctx,
		base.TestUserPhoneNumber,
		base.TestUserPin,
		flavour,
	)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("nil user response")
	}

	return u, nil
}
func TestLoginUseCasesImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()
	flavour := base.FlavourConsumer
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	u, err := createTestPhoneUser(ctx, flavour)
	if err != nil {
		t.Errorf("failed to create test phone user: %v", err)
		return
	}
	if u == nil {
		t.Errorf("nil user response")
		return
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid login",
			args: args{
				ctx:     ctx,
				phone:   base.TestUserPhoneNumber,
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "sad case: wrong pin number supplied",
			args: args{
				ctx:     ctx,
				phone:   base.TestUserPhoneNumber,
				PIN:     "4567",
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: user profile without a primary phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2547900900", // not a primary phone number
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2541234",
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := s
			authResponse, err := l.Login.LoginByPhone(tt.args.ctx, tt.args.phone, tt.args.PIN, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUseCasesImpl.LoginByPhone() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && authResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}
