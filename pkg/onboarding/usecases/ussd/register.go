package ussd

import (
	"context"
	"strconv"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// InitialState ...
	InitialState = 0
	// GetFirstNameState ...
	GetFirstNameState = 1
	// GetLastNameState ...
	GetLastNameState = 2
	// GetDOBState ...
	GetDOBState = 3
	// GetPINState ...
	GetPINState = 4
	// SaveRecordState ...
	SaveRecordState = 5
	// RegisterInput ...
	RegisterInput = "1"
	//RegOptOutInput ...
	RegOptOutInput = "1"
	//RegChangePINInput ...
	RegChangePINInput = "2"
)

var userFirstName string
var userLastName string
var date string

// HandleUserRegistration ...
func (u *Impl) HandleUserRegistration(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	//Creating contact stub on first USSD Dial
	if err := u.onboardingRepository.StageCRMPayload(ctx, dto.ContactLeadInput{
		ContactType:    "phone",
		ContactValue:   session.PhoneNumber,
		IsSync:         false,
		TimeSync:       &time.Time{},
		OptOut:         "NO",
		WantCover:      false,
		ContactChannel: "USSD",
		IsRegistered:   false,
	}); err != nil {
		return "END Something went wrong. Please try again."
	}

	if userResponse == EmptyInput || userResponse == GoBackHomeInput && session.Level == InitialState {
		resp := "CON Welcome to Be.Well\r\n"
		resp += "1. Register\r\n"
		return resp
	}

	if userResponse == RegisterInput && session.Level == InitialState {
		err := u.UpdateSessionLevel(ctx, GetFirstNameState, session.SessionID)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		resp := "CON Please enter your firstname(e.g.\r\n"
		resp += "John).\r\n"
		return resp
	}

	if session.Level == GetFirstNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}
		userFirstName = userResponse

		err = u.UpdateSessionLevel(ctx, GetLastNameState, session.SessionID)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		resp := "CON Please enter your lastname(e.g.\r\n"
		resp += "Doe)\r\n"
		return resp

	}

	if session.Level == GetLastNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		userLastName = userResponse

		err = u.UpdateSessionLevel(ctx, GetDOBState, session.SessionID)
		if err != nil {
			return err.Error()
		}

		resp := "CON Please enter your date of birth in\r\n"
		resp += "DDMMYYYY format e.g 14031996 for\r\n"
		resp += "14th March 1992\r\n"
		return resp
	}

	if session.Level == GetDOBState {
		err := utils.ValidateDateDigits(userResponse)
		if err != nil {
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}

		err = utils.ValidateDateLength(userResponse)
		if err != nil {
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}
		resp := utils.ValidateYearOfBirth(userResponse)
		if resp != "" {
			return resp
		}

		date = userResponse

		err = u.UpdateSessionLevel(ctx, GetPINState, session.SessionID)
		if err != nil {
			return err.Error()
		}
		return "CON Please enter a 4 digit PIN to secure your account"
	}

	if session.Level == GetPINState {
		// TODO FIXME check for empty response
		err := utils.ValidatePIN(userResponse)
		if err != nil {
			return "CON The PIN you entered in not correct please enter a 4 digit PIN"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		err = u.UpdateSessionLevel(ctx, SaveRecordState, session.SessionID)
		if err != nil {
			return err.Error()
		}
		return "CON Please enter a 4 digit PIN again to confirm"

	}

	if session.Level == SaveRecordState {
		if userResponse != session.PIN {
			resp := "CON The PIN you entered does not match\r\n"
			resp += "Please enter a 4 digit PIN to secure your account\r\n"
			return resp
		}
		day, _ := strconv.Atoi(date[0:2])
		month, _ := strconv.Atoi(date[2:4])
		year, _ := strconv.Atoi(date[4:8])
		dateofBirth := &base.Date{
			Month: month,
			Day:   day,
			Year:  year,
		}
		updateInput := &dto.UserProfileInput{
			DateOfBirth: dateofBirth,
			FirstName:   &userFirstName,
			LastName:    &userLastName,
		}

		err := u.CreateUsddUserProfile(ctx, session.PhoneNumber, session.PIN, updateInput)
		if err != nil {
			return "END Something went wrong. Please try again."
		}

		contactLead := &dto.ContactLeadInput{
			FirstName:    userFirstName,
			LastName:     userLastName,
			DateOfBirth:  *dateofBirth,
			IsRegistered: true,
		}
		_ = u.onboardingRepository.UpdateStageCRMPayload(ctx, session.PhoneNumber, contactLead)

		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}
	resp := "CON Invalid choice. Try again.\r\n"
	resp += "1. Register\r\n"
	return resp
}
