package ussd

import (
	"context"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// ChangePINInput indicates the user intention to change their PIN
	ChangePINInput = "3"
	// ForgotPINInput indicates the user has forgotten their PIN and would like to reset it
	ForgotPINInput = "00"
	// ChangePINEnterNewPINState indicates the state at which user wants to set a new PIN
	ChangePINEnterNewPINState = 51
	// ChangePINProcessNewPINState indicates the state when the supplied PIN is being processed
	ChangePINProcessNewPINState = 52
	// PINResetEnterNewPINState indicates the state when the user wants to reset their PIN
	PINResetEnterNewPINState = 10
	// PINResetProcessState represents the state when the user has provided a wrong PIN
	PINResetProcessState = 11
)

// HandleChangePIN represents workflow used to change a user PIN
func (u *Impl) HandleChangePIN(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == EmptyInput || userResponse == ChangePINInput {
		// TODO FIXME validate/check if supplied PIN is correct
		err := u.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		resp := "CON Enter your old PIN to continue\r\n"
		resp += "0. Go back home\r\n"
		return resp
	}

	if userResponse == GoBackHomeInput {
		err := u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}

	if session.Level == ChangePINEnterNewPINState {
		err := u.UpdateSessionLevel(ctx, ChangePINProcessNewPINState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		resp := "CON Enter a new four digit PIN\r\n"
		return resp
	}

	if session.Level == ChangePINProcessNewPINState {
		_, err := u.ChangeUSSDUserPIN(ctx, session.PhoneNumber, userResponse)
		if err != nil {
			return "END something is wrong"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}
	return "END invalid input"
}

// HandlePINReset represents workflow used to reset to a user PIN
func (u *Impl) HandlePINReset(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == ForgotPINInput {
		resp := "CON Please enter a 4 digit PIN to\r\n"
		resp += "secure your account\r\n"
		return resp
	}

	if session.Level == PINResetEnterNewPINState {
		_, err := u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			return "END something wrong it happened"
		}
		err = u.UpdateSessionLevel(ctx, PINResetProcessState, session.SessionID)
		if err != nil {
			return "END something is wrong"
		}
		resp := "CON Please enter a 4 digit PIN again to\r\n"
		resp += "confirm.\r\n"
		return resp
	}
	if session.Level == PINResetProcessState {
		if userResponse != session.PIN {
			err := u.UpdateSessionLevel(ctx, PINResetEnterNewPINState, session.SessionID)
			if err != nil {
				return "END something wrong happened"
			}
			resp := "CON The PIN you entered does not match\r\n"
			resp += "Please enter a 4 digit PIN to\r\n"
			resp += "secure your account\r\n"
			return resp
		}
		_, err := u.ChangeUSSDUserPIN(ctx, session.PhoneNumber, userResponse)
		if err != nil {
			return "END something is wrong"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END something wrong happened"
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}
	return "END something went wrong"
}

//SetUSSDUserPin sets user pin when a user registers via USSD
func (u *Impl) SetUSSDUserPin(ctx context.Context, phoneNumber string, PIN string) error {
	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		phoneNumber,
		false,
	)
	if err != nil {
		return err
	}

	_, err = u.pinUsecase.SetUserPIN(
		ctx,
		PIN,
		profile.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// ChangeUSSDUserPIN updates user's pin with the newly supplied pin via USSD
func (u *Impl) ChangeUSSDUserPIN(
	ctx context.Context,
	phone string,
	pin string,
) (bool, error) {
	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		phone,
		false,
	)
	if err != nil {
		return false, err
	}

	salt, encryptedPin := u.pinExt.EncryptPIN(pin, nil)
	pinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: profile.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	_, err = u.onboardingRepository.UpdatePIN(ctx, profile.ID, pinPayload)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}
