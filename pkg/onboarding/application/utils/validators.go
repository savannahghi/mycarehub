package utils

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

// ValidateUID checks that the uid supplied in the indicated request is valid
func ValidateUID(w http.ResponseWriter, r *http.Request) (*resources.UIDPayload, error) {
	p := &resources.UIDPayload{}
	base.DecodeJSONToTargetStruct(w, r, p)
	if p.UID == nil {
		err := fmt.Errorf("invalid credentials, expected a uid")
		return nil, err
	}
	if p == nil {
		err := fmt.Errorf(
			"nil business partner UID struct after decoding input")
		return nil, err
	}

	return p, nil
}

// ValidateSignUpInput returns a valid sign up input
func ValidateSignUpInput(input *resources.SignUpInput) (*resources.SignUpInput, error) {
	if !input.Flavour.IsValid() {
		return nil, exceptions.WrongEnumTypeError(input.Flavour.String(), nil)
	}

	phone, err := base.NormalizeMSISDN(*input.PhoneNumber)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	err = ValidatePINLength(*input.PIN)
	if err != nil {
		return nil, exceptions.ValidatePINLengthError(err)
	}

	err = ValidatePINDigits(*input.PIN)
	if err != nil {
		return nil, exceptions.ValidatePINDigitsError(err)
	}

	if input.OTP == nil {
		return nil, exceptions.MissingInputError("otp")
	}

	return &resources.SignUpInput{
		PhoneNumber: &phone,
		PIN:         input.PIN,
		Flavour:     input.Flavour,
		OTP:         input.OTP,
	}, nil
}
