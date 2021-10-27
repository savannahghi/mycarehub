package staff_test

import (
	"context"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestUseCaseStaffProfileImpl_RegisterStaffUser(t *testing.T) {

	f := testInfrastructureInteractor
	ctx := context.Background()

	testFacilityID := uuid.New().String()
	i := testInteractor

	code := ksuid.New().String()
	facilityInput := dto.FacilityInput{
		Name:        "test",
		Code:        code,
		Active:      true,
		County:      "test",
		Description: "test description",
	}

	//valid: Create a facility
	facility, err := f.GetOrCreateFacility(ctx, facilityInput)
	assert.Nil(t, err)
	assert.NotNil(t, facility)

	// First Set of Valid Input
	contactInput := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: randomdata.PhoneNumber(),
		Active:  true,
		OptedIn: true,
	}

	userInput := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Gender:      enumutils.GenderMale,
		UserType:    enums.HealthcareWorkerUser,
		Contacts:    []*dto.ContactInput{contactInput},
		Languages:   []enumutils.Language{enumutils.LanguageEn},
		Flavour:     feedlib.FlavourPro,
	}

	staffID := ksuid.New().String()
	staffInput := &dto.StaffProfileInput{
		StaffNumber:       staffID,
		DefaultFacilityID: facility.ID,
		Addresses: []*dto.AddressesInput{
			{
				Type:       enums.AddressesTypePhysical,
				Text:       "test",
				Country:    enums.CountryTypeKenya,
				PostalCode: "test code",
				County:     enums.CountyTypeNakuru,
				Active:     true,
			},
		},
		Roles: []enums.RolesType{enums.RolesTypeCanInviteClient},
	}

	// Second set of valid Inputs
	contactInput2 := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: randomdata.PhoneNumber(),
		Active:  true,
		OptedIn: true,
	}

	userInput2 := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Gender:      enumutils.GenderMale,
		UserType:    enums.HealthcareWorkerUser,
		Contacts:    []*dto.ContactInput{contactInput2},
		Languages:   []enumutils.Language{enumutils.LanguageEn},
		Flavour:     feedlib.FlavourPro,
	}

	staffID2 := ksuid.New().String()
	staffInpu2 := &dto.StaffProfileInput{
		StaffNumber:       staffID2,
		DefaultFacilityID: facility.ID,
		Addresses: []*dto.AddressesInput{
			{
				Type:       enums.AddressesTypePhysical,
				Text:       "test",
				Country:    enums.CountryTypeKenya,
				PostalCode: "test code",
				County:     enums.CountyTypeBaringo,
				Active:     true,
			},
		},
		Roles: []enums.RolesType{enums.RolesTypeCanInviteClient},
	}

	// Invalid facility id
	staffInputNoFacility := &dto.StaffProfileInput{
		StaffNumber:       ksuid.New().String(),
		DefaultFacilityID: &testFacilityID,
	}

	// Invalid country input
	staffInputInvalidCountry := &dto.StaffProfileInput{
		StaffNumber:       staffID2,
		DefaultFacilityID: facility.ID,
		Addresses: []*dto.AddressesInput{
			{
				Type:       enums.AddressesTypePhysical,
				Text:       "test",
				Country:    "Invalid",
				PostalCode: "test code",
				County:     enums.CountyTypeBaringo,
				Active:     true,
			},
		},
	}

	// Invalid role
	staffInputInvalidRole := &dto.StaffProfileInput{
		StaffNumber:       ksuid.New().String(),
		DefaultFacilityID: &testFacilityID,
		Roles:             []enums.RolesType{"invalid"},
	}

	// Invalid county input
	staffInputInvalidCounty := &dto.StaffProfileInput{
		StaffNumber:       staffID2,
		DefaultFacilityID: facility.ID,
		Addresses: []*dto.AddressesInput{
			{
				Type:       enums.AddressesTypePhysical,
				Text:       "test",
				Country:    enums.CountryTypeKenya,
				PostalCode: "test code",
				County:     "Invalid",
				Active:     true,
			},
		},
	}

	//  invalid: non existent facility assignment
	useStaffProfile, err := f.GetOrCreateStaffUser(ctx, userInput, staffInputNoFacility)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//  invalid: non existent country
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInputInvalidCountry)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//  invalid: non existent county
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInputInvalidCounty)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	// TODO:add case where county is valid but does not belong to country after another country is available

	//  invalid: non existent facility assignment
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInputNoFacility)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//  invalid: invalid role provided
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInputInvalidRole)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//valid: create a staff user with valid parameters
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, err)
	assert.NotNil(t, useStaffProfile)

	//Valid: creating a user with duplicate staff number and contact
	useStaffProfile, err = i.StaffUsecase.GetOrCreateStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, err)
	assert.NotNil(t, useStaffProfile)

	//Invalid: creating a user with duplicate staff number (changed contact only)
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput2, staffInput)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//Invalid: creating a user with duplicate Contact (changed staff number only)
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput, staffInpu2)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	// Valid: saves again if the duplicate Contact and Staff number are rectified
	useStaffProfile, err = f.GetOrCreateStaffUser(ctx, userInput2, staffInpu2)
	assert.Nil(t, err)
	assert.NotNil(t, useStaffProfile)

	// TODO: teardown the user and replace randomdata with gofakeit

}

func TestUsecasesStaffProfileImpl_UpdateStaffUser_Integration(t *testing.T) {
	ctx := context.Background()

	f := testInfrastructureInteractor
	i := testInteractor

	testFacilityID := uuid.New().String()

	code := ksuid.New().String()
	facilityInput := dto.FacilityInput{
		Name:        "test",
		Code:        code,
		Active:      true,
		County:      "test",
		Description: "test description",
	}

	//valid: Create a facility
	facility, err := f.GetOrCreateFacility(ctx, facilityInput)
	assert.Nil(t, err)
	assert.NotNil(t, facility)

	// First Set of Valid Input
	contactInput := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: randomdata.PhoneNumber(),
		Active:  true,
		OptedIn: true,
	}

	userInput := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Gender:      enumutils.GenderMale,
		UserType:    enums.HealthcareWorkerUser,
		Contacts:    []*dto.ContactInput{contactInput},
		Languages:   []enumutils.Language{enumutils.LanguageEn},
		Flavour:     feedlib.FlavourPro,
	}

	StaffNumber := ksuid.New().String()
	staffInput := &dto.StaffProfileInput{
		StaffNumber:       StaffNumber,
		DefaultFacilityID: facility.ID,
	}

	// Second set of valid Inputs
	contactInput2 := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: randomdata.PhoneNumber(),
		Active:  true,
		OptedIn: true,
	}

	userInput2 := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Gender:      enumutils.GenderMale,
		UserType:    enums.HealthcareWorkerUser,
		Contacts:    []*dto.ContactInput{contactInput2},
		Languages:   []enumutils.Language{enumutils.LanguageSw},
		Flavour:     feedlib.FlavourPro,
	}

	staffInput2 := &dto.StaffProfileInput{
		StaffNumber:       StaffNumber,
		DefaultFacilityID: facility.ID,
	}

	// Invalid facility id
	staffInputNoFacility := &dto.StaffProfileInput{
		StaffNumber:       ksuid.New().String(),
		DefaultFacilityID: &testFacilityID,
	}

	//valid: create a staff user with valid parameters
	userStaffProfile, err := f.GetOrCreateStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, err)
	assert.NotNil(t, userStaffProfile)

	// Valid userID
	userProfile, err := f.GetUserProfileByUserID(ctx, *userStaffProfile.User.ID, userStaffProfile.User.Flavour)
	assert.Nil(t, err)
	assert.NotNil(t, userProfile)

	staffProfile, err5 := f.GetStaffProfileByStaffNumber(ctx, StaffNumber)
	assert.Nil(t, err5)
	assert.NotNil(t, staffProfile)

	updated, err := i.StaffUsecase.UpdateStaffUserProfile(ctx, *userStaffProfile.User.ID, userInput2, staffInput2)
	assert.Nil(t, err)
	assert.Equal(t, true, updated)

	//Invalid: update user with wrong data
	userProfile2, err := f.GetUserProfileByUserID(ctx, *userStaffProfile.User.ID, userStaffProfile.User.Flavour)
	assert.Nil(t, err)
	assert.NotNil(t, userProfile2)

	staffProfile3, err6 := f.GetStaffProfileByStaffNumber(ctx, "StaffNumber")
	assert.NotNil(t, err6)
	assert.Nil(t, staffProfile3)

	updated2, err := i.StaffUsecase.UpdateStaffUserProfile(ctx, *userStaffProfile.User.ID, userInput2, staffInputNoFacility)
	assert.NotNil(t, err)
	assert.Equal(t, false, updated2)

}
