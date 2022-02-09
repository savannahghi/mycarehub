package domain

import "time"

// ServiceRequest is a domain entity that represents a service request.
type ServiceRequest struct {
	ID           string     `json:"id"`
	RequestType  string     `json:"requestType"`
	Request      string     `json:"request"`
	Status       string     `json:"status"`
	ClientID     string     `json:"clientID"`
	InProgressAt *time.Time `json:"inProgressAt"`
	InProgressBy *string    `json:"inProgressBy"`
	ResolvedAt   *time.Time `json:"resolvedAt"`
	ResolvedBy   *string    `json:"resolvedBy"`
	FacilityID   *string    `json:"facility_id"`
}
