package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
)

// Branch represents a Slade 360 Charge Master branch
type Branch struct {
	ID                    string `json:"id"                    firestore:"id"`
	Name                  string `json:"name"                  firestore:"name"`
	OrganizationSladeCode string `json:"organizationSladeCode" firestore:"organizationSladeCode"`
	BranchSladeCode       string `json:"branchSladeCode"       firestore:"branchSladeCode"`
	// this won' be saved in the repository. it will be computed when fetching the supplier's allowed locations
	Default bool `json:"default"`
}

// KYCRequest represent payload required to stage kyc processing request
type KYCRequest struct {
	ID                  string                   `json:"id"                  firestore:"id"`
	ReqPartnerType      profileutils.PartnerType `json:"reqPartnerType"      firestore:"reqPartnerType"`
	ReqOrganizationType OrganizationType         `json:"reqOrganizationType" firestore:"reqOrganizationType"`
	ReqRaw              map[string]interface{}   `json:"reqRaw"              firestore:"reqRaw"`
	Processed           bool                     `json:"processed"           firestore:"processed"`
	SupplierRecord      *profileutils.Supplier   `json:"supplierRecord"      firestore:"supplierRecord"`
	Status              KYCProcessStatus         `json:"status"              firestore:"status"`
	RejectionReason     *string                  `json:"rejectionRejection"  firestore:"rejectionRejection"`
	FiledTimestamp      time.Time                `json:"filedTimeStamp"      firestore:"filedTimeStamp"`
	ProcessedTimestamp  time.Time                `json:"processedTimeStamp"  firestore:"processedTimeStamp"`

	// points to the userProfile id of the addmin that has processed the KYC
	ProcessedBy string `json:"processedBy" firestore:"processedBy"`
}

// BusinessPartner represents a Slade 360 Charge Master business partner
type BusinessPartner struct {
	ID        string  `json:"id"         firestore:"id"`
	Name      string  `json:"name"       firestore:"name"`
	SladeCode string  `json:"slade_code" firestore:"sladeCode"`
	Parent    *string `json:"parent"     firestore:"parent"`
}

// PIN represents a user's PIN information
type PIN struct {
	ID        string `json:"id"        firestore:"id"`
	ProfileID string `json:"profileID" firestore:"profileID"`
	PINNumber string `json:"pinNumber" firestore:"pinNumber"`
	Salt      string `json:"salt"      firestore:"salt"`

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
	Criticism         string    `json:"criticism"         firestore:"criticism"`
	Suggestions       string    `json:"suggestions"       firestore:"suggestions"`
	UID               string    `json:"uid"               firestore:"uid"`
	Timestamp         time.Time `json:"timestamp"         firestore:"timestamp"`
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
	ID                        string                          `json:"id"                        firestore:"id"`
	ProfileID                 string                          `json:"profileID"                 firestore:"profileID"`
	MembershipNumber          string                          `json:"membershipNumber"          firestore:"membershipNumber"`
	Employment                EmploymentType                  `json:"employmentType"`
	IDDocType                 enumutils.IdentificationDocType `json:"IDDocType"`
	IDNumber                  string                          `json:"IDNumber"                  firestore:"IDNumber"`
	IdentificationCardPhotoID string                          `json:"identificationCardPhotoID" firestore:"identificationCardPhotoID"`
	NHIFCardPhotoID           string                          `json:"nhifCardPhotoID"           firestore:"nhifCardPhotoID"`
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

// Microservice identifies a micro-service that conforms to the Apollo Graphqql
// federation specification. These microservices are composed by an Apollo
// Gateway into a single data graph.
type Microservice struct {
	ID          string `json:"id"          firestore:"id"`
	Name        string `json:"name"        firestore:"name"`
	URL         string `json:"url"         firestore:"url"`
	Description string `json:"description" firestore:"description"`
}

// IsNode marks this model as a GraphQL Relay Node
func (m *Microservice) IsNode() {}

// GetID returns the micro-service's ID
func (m *Microservice) GetID() firebasetools.ID {
	return firebasetools.IDValue(m.ID)
}

// SetID sets the microservice's ID
func (m *Microservice) SetID(id string) {
	m.ID = id
}

// IsEntity marks the struct as an Apollo Federation entity
func (m *Microservice) IsEntity() {}

// MicroserviceStatus denotes the status of a deployed microservice
// shows if the revision is serving HTTP request
type MicroserviceStatus struct {
	Service *Microservice `json:"service"`
	Active  bool          `json:"active"`
}

// NavigationGroup is the grouping of related navigation actions based on resource
type NavigationGroup string

// NavigationAction is the menu rendered to PRO users for navigating the app
type NavigationAction struct {
	Group              NavigationGroup          `json:"code"`
	Title              string                   `json:"title"`
	OnTapRoute         string                   `json:"onTapRoute"`
	Icon               string                   `json:"icon"`
	Favorite           bool                     `json:"favorite"`
	HasParent          bool                     `json:"isParent"`
	Nested             []interface{}            `json:"nested"`
	RequiredPermission *profileutils.Permission `json:"requires"`

	// Sequence Number assigns a priority to an action
	// the number is used when sorting/ordering navigation actions
	// Actions with a higher sequence number appear at the top i.e ascending order
	SequenceNumber int `json:"sequenceNumber"`
}

// RoleRevocationLog represents a log for revoking a users role
// used when removing a role from a user i.e user deactivation
type RoleRevocationLog struct {
	// Unique identifier for a revocation
	ID string `json:"id" firestore:"id"`

	// profile of user whose role is being revoked
	ProfileID string `json:"profileID" firestore:"profileID"`

	// ID of role being revoked
	RoleID string `json:"roleID" firestore:"roleID"`

	// Reason role is being revoked
	Reason string `json:"reason" firestore:"reason"`

	// CreatedBy is the Profile ID of the user removing the role.
	CreatedBy string `json:"createdBy,omitempty" firestore:"createdBy"`

	// Created is the timestamp indicating when the role was created
	Created time.Time `json:"created" firestore:"created"`
}

//AdminProfile is the profile of all users who permform admin tasks
type AdminProfile struct {
	ID             string `json:"id,omitempty"`
	ProfileID      string `json:"profileID,omitempty"`
	OrganizationID string `json:"organizationID,omitempty"`
}

//AgentProfile is the profile of all users who permform agent tasks
type AgentProfile struct {
	ID        string    `json:"id,omitempty"`
	ProfileID string    `json:"profileID,omitempty"`
	AgentType AgentType `json:"agentType,omitempty"`
}
