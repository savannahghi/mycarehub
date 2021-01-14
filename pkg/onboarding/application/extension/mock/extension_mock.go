package mock

import (
	"context"
)

// FakeBaseExtensionImpl is a `base` library mock .
type FakeBaseExtensionImpl struct {
	GetLoggedInUserUIDFn func(ctx context.Context) (string, error)
	NormalizeMSISDNFn    func(msisdn string) (*string, error)
}

// GetLoggedInUserUID ...
func (b *FakeBaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return b.GetLoggedInUserUIDFn(ctx)
}

// NormalizeMSISDN ...
func (b *FakeBaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return b.NormalizeMSISDNFn(msisdn)
}
