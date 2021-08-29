package dto

import (
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
)

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
