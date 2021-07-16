package ussd

import (
	"context"

	"github.com/savannahghi/feedlib"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases/ussd")

// HandleLogin represents the workflow for authenticating a user
func (u *Impl) HandleLogin(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	ctx, span := tracer.Start(ctx, "HandleLogin")
	defer span.End()

	switch userResponse {
	case EmptyInput:
		resp := "CON Welcome to Be.Well.Please enter\r\n"
		resp += "your PIN to continue(enter 00 if\r\n"
		resp += "you forgot your PIN)\r\n"
		return resp

	case ForgotPINInput:
		err := u.UpdateSessionLevel(ctx, ForgotPINVerifyDate, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		resp := "CON Please enter your date of birth in\r\n"
		resp += "DDMMYYYY format e.g 14031996 for\r\n"
		resp += "14th March 1996\r\n"
		resp += "to be able to reset PIN\r\n"
		return resp

	default:
		isLoggedIn, err := u.LoginInUser(ctx, session.PhoneNumber, userResponse, feedlib.FlavourConsumer)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		if !isLoggedIn {
			resp := "CON The PIN you entered is not correct\r\n"
			resp += "Please try again (enter 00 if you\r\n"
			resp += "forgot your PIN)"
			return resp
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}

}

// LoginInUser authenticates a user to allow them proceed to the home menu
func (u *Impl) LoginInUser(
	ctx context.Context,
	phone string,
	PIN string,
	flavour feedlib.Flavour,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "LoginInUser")
	defer span.End()

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		phone,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	PINData, err := u.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if PINData == nil {
		return false, err
	}
	matched := u.pinExt.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		return false, nil

	}
	return true, nil

}
