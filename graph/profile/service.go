// Package profile maintains user (consumer and practitioner) profiles
package profile

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
	"gitlab.slade360emr.com/go/authorization/graph/authorization"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/mailgun"
)

// configuration constants
const (
	UserProfileCollectionName           = "user_profiles"
	PractitionerCollectionName          = "practitioners"
	SurveyCollectionName                = "post_visit_survey"
	HealthcashRootCollectionName        = "healthcash"
	HealthcashDepositsCollectionName    = "healthcash_deposits"
	HealthcashWithdrawalsCollectionName = "healthcash_withdrawals"
	HealthcashWelcomeBonusAmount        = 1000
	HealthcashCurrency                  = "KES"
	EmailSignupSubject                  = "Thank you for signing up"
	EmailWelcomeSubject                 = "Welcome to Slade 360 HealthCloud"
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

	return &Service{
		firestoreClient: firestore,
		firebaseAuth:    auth,
		emailService:    mailgun.NewService(),
	}
}

// Service is an authentication service. It handles authentication related
// issues e.g user profiles
type Service struct {
	firestoreClient *firestore.Client
	firebaseAuth    *auth.Client
	emailService    *mailgun.Service
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
}

// RetrieveUserProfileFirebaseDocSnapshotByUID retrieves the user profile of a
// specified user
func (s Service) RetrieveUserProfileFirebaseDocSnapshotByUID(
	ctx context.Context,
	uid string,
) (*firestore.DocumentSnapshot, error) {
	collection := s.firestoreClient.Collection(UserProfileCollectionName)
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
			s.firestoreClient, UserProfileCollectionName, newProfile)
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
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
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}
	userProfile.IsTester = isTester(ctx, userProfile.Emails)
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
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
		s.firestoreClient, PractitionerCollectionName, practitioner)
	if err != nil {
		return false, fmt.Errorf("unable to save practitioner info: %w", err)
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

	_, _, err := s.emailService.SendEmail(EmailSignupSubject, text, emailaddress)
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return false, fmt.Errorf("unable to update user profile: %v", err)
	}
	return true, nil
}

// CompleteSignup allocates the sign-up bonus
func (s Service) CompleteSignup(ctx context.Context) (*base.Decimal, error) {
	s.checkPreconditions()
	profile, err := s.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	// do not re-process approved profiles
	if profile.IsApproved {
		return s.HealthcashBalance(ctx)
	}

	// do not re-allocate HealthCash balance to those that already have a balance
	currentBalance, err := s.HealthcashBalance(ctx)
	if err != nil {
		return nil, err
	}
	if currentBalance.Decimal().Equal(decimal.Zero) {
		rootCollection := s.firestoreClient.Collection(HealthcashRootCollectionName)
		ts := time.Now()

		depositsCollection := rootCollection.Doc(profile.UID).Collection(HealthcashDepositsCollectionName)
		deposit := HealthcashTransaction{
			At:       ts,
			Amount:   HealthcashWelcomeBonusAmount,
			Reason:   "Welcome bonus",
			Currency: HealthcashCurrency,
		}
		_, _, err = depositsCollection.Add(ctx, deposit)
		if err != nil {
			return nil, fmt.Errorf("unable to save HealthCash deposit opening balance: %w", err)
		}

		withdrawalsCollection := rootCollection.Doc(profile.UID).Collection(HealthcashWithdrawalsCollectionName)
		withdrawal := HealthcashTransaction{
			At:       ts,
			Amount:   0.0,
			Reason:   "Opening balance",
			Currency: HealthcashCurrency,
		}
		_, _, err = withdrawalsCollection.Add(ctx, withdrawal)
		if err != nil {
			return nil, fmt.Errorf("unable to save HealthCash withdrawal opening balance: %w", err)
		}
	}

	// update the profile and mark it as approved
	dsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		return nil, err
	}
	profile.IsApproved = true
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, profile)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}

	for _, practitionerEmail := range profile.Emails {
		err = s.SendPractitionerWelcomeEmail(ctx, practitionerEmail)
		if err != nil {
			return nil, fmt.Errorf("unable to send welcome email: %w", err)
		}
	}

	bal, err := s.HealthcashBalance(ctx)
	if err != nil {
		return nil, err
	}
	return bal, nil
}

// SendPractitionerWelcomeEmail will send a welcome email to the practitioner
func (s Service) SendPractitionerWelcomeEmail(ctx context.Context, emailaddress string) error {
	s.checkPreconditions()

	text := generatePractitionerWelcomeEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	_, _, err := s.emailService.SendEmail(EmailWelcomeSubject, text, emailaddress)
	if err != nil {
		return fmt.Errorf("unable to send welcome email: %w", err)
	}

	return nil
}

// HealthcashBalance returns the logged in user's HealthCash balance
func (s Service) HealthcashBalance(ctx context.Context) (*base.Decimal, error) {
	s.checkPreconditions()
	uid, err := authorization.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	rootCollection := s.firestoreClient.Collection(HealthcashRootCollectionName)

	depositsCollection := rootCollection.Doc(uid).Collection(
		HealthcashDepositsCollectionName)
	deposits, err := depositsCollection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("can't retrieve deposits: %w", err)
	}
	depositsTotal := 0.0
	for _, deposit := range deposits {
		trans := HealthcashTransaction{}
		err = deposit.DataTo(&trans)
		if err != nil {
			return nil, fmt.Errorf(
				"%#v is not a valid healthcash transaction: %w", deposit, err)
		}
		depositsTotal += trans.Amount
	}

	withdrawalsCollection := rootCollection.Doc(uid).Collection(
		HealthcashWithdrawalsCollectionName)
	withdrawals, err := withdrawalsCollection.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("can't retrieve withdrawals: %w", err)
	}
	withdrawalsTotal := 0.0
	for _, withdrawal := range withdrawals {
		trans := HealthcashTransaction{}
		err = withdrawal.DataTo(&trans)
		if err != nil {
			return nil, fmt.Errorf(
				"%#v is not a valid healthcash transaction: %w", withdrawal, err)
		}
		depositsTotal += trans.Amount
	}

	balance := depositsTotal - withdrawalsTotal
	balanceDecimal := decimal.NewFromFloat(balance)
	balanceAPIDecimal := base.Decimal(balanceDecimal)
	return &balanceAPIDecimal, nil
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
	feedbackCollection := s.firestoreClient.Collection(SurveyCollectionName)
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

// TODO Separate practitioner and consumer profiles - isApproved
// TODO practitionerTermsOfServiceAccepted
