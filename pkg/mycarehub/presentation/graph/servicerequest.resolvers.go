package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) SetInProgressBy(ctx context.Context, serviceRequestID string, staffID string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.SetInProgressBy(ctx, serviceRequestID, staffID)
}

func (r *mutationResolver) CreateServiceRequest(ctx context.Context, input dto.ServiceRequestInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.CreateServiceRequest(ctx, &input)
}

func (r *mutationResolver) ResolveServiceRequest(ctx context.Context, staffID string, requestID string, action string, comment *string) (bool, error) {
	return r.mycarehub.ServiceRequest.ResolveServiceRequest(ctx, &staffID, &requestID, action, comment)
}

func (r *mutationResolver) VerifyClientPinResetServiceRequest(ctx context.Context, clientID string, serviceRequestID string, cccNumber string, phoneNumber string, physicalIdentityVerified bool, state string) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyClientPinResetServiceRequest(ctx, clientID, serviceRequestID, cccNumber, phoneNumber, physicalIdentityVerified, state)
}

func (r *mutationResolver) VerifyStaffPinResetServiceRequest(ctx context.Context, phoneNumber string, serviceRequestID string, verificationStatus string) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyStaffPinResetServiceRequest(ctx, phoneNumber, serviceRequestID, verificationStatus)
}

func (r *queryResolver) GetServiceRequests(ctx context.Context, requestType *string, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	return r.mycarehub.ServiceRequest.GetServiceRequests(ctx, requestType, requestStatus, facilityID, flavour)
}

func (r *queryResolver) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
	return r.mycarehub.ServiceRequest.GetPendingServiceRequestsCount(ctx, facilityID)
}
