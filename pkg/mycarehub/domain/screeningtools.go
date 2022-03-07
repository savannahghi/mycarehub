package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// ScreeningToolQuestion defines a question for a screening tool
type ScreeningToolQuestion struct {
	ID               string                              `json:"questionID"`
	Question         string                              `json:"question"`
	ToolType         enums.ScreeningToolType             `json:"toolType"`
	ResponseChoices  map[string]interface{}              `json:"responseChoices"`
	ResponseType     enums.ScreeningToolResponseType     `json:"responseType"`
	ResponseCategory enums.ScreeningToolResponseCategory `json:"responseCategory"`
	Sequence         int                                 `json:"sequence"`
	Meta             map[string]interface{}              `json:"meta"`
	Active           bool                                `json:"active"`
}
