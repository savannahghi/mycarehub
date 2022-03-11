package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *queryResolver) FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filterInput []*dto.FiltersInput) (*domain.AppointmentsPage, error) {
	return r.mycarehub.Appointment.FetchClientAppointments(ctx, clientID, paginationInput, filterInput)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
