package getstream_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

var (
	channelID            = "testChannelJnJ"
	channelType          = "messaging"
	testChannelOwner     = "256b9f95-c53e-44c3-81f7-4c37cf4e6510"
	member1              = "422a6d86-7f01-4c63-8ebd-51a343775f1b"
	member2              = "310145d2-95ad-4a2c-ac88-1e60bacfd37c"
	moderator1           = "3cad8cad-7623-4e78-b46e-06c150552aa3"
	userToAcceptInviteID = "f3926591-27f8-4756-90d9-fa88e228c582"
	userToRejectInviteID = "f3926591-27f8-4756-90d9-fa88e228c583"
	c                    getstream.ServiceGetStream
	ch                   *stream.CreateChannelResponse
	err                  error
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
	createTestUsers()

	run := m.Run()

	// Clean up
	deleteTestUsers()
	deleteTestChannel()

	os.Exit(run)

}

func createTestUsers() {
	ctx := context.Background()
	user1 := stream.User{
		ID:        member1,
		Name:      "member1",
		Invisible: false,
	}
	user2 := stream.User{
		ID:        member2,
		Name:      "member2",
		Invisible: false,
	}

	userModerator1 := stream.User{
		ID:        moderator1,
		Name:      "moderator1",
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

	_, err = c.CreateGetStreamUser(ctx, &userModerator1)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	_, err = c.AddModeratorsWithMessage(ctx, []string{moderator1}, channelID, nil)
	if err != nil {
		fmt.Printf("ChatClient.AddModeratorsWithMessage() error = %v", err)
		return
	}
}

func deleteTestUsers() {
	ctx := context.Background()
	_, err := c.DeleteUsers(
		ctx,
		[]string{member1, member2, moderator1, testChannelOwner, userToAcceptInviteID, userToRejectInviteID},
		stream.DeleteUserOptions{
			User:     stream.HardDelete,
			Messages: stream.HardDelete,
		},
	)
	if err != nil {
		fmt.Printf("ChatClient.DeleteUsers() error = %v", err)
	}
}

func deleteTestChannel() {
	ctx := context.Background()
	_, err = c.DeleteChannels(ctx, []string{channelType + ":" + channelID}, true)
	if err != nil {
		fmt.Printf("ChatClient.DeleteChannels() error = %v", err)
	}

}
