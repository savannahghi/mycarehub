package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
)

const (
	// ServiceName ..
	ServiceName = "onboarding"

	// TopicVersion ...
	TopicVersion = "v1"

	// CreateCustomerTopic is the TopicID for customer creation Topic
	CreateCustomerTopic = "customers.create"

	// CreateSupplierTopic is the TopicID for supplier creation Topic
	CreateSupplierTopic = "suppliers.create"

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
}

// ServicePubSubMessaging sends "real" (production) notifications
type ServicePubSubMessaging struct {
	client  *pubsub.Client
	baseExt extension.BaseExtension
	erp     erp.ServiceERP
}

// NewServicePubSubMessaging ...
func NewServicePubSubMessaging(
	client *pubsub.Client,
	ext extension.BaseExtension,
	erp erp.ServiceERP,
) (*ServicePubSubMessaging, error) {
	s := &ServicePubSubMessaging{
		client:  client,
		baseExt: ext,
		erp:     erp,
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
func (ps ServicePubSubMessaging) AddEngagementPubsubNameSpace(topic string) string {
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
		ps.AddPubSubNamespace(CreateCustomerTopic),
		ps.AddPubSubNamespace(CreateSupplierTopic),
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

// ReceivePubSubPushMessages receives and processes a Pub/Sub push message.
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	message, err := ps.baseExt.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	topicID, err := ps.baseExt.GetPubSubTopic(message)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	ctx := r.Context()
	switch topicID {
	case ps.AddPubSubNamespace(CreateCustomerTopic):
		var data resources.CustomerPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err := ps.erp.CreateERPCustomer(
			ctx,
			data.CustomerPayload,
			data.UID,
		); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	case ps.AddPubSubNamespace(CreateSupplierTopic):
		var data resources.SupplierPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err := ps.erp.CreateERPSupplier(
			ctx,
			data.SupplierPayload,
			data.UID,
		); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	default:
		errMsg := fmt.Sprintf(
			"pub sub handler error: unknown topic `%s`",
			topicID,
		)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	resp := map[string]string{"status": "success"}
	marshalledSuccessMsg, err := json.Marshal(resp)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusInternalServerError,
		)
		return
	}
	_, _ = w.Write(marshalledSuccessMsg)
}
