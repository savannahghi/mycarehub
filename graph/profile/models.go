package profile

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

// KMPDUPractitioner is used to serialize a records of a particular practitioner registered with KMPDU
type KMPDUPractitioner struct {
	base.Model

	Name           string `json:"name"`
	Regno          string `json:"regno"`
	Address        string `json:"address"`
	Qualifications string `json:"qualifications"`
	Speciality     string `json:"specialty"`
	Subspeciality  string `json:"subspeciality"`
	Licensetype    string `json:"licensetype"`
	Active         string `json:"active"`
}

// KMPDUPractitionerConnection is used to return lists of practitioners registered with KMPDU.
type KMPDUPractitionerConnection struct {
	Edges    []*KMPDUPractitionerEdge `json:"edges"`
	PageInfo *base.PageInfo           `json:"pageInfo"`
}

// KMPDUPractitionerEdge is used to represent practitioners in Relay type lists.
type KMPDUPractitionerEdge struct {
	Cursor *string            `json:"cursor"`
	Node   *KMPDUPractitioner `json:"node"`
}

// Practitioner is used to serialize practitioner profile details.
// These details are in addition to the user profile that all users get.
type Practitioner struct {
	Profile                  UserProfile                `json:"profile" firestore:"profile"`
	License                  string                     `json:"license"`
	Cadre                    PractitionerCadre          `json:"cadre"`
	Specialty                base.PractitionerSpecialty `json:"specialty"`
	ProfessionalProfile      base.Markdown              `json:"professionalProfile"`
	AverageConsultationPrice float64                    `json:"averageConsultationPrice"`
	Services                 ServicesOffered            `json:"services_offered"`
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

// IsEntity ...
func (c Cover) IsEntity() {}

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

	IsTester      bool          `json:"isTester" firestore:"isTester"`
	CanExperiment bool          `json:"canExperiment" firestore:"canExperiment"`
	Language      base.Language `json:"language" firestore:"language"`

	// used to determine whether to persist asking the user on the UI
	AskAgainToSetIsTester      bool `json:"askAgainToSetIsTester" firestore:"askAgainToSetIsTester"`
	AskAgainToSetCanExperiment bool `json:"askAgainToSetCanExperiment" firestore:"askAgainToSetCanExperiment"`
	IsEmailVerified            bool `json:"isEmailVerified" firestore:"isEmailVerified"`
	IsMsisdnVerified           bool `json:"isMsisdnVerified" firestore:"isMsisdnVerified"`
	HasPin                     bool `json:"hasPin" firestore:"hasPin"`
	HasSupplierAccount         bool `json:"hasSupplierAccount" firestore:"hasSupplierAccount"`
	HasCustomerAccount         bool `json:"hasCustomerAccount" firestore:"hasCustomerAccount"`
	PractitionerHasServices    bool `json:"practitionerHasServices" firestore:"practitionerHasServices"`
}

// IsEntity ...
func (u UserProfile) IsEntity() {}

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
	CanExperiment                      bool    `json:"canExperiment" firestore:"canExperiment"`

	// used to determine whether to persist asking the user on the UI
	AskAgainToSetIsTester      bool `json:"askAgainToSetIsTester" firestore:"askAgainToSetIsTester"`
	AskAgainToSetCanExperiment bool `json:"askAgainToSetCanExperiment" firestore:"askAgainToSetCanExperiment"`
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

// PIN is used to store a PIN (Personal Identifiation Number) associated
// to a phone number sign up to Firebase
type PIN struct {
	UID     string `json:"uid" firestore:"uid"`
	MSISDN  string `json:"msisdn,omitempty" firestore:"msisdn"`
	PIN     string `json:"pin,omitempty" firestore:"pin"`
	IsValid bool   `json:"isValid,omitempty" firestore:"isValid"`
}

// PinRecovery stores information required in resetting and updating a forgotten pin
type PinRecovery struct {
	MSISDN string `json:"msisdn"`
	PIN    string `json:"pin"`
	OTP    string `json:"otp"`
}

// OtpResponse returns an otp
type OtpResponse struct {
	OTP string `json:"otp"`
}

// SignUpInfo stores the user UID and their sign up method
type SignUpInfo struct {
	UID          string       `json:"uid" firestore:"uid"`
	SignUpMethod SignUpMethod `json:"signupmethod" firestore:"signupmethod"`
}

// Customer used to create a customer request payload
type Customer struct {
	UserProfile        UserProfile        `json:"userprofile,omitempty" firestore:"userprofile"`
	CustomerID         string             `json:"id,omitempty" firestore:"customerid"`
	ReceivablesAccount ReceivablesAccount `json:"receivables_account,omitempty"`
	CustomerKYC        CustomerKYC        `json:"customer_kyc,omitempty"`
}

// ReceivablesAccount stores a customer's receivables account info
type ReceivablesAccount struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	IsActive    bool   `json:"is_active,omitempty"`
	Number      string `json:"number,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Description string `json:"description,omitempty"`
}

// CustomerKYC stores information required to know your customer
type CustomerKYC struct {
	KRAPin     string `json:"kra_pin,omitempty"`
	Occupation string `json:"occupation,omitempty"` // Should this be an enum?
	IDNumber   string `json:"id_number,omitempty"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	// Beneficiary
}

// CustomerKYCInput stores customerKYC input details
type CustomerKYCInput struct {
	KRAPin     string `json:"KRAPin"`
	Occupation string `json:"occupation"`
	IDNumber   string `json:"idNumber"`
	Address    string `json:"address"`
	City       string `json:"city"`
}

// OtherPractitionerServiceInput ..
type OtherPractitionerServiceInput struct {
	OtherServices []string `json:"otherServices"`
}

// PractitionerServiceInput ..
type PractitionerServiceInput struct {
	Services []PractitionerService `json:"services"`
}

// ServicesOffered ..
type ServicesOffered struct {
	Services      []PractitionerService `json:"services"`
	OtherServices []string              `json:"otherServices"`
}

// PayablesAccount stores a supplier's payables account info
type PayablesAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	Number      string `json:"number"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

// Supplier used to create a supplier request payload
type Supplier struct {
	UserProfile     UserProfile     `json:"userProfile" firestore:"userprofile"`
	SupplierID      string          `json:"id" firestore:"supplierid"`
	PayablesAccount PayablesAccount `json:"payables_account"`
}

// StatusResponse creates a status response for requests
type StatusResponse struct {
	Status string `json:"status"`
}

// BusinessPartnerUID is the business partner uid used in requests
type BusinessPartnerUID struct {
	UID string `json:"uid"`
}

// CustomerResponse returns customer accounts for interservice communication
type CustomerResponse struct {
	CustomerID         string             `json:"customer_id"`
	ReceivablesAccount ReceivablesAccount `json:"receivables_account"`
}

// SupplierResponse returns supplier accounts for interservice communication
type SupplierResponse struct {
	SupplierID      string          `json:"supplier_id"`
	PayablesAccount PayablesAccount `json:"payables_account"`
}
