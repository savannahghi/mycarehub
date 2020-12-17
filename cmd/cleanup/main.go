package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/profile"
	"google.golang.org/api/iterator"
)

func main() {
	log.Printf("starting clean up...")
	err := deleteUnusedUserProfiles()
	if err != nil {
		log.Fatalf("can't delete collections %v:", err)
		os.Exit(-1)
	}
	log.Printf("clean up finished")
}

// An unused profile is a profile that exists in our firestore collections
// but the user UID does not exist in FireBase. The data cleaned is wrt the
// root collection suffix and project you are using in your env.
// To run this script, navigate to the file (cmd/cleanup) and run
// `go run main.go`
// It extracts uids from document data and checks if they exist in Firebase.
// If not the collection is then deleted to clean the collections

func deleteUnusedUserProfiles() error {
	fc := &base.FirebaseClient{}
	ctx := context.Background()

	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("can't initialize Firebase app when setting up profile service: %s", err)
	}

	firebaseAuth, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}

	firestore, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("can't initialize Firestore client when setting up profile service: %s", err)
	}

	service := profile.NewService()
	collection := service.GetUserProfileCollectionName()
	iter := firestore.Collection(collection).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		data := doc.Data()

		uid, present := data["uid"]
		if present {
			_, err = firebaseAuth.GetUser(ctx, uid.(string))
			if auth.IsUserNotFound(err) {
				_, err := firestore.Collection(collection).Doc(doc.Ref.ID).Delete(ctx)
				if err != nil {
					log.Printf("Can't delete a document with Ref ID %s: %v", doc.Ref.ID, err)
					continue
				}
			}
		}

		uids, present := data["uids"]
		if present {
			sliceOfUids := uids.([]interface{})
			deleteDocument(
				ctx,
				sliceOfUids,
				collection,
				doc.Ref.ID,
				firebaseAuth,
				firestore,
			)
		}

		verifiedIDs, present := data["verifiedIdentifiers"]
		if present {
			sliceOfVerifiedIDs := verifiedIDs.([]interface{})
			deleteDocument(
				ctx,
				sliceOfVerifiedIDs,
				collection,
				doc.Ref.ID,
				firebaseAuth,
				firestore,
			)
		}

	}
	return nil
}

func deleteDocument(
	ctx context.Context,
	uids []interface{},
	collection string,
	docRefID string,
	fbAuth *auth.Client,
	fsClient *firestore.Client,
) {
	if len(uids) == 0 {
		_, err := fsClient.Collection(collection).Doc(docRefID).Delete(ctx)
		if err != nil {
			log.Printf("Can't delete a document with Ref ID %s: %v", docRefID, err)
		}
	}
	var notFoundUids []string

	for _, uid := range uids {
		uidStr := uid.(string)
		_, err := fbAuth.GetUser(ctx, uidStr)
		if auth.IsUserNotFound(err) {
			notFoundUids = append(notFoundUids, uidStr)

			if len(notFoundUids) == len(uids) {
				_, err := fsClient.Collection(collection).Doc(docRefID).Delete(ctx)

				if err != nil {
					log.Printf("Can't delete a document with Ref ID %s: %v", docRefID, err)
					continue
				}
			}
			continue
		}
	}
}
