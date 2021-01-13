package domain

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

// Branch represents a Slade 360 Charge Master branch
type Branch struct {
	base.Model

	ID                    string `json:"id"`
	Name                  string `json:"name"`
	OrganizationSladeCode string `json:"organizationSladeCode"`
	BranchSladeCode       string `json:"branchSladeCode"`
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
}

// BusinessPartner represents a Slade 360 Charge Master business partner
type BusinessPartner struct {
	base.Model

	ID        string  `json:"id"`
	Name      string  `json:"name"`
	SladeCode string  `json:"slade_code"`
	Parent    *string `json:"parent"`
}

// PIN represents a user's PIN information
type PIN struct {
	ID        string `json:"id" firestore:"id"`
	ProfileID string `json:"profileID" firestore:"profileID"`
	PINNumber string `json:"pinNumber" firestore:"pinNumber"`
	Salt      string `json:"salt" firestore:"salt"`
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
