package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.CreateHealthDiaryEntry(ctx, clientID, note, mood, reportToStaff)
}

func (r *mutationResolver) ShareHealthDiaryEntry(ctx context.Context, healthDiaryEntryID string, shareEntireHealthDiary bool) (bool, error) {
	return r.mycarehub.HealthDiary.ShareHealthDiaryEntry(ctx, healthDiaryEntryID, shareEntireHealthDiary)
}

func (r *queryResolver) CanRecordMood(ctx context.Context, clientID string) (bool, error) {
	return r.mycarehub.HealthDiary.CanRecordHeathDiary(ctx, clientID)
}

func (r *queryResolver) GetHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.GetClientHealthDiaryQuote(ctx, limit)
}

func (r *queryResolver) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.GetClientHealthDiaryEntries(ctx, clientID)
}

func (r *queryResolver) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return r.mycarehub.HealthDiary.GetSharedHealthDiaryEntries(ctx, clientID, facilityID)
}
