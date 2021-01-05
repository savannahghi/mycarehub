package domain

import (
	"time"

	"gitlab.slade360emr.com/go/base"
)

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
	ID                     string                 `json:"id" firestore:"id"`
	ProfileID              *string                `json:"profileID" firestore:"profileID"`
	SupplierID             string                 `json:"supplierID" firestore:"supplierID"`
	SupplierName           string                 `json:"supplierName" firestore:"supplierName"`
	PayablesAccount        *PayablesAccount       `json:"payablesAccount"`
	SupplierKYC            map[string]interface{} `json:"supplierKYC"`
	Active                 bool                   `json:"active" firestore:"active"`
	AccountType            AccountType            `json:"accountType"`
	UnderOrganization      bool                   `json:"underOrganization"`
	IsOrganizationVerified bool                   `json:"isOrganizationVerified"`
	SladeCode              string                 `json:"sladeCode"`
	ParentOrganizationID   string                 `json:"parentOrganizationID"`
	HasBranches            bool                   `json:"hasBranches,omitempty"`
	Location               *Location              `json:"location,omitempty"`
	PartnerType            PartnerType            `json:"partnerType"`
	EDIUserProfile         *base.EDIUserProfile   `json:"ediuserprofile" firestore:"ediuserprofile"`
	PartnerSetupComplete   bool                   `json:"partnerSetupComplete" firestore:"partnerSetupComplete"`
	KYCSubmitted           bool                   `json:"kycSubmitted" firestore:"kycSubmitted"`
}

// Location is used to store a user's branch or organisation
type Location struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BranchSladeCode *string `json:"branchSladeCode"`
}

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
	ReqPartnerType      PartnerType            `json:"reqPartnerType" firestore:"reqPartnerType"`
	ReqOrganizationType OrganizationType       `json:"reqOrganizationType" firestore:"reqOrganizationType"`
	ReqRaw              map[string]interface{} `json:"reqRaw" firestore:"reqRaw"`
	Proceseed           bool                   `json:"proceseed" firestore:"proceseed"`
	SupplierRecord      *Supplier              `json:"supplierRecord" firestore:"supplierRecord"`
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

// Customer used to create a customer request payload
type Customer struct {
	ID                 string             `json:"id" firestore:"id"`
	ProfileID          *string            `json:"profileID,omitempty" firestore:"profileID"`
	CustomerID         string             `json:"customerID,omitempty" firestore:"customerID"`
	ReceivablesAccount ReceivablesAccount `json:"receivablesAccount" firestore:"profileID"`
	Active             bool               `json:"active" firestore:"active"`
}

// ReceivablesAccount stores a customer's receivables account info
type ReceivablesAccount struct {
	ID          string `json:"id" firestore:"id"`
	Name        string `json:"name" firestore:"name"`
	IsActive    bool   `json:"isActive" firestore:"isActive"`
	Number      string `json:"number" firestore:"number"`
	Tag         string `json:"tag" firestore:"tag"`
	Description string `json:"description" firestore:"description"`
}

// PIN represents a user's PIN information
type PIN struct {
	ID        string `json:"id" firestore:"id"`
	ProfileID string `json:"profileID" firestore:"profileID"`
	PINNumber string `json:"pinNumber" firestore:"pinNumber"`
	Salt      string `json:"salt" firestore:"salt"`
}

// PostVisitSurvey is used to record and retrieve post visit surveys from Firebase
type PostVisitSurvey struct {
	LikelyToRecommend int       `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string    `json:"criticism" firestore:"criticism"`
	Suggestions       string    `json:"suggestions" firestore:"suggestions"`
	UID               string    `json:"uid" firestore:"uid"`
	Timestamp         time.Time `json:"timestamp" firestore:"timestamp"`
}

// BusinessPartnerUID is the user ID used in some inter-service requests
type BusinessPartnerUID struct {
	UID *string `json:"uid"`
}

//TODO: restore commented structs when implementing profile missing methods
// // PIN is used to store a PIN (Personal Identifiation Number) associated
// // to a phone number sign up to Firebase
// type PIN struct {
// 	ProfileID string `json:"profile_id" firestore:"profileID"`
// 	MSISDN    string `json:"msisdn,omitempty" firestore:"msisdn"`
// 	PINNumber string `json:"pin_number" firestore:"pin"`
// 	IsValid   bool   `json:"isValid,omitempty" firestore:"isValid"`
// }

// // PinRecovery stores information required in resetting and updating a forgotten pin
// type PinRecovery struct {
// 	MSISDN    string `json:"msisdn" firestore:"msisdn"`
// 	PINNumber string `json:"pin_number" firestore:"PINNumber"`
// 	OTP       string `json:"otp" firestore:"otp"`
// }

// // Beneficiary stores a customer's beneficiary details
// type Beneficiary struct {
// 	Name         string                  `json:"name"`
// 	Msisdns      []string                `json:"msisdns"`
// 	Emails       []string                `json:"emails"`
// 	Relationship BeneficiaryRelationship `json:"relationship"`
// 	DateOfBirth  base.Date               `json:"dateOfBirth"`
// }

// // BeneficiaryInput stores beneficiary input details
// type BeneficiaryInput struct {
// 	Name         string                  `json:"name"`
// 	Msisdns      []string                `json:"msisdns"`
// 	Emails       []string                `json:"emails"`
// 	Relationship BeneficiaryRelationship `json:"relationship"`
// 	DateOfBirth  base.Date               `json:"dateOfBirth"`
// }

// // OtherPractitionerServiceInput ..
// type OtherPractitionerServiceInput struct {
// 	OtherServices []string `json:"otherServices"`
// }

// // PractitionerServiceInput ..
// type PractitionerServiceInput struct {
// 	Services []PractitionerService `json:"services"`
// }

// // ServicesOffered ..
// type ServicesOffered struct {
// 	Services      []PractitionerService `json:"services"`
// 	OtherServices []string              `json:"otherServices"`
// }

// // StatusResponse creates a status response for requests
// type StatusResponse struct {
// 	Status string `json:"status"`
// }

// // SendRetryOTP is an input struct for generating and
// // sending fallback otp
// type SendRetryOTP struct {
// 	Msisdn    string `json:"msisdn"`
// 	RetryStep int    `json:"retryStep"`
// }

// // UserUIDs is an input of a list of user uids for isc requests
// type UserUIDs struct {
// 	UIDs []string `json:"uids"`
// }

// // CreatedUserResponse represents payload returned after creating a user
// type CreatedUserResponse struct {
// 	UserProfile *base.UserProfile `json:"user_profile"`
// 	CustomToken *string           `json:"custom_token"`
// }

// // CreateUserViaPhoneInput represents input required to create a user via phoneNumber
// type CreateUserViaPhoneInput struct {
// 	MSISDN string `json:"msisdn"`
// }

// // PhoneSignInInput represents input required to sign in a user via phoneNumber
// type PhoneSignInInput struct {
// 	PhoneNumber string `json:"phonenumber"`
// 	Pin         string `json:"pin"`
// }

// // PhoneSignInResponse is a thin payload returned when a user signs in
// // with their phone number
// type PhoneSignInResponse struct {
// 	CustomToken  string `json:"custom_token"`
// 	IDToken      string `json:"id_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// // OKResp is used to return OK responses in inter-service calls
// type OKResp struct {
// 	Status string `json:"status"`
// }

// // SaveMemberCoverPayload deserializes inter-service requests to save
// // member covers
// type SaveMemberCoverPayload struct {
// 	PayerName      string `json:"payerName"`
// 	MemberName     string `json:"memberName"`
// 	MemberNumber   string `json:"memberNumber"`
// 	PayerSladeCode int    `json:"payerSladeCode"`
// 	UID            string `json:"uid"`
// }

// // SaveResponsePayload is used to return successful save feedback for
// // inter-service calls
// type SaveResponsePayload struct {
// 	SuccessfullySaved bool `json:"successfullySaved"`
// }

// // OTPResponse is used to return the results of requesting an OTP
// // or OTP retry.
// type OTPResponse struct {
// 	OTP string `json:"otp"`
// }

// // SupplierAccountInput is used when setting up basic/"key" supplier
// // account during onboarding
// type SupplierAccountInput struct {
// 	AccountType            AccountType    `json:"accountType"`
// 	UnderOrganization      bool           `json:"underOrganization"`
// 	IsOrganizationVerified *bool          `json:"isOrganizationVerified"`
// 	SladeCode              *string        `json:"sladeCode"`
// 	ParentOrganizationID   *string        `json:"parentOrganizationID"`
// 	Location               *LocationInput `json:"location,omitempty"`
// }

// // LocationInput is used when setting up a location (branch or parent) for a user
// type LocationInput struct {
// 	ID              string  `json:"id"`
// 	Name            string  `json:"name"`
// 	BranchSladeCode *string `json:"branchSladeCode"`
// }
