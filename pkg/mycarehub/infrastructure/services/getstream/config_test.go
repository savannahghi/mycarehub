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
	channelID        = "testChannelJnJ"
	channelType      = "messaging"
	channelCID       = channelType + ":" + channelID
	testChannelOwner = "256b9f95-c53e-44c3-81f7-4c37cf4e6510"

	// Channel members
	defaultMemberID              = "422a6d86-7f01-4c63-8ebd-51a343775f1b"
	defaultModeratorID           = "3cad8cad-7623-4e78-b46e-06c150552aa3"
	userToAcceptInviteID         = "f3926591-27f8-4756-90d9-fa88e228c582"
	userToAcceptInviteName       string
	userToAddToNewChannelID      = "62605ecd-5136-446c-876b-89ffa17335fe"
	userToRejectInviteID         = "f3926591-27f8-4756-90d9-fa88e228c583"
	userToBanID                  = "19d7d30d-6cc2-483c-98cb-c261fb2cae54"
	userToUnbanID                = "114b857c-4365-4a68-9299-ed33975d9ddc"
	userRemoveFromCommunityID    = "71c90bc3-bb6d-4a6e-b9e1-a8dee8059431"
	moderatorToDemoteID          = "627f9740-41cb-409d-8f1a-fa9b05733609"
	userToRevokeGetstreamTokenID = "e75f3cc5-d085-4df4-b870-1d3aa3430d82"
	userToUpsertID               = "e3620079-1e98-4d48-8d89-530ad5c1978a"
	userToDeleteID               = "71c90bc3-bb6d-4a6e-b9e1-a8dee8059432"

	c   getstream.ServiceGetStream
	ch  *stream.CreateChannelResponse
	err error
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

	defaultMember := stream.User{
		ID:        defaultMemberID,
		Name:      "defaultMember",
		Invisible: false,
	}
	_, err := c.CreateGetStreamUser(ctx, &defaultMember)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	defaultModerator := stream.User{
		ID:        defaultModeratorID,
		Name:      "defaultModerator",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &defaultModerator)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}
	_, err = c.AddModeratorsWithMessage(ctx, []string{defaultModeratorID}, channelID, nil)
	if err != nil {
		fmt.Printf("ChatClient.AddModeratorsWithMessage() error = %v", err)
		return
	}

	moderatorToDemote := stream.User{
		ID:        moderatorToDemoteID,
		Name:      "moderatorToDemote",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &moderatorToDemote)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}
	_, err = c.AddModeratorsWithMessage(ctx, []string{moderatorToDemoteID}, channelID, nil)
	if err != nil {
		fmt.Printf("ChatClient.AddModeratorsWithMessage() error = %v", err)
		return
	}

	userToBan := &stream.User{
		ID:        userToBanID,
		Name:      "userToBan",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, userToBan)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToUnban := &stream.User{
		ID:        userToUnbanID,
		Name:      "userToUnban",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, userToUnban)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToRejectInvite := stream.User{
		ID:        userToRejectInviteID,
		Name:      "userToRejectInvite",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &userToRejectInvite)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToAcceptInvite := stream.User{
		ID:        userToAcceptInviteID,
		Name:      "userToAcceptInvite",
		Invisible: false,
	}
	userToAcceptInviteModel, err := c.CreateGetStreamUser(ctx, &userToAcceptInvite)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}
	userToAcceptInviteName = userToAcceptInviteModel.User.Name

	userRemoveFromCommunity := stream.User{
		ID:        userRemoveFromCommunityID,
		Name:      "userRemoveFromCommunity",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &userRemoveFromCommunity)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToRevokeGetstreamToken := stream.User{
		ID:        userToRevokeGetstreamTokenID,
		Name:      "userToRevokeGetstreamToken",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &userToRevokeGetstreamToken)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToUpsert := stream.User{
		ID:        userToUpsertID,
		Name:      "userToUpsert",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &userToUpsert)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	userToDelete := stream.User{
		ID:        userToDeleteID,
		Name:      "userToDelete",
		Invisible: false,
	}
	_, err = c.CreateGetStreamUser(ctx, &userToDelete)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	// Add them to a community
	_, err = c.AddMembersToCommunity(ctx, []string{
		defaultMemberID,
		defaultModeratorID,
		userToBanID,
		userToUnbanID,
		userRemoveFromCommunityID,
		moderatorToDemoteID,
		userToRevokeGetstreamTokenID,
		userToUpsertID,
		userToDeleteID,
	}, channelID)
	if err != nil {
		fmt.Printf("ChatClient.AddMembersToCommunity() error = %v", err)
		return
	}

	// Invite members to accept invite and members to reject invite to community
	_, err = c.InviteMembers(ctx, []string{userToRejectInviteID, userToAcceptInviteID}, channelID, nil)
	if err != nil {
		fmt.Printf("ChatClient.InviteMembers() error = %v", err)
		return
	}

	// ban  user who should be unbanned
	_, err = c.BanUser(ctx, userToUnbanID, defaultModeratorID, channelID)
	if err != nil {
		fmt.Printf("unable to ban user: %v", err)
		return
	}

}

func deleteTestUsers() {
	ctx := context.Background()
	_, err := c.DeleteUsers(
		ctx,
		[]string{
			testChannelOwner,
			defaultMemberID,
			defaultModeratorID,
			userToAcceptInviteID,
			userToRejectInviteID,
			userToBanID,
			userToUnbanID,
			userRemoveFromCommunityID,
			moderatorToDemoteID,
			userToRevokeGetstreamTokenID,
			userToUpsertID,
		},
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
	_, err = c.DeleteChannels(ctx, []string{channelCID}, true)
	if err != nil {
		fmt.Printf("ChatClient.DeleteChannels() error = %v", err)
	}

}
