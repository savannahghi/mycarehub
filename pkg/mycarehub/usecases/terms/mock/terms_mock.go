package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// TermsUseCaseMock mocks the implementation of terms usecase methods.
type TermsUseCaseMock struct {
	MockGetCurrentTermsFn func(ctx context.Context) (*domain.TermsOfService, error)
	MockAcceptTermsFn     func(ctx context.Context, userID *string, termsID *int) (bool, error)
}

// NewTermsUseCaseMock creates in itializes create type mocks
func NewTermsUseCaseMock() *TermsUseCaseMock {
	return &TermsUseCaseMock{
		MockGetCurrentTermsFn: func(ctx context.Context) (*domain.TermsOfService, error) {
			termsID := gofakeit.Number(1, 1000)
			testText := "test"
			terms := &domain.TermsOfService{
				TermsID: termsID,
				Text:    &testText,
			}
			return terms, nil
		},
		MockAcceptTermsFn: func(ctx context.Context, userID *string, termsID *int) (bool, error) {
			return true, nil
		},
	}
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *TermsUseCaseMock) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}

// AcceptTerms mocks the implementation of accept current terms of service
func (gm *TermsUseCaseMock) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	return gm.MockAcceptTermsFn(ctx, userID, termsID)
}
