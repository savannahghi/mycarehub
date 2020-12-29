package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	userProfileCollectionName = "user_profiles"
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	firestoreClient *firestore.Client
}

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(ctx context.Context) (repository.OnboardingRepository, error) {
	fc := base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Fatalf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}

	ff := &Repository{
		firestoreClient: fsc,
	}
	err = ff.checkPreconditions()
	if err != nil {
		log.Fatalf("firebase repository precondition check failed: %s", err)
	}
	return ff, nil
}

func (fr Repository) checkPreconditions() error {
	if fr.firestoreClient == nil {
		return fmt.Errorf("nil firestore client in feed firebase repository")
	}

	return nil
}

// GetUserProfileCollectionName ...
func (fr Repository) GetUserProfileCollectionName() string {
	suffixed := base.SuffixCollection(userProfileCollectionName)
	return suffixed
}

// GetUserProfile retrieves the user profile
func (fr *Repository) GetUserProfile(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	// Retrieve the user profile
	uids := []string{uid}
	collection := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	query := collection.Where("verifiedIdentifiers", "array-contains-any", uids)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("user with uids %s has > 1 profile (they have %d)", uids, len(docs))
	}
	// read and unpack profile
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	return userProfile, nil
}

// GetUserProfileByID retrieves a user profile by ID
func (fr *Repository) GetUserProfileByID(
	ctx context.Context,
	id string,
) (*base.UserProfile, error) {
	collection := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	query := collection.Where("id", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("user profile not found: %w", err)
	}
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	return userProfile, nil
}

// CreateUserProfile creates a user profile of using the provided phone number and uid
func (fr *Repository) CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {

	newProfile := &base.UserProfile{
		ID:           uuid.New().String(),
		PrimaryPhone: phoneNumber,
		VerifiedIdentifiers: []base.VerifiedIdentifier{{
			UID:           uid,
			LoginProvider: base.LoginProviderTypePhone,
			Timestamp:     time.Now().In(base.TimeLocation),
		}},
		TermsAccepted: true,
		Suspended:     false,
	}

	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(fr.firestoreClient, fr.GetUserProfileCollectionName(), newProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	dsnap, err := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName()).Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
	}
	// return the newly created user profile
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	return userProfile, nil

}

// CheckIfPhoneNumberExists checks both PRIMARY PHONE NUMBERs and SECONDARY PHONE NUMBERs for the
// existance of the argument phoneNnumber.
func (fr *Repository) CheckIfPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	// check first primary phone numbers
	collection1 := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs1, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	if len(docs1) == 1 {
		return true, nil
	}

	// then check in secondary phone numbers
	collection2 := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs2, err := collection2.Where("secondaryPhoneNumbers", "array-contains", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	if len(docs2) == 1 {
		return true, nil
	}

	return false, nil
}

// UpdatePrimaryPhoneNumber ...
func (fr *Repository) UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PrimaryPhone = phoneNumber
	// todo : save
	return nil
}

// UpdatePrimaryEmailAddress ...
func (fr *Repository) UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error {
	return nil
}

// UpdateSecondaryPhoneNumbers ...
func (fr *Repository) UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error {
	return nil
}

// UpdateSecondaryEmailAddresses ...
func (fr *Repository) UpdateSecondaryEmailAddresses(ctx context.Context, id string, emailAddresses []string) error {
	return nil
}

// UpdateSuspended ...
func (fr *Repository) UpdateSuspended(ctx context.Context, id string) bool {
	return false
}

// UpdatePhotoUploadID ...
func (fr *Repository) UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error {
	return nil
}

// UpdateCovers ...
func (fr *Repository) UpdateCovers(ctx context.Context, id string, covers []base.Cover) error {
	return nil
}

// UpdatePushTokens ...
func (fr *Repository) UpdatePushTokens(ctx context.Context, id string, pushToken string) error {
	return nil
}

// UpdateBioData ...
func (fr *Repository) UpdateBioData(ctx context.Context, id string, data base.BioData) error {
	return nil
}
