package ussd

import (
	"context"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// WantCoverInput indicates user intention to buy a cover
	WantCoverInput = "1"
	// OptOutFromMarketingInput indicates users who don't want to be send marketing sms(messages)
	OptOutFromMarketingInput = "2"

	layoutISO = "01-02-2006"
)

// WelcomeMenu represents  the default welcome submenu
func (u *Impl) WelcomeMenu() string {
	resp := "CON Welcome to Be.Well\r\n"
	resp += "1. I want a cover\r\n"
	resp += "2. Opt out from marketing messages\r\n"
	resp += "3. Change PIN"
	return resp
}

func (u *Impl) ResetPinMenu() string {
	resp := "CON Your PIN was reset successfully.\r\n"
	resp += "1. I want a cover\r\n"
	resp += "2. Opt out from marketing messages\r\n"
	resp += "3. Change PIN"
	return resp
}

// HandleHomeMenu represents the default home menu
func (u *Impl) HandleHomeMenu(ctx context.Context, level int, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == EmptyInput || userResponse == GoBackHomeInput {
		return u.WelcomeMenu()
	} else if userResponse == WantCoverInput {
		// TODO FIXME asynchronously send request to CRM
		resp := "CON We have recorded your request\r\n"
		resp += "and one of the representatives will\r\n"
		resp += "reach out to you. Thank you\r\n"
		resp += "0. Go back home"

		validDate := utils.ParseUSSDDateInput(session.DateOfBirth)
		DOB, _ := time.Parse(layoutISO, validDate)
		payload := dto.ContactLeadInput{
			FirstName: session.FirstName,
			LastName:  session.LastName,
			DateOfBirth: base.Date{
				Year:  DOB.Year(),
				Month: int(DOB.Month()),
				Day:   DOB.Day(),
			},
			WantCover: true,
		}
		//Error shouldn't break USSD flow
		_ = u.onboardingRepository.StageCRMPayload(ctx, payload)

		return resp

	} else if userResponse == OptOutFromMarketingInput {
		option := "STOP"
		err := u.profile.SetOptOut(ctx, option, session.PhoneNumber)
		if err != nil {
			return "END Something wrong happened. Please try again."
		}
		resp := "CON We have successfully opted you\r\n"
		resp += "out of marketing messages\r\n"
		resp += "0. Go back home"
		return resp

	} else if userResponse == ChangePINInput {
		err := u.UpdateSessionLevel(ctx, ChangeUserPINState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		return u.HandleChangePIN(ctx, session, userResponse)

	} else {
		resp := "CON Invalid choice. Try again.\r\n"
		resp += "1. I want a cover\r\n"
		resp += "2. Opt out from marketing messages\r\n"
		resp += "3. Change PIN"
		return resp
	}
}
