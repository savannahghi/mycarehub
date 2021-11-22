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
	"github.com/sirupsen/logrus"
)

const (
	contentAPIEndpoint = "https://mycarehub-test.savannahghi.org/contentapi/pages/?"
)

// IGetContent is used to fetch content from the CMS
type IGetContent interface {
	GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error)
}

// UseCasesContent holds the interfaces that are implemented within the content service
type UseCasesContent interface {
	IGetContent
}

// UseCasesContentImpl represents content implementation
type UseCasesContentImpl struct {
}

// NewUseCasesContentImplementation initializes a new contents service
func NewUseCasesContentImplementation() *UseCasesContentImpl {
	return &UseCasesContentImpl{}
}

// GetContent fetches content from wagtail CMS. The category ID is optional and it is used to return content based
// on the category it belongs to. The limit field describes how many items will be rendered on the front end side.
func (u UseCasesContentImpl) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
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
