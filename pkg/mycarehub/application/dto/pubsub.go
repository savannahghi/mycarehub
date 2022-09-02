package dto

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
)

// PubSubCMSClientInput is the subscribers input to make an api call to the cms service
type PubSubCMSClientInput struct {
	// user details
	UserID      string           `json:"user_id"`
	Name        string           `json:"name"`
	Gender      enumutils.Gender `json:"gender"`
	UserType    enums.UsersType  `json:"user_type"`
	PhoneNumber string           `json:"phone_number"`
	Handle      string           `json:"handle"`
	Flavour     feedlib.Flavour  `json:"flavour"`
	DateOfBirth scalarutils.Date `json:"date_of_birth"`

	// client details
	ClientID       string             `json:"client_id"`
	ClientTypes    []enums.ClientType `json:"client_types"`
	EnrollmentDate scalarutils.Date   `json:"enrollment_date"`

	// facility details
	FacilityID   string `json:"facility_id"`
	FacilityName string `json:"facility_name"`

	// organisation details
	OrganisationID string `json:"organisation_id"`
}
