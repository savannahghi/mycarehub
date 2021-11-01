package utils

import (
	"fmt"
	"math"
	"time"

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
