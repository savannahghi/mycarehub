package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// FakeCMSService contains the mock implementation of the cms service
type FakeCMSService struct {
	MockRegisterClientFn func(ctx context.Context, client *dto.PubSubCMSClientInput) error
}

// NewFakeCMSService creates a new instance of the FakeCMSService
func NewFakeCMSService() *FakeCMSService {
	return &FakeCMSService{
		MockRegisterClientFn: func(ctx context.Context, client *dto.PubSubCMSClientInput) error {
			return nil
		},
	}
}

// RegisterClient registers a new client in the cms service
func (s *FakeCMSService) RegisterClient(ctx context.Context, client *dto.PubSubCMSClientInput) error {
	return s.MockRegisterClientFn(ctx, client)
}
