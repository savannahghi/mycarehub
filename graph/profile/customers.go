package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

const (
	customerAPIPath        = "/api/business_partners/customers/"
	active                 = true
	country                = "KEN" // Anticipate worldwide expansion
	isCustomer             = true
	customerType           = "PATIENT" // Further Discussions
	customerCollectionName = "customers"

	// Fetch the orgnisation's default currency from the env
	// Currency is used in the creation of a business partner in the ERP
	erpCurrencyEnvName = "ERP_DEFAULT_CURRENCY"
)

// SaveCustomerToFireStore persists customer data to firestore
func (s Service) SaveCustomerToFireStore(customer Customer) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetCustomerCollectionName()).Add(ctx, customer)
	return err
}

// GetCustomerCollectionName creates a suffixed customer collection name
func (s Service) GetCustomerCollectionName() string {
	suffixed := base.SuffixCollection(customerCollectionName)
	return suffixed
}

// AddCustomer creates a customer on the ERP when a user signs up in our Be.Well Consumer
func (s Service) AddCustomer(ctx context.Context) (*Customer, error) {
	s.checkPreconditions()

	profile, profileErr := s.UserProfile(ctx)
	if profileErr != nil {
		return nil, profileErr
	}

	customer, err := s.FindCustomer(ctx, profile.UID)
	if err != nil {
		return nil, fmt.Errorf("unable to get customer: %v", err)
	}

	if customer != nil {
		return customer, nil
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
		"is_customer":   isCustomer,
		"customer_type": customerType,
	}

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}
	newCustomer := Customer{
		UserProfile: *profile,
	}

	err = base.ReadRequestToTarget(s.client, "POST", customerAPIPath, "", content, &newCustomer)
	if err != nil {
		return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	err = s.SaveCustomerToFireStore(newCustomer)
	if err != nil {
		return nil, fmt.Errorf("unable to add customer to firestore: %v", err)
	}

	profile.HasCustomerAccount = true
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

	return &newCustomer, nil
}

// AddCustomerKYC persists information to know your customer
func (s Service) AddCustomerKYC(ctx context.Context, input CustomerKYCInput) (*CustomerKYC, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch user profile: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, profile.UID, s.GetCustomerCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve customer: %v", err)
	}
	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer data: %v", err)
	}

	customer.CustomerKYC.KRAPin = input.KRAPin
	customer.CustomerKYC.Occupation = input.Occupation
	customer.CustomerKYC.IDNumber = input.IDNumber
	customer.CustomerKYC.Address = input.Address
	customer.CustomerKYC.City = input.City

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetCustomerCollectionName(), dsnap.Ref.ID, customer,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update customer with customer KYC info: %v", err)
	}

	customerKYC := customer.CustomerKYC
	return &customerKYC, nil
}

// UpdateCustomer updates a customerKYC information in firestore
func (s Service) UpdateCustomer(ctx context.Context, input CustomerKYCInput) (*Customer, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch user profile: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, profile.UID, s.GetCustomerCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve customer: %v", err)
	}

	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer: %w", err)
	}

	if input.KRAPin != "" {
		customer.CustomerKYC.KRAPin = input.KRAPin
	}

	if input.Occupation != "" {
		customer.CustomerKYC.Occupation = input.Occupation
	}

	if input.IDNumber != "" {
		customer.CustomerKYC.IDNumber = input.IDNumber
	}

	if input.City != "" {
		customer.CustomerKYC.City = input.City
	}

	if input.Address != "" {
		customer.CustomerKYC.Address = input.Address
	}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetCustomerCollectionName(), dsnap.Ref.ID, customer,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update customer with customer KYC info: %v", err)
	}

	return customer, nil
}

// FindCustomer fetches a customer by their UID
func (s Service) FindCustomer(ctx context.Context, uid string) (*Customer, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetCustomerCollectionName(), "userprofile.uid")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %v", err)
	}

	if dsnap == nil {
		return nil, nil
	}

	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer: %v", err)
	}

	return customer, nil
}

// FindCustomerByUIDHandler is a used for inter service communication to return details about a customer
func FindCustomerByUIDHandler(ctx context.Context, service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerUID, err := ValidateUID(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		customer, err := service.FindCustomer(ctx, customerUID)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if customer == nil {
			base.WriteJSONResponse(w, StatusResponse{Status: "not found"}, http.StatusNotFound)
			return
		}

		customerResponse := CustomerResponse{
			CustomerID: customer.CustomerID,
			ReceivablesAccount: ReceivablesAccount{
				ID:          customer.ReceivablesAccount.ID,
				Name:        customer.ReceivablesAccount.Name,
				IsActive:    customer.ReceivablesAccount.IsActive,
				Number:      customer.ReceivablesAccount.Number,
				Tag:         customer.ReceivablesAccount.Tag,
				Description: customer.ReceivablesAccount.Description,
			},
		}

		base.WriteJSONResponse(w, customerResponse, http.StatusOK)
	}
}
