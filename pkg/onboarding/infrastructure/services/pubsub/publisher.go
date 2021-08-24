package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/common"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

func (ps *ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data received: %v", err)
	}
	return ps.PublishToPubsub(
		ctx,
		ps.AddPubSubNamespace(topic),
		payload,
	)
}

// NotifyCreateContact publishes to crm.contact.create topic
func (ps *ServicePubSubMessaging) NotifyCreateContact(
	ctx context.Context,
	contact domain.CRMContact,
) error {
	return ps.newPublish(ctx, contact, common.CreateCRMContact)
}

// NotifyUpdateContact publishes to crm.contact.update topic
func (ps *ServicePubSubMessaging) NotifyUpdateContact(
	ctx context.Context,
	contact domain.CRMContact,
) error {
	return ps.newPublish(ctx, contact, common.UpdateCRMContact)
}

// NotifyCoverLinking pushes to covers.link topic
func (ps *ServicePubSubMessaging) NotifyCoverLinking(
	ctx context.Context,
	data dto.LinkCoverPubSubMessage,
) error {
	return ps.newPublish(ctx, data, common.LinkCoverTopic)
}

// EDIMemberCoverLinking publishes to the edi.covers.link topic. The reason for this is
// to Auto-link the Sladers who get text messages from EDI. If a slader is converted
// and creates an account on Be.Well app, we should automatically append a cover to their profile.
func (ps *ServicePubSubMessaging) EDIMemberCoverLinking(
	ctx context.Context,
	data dto.LinkCoverPubSubMessage,
) error {
	return ps.newPublish(ctx, data, common.LinkEDIMemberCoverTopic)
}
