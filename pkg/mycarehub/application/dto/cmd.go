package dto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	validator "gopkg.in/go-playground/validator.v9"
)

var v = validator.New()

// CMDUsernameInput describes the username format to be used on the CLI
type CMDUsernameInput struct {
	Username string `json:"username" validate:"required,alphanum"`
}

// ParseUsername ensures the username entered is alphanumeric and not empty
func (u *CMDUsernameInput) ParseUsername() (string, error) {
	u.Username = strings.TrimSpace(u.Username)

	err := v.Struct(u)
	if err != nil {
		return "", err
	}

	return u.Username, nil
}

// CMDNameInput describes the user's name format to be used on the CLI
type CMDNameInput struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// ParseName ensures that the input is valid
func (s *CMDNameInput) ParseName() (string, error) {
	s.FirstName = strings.TrimSpace(s.FirstName)
	s.LastName = strings.TrimSpace(s.LastName)

	err := v.Struct(s)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%v %v", s.FirstName, s.LastName)

	return name, err
}

// CMDDateInput defines the input of the date to be used in the CLI
type CMDDateInput struct {
	Year  string `json:"year" validate:"required"`
	Month string `json:"month" validate:"required"`
	Day   string `json:"day" validate:"required"`
}

// ParseDate converts the CMD input to a date and also validates it
func (d *CMDDateInput) ParseDate() (scalarutils.Date, error) {
	d.Year = strings.TrimSpace(d.Year)
	d.Month = strings.TrimSpace(d.Month)
	d.Day = strings.TrimSpace(d.Day)

	err := v.Struct(d)
	if err != nil {
		return scalarutils.Date{}, err
	}

	year, err := strconv.Atoi(d.Year)
	if err != nil {
		return scalarutils.Date{}, err
	}
	month, err := strconv.Atoi(d.Month)
	if err != nil {
		return scalarutils.Date{}, err
	}
	day, err := strconv.Atoi(d.Day)
	if err != nil {
		return scalarutils.Date{}, err
	}
	newDate := scalarutils.Date{Year: year, Month: month, Day: day}
	if err := newDate.Validate(); err != nil {
		return scalarutils.Date{}, err
	}

	return newDate, nil
}

// CMDGenderInput defines the input of gender that will be used in CLI
type CMDGenderInput struct {
	Gender string `json:"gender" validate:"required"`
}

// ParseGender converts the CMD input to a gender and also validates it
func (g *CMDGenderInput) ParseGender() (enumutils.Gender, error) {
	g.Gender = strings.TrimSpace(strings.ToLower(g.Gender))

	err := v.Struct(g)
	if err != nil {
		return enumutils.Gender(""), err
	}

	if ok := enumutils.Gender(g.Gender).IsValid(); !ok {
		return enumutils.Gender(""), fmt.Errorf("invalid gender: %v", g.Gender)
	}

	gender := enumutils.Gender(g.Gender)

	return gender, nil
}

// CMDPhoneInput defines the input of phone that will be used in CLI
type CMDPhoneInput struct {
	Phone string `json:"phone" validate:"required"`
}

// ParsePhone converts the CMD input to a phone and also validates it
func (p *CMDPhoneInput) ParsePhone() (string, error) {
	p.Phone = strings.TrimSpace(p.Phone)

	err := v.Struct(p)
	if err != nil {
		return "", err
	}

	_, err = converterandformatter.NormalizeMSISDN(p.Phone)
	if err != nil {
		return "", err
	}

	return p.Phone, nil

}

// CMDSendInviteInput defines the input of send invite that will be used in CLI
type CMDSendInviteInput struct {
	SendInvite string `json:"sendInvite" validate:"required"`
}

// ParseSendInvite converts yes/no responses from the CMD to boolean values and validates them
func (s *CMDSendInviteInput) ParseSendInvite() (bool, error) {
	s.SendInvite = strings.TrimSpace(s.SendInvite)
	s.SendInvite = strings.ToLower(s.SendInvite)

	err := v.Struct(s)
	if err != nil {
		return false, err
	}

	switch s.SendInvite {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input, expected 'yes' or 'no'")
	}
}

// CMDIDNumberInput defines the input of ID number that will be used in CLI
type CMDIDNumberInput struct {
	IDNumber string `json:"idNumber" validate:"required"`
}

// ParseIDNumber converts the CMD input to aN ID number and also validates it
func (i *CMDIDNumberInput) ParseIDNumber() (string, error) {
	i.IDNumber = strings.TrimSpace(i.IDNumber)

	err := v.Struct(i)
	if err != nil {
		return "", err
	}

	return i.IDNumber, nil

}

// CMDStaffInput defines the input of staff number that will be used in CLI
type CMDStaffInput struct {
	StaffNumber string `json:"staffNumber" validate:"required"`
}

// CMDStaffInput converts the CMD input to a staff number and also validates it
func (s *CMDStaffInput) ParseStaffNumber() (string, error) {
	s.StaffNumber = strings.TrimSpace(s.StaffNumber)

	err := v.Struct(s)
	if err != nil {
		return "", err
	}

	return s.StaffNumber, nil

}

// CMDOrganisationInput defines the input of organisation index from the list of organisations that will be used in CLI
type CMDOrganisationInput struct {
	OrganisationIndex string `json:"organisationIndex" validate:"required"`
}

// ParseOrganisation converts the CMD input to an organisation index of a list of organisations and also validates it
func (p *CMDOrganisationInput) ParseOrganisation(organisationLength int) (*int, error) {
	p.OrganisationIndex = strings.TrimSpace(p.OrganisationIndex)

	err := v.Struct(p)
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(p.OrganisationIndex)
	if err != nil {
		return nil, err
	}
	if index > organisationLength-1 {
		err := fmt.Errorf("invalid choice: %v", p)
		return nil, err
	}

	return &index, nil

}

// CMDProgramInput defines the input of program index from the list of programs that will be used in CLI
type CMDProgramInput struct {
	ProgramIndex string `json:"programIndex" validate:"required"`
}

// ParseProgram converts the CMD input to a program index of a list of programs and also validates it
func (p *CMDProgramInput) ParseProgram(programsLength int) (*int, error) {
	p.ProgramIndex = strings.TrimSpace(p.ProgramIndex)

	err := v.Struct(p)
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(p.ProgramIndex)
	if err != nil {
		return nil, err
	}
	if index > programsLength-1 {
		err := fmt.Errorf("invalid choice: %v", p)
		return nil, err
	}

	return &index, nil

}

// CMDFacilityInput defines the input of a facility index of a list of facilities that will be used in CLI
type CMDFacilityInput struct {
	FacilityIndex string `json:"facilityIndex" validate:"required"`
}

// ParsePhone converts the CMD input to a facility index of a list of facilities and also validates it
func (f *CMDFacilityInput) ParseFacility(facilitiesLength int) (*int, error) {
	f.FacilityIndex = strings.TrimSpace(f.FacilityIndex)

	err := v.Struct(f)
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(f.FacilityIndex)
	if err != nil {
		return nil, err
	}
	if index > facilitiesLength-1 {
		err := fmt.Errorf("invalid choice: %v", f)
		return nil, err
	}

	return &index, nil

}
