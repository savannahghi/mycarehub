package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/segmentio/ksuid"
)

// TermsUseCaseMock mocks the implementation of terms usecase methods.
type TermsUseCaseMock struct {
	MockGetCurrentTermsFn func(ctx context.Context) (*domain.TermsOfService, error)
}

// NewTermsUseCaseMock creates in itializes create type mocks
func NewTermsUseCaseMock() *TermsUseCaseMock {
	return &TermsUseCaseMock{
		MockGetCurrentTermsFn: func(ctx context.Context) (*domain.TermsOfService, error) {
			termsID := ksuid.New().String()
			testText := "test"
			terms := &domain.TermsOfService{
				TermsID: &termsID,
				Text:    &testText,
			}
			return terms, nil
		},
	}
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *TermsUseCaseMock) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}
