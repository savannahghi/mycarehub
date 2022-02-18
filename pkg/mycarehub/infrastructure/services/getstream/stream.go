package getstream

import (
	"context"
	"log"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/serverutils"
)

var (
	getStreamAPIKey          = serverutils.MustGetEnvVar("GET_STREAM_KEY")
	getStreamAPISecret       = serverutils.MustGetEnvVar("GET_STREAM_SECRET")
	getStreamTokenExpiryTime = time.Now().UTC().Add(time.Hour * 12)
)

// ServiceGetStream represents the various Getstream usecases
type ServiceGetStream interface {
	CreateGetStreamUserToken(ctx context.Context, userID string) (string, error)
	CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
}

// ChatClient is the service's struct implementation
type ChatClient struct {
	client *stream.Client
}

// NewServiceGetStream initializes a new getstream service
func NewServiceGetStream() ServiceGetStream {
	client, err := stream.NewClient(getStreamAPIKey, getStreamAPISecret)
	if err != nil {
		log.Fatalf("failed to start getstream client: %v", err)
	}

	return &ChatClient{
		client: client,
	}
}

// CreateGetStreamUserToken creates a new token for a user with optional expire time. This token is handed
// to the client side during login. It allows the client side to connect to the chat API for that user.
func (c *ChatClient) CreateGetStreamUserToken(ctx context.Context, userID string) (string, error) {
	return c.client.CreateToken(userID, getStreamTokenExpiryTime)
}

// CreateGetStreamUser creates or updates a user
func (c *ChatClient) CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return c.client.UpsertUser(ctx, user)
}
