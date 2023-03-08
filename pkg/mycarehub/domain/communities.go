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
