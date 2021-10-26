package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *OnboardingDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
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

// SavePin does the actual saving of the users PIN in the database
func (d *OnboardingDb) SavePin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	if pinData.UserID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	pinObj := &gorm.PINData{
		Base:      gorm.Base{},
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		IsValid:   pinData.IsValid,
		Flavour:   pinData.Flavour,
		Salt:      pinData.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to set user pin: %v", err)
	}

	return true, nil
}

// CollectMetrics is responsible for cretating a representation of metrics data.
func (d *OnboardingDb) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {

	if metric.Type == "" {
		return nil, fmt.Errorf("metric type must be specified")
	}

	metricObj := &gorm.Metric{
		Type:      metric.Type,
		Payload:   metric.Payload,
		Timestamp: metric.Timestamp,
		UID:       metric.UID,
	}

	metricSession, err := d.create.CollectMetrics(ctx, metricObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapMetricObjectToDomain(metricSession), nil
}

// RegisterStaffUser creates both the user profile and the staff profile.
func (d *OnboardingDb) RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	if staff.DefaultFacilityID == nil {
		return nil, fmt.Errorf("expected default facility ID to be provided")
	}

	defaultFacilityID := staff.DefaultFacilityID

	facilities := []*gorm.Facility{}

	facilitiesInput := []string{}

	// append default facility at index 0 in the slice of facilitiesInput
	facilitiesInput = append(facilitiesInput, *defaultFacilityID)

	// Append other facilities in the facilitiesInput
	for _, staffFacility := range staff.Facilities {
		facilitiesInput = append(facilitiesInput, *staffFacility)
	}

	// ensure we don't assign duplicate facilities
	sanitizedFacilities := uniqueSliceOfString(facilitiesInput)

	for _, f := range sanitizedFacilities {
		facility, err := d.query.RetrieveFacility(ctx, &f, true)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve facility with id %v: %v", f, err)
		}
		facilities = append(facilities, facility)
	}

	userObject := createUserObject(user)

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

	roles := []string{}
	for _, r := range staff.Roles {
		if !r.IsValid() {
			return nil, fmt.Errorf("role %s is not valid", r)
		}

		roles = append(roles, r.String())
	}

	staffObject := &gorm.StaffProfile{
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
		Addresses:         addresses,
		Roles:             roles,
		Facilities:        facilities,
	}

	userStaffProfile, err := d.create.RegisterStaffUser(ctx, userObject, staffObject)
	if err != nil {

		return nil, fmt.Errorf("failed to create user session: %v", err)
	}

	return d.mapRegisterStaffObjectToDomain(userStaffProfile), nil

}

// AddIdentifier is responsible for creating an identifier and associating it with a specific client
func (d *OnboardingDb) AddIdentifier(
	ctx context.Context,
	clientID string,
	idType enums.IdentifierType,
	idValue string,
	isPrimary bool,
) (*domain.Identifier, error) {
	identifierPayload := &gorm.Identifier{
		ClientID:            clientID,
		IdentifierType:      idType,
		IdentifierValue:     idValue,
		IdentifierUse:       enums.IdentifierUseOfficial,
		IsPrimaryIdentifier: isPrimary,
		Active:              true,
	}

	identifier, err := d.create.AddIdentifier(ctx, identifierPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create identifier: %v", err)
	}

	return d.mapIdentifierObjectToDomain(identifier), nil
}

// RegisterClient is responsible for registering and saving the client's data to the database
func (d *OnboardingDb) RegisterClient(
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

// A helper function that ensures sanitizes the facility ID inputs for staff
func uniqueSliceOfString(slice []string) []string {
	var unique []string
sampleLoop:
	for _, v := range slice {
		for i, u := range unique {
			if v == u {
				unique[i] = v
				continue sampleLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}
