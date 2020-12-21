// Package profile maintains user (consumer and practitioner) profiles
package profile

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gopkg.in/yaml.v2"
)

// configuration constants
const (
	userProfileCollectionName      = "user_profiles"
	practitionerCollectionName     = "practitioners"
	surveyCollectionName           = "post_visit_survey"
	emailSignupSubject             = "Thank you for signing up"
	emailWelcomeSubject            = "Welcome to Slade 360 HealthCloud"
	emailRejectionSubject          = "Your Account was not Approved"
	appleTesterPractitionerLicense = "A1B4C6"
	legalAge                       = 18
	PINCollectionName              = "pins"
	ProfileNudgesCollectionName    = "profile_nudges"
	KCYProcessCollectionName       = "kyc_processing"
	signUpInfoCollectionName       = "sign_up_info"
)

const (
	mailgunService    = "mailgun"
	otpService        = "otp"
	engagementService = "engagement"
	smsService        = "sms"
)

// OTP service endpoints
const (
	SendEmail    = "internal/send_email"
	SendRetryOtp = "internal/send_retry_otp/"
	SendOtp      = "internal/send_otp/"
)

// NewService returns a new authentication service
func NewService() *Service {

	var config base.DepsConfig

	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(base.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occured while opening deps file %v", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

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

	chargemasterClient, err := NewChargeMasterClient()
	if err != nil {
		log.Panicf("unable to initialize chargemaster client for profile service: %s", err)
	}
	if !chargemasterClient.IsInitialized() {
		log.Panicf("uninitialized chargemaster client")
	}

	var mailgunClient *base.InterServiceClient
	mailgunClient, err = base.SetupISCclient(config, mailgunService)
	if err != nil {
		log.Panicf("unable to initialize mailgun inter service client: %s", err)
	}

	var otpClient *base.InterServiceClient
	otpClient, err = base.SetupISCclient(config, otpService)
	if err != nil {
		log.Panicf("unable to initialize otp inter service client: %s", err)

	}

	var engagementClient *base.InterServiceClient
	engagementClient, err = base.SetupISCclient(config, engagementService)
	if err != nil {
		log.Panicf("unable to initialize engagement inter service client: %s", err)

	}

	var smsClient *base.InterServiceClient
	smsClient, err = base.SetupISCclient(config, smsService)
	if err != nil {
		log.Panicf("unable to initialize sms inter service client: %s", err)

	}

	return &Service{
		firestoreClient:    firestore,
		firebaseAuth:       auth,
		erpClient:          erpClient,
		Mailgun:            mailgunClient,
		Otp:                otpClient,
		chargemasterClient: chargemasterClient,
		engagement:         engagementClient,
		sms:                smsClient,
	}
}

// Service is an authentication service. It handles authentication related
// issues e.g user profiles
type Service struct {
	Mailgun    *base.InterServiceClient
	Otp        *base.InterServiceClient
	engagement *base.InterServiceClient
	sms        *base.InterServiceClient

	firestoreClient    *firestore.Client
	firebaseAuth       *auth.Client
	erpClient          *base.ServerClient
	chargemasterClient *base.ServerClient
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("profile service does not have an initialized firestoreClient")
	}

	if s.firebaseAuth == nil {
		log.Panicf("profile service does not have an initialized firebaseAuth")
	}

	if s.Mailgun == nil {
		log.Panicf("profile service does not have an initialized mailgun ISC Client")
	}

	if s.erpClient == nil {
		log.Panicf("profile service does not have an initialized ERP client")
	}

	if s.Otp == nil {
		log.Panicf("profile service does not have an initialized otp ISC Client")
	}

	if s.engagement == nil {
		log.Panicf("profile service does not have an initialized engagement ISC Client")
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

// GetProfileNudgesCollectionName return the storage location of profile nudges
func (s Service) GetProfileNudgesCollectionName() string {
	suffixed := base.SuffixCollection(ProfileNudgesCollectionName)
	return suffixed
}

// GetKCYProcessCollectionName fetches location where kyc processing request will be saved
func (s Service) GetKCYProcessCollectionName() string {
	suffixed := base.SuffixCollection(KCYProcessCollectionName)
	return suffixed
}

// SaveSignUpInfoToFirestore persists the supplied sign up info
func (s Service) SaveSignUpInfoToFirestore(info SignUpInfo) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetSignUpInfoCollectionName()).Add(ctx, info)
	return err
}

// SaveProfileNudge stages nudges published from this service. These nudges will be
// referenced later to support some specialized use-case. A nudge will be uniquely
// identified by its id and sequenceNumber
func (s Service) SaveProfileNudge(nudge map[string]interface{}) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetProfileNudgesCollectionName()).Add(ctx, nudge)
	return err
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (s Service) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	s.checkPreconditions()
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetUserProfile(ctx, uid)
}

// GetOrCreateUserProfile retrieves the user profile of a
// specified user using either their uid or phone number.
// If the user perofile does not exist then a new one is created
func (s Service) GetOrCreateUserProfile(ctx context.Context, phone string) (*base.UserProfile, error) {
	s.checkPreconditions()

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return nil, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	dsnap, err := s.RetrieveOrCreateUserProfileFirebaseDocSnapshot(ctx, uid, phoneNumber)
	if err != nil {
		return nil, err
	}
	userProfile := &base.UserProfile{}

	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	if !base.StringSliceContains(userProfile.VerifiedIdentifiers, uid) {
		userProfile.VerifiedIdentifiers = append(userProfile.VerifiedIdentifiers, uid)
		err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to update user profile: %v", err)
		}
	}

	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// GetProfile returns the profile of the user with the supplied uid
func (s Service) GetProfile(ctx context.Context, uid string) (*base.UserProfile, error) {
	s.checkPreconditions()
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshotByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// GetProfileByID returns the profile identified by the indicated ID
func (s Service) GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshotByID(ctx, id)
	if err != nil {
		return nil, err
	}
	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
	return userProfile, nil
}

// FindProfile returns a user profile if it exists and returns a nil if the
// profile does not exist instead of creating a new default profile
// This purely handles the issue of backwards compatibility and should be depreciated
// once the side effects are handled.
func (s Service) FindProfile(ctx context.Context) (*base.UserProfile, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetUserProfileCollectionName(), "verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to get a profile dsnap for this user: %v", err)
	}

	if dsnap == nil {
		return nil, fmt.Errorf("the user's profile has not been found")
	}

	userProfile := &base.UserProfile{}
	err = dsnap.DataTo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

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
	valid, err := base.VerifyOTP(phone, "", s.Otp)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, fmt.Errorf("phonenumber given is not valid")
	}

	if !base.StringSliceContains(userProfile.Msisdns, phone) {
		userProfile.Msisdns = append(userProfile.Msisdns, phone)
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
	ctx context.Context, input UserProfileInput) (*base.UserProfile, error) {
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

	msisdns := userProfile.Msisdns
	verifiedPhones := userProfile.VerifiedPhones
	if input.Msisdns != nil {
		for _, msisdnInp := range input.Msisdns {
			valid, err := base.VerifyOTP(
				msisdnInp.Phone,
				msisdnInp.Otp,
				s.Otp,
			)
			if err != nil {
				return nil, fmt.Errorf("invalid phone/OTP: %s", err)
			}
			if !valid {
				return nil, fmt.Errorf("phonenumber given is not valid")
			}
			if !base.StringSliceContains(msisdns, msisdnInp.Phone) {
				msisdns = append(msisdns, msisdnInp.Phone)
			}
			verifyPhone := base.VerifiedMsisdn{
				Msisdn:   msisdnInp.Phone,
				Verified: true,
			}
			verifiedPhones = append(verifiedPhones, verifyPhone)
		}
	}

	emails := userProfile.Emails
	verifiedEmails := userProfile.VerifiedEmails
	if input.Emails != nil {
		for _, email := range input.Emails {
			if base.StringSliceContains(emails, email) {
				continue
			}
			if !govalidator.IsEmail(email) {
				return nil, fmt.Errorf("%s is not a valid email", email)
			}
			emails = append(emails, email)
			verifyEmail := base.VerifiedEmail{
				Email:    email,
				Verified: true,
			}
			verifiedEmails = append(verifiedEmails, verifyEmail)
		}
	}

	userProfile.Msisdns = msisdns
	userProfile.Emails = emails
	userProfile.VerifiedPhones = verifiedPhones
	userProfile.VerifiedEmails = verifiedEmails

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
func (s Service) ConfirmEmail(ctx context.Context, email string) (*base.UserProfile, error) {
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

	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": emailSignupSubject,
	}

	resp, err := s.Mailgun.MakeRequest(http.MethodPost, SendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send Practitioner signup email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send Practitioner signup email : %w, with status code %v", err, resp.StatusCode)
	}

	return nil
}

// UpdateBiodata updates the profile of the logged in user with the supplied
// bio-data
func (s Service) UpdateBiodata(
	ctx context.Context, input BiodataInput) (*base.UserProfile, error) {
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

	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": emailWelcomeSubject,
	}

	resp, err := s.Mailgun.MakeRequest(http.MethodPost, SendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send welcome email: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send welcome email : %w, with status code %v", err, resp.StatusCode)
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

	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": emailRejectionSubject,
	}

	resp, err := s.Mailgun.MakeRequest(http.MethodPost, SendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send rejection email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send rejection email : %w, with status code %v", err, resp.StatusCode)
	}

	return nil
}

// RecordPostVisitSurvey records the survey input supplied by the user
func (s Service) RecordPostVisitSurvey(ctx context.Context, input PostVisitSurveyInput) (bool, error) {
	s.checkPreconditions()

	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return false, fmt.Errorf("the likelihood of recommending should be an int between 0 and 10")
	}

	uid, err := base.GetLoggedInUserUID(ctx)
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
	userProfile.Covers = []base.Cover{}

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
	userProfile.Covers = []base.Cover{}

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

// IsUnderAge checks if the user in context is an underage or not
func (s Service) IsUnderAge(ctx context.Context) (bool, error) {
	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't retrieve user profile when getting the age: %w", err)
	}
	dob := userProfile.DateOfBirth
	if dob == nil {
		log.Printf("user %s does not have a date of birth", userProfile.ID)
		return true, nil
	}
	dateOfBirth := dob.AsTime()
	today := time.Now()
	age := math.Floor(today.Sub(dateOfBirth).Hours() / 24 / 365)
	if age >= legalAge {
		return false, nil
	}

	return true, nil
}

//SetUserPIN receives phone number and pin from phonenumber sign up
// and save them to Firestore
func (s Service) SetUserPIN(ctx context.Context, msisdn string, pin string) (bool, error) {
	s.checkPreconditions()
	// retrieve profile linked to this user
	profile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get a user profile: %v", err)
	}

	err = validatePIN(pin)
	if err != nil {
		return false, fmt.Errorf("invalid pin: %w", err)
	}

	// ensure the phone number is valid
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}
	// check if user has existing PIN
	exists, err := s.CheckHasPIN(ctx, msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to check if the user has a PIN: %v", err)
	}
	// return true if the user already have one
	if exists {
		if base.IsDebug() {
			log.Printf("user with msisdn %s has more than one PINs)", msisdn)
		}
		return true, nil
	}
	// EncryptPIN the PIN
	encryptedPin, err := EncryptPIN(pin)
	if err != nil {
		return false, fmt.Errorf("unable to encrypt PIN: %w", err)
	}

	// we link the PIN to their profile
	// one profile should have one PIN
	PINPayload := PIN{
		ProfileID: profile.ID,
		MSISDN:    phoneNumber,
		PINNumber: encryptedPin,
		IsValid:   true,
	}

	err = s.SavePINToFirestore(PINPayload)
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

// VerifyMSISDNandPIN verifies a given msisdn and pin match.
func (s Service) VerifyMSISDNandPIN(ctx context.Context, msisdn string, pinNumber string) (bool, error) {
	s.checkPreconditions()
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}
	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve pin given the msisdn: %v", err)
	}
	if dsnap == nil {
		return false, fmt.Errorf("VerifyMSISDNandPIN: unable to retrieve user PIN")
	}

	msisdnPin := &PIN{}
	err = dsnap.DataTo(msisdnPin)
	if err != nil {
		return false, fmt.Errorf("unable to read PIN: %w", err)
	}

	err = s.encryptExistingPin(msisdnPin, dsnap)
	if err != nil {
		return false, err
	}

	// compare if the two PINS match
	isMatch, err := ComparePIN(msisdnPin.PINNumber, pinNumber)
	if err != nil {
		return false, fmt.Errorf("unable to match PIN Number provided: %w", err)
	}
	if !isMatch {
		return false, nil
	}

	return true, nil
}

// CheckHasPIN given a phone number checks if the phonenumber is present in our collections
// which essentially means that the number has an already existing PIN
func (s Service) CheckHasPIN(ctx context.Context, msisdn string) (bool, error) {
	s.checkPreconditions()
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("unable to fetch PINs dsnap: %v", err)
	}
	if dsnap == nil {
		return false, nil
	}

	return true, nil
}

// SendRetryOTP generates fallback OTPs when Africa is talking sms fails
func (s Service) SendRetryOTP(ctx context.Context, msisdn string, retryStep int) (string, error) {
	s.checkPreconditions()

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
	}
	resp, err := s.Otp.MakeRequest(http.MethodPost, SendRetryOtp, body)
	if err != nil {
		return "", fmt.Errorf("unable to generate and send fallback otp: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to generate and send fallback otp, with status code %v", resp.StatusCode)
	}

	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to convert response to string: %v", err)
	}

	return string(code), nil
}

// RequestPINReset given an existing user's phone number, sends an otp to the phone number
// that is then used in the process of updating their old PIN to a new one
func (s Service) RequestPINReset(ctx context.Context, msisdn string) (string, error) {
	s.checkPreconditions()
	// if they are a firebase user
	exists, err := s.CheckHasPIN(ctx, msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to check if the user has a PIN: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("request for a PIN reset failed. User does not have an existing PIN")
	}

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	code, err := s.generateAndSendOTP(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("can't generate and send an otp to %s: %v", phoneNumber, err)
	}

	return code, nil
}

// UpdateUserPIN resets a user's pin
func (s Service) UpdateUserPIN(ctx context.Context, msisdn string, pin string) (bool, error) {
	s.checkPreconditions()
	exists, err := s.CheckHasPIN(ctx, msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to check if the user has a PIN: %v", err)
	}
	if !exists {
		return false, fmt.Errorf("request for a PIN update failed. User does not have an existing PIN")
	}

	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return false, err
	}
	if dsnap == nil {
		return false, fmt.Errorf("UpdateUserPIN: unable to retrieve user PIN")
	}
	msisdnPIN := &PIN{}
	err = dsnap.DataTo(msisdnPIN)
	if err != nil {
		return false, fmt.Errorf("unable to read PIN: %w", err)
	}
	// encrypt the PIN
	encryptedPin, err := EncryptPIN(pin)
	if err != nil {
		return false, fmt.Errorf("unable to encrypt PIN: %w", err)
	}

	msisdnPIN.PINNumber = encryptedPin

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetPINCollectionName(), dsnap.Ref.ID, msisdnPIN,
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

// VerifyEmailOtp checks for the validity of the supplied OTP but does not invalidate it
func (s Service) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return false, err
	}

	_, emailErr := ValidateEmail(email, otp, s.firestoreClient)
	if emailErr != nil {
		return false, fmt.Errorf("email failed verification: %w", emailErr)
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		return false, fmt.Errorf("can't fetch user profile: %v", err)
	}

	verifiedEmails := userProfile.VerifiedEmails

	var foundVerifiedEmails []string
	for _, verifiedEmail := range verifiedEmails {
		foundVerifiedEmails = append(foundVerifiedEmails, verifiedEmail.Email)
	}

	if !base.StringSliceContains(foundVerifiedEmails, email) {
		verifyEmail := base.VerifiedEmail{
			Email:    email,
			Verified: true,
		}
		verifiedEmails = append(verifiedEmails, verifyEmail)
		userProfile.VerifiedEmails = verifiedEmails

		err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile,
		)
		if err != nil {
			return false, fmt.Errorf("unable to update user profile: %v", err)
		}
	}

	return true, nil
}

// CreateSignUpMethod attahces a users sign up method to a user's UID
func (s Service) CreateSignUpMethod(ctx context.Context, signUpMethod SignUpMethod) (bool, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	validSignUpMethod := signUpMethod.IsValid()
	if !validSignUpMethod {
		return false, fmt.Errorf("%v is not an allowed sign up method choice", signUpMethod.String())
	}

	signUpInfo := SignUpInfo{
		UID:          uid,
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
	s.checkPreconditions()

	collection := s.firestoreClient.Collection(s.GetSignUpInfoCollectionName())
	query := collection.Where("uid", "==", id)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return "", fmt.Errorf("unable to fetch sign up info: %w", err)
	}
	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("more than one snapshot found (they have %d)", len(docs))
		}
	}
	if len(docs) == 0 {
		return "", nil
	}
	dsnap := docs[0]

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
	ctx context.Context,
	services PractitionerServiceInput,
	otherServices *OtherPractitionerServiceInput) (bool, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't fetch the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetPractitionerCollectionName(), "profile.verifiedIdentifiers")
	if err != nil {
		return false, fmt.Errorf("unable to retreive practitioner: %v", err)
	}

	if dsnap == nil {
		return false, fmt.Errorf("no practitioner record found for profile with identifier %v", uid)
	}

	practitioner := &Practitioner{}
	err = dsnap.DataTo(practitioner)
	if err != nil {
		return false, fmt.Errorf("unable to read practitioner information: %v", err)
	}

	offeredServices := practitioner.Services.Services
	otherOfferedServices := practitioner.Services.OtherServices

	var foundServices []string
	for _, offeredService := range offeredServices {
		foundServices = append(foundServices, offeredService.String())
	}

	var otherFoundServices []string
	otherFoundServices = append(otherFoundServices, otherOfferedServices...)

	for _, service := range services.Services {
		validservice := service.IsValid()
		if !validservice {
			return false, fmt.Errorf("%v is not an allowed service enum", service.String())
		}

		if service == "OTHER" {
			if otherServices == nil {
				return false, fmt.Errorf("specify other services after selecting Others")
			}
			if !base.StringSliceContains(foundServices, service.String()) {
				offeredServices = append(offeredServices, service)
				for i, v := range offeredServices {
					if v == "OTHER" {
						offeredServices = append(offeredServices[:i], offeredServices[i+1:]...)
					}
				}
			}

			for _, otherService := range otherServices.OtherServices {
				if !base.StringSliceContains(otherFoundServices, otherService) {
					otherOfferedServices = append(otherOfferedServices, otherService)
				}
			}

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
		if !base.StringSliceContains(foundServices, service.String()) {
			offeredServices = append(offeredServices, service)
		}
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

// ParseUserProfileFromContextOrUID parses a user profile from either the logged-in context or uid
func (s Service) ParseUserProfileFromContextOrUID(ctx context.Context, uid *string) (*base.UserProfile, error) {
	if uid != nil {
		return s.GetProfile(ctx, *uid)
	}
	return s.UserProfile(ctx)
}

// SaveMemberCoverToFirestore saves users cover details to firebase
func (s Service) SaveMemberCoverToFirestore(ctx context.Context, payerName, memberNumber, memberName string, PayerSladeCode int) error {
	cover := base.Cover{
		PayerName:      payerName,
		MemberName:     memberName,
		MemberNumber:   memberNumber,
		PayerSladeCode: PayerSladeCode,
	}

	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		log.Printf("unable to retrieve user profile snapshot for the logged in user: %v", err)
		return fmt.Errorf("system error: unable to retrieve user profile. Please report a bug")
	}

	userProfile, err := s.UserProfile(ctx)
	if err != nil {
		log.Printf("unable to retrieve user profile snapshot for the logged in user: %v", err)
		return fmt.Errorf("system error: unable to retrieve user profile. Please report a bug")
	}

	existingCovers := userProfile.Covers
	exist := false

	for _, profileCover := range existingCovers {
		if profileCover.MemberNumber == cover.MemberNumber && profileCover.PayerSladeCode == cover.PayerSladeCode {
			exist = true
		}
	}

	if !exist {
		existingCovers = append(existingCovers, cover)
		userProfile.Covers = existingCovers
		err := base.UpdateRecordOnFirestore(s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, userProfile)
		if err != nil {
			log.Printf("unable to save update on user profile %v for member %v: %v", userProfile, memberNumber, err)
			return fmt.Errorf("system error: unable to update user profile. Please report a bug")
		}
	}

	return nil
}

// DeleteUser deletes a user records given their uid
func (s Service) DeleteUser(ctx context.Context, uid string) error {
	s.checkPreconditions()

	profile, err := s.GetProfile(ctx, uid)
	if err != nil {
		return fmt.Errorf("unable to get user profile: %v", err)
	}

	for _, profileUID := range profile.VerifiedIdentifiers {
		params := (&auth.UserToUpdate{}).
			Disabled(true)
		_, err := s.firebaseAuth.UpdateUser(ctx, profileUID, params)
		if err != nil {
			return fmt.Errorf("error updating user: %v", err)
		}
	}

	profile.Active = false
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrieve user profile doc snapshot: %v", err)
	}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetUserProfileCollectionName(), dsnap.Ref.ID, profile,
	)
	if err != nil {
		return fmt.Errorf("unable to update user profile: %v", err)
	}

	return nil
}

// VerifySignUpPhoneNumber does a sanity check on the supplied phone number, that is,
// it checks if a record of the phone number exists in both our collection and
// Firebase accounts. If it doesn't then an otp is generated and sent to the phone number.
func (s Service) VerifySignUpPhoneNumber(ctx context.Context, phone string) (map[string]interface{}, error) {
	s.checkPreconditions()

	defaultData := map[string]interface{}{
		"isNewUser": false,
		"OTP":       "",
	}

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return defaultData, fmt.Errorf("can't normalize the phone number: %v", err)
	}

	collection := s.firestoreClient.Collection(s.GetUserProfileCollectionName())
	query := collection.Where("msisdns", "array-contains", phoneNumber)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return defaultData, fmt.Errorf("can't fetch user profile: %v", err)
	}
	if len(docs) == 0 && base.IsDebug() {
		log.Printf("user with phone number %s has no user profile", phoneNumber)
	}

	_, userErr := s.firebaseAuth.GetUserByPhoneNumber(ctx, phoneNumber)
	if userErr != nil {
		newUserData := make(map[string]interface{})
		code, err := s.generateAndSendOTP(phoneNumber)
		if err != nil {
			return nil, fmt.Errorf("can't generate and send an otp to %s: %v", phoneNumber, err)
		}
		newUserData["OTP"] = code
		newUserData["isNewUser"] = true

		return newUserData, nil
	}

	return defaultData, nil
}

func (s Service) generateAndSendOTP(phone string) (string, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	defaultOTP := ""
	resp, err := s.Otp.MakeRequest(http.MethodPost, SendOtp, body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to generate and send otp: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return defaultOTP, fmt.Errorf("unable to generate and send otp, with status code %v", resp.StatusCode)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to convert response to string: %v", err)
	}

	return string(code), nil
}

// CreateUserByPhone represents logic to create a user via their phoneNumber
func (s Service) CreateUserByPhone(ctx context.Context, phoneNumber string) (*CreatedUserResponse, error) {
	// validate phone number
	phone, err := base.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("CreateUserByPhone: invalid phone number: %w", err)
	}
	// get or create user via thier phone number
	user, err := base.GetOrCreatePhoneNumberUser(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("CreateUserByPhone: unable to create firebase user: %w", err)
	}
	// generate a token for the user
	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	if tokenErr != nil {
		return nil, fmt.Errorf("CreateUserByPhone: unable to get or create custom token: %w", tokenErr)
	}
	//  create a profile for the user
	userProfile, err := s.CreateUserProfile(ctx, phone, user.UID)
	if err != nil {
		return nil, fmt.Errorf("CreateUserByPhone: unable to get or create a profile for the user: %w", err)
	}
	// prepare payload to return as response
	createdUser := &CreatedUserResponse{
		UserProfile: userProfile,
		CustomToken: &customToken,
	}
	return createdUser, nil
}

// PhoneSignIn authenticates a phone number and a pin and generates a firebase token
// that allows a user to access our application
func (s Service) PhoneSignIn(ctx context.Context, phoneNumber, pin string) (*PhoneSignInResponse, error) {
	s.checkPreconditions()

	u, err := s.firebaseAuth.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("PhoneSignIn: error getting user with phone %s: %w", phoneNumber, err)
	}

	dsnap, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("PhoneSignIn: unable to retrieve pin given the msisdn %s: %v", phoneNumber, err)
	}
	if dsnap == nil {
		return nil, fmt.Errorf("PhoneSignIn: unable to retrieve user PIN")
	}

	msisdnPin := &PIN{}
	err = dsnap.DataTo(msisdnPin)
	if err != nil {
		return nil, fmt.Errorf("PhoneSignIn: unable to read PIN: %w", err)
	}

	err = s.encryptExistingPin(msisdnPin, dsnap)
	if err != nil {
		return nil, err
	}

	isMatch, err := ComparePIN(msisdnPin.PINNumber, pin)
	if err != nil {
		return nil, fmt.Errorf("unable to match PIN Number provided. A wrong pin has been supplied")
	}
	if !isMatch {
		return nil, fmt.Errorf("a wrong pin has been supplied")
	}

	customToken, err := base.CreateFirebaseCustomToken(ctx, u.UID)
	if err != nil {
		return nil, fmt.Errorf("can't to generate a custom token. Please contact Slade360 for assistance")
	}

	userTokens, err := base.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, fmt.Errorf("can't authenticate the custom token. Please contact Slade360 for assistance")
	}

	loginResp := PhoneSignInResponse{
		CustomToken:  customToken,
		IDToken:      userTokens.IDToken,
		RefreshToken: userTokens.RefreshToken,
	}

	return &loginResp, nil
}
