package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lib/pq"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *MyCareHubDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if facility.Code == "" {
		return nil, fmt.Errorf("`code` should be defined")
	}

	facilityObj := &gorm.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      strconv.FormatBool(facility.Active),
		County:      facility.County,
		Description: facility.Description,
	}

	facilitySession, err := d.create.GetOrCreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// RegisterClient is responsible for registering and saving the client's data to the database
func (d *MyCareHubDb) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	if clientInput == nil {
		return nil, fmt.Errorf("expected client input to be provided")
	}

	if userInput == nil {
		return nil, fmt.Errorf("expected user input to be provided")
	}

	userObject := createUserObject(userInput)

	clientObject := &gorm.ClientProfile{
		ClientType: clientInput.ClientType,
	}

	clientUserProfile, err := d.create.RegisterClient(ctx, userObject, clientObject)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return d.mapRegisterClientObjectToDomain(clientUserProfile), nil
}

// A helper method to create a user object
func createUserObject(user *dto.UserInput) *gorm.User {
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

// SavePin gets the pin details from the user and saves it in the database
func (d *MyCareHubDb) SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error) {
	pinObj := &gorm.PINData{
		UserID:    pinInput.UserID,
		HashedPIN: pinInput.HashedPIN,
		ValidFrom: pinInput.ValidFrom,
		ValidTo:   pinInput.ValidTo,
		IsValid:   pinInput.IsValid,
		Flavour:   pinInput.Flavour,
		Salt:      pinInput.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}
