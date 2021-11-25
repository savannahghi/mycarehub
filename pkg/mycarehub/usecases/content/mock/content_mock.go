package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ContentUsecaseMock contains the mock of contentusecase methods
type ContentUsecaseMock struct {
	MockListContentCategoriesFn     func(ctx context.Context) ([]*domain.ContentItemCategory, error)
	MockShareContentFn              func(ctx context.Context, input dto.ShareContentInput) (bool, error)
	MockGetContentFn                func(ctx context.Context, categoryID *int, limit string) (*domain.Content, error)
	MockGetUserBookmarkedContentFn  func(ctx context.Context, userID string) (*domain.Content, error)
	MockGetContentByContentItemIDFn func(ctx context.Context, contentID int) (*domain.Content, error)
}

// NewContentUsecaseMock instantiates all the content usecase mock methods
func NewContentUsecaseMock() *ContentUsecaseMock {
	contentItemCategory := &domain.ContentItemCategory{
		ID:      1,
		Name:    "name",
		IconURL: "test",
	}

	content := &domain.Content{
		Items: []domain.ContentItem{
			{
				ID:    int(uuid.New()[9]),
				Title: gofakeit.Name(),
			},
		},
	}

	return &ContentUsecaseMock{
		MockListContentCategoriesFn: func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
			return []*domain.ContentItemCategory{contentItemCategory}, nil
		},
		MockShareContentFn: func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
			return true, nil
		},
		MockGetContentFn: func(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
			return &domain.Content{
				Items: []domain.ContentItem{
					{
						ID: int(uuid.New()[9]),
					},
				},
			}, nil
		},
		MockGetUserBookmarkedContentFn: func(ctx context.Context, userID string) (*domain.Content, error) {
			return content, nil
		},
		MockGetContentByContentItemIDFn: func(ctx context.Context, contentID int) (*domain.Content, error) {
			return content, nil
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

// GetContent mocks the implementation of making an API call to fetch content from our APIs
func (cm *ContentUsecaseMock) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	return cm.MockGetContentFn(ctx, categoryID, limit)
}

// GetUserBookmarkedContent mocks the implementation of getting a users bookmarked content
func (cm *ContentUsecaseMock) GetUserBookmarkedContent(ctx context.Context, userID string) (*domain.Content, error) {
	return cm.MockGetUserBookmarkedContentFn(ctx, userID)
}

// GetContentByContentItemID mocks fetching content using it's item ID
func (cm *ContentUsecaseMock) GetContentByContentItemID(ctx context.Context, contentID int) (*domain.Content, error) {
	return cm.MockGetContentByContentItemIDFn(ctx, contentID)
}
