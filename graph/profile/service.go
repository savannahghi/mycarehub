// Package profile maintains user (consumer and practitioner) profiles
package profile

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
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
	}
}

// Service is an authentication service. It handles authentication related
// issues e.g user profiles
type Service struct {
	firestoreClient *firestore.Client
	firebaseAuth    *auth.Client
}

func (s Service) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("profile service does not have an initialized firestoreClient")
	}

	if s.firebaseAuth == nil {
		log.Panicf("profile service does not have an initialized firebaseAuth")
	}
}

// RetrieveUserProfileFirebaseDocSnapshot retrievs a raw Firebase doc snapshot
// for the logged in user's user profile or creates one if it does not exist
func (s Service) RetrieveUserProfileFirebaseDocSnapshot(
	ctx context.Context) (*firestore.DocumentSnapshot, error) {
	uid, err := authorization.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	collection := s.firestoreClient.Collection(UserProfileCollectionName)
	query := collection.Where("uid", "==", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		return nil, fmt.Errorf("system error: there is more than one user profile for this user")
	}
	if len(docs) == 0 {
		newProfile := &UserProfile{
			UID:           uid,
			IsApproved:    false,
			TermsAccepted: false,
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
	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
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
		s.firestoreClient, UserProfileCollectionName, dsnap.Ref.ID, userProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update user profile: %v", err)
	}
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

// generatePractitionerSignupEmailTemplate generates an signup email
func generatePractitionerSignupEmailTemplate() string {
	t := template.Must(template.New("signupemail").Parse(practitionerSignupEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		logger.Errorf("Error while generating template")
	}
	return buf.String()
}

// SendPractitionerSignUpEmail will send a signup email to the practitioner
func (s Service) SendPractitionerSignUpEmail(ctx context.Context, emailaddress string) error {
	text := generatePractitionerSignupEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	emailService := mailgun.NewService()
	_, _, err := emailService.SendEmail(EmailSignupSubject, text, emailaddress)
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

	bal, err := s.HealthcashBalance(ctx)
	if err != nil {
		return nil, err
	}
	return bal, nil
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
