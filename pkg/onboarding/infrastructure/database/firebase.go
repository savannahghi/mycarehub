package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	userProfileCollectionName     = "user_profiles"
	supplierProfileCollectionName = "supplier_profiles"
	customerProfileCollectionName = "customer_profiles"
	pinsCollectionName            = "pins"
	surveyCollectionName          = "post_visit_survey"
	profileNudgesCollectionName   = "profile_nudges"
	kycProcessCollectionName      = "kyc_processing"

	firebaseExchangeRefreshTokenURL = "https://securetoken.googleapis.com/v1/token?key="
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	FirestoreClient *firestore.Client
	FirebaseClient  *auth.Client
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

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}

	ff := &Repository{
		FirestoreClient: fsc,
		FirebaseClient:  fbc,
	}
	err = ff.checkPreconditions()
	if err != nil {
		log.Fatalf("firebase repository precondition check failed: %s", err)
	}
	return ff, nil
}

func (fr Repository) checkPreconditions() error {
	if fr.FirestoreClient == nil {
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

// GetSurveyCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetSurveyCollectionName() string {
	suffixed := base.SuffixCollection(surveyCollectionName)
	return suffixed
}

// GetPINsCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetPINsCollectionName() string {
	suffixed := base.SuffixCollection(pinsCollectionName)
	return suffixed
}

// GetProfileNudgesCollectionName return the storage location of profile nudges
func (fr Repository) GetProfileNudgesCollectionName() string {
	suffixed := base.SuffixCollection(profileNudgesCollectionName)
	return suffixed
}

// GetKCYProcessCollectionName fetches location where kyc processing request will be saved
func (fr Repository) GetKCYProcessCollectionName() string {
	suffixed := base.SuffixCollection(kycProcessCollectionName)
	return suffixed
}

// ParseRecordAsSnapshot parses a record to firebase snapshot
func (fr Repository) ParseRecordAsSnapshot(ctx context.Context, collection string, id string) (*firestore.DocumentSnapshot, error) {
	var doc []*firestore.DocumentSnapshot
	var err error
	switch collection {
	case fr.GetUserProfileCollectionName():
		doc, err = fr.FirestoreClient.Collection(collection).Where("id", "==", id).Documents(ctx).GetAll()
	case fr.GetSupplierProfileCollectionName():
		doc, err = fr.FirestoreClient.Collection(collection).Where("id", "==", id).Documents(ctx).GetAll()
	case fr.GetCustomerProfileCollectionName():
		doc, err = fr.FirestoreClient.Collection(collection).Where("id", "==", id).Documents(ctx).GetAll()
	}
	return doc[0], err
}

// GetUserProfileByUID retrieves the user profile bu UID
func (fr *Repository) GetUserProfileByUID(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	// Retrieve the user profile
	uids := []string{uid}
	collection := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	query := collection.Where("verifiedIdentifiers", "array-contains-any", uids)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("user profile not found: %w", err)
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
	collection := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
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

	v, err := fr.CheckIfPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check the phone number: %v", err)
	}

	if v {
		// this phone is number is associated with another user profile, hence can not create an profile with the same phone number
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.PhoneNUmberInUseErrMsg,
			Code:    int(base.PhoneNumberInUse),
		}
	}

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
	docID, err := base.SaveDataToFirestore(fr.FirestoreClient, fr.GetUserProfileCollectionName(), pr)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	dsnap, err := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName()).Doc(docID).Get(ctx)
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
	docID, err := base.SaveDataToFirestore(fr.FirestoreClient, fr.GetSupplierProfileCollectionName(), sup)
	if err != nil {
		return nil, fmt.Errorf("unable to create new supplier empty profile: %w", err)
	}
	dsnap, err := fr.FirestoreClient.Collection(fr.GetSupplierProfileCollectionName()).Doc(docID).Get(ctx)
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
	docID, err := base.SaveDataToFirestore(fr.FirestoreClient, fr.GetCustomerProfileCollectionName(), cus)
	if err != nil {
		return nil, fmt.Errorf("unable to create new customer empty profile: %w", err)
	}
	dsnap, err := fr.FirestoreClient.Collection(fr.GetCustomerProfileCollectionName()).Doc(docID).Get(ctx)
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

//GetUserProfileByPrimaryPhoneNumber fetches a user profile by primary phone number
func (fr *Repository) GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
	collection1 := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("%v", base.ProfileNotFound)
	}
	dsnap := docs[0]
	profile := &base.UserProfile{}
	err = dsnap.DataTo(profile)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer profile: %w", err)
	}
	return profile, nil
}

// GetUserProfileByPhoneNumber fetches a user profile by phone number. This method traverses both PRIMARY PHONE numbers
// and SECONDARY PHONE numbers.
func (fr *Repository) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
	// check first primary phone numbers
	collection1 := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs1, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs1) == 1 {
		dsnap := docs1[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, fmt.Errorf("unable to read customer profile: %w", err)
		}
		return pr, nil
	}

	// then check in secondary phone numbers
	collection2 := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs2, err := collection2.Where("secondaryPhoneNumbers", "array-contains", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs2) == 1 {
		dsnap := docs2[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, fmt.Errorf("unable to read customer profile: %w", err)
		}
		return pr, nil
	}

	return nil, fmt.Errorf("%v", base.ProfileNotFound)

}

// CheckIfPhoneNumberExists checks both PRIMARY PHONE NUMBERs and SECONDARY PHONE NUMBERs for the
// existance of the argument phoneNnumber.
func (fr *Repository) CheckIfPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	// check first primary phone numbers
	collection1 := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs1, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	if len(docs1) == 1 {
		return true, nil
	}

	// then check in secondary phone numbers
	collection2 := fr.FirestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs2, err := collection2.Where("secondaryPhoneNumbers", "array-contains", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	if len(docs2) == 1 {
		return true, nil
	}

	return false, nil
}

// GetPINByProfileID gets a user's PIN by their profile ID
func (fr *Repository) GetPINByProfileID(ctx context.Context, profileID string) (*domain.PIN, error) {
	collection := fr.FirestoreClient.Collection(fr.GetPINsCollectionName())
	docs, err := collection.Where("profileID", "==", profileID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// this should never run. If it does, it means we are doing something wrong.
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 PINs with profile ID %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, &domain.CustomError{
			Err:     nil,
			Message: exceptions.PINNotFoundErrMsg,
			Code:    int(base.PINNotFound),
		}
	}

	dsnap := docs[0]
	PIN := &domain.PIN{}
	err = dsnap.DataTo(PIN)
	if err != nil {
		return nil, err
	}

	return PIN, nil
}

// GenerateAuthCredentials gets a Firebase user by phone and creates their tokens
func (fr *Repository) GenerateAuthCredentials(
	ctx context.Context,
	phone string,
) (*domain.AuthCredentialResponse, error) {
	u, err := fr.FirebaseClient.GetUserByPhoneNumber(ctx, phone)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, &domain.CustomError{
				Err:     err,
				Message: exceptions.UserNotFoundErrMsg,
				Code:    int(base.UserNotFound),
			}
		}
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.UserNotFoundErrMsg,
			Code:    int(base.Internal),
		}
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, u.UID)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.CustomTokenErrMsg,
			Code:    int(base.Internal),
		}
	}
	userTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.AuthenticateTokenErrMsg,
			Code:    int(base.Internal),
		}
	}
	pr, err := fr.GetUserProfileByPrimaryPhoneNumber(ctx, phone)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.AuthenticateTokenErrMsg,
			Code:    int(base.ProfileNotFound),
		}
	}

	err = fr.UpdateVerifiedIdentifiers(ctx, pr.ID, []base.VerifiedIdentifier{base.VerifiedIdentifier{
		UID:           u.UID,
		LoginProvider: base.LoginProviderTypePhone,
		Timestamp:     time.Now().In(base.TimeLocation),
	}})
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.UpdateProfileErrMsg,
			Code:    int(base.Internal),
		}
	}

	return &domain.AuthCredentialResponse{
		CustomToken:  &customToken,
		IDToken:      &userTokens.IDToken,
		ExpiresIn:    userTokens.ExpiresIn,
		RefreshToken: userTokens.RefreshToken,
	}, nil
}

// UpdateUserName updates the username of a profile that matches the id
// this method should be called after asserting the username is unique and not associated with another userProfile
func (fr *Repository) UpdateUserName(ctx context.Context, id string, userName string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.UserName = userName

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary phone number: %v", err)
	}

	return nil
}

// UpdatePrimaryPhoneNumber append a new primary phone number to the user profile
// this method should be called after asserting the phone number is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PrimaryPhone = phoneNumber

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary phone number: %v", err)
	}

	return nil
}

// UpdatePrimaryEmailAddress the primary email addresse of the profile that matches the id
// this method should be called after asserting the emailAddress is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PrimaryEmailAddress = emailAddress

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary email address: %v", err)
	}

	return nil
}

// UpdateSecondaryPhoneNumbers the secondary phone numbers of the profile that matches the id
// this method should be called after asserting the phone numbers are unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.SecondaryPhoneNumbers = phoneNumbers

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile secondary phone numbers: %v", err)
	}

	return nil
}

// UpdateSecondaryEmailAddresses the secondary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddresses  as unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryEmailAddresses(ctx context.Context, id string, emailAddresses []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.SecondaryEmailAddresses = emailAddresses

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile secondary email addresses: %v", err)
	}
	return nil
}

// UpdateSuspended updates the suspend attribute of the profile that matches the id
func (fr *Repository) UpdateSuspended(ctx context.Context, id string, status bool) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.Suspended = status

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return err
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	return err
}

// UpdatePhotoUploadID updates the photoUploadID attribute of the profile that matches the id
func (fr *Repository) UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PhotoUploadID = uploadID

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile photo upload id: %v", err)
	}
	return nil
}

// UpdateCovers updates the covers attribute of the profile that matches the id
func (fr *Repository) UpdateCovers(ctx context.Context, id string, covers []base.Cover) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.Covers = covers

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile covers: %v", err)
	}
	return nil
}

// UpdatePushTokens updates the pushTokens attribute of the profile that matches the id
func (fr *Repository) UpdatePushTokens(ctx context.Context, id string, pushToken []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	tokens := profile.PushTokens
	tokens = append(tokens, pushToken...)
	profile.PushTokens = tokens

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile push tokens: %v", err)
	}
	return nil
}

// UpdateBioData updates the biodate of the profile that matches the id
func (fr *Repository) UpdateBioData(ctx context.Context, id string, data base.BioData) error {
	profile, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}

	profile.UserBioData.FirstName = func(pr *base.UserProfile, dt base.BioData) string {
		if dt.FirstName != "" {
			return dt.FirstName
		}
		return pr.UserBioData.FirstName
	}(profile, data)

	profile.UserBioData.LastName = func(pr *base.UserProfile, dt base.BioData) string {
		if dt.LastName != "" {
			return dt.LastName
		}
		return pr.UserBioData.LastName
	}(profile, data)

	profile.UserBioData.Gender = func(pr *base.UserProfile, dt base.BioData) base.Gender {
		if dt.Gender.String() != "" {
			return dt.Gender
		}
		return pr.UserBioData.Gender
	}(profile, data)

	profile.UserBioData.DateOfBirth = func(pr *base.UserProfile, dt base.BioData) *base.Date {
		if dt.DateOfBirth == nil {
			return dt.DateOfBirth
		}
		return pr.UserBioData.DateOfBirth
	}(profile, data)

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
	if err != nil {
		return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile push tokens: %v", err)
	}
	return nil
}

// UpdateVerifiedIdentifiers adds a UID to a user profile during login if it does not exist
func (fr *Repository) UpdateVerifiedIdentifiers(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error {

	for _, identifier := range identifiers {
		profile, err := fr.GetUserProfileByID(ctx, id)
		if err != nil {
			return err
		}

		if !checkIdentifierExists(profile, identifier.UID) {
			uids := profile.VerifiedIdentifiers

			uids = append(uids, identifier)

			profile.VerifiedIdentifiers = uids

			record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetUserProfileCollectionName(), profile.ID)
			if err != nil {
				return fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err)
			}

			err = base.UpdateRecordOnFirestore(
				fr.FirestoreClient, fr.GetUserProfileCollectionName(), record.Ref.ID, profile,
			)
			if err != nil {
				return fmt.Errorf("unable to update user profile push tokens: %v", err)
			}
			return nil

		}
	}

	return nil
}

func checkIdentifierExists(profile *base.UserProfile, UID string) bool {
	foundVerifiedUIDs := []string{}
	verifiedIDs := profile.VerifiedIdentifiers
	for _, verifiedID := range verifiedIDs {
		foundVerifiedUIDs = append(foundVerifiedUIDs, verifiedID.UID)
	}
	return base.StringSliceContains(foundVerifiedUIDs, UID)
}

// RecordPostVisitSurvey records an end of visit survey
func (fr *Repository) RecordPostVisitSurvey(
	ctx context.Context,
	input *domain.PostVisitSurveyInput,
	UID string,
) (bool, error) {
	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return false, &domain.CustomError{
			Err:     nil,
			Message: exceptions.LikelyToRecommendErrMsg,
			Code:    0, // TODO: Add a code for this error
		}

	}
	feedbackCollection := fr.FirestoreClient.Collection(fr.GetSurveyCollectionName())
	feedback := domain.PostVisitSurvey{
		LikelyToRecommend: input.LikelyToRecommend,
		Criticism:         input.Criticism,
		Suggestions:       input.Suggestions,
		UID:               UID,
		Timestamp:         time.Now(),
	}
	_, _, err := feedbackCollection.Add(ctx, feedback)
	if err != nil {
		return false, &domain.CustomError{
			Err:     err,
			Message: exceptions.AddRecordErrMsg,
			Code:    int(base.Internal),
		}

	}
	return true, nil
}

// SavePIN  persist the data of the newly created PIN to a datastore
func (fr *Repository) SavePIN(ctx context.Context, pin *domain.PIN) (*domain.PIN, error) {
	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(fr.FirestoreClient, fr.GetPINsCollectionName(), pin)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.AddRecordErrMsg,
			Code:    int(base.Internal),
		}
	}
	dsnap, err := fr.FirestoreClient.Collection(fr.GetPINsCollectionName()).Doc(docID).Get(ctx)

	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.RetrieveRecordErrMsg,
			Code:    int(base.Internal),
		}
	}

	PIN := &domain.PIN{}
	err = dsnap.DataTo(PIN)
	if err != nil {
		return nil, err
	}

	return PIN, nil

}

// UpdatePIN  persist the data of the updated PIN to a datastore
func (fr *Repository) UpdatePIN(ctx context.Context, pin *domain.PIN) (*domain.PIN, error) {
	profile, err := fr.GetPINByProfileID(ctx, pin.ProfileID)
	if err != nil {
		return nil, err
	}
	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetPINsCollectionName(), profile.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user pin as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetPINsCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.UpdateProfileErrMsg,
			Code:    int(base.Internal),
		}
	}
	return profile, nil

}

// ExchangeRefreshTokenForIDToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (fr Repository) ExchangeRefreshTokenForIDToken(refreshToken string) (*domain.AuthCredentialResponse, error) {
	apiKey, err := base.GetEnvVar(base.FirebaseWebAPIKeyEnvVarName)
	if err != nil {
		return nil, err
	}

	payload := domain.RefreshTokenExchangePayload{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := firebaseExchangeRefreshTokenURL + apiKey
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * base.HTTPClientTimeoutSecs
	resp, err := httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(payloadBytes),
	)

	defer base.CloseRespBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"firebase HTTP error, status code %d\nBody: %s\nBody read error: %s",
			resp.StatusCode,
			string(bs),
			err,
		)
	}

	var tokenResp domain.AuthCredentialResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetSupplierProfileByProfileID fetch the supplier profile by profile id.
// since this same supplierProfile can be used for updating, a companion snapshot record is returned as well
func (fr *Repository) GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*domain.Supplier, error) {
	collection := fr.FirestoreClient.Collection(fr.GetSupplierProfileCollectionName())
	query := collection.Where("profileID", "==", profileID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("supplier profile not found: %w", err)
	}
	dsnap := docs[0]
	sup := &domain.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier profile: %w", err)
	}
	return sup, nil
}

// GetSupplierProfileByID fetches supplier profile by given ID
func (fr *Repository) GetSupplierProfileByID(ctx context.Context, id string) (*domain.Supplier, error) {
	collection := fr.FirestoreClient.Collection(fr.GetSupplierProfileCollectionName())
	query := collection.Where("id", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("supplier profile not found: %w", err)
	}
	dsnap := docs[0]
	sup := &domain.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier profile: %w", err)
	}
	return sup, nil
}

// AddPartnerType updates the suppier profile with the provided name and  partner type.
func (fr *Repository) AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *domain.PartnerType) (bool, error) {

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return false, err
	}

	sup.SupplierName = *name
	sup.PartnerType = *partnerType
	sup.PartnerSetupComplete = true

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetSupplierProfileCollectionName(), sup.ID)
	if err != nil {
		return false, fmt.Errorf("unable to parse supplier profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, sup,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil

}

// ActivateSupplierProfile sets the active attribute of supplier profile to true
func (fr *Repository) ActivateSupplierProfile(ctx context.Context, profileID string) (*domain.Supplier, error) {
	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	sup.Active = true

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetSupplierProfileCollectionName(), sup.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse supplier profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, sup,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier profile: %v", err)
	}

	return sup, nil
}

// UpdateSupplierProfile update the supplier profile
func (fr *Repository) UpdateSupplierProfile(ctx context.Context, data *domain.Supplier) (*domain.Supplier, error) {
	sup, err := fr.GetSupplierProfileByProfileID(ctx, *data.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch supplier profile: %v", err)
	}

	record, err := fr.ParseRecordAsSnapshot(ctx, fr.GetSupplierProfileCollectionName(), sup.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse supplier profile as firebase snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		fr.FirestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, sup,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier profile: %v", err)
	}

	return sup, nil

}

// StageProfileNudge ...
func (fr *Repository) StageProfileNudge(ctx context.Context, nudge map[string]interface{}) error {
	_, _, err := fr.FirestoreClient.Collection(fr.GetProfileNudgesCollectionName()).Add(ctx, nudge)
	return err
}

// FetchKYCProcessingRequests retrieves all unprocessed kycs for admins
func (fr *Repository) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	collection := fr.FirestoreClient.Collection(fr.GetKCYProcessCollectionName())
	query := collection.Where("proceseed", "==", false)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch kyc request documents: %v", err)
	}

	res := []*domain.KYCRequest{}

	for _, doc := range docs {
		req := &domain.KYCRequest{}
		err = doc.DataTo(req)
		if err != nil {
			return nil, fmt.Errorf("unable to read supplier: %w", err)
		}
		res = append(res, req)
	}
	return res, nil
}
