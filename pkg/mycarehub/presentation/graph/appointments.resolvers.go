package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
	"github.com/savannahghi/scalarutils"
)

func (r *mutationResolver) RescheduleAppointment(ctx context.Context, appointmentID string, date scalarutils.Date) (bool, error) {
	return r.mycarehub.Appointment.RescheduleClientAppointment(ctx, appointmentID, date)
}

func (r *queryResolver) FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error) {
	return r.mycarehub.Appointment.FetchClientAppointments(ctx, clientID, paginationInput, filters)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
