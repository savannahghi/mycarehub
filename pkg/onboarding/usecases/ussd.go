package usecases

import (
	"context"
	"strings"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	firstName           = "CON Please enter your first name (e.g John)"
	invalidName         = "CON Invalid name. Please enter a valid name (e.g John)"
	lastName            = "CON Please enter your last name (eg Doe)"
	dateOfBirth         = "CON Please enter your date of birth in DDMMYYYY format e.g 14031996 for 14th March 1992"
	invalidDateOfBirth  = "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
	pin                 = "CON Please enter a 4 digit PIN to secure your account"
	invalidPIN          = "CON Invalid PIN. Please enter a 4 digit PIN to secure your account"
	confirmPIN          = "CON Please enter a 4 digit PIN again to confirm"
	currentPIN          = "CON Enter your old PIN to continue"
	mismatchedPIN       = "CON PIN mismatch. Please enter a PIN that matches the first PIN"
	mismatchedChangePIN = "CON PIN mismatch. Please enter a PIN that matches your current PIN"
	newPIN              = "CON Enter a new four digit PIN"
	confirmNewPIN       = "CON Please enter the 4 digit PIN again to confirm"
)

var temporaryPINHolder string

//UssdUsecase represent the logic involved in receiving a USSD
type UssdUsecase interface {
	GenerateUSSD(context context.Context, input *dto.SessionDetails) string
}

//UssdImpl represents usecase implementation object
type UssdImpl struct {
	baseExt              extension.BaseExtension
	onboardingRepository repository.OnboardingRepository
	profile              ProfileUseCase
}

//NewUssdUsecases returns a new Ussd usecase
func NewUssdUsecases(
	repository repository.OnboardingRepository,
	ext extension.BaseExtension,
	profileUsecase ProfileUseCase,
) UssdUsecase {
	return &UssdImpl{
		baseExt:              ext,
		onboardingRepository: repository,
		profile:              profileUsecase,
	}
}

//GenerateUSSD generates the USSD response
func (u *UssdImpl) GenerateUSSD(ctx context.Context, payload *dto.SessionDetails) string {
	// convert text into an array of user responses
	ussdTextArray := strings.Split(payload.Text, "*")

	return u.userWithNoAccountMenu(ctx, payload, ussdTextArray)
}

// UpdateSessionLevel updates user current level of interaction with USSD
func (u *UssdImpl) UpdateSessionLevel(ctx context.Context, level int, sessionID string) error {
	//increase level by 1
	c := level + 1

	_, err := u.onboardingRepository.UpdateSessionLevel(ctx, sessionID, c)
	if err != nil {
		return err
	}
	return nil

}

func (u *UssdImpl) userWithNoAccountMenu(ctx context.Context, payload *dto.SessionDetails, textArray []string) string {
	// if the text field is empty, this indicates that this is the begining of a session
	if len(payload.Text) == 0 {
		payload.Level = 0
		err := u.AddAITSessionDetails(ctx, payload)
		if err != nil {
			return err.Error()
		}

		return "CON Welcome to Be.Well \n1. Register"
	}

	if textArray[0] == "1" {
		//Getting level
		sessionDetails, err := u.onboardingRepository.GetAITSessionDetails(ctx, payload.SessionID)
		if err != nil {
			return err.Error()
		}

		return u.USSDSignupFlow(ctx, payload.Text, sessionDetails.Level, payload.SessionID, textArray)
	}

	userChoice := utils.GetUserChoice(payload.Text, 1)

	sessionDetails, err := u.onboardingRepository.GetAITSessionDetails(ctx, payload.SessionID)
	if err != nil {
		return err.Error()
	}

	if sessionDetails.Level == 0 && userChoice != "1" {
		return "CON Invalid choice.Please try again. \n1. Register"
	}

	return u.USSDSignupFlow(ctx, userChoice, sessionDetails.Level, payload.SessionID, textArray)
}

// USSDSignupFlow ...
func (u *UssdImpl) USSDSignupFlow(ctx context.Context, text string, level int, sessionID string, textArray []string) string {
	var response string
	if text == "1" {
		response = firstName
		err := u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		return response
	}

	if level == 1 {
		firstname := utils.GetUserChoice(text, 1)

		err := utils.ValidateUSSDInput(firstname)
		if err != nil {
			return invalidName
		}

		isLetter := utils.IsLetter(firstname)
		if !isLetter {
			return invalidName
		}

		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		response = lastName
		return response
	}

	if level == 2 {
		lastname := utils.GetUserChoice(text, 1)

		err := utils.ValidateUSSDInput(lastName)
		if err != nil {
			return invalidName
		}

		isLetter := utils.IsLetter(lastname)
		if !isLetter {
			return invalidName
		}

		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}

		response = "CON Please enter your date of birth in DDMMYYYY format e.g 14031996 for 14th March 1992"
		return response
	}

	if level == 3 {
		dob := utils.GetUserChoice(text, 1)

		err := utils.ValidateDateDigits(dob)
		if err != nil {
			return invalidDateOfBirth
		}

		err = utils.ValidateDateLength(dob)
		if err != nil {
			return invalidDateOfBirth
		}

		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		response = pin
		return response
	}

	if level == 4 {
		pin := utils.GetUserChoice(text, 1)

		err := utils.ValidatePIN(pin)
		if err != nil {
			return invalidPIN
		}

		temporaryPINHolder = pin

		response = confirmPIN

		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}

		return response
	}

	if level == 5 {

		pin := temporaryPINHolder
		confirmedPIN := utils.GetUserChoice(text, 1)

		if pin != confirmedPIN {
			return mismatchedPIN
		}

		_, err := u.onboardingRepository.UpdateSessionPIN(ctx, sessionID, confirmedPIN)
		if err != nil {
			return err.Error()
		}

		response = "CON Thanks for signing up for Be.Well \n1. Opt out from marketing messages \n2. Change PIN"

		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		return response
	}

	if level == 6 {
		userOption := utils.GetUserChoice(text, 1)
		switch userOption {
		case "1":
			response = "END We have successfully opted you out of marketing messages"
			return response
		case "2":
			response = currentPIN
			err := u.UpdateSessionLevel(ctx, level, sessionID)
			if err != nil {
				return err.Error()
			}
			return response
		default:
			resp := "CON Invalid choice. Please try again.\n"
			response = "1.Opt out from marketing messages \n2. Change PIN"
			resp += response
			return resp
		}

	}

	if level == 7 {
		//Fetching old PIN for comparison
		sessionDetails, err := u.onboardingRepository.GetAITSessionDetails(ctx, sessionID)
		if err != nil {
			return err.Error()
		}
		// check if old PIN is correct here
		currentPIN := utils.GetUserChoice(text, 1)

		if sessionDetails.PIN != currentPIN {
			return mismatchedChangePIN
		}
		response = newPIN
		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		return response
	}

	if level == 8 {
		// pin preprocessing
		newPIN := utils.GetUserChoice(text, 1)

		err := utils.ValidatePIN(newPIN)
		if err != nil {
			return invalidPIN
		}

		temporaryPINHolder = newPIN

		response = "CON Please enter the 4 digit PIN again to confirm"
		err = u.UpdateSessionLevel(ctx, level, sessionID)
		if err != nil {
			return err.Error()
		}
		return response
	}

	if level == 9 {
		newPIN := temporaryPINHolder
		newConfirmedPIN := utils.GetUserChoice(text, 1)

		if newPIN != newConfirmedPIN {
			return "CON PIN mismatch. Please enter a PIN that matches the first PIN"
		}

		_, err := u.onboardingRepository.UpdateSessionPIN(ctx, sessionID, newConfirmedPIN)
		if err != nil {
			return err.Error()
		}

		response = "END Your PIN was changed successfully"
		return response
	}

	return "END Invalid choice."
}

//AddAITSessionDetails persists USSD details
func (u *UssdImpl) AddAITSessionDetails(ctx context.Context, input *dto.SessionDetails) error {
	phone, err := base.NormalizeMSISDN(*input.PhoneNumber)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	sessionDetails := &dto.SessionDetails{
		PhoneNumber: phone,
		SessionID:   input.SessionID,
		Level:       input.Level,
	}
	err = u.onboardingRepository.AddAITSessionDetails(ctx, sessionDetails)
	if err != nil {
		return err
	}
	return nil
}
