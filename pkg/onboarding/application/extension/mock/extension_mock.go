package mock

import (
	"context"
	"net/http"

	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

// FakeBaseExtensionImpl is a `base` library fake  .
type FakeBaseExtensionImpl struct {
	GetLoggedInUserFn      func(ctx context.Context) (*dto.UserInfo, error)
	GetLoggedInUserUIDFn   func(ctx context.Context) (string, error)
	NormalizeMSISDNFn      func(msisdn string) (*string, error)
	FetchDefaultCurrencyFn func(c base.Client) (*base.FinancialYearAndCurrency, error)
	FetchUserProfileFn     func(authClient base.Client) (*base.EDIUserProfile, error)
	LoginClientFn          func(username string, password string) (base.Client, error)
	LoadDepsFromYAMLFn     func() (*base.DepsConfig, error)
	SetupISCclientFn       func(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error)
	GetEnvVarFn            func(envName string) (string, error)
	NewServerClientFn      func(
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
	EnsureTopicsExistFn func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicIDs []string,
	) error
	GetRunningEnvironmentFn     func() string
	NamespacePubsubIdentifierFn func(
		serviceName string,
		topicID string,
		environment string,
		version string,
	) string
	PublishToPubsubFn func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicID string,
		environment string,
		serviceName string,
		version string,
		payload []byte,
	) error
	GoogleCloudProjectIDEnvVarNameFn func() (string, error)
	EnsureSubscriptionsExistFn       func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicSubscriptionMap map[string]string,
		callbackURL string,
	) error
	SubscriptionIDsFn                 func(topicIDs []string) map[string]string
	PubSubHandlerPathFn               func() string
	VerifyPubSubJWTAndDecodePayloadFn func(
		w http.ResponseWriter,
		r *http.Request,
	) (*base.PubSubPayload, error)
	GetPubSubTopicFn    func(m *base.PubSubPayload) (string, error)
	ErrorMapFn          func(err error) map[string]string
	WriteJSONResponseFn func(
		w http.ResponseWriter,
		source interface{},
		status int,
	)
	GetLoginFuncFn                       func(ctx context.Context) http.HandlerFunc
	GetLogoutFuncFn                      func(ctx context.Context) http.HandlerFunc
	GetRefreshFuncFn                     func() http.HandlerFunc
	GetVerifyTokenFuncFn                 func(ctx context.Context) http.HandlerFunc
	GetUserProfileByPrimaryPhoneNumberFn func(ctx context.Context, phone string, suspended bool) (*base.UserProfile, error)
}

// GetLoggedInUser retrieves logged in user information
func (b *FakeBaseExtensionImpl) GetLoggedInUser(ctx context.Context) (*dto.UserInfo, error) {
	return b.GetLoggedInUserFn(ctx)
}

// GetLoggedInUserUID ...
func (b *FakeBaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return b.GetLoggedInUserUIDFn(ctx)
}

// NormalizeMSISDN ...
func (b *FakeBaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return b.NormalizeMSISDNFn(msisdn)
}

// FetchDefaultCurrency ...
func (b *FakeBaseExtensionImpl) FetchDefaultCurrency(c base.Client,
) (*base.FinancialYearAndCurrency, error) {
	return b.FetchDefaultCurrencyFn(c)
}

// FetchUserProfile ...
func (b *FakeBaseExtensionImpl) FetchUserProfile(authClient base.Client) (*base.EDIUserProfile, error) {
	return b.FetchUserProfileFn(authClient)
}

// LoginClient returns a logged in client with the supplied username and password
func (b *FakeBaseExtensionImpl) LoginClient(username, password string) (base.Client, error) {
	return b.LoginClientFn(username, password)
}

// LoadDepsFromYAML ...
func (b *FakeBaseExtensionImpl) LoadDepsFromYAML() (*base.DepsConfig, error) {
	return b.LoadDepsFromYAMLFn()
}

// SetupISCclient ...
func (b *FakeBaseExtensionImpl) SetupISCclient(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error) {
	return b.SetupISCclientFn(config, serviceName)
}

// GetEnvVar ...
func (b *FakeBaseExtensionImpl) GetEnvVar(envName string) (string, error) {
	return b.GetEnvVarFn(envName)
}

// GetLoginFunc ..
func (b *FakeBaseExtensionImpl) GetLoginFunc(ctx context.Context) http.HandlerFunc {
	return b.GetLoginFuncFn(ctx)
}

// GetLogoutFunc ..
func (b *FakeBaseExtensionImpl) GetLogoutFunc(ctx context.Context) http.HandlerFunc {
	return b.GetLogoutFuncFn(ctx)
}

// GetRefreshFunc ..
func (b *FakeBaseExtensionImpl) GetRefreshFunc() http.HandlerFunc {
	return b.GetRefreshFuncFn()
}

// GetVerifyTokenFunc ..
func (b *FakeBaseExtensionImpl) GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc {
	return b.GetVerifyTokenFuncFn(ctx)
}

// NewServerClient ...
func (b *FakeBaseExtensionImpl) NewServerClient(
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
	return b.NewServerClientFn(clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
}

// EnsureTopicsExist ...
func (b *FakeBaseExtensionImpl) EnsureTopicsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicIDs []string,
) error {
	return b.EnsureTopicsExistFn(ctx, pubsubClient, topicIDs)
}

// GetRunningEnvironment ..
func (b *FakeBaseExtensionImpl) GetRunningEnvironment() string {
	return b.GetRunningEnvironmentFn()
}

// NamespacePubsubIdentifier ..
func (b *FakeBaseExtensionImpl) NamespacePubsubIdentifier(
	serviceName string,
	topicID string,
	environment string,
	version string,
) string {
	return b.NamespacePubsubIdentifierFn(
		serviceName,
		topicID,
		environment,
		version,
	)
}

// PublishToPubsub ..
func (b *FakeBaseExtensionImpl) PublishToPubsub(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicID string,
	environment string,
	serviceName string,
	version string,
	payload []byte,
) error {
	return b.PublishToPubsubFn(
		ctx,
		pubsubClient,
		topicID,
		environment,
		serviceName,
		version,
		payload,
	)
}

// GoogleCloudProjectIDEnvVarName ..
func (b *FakeBaseExtensionImpl) GoogleCloudProjectIDEnvVarName() (string, error) {
	return b.GoogleCloudProjectIDEnvVarNameFn()
}

// EnsureSubscriptionsExist ...
func (b *FakeBaseExtensionImpl) EnsureSubscriptionsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicSubscriptionMap map[string]string,
	callbackURL string,
) error {
	return b.EnsureSubscriptionsExistFn(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// SubscriptionIDs ..
func (b *FakeBaseExtensionImpl) SubscriptionIDs(topicIDs []string) map[string]string {
	return b.SubscriptionIDsFn(topicIDs)
}

// PubSubHandlerPath ..
func (b *FakeBaseExtensionImpl) PubSubHandlerPath() string {
	return b.PubSubHandlerPathFn()
}

// VerifyPubSubJWTAndDecodePayload ..
func (b *FakeBaseExtensionImpl) VerifyPubSubJWTAndDecodePayload(
	w http.ResponseWriter,
	r *http.Request,
) (*base.PubSubPayload, error) {
	return b.VerifyPubSubJWTAndDecodePayloadFn(w, r)
}

// GetPubSubTopic ..
func (b *FakeBaseExtensionImpl) GetPubSubTopic(m *base.PubSubPayload) (string, error) {
	return b.GetPubSubTopicFn(m)
}

// ErrorMap ..
func (b *FakeBaseExtensionImpl) ErrorMap(err error) map[string]string {
	return b.ErrorMapFn(err)
}

// WriteJSONResponse ..
func (b *FakeBaseExtensionImpl) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
	b.WriteJSONResponseFn(w, source, status)
}

// PINExtensionImpl is a `PIN` fake  .
type PINExtensionImpl struct {
	EncryptPINFn      func(rawPwd string, options *extension.Options) (string, string)
	ComparePINFn      func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	GenerateTempPINFn func(ctx context.Context) (string, error)
}

// EncryptPIN ...
func (p *PINExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return p.EncryptPINFn(rawPwd, options)
}

// ComparePIN ...
func (p *PINExtensionImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return p.ComparePINFn(rawPwd, salt, encodedPwd, options)
}

// GenerateTempPIN ...
func (p *PINExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return p.GenerateTempPINFn(ctx)
}

// ISCClientExtension is an ISC fake
type ISCClientExtension struct {
	MakeRequestFn func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// MakeRequest ...
func (i *ISCClientExtension) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	return i.MakeRequestFn(ctx, method, path, body)
}

// GetUserProfileByPrimaryPhoneNumber ..
func (b *FakeBaseExtensionImpl) GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phone string, suspended bool) (*base.UserProfile, error) {
	return b.GetUserProfileByPrimaryPhoneNumberFn(ctx, phone, suspended)
}

// CRMExtensionImpl is a fake CRM extension implementation
type CRMExtensionImpl struct {
	CreateContactFn func(contact domain.CRMContact) (*domain.CRMContact, error)
	UpdateContactFn func(
		phone string,
		properties domain.ContactProperties,
	) (*domain.CRMContact, error)
}

// CreateContact ..
func (c *CRMExtensionImpl) CreateContact(contact domain.CRMContact) (*domain.CRMContact, error) {
	return c.CreateContactFn(contact)
}

// UpdateContact ..
func (c *CRMExtensionImpl) UpdateContact(
	phone string,
	properties domain.ContactProperties,
) (*domain.CRMContact, error) {
	return c.UpdateContactFn(phone, properties)
}
