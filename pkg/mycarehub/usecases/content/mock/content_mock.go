package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ContentUsecaseMock contains the mock of contentusecase methods
type ContentUsecaseMock struct {
	MockListContentCategoriesFn func(ctx context.Context) ([]*domain.ContentItemCategory, error)
	MockShareContentFn          func(ctx context.Context, input dto.ShareContentInput) (bool, error)
}

// NewContentUsecaseMock instantiates all the content usecase mock methods
func NewContentUsecaseMock() *ContentUsecaseMock {
	contentItemCategory := &domain.ContentItemCategory{
		ID:      1,
		Name:    "name",
		IconURL: "test",
	}

	return &ContentUsecaseMock{
		MockListContentCategoriesFn: func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
			return []*domain.ContentItemCategory{contentItemCategory}, nil
		},
		MockShareContentFn: func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
			return true, nil
		},
	}
}

//ListContentCategories mocks the implementation listing content categories
func (cm *ContentUsecaseMock) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return cm.MockListContentCategoriesFn(ctx)
}

// ShareContent mocks the implementation of `gorm's` ShareContent method.
func (cm *ContentUsecaseMock) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return cm.MockShareContentFn(ctx, input)
}
