package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/serverutils"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
)

const (
	// Default min length of the date
	minDateLength = 8
	// Default max length of the date
	maxDateLength = 8

	// Default min length of the date
	minPINLength = 4
	// Default max length of the date
	maxPINLength = 4
)

// ValidateUID checks that the uid supplied in the indicated request is valid
func ValidateUID(w http.ResponseWriter, r *http.Request) (*dto.UIDPayload, error) {
	p := &dto.UIDPayload{}
	serverutils.DecodeJSONToTargetStruct(w, r, p)
	if p.UID == nil {
		err := fmt.Errorf("invalid credentials, expected a uid")
		return nil, err
	}
	return p, nil
}

// ValidateSignUpInput returns a valid sign up input
func ValidateSignUpInput(input *dto.SignUpInput) (*dto.SignUpInput, error) {
	if !input.Flavour.IsValid() {
		return nil, exceptions.WrongEnumTypeError(input.Flavour.String())
	}

	phone, err := converterandformatter.NormalizeMSISDN(*input.PhoneNumber)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	err = extension.ValidatePINLength(*input.PIN)
	if err != nil {
		return nil, exceptions.ValidatePINLengthError(err)
	}

	err = extension.ValidatePINDigits(*input.PIN)
	if err != nil {
		return nil, exceptions.ValidatePINDigitsError(err)
	}

	if input.OTP == nil {
		return nil, exceptions.MissingInputError("otp")
	}

	return &dto.SignUpInput{
		PhoneNumber: phone,
		PIN:         input.PIN,
		Flavour:     input.Flavour,
		OTP:         input.OTP,
	}, nil
}

// ValidatePIN ...
func ValidatePIN(pin string) error {
	validatePINErr := ValidatePINLength(pin)
	if validatePINErr != nil {
		return validatePINErr
	}

	pinDigitsErr := extension.ValidatePINDigits(pin)
	if pinDigitsErr != nil {
		return pinDigitsErr
	}
	return nil
}

// ValidatePINLength ...
func ValidatePINLength(pin string) error {
	// make sure pin length is [4]
	if len(pin) < minPINLength || len(pin) > maxPINLength {
		return exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
	}
	return nil
}

// IsLetter ...
func IsLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

//ValidateDateLength ensures that the dates are of only 8 numbers
func ValidateDateLength(date string) error {
	// make sure date length is [8]
	if len(date) < minDateLength || len(date) > maxDateLength {
		return fmt.Errorf("date should be of only 8 digits")
	}
	return nil
}

// ValidateDateDigits validates user pin to ensure a PIN only contains digits
func ValidateDateDigits(pin string) error {
	// ensure pin is only digits
	_, err := strconv.ParseUint(pin, 10, 64)
	if err != nil {
		return fmt.Errorf("date can only be numbers")
	}
	return nil
}

//GetUserResponse gets the concatenated text from Africans Talking and splits it to get the current user input
func GetUserResponse(text string) string {
	response := strings.Split(text, "*")
	lastUserInput := response[len(response)-1]
	return lastUserInput
}

//ValidateYearOfBirth validates that the year enter is 18 years and above
func ValidateYearOfBirth(date string) string {
	year, _, _ := time.Now().Date()
	dayEntered, _ := strconv.Atoi(date[0:2])
	monthEntered, _ := strconv.Atoi(date[2:4])
	yearEntered, _ := strconv.Atoi(date[4:8])
	if dayEntered <= 0 || dayEntered > 31 {
		return "CON Wrong date value. Please enter a valid date"
	}
	if monthEntered <= 0 || monthEntered > 12 {
		return "CON Wrong month value. Please enter a valid month"
	}

	age := year - yearEntered
	if age < 18 {
		return "END Your age needs to be 18 years and above"
	}
	return ""

}
