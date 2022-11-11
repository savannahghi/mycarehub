package mock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
)

// FakeExtensionImpl mocks the external calls logic
type FakeExtensionImpl struct {
	MockCreateFirebaseCustomTokenFn           func(ctx context.Context, uid string) (string, error)
	MockAuthenticateCustomFirebaseTokenFn     func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	MockGetLoggedInUserUIDFn                  func(ctx context.Context) (string, error)
	MockMakeRequestFn                         func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
	MockLoginFn                               func(ctx context.Context) http.HandlerFunc
	MockNamespacePubsubIdentifierFn           func(serviceName string, topicID string, environment string, version string) string
	MockPublishToPubsubFn                     func(ctx context.Context, pubsubClient *pubsub.Client, topicID string, environment string, serviceName string, version string, payload []byte) error
	MockEnsureTopicsExistFn                   func(ctx context.Context, pubsubClient *pubsub.Client, topicIDs []string) error
	MockEnsureSubscriptionsExistFn            func(ctx context.Context, pubsubClient *pubsub.Client, topicSubscriptionMap map[string]string, callbackURL string) error
	MockVerifyPubSubJWTAndDecodePayloadFn     func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)
	MockLoadDepsFromYAMLFn                    func() (*interserviceclient.DepsConfig, error)
	MockSetupISCclientFn                      func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	MockCreateFirebaseCustomTokenWithClaimsFn func(ctx context.Context, uid string, claims map[string]interface{}) (string, error)
}

// NewFakeExtension initializes a new instance of the external calls mock
func NewFakeExtension() *FakeExtensionImpl {
	return &FakeExtensionImpl{
		MockLoadDepsFromYAMLFn: func() (*interserviceclient.DepsConfig, error) {
			return &interserviceclient.DepsConfig{
				Staging: []interserviceclient.Dep{
					{
						DepName:       "clinical",
						DepRootDomain: "https://clinical",
					},
				},
				Testing: []interserviceclient.Dep{
					{
						DepName:       "clinical",
						DepRootDomain: "https://clinical",
					},
				},
				Production: []interserviceclient.Dep{
					{
						DepName:       "clinical",
						DepRootDomain: "https://clinical",
					},
				},
			}, nil
		},
		MockSetupISCclientFn: func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
			return &interserviceclient.InterServiceClient{
				Name:              serviceName,
				RequestRootDomain: "https://clinical",
			}, nil
		},

		MockCreateFirebaseCustomTokenFn: func(ctx context.Context, uid string) (string, error) {
			return uuid.New().String(), nil
		},
		MockCreateFirebaseCustomTokenWithClaimsFn: func(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
			return uuid.New().String(), nil
		},
		MockAuthenticateCustomFirebaseTokenFn: func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
			return &firebasetools.FirebaseUserTokens{
				IDToken:      uuid.New().String(),
				RefreshToken: uuid.NewString(),
				ExpiresIn:    "1000",
			}, nil
		},
		MockGetLoggedInUserUIDFn: func(ctx context.Context) (string, error) {
			return uuid.New().String(), nil
		},
		MockMakeRequestFn: func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
			msg := struct {
				Message string `json:"message"`
			}{
				Message: "success",
			}

			payload, _ := json.Marshal(msg)

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBuffer(payload)),
			}, nil
		},
		MockPublishToPubsubFn: func(ctx context.Context, pubsubClient *pubsub.Client, topicID string, environment string, serviceName string, version string, payload []byte) error {
			return nil
		},
		MockEnsureTopicsExistFn: func(ctx context.Context, pubsubClient *pubsub.Client, topicIDs []string) error {
			return nil
		},
		MockEnsureSubscriptionsExistFn: func(ctx context.Context, pubsubClient *pubsub.Client, topicSubscriptionMap map[string]string, callbackURL string) error {
			return nil
		},
		MockVerifyPubSubJWTAndDecodePayloadFn: func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
			return &pubsubtools.PubSubPayload{
				Message: pubsubtools.PubSubMessage{
					Attributes: map[string]string{
						"topicID": "test-id",
					},
				},
				Subscription: "test-subscription",
			}, nil
		},
	}
}

// CreateFirebaseCustomToken mocks the create firebase custom token method
func (f *FakeExtensionImpl) CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error) {
	return f.MockCreateFirebaseCustomTokenFn(ctx, uid)
}

// AuthenticateCustomFirebaseToken mocks the authenticate custom firebase token method
func (f *FakeExtensionImpl) AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
	return f.MockAuthenticateCustomFirebaseTokenFn(customAuthToken)
}

// GetLoggedInUserUID mocks the implementation of getting a logged in user
func (f *FakeExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return f.MockGetLoggedInUserUIDFn(ctx)
}

// MakeRequest mocks the implementation of making a http request
func (f *FakeExtensionImpl) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	return f.MockMakeRequestFn(ctx, method, path, body)
}

// Login mocks the login implementation to retrieve a token
func (f *FakeExtensionImpl) Login(ctx context.Context) http.HandlerFunc {
	return f.MockLoginFn(ctx)
}

// PublishToPubsub sends the supplied payload to the indicated topic
func (f *FakeExtensionImpl) PublishToPubsub(ctx context.Context, pubsubClient *pubsub.Client, topicID string, environment string, serviceName string, version string, payload []byte) error {
	return f.MockPublishToPubsubFn(
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
func (f *FakeExtensionImpl) EnsureTopicsExist(ctx context.Context, pubsubClient *pubsub.Client, topicIDs []string) error {
	return f.MockEnsureTopicsExistFn(ctx, pubsubClient, topicIDs)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (f *FakeExtensionImpl) EnsureSubscriptionsExist(ctx context.Context, pubsubClient *pubsub.Client, topicSubscriptionMap map[string]string, callbackURL string) error {
	return f.MockEnsureSubscriptionsExistFn(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (f *FakeExtensionImpl) VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
	return f.MockVerifyPubSubJWTAndDecodePayloadFn(w, r)
}

// LoadDepsFromYAML loads the dependency config
func (f *FakeExtensionImpl) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return f.MockLoadDepsFromYAMLFn()
}

// SetupISCclient creates an isc client
func (f *FakeExtensionImpl) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return f.MockSetupISCclientFn(config, serviceName)
}

// CreateFirebaseCustomTokenWithClaims creates a custom auth token for the user with the
// indicated UID and additional claims
func (f *FakeExtensionImpl) CreateFirebaseCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	return f.MockCreateFirebaseCustomTokenWithClaimsFn(ctx, uid, claims)
}
