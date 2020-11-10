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
func (s Service) AddSupplier(ctx context.Context) (*Supplier, error) {
	s.checkPreconditions()

	profile, profileErr := s.UserProfile(ctx)
	if profileErr != nil {
		return nil, profileErr
	}

	supplier, err := s.FindSupplier(ctx, profile.UID)
	if err != nil {
		return nil, fmt.Errorf("unable to get supplier: %v", err)
	}

	if supplier != nil {
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
		UserProfile: *profile,
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

// FindSupplier fetches a supplier by their UID
func (s Service) FindSupplier(ctx context.Context, uid string) (*Supplier, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetSupplierCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %v", err)
	}

	if dsnap == nil {
		return nil, nil
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
		supplierUID, err := ValidateUID(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		supplier, err := service.FindSupplier(ctx, supplierUID)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if supplier == nil {
			base.WriteJSONResponse(w, StatusResponse{Status: "not found"}, http.StatusNotFound)
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
		}

		base.WriteJSONResponse(w, supplierResponse, http.StatusOK)
	}
}
