package mock

import "context"

// ServiceRequestUseCaseMock mocks the service request instance
type ServiceRequestUseCaseMock struct {
	MockCreateServiceRequestFn func(
		ctx context.Context,
		clientID string,
		requestType, request string,
	) (bool, error)
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
