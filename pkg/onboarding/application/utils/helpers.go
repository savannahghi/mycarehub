package utils

import (
	"fmt"
	"hash"
	"math"
	"time"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

const (
	// Default min length of the date
	minPINLength = 4
	// Default max length of the date
	maxPINLength = 4
	//BackOffWaitTime is the default time to wait
	BackOffWaitTime = 3
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

// ValidatePIN ...
func ValidatePIN(pin string) error {
	validatePINErr := ValidatePINLength(pin)
	if validatePINErr != nil {
		return validatePINErr
	}

	pinDigitsErr := extension.ValidatePINDigits(pin)
	if pinDigitsErr != nil {
		return pinDigitsErr
	}
	return nil
}

// ValidatePINLength ...
func ValidatePINLength(pin string) error {
	// make sure pin length is [4]
	if len(pin) < minPINLength || len(pin) > maxPINLength {
		return exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
	}
	return nil
}

// GetHourMinuteSecond will be used to calculate time in the past or the future
func GetHourMinuteSecond(hour, minute, second time.Duration) time.Time {

	return time.Now().Add(time.Hour*hour + time.Minute*minute + time.Second*second)

}

// NextAllowedLoginTime calculates the next allowed time to login.
// It depends on the number of user's failed login attempts.
func NextAllowedLoginTime(trials int) time.Time {
	baseValue := float64(trials)
	result := math.Pow(baseValue, BackOffWaitTime)
	nextAllowedLoginTime := GetHourMinuteSecond(0, 0, time.Duration(result))
	return nextAllowedLoginTime
}
