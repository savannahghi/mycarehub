package utils

import (
	"hash"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

// Options is a struct for custom values of salt length, number of iterations, the encoded key's length,
// and the hash function being used. If set to `nil`, default options are used:
// &Options{ 256, 10000, 512, "sha512" }
type Options struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

// EncryptUID takes two arguments, a raw uid, and a pointer to an Options struct.
// In order to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func EncryptUID(rawUID string, options *Options) (string, string) {
	interactor := infrastructure.NewInteractor()
	return interactor.PINExtension.EncryptPIN(rawUID, nil)
}

// CompareUID takes four arguments, the raw UID, its generated salt, the encoded UID,
// and a pointer to the Options struct, and returns a boolean value determining whether the UID is the correct one or not.
// Passing `nil` as the last argument resorts to default options.
func CompareUID(rawUID string, salt string, encodedUID string, options *Options) bool {

	interactor := infrastructure.NewInteractor()
	return interactor.PINExtension.ComparePIN(rawUID, salt, encodedUID, nil)
}
