package postgres

import (
	"context"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// UpdateUserLastSuccessfulLogin update the user with the last login time
func (d *OnboardingDb) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
	return d.update.UpdateUserLastSuccessfulLogin(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin updates the user's last failed login time
func (d *OnboardingDb) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
	return d.update.UpdateUserLastFailedLogin(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount updates the user's failed login count.
func (d *OnboardingDb) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error {
	return d.update.UpdateUserFailedLoginCount(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin updates the user's next allowed login time
func (d *OnboardingDb) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
	return d.update.UpdateUserNextAllowedLogin(ctx, userID, nextAllowedLoginTime, flavour)
}

// UpdateStaffUserProfile updates the staff user details
func (d *OnboardingDb) UpdateStaffUserProfile(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error) {
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

	addresses := []*gorm.Addresses{}
	if len(staff.Addresses) > 0 {
		for _, a := range staff.Addresses {
			address := gorm.Addresses{
				Type:       a.Type,
				Text:       a.Text,
				Country:    a.Country,
				PostalCode: a.PostalCode,
				County:     a.County,
				Active:     a.Active,
			}
			addresses = append(addresses, &address)
		}
	}

	languages := []string{}
	for _, l := range user.Languages {
		languages = append(languages, l.String())
	}

	userData := &gorm.User{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		UserType:    user.UserType,
		Gender:      user.Gender,
		Contacts:    contacts,
		Languages:   languages,
		Flavour:     user.Flavour,
	}

	staffData := &gorm.StaffProfile{
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
		Facility:          gorm.Facility{},
		Addresses:         addresses,
	}

	return d.update.UpdateStaffUserProfile(ctx, userID, userData, staffData)

}
