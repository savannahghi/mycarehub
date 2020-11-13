// Package profile maintains user (consumer and practitioner) profiles
package profile

import (
	"context"
	"fmt"
	"log"
	"math"

	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/authorization/graph/authorization"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/mailgun/graph/mailgun"
	"gitlab.slade360emr.com/go/otp/graph/otp"
)

// configuration constants
const (
	userProfileCollectionName           = "user_profiles"
	practitionerCollectionName          = "practitioners"
	surveyCollectionName                = "post_visit_survey"
	healthcashRootCollectionName        = "healthcash"
	healthcashDepositsCollectionName    = "healthcash_deposits"
	healthcashWithdrawalsCollectionName = "healthcash_withdrawals"
	healthcashCurrency                  = "KES"
	emailSignupSubject                  = "Thank you for signing up"
	emailWelcomeSubject                 = "Welcome to Slade 360 HealthCloud"
	emailRejectionSubject               = "Your Account was not Approved"
	appleTesterPractitionerLicense      = "A1B4C6"
	legalAge                            = 18
	PINCollectionName                   = "pins"
	signUpInfoCollectionName            = "sign_up_info"
)

// NewService returns a new authentication service
func NewService() *Service {
	fc := &base.FirebaseClient{}
	ctx := context.Background()

	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("can't initialize Firebase app when setting up profile service: %s", err)
	}

	auth, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}

	firestore, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("can't initialize Firestore client when setting up profile service: %s", err)
	}

	erpClient, err := base.NewERPClient()
	if err != nil {
		log.Panicf("unable to initialize ERP client for profile service: %s", err)
	}
	if !erpClient.IsInitialized() {
		log.Panicf("uninitialized ERP client")
	}

	return &Service{
		firestoreClient: firestore,
		firebaseAuth:    auth,
		emailService:    mailgun.NewService(),
		otpService:      otp.NewService(),
		client:          erpClient,
	}
}

// Service is an authentication service. It handles authentication related
// issues e.g user profiles
type Service struct {
	firestoreClient *firestore.Client
	firebaseAuth    *auth.Client
	emailService    *mailgun.Service
	otpService      *otp.Service
	client          *base.ServerClient
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("profile service does not have an initialized firestoreClient")
	}

	if s.firebaseAuth == nil {
		log.Panicf("profile service does not have an initialized firebaseAuth")
	}

	if s.emailService == nil {
		log.Panicf("profile service does not have an initialized emailService")
	}

	if s.client == nil {
		log.Panicf("profile service does not have an initialized ERP client")
	}
}

// GetUserProfileCollectionName ...
func (s Service) GetUserProfileCollectionName() string {
	suffixed := base.SuffixCollection(userProfileCollectionName)
	return suffixed
}

// GetPractitionerCollectionName ...
func (s Service) GetPractitionerCollectionName() string {
	// add env suffix
	suffixed := base.SuffixCollection(practitionerCollectionName)
	return suffixed
}

// GetSurveyCollectionName ..
func (s Service) GetSurveyCollectionName() string {
	// add env suffix
	suffixed := base.SuffixCollection(surveyCollectionName)
	return suffixed
}

// GetHealthcashRootCollectionName ..
func (s Service) GetHealthcashRootCollectionName() string {
	// add env suffix
	suffixed := base.SuffixCollection(healthcashRootCollectionName)
	return suffixed
}

// GetHealthcashDepositsCollectionName ..
func (s Service) GetHealthcashDepositsCollectionName() string {
	// add env suffix
	suffixed := base.SuffixCollection(healthcashDepositsCollectionName)
	return suffixed
}

// GetHealthcashWithdrawalsCollectionName ..
func (s Service) GetHealthcashWithdrawalsCollectionName() string {
	// add env suffix
	suffixed := base.SuffixCollection(healthcashWithdrawalsCollectionName)
	return suffixed
}

// GetPINCollectionName ..
func (s Service) GetPINCollectionName() string {
	suffixed := base.SuffixCollection(PINCollectionName)
	return suffixed
}

// SavePINToFirestore persists the supplied OTP
func (s Service) SavePINToFirestore(personalIDNumber PIN) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetPINCollectionName()).Add(ctx, personalIDNumber)
	return err
}

// GetSignUpInfoCollectionName ..
func (s Service) GetSignUpInfoCollectionName() string {
	suffixed := base.SuffixCollection(signUpInfoCollectionName)
	return suffixed
}

// SaveSignUpInfoToFirestore persists the supplied sign up info
func (s Service) SaveSignUpInfoToFirestore(info SignUpInfo) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetSignUpInfoCollectionName()).Add(ctx, info)
	return err
}

// RetrieveUserProfileFirebaseDocSnapshotByUID retrieves the user profile of a
// specified user
func (s Service) RetrieveUserProfileFirebaseDocSnapshotByUID(
	ctx context.Context,
	uid string,
) (*firestore.DocumentSnapshot, error) {

	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	// the ordering is necessary in order to provide a stable sort order
	query := collection.Where("uid", "==", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		log.Printf("user %s has > 1 profile (they have %d)", uid, len(docs))
	}
	if len(docs) == 0 {
		newProfile := &UserProfile{
			UID:           uid,
			IsApproved:    false,
			TermsAccepted: false,
			CanExperiment: false,
		}
		docID, err := base.SaveDataToFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), newProfile)
		if err != nil {
			return nil, fmt.Errorf("unable to create new user profile: %w", err)
		}
		dsnap, err := collection.Doc(docID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve newly created user profile: %w", err)
		}
		return dsnap, nil
	}
	dsnap := docs[0]
	return dsnap, nil
}

// RetrieveUserProfileFirebaseDocSnapshot retrievs a raw Firebase doc snapshot
// for the logged in user's user profile or creates one if it does not exist
func (s Service) RetrieveUserProfileFirebaseDocSnapshot(
	ctx context.Context) (*firestore.DocumentSnapshot, error) {
	uid, err := authorization.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return s.RetrieveUserProfileFirebaseDocSnapshotByUID(ctx, uid)
}

// RetrieveFireStoreSnapshotByUID retrieves a specified Firestore document snapshot by its UID
func (s Service) RetrieveFireStoreSnapshotByUID(
	ctx context.Context, uid string, collectionName string,
	field string) (*firestore.DocumentSnapshot, error) {
	collection := s.firestoreClient.Collection(collectionName)
	query := collection.Where(field, "==", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		log.Printf("more than one snapshot found (they have %d)", len(docs))
	}
	if len(docs) == 0 {
		return nil, nil
	}
	dsnap := docs[0]

	return dsnap, nil
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (s Service) UserProfile(ctx context.Context) (*UserProfile, error) {
	s.checkPreconditions()
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	userProfile := &UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// GetProfile returns the profile of the user with the supplied uid
func (s Service) GetProfile(ctx context.Context, uid string) (*UserProfile, error) {
	s.checkPreconditions()
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshotByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	userProfile := &UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// AcceptTermsAndConditions updates the profile of the logged in user to indicate that they
// have accepted the terms and conditions
func (s Service) AcceptTermsAndConditions(
	ctx context.Context, accept bool) (bool, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}
	userProfile.TermsAccepted = accept
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// UpdatePhoneNumber updates the profile with the supplied phone number, if it
// was not already there
func (s Service) UpdatePhoneNumber(
	ctx context.Context, phone string) (bool, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}
	validatedPhone, err := base.ValidateMSISDN(phone, "", true, s.firestoreClient)
	if err != nil {
		return false, err
	}
	if !base.StringSliceContains(userProfile.Msisdns, phone) {
		userProfile.Msisdns = append(userProfile.Msisdns, validatedPhone)
	}
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// UpdateUserProfile updates a practitioner's user profile with the supplied
// data
func (s Service) UpdateUserProfile(
	ctx context.Context, input UserProfileInput) (*UserProfile, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	userProfile.PhotoBase64 = input.PhotoBase64
	userProfile.PhotoContentType = input.PhotoContentType

	verifiedMSISDNS := userProfile.Msisdns
	for _, msisdnInp := range input.Msisdns {
		validPhone, err := base.ValidateMSISDN(
			msisdnInp.Phone,
			msisdnInp.Otp,
			false,
			s.firestoreClient,
		)
		if err != nil {
			return nil, fmt.Errorf("invalid phone/OTP: %s", err)
		}
		if !base.StringSliceContains(verifiedMSISDNS, validPhone) {
			verifiedMSISDNS = append(verifiedMSISDNS, validPhone)
		}
	}

	verifiedEmails := userProfile.Emails
	if input.Emails != nil {
		for _, email := range input.Emails {
			if base.StringSliceContains(verifiedEmails, email) {
				continue
			}
			if !govalidator.IsEmail(email) {
				return nil, fmt.Errorf("%s is not a valid email", email)
			}
			verifiedEmails = append(verifiedEmails, email)
		}
	}

	userProfile.Msisdns = verifiedMSISDNS
	userProfile.Emails = verifiedEmails
	if input.PushTokens != nil && len(input.PushTokens) > 0 {
		// facilitate updating of push tokens e.g retire older ones
		for _, token := range input.PushTokens {
			if token != nil {
				userProfile.PushTokens = append(userProfile.PushTokens, *token)
			}
		}
	}
	userProfile.Name = input.Name
	userProfile.Bio = input.Bio
	userProfile.CanExperiment = input.CanExperiment
	userProfile.AskAgainToSetIsTester = input.AskAgainToSetIsTester
	userProfile.AskAgainToSetCanExperiment = input.AskAgainToSetCanExperiment
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}
	return userProfile, nil
}

// ConfirmEmail updates the profile of the logged in user with an email address
func (s Service) ConfirmEmail(ctx context.Context, email string) (*UserProfile, error) {
	s.checkPreconditions()

	if !govalidator.IsEmail(email) {
		return nil, fmt.Errorf("%s is not a valid email", email)
	}

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	verifiedEmails := userProfile.Emails
	if !base.StringSliceContains(verifiedEmails, email) {
		verifiedEmails = append(verifiedEmails, email)
	}
	userProfile.Emails = verifiedEmails

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// PractitionerSignUp is used to receive/record a practitioner's sign-up details
func (s Service) PractitionerSignUp(
	ctx context.Context, input PractitionerSignupInput) (bool, error) {
	s.checkPreconditions()
	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}
	practitioner := Practitioner{
		Profile:   *profile,
		License:   input.License,
		Cadre:     input.Cadre,
		Specialty: input.Specialty,
	}
	_, err = base.SaveDataToFirestore(
		s.firestoreClient, s.GetPractitionerCollectionName(), practitioner)
	if err != nil {
		return false, fmt.Errorf("unable to save practitioner info: %w", err)
	}

	// is the license belongs to the once expected from apple tester, approve their
	//profile automatically
	if input.License == appleTesterPractitionerLicense {
		return s.ApprovePractitionerSignup(ctx)
	}

	for _, practitionerEmail := range profile.Emails {
		err = s.SendPractitionerSignUpEmail(ctx, practitionerEmail)
		if err != nil {
			return false, fmt.Errorf("unable to send signup email: %w", err)
		}
	}

	if err != nil {
		return false, fmt.Errorf("unable to save practitioner details: %w", err)
	}
	return true, nil
}

// SendPractitionerSignUpEmail will send a signup email to the practitioner
func (s Service) SendPractitionerSignUpEmail(ctx context.Context, emailaddress string) error {
	text := generatePractitionerSignupEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}

	_, _, err := s.emailService.SendEmail(emailSignupSubject, text, emailaddress)
	if err != nil {
		return nil
	}

	return nil
}

// UpdateBiodata updates the profile of the logged in user with the supplied
// bio-data
func (s Service) UpdateBiodata(
	ctx context.Context, input BiodataInput) (*UserProfile, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	dob := input.DateOfBirth
	gender := input.Gender

	userProfile.DateOfBirth = &dob
	userProfile.Gender = &gender
	userProfile.Name = input.Name
	userProfile.Bio = input.Bio

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// RegisterPushToken registers a user's device (push) token
func (s Service) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	s.checkPreconditions()
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't register push token: %w", err)
	}
	if base.StringSliceContains(userProfile.PushTokens, token) {
		// don't add a token that already exists
		return true, nil
	}
	userProfile.PushTokens = append(userProfile.PushTokens, token)
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// CompleteSignup completes the sign-up
func (s Service) CompleteSignup(ctx context.Context) (bool, error) {
	s.checkPreconditions()
	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}

	// do not re-process approved profiles
	if profile.IsApproved {
		return true, nil
	}

	// update the profile and mark it as approved
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}
	profile.IsApproved = true
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, profile)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil
}

//ApprovePractitionerSignup is used to approve the practitioner signup
func (s Service) ApprovePractitionerSignup(ctx context.Context) (bool, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve user profile: %w", err)
	}

	if profile.IsApproved {
		return true, nil
	}

	if !profile.IsApproved {
		profile.IsApproved = true

		for _, practitionerEmail := range profile.Emails {
			err = s.SendPractitionerWelcomeEmail(ctx, practitionerEmail)
			if err != nil {
				return false, fmt.Errorf("unable to send welcome email: %w", err)
			}
		}
	}
	return true, nil
}

//RejectPractitionerSignup is used to reject the practitioner signup
func (s Service) RejectPractitionerSignup(ctx context.Context) (bool, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve user profile: %w", err)
	}

	if !profile.IsApproved {
		return false, nil
	}

	profile.IsApproved = false
	for _, practitionerEmail := range profile.Emails {
		err = s.SendPractitionerRejectionEmail(ctx, practitionerEmail)
		if err != nil {
			return false, fmt.Errorf("unable to send rejection email: %w", err)
		}
	}
	return false, nil
}

// SendPractitionerWelcomeEmail will send a welcome email to the practitioner
func (s Service) SendPractitionerWelcomeEmail(ctx context.Context, emailaddress string) error {
	s.checkPreconditions()

	text := generatePractitionerWelcomeEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	_, _, err := s.emailService.SendEmail(emailWelcomeSubject, text, emailaddress)
	if err != nil {
		return fmt.Errorf("unable to send welcome email: %w", err)
	}

	return nil
}

//SendPractitionerRejectionEmail will send a rejection email to the practitioner
func (s Service) SendPractitionerRejectionEmail(ctx context.Context, emailaddress string) error {
	s.checkPreconditions()
	text := generatePractitionerRejectionEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	_, _, err := s.emailService.SendEmail(emailRejectionSubject, text, emailaddress)
	if err != nil {
		return fmt.Errorf("unable to send rejection email: %w", err)
	}
	return nil
}

// RecordPostVisitSurvey records the survey input supplied by the user
func (s Service) RecordPostVisitSurvey(ctx context.Context, input PostVisitSurveyInput) (bool, error) {
	s.checkPreconditions()

	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return false, fmt.Errorf("the likelihood of recommending should be an int between 0 and 10")
	}

	uid, err := authorization.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, err
	}
	feedbackCollection := s.firestoreClient.Collection(s.GetSurveyCollectionName())
	feedback := PostVisitSurvey{
		LikelyToRecommend: input.LikelyToRecommend,
		Criticism:         input.Criticism,
		Suggestions:       input.Suggestions,
		UID:               uid,
		Timestamp:         time.Now(),
	}
	_, _, err = feedbackCollection.Add(ctx, feedback)
	if err != nil {
		return false, fmt.Errorf("unable to save feedback: %w", err)
	}
	return true, nil
}

// AddTester enrolls a user's email into the test group
func (s Service) AddTester(ctx context.Context, email string) (bool, error) {
	s.checkPreconditions()

	if !govalidator.IsEmail(email) {
		return false, fmt.Errorf("%s is not a valid email", email)
	}
	emails := []string{email}
	if isTester(ctx, emails) {
		return true, nil // add only once
	}
	tester := &TesterWhitelist{Email: email}
	_, _, err := base.CreateNode(ctx, tester)
	if err != nil {
		return false, fmt.Errorf("can't save whitelist entry: %s", err)
	}

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}
	userProfile.IsTester = true
	// reset covers
	userProfile.Covers = []Cover{}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil
}

// RemoveTester removes a user's email from the test group
func (s Service) RemoveTester(ctx context.Context, email string) (bool, error) {
	s.checkPreconditions()

	tester, err := getTester(ctx, email)
	if err != nil {
		return false, fmt.Errorf("can't get tester with email %s: %w", email, err)
	}
	if tester == nil {
		return true, nil // idempotent...you can safely "re-delete"
	}

	_, err = base.DeleteNode(ctx, tester.ID, &TesterWhitelist{})
	if err != nil {
		return false, fmt.Errorf("can't delete tester with email %s: %w", email, err)
	}

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}
	userProfile.IsTester = false
	// reset covers
	userProfile.Covers = []Cover{}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil
}

// ListTesters returns the emails of new testers
func (s Service) ListTesters(ctx context.Context) ([]string, error) {
	s.checkPreconditions()

	testerDocs, _, err := base.QueryNodes(
		ctx,
		nil,
		nil,
		nil,
		&TesterWhitelist{},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve testers: %w", err)
	}

	testers := []*TesterWhitelist{}
	for _, doc := range testerDocs {
		tester := &TesterWhitelist{}
		err := doc.DataTo(tester)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal tester doc to struct: %w", err)
		}
		testers = append(testers, tester)
	}

	emails := []string{}
	for _, tester := range testers {
		emails = append(emails, tester.Email)
	}

	return emails, nil
}

// ListKMPDURegisteredPractitioners lists the practitioners registered with KMPDU
func (s Service) ListKMPDURegisteredPractitioners(ctx context.Context, pagination *base.PaginationInput, filter *base.FilterInput, sort *base.SortInput) (*KMPDUPractitionerConnection, error) {
	s.checkPreconditions()

	registeredPractitioners, pageInfo, err := base.QueryNodes(
		ctx, pagination, filter, sort, &KMPDUPractitioner{},
	)
	if err != nil {
		return nil, err
	}
	edges := []*KMPDUPractitionerEdge{}
	for _, doc := range registeredPractitioners {
		practitioner := &KMPDUPractitioner{}
		err := doc.DataTo(practitioner)
		if err != nil {
			return nil, err
		}
		practitioner.ID = doc.Ref.ID
		edges = append(edges, &KMPDUPractitionerEdge{
			Cursor: &practitioner.ID,
			Node:   practitioner,
		})
	}

	connection := &KMPDUPractitionerConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
	return connection, nil
}

// GetRegisteredPractitionerByLicense retrieves a single practitioners records
func (s Service) GetRegisteredPractitionerByLicense(
	ctx context.Context, license string,
) (*KMPDUPractitioner, error) {
	s.checkPreconditions()
	dsnap, err := s.firestoreClient.Collection("RegisteredPractitioners").Doc(license).Get(ctx)
	if err != nil {
		return nil, err
	}
	practitioner := &KMPDUPractitioner{}
	err = dsnap.DataTo(practitioner)
	if err != nil {
		return nil, fmt.Errorf("unable to read practitioner records: %w", err)
	}

	return practitioner, nil
}

// TODO Separate practitioner and consumer profiles - isApproved
// TODO practitionerTermsOfServiceAccepted

// IsUnderAge checks if the user in context is an underage or not
func (s Service) IsUnderAge(ctx context.Context) (bool, error) {
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't retrieve user profile when getting the age: %w", err)
	}
	dob := userProfile.DateOfBirth
	if dob == nil {
		return false, fmt.Errorf("user should have a date of birth")
	}
	dateOfBirth := dob.AsTime()
	today := time.Now()
	age := math.Floor(today.Sub(dateOfBirth).Hours() / 24 / 365)

	if age >= legalAge {
		return false, nil
	}

	return true, nil
}

//SetUserPin receives phone number and pin from phonenumber sign up
// and save them to Firestore
func (s Service) SetUserPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	s.checkPreconditions()

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get or create a user profile: %v", err)
	}

	personalIDNumber := PIN{
		UID:     profile.UID,
		MSISDN:  phoneNumber,
		PIN:     pin,
		IsValid: true,
	}

	err = s.SavePINToFirestore(personalIDNumber)
	if err != nil {
		return false, fmt.Errorf("unable to save PIN: %v", err)
	}

	profile.HasPin = true
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve firebase user profile: %v", err)
	}
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, profile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil
}

// VerifyMSISDNandPin verifies a given msisdn and pin match.
func (s Service) VerifyMSISDNandPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	s.checkPreconditions()
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}
	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve pin given the msisdn: %v", err)
	}

	msisdnPin := &PIN{}
	err = dsnap.DataTo(msisdnPin)
	if err != nil {
		return false, fmt.Errorf("unable to read PIN: %w", err)
	}

	if msisdnPin.PIN != pin {
		return false, nil
	}

	return true, nil
}

// RetrievePINFirebaseDocSnapshotByMSISDN retrieves the user profile of a
// specified user
func (s Service) RetrievePINFirebaseDocSnapshotByMSISDN(
	ctx context.Context,
	msisdn string,
) (*firestore.DocumentSnapshot, error) {

	collection := s.firestoreClient.Collection(s.GetPINCollectionName())
	query := collection.Where("msisdn", "==", msisdn).Where("isValid", "==", true)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		log.Printf("msisdn %s has more than one pin (it has %d)", msisdn, len(docs))
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("pin can't be retrieved beacuse it does not exist")
	}
	dsnap := docs[0]
	return dsnap, nil
}

// CheckUserWithMsisdn checks if a user msisdn is present in pins firestore collection
// which essentially means that the number was used during user registration
func (s Service) CheckUserWithMsisdn(ctx context.Context, msisdn string) (bool, error) {
	s.checkPreconditions()
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}
	_, checkErr := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if checkErr != nil {
		return false, fmt.Errorf("user does not exist: %v", checkErr)
	}
	return true, nil
}

// RequestPinReset sends an otp to an exisitng user that is used to update their pin
func (s Service) RequestPinReset(ctx context.Context, msisdn string) (string, error) {
	s.checkPreconditions()

	_, err := s.CheckUserWithMsisdn(ctx, msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to get user profile: %v", err)
	}

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	otp, err := s.otpService.GenerateAndSendOTP(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("unable to generate and send otp: %v", err)
	}

	return otp, nil
}

// UpdateUserPin resets a user's pin
func (s Service) UpdateUserPin(ctx context.Context, msisdn string, pin string, otp string) (bool, error) {
	s.checkPreconditions()

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	_, checkErr := s.CheckUserWithMsisdn(ctx, phoneNumber)
	if checkErr != nil {
		return false, fmt.Errorf("unable to get user profile: %v", checkErr)
	}

	// verify OTP
	_, validateErr := base.ValidateMSISDN(phoneNumber, otp, false, s.firestoreClient)
	if validateErr != nil {
		return false, fmt.Errorf("OTP failed verification: %w", validateErr)
	}

	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return false, err
	}

	msisdnPin := &PIN{}
	err = dsnap.DataTo(msisdnPin)
	if err != nil {
		return false, fmt.Errorf("unable to read PIN: %w", err)
	}

	msisdnPin.PIN = pin

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetPINCollectionName(), dsnap.Ref.ID, msisdnPin,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// SetLanguagePreference sets the language a user prefers for using/interacting in be.well
func (s Service) SetLanguagePreference(ctx context.Context, language base.Language) (bool, error) {
	s.checkPreconditions()

	validLanguage := language.IsValid()
	if !validLanguage {
		return false, fmt.Errorf("%v is not an allowed language choice", language.String())
	}

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, err
	}

	userProfile.Language = language

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// CheckEmailVerified checks if the logged in user's email is verified
func (s Service) CheckEmailVerified(ctx context.Context) (bool, error) {
	s.checkPreconditions()

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't fetch user profile: %v", err)
	}

	emailVerified := userProfile.IsEmailVerified

	if !emailVerified {
		return false, nil
	}

	return true, nil
}

// CheckPhoneNumberVerified checks if the logged in user's phone number is verified
func (s Service) CheckPhoneNumberVerified(ctx context.Context) (bool, error) {
	s.checkPreconditions()

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't fetch user profile: %v", err)
	}

	msisdnVerified := userProfile.IsMsisdnVerified

	if !msisdnVerified {
		return false, nil
	}

	return true, nil
}

// VerifyEmailOtp checks for the validity of the supplied OTP but does not invalidate it
func (s Service) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}

	_, emailErr := ValidateEmail(email, otp, s.firestoreClient)
	if emailErr != nil {
		return false, fmt.Errorf("email failed verification: %w", err)
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't fetch user profile: %v", err)
	}

	userProfile.IsEmailVerified = true

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}

	return true, nil
}

// CreateSignUpMethod attahces a users sign up method to a user's UID
func (s Service) CreateSignUpMethod(ctx context.Context, signUpMethod SignUpMethod) (bool, error) {
	s.checkPreconditions()
	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get or create a user profile: %v", err)
	}

	validSignUpMethod := signUpMethod.IsValid()
	if !validSignUpMethod {
		return false, fmt.Errorf("%v is not an allowed sign up method choice", signUpMethod.String())
	}

	signUpInfo := SignUpInfo{
		UID:          profile.UID,
		SignUpMethod: signUpMethod,
	}

	err = s.SaveSignUpInfoToFirestore(signUpInfo)
	if err != nil {
		return false, fmt.Errorf("unable to save user sign up info: %v", err)
	}

	return true, nil
}

// GetSignUpMethod returns a user's sign up method
func (s Service) GetSignUpMethod(ctx context.Context, id string) (SignUpMethod, error) {
	dsnap, err := s.RetrieveFireStoreSnapshotByUID(ctx, id, s.GetSignUpInfoCollectionName(), "uid")
	if err != nil {
		return "", fmt.Errorf("unable to fetch sign up info: %v", err)
	}

	if dsnap == nil {
		return "", nil
	}

	info := &SignUpInfo{}
	err = dsnap.DataTo(info)
	if err != nil {
		return "", fmt.Errorf("unable to read sign up info: %w", err)
	}

	signUpMethod := info.SignUpMethod

	return signUpMethod, nil
}

// AddPractitionerServices persists a practitioner services to firestore
func (s Service) AddPractitionerServices(
	ctx context.Context, services PractitionerServiceInput,
	otherServices *OtherPractitionerServiceInput) (bool, error) {
	s.checkPreconditions()

	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to fetch user profile: %v", err)
	}
	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, profile.UID, s.GetPractitionerCollectionName(), "profile.uid")
	if err != nil {
		return false, fmt.Errorf("unable to retreive practitioner: %v", err)
	}

	if dsnap == nil {
		return false, nil
	}
	practitioner := &Practitioner{}
	err = dsnap.DataTo(practitioner)
	if err != nil {
		return false, fmt.Errorf("unable to read practitioner information: %v", err)
	}
	offeredServices := practitioner.Services.Services
	otherOfferedServices := practitioner.Services.OtherServices

	for _, service := range services.Services {
		validservice := service.IsValid()
		if !validservice {
			return false, fmt.Errorf("%v is not an allowed service enum", service.String())
		}

		if service == "OTHER" {
			if otherServices == nil {
				return false, fmt.Errorf("specify other services after selecting Others")
			}
			//TODO Pop the "others"
			offeredServices = append(offeredServices, service)
			otherOfferedServices = append(otherOfferedServices, otherServices.OtherServices...)

			practitioner.Services.Services = offeredServices
			practitioner.Services.OtherServices = otherOfferedServices

			practitioner.Profile.PractitionerHasServices = true

			err = base.UpdateRecordOnFirestore(
				s.firestoreClient, s.GetPractitionerCollectionName(), dsnap.Ref.ID, practitioner,
			)
			if err != nil {
				return false, fmt.Errorf("unable to update practitioner info: %v", err)
			}

			return true, nil
		}
		offeredServices = append(offeredServices, service)
		practitioner.Services.Services = offeredServices

		practitioner.Profile.PractitionerHasServices = true
	}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetPractitionerCollectionName(), dsnap.Ref.ID, practitioner,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update practitioner info: %v", err)
	}

	return true, nil
}
