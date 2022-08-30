package cms

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/serverutils"
)

var (
	registerClientAPIEndpoint = serverutils.MustGetEnvVar("CLIENT_REGISTRATION_URL")
)

// IServiceCMS contains all the methods that are required to interact with the cms service
type IServiceCMS interface {
	RegisterClient(ctx context.Context, client *dto.PubSubCMSClientInput) error
}

// ServiceCMSImpl is the implementation of the IServiceCMS interface
type ServiceCMSImpl struct {
	CMSClient extension.ISCClientExtension
	baseExt   extension.ExternalMethodsExtension
}

// NewServiceCMS creates a new instance of the ServiceCMSImpl
func NewServiceCMS(
	cmsext extension.ISCClientExtension,
	ext extension.ExternalMethodsExtension,
) IServiceCMS {
	return &ServiceCMSImpl{
		CMSClient: cmsext,
		baseExt:   ext,
	}
}

// RegisterClient registers a new client in the cms service
func (s *ServiceCMSImpl) RegisterClient(ctx context.Context, client *dto.PubSubCMSClientInput) error {
	resp, err := s.CMSClient.MakeRequest(ctx, http.MethodPost, registerClientAPIEndpoint, client)
	if err != nil {
		return fmt.Errorf("error registering client: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register client : %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	return nil
}
