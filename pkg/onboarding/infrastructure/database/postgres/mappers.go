package postgres

import (
	"strconv"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
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

// mapMetricObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapMetricObjectToDomain(metricObject *gorm.Metric) *domain.Metric {
	if metricObject == nil {
		return nil
	}

	return &domain.Metric{
		MetricID:  metricObject.MetricID,
		Type:      metricObject.Type,
		Payload:   metricObject.Payload,
		Timestamp: metricObject.Timestamp,
		UID:       metricObject.UID,
	}
}

// mapProfileObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapProfileObjectToDomain(profileObject *gorm.User) *domain.User {
	if profileObject == nil {
		return nil
	}

	return &domain.User{
		ID:          profileObject.UserID,
		Username:    profileObject.Username,
		DisplayName: profileObject.DisplayName,
		FirstName:   profileObject.FirstName,
		MiddleName:  profileObject.MiddleName,
		LastName:    profileObject.LastName,
		Flavour:     profileObject.Flavour,
		// UserType:            profileObject.UserType,
		// Gender:              profileObject.Gender,
		Active: profileObject.Active,
		// Languages:           profileObject.Languages,
		// PushTokens:          profileObject.PushTokens,
		LastSuccessfulLogin: profileObject.LastSuccessfulLogin,
		LastFailedLogin:     profileObject.LastFailedLogin,
		FailedLoginCount:    profileObject.FailedLoginCount,
		NextAllowedLogin:    profileObject.NextAllowedLogin,
		TermsAccepted:       profileObject.TermsAccepted,
		AcceptedTermsID:     profileObject.AcceptedTermsID,
	}
}

// mapPINObjectToDomain maps the db pin data to a domain model.
// It searches the database to fetch items specific to the PIN
func (d *OnboardingDb) mapPINObjectToDomain(pinObj *gorm.PINData) *domain.UserPIN {
	if pinObj == nil {
		return nil
	}

	return &domain.UserPIN{
		UserID:    pinObj.UserID,
		HashedPIN: pinObj.HashedPIN,
		ValidFrom: pinObj.ValidFrom,
		ValidTo:   pinObj.ValidTo,
		Flavour:   pinObj.Flavour,
		IsValid:   pinObj.IsValid,
		Salt:      pinObj.Salt,
	}
}

// mapUserObjectToDomain maps the db user to a domain model.
// It searches the database to fetch items specific to the user
func (d *OnboardingDb) mapRegisterStaffObjectToDomain(userStaffObject *gorm.StaffUserProfile) *domain.StaffUserProfile {

	userObject := userStaffObject.User
	staffObject := userStaffObject.Staff

	user := createMapUser(userObject)

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
		UserID:      userObject.UserID,
		StaffNumber: staffObject.StaffNumber,
		// Facilities:        staffObject.Facilities,
		DefaultFacilityID: staffObject.DefaultFacilityID,
		// Roles:             staffObject.Roles,
		Addresses: addresses,
	}
	return &domain.StaffUserProfile{
		User:  user,
		Staff: staffProfile,
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

// mapIdentifierObjectToDomain maps the identifier object to our domain defined type
func (d *OnboardingDb) mapIdentifierObjectToDomain(identifierObject *gorm.Identifier) *domain.Identifier {
	if identifierObject == nil {
		return nil
	}

	return &domain.Identifier{
		ID:                  identifierObject.ID,
		ClientID:            identifierObject.ClientID,
		IdentifierType:      identifierObject.IdentifierType,
		IdentifierUse:       identifierObject.IdentifierUse,
		IdentifierValue:     identifierObject.IdentifierValue,
		Description:         identifierObject.Description,
		ValidFrom:           identifierObject.ValidFrom,
		ValidTo:             identifierObject.ValidTo,
		Active:              identifierObject.Active,
		IsPrimaryIdentifier: identifierObject.IsPrimaryIdentifier,
	}
}

// mapClientObjectToDomain maps the client object to the domain defined type
func (d *OnboardingDb) mapClientObjectToDomain(client *gorm.ClientProfile) *domain.ClientProfile {
	if client == nil {
		return nil
	}

	return &domain.ClientProfile{
		ID:                      client.ID,
		UserID:                  client.UserID,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
		ClientType:              client.ClientType,
		Active:                  client.Active,
		HealthRecordID:          client.HealthRecordID,
		// Identifiers:             client.Identifiers,
		FacilityID:           client.FacilityID,
		TreatmentBuddyUserID: client.TreatmentBuddy,
		CHVUserID:            client.CHVUserID,
		ClientCounselled:     client.ClientCounselled,
	}
}
