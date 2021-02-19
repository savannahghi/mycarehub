package mock

import (
	"context"
	"net/http"
)

// FakeServicePubSub ...
type FakeServicePubSub struct {
	TopicIDsFn           func() []string
	AddPubSubNamespaceFn func(topicName string) string
	PublishToPubsubFn    func(
		ctx context.Context,
		topicID string,
		payload []byte,
	) error
	EnsureTopicsExistFn func(
		ctx context.Context,
		topicIDs []string,
	) error
	EnsureSubscriptionsExistFn func(
		ctx context.Context,
	) error
	SubscriptionIDsFn           func() map[string]string
	ReceivePubSubPushMessagesFn func(
		w http.ResponseWriter,
		r *http.Request,
	)
}

// AddPubSubNamespace ...
func (m *FakeServicePubSub) AddPubSubNamespace(topicName string) string {
	return m.AddPubSubNamespaceFn(topicName)
}

// TopicIDs returns the known (registered) topic IDs
func (m *FakeServicePubSub) TopicIDs() []string {
	return m.TopicIDsFn()
}

// PublishToPubsub ..
func (m *FakeServicePubSub) PublishToPubsub(
	ctx context.Context,
	topicID string,
	payload []byte,
) error {
	return m.PublishToPubsubFn(
		ctx,
		topicID,
		payload,
	)
}

// EnsureTopicsExist ...
func (m *FakeServicePubSub) EnsureTopicsExist(
	ctx context.Context,
	topicIDs []string,
) error {
	return m.EnsureTopicsExistFn(ctx, topicIDs)
}

// EnsureSubscriptionsExist ..
func (m *FakeServicePubSub) EnsureSubscriptionsExist(
	ctx context.Context,
) error {
	return m.EnsureSubscriptionsExistFn(ctx)
}

// SubscriptionIDs ..
func (m *FakeServicePubSub) SubscriptionIDs() map[string]string {
	return m.SubscriptionIDsFn()
}

// ReceivePubSubPushMessages ...
func (m *FakeServicePubSub) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	m.ReceivePubSubPushMessagesFn(w, r)
}
