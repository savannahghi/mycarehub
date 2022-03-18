package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ServiceRequestUseCaseMock mocks the service request instance
type ServiceRequestUseCaseMock struct {
	MockCreateServiceRequestFn func(
		ctx context.Context,
		clientID string,
		requestType, request, cccNumber string,
	) (bool, error)
	MockApprovePinResetServiceRequestFn func(
		ctx context.Context,
		clientID string,
		serviceRequestID string,
		cccNumber string,
		phoneNumber string,
		physicalIdentityVerified bool,
	) (bool, error)
	MockGetPendingServiceRequestsCountFn    func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	MockGetServiceRequestsFn                func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error)
	MockResolveServiceRequestFn             func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	MockSetInProgressByFn                   func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetServiceRequestsForKenyaEMRFn     func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error)
	MockUpdateServiceRequestsFromKenyaEMRFn func(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error)
	MockCreatePinResetServiceRequestFn      func(ctx context.Context, cccNumber string) (bool, error)
}

// NewServiceRequestUseCaseMock initializes a new service request instance mock
func NewServiceRequestUseCaseMock() *ServiceRequestUseCaseMock {
	return &ServiceRequestUseCaseMock{
		MockCreateServiceRequestFn: func(
			ctx context.Context,
			clientID string,
			requestType, request, cccNumber string,
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
		MockUpdateServiceRequestsFromKenyaEMRFn: func(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error) {
			return true, nil
		},
		MockGetServiceRequestsForKenyaEMRFn: func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
			currentTime := time.Now()
			staffID := uuid.New().String()
			facilityID := uuid.New().String()
			contact := "123454323"
			serviceReq := &domain.ServiceRequest{
				ID:            staffID,
				RequestType:   "SERVICE_REQUEST",
				Request:       "SERVICE_REQUEST",
				Status:        "PENDING",
				ClientID:      uuid.New().String(),
				InProgressAt:  &currentTime,
				InProgressBy:  &staffID,
				ResolvedAt:    &currentTime,
				ResolvedBy:    &staffID,
				FacilityID:    &facilityID,
				ClientName:    &staffID,
				ClientContact: &contact,
			}
			return []*domain.ServiceRequest{serviceReq}, nil
		},
		MockCreatePinResetServiceRequestFn: func(ctx context.Context, cccNumber string) (bool, error) {
			return true, nil
		},
		MockApprovePinResetServiceRequestFn: func(
			ctx context.Context,
			clientID string,
			serviceRequestID string,
			cccNumber string,
			phoneNumber string,
			physicalIdentityVerified bool,
		) (bool, error) {
			return true, nil
		},
	}
}

// CreateServiceRequest mocks the implementation for creating a service request
func (s *ServiceRequestUseCaseMock) CreateServiceRequest(
	ctx context.Context,
	clientID string,
	requestType, request, cccNumber string,
) (bool, error) {
	return s.MockCreateServiceRequestFn(ctx, clientID, requestType, request, cccNumber)
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

// GetServiceRequestsForKenyaEMR mocks getting service requests attached to a specific facility to be used by KenyaEMR
func (s *ServiceRequestUseCaseMock) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
	return s.MockGetServiceRequestsForKenyaEMRFn(ctx, payload)
}

// UpdateServiceRequestsFromKenyaEMR mocks the implementation of updating service requests from KenyaEMR to MyCareHub
func (s *ServiceRequestUseCaseMock) UpdateServiceRequestsFromKenyaEMR(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error) {
	return s.MockUpdateServiceRequestsFromKenyaEMRFn(ctx, payload)
}

// CreatePinResetServiceRequest mocks the implementation of creating a pin reset service request
func (s *ServiceRequestUseCaseMock) CreatePinResetServiceRequest(ctx context.Context, cccNumber string) (bool, error) {
	return s.MockCreatePinResetServiceRequestFn(ctx, cccNumber)
}

// ApprovePinResetServiceRequest mocks the implementation of approving a pin reset service request
func (s *ServiceRequestUseCaseMock) ApprovePinResetServiceRequest(
	ctx context.Context,
	clientID string,
	serviceRequestID string,
	cccNumber string,
	phoneNumber string,
	physicalIdentityVerified bool,
) (bool, error) {
	return s.MockApprovePinResetServiceRequestFn(ctx, clientID, serviceRequestID, cccNumber, phoneNumber, physicalIdentityVerified)
}
