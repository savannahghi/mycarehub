package dto

import (
	"net/url"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"gorm.io/datatypes"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Active      bool   `json:"active"`
	County      string `json:"county"`
	Description string `json:"description"`
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

// MetricInput reprents the metrics data structure input
type MetricInput struct {

	// TODO Metric types should be a controlled list i.e enum
	Type domain.MetricType `json:"metric_type"`

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time `json:"time"`

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string `json:"uid"`
}
