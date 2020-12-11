package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/base"
)

const (
	supplierAPIPath        = "/api/business_partners/suppliers/"
	supplierCollectionName = "suppliers"
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

// AddSupplier makes a call to our own ERP and creates a supplier account for the pro users based
// on their correct partner types that is used for transacting on Be.Well
func (s Service) AddSupplier(
	ctx context.Context,
	uid *string,
	name string,
	partnerType PartnerType,
) (*Supplier, error) {
	s.checkPreconditions()

	userUID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", userUID)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("uid %s has more than one supplier records (it has %d)", userUID, len(docs))
		}
	}
	if len(docs) == 0 {
		currency, err := base.FetchDefaultCurrency(s.erpClient)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
		}

		validPartnerType := partnerType.IsValid()
		if !validPartnerType {
			return nil, fmt.Errorf("%v is not an allowed partner type choice", partnerType.String())
		}

		payload := map[string]interface{}{
			"active":        active,
			"partner_name":  name,
			"country":       country,
			"currency":      *currency.ID,
			"is_supplier":   isSupplier,
			"supplier_type": partnerType,
		}

		content, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
		}
		newSupplier := Supplier{
			UserProfile: profile,
			PartnerType: partnerType,
		}

		if err := base.ReadRequestToTarget(s.erpClient, "POST", supplierAPIPath, "", content, &newSupplier); err != nil {
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

	dsnap := docs[0]
	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %w", err)
	}

	return supplier, nil
}

// FindSupplier fetches a supplier by their UID
func (s Service) FindSupplier(ctx context.Context, uid string) (*Supplier, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx,
		uid,
		s.GetSupplierCollectionName(),
		"userprofile.verifiedIdentifiers",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retreive doc snapshot by uid: %v", err)
	}

	if dsnap == nil {
		return nil, fmt.Errorf("a user with the UID %s does not have a supplier's account", uid)
	}

	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %v", err)
	}

	return supplier, nil
}

// AddSupplierKyc persists a supplier KYC information to firestore
func (s Service) AddSupplierKyc(
	ctx context.Context,
	input SupplierKYCInput) (*SupplierKYC, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}
	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
	}
	if dsnap == nil {
		return nil, fmt.Errorf("the supplier does not exist in our records")
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
	return &supplierKYC, nil
}

// SuspendSupplier flips the active boolean on the erp partner from true to false
// consequently logically deleting the account
func (s Service) SuspendSupplier(ctx context.Context, uid string) (bool, error) {
	s.checkPreconditions()

	err := s.DeleteUser(ctx, uid)
	if err != nil {
		return false, fmt.Errorf("error deleting user: %v", err)
	}

	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	if len(docs) == 0 {
		return false, nil
	}

	dsnap := docs[0]
	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return false, fmt.Errorf("unable to read supplier: %w", err)
	}

	payload := map[string]interface{}{
		"active": false,
	}

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return false, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}

	supplierPath := fmt.Sprintf("%s%s/", customerAPIPath, supplier.SupplierID)
	if err := base.ReadRequestToTarget(s.erpClient, "PATCH", supplierPath, "", content, &supplier); err != nil {
		return false, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	if err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier,
	); err != nil {
		return false, fmt.Errorf("unable to update supplier: %v", err)
	}
	return true, nil
}

// SetUpSupplier performs initial account set up during onboarding
func (s Service) SetUpSupplier(ctx context.Context, input SupplierAccountInput) (*Supplier, error) {
	s.checkPreconditions()

	validAccountType := input.AccountType.IsValid()
	if !validAccountType {
		return nil, fmt.Errorf("%v is not an allowed AccountType choice", input.AccountType.String())
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %w", err)
	}
	supplier := &Supplier{}

	if dsnap != nil {
		err = dsnap.DataTo(supplier)
		if err != nil {
			return nil, fmt.Errorf("unable to read supplier: %v", err)
		}
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	supplier.UserProfile = profile
	supplier.AccountType = input.AccountType
	supplier.UnderOrganization = input.UnderOrganization

	if input.IsOrganizationVerified != nil {
		// TODO:
		// verified := *input.IsOrganizationVerified
		// if !verified {
		// 	// Successful/unsuccessful authserver login
		// 	// if false requires support follow up
		// 	// send an email, send a nudge...??
		// }
		supplier.IsOrganizationVerified = *input.IsOrganizationVerified
	}

	if input.SladeCode == nil {
		// TODO:
		// User without slade code requires approval i.e creation of an account
		// send an email to support for follow up, send a nudge ..?
		// full KYC based on account type (individual or organisation)
	} else {
		supplier.SladeCode = *input.SladeCode
	}

	if input.ParentOrganizationID != nil {
		supplier.ParentOrganizationID = *input.ParentOrganizationID
		supplier.HasBranches = true
	}

	if input.Location != nil {
		supplier.Location.ID = input.Location.ID
		supplier.Location.Name = input.Location.Name
		supplier.Location.BranchSladeCode = input.Location.BranchSladeCode
	}

	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}

	return supplier, nil
}
