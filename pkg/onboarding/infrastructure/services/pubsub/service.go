package pubsubmessaging

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/common"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	erp "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases"
)

const (
	// ServiceName ..
	ServiceName = "onboarding"

	// TopicVersion ...
	TopicVersion = "v1"

	hostNameEnvVarName = "SERVICE_HOST"

	engagementService = "engagement"
)

// ServicePubSub represents logic required to communicate with pubsub
type ServicePubSub interface {
	// Ask the notification service about the topics that it knows about
	TopicIDs() []string
	AddPubSubNamespace(topicName string) string
	PublishToPubsub(
		ctx context.Context,
		topicID string,
		payload []byte,
	) error
	EnsureTopicsExist(
		ctx context.Context,
		topicIDs []string,
	) error
	EnsureSubscriptionsExist(
		ctx context.Context,
	) error
	SubscriptionIDs() map[string]string
	ReceivePubSubPushMessages(
		w http.ResponseWriter,
		r *http.Request,
	)
	AddEngagementPubsubNameSpace(topic string) string

	// Publishers
	NotifyCreateContact(ctx context.Context, contact domain.CRMContact) error
	NotifyCoverLinking(ctx context.Context, data dto.LinkCoverPubSubMessage) error
	NotifyUpdateContact(
		ctx context.Context,
		updateData dto.UpdateContactPSMessage,
	) error
	NotifyCreateCustomer(
		ctx context.Context,
		data dto.CustomerPubSubMessage,
	) error
	NotifyCreateSupplier(
		ctx context.Context,
		data dto.SupplierPubSubMessage,
	) error
}

// ServicePubSubMessaging sends "real" (production) notifications
type ServicePubSubMessaging struct {
	client  *pubsub.Client
	baseExt extension.BaseExtension
	erp     erp.AccountingUsecase
	crm     extension.CRMExtension
	edi     edi.ServiceEdi
	repo    repository.OnboardingRepository
}

// NewServicePubSubMessaging ...
func NewServicePubSubMessaging(
	client *pubsub.Client,
	ext extension.BaseExtension,
	erp erp.AccountingUsecase,
	crm extension.CRMExtension,
	edi edi.ServiceEdi,
	repo repository.OnboardingRepository,
) (*ServicePubSubMessaging, error) {
	s := &ServicePubSubMessaging{
		client:  client,
		baseExt: ext,
		erp:     erp,
		crm:     crm,
		edi:     edi,
		repo:    repo,
	}

	ctx := context.Background()
	if err := s.EnsureTopicsExist(
		ctx,
		s.TopicIDs(),
	); err != nil {
		return nil, err
	}

	if err := s.EnsureSubscriptionsExist(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// AddEngagementPubsubNameSpace creates a namespaced topic that resembles the one in
// engagement service, which is prepended with the word "engagement". This solves the problem
// where namespaced topics from "onboarding" are different from the ones in engagement.
// This fix allows for uniformity of topic names between the engagement and onboarding services.
func (ps ServicePubSubMessaging) AddEngagementPubsubNameSpace(
	topic string,
) string {
	environment := ps.baseExt.GetRunningEnvironment()
	return ps.baseExt.NamespacePubsubIdentifier(
		engagementService,
		topic,
		environment,
		TopicVersion,
	)
}

// AddPubSubNamespace creates a namespaced topic name
func (ps ServicePubSubMessaging) AddPubSubNamespace(topicName string) string {
	environment := ps.baseExt.GetRunningEnvironment()
	return ps.baseExt.NamespacePubsubIdentifier(
		ServiceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps ServicePubSubMessaging) TopicIDs() []string {
	return []string{
		ps.AddPubSubNamespace(common.CreateCustomerTopic),
		ps.AddPubSubNamespace(common.CreateSupplierTopic),
		ps.AddPubSubNamespace(common.CreateCRMContact),
		ps.AddPubSubNamespace(common.UpdateCRMContact),
		ps.AddPubSubNamespace(common.LinkCoverTopic),
	}
}

// PublishToPubsub sends a message to a specifeid Topic
func (ps ServicePubSubMessaging) PublishToPubsub(
	ctx context.Context,
	topicID string,
	payload []byte,
) error {
	environment, err := ps.baseExt.GoogleCloudProjectIDEnvVarName()
	if err != nil {
		return err
	}
	return ps.baseExt.PublishToPubsub(
		ctx,
		ps.client,
		topicID,
		environment,
		ServiceName,
		TopicVersion,
		payload,
	)
}

// EnsureTopicsExist creates the topic(s) in the suppplied list if they do not
// already exist.
func (ps ServicePubSubMessaging) EnsureTopicsExist(
	ctx context.Context,
	topicIDs []string,
) error {
	return ps.baseExt.EnsureTopicsExist(
		ctx,
		ps.client,
		topicIDs,
	)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (ps ServicePubSubMessaging) EnsureSubscriptionsExist(
	ctx context.Context,
) error {
	hostName, err := ps.baseExt.GetEnvVar(hostNameEnvVarName)
	if err != nil {
		return err
	}

	callbackURL := fmt.Sprintf(
		"%s%s",
		hostName,
		ps.baseExt.PubSubHandlerPath(),
	)

	return ps.baseExt.EnsureSubscriptionsExist(
		ctx,
		ps.client,
		ps.SubscriptionIDs(),
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (ps ServicePubSubMessaging) SubscriptionIDs() map[string]string {
	return ps.baseExt.SubscriptionIDs(ps.TopicIDs())
}
