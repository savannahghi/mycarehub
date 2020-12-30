package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	userProfileCollectionName     = "user_profiles"
	supplierProfileCollectionName = "supplier_profiles"
	customerProfileCollectionName = "customer_profiles"
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

// GetSupplierProfileCollectionName ...
func (fr Repository) GetSupplierProfileCollectionName() string {
	suffixed := base.SuffixCollection(supplierProfileCollectionName)
	return suffixed
}

// GetCustomerProfileCollectionName ...
func (fr Repository) GetCustomerProfileCollectionName() string {
	suffixed := base.SuffixCollection(customerProfileCollectionName)
	return suffixed
}

// GetUserProfileByUID retrieves the user profile bu UID
func (fr *Repository) GetUserProfileByUID(
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

// GetSupplierProfileByProfileID fetch the supplier profile by profile id.
// since this same supplierProfile can be used for updating, a companion snapshot record is returned as well
func (fr *Repository) GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*domain.Supplier, *firestore.DocumentSnapshot, error) {
	collection := fr.firestoreClient.Collection(fr.GetSupplierProfileCollectionName())
	query := collection.Where("profileID", "==", profileID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("supplier profile not found: %w", err)
	}
	dsnap := docs[0]
	sup := &domain.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read supplier profile: %w", err)
	}
	return sup, dsnap, nil
}

// CreateUserProfile creates a user profile of using the provided phone number and uid
func (fr *Repository) CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {

	pr := &base.UserProfile{
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
	docID, err := base.SaveDataToFirestore(fr.firestoreClient, fr.GetUserProfileCollectionName(), pr)
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

// CreateEmptySupplierProfile creates an empty supplier profile
func (fr *Repository) CreateEmptySupplierProfile(ctx context.Context, profileID string) (*domain.Supplier, error) {
	sup := &domain.Supplier{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(fr.firestoreClient, fr.GetSupplierProfileCollectionName(), sup)
	if err != nil {
		return nil, fmt.Errorf("unable to create new supplier empty profile: %w", err)
	}
	dsnap, err := fr.firestoreClient.Collection(fr.GetSupplierProfileCollectionName()).Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve newly created supplier profile: %w", err)
	}
	// return the newly created supplier profile
	supplier := &domain.Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier profile: %w", err)
	}
	return supplier, nil

}

// CreateEmptyCustomerProfile creates an empty customer profile
func (fr *Repository) CreateEmptyCustomerProfile(ctx context.Context, profileID string) (*domain.Customer, error) {
	cus := &domain.Customer{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(fr.firestoreClient, fr.GetCustomerProfileCollectionName(), cus)
	if err != nil {
		return nil, fmt.Errorf("unable to create new customer empty profile: %w", err)
	}
	dsnap, err := fr.firestoreClient.Collection(fr.GetCustomerProfileCollectionName()).Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve newly created customer profile: %w", err)
	}
	// return the newly created customer profile
	customer := &domain.Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer profile: %w", err)
	}
	return customer, nil
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

// AddPartnerType updates the suppier profile with the provided name and  partner type.
func (fr *Repository) AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *domain.PartnerType) (bool, error) {

	// get the suppier profile
	sup, record, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return false, err
	}

	sup.SupplierName = *name
	sup.PartnerType = *partnerType
	sup.PartnerSetupComplete = true

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, sup,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil

}
