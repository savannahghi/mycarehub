package profile

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

// Practitioner is used to serialize practitioner profile details.
// These details are in addition to the user profile that all users get.
type Practitioner struct {
	Profile                  UserProfile                `json:"profile" firestore:"profile"`
	License                  string                     `json:"license"`
	Cadre                    PractitionerCadre          `json:"cadre"`
	Specialty                base.PractitionerSpecialty `json:"specialty"`
	ProfessionalProfile      base.Markdown              `json:"professionalProfile"`
	AverageConsultationPrice float64                    `json:"averageConsultationPrice"`
}

// PractitionerConnection is used to return lists of practitioners.
type PractitionerConnection struct {
	Edges    []*PractitionerEdge `json:"edges"`
	PageInfo *base.PageInfo      `json:"pageInfo"`
}

// PractitionerEdge is used to represent practitioners in Relay type lists.
type PractitionerEdge struct {
	Cursor *string       `json:"cursor"`
	Node   *Practitioner `json:"node"`
}

// PractitionerSignupInput is used to sign up practitioners.
//
// The `uid` is obtained from the logged in user.
type PractitionerSignupInput struct {
	License   string                     `json:"license"`
	Cadre     PractitionerCadre          `json:"cadre"`
	Specialty base.PractitionerSpecialty `json:"specialty"`
	Emails    []*string                  `json:"emails"`
}

// Cover is used to save a user's insurance details.
type Cover struct {
	PayerName      string `json:"payerName,omitempty" firestore:"payerName"`
	PayerSladeCode int    `json:"payerSladeCode,omitempty" firestore:"payerSladeCode"`
	MemberNumber   string `json:"memberNumber,omitempty" firestore:"memberNumber"`
	MemberName     string `json:"memberName,omitempty" firestore:"memberName"`
}

// TesterWhitelist is used to maintain
type TesterWhitelist struct {
	base.Model

	Email string `json:"email" firestore:"email"`
}

// UserProfile serializes the profile of the logged in user.
type UserProfile struct {
	UID              string           `json:"uid" firestore:"uid"`
	TermsAccepted    bool             `json:"termsAccepted" firestore:"termsAccepted"`
	IsApproved       bool             `json:"isApproved" firestore:"isApproved"`
	Msisdns          []string         `json:"msisdns" firestore:"msisdns"`
	Emails           []string         `json:"emails" firestore:"emails"`
	PhotoBase64      string           `json:"photoBase64" firestore:"photoBase64"`
	PhotoContentType base.ContentType `json:"photoContentType" firestore:"photoContentType"`
	Covers           []Cover          `json:"covers" firestore:"covers"`

	DateOfBirth *base.Date   `json:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
	Gender      *base.Gender `json:"gender,omitempty" firestore:"gender,omitempty"`
	PatientID   *string      `json:"patientID,omitempty" firestore:"patientID"`
	PushTokens  []string     `json:"pushTokens" firestore:"pushTokens"`

	Name                               *string `json:"name" firestore:"name"`
	Bio                                *string `json:"bio" firestore:"bio"`
	PractitionerApproved               *bool   `json:"practitionerApproved" firestore:"practitionerApproved"`
	PractitionerTermsOfServiceAccepted *bool   `json:"practitionerTermsOfServiceAccepted" firestore:"practitionerTermsOfServiceAccepted"`

	IsTester bool `json:"isTester" firestore:"isTester"`
}

// UserProfileInput is used to create or update a user's profile.
type UserProfileInput struct {
	PhotoBase64      string              `json:"photoBase64"`
	PhotoContentType base.ContentType    `json:"photoContentType"`
	Msisdns          []*UserProfilePhone `json:"msisdns"`
	Emails           []string            `json:"emails"`

	DateOfBirth *base.Date   `json:"dateOfBirth,omitempty"`
	Gender      *base.Gender `json:"gender,omitempty"`
	PushTokens  []*string    `json:"pushTokens"`

	Name                               *string `json:"name"`
	Bio                                *string `json:"bio"`
	PractitionerApproved               *bool   `json:"practitionerApproved" firestore:"practitionerApproved"`
	PractitionerTermsOfServiceAccepted *bool   `json:"practitionerTermsOfServiceAccepted" firestore:"practitionerTermsOfServiceAccepted"`
}

// UserProfilePhone is used to input a user's phone and the corresponding OTP
// confirmation code.
type UserProfilePhone struct {
	Phone string `json:"phone"`
	Otp   string `json:"otp"`
}

// BiodataInput is used to update a user's bio-data.
type BiodataInput struct {
	DateOfBirth base.Date   `json:"dateOfBirth"`
	Gender      base.Gender `json:"gender"`

	Name *string `json:"name"`
	Bio  *string `json:"bio"`
}

// HealthcashTransaction is used to record increases in a user's HealthCash.
// In order for the aggregations to be manageable, instances of this
// in firestore MUST be nested under users' uids.
type HealthcashTransaction struct {
	At       time.Time `json:"at,omitempty" firestore:"at,omitempty"`
	Amount   float64   `json:"amount,omitempty" firestore:"amount,omitempty"`
	Reason   string    `json:"reason,omitempty" firestore:"reason,omitempty"`
	Currency string    `json:"currency,omitempty" firestore:"currency,omitempty"`
}

// PostVisitSurveyInput is used to send the results of post-visit surveys to the
// server.
type PostVisitSurveyInput struct {
	LikelyToRecommend int    `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string `json:"criticism" firestore:"criticism"`
	Suggestions       string `json:"suggestions" firestore:"suggestions"`
}

// PostVisitSurvey is used to record and retrieve post visit surveys from Firebase
type PostVisitSurvey struct {
	LikelyToRecommend int    `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string `json:"criticism" firestore:"criticism"`
	Suggestions       string `json:"suggestions" firestore:"suggestions"`

	UID       string    `json:"uid" firestore:"uid"`
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}
