package postgres

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// a helper method to create mapped user
func createMapUser(userObject *gorm.User) *domain.User {
	contact := &domain.Contact{
		ID:           userObject.Contacts.ContactID,
		ContactType:  userObject.Contacts.ContactType,
		ContactValue: userObject.Contacts.ContactValue,
		Active:       userObject.Contacts.Active,
		OptedIn:      userObject.Contacts.OptedIn,
	}

	user := &domain.User{
		ID:                     userObject.UserID,
		Username:               userObject.Username,
		FirstName:              userObject.FirstName,
		MiddleName:             userObject.MiddleName,
		LastName:               userObject.LastName,
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
	}
	return user
}

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *MyCareHubDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility) *domain.Facility {
	if facilityObject == nil {
		return nil
	}

	return &domain.Facility{
		ID:          facilityObject.FacilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      facilityObject.Active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}

// mapProfileObjectToDomain maps the db user profile to a domain model.
func (d *MyCareHubDb) mapProfileObjectToDomain(profileObject *gorm.User) *domain.User {
	if profileObject == nil {
		return nil
	}
	return createMapUser(profileObject)
}

// mapPINObjectToDomain maps the db pin data to a domain model.
func (d *MyCareHubDb) mapPINObjectToDomain(pinObj *gorm.PINData) *domain.UserPIN {
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

// mapOTPObjectToDomain maps the db otp data to a domain model.
func (d *MyCareHubDb) mapOTPObjectToDomain(otpObj *gorm.UserOTP) *domain.OTP {
	if otpObj == nil {
		return nil
	}

	return &domain.OTP{
		UserID:      otpObj.UserID,
		OTP:         otpObj.OTP,
		GeneratedAt: otpObj.GeneratedAt,
		ValidUntil:  otpObj.ValidUntil,
		Flavour:     otpObj.Flavour,
		Valid:       otpObj.Valid,
	}
}

// mapContactObjectToDomain maps the db contact data to a domain model.
func (d *MyCareHubDb) mapContactObjectToDomain(contactObj *gorm.Contact) *domain.Contact {
	if contactObj == nil {
		return nil
	}

	return &domain.Contact{
		ID:           contactObj.ContactID,
		ContactType:  contactObj.ContactType,
		ContactValue: contactObj.ContactValue,
		Active:       contactObj.Active,
		OptedIn:      contactObj.OptedIn,
	}
}
