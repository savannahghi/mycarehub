package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/profile/domain"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
)

// SupplierUseCases represent the business logic required for management of suppliers
type SupplierUseCases interface {
	AddPartnerType(ctx context.Context, name *string, partnerType *domain.PartnerType) (bool, error)
	AddSupplier(ctx context.Context, uid *string, name string) (*domain.Supplier, error)
	FindSupplier(ctx context.Context, uid string) (*domain.Supplier, error)
	// AddSupplierKyc(ctx context.Context, input domain.SupplierKYCInput) (*domain.SupplierKYC, error)
	SetUpSupplier(ctx context.Context, accountType domain.AccountType) (*domain.Supplier, error)
	DeleteUser(ctx context.Context, uid string) error
	SuspendSupplier(ctx context.Context, uid string) (bool, error)
	EDIUserLogin(username, password string) (*base.EDIUserProfile, error)
	CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error)
	SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*domain.BranchConnection, error)
	SupplierSetDefaultLocation(ctx context.Context, locationID string) (bool, error)
	FetchSupplierAllowedLocations(ctx context.Context) (*domain.BranchConnection, error)
	SaveProfileNudge(nudge map[string]interface{}) error
	PublishKYCNudge(uid string, partner *domain.PartnerType, account *domain.AccountType) error
	PublishKYCFeedItem(ctx context.Context, uids ...string) error
	StageKYCProcessingRequest(sup *domain.Supplier) error
	AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error)
	AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error)
	AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error)
	AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error)
	AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error)
	AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error)
	AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error)
	AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error)
	AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error)
	AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error)
	AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error)
	SaveKYCResponse(ctx context.Context, kycJSON []byte, supplier *domain.Supplier, dsnap *firestore.DocumentSnapshot) error
	FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error)
	ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error)
	SendKYCEmail(ctx context.Context, text, emailaddress string) error
}

// // AddPartnerType create the initial supplier record
// func (s Service) AddPartnerType(ctx context.Context, name *string,
// 	partnerType *PartnerType) (bool, error) {

// 	s.checkPreconditions()

// 	if name == nil || partnerType == nil || *name == " " || !partnerType.IsValid() {
// 		return false, fmt.Errorf("expected `name` to be defined and `partnerType` to be valid")
// 	}

// 	if *partnerType == PartnerTypeConsumer {
// 		return false, fmt.Errorf("invalid `partnerType`. cannot use CONSUMER in this context")
// 	}

// 	userUID, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &userUID)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to read user profile: %w", err)
// 	}

// 	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
// 	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", userUID)

// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, err
// 	}

// 	// if record length is equal to on 1, update otherwise create
// 	if len(docs) == 1 {
// 		// update
// 		supplier := &Supplier{}
// 		err = docs[0].DataTo(supplier)
// 		if err != nil {
// 			return false, fmt.Errorf("unable to read supplier: %v", err)
// 		}

// 		supplier.UserProfile.Name = name
// 		supplier.PartnerType = *partnerType
// 		supplier.PartnerSetupComplete = true

// 		if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 			return false, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 		}

// 		return true, nil
// 	}

// 	// create new record
// 	profile.Name = name
// 	newSupplier := Supplier{
// 		UserProfile:          profile,
// 		PartnerType:          *partnerType,
// 		PartnerSetupComplete: true,
// 	}

// 	if err := s.SaveSupplierToFireStore(newSupplier); err != nil {
// 		return false, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 	}

// 	return true, nil
// }

// // AddSupplier makes a call to our own ERP and creates a supplier account for the pro users based
// // on their correct partner types that is used for transacting on Be.Well
// func (s Service) AddSupplier(
// 	ctx context.Context,
// 	name string,
// 	partnerType PartnerType,
// ) (*Supplier, error) {
// 	s.checkPreconditions()

// 	userUID, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &userUID)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read user profile: %w", err)
// 	}

// 	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
// 	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", userUID)

// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(docs) > 1 {
// 		if base.IsDebug() {
// 			log.Printf("uid %s has more than one supplier records (it has %d)", userUID, len(docs))
// 		}
// 	}

// 	if len(docs) == 0 {
// 		return nil, fmt.Errorf("expected user to have a supplier account : %w", err)
// 	}

// 	dsnap := docs[0]
// 	supplier := &Supplier{}
// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier: %w", err)
// 	}

// 	currency, err := base.FetchDefaultCurrency(s.erpClient)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
// 	}

// 	validPartnerType := partnerType.IsValid()
// 	if !validPartnerType {
// 		return nil, fmt.Errorf("%v is not an valid partner type choice", partnerType.String())
// 	}

// 	payload := map[string]interface{}{
// 		"active":        active,
// 		"partner_name":  name,
// 		"country":       country,
// 		"currency":      *currency.ID,
// 		"is_supplier":   isSupplier,
// 		"supplier_type": partnerType,
// 	}

// 	content, marshalErr := json.Marshal(payload)
// 	if marshalErr != nil {
// 		return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
// 	}

// 	if err := base.ReadRequestToTarget(s.erpClient, "POST", supplierAPIPath, "", content, &Supplier{
// 		UserProfile: profile,
// 		PartnerType: partnerType,
// 	}); err != nil {
// 		return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
// 	}

// 	supplier.Active = true

// 	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 	}

// 	profile.HasSupplierAccount = true
// 	profileDsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve firebase user profile: %v", err)
// 	}

// 	if err := base.UpdateRecordOnFirestore(
// 		s.firestoreClient, s.GetUserProfileCollectionName(), profileDsnap.Ref.ID, profile,
// 	); err != nil {
// 		return nil, fmt.Errorf("unable to update user profile: %v", err)
// 	}

// 	return supplier, nil
// }

// // FindSupplier fetches a supplier by their UID
// func (s Service) FindSupplier(ctx context.Context, uid string) (*Supplier, error) {
// 	s.checkPreconditions()

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
// 		ctx,
// 		uid,
// 		s.GetSupplierCollectionName(),
// 		"userprofile.verifiedIdentifiers",
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"unable to retreive doc snapshot by uid: %v", err)
// 	}

// 	if dsnap == nil {
// 		// create a default supplier account instead of throwing an error
// 		// this is for backwards compatibility for dev accounts that don't have supplier account.
// 		// there is no harm in doing this since the same supplier will be updated
// 		// We default to PartnerTypeProvider since we expect majority of our suppliers to be providers.
// 		// In any case, this can always be changed.

// 		pr, err := s.UserProfile(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to retreive userprofile of logged in user: %v", err)
// 		}
// 		pty := PartnerTypeProvider
// 		if _, err := s.AddPartnerType(ctx, pr.Name, &pty); err != nil {
// 			return nil, fmt.Errorf("unable to create default supplier account : %v", err)
// 		}

// 		return nil, fmt.Errorf("a user with the UID %s does not have a supplier's account", uid)
// 	}

// 	supplier := &Supplier{}
// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier: %v", err)
// 	}

// 	return supplier, nil
// }

// // FindSupplierByUIDHandler is a used for inter service communication to return details about a supplier
// func FindSupplierByUIDHandler(ctx context.Context, service *Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		s, err := ValidateUID(w, r)
// 		if err != nil {
// 			base.ReportErr(w, err, http.StatusBadRequest)
// 			return
// 		}

// 		var supplier *Supplier

// 		if s.Token != nil {
// 			newContext := context.WithValue(ctx, base.AuthTokenContextKey, s.Token)
// 			supplier, err = service.FindSupplier(newContext, *s.UID)
// 		} else {
// 			supplier, err = service.FindSupplier(ctx, *s.UID)
// 		}

// 		if supplier == nil || err != nil {
// 			base.ReportErr(w, err, http.StatusNotFound)
// 			return
// 		}

// 		supplierResponse := SupplierResponse{
// 			SupplierID:      supplier.SupplierID,
// 			PayablesAccount: *supplier.PayablesAccount,
// 			Profile: BioData{
// 				UID:        supplier.UserProfile.UID,
// 				Name:       supplier.UserProfile.Name,
// 				Gender:     supplier.UserProfile.Gender,
// 				Msisdns:    supplier.UserProfile.Msisdns,
// 				Emails:     supplier.UserProfile.Emails,
// 				PushTokens: supplier.UserProfile.PushTokens,
// 				Bio:        supplier.UserProfile.Bio,
// 			},
// 		}

// 		base.WriteJSONResponse(w, supplierResponse, http.StatusOK)
// 	}
// }

// // AddSupplierKyc persists a supplier KYC information to firestore
// func (s Service) AddSupplierKyc(
// 	ctx context.Context,
// 	input SupplierKYCInput) (*SupplierKYC, error) {
// 	s.checkPreconditions()

// 	profile, err := s.UserProfile(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch user profile: %v", err)
// 	}
// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
// 		ctx, profile.UID, s.GetSupplierCollectionName(), "userprofile.uid")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in out records")
// 	}
// 	supplier := &Supplier{}
// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	supplier.SupplierKYC.AccountType = input.AccountType
// 	supplier.SupplierKYC.IdentificationDocType = input.IdentificationDocType
// 	supplier.SupplierKYC.IdentificationDocNumber = input.IdentificationDocNumber
// 	supplier.SupplierKYC.IdentificationDocPhotoBase64 = input.IdentificationDocPhotoBase64
// 	supplier.SupplierKYC.IdentificationDocPhotoContentType = input.IdentificationDocPhotoContentType
// 	supplier.SupplierKYC.License = input.License
// 	supplier.SupplierKYC.Cadre = input.Cadre
// 	supplier.SupplierKYC.Profession = input.Profession
// 	supplier.SupplierKYC.KraPin = input.KraPin
// 	supplier.SupplierKYC.KraPINDocPhoto = input.KraPINDocPhoto
// 	supplier.SupplierKYC.BusinessNumber = input.BusinessNumber
// 	supplier.SupplierKYC.BusinessNumberDocPhotoBase64 = input.BusinessNumberDocPhotoBase64
// 	supplier.SupplierKYC.BusinessNumberDocPhotoContentType = input.BusinessNumberDocPhotoContentType

// 	err = base.UpdateRecordOnFirestore(
// 		s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
// 	}
// 	supplierKYC := supplier.SupplierKYC
// 	return &supplierKYC, nil
// }

// // SetUpSupplier performs initial account set up during onboarding
// func (s Service) SetUpSupplier(ctx context.Context, accountType AccountType) (*Supplier, error) {
// 	s.checkPreconditions()

// 	validAccountType := accountType.IsValid()
// 	if !validAccountType {
// 		return nil, fmt.Errorf("%v is not an allowed AccountType choice", accountType.String())
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
// 		ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %w", err)
// 	}
// 	supplier := &Supplier{}

// 	if dsnap == nil {
// 		return nil, fmt.Errorf("cannot find supplier record")
// 	}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier: %v", err)
// 	}

// 	profile, err := s.ParseUserProfileFromContextOrUID(ctx, &uid)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read user profile: %w", err)
// 	}

// 	supplier.UserProfile = profile
// 	supplier.AccountType = accountType
// 	supplier.UnderOrganization = false
// 	supplier.IsOrganizationVerified = false
// 	supplier.HasBranches = false

// 	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 	}

// 	go func() {
// 		op := func() error {
// 			return s.PublishKYCNudge(uid, &supplier.PartnerType, &supplier.AccountType)
// 		}

// 		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
// 			logrus.Error(err)
// 		}
// 	}()

// 	return supplier, nil
// }

// // SuspendSupplier flips the active boolean on the erp partner from true to false
// // consequently logically deleting the account
// func (s Service) SuspendSupplier(ctx context.Context, uid string) (bool, error) {
// 	s.checkPreconditions()

// 	err := s.DeleteUser(ctx, uid)
// 	if err != nil {
// 		return false, fmt.Errorf("error deleting user: %v", err)
// 	}

// 	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
// 	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, err
// 	}
// 	if len(docs) == 0 {
// 		return false, nil
// 	}

// 	dsnap := docs[0]
// 	supplier := &Supplier{}
// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to read supplier: %w", err)
// 	}

// 	payload := map[string]interface{}{
// 		"active": false,
// 	}

// 	content, marshalErr := json.Marshal(payload)
// 	if marshalErr != nil {
// 		return false, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
// 	}

// 	supplierPath := fmt.Sprintf("%s%s", customerAPIPath, supplier.SupplierID)
// 	if err := base.ReadRequestToTarget(s.erpClient, "PATCH", supplierPath, "", content, &supplier); err != nil {
// 		return false, fmt.Errorf("unable to make request to the ERP: %v", err)
// 	}

// 	if err = base.UpdateRecordOnFirestore(
// 		s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier,
// 	); err != nil {
// 		return false, fmt.Errorf("unable to update supplier: %v", err)
// 	}
// 	return true, nil
// }

// // EDIUserLogin used to login a user to EDI (Portal Authserver) and return their
// // EDI (Portal Authserver) profile
// func EDIUserLogin(username, password string) (*base.EDIUserProfile, error) {

// 	if username == "" || password == "" {
// 		return nil, fmt.Errorf("invalid credentials, expected a username AND password")
// 	}

// 	ediClient, err := base.LoginClient(username, password)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
// 	}

// 	userProfile, err := base.FetchUserProfile(ediClient)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
// 	}

// 	return userProfile, nil

// }

// // CoreEDIUserLogin used to login a user to EDI (Core Authserver) and return their EDI
// // EDI (Core Authserver) profile
// func CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error) {

// 	if username == "" || password == "" {
// 		return nil, fmt.Errorf("invalid credentials, expected a username AND password")
// 	}

// 	ediClient, err := LoginClient(username, password)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
// 	}

// 	userProfile, err := base.FetchUserProfile(ediClient)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
// 	}

// 	return userProfile, nil

// }

// // SupplierEDILogin it used to instantiate as call when setting up a supplier's account's who
// // has an affliation to a provider with the slade ecosystem. The logic is as follows;
// // 1 . login to the relevant edi to assert the user has an account
// // 2 . fetch the branches of the provider given the slade code which we have
// // 3 . update the user's supplier record
// // 4. return the list of branches to the frontend so that a default location can be set
// func (s Service) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*BranchConnection, error) {
// 	s.checkPreconditions()
// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %w", err)
// 	}

// 	supplier := &Supplier{}

// 	if dsnap == nil {
// 		return nil, fmt.Errorf("cannot find supplier record")
// 	}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier: %v", err)
// 	}

// 	profiledsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// return the newly created user profile
// 	profile := &base.UserProfile{}
// 	err = profiledsnap.DataTo(profile)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read user profile: %w", err)
// 	}

// 	supplier.UserProfile = profile
// 	supplier.AccountType = AccountTypeIndividual
// 	supplier.UnderOrganization = true

// 	ediUserProfile, err := func(sladeCode string) (*base.EDIUserProfile, error) {
// 		var ediUserProfile *base.EDIUserProfile
// 		var err error

// 		switch sladeCode {
// 		case savannahSladeCode:
// 			// login to core
// 			//TODO(calvine) add login to core
// 			ediUserProfile, err = CoreEDIUserLogin(username, password)
// 			if err != nil {
// 				supplier.IsOrganizationVerified = false
// 				return nil, fmt.Errorf("cannot get edi user profile: %w", err)
// 			}

// 			if ediUserProfile == nil {
// 				return nil, fmt.Errorf("edi user profile not found")
// 			}

// 		default:
// 			//Login to portal edi
// 			ediUserProfile, err = EDIUserLogin(username, password)
// 			if err != nil {
// 				supplier.IsOrganizationVerified = false
// 				return nil, fmt.Errorf("cannot get edi user profile: %w", err)
// 			}

// 			if ediUserProfile == nil {
// 				return nil, fmt.Errorf("edi user profile not found")
// 			}

// 		}
// 		return ediUserProfile, nil
// 	}(sladeCode)

// 	if err != nil {
// 		return nil, err
// 	}

// 	pageInfo := &base.PageInfo{
// 		HasNextPage:     false,
// 		HasPreviousPage: false,
// 		StartCursor:     nil,
// 		EndCursor:       nil,
// 	}

// 	// The slade code comes in the form 'PRO-1234' or 'BRA-PRO-1234-1'
// 	// or a single code '1234'
// 	// we split it to get the interger part of the slade code.
// 	var orgSladeCode string
// 	if strings.HasPrefix(sladeCode, "BRA-") {
// 		orgSladeCode = strings.Split(sladeCode, "-")[2]
// 	} else if strings.HasPrefix(sladeCode, "PRO-") {
// 		orgSladeCode = strings.Split(sladeCode, "-")[1]
// 	} else {
// 		orgSladeCode = sladeCode
// 	}

// 	if orgSladeCode == savannahSladeCode {
// 		profile.Permissions = base.DefaultAdminPermissions

// 		supplier.EDIUserProfile = ediUserProfile
// 		supplier.IsOrganizationVerified = true
// 		supplier.SladeCode = sladeCode
// 		supplier.Active = true
// 		supplier.KYCSubmitted = true
// 		supplier.PartnerSetupComplete = true

// 		if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 			return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 		}
// 		err = base.UpdateRecordOnFirestore(
// 			s.firestoreClient, s.GetUserProfileCollectionName(), profiledsnap.Ref.ID, profile,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to update user profile: %v", err)
// 		}

// 		return &BranchConnection{PageInfo: pageInfo}, nil
// 	}

// 	// verify slade code.
// 	if ediUserProfile.BusinessPartner != orgSladeCode {
// 		supplier.IsOrganizationVerified = false
// 		return nil, fmt.Errorf("invalid slade code for selected provider: %v, got: %v", sladeCode, ediUserProfile.BusinessPartner)
// 	}

// 	supplier.EDIUserProfile = ediUserProfile
// 	supplier.IsOrganizationVerified = true
// 	supplier.SladeCode = sladeCode

// 	filter := []*BusinessPartnerFilterInput{
// 		{
// 			SladeCode: &sladeCode,
// 		},
// 	}

// 	partner, err := s.FindProvider(ctx, nil, filter, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
// 	}

// 	var businessPartner BusinessPartner

// 	if len(partner.Edges) != 1 {
// 		return nil, fmt.Errorf("expected one business partner, found: %v", len(partner.Edges))
// 	}

// 	businessPartner = *partner.Edges[0].Node
// 	var brFilter []*BranchFilterInput

// 	go func() {
// 		op := func() error {
// 			return s.PublishKYCNudge(uid, &supplier.PartnerType, &supplier.AccountType)
// 		}

// 		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
// 			logrus.Error(err)
// 		}
// 	}()

// 	if businessPartner.Parent != nil {
// 		supplier.HasBranches = true
// 		supplier.ParentOrganizationID = *businessPartner.Parent
// 		filter := &BranchFilterInput{
// 			ParentOrganizationID: businessPartner.Parent,
// 		}

// 		brFilter = append(brFilter, filter)
// 		if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 			return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 		}

// 		return s.FindBranch(ctx, nil, brFilter, nil)
// 	}
// 	loc := Location{
// 		ID:   businessPartner.ID,
// 		Name: businessPartner.Name,
// 	}
// 	supplier.Location = &loc

// 	if err := s.SaveSupplierToFireStore(*supplier); err != nil {
// 		return nil, fmt.Errorf("unable to add supplier to firestore: %v", err)
// 	}

// 	return &BranchConnection{PageInfo: pageInfo}, nil
// }

// // SupplierSetDefaultLocation updates the default location ot the supplier by the given location id
// func (s Service) SupplierSetDefaultLocation(ctx context.Context, locationID string) (bool, error) {
// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to get the logged in user: %w", err)
// 	}

// 	// fetch the supplier records
// 	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
// 	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, fmt.Errorf("unable to fetch supplier record: %w", err)
// 	}
// 	if len(docs) == 0 {
// 		return false, fmt.Errorf("unable to find supplier record: %w", err)
// 	}

// 	dsnap := docs[0]
// 	sup := &Supplier{}
// 	err = dsnap.DataTo(sup)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to read supplier: %w", err)
// 	}

// 	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
// 	filter := []*BranchFilterInput{
// 		{
// 			SladeCode:            &sup.SladeCode,
// 			ParentOrganizationID: &sup.ParentOrganizationID,
// 		},
// 	}

// 	brs, err := s.FindBranch(ctx, nil, filter, nil)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to fetch organization branches location: %v", err)
// 	}

// 	branch := func(brs *BranchConnection, location string) *BranchEdge {
// 		for _, b := range brs.Edges {
// 			if b.Node.ID == location {
// 				return b
// 			}
// 		}
// 		return nil
// 	}(brs, locationID)

// 	if branch != nil {
// 		loc := Location{
// 			ID:              branch.Node.ID,
// 			Name:            branch.Node.Name,
// 			BranchSladeCode: &branch.Node.BranchSladeCode,
// 		}
// 		sup.Location = &loc

// 		// update the supplier record with new location
// 		if err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, sup); err != nil {
// 			return false, fmt.Errorf("unable to update supplier default location: %v", err)
// 		}
// 	}

// 	return false, fmt.Errorf("unable to get location of id %v : %v", locationID, err)
// }

// // FetchSupplierAllowedLocations retrieves all the locations that the user in context can work on.
// func (s *Service) FetchSupplierAllowedLocations(ctx context.Context) (*BranchConnection, error) {

// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
// 	}

// 	// fetch the supplier records
// 	collection := s.firestoreClient.Collection(s.GetSupplierCollectionName())
// 	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch supplier record: %w", err)
// 	}
// 	if len(docs) == 0 {
// 		return nil, fmt.Errorf("unable to find supplier record: %w", err)
// 	}

// 	dsnap := docs[0]
// 	sup := &Supplier{}
// 	err = dsnap.DataTo(sup)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier record: %w", err)
// 	}

// 	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
// 	filter := []*BranchFilterInput{
// 		{
// 			SladeCode:            &sup.SladeCode,
// 			ParentOrganizationID: &sup.ParentOrganizationID,
// 		},
// 	}

// 	brs, err := s.FindBranch(ctx, nil, filter, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
// 	}

// 	return brs, nil
// }

// // PublishKYCNudge pushes a kyc nudge to the user feed
// func (s *Service) PublishKYCNudge(uid string, partner *PartnerType, account *AccountType) error {

// 	s.checkPreconditions()

// 	if partner == nil || !partner.IsValid() {
// 		return fmt.Errorf("expected `partner` to be defined and to be valid")
// 	}

// 	if *partner == PartnerTypeConsumer {
// 		return fmt.Errorf("invalid `partner`. cannot use CONSUMER in this context")
// 	}

// 	if !account.IsValid() {
// 		return fmt.Errorf("provided `account` is not valid")
// 	}

// 	payload := base.Nudge{
// 		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
// 		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
// 		Visibility:     "SHOW",
// 		Status:         "PENDING",
// 		Expiry:         time.Now().Add(time.Hour * futureHours),
// 		Title:          fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
// 		Text:           "Fill in your Be.Well business KYC in order to start transacting",
// 		Links: []base.Link{
// 			{
// 				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
// 				URL:         base.LogoURL,
// 				LinkType:    base.LinkTypePngImage,
// 				Title:       "KYC",
// 				Description: fmt.Sprintf("KYC for %v", partner.String()),
// 				Thumbnail:   base.LogoURL,
// 			},
// 		},
// 		Actions: []base.Action{
// 			{
// 				ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
// 				SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
// 				Name:           strings.ToUpper(fmt.Sprintf("COMPLETE_%v_%v_KYC", account.String(), partner.String())),
// 				ActionType:     base.ActionTypePrimary,
// 				Handling:       base.HandlingFullPage,
// 				AllowAnonymous: false,
// 				Icon: base.Link{
// 					ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
// 					URL:         base.LogoURL,
// 					LinkType:    base.LinkTypePngImage,
// 					Title:       fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
// 					Description: "Fill in your Be.Well business KYC in order to start transacting",
// 					Thumbnail:   base.LogoURL,
// 				},
// 			},
// 		},
// 		Users:                []string{uid},
// 		Groups:               []string{uid},
// 		NotificationChannels: []base.Channel{base.ChannelEmail, base.ChannelFcm},
// 	}

// 	resp, err := s.engagement.MakeRequest("POST", fmt.Sprintf(publishNudge, uid), payload)
// 	if err != nil {
// 		return fmt.Errorf("unable to publish kyc nudge : %v", err)
// 	}

// 	//TODO(dexter) to be removed. Just here for debug
// 	res, _ := httputil.DumpResponse(resp, true)
// 	log.Println(string(res))

// 	if resp.StatusCode != http.StatusOK {
// 		// stage the nudge
// 		stage := func(pl base.Nudge) error {
// 			k, err := json.Marshal(payload)
// 			if err != nil {
// 				return fmt.Errorf("cannot marshal payload to json")
// 			}

// 			var kMap map[string]interface{}
// 			err = json.Unmarshal(k, &kMap)
// 			if err != nil {
// 				return fmt.Errorf("cannot unmarshal payload from json")
// 			}

// 			if err := s.SaveProfileNudge(kMap); err != nil {
// 				logrus.Errorf("failed to stage nudge : %v", err)
// 			}
// 			return nil

// 		}(payload)

// 		if err := stage; err != nil {
// 			logrus.Errorf("failed to stage nudge : %v", err)
// 		}
// 		return fmt.Errorf("unable to publish kyc nudge. unexpected status code  %v", resp.StatusCode)
// 	}

// 	return nil
// }

// // PublishKYCFeedItem notifies admin users of a KYC approval request
// func (s Service) PublishKYCFeedItem(ctx context.Context, uids ...string) error {

// 	s.checkPreconditions()

// 	for _, uid := range uids {
// 		payload := base.Item{
// 			ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
// 			SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
// 			Expiry:         time.Now().Add(time.Hour * futureHours),
// 			Persistent:     true,
// 			Status:         base.StatusPending,
// 			Visibility:     base.VisibilityShow,
// 			Author:         "Be.Well Team",
// 			Label:          "KYC",
// 			Tagline:        "Process incoming KYC",
// 			Text:           "Review KYC for the partner and either approve or reject",
// 			TextType:       base.TextTypeMarkdown,
// 			Icon: base.Link{
// 				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
// 				URL:         base.LogoURL,
// 				LinkType:    base.LinkTypePngImage,
// 				Title:       "KYC Review",
// 				Description: "Review KYC for the partner and either approve or reject",
// 				Thumbnail:   base.LogoURL,
// 			},
// 			Timestamp: time.Now(),
// 			Actions: []base.Action{
// 				{
// 					ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
// 					SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
// 					Name:           "Review KYC details",
// 					Icon: base.Link{
// 						ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
// 						URL:         base.LogoURL,
// 						LinkType:    base.LinkTypePngImage,
// 						Title:       "Review KYC details",
// 						Description: "Review and approve or reject KYC details for the supplier",
// 						Thumbnail:   base.LogoURL,
// 					},
// 					ActionType:     base.ActionTypePrimary,
// 					Handling:       base.HandlingFullPage,
// 					AllowAnonymous: false,
// 				},
// 			},
// 			Links: []base.Link{
// 				{
// 					ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
// 					URL:         base.LogoURL,
// 					LinkType:    base.LinkTypePngImage,
// 					Title:       "KYC process request",
// 					Description: "Process KYC request",
// 					Thumbnail:   base.LogoURL,
// 				},
// 			},

// 			Summary: "Process incoming KYC",
// 			Users:   uids,
// 			NotificationChannels: []base.Channel{
// 				base.ChannelFcm,
// 				base.ChannelEmail,
// 				base.ChannelSms,
// 			},
// 		}

// 		resp, err := s.engagement.MakeRequest("POST", fmt.Sprintf(publishItem, uid), payload)
// 		if err != nil {
// 			return fmt.Errorf("unable to publish kyc admin notification feed item : %v", err)
// 		}

// 		//TODO(dexter) to be removed. Just here for debug
// 		res, _ := httputil.DumpResponse(resp, true)
// 		log.Println(string(res))

// 		if resp.StatusCode != http.StatusOK {
// 			return fmt.Errorf("unable to publish kyc admin notification feed item. unexpected status code  %v", resp.StatusCode)
// 		}
// 	}

// 	return nil
// }

// // StageKYCProcessingRequest saves kyc processing requests
// func (s *Service) StageKYCProcessingRequest(sup *Supplier) error {
// 	r := KYCRequest{
// 		ID:                  uuid.New().String(),
// 		ReqPartnerType:      sup.PartnerType,
// 		ReqOrganizationType: OrganizationType(sup.AccountType),
// 		ReqRaw:              sup.SupplierKYC,
// 		Proceseed:           false,
// 		SupplierRecord:      sup,
// 		Status:              KYCProcessStatusPending,
// 	}

// 	_, err := base.SaveDataToFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), r)
// 	if err != nil {
// 		return fmt.Errorf("unable to save kyc processing request: %w", err)
// 	}
// 	return nil
// }

// // AddIndividualRiderKyc adds KYC for an individual rider
// func (s *Service) AddIndividualRiderKyc(ctx context.Context, input IndividualRider) (*IndividualRider, error) {

// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := IndividualRider{
// 		IdentificationDoc: Identification{
// 			IdentificationDocType:           input.IdentificationDoc.IdentificationDocType,
// 			IdentificationDocNumber:         input.IdentificationDoc.IdentificationDocNumber,
// 			IdentificationDocNumberUploadID: input.IdentificationDoc.IdentificationDocNumberUploadID,
// 		},
// 		KRAPIN:                         input.KRAPIN,
// 		KRAPINUploadID:                 input.KRAPINUploadID,
// 		DrivingLicenseID:               input.DrivingLicenseID,
// 		DrivingLicenseUploadID:         input.DrivingLicenseUploadID,
// 		CertificateGoodConductUploadID: input.CertificateGoodConductUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationRiderKyc adds KYC for an organization rider
// func (s *Service) AddOrganizationRiderKyc(ctx context.Context, input OrganizationRider) (*OrganizationRider, error) {

// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationRider{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,

// 		KRAPIN:                      input.KRAPIN,
// 		KRAPINUploadID:              input.KRAPINUploadID,
// 		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddIndividualPractitionerKyc adds KYC for an individual pratitioner
// func (s *Service) AddIndividualPractitionerKyc(ctx context.Context, input IndividualPractitioner) (*IndividualPractitioner, error) {
// 	s.checkPreconditions()

// 	for _, p := range input.PracticeServices {
// 		if !p.IsValid() {
// 			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
// 		}
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := IndividualPractitioner{

// 		IdentificationDoc: func(p Identification) Identification {
// 			return Identification(p)
// 		}(input.IdentificationDoc),

// 		KRAPIN:                      input.KRAPIN,
// 		KRAPINUploadID:              input.KRAPINUploadID,
// 		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
// 		RegistrationNumber:          input.RegistrationNumber,
// 		PracticeLicenseID:           input.PracticeLicenseID,
// 		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
// 		PracticeServices:            input.PracticeServices,
// 		Cadre:                       input.Cadre,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationPractitionerKyc adds KYC for an organization pratitioner
// func (s *Service) AddOrganizationPractitionerKyc(ctx context.Context, input OrganizationPractitioner) (*OrganizationPractitioner, error) {

// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	for _, p := range input.PracticeServices {
// 		if !p.IsValid() {
// 			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
// 		}
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationPractitioner{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		KRAPIN:                             input.KRAPIN,
// 		KRAPINUploadID:                     input.KRAPINUploadID,
// 		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
// 		RegistrationNumber:                 input.RegistrationNumber,
// 		PracticeLicenseID:                  input.PracticeLicenseID,
// 		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
// 		PracticeServices:                   input.PracticeServices,
// 		Cadre:                              input.Cadre,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationProviderKyc adds KYC for an organization provider
// func (s *Service) AddOrganizationProviderKyc(ctx context.Context, input OrganizationProvider) (*OrganizationProvider, error) {

// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	for _, p := range input.PracticeServices {
// 		if !p.IsValid() {
// 			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
// 		}
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationProvider{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		KRAPIN:                             input.KRAPIN,
// 		KRAPINUploadID:                     input.KRAPINUploadID,
// 		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
// 		RegistrationNumber:                 input.RegistrationNumber,
// 		PracticeLicenseID:                  input.PracticeLicenseID,
// 		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
// 		PracticeServices:                   input.PracticeServices,
// 		Cadre:                              input.Cadre,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddIndividualPharmaceuticalKyc adds KYC for an individual Pharmaceutical kyc
// func (s *Service) AddIndividualPharmaceuticalKyc(ctx context.Context, input IndividualPharmaceutical) (*IndividualPharmaceutical, error) {

// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := IndividualPharmaceutical{
// 		IdentificationDoc: func(p Identification) Identification {
// 			return Identification(p)
// 		}(input.IdentificationDoc),
// 		KRAPIN:                      input.KRAPIN,
// 		KRAPINUploadID:              input.KRAPINUploadID,
// 		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
// 		RegistrationNumber:          input.RegistrationNumber,
// 		PracticeLicenseID:           input.PracticeLicenseID,
// 		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}
// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationPharmaceuticalKyc adds KYC for a pharmacy organization
// func (s *Service) AddOrganizationPharmaceuticalKyc(ctx context.Context, input OrganizationPharmaceutical) (*OrganizationPharmaceutical, error) {
// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationPharmaceutical{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		KRAPIN:                             input.KRAPIN,
// 		KRAPINUploadID:                     input.KRAPINUploadID,
// 		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,
// 		RegistrationNumber:      input.RegistrationNumber,
// 		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}

// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddIndividualCoachKyc adds KYC for an individual coach
// func (s *Service) AddIndividualCoachKyc(ctx context.Context, input IndividualCoach) (*IndividualCoach, error) {
// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := IndividualCoach{
// 		IdentificationDoc: func(p Identification) Identification {
// 			return Identification(p)
// 		}(input.IdentificationDoc),
// 		KRAPIN:                      input.KRAPIN,
// 		KRAPINUploadID:              input.KRAPINUploadID,
// 		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
// 		PracticeLicenseID:           input.PracticeLicenseID,
// 		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}

// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationCoachKyc adds KYC for an organization coach
// func (s *Service) AddOrganizationCoachKyc(ctx context.Context, input OrganizationCoach) (*OrganizationCoach, error) {
// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationCoach{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		KRAPIN:                             input.KRAPIN,
// 		KRAPINUploadID:                     input.KRAPINUploadID,
// 		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,
// 		RegistrationNumber:      input.RegistrationNumber,
// 		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}

// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddIndividualNutritionKyc adds KYC for an individual nutritionist
// func (s *Service) AddIndividualNutritionKyc(ctx context.Context, input IndividualNutrition) (*IndividualNutrition, error) {
// 	s.checkPreconditions()

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := IndividualNutrition{
// 		IdentificationDoc: func(p Identification) Identification {
// 			return Identification(p)
// 		}(input.IdentificationDoc),
// 		KRAPIN:                      input.KRAPIN,
// 		KRAPINUploadID:              input.KRAPINUploadID,
// 		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
// 		PracticeLicenseID:           input.PracticeLicenseID,
// 		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}

// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // AddOrganizationNutritionKyc adds kyc for a nutritionist organisation
// func (s *Service) AddOrganizationNutritionKyc(ctx context.Context, input OrganizationNutrition) (*OrganizationNutrition, error) {
// 	s.checkPreconditions()

// 	if !input.OrganizationTypeName.IsValid() {
// 		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
// 	}

// 	uid, err := base.GetLoggedInUserUID(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
// 	}

// 	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, uid, s.GetSupplierCollectionName(), "userprofile.verifiedIdentifiers")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to retrieve supplier from collections: %v", err)
// 	}
// 	if dsnap == nil {
// 		return nil, fmt.Errorf("the supplier does not exist in our records")
// 	}

// 	supplier := &Supplier{}

// 	err = dsnap.DataTo(supplier)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read supplier data: %v", err)
// 	}

// 	kyc := OrganizationNutrition{
// 		OrganizationTypeName:               input.OrganizationTypeName,
// 		KRAPIN:                             input.KRAPIN,
// 		KRAPINUploadID:                     input.KRAPINUploadID,
// 		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
// 		CertificateOfIncorporation:         input.CertificateOfIncorporation,
// 		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
// 		DirectorIdentifications: func(p []Identification) []Identification {
// 			pl := []Identification{}
// 			for _, i := range p {
// 				pl = append(pl, Identification(i))
// 			}
// 			return pl
// 		}(input.DirectorIdentifications),
// 		OrganizationCertificate: input.OrganizationCertificate,
// 		RegistrationNumber:      input.RegistrationNumber,
// 		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
// 	}

// 	if len(input.SupportingDocumentsUploadID) != 0 {
// 		ids := []string{}
// 		ids = append(ids, input.SupportingDocumentsUploadID...)

// 		kyc.SupportingDocumentsUploadID = ids
// 	}

// 	k, err := json.Marshal(kyc)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot marshal kyc to json")
// 	}

// 	err = s.SaveKYCResponse(ctx, k, supplier, dsnap)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot save KYC request: %v", err)
// 	}

// 	return &kyc, nil
// }

// // SaveKYCResponse updates the record of a supplier with the provided KYC, stages the request for
// // approval by Savannah admins and sends a notification of the request to admins
// func (s Service) SaveKYCResponse(ctx context.Context, kycJSON []byte, supplier *Supplier, dsnap *firestore.DocumentSnapshot) error {
// 	var kycAsMap map[string]interface{}

// 	err := json.Unmarshal(kycJSON, &kycAsMap)
// 	if err != nil {
// 		return fmt.Errorf("cannot unmarshal kyc from json")
// 	}

// 	supplier.SupplierKYC = kycAsMap
// 	supplier.KYCSubmitted = true

// 	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetSupplierCollectionName(), dsnap.Ref.ID, supplier)
// 	if err != nil {
// 		return fmt.Errorf("unable to update supplier with supplier KYC info: %v", err)
// 	}

// 	if err := s.StageKYCProcessingRequest(supplier); err != nil {
// 		logrus.Errorf("unable to stage kyc processing request: %v", err)
// 	}

// 	go func() {
// 		op := func() error {
// 			a, err := s.FetchAdminUsers(ctx)
// 			if err != nil {
// 				return err
// 			}
// 			var uids []string
// 			for _, u := range a {
// 				uids = append(uids, u.ID)
// 			}

// 			return s.PublishKYCFeedItem(ctx, uids...)
// 		}

// 		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
// 			logrus.Error(err)
// 		}
// 	}()

// 	return nil
// }

// // FetchKYCProcessingRequests fetches a list of all unprocessed kyc approval requests
// func (s *Service) FetchKYCProcessingRequests(ctx context.Context) ([]*KYCRequest, error) {
// 	collection := s.firestoreClient.Collection(s.GetKCYProcessCollectionName())
// 	query := collection.Where("proceseed", "==", false)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to fetch kyc request documents: %v", err)
// 	}

// 	res := []*KYCRequest{}

// 	for _, doc := range docs {
// 		req := &KYCRequest{}
// 		err = doc.DataTo(req)
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to read supplier: %w", err)
// 		}
// 		res = append(res, req)
// 	}
// 	return res, nil
// }

// // ProcessKYCRequest transitions a kyc request to a given state
// func (s *Service) ProcessKYCRequest(ctx context.Context, id string, status KYCProcessStatus, rejectionReason *string) (bool, error) {
// 	collection := s.firestoreClient.Collection(s.GetKCYProcessCollectionName())
// 	query := collection.Where("id", "==", id)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, fmt.Errorf("unable to fetch kyc request documents: %v", err)
// 	}

// 	doc := docs[0]
// 	req := &KYCRequest{}
// 	err = doc.DataTo(req)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to read supplier: %w", err)
// 	}

// 	req.Status = status
// 	req.Proceseed = true
// 	req.RejectionReason = rejectionReason

// 	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), doc.Ref.ID, req)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to update KYC request record: %v", err)
// 	}

// 	var email string
// 	var message string

// 	switch status {
// 	case KYCProcessStatusApproved:
// 		// create supplier erp account
// 		if _, err := s.AddSupplier(ctx, *req.SupplierRecord.UserProfile.Name, req.ReqPartnerType); err != nil {
// 			return false, fmt.Errorf("unable to create erp supplier account: %v", err)
// 		}

// 		email = generateProcessKYCApprovalEmailTemplate()
// 		message = "Your KYC details have been reviewed and approved. We look forward to working with you."

// 	case KYCProcessStatusRejected:
// 		email = generateProcessKYCRejectionEmailTemplate()
// 		message = "Your KYC details have been reviewed and not verified. Incase of any queries, please contact us via +254 790 360 360"

// 	}

// 	for _, supplierEmail := range req.SupplierRecord.UserProfile.Emails {
// 		err = s.SendKYCEmail(ctx, email, supplierEmail)
// 		if err != nil {
// 			return false, fmt.Errorf("unable to send KYC processing email: %w", err)
// 		}
// 	}

// 	smsISC := base.SmsISC{
// 		Isc:      s.sms,
// 		EndPoint: sendSMS,
// 	}

// 	twilioISC := base.SmsISC{
// 		Isc:      s.twilio,
// 		EndPoint: sendTwilioSMS,
// 	}

// 	err = base.SendSMS(req.SupplierRecord.UserProfile.Msisdns, message, smsISC, twilioISC)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to send KYC processing message: %w", err)
// 	}

// 	return true, nil

// }

// // SendKYCEmail will send a KYC processing request email to the supplier
// func (s Service) SendKYCEmail(ctx context.Context, text, emailaddress string) error {
// 	if !govalidator.IsEmail(emailaddress) {
// 		return nil
// 	}

// 	body := map[string]interface{}{
// 		"to":      []string{emailaddress},
// 		"text":    text,
// 		"subject": emailSignupSubject,
// 	}

// 	resp, err := s.Mailgun.MakeRequest(http.MethodPost, SendEmail, body)
// 	if err != nil {
// 		return fmt.Errorf("unable to send KYC email: %w", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("unable to send KYC email : %w, with status code %v", err, resp.StatusCode)
// 	}

// 	return nil
// }

// // DeleteUser deletes a user records given their uid
// func (s Service) DeleteUser(ctx context.Context, uid string) error {
// 	s.checkPreconditions()

// 	profile, err := s.GetProfile(ctx, uid)
// 	if err != nil {
// 		return fmt.Errorf("unable to get user profile: %v", err)
// 	}

// 	for _, profileUID := range profile.VerifiedIdentifiers {
// 		params := (&auth.UserToUpdate{}).
// 			Disabled(true)
// 		_, err := s.firebaseAuth.UpdateUser(ctx, profileUID, params)
// 		if err != nil {
// 			return fmt.Errorf("error updating user: %v", err)
// 		}
// 	}

// 	profile.Active = false
// 	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
// 	if err != nil {
// 		return fmt.Errorf("unable to retrieve user profile doc snapshot: %v", err)
// 	}

// 	err = base.UpdateRecordOnFirestore(
// 		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, profile,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("unable to update user profile: %v", err)
// 	}

// 	return nil
// }

// // SaveProfileNudge stages nudges published from this service. These nudges will be
// // referenced later to support some specialized use-case. A nudge will be uniquely
// // identified by its id and sequenceNumber
// func (s Service) SaveProfileNudge(nudge map[string]interface{}) error {
// 	ctx := context.Background()
// 	_, _, err := s.firestoreClient.Collection(s.GetProfileNudgesCollectionName()).Add(ctx, nudge)
// 	return err
// }
