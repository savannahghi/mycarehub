package utils

import (
	"math"
	"time"
)

const (
	// BackOffWaitTime is the base value used for exponential backoff i.e 3**i where i is
	// number of failed login counts
	BackOffWaitTime = 3
)

// CalculateNextAllowedLoginTime will be used to calculate the next allowed login time in cases where
// a login failure occurs. If the failed login count is 1, then the next allowed login time
// will be after 3 seconds due to exponential backoff.
func CalculateNextAllowedLoginTime(hour, minute, second time.Duration) time.Time {
	return time.Now().Add(time.Hour*hour + time.Minute*minute + time.Second*second)
}

// NextAllowedLoginTime calculates the user's next allowed login time.
// This implements an exponential backoff that will help to mitigate brute force attacks.
// Each incorrect PIN entry will add a delay, hence increasing the time to the next allowed login
func NextAllowedLoginTime(loginAttempts int) time.Time {
	baseValue := float64(loginAttempts)
	result := math.Pow(BackOffWaitTime, baseValue)
	nextAllowedLoginTime := CalculateNextAllowedLoginTime(0, 0, time.Duration(result))
	return nextAllowedLoginTime
}
