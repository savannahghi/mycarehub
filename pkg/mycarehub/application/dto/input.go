package dto

import (
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"gopkg.in/go-playground/validator.v9"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Code        int    `json:"code" validate:"required"`
	Active      bool   `json:"active"`
	County      string `json:"county" validate:"required"`
	Description string `json:"description" validate:"required,min=3,max=256"`
}

// Validate helps with validation of facility input fields
func (f *FacilityInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// PaginationsInput contains fields required for pagination
type PaginationsInput struct {
	Limit       int        `json:"limit"`
	CurrentPage int        `json:"currentPage" validate:"required"`
	Sort        SortsInput `json:"sort"`
}

// Validate helps with validation of PaginationsInput fields
func (f *PaginationsInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// FiltersInput contains fields required for filtering
type FiltersInput struct {
	DataType enums.FilterSortDataType `json:"dataType" validate:"required"`
	Value    string                   `json:"value" validate:"required"` // TODO: Clear spec on validation e.g dates must be ISO 8601. This is the actual data being filtered
}

// Validate helps with validation of FiltersInput fields
func (f *FiltersInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// SortsInput includes the fields required for sorting the different types of fields
type SortsInput struct {
	Direction enums.SortDataType       `json:"direction"`
	Field     enums.FilterSortDataType `json:"field"`
}

// LoginInput represents the Login input data structure
type LoginInput struct {
	PhoneNumber *string         `json:"phoneNumber" validate:"required"`
	PIN         *string         `json:"pin" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of LoginInput fields
func (f *LoginInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// PINInput represents the Pin input data structure
type PINInput struct {
	UserID     *string         `json:"id" validate:"required"`
	PIN        *string         `json:"pin" validate:"required"`
	ConfirmPIN *string         `json:"confirmPin" validate:"required"`
	Flavour    feedlib.Flavour `json:"flavour" validate:"required"`
}

// VerifyOTPInput represents the verify OTP input data structure
type VerifyOTPInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// SendOTPInput represents the send OTP input data structure
type SendOTPInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// SendRetryOTPPayload is used to define the inputs passed when calling the endpoint
// that resends an otp
type SendRetryOTPPayload struct {
	Phone   string          `json:"phoneNumber" validate:"required"`
	Flavour feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of PINInput fields
func (f *PINInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// SecurityQuestionResponseInput represents the SecurityQuestionResponse input data structure
type SecurityQuestionResponseInput struct {
	UserID             string `json:"userID" validate:"required"`
	SecurityQuestionID string `json:"securityQuestionID" validate:"required"`
	Response           string `json:"Response" validate:"required"`
}

// Validate helps with validation of SecurityQuestionResponseInput fields
func (f *SecurityQuestionResponseInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// VerifySecurityQuestionInput defines the field passed when verifying the set security questions
type VerifySecurityQuestionInput struct {
	QuestionID string          `json:"questionID" validate:"required"`
	Flavour    feedlib.Flavour `json:"flavour" validate:"required"`
	Response   string          `json:"response" validate:"required"`
	UserID     string          `json:"userID" validate:"required"`
}

// Validate checks to validate whether the field inputs for verifying a security question
// are filled
func (f *VerifySecurityQuestionInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// VerifyPhoneInput carries the OTP data used to send OTP messages to a particular phone number
type VerifyPhoneInput struct {
	PhoneNumber string          `json:"phoneNumber"`
	Flavour     feedlib.Flavour `json:"flavour"`
}

// GetSecurityQuestionsInput defines the field passed when getting the security questions
type GetSecurityQuestionsInput struct {
	Flavour feedlib.Flavour `json:"flavour" validate:"required"`
}

// GetUserRespondedSecurityQuestionsInput defines the field passed when getting the security questions
type GetUserRespondedSecurityQuestionsInput struct {
	PhoneNumber string          `json:"phonenumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
}

// Validate helps with validation of GetUserRespondedSecurityQuestionsInput fields
func (f *GetUserRespondedSecurityQuestionsInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// UserResetPinInput contains the fields requires when a user is resetting a pin
type UserResetPinInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
	PIN         string          `json:"pin" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
}

// Validate checks to validate whether the field inputs for verifying user pin
func (f *UserResetPinInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// ShareContentInput defines the field passed when sharing content
type ShareContentInput struct {
	UserID    string `json:"userID" validate:"required"`
	ContentID int    `json:"contentID" validate:"required"`
	Channel   string `json:"channel" validate:"required"`
}

// Validate helps with validation of ShareContentInput fields
func (f *ShareContentInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}
