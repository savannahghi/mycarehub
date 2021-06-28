package domain

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

// Branch represents a Slade 360 Charge Master branch
type Branch struct {
	ID                    string `json:"id" firestore:"id"`
	Name                  string `json:"name" firestore:"name"`
	OrganizationSladeCode string `json:"organizationSladeCode" firestore:"organizationSladeCode"`
	BranchSladeCode       string `json:"branchSladeCode" firestore:"branchSladeCode"`
	// this won' be saved in the repository. it will be computed when fetching the supplier's allowed locations
	Default bool `json:"default"`
}

// KYCRequest represent payload required to stage kyc processing request
type KYCRequest struct {
	ID                  string                 `json:"id" firestore:"id"`
	ReqPartnerType      base.PartnerType       `json:"reqPartnerType" firestore:"reqPartnerType"`
	ReqOrganizationType OrganizationType       `json:"reqOrganizationType" firestore:"reqOrganizationType"`
	ReqRaw              map[string]interface{} `json:"reqRaw" firestore:"reqRaw"`
	Processed           bool                   `json:"processed" firestore:"processed"`
	SupplierRecord      *base.Supplier         `json:"supplierRecord" firestore:"supplierRecord"`
	Status              KYCProcessStatus       `json:"status" firestore:"status"`
	RejectionReason     *string                `json:"rejectionRejection" firestore:"rejectionRejection"`
	FiledTimestamp      time.Time              `json:"filedTimeStamp" firestore:"filedTimeStamp"`
	ProcessedTimestamp  time.Time              `json:"processedTimeStamp" firestore:"processedTimeStamp"`

	// points to the userProfile id of the addmin that has processed the KYC
	ProcessedBy string `json:"processedBy" firestore:"processedBy"`
}

// BusinessPartner represents a Slade 360 Charge Master business partner
type BusinessPartner struct {
	ID        string  `json:"id" firestore:"id"`
	Name      string  `json:"name" firestore:"name"`
	SladeCode string  `json:"slade_code" firestore:"sladeCode"`
	Parent    *string `json:"parent" firestore:"parent"`
}

// PIN represents a user's PIN information
type PIN struct {
	ID        string `json:"id" firestore:"id"`
	ProfileID string `json:"profileID" firestore:"profileID"`
	PINNumber string `json:"pinNumber" firestore:"pinNumber"`
	Salt      string `json:"salt" firestore:"salt"`

	// Flags the PIN as temporary and should be changed by user
	IsOTP bool `json:"isOTP" firestore:"isOTP"`
}

// SetPINRequest payload to set PIN information
type SetPINRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	PIN         string `json:"pin"`
}

// ChangePINRequest payload to set or change PIN information
type ChangePINRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	PIN         string `json:"pin"`
	OTP         string `json:"otp"`
}

// PostVisitSurvey is used to record and retrieve post visit surveys from Firebase
type PostVisitSurvey struct {
	LikelyToRecommend int       `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string    `json:"criticism" firestore:"criticism"`
	Suggestions       string    `json:"suggestions" firestore:"suggestions"`
	UID               string    `json:"uid" firestore:"uid"`
	Timestamp         time.Time `json:"timestamp" firestore:"timestamp"`
}

// UserAddresses represents a user's home and work addresses
type UserAddresses struct {
	HomeAddress ThinAddress `json:"homeAddress"`
	WorkAddress ThinAddress `json:"workAddress"`
}

// ThinAddress represents an addresses lat-long
type ThinAddress struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NHIFDetails represents a user's thin NHIF details
type NHIFDetails struct {
	ID                        string                     `json:"id" firestore:"id"`
	ProfileID                 string                     `json:"profileID" firestore:"profileID"`
	MembershipNumber          string                     `json:"membershipNumber" firestore:"membershipNumber"`
	Employment                EmploymentType             `json:"employmentType"`
	IDDocType                 base.IdentificationDocType `json:"IDDocType"`
	IDNumber                  string                     `json:"IDNumber" firestore:"IDNumber"`
	IdentificationCardPhotoID string                     `json:"identificationCardPhotoID" firestore:"identificationCardPhotoID"`
	NHIFCardPhotoID           string                     `json:"nhifCardPhotoID" firestore:"nhifCardPhotoID"`
}

//USSDLeadDetails represents ussd details from AIT
type USSDLeadDetails struct {
	ID          string `json:"id" firestore:"id"`
	Level       int    `json:"level" firestore:"level"`
	PhoneNumber string `json:"phoneNumber" firestore:"phoneNumber"`
	SessionID   string `json:"sessionID" firestore:"sessionID"`
	FirstName   string `json:"firstName" firestore:"firstName"`
	LastName    string `json:"lastName" firestore:"lastName"`
	DateOfBirth string `json:"dob" firestore:"dob"`
	PIN         string `json:"pin" firestore:"pin"`
}

// CRMContact represents a stored CRM contact
type CRMContact struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DOB         string `json:"dob,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	OptOut      string `json:"opt_out,omitempty"`
	TimeStamp   string `json:"time_stamp,omitempty"`
	IsSynced    string `json:"is_synced,omitempty"`
}
