package dto

import (
	"net/url"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string           `json:"name"`
	Code        string           `json:"code"`
	Active      bool             `json:"active"`
	County      enums.CountyType `json:"county"`
	Description string           `json:"description"`
}

// FacilityFilterInput is used to supply filter parameters for healthcare facility filter inputs
type FacilityFilterInput struct {
	Search  *string `json:"search"`
	Name    *string `json:"name"`
	MFLCode *string `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilityFilterInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.Name != nil {
		vals.Add("name", *i.Name)
	}
	if i.MFLCode != nil {
		vals.Add("code", *i.MFLCode)
	}
	return vals
}

// FacilitySortInput is used to supply sort input for healthcare facility list queries
type FacilitySortInput struct {
	Name    *enumutils.SortOrder `json:"name"`
	MFLCode *enumutils.SortOrder `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilitySortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.MFLCode != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("code", "number")
		} else {
			vals.Add("code", "-number")
		}
	}
	return vals
}

// PaginationsInput contains fields required for pagination
type PaginationsInput struct {
	Limit       int        `json:"limit"`
	CurrentPage int        `json:"currentPage"`
	Sort        SortsInput `json:"sort"`
}

// FiltersInput contains fields required for filtering
type FiltersInput struct {
	DataType enums.FilterSortDataType
	Value    string // TODO: Clear spec on validation e.g dates must be ISO 8601. This is the actual data being filtered
}

// SortsInput includes the fields required for sorting the different types of fields
type SortsInput struct {
	Direction enums.SortDataType
	Field     enums.FilterSortDataType
}

// LoginInput represents the Login input data structure
type LoginInput struct {
	PhoneNumber *string         `json:"phoneNumber"`
	PIN         *string         `json:"pin"`
	Flavour     feedlib.Flavour `json:"flavour"`
}
