package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ServiceRequestUseCaseMock mocks the service request instance
type ServiceRequestUseCaseMock struct {
	MockCreateServiceRequestFn               func(ctx context.Context, input *dto.ServiceRequestInput) (bool, error)
	MockVerifyClientPinResetServiceRequestFn func(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error)
	MockGetPendingServiceRequestsCountFn     func(ctx context.Context) (*domain.ServiceRequestsCountResponse, error)
	MockGetServiceRequestsFn                 func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour, paginationInput *dto.PaginationsInput) (*domain.ServiceRequestPage, error)
	MockResolveServiceRequestFn              func(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error)
	MockSetInProgressByFn                    func(ctx context.Context, requestID string, staffID string) (bool, error)
	MockGetServiceRequestsForKenyaEMRFn      func(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error)
	MockUpdateServiceRequestsFromKenyaEMRFn  func(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error)
	MockCreatePinResetServiceRequestFn       func(ctx context.Context, username string, cccNumber string, flavour feedlib.Flavour) (bool, error)
	MockVerifyStaffPinResetServiceRequestFn  func(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error)
	MockSearchServiceRequestsFn              func(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
	MockCompleteVisitFn                      func(ctx context.Context, staffID string, serviceRequestID string, bookingID string, notes string) (bool, error)
}

// NewServiceRequestUseCaseMock initializes a new service request instance mock
func NewServiceRequestUseCaseMock() *ServiceRequestUseCaseMock {
	return &ServiceRequestUseCaseMock{
		MockCreateServiceRequestFn: func(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
			return true, nil
		},
		MockGetPendingServiceRequestsCountFn: func(ctx context.Context) (*domain.ServiceRequestsCountResponse, error) {
			return &domain.ServiceRequestsCountResponse{
				ClientsServiceRequestCount: &domain.ServiceRequestsCount{
					Total: 0,
					RequestsTypeCount: []*domain.RequestTypeCount{
						{
							RequestType: "client",
							Total:       0,
						},
					},
				},
				StaffServiceRequestCount: &domain.ServiceRequestsCount{
					Total: 0,
					RequestsTypeCount: []*domain.RequestTypeCount{
						{
							RequestType: "staff",
							Total:       0,
						},
					},
				},
			}, nil
		},
		MockGetServiceRequestsFn: func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour, paginationInput *dto.PaginationsInput) (*domain.ServiceRequestPage, error) {
			return &domain.ServiceRequestPage{
				Results: []*domain.ServiceRequest{
					{
						ID: uuid.New().String(),
					},
				},
				Pagination: domain.Pagination{
					CurrentPage: 1,
					Limit:       10,
				},
			}, nil
		},
		MockResolveServiceRequestFn: func(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error) {
			return true, nil
		},
		MockSetInProgressByFn: func(ctx context.Context, requestID string, staffID string) (bool, error) {
			return true, nil
		},
		MockUpdateServiceRequestsFromKenyaEMRFn: func(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error) {
			return true, nil
		},
		MockGetServiceRequestsForKenyaEMRFn: func(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error) {
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
				FacilityID:    facilityID,
				ClientName:    &staffID,
				ClientContact: &contact,
				Meta:          map[string]interface{}{"meta": "data"},
			}
			return &dto.RedFlagServiceRequestResponse{
				RedFlagServiceRequests: []*domain.ServiceRequest{serviceReq},
			}, nil
		},
		MockCreatePinResetServiceRequestFn: func(ctx context.Context, username string, cccNumber string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockVerifyStaffPinResetServiceRequestFn: func(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error) {
			return true, nil
		},
		MockVerifyClientPinResetServiceRequestFn: func(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error) {
			return true, nil
		},
		MockSearchServiceRequestsFn: func(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
			return []*domain.ServiceRequest{
				{
					ID:          uuid.New().String(),
					RequestType: "RED_FLAG",
					Request:     "TEST",
					Status:      "PENDING",
					Active:      true,
					ClientID:    "",
					StaffID:     "",
					FacilityID:  uuid.NewString(),
					Meta: map[string]interface{}{
						"meta": "data",
					},
				},
			}, nil
		},
		MockCompleteVisitFn: func(ctx context.Context, staffID, serviceRequestID, bookingID string, notes string) (bool, error) {
			return true, nil
		},
	}
}

// CreateServiceRequest mocks the implementation for creating a service request
func (s *ServiceRequestUseCaseMock) CreateServiceRequest(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
	return s.MockCreateServiceRequestFn(ctx, input)
}

// GetPendingServiceRequestsCount mocks the method of getting the number of pending service requests count
func (s *ServiceRequestUseCaseMock) GetPendingServiceRequestsCount(ctx context.Context) (*domain.ServiceRequestsCountResponse, error) {
	return s.MockGetPendingServiceRequestsCountFn(ctx)
}

// VerifyStaffPinResetServiceRequest mocks the implementation of getting the number of staff pending service requests count
func (s *ServiceRequestUseCaseMock) VerifyStaffPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error) {
	return s.MockVerifyStaffPinResetServiceRequestFn(ctx, serviceRequestID, status)
}

// GetServiceRequests mocks the method for fetching service requests
func (s *ServiceRequestUseCaseMock) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour, paginationInput *dto.PaginationsInput) (*domain.ServiceRequestPage, error) {
	return s.MockGetServiceRequestsFn(ctx, requestType, requestStatus, facilityID, flavour, paginationInput)
}

// ResolveServiceRequest mocks resolving a service request
func (s *ServiceRequestUseCaseMock) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error) {
	return s.MockResolveServiceRequestFn(ctx, staffID, serviceRequestID, action, comment)
}

// SetInProgressBy mocks the implementation of marking a service request as in progress
func (s *ServiceRequestUseCaseMock) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	return s.MockSetInProgressByFn(ctx, requestID, staffID)
}

// GetServiceRequestsForKenyaEMR mocks getting service requests attached to a specific facility to be used by KenyaEMR
func (s *ServiceRequestUseCaseMock) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error) {
	return s.MockGetServiceRequestsForKenyaEMRFn(ctx, payload)
}

// UpdateServiceRequestsFromKenyaEMR mocks the implementation of updating service requests from KenyaEMR to MyCareHub
func (s *ServiceRequestUseCaseMock) UpdateServiceRequestsFromKenyaEMR(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error) {
	return s.MockUpdateServiceRequestsFromKenyaEMRFn(ctx, payload)
}

// CreatePinResetServiceRequest mocks the implementation of creating a pin reset service request
func (s *ServiceRequestUseCaseMock) CreatePinResetServiceRequest(ctx context.Context, username string, cccNumber string, flavour feedlib.Flavour) (bool, error) {
	return s.MockCreatePinResetServiceRequestFn(ctx, username, cccNumber, flavour)
}

// VerifyClientPinResetServiceRequest mocks the implementation of approving a pin reset service request
func (s *ServiceRequestUseCaseMock) VerifyClientPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error) {
	return s.MockVerifyClientPinResetServiceRequestFn(ctx, serviceRequestID, status, physicalIdentityVerified)
}

// SearchServiceRequests mocks the implementation of searching service requests
func (s *ServiceRequestUseCaseMock) SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	return s.MockSearchServiceRequestsFn(ctx, searchTerm, flavour, requestType, facilityID)
}

// CompleteVisit mocks the implementation of fulfilling a booking
func (s *ServiceRequestUseCaseMock) CompleteVisit(ctx context.Context, staffID string, serviceRequestID string, bookingID string, notes string) (bool, error) {
	return s.MockCompleteVisitFn(ctx, staffID, serviceRequestID, bookingID, notes)
}
