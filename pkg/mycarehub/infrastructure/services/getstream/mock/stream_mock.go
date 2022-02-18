package mock

import (
	"context"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

// GetStreamServiceMock mocks the GetStream service library implementations
type GetStreamServiceMock struct {
	MockCreateGetStreamUserTokenFn func(ctx context.Context, userID string) (string, error)
	MockCreateGetStreamUserFn      func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
}

// NewGetStreamServiceMock initializes the mock service
func NewGetStreamServiceMock() *GetStreamServiceMock {
	return &GetStreamServiceMock{
		MockCreateGetStreamUserTokenFn: func(ctx context.Context, userID string) (string, error) {
			return uuid.New().String(), nil
		},
		MockCreateGetStreamUserFn: func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
			return &stream.UpsertUserResponse{
				User: &stream.User{
					ID:   uuid.New().String(),
					Name: gofakeit.Name(),
				},
			}, nil
		},
	}
}

// CreateGetStreamUserToken mocks creating a getstream user token
func (g GetStreamServiceMock) CreateGetStreamUserToken(ctx context.Context, userID string) (string, error) {
	return g.MockCreateGetStreamUserTokenFn(ctx, userID)
}

// CreateGetStreamUser mocks creating a getstream user
func (g GetStreamServiceMock) CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return g.MockCreateGetStreamUserFn(ctx, user)
}
