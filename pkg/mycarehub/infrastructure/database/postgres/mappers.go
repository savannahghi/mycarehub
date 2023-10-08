package postgres

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// a helper method to create mapped user
func createMapUser(userObject *gorm.User) *domain.User {
	contact := &domain.Contact{
		ID:           &userObject.Contacts.ID,
		ContactType:  userObject.Contacts.Type,
		ContactValue: userObject.Contacts.Value,
		Active:       userObject.Contacts.Active,
		OptedIn:      userObject.Contacts.OptedIn,
	}

	user := &domain.User{
		ID:                     userObject.UserID,
		Username:               userObject.Username,
		Email:                  userObject.Email,
		Name:                   userObject.Name,
		Gender:                 userObject.Gender,
		Contacts:               contact,
		Active:                 userObject.Active,
		LastSuccessfulLogin:    userObject.LastSuccessfulLogin,
		LastFailedLogin:        userObject.LastFailedLogin,
		FailedLoginCount:       userObject.FailedLoginCount,
		NextAllowedLogin:       userObject.NextAllowedLogin,
		TermsAccepted:          userObject.TermsAccepted,
		PinChangeRequired:      userObject.PinChangeRequired,
		HasSetPin:              userObject.HasSetPin,
		HasSetSecurityQuestion: userObject.HasSetSecurityQuestion,
		IsPhoneVerified:        userObject.IsPhoneVerified,
		DateOfBirth:            userObject.DateOfBirth,
		FailedSecurityCount:    userObject.FailedSecurityCount,
		PinUpdateRequired:      userObject.PinUpdateRequired,
		PushTokens:             userObject.PushTokens,
		CurrentOrganizationID:  userObject.CurrentOrganisationID,
		CurrentProgramID:       userObject.CurrentProgramID,
		HasSetNickname:         userObject.HasSetUsername,
		IsSuperuser:            userObject.IsSuperuser,
	}
	return user
}

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *MyCareHubDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility, identifierObject []*gorm.FacilityIdentifier) *domain.Facility {
	if facilityObject == nil {
		return nil
	}
	if identifierObject == nil {
		return nil
	}

	var ids []*domain.FacilityIdentifier

	for _, identifier := range identifierObject {
		ids = append(ids, &domain.FacilityIdentifier{
			ID:     identifier.ID,
			Active: identifier.Active,
			Type:   enums.FacilityIdentifierType(identifier.Type),
			Value:  identifier.Value,
		})
	}

	return &domain.Facility{
		ID:                 facilityObject.FacilityID,
		Name:               facilityObject.Name,
		Phone:              facilityObject.Phone,
		Active:             facilityObject.Active,
		Country:            facilityObject.Country,
		Description:        facilityObject.Description,
		FHIROrganisationID: facilityObject.FHIROrganisationID,
		Identifiers:        ids,
		WorkStationDetails: domain.WorkStationDetails{},
		Coordinates: &domain.Coordinates{
			Lat: facilityObject.Latitude,
			Lng: facilityObject.Longitude,
		},
	}
}

// mapProfileObjectToDomain maps the db user profile to a domain model.
func (d *MyCareHubDb) mapProfileObjectToDomain(profileObject *gorm.User) *domain.User {
	if profileObject == nil {
		return nil
	}
	return createMapUser(profileObject)
}
