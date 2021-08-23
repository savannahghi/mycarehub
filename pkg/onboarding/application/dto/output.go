package dto

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
	crmDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

// BusinessPartnerEdge is used to serialize GraphQL Relay edges for organization
type BusinessPartnerEdge struct {
	Cursor *string                 `json:"cursor"`
	Node   *domain.BusinessPartner `json:"node"`
}

// BusinessPartnerConnection is used to serialize GraphQL Relay connections for organizations
type BusinessPartnerConnection struct {
	Edges    []*BusinessPartnerEdge  `json:"edges"`
	PageInfo *firebasetools.PageInfo `json:"pageInfo"`
}

// AgentConnection is used to serialize GraphQL Relay connections for agents
type AgentConnection struct {
	Edges    []AgentEdge            `json:"edges"`
	PageInfo firebasetools.PageInfo `json:"pageInfo"`
}

// AgentEdge is used to serialize GraphQL Edges for an agent
type AgentEdge struct {
	Cursor string `json:"cursor"`
	Node   Agent  `json:"node"`
}

// Agent represents agent with details inferred from their user profile
type Agent struct {
	ID string `json:"id"`

	PrimaryPhone string `json:"primaryPhone"`

	PrimaryEmailAddress *string `json:"primaryEmailAddress"`

	SecondaryPhoneNumbers []string `json:"secondaryPhoneNumbers"`

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses"`

	TermsAccepted bool `json:"termsAccepted,omitempty"`

	Suspended bool `json:"suspended"`

	PhotoUploadID string `json:"photoUploadID,omitempty"`

	UserBioData profileutils.BioData `json:"userBioData,omitempty"`

	// Resend PIN helps inform the whether a send new temporary PIN
	// True when the user hasn't performed the initial sign up to change PIN
	ResendPIN bool `json:"resendPIN"`

	Roles []RoleOutput `json:"roles"`
}

// Admin represents agent with details inferred from their user profile
type Admin struct {
	ID string `json:"id"`

	PrimaryPhone string `json:"primaryPhone"`

	PrimaryEmailAddress *string `json:"primaryEmailAddress"`

	SecondaryPhoneNumbers []string `json:"secondaryPhoneNumbers"`

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses"`

	TermsAccepted bool `json:"termsAccepted,omitempty"`

	Suspended bool `json:"suspended"`

	PhotoUploadID string `json:"photoUploadID,omitempty"`

	UserBioData profileutils.BioData `json:"userBioData,omitempty"`

	// Resend PIN helps inform the whether a send new temporary PIN
	// True when the user hasn't performed the initial sign up to change PIN
	ResendPIN bool `json:"resendPIN"`

	Roles []RoleOutput `json:"roles"`
}

// AccountRecoveryPhonesResponse  payload sent back to the frontend when recovery an account
type AccountRecoveryPhonesResponse struct {
	MaskedPhoneNumbers   []string `json:"maskedPhoneNumbers"`
	UnMaskedPhoneNumbers []string `json:"unmaskedPhoneNumbers"`
}

// OKResp is used to return OK response in inter-service calls
type OKResp struct {
	Status   string      `json:"status,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// NewOKResp a shortcut to create an instance of OKResp
func NewOKResp(rawResponse interface{}) *OKResp {
	return &OKResp{
		Status:   "OK",
		Response: rawResponse,
	}
}

// CreatedUserResponse is used to return a created user
type CreatedUserResponse struct {
	UID         string `json:"uid,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	ProviderID  string `json:"provider_id,omitempty"`
}

// SupplierLogin is the response returned after the user has successfully login to edi
type SupplierLogin struct {
	Branches *BranchConnection      `json:"branches,omitempty"`
	Supplier *profileutils.Supplier `json:"supplier,omitempty"`
}

// UserInfo is a collection of standard profile information for a user.
type UserInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	// In the ProviderUserInfo[] ProviderID can be a short domain name (e.g. google.com),
	// or the identity of an OpenID identity provider.
	// In UserRecord.UserInfo it will return the constant string "firebase".
	ProviderID string `json:"providerId,omitempty"`
	UID        string `json:"rawId,omitempty"`
}

// Segment represents the Segments data
type Segment struct {
	Properties     crmDomain.ContactProperties `json:"properties"       firestore:"properties"`
	Wing           string                      `json:"wing"             firestore:"wing"`
	MessageSent    string                      `json:"message_sent"     firestore:"message_sent"`
	IsSynced       string                      `json:"is_synced"        firestore:"is_synced"`
	TimeSynced     string                      `json:"time_synced"      firestore:"time_synced"`
	PayerSladeCode string                      `json:"payer_slade_code" firestore:"payersladecode"`
	MemberNumber   string                      `json:"member_number"    firestore:"membernumber"`
}

// RoleOutput is the formatted output with scopes and permissions
type RoleOutput struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Active      bool                        `json:"active"`
	Scopes      []string                    `json:"scopes"`
	Permissions []profileutils.Permission   `json:"permissions"`
	Users       []*profileutils.UserProfile `json:"users"`
}

// GroupedNavigationActions is the list of Navigation Actions sorted into primary and secondary actions
type GroupedNavigationActions struct {
	Primary   []domain.NavigationAction `json:"primary,omitempty"`
	Secondary []domain.NavigationAction `json:"secondary,omitempty"`
}
