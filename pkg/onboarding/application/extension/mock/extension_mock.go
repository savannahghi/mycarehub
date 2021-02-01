package mock

import (
	"context"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

// FakeBaseExtensionImpl is a `base` library fake  .
type FakeBaseExtensionImpl struct {
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

// PINExtensionImpl is a `PIN` fake  .
type PINExtensionImpl struct {
	EncryptPINFn func(rawPwd string, options *extension.Options) (string, string)
	ComparePINFn func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
}

// EncryptPIN ...
func (p *PINExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return p.EncryptPINFn(rawPwd, options)
}

// ComparePIN ...
func (p *PINExtensionImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return p.ComparePINFn(rawPwd, salt, encodedPwd, options)
}

// ISCClientExtension is an ISC fake
type ISCClientExtension struct {
	MakeRequestFn func(method string, path string, body interface{}) (*http.Response, error)
}

// MakeRequest ...
func (i *ISCClientExtension) MakeRequest(method string, path string, body interface{}) (*http.Response, error) {
	return i.MakeRequestFn(method, path, body)
}
