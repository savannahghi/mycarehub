package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

func (ps *ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic, serviceName string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data received: %v", err)
	}
	return ps.PublishToPubsub(
		ctx,
		ps.AddPubSubNamespace(topic, serviceName),
		serviceName,
		payload,
	)
}

// NotifyCreatePatient publishes to the create patient topic
func (ps ServicePubSubMessaging) NotifyCreatePatient(ctx context.Context, client *dto.ClientRegistrationOutput) error {
	return ps.newPublish(ctx, client, common.CreatePatientTopic, ClinicalServiceName)
}
