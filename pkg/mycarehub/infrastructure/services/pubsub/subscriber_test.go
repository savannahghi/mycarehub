package pubsubmessaging_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"google.golang.org/api/idtoken"
)

var (
	srv                *http.Server
	baseURL            string
	serverErr          error
	getStreamAPIKey    = serverutils.MustGetEnvVar("GET_STREAM_KEY")
	getStreamAPISecret = serverutils.MustGetEnvVar("GET_STREAM_SECRET")
	channelID          = "testChannelJnJ"
	channelType        = "messaging"
	channelCID         = channelType + ":" + channelID
	testChannelOwner   = uuid.New().String()
)

func TestMain(m *testing.M) {
	initialEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")

	ctx := context.Background()
	srv, baseURL, serverErr = serverutils.StartTestServer(ctx, presentation.PrepareServer, presentation.AllowedOrigins)
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	code := m.Run()

	// restore envs
	os.Setenv("ENVIRONMENT", initialEnv)
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()

	os.Exit(code)
}

func createStreamChannel(streamSvc getstream.ServiceGetStream) {
	ctx := context.Background()
	// create origin test user
	testChannelMember := stream.User{
		ID:        testChannelOwner,
		Name:      "testChannelOwner",
		Invisible: false,
		ExtraData: map[string]interface{}{
			"userType": "STAFF",
		},
	}
	_, err := streamSvc.CreateGetStreamUser(ctx, &testChannelMember)
	if err != nil {
		fmt.Printf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	// Create a test channel
	_, err = streamSvc.CreateChannel(ctx, channelType, channelID, testChannelOwner, map[string]interface{}{
		"description": "This is just a test channel",
		"inviteOnly":  true,
	})
	if err != nil {
		fmt.Printf("ChatClient.CreateCommunity() error = %v", err)
	}

	_, err = streamSvc.AddMembersToCommunity(ctx, []string{testChannelOwner}, channelID)
	if err != nil {
		fmt.Printf("ChatClient.AddMembersToCommunity() error = %v", err)
	}
}

func deleteTestData(streamSvc getstream.ServiceGetStream) {
	ctx := context.Background()
	_, err := streamSvc.DeleteChannels(ctx, []string{channelCID}, true)
	if err != nil {
		fmt.Printf("ChatClient.DeleteChannels() error = %v", err)
	}

	_, err = streamSvc.DeleteUsers(
		ctx,
		[]string{
			testChannelOwner,
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

func composeInvalidPubsubTestPayload(t *testing.T, topic string) (*bytes.Buffer, error) {
	// Compose the payload
	pubsubPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: nil,
			Attributes: map[string]string{
				"invalid": "invalid",
			},
		},
		Subscription: ksuid.New().String(),
	}

	payload, err := json.Marshal(pubsubPayload)
	if err != nil {
		return nil, err
	}
	bs := bytes.NewBuffer(payload)

	return bs, nil
}

func composeGetStreamEventPayload(
	t *testing.T,
	eventType stream.EventType,
	topic string,
	channelID string,
) *bytes.Buffer {
	eventPayload := &dto.GetStreamEvent{
		CID:      channelCID,
		Type:     eventType,
		Message:  &stream.Message{},
		Reaction: &stream.Reaction{},
		Channel:  &stream.Channel{},
		Member:   &stream.ChannelMember{},
		Members: []*stream.ChannelMember{
			{
				UserID: testChannelOwner,
				User: &stream.User{
					ID:        testChannelOwner,
					Name:      "testChannelOwner",
					Invisible: false,
					ExtraData: map[string]interface{}{
						"userType": "STAFF",
					},
				},
				IsModerator: true,
			},
		},
		User:         &stream.User{},
		UserID:       "",
		OwnUser:      &stream.User{},
		WatcherCount: 0,
		CreatedAt:    time.Time{},
		ChannelID:    channelID,
	}
	payload, err := json.Marshal(eventPayload)
	if err != nil {
		log.Printf("an error occurred: %v", err)
		return nil
	}

	pubsubPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: payload,
			Attributes: map[string]string{
				"topicID": namespacePubsubIdentifier(pubsubmessaging.MyCareHubServiceName, topic),
			},
		},
		Subscription: ksuid.New().String(),
	}

	testDataJSON, err := json.Marshal(pubsubPayload)
	if err != nil {
		t.Errorf("can't marshal JSON: %v", err)
		return nil
	}

	return bytes.NewBuffer(testDataJSON)
}

func namespacePubsubIdentifier(
	serviceName string,
	topicID string,
) string {
	return fmt.Sprintf(
		"%s-%s-staging-v1",
		serviceName,
		topicID,
	)
}

func TestPubsub(t *testing.T) {
	ctx := context.Background()
	streamClient, err := stream.NewClient(getStreamAPIKey, getStreamAPISecret)
	if err != nil {
		log.Fatalf("failed to start getstream client: %v", err)
	}
	c := getstream.NewServiceGetStream(streamClient)
	createStreamChannel(c)

	invalidPayload, err := composeInvalidPubsubTestPayload(t, "invalidTopic")
	if err != nil {
		t.Errorf("failed to compose invalid payload")
		return
	}
	newMessageEventPayload := composeGetStreamEventPayload(t, stream.EventMessageNew, common.CreateGetstreamEventTopicName, channelID)
	invalidMessagePayload := composeGetStreamEventPayload(t, stream.EventMessageNew, common.CreateGetstreamEventTopicName, "invalidChannelID")

	newMessageFlaggedPayload := composeGetStreamEventPayload(t, "message.flagged", common.CreateGetstreamEventTopicName, channelID)
	invalidMessageFlaggedPayload := composeGetStreamEventPayload(t, "message.flagged", common.CreateGetstreamEventTopicName, "invalidChannelID")

	newMemberAddedPayload := composeGetStreamEventPayload(t, stream.EventMemberAdded, common.CreateGetstreamEventTopicName, channelID)
	invalidMemberAddedPayload := composeGetStreamEventPayload(t, stream.EventMemberAdded, common.CreateGetstreamEventTopicName, "invalidChannelID")

	memberRemovedPayload := composeGetStreamEventPayload(t, stream.EventMemberRemoved, common.CreateGetstreamEventTopicName, channelID)
	invalidMemberRemovedPayload := composeGetStreamEventPayload(t, stream.EventMemberRemoved, common.CreateGetstreamEventTopicName, "invalidChannelID")

	userBannedPayload := composeGetStreamEventPayload(t, "user.banned", common.CreateGetstreamEventTopicName, channelID)
	invalidUserBannedPayload := composeGetStreamEventPayload(t, "user.banned", common.CreateGetstreamEventTopicName, "invalidChannelID")

	userUnbannedPayload := composeGetStreamEventPayload(t, "user.unbanned", common.CreateGetstreamEventTopicName, channelID)
	invalidChannelPayload := composeGetStreamEventPayload(t, "user.unbanned", common.CreateGetstreamEventTopicName, "invalidChannelID")

	header := req.Header{
		"Content-Type": "application/json",
	}
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
		headers    map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid payload",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Happy Case - Successfully publish to new message event topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       newMessageEventPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy Case - Successfully publish to new message flagged topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       newMessageFlaggedPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy Case - Successfully publish to new member added topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       newMemberAddedPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy Case - Successfully publish to member removed topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       memberRemovedPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy Case - Successfully publish to user banned topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       userBannedPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy Case - Successfully publish to user unbanned topic",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       userUnbannedPayload,
				headers:    header,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad Case - Use an invalid channel ID - user unbanned event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidChannelPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Use an invalid channel ID - new message event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidMessagePayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Use an invalid channel ID - flagged message event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidMessageFlaggedPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Use an invalid channel ID - new member added event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidMemberAddedPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Use an invalid channel ID - member removed event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidMemberRemovedPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Use an invalid channel ID - user banned event",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidUserBannedPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		r, err := http.NewRequest(
			tt.args.httpMethod,
			tt.args.url,
			tt.args.body,
		)
		if err != nil {
			t.Errorf("unable to compose request: %s", err)
			return
		}

		if r == nil {
			t.Errorf("nil request")
			return
		}

		for k, v := range tt.args.headers {
			r.Header.Add(k, v)
		}

		client, err := idtoken.NewClient(ctx, pubsubtools.Aud)
		if err != nil {
			t.Errorf("can't initialize client: %s", err)
			return
		}
		resp, err := client.Do(r)
		if err != nil {
			t.Errorf("request error: %s", err)
			return
		}

		dataResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("can't read request body: %s", err)
			return
		}
		if dataResponse == nil {
			t.Errorf("nil response data")
			return
		}

		if tt.wantStatus != resp.StatusCode {
			t.Errorf(
				"expected status %d, got %d and response %s",
				tt.wantStatus,
				resp.StatusCode,
				string(dataResponse),
			)
			return
		}
	}
	deleteTestData(c)
}
