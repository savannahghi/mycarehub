package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.26

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SetInProgressBy is the resolver for the setInProgressBy field.
func (r *mutationResolver) SetInProgressBy(ctx context.Context, serviceRequestID string, staffID string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.SetInProgressBy(ctx, serviceRequestID, staffID)
}

// CreateServiceRequest is the resolver for the createServiceRequest field.
func (r *mutationResolver) CreateServiceRequest(ctx context.Context, input dto.ServiceRequestInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.CreateServiceRequest(ctx, &input)
}

// ResolveServiceRequest is the resolver for the resolveServiceRequest field.
func (r *mutationResolver) ResolveServiceRequest(ctx context.Context, staffID string, requestID string, action []string, comment *string) (bool, error) {
	return r.mycarehub.ServiceRequest.ResolveServiceRequest(ctx, &staffID, &requestID, action, comment)
}

// VerifyClientPinResetServiceRequest is the resolver for the verifyClientPinResetServiceRequest field.
func (r *mutationResolver) VerifyClientPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyClientPinResetServiceRequest(ctx, serviceRequestID, status, physicalIdentityVerified)
}

// VerifyStaffPinResetServiceRequest is the resolver for the verifyStaffPinResetServiceRequest field.
func (r *mutationResolver) VerifyStaffPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error) {
	return r.mycarehub.ServiceRequest.VerifyStaffPinResetServiceRequest(ctx, serviceRequestID, status)
}

// GetServiceRequests is the resolver for the getServiceRequests field.
func (r *queryResolver) GetServiceRequests(ctx context.Context, requestType *string, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	return r.mycarehub.ServiceRequest.GetServiceRequests(ctx, requestType, requestStatus, facilityID, flavour)
}

// GetPendingServiceRequestsCount is the resolver for the getPendingServiceRequestsCount field.
func (r *queryResolver) GetPendingServiceRequestsCount(ctx context.Context) (*domain.ServiceRequestsCountResponse, error) {
	return r.mycarehub.ServiceRequest.GetPendingServiceRequestsCount(ctx)
}

// SearchServiceRequests is the resolver for the searchServiceRequests field.
func (r *queryResolver) SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	return r.mycarehub.ServiceRequest.SearchServiceRequests(ctx, searchTerm, flavour, requestType, facilityID)
}
