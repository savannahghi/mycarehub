package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Community defines the payload to create a channel
type Community struct {
	ID          string    `json:"id"`
	CID         string    `json:"cid"`
	Name        string    `json:"name"`
	Disabled    bool      `json:"disabled"`
	Frozen      bool      `json:"frozen"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// The fields below are custom to our implementation
	Description string             `json:"description"`
	AgeRange    *AgeRange          `json:"ageRange"`
	Gender      []enumutils.Gender `json:"gender"`
	ClientType  []enums.ClientType `json:"clientType"`
	InviteOnly  bool               `json:"inviteOnly"`
	Members     []CommunityMember  `json:"members"`
	CreatedBy   *Member            `json:"created_by"`
}

// AgeRange defines the channel users age input
type AgeRange struct {
	LowerBound int `json:"lowerBound"`
	UpperBound int `json:"upperBound"`
}

// PostingHours defines the channel posting hours
type PostingHours struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Member represents a user and is specific to use in the context of communities
type Member struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
	Name   string `json:"name"`
	Role   string `json:"role"`

	Username string           `json:"username"`
	Gender   enumutils.Gender `json:"gender"`
}

// CommunityMember represents a user in a community and their associated additional details.
type CommunityMember struct {
	UserID           string     `json:"userID"`
	User             Member     `json:"user"`
	Role             string     `json:"role"`
	IsModerator      bool       `json:"isModerator"`
	UserType         string     `json:"userType"`
	Invited          bool       `json:"invited"`
	InviteAcceptedAt *time.Time `json:"invite_accepted_at"`
	InviteRejectedAt *time.Time `json:"invite_rejected_at"`
}
