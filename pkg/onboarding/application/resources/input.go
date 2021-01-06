package resources

import (
	"net/url"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// UserProfileInput is used to create or update a user's profile.
type UserProfileInput struct {
	PhotoUploadID *string      `json:"photoUploadID"`
	DateOfBirth   *base.Date   `json:"dateOfBirth,omitempty"`
	Gender        *base.Gender `json:"gender,omitempty"`
	FirstName     *string      `json:"lastName"`
	LastName      *string      `json:"firstName"`
}

// PostVisitSurveyInput is used to send the results of post-visit surveys to the
// server.
type PostVisitSurveyInput struct {
	LikelyToRecommend int    `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string `json:"criticism" firestore:"criticism"`
	Suggestions       string `json:"suggestions" firestore:"suggestions"`
}

// BusinessPartnerFilterInput is used to supply filter parameters for organizatiom filter inputs
type BusinessPartnerFilterInput struct {
	Search    *string `json:"search"`
	Name      *string `json:"name"`
	SladeCode *string `json:"slade_code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BusinessPartnerFilterInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.Name != nil {
		vals.Add("name", *i.Name)
	}
	if i.SladeCode != nil {
		vals.Add("slade_code", *i.SladeCode)
	}
	return vals
}

// BusinessPartnerSortInput is used to supply sort input for organization list queries
type BusinessPartnerSortInput struct {
	Name      *base.SortOrder `json:"name"`
	SladeCode *base.SortOrder `json:"slade_code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BusinessPartnerSortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.SladeCode != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("slade_code", "number")
		} else {
			vals.Add("slade_code", "-number")
		}
	}
	return vals
}

// BranchSortInput is used to supply sorting input for location list queries
type BranchSortInput struct {
	Name      *base.SortOrder `json:"name"`
	SladeCode *base.SortOrder `json:"slade_code"`
}

// ToURLValues transforms the sort input to `url.Values`
func (i *BranchSortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.SladeCode != nil {
		if *i.SladeCode == base.SortOrderAsc {
			vals.Add("slade_code", "number")
		} else {
			vals.Add("slade_code", "-number")
		}
	}
	return vals
}

// BranchEdge is used to serialize GraphQL Relay edges for locations
type BranchEdge struct {
	Cursor *string        `json:"cursor"`
	Node   *domain.Branch `json:"node"`
}

// BranchConnection is used tu serialize GraphQL Relay connections for locations
type BranchConnection struct {
	Edges    []*BranchEdge  `json:"edges"`
	PageInfo *base.PageInfo `json:"pageInfo"`
}

// BranchFilterInput is used to supply filter parameters for locatioon list queries
type BranchFilterInput struct {
	Search               *string `json:"search"`
	SladeCode            *string `json:"sladeCode"`
	ParentOrganizationID *string `json:"parentOrganizationID"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BranchFilterInput) ToURLValues() url.Values {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.SladeCode != nil {
		vals.Add("slade_code", *i.SladeCode)
	}
	if i.ParentOrganizationID != nil {
		vals.Add("parent", *i.ParentOrganizationID)
	}
	return vals
}

// PhoneNumberPayload used when verifying a phone number.
type PhoneNumberPayload struct {
	PhoneNumber *string `json:"phoneNumber"`
}

// SignUpPayload used when calling the REST API to create a new account
type SignUpPayload struct {
	PhoneNumber *string      `json:"phoneNumber"`
	PIN         *string      `json:"pin"`
	Flavour     base.Flavour `json:"flavour"`
}

// ChangePINRequest payload to set or change PIN information
type ChangePINRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	PIN         string `json:"pin"`
	OTP         string `json:"otp"`
}

// LoginPayload used when calling the REST API to log a user in
type LoginPayload struct {
	PhoneNumber *string      `json:"phoneNumber"`
	PIN         *string      `json:"pin"`
	Flavour     base.Flavour `json:"flavour"`
}

// SendRetryOTPPayload is used when calling the REST API to resend an otp
type SendRetryOTPPayload struct {
	Phone     *string `json:"phoneNumber"`
	RetryStep *int    `json:"retryStep"`
}

// RefreshTokenExchangePayload is marshalled into JSON
// and sent to the Firebase Auth REST API when exchanging a
// refresh token for an ID token that can be used to make API calls
type RefreshTokenExchangePayload struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenPayload is used when calling the REST API to
// exchange a Refresh Token for new ID Token
type RefreshTokenPayload struct {
	RefreshToken *string `json:"refreshToken"`
}

// UIDPayload is the user ID used in some inter-service requests
type UIDPayload struct {
	UID *string `json:"uid"`
}
