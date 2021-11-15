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

// Validate helps with validation of PINInput fields
func (f *PINInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}
