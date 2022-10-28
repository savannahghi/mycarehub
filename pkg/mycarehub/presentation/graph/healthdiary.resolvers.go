package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateHealthDiaryEntry is the resolver for the createHealthDiaryEntry field.
func (r *mutationResolver) CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.CreateHealthDiaryEntry(ctx, clientID, note, mood, reportToStaff)
}

// ShareHealthDiaryEntry is the resolver for the shareHealthDiaryEntry field.
func (r *mutationResolver) ShareHealthDiaryEntry(ctx context.Context, healthDiaryEntryID string, shareEntireHealthDiary bool) (bool, error) {
	return r.mycarehub.HealthDiary.ShareHealthDiaryEntry(ctx, healthDiaryEntryID, shareEntireHealthDiary)
}

// CanRecordMood is the resolver for the canRecordMood field.
func (r *queryResolver) CanRecordMood(ctx context.Context, clientID string) (bool, error) {
	return r.mycarehub.HealthDiary.CanRecordHeathDiary(ctx, clientID)
}

// GetHealthDiaryQuote is the resolver for the getHealthDiaryQuote field.
func (r *queryResolver) GetHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.GetClientHealthDiaryQuote(ctx, limit)
}

// GetClientHealthDiaryEntries is the resolver for the getClientHealthDiaryEntries field.
func (r *queryResolver) GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
	r.checkPreconditions()
	return r.mycarehub.HealthDiary.GetClientHealthDiaryEntries(ctx, clientID, moodType, shared)
}

// GetSharedHealthDiaryEntries is the resolver for the getSharedHealthDiaryEntries field.
func (r *queryResolver) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return r.mycarehub.HealthDiary.GetSharedHealthDiaryEntries(ctx, clientID, facilityID)
}
