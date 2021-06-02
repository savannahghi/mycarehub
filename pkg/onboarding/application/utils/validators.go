package utils

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
)

// ValidateUID checks that the uid supplied in the indicated request is valid
func ValidateUID(w http.ResponseWriter, r *http.Request) (*dto.UIDPayload, error) {
	p := &dto.UIDPayload{}
	base.DecodeJSONToTargetStruct(w, r, p)
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

	phone, err := base.NormalizeMSISDN(*input.PhoneNumber)
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

// ValidateAficasTalkingSMSData returns AIT validated SMS data
func ValidateAficasTalkingSMSData(input *dto.AfricasTalkingMessage) (*dto.AfricasTalkingMessage, error) {
	if input.LinkID == " " {
		return nil, fmt.Errorf("message `linkID` cannot be empty")
	}

	if input.Text == " " {
		return nil, fmt.Errorf("`text` message cannot be empty")
	}

	if input.To == " " {
		return nil, fmt.Errorf("`to` cannot be empty")
	}

	if input.ID == " " {
		return nil, fmt.Errorf("message `ID` cannot be empty")
	}

	if input.Date == " " {
		return nil, fmt.Errorf("`date` of sending cannot be empty")
	}

	if input.From == " " {
		return nil, fmt.Errorf("`phone` number cannot be empty")
	}

	_, err := base.NormalizeMSISDN(input.From)
	if err != nil {
		return nil, err
	}

	return &dto.AfricasTalkingMessage{
		Date:   input.Date,
		From:   input.From,
		ID:     input.ID,
		LinkID: input.LinkID,
		Text:   input.Text,
		To:     input.To,
	}, nil
}
