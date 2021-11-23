package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/sirupsen/logrus"
)

const (
	contentAPIEndpoint = "https://mycarehub-test.savannahghi.org/contentapi/pages/?"
)

// IGetContent is used to fetch content from the CMS
type IGetContent interface {
	GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error)
}

// IContentCategoryList groups allthe content category listing methods
type IContentCategoryList interface {
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
}

// UseCasesContent holds the interfaces that are implemented within the content service
type UseCasesContent interface {
	IGetContent
	IContentCategoryList
}

// UsecaseContentImpl represents content implementation object
type UsecaseContentImpl struct {
	Query infrastructure.Query
}

// NewUsecaseContent returns a new content service
func NewUsecaseContent(query infrastructure.Query) *UsecaseContentImpl {
	return &UsecaseContentImpl{
		Query: query,
	}
}

// GetContent fetches content from wagtail CMS. The category ID is optional and it is used to return content based
// on the category it belongs to. The limit field describes how many items will be rendered on the front end side.
func (u UsecaseContentImpl) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	params := url.Values{}
	params.Add("type", "content.ContentItem")
	params.Add("limit", limit)
	params.Add("fields", "*")
	params.Add("order", "-first_published_at")

	getContentEndpoint := fmt.Sprintf(contentAPIEndpoint + params.Encode())
	logrus.Printf("the url is %v", getContentEndpoint)
	var contentItems *domain.Content
	resp, err := utils.MakeRequest(ctx, http.MethodGet, "https://mycarehub-test.savannahghi.org/contentapi/pages/?fields=*&limit=10&order=-first_published_at&type=content.ContentItem", nil)
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
func (u *UsecaseContentImpl) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return u.Query.ListContentCategories(ctx)
}
