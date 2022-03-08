package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ServiceRequestUseCaseMock mocks the service request instance
type ServiceRequestUseCaseMock struct {
	MockCreateServiceRequestFn func(
		ctx context.Context,
		clientID string,
		requestType, request string,
	) (bool, error)
	MockGetPendingServiceRequestsCountFn func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockGetServiceRequestsFn             func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error)
	MockResolveServiceRequestFn          func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	MockSetInProgressByFn                func(ctx context.Context, requestID string, staffID string) (bool, error)
}

// NewServiceRequestUseCaseMock initializes a new service request instance mock
func NewServiceRequestUseCaseMock() *ServiceRequestUseCaseMock {
	return &ServiceRequestUseCaseMock{
		MockCreateServiceRequestFn: func(
			ctx context.Context,
			clientID string,
			requestType, request string,
		) (bool, error) {
			return true, nil
		},
		MockGetPendingServiceRequestsCountFn: func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
			return &domain.ServiceRequestsCount{Total: 10}, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
			return []*domain.ServiceRequest{
				{
					ID: uuid.New().String(),
				},
			}, nil
		},
		MockResolveServiceRequestFn: func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
			return true, nil
		},
		MockSetInProgressByFn: func(ctx context.Context, requestID string, staffID string) (bool, error) {
			return true, nil
		},
	}
}

// CreateServiceRequest mocks the implementation for creating a service request
func (s *ServiceRequestUseCaseMock) CreateServiceRequest(
	ctx context.Context,
	clientID string,
	requestType, request string,
) (bool, error) {
	return s.MockCreateServiceRequestFn(ctx, clientID, requestType, request)
}

// GetPendingServiceRequestsCount mocks the method of getting the number of pending service requests count
func (s *ServiceRequestUseCaseMock) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return s.MockGetPendingServiceRequestsCountFn(ctx, facilityID)
}

// GetServiceRequests mocks the method for fetching service requests
func (s *ServiceRequestUseCaseMock) GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
	return s.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID)
}

// ResolveServiceRequest mocks resolving a service request
func (s *ServiceRequestUseCaseMock) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
	return s.MockResolveServiceRequestFn(ctx, staffID, serviceRequestID)
}

// SetInProgressBy mocks the implementation of marking a service request as in progress
func (s *ServiceRequestUseCaseMock) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	return s.MockSetInProgressByFn(ctx, requestID, staffID)
}
