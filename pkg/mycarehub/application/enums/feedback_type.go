package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FeedbackType is a list of all the feedback types.
type FeedbackType string

const (
	//GeneralFeedbackType represents the general feedback type
	GeneralFeedbackType FeedbackType = "GENERAL_FEEDBACK"
	// ServiceFeedbackType represents the service feedback type
	ServiceFeedbackType FeedbackType = "SERVICES_OFFERED"
)

// AllFeedbackTypes is a set of a  valid and known feedback types.
var AllFeedbackTypes = []FeedbackType{
	GeneralFeedbackType,
	ServiceFeedbackType,
}

//IsValid returns true if a feedback type is valid
func (f FeedbackType) IsValid() bool {
	switch f {
	case GeneralFeedbackType, ServiceFeedbackType:
		return true
	}
	return false
}

func (f FeedbackType) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a feedback type.
func (f *FeedbackType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = FeedbackType(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid FeedbackType", str)
	}
	return nil
}

//MarshalGQL writes the feedback type to the supplied
func (f FeedbackType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}
