package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// ServiceRequest is a domain entity that represents a service request.
type ServiceRequest struct {
	ID                 string                 `json:"id"`
	RequestType        string                 `json:"requestType"`
	Request            string                 `json:"request"`
	Status             string                 `json:"status"`
	Active             bool                   `json:"active"`
	ClientID           string                 `json:"clientID,omitempty"`
	StaffID            string                 `json:"staffID,omitempty"`
	CreatedAt          time.Time              `json:"created"`
	InProgressAt       *time.Time             `json:"inProgressAt"`
	InProgressBy       *string                `json:"inProgressBy"`
	ResolvedAt         *time.Time             `json:"resolvedAt"`
	ResolvedBy         *string                `json:"resolvedBy"`
	ResolvedByName     *string                `string:"resolvedByName"`
	FacilityID         string                 `json:"facilityID,omitempty"`
	ClientName         *string                `json:"clientName,omitempty"`
	StaffName          *string                `json:"staffName,omitempty"`
	StaffContact       *string                `json:"staffContact,omitempty"`
	ClientContact      *string                `json:"clientContact,omitempty"`
	CCCNumber          *string                `json:"cccNumber,omitempty"`
	ScreeningToolName  string                 `json:"screeningToolName"`
	ScreeningToolScore string                 `json:"screeningToolScore"`
	ProgramID          string                 `json:"programID,omitempty"`
	OrganisationID     string                 `json:"organisationID,omitempty"`
	Meta               map[string]interface{} `json:"meta"`
	CaregiverID        string                 `json:"caregiverID"`
	CaregiverName      *string                `json:"caregiverName"`
	CaregiverContact   *string                `json:"caregiverContact"`
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

// ServiceRequestsCountResponse returns both clients and staff service requests
type ServiceRequestsCountResponse struct {
	ClientsServiceRequestCount *ServiceRequestsCount `json:"clientServiceRequestCount"`
	StaffServiceRequestCount   *ServiceRequestsCount `json:"staffServiceRequestCount"`
}

// UpdateServiceRequestsPayload defined a list of service requests to synchronize MyCareHub with.
type UpdateServiceRequestsPayload struct {
	ServiceRequests []ServiceRequest `json:"serviceRequests" validate:"required"`
}
