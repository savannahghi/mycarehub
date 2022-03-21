package mock

import (
	"context"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// HealthDiaryUseCaseMock mocks the implementation of HealthDiary usecase
type HealthDiaryUseCaseMock struct {
	MockCreateHealthDiaryEntryFn        func(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error)
	MockCanRecordHeathDiaryFn           func(ctx context.Context, clientID string) (bool, error)
	MockGetClientHealthDiaryQuoteFn     func(ctx context.Context) (*domain.ClientHealthDiaryQuote, error)
	MockGetClientHealthDiaryEntriesFn   func(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockGetFacilityHealthDiaryEntriesFn func(ctx context.Context, input dto.FetchHealthDiaryEntries) (*dto.HealthDiaryEntriesResponse, error)
	MockGetRecentHealthDiaryEntriesFn   func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
	MockShareHealthDiaryEntryFn         func(ctx context.Context, clientID string, shareWithStaff bool) (bool, error)
}

// NewHealthDiaryUseCaseMock initializes a new instance mock of the HealthDiary usecase
func NewHealthDiaryUseCaseMock() *HealthDiaryUseCaseMock {
	return &HealthDiaryUseCaseMock{
		MockCreateHealthDiaryEntryFn: func(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error) {
			return true, nil
		},
		MockCanRecordHeathDiaryFn: func(ctx context.Context, clientID string) (bool, error) {
			return true, nil
		},
		MockGetClientHealthDiaryQuoteFn: func(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
			return &domain.ClientHealthDiaryQuote{
				Author: "test",
				Quote:  "test",
			}, nil
		},
		MockGetClientHealthDiaryEntriesFn: func(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{
				{
					Active: true,
				},
			}, nil
		},
		MockGetFacilityHealthDiaryEntriesFn: func(ctx context.Context, input dto.FetchHealthDiaryEntries) (*dto.HealthDiaryEntriesResponse, error) {
			return &dto.HealthDiaryEntriesResponse{
				MFLCode: 1234,
				HealthDiaryEntries: []*domain.ClientHealthDiaryEntry{
					{
						Active: true,
					},
				},
			}, nil
		},
		MockShareHealthDiaryEntryFn: func(ctx context.Context, clientID string, shareWithStaff bool) (bool, error) {
			return true, nil
		},
		MockGetRecentHealthDiaryEntriesFn: func(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
			return []*domain.ClientHealthDiaryEntry{}, nil
		},
	}
}

// CreateHealthDiaryEntry mocks the method for creating a new health diary entry
func (h *HealthDiaryUseCaseMock) CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error) {
	return h.MockCreateHealthDiaryEntryFn(ctx, clientID, note, mood, reportToStaff)
}

// CanRecordHeathDiary implements check for eligibility of a health diary to be shown to a user
func (h *HealthDiaryUseCaseMock) CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error) {
	return h.MockCanRecordHeathDiaryFn(ctx, clientID)
}

// GetClientHealthDiaryQuote mocks the method for getting a random health diary quote
func (h *HealthDiaryUseCaseMock) GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
	return h.MockGetClientHealthDiaryQuoteFn(ctx)
}

// GetClientHealthDiaryEntries mocks the method for fetching a client's health record entries
func (h *HealthDiaryUseCaseMock) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return h.MockGetClientHealthDiaryEntriesFn(ctx, clientID)
}

// GetFacilityHealthDiaryEntries mocks getting health diary entries per facility
func (h *HealthDiaryUseCaseMock) GetFacilityHealthDiaryEntries(ctx context.Context, input dto.FetchHealthDiaryEntries) (*dto.HealthDiaryEntriesResponse, error) {
	return h.MockGetFacilityHealthDiaryEntriesFn(ctx, input)
}

// GetRecentHealthDiaryEntries mocks getting the most recent health diary entries
func (h *HealthDiaryUseCaseMock) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	return h.MockGetRecentHealthDiaryEntriesFn(ctx, lastSyncTime, clientID)
}

// ShareHealthDiaryEntry mocks the implementation of sharing a health diary entry when the client opts to share it with the health care worker
func (h *HealthDiaryUseCaseMock) ShareHealthDiaryEntry(ctx context.Context, clientID string, shareWithStaff bool) (bool, error) {
	return h.MockShareHealthDiaryEntryFn(ctx, clientID, shareWithStaff)
}
