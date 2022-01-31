package servicerequest

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
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

// IGetServiceRequests is an interface that holds the method signature for getting service requests
type IGetServiceRequests interface {
	GetServiceRequests(
		ctx context.Context,
		requestType *string,
		requestStatus *string,
	) ([]*domain.ServiceRequest, error)
}

// UseCaseServiceRequest holds all the interfaces that represent the service request business logic
type UseCaseServiceRequest interface {
	ICreateServiceRequest
	IGetServiceRequests
}

// UseCasesServiceRequestImpl embeds the service request logic
type UseCasesServiceRequestImpl struct {
	Create infrastructure.Create
	Query  infrastructure.Query
}

// NewUseCaseServiceRequestImpl creates a new service request instance
func NewUseCaseServiceRequestImpl(
	create infrastructure.Create,
	query infrastructure.Query,
) *UseCasesServiceRequestImpl {
	return &UseCasesServiceRequestImpl{
		Create: create,
		Query:  query,
	}
}

// CreateServiceRequest creates a service request
func (u *UseCasesServiceRequestImpl) CreateServiceRequest(
	ctx context.Context,
	clientID string,
	requestType, request string,
) (bool, error) {
	serviceRequest := &domain.ClientServiceRequest{
		Active:       true,
		RequestType:  requestType,
		Request:      request,
		Status:       "PENDING",
		InProgressAt: time.Now(),
		ClientID:     clientID,
	}
	err := u.Create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to create service request: %v", err)
	}
	return true, nil
}

// GetServiceRequests gets service requests based on the parameters provided
func (u *UseCasesServiceRequestImpl) GetServiceRequests(
	ctx context.Context,
	requestType *string,
	requestStatus *string,
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

	return u.Query.GetServiceRequests(ctx, requestType, requestStatus)
}
