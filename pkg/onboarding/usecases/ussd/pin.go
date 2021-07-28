package ussd

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/scalarutils"
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
	// ConfirmNewPINState indicates the state when a user is confirming a pin update
	ConfirmNewPINState = 53
	// PINResetEnterNewPINState indicates the state when the user wants to reset their PIN
	PINResetEnterNewPINState = 10
	// PINResetProcessState represents the state when the user has provided a wrong PIN
	PINResetProcessState = 11
	//ForgetPINResetState indicates the state when a use wants to reset PIN
	ForgetPINResetState = 13
	//ForgotPINVerifyDate indicates the state when a use wants to reset PIN
	ForgotPINVerifyDate = 15

	//USSDChooseToChangePIN indicates user chose to change PIN
	USSDChooseToChangePIN = "chose to change PIN"
	//USSDEnterOldPIN is the event when user enters their old PIN
	USSDEnterOldPIN = "entered old PIN"
	//USSDEnterNewPIN ...
	USSDEnterNewPIN = "entered a new 4 digit PIN"
	//USSDConfirmChangePIN ...
	USSDConfirmChangePIN = "confirmed new PIN"
	//USSDChooseToGoBackHome ...
	USSDChooseToGoBackHome = "chose to go back home"

	//USSDChooseToResetPIN ...
	USSDChooseToResetPIN = "chose to reset PIN"
	//USSDChooseToConfirmResetPIN ...
	USSDChooseToConfirmResetPIN = "confirm reset PIN"
	//USSDPINResetVerifyDate ...
	USSDPINResetVerifyDate = "verify date of birth"
)

// HandleChangePIN represents workflow used to change a user PIN
func (u *Impl) HandleChangePIN(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	ctx, span := tracer.Start(ctx, "HandleChangePIN")
	defer span.End()

	time := time.Now()

	if userResponse == EmptyInput || userResponse == ChangePINInput {
		err := u.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		// Capture choose to change PIN
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             ChangePINEnterNewPINState,
			USSDEventName:     USSDChooseToChangePIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}
		resp := "CON Enter your old PIN to continue\r\n"
		resp += "0. Go back home\r\n"

		return resp
	}

	if userResponse == GoBackHomeInput {
		isLoggedInUser, err := u.LoginInUser(ctx, session.PhoneNumber, session.PIN, feedlib.FlavourConsumer)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		if !isLoggedInUser {
			return "CON Invalid PIN. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		// Capture go back home event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             HomeMenuState,
			USSDEventName:     USSDChooseToGoBackHome,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}

	if session.Level == ChangePINEnterNewPINState {
		isLoggedInUser, err := u.LoginInUser(ctx, session.PhoneNumber, userResponse, feedlib.FlavourConsumer)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		if !isLoggedInUser {
			return "CON Invalid PIN. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, ConfirmNewPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}

		resp := "CON Enter a new four digit PIN\r\n"
		return resp
	}
	if session.Level == ConfirmNewPINState {
		err := utils.ValidatePIN(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON The PIN you entered is invalid. Please try again"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}
		err = u.UpdateSessionLevel(ctx, ChangePINProcessNewPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return err.Error()
		}

		// Capture enter new PIN event
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             ChangePINEnterNewPINState,
			USSDEventName:     USSDEnterNewPIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		return "CON Please enter a 4 digit PIN again to confirm"
	}

	if session.Level == ChangePINProcessNewPINState {
		if userResponse != session.PIN {
			resp := "CON The PIN you entered does not match\r\n"
			resp += "Please enter a 4 digit PIN that matches your PIN\r\n"
			return resp
		}
		_, err := u.ChangeUSSDUserPIN(ctx, session.PhoneNumber, userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}

		// Capture confirm new PIN
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             ConfirmNewPINState,
			USSDEventName:     USSDConfirmChangePIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		return u.ResetPinMenu()
	}

	if userResponse != GoBackHomeInput && userResponse != EmptyInput && userResponse != ChangePINInput {
		resp := "CON Invalid choice. Please try again."
		return resp
	}

	return "END invalid input"
}

// HandlePINReset represents workflow used to reset to a user PIN
func (u *Impl) HandlePINReset(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	ctx, span := tracer.Start(ctx, "HandlePINReset")
	defer span.End()

	time := time.Now()

	if session.Level == ForgetPINResetState {
		resp := "CON Please enter a new 4 digit PIN to\r\n"
		resp += "secure your account\r\n"
		err := u.UpdateSessionLevel(ctx, PINResetEnterNewPINState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		return resp
	}

	if session.Level == PINResetEnterNewPINState {

		err := utils.ValidatePIN(userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "CON The PIN you entered is invalid. Please try again"
		}
		_, err = u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, PINResetProcessState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}

		// Capture reset PIN
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             PINResetEnterNewPINState,
			USSDEventName:     USSDChooseToResetPIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
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
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again"
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture confirm reset PIN
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             PINResetProcessState,
			USSDEventName:     USSDChooseToConfirmResetPIN,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		return u.ResetPinMenu()
	}
	if session.Level == ForgotPINVerifyDate {
		profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, session.PhoneNumber, false)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END something wrong it happened"
		}
		date := userResponse
		day, _ := strconv.Atoi(date[0:2])
		month, _ := strconv.Atoi(date[2:4])
		year, _ := strconv.Atoi(date[4:8])
		dateofBirth := &scalarutils.Date{
			Month: month,
			Day:   day,
			Year:  year,
		}
		if !reflect.DeepEqual(profile.UserBioData.DateOfBirth, dateofBirth) {
			return "CON Date of birth entered does not match the date of birth on record. Please enter your valid date of birth"
		}
		err = u.UpdateSessionLevel(ctx, UserPINResetState, session.SessionID)
		if err != nil {
			utils.RecordSpanError(span, err)
			return "END Something went wrong. Please try again."
		}

		// Capture verify DOB
		if _, err := u.onboardingRepository.SaveUSSDEvent(ctx, &dto.USSDEvent{
			SessionID:         session.SessionID,
			PhoneNumber:       session.PhoneNumber,
			USSDEventDateTime: &time,
			Level:             ForgotPINVerifyDate,
			USSDEventName:     USSDPINResetVerifyDate,
		}); err != nil {
			return "END Something went wrong. Please try again."
		}

		session.Level = ForgetPINResetState
		return u.HandlePINReset(ctx, session, userResponse)
	}
	return "END something went wrong"
}

//SetUSSDUserPin sets user pin when a user registers via USSD
func (u *Impl) SetUSSDUserPin(ctx context.Context, phoneNumber string, PIN string) error {
	ctx, span := tracer.Start(ctx, "SetUSSDUserPin")
	defer span.End()

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		phoneNumber,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	_, err = u.pinUsecase.SetUserPIN(
		ctx,
		PIN,
		profile.ID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
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
	ctx, span := tracer.Start(ctx, "ChangeUSSDUserPIN")
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

	salt, encryptedPin := u.pinExt.EncryptPIN(pin, nil)
	pinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: profile.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	_, err = u.onboardingRepository.UpdatePIN(ctx, profile.ID, pinPayload)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}
