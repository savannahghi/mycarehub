package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
)

const (
	supplierAPIPath        = "/api/business_partners/suppliers/"
	supplierCollectionName = "suppliers"
	supplierType           = "PHARMACEUTICAL" // TODO
	isSupplier             = true
)

// SaveSupplierToFireStore persists supplier data to firestore
func (s Service) SaveSupplierToFireStore(supplier Supplier) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetSupplierCollectionName()).Add(ctx, supplier)
	return err
}

// GetSupplierCollectionName creates a suffixed supplier collection name
func (s Service) GetSupplierCollectionName() string {
	suffixed := base.SuffixCollection(supplierCollectionName)
	return suffixed
}

// RetrieveSupplierFirebaseDocSnapshotByUID retrieves a raw supplier Firebase doc snapshot
func (s Service) RetrieveSupplierFirebaseDocSnapshotByUID(
	ctx context.Context, uid string) (*firestore.DocumentSnapshot, bool, error) {
	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("UID", "==", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, false, fmt.Errorf("unable to retrieve supplier snapshot: %v", err)
	}
	if len(docs) > 1 {
		log.Printf("user %s has > 1 supplier profile (they have %d)", uid, len(docs))
	}
	if len(docs) == 0 {
		return nil, false, nil
	}
	dsnap := docs[0]
	return dsnap, true, nil
}

// AddSupplier creates a supplier on the ERP when a user signs up in our Be.Well Pro
func (s Service) AddSupplier(ctx context.Context) (*Supplier, error) {
	s.checkPreconditions()

	supplier := &Supplier{}
	profile, profileErr := s.UserProfile(ctx)
	if profileErr != nil {
		return nil, profileErr
	}

	dsnap, exists, err := s.RetrieveSupplierFirebaseDocSnapshotByUID(ctx, profile.UID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve supplier: %v", err)
	}

	if exists {
		err = dsnap.DataTo(supplier)
		if err != nil {
			return nil, fmt.Errorf("unable to read supplier data: %v", err)
		}
		return supplier, nil
	}

	fireBaseClient, clientErr := base.GetFirebaseAuthClient(ctx)
	if clientErr != nil {
		return nil, fmt.Errorf("unable to initialize Firebase auth client: %w", clientErr)
	}
	user, userErr := fireBaseClient.GetUser(ctx, profile.UID)
	if userErr != nil {
		return nil, fmt.Errorf("unable to get Firebase user with UID %s: %w", profile.UID, userErr)
	}

	if user.DisplayName == "" {
		return nil, fmt.Errorf("user does not have a DisplayName")
	}

	currency := base.MustGetEnvVar(erpCurrencyEnvName)
	payload := map[string]interface{}{
		"active":        active,
		"partner_name":  user.DisplayName,
		"country":       country,
		"currency":      currency,
		"is_supplier":   isSupplier,
		"supplier_type": supplierType,
	}

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}
	newSupplier := Supplier{
		UID:         profile.UID,
		UserProfile: profile,
	}

	err = base.ReadRequestToTarget(s.client, "POST", supplierAPIPath, "", content, &newSupplier)
	if err != nil {
		return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	err = s.SaveSupplierToFireStore(newSupplier)
	if err != nil {
		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}

	profile.HasSupplierAccount = true
	profileDsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve firebase user profile: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), profileDsnap.Ref.ID, profile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}

	return &newSupplier, nil
}
