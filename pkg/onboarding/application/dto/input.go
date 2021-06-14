package dto

import (
	"net/url"
	"time"

	"gitlab.slade360emr.com/go/base"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
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

// SignUpInput represents the user information required to create a new account
type SignUpInput struct {
	PhoneNumber *string      `json:"phoneNumber"`
	PIN         *string      `json:"pin"`
	Flavour     base.Flavour `json:"flavour"`
	OTP         *string      `json:"otp"`
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

// SetPrimaryPhoneNumberPayload used when veriying and setting a user's primary phone number via REST
type SetPrimaryPhoneNumberPayload struct {
	PhoneNumber *string `json:"phoneNumber"`
	OTP         *string `json:"otp"`
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

// UpdateCoversPayload is used to make a REST
// request to update a user's covers in their user profile
type UpdateCoversPayload struct {
	UID                   *string    `json:"uid"`
	PayerName             *string    `json:"payerName"`
	MemberName            *string    `json:"memberName"`
	MemberNumber          *string    `json:"memberNumber"`
	PayerSladeCode        *int       `json:"payerSladeCode"`
	BeneficiaryID         *int       `json:"beneficiaryID"`
	EffectivePolicyNumber *string    `json:"effectivePolicyNumber"`
	ValidFrom             *time.Time `json:"validFrom"`
	ValidTo               *time.Time `json:"validTo"`
}

// UIDsPayload is an input of a slice of users' UIDs used
// for ISC requests to retrieve contact details of the users
type UIDsPayload struct {
	UIDs []string `json:"uids"`
}

// UserAddressInput represents a user's geo location input
type UserAddressInput struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Locality         *string `json:"locality"`
	Name             *string `json:"name"`
	PlaceID          *string `json:"placeID"`
	FormattedAddress *string `json:"formattedAddress"`
}

// NHIFDetailsInput represents a user's thin NHIF input details
type NHIFDetailsInput struct {
	MembershipNumber          string                     `json:"membershipNumber"`
	Employment                domain.EmploymentType      `json:"employmentType"`
	IDDocType                 base.IdentificationDocType `json:"IDDocType"`
	IDNumber                  string                     `json:"IDNumber"`
	IdentificationCardPhotoID string                     `json:"identificationCardPhotoID"`
	NHIFCardPhotoID           string                     `json:"nhifCardPhotoID"`
}

// PushTokenPayload represents user device push token
type PushTokenPayload struct {
	PushToken string `json:"pushTokens"`
	UID       string `json:"uid"`
}

// CustomerPubSubMessage is an `onboarding` PubSub message struct
type CustomerPubSubMessage struct {
	CustomerPayload CustomerPayload `json:"customerPayload"`
	UID             string          `json:"uid"`
}

// CustomerPayload is the customer data used to create a customer
// business partner in the ERP
type CustomerPayload struct {
	Active       bool             `json:"active"`
	PartnerName  string           `json:"partner_name"`
	Country      string           `json:"country"`
	Currency     string           `json:"currency"`
	IsCustomer   bool             `json:"is_customer"`
	CustomerType base.PartnerType `json:"customer_type"`
}

// SupplierPubSubMessage is an `onboarding` PubSub message struct
type SupplierPubSubMessage struct {
	SupplierPayload SupplierPayload `json:"supplierPayload"`
	UID             string          `json:"uid"`
}

// SupplierPayload is the supplier data used to create a supplier
// business partner in the ERP
type SupplierPayload struct {
	Active       bool             `json:"active"`
	PartnerName  string           `json:"partner_name"`
	Country      string           `json:"country"`
	Currency     string           `json:"currency"`
	IsSupplier   bool             `json:"is_supplier"`
	SupplierType base.PartnerType `json:"supplier_type"`
}

// EmailNotificationPayload is the email payload used to send email
// supplier and admins for KYC requests
type EmailNotificationPayload struct {
	SupplierName string `json:"supplier_name"`
	PartnerType  string `json:"partner_type"`
	AccountType  string `json:"account_type"`
	SubjectTitle string `json:"subject_title"`
	EmailBody    string `json:"email_body"`
	EmailAddress string `json:"email_address"`
	PrimaryPhone string `json:"primary_phone"`
}

// UserProfilePayload is used to update a user's profile.
// This payload is used for REST endpoints
type UserProfilePayload struct {
	UID           *string      `json:"uid"`
	PhotoUploadID *string      `json:"photoUploadID"`
	DateOfBirth   *base.Date   `json:"dateOfBirth,omitempty"`
	Gender        *base.Gender `json:"gender,omitempty"`
	FirstName     *string      `json:"lastName"`
	LastName      *string      `json:"firstName"`
}

// PermissionInput input required to create a permission
type PermissionInput struct {
	Action   string
	Resource string
}

// UpdateContactPSMessage represents CRM update contact Pub/Sub message
type UpdateContactPSMessage struct {
	Properties CRMDomain.ContactProperties `json:"properties"`
	Phone      string                      `json:"phone"`
}
