package utils

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jftuga/geodist"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

const (
	// BackOffWaitTime is the base value used for exponential backoff i.e 3**i where i is
	// number of failed login counts
	BackOffWaitTime = 3

	issuer      = "Savannah Informatics Limited"
	accountName = "info@healthcloud.co.ke"

	// OrganisationContextKey is used to add/retrieve the organisation ID in the context
	OrganisationContextKey = firebasetools.ContextKey("OrganisationID")

	// ProgramContextKey is used to add/retrieve the program ID that is passed around in the context
	ProgramContextKey = firebasetools.ContextKey("ProgramID")
)

// CalculateNextAllowedLoginTime will be used to calculate the next allowed login time in cases where
// a login failure occurs. If the failed login count is 1, then the next allowed login time
// will be after 3 seconds due to exponential backoff.
func CalculateNextAllowedLoginTime(hour, minute, second time.Duration) time.Time {
	return time.Now().UTC().Add(time.Hour*hour + time.Minute*minute + time.Second*second)
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
		"nin":          "$nin",
		"autocomplete": "$autocomplete",
		"q":            "$q",
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
	if jsonString == "" {
		return jsonMap, nil
	}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

// ConvertTimeToScalarDate converts a time to a scalar date
func ConvertTimeToScalarDate(t time.Time) (scalarutils.Date, error) {
	if t.IsZero() {
		return scalarutils.Date{}, fmt.Errorf("invalid time format")
	}
	return scalarutils.Date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}, nil
}

// InterfaceToInt converts an interface to an int
func InterfaceToInt(n interface{}) int {
	if i, ok := n.(int); ok {
		return i
	}
	return 0
}

// InterfaceToFloat64 converts an interface to float64
func InterfaceToFloat64(n interface{}) float64 {
	if f, ok := n.(float64); ok {
		return f
	}
	return 0.0
}

// InterfaceToString converts an interface to a string
func InterfaceToString(n interface{}) string {
	if n == nil {
		return ""
	}
	return n.(string)
}

// GenerateOTP is used to generate a one time password
func GenerateOTP() (string, error) {
	opts := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	}
	key, err := totp.Generate(opts)
	if err != nil {
		return "", errors.Wrap(err, "generateOTP")
	}

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", errors.Wrap(err, "generateOTP > GenerateCode")
	}

	return code, nil
}

// CheckIfCurrentDBIsLocal checks whether the database used to run the test is a test/local database. If not, the setup exits
func CheckIfCurrentDBIsLocal() bool {
	isLocal, err := strconv.ParseBool(os.Getenv("IS_LOCAL_DB"))
	if err != nil {
		return false
	}

	return isLocal
}

// TruncateMatrixUserID truncates the Matrix's user ID and returns the name
func TruncateMatrixUserID(userID string) string {
	parts := strings.Split(userID, ":")
	username := strings.TrimLeft(parts[0], "@")

	return username
}

// GetValueFromContext retrieves a value from the supplied context
func GetValueFromContext(ctx context.Context, key firebasetools.ContextKey) (string, error) {
	val := ctx.Value(key)
	if val == nil {
		return "", fmt.Errorf(
			"unable to get value from context with key %#v", key)
	}

	token, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("wrong value type, got %#v, expected a string", val)
	}

	return token, nil
}

// CalculateDistance is used to calculate the distance between two points on the earth give the starting point coordinates and destination point coordinates
func CalculateDistance(startPoint geodist.Coord, destination geodist.Coord) (float64, error) {
	_, km, err := geodist.VincentyDistance(startPoint, destination)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return 0, err
	}

	return km, nil
}

// ShuffleSecurityQuestionResponses is used to randomize users security questions.
// It uses crypto/rand package which makes use of a shuffling algorithm that provides higher level of randomness and security.
func ShuffleSecurityQuestionResponses(slice []*domain.SecurityQuestionResponse) ([]*domain.SecurityQuestionResponse, error) {
	n := len(slice)
	shuffled := make([]*domain.SecurityQuestionResponse, n)
	copy(shuffled, slice)

	for i := n - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return nil, err
		}

		randomIndex := int(j.Int64())
		shuffled[i], shuffled[randomIndex] = shuffled[randomIndex], shuffled[i]
	}

	return shuffled, nil
}
