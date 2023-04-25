package pubsubmessaging

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (
	// ClinicalServiceName defines the service where the topic is created
	ClinicalServiceName = "clinical"

	// MyCareHubServiceName defines the service where some of the topics have been created
	MyCareHubServiceName = "mycarehub"

	// TopicVersion defines the topic version. That standard one is `v1`
	TopicVersion = "v2"

	// HostNameEnvVarName defines the host name
	HostNameEnvVarName = "SERVICE_HOST"

	// TestTopicName is a topic that is used for testing purposes
	TestTopicName = "pubsub.mycarehub.testing.topic"
)

// ServicePubsub represent all the logic required to interact with pubsub
type ServicePubsub interface {
	PublishToPubsub(
		ctx context.Context,
		topicID string,
		serviceName string,
		payload []byte,
	) error
	ReceivePubSubPushMessages(
		w http.ResponseWriter,
		r *http.Request,
	)

	NotifyCreatePatient(ctx context.Context, client *dto.PatientCreationOutput) error
	NotifyCreateVitals(ctx context.Context, vitals *dto.PatientVitalSignOutput) error
	NotifyCreateAllergy(ctx context.Context, allergy *dto.PatientAllergyOutput) error
	NotifyCreateMedication(ctx context.Context, medication *dto.PatientMedicationOutput) error
	NotifyCreateTestOrder(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error
	NotifyCreateTestResult(ctx context.Context, testResult *dto.PatientTestResultOutput) error
	NotifyCreateOrganization(ctx context.Context, facility *domain.Facility) error

	NotifyCreateCMSClient(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error
	NotifyDeleteCMSClient(ctx context.Context, user *dto.DeleteCMSUserPayload) error
	NotifyCreateCMSProgram(ctx context.Context, program *dto.CreateCMSProgramPayload) error
	NotifyCreateCMSOrganisation(ctx context.Context, program *dto.CreateCMSOrganisationPayload) error
	NotifyCreateCMSFacility(ctx context.Context, facility *dto.CreateCMSFacilityPayload) error
	NotifyCMSAddFacilityToProgram(ctx context.Context, program *dto.CMSLinkFacilityToProgramPayload) error

	NotifyCreateClinicalTenant(ctx context.Context, tenant *dto.ClinicalTenantPayload) error
}

// ServicePubSubMessaging is used to send and receive pubsub notifications
type ServicePubSubMessaging struct {
	client  *pubsub.Client
	BaseExt extension.ExternalMethodsExtension
	Query   infrastructure.Query
	FCM     fcm.ServiceFCM
}

// NewServicePubSubMessaging returns a new instance of pubsub
func NewServicePubSubMessaging(
	baseExt extension.ExternalMethodsExtension,
	query infrastructure.Query,
	fcm fcm.ServiceFCM,
) (*ServicePubSubMessaging, error) {
	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}

	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	s := &ServicePubSubMessaging{
		client:  client,
		BaseExt: baseExt,
		Query:   query,
		FCM:     fcm,
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

// AddPubSubNamespace creates unique topics. The topics will be in the form
// <service name>-<topicName>-<environment>-v1
func (ps ServicePubSubMessaging) AddPubSubNamespace(topicName string, ServiceName string) string {
	environment := serverutils.GetRunningEnvironment()
	return pubsubtools.NamespacePubsubIdentifier(
		ServiceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps ServicePubSubMessaging) TopicIDs() []string {
	return []string{
		ps.AddPubSubNamespace(TestTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSClientTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.DeleteCMSClientTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.DeleteCMSStaffTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSStaffTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSProgramTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSOrganisationTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSFacilityTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.CreateCMSProgramFacilityTopicName, MyCareHubServiceName),
		ps.AddPubSubNamespace(common.TenantTopicName, ClinicalServiceName),
	}
}

// PublishToPubsub publishes a message to a specified topic
func (ps ServicePubSubMessaging) PublishToPubsub(
	ctx context.Context,
	topicID, serviceName string,
	payload []byte,
) error {
	environment, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return err
	}
	return ps.BaseExt.PublishToPubsub(
		ctx,
		ps.client,
		topicID,
		environment,
		serviceName,
		TopicVersion,
		payload,
	)
}

// EnsureTopicsExist creates the topic(s) in the supplied list if they do not
// exist
func (ps ServicePubSubMessaging) EnsureTopicsExist(
	ctx context.Context,
	topicIDs []string,
) error {
	return ps.BaseExt.EnsureTopicsExist(
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
	hostName, err := serverutils.GetEnvVar(HostNameEnvVarName)
	if err != nil {
		return err
	}

	callbackURL := fmt.Sprintf(
		"%s%s",
		hostName,
		pubsubtools.PubSubHandlerPath,
	)

	return ps.BaseExt.EnsureSubscriptionsExist(
		ctx,
		ps.client,
		ps.SubscriptionIDs(),
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (ps ServicePubSubMessaging) SubscriptionIDs() map[string]string {
	return pubsubtools.SubscriptionIDs(ps.TopicIDs())
}
