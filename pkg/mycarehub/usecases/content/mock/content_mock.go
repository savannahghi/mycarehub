package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ContentUsecaseMock contains the mock of contentusecase methods
type ContentUsecaseMock struct {
	MockListContentCategoriesFn           func(ctx context.Context) ([]*domain.ContentItemCategory, error)
	MockShareContentFn                    func(ctx context.Context, input dto.ShareContentInput) (bool, error)
	MockGetContentFn                      func(ctx context.Context, categoryIDs []int, categoryNames []string, limit string, clientID *string) (*domain.Content, error)
	MockGetUserBookmarkedContentFn        func(ctx context.Context, userID string) (*domain.Content, error)
	MockGetContentItemByIDFn              func(ctx context.Context, contentID int) (*domain.ContentItem, error)
	MockLikeContentFn                     func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCheckWhetherUserHasLikedContentFn func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnlikeContentFn                   func(ctx context.Context, userID string, contentID int) (bool, error)
	MockGetFAQsFn                         func(ctx context.Context, flavor feedlib.Flavour) (*domain.Content, error)
	MockBookmarkContentFn                 func(ctx context.Context, userID string, contentID int) (bool, error)
	MockUnBookmarkContentFn               func(ctx context.Context, userID string, contentID int) (bool, error)
	MockCheckIfUserBookmarkedContentFn    func(ctx context.Context, userID string, contentID int) (bool, error)
	MockViewContentFn                     func(ctx context.Context, userID string, contentID int) (bool, error)
}

// NewContentUsecaseMock instantiates all the content usecase mock methods
func NewContentUsecaseMock() *ContentUsecaseMock {
	contentItemCategory := &domain.ContentItemCategory{
		ID:      1,
		Name:    "name",
		IconURL: "test",
	}

	now := time.Now()

	content := &domain.Content{
		Meta: domain.Meta{},
		Items: []domain.ContentItem{
			{
				ID:    1,
				Meta:  domain.ContentMeta{},
				Title: gofakeit.Name(),
				Date:  now.String(),
				Intro: gofakeit.Sentence(2),
				Author: domain.Author{
					ID: uuid.New().String(),
					Meta: domain.AuthorMeta{
						Type: gofakeit.Name(),
					},
				},
				AuthorName:          gofakeit.Name(),
				ItemType:            gofakeit.Name(),
				TimeEstimateSeconds: 30,
				Body:                gofakeit.Name(),
				TagNames:            []string{},
				HeroImage:           domain.HeroImage{},
				HeroImageRendition:  domain.HeroImageRendition{},
				LikeCount:           0,
				BookmarkCount:       0,
				ViewCount:           0,
				ShareCount:          0,
				Documents:           []domain.Document{},
				CategoryDetails:     []domain.CategoryDetail{},
				FeaturedMedia:       []domain.FeaturedMedia{},
				GalleryImages:       []domain.GalleryImage{},
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
		MockGetContentFn: func(ctx context.Context, categoryIDs []int, categoryNames []string, limit string, clientID *string) (*domain.Content, error) {
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
		MockGetContentItemByIDFn: func(ctx context.Context, contentID int) (*domain.ContentItem, error) {
			return &content.Items[0], nil
		},
		MockLikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockCheckWhetherUserHasLikedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnlikeContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockGetFAQsFn: func(ctx context.Context, flavor feedlib.Flavour) (*domain.Content, error) {
			return content, nil
		},
		MockBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockUnBookmarkContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockCheckIfUserBookmarkedContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
		MockViewContentFn: func(ctx context.Context, userID string, contentID int) (bool, error) {
			return true, nil
		},
	}
}

// ListContentCategories mocks the implementation listing content categories
func (cm *ContentUsecaseMock) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return cm.MockListContentCategoriesFn(ctx)
}

// ShareContent mocks the implementation of `gorm's` ShareContent method.
func (cm *ContentUsecaseMock) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return cm.MockShareContentFn(ctx, input)
}

// GetContent mocks the implementation of making an API call to fetch content from our APIs
func (cm *ContentUsecaseMock) GetContent(ctx context.Context, categoryIDs []int, categoryNames []string, limit string, clientID *string) (*domain.Content, error) {
	return cm.MockGetContentFn(ctx, categoryIDs, categoryNames, limit, clientID)
}

// GetUserBookmarkedContent mocks the implementation of getting a users bookmarked content
func (cm *ContentUsecaseMock) GetUserBookmarkedContent(ctx context.Context, userID string) (*domain.Content, error) {
	return cm.MockGetUserBookmarkedContentFn(ctx, userID)
}

// GetContentItemByID mocks fetching content using it's item ID
func (cm *ContentUsecaseMock) GetContentItemByID(ctx context.Context, contentID int) (*domain.ContentItem, error) {
	return cm.MockGetContentItemByIDFn(ctx, contentID)
}

// LikeContent mocks the implementation liking a feed content
func (cm *ContentUsecaseMock) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockLikeContentFn(ctx, userID, contentID)
}

// CheckWhetherUserHasLikedContent mocks the implementation of `gorm's` CheckWhetherUserHasLikedContent method.
func (cm *ContentUsecaseMock) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {

	return cm.MockCheckWhetherUserHasLikedContentFn(ctx, userID, contentID)
}

// UnlikeContent mocks the implementation liking a feed content
func (cm *ContentUsecaseMock) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockUnlikeContentFn(ctx, userID, contentID)
}

// GetFAQs mocks the implementation of getting FAQs
func (cm *ContentUsecaseMock) GetFAQs(ctx context.Context, flavor feedlib.Flavour) (*domain.Content, error) {
	return cm.MockGetFAQsFn(ctx, flavor)
}

// BookmarkContent mock the implementation of the BookmarkContent method
func (cm *ContentUsecaseMock) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockBookmarkContentFn(ctx, userID, contentID)
}

// UnBookmarkContent mock the implementation of the UnBookmarkContent method
func (cm *ContentUsecaseMock) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockUnBookmarkContentFn(ctx, userID, contentID)
}

// CheckIfUserBookmarkedContent mock the implementation of the CheckIfUserBookmarkedContent method
func (cm *ContentUsecaseMock) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockCheckIfUserBookmarkedContentFn(ctx, userID, contentID)
}

// ViewContent mock the implementation of the ViewContent method
func (cm *ContentUsecaseMock) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return cm.MockViewContentFn(ctx, userID, contentID)
}
