package utils

import (
	"encoding/json"
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

// FormatFilterParamsHelper is a helper function that formats the filter params from getstream
func FormatFilterParamsHelper(a map[string]interface{}) map[string]interface{} {
	supportedOperations := map[string]string{
		"lt":           "$lt",
		"lte":          "$lte",
		"gt":           "$gt",
		"gte":          "$gte",
		"eq":           "$eq",
		"ne":           "$ne",
		"in":           "$in",
		"autocomplete": "$autocomplete",
	}

	newMap := map[string]interface{}{}
	for k, v := range a {

		switch v := v.(type) {
		case map[string]interface{}:
			newMap[k] = FormatFilterParamsHelper(v)
		default:
			if supportedOperations[k] != "" {
				operation := supportedOperations[k]
				newMap[operation] = a[k]
			} else {
				newMap[k] = v
			}
		}
	}
	return newMap
}

// CalculateAge calculates the age of a community member
func CalculateAge(birthday time.Time) int {
	if birthday.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}
	return age
}

// ConvertJSONStringToMap converts a json string to a map of string interface{}
func ConvertJSONStringToMap(jsonString string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
