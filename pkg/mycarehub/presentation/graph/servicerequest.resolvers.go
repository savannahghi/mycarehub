package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) CreateServiceRequest(ctx context.Context, clientID string, requestType string, request *string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.ServiceRequest.CreateServiceRequest(ctx, clientID, requestType, *request)
}
