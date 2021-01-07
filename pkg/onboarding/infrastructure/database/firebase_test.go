package database_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

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
				t.Errorf("error expected but returned no erro")
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
