package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Community defines the payload to create a channel
type Community struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	AgeRange    *AgeRange          `json:"ageRange"`
	Gender      []enumutils.Gender `json:"gender"`
	ClientType  []enums.ClientType `json:"clientType"`
	InviteOnly  bool               `json:"inviteOnly"`
	Members     []CommunityMember  `json:"members"`
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
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`

	Username string           `json:"username"`
	Gender   enumutils.Gender `json:"gender"`
}

// CommunityMember represents a user in a community and their associated additional details.
type CommunityMember struct {
	UserID      string `json:"userID"`
	User        Member `json:"user"`
	Role        string `json:"role"`
	IsModerator bool   `json:"isModerator"`
	UserType    string `json:"userType"`
}
