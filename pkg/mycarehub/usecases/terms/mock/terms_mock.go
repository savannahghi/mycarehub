package mock

import (
	"context"
)

// TermsUseCaseMock mocks the implementation of terms usecase methods.
type TermsUseCaseMock struct {
	MockGetCurrentTermsFn func(ctx context.Context) (string, error)
}

// NewTermsUseCaseMock creates in itializes create type mocks
func NewTermsUseCaseMock() *TermsUseCaseMock {
	return &TermsUseCaseMock{
		MockGetCurrentTermsFn: func(ctx context.Context) (string, error) {
			terms := "terms"
			return terms, nil
		},
	}
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *TermsUseCaseMock) GetCurrentTerms(ctx context.Context) (string, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}
