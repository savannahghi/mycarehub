package extension

import (
	"context"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// BaseExtension is an interface that represents some methods in base
// The `onboarding` service has a dependency on `base` library.
// Our first step to making some functions are testable is to remove the base dependency.
// This can be achieved with the below interface.
type BaseExtension interface {
	// functions that we use from base
	GetLoggedInUserUID(ctx context.Context) (string, error)
	NormalizeMSISDN(msisdn string) (*string, error)
	FetchDefaultCurrency(c base.Client,
	) (*base.FinancialYearAndCurrency, error)
	LoginClient(username string, password string) (base.Client, error)
	FetchUserProfile(authClient base.Client) (*base.EDIUserProfile, error)
}

// BaseExtensionImpl ...
type BaseExtensionImpl struct {
}

// NewBaseExtensionImpl ...
func NewBaseExtensionImpl() BaseExtension {
	return &BaseExtensionImpl{}
}

// GetLoggedInUserUID get the logged in user uid
func (b *BaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return base.GetLoggedInUserUID(ctx)
}

// NormalizeMSISDN validates the input phone number.
func (b *BaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return base.NormalizeMSISDN(msisdn)
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

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(method string, path string, body interface{}) (*http.Response, error)
}

// ISCExtensionImpl ...
type ISCExtensionImpl struct{}

// NewISCExtension initializes an ISC extension
func NewISCExtension() ISCClientExtension {
	return &ISCExtensionImpl{}
}

// MakeRequest performs an inter service http request and returns a response
func (i *ISCExtensionImpl) MakeRequest(method string, path string, body interface{}) (*http.Response, error) {
	var isc base.InterServiceClient
	return isc.MakeRequest(method, path, body)
}
