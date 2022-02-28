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
		fmt.Printf("ChatClient.CreateChannel() error = %v", err)
		os.Exit(1)
	}

	run := m.Run()
	os.Exit(run)

	// Clean up
	_, err = c.DeleteChannels(ctx, []string{channelType + ":" + channelID}, true)
	if err != nil {
		fmt.Printf("ChatClient.DeleteChannels() error = %v", err)
		os.Exit(1)
	}

}
