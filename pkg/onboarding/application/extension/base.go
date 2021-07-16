package extension

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/base"
)

// BaseExtension is an interface that represents some methods in base
// The `onboarding` service has a dependency on `base` library.
// Our first step to making some functions are testable is to remove the base dependency.
// This can be achieved with the below interface.
type BaseExtension interface {
	GetLoggedInUser(ctx context.Context) (*dto.UserInfo, error)
	GetLoggedInUserUID(ctx context.Context) (string, error)
	NormalizeMSISDN(msisdn string) (*string, error)
	FetchDefaultCurrency(c base.Client,
	) (*base.FinancialYearAndCurrency, error)
	LoginClient(username string, password string) (base.Client, error)
	FetchUserProfile(authClient base.Client) (*base.EDIUserProfile, error)
	LoadDepsFromYAML() (*base.DepsConfig, error)
	SetupISCclient(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error)
	GetEnvVar(envName string) (string, error)
	NewServerClient(
		clientID string,
		clientSecret string,
		apiTokenURL string,
		apiHost string,
		apiScheme string,
		grantType string,
		username string,
		password string,
		extraHeaders map[string]string,
	) (*base.ServerClient, error)

	// PubSub
	EnsureTopicsExist(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicIDs []string,
	) error
	GetRunningEnvironment() string
	NamespacePubsubIdentifier(
		serviceName string,
		topicID string,
		environment string,
		version string,
	) string
	PublishToPubsub(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicID string,
		environment string,
		serviceName string,
		version string,
		payload []byte,
	) error
	GoogleCloudProjectIDEnvVarName() (string, error)
	EnsureSubscriptionsExist(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicSubscriptionMap map[string]string,
		callbackURL string,
	) error
	SubscriptionIDs(topicIDs []string) map[string]string
	PubSubHandlerPath() string
	VerifyPubSubJWTAndDecodePayload(
		w http.ResponseWriter,
		r *http.Request,
	) (*base.PubSubPayload, error)
	GetPubSubTopic(m *base.PubSubPayload) (string, error)
	ErrorMap(err error) map[string]string
	WriteJSONResponse(
		w http.ResponseWriter,
		source interface{},
		status int,
	)

	// Login
	GetLoginFunc(ctx context.Context) http.HandlerFunc
	GetLogoutFunc(ctx context.Context) http.HandlerFunc
	GetRefreshFunc() http.HandlerFunc
	GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc
}

// BaseExtensionImpl ...
type BaseExtensionImpl struct {
	fc base.IFirebaseClient
}

// NewBaseExtensionImpl ...
func NewBaseExtensionImpl(fc base.IFirebaseClient) BaseExtension {
	return &BaseExtensionImpl{
		fc: fc,
	}
}

// GetLoggedInUser retrieves logged in user information
func (b *BaseExtensionImpl) GetLoggedInUser(ctx context.Context) (*dto.UserInfo, error) {
	authToken, err := base.GetUserTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("user auth token not found in context: %w", err)
	}

	authClient, err := base.GetFirebaseAuthClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get or create Firebase client: %w", err)
	}

	user, err := authClient.GetUser(ctx, authToken.UID)
	if err != nil {

		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	return &dto.UserInfo{
		UID:         user.UID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		DisplayName: user.DisplayName,
		ProviderID:  user.ProviderID,
		PhotoURL:    user.PhotoURL,
	}, nil
}

// GetLoggedInUserUID get the logged in user uid
func (b *BaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return base.GetLoggedInUserUID(ctx)
}

// NormalizeMSISDN validates the input phone number.
func (b *BaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return converterandformatter.NormalizeMSISDN(msisdn)
}

// FetchDefaultCurrency fetched an ERP's organization's default
// current currency
func (b *BaseExtensionImpl) FetchDefaultCurrency(c base.Client,
) (*base.FinancialYearAndCurrency, error) {
	return base.FetchDefaultCurrency(c)
}

// LoginClient returns a logged in client with the supplied username and password
func (b *BaseExtensionImpl) LoginClient(username, password string) (base.Client, error) {
	return base.LoginClient(username, password)
}

// FetchUserProfile ...
func (b *BaseExtensionImpl) FetchUserProfile(authClient base.Client) (*base.EDIUserProfile, error) {
	return base.FetchUserProfile(authClient)
}

// LoadDepsFromYAML ...
func (b *BaseExtensionImpl) LoadDepsFromYAML() (*base.DepsConfig, error) {
	return base.LoadDepsFromYAML()
}

// SetupISCclient ...
func (b *BaseExtensionImpl) SetupISCclient(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error) {
	return base.SetupISCclient(config, serviceName)
}

// GetEnvVar ...
func (b *BaseExtensionImpl) GetEnvVar(envName string) (string, error) {
	return serverutils.GetEnvVar(envName)
}

// GetLoginFunc returns a function that can authenticate against both Slade 360 and Firebase
func (b *BaseExtensionImpl) GetLoginFunc(ctx context.Context) http.HandlerFunc {
	return base.GetLoginFunc(ctx, b.fc)
}

// GetLogoutFunc logs the user out of Firebase
func (b *BaseExtensionImpl) GetLogoutFunc(ctx context.Context) http.HandlerFunc {
	return base.GetLogoutFunc(ctx, b.fc)
}

// GetRefreshFunc is used to refresh OAuth tokens
func (b *BaseExtensionImpl) GetRefreshFunc() http.HandlerFunc {
	return base.GetRefreshFunc()
}

// GetVerifyTokenFunc confirms that an EDI access token (supplied) is valid.
// If it is valid, it exchanges it for a Firebase ID token.
func (b *BaseExtensionImpl) GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc {
	return base.GetVerifyTokenFunc(ctx, b.fc)
}

// NewServerClient ...
func (b *BaseExtensionImpl) NewServerClient(
	clientID string,
	clientSecret string,
	apiTokenURL string,
	apiHost string,
	apiScheme string,
	grantType string,
	username string,
	password string,
	extraHeaders map[string]string,
) (*base.ServerClient, error) {
	return base.NewServerClient(
		clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
}

// EnsureTopicsExist creates the topic(s) in the suppplied list if they do not
// already exist.
func (b *BaseExtensionImpl) EnsureTopicsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicIDs []string,
) error {
	return base.EnsureTopicsExist(ctx, pubsubClient, topicIDs)
}

// GetRunningEnvironment returns the environment wheere the service is running. Importannt
// so as to point to the correct deps
func (b *BaseExtensionImpl) GetRunningEnvironment() string {
	return serverutils.GetRunningEnvironment()
}

// NamespacePubsubIdentifier uses the service name, environment and version to
// create a "namespaced" pubsub identifier. This could be a topicID or
// subscriptionID.
func (b *BaseExtensionImpl) NamespacePubsubIdentifier(
	serviceName string,
	topicID string,
	environment string,
	version string,
) string {
	return base.NamespacePubsubIdentifier(
		serviceName,
		topicID,
		environment,
		version,
	)
}

// PublishToPubsub sends the supplied payload to the indicated topic
func (b *BaseExtensionImpl) PublishToPubsub(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicID string,
	environment string,
	serviceName string,
	version string,
	payload []byte,
) error {
	return base.PublishToPubsub(
		ctx,
		pubsubClient,
		topicID,
		environment,
		serviceName,
		version,
		payload,
	)
}

// GoogleCloudProjectIDEnvVarName returns `GOOGLE_CLOUD_PROJECT` env
func (b *BaseExtensionImpl) GoogleCloudProjectIDEnvVarName() (string, error) {
	return b.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (b *BaseExtensionImpl) EnsureSubscriptionsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicSubscriptionMap map[string]string,
	callbackURL string,
) error {
	return base.EnsureSubscriptionsExist(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (b *BaseExtensionImpl) SubscriptionIDs(topicIDs []string) map[string]string {
	return base.SubscriptionIDs(topicIDs)
}

// PubSubHandlerPath returns pubsub hander path `/pubsub`
func (b *BaseExtensionImpl) PubSubHandlerPath() string {
	return base.PubSubHandlerPath
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (b *BaseExtensionImpl) VerifyPubSubJWTAndDecodePayload(
	w http.ResponseWriter,
	r *http.Request,
) (*base.PubSubPayload, error) {
	return base.VerifyPubSubJWTAndDecodePayload(
		w,
		r,
	)
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
func (b *BaseExtensionImpl) GetPubSubTopic(m *base.PubSubPayload) (string, error) {
	return base.GetPubSubTopic(m)
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.
func (b *BaseExtensionImpl) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
	serverutils.WriteJSONResponse(w, source, status)
}

// ErrorMap turns the supplied error into a map with "error" as the key
func (b *BaseExtensionImpl) ErrorMap(err error) map[string]string {
	return serverutils.ErrorMap(err)
}

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// ISCExtensionImpl ...
type ISCExtensionImpl struct{}

// NewISCExtension initializes an ISC extension
func NewISCExtension() ISCClientExtension {
	return &ISCExtensionImpl{}
}

// MakeRequest performs an inter service http request and returns a response
func (i *ISCExtensionImpl) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	var isc base.InterServiceClient
	return isc.MakeRequest(ctx, method, path, body)
}
