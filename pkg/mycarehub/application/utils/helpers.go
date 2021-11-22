package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/savannahghi/serverutils"
)

const (
	// BackOffWaitTime is the base value used for exponential backoff i.e 3**i where i is
	// number of failed login counts
	BackOffWaitTime = 3

	// DjangoAuthorizationToken is used as an authorization token for making request to our
	// django backend service
	DjangoAuthorizationToken = "DJANGO_AUTHORIZATION_TOKEN"
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

// MakeRequest performs a http request and returns a response
func MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	token := serverutils.MustGetEnvVar(DjangoAuthorizationToken)
	client := http.Client{}
	// A GET request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if method == http.MethodGet {
		req, reqErr := http.NewRequestWithContext(ctx, method, path, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Authorization", "Token "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		return client.Do(req)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(encoded)
	req, reqErr := http.NewRequestWithContext(ctx, method, path, payload)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
