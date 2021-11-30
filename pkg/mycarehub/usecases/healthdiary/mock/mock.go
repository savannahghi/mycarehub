package mock

import "context"

// HealthDiaryUseCaseMock mocks the implementation of HealthDiary usecase
type HealthDiaryUseCaseMock struct {
	MockCreateHealthDiaryEntryFn func(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error)
	MockCanRecordHeathDiaryFn    func(ctx context.Context, clientID string) (bool, error)
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
