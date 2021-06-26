package ussd

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// WantCoverInput indicates user intention to buy a cover
	WantCoverInput = "1"
	// OptOutFromMarketingInput indicates users who don't want to be send marketing sms(messages)
	OptOutFromMarketingInput = "2"
)

// WelcomeMenu represents  the default welcome submenu
func (u *Impl) WelcomeMenu() string {
	resp := "CON Welcome to Be.Well.\r\n"
	resp += "1. I want a cover.\r\n"
	resp += "2. Opt out from marketing messages.\r\n"
	resp += "3. Change PIN."
	return resp
}

// HandleHomeMenu represents the default home menu
func (u *Impl) HandleHomeMenu(ctx context.Context, level int, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == EmptyInput || userResponse == GoBackHomeInput {
		return u.WelcomeMenu()
	} else if userResponse == WantCoverInput {
		// TODO FIXME asynchronously send request to CRM
		resp := "CON We have recorded your request.\r\n"
		resp += "and one of the representatives will.\r\n"
		resp += "reach out to you. Thank you.\r\n"
		resp += "0. Go back home."
		return resp
	} else if userResponse == OptOutFromMarketingInput {
		// TODO FIXME asynchronously hook in opt out logic
		resp := "CON We have successfully opted you.\r\n"
		resp += "marketing messages.\r\n"
		resp += "0. Go back home."
		return resp
	} else if userResponse == ChangePINInput {
		err := u.UpdateSessionLevel(ctx, UserPINState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		return u.HandleChangePIN(ctx, session, userResponse)
	} else {
		// TODO FIXME return user to home
		return "END invalid input try again"
	}
}
