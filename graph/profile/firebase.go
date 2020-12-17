package profile

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
)

// RetrieveFireStoreSnapshotByUID retrieves a specified Firestore document snapshot by its UID
func (s Service) RetrieveFireStoreSnapshotByUID(
	ctx context.Context, uid string, collectionName string,
	field string) (*firestore.DocumentSnapshot, error) {
	collection := s.firestoreClient.Collection(collectionName)
	query := collection.Where(field, "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("more than one snapshot found (they have %d)", len(docs))
		}
	}
	if len(docs) == 0 {
		return nil, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// RetrieveUserProfileFirebaseDocSnapshotByUID retrieves the user profile
func (s Service) RetrieveUserProfileFirebaseDocSnapshotByUID(
	ctx context.Context,
	uid string,
) (*firestore.DocumentSnapshot, error) {
	uids := []string{uid}
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("verifiedIdentifiers", "array-contains-any", uids)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {

		log.Printf("user with uids %s has > 1 profile (they have %d)", uids, len(docs))
	}
	if len(docs) > 1 {
		for i, doc := range docs {
			if i != 0 {
				_, err := doc.Ref.Delete(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to delete profile to avoid multiple user profile: %w", err)
				}

			}
		}
	}
	if len(docs) == 0 {
		newProfile := &base.UserProfile{
			ID:                  uuid.New().String(),
			VerifiedIdentifiers: uids,
			IsApproved:          false,
			TermsAccepted:       false,
			CanExperiment:       false,
			Active:              true,
		}
		docID, err := base.SaveDataToFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
		if err != nil {
			return nil, fmt.Errorf("unable to create new user profile: %w", err)
		}
		dsnap, err := collection.Doc(docID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
		}
		return dsnap, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// RetrieveUserProfileFirebaseDocSnapshotByID retrieves a user profile by ID
func (s Service) RetrieveUserProfileFirebaseDocSnapshotByID(
	ctx context.Context,
	id string,
) (*firestore.DocumentSnapshot, error) {
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("id", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	// allow user to have one profile by deleting the other profiles.
	if len(docs) > 1 {
		for i, doc := range docs {
			if i != 0 {
				_, err := doc.Ref.Delete(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to delete profile to avoid multiple user profile: %w", err)
				}

			}
		}
	}
	if len(docs) == 0 {
		newProfile := &base.UserProfile{
			ID:                  uuid.New().String(),
			VerifiedIdentifiers: []string{},
			IsApproved:          false,
			TermsAccepted:       false,
			CanExperiment:       false,
		}
		docID, err := base.SaveDataToFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
		if err != nil {
			return nil, fmt.Errorf("unable to create new user profile: %w", err)
		}
		dsnap, err := collection.Doc(docID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
		}
		return dsnap, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// RetrieveOrCreateUserProfileFirebaseDocSnapshot retrieves the user profile of a
// specified user using either their uid or phone number.
// If the user profile does not exist then a new one is created
func (s Service) RetrieveOrCreateUserProfileFirebaseDocSnapshot(
	ctx context.Context,
	uid string,
	phone string,
) (*firestore.DocumentSnapshot, error) {
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	// the ordering is necessary in order to provide a stable sort order
	query := collection.Where("verifiedIdentifiers", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("user %s has > 1 profile (they have %d)", uid, len(docs))
		}
	}
	// remove the other profiles and retain the first one that was initially created
	if len(docs) > 1 {
		for i, doc := range docs {
			if i != 0 {
				_, err := doc.Ref.Delete(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to delete profile to avoid multiple user profile: %w", err)
				}

			}
		}
	}
	var uids []string
	var msisdns []string

	if len(docs) == 0 {
		collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
		query := collection.Where("msisdns", "array-contains", phone)
		docs, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
		if len(docs) > 1 {
			if base.IsDebug() {
				log.Printf("phone number %s is in > 1 profile (%d)", phone, len(docs))
			}
		}

		if len(docs) == 0 {
			uids = append(uids, uid)
			msisdns = append(msisdns, phone)
			// generate a new internal ID for the profile
			newProfile := &base.UserProfile{
				ID:                  uuid.New().String(),
				VerifiedIdentifiers: uids,
				IsApproved:          false,
				TermsAccepted:       false,
				CanExperiment:       false,
				Msisdns:             msisdns,
				Active:              true,
			}
			docID, err := base.SaveDataToFirestore(
				s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
			if err != nil {
				return nil, fmt.Errorf("unable to create new user profile: %w", err)
			}
			dsnap, err := collection.Doc(docID).Get(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
			}
			return dsnap, nil
		}

		dsnap := docs[0]
		return dsnap, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// RetrieveUserProfileFirebaseDocSnapshot retrievs a raw Firebase doc snapshot
// for the logged in user's user profile or creates one if it does not exist
func (s Service) RetrieveUserProfileFirebaseDocSnapshot(
	ctx context.Context) (*firestore.DocumentSnapshot, error) {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return s.RetrieveUserProfileFirebaseDocSnapshotByUID(ctx, uid)
}

// RetrievePINFirebaseDocSnapshotByMSISDN retrieves the user profile of a
// specified user
func (s Service) RetrievePINFirebaseDocSnapshotByMSISDN(
	ctx context.Context,
	msisdn string,
) (*firestore.DocumentSnapshot, error) {

	collection := s.firestoreClient.Collection(s.GetPINCollectionName())
	query := collection.Where("msisdn", "==", msisdn).Where("isValid", "==", true)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("msisdn %s has more than one PIN (it has %d)", msisdn, len(docs))
		}
	}
	// remove the other PINs and retain the first one that was initially created
	if len(docs) > 1 {
		for i, doc := range docs {
			if i != 0 {
				_, err := doc.Ref.Delete(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to delete PIN to avoid multiple user PINs: %w", err)
				}

			}
		}
	}
	if len(docs) == 0 {
		return nil, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// AddPractitionerHelper helper to add a practitioner to use in tests
func (s Service) AddPractitionerHelper(practitioner *Practitioner) (*string, error) {
	docID, err := base.SaveDataToFirestore(
		s.firestoreClient, s.GetPractitionerCollectionName(), practitioner)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	return &docID, nil
}

// CreateUserProfile creates a user profile in the database given a phone number and verified firebase auth uid
func (s Service) CreateUserProfile(ctx context.Context, msisdn, uid string) (*base.UserProfile, error) {
	// prepare the user payload
	var uids []string
	var msisdns []string
	uids = append(uids, uid)
	msisdns = append(msisdns, msisdn)
	newProfile := &base.UserProfile{
		ID:                  uuid.New().String(),
		VerifiedIdentifiers: uids,
		IsApproved:          false,
		TermsAccepted:       false,
		CanExperiment:       false,
		Msisdns:             msisdns,
		Active:              true,
	}
	// persist the data to a datastore
	docID, err := base.SaveDataToFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	dsnap, err := collection.Doc(docID).Get(ctx)
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

// RetrieveAndUpdateOldProfile retrieve old profile and update with new
func (s Service) RetrieveAndUpdateOldProfile(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("uid", "==", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	// No profile found with uid check with uids
	if len(docs) == 0 {
		// check with uids
		return s.RetrieveOldProfileByUIDS(ctx, uid)
	}
	return s.AssignNewProfile(ctx, uid)
}

// RetrieveOldProfileByUIDS retrieve old profiles  that used uids as the primary reference
func (s Service) RetrieveOldProfileByUIDS(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("uids", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	// No profile found assign new one
	if len(docs) == 0 {
		return s.AssignNewProfile(ctx, uid)
	}
	// user has an old profile give them a new one
	return s.AssignNewProfile(ctx, uid)
}

// AssignNewProfile upgrades/refreshes/creates a new profile for the user
// this will ensure the user has a profile that has fields that reflect
// new changes that were introduced to the new profile model.
func (s Service) AssignNewProfile(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	// retrieve the user from firebase auth
	uids := []string{uid}
	user, userErr := s.firebaseAuth.GetUser(ctx, uid)
	if userErr != nil {
		return nil, fmt.Errorf("unable to get Firebase user with UID %s: %w", uid, userErr)
	}
	// update their new profile with firebase verified phone number
	var msisdns []string
	msisdns = append(msisdns, user.PhoneNumber)
	newProfile := &base.UserProfile{
		ID:                  uuid.New().String(),
		VerifiedIdentifiers: uids,
		IsApproved:          false,
		TermsAccepted:       false,
		CanExperiment:       false,
		Active:              true,
		Msisdns:             msisdns,
		// for backward compatibility old user profiles  will have to reset thier PIN
		HasPin: true,
	}
	// persist the profile to the datastore
	docID, err := base.SaveDataToFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	dsnap, err := collection.Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
	}
	// read and unpack profile
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	// check whether they are Testers and update accordingly
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// GetUserProfile retrieves the user profile
func (s Service) GetUserProfile(
	ctx context.Context,
	uid string,
) (*base.UserProfile, error) {
	// Retrieve the user profile
	uids := []string{uid}
	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("verifiedIdentifiers", "array-contains-any", uids)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("user with uids %s has > 1 profile (they have %d)", uids, len(docs))
	}
	// if user has more than one profile retain the first one initially created
	if len(docs) > 1 {
		for i, doc := range docs {
			if i != 0 {
				_, err := doc.Ref.Delete(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to delete profile to avoid multiple user profile: %w", err)
				}

			}
		}
	}
	// no profile found check for old ones and link/assign a new one
	// this ensures we port old profiles to new for backward compatibility
	if len(docs) == 0 {
		// check if the user has any old profile and update with new
		return s.RetrieveAndUpdateOldProfile(ctx, uid)
	}
	// read and unpack profile
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	// check whether they are Testers and update accordingly
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	// all userprofiles should have a PIN set for backward compatibility
	userProfile.HasPin = true
	return userProfile, nil
}
