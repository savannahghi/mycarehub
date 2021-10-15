package mock

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

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
