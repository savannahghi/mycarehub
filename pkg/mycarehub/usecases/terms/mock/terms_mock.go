package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// TermsUseCaseMock mocks the implementation of terms usecase methods.
type TermsUseCaseMock struct {
	MockGetCurrentTermsFn func(ctx context.Context, flavour enums.Flavour) (string, error)
}

// NewTermsUseCaseMock creates in itializes create type mocks
func NewTermsUseCaseMock() *TermsUseCaseMock {
	return &TermsUseCaseMock{
		MockGetCurrentTermsFn: func(ctx context.Context, flavour enums.Flavour) (string, error) {
			terms := "terms"
			return terms, nil
		},
	}
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *TermsUseCaseMock) GetCurrentTerms(ctx context.Context, flavour enums.Flavour) (string, error) {
	return gm.MockGetCurrentTermsFn(ctx, flavour)
}
