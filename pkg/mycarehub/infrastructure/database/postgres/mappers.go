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
		Name:                   userObject.Name,
		Gender:                 userObject.Gender,
		UserType:               userObject.UserType,
		Contacts:               contact,
		Active:                 userObject.Active,
		LastSuccessfulLogin:    userObject.LastSuccessfulLogin,
		LastFailedLogin:        userObject.LastFailedLogin,
		FailedLoginCount:       userObject.FailedLoginCount,
		NextAllowedLogin:       userObject.NextAllowedLogin,
		Flavour:                userObject.Flavour,
		TermsAccepted:          userObject.TermsAccepted,
		PinChangeRequired:      userObject.PinChangeRequired,
		HasSetPin:              userObject.HasSetPin,
		HasSetSecurityQuestion: userObject.HasSetSecurityQuestion,
		IsPhoneVerified:        userObject.IsPhoneVerified,
		DateOfBirth:            userObject.DateOfBirth,
		FailedSecurityCount:    userObject.FailedSecurityCount,
		PinUpdateRequired:      userObject.PinUpdateRequired,
		PushTokens:             userObject.PushTokens,
		OrganizationID:         userObject.OrganisationID,
		CurrentProgramID:       userObject.CurrentProgramID,
	}
	return user
}

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *MyCareHubDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility, identifierObject *gorm.FacilityIdentifier) *domain.Facility {
	if facilityObject == nil {
		return nil
	}
	if identifierObject == nil {
		return nil
	}

	return &domain.Facility{
		ID:                 facilityObject.FacilityID,
		Name:               facilityObject.Name,
		Phone:              facilityObject.Phone,
		Active:             facilityObject.Active,
		County:             facilityObject.Country,
		Description:        facilityObject.Description,
		FHIROrganisationID: facilityObject.FHIROrganisationID,
		Identifier: domain.FacilityIdentifier{
			ID:     identifierObject.ID,
			Active: identifierObject.Active,
			Type:   enums.FacilityIdentifierType(identifierObject.Type),
			Value:  identifierObject.Value,
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
