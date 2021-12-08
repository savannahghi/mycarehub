package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// FAQUsecaseMock contains the mock of FAQ usecase methods
type FAQUsecaseMock struct {
	MockGetFAQContentFn func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error)
}

// NewFAQUsecaseMock instantiates all the FAQ usecase mock methods
func NewFAQUsecaseMock() *FAQUsecaseMock {

	return &FAQUsecaseMock{

		MockGetFAQContentFn: func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
			ID := uuid.New().String()
			return []*domain.FAQ{
				{
					ID:          &ID,
					Active:      true,
					Title:       gofakeit.Name(),
					Description: gofakeit.Name(),
					Body:        gofakeit.Name(),
				},
			}, nil
		},
	}
}

// GetFAQContent mock method FAQ usecase
func (f *FAQUsecaseMock) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
	return f.MockGetFAQContentFn(ctx, flavour, limit)
}
