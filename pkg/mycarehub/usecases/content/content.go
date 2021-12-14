package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/serverutils"
)

var (
	contentAPIEndpoint = serverutils.MustGetEnvVar("CONTENT_API_URL")
)

// IGetContent is used to fetch content from the CMS
type IGetContent interface {
	GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error)
	GetContentByContentItemID(ctx context.Context, contentID int) (*domain.Content, error)
}

// IGetBookmarkedContent holds the method signature used to return a user's bookmarked content
type IGetBookmarkedContent interface {
	GetUserBookmarkedContent(ctx context.Context, userID string) (*domain.Content, error)
}

// ICheckIfUserBookmarkedContent is used to check if a user has bookmarked a content item
type ICheckIfUserBookmarkedContent interface {
	CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// IContentCategoryList groups all the content category listing methods
type IContentCategoryList interface {
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
}

// IShareContent is the interface for the ShareContent
type IShareContent interface {
	// TODO: update share count (increment)
	// TODO: add / check entry in ContentShares table
	// TODO: metrics
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
}

// IBookmarkContent is used to bookmark content
type IBookmarkContent interface {
	// TODO: update bookmark count (increment)
	// TODO: idempotence, with user ID i.e a user can only bookmark once
	// TODO: add / check entry in ContentBookmarks table
	// TODO: metrics
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// IUnBookmarkContent is used to unbookmark content
type IUnBookmarkContent interface {
	// TODO: update bookmark count (decrement)
	// TODO: idempotence, with user ID i.e a user can only remove something they bookmarked
	// TODO: remove entry from ContentBookmarks table if it exists...be forgiving (idempotence)
	// TODO: metrics
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// ILikeContent groups the like feature methods
type ILikeContent interface {
	// TODO: update like count (increment)
	// TODO: idempotence, with user ID i.e a user can only like once
	// TODO: add / check entry in ContentLikes table
	// TODO: metrics
	LikeContent(ctx context.Context, userID string, contentID int) (bool, error)
	CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// IUnlikeContent groups the unllike feature methods
type IUnlikeContent interface {
	// TODO: update like count (decrement)
	// TODO: idempotence, with user ID i.e a user can only unlike something they liked
	// TODO: remove entry from ContentLikes table if it exists...be forgiving (idempotence)
	// TODO: metrics
	UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// UseCasesContent holds the interfaces that are implemented within the content service
type UseCasesContent interface {
	IGetContent
	IContentCategoryList
	IGetBookmarkedContent
	IShareContent
	IBookmarkContent
	IUnBookmarkContent
	ILikeContent
	IUnlikeContent
	IViewContent
	ICheckIfUserBookmarkedContent
}

// IViewContent gets a content ite and updates the view count
type IViewContent interface {
	// TODO Update view metrics each time a user views a piece
	// TODO Increment view count, idempotent
	ViewContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// UseCasesContentImpl represents content implementation
type UseCasesContentImpl struct {
	Update infrastructure.Update
	Query  infrastructure.Query
}

// NewUseCasesContentImplementation initializes a new contents service
func NewUseCasesContentImplementation(
	update infrastructure.Update,
	query infrastructure.Query,
) *UseCasesContentImpl {
	return &UseCasesContentImpl{
		Update: update,
		Query:  query,
	}

}

// LikeContent implements the content liking api
func (u UseCasesContentImpl) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.LikeContent(ctx, userID, contentID)
}

// CheckWhetherUserHasLikedContent implements action of checking whether a user has liked a particular content
func (u UseCasesContentImpl) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Query.CheckWhetherUserHasLikedContent(ctx, userID, contentID)
}

// UnlikeContent implements the content liking api
func (u UseCasesContentImpl) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.UnlikeContent(ctx, userID, contentID)
}

// GetContent fetches content from wagtail CMS. The category ID is optional and it is used to return content based
// on the category it belongs to. The limit field describes how many items will be rendered on the front end side.
func (u UseCasesContentImpl) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	params := url.Values{}
	params.Add("type", "content.ContentItem")
	params.Add("limit", limit)
	params.Add("order", "-first_published_at")
	params.Add("fields", "'*")
	if categoryID != nil {
		params.Add("category", strconv.Itoa(*categoryID))
	}
	categoryIDs, err := u.Query.GetContentItemCategoryID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category IDs: %v", err)
	}
	if categoryID == nil {
		for _, v := range categoryIDs {
			params.Add("category", fmt.Sprintf("%v", v))
		}
	}

	getContentEndpoint := fmt.Sprintf(contentAPIEndpoint + "/?" + params.Encode())
	var contentItems *domain.Content
	resp, err := utils.MakeRequest(ctx, http.MethodGet, getContentEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request")
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	err = json.Unmarshal(dataResponse, &contentItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contentItems, nil
}

// ListContentCategories gets the list of all content categories
func (u *UseCasesContentImpl) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return u.Query.ListContentCategories(ctx)
}

// ShareContent enables a user to share a content
func (u *UseCasesContentImpl) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return u.Update.ShareContent(ctx, input)
}

// BookmarkContent increments the bookmark count for a content item
func (u UseCasesContentImpl) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.BookmarkContent(ctx, userID, contentID)
}

// UnBookmarkContent decrements the bookmark count for a content item
func (u UseCasesContentImpl) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.UnBookmarkContent(ctx, userID, contentID)
}

// GetUserBookmarkedContent gets the user's pinned/bookmarked content and displays it on their profile
func (u *UseCasesContentImpl) GetUserBookmarkedContent(ctx context.Context, userID string) (*domain.Content, error) {
	if userID == "" {
		return nil, exceptions.EmptyInputErr(fmt.Errorf("user ID must be defined"))
	}

	user, err := u.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, exceptions.ProfileNotFoundErr(err)
	}

	content, err := u.Query.GetUserBookmarkedContent(ctx, *user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarked content")
	}

	bookmarkedContent := &domain.Content{}

	items := make([]domain.ContentItem, len(content))

	for i, contentItem := range content {

		bookmarkedContent, err = u.GetContentByContentItemID(ctx, contentItem.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bookmarked content")
		}
		items[i] = *contentItem
	}

	bookmarkedContent.Meta.TotalCount = len(items)
	bookmarkedContent.Items = items

	return bookmarkedContent, nil
}

// GetContentByContentItemID fetches a specific content using the specific content item ID. This will be important
// when fetching content that a user bookmarked. The data returned directly from the database does not contain all the
// information regarding a content item hence why this method has been chosen.
func (u *UseCasesContentImpl) GetContentByContentItemID(ctx context.Context, contentID int) (*domain.Content, error) {
	params := url.Values{}
	params.Add("type", "content.ContentItem")
	params.Add("id", strconv.Itoa(contentID))
	params.Add("order", "-first_published_at")
	params.Add("fields", "'*")

	getContentEndpoint := fmt.Sprintf(contentAPIEndpoint + "/?" + params.Encode())
	var contentItems *domain.Content
	resp, err := utils.MakeRequest(ctx, http.MethodGet, getContentEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request")
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	err = json.Unmarshal(dataResponse, &contentItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contentItems, nil
}

// ViewContent gets a content item and updates the view count
func (u *UseCasesContentImpl) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID and contentID cannot be empty")
	}
	return u.Update.ViewContent(ctx, userID, contentID)
}

// CheckIfUserBookmarkedContent checks if a user has bookmarked a specific content item
func (u *UseCasesContentImpl) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID and contentID cannot be empty")
	}
	return u.Query.CheckIfUserBookmarkedContent(ctx, userID, contentID)
}
