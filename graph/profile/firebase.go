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

// RetrieveUserProfileFirebaseDocSnapshotByUID retrieves the user profile of a
// specified user
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
	if len(docs) == 0 {
		newProfile := &base.UserProfile{
			ID:                  uuid.New().String(),
			VerifiedIdentifiers: uids,
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

// AddPractitionerHelper helper to add a practitioner to use in tests
func (s Service) AddPractitionerHelper(practitioner *Practitioner) (*string, error) {
	docID, err := base.SaveDataToFirestore(
		s.firestoreClient, s.GetPractitionerCollectionName(), practitioner)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user profile: %w", err)
	}
	return &docID, nil
}
