package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// ServiceRequest is a domain entity that represents a service request.
type ServiceRequest struct {
	ID             string                 `json:"id"`
	RequestType    string                 `json:"requestType"`
	Request        string                 `json:"request"`
	Status         string                 `json:"status"`
	Active         bool                   `json:"active"`
	ClientID       string                 `json:"clientID"`
	StaffID        string                 `json:"staffID"`
	CreatedAt      time.Time              `json:"created"`
	InProgressAt   *time.Time             `json:"inProgressAt"`
	InProgressBy   *string                `json:"inProgressBy"`
	ResolvedAt     *time.Time             `json:"resolvedAt"`
	ResolvedBy     *string                `json:"resolvedBy"`
	ResolvedByName *string                `string:"resolvedByName"`
	FacilityID     string                 `json:"facility_id"`
	ClientName     *string                `json:"client_name"`
	ClientContact  *string                `json:"client_contact"`
	Meta           map[string]interface{} `json:"meta"`
}

// RequestTypeCount ...
type RequestTypeCount struct {
	RequestType enums.ServiceRequestType `json:"requestType"`
	Total       int                      `json:"total"`
}

// ServiceRequestsCount ...
type ServiceRequestsCount struct {
	Total             int                 `json:"total"`
	RequestsTypeCount []*RequestTypeCount `json:"requestsTypeCount"`
}

// UpdateServiceRequestsPayload defined a list of service requests to synchronize MyCareHub with.
type UpdateServiceRequestsPayload struct {
	ServiceRequests []ServiceRequest `json:"serviceRequests" validate:"required"`
}
