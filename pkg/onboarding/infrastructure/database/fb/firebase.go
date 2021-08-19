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

	"go.opentelemetry.io/otel"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

// Package that generates trace information
var tracer = otel.Tracer(
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb",
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
	smsCollectionName                    = "incoming_sms"
	ussdDataCollectioName                = "ussd_data"
	firebaseExchangeRefreshTokenURL      = "https://securetoken.googleapis.com/v1/token?key="
	marketingDataCollectionName          = "marketing_data"
	ussdEventsCollectionName             = "ussd_events"
	coverLinkingEventsCollectionName     = "coverlinking_events"
	rolesCollectionName                  = "user_roles"
)

// Repository accesses and updates an item that is stored on Firebase
type Repository struct {
	FirestoreClient FirestoreClientExtension
	FirebaseClient  FirebaseClientExtension
}

// NewFirebaseRepository initializes a Firebase repository
func NewFirebaseRepository(
	firestoreClient FirestoreClientExtension,
	firebaseClient FirebaseClientExtension,
) repository.OnboardingRepository {
	return &Repository{
		FirestoreClient: firestoreClient,
		FirebaseClient:  firebaseClient,
	}
}

// GetUserProfileCollectionName ...
func (fr Repository) GetUserProfileCollectionName() string {
	suffixed := firebasetools.SuffixCollection(userProfileCollectionName)
	return suffixed
}

// GetSupplierProfileCollectionName ...
func (fr Repository) GetSupplierProfileCollectionName() string {
	suffixed := firebasetools.SuffixCollection(supplierProfileCollectionName)
	return suffixed
}

// GetCustomerProfileCollectionName ...
func (fr Repository) GetCustomerProfileCollectionName() string {
	suffixed := firebasetools.SuffixCollection(customerProfileCollectionName)
	return suffixed
}

// GetSurveyCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetSurveyCollectionName() string {
	suffixed := firebasetools.SuffixCollection(surveyCollectionName)
	return suffixed
}

// GetPINsCollectionName returns a well suffixed PINs collection name
func (fr Repository) GetPINsCollectionName() string {
	suffixed := firebasetools.SuffixCollection(pinsCollectionName)
	return suffixed
}

// GetProfileNudgesCollectionName return the storage location of profile nudges
func (fr Repository) GetProfileNudgesCollectionName() string {
	suffixed := firebasetools.SuffixCollection(profileNudgesCollectionName)
	return suffixed
}

// GetKCYProcessCollectionName fetches collection where kyc processing request will be saved
func (fr Repository) GetKCYProcessCollectionName() string {
	suffixed := firebasetools.SuffixCollection(kycProcessCollectionName)
	return suffixed
}

// GetExperimentParticipantCollectionName fetches the collection where experiment participant will be saved
func (fr *Repository) GetExperimentParticipantCollectionName() string {
	suffixed := firebasetools.SuffixCollection(experimentParticipantCollectionName)
	return suffixed
}

// GetNHIFDetailsCollectionName ...
func (fr Repository) GetNHIFDetailsCollectionName() string {
	suffixed := firebasetools.SuffixCollection(nhifDetailsCollectionName)
	return suffixed
}

// GetCommunicationsSettingsCollectionName ...
func (fr Repository) GetCommunicationsSettingsCollectionName() string {
	suffixed := firebasetools.SuffixCollection(communicationsSettingsCollectionName)
	return suffixed
}

// GetSMSCollectionName gets the collection name from firestore
func (fr Repository) GetSMSCollectionName() string {
	suffixed := firebasetools.SuffixCollection(smsCollectionName)
	return suffixed
}

//GetUSSDDataCollectionName gets the collection from firestore
func (fr Repository) GetUSSDDataCollectionName() string {
	suffixed := firebasetools.SuffixCollection(ussdDataCollectioName)
	return suffixed
}

//GetUSSDEventsCollectionName ...
func (fr Repository) GetUSSDEventsCollectionName() string {
	suffixed := firebasetools.SuffixCollection(ussdEventsCollectionName)
	return suffixed
}

// GetMarketingDataCollectionName ...
func (fr Repository) GetMarketingDataCollectionName() string {
	suffixed := firebasetools.SuffixCollection(marketingDataCollectionName)
	return suffixed
}

// GetCoverLinkingEventsCollectionName ...
func (fr Repository) GetCoverLinkingEventsCollectionName() string {
	suffixed := firebasetools.SuffixCollection(coverLinkingEventsCollectionName)
	return suffixed
}

// GetRolesCollectionName ...
func (fr Repository) GetRolesCollectionName() string {
	suffixed := firebasetools.SuffixCollection(rolesCollectionName)
	return suffixed
}

// GetUserProfileByUID retrieves the user profile by UID
func (fr *Repository) GetUserProfileByUID(
	ctx context.Context,
	uid string,
	suspended bool,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByUID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "verifiedUIDS",
		Value:          uid,
		Operator:       "array-contains",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) == 0 {
		err = exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))

		utils.RecordSpanError(span, err)
		return nil, err
	}

	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("user with uids %s has > 1 profile (they have %d)",
			uid,
			len(docs),
		)
	}

	dsnap := docs[0]
	userProfile := &profileutils.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		utils.RecordSpanError(span, err)
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

//GetUserProfileByPhoneOrEmail retrieves user profile by email adddress
func (fr *Repository) GetUserProfileByPhoneOrEmail(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByPhoneOrEmail")
	defer span.End()

	if payload.PhoneNumber == nil {
		query := &GetAllQuery{
			CollectionName: fr.GetUserProfileCollectionName(),
			FieldName:      "primaryEmailAddress",
			Value:          payload.Email,
			Operator:       "==",
		}

		docs, err := fr.FirestoreClient.GetAll(ctx, query)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, exceptions.InternalServerError(err)
		}

		if len(docs) == 0 {
			query := &GetAllQuery{
				CollectionName: fr.GetUserProfileCollectionName(),
				FieldName:      "secondaryEmailAddresses",
				Value:          payload.Email,
				Operator:       "array-contains",
			}

			docs, err := fr.FirestoreClient.GetAll(ctx, query)
			if err != nil {
				utils.RecordSpanError(span, err)
				return nil, exceptions.InternalServerError(err)
			}

			if len(docs) == 0 {
				err = exceptions.ProfileNotFoundError(err)

				utils.RecordSpanError(span, err)
				return nil, err
			}

			dsnap := docs[0]
			userProfile := &profileutils.UserProfile{}
			err = dsnap.DataTo(userProfile)
			if err != nil {
				utils.RecordSpanError(span, err)
				err = fmt.Errorf("unable to read user profile")
				return nil, exceptions.InternalServerError(err)
			}

			return userProfile, nil
		}

		dsnap := docs[0]
		userProfile := &profileutils.UserProfile{}
		err = dsnap.DataTo(userProfile)
		if err != nil {
			utils.RecordSpanError(span, err)
			err = fmt.Errorf("unable to read user profile")
			return nil, exceptions.InternalServerError(err)
		}

		return userProfile, nil
	}

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          payload.PhoneNumber,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		query := &GetAllQuery{
			CollectionName: fr.GetUserProfileCollectionName(),
			FieldName:      "secondaryPhoneNumbers",
			Value:          payload.PhoneNumber,
			Operator:       "array-contains",
		}

		docs, err := fr.FirestoreClient.GetAll(ctx, query)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, exceptions.InternalServerError(err)
		}

		if len(docs) == 0 {
			err = exceptions.ProfileNotFoundError(err)

			utils.RecordSpanError(span, err)
			return nil, err
		}

		dsnap := docs[0]
		userProfile := &profileutils.UserProfile{}
		err = dsnap.DataTo(userProfile)
		if err != nil {
			utils.RecordSpanError(span, err)
			err = fmt.Errorf("unable to read user profile")
			return nil, exceptions.InternalServerError(err)
		}

		return userProfile, nil
	}

	dsnap := docs[0]
	userProfile := &profileutils.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		utils.RecordSpanError(span, err)
		err = fmt.Errorf("unable to read user profile")
		return nil, exceptions.InternalServerError(err)
	}

	return userProfile, nil
}

// UpdateUserProfileEmail updates user profile's email
func (fr *Repository) UpdateUserProfileEmail(
	ctx context.Context,
	phone string,
	email string,
) error {
	ctx, span := tracer.Start(ctx, "UpdateUserProfileEmail")
	defer span.End()

	payload := &dto.RetrieveUserProfileInput{
		PhoneNumber: &phone,
	}

	profile, err := fr.GetUserProfileByPhoneOrEmail(ctx, payload)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	profile.PrimaryEmailAddress = &email

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          profile.PrimaryPhone,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile primary phone number: %v", err),
		)
	}

	return nil
}

// GetUserProfileByID retrieves a user profile by ID
func (fr *Repository) GetUserProfileByID(
	ctx context.Context,
	id string,
	suspended bool,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
	}
	dsnap := docs[0]
	userProfile := &profileutils.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read user profile: %w", err),
		)
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
func (fr *Repository) CreateUserProfile(
	ctx context.Context,
	phoneNumber, uid string,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "CreateUserProfile")
	defer span.End()

	v, err := fr.CheckIfPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to check the phone number: %v", err),
		)
	}

	if v {
		// this phone is number is associated with another user profile, hence can not create an profile with the same phone number
		return nil, exceptions.CheckPhoneNumberExistError()
	}

	profileID := uuid.New().String()
	pr := &profileutils.UserProfile{
		ID:           profileID,
		UserName:     fr.fetchUserRandomName(ctx),
		PrimaryPhone: &phoneNumber,
		VerifiedIdentifiers: []profileutils.VerifiedIdentifier{{
			UID:           uid,
			LoginProvider: profileutils.LoginProviderTypePhone,
			Timestamp:     time.Now().In(pubsubtools.TimeLocation),
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to create new user profile: %w", err),
		)
	}
	query := &GetSingleQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to retrieve newly created user profile: %w", err),
		)
	}
	// return the newly created user profile
	userProfile := &profileutils.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read user profile: %w", err),
		)
	}
	return userProfile, nil

}

// CreateDetailedUserProfile creates a new user profile that is pre-filled using the provided phone number
func (fr *Repository) CreateDetailedUserProfile(
	ctx context.Context,
	phoneNumber string,
	profile profileutils.UserProfile,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "CreateDetailedUserProfile")
	defer span.End()

	exists, err := fr.CheckIfPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("failed to check the phone number: %v", err),
		)
	}

	if exists {
		// this phone is number is associated with another user profile, hence can not create an profile with the same phone number
		err = exceptions.CheckPhoneNumberExistError()
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// create user via their phone number on firebase
	user, err := fr.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	phoneIdentifier := profileutils.VerifiedIdentifier{
		UID:           user.UID,
		LoginProvider: profileutils.LoginProviderTypePhone,
		Timestamp:     time.Now().In(pubsubtools.TimeLocation),
	}

	profile.VerifiedIdentifiers = append(profile.VerifiedIdentifiers, phoneIdentifier)
	profile.VerifiedUIDS = append(profile.VerifiedUIDS, user.UID)

	profileID := uuid.New().String()
	profile.ID = profileID
	profile.PrimaryPhone = &phoneNumber
	profile.UserName = fr.fetchUserRandomName(ctx)
	profile.TermsAccepted = true
	profile.Suspended = false

	command := &CreateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		Data:           profile,
	}

	_, err = fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to create new user profile: %w", err),
		)
	}

	return &profile, nil
}

// CreateEmptySupplierProfile creates an empty supplier profile
func (fr *Repository) CreateEmptySupplierProfile(
	ctx context.Context,
	profileID string,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "CreateEmptySupplierProfile")
	defer span.End()

	sup := &profileutils.Supplier{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		Data:           sup,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to create new supplier empty profile: %w", err),
		)
	}
	getSupplierquery := &GetSingleQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getSupplierquery)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to retrieve newly created supplier profile: %w", err),
		)
	}
	// return the newly created supplier profile
	supplier := &profileutils.Supplier{}
	err = dsnap.DataTo(supplier)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read supplier profile: %w", err),
		)
	}
	return supplier, nil

}

// CreateDetailedSupplierProfile create a new supplier profile that is pre-filled using the provided profile ID
func (fr *Repository) CreateDetailedSupplierProfile(
	ctx context.Context,
	profileID string,
	supplier profileutils.Supplier,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "CreateDetailedSupplierProfile")
	defer span.End()

	supplierID := uuid.New().String()
	supplier.ID = supplierID
	supplier.ProfileID = &profileID
	supplier.Active = true

	createCommand := &CreateCommand{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		Data:           supplier,
	}

	_, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to create new supplier empty profile: %w", err),
		)
	}

	return &supplier, nil
}

// CreateEmptyCustomerProfile creates an empty customer profile
func (fr *Repository) CreateEmptyCustomerProfile(
	ctx context.Context,
	profileID string,
) (*profileutils.Customer, error) {
	ctx, span := tracer.Start(ctx, "CreateEmptyCustomerProfile")
	defer span.End()

	cus := &profileutils.Customer{
		ID:        uuid.New().String(),
		ProfileID: &profileID,
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		Data:           cus,
	}
	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to create new customer empty profile: %w", err),
		)
	}

	getSupplierquery := &GetSingleQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getSupplierquery)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to retrieve newly created customer profile: %w", err),
		)
	}

	// return the newly created customer profile
	customer := &profileutils.Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read customer profile: %w", err),
		)
	}
	return customer, nil
}

//GetUserProfileByPrimaryPhoneNumber fetches a user profile by primary phone number
func (fr *Repository) GetUserProfileByPrimaryPhoneNumber(
	ctx context.Context,
	phoneNumber string,
	suspended bool,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByPrimaryPhoneNumber")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) == 0 {
		return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
	}
	dsnap := docs[0]
	profile := &profileutils.UserProfile{}
	err = dsnap.DataTo(profile)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read user profile: %w", err),
		)
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
func (fr *Repository) GetUserProfileByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
	suspended bool,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfileByPhoneNumber")
	defer span.End()

	// check first primary phone numbers
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) == 1 {
		dsnap := docs[0]
		pr := &profileutils.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, exceptions.InternalServerError(
				fmt.Errorf("unable to read customer profile: %w", err),
			)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs1) == 1 {
		dsnap := docs1[0]
		pr := &profileutils.UserProfile{}
		if err := dsnap.DataTo(pr); err != nil {
			return nil, exceptions.InternalServerError(
				fmt.Errorf("unable to read customer profile: %w", err),
			)
		}

		if !suspended {
			// never return a suspended user profile
			if pr.Suspended {
				return nil, exceptions.ProfileSuspendFoundError()
			}
		}

		return pr, nil
	}

	return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))

}

// CheckIfPhoneNumberExists checks both PRIMARY PHONE NUMBERs and SECONDARY PHONE NUMBERs for the
// existence of the argument phoneNumber.
func (fr *Repository) CheckIfPhoneNumberExists(
	ctx context.Context,
	phoneNumber string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckIfPhoneNumberExists")
	defer span.End()

	// check first primary phone numbers
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryPhone",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
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
	ctx, span := tracer.Start(ctx, "CheckIfEmailExists")
	defer span.End()

	// check first primary email
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "primaryEmailAddress",
		Value:          email,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
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
	ctx, span := tracer.Start(ctx, "CheckIfUsernameExists")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "userName",
		Value:          userName,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}
	if len(docs) == 1 {
		return true, nil
	}

	return false, nil
}

// GetPINByProfileID gets a user's PIN by their profile ID
func (fr *Repository) GetPINByProfileID(
	ctx context.Context,
	profileID string,
) (*domain.PIN, error) {
	ctx, span := tracer.Start(ctx, "GetPINByProfileID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetPINsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	// this should never run. If it does, it means we are doing something wrong.
	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("> 1 PINs with profile ID %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.PinNotFoundError(fmt.Errorf("failed to get a user pin"))
	}

	dsnap := docs[0]
	PIN := &domain.PIN{}
	err = dsnap.DataTo(PIN)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return PIN, nil
}

// GenerateAuthCredentialsForAnonymousUser generates auth credentials for the anonymous user. This method is here since we don't
// want to delegate sign-in of anonymous users to the frontend. This is an effort the over reliance on firebase and lettin us
// handle all the heavy lifting
func (fr *Repository) GenerateAuthCredentialsForAnonymousUser(
	ctx context.Context,
) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "GenerateAuthCredentialsForAnonymousUser")
	defer span.End()

	// todo(dexter) : move anonymousPhoneNumber to base. AnonymousPhoneNumber should NEVER NEVER have a user profile
	anonymousPhoneNumber := "+254700000000"

	u, err := fr.GetOrCreatePhoneNumberUser(ctx, anonymousPhoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	customToken, err := firebasetools.CreateFirebaseCustomToken(ctx, u.UID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.CustomTokenError(err)
	}
	userTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.AuthenticateTokenError(err)
	}

	return &profileutils.AuthCredentialResponse{
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
	profile *profileutils.UserProfile,
) (*profileutils.AuthCredentialResponse, error) {
	ctx, span := tracer.Start(ctx, "GenerateAuthCredentials")
	defer span.End()

	resp, err := fr.GetOrCreatePhoneNumberUser(ctx, phone)
	if err != nil {
		utils.RecordSpanError(span, err)
		if auth.IsUserNotFound(err) {
			return nil, exceptions.UserNotFoundError(err)
		}
		return nil, exceptions.UserNotFoundError(err)
	}

	customToken, err := firebasetools.CreateFirebaseCustomToken(ctx, resp.UID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.CustomTokenError(err)
	}
	userTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.AuthenticateTokenError(err)
	}

	if err := fr.UpdateVerifiedIdentifiers(ctx, profile.ID, []profileutils.VerifiedIdentifier{{
		UID:           resp.UID,
		LoginProvider: profileutils.LoginProviderTypePhone,
		Timestamp:     time.Now().In(pubsubtools.TimeLocation),
	}}); err != nil {
		return nil, exceptions.UpdateProfileError(err)
	}

	if err := fr.UpdateVerifiedUIDS(ctx, profile.ID, []string{resp.UID}); err != nil {
		return nil, exceptions.UpdateProfileError(err)
	}

	canExperiment, err := fr.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	return &profileutils.AuthCredentialResponse{
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
func (fr *Repository) CheckIfAdmin(profile *profileutils.UserProfile) bool {
	if len(profile.Permissions) == 0 {
		return false
	}
	exists := false
	for _, p := range profile.Permissions {
		if p == profileutils.PermissionTypeSuperAdmin || p == profileutils.PermissionTypeAdmin {
			exists = true
			break
		}
	}
	return exists
}

// UpdateUserName updates the username of a profile that matches the id
// this method should be called after asserting the username is unique and not associated with another userProfile
func (fr *Repository) UpdateUserName(ctx context.Context, id string, userName string) error {
	ctx, span := tracer.Start(ctx, "UpdateUserName")
	defer span.End()

	v, err := fr.CheckIfUsernameExists(ctx, userName)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	if v {
		return exceptions.InternalServerError(fmt.Errorf("%v", exceptions.UsernameInUseErrMsg))
	}
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile primary phone number: %v", err),
		)
	}

	return nil
}

// UpdatePrimaryPhoneNumber append a new primary phone number to the user profile
// this method should be called after asserting the phone number is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryPhoneNumber(
	ctx context.Context,
	id string,
	phoneNumber string,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePrimaryPhoneNumber")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile primary phone number: %v", err),
		)
	}

	return nil
}

// UpdateUserRoleIDs updates the roles for a user
func (fr Repository) UpdateUserRoleIDs(ctx context.Context, id string, roleIDs []string) error {
	ctx, span := tracer.Start(ctx, "UpdateUserRoleIDs")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	// Add the roles
	profile.Roles = roleIDs

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile primary email address: %v", err),
		)
	}

	return nil
}

// UpdatePrimaryEmailAddress the primary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddress is unique and not associated with another userProfile
func (fr *Repository) UpdatePrimaryEmailAddress(
	ctx context.Context,
	id string,
	emailAddress string,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePrimaryEmailAddress")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile primary email address: %v", err),
		)
	}

	return nil
}

// UpdateSecondaryPhoneNumbers updates the secondary phone numbers of the profile that matches the id
// this method should be called after asserting the phone numbers are unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryPhoneNumbers(
	ctx context.Context,
	id string,
	phoneNumbers []string,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSecondaryPhoneNumbers")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	// Check if the former primary phone exists in the phoneNumber list
	index, exist := utils.FindItem(profile.SecondaryPhoneNumbers, *profile.PrimaryPhone)
	if exist {
		// Remove the former secondary phone from the list since it's now primary
		profile.SecondaryPhoneNumbers = append(
			profile.SecondaryPhoneNumbers[:index],
			profile.SecondaryPhoneNumbers[index+1:]...,
		)
	}

	for _, phone := range phoneNumbers {
		index, exist := utils.FindItem(profile.SecondaryPhoneNumbers, phone)
		if exist {
			profile.SecondaryPhoneNumbers = append(
				profile.SecondaryPhoneNumbers[:index],
				profile.SecondaryPhoneNumbers[index+1:]...)
		}
	}

	profile.SecondaryPhoneNumbers = append(profile.SecondaryPhoneNumbers, phoneNumbers...)

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile secondary phone numbers: %v", err),
		)
	}

	return nil
}

// UpdateSecondaryEmailAddresses the secondary email addresses of the profile that matches the id
// this method should be called after asserting the emailAddresses  as unique and not associated with another userProfile
func (fr *Repository) UpdateSecondaryEmailAddresses(
	ctx context.Context,
	id string,
	uniqueEmailAddresses []string,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSecondaryEmailAddresses")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	// check if former primary email still exists in the
	// secondary emails list
	if profile.PrimaryEmailAddress != nil {
		index, exist := utils.FindItem(
			profile.SecondaryEmailAddresses,
			*profile.PrimaryEmailAddress,
		)
		if exist {
			// remove the former secondary email from the list
			profile.SecondaryEmailAddresses = append(
				profile.SecondaryEmailAddresses[:index],
				profile.SecondaryEmailAddresses[index+1:]...,
			)
		}
	}

	// Check to see whether the new emails exist in list of secondary emails
	for _, email := range uniqueEmailAddresses {
		index, exist := utils.FindItem(profile.SecondaryEmailAddresses, email)
		if exist {
			profile.SecondaryEmailAddresses = append(
				profile.SecondaryEmailAddresses[:index],
				profile.SecondaryEmailAddresses[index+1:]...)
		}
	}

	profile.SecondaryEmailAddresses = append(
		profile.SecondaryEmailAddresses,
		uniqueEmailAddresses...)

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile secondary email address: %v", err),
		)
	}
	return nil
}

// UpdateSuspended updates the suspend attribute of the profile that matches the id
func (fr *Repository) UpdateSuspended(ctx context.Context, id string, status bool) error {
	ctx, span := tracer.Start(ctx, "UpdateSuspended")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, true)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}

	return nil

}

// UpdatePhotoUploadID updates the photoUploadID attribute of the profile that matches the id
func (fr *Repository) UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error {
	ctx, span := tracer.Start(ctx, "UpdatePhotoUploadID")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile photo upload id: %v", err),
		)
	}

	return nil
}

// UpdateCovers updates the covers attribute of the profile that matches the id
func (fr *Repository) UpdateCovers(
	ctx context.Context,
	id string,
	covers []profileutils.Cover,
) error {
	ctx, span := tracer.Start(ctx, "UpdateCovers")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	// check that the new cover been added is unique and does not currently exist in the user's profile.
	newCovers := []profileutils.Cover{}
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile covers: %v", err),
		)
	}

	return nil
}

// UpdatePushTokens updates the pushTokens attribute of the profile that matches the id. This function does a hard reset instead of prior
// matching
func (fr *Repository) UpdatePushTokens(ctx context.Context, id string, pushTokens []string) error {
	ctx, span := tracer.Start(ctx, "UpdatePushTokens")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile push tokens: %v", err),
		)
	}
	return nil
}

// UpdatePermissions update the permissions of the user profile
func (fr *Repository) UpdatePermissions(
	ctx context.Context,
	id string,
	perms []profileutils.PermissionType,
) error {
	ctx, span := tracer.Start(ctx, "UpdatePermissions")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	// Removes duplicate permissions from array
	// Used for cleaning existing records
	profile.Permissions = utils.UniquePermissionsArray(profile.Permissions)

	newPerms := []profileutils.PermissionType{}
	// Check if has perms
	if len(profile.Permissions) >= 1 {
		// copy the existing perms
		newPerms = append(newPerms, profile.Permissions...)

		for _, perm := range perms {
			// add permission if it doesn't exist
			if !profile.HasPermission(perm) {
				newPerms = append(newPerms, perm)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile permissions: %v", err),
		)
	}
	return nil

}

// UpdateRole update the permissions of the user profile
func (fr *Repository) UpdateRole(ctx context.Context, id string, role profileutils.RoleType) error {
	ctx, span := tracer.Start(ctx, "UpdateRole")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	profile.Role = role
	profile.Permissions = role.Permissions()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user role and permissions: %v", err),
		)
	}
	return nil

}

// UpdateFavNavActions update the permissions of the user profile
func (fr *Repository) UpdateFavNavActions(
	ctx context.Context,
	id string,
	favActions []string,
) error {
	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	profile.FavNavActions = favActions

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user favorite actions: %v", err),
		)
	}
	return nil
}

// UpdateBioData updates the biodate of the profile that matches the id
func (fr *Repository) UpdateBioData(
	ctx context.Context,
	id string,
	data profileutils.BioData,
) error {
	ctx, span := tracer.Start(ctx, "UpdateBioData")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	profile.UserBioData.FirstName = func(pr *profileutils.UserProfile, dt profileutils.BioData) *string {
		if dt.FirstName != nil {
			return dt.FirstName
		}
		return pr.UserBioData.FirstName
	}(
		profile,
		data,
	)
	profile.UserBioData.LastName = func(pr *profileutils.UserProfile, dt profileutils.BioData) *string {
		if dt.LastName != nil {
			return dt.LastName
		}
		return pr.UserBioData.LastName
	}(
		profile,
		data,
	)
	profile.UserBioData.Gender = func(pr *profileutils.UserProfile, dt profileutils.BioData) enumutils.Gender {
		if dt.Gender.String() != "" {
			return dt.Gender
		}
		return pr.UserBioData.Gender
	}(
		profile,
		data,
	)
	profile.UserBioData.DateOfBirth = func(pr *profileutils.UserProfile, dt profileutils.BioData) *scalarutils.Date {
		if dt.DateOfBirth != nil {
			return dt.DateOfBirth
		}

		return pr.UserBioData.DateOfBirth
	}(
		profile,
		data,
	)
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile bio data: %v", err),
		)
	}
	return nil
}

// UpdateVerifiedIdentifiers adds a UID to a user profile during login if it does not exist
func (fr *Repository) UpdateVerifiedIdentifiers(
	ctx context.Context,
	id string,
	identifiers []profileutils.VerifiedIdentifier,
) error {
	ctx, span := tracer.Start(ctx, "UpdateVerifiedIdentifiers")
	defer span.End()

	for _, identifier := range identifiers {
		// for each run, get the user profile. this will ensure the fetch profile always has the latest data
		profile, err := fr.GetUserProfileByID(ctx, id, false)
		if err != nil {
			utils.RecordSpanError(span, err)
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
				utils.RecordSpanError(span, err)
				return exceptions.InternalServerError(
					fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
				)
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
				utils.RecordSpanError(span, err)
				return exceptions.InternalServerError(
					fmt.Errorf("unable to update user profile verified identifiers: %v", err),
				)
			}
			return nil

		}
	}

	return nil
}

// UpdateVerifiedUIDS adds a UID to a user profile during login if it does not exist
func (fr *Repository) UpdateVerifiedUIDS(ctx context.Context, id string, uids []string) error {
	ctx, span := tracer.Start(ctx, "UpdateVerifiedUIDS")
	defer span.End()

	for _, uid := range uids {
		// for each run, get the user profile. this will ensure the fetch profile always has the latest data
		profile, err := fr.GetUserProfileByID(ctx, id, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			// this is a wrapped error. No need to wrap it again
			return err
		}

		if !converterandformatter.StringSliceContains(profile.VerifiedUIDS, uid) {
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
				utils.RecordSpanError(span, err)
				return exceptions.InternalServerError(
					fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
				)
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
				utils.RecordSpanError(span, err)
				return exceptions.InternalServerError(
					fmt.Errorf("unable to update user profile verified UIDS: %v", err),
				)
			}
			return nil

		}
	}

	return nil
}

// RecordPostVisitSurvey records an end of visit survey
func (fr *Repository) RecordPostVisitSurvey(
	ctx context.Context,
	input dto.PostVisitSurveyInput,
	UID string,
) error {
	ctx, span := tracer.Start(ctx, "RecordPostVisitSurvey")
	defer span.End()

	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return exceptions.LikelyToRecommendError(
			fmt.Errorf("the likelihood of recommending should be an int between 0 and 10"),
		)

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
		utils.RecordSpanError(span, err)
		return exceptions.AddRecordError(err)

	}
	return nil
}

// SavePIN  persist the data of the newly created PIN to a datastore
func (fr *Repository) SavePIN(ctx context.Context, pin *domain.PIN) (bool, error) {
	ctx, span := tracer.Start(ctx, "SavePin")
	defer span.End()

	// persist the data to a datastore
	command := &CreateCommand{
		CollectionName: fr.GetPINsCollectionName(),
		Data:           pin,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.AddRecordError(err)
	}
	return true, nil

}

// UpdatePIN  persist the data of the updated PIN to a datastore
func (fr *Repository) UpdatePIN(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
	ctx, span := tracer.Start(ctx, "UpdatePIN")
	defer span.End()

	pinData, err := fr.GetPINByProfileID(ctx, id)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf("unable to parse user pin as firebase snapshot: %v", err),
		)
	}
	if len(docs) == 0 {
		return false, exceptions.InternalServerError(fmt.Errorf("user pin not found"))
	}

	// Check if PIN being updated is a Temporary PIN
	if pinData.IsOTP {
		// Set New PIN flag as false
		pin.IsOTP = false
	}

	updateCommand := &UpdateCommand{
		CollectionName: fr.GetPINsCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           pin,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.UpdateProfileError(err)
	}

	return true, nil

}

// ExchangeRefreshTokenForIDToken takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (fr Repository) ExchangeRefreshTokenForIDToken(
	ctx context.Context,
	refreshToken string,
) (*profileutils.AuthCredentialResponse, error) {
	_, span := tracer.Start(ctx, "ExchangeRefreshTokenForIDToken")
	defer span.End()

	apiKey, err := serverutils.GetEnvVar(firebasetools.FirebaseWebAPIKeyEnvVarName)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	payload := dto.RefreshTokenExchangePayload{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	url := firebaseExchangeRefreshTokenURL + apiKey
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * firebasetools.HTTPClientTimeoutSecs
	resp, err := httpClient.Post(
		url,
		"application/json",
		bytes.NewReader(payloadBytes),
	)

	defer firebasetools.CloseRespBody(resp)
	if err != nil {
		utils.RecordSpanError(span, err)
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

	type refreshTokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		UserID       string `json:"user_id"`
		ProjectID    string `json:"project_id"`
	}

	var tokenResponse refreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(fmt.Errorf(
			"failed to decode refresh token response: %s", err,
		))
	}

	profile, err := fr.GetUserProfileByUID(ctx, tokenResponse.UserID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(fmt.Errorf(
			"failed to retrieve user profile: %s", err,
		))
	}

	canExperiment, err := fr.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(fmt.Errorf(
			"failed to check if the logged in user is an experimental participant: %s", err,
		))
	}

	return &profileutils.AuthCredentialResponse{
		IDToken:       &tokenResponse.IDToken,
		ExpiresIn:     tokenResponse.ExpiresIn,
		RefreshToken:  tokenResponse.RefreshToken,
		UID:           tokenResponse.UserID,
		CanExperiment: canExperiment,
	}, nil
}

// GetCustomerProfileByID fetch the customer profile by profile id.
func (fr *Repository) GetCustomerProfileByID(
	ctx context.Context,
	id string,
) (*profileutils.Customer, error) {
	ctx, span := tracer.Start(ctx, "GetCustomerProfileByID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", id, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("customer profile not found: %w", err),
		)
	}
	dsnap := docs[0]
	cus := &profileutils.Customer{}
	err = dsnap.DataTo(cus)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read customer profile: %w", err),
		)
	}
	return cus, nil
}

// GetCustomerProfileByProfileID fetches customer profile by given ID
func (fr *Repository) GetCustomerProfileByProfileID(
	ctx context.Context,
	profileID string,
) (*profileutils.Customer, error) {
	ctx, span := tracer.Start(ctx, "GetCustomerProfileByProfileID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetCustomerProfileCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, exceptions.CustomerNotFoundError()
	}
	dsnap := docs[0]
	cus := &profileutils.Customer{}
	err = dsnap.DataTo(cus)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read customer profile: %w", err),
		)
	}
	return cus, nil
}

// GetSupplierProfileByProfileID fetch the supplier profile by profile id.
// since this same supplierProfile can be used for updating, a companion snapshot record is returned as well
func (fr *Repository) GetSupplierProfileByProfileID(
	ctx context.Context,
	profileID string,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "GetSupplierProfileByProfileID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("> 1 profile with id %s (count: %d)", profileID, len(docs))
	}

	if len(docs) == 0 {
		return nil, exceptions.SupplierNotFoundError()
	}
	dsnap := docs[0]
	sup := &profileutils.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return sup, nil
}

// GetSupplierProfileByID fetches supplier profile by given ID
func (fr *Repository) GetSupplierProfileByID(
	ctx context.Context,
	id string,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "GetSupplierProfileByID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetSupplierProfileCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, exceptions.InternalServerError(
			fmt.Errorf("supplier profile not found: %w", err),
		)
	}
	dsnap := docs[0]
	sup := &profileutils.Supplier{}
	err = dsnap.DataTo(sup)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to read supplier profile: %w", err),
		)
	}
	return sup, nil
}

// UpdateSupplierProfile does a generic update of supplier profile.
func (fr *Repository) UpdateSupplierProfile(
	ctx context.Context,
	profileID string,
	data *profileutils.Supplier,
) error {
	ctx, span := tracer.Start(ctx, "UpdateSupplierProfile")
	defer span.End()

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(fmt.Errorf("unable to update user profile: %v", err))
	}
	return nil

}

// AddSupplierAccountType update the supplier profile with the correct account type
func (fr *Repository) AddSupplierAccountType(
	ctx context.Context,
	profileID string,
	accountType profileutils.AccountType,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "AddSupplierAccountType")
	defer span.End()

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile: %v", err),
		)
	}

	return sup, nil

}

// AddPartnerType updates the suppier profile with the provided name and  partner type.
func (fr *Repository) AddPartnerType(
	ctx context.Context,
	profileID string,
	name *string,
	partnerType *profileutils.PartnerType,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "AddPartnerType")
	defer span.End()

	// get the suppier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile: %v", err),
		)
	}

	return true, nil

}

// ActivateSupplierProfile sets the active attribute of supplier profile to true
func (fr *Repository) ActivateSupplierProfile(
	ctx context.Context,
	profileID string,
	supplier profileutils.Supplier,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "AddPartnerType")
	defer span.End()

	sup, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return sup, nil
}

// StageProfileNudge stages nudges published from this service.
func (fr *Repository) StageProfileNudge(
	ctx context.Context,
	nudge *feedlib.Nudge,
) error {
	ctx, span := tracer.Start(ctx, "StageProfileNudge")
	defer span.End()

	command := &CreateCommand{
		CollectionName: fr.GetProfileNudgesCollectionName(),
		Data:           nudge,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	return nil
}

// StageKYCProcessingRequest stages the request which will be retrieved later for admins
func (fr *Repository) StageKYCProcessingRequest(
	ctx context.Context,
	data *domain.KYCRequest,
) error {
	ctx, span := tracer.Start(ctx, "StageKYCProcessingRequest")
	defer span.End()

	command := &CreateCommand{
		CollectionName: fr.GetKCYProcessCollectionName(),
		Data:           data,
	}
	_, err := fr.FirestoreClient.Create(ctx, command)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	return nil
}

// RemoveKYCProcessingRequest removes the supplier's kyc processing request
func (fr *Repository) RemoveKYCProcessingRequest(
	ctx context.Context,
	supplierProfileID string,
) error {
	ctx, span := tracer.Start(ctx, "RemoveKYCProcessingRequest")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "supplierRecord.id",
		Value:          supplierProfileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to fetch kyc request documents: %v", err),
		)
	}

	if len(docs) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("no kyc processing record found: %v", err))
	}

	req := &domain.KYCRequest{}
	if err := docs[0].DataTo(req); err != nil {
		return exceptions.InternalServerError(
			fmt.Errorf("unable to read supplier kyc record: %w", err),
		)
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
func (fr *Repository) FetchKYCProcessingRequests(
	ctx context.Context,
) ([]*domain.KYCRequest, error) {
	ctx, span := tracer.Start(ctx, "FetchKYCProcessingRequests")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "status",
		Value:          "",
		Operator:       "!=",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to fetch kyc request documents: %v", err),
		)
	}

	res := []*domain.KYCRequest{}

	for _, doc := range docs {
		req := &domain.KYCRequest{}
		err = doc.DataTo(req)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, exceptions.InternalServerError(
				fmt.Errorf("unable to read supplier: %w", err),
			)
		}
		res = append(res, req)
	}

	return res, nil
}

// FetchKYCProcessingRequestByID retrieves a specific kyc processing request
func (fr *Repository) FetchKYCProcessingRequestByID(
	ctx context.Context,
	id string,
) (*domain.KYCRequest, error) {
	ctx, span := tracer.Start(ctx, "FetchKYCProcessingRequestByID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "id",
		Value:          id,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(
			fmt.Errorf("unable to fetch kyc request documents: %v", err),
		)
	}

	req := &domain.KYCRequest{}
	err = docs[0].DataTo(req)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(fmt.Errorf("unable to read supplier: %w", err))
	}

	return req, nil
}

// UpdateKYCProcessingRequest update the supplier profile
func (fr *Repository) UpdateKYCProcessingRequest(
	ctx context.Context,
	kycRequest *domain.KYCRequest,
) error {
	ctx, span := tracer.Start(ctx, "UpdateKYCProcessingRequest")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetKCYProcessCollectionName(),
		FieldName:      "id",
		Value:          kycRequest.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse kyc processing request as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update kyc processing request profile: %v", err),
		)
	}
	return nil
}

// FetchAdminUsers fetches all admins
func (fr *Repository) FetchAdminUsers(ctx context.Context) ([]*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "FetchAdminUsers")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "permissions",
		Value:          profileutils.DefaultAdminPermissions,
		Operator:       "array-contains-any",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	var admins []*profileutils.UserProfile
	for _, doc := range docs {
		u := &profileutils.UserProfile{}
		err = doc.DataTo(u)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, exceptions.InternalServerError(
				fmt.Errorf("unable to read user profile: %w", err),
			)
		}
		admins = append(admins, u)
	}
	return admins, nil
}

// PurgeUserByPhoneNumber removes the record of a user given a phone number.
func (fr *Repository) PurgeUserByPhoneNumber(ctx context.Context, phone string) error {
	ctx, span := tracer.Start(ctx, "PurgeUserByPhoneNumber")
	defer span.End()

	profile, err := fr.GetUserProfileByPhoneNumber(ctx, phone, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}

	// delete pin of the user
	pin, err := fr.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		// Should not panic but allow for deletion of the profile
		log.Printf("failed to get a user pin %v", err)
	}
	// Remove user profile with or without PIN
	if pin != nil {
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
	}
	// delete user supplier profile
	// some old profiles may not have a supplier profile since the original implementation
	// created a supplier profile only for PRO.
	// However the current and correct logic creates a supplier profile regardless of flavour.
	// Hence, the deletion of supplier
	// profile should only occur if a supplier profile exists and not throw an error.
	supplier, err := fr.GetSupplierProfileByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
	} else {
		if supplier != nil {
			err = fr.RemoveKYCProcessingRequest(ctx, supplier.ID)
			if err != nil {
				utils.RecordSpanError(span, err)
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
	}

	// delete user customer profile
	// some old profiles may not have a customer profile since the original implementation
	// created a customer profile only for CONSUMER.
	// However the current and correct logic creates a customer profile regardless of flavour.
	// Hence, the deletion of customer
	// profile should only occur if a customer profile exists and not throw an error.
	customer, err := fr.GetCustomerProfileByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
	} else {
		if customer != nil {
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
	flavour feedlib.Flavour,
	profileID string,
) (*profileutils.Customer, *profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "GetCustomerOrSupplierProfileByProfileID")
	defer span.End()

	var customer *profileutils.Customer
	var supplier *profileutils.Supplier

	switch flavour {
	case feedlib.FlavourConsumer:
		{
			customerProfile, err := fr.GetCustomerProfileByProfileID(ctx, profileID)
			if err != nil {
				utils.RecordSpanError(span, err)
				return nil, nil, fmt.Errorf("failed to get customer profile")
			}
			customer = customerProfile
		}
	case feedlib.FlavourPro:
		{
			supplierProfile, err := fr.GetSupplierProfileByProfileID(ctx, profileID)
			if err != nil {
				utils.RecordSpanError(span, err)
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
) (*dto.CreatedUserResponse, error) {
	ctx, span := tracer.Start(ctx, "GetOrCreatePhoneNumberUser")
	defer span.End()

	user, err := fr.FirebaseClient.GetUserByPhoneNumber(
		ctx,
		phone,
	)
	if err == nil {
		return &dto.CreatedUserResponse{
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return &dto.CreatedUserResponse{
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
	profile *profileutils.UserProfile,
	newSecondaryPhoneNumbers []string,
) error {
	ctx, span := tracer.Start(ctx, "HardResetSecondaryPhoneNumbers")
	defer span.End()

	profile.SecondaryPhoneNumbers = newSecondaryPhoneNumbers

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile secondary phone numbers: %v", err),
		)
	}

	return nil
}

// HardResetSecondaryEmailAddress does a hard reset of user secondary email addresses. This should be called when retiring specific
// secondary email addresses and passing in the new secondary email address as an argument.
func (fr *Repository) HardResetSecondaryEmailAddress(
	ctx context.Context,
	profile *profileutils.UserProfile,
	newSecondaryEmails []string,
) error {
	ctx, span := tracer.Start(ctx, "HardResetSecondaryEmailAddress")
	defer span.End()

	profile.SecondaryEmailAddresses = newSecondaryEmails

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(
			fmt.Errorf("unable to update user profile secondary phone numbers: %v", err),
		)
	}

	return nil
}

// CheckIfExperimentParticipant check if a user has subscribed to be an experiment participant
func (fr *Repository) CheckIfExperimentParticipant(
	ctx context.Context,
	profileID string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckIfExperimentParticipant")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetExperimentParticipantCollectionName(),
		FieldName:      "id",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
	}

	if len(docs) == 0 {
		return false, nil
	}
	return true, nil
}

// AddUserAsExperimentParticipant adds the provided user profile as an experiment participant if does not already exist.
// this method is idempotent.
func (fr *Repository) AddUserAsExperimentParticipant(
	ctx context.Context,
	profile *profileutils.UserProfile,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "AddUserAsExperimentParticipant")
	defer span.End()

	exists, err := fr.CheckIfExperimentParticipant(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	if !exists {
		createCommand := &CreateCommand{
			CollectionName: fr.GetExperimentParticipantCollectionName(),
			Data:           profile,
		}
		_, err = fr.FirestoreClient.Create(ctx, createCommand)
		if err != nil {
			utils.RecordSpanError(span, err)
			return false, exceptions.InternalServerError(
				fmt.Errorf(
					"unable to add user profile of ID %v in experiment_participant: %v",
					profile.ID,
					err,
				),
			)
		}
		return true, nil
	}
	// the user already exists as an experiment participant
	return true, nil

}

// RemoveUserAsExperimentParticipant removes the provide user profile as an experiment participant. This methold does not check
// for existence before deletion since non-existence is relatively equivalent to a removal
func (fr *Repository) RemoveUserAsExperimentParticipant(
	ctx context.Context,
	profile *profileutils.UserProfile,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "RemoveUserAsExperimentParticipant")
	defer span.End()

	// fetch the document References
	query := &GetAllQuery{
		CollectionName: fr.GetExperimentParticipantCollectionName(),
		FieldName:      "id",
		Value:          profile.ID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf("unable to parse user profile as firebase snapshot: %v", err),
		)
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
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(
			fmt.Errorf(
				"unable to remove user profile of ID %v from experiment_participant: %v",
				profile.ID,
				err,
			),
		)
	}

	return true, nil
}

// UpdateAddresses persists a user's home or work address information to the database
func (fr *Repository) UpdateAddresses(
	ctx context.Context,
	id string,
	address profileutils.Address,
	addressType enumutils.AddressType,
) error {
	ctx, span := tracer.Start(ctx, "UpdateAddresses")
	defer span.End()

	profile, err := fr.GetUserProfileByID(ctx, id, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	switch addressType {
	case enumutils.AddressTypeHome:
		{
			profile.HomeAddress = &address
		}
	case enumutils.AddressTypeWork:
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
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	updateCommand := &UpdateCommand{
		CollectionName: fr.GetUserProfileCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           profile,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	return nil
}

// AddNHIFDetails persists a user's NHIF details
func (fr *Repository) AddNHIFDetails(
	ctx context.Context,
	input dto.NHIFDetailsInput,
	profileID string,
) (*domain.NHIFDetails, error) {
	ctx, span := tracer.Start(ctx, "AddNHIFDetails")
	defer span.End()

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
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	getNhifQuery := &GetSingleQuery{
		CollectionName: collectionName,
		Value:          docRef.ID,
	}
	dsnap, err := fr.FirestoreClient.Get(ctx, getNhifQuery)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	nhif := &domain.NHIFDetails{}
	err = dsnap.DataTo(nhif)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return nhif, nil
}

// GetNHIFDetailsByProfileID fetches a user's NHIF details given their profile ID
func (fr *Repository) GetNHIFDetailsByProfileID(
	ctx context.Context,
	profileID string,
) (*domain.NHIFDetails, error) {
	ctx, span := tracer.Start(ctx, "GetNHIFDetailsByProfileID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetNHIFDetailsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) > 1 && serverutils.IsDebug() {
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
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return nhif, nil
}

// GetUserCommunicationsSettings fetches the communication settings of a specific user.
func (fr *Repository) GetUserCommunicationsSettings(
	ctx context.Context,
	profileID string,
) (*profileutils.UserCommunicationsSetting, error) {
	ctx, span := tracer.Start(ctx, "GetUserCommunicationsSettings")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetCommunicationsSettingsCollectionName(),
		FieldName:      "profileID",
		Value:          profileID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) > 1 && serverutils.IsDebug() {
		log.Printf("> 1 communications settings with profile ID %s (count: %d)",
			profileID,
			len(docs),
		)
	}

	if len(docs) == 0 {
		return &profileutils.UserCommunicationsSetting{ProfileID: profileID}, nil
	}

	comms := &profileutils.UserCommunicationsSetting{}
	err = docs[0].DataTo(comms)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return comms, nil
}

// SetUserCommunicationsSettings sets communication settings for a specific user
func (fr *Repository) SetUserCommunicationsSettings(
	ctx context.Context,
	profileID string,
	allowWhatsApp *bool,
	allowTextSms *bool,
	allowPush *bool,
	allowEmail *bool,
) (*profileutils.UserCommunicationsSetting, error) {

	ctx, span := tracer.Start(ctx, "SetUserCommunicationsSettings")
	defer span.End()

	// get the previous communications_settings
	comms, err := fr.GetUserCommunicationsSettings(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	setCommsSettings := profileutils.UserCommunicationsSetting{
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	// fetch the now set communications_settings and return it
	return fr.GetUserCommunicationsSettings(ctx, profileID)
}

// UpdateCustomerProfile does a generic update of the customer profile
// to add the data received from the ERP.
func (fr *Repository) UpdateCustomerProfile(
	ctx context.Context,
	profileID string,
	cus profileutils.Customer,
) (*profileutils.Customer, error) {
	ctx, span := tracer.Start(ctx, "UpdateCustomerProfile")
	defer span.End()

	customer, err := fr.GetCustomerProfileByProfileID(ctx, profileID)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return customer, nil
}

// PersistIncomingSMSData persists SMS data
func (fr *Repository) PersistIncomingSMSData(
	ctx context.Context,
	input *dto.AfricasTalkingMessage,
) error {
	ctx, span := tracer.Start(ctx, "PersistIncomingSMSData")
	defer span.End()

	message := &dto.AfricasTalkingMessage{
		Date:   input.Date,
		From:   input.From,
		ID:     input.ID,
		LinkID: input.LinkID,
		Text:   input.Text,
		To:     input.To,
	}

	validatedMessage, err := utils.ValidateAficasTalkingSMSData(message)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetSMSCollectionName(),
		Data:           validatedMessage,
	}

	_, err = fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}

	return nil

}

// SaveUSSDEvent saves the USSD event that has taken place while interacting with the USSD
func (fr *Repository) SaveUSSDEvent(
	ctx context.Context,
	input *dto.USSDEvent,
) (*dto.USSDEvent, error) {
	ctx, span := tracer.Start(ctx, "SaveUSSDEvent")
	defer span.End()

	ussdEvent := &dto.USSDEvent{
		SessionID:         input.SessionID,
		PhoneNumber:       input.PhoneNumber,
		USSDEventDateTime: input.USSDEventDateTime,
		Level:             input.Level,
		USSDEventName:     input.USSDEventName,
	}

	validatesUSSDEvent, err := utils.ValidateUSSDEvent(ussdEvent)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetUSSDEventsCollectionName(),
		Data:           validatesUSSDEvent,
	}

	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	getUSSDEventsQuery := &GetSingleQuery{
		CollectionName: fr.GetUSSDEventsCollectionName(),
		Value:          docRef.ID,
	}

	docsnapshot, err := fr.FirestoreClient.Get(ctx, getUSSDEventsQuery)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	event := &dto.USSDEvent{}
	err = docsnapshot.DataTo(event)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return event, nil
}

// SaveCoverAutolinkingEvents saves cover linking events into the database
func (fr *Repository) SaveCoverAutolinkingEvents(
	ctx context.Context,
	input *dto.CoverLinkingEvent,
) (*dto.CoverLinkingEvent, error) {
	ctx, span := tracer.Start(ctx, "SaveCoverAutolinkingEvents")
	defer span.End()

	coverLinkingEvent := &dto.CoverLinkingEvent{
		ID:                    input.ID,
		CoverLinkingEventTime: input.CoverLinkingEventTime,
		CoverStatus:           input.CoverStatus,
		MemberNumber:          input.MemberNumber,
		PhoneNumber:           input.PhoneNumber,
	}

	validatedCoverLinkingEvent, err := utils.ValidateCoverLinkingEvent(coverLinkingEvent)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetCoverLinkingEventsCollectionName(),
		Data:           validatedCoverLinkingEvent,
	}

	docRef, err := fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	getCoverlinkingEventQuery := &GetSingleQuery{
		CollectionName: fr.GetCoverLinkingEventsCollectionName(),
		Value:          docRef.ID,
	}

	docsnapshot, err := fr.FirestoreClient.Get(ctx, getCoverlinkingEventQuery)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	event := &dto.CoverLinkingEvent{}
	err = docsnapshot.DataTo(event)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return event, nil
}

// AddAITSessionDetails saves diallers session details in the database
func (fr *Repository) AddAITSessionDetails(
	ctx context.Context,
	input *dto.SessionDetails,
) (*domain.USSDLeadDetails, error) {
	ctx, span := tracer.Start(ctx, "AddAITSessionDetails")
	defer span.End()

	validateDetails, err := utils.ValidateUSSDDetails(input)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	sessionDetails := &domain.USSDLeadDetails{
		ID:             uuid.New().String(),
		Level:          validateDetails.Level,
		PhoneNumber:    *validateDetails.PhoneNumber,
		SessionID:      validateDetails.SessionID,
		IsRegistered:   false,
		ContactChannel: "USSD",
		WantCover:      false,
	}

	createCommand := &CreateCommand{
		CollectionName: fr.GetUSSDDataCollectionName(),
		Data:           sessionDetails,
	}

	_, err = fr.FirestoreClient.Create(ctx, createCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	sessionDetails, err = fr.GetAITSessionDetails(ctx, validateDetails.SessionID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return sessionDetails, nil
}

// ListUserProfiles fetches all users with the specified role from the database
func (fr *Repository) ListUserProfiles(
	ctx context.Context,
	role profileutils.RoleType,
) ([]*profileutils.UserProfile, error) {
	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "role",
		Value:          role,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	profiles := []*profileutils.UserProfile{}

	for _, doc := range docs {
		profile := &profileutils.UserProfile{}
		err = doc.DataTo(profile)
		if err != nil {
			return nil, exceptions.InternalServerError(
				fmt.Errorf("unable to read agent user profile: %w", err),
			)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetAITSessionDetails gets Africa's Talking session details
func (fr *Repository) GetAITSessionDetails(
	ctx context.Context,
	sessionID string,
) (*domain.USSDLeadDetails, error) {
	ctx, span := tracer.Start(ctx, "GetAITSessionDetails")
	defer span.End()

	validatedSessionID, err := utils.CheckEmptyString(sessionID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	query := &GetAllQuery{
		CollectionName: fr.GetUSSDDataCollectionName(),
		FieldName:      "sessionID",
		Value:          validatedSessionID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	sessionDet := &domain.USSDLeadDetails{}
	err = docs[0].DataTo(sessionDet)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return sessionDet, nil

}

// UpdateSessionLevel updates user interaction level whike they interact with USSD
func (fr *Repository) UpdateSessionLevel(
	ctx context.Context,
	sessionID string,
	level int,
) (*domain.USSDLeadDetails, error) {
	ctx, span := tracer.Start(ctx, "UpdateSessionLevel")
	defer span.End()

	validSessionID, err := utils.CheckEmptyString(sessionID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	sessionDetails, err := fr.GetAITSessionDetails(ctx, *validSessionID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	collectionName := fr.GetUSSDDataCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "sessionID",
		Value:          sessionID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	sessionDetails.Level = level

	updateCommand := &UpdateCommand{
		CollectionName: collectionName,
		ID:             docs[0].Ref.ID,
		Data:           sessionDetails,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return sessionDetails, nil
}

// UpdateSessionPIN updates current user's session PIN when signing up or changing PIN
func (fr *Repository) UpdateSessionPIN(
	ctx context.Context,
	sessionID string,
	pin string,
) (*domain.USSDLeadDetails, error) {
	ctx, span := tracer.Start(ctx, "UpdateSessionPIN")
	defer span.End()

	sessionDetails, err := fr.GetAITSessionDetails(ctx, sessionID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	collectionName := fr.GetUSSDDataCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "sessionID",
		Value:          sessionID,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	sessionDetails.PIN = pin

	updateCommand := &UpdateCommand{
		CollectionName: collectionName,
		ID:             docs[0].Ref.ID,
		Data:           sessionDetails,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}
	return sessionDetails, nil
}

// GetAITDetails retrieves session details from the database
func (fr *Repository) GetAITDetails(
	ctx context.Context,
	phoneNumber string,
) (*domain.USSDLeadDetails, error) {
	ctx, span := tracer.Start(ctx, "GetAITDetails")
	defer span.End()

	validPhoneNumber, err := utils.CheckEmptyString(phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	query := &GetAllQuery{
		CollectionName: fr.GetUSSDDataCollectionName(),
		FieldName:      "phoneNumber",
		Value:          validPhoneNumber,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	ussdLead := &domain.USSDLeadDetails{}
	err = docs[0].DataTo(ussdLead)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return ussdLead, nil

}

// UpdateAITSessionDetails updates session details using phone number
func (fr *Repository) UpdateAITSessionDetails(
	ctx context.Context,
	phoneNumber string,
	contactLead *domain.USSDLeadDetails,
) error {
	ctx, span := tracer.Start(ctx, "UpdateAITSessionDetails")
	defer span.End()

	validPhoneNumber, err := utils.CheckEmptyString(phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	contactDetails, err := fr.GetAITDetails(ctx, *validPhoneNumber)
	if err != nil {
		return err
	}

	collectionName := fr.GetUSSDDataCollectionName()
	query := &GetAllQuery{
		CollectionName: collectionName,
		FieldName:      "phoneNumber",
		Value:          phoneNumber,
		Operator:       "==",
	}
	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	contactDetails.FirstName = contactLead.FirstName
	contactDetails.LastName = contactLead.LastName
	contactDetails.DateOfBirth = contactLead.DateOfBirth
	contactDetails.IsRegistered = contactLead.IsRegistered

	updateCommand := &UpdateCommand{
		CollectionName: collectionName,
		ID:             docs[0].Ref.ID,
		Data:           contactDetails,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	return nil
}

// CreateRole creates a new role and persists it to the database
func (fr *Repository) CreateRole(
	ctx context.Context,
	profileID string,
	input dto.RoleInput,
) (*profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "CreateRole")
	defer span.End()

	exists, err := fr.CheckIfRoleNameExists(ctx, input.Name)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	if exists {
		err := fmt.Errorf("role with similar name exists:%v", input.Name)
		utils.RecordSpanError(span, err)
		return nil, err
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)

	role := profileutils.Role{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   profileID,
		Created:     timestamp,
		Active:      true,
		Scopes:      input.Scopes,
	}

	createCommad := &CreateCommand{
		CollectionName: fr.GetRolesCollectionName(),
		Data:           role,
	}

	_, err = fr.FirestoreClient.Create(ctx, createCommad)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return &role, nil
}

// GetAllRoles returns a list of all created roles
func (fr *Repository) GetAllRoles(ctx context.Context) (*[]profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "GetAllRoles")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)

	if err != nil {
		utils.RecordSpanError(span, err)
		err = fmt.Errorf("unable to read role")
		return nil, exceptions.InternalServerError(err)
	}

	roles := []profileutils.Role{}
	for _, doc := range docs {
		role := &profileutils.Role{}

		err := doc.DataTo(role)
		if err != nil {
			utils.RecordSpanError(span, err)
			err = fmt.Errorf("unable to read role")
			return nil, exceptions.InternalServerError(err)
		}
		roles = append(roles, *role)
	}

	return &roles, nil
}

// UpdateRoleDetails  updates the details of a role
func (fr *Repository) UpdateRoleDetails(
	ctx context.Context,
	profileID string,
	role profileutils.Role,
) (*profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "UpdateRoleDetails")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
		Value:          role.ID,
		FieldName:      "id",
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)

	updatedRole := profileutils.Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Active:      role.Active,
		Scopes:      role.Scopes,
		CreatedBy:   role.CreatedBy,
		Created:     role.Created,
		UpdatedBy:   profileID,
		Updated:     timestamp,
	}

	updateCommand := &UpdateCommand{
		CollectionName: fr.GetRolesCollectionName(),
		ID:             docs[0].Ref.ID,
		Data:           updatedRole,
	}
	err = fr.FirestoreClient.Update(ctx, updateCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	return &updatedRole, nil
}

// GetRoleByID gets role with matching id
func (fr *Repository) GetRoleByID(ctx context.Context, roleID string) (*profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "GetRoleByID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
		FieldName:      "id",
		Value:          roleID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) != 1 {
		err = fmt.Errorf("role not found: %v", roleID)
		utils.RecordSpanError(span, err)
		return nil, err
	}

	role := &profileutils.Role{}

	err = docs[0].DataTo(role)
	if err != nil {
		utils.RecordSpanError(span, err)
		err = fmt.Errorf("unable to read role")
		return nil, exceptions.InternalServerError(err)
	}

	return role, nil
}

// GetRolesByIDs gets all roles matching provided roleIDs if specified otherwise all roles
func (fr *Repository) GetRolesByIDs(
	ctx context.Context,
	roleIDs []string,
) (*[]profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "GetRoleByID")
	defer span.End()
	roles := []profileutils.Role{}
	// role ids provided
	for _, id := range roleIDs {
		role, err := fr.GetRoleByID(ctx, id)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return &roles, nil
}

// DeleteRole removes a role permanently from the database
func (fr *Repository) DeleteRole(
	ctx context.Context,
	roleID string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "DeleteRole")
	defer span.End()

	// remove this role for all users who has it assigned
	query1 := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "roles",
		Operator:       "array-contains",
		Value:          roleID,
	}

	docs1, err := fr.FirestoreClient.GetAll(ctx, query1)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	for _, doc := range docs1 {
		user := &profileutils.UserProfile{}
		err = doc.DataTo(user)
		if err != nil {
			return false, fmt.Errorf("unable to parse userprofile")
		}
		newRoles := []string{}
		for _, userRole := range user.Roles {
			if userRole != roleID {
				newRoles = append(newRoles, userRole)
			}
		}
		err = fr.UpdateUserRoleIDs(ctx, user.ID, newRoles)
		if err != nil {
			utils.RecordSpanError(span, err)
			return false, exceptions.InternalServerError(err)
		}
	}

	// delete the role
	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
		FieldName:      "id",
		Value:          roleID,
		Operator:       "==",
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	// means the document was removed or does not exist
	if len(docs) == 0 {
		return false, fmt.Errorf("error role does not exist")
	}
	deleteCommand := &DeleteCommand{
		CollectionName: fr.GetRolesCollectionName(),
		ID:             docs[0].Ref.ID,
	}
	err = fr.FirestoreClient.Delete(ctx, deleteCommand)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf(
			"unable to remove role of ID %v, error: %v",
			roleID,
			err,
		)
	}
	return true, nil
}

// CheckIfRoleNameExists checks if a role with a similar name exists
// Ensures unique name for each role during creation
func (fr *Repository) CheckIfRoleNameExists(ctx context.Context, name string) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckIfRoleNameExists")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
		FieldName:      "name",
		Operator:       "==",
		Value:          name,
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	if len(docs) == 1 {
		return true, nil
	}

	return false, nil
}

// GetUserProfilesByRoleID returns a list of user profiles with the role ID
// i.e users assigned a particular role
func (fr *Repository) GetUserProfilesByRoleID(ctx context.Context, roleID string) ([]*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "GetUserProfilesByRoleID")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetUserProfileCollectionName(),
		FieldName:      "roles",
		Operator:       "array-contains",
		Value:          roleID,
	}

	users := []*profileutils.UserProfile{}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	for _, doc := range docs {
		user := &profileutils.UserProfile{}

		err = doc.DataTo(user)
		if err != nil {
			return nil, fmt.Errorf("unable to parse userprofile")
		}

		users = append(users, user)
	}

	return users, nil
}

// GetRoleByName retrieves a role using it's name
func (fr *Repository) GetRoleByName(ctx context.Context, roleName string) (*profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "GetRoleByName")
	defer span.End()

	query := &GetAllQuery{
		CollectionName: fr.GetRolesCollectionName(),
		FieldName:      "name",
		Operator:       "==",
		Value:          roleName,
	}

	docs, err := fr.FirestoreClient.GetAll(ctx, query)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	if len(docs) != 1 {
		err = fmt.Errorf("role with name %v not found", roleName)
		utils.RecordSpanError(span, err)
		return nil, err
	}

	role := &profileutils.Role{}

	err = docs[0].DataTo(role)
	if err != nil {
		utils.RecordSpanError(span, err)
		err = fmt.Errorf("unable to read role")
		return nil, exceptions.InternalServerError(err)
	}

	return role, nil
}

//CheckIfUserHasPermission checks if a user has the required permission
func (fr *Repository) CheckIfUserHasPermission(
	ctx context.Context,
	UID string,
	requiredPermission profileutils.Permission,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckIfUserHasPermission")
	defer span.End()

	userprofile, err := fr.GetUserProfileByUID(ctx, UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	roles, err := fr.GetRolesByIDs(ctx, userprofile.Roles)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	for _, role := range *roles {
		if role.Active && role.HasPermission(ctx, requiredPermission.Scope) {
			return true, nil
		}
	}

	return false, nil
}
