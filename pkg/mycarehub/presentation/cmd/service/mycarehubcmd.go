package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
)

type MyCareHubCmdInterfaces interface {
	CreateSuperUser(ctx context.Context, stdin io.Reader) error
	LoadOrganisation(ctx context.Context, organisationPath string) error
	LoadProgram(ctx context.Context, programPath string, stdin io.Reader) error
	LoadFacilities(ctx context.Context, absoluteFilePath string) error
	LinkFacilityToProgram(ctx context.Context, stdin io.Reader) error
	LoadSecurityQuestions(ctx context.Context, absoluteFilePath string) error
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

	reader := bufio.NewReader(stdin)

	organisation, err := m.SelectOrganisation(ctx, reader)
	if err != nil {
		return err
	}

	program, err := m.SelectProgram(ctx, organisation.ID, reader)
	if err != nil {
		return err
	}

	facilities, err := m.usecase.Programs.GetProgramFacilities(ctx, program.ID)
	if err != nil {
		return err
	}
	if len(facilities) < 1 {
		err = fmt.Errorf(`facilities not found: try running the following commands in order
		"mycarehub loadfacilities"
		"mycarehub linkfacilitytoprogram"
		`)
		return err
	}

	var registrationInput dto.StaffRegistrationInput

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

	print("Gender (male, female, other): ")
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

	facilities, err := utils.ParseFacilitiesFromCSV(path)
	if err != nil {
		return err
	}

	facilities, err = m.usecase.Facility.CreateFacilities(ctx, facilities)
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

// LoadOrganisation reads the organisation json files and saves the details in the database
func (m *MyCareHubCmdInterfacesImpl) LoadOrganisation(ctx context.Context, organisationPath string) error {
	fmt.Println("Loading organisation...")

	organisationOutput := dto.OrganisationOutput{}
	orgBs, err := utils.ReadFile(organisationPath)
	if err != nil {
		return err
	}
	organisationInput, err := organisationOutput.ParseValues(orgBs)
	if err != nil {
		return err
	}

	_, err = m.usecase.Organisation.CreateOrganisation(ctx, *organisationInput, nil)
	if err != nil {
		return err
	}

	fmt.Println("Successfully loaded organisation")

	return nil
}

// LoadProgram reads the program json files and saves the details in the database
func (m *MyCareHubCmdInterfacesImpl) LoadProgram(ctx context.Context, programPath string, stdin io.Reader) error {
	fmt.Println("Loading program...")

	reader := bufio.NewReader(stdin)

	organisation, err := m.SelectOrganisation(ctx, reader)
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

	programInput.OrganisationID = organisation.ID

	_, err = m.usecase.Programs.CreateProgram(ctx, programInput)
	if err != nil {
		return err
	}

	fmt.Println("Successfully loaded program")

	return nil
}

// LinkFacilityToProgram enables a user to select facilities and link it to a program
func (m *MyCareHubCmdInterfacesImpl) LinkFacilityToProgram(ctx context.Context, stdin io.Reader) error {
	println("Linking Facility to Program...")

	reader := bufio.NewReader(stdin)

	organisation, err := m.SelectOrganisation(ctx, reader)
	if err != nil {
		return err
	}

	program, err := m.SelectProgram(ctx, organisation.ID, reader)
	if err != nil {
		return err
	}

	facility, err := m.SelectFacility(ctx, reader)
	if err != nil {
		return err
	}

	_, err = m.usecase.Facility.CmdAddFacilityToProgram(ctx, []string{*facility.ID}, program.ID)
	if err != nil {
		return err
	}

	fmt.Println("Successfully linked facility to program")

	return nil
}

// SelectOrganisation enables a user to select an organisation
func (m *MyCareHubCmdInterfacesImpl) SelectOrganisation(ctx context.Context, reader *bufio.Reader) (*domain.Organisation, error) {
	print("Organisations: ")
	organisationsPage, err := m.usecase.Organisation.ListOrganisations(ctx, &dto.PaginationsInput{Limit: 2, CurrentPage: 1})
	if err != nil {
		return nil, err
	}
	if organisationsPage == nil {
		err := fmt.Errorf(`organisation not found: try running the following command
		"mycarehub loadorganisation"
		`)
		return nil, err
	}
	for i, v := range organisationsPage.Organisations {
		fmt.Printf("\t%v: %v\n", i, v.Name)
	}

	print("Select Organisation: ")
	var organisationIndexInput string
	organisationIndexInput, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	organisation, err := utils.ParseChoice(organisationsPage.Organisations, organisationIndexInput)
	if err != nil {
		return nil, err
	}

	return *organisation, nil
}

// SelectProgram enables a user to select a program that belongs to an organisation
func (m *MyCareHubCmdInterfacesImpl) SelectProgram(ctx context.Context, organisationID string, reader *bufio.Reader) (*domain.Program, error) {
	print("Programs: ")
	programsPage, err := m.usecase.Programs.ListOrganisationPrograms(ctx, organisationID, &dto.PaginationsInput{Limit: 2, CurrentPage: 1})
	if err != nil {
		return nil, err
	}
	if programsPage == nil {
		err := fmt.Errorf(`program not found: try running the following commands
		"mycarehub loadprogram
		"mycarehub loadorganisation"
		`)
		return nil, err
	}
	for i, v := range programsPage.Programs {
		fmt.Printf("\t%v: %v\n", i, v.Name)
	}

	print("Select Program: ")
	var programIndexInput string
	programIndexInput, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	program, err := utils.ParseChoice(programsPage.Programs, programIndexInput)
	if err != nil {
		return nil, err
	}

	return *program, nil
}

// SelectFacility enables a user to select a facility
func (m *MyCareHubCmdInterfacesImpl) SelectFacility(ctx context.Context, reader *bufio.Reader) (*domain.Facility, error) {
	print("Facilities: ")
	facilitiesPage, err := m.usecase.Facility.ListFacilities(ctx, nil, nil, &dto.PaginationsInput{Limit: 2, CurrentPage: 1})
	if err != nil {
		return nil, err
	}
	if facilitiesPage == nil {
		err := fmt.Errorf(`facilities not found: try running the following command
		"mycarehub loadfacilities"
		`)
		return nil, err
	}
	for i, v := range facilitiesPage.Facilities {
		fmt.Printf("\t%v: %v\n", i, v.Name)
	}

	print("Select Facility: ")
	var facilityIndexInput string
	facilityIndexInput, err = reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	facility, err := utils.ParseChoice(facilitiesPage.Facilities, facilityIndexInput)
	if err != nil {
		return nil, err
	}

	return *facility, nil
}

// LoadSecurityQuestions enables a user to load security questions
func (m *MyCareHubCmdInterfacesImpl) LoadSecurityQuestions(ctx context.Context, absoluteFilePath string) error {
	println("Loading security questions...")
	securityQuestions := []*domain.SecurityQuestion{}
	bs, err := utils.ReadFile(absoluteFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &securityQuestions)
	if err != nil {
		return err
	}

	for _, securityQuestion := range securityQuestions {
		if valid := securityQuestion.ResponseType.IsValid(); !valid {
			return fmt.Errorf("invalid response type: %v", securityQuestion.ResponseType)
		}
		if valid := securityQuestion.Flavour.IsValid(); !valid {
			return fmt.Errorf("invalid flavour: %v", securityQuestion.Flavour)
		}
	}

	_, err = m.usecase.SecurityQuestions.CreateSecurityQuestions(ctx, securityQuestions)
	if err != nil {
		return err
	}

	println("Successfully loaded security questions")
	return nil
}
