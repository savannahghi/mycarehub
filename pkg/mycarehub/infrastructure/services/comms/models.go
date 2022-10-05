package comms

// APIResponse
type APIResponse struct {
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Token response
// The refresh token has a lifetime of 24 hours
// The access token has a lifetime of 30 minutes.
type TokenResponse struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

// ErrorMessage
type ErrorMessage struct {
	Message string `json:"message,omitempty"`
}

// ErrorResponse is the data response when an api request fails
type ErrorResponse struct {
	Detail  string         `json:"detail,omitempty"`
	Code    string         `json:"code,omitempty"`
	Message []ErrorMessage `json:"message,omitempty"`
}

// BulkSMSResponse
type BulkSMSResponse struct {
	GUID       string   `json:"guid"`
	Sender     string   `json:"sender"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
	State      string   `json:"state"`
	Sms        []string `json:"sms"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
}
