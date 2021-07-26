package mock

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/profileutils"

	"cloud.google.com/go/firestore"
)

// FirestoreClientExtension represents a `firestore.Client` fake
type FirestoreClientExtension struct {
	CollectionFn func(path string) *firestore.CollectionRef
	GetAllFn     func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error)
	CreateFn     func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error)
	UpdateFn     func(ctx context.Context, command *fb.UpdateCommand) error
	DeleteFn     func(ctx context.Context, command *fb.DeleteCommand) error
	GetFn        func(ctx context.Context, query *fb.GetSingleQuery) (*firestore.DocumentSnapshot, error)
}

// Collection ...
func (f *FirestoreClientExtension) Collection(path string) *firestore.CollectionRef {
	return f.CollectionFn(path)
}

// GetAll retrieve a value from the store
func (f *FirestoreClientExtension) GetAll(ctx context.Context, getQuery *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
	return f.GetAllFn(ctx, getQuery)
}

// Create persists data to a firestore collection
func (f *FirestoreClientExtension) Create(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
	return f.CreateFn(ctx, command)
}

// Update updates data to a firestore collection
func (f *FirestoreClientExtension) Update(ctx context.Context, command *fb.UpdateCommand) error {
	return f.UpdateFn(ctx, command)
}

// Delete deletes data to a firestore collection
func (f *FirestoreClientExtension) Delete(ctx context.Context, command *fb.DeleteCommand) error {
	return f.DeleteFn(ctx, command)
}

// Get retrieves data to a firestore collection
func (f *FirestoreClientExtension) Get(ctx context.Context, query *fb.GetSingleQuery) (*firestore.DocumentSnapshot, error) {
	return f.GetFn(ctx, query)
}

// FirebaseClientExtension represents `auth.Client` fake
type FirebaseClientExtension struct {
	GetUserByPhoneNumberFn func(ctx context.Context, phone string) (*auth.UserRecord, error)
	CreateUserFn           func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error)
	DeleteUserFn           func(ctx context.Context, uid string) error
	GetUserProfileByIDFn   func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error)
}

// GetUserByPhoneNumber ...
func (f *FirebaseClientExtension) GetUserByPhoneNumber(ctx context.Context, phone string) (*auth.UserRecord, error) {
	return f.GetUserByPhoneNumberFn(ctx, phone)
}

// CreateUser ...
func (f *FirebaseClientExtension) CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	return f.CreateUserFn(ctx, user)
}

// DeleteUser ...
func (f *FirebaseClientExtension) DeleteUser(ctx context.Context, uid string) error {
	return f.DeleteUserFn(ctx, uid)
}

// GetUserProfileByID ...
func (f *FirebaseClientExtension) GetUserProfileByID(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
	return f.GetUserProfileByIDFn(ctx, id, suspended)
}
