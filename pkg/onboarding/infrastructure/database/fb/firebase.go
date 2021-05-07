package fb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	userProfileCollectionName            = "user_profiles"
	supplierProfileCollectionName        = "supplier_profiles"
	customerProfileCollectionName        = "customer_profiles"
	pinsCollectionName                   = "pins"
	surveyCollectionName                 = "post_visit_survey"
	profileNudgesCollectionName          = "profile_nudges"
	kycProcessCollectionName             = "kyc_processing"
	experimentParticipantCollectionName  = "experiment_participants"
	nhifDetailsCollectionName            = "nhif_details"
	communicationsSettingsCollectionName = "communications_settings"

	firebaseExchangeRefreshTokenURL = "https://securetoken.googleapis.com/v1/token?key="
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	FirestoreClient FirestoreClientExtension
	FirebaseClient  FirebaseClientExtension
}

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(firestoreClient FirestoreClientExtension, firebaseClient FirebaseClientExtension) repository.OnboardingRepository {
	return &Repository{
		FirestoreClient: firestoreClient,
		FirebaseClient:  firebaseClient,
	}
}

// GetUserProfileCollectionName ...
func (fr Repository) GetUserProfileCollectionName() string {
	suffixed := base.SuffixCollection(userProfileCollectionName)
	return suffixed
}

// GetSupplierProfileCollectionName ...
func (fr Repository) GetSupplierProfileCollectionName() string {
	suffixed := base.SuffixCollection(supplierProfileCollectionName)
	return suffixed
}

// GetCustomerProfileCollectionName ...
func (fr Repository) GetCustomerProfileCollectionName() string {
	suffixed := base.SuffixCollection(customerProfileCollectionName)
	return suffixed
}

// GetSurveyCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetSurveyCollectionName() string {
	suffixed := base.SuffixCollection(surveyCollectionName)
	return suffixed
}

// GetPINsCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetPINsCollectionName() string {
	suffixed := base.SuffixCollection(pinsCollectionName)
	return suffixed
}

// GetProfileNudgesCollectionName return the storage location of profile nudges
func (fr Repository) GetProfileNudgesCollectionName() string {
	suffixed := base.SuffixCollection(profileNudgesCollectionName)
	return suffixed
}

// GetKCYProcessCollectionName fetches collection where kyc processing request will be saved
func (fr Repository) GetKCYProcessCollectionName() string {
	suffixed := base.SuffixCollection(kycProcessCollectionName)
	return suffixed
}

// GetExperimentParticipantCollectionName fetches the collection where experiment participant will be saved
func (fr *Repository) GetExperimentParticipantCollectionName() string {
	suffixed := base.SuffixCollection(experimentParticipantCollectionName)
	return suffixed
}

// GetNHIFDetailsCollectionName ...
func (fr Repository) GetNHIFDetailsCollectionName() string {
	suffixed := base.SuffixCollection(nhifDetailsCollectionName)
	return suffixed
}

// GetCommunicationsSettingsCollectionName ...
func (fr Repository) GetCommunicationsSettingsCollectionName() string {
	suffixed := base.SuffixCollection(communicationsSettingsCollectionName)
	return suffixed
}

// GetUserProfileByUID retrieves the user profile by UID
func (fr *Repository) GetUserProfileByUID(
	ctx context.Context,
	uid string,
	suspended bool,
) (*base.UserProfile, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "verifiedUIDS",
		Value:          uid,
		Operator:       "array-contains",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, exceptions.ProfileNotFoundError()
	}

	if len(docs) > 1 && base.IsDebug() {
		log.Printf("user with uids %s has > 1 profile (they have %d)",
			uid,
			len(docs),
		)
	}

	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		err = fmt.Errorf("unable to read user profile")
		return nil, exceptions.InternalServerError(err)
	}

	if !suspended {
		// never return a suspended user profile
		if userProfile.Suspended {
			return nil, exceptions.ProfileSuspendFoundError()
		}
	}

	return userProfile, nil
}

// GetUserProfileByID retrieves a user profile by ID
func (fr *Repository) GetUserProfileByID(
	ctx context.Context,
	id string,
	suspended bool,
) (*base.UserProfile, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.ProfileNotFoundError()
	}
	dsnap := docs[0]
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read user profile: %w", err))
	}

	if !suspended {
		// never return a suspended user profile
		if userProfile.Suspended {
			return nil, exceptions.ProfileSuspendFoundError()
		}
	}
	return userProfile, nil
}

func (fr *Repository) fetchUserRandomName(ctx context.Context) *string {
	n := utils.GetRandomName()
	if v, err := fr.CheckIfUsernameExists(ctx, *n); v && (err == nil) {
		return fr.fetchUserRandomName(ctx)
	}
	return n
}

// CreateUserProfile creates a user profile of using the provided phone number and uid
func (fr *Repository) CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {

	v, err := fr.CheckIfPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("failed to check the phone number: %v", err))
	}

	if v {
		// this phone is number is associated with another user profile, hence can not create an profile with the same phone number
		return nil, exceptions.CheckPhoneNumberExistError()
	}

	profileID := uuid.New().String()
	pr := &base.UserProfile{
		ID:           profileID,
		UserName:     fr.fetchUserRandomName(ctx),
		PrimaryPhone: &phoneNumber,
		VerifiedIdentifiers: []base.VerifiedIdentifier{{
			UID:           uid,
			LoginProvider: base.LoginProviderTypePhone,
			Timestamp:     time.Now().In(base.TimeLocation),
		}},
		VerifiedUIDS:  []string{uid},
		TermsAccepted: true,
		Suspended:     false,
	}

	command := &CreateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		Data:           pr,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to create new user profile: %w", err))
	}
	query := &GetSingleQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to retrieve newly created user profile: %w", err))
	}
	// return the newly created user profile
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read user profile: %w", err))
	}
	return userProfile, nil

}

// CreateEmptySupplierProfile creates an empty supplier profile
func (fr *Repository) CreateEmptySupplierProfile(ctx context.Context, profileID string) (*base.Supplier, error) {
	sup := &base.Supplier{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		Data:           sup,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to create new supplier empty profile: %w", err))
	}
	getSupplierquery := &GetSingleQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getSupplierquery)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to retrieve newly created supplier profile: %w", err))
	}
	// return the newly created supplier profile
	supplier := &base.Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read supplier profile: %w", err))
	}
	return supplier, nil

}

// CreateEmptyCustomerProfile creates an empty customer profile
func (fr *Repository) CreateEmptyCustomerProfile(ctx context.Context, profileID string) (*base.Customer, error) {
	cus := &base.Customer{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		Data:           cus,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to create new customer empty profile: %w", err))
	}
	getSupplierquery := &GetSingleQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getSupplierquery)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to retrieve newly created customer profile: %w", err))
	}
	// return the newly created customer profile
	customer := &base.Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read customer profile: %w", err))
	}
	return customer, nil
}

//GetUserProfileByPrimaryPhoneNumber fetches a user profile by primary phone number
func (fr *Repository) GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) == 0 {
		return nil, exceptions.ProfileNotFoundError()
	}
	dsnap := docs[0]
	profile := &base.UserProfile{}
	err = dsnap.DataTo(profile)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read user profile: %w", err))
	}

	if !suspended {
		// never return a suspended user profile
		if profile.Suspended {
			return nil, exceptions.ProfileSuspendFoundError()
		}
	}
	return profile, nil
}

// GetUserProfileByPhoneNumber fetches a user profile by phone number. This method traverses both PRIMARY PHONE numbers
// and SECONDARY PHONE numbers.
func (fr *Repository) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
	// check first primary phone numbers
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) == 1 {
		dsnap := docs[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, exceptions.InternalServerError(fmt.Errorf("unable to read customer profile: %w", err))
		}
		return pr, nil
	}

	// then check in secondary phone numbers
	query1 := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "secondaryPhoneNumbers",
		Value:          phoneNumber,
		Operator:       "array-contains",
	}
	docs1, err := fr.FirestoreClient.GetAll(ctx, query1)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs1) == 1 {
		dsnap := docs1[0]
		pr := &base.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, exceptions.InternalServerError(fmt.Errorf("unable to read customer profile: %w", err))
		}

		if !suspended {
			// never return a suspended user profile
			if pr.Suspended {
				return nil, exceptions.ProfileSuspendFoundError()
			}
		}

		return pr, nil
	}

	return nil, exceptions.ProfileNotFoundError()

}

// CheckIfPhoneNumberExists checks both PRIMARY PHONE NUMBERs and SECONDARY PHONE NUMBERs for the
// existence of the argument phoneNumber.
func (fr *Repository) CheckIfPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	// check first primary phone numbers
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}

	if len(docs) > 0 {
		return true, nil
	}

	// then check in secondary phone numbers
	query1 := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "secondaryPhoneNumbers",
		Value:          phoneNumber,
		Operator:       "array-contains",
	}
	docs1, err := fr.FirestoreClient.GetAll(ctx, query1)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	if len(docs1) > 0 {
		return true, nil
	}

	return false, nil
}

// CheckIfEmailExists checks in both PRIMARY EMAIL and SECONDARY EMAIL for the
// existence of the argument email
func (fr *Repository) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	// check first primary email
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryEmailAddress",
		Value:          email,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	if len(docs) == 1 {
		return true, nil
	}

	// then check in secondary email
	query1 := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "secondaryEmailAddresses",
		Value:          email,
		Operator:       "array-contains",
	}
	docs1, err := fr.FirestoreClient.GetAll(ctx, query1)
	if err != nil {
		return false, err
	}
	if len(docs1) == 1 {
		return true, nil
	}
	return false, nil
}

// CheckIfUsernameExists checks if the provided username exists. If true, it means its has already been associated with
// another user
func (fr *Repository) CheckIfUsernameExists(ctx context.Context, userName string) (bool, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "userName",
		Value:          userName,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	if len(docs) == 1 {
		return true, nil
	}

	return false, nil
}

// GetPINByProfileID gets a user's PIN by their profile ID
func (fr *Repository) GetPINByProfileID(ctx context.Context, profileID string) (*domain.PIN, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetPINsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	// this should never run. If it does, it means we are doing something wrong.
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 PINs with profile ID %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.PinNotFoundError(fmt.Errorf("failed to get a user pin"))
	}

	dsnap := docs[0]
	PIN := &domain.PIN{}
	err = dsnap.DataTo(PIN)
	if err != nil {
		return nil, err
	}

	return PIN, nil
}

// GenerateAuthCredentialsForAnonymousUser generates auth credentials for the anonymous user. This method is here since we don't
// want to delegate sign-in of anonymous users to the frontend. This is an effort the over reliance on firebase and lettin us
// handle all the heavy lifting
func (fr *Repository) GenerateAuthCredentialsForAnonymousUser(ctx context.Context) (*base.AuthCredentialResponse, error) {
	// todo(dexter) : move anonymousPhoneNumber to base. AnonymousPhoneNumber should NEVER NEVER have a user profile
	anonymousPhoneNumber := "+254700000000"

	u, err := fr.GetOrCreatePhoneNumberUser(ctx, anonymousPhoneNumber)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, u.UID)
	if err != nil {
		return nil, exceptions.CustomTokenError(err)
	}
	userTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, exceptions.AuthenticateTokenError(err)
	}

	return &base.AuthCredentialResponse{
		CustomToken:  &customToken,
		IDToken:      &userTokens.IDToken,
		ExpiresIn:    userTokens.ExpiresIn,
		RefreshToken: userTokens.RefreshToken,
		UID:          u.UID,
		IsAnonymous:  true,
		IsAdmin:      false,
	}, nil
}

// GenerateAuthCredentials gets a Firebase user by phone and creates their tokens
func (fr *Repository) GenerateAuthCredentials(
	ctx context.Context,
	phone string,
	profile *base.UserProfile,
) (*base.AuthCredentialResponse, error) {
	resp, err := fr.GetOrCreatePhoneNumberUser(ctx, phone)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, exceptions.UserNotFoundError(err)
		}
		return nil, exceptions.UserNotFoundError(err)
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, resp.UID)
	if err != nil {
		return nil, exceptions.CustomTokenError(err)
	}
	userTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, exceptions.AuthenticateTokenError(err)
	}

	if err := fr.UpdateVerifiedIdentifiers(ctx, profile.ID, []base.VerifiedIdentifier{{
		UID:           resp.UID,
		LoginProvider: base.LoginProviderTypePhone,
		Timestamp:     time.Now().In(base.TimeLocation),
	}}); err != nil {
		return nil, exceptions.UpdateProfileError(err)
	}

	if err := fr.UpdateVerifiedUIDS(ctx, profile.ID, []string{resp.UID}); err != nil {
		return nil, exceptions.UpdateProfileError(err)
	}

	canExperiment, err := fr.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	return &base.AuthCredentialResponse{
		CustomToken:   &customToken,
		IDToken:       &userTokens.IDToken,
		ExpiresIn:     userTokens.ExpiresIn,
		RefreshToken:  userTokens.RefreshToken,
		UID:           resp.UID,
		IsAnonymous:   false,
		IsAdmin:       fr.CheckIfAdmin(profile),
		CanExperiment: canExperiment,
	}, nil
}

// CheckIfAdmin checks if a user has admin permissions
func (fr *Repository) CheckIfAdmin(profile *base.UserProfile) bool {
	if len(profile.Permissions) == 0 {
		return false
	}
	exists := false
	for _, p := range profile.Permissions {
		if p == base.PermissionTypeSuperAdmin || p == base.PermissionTypeAdmin {
			exists = true
			break
		}
	}
	return exists
}

// UpdateUserName updates the username of a profile that matches the id
// this method should be called after asserting the username is unique and not associated with another userProfile
func (fr *Repository) UpdateUserName(ctx context.Context, id string, userName string) error {
	v, err := fr.CheckIfUsernameExists(ctx, userName)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	if v {
		return exceptions.InternalServerError(fmt.Errorf("%v", exceptions.UsernameInUseErrMsg))
	}
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.UserName = &userName
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile primary phone number: %v", err))
	}

	return nil
}

// UpdatePrimaryPhoneNumber append a new primary phone number to the user profile
// this method should be called after asserting the phone number is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.PrimaryPhone = &phoneNumber

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}

	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile primary phone number: %v", err))
	}

	return nil
}

// UpdatePrimaryEmailAddress the primary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddress is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.PrimaryEmailAddress = &emailAddress

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile primary email address: %v", err))
	}

	return nil
}

// UpdateSecondaryPhoneNumbers updates the secondary phone numbers of the profile that matches the id
// this method should be called after asserting the phone numbers are unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	newSecondaryPhoneNumber := []string{}
	// Check if the former primary phone exists in the phoneNumber list
	index, exist := utils.FindItem(profile.SecondaryPhoneNumbers, *profile.PrimaryPhone)
	if exist {
		// Remove the former secondary phone from the list since it's now primary
		profile.SecondaryPhoneNumbers = append(
			profile.SecondaryPhoneNumbers[:index],
			profile.SecondaryPhoneNumbers[index+1:]...,
		)
	}

	profile.SecondaryPhoneNumbers = append(newSecondaryPhoneNumber, phoneNumbers...)

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile secondary phone numbers: %v", err))
	}

	return nil
}

// UpdateSecondaryEmailAddresses the secondary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddresses  as unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryEmailAddresses(ctx context.Context, id string, uniqueEmailAddresses []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	newSecondaryEmail := []string{}
	// check if former primary email still exists in the
	// secondary emails list
	if profile.PrimaryEmailAddress != nil {
		index, exist := utils.FindItem(profile.SecondaryEmailAddresses, *profile.PrimaryEmailAddress)
		if exist {
			// remove the former secondary email from the list
			profile.SecondaryEmailAddresses = append(
				profile.SecondaryEmailAddresses[:index],
				profile.SecondaryEmailAddresses[index+1:]...,
			)
		}
	}

	profile.SecondaryEmailAddresses = append(newSecondaryEmail, uniqueEmailAddresses...)

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile secondary email address: %v", err))
	}
	return nil
}

// UpdateSuspended updates the suspend attribute of the profile that matches the id
func (fr *Repository) UpdateSuspended(ctx context.Context, id string, status bool) error {
	profile, err := fr.GetUserProfileByID(ctx, id, true)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.Suspended = status

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil

}

// UpdatePhotoUploadID updates the photoUploadID attribute of the profile that matches the id
func (fr *Repository) UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.PhotoUploadID = uploadID

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile photo upload id: %v", err))
	}

	return nil
}

// UpdateCovers updates the covers attribute of the profile that matches the id
func (fr *Repository) UpdateCovers(ctx context.Context, id string, covers []base.Cover) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	// check that the new cover been added is unique and does not currently exist in the user's profile.
	newCovers := []base.Cover{}
	if len(profile.Covers) >= 1 {
		for _, cover := range covers {
			if !utils.IfCoverExistsInSlice(profile.Covers, cover) {
				if !utils.IfCoverExistsInSlice(newCovers, cover) {
					newCovers = append(newCovers, cover)
				}
			}
		}
	} else {
		newCovers = append(newCovers, covers...)
	}

	profile.Covers = append(profile.Covers, newCovers...)

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile covers: %v", err))
	}

	return nil
}

// UpdatePushTokens updates the pushTokens attribute of the profile that matches the id. This function does a hard reset instead of prior
// matching
func (fr *Repository) UpdatePushTokens(ctx context.Context, id string, pushTokens []string) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	profile.PushTokens = pushTokens

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile push tokens: %v", err))
	}
	return nil
}

// UpdatePermissions update the permissions of the user profile
func (fr *Repository) UpdatePermissions(ctx context.Context, id string, perms []base.PermissionType) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	newPerms := []base.PermissionType{}
	if len(profile.Permissions) >= 1 {
		for _, perm := range perms {
			for _, current := range profile.Permissions {
				if string(perm) != string(current) {
					newPerms = append(newPerms, perm)
				}
			}
		}
	} else {
		newPerms = append(newPerms, perms...)
	}

	profile.Permissions = newPerms

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile permissions: %v", err))
	}
	return nil

}

// UpdateBioData updates the biodate of the profile that matches the id
func (fr *Repository) UpdateBioData(ctx context.Context, id string, data base.BioData) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	profile.UserBioData.FirstName = func(pr *base.UserProfile, dt base.BioData) *string {
		if dt.FirstName != nil {
			return dt.FirstName
		}
		return pr.UserBioData.FirstName
	}(profile, data)
	profile.UserBioData.LastName = func(pr *base.UserProfile, dt base.BioData) *string {
		if dt.LastName != nil {
			return dt.LastName
		}
		return pr.UserBioData.LastName
	}(profile, data)
	profile.UserBioData.Gender = func(pr *base.UserProfile, dt base.BioData) base.Gender {
		if dt.Gender.String() != "" {
			return dt.Gender
		}
		return pr.UserBioData.Gender
	}(profile, data)
	profile.UserBioData.DateOfBirth = func(pr *base.UserProfile, dt base.BioData) *base.Date {
		if dt.DateOfBirth != nil {
			return dt.DateOfBirth
		}

		return pr.UserBioData.DateOfBirth
	}(profile, data)
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile bio data: %v", err))
	}
	return nil
}

// UpdateVerifiedIdentifiers adds a UID to a user profile during login if it does not exist
func (fr *Repository) UpdateVerifiedIdentifiers(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error {

	for _, identifier := range identifiers {
		// for each run, get the user profile. this will ensure the fetch profile always has the latest data
		profile, err := fr.GetUserProfileByID(ctx, id, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

		if !utils.CheckIdentifierExists(profile, identifier.UID) {
			uids := profile.VerifiedIdentifiers

			uids = append(uids, identifier)

			profile.VerifiedIdentifiers = append(profile.VerifiedIdentifiers, uids...)

			query := &GetAllQuery{
				CollectionName: fr.GetUserProfileCollectionName(),
				FieldName:      "id",
				Value:          profile.ID,
				Operator:       "==",
			}
			docs, err := fr.FirestoreClient.GetAll(ctx, query)
			if err != nil {
				return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
			}
			if len(docs) == 0 {
				return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
			}
			updateCommand := &UpdateCommand{
				CollectionName: fr.GetUserProfileCollectionName(),
				ID:             docs[0].Ref.ID,
				Data:           profile,
			}
			err = fr.FirestoreClient.Update(ctx, updateCommand)
			if err != nil {
				return exceptions.InternalServerError(fmt.Errorf("unable to update user profile verified identifiers: %v", err))
			}
			return nil

		}
	}

	return nil
}

// UpdateVerifiedUIDS adds a UID to a user profile during login if it does not exist
func (fr *Repository) UpdateVerifiedUIDS(ctx context.Context, id string, uids []string) error {

	for _, uid := range uids {
		// for each run, get the user profile. this will ensure the fetch profile always has the latest data
		profile, err := fr.GetUserProfileByID(ctx, id, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

		if !base.StringSliceContains(profile.VerifiedUIDS, uid) {
			uids := []string{}

			uids = append(uids, uid)

			profile.VerifiedUIDS = append(profile.VerifiedUIDS, uids...)

			query := &GetAllQuery{
				CollectionName: fr.GetUserProfileCollectionName(),
				FieldName:      "id",
				Value:          profile.ID,
				Operator:       "==",
			}
			docs, err := fr.FirestoreClient.GetAll(ctx, query)
			if err != nil {
				return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
			}
			if len(docs) == 0 {
				return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
			}
			updateCommand := &UpdateCommand{
				CollectionName: fr.GetUserProfileCollectionName(),
				ID:             docs[0].Ref.ID,
				Data:           profile,
			}
			err = fr.FirestoreClient.Update(ctx, updateCommand)
			if err != nil {
				return exceptions.InternalServerError(fmt.Errorf("unable to update user profile verified UIDS: %v", err))
			}
			return nil

		}
	}

	return nil
}

// RecordPostVisitSurvey records an end of visit survey
func (fr *Repository) RecordPostVisitSurvey(
	ctx context.Context,
	input resources.PostVisitSurveyInput,
	UID string,
) error {
	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return exceptions.LikelyToRecommendError(fmt.Errorf("the likelihood of recommending should be an int between 0 and 10"))

	}
	feedback := domain.PostVisitSurvey{
		LikelyToRecommend: input.LikelyToRecommend,
		Criticism:         input.Criticism,
		Suggestions:       input.Suggestions,
		UID:               UID,
		Timestamp:         time.Now(),
	}
	command := &CreateCommand{
		CollectionName: fr.GetSurveyCollectionName(),
		Data:           feedback,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return exceptions.AddRecordError(err)

	}
	return nil
}

// SavePIN  persist the data of the newly created PIN to a datastore
func (fr *Repository) SavePIN(ctx context.Context, pin *domain.PIN) (bool, error) {
	// persist the data to a datastore
	command := &CreateCommand{
		CollectionName: fr.GetPINsCollectionName(),
		Data:           pin,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return false, exceptions.AddRecordError(err)
	}
	return true, nil

}

// UpdatePIN  persist the data of the updated PIN to a datastore
func (fr *Repository) UpdatePIN(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
	pinData, err := fr.GetPINByProfileID(ctx, id)
	if err != nil {
		return false, exceptions.PinNotFoundError(err)
	}
	query := &GetAllQuery{
		CollectionName: fr.GetPINsCollectionName(),
		FieldName:      "id",
		Value:          pinData.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to parse user pin as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return false, exceptions.InternalServerError(fmt.Errorf("user pin not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetPINsCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           pin,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return false, exceptions.UpdateProfileError(err)
	}

	return true, nil

}

// ExchangeRefreshTokenForIDToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (fr Repository) ExchangeRefreshTokenForIDToken(refreshToken string) (*base.AuthCredentialResponse, error) {
	apiKey, err := base.GetEnvVar(base.FirebaseWebAPIKeyEnvVarName)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	payload := resources.RefreshTokenExchangePayload{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	url := firebaseExchangeRefreshTokenURL + apiKey
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * base.HTTPClientTimeoutSecs
	resp, err := httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(payloadBytes),
	)

	defer base.CloseRespBody(resp)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		return nil,
			exceptions.InternalServerError(fmt.Errorf(
				"firebase HTTP error, status code %d\nBody: %s\nBody read error: %s",
				resp.StatusCode,
				string(bs),
				err,
			))
	}

	var tokenResp base.AuthCredentialResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	return &tokenResp, nil
}

// GetCustomerProfileByID fetch the customer profile by profile id.
func (fr *Repository) GetCustomerProfileByID(ctx context.Context, id string) (*base.Customer, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.InternalServerError(fmt.Errorf("customer profile not found: %w", err))
	}
	dsnap := docs[0]
	cus := &base.Customer{}
	err = dsnap.DataTo(cus)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read customer profile: %w", err))
	}
	return cus, nil
}

// GetCustomerProfileByProfileID fetches customer profile by given ID
func (fr *Repository) GetCustomerProfileByProfileID(ctx context.Context, profileID string) (*base.Customer, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, exceptions.CustomerNotFoundError()
	}
	dsnap := docs[0]
	cus := &base.Customer{}
	err = dsnap.DataTo(cus)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read customer profile: %w", err))
	}
	return cus, nil
}

// GetSupplierProfileByProfileID fetch the supplier profile by profile id.
// since this same supplierProfile can be used for updating, a companion snapshot record is returned as well
func (fr *Repository) GetSupplierProfileByProfileID(
	ctx context.Context,
	profileID string,
) (*base.Supplier, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.SupplierNotFoundError()
	}
	dsnap := docs[0]
	sup := &base.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return sup, nil
}

// GetSupplierProfileByID fetches supplier profile by given ID
func (fr *Repository) GetSupplierProfileByID(ctx context.Context, id string) (*base.Supplier, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, exceptions.InternalServerError(fmt.Errorf("supplier profile not found: %w", err))
	}
	dsnap := docs[0]
	sup := &base.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read supplier profile: %w", err))
	}
	return sup, nil
}

// UpdateSupplierProfile does a generic update of supplier profile.
func (fr *Repository) UpdateSupplierProfile(ctx context.Context, profileID string, data *base.Supplier) error {
	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	sup.PayablesAccount = data.PayablesAccount
	sup.SupplierKYC = data.SupplierKYC
	sup.Active = data.Active
	sup.AccountType = data.AccountType
	sup.UnderOrganization = data.UnderOrganization
	sup.IsOrganizationVerified = data.IsOrganizationVerified
	sup.SladeCode = data.SladeCode
	sup.ParentOrganizationID = data.ParentOrganizationID
	sup.HasBranches = data.HasBranches
	sup.Location = data.Location
	sup.PartnerType = data.PartnerType
	sup.EDIUserProfile = data.EDIUserProfile
	sup.PartnerSetupComplete = data.PartnerSetupComplete
	sup.KYCSubmitted = data.KYCSubmitted
	sup.OrganizationName = data.OrganizationName

	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "id",
		Value:          sup.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           sup,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile: %v", err))
	}
	return nil

}

// AddSupplierAccountType update the supplier profile with the correct account type
func (fr *Repository) AddSupplierAccountType(ctx context.Context, profileID string, accountType base.AccountType) (*base.Supplier, error) {

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	sup.AccountType = &accountType
	sup.UnderOrganization = false
	sup.IsOrganizationVerified = false
	sup.HasBranches = false
	sup.Active = false

	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "id",
		Value:          sup.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return nil, exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           sup,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to update user profile: %v", err))
	}

	return sup, nil

}

// AddPartnerType updates the suppier profile with the provided name and  partner type.
func (fr *Repository) AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *base.PartnerType) (bool, error) {

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}

	sup.SupplierName = *name
	sup.PartnerType = *partnerType
	sup.PartnerSetupComplete = true

	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "id",
		Value:          sup.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return false, exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}

	updateCommand := &UpdateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           sup,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to update user profile: %v", err))
	}

	return true, nil

}

// ActivateSupplierProfile sets the active attribute of supplier profile to true
func (fr *Repository) ActivateSupplierProfile(
	profileID string,
	supplier base.Supplier,
) (*base.Supplier, error) {
	ctx := context.Background()
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	collectionName := fr.GetSupplierProfileCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "id",
		Value:          sup.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, err
	}

	sup.Active = supplier.Active
	sup.PayablesAccount = supplier.PayablesAccount
	sup.SupplierID = supplier.SupplierID

	updateCommand := &UpdateCommand{
		CollectionName: collectionName,
		ID:             docs[0].Ref.ID,
		Data:           sup,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return sup, nil
}

// StageProfileNudge stages nudges published from this service.
func (fr *Repository) StageProfileNudge(
	ctx context.Context,
	nudge *base.Nudge,
) error {
	command := &CreateCommand{
		CollectionName: fr.GetProfileNudgesCollectionName(),
		Data:           nudge,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

// StageKYCProcessingRequest stages the request which will be retrieved later for admins
func (fr *Repository) StageKYCProcessingRequest(ctx context.Context, data *domain.KYCRequest) error {
	command := &CreateCommand{
		CollectionName: fr.GetKCYProcessCollectionName(),
		Data:           data,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

// RemoveKYCProcessingRequest removes the supplier's kyc processing request
func (fr *Repository) RemoveKYCProcessingRequest(ctx context.Context, supplierProfileID string) error {
	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "supplierRecord.id",
		Value:          supplierProfileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to fetch kyc request documents: %v", err))
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("no kyc processing record found: %v", err))
	}

	req := &domain.KYCRequest{}
	if err := docs[0].DataTo(req); err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to read supplier kyc record: %w", err))
	}
	getKYCQuery := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "id",
		Value:          req.ID,
		Operator:       "==",
	}
	if docs, err := fr.FirestoreClient.GetAll(ctx, getKYCQuery); err == nil {
		command := &DeleteCommand{
			CollectionName: fr.GetKCYProcessCollectionName(),
			ID:             docs[0].Ref.ID,
		}
		if err = fr.FirestoreClient.Delete(ctx, command); err != nil {
			return exceptions.InternalServerError(err)
		}
	}
	return nil
}

// FetchKYCProcessingRequests retrieves all unprocessed kycs for admins
func (fr *Repository) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "processed",
		Value:          false,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to fetch kyc request documents: %v", err))
	}

	res := []*domain.KYCRequest{}

	for _, doc := range docs {
		req := &domain.KYCRequest{}
		err = doc.DataTo(req)
		if err != nil {
			return nil, exceptions.InternalServerError(fmt.Errorf("unable to read supplier: %w", err))
		}
		res = append(res, req)
	}

	return res, nil
}

// FetchKYCProcessingRequestByID retrieves a specific kyc processing request
func (fr *Repository) FetchKYCProcessingRequestByID(ctx context.Context, id string) (*domain.KYCRequest, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to fetch kyc request documents: %v", err))
	}

	req := &domain.KYCRequest{}
	err = docs[0].DataTo(req)
	if err != nil {
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read supplier: %w", err))
	}

	return req, nil
}

// UpdateKYCProcessingRequest update the supplier profile
func (fr *Repository) UpdateKYCProcessingRequest(ctx context.Context, kycRequest *domain.KYCRequest) error {
	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "id",
		Value:          kycRequest.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse kyc processing request as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("kyc processing request not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetKCYProcessCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           kycRequest,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update kyc processing request profile: %v", err))
	}
	return nil
}

// FetchAdminUsers fetches all admins
func (fr *Repository) FetchAdminUsers(ctx context.Context) ([]*base.UserProfile, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "permissions",
		Value:          base.DefaultAdminPermissions,
		Operator:       "array-contains-any",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	var admins []*base.UserProfile
	for _, doc := range docs {
		u := &base.UserProfile{}
		err = doc.DataTo(u)
		if err != nil {
			return nil, exceptions.InternalServerError(fmt.Errorf("unable to read user profile: %w", err))
		}
		admins = append(admins, u)
	}
	return admins, nil
}

// PurgeUserByPhoneNumber removes the record of a user given a phone number.
func (fr *Repository) PurgeUserByPhoneNumber(ctx context.Context, phone string) error {
	profile, err := fr.GetUserProfileByPhoneNumber(ctx, phone, false)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	// delete pin of the user
	pin, err := fr.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	query := &GetAllQuery{
		CollectionName: fr.GetPINsCollectionName(),
		FieldName:      "id",
		Value:          pin.ID,
		Operator:       "==",
	}
	if docs, err := fr.FirestoreClient.GetAll(ctx, query); err == nil {
		command := &DeleteCommand{
			CollectionName: fr.GetPINsCollectionName(),
			ID:             docs[0].Ref.ID,
		}
		if err = fr.FirestoreClient.Delete(ctx, command); err != nil {
			return exceptions.InternalServerError(err)
		}
	}
	// delete user supplier profile
	// some old profiles may not have a supplier profile since the original implementation
	// created a supplier profile only for PRO.
	// However the current and correct logic creates a supplier profile regardless of flavour.
	// Hence, the deletion of supplier
	// profile should only occur if a supplier profile exists and not throw an error.
	supplier, err := fr.GetSupplierProfileByProfileID(ctx, profile.ID)
	if err != nil {
		log.Printf("Supplier record was not found: %v", err)
	} else {
		err = fr.RemoveKYCProcessingRequest(ctx, supplier.ID)
		if err != nil {
			log.Printf("KYC request information was not removed %v", err)
		}

		query := &GetAllQuery{
			CollectionName: fr.GetSupplierProfileCollectionName(),
			FieldName:      "id",
			Value:          supplier.ID,
			Operator:       "==",
		}
		if docs, err := fr.FirestoreClient.GetAll(ctx, query); err == nil {
			command := &DeleteCommand{
				CollectionName: fr.GetSupplierProfileCollectionName(),
				ID:             docs[0].Ref.ID,
			}
			if err = fr.FirestoreClient.Delete(ctx, command); err != nil {
				return exceptions.InternalServerError(err)
			}
		}
	}

	// delete user customer profile
	// some old profiles may not have a customer profile since the original implementation
	// created a customer profile only for CONSUMER.
	// However the current and correct logic creates a customer profile regardless of flavour.
	// Hence, the deletion of customer
	// profile should only occur if a customer profile exists and not throw an error.
	customer, err := fr.GetCustomerProfileByProfileID(ctx, profile.ID)
	if err != nil {
		log.Printf("Customer record was not found: %v", err)
	} else {
		query := &GetAllQuery{
			CollectionName: fr.GetCustomerProfileCollectionName(),
			FieldName:      "id",
			Value:          customer.ID,
			Operator:       "==",
		}
		if docs, err := fr.FirestoreClient.GetAll(ctx, query); err == nil {
			command := &DeleteCommand{
				CollectionName: fr.GetCustomerProfileCollectionName(),
				ID:             docs[0].Ref.ID,
			}
			if err = fr.FirestoreClient.Delete(ctx, command); err != nil {
				return exceptions.InternalServerError(err)
			}
		}
	}

	// delete the user profile
	query1 := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	if docs, err := fr.FirestoreClient.GetAll(ctx, query1); err == nil {
		command := &DeleteCommand{
			CollectionName: fr.GetUserProfileCollectionName(),
			ID:             docs[0].Ref.ID,
		}
		if err = fr.FirestoreClient.Delete(ctx, command); err != nil {
			return exceptions.InternalServerError(err)
		}
	}

	// delete the user from firebase
	u, err := fr.FirebaseClient.GetUserByPhoneNumber(ctx, phone)
	if err == nil {
		// only run firebase delete if firebase manages to find the user. It's not fatal if firebase fails to find the user
		if err := fr.FirebaseClient.DeleteUser(ctx, u.UID); err != nil {
			return exceptions.InternalServerError(err)
		}
	}

	return nil
}

// GetCustomerOrSupplierProfileByProfileID returns either a customer or supplier profile
// given the flavour and the profile ID that has been provided
func (fr *Repository) GetCustomerOrSupplierProfileByProfileID(
	ctx context.Context,
	flavour base.Flavour,
	profileID string,
) (*base.Customer, *base.Supplier, error) {
	var customer *base.Customer
	var supplier *base.Supplier

	switch flavour {
	case base.FlavourConsumer:
		{
			customerProfile, err := fr.GetCustomerProfileByProfileID(ctx, profileID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get customer profile")
			}
			customer = customerProfile
		}
	case base.FlavourPro:
		{
			supplierProfile, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get customer profile")
			}
			supplier = supplierProfile
		}
	default:
		return customer, supplier, exceptions.WrongEnumTypeError(flavour.String())
	}

	return customer, supplier, nil
}

// GetOrCreatePhoneNumberUser retrieves or creates an phone number user
// account in Firebase Authentication
func (fr *Repository) GetOrCreatePhoneNumberUser(
	ctx context.Context,
	phone string,
) (*resources.CreatedUserResponse, error) {
	user, err := fr.FirebaseClient.GetUserByPhoneNumber(
		ctx,
		phone,
	)
	if err == nil {
		return &resources.CreatedUserResponse{
			UID:         user.UID,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			PhotoURL:    user.PhotoURL,
			ProviderID:  user.ProviderID,
		}, nil
	}

	params := (&auth.UserToCreate{}).
		PhoneNumber(phone)
	newUser, err := fr.FirebaseClient.CreateUser(
		ctx,
		params,
	)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return &resources.CreatedUserResponse{
		UID:         newUser.UID,
		DisplayName: newUser.DisplayName,
		Email:       newUser.Email,
		PhoneNumber: newUser.PhoneNumber,
		PhotoURL:    newUser.PhotoURL,
		ProviderID:  newUser.ProviderID,
	}, nil
}

// HardResetSecondaryPhoneNumbers does a hard reset of user secondary phone numbers.
// This should be called when retiring specific secondary phone number and passing in
// the new secondary phone numbers as an argument.
func (fr *Repository) HardResetSecondaryPhoneNumbers(
	ctx context.Context,
	profile *base.UserProfile,
	newSecondaryPhoneNumbers []string,
) error {
	profile.SecondaryPhoneNumbers = newSecondaryPhoneNumbers

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile secondary phone numbers: %v", err))
	}

	return nil
}

// HardResetSecondaryEmailAddress does a hard reset of user secondary email addresses. This should be called when retiring specific
// secondary email addresses and passing in the new secondary email address as an argument.
func (fr *Repository) HardResetSecondaryEmailAddress(
	ctx context.Context,
	profile *base.UserProfile,
	newSecondaryEmails []string,
) error {
	profile.SecondaryEmailAddresses = newSecondaryEmails

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("user profile not found"))
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile secondary phone numbers: %v", err))
	}

	return nil
}

// CheckIfExperimentParticipant check if a user has subscribed to be an experiment participant
func (fr *Repository) CheckIfExperimentParticipant(ctx context.Context, profileID string) (bool, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetExperimentParticipantCollectionName(),
		FieldName:      "id",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}

	if len(docs) == 0 {
		return false, nil
	}
	return true, nil
}

// AddUserAsExperimentParticipant adds the provided user profile as an experiment participant if does not already exist.
// this method is idempotent.
func (fr *Repository) AddUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error) {
	exists, err := fr.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		return false, err
	}

	if !exists {
		createCommand := &CreateCommand{
			CollectionName: fr.GetExperimentParticipantCollectionName(),
			Data:           profile,
		}
		_, err = fr.FirestoreClient.Create(ctx, createCommand)
		if err != nil {
			return false, exceptions.InternalServerError(fmt.Errorf("unable to add user profile of ID %v in experiment_participant: %v", profile.ID, err))
		}
		return true, nil
	}
	// the user already exists as an experiment participant
	return true, nil

}

// RemoveUserAsExperimentParticipant removes the provide user profile as an experiment participant. This methold does not check
// for existence before deletion since non-existence is relatively equivalent to a removal
func (fr *Repository) RemoveUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error) {
	// fetch the document References
	query := &GetAllQuery{
		CollectionName: fr.GetExperimentParticipantCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err))
	}
	// means the document was removed or does not exist
	if len(docs) == 0 {
		return true, nil
	}
	deleteCommand := &DeleteCommand{
		CollectionName: fr.GetExperimentParticipantCollectionName(),
		ID:             docs[0].Ref.ID,
	}
	err = fr.FirestoreClient.Delete(ctx, deleteCommand)
	if err != nil {
		return false, exceptions.InternalServerError(fmt.Errorf("unable to remove user profile of ID %v from experiment_participant: %v", profile.ID, err))
	}

	return true, nil
}

// UpdateAddresses persists a user's home or work address information to the database
func (fr *Repository) UpdateAddresses(
	ctx context.Context,
	id string,
	address base.Address,
	addressType base.AddressType,
) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		return err
	}

	switch addressType {
	case base.AddressTypeHome:
		{
			profile.HomeAddress = &address
		}
	case base.AddressTypeWork:
		{
			profile.WorkAddress = &address
		}
	default:
		return exceptions.WrongEnumTypeError(addressType.String())
	}

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

// AddNHIFDetails persists a user's NHIF details
func (fr *Repository) AddNHIFDetails(
	ctx context.Context,
	input resources.NHIFDetailsInput,
	profileID string,
) (*domain.NHIFDetails, error) {
	// Do a check if the item exists
	collectionName := fr.GetNHIFDetailsCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "membershipNumber",
		Value:          input.MembershipNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) > 0 {
		return nil, exceptions.RecordExistsError(fmt.Errorf("record exists"))
	}

	nhifDetails := domain.NHIFDetails{
		ID:                        uuid.New().String(),
		ProfileID:                 profileID,
		MembershipNumber:          input.MembershipNumber,
		Employment:                input.Employment,
		IDDocType:                 input.IDDocType,
		IDNumber:                  input.IDNumber,
		IdentificationCardPhotoID: input.IdentificationCardPhotoID,
		NHIFCardPhotoID:           input.NHIFCardPhotoID,
	}

	createCommand := &CreateCommand{
		CollectionName: collectionName,
		Data:           nhifDetails,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	getNhifQuery := &GetSingleQuery{
		CollectionName: collectionName,
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getNhifQuery)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	nhif := &domain.NHIFDetails{}
	err = dsnap.DataTo(nhif)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	return nhif, nil
}

// GetNHIFDetailsByProfileID fetches a user's NHIF details given their profile ID
func (fr *Repository) GetNHIFDetailsByProfileID(
	ctx context.Context,
	profileID string,
) (*domain.NHIFDetails, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetNHIFDetailsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 NHIF details with profile ID %s (count: %d)",
			profileID,
			len(docs),
		)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	nhif := &domain.NHIFDetails{}
	err = docs[0].DataTo(nhif)
	if err != nil {
		return nil, err
	}

	return nhif, nil
}

// GetUserCommunicationsSettings fetches the communication settings of a specific user.
func (fr *Repository) GetUserCommunicationsSettings(ctx context.Context, profileID string) (*base.UserCommunicationsSetting, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetCommunicationsSettingsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) > 1 && base.IsDebug() {
		log.Printf("> 1 communications settings with profile ID %s (count: %d)",
			profileID,
			len(docs),
		)
	}

	if len(docs) == 0 {
		return &base.UserCommunicationsSetting{ProfileID: profileID}, nil
	}

	comms := &base.UserCommunicationsSetting{}
	err = docs[0].DataTo(comms)
	if err != nil {
		return nil, err
	}
	return comms, nil
}

// SetUserCommunicationsSettings sets communication settings for a specific user
func (fr *Repository) SetUserCommunicationsSettings(ctx context.Context, profileID string,
	allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {

	// get the previous communications_settings
	comms, err := fr.GetUserCommunicationsSettings(ctx, profileID)
	if err != nil {
		return nil, err
	}

	setCommsSettings := base.UserCommunicationsSetting{
		ID:            uuid.New().String(),
		ProfileID:     profileID,
		AllowWhatsApp: utils.MatchAndReturn(comms.AllowWhatsApp, *allowWhatsApp),
		AllowTextSMS:  utils.MatchAndReturn(comms.AllowWhatsApp, *allowTextSms),
		AllowPush:     utils.MatchAndReturn(comms.AllowWhatsApp, *allowPush),
		AllowEmail:    utils.MatchAndReturn(comms.AllowWhatsApp, *allowEmail),
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetCommunicationsSettingsCollectionName(),
		Data:           setCommsSettings,
	}
	_, err = fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	// fetch the now set communications_settings and return it
	return fr.GetUserCommunicationsSettings(ctx, profileID)
}

// UpdateCustomerProfile does a generic update of the customer profile
// to add the data recieved from the ERP.
func (fr *Repository) UpdateCustomerProfile(
	ctx context.Context,
	profileID string,
	cus base.Customer,
) (*base.Customer, error) {
	customer, err := fr.GetCustomerProfileByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	collectionName := fr.GetCustomerProfileCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "id",
		Value:          customer.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, err
	}

	customer.CustomerID = cus.CustomerID
	customer.ReceivablesAccount = cus.ReceivablesAccount
	customer.Active = cus.Active

	updateCommand := &UpdateCommand{
		CollectionName: collectionName,
		ID:             docs[0].Ref.ID,
		Data:           customer,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return customer, nil
}
