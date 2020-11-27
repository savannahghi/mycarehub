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
func (s Service) AddSupplier(ctx context.Context, uid *string, name string) (*Supplier, error) {
	s.checkPreconditions()

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	currency, err := base.FetchDefaultCurrency(s.client)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
	}
	payload := map[string]interface{}{
		"active":        active,
		"partner_name":  name,
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
		UserProfile: profile,
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
		// If supplier is not found,
		// and the user exists
		// then create one using their UID
		user, userErr := s.firebaseAuth.GetUser(ctx, uid)
		if userErr != nil {
			return nil, fmt.Errorf("unable to get Firebase user with UID %s: %w", uid, userErr)
		}
		return s.AddSupplier(ctx, &uid, user.UID)
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
			SupplierID:      supplier.SupplierID,
			PayablesAccount: *supplier.PayablesAccount,
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

// AddSupplierKyc persists a supplier KYC information to firestore
func (s Service) AddSupplierKyc(
	ctx context.Context,
	input SupplierKYCInput) (*SupplierKYC, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch user profile: %v", err)
	}
	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, profile.UID, s.GetSupplierCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
	}
	if dsnap == nil {
		return nil, fmt.Errorf("the supplier does not exist in out records")
	}
	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier data: %v", err)
	}

	supplier.SupplierKYC.AccountType = input.AccountType
	supplier.SupplierKYC.IdentificationDocType = input.IdentificationDocType
	supplier.SupplierKYC.IdentificationDocNumber = input.IdentificationDocNumber
	supplier.SupplierKYC.IdentificationDocPhotoBase64 = input.IdentificationDocPhotoBase64
	supplier.SupplierKYC.IdentificationDocPhotoContentType = input.IdentificationDocPhotoContentType
	supplier.SupplierKYC.License = input.License
	supplier.SupplierKYC.Cadre = input.Cadre
	supplier.SupplierKYC.Profession = input.Profession
	supplier.SupplierKYC.KraPin = input.KraPin
	supplier.SupplierKYC.KraPINDocPhoto = input.KraPINDocPhoto
	supplier.SupplierKYC.BusinessNumber = input.BusinessNumber
	supplier.SupplierKYC.BusinessNumberDocPhotoBase64 = input.BusinessNumberDocPhotoBase64
	supplier.SupplierKYC.BusinessNumberDocPhotoContentType = input.BusinessNumberDocPhotoContentType

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}
	supplierKYC := supplier.SupplierKYC
	return supplierKYC, nil
}
