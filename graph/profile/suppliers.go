package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

const (
	supplierAPIPath        = "/api/business_partners/suppliers/"
	supplierCollectionName = "suppliers"
	isSupplier             = true
	futureHours            = 878400
)

const (
	// engagement ISC paths
	publishNudge = "feed/%s/PRO/false/nudges/"
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

// AddPartnerType create the initial supplier record
func (s Service) AddPartnerType(ctx context.Context, name *string,
	partnerType *PartnerType) (bool, error) {

	s.checkPreconditions()

	if name == nil || partnerType == nil || *name == " " || !partnerType.IsValid() {
		return false, fmt.Errorf("expected `name` to be defined and `partnerType` to be valid")
	}

	if *partnerType == PartnerTypeConsumer {
		return false, fmt.Errorf("invalid `partnerType`. cannot use CONSUMER in this context")
	}

	userUID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &userUID)
	if err != nil {
		return false, fmt.Errorf("unable to read user profile: %w", err)
	}

	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", userUID)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	// if record length is equal to on 1, update otherwise create
	if len(docs) == 1 {
		// update
		supplier := &Supplier{}
		err = docs[0].DataTo(supplier)
		if err != nil {
			return false, fmt.Errorf("unable to read supplier: %v", err)
		}

		supplier.UserProfile.Name = name
		supplier.PartnerType = *partnerType
		supplier.PartnerSetupComplete = true

		if err := s.SaveSupplierToFireStore(*supplier); err != nil {
			return false, fmt.Errorf("unable to add supplier to firestore: %v", err)
		}

		return true, nil
	}

	// create new record
	profile.Name = name
	newSupplier := Supplier{
		UserProfile:          profile,
		PartnerType:          *partnerType,
		PartnerSetupComplete: true,
	}

	if err := s.SaveSupplierToFireStore(newSupplier); err != nil {
		return false, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}

	return true, nil
}

// AddSupplier makes a call to our own ERP and creates a supplier account for the pro users based
// on their correct partner types that is used for transacting on Be.Well
func (s Service) AddSupplier(
	ctx context.Context,
	name string,
	partnerType PartnerType,
) (*Supplier, error) {
	s.checkPreconditions()

	userUID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &userUID)
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
		return nil, fmt.Errorf("expected user to have a supplier account : %w", err)
	}

	dsnap := docs[0]
	supplier := &Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %w", err)
	}

	currency, err := base.FetchDefaultCurrency(s.erpClient)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
	}

	validPartnerType := partnerType.IsValid()
	if !validPartnerType {
		return nil, fmt.Errorf("%v is not an valid partner type choice", partnerType.String())
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

	if err := base.ReadRequestToTarget(s.erpClient, "POST", supplierAPIPath, "", content, &Supplier{
		UserProfile: profile,
		PartnerType: partnerType,
	}); err != nil {
		return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	supplier.Active = true

	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
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

	supplierPath := fmt.Sprintf("%s%s", customerAPIPath, supplier.SupplierID)
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
func (s Service) SetUpSupplier(ctx context.Context, accountType AccountType) (*Supplier, error) {
	s.checkPreconditions()

	validAccountType := accountType.IsValid()
	if !validAccountType {
		return nil, fmt.Errorf("%v is not an allowed AccountType choice", accountType.String())
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

	if dsnap == nil {
		return nil, fmt.Errorf("cannot find supplier record")
	}

	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	supplier.UserProfile = profile
	supplier.AccountType = accountType
	supplier.UnderOrganization = false
	supplier.IsOrganizationVerified = false
	supplier.HasBranches = false

	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}

	go func() {
		op := func() error {
			return s.PublishKYCNudge(uid, &supplier.PartnerType, &supplier.AccountType)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	return supplier, nil
}

// EDIUserLogin used to login a user to EDI and return their EDI profile
func EDIUserLogin(username, password string) (*base.EDIUserProfile, error) {

	if username == "" || password == "" {
		return nil, fmt.Errorf("invalid credentials, expected a username AND password")
	}

	ediClient, err := base.LoginClient(username, password)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := base.FetchUserProfile(ediClient)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// SupplierEDILogin it used to instantiate as call when setting up a supplier's account's who
// has an affliation to a provider with the slade ecosystem. The logic is as follows;
// 1 . login to the relevant edi to assert the user has an account
// 2 . fetch the branches of the provider given the slade code which we have
// 3 . update the user's supplier record
// 4. return the list of branches to the frontend so that a default location can be set
func (s Service) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*BranchConnection, error) {
	s.checkPreconditions()
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %w", err)
	}

	supplier := &Supplier{}

	if dsnap == nil {
		return nil, fmt.Errorf("cannot find supplier record")
	}

	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	supplier.UserProfile = profile
	supplier.AccountType = AccountTypeIndividual
	supplier.UnderOrganization = true

	ediUserProfile, err := func(sladeCode string) (*base.EDIUserProfile, error) {
		var ediUserProfile *base.EDIUserProfile
		var err error

		switch sladeCode {
		case "1":
			// login to core
			//TODO(calvine) add login to core
			//TODO(calvine) if login passes, update userprofile to admin and assign admin permissions
		default:
			//Login to portal edi
			ediUserProfile, err = EDIUserLogin(username, password)
			if err != nil {
				supplier.IsOrganizationVerified = false
				return nil, fmt.Errorf("cannot get edi user profile: %w", err)
			}

			if ediUserProfile == nil {
				return nil, fmt.Errorf("edi user profile not found")
			}

		}
		return ediUserProfile, nil
	}(sladeCode)

	if err != nil {
		return nil, err
	}

	// The slade code comes in the form 'PRO-1234' or 'BRA-PRO-1234-1'
	// we split it to get the interger part of the slade code.
	var orgSladeCode string
	if strings.HasPrefix(sladeCode, "BRA-") {
		orgSladeCode = strings.Split(sladeCode, "-")[2]
	} else {
		orgSladeCode = strings.Split(sladeCode, "-")[1]
	}

	// verify slade code.
	if ediUserProfile.BusinessPartner != orgSladeCode {
		supplier.IsOrganizationVerified = false
		return nil, fmt.Errorf("invalid slade code for selected provider: %v, got: %v", sladeCode, ediUserProfile.BusinessPartner)
	}

	supplier.EDIUserProfile = ediUserProfile
	supplier.IsOrganizationVerified = true
	supplier.SladeCode = sladeCode

	filter := []*BusinessPartnerFilterInput{
		{
			SladeCode: &sladeCode,
		},
	}

	partner, err := s.FindProvider(ctx, nil, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	var businessPartner BusinessPartner

	if len(partner.Edges) != 1 {
		return nil, fmt.Errorf("expected one business partner, found: %v", len(partner.Edges))
	}

	businessPartner = *partner.Edges[0].Node
	var brFilter []*BranchFilterInput

	go func() {
		op := func() error {
			return s.PublishKYCNudge(uid, &supplier.PartnerType, &supplier.AccountType)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	if businessPartner.Parent != nil {
		supplier.HasBranches = true
		supplier.ParentOrganizationID = *businessPartner.Parent
		filter := &BranchFilterInput{
			ParentOrganizationID: businessPartner.Parent,
		}

		brFilter = append(brFilter, filter)
		if err := s.SaveSupplierToFireStore(*supplier); err != nil {
			return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
		}

		return s.FindBranch(ctx, nil, brFilter, nil)
	}
	loc := Location{
		ID:   businessPartner.ID,
		Name: businessPartner.Name,
	}
	supplier.Location = &loc

	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
	}
	pageInfo := &base.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     nil,
		EndCursor:       nil,
	}

	return &BranchConnection{PageInfo: pageInfo}, nil
}

// SupplierSetDefaultLocation updates the default location ot the supplier by the given location id
func (s Service) SupplierSetDefaultLocation(ctx context.Context, locationID string) (bool, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	// fetch the supplier records
	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("unable to fetch supplier record: %w", err)
	}
	if len(docs) == 0 {
		return false, fmt.Errorf("unable to find supplier record: %w", err)
	}

	dsnap := docs[0]
	sup := &Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return false, fmt.Errorf("unable to read supplier: %w", err)
	}

	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
	filter := []*BranchFilterInput{
		{
			SladeCode:            &sup.SladeCode,
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return false, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	branch := func(brs *BranchConnection, location string) *BranchEdge {
		for _, b := range brs.Edges {
			if b.Node.ID == location {
				return b
			}
		}
		return nil
	}(brs, locationID)

	if branch != nil {
		loc := Location{
			ID:              branch.Node.ID,
			Name:            branch.Node.Name,
			BranchSladeCode: &branch.Node.BranchSladeCode,
		}
		sup.Location = &loc

		// update the supplier record with new location
		if err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, sup); err != nil {
			return false, fmt.Errorf("unable to update supplier default location: %v", err)
		}
	}

	return false, fmt.Errorf("unable to get location of id %v : %v", locationID, err)
}

// FetchSupplierAllowedLocations retrieves all the locations that the user in context can work on.
func (s *Service) FetchSupplierAllowedLocations(ctx context.Context) (*BranchConnection, error) {

	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	// fetch the supplier records
	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch supplier record: %w", err)
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("unable to find supplier record: %w", err)
	}

	dsnap := docs[0]
	sup := &Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, fmt.Errorf("unable to read supplier record: %w", err)
	}

	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
	filter := []*BranchFilterInput{
		{
			SladeCode:            &sup.SladeCode,
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	return brs, nil
}

// PublishKYCNudge pushes a kyc nudge to the user feed
func (s *Service) PublishKYCNudge(uid string, partner *PartnerType, account *AccountType) error {

	s.checkPreconditions()

	if partner == nil || !partner.IsValid() {
		return fmt.Errorf("expected `partner` to be defined and to be valid")
	}

	if *partner == PartnerTypeConsumer {
		return fmt.Errorf("invalid `partner`. cannot use CONSUMER in this context")
	}

	if !account.IsValid() {
		return fmt.Errorf("provided `account` is not valid")
	}

	payload := base.Nudge{
		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
		Visibility:     "SHOW",
		Status:         "PENDING",
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
		Text:           "Fill in your Be.Well business KYC in order to start transacting",
		Links: []base.Link{
			{
				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC",
				Description: fmt.Sprintf("KYC for %v", partner.String()),
				Thumbnail:   base.LogoURL,
			},
		},
		Actions: []base.Action{
			{
				ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
				SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
				Name:           strings.ToUpper(fmt.Sprintf("COMPLETE_%v_%v_KYC", account.String(), partner.String())),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
				Icon: base.Link{
					ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
					Description: "Fill in your Be.Well business KYC in order to start transacting",
					Thumbnail:   base.LogoURL,
				},
			},
		},
		Users:                []string{uid},
		Groups:               []string{uid},
		NotificationChannels: []base.Channel{base.ChannelEmail, base.ChannelFcm},
	}

	resp, err := s.engagement.MakeRequest("POST", fmt.Sprintf(publishNudge, uid), payload)
	if err != nil {
		return fmt.Errorf("unable to publish kyc nudge : %v", err)
	}

	//TODO(dexter) to be removed. Just here for debug
	res, _ := httputil.DumpResponse(resp, true)
	log.Println(string(res))

	if resp.StatusCode != http.StatusOK {
		// stage the nudge
		stage := func(pl base.Nudge) error {
			k, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("cannot marshal payload to json")
			}

			var kMap map[string]interface{}
			err = json.Unmarshal(k, &kMap)
			if err != nil {
				return fmt.Errorf("cannot unmarshal payload from json")
			}

			if err := s.SaveProfileNudge(kMap); err != nil {
				logrus.Errorf("failed to stage nudge : %v", err)
			}
			return nil

		}(payload)

		if err := stage; err != nil {
			logrus.Errorf("failed to stage nudge : %v", err)
		}
		return fmt.Errorf("unable to publish kyc nudge. unexpected status code  %v", resp.StatusCode)
	}

	return nil
}

// StageKYCProcessingRequest saves kyc processing requests
func (s *Service) StageKYCProcessingRequest(sup *Supplier) error {
	r := KYCRequest{
		ID:                  uuid.New().String(),
		ReqPartnerType:      sup.PartnerType,
		ReqOrganizationType: OrganizationType(sup.AccountType),
		ReqRaw:              sup.SupplierKYC,
		Proceseed:           false,
		SupplierRecord:      sup,
		Status:              KYCProcessStatusPending,
	}

	_, err := base.SaveDataToFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), r)
	if err != nil {
		return fmt.Errorf("unable to save kyc processing request: %w", err)
	}
	return nil
}

// AddIndividualRiderKyc adds KYC for an individual rider
func (s *Service) AddIndividualRiderKyc(ctx context.Context, input IndividualRider) (*IndividualRider, error) {

	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := IndividualRider{
		IdentificationDoc: Identification{
			IdentificationDocType:           input.IdentificationDoc.IdentificationDocType,
			IdentificationDocNumber:         input.IdentificationDoc.IdentificationDocNumber,
			IdentificationDocNumberUploadID: input.IdentificationDoc.IdentificationDocNumberUploadID,
		},
		KRAPIN:                         input.KRAPIN,
		KRAPINUploadID:                 input.KRAPINUploadID,
		DrivingLicenseID:               input.DrivingLicenseID,
		DrivingLicenseUploadID:         input.DrivingLicenseUploadID,
		CertificateGoodConductUploadID: input.CertificateGoodConductUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationRiderKyc adds KYC for an organization rider
func (s *Service) AddOrganizationRiderKyc(ctx context.Context, input OrganizationRider) (*OrganizationRider, error) {

	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationRider{
		OrganizationTypeName:               input.OrganizationTypeName,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,

		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddIndividualPractitionerKyc adds KYC for an individual pratitioner
func (s *Service) AddIndividualPractitionerKyc(ctx context.Context, input IndividualPractitioner) (*IndividualPractitioner, error) {
	s.checkPreconditions()

	for _, p := range input.PracticeServices {
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
		}
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := IndividualPractitioner{

		IdentificationDoc: func(p Identification) Identification {
			return Identification(p)
		}(input.IdentificationDoc),

		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		RegistrationNumber:          input.RegistrationNumber,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
		PracticeServices:            input.PracticeServices,
		Cadre:                       input.Cadre,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationPractitionerKyc adds KYC for an organization pratitioner
func (s *Service) AddOrganizationPractitionerKyc(ctx context.Context, input OrganizationPractitioner) (*OrganizationPractitioner, error) {

	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	for _, p := range input.PracticeServices {
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
		}
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationPractitioner{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		RegistrationNumber:                 input.RegistrationNumber,
		PracticeLicenseID:                  input.PracticeLicenseID,
		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
		PracticeServices:                   input.PracticeServices,
		Cadre:                              input.Cadre,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationProviderKyc adds KYC for an organization provider
func (s *Service) AddOrganizationProviderKyc(ctx context.Context, input OrganizationProvider) (*OrganizationProvider, error) {

	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	for _, p := range input.PracticeServices {
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
		}
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationProvider{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		RegistrationNumber:                 input.RegistrationNumber,
		PracticeLicenseID:                  input.PracticeLicenseID,
		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
		PracticeServices:                   input.PracticeServices,
		Cadre:                              input.Cadre,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddIndividualPharmaceuticalKyc adds KYC for an individual Pharmaceutical kyc
func (s *Service) AddIndividualPharmaceuticalKyc(ctx context.Context, input IndividualPharmaceutical) (*IndividualPharmaceutical, error) {

	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := IndividualPharmaceutical{
		IdentificationDoc: func(p Identification) Identification {
			return Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		RegistrationNumber:          input.RegistrationNumber,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationPharmaceuticalKyc adds KYC for a pharmacy organization
func (s *Service) AddOrganizationPharmaceuticalKyc(ctx context.Context, input OrganizationPharmaceutical) (*OrganizationPharmaceutical, error) {
	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationPharmaceutical{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddIndividualCoachKyc adds KYC for an individual coach
func (s *Service) AddIndividualCoachKyc(ctx context.Context, input IndividualCoach) (*IndividualCoach, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := IndividualCoach{
		IdentificationDoc: func(p Identification) Identification {
			return Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationCoachKyc adds KYC for an organization coach
func (s *Service) AddOrganizationCoachKyc(ctx context.Context, input OrganizationCoach) (*OrganizationCoach, error) {
	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationCoach{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddIndividualNutritionKyc adds KYC for an individual nutritionist
func (s *Service) AddIndividualNutritionKyc(ctx context.Context, input IndividualNutrition) (*IndividualNutrition, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := IndividualNutrition{
		IdentificationDoc: func(p Identification) Identification {
			return Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// AddOrganizationNutritionKyc adds kyc for a nutritionist organisation
func (s *Service) AddOrganizationNutritionKyc(ctx context.Context, input OrganizationNutrition) (*OrganizationNutrition, error) {
	s.checkPreconditions()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
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

	kyc := OrganizationNutrition{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []Identification) []Identification {
			pl := []Identification{}
			for _, i := range p {
				pl = append(pl, Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	k, err := json.Marshal(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}
	var kycAsMap map[string]interface{}
	err = json.Unmarshal(k, &kycAsMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	supplier.SupplierKYC = kycAsMap
	supplier.KYCSubmitted = true

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
	if err != nil {
		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
	}

	if err := s.StageKYCProcessingRequest(supplier); err != nil {
		logrus.Errorf("unable to stage kyc processing request: %v", err)
	}

	return &kyc, nil
}

// FetchKYCProcessingRequests fetches a list of all unprocessed kyc approval requests
func (s *Service) FetchKYCProcessingRequests(ctx context.Context) ([]*KYCRequest, error) {
	collection := s.firestoreClient.Collection(s.GetKCYProcessCollectionName())
	query := collection.Where("proceseed", "==", false)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch kyc request documents: %v", err)
	}

	res := []*KYCRequest{}

	for _, doc := range docs {
		req := &KYCRequest{}
		err = doc.DataTo(req)
		if err != nil {
			return nil, fmt.Errorf("unable to read supplier: %w", err)
		}
		res = append(res, req)
	}
	// todo : amend this to include documents/attachments urls. This will be used in the frontend to pull the images for viewing e.g ID, DL

	return res, nil
}

// ProcessKYCRequest transitions a kyc request to a given state
func (s *Service) ProcessKYCRequest(ctx context.Context, id string, status KYCProcessStatus, rejectionReason *string) (bool, error) {
	collection := s.firestoreClient.Collection(s.GetKCYProcessCollectionName())
	query := collection.Where("id", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, fmt.Errorf("unable to fetch kyc request documents: %v", err)
	}

	doc := docs[0]
	req := &KYCRequest{}
	err = doc.DataTo(req)
	if err != nil {
		return false, fmt.Errorf("unable to read supplier: %w", err)
	}

	req.Status = status
	req.Proceseed = true
	req.RejectionReason = rejectionReason

	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), doc.Ref.ID, req)
	if err != nil {
		return false, fmt.Errorf("unable to update KYC request record: %v", err)
	}

	var email string

	switch status {
	case KYCProcessStatusApproved:
		// create supplier erp account
		if _, err := s.AddSupplier(ctx, *req.SupplierRecord.UserProfile.Name, req.ReqPartnerType); err != nil {
			return false, fmt.Errorf("unable to create erp supplier account: %v", err)
		}

		email = generateProcessKYCApprovalEmailTemplate()
		// todo: send text to the supplier on approval
	case KYCProcessStatusRejected:
		email = generateProcessKYCRejectionEmailTemplate()
		// todo: send text to the supplier on rejection

	}

	for _, supplierEmail := range req.SupplierRecord.UserProfile.Emails {
		err = s.SendKYCEmail(ctx, email, supplierEmail)
		if err != nil {
			return false, fmt.Errorf("unable to send KYC processing email: %w", err)
		}
	}

	return true, nil

}

// SendKYCEmail will send a KYC processing request email to the supplier
func (s Service) SendKYCEmail(ctx context.Context, text, emailaddress string) error {
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}

	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": emailSignupSubject,
	}

	resp, err := s.Mailgun.MakeRequest(http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send KYC email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send KYC email : %w, with status code %v", err, resp.StatusCode)
	}

	return nil
}
