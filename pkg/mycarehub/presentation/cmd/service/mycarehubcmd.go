package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"gorm.io/gorm"
)

type MyCareHubCmdInterfaces interface {
	CreateSuperUser(ctx context.Context, stdin io.Reader) error
	LoadFacilities(ctx context.Context, absoluteFilePath string) error
	LoadOrganisatioAndProgram(ctx context.Context, organisationPath, programPath string) error
}

// MyCareHubCmdInterfacesImpl represents the usecase implementation object
type MyCareHubCmdInterfacesImpl struct {
	usecase usecases.MyCareHub
}

// NewMyCareHubCmdInterfaces initializes a new rest handlers usecase
func NewMyCareHubCmdInterfaces(usecase usecases.MyCareHub) MyCareHubCmdInterfaces {
	return &MyCareHubCmdInterfacesImpl{usecase}
}

// CreateSuperUser is a method that enables the creation of the initial user of the platform
func (m *MyCareHubCmdInterfacesImpl) CreateSuperUser(ctx context.Context, stdin io.Reader) error {
	superuserExists, err := m.usecase.User.CheckSuperUserExists(ctx)
	if err != nil {
		return err
	}
	if superuserExists {
		err := fmt.Errorf("superuser already exists")
		return err
	}

	program, err := m.usecase.Programs.GetProgramByNameAndOrgName(ctx, "Mycarehub", "Savannah Global Health Institute")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("program not found: try running \"mycarehub loadorganisationandprogram\"")
		}
		return err
	}

	facilities, err := m.usecase.Programs.GetProgramFacilities(ctx, program.ID)
	if err != nil {
		return err
	}
	if len(facilities) < 1 {
		err = fmt.Errorf("facilities not found: try running \"mycarehub loadfacilities\"")
		return err
	}

	var registrationInput dto.StaffRegistrationInput
	reader := bufio.NewReader(stdin)

	fmt.Println("Create superuser")
	print("Username: ")
	var usernameInput dto.CMDUsernameInput
	usernameInput.Username, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.Username, err = usernameInput.ParseUsername()
	if err != nil {
		return err
	}

	print("Staff Name:\n")
	var name dto.CMDNameInput
	print("\tFirst Name: ")
	name.FirstName, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	print("\tLast Name: ")
	name.LastName, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.StaffName, err = name.ParseName()
	if err != nil {
		return err
	}

	print("Date of Birth:\n")
	var d dto.CMDDateInput
	print("\tYear: ")
	d.Year, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	print("\tMonth: ")
	d.Month, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	print("\tDay: ")
	d.Day, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.DateOfBirth, err = d.ParseDate()
	if err != nil {
		return err
	}

	print("Gender: ")
	var genderInput dto.CMDGenderInput
	genderInput.Gender, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.Gender, err = genderInput.ParseGender()
	if err != nil {
		return err
	}

	print("Phone Number: ")
	var phoneInput dto.CMDPhoneInput
	phoneInput.Phone, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.PhoneNumber, err = phoneInput.ParsePhone()
	if err != nil {
		return err
	}

	print("Send Invite to This Phone (yes/no): ")
	var sendInviteInput dto.CMDSendInviteInput
	sendInviteInput.SendInvite, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.InviteStaff, err = sendInviteInput.ParseSendInvite()
	if err != nil {
		return err
	}

	print("ID Number: ")
	var idNumberInput dto.CMDIDNumberInput
	idNumberInput.IDNumber, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.IDNumber, err = idNumberInput.ParseIDNumber()
	if err != nil {
		return err
	}

	print("Staff Number: ")
	var staffNumberInput dto.CMDStaffInput
	staffNumberInput.StaffNumber, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	registrationInput.StaffNumber, err = staffNumberInput.ParseStaffNumber()
	if err != nil {
		return err
	}

	registrationInput.ProgramID = program.ID
	registrationInput.OrganisationID = program.Organisation.ID
	registrationInput.Facility = facilities[0].Identifier.Value

	_, err = m.usecase.User.CreateSuperUser(ctx, registrationInput)
	if err != nil {
		return err
	}
	fmt.Println("Successfully created superuser")
	return nil
}

// LoadFacilities reads the facilities file and saves the initial facilities to the database
func (m *MyCareHubCmdInterfacesImpl) LoadFacilities(ctx context.Context, path string) error {
	fmt.Println("Loading Facilities...")

	program, err := m.usecase.Programs.GetProgramByNameAndOrgName(ctx, "Mycarehub", "Savannah Global Health Institute")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("program not found: try running \"mycarehub loadorganisationandprogram\"")
		}
		return err
	}

	facilities, err := utils.ParseFacilitiesFromCSV(path)
	if err != nil {
		return err
	}

	facilities, err = m.usecase.Facility.CreateFacilities(ctx, facilities)
	if err != nil {
		return err
	}

	var facilityIDs []string

	for _, facility := range facilities {
		facilityIDs = append(facilityIDs, *facility.ID)
	}

	_, err = m.usecase.Facility.CmdAddFacilityToProgram(ctx, facilityIDs, program.ID)
	if err != nil {
		return err
	}

	err = m.usecase.Facility.PublishFacilitiesToCMS(ctx, facilities)
	if err != nil {
		return err
	}

	fmt.Println("Successfully loaded facilities")

	return nil
}

// LoadOrganisatioAndProgram reads the organisation and program json files and saves the details in the database
func (m *MyCareHubCmdInterfacesImpl) LoadOrganisatioAndProgram(ctx context.Context, organisationPath, programPath string) error {
	fmt.Println("Loading organisation and program...")

	organisationOutput := dto.OrganisationOutput{}
	orgBs, err := utils.ReadFile(organisationPath)
	if err != nil {
		return err
	}
	organisationInput, err := organisationOutput.ParseValues(orgBs)
	if err != nil {
		return err
	}

	programOutput := dto.ProgramJsonOutput{}
	progBs, err := utils.ReadFile(programPath)
	if err != nil {
		return err
	}
	programInput, err := programOutput.ParseValues(progBs)
	if err != nil {
		return err
	}

	_, err = m.usecase.Organisation.CreateOrganisation(ctx, *organisationInput, []*dto.ProgramInput{programInput})
	if err != nil {
		return err
	}

	fmt.Println("Successfully loaded organisation and program")

	return nil

}
