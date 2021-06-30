package dto

import (
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// BusinessPartnerEdge is used to serialize GraphQL Relay edges for organization
type BusinessPartnerEdge struct {
	Cursor *string                 `json:"cursor"`
	Node   *domain.BusinessPartner `json:"node"`
}

// BusinessPartnerConnection is used to serialize GraphQL Relay connections for organizations
type BusinessPartnerConnection struct {
	Edges    []*BusinessPartnerEdge `json:"edges"`
	PageInfo *base.PageInfo         `json:"pageInfo"`
}

type AgentConnection struct {
	Edges    []AgentEdge
	PageInfo base.PageInfo `json:"pageInfo"`
}

type AgentEdge struct {
	Cursor string
	Node   Agent
}

// Agent represents agent with details inferred from their user profile
type Agent struct {
	ID string `json:"id"`

	PrimaryPhone string `json:"primaryPhone" `

	PrimaryEmailAddress string `json:"primaryEmailAddress" `

	SecondaryPhoneNumbers []string `json:"secondaryPhoneNumbers" `

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses" `

	TermsAccepted bool `json:"terms_accepted,omitempty" `

	Suspended bool `json:"suspended"`

	PhotoUploadID string `json:"photoUploadID,omitempty" `

	UserBioData base.BioData `json:"userBioData,omitempty" `
}

type AgentFilterInput struct {
}

type AgentSortInput struct {
}

// AccountRecoveryPhonesResponse  payload sent back to the frontend when recovery an account
type AccountRecoveryPhonesResponse struct {
	MaskedPhoneNumbers   []string `json:"maskedPhoneNumbers"`
	UnMaskedPhoneNumbers []string `json:"unmaskedPhoneNumbers"`
}

// OKResp is used to return OK responses in inter-service calls
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
	Branches *BranchConnection `json:"branches,omitempty"`
	Supplier *base.Supplier    `json:"supplier,omitempty"`
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
