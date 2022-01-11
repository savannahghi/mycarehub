package servicerequest

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
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

// UseCaseServiceRequest holds all the interfaces that represent the service request business logic
type UseCaseServiceRequest interface {
	ICreateServiceRequest
}

// UseCasesServiceRequestImpl embeds the service request logic
type UseCasesServiceRequestImpl struct {
	Create infrastructure.Create
}

// NewUseCaseServiceRequestImpl creates a new service request instance
func NewUseCaseServiceRequestImpl(
	create infrastructure.Create,
) *UseCasesServiceRequestImpl {
	return &UseCasesServiceRequestImpl{
		Create: create,
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
