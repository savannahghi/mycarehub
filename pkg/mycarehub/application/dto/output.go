package dto

import "gopkg.in/go-playground/validator.v9"

// RestEndpointResponses represents the rest endpoints response(s) output
type RestEndpointResponses struct {
	Data map[string]interface{} `json:"data"`
}

// Validate helps with validation of ShareContentInput fields
func (f *RestEndpointResponses) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}
