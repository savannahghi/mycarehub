package database

import (
	"context"
	"log"

	"github.com/savannahghi/firebasetools"
	libDomain "github.com/savannahghi/onboarding/pkg/onboarding/domain"
	libDatabase "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database"
	libFirestore "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/serverutils"
)

// Repository ...
type Repository interface {
	libDatabase.Repository
}

// DbService is an implementation of the database repository
// It is implementation agnostic i.e logic should be handled using
// the preferred database
type DbService struct {
	Repository
}

// NewDbService creates a new database service
func NewDbService() *DbService {
	ctx := context.Background()

	var repo Repository

	if serverutils.MustGetEnvVar(libDomain.Repo) == libDomain.FirebaseRepository {

		fc := &firebasetools.FirebaseClient{}
		firebaseApp, err := fc.InitFirebase()
		if err != nil {
			return nil
		}
		fbc, err := firebaseApp.Auth(ctx)
		if err != nil {
			log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
		}
		fsc, err := firebaseApp.Firestore(ctx)
		if err != nil {
			log.Fatalf("unable to initialize Firestore: %s", err)
		}
		firestoreExtension := libFirestore.NewFirestoreClientExtension(fsc)

		repo = libFirestore.NewFirebaseRepository(firestoreExtension, fbc)
	}

	return &DbService{
		repo,
	}
}
