package pubsubmessaging_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
)

func InitializeTestPubSub(t *testing.T) (*pubsubmessaging.ServicePubSubMessaging, error) {
	// Initialize base (common) extension
	baseExt := extension.NewExternalMethodsImpl()
	getStream := streamService.NewServiceGetStream(&stream.Client{})
	fcm := fcm.NewService()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate test repository: %v", err)
	}

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(baseExt, getStream, db, fcm)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pubsub messaging service: %w", err)
	}
	return pubSub, nil
}

func TestServicePubSubMessaging_AddPubSubNamespace(t *testing.T) {
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	topicName := pubsubmessaging.TestTopicName
	environment := serverutils.GetRunningEnvironment()

	type args struct {
		topicName   string
		serviceName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case -> Correct pubsub namespace",
			args: args{
				topicName:   topicName,
				serviceName: pubsubmessaging.MyCareHubServiceName,
			},
			want: fmt.Sprintf("%s-%s-%s-%s",
				pubsubmessaging.MyCareHubServiceName,
				topicName,
				environment,
				pubsubmessaging.TopicVersion,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ps.AddPubSubNamespace(tt.args.topicName, tt.args.serviceName)
			logrus.Printf("we got %v", got)
			if got != tt.want {
				t.Errorf("ServicePubSubMessaging.AddPubSubNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServicePubSubMessaging_PublishToPubsub(t *testing.T) {
	ctx := context.Background()
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	topic := ps.AddPubSubNamespace(pubsubmessaging.TestTopicName, pubsubmessaging.MyCareHubServiceName)
	// Create the test topic
	topics := ps.TopicIDs()
	topics = append(topics, topic)

	err = ps.EnsureTopicsExist(ctx, topics)
	if err != nil {
		t.Errorf("failed to create test topic")
		return
	}

	payload := map[string]interface{}{
		"name": "Test PubsubPayload",
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		topicID     string
		serviceName string
		payload     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad Case -> Fail to publish to pubsub - nil payload",
			args: args{
				ctx:         ctx,
				topicID:     topic,
				serviceName: pubsubmessaging.MyCareHubServiceName,
				payload:     nil,
			},
			wantErr: true,
		},
		{
			name: "Sad Case -> Fail to publish to pubsub - unknown topic",
			args: args{
				ctx:         ctx,
				topicID:     "invalid",
				serviceName: pubsubmessaging.MyCareHubServiceName,
				payload:     marshalled,
			},
			wantErr: true,
		},
		{
			name: "Happy Case-> Publish to pubsub",
			args: args{
				ctx:         ctx,
				topicID:     topic,
				serviceName: pubsubmessaging.MyCareHubServiceName,
				payload:     marshalled,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.PublishToPubsub(tt.args.ctx, tt.args.topicID, tt.args.serviceName, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.PublishToPubsub() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_EnsureTopicsExist(t *testing.T) {
	ctx := context.Background()
	ps, err := InitializeTestPubSub(t)
	if err != nil {
		t.Errorf("failed to initialize test pubsub: %v", err)
		return
	}

	type args struct {
		ctx      context.Context
		topicIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Create a topic",
			args: args{
				ctx:      ctx,
				topicIDs: []string{pubsubmessaging.TestTopicName},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.EnsureTopicsExist(tt.args.ctx, tt.args.topicIDs); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.EnsureTopicsExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
