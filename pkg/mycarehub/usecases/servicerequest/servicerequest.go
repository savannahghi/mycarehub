package servicerequest

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// Service requests are tasks for the healthcare staff on the platform. Some examples are:
// red flags raised by content posted to a group or entries into health diaries
// appointment reschedule requests
// These tasks will be presented on a list and notified (e.g., via push notifications).
// Each task will have a status. When first created, the tasks will be marked as “PENDING”.
// Once the relevant actions are taken, it will be possible to mark them as “IN PROGRESS”, “RESOLVED” and add relevant notes.
// In order to ensure that a task is not addressed by multiple people at the same time, each task will be updated with a record of the user and timestamp each time the status is changed.

// ICreateServiceRequest is an interface that holds the method signature for creating a service request
type ICreateServiceRequest interface {
	CreateServiceRequest(
		ctx context.Context,
		clientID string,
		requestType, request string,
	) (bool, error)
}

// ISetInProgresssBy is an interface that contains the method signature for assigning the staff currently working on a request
type ISetInProgresssBy interface {
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
}

// IGetServiceRequests is an interface that holds the method signature for getting service requests
type IGetServiceRequests interface {
	GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error)
	GetServiceRequestsCount(
		ctx context.Context,
		requestType *string,
		facilityID string,
	) (int, error)
}

// IResolveServiceRequest is an interface that holds the method signature for resolving a service request
type IResolveServiceRequest interface {
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
}

// UseCaseServiceRequest holds all the interfaces that represent the service request business logic
type UseCaseServiceRequest interface {
	ICreateServiceRequest
	IGetServiceRequests
	ISetInProgresssBy
	IResolveServiceRequest
}

// UseCasesServiceRequestImpl embeds the service request logic
type UseCasesServiceRequestImpl struct {
	Create infrastructure.Create
	Query  infrastructure.Query
	Update infrastructure.Update
}

// NewUseCaseServiceRequestImpl creates a new service request instance
func NewUseCaseServiceRequestImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
) *UseCasesServiceRequestImpl {
	return &UseCasesServiceRequestImpl{
		Create: create,
		Query:  query,
		Update: update,
	}
}

// CreateServiceRequest creates a service request
func (u *UseCasesServiceRequestImpl) CreateServiceRequest(
	ctx context.Context,
	clientID string,
	requestType, request string,
) (bool, error) {
	clientProfile, err := u.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ClientProfileNotFoundErr(err)
	}
	serviceRequest := &domain.ClientServiceRequest{
		Active:      true,
		RequestType: requestType,
		Request:     request,
		Status:      "PENDING",
		ClientID:    clientID,
		FacilityID:  clientProfile.FacilityID,
	}
	err = u.Create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to create service request: %v", err)
	}
	return true, nil
}

// SetInProgressBy assigns to a service request, staff currently working on the service request
func (u *UseCasesServiceRequestImpl) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	if requestID == "" || staffID == "" {
		return false, fmt.Errorf("request ID or staff ID cannot be empty")
	}
	return u.Update.SetInProgressBy(ctx, requestID, staffID)
}

// GetServiceRequests gets service requests based on the parameters provided
func (u *UseCasesServiceRequestImpl) GetServiceRequests(
	ctx context.Context,
	requestType *string,
	requestStatus *string,
	facilityID *string,
) ([]*domain.ServiceRequest, error) {
	if requestType != nil {
		if !enums.ServiceRequestType(*requestType).IsValid() {
			return nil, fmt.Errorf("invalid request type: %v", *requestType)
		}
	}
	if requestStatus != nil {
		if !enums.ServiceRequestStatus(*requestStatus).IsValid() {
			return nil, fmt.Errorf("invalid request status: %v", *requestStatus)
		}
	}

	return u.Query.GetServiceRequests(ctx, requestType, requestStatus, facilityID)
}

// GetServiceRequestsCount gets service requests count based on the parameters provided
func (u *UseCasesServiceRequestImpl) GetServiceRequestsCount(
	ctx context.Context,
	requestType *string,
	facilityID string,
) (int, error) {

	return u.Query.GetServiceRequestsCount(ctx, requestType, facilityID)
}

// ResolveServiceRequest resolves a service request
func (u *UseCasesServiceRequestImpl) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
	if staffID == nil {
		return false, fmt.Errorf("staff ID is required")
	}
	if serviceRequestID == nil {
		return false, fmt.Errorf("service request ID is required")
	}
	ok, err := u.Update.ResolveServiceRequest(ctx, staffID, serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update service request: %v", err)
	}

	if !ok {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to resolve service request")
	}

	return ok, nil
}
