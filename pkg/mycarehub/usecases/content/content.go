package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/serverutils"
)

var (
	contentBaseURL = serverutils.MustGetEnvVar("CONTENT_SERVICE_BASE_URL")
)

// IGetContent is used to fetch content from the CMS
type IGetContent interface {
	GetContent(ctx context.Context, categoryIDs []int, categoryNames []string, limit string, clientID *string) (*domain.Content, error)
	GetContentItemByID(ctx context.Context, contentID int) (*domain.ContentItem, error)
	GetFAQs(ctx context.Context, flavour feedlib.Flavour) (*domain.Content, error)
	FetchContent(ctx context.Context, limit string) (*domain.Content, error)
}

// IGetBookmarkedContent holds the method signature used to return a user's bookmarked content
type IGetBookmarkedContent interface {
	GetUserBookmarkedContent(ctx context.Context, clientID string) (*domain.Content, error)
}

// ICheckIfUserBookmarkedContent is used to check if a user has bookmarked a content item
type ICheckIfUserBookmarkedContent interface {
	CheckIfUserBookmarkedContent(ctx context.Context, clientID string, contentID int) (bool, error)
}

// IContentCategoryList groups all the content category listing methods
type IContentCategoryList interface {
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
}

// IShareContent is the interface for the ShareContent
type IShareContent interface {
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
}

// IBookmarkContent is used to bookmark content
type IBookmarkContent interface {
	BookmarkContent(ctx context.Context, clientID string, contentID int) (bool, error)
}

// IUnBookmarkContent is used to unbookmark content
type IUnBookmarkContent interface {
	UnBookmarkContent(ctx context.Context, clientID string, contentID int) (bool, error)
}

// ILikeContent groups the like feature methods
type ILikeContent interface {
	LikeContent(ctx context.Context, clientID string, contentID int) (bool, error)
	CheckWhetherUserHasLikedContent(ctx context.Context, clientID string, contentID int) (bool, error)
}

// IUnlikeContent groups the unllike feature methods
type IUnlikeContent interface {
	UnlikeContent(ctx context.Context, clientID string, contentID int) (bool, error)
}

// IViewContent gets a content ite and updates the view count
type IViewContent interface {
	ViewContent(ctx context.Context, clientID string, contentID int) (bool, error)
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

// UseCasesContentImpl represents content implementation
type UseCasesContentImpl struct {
	Update      infrastructure.Update
	Query       infrastructure.Query
	ExternalExt extension.ExternalMethodsExtension
}

// NewUseCasesContentImplementation initializes a new contents service
func NewUseCasesContentImplementation(
	update infrastructure.Update,
	query infrastructure.Query,
	externalExt extension.ExternalMethodsExtension,
) *UseCasesContentImpl {
	return &UseCasesContentImpl{
		Update:      update,
		Query:       query,
		ExternalExt: externalExt,
	}

}

// LikeContent implements the content liking api
func (u UseCasesContentImpl) LikeContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	if clientID == "" || contentID == 0 {
		return false, fmt.Errorf("client id an content id are required")
	}

	contentLikeAPI := fmt.Sprintf("%s/api/content_like/", contentBaseURL)

	payload := struct {
		Active      bool   `json:"active"`
		Client      string `json:"client"`
		ContentItem int    `json:"content_item"`
	}{
		Active:      true,
		Client:      clientID,
		ContentItem: contentID,
	}

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodPost, contentLikeAPI, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("failed to like content")
	}

	return true, nil
}

// CheckWhetherUserHasLikedContent implements action of checking whether a user has liked a particular content
func (u UseCasesContentImpl) CheckWhetherUserHasLikedContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	if clientID == "" || contentID <= 0 {
		return false, fmt.Errorf("client id and content id are required")
	}

	params := url.Values{}
	params.Add("client", clientID)
	params.Add("content_item", strconv.Itoa(contentID))
	params.Add("active", "True")

	contentLikeAPI := fmt.Sprintf("%s/api/content_like/?%s", contentBaseURL, params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentLikeAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to check if client liked content")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read request body: %v", err)
	}

	count := struct {
		Count int `json:"count"`
	}{}

	err = json.Unmarshal(body, &count)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if count.Count == 0 {
		return false, nil
	}

	return true, nil
}

// UnlikeContent implements the content liking api
func (u UseCasesContentImpl) UnlikeContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	if clientID == "" || contentID == 0 {
		return false, fmt.Errorf("client id and content id are required")
	}
	params := url.Values{}
	params.Add("client", clientID)
	params.Add("content_item", strconv.Itoa(contentID))
	params.Add("active", "True")

	contentLikeAPI := fmt.Sprintf("%s/api/content_like/?%s", contentBaseURL, params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentLikeAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read request body: %v", err)
	}

	result := struct {
		Count   int `json:"count"`
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if result.Count == 0 {
		return false, nil
	}

	deleteContentLikeAPI := fmt.Sprintf("%s/api/content_like/%s/", contentBaseURL, result.Results[0].ID)

	resp, err := u.ExternalExt.MakeRequest(ctx, http.MethodDelete, deleteContentLikeAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if resp.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("failed to unlike content")
	}

	return true, nil
}

// GetContent fetches content from wagtail CMS. The category ID is optional and it is used to return content based
// on the category it belongs to. The limit field describes how many items will be rendered on the front end side.
func (u UseCasesContentImpl) GetContent(ctx context.Context, categoryIDs []int, categoryNames []string, limit string, clientID *string) (*domain.Content, error) {
	var clientProfile *domain.ClientProfile
	var err error

	if clientID == nil {
		uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
		if err != nil {
			return nil, err
		}

		userProfile, err := u.Query.GetUserProfileByUserID(ctx, uid)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		clientProfile, err = u.Query.GetClientProfile(ctx, *userProfile.ID, userProfile.CurrentProgramID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
	} else {
		clientProfile, err = u.Query.GetClientProfileByClientID(ctx, *clientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
	}

	params := url.Values{}
	params.Add("client_id", *clientProfile.ID)
	params.Add("facility_id", *clientProfile.DefaultFacility.ID)
	params.Add("limit", limit)

	if len(categoryIDs) > 0 {
		var categories []string
		for _, categoryID := range categoryIDs {
			categories = append(categories, strconv.Itoa(categoryID))
		}

		categoryIDList := strings.Join(categories, ",")
		params.Add("category", categoryIDList)
	}

	if len(categoryNames) > 0 {
		categoryNameList := strings.Join(categoryNames, ",")
		params.Add("category_name", categoryNameList)
	}

	contentItems, err := u.fetchContent(ctx, params)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return contentItems, nil
}

// FetchContent is used to fetch content from the CMS to serve the content landing page for non-registered users
func (u *UseCasesContentImpl) FetchContent(ctx context.Context, limit string) (*domain.Content, error) {
	params := url.Values{}
	params.Add("limit", limit)

	return u.fetchContent(ctx, params)
}

// fetchContent is a helper method to make a request to CMS and fetch content. It's purpose is majorly to allow fetching content for both landing page
// and service other authenticated APIs that fetch content.
func (u *UseCasesContentImpl) fetchContent(ctx context.Context, params url.Values) (*domain.Content, error) {
	var (
		consumerFAQs = "consumer-faqs"
		proFAQs      = "pro-faqs"
	)

	params.Add("type", "content.ContentItem")
	params.Add("order", "-first_published_at")
	params.Add("fields", "*")

	params.Add("exclude_category", consumerFAQs)
	params.Add("exclude_category", proFAQs)

	getContentEndpoint := fmt.Sprintf("%s/contentapi/pages/?%s", contentBaseURL, params.Encode())
	resp, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, getContentEndpoint, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request")
	}

	dataResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	var contentItems *domain.Content
	err = json.Unmarshal(dataResponse, &contentItems)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return contentItems, nil
}

// ListContentCategories gets the list of all content categories
func (u *UseCasesContentImpl) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	uid, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	userProfile, err := u.Query.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	params := url.Values{}
	params.Add("has_content", "True")
	params.Add("program_id", userProfile.CurrentProgramID)
	categoryAPI := fmt.Sprintf("%s/api/content_item_category/?%s", contentBaseURL, params.Encode())

	resp, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, categoryAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	results := struct {
		Results []*domain.ContentItemCategory
	}{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return results.Results, nil
}

// ShareContent enables a user to share a content
func (u *UseCasesContentImpl) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	if input.ClientID == "" || input.ContentID <= 0 {
		return false, fmt.Errorf("client id and content id are required")
	}
	contentShareAPI := fmt.Sprintf("%s/api/content_share/", contentBaseURL)

	payload := struct {
		Active      bool   `json:"active"`
		Client      string `json:"client"`
		ContentItem int    `json:"content_item"`
	}{
		Active:      true,
		Client:      input.ClientID,
		ContentItem: input.ContentID,
	}

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodPost, contentShareAPI, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("failed to share content")
	}

	return true, nil
}

// BookmarkContent increments the bookmark count for a content item
func (u UseCasesContentImpl) BookmarkContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	contentBookmarkAPI := fmt.Sprintf("%s/api/content_bookmark/", contentBaseURL)

	payload := struct {
		Active      bool   `json:"active"`
		Client      string `json:"client"`
		ContentItem int    `json:"content_item"`
	}{
		Active:      true,
		Client:      clientID,
		ContentItem: contentID,
	}

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodPost, contentBookmarkAPI, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("failed to bookmark content")
	}

	return true, nil
}

// UnBookmarkContent decrements the bookmark count for a content item
func (u UseCasesContentImpl) UnBookmarkContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	if clientID == "" || contentID == 0 {
		return false, fmt.Errorf("client id and content id are required")
	}

	params := url.Values{}
	params.Add("client", clientID)
	params.Add("content_item", strconv.Itoa(contentID))
	params.Add("active", "True")

	contentBookmarkAPI := fmt.Sprintf("%s/api/content_bookmark/?%s", contentBaseURL, params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentBookmarkAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read request body: %v", err)
	}

	result := struct {
		Count   int `json:"count"`
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if result.Count == 0 {
		return false, nil
	}

	deleteContentBookmarkAPI := fmt.Sprintf("%s/api/content_bookmark/%s/", contentBaseURL, result.Results[0].ID)

	resp, err := u.ExternalExt.MakeRequest(ctx, http.MethodDelete, deleteContentBookmarkAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if resp.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("failed to remove bookmark")
	}

	return true, nil
}

// GetUserBookmarkedContent gets the user's pinned/bookmarked content and displays it on their profile
func (u *UseCasesContentImpl) GetUserBookmarkedContent(ctx context.Context, clientID string) (*domain.Content, error) {
	if clientID == "" {
		return nil, exceptions.EmptyInputErr(fmt.Errorf("client ID must be defined"))
	}

	params := url.Values{}
	params.Add("client", clientID)
	params.Add("active", "True")

	contentBookmarkAPI := fmt.Sprintf("%s/api/content_bookmark/?%s", contentBaseURL, params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentBookmarkAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	result := struct {
		Count   int `json:"count"`
		Results []struct {
			ID          string `json:"id"`
			Client      string `json:"client"`
			ContentItem int    `json:"content_item"`
		} `json:"results"`
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	content := &domain.Content{
		Meta: domain.Meta{
			TotalCount: 0,
		},
		Items: []domain.ContentItem{},
	}

	for _, item := range result.Results {
		contentItem, err := u.GetContentItemByID(ctx, item.ContentItem)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			// should not fail to render bookmarked content
			continue
		}
		content.Items = append(content.Items, *contentItem)
		content.Meta.TotalCount++
	}

	return content, nil
}

// GetContentItemByID fetches a specific content using the specific content item ID. This will be important
// when fetching content that a user bookmarked. The data returned directly from the database does not contain all the
// information regarding a content item hence why this method has been chosen.
func (u *UseCasesContentImpl) GetContentItemByID(ctx context.Context, contentID int) (*domain.ContentItem, error) {
	getContentEndpoint := fmt.Sprintf(contentBaseURL+"/contentapi/pages/%s/", strconv.Itoa(contentID))

	var contentItem *domain.ContentItem
	resp, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, getContentEndpoint, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user bookmarked content")
	}

	dataResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	err = json.Unmarshal(dataResponse, &contentItem)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contentItem, nil
}

// ViewContent gets a content item and updates the view count
func (u *UseCasesContentImpl) ViewContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	contentViewAPI := fmt.Sprintf("%s/api/content_view/", contentBaseURL)

	payload := struct {
		Active      bool   `json:"active"`
		Client      string `json:"client"`
		ContentItem int    `json:"content_item"`
	}{
		Active:      true,
		Client:      clientID,
		ContentItem: contentID,
	}

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodPost, contentViewAPI, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("failed to view content")
	}

	return true, nil
}

// CheckIfUserBookmarkedContent checks if a user has bookmarked a specific content item
func (u *UseCasesContentImpl) CheckIfUserBookmarkedContent(ctx context.Context, clientID string, contentID int) (bool, error) {
	if clientID == "" || contentID <= 0 {
		return false, fmt.Errorf("clientID and contentID cannot be empty")
	}

	params := url.Values{}
	params.Add("client", clientID)
	params.Add("content_item", strconv.Itoa(contentID))
	params.Add("active", "True")

	contentBookmarkAPI := fmt.Sprintf("%s/api/content_bookmark/?%s", contentBaseURL, params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentBookmarkAPI, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to make request")
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to check bookmarked content")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read request body: %v", err)
	}

	count := struct {
		Count int `json:"count"`
	}{}

	err = json.Unmarshal(body, &count)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if count.Count == 0 {
		return false, nil
	}

	return true, nil
}

// GetFAQs retrieves the faqs depending on the provided flavour
func (u *UseCasesContentImpl) GetFAQs(ctx context.Context, flavour feedlib.Flavour) (*domain.Content, error) {
	// 'consumer-faqs' and 'pro-faqs' are CMS category names for FAQs contents
	var (
		consumerFAQs = "consumer-faqs"
		proFAQs      = "pro-faqs"
	)

	params := url.Values{}

	switch flavour {
	case feedlib.FlavourConsumer:
		params.Add("category_name", consumerFAQs)

	case feedlib.FlavourPro:
		params.Add("category_name", proFAQs)
	}

	params.Add("type", "content.ContentItem")
	params.Add("limit", "20")
	params.Add("order", "-first_published_at")
	params.Add("fields", "'*")

	contentURL := fmt.Sprintf(contentBaseURL + "/contentapi/pages/?" + params.Encode())

	response, err := u.ExternalExt.MakeRequest(ctx, http.MethodGet, contentURL, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to make request")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	var contentItems *domain.Content
	err = json.Unmarshal(data, &contentItems)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contentItems, nil
}
