package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/config/errors"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	userProfileCollectionName     = "user_profiles"
	supplierProfileCollectionName = "supplier_profiles"
	customerProfileCollectionName = "customer_profiles"
	pinsCollectionName            = "pins"
	surveyCollectionName          = "post_visit_survey"
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	firestoreClient *firestore.Client
	firebaseClient  *auth.Client
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
		firestoreClient: fsc,
		firebaseClient:  fbc,
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

// GetUserProfileByUID retrieves the user profile bu UID
func (fr *Repository) GetUserProfileByUID(
	ctx context.Context,
	uid string,
) (*base.UserProfile, *firestore.DocumentSnapshot, error) {
	// Retrieve the user profile
	uids := []string{uid}
	collection := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	query := collection.Where("verifiedIdentifiers", "array-contains-any", uids)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("user with uids %s has > 1 profile (they have %d)", uids, len(docs))
	}
	// read and unpack profile
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	return userProfile, dsnap, nil
}

// GetUserProfileByID retrieves a user profile by ID
func (fr *Repository) GetUserProfileByID(
	ctx context.Context,
	id string,
) (*base.UserProfile, *firestore.DocumentSnapshot, error) {
	collection := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	query := collection.Where("id", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, nil, fmt.Errorf("user profile not found: %w", err)
	}
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	return userProfile, dsnap, nil
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

//GetUserProfileByPrimaryPhoneNumber fetches a user profile by primary phone number
func (fr *Repository) GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, *firestore.DocumentSnapshot, error) {
	collection1 := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}
	if len(docs) == 1 {
		dsnap := docs[0]
		profile := &base.UserProfile{}
		if err := dsnap.DataTo(profile); err != nil {
			return nil, nil, fmt.Errorf("unable to read customer profile: %w", err)
		}
		return profile, dsnap, nil
	}
	return nil, nil, fmt.Errorf("%v", base.ProfileNotFound)
}

// GetUserProfileByPhoneNumber fetches a user profile by phone number. This method traverses both PRIMARY PHONE numbers
// and SECONDARY PHONE numbers.
func (fr *Repository) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, *firestore.DocumentSnapshot, error) {
	// check first primary phone numbers
	collection1 := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs1, err := collection1.Where("primaryPhone", "==", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}
	if len(docs1) == 1 {
		dsnap := docs1[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, nil, fmt.Errorf("unable to read customer profile: %w", err)
		}
		return pr, dsnap, nil
	}

	// then check in secondary phone numbers
	collection2 := fr.firestoreClient.Collection(fr.GetUserProfileCollectionName())
	docs2, err := collection2.Where("secondaryPhoneNumbers", "array-contains", phoneNumber).Documents(ctx).GetAll()
	if err != nil {
		return nil, nil, err
	}

	if len(docs2) == 1 {
		dsnap := docs2[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, nil, fmt.Errorf("unable to read customer profile: %w", err)
		}
		return pr, dsnap, nil
	}

	return nil, nil, fmt.Errorf("%v", base.ProfileNotFound)

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

// GetPINByProfileID gets a user's PIN by their profile ID
func (fr *Repository) GetPINByProfileID(ctx context.Context, profileID string) (*domain.PIN, error) {
	collection := fr.firestoreClient.Collection(fr.GetPINsCollectionName())
	query := collection.Where("profileID", "==", profileID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// this should never run. If it does, it means we are doing something wrong.
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 PINs with profile ID %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("%v", base.PINNotFound)
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
	u, err := fr.firebaseClient.GetUserByPhoneNumber(ctx, phone)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, &domain.CustomError{
				Err:     err,
				Message: errors.UserNotFoundErrMsg,
				Code:    int(base.UserNotFound),
			}
		}
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.UserNotFoundErrMsg,
			Code:    int(base.Internal),
		}
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, u.UID)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.CustomTokenErrMsg,
			Code:    int(base.Internal),
		}
	}

	userTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.AuthenticateTokenErrMsg,
			Code:    int(base.Internal),
		}
	}

	pr, _, err := fr.GetUserProfileByPrimaryPhoneNumber(ctx, phone)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.AuthenticateTokenErrMsg,
			Code:    int(base.ProfileNotFound),
		}
	}

	err = fr.UpdateVerifiedIdentifiers(ctx, pr.ID, u.UID)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.UpdateProfileErrMsg,
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
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.UserName = userName

	err = base.UpdateRecordOnFirestore(fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary phone number: %v", err)
	}

	return nil
}

// UpdatePrimaryPhoneNumber append a new primary phone number to the user profile
// this method should be called after asserting the phone number is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PrimaryPhone = phoneNumber

	err = base.UpdateRecordOnFirestore(fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary phone number: %v", err)
	}

	return nil
}

// UpdatePrimaryEmailAddress the primary email addresse of the profile that matches the id
// this method should be called after asserting the emailAddress is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PrimaryEmailAddress = emailAddress

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile primary email address: %v", err)
	}

	return nil
}

// UpdateSecondaryPhoneNumbers the secondary phone numbers of the profile that matches the id
// this method should be called after asserting the phone numbers are unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.SecondaryPhoneNumbers = phoneNumbers

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile secondary phone numbers: %v", err)
	}

	return nil
}

// SavePIN  persist the data of the newly created PIN to a datastore
func (fr *Repository) SavePIN(ctx context.Context, pin *domain.PIN) (*domain.PIN, error) {
	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(fr.firestoreClient, fr.GetPINsCollectionName(), pin)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.AddRecordErrMsg,
			Code:    int(base.Internal),
		}
	}
	dsnap, err := fr.firestoreClient.Collection(fr.GetPINsCollectionName()).Doc(docID).Get(ctx)

	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: errors.RetrieveRecordErrMsg,
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

// UpdateSecondaryEmailAddresses the secondary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddresses  as unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryEmailAddresses(ctx context.Context, id string, emailAddresses []string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.SecondaryEmailAddresses = emailAddresses

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile secondary email addresses: %v", err)
	}
	return nil
}

// UpdateSuspended updates the suspend attribute of the profile that matches the id
// todo: (dexter) : amend interface method in base to return error instead of bool
func (fr *Repository) UpdateSuspended(ctx context.Context, id string, status bool) bool {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return false
	}
	profile.Suspended = status

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	return err == nil
}

// UpdatePhotoUploadID updates the photoUploadID attribute of the profile that matches the id
func (fr *Repository) UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.PhotoUploadID = uploadID

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile photo upload id: %v", err)
	}
	return nil
}

// UpdateCovers updates the covers attribute of the profile that matches the id
func (fr *Repository) UpdateCovers(ctx context.Context, id string, covers []base.Cover) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	profile.Covers = covers

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile covers: %v", err)
	}
	return nil
}

// UpdatePushTokens updates the pushTokens attribute of the profile that matches the id
// todo(dexter) : amend the interface method in base to take pushToken as []string instead of string
func (fr *Repository) UpdatePushTokens(ctx context.Context, id string, pushToken string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}
	tokens := profile.PushTokens
	tokens = append(tokens, pushToken)
	profile.PushTokens = tokens

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile push tokens: %v", err)
	}
	return nil
}

// UpdateBioData updates the biodate of the profile that matches the id
func (fr *Repository) UpdateBioData(ctx context.Context, id string, data base.BioData) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
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

	err = base.UpdateRecordOnFirestore(
		fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile push tokens: %v", err)
	}
	return nil
}

// UpdateVerifiedIdentifiers adds a UID to a user profile during login if it does not exist
// todo: (dexter define this in the base interface)
func (fr *Repository) UpdateVerifiedIdentifiers(ctx context.Context, id string, UID string) error {
	profile, record, err := fr.GetUserProfileByID(ctx, id)
	if err != nil {
		return err
	}

	if !checkIdentifierExists(profile, UID) {
		uids := profile.VerifiedIdentifiers
		uids = append(uids, base.VerifiedIdentifier{
			UID:           UID,
			LoginProvider: base.LoginProviderTypePhone,
			Timestamp:     time.Now().In(base.TimeLocation),
		})
		profile.VerifiedIdentifiers = uids

		err = base.UpdateRecordOnFirestore(
			fr.firestoreClient, fr.GetSupplierProfileCollectionName(), record.Ref.ID, profile,
		)
		if err != nil {
			return fmt.Errorf("unable to update user profile push tokens: %v", err)
		}
		return nil

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
			Message: errors.LikelyToRecommendErrMsg,
			Code:    0, // TODO: Add a code for this error
		}

	}
	feedbackCollection := fr.firestoreClient.Collection(fr.GetSurveyCollectionName())
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
			Message: errors.AddRecordErrMsg,
			Code:    int(base.Internal),
		}

	}
	return true, nil
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
