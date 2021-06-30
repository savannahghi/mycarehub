package ussd

import (
	"context"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// ChangePINInput indicates the user intention to change their PIN
	ChangePINInput = "2"
	// ForgotPINInput indicates the user has forgotten their PIN and would like to reset it
	ForgotPINInput = "00"
	// ChangePINEnterNewPINState indicates the state at which user wants to set a new PIN
	ChangePINEnterNewPINState = 51
	// ChangePINProcessNewPINState indicates the state when the supplied PIN is being processed
	ChangePINProcessNewPINState = 52
	//ConfirmNewPInState indicates the state when a user is confirming a pin update
	ConfirmNewPInState = 53
	// PINResetEnterNewPINState indicates the state when the user wants to reset their PIN
	PINResetEnterNewPINState = 10
	// PINResetProcessState represents the state when the user has provided a wrong PIN
	PINResetProcessState = 11
	//ForgetPINResetState indicates the state when a use wants to reset PIN
	ForgetPINResetState = 13
	//ForgotPINVerifyDate indicates the state when a use wants to reset PIN
	ForgotPINVerifyDate = 15
)

// HandleChangePIN represents workflow used to change a user PIN
func (u *Impl) HandleChangePIN(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == EmptyInput || userResponse == ChangePINInput {
		// TODO FIXME validate/check if supplied PIN is correct
		err := u.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		resp := "CON Enter your old PIN to continue\r\n"
		resp += "0. Go back home\r\n"
		return resp
	}

	if userResponse == GoBackHomeInput {
		correctPin, err := u.LoginInUser(ctx, session.PhoneNumber, session.PIN, base.FlavourConsumer)
		if err != nil {
			return "END something went wrong. Please try again"
		}
		if !correctPin {
			return "CON Invalid PIN. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}

	if session.Level == ChangePINEnterNewPINState {
		correctPin, err := u.LoginInUser(ctx, session.PhoneNumber, userResponse, base.FlavourConsumer)
		if err != nil {
			return "END something went wrong. Please try again"
		}
		if !correctPin {
			return "CON Invalid PIN. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, ConfirmNewPInState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		resp := "CON Enter a new four digit PIN\r\n"
		return resp
	}
	if session.Level == ConfirmNewPInState {
		err := utils.ValidatePIN(userResponse)
		if err != nil {
			return "CON The PIN you entered in not correct. Please enter a 4 digit PIN"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			return "END Something wrong happened. Please try again. please retry again"
		}
		err = u.UpdateSessionLevel(ctx, ChangePINProcessNewPINState, session.SessionID)
		if err != nil {
			return err.Error()
		}
		return "CON Please enter a 4 digit PIN again to confirm"
	}

	if session.Level == ChangePINProcessNewPINState {
		if userResponse != session.PIN {
			resp := "CON The PIN you entered does not match\r\n"
			resp += "Please enter a 4 digit PIN that matches you PIN\r\n"
			return resp
		}
		_, err := u.ChangeUSSDUserPIN(ctx, session.PhoneNumber, userResponse)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		return u.ResetPinMenu()
	}
	return "END invalid input"
}

// HandlePINReset represents workflow used to reset to a user PIN
func (u *Impl) HandlePINReset(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	if session.Level == ForgetPINResetState {
		resp := "CON Please enter a new  4 digit PIN to\r\n"
		resp += "secure your account\r\n"
		err := u.UpdateSessionLevel(ctx, PINResetEnterNewPINState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		return resp
	}

	if session.Level == PINResetEnterNewPINState {

		err := utils.ValidatePIN(userResponse)
		if err != nil {
			return "CON The PIN you entered in not correct please enter a 4 digit PIN"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			return "END something wrong it happened"
		}
		err = u.UpdateSessionLevel(ctx, PINResetProcessState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		resp := "CON Please enter a 4 digit PIN again to\r\n"
		resp += "confirm.\r\n"
		return resp
	}
	if session.Level == PINResetProcessState {
		if userResponse != session.PIN {
			resp := "CON The PIN you entered does not match\r\n"
			resp += "Please enter a 4 digit PIN to\r\n"
			resp += "secure your account\r\n"
			return resp
		}
		_, err := u.ChangeUSSDUserPIN(ctx, session.PhoneNumber, userResponse)
		if err != nil {
			return "END Something wrong happened. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again."
		}
		return u.ResetPinMenu()
	}
	if session.Level == ForgotPINVerifyDate {
		profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, session.PhoneNumber, false)
		if err != nil {
			return "END something wrong it happened"
		}
		date := userResponse
		day, _ := strconv.Atoi(date[0:2])
		month, _ := strconv.Atoi(date[2:4])
		year, _ := strconv.Atoi(date[4:8])
		dateofBirth := &base.Date{
			Month: month,
			Day:   day,
			Year:  year,
		}
		if !reflect.DeepEqual(profile.UserBioData.DateOfBirth, dateofBirth) {
			return "CON Date of birth entered does not match the date of birth on record. Please enter your valid date of birth"
		}
		err = u.UpdateSessionLevel(ctx, UserPINResetState, session.SessionID)
		if err != nil {
			return "END Something wrong happened. Please try again."
		}

		session.Level = ForgetPINResetState
		return u.HandlePINReset(ctx, session, userResponse)
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
