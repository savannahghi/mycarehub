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

	Username  string                 `json:"username"`
	Gender    enumutils.Gender       `json:"gender"`
	UserType  string                 `json:"userType"`
	ExtraData map[string]interface{} `json:"extraData"`
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

// CommunityMetadata is extra data that is added to the communities. This data will
// be used in showing recommended channels to the users
type CommunityMetadata struct {
	MinimumAge  int                `json:"minimumAge"`
	MaximumAge  int                `json:"maximumAge"`
	Gender      []enumutils.Gender `json:"gender"`
	ClientType  []enums.ClientType `json:"clientType"`
	InviteOnly  bool               `json:"inviteOnly"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

// MemberMetadata is extra data that is associated with a getstream user.
type MemberMetadata struct {
	UserID   string `json:"userID"`
	UserType string `json:"userType"`
	NickName string `json:"nickName"`
}

// AIModerationResponse is the response from the AIModeration service
type AIModerationResponse struct {
	Toxic    *float64 `json:"toxic"`
	Explicit *float64 `json:"explicit"`
	Spam     *float64 `json:"spam"`
}

// Explicit gets the explicit response from the AIModeration service
type Explicit struct {
	Flag  *float64 `json:"flag"`
	Block *float64 `json:"block"`
}

// ModerationResult is the response from the AIModeration service or the moderation by the user
type ModerationResult struct {
	MessageID            *string               `json:"messageID"`
	Action               *string               `json:"action"`
	ModeratedBy          *string               `json:"moderatedBy"`
	BlockedWord          *string               `json:"blockedWord"`
	BlocklistName        *string               `json:"blocklistName"`
	ModerationThresholds *ModerationThresholds `json:"moderationThresholds"`
	AiModerationResponse *AIModerationResponse `json:"AIModerationResponse"`
	UserKarma            *float64              `json:"userKarma"`
	UserBadKarma         *bool                 `json:"userBadKarma"`
	CreatedAt            *time.Time            `json:"createdAt"`
	UpdatedAt            *time.Time            `json:"updatedAt"`
}

// ModerationThresholds gets the moderation thresholds from the moderation result
type ModerationThresholds struct {
	Explicit *Explicit `json:"explicit"`
	Spam     *Spam     `json:"spam"`
	Toxic    *Toxic    `json:"toxic"`
}

// Spam gets the spam response from the AIModeration service or the moderation by the user
type Spam struct {
	Flag  *float64 `json:"flag"`
	Block *float64 `json:"block"`
}

// Toxic gets the toxic response from the AIModeration service or the moderation by the user
type Toxic struct {
	Flag  *float64 `json:"flag"`
	Block *float64 `json:"block"`
}

// Attachment is the attachment payload for a message
type Attachment struct {
	Type        *string `json:"type"`
	AuthorName  *string `json:"authorName"`
	Title       *string `json:"title"`
	TitleLink   *string `json:"titleLink"`
	Text        *string `json:"text"`
	ImageURL    *string `json:"imageUrl"`
	ThumbURL    *string `json:"thumbUrl"`
	AssetURL    *string `json:"assetUrl"`
	OgScrapeURL *string `json:"ogScrapeUrl"`
}

// MessageFlag is the payload for a message flag
type MessageFlag struct {
	CreatedByAutomod *bool             `json:"createdByAutomod"`
	ModerationResult *ModerationResult `json:"moderationResult"`
	Message          *GetstreamMessage `json:"message"`
	User             *Member           `json:"user"`
	CreatedAt        *time.Time        `json:"createdAt"`
	UpdatedAt        *time.Time        `json:"updatedAt"`
	ReviewedAt       *time.Time        `json:"reviewedAt"`
	ReviewedBy       Member            `json:"reviewedBy"`
	ApprovedAt       *time.Time        `json:"approvedAt"`
	RejectedAt       *time.Time        `json:"rejectedAt"`
}

// Reaction is the payload for a message reaction
type Reaction struct {
	MessageID *string `json:"messageID"`
	UserID    *string `json:"userID"`
	Type      *string `json:"type"`
}

// GetstreamMessage is the payload for a message
type GetstreamMessage struct {
	ID string `json:"id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type            enums.MessageType `json:"type,omitempty"` // one of MessageType* constants
	Silent          bool              `json:"silent,omitempty"`
	User            *Member           `json:"user"`
	Attachments     []*Attachment     `json:"attachments"`
	LatestReactions []*Reaction       `json:"latestReactions"` // last reactions
	OwnReactions    []*Reaction       `json:"ownReactions"`
	ReactionCounts  map[string]int    `json:"reactionCounts"`

	ParentID      string `json:"parentID"`      // id of parent message if it's reply
	ShowInChannel bool   `json:"showInChannel"` // show reply message also in channel

	ReplyCount int `json:"replyCount,omitempty"`

	MentionedUsers []*Member `json:"mentionedUsers"`

	Shadowed bool       `json:"shadowed,omitempty"`
	PinnedAt *time.Time `json:"pinnedAt,omitempty"`
	PinnedBy *Member    `json:"pinnedBy,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`

	ExtraData map[string]interface{} `json:"-"`
}
