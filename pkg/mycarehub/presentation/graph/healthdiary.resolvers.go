package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.CreateHealthDiaryEntry(ctx, clientID, note, mood, reportToStaff)
}
