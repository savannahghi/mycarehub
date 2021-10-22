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
	}

	// Invalid facility id
	staffInputNoFacility := &dto.StaffProfileInput{
		StaffNumber:       ksuid.New().String(),
		DefaultFacilityID: &testFacilityID,
	}

	//valid: create a staff user with valid parameters
	useStaffProfile, err := f.RegisterStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, err)
	assert.NotNil(t, useStaffProfile)

	//Invalid: creating a user with duplicate staff number and contact
	useStaffProfile, err = f.RegisterStaffUser(ctx, userInput, staffInput)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//Invalid: creating a user with duplicate staff number (changed contact only)
	useStaffProfile, err = f.RegisterStaffUser(ctx, userInput2, staffInput)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	//Invalid: creating a user with duplicate Contact (changed staff number only)
	useStaffProfile, err = f.RegisterStaffUser(ctx, userInput, staffInpu2)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	// Valid: saves again if the duplicate Contact and Staff number are rectified
	useStaffProfile, err = f.RegisterStaffUser(ctx, userInput2, staffInpu2)
	assert.Nil(t, err)
	assert.NotNil(t, useStaffProfile)

	//  invalid: non existent facility assignment
	useStaffProfile, err = f.RegisterStaffUser(ctx, userInput, staffInputNoFacility)
	assert.Nil(t, useStaffProfile)
	assert.NotNil(t, err)

	// TODO: teardown the user and replace randomdata with gofakeit

}
