package domain

import "time"

// QueryUsersResponse models the response that is returned by getstream API when fetching users
type QueryUsersResponse struct {
	Users []*GetStreamUser `json:"users"`
}

// GetStreamUser models the getstream user data structure
type GetStreamUser struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	// Image string   `json:"image,omitempty"`
	Role string `json:"role,omitempty"`
	// Teams []string `json:"teams,omitempty"`

	// Online    bool `json:"online,omitempty"`
	// Invisible bool `json:"invisible,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// LastActive *time.Time `json:"last_active,omitempty"`

	// ExtraData                map[string]interface{} `json:"-"`
	// RevokeTokensIssuedBefore *time.Time             `json:"revoke_tokens_issued_before,omitempty"`
}

// QueryOption are optional parameters to pass to the API. It helps in filtering. The 'Filter' value is  required.
type QueryOption struct {
	// https://getstream.io/chat/docs/#query_syntax
	Filter map[string]interface{} `json:"filter_conditions,omitempty"`
	Sort   []*SortOption          `json:"sort,omitempty"`

	UserID       string `json:"user_id,omitempty"`
	Limit        int    `json:"limit,omitempty"`  // pagination option: limit number of results
	Offset       int    `json:"offset,omitempty"` // pagination option: offset to return items from
	MessageLimit *int   `json:"message_limit,omitempty"`
	MemberLimit  *int   `json:"member_limit,omitempty"`
}

// SortOption structures the sorting method. It has a field name to sort by and takes direction as integers
// [-1 or 1]
type SortOption struct {
	Field     string `json:"field"`
	Direction int    `json:"direction"`
}
