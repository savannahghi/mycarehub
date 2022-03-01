package getstream_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

var (
	channelID        = "testChannelJnJ"
	channelType      = "messaging"
	testChannelOwner = uuid.New().String()
	member1          = uuid.New().String()
	member2          = uuid.New().String()
	c                getstream.ServiceGetStream
	ch               *stream.CreateChannelResponse
	err              error
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	c = getstream.NewServiceGetStream()

	// Create a test channel
	ch, err = c.CreateChannel(ctx, channelType, channelID, testChannelOwner, nil)
	if err != nil {
		fmt.Printf("ChatClient.CreateCommunity() error = %v", err)
		os.Exit(1)
	}

	// Add members to the channel
	addMembersToAChannel()

	run := m.Run()
	os.Exit(run)

	// Clean up
	_, err = c.DeleteChannels(ctx, []string{channelType + ":" + channelID}, true)
	if err != nil {
		fmt.Printf("ChatClient.DeleteChannels() error = %v", err)
		os.Exit(1)
	}

}

func addMembersToAChannel() {
	ctx := context.Background()
	user1 := stream.User{
		ID:        member1,
		Name:      "test",
		Invisible: false,
	}
	user2 := stream.User{
		ID:        member2,
		Name:      "test",
		Invisible: false,
	}
	_, err := c.CreateGetStreamUser(ctx, &user1)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	_, err = c.CreateGetStreamUser(ctx, &user2)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}
}
