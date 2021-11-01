package postgres

import (
	"strconv"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *OnboardingDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility) *domain.Facility {
	if facilityObject == nil {
		return nil
	}

	active, err := strconv.ParseBool(facilityObject.Active)
	if err != nil {
		return nil
	}

	return &domain.Facility{
		ID:          facilityObject.FacilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}

// mapRegisterClientObjectToDomain maps the database client object to our custom domain type
func (d *OnboardingDb) mapRegisterClientObjectToDomain(clientObject *gorm.ClientUserProfile) *domain.ClientUserProfile {
	userObject := clientObject.User
	client := clientObject.Client

	user := createMapUser(userObject)

	clientProfile := &domain.ClientProfile{
		ID:         client.ID,
		UserID:     client.UserID,
		ClientType: client.ClientType,
	}

	return &domain.ClientUserProfile{
		User:   user,
		Client: clientProfile,
	}
}

// a helper method to create mapped user
func createMapUser(userObject *gorm.User) *domain.User {
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

// mapProfileObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapProfileObjectToDomain(profileObject *gorm.User) *domain.User {
	if profileObject == nil {
		return nil
	}

	user := createMapUser(profileObject)

	return user
}
