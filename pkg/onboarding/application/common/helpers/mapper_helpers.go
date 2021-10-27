package helpers

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// CreateMapUser is a helper method that combines staff and user objects and maps it to a domain object
// that is used in the output when displaying data in the graph output
func CreateMapUser(userObject *gorm.User) *domain.User {
	contacts := []*domain.Contact{}
	if len(userObject.Contacts) > 0 {
		for _, u := range userObject.Contacts {
			contact := &domain.Contact{
				ID:      u.ContactID,
				Type:    u.Type,
				Contact: u.Contact,
				Active:  u.Active,
				OptedIn: u.OptedIn,
			}
			contacts = append(contacts, contact)
		}
	}

	languages := []enumutils.Language{}
	for _, l := range userObject.Languages {
		languages = append(languages, enumutils.Language(l))

	}

	user := &domain.User{
		ID:                  userObject.UserID,
		Username:            userObject.Username,
		DisplayName:         userObject.DisplayName,
		FirstName:           userObject.FirstName,
		MiddleName:          userObject.MiddleName,
		LastName:            userObject.LastName,
		Gender:              userObject.Gender,
		UserType:            userObject.UserType,
		Contacts:            contacts,
		Languages:           languages,
		Active:              userObject.Active,
		LastSuccessfulLogin: userObject.LastSuccessfulLogin,
		LastFailedLogin:     userObject.LastFailedLogin,
		FailedLoginCount:    userObject.FailedLoginCount,
		NextAllowedLogin:    userObject.NextAllowedLogin,
		TermsAccepted:       userObject.TermsAccepted,
		AcceptedTermsID:     userObject.AcceptedTermsID,
		Flavour:             userObject.Flavour,
	}
	return user
}

// CreateMapStaff is a helper that maps the staff gorm object to domain object
func CreateMapStaff(staffObject *gorm.StaffProfile) *domain.StaffProfile {

	addresses := []*domain.Addresses{}
	for _, a := range staffObject.Addresses {
		address := &domain.Addresses{
			Type:       a.Type,
			Text:       a.Text,
			Country:    a.Country,
			PostalCode: a.PostalCode,
			County:     a.County,
			Active:     a.Active,
		}
		addresses = append(addresses, address)
	}

	staffProfile := &domain.StaffProfile{
		ID:          staffObject.StaffProfileID,
		UserID:      staffObject.UserID,
		StaffNumber: staffObject.StaffNumber,
		// Facilities:        staffObject.Facilities,
		DefaultFacilityID: staffObject.DefaultFacilityID,
		// Roles:             staffObject.Roles,
		Addresses: addresses,
	}
	return staffProfile
}
