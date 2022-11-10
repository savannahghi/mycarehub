package extension

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
)

const (
	// DjangoAuthorizationToken is used as an authorization token for making request to our
	// django backend service
	DjangoAuthorizationToken = "DJANGO_AUTHORIZATION_TOKEN"
)

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// ExternalMethodsExtension is an interface that represents methods that are
// called from external libraries. Adding this layer will help write unit tests
type ExternalMethodsExtension interface {
	CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error)
	AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	GetLoggedInUserUID(ctx context.Context) (string, error)
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)

	Login(ctx context.Context) http.HandlerFunc

	PublishToPubsub(ctx context.Context, pubsubClient *pubsub.Client, topicID string, environment string, serviceName string, version string, payload []byte) error
	EnsureTopicsExist(ctx context.Context, pubsubClient *pubsub.Client, topicIDs []string) error
	EnsureSubscriptionsExist(ctx context.Context, pubsubClient *pubsub.Client, topicSubscriptionMap map[string]string, callbackURL string) error
	VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)

	LoadDepsFromYAML() (*interserviceclient.DepsConfig, error)
	SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
}

// External type implements external methods
type External struct {
}

// NewExternalMethodsImpl creates a new instance of the external methods
func NewExternalMethodsImpl() ExternalMethodsExtension {
	return &External{}
}

// CreateFirebaseCustomToken creates a custom auth token for the user with the
// indicated UID
func (e *External) CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error) {
	return firebasetools.CreateFirebaseCustomToken(ctx, uid)
}

// AuthenticateCustomFirebaseToken takes a custom Firebase auth token and tries to fetch an ID token
// If successful, a pointer to the ID token is returned
// Otherwise, an error is returned
func (e *External) AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
	return firebasetools.AuthenticateCustomFirebaseToken(customAuthToken)
}

// GetLoggedInUserUID get the logged in user uid
func (e *External) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return firebasetools.GetLoggedInUserUID(ctx)
}

// MakeRequest performs a http request and returns a response
func (e *External) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	token := serverutils.MustGetEnvVar(DjangoAuthorizationToken)
	client := http.Client{}
	// A GET request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if method == http.MethodGet {
		req, reqErr := http.NewRequestWithContext(ctx, method, path, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Authorization", "Token "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		return client.Do(req)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(encoded)
	req, reqErr := http.NewRequestWithContext(ctx, method, path, payload)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

// Login authenticates against firebase to return a valid token
func (e *External) Login(ctx context.Context) http.HandlerFunc {
	return firebasetools.GetLoginFunc(ctx, &firebasetools.FirebaseClient{})
}

// PublishToPubsub sends the supplied payload to the indicated topic
func (e *External) PublishToPubsub(ctx context.Context, pubsubClient *pubsub.Client, topicID string, environment string, serviceName string, version string, payload []byte) error {
	return pubsubtools.PublishToPubsub(
		ctx,
		pubsubClient,
		topicID,
		environment,
		serviceName,
		version,
		payload,
	)
}

// EnsureTopicsExist creates the topic(s) in the suppplied list if they do not
// already exist.
func (e *External) EnsureTopicsExist(ctx context.Context, pubsubClient *pubsub.Client, topicIDs []string) error {
	return pubsubtools.EnsureTopicsExist(ctx, pubsubClient, topicIDs)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (e *External) EnsureSubscriptionsExist(ctx context.Context, pubsubClient *pubsub.Client, topicSubscriptionMap map[string]string, callbackURL string) error {
	return pubsubtools.EnsureSubscriptionsExist(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// LoadDepsFromYAML loads the dependency config
func (*External) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return interserviceclient.LoadDepsFromYAML()
}

// SetupISCclient creates an isc client
func (*External) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return interserviceclient.SetupISCclient(config, serviceName)
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (e *External) VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
	return pubsubtools.VerifyPubSubJWTAndDecodePayload(w, r)
}

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string, ext ExternalMethodsExtension) *interserviceclient.InterServiceClient {
	config, err := ext.LoadDepsFromYAML()
	if err != nil {
		logrus.Panicf("occurred while opening deps file %v", err)
		return nil
	}

	client, err := ext.SetupISCclient(*config, serviceName)
	if err != nil {
		logrus.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
		return nil
	}
	return client
}
