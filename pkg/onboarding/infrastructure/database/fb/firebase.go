package fb

// import (
// 	libFirestore "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
// )

// // Repository accesses and updates an item that is stored on Firebase
// type Repository struct {
// 	FirestoreClient libFirestore.FirestoreClientExtension
// 	FirebaseClient  libFirestore.FirebaseClientExtension
// 	library         *libFirestore.Repository
// }

// // NewFirebaseRepository initializes a Firebase repository
// func NewFirebaseRepository(
// 	firestoreClient libFirestore.FirestoreClientExtension,
// 	firebaseClient libFirestore.FirebaseClientExtension,
// ) *Repository {

// 	lib := libFirestore.NewFirebaseRepository(firestoreClient, firebaseClient)

// 	return &Repository{
// 		FirestoreClient: firestoreClient,
// 		FirebaseClient:  firebaseClient,
// 		library:         lib,
// 	}
// }
