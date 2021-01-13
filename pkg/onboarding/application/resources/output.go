package resources

import (
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// PINOutput represents a user's PIN information
type PINOutput struct {
	ProfileID string `json:"profileID" firestore:"profileID"`
	PINNumber string `json:"pinNumber" firestore:"pinNumber"`
}

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

// AccountRecoveryPhonesResponse  payload sent back to the frontend when recovery an account
type AccountRecoveryPhonesResponse struct {
	MaskedPhoneNumbers   []string `json:"maskedPhoneNumbers"`
	UnMaskedPhoneNumbers []string `json:"unMaskedPhoneNumbers"`
}

// OKResp is used to return OK responses in inter-service calls
type OKResp struct {
	Status string `json:"status"`
}
