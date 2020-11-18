package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

// AddSupplier creates a supplier on the ERP when a user signs up in our Be.Well Pro
func (s Service) AddSupplier(ctx context.Context, uid *string) (*Supplier, error) {
	s.checkPreconditions()

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	user, userErr := s.firebaseAuth.GetUser(ctx, profile.UID)

	if userErr != nil {
		return nil, fmt.Errorf("unable to get Firebase user with UID %s: %w", profile.UID, userErr)
	}

	if user.DisplayName == "" {
		return nil, fmt.Errorf("user does not have a DisplayName")
	}

	currency, err := base.FetchDefaultCurrency(s.client)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
	}
	payload := map[string]interface{}{
		"active":        active,
		"partner_name":  user.DisplayName,
		"country":       country,
		"currency":      *currency.ID,
		"is_supplier":   isSupplier,
		"supplier_type": supplierType,
	}

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}
	newSupplier := Supplier{
		UserProfile: *profile,
	}

	if err := base.ReadRequestToTarget(s.client, "POST", supplierAPIPath, "", content, &newSupplier); err != nil {
		return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	if err := s.SaveSupplierToFireStore(newSupplier); err != nil {
		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}

	profile.HasSupplierAccount = true
	profileDsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve firebase user profile: %v", err)
	}

	if err := base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), profileDsnap.Ref.ID, profile,
	); err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}

	return &newSupplier, nil
}

// FindSupplier fetches a supplier by their UID
func (s Service) FindSupplier(ctx context.Context, uid string) (*Supplier, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetSupplierCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %v", err)
	}

	if dsnap == nil {
		return s.AddSupplier(ctx, &uid)
	}
	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %v", err)
	}

	return supplier, nil
}

// FindSupplierByUIDHandler is a used for inter service communication to return details about a supplier
func FindSupplierByUIDHandler(ctx context.Context, service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := ValidateUID(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		var supplier *Supplier

		if s.Token != nil {
			newContext := context.WithValue(ctx, base.AuthTokenContextKey, s.Token)
			supplier, err = service.FindSupplier(newContext, *s.UID)
		} else {
			supplier, err = service.FindSupplier(ctx, *s.UID)
		}

		if supplier == nil || err != nil {
			base.ReportErr(w, err, http.StatusNotFound)
			return
		}

		supplierResponse := SupplierResponse{
			SupplierID: supplier.SupplierID,
			PayablesAccount: PayablesAccount{
				ID:          supplier.PayablesAccount.ID,
				Name:        supplier.PayablesAccount.Name,
				IsActive:    supplier.PayablesAccount.IsActive,
				Number:      supplier.PayablesAccount.Number,
				Tag:         supplier.PayablesAccount.Tag,
				Description: supplier.PayablesAccount.Description,
			},
			Profile: BioData{
				UID:        supplier.UserProfile.UID,
				Name:       supplier.UserProfile.Name,
				Gender:     supplier.UserProfile.Gender,
				Msisdns:    supplier.UserProfile.Msisdns,
				Emails:     supplier.UserProfile.Emails,
				PushTokens: supplier.UserProfile.PushTokens,
				Bio:        supplier.UserProfile.Bio,
			},
		}

		base.WriteJSONResponse(w, supplierResponse, http.StatusOK)
	}
}
