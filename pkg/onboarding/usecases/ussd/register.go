package ussd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/scalarutils"
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
	RegOptOutInput = "2"
	//RegChangePINInput ...
	RegChangePINInput = "2"

	//USSDDialled ...
	USSDDialled = "dialled"
	//USSDSelectRegister ...
	USSDSelectRegister = "select register"
	//USSDEnterFirstname ...
	USSDEnterFirstname = "entered firstname"
	//USSDEnterLastname ...
	USSDEnterLastname = "entered lastname"
	//USSDEnterDOB ...
	USSDEnterDOB = "entered date of birth"
	//USSDEnterPIN ...
	USSDEnterPIN = "entered PIN"
	//USSDConfirmPIN ...
	USSDConfirmPIN = "confirmed PIN"
	//USSDOptOut ...
	USSDOptOut = "opted out"
)

var userFirstName string
var userLastName string
var date string

// HandleUserRegistration ...
func (u *Impl) HandleUserRegistration(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	ctx, span := tracer.Start(ctx, "HandleUserRegistration")
	defer span.End()

	time := time.Now()

	if userResponse == EmptyInput || userResponse == GoBackHomeInput && session.Level == InitialState {
		//Creating contact stub on first USSD Dial
		if err := u.onboardingRepository.StageCRMPayload(ctx, &dto.ContactLeadInput{
			ContactType:    "phone",
			ContactValue:   session.PhoneNumber,
			IsSync:         false,
			TimeSync:       &time,
			OptOut:         "NO",
			WantCover:      false,
			ContactChannel: "USSD",
			IsRegistered:   false,
		}); err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture dialling event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             InitialState,
			USSDEventName:     USSDDialled,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		resp := "CON Welcome to Be.Well\r\n"
		resp += "1. Register\r\n"
		resp += "2. Opt Out\r\n"
		return resp
	}

	if userResponse == RegOptOutInput && session.Level == InitialState {
		option := "STOP"
		err := u.profile.SetOptOut(ctx, option, session.PhoneNumber)
		if err != nil {
			return "END Something went wrong. Please try again."
		}

		// Capture opt out event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             InitialState,
			USSDEventName:     USSDOptOut,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		resp := "CON We have successfully opted you\r\n"
		resp += "out of marketing messages\r\n"
		resp += "0. Go back home"
		return resp
	}

	if userResponse == RegisterInput && session.Level == InitialState {
		// Capture opt out event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             InitialState,
			USSDEventName:     USSDSelectRegister,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}
		err := u.UpdateSessionLevel(ctx, GetFirstNameState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		resp := "CON Please enter your firstname(e.g.\r\n"
		resp += "John).\r\n"
		return resp
	}

	if session.Level == GetFirstNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}
		userFirstName = userResponse

		err = u.UpdateSessionLevel(ctx, GetLastNameState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture first name event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             GetFirstNameState,
			USSDEventName:     USSDEnterFirstname,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		resp := "CON Please enter your lastname(e.g.\r\n"
		resp += "Doe)\r\n"
		return resp

	}
	fmt.Println("the last name level", session.Level)
	if session.Level == GetLastNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON Invalid name. Please enter a valid name (e.g Doe)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g Doe)"
		}

		userLastName = userResponse

		err = u.UpdateSessionLevel(ctx, GetDOBState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err.Error()
		}

		// Capture last name event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             GetLastNameState,
			USSDEventName:     USSDEnterLastname,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		resp := "CON Please enter your date of birth in\r\n"
		resp += "DDMMYYYY format e.g 14031996 for\r\n"
		resp += "14th March 1996\r\n"
		return resp
	}

	if session.Level == GetDOBState {
		err := utils.ValidateDateDigits(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}

		err = utils.ValidateDateLength(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}
		resp := utils.ValidateYearOfBirth(userResponse)
		if resp != "" {
			return resp
		}

		date = userResponse

		err = u.UpdateSessionLevel(ctx, GetPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err.Error()
		}

		// Capture date of birth event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             GetDOBState,
			USSDEventName:     USSDEnterDOB,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		return "CON Please enter a 4 digit PIN to secure your account"
	}

	if session.Level == GetPINState {
		err := utils.ValidatePIN(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON The PIN you entered in not correct please enter a 4 digit PIN"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		err = u.UpdateSessionLevel(ctx, SaveRecordState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err.Error()
		}

		// Capture PIN entry event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             GetPINState,
			USSDEventName:     USSDEnterPIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
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
		dateofBirth := &scalarutils.Date{
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
			utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture confirm PIN entry event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             SaveRecordState,
			USSDEventName:     USSDConfirmPIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}
	resp := "CON Invalid choice. Try again.\r\n"
	resp += "1. Register\r\n"
	resp += "2. Opt Out\r\n"
	return resp
}
