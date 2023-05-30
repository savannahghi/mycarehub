package domain

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Community defines the payload to create a channel
type Community struct {
	ID          string `json:"id"`
	RoomID      string `json:"roomID"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// The fields below are custom to our implementation
	AgeRange       *AgeRange          `json:"ageRange"`
	Gender         []enumutils.Gender `json:"gender"`
	ClientType     []enums.ClientType `json:"clientType"`
	OrganisationID string             `json:"organisationID"`
	ProgramID      string             `json:"programID"`
	FacilityID     string             `json:"facilityID"`
}

// AgeRange defines the channel users age input
type AgeRange struct {
	LowerBound int `json:"lowerBound"`
	UpperBound int `json:"upperBound"`
}

// MatrixAuth models the Matrix's user authentication data
type MatrixAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MatrixUserRegistration defines the structure of the input to be used when registering a Matrix user
type MatrixUserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// MatrixUserSearchResult defines the structure of the users search output
type MatrixUserSearchResult struct {
	Limited bool     `json:"limited"`
	Results []Result `json:"results"`
}

// Results contain the details of the searched user
type Result struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// PusherPayload models the data class to be used in configuration of Matrix's pusher data
type PusherPayload struct {
	AppDisplayName    string     `json:"app_display_name,omitempty"`
	AppID             string     `json:"app_id,omitempty"`
	Append            bool       `json:"append,omitempty"`
	PusherData        PusherData `json:"data,omitempty"`
	DeviceDisplayName string     `json:"device_display_name,omitempty"`
	Kind              *string    `json:"kind,omitempty"`
	Lang              string     `json:"lang,omitempty"`
	ProfileTag        string     `json:"profile_tag,omitempty"`
	Pushkey           string     `json:"pushkey,omitempty"`
}

// PusherData dictionary of information for the pusher implementation itself
type PusherData struct {
	Format string `json:"format,omitempty"`
	URL    string `json:"url,omitempty"`
}
