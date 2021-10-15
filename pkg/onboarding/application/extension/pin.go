package extension

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// EncryptPIN encrypts the clients provided PIN
func EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	interactor := infrastructure.NewInteractor()
	return interactor.PINExtension.EncryptPIN(rawPwd, nil)
}

// ComparePIN compares the the entered PINs during sign up to check whether they are a match
func ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	interactor := infrastructure.NewInteractor()
	return interactor.PINExtension.ComparePIN(rawPwd, salt, encodedPwd, nil)
}
