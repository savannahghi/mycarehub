package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/common"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func (ps *ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data recieved: %v", err)
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
	updateData dto.UpdateContactPSMessage,
) error {
	return ps.newPublish(ctx, updateData, common.UpdateCRMContact)
}

// NotifyCreateCustomer publishes to customers.create topic
func (ps *ServicePubSubMessaging) NotifyCreateCustomer(
	ctx context.Context,
	data dto.CustomerPubSubMessage,
) error {
	return ps.newPublish(ctx, data, common.CreateCustomerTopic)
}

// NotifyCreateSupplier publishes to suppliers.create topic
func (ps *ServicePubSubMessaging) NotifyCreateSupplier(
	ctx context.Context,
	data dto.SupplierPubSubMessage,
) error {
	return ps.newPublish(ctx, data, common.CreateCustomerTopic)
}
