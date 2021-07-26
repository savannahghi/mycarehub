package ussd

import (
	"context"
	"time"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
)

const (
	// OptOutFromMarketingInput indicates users who don't want to be send marketing sms(messages)
	OptOutFromMarketingInput = "1"
)

// WelcomeMenu represents  the default welcome submenu
func (u *Impl) WelcomeMenu() string {
	resp := "CON Welcome to Be.Well\r\n"
	resp += "1. Opt out from marketing messages\r\n"
	resp += "2. Change PIN"
	return resp
}

// ResetPinMenu ...
func (u *Impl) ResetPinMenu() string {
	resp := "CON Your PIN was reset successfully.\r\n"
	resp += "1. Opt out from marketing messages\r\n"
	resp += "2. Change PIN"
	return resp
}

// HandleHomeMenu represents the default home menu
func (u *Impl) HandleHomeMenu(ctx context.Context, level int, session *domain.USSDLeadDetails, userResponse string) string {
	ctx, span := tracer.Start(ctx, "HandleHomeMenu")
	defer span.End()

	time := time.Now()

	if userResponse == EmptyInput || userResponse == GoBackHomeInput {
		return u.WelcomeMenu()

	} else if userResponse == OptOutFromMarketingInput {
		option := "STOP"
		err := u.profile.SetOptOut(ctx, option, session.PhoneNumber)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture enter old PIN event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			USSDEventName:     USSDOptOut,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		resp := "CON We have successfully opted you\r\n"
		resp += "out of marketing messages\r\n"
		resp += "0. Go back home"
		return resp

	} else if userResponse == ChangePINInput {
		err := u.UpdateSessionLevel(ctx, ChangeUserPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		return u.HandleChangePIN(ctx, session, userResponse)

	} else {
		resp := "CON Invalid choice. Try again.\r\n"
		resp += "1. Opt out from marketing messages\r\n"
		resp += "2. Change PIN"
		return resp
	}
}
