package helpers

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// CreateUserObject is a helper function to  create a user object.
func CreateUserObject(user *dto.UserInput) *gorm.User {
	contacts := []gorm.Contact{}
	if len(user.Contacts) > 0 {
		for _, u := range user.Contacts {
			contact := gorm.Contact{
				Type:    u.Type,
				Contact: u.Contact,
				Active:  u.Active,
				OptedIn: u.OptedIn,
			}
			contacts = append(contacts, contact)
		}
	}

	languages := []string{}
	for _, l := range user.Languages {
		languages = append(languages, l.String())
	}

	userObject := &gorm.User{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		Contacts:    contacts,
		UserType:    user.UserType,
		Languages:   pq.StringArray(languages),
		Flavour:     user.Flavour,
	}
	return userObject
}

// CreateStaffObject is a helper function to  create a staff object.
func CreateStaffObject(staff *dto.StaffProfileInput) (*gorm.StaffProfile, error) {
	addresses := []*gorm.Addresses{}
	if len(staff.Addresses) > 0 {
		for _, a := range staff.Addresses {
			// ensure counties belong to a country
			err := enums.ValidateCountiesOfCountries(enums.CountryType(a.Country), enums.CountyType(a.County))
			if err != nil {
				return nil, fmt.Errorf("failed to validate %v county belongs to %v: %v", a.County, a.Country, err)
			}
			address := &gorm.Addresses{
				Type:       a.Type,
				Text:       a.Text,
				Country:    a.Country,
				PostalCode: a.PostalCode,
				County:     a.County,
				Active:     a.Active,
			}
			addresses = append(addresses, address)
		}
	}

	staffObject := &gorm.StaffProfile{
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
		Addresses:         addresses,
	}
	return staffObject, nil
}
