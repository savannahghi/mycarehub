package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

const (
	// ProInviteLink will store the pro app invite link at the settings
	ProInviteLink = "PRO_INVITE_LINK"

	// ConsumerInviteLink will store the consumer invite link at the settings
	ConsumerInviteLink = "CONSUMER_INVITE_LINK"

	// GoogleCloudStorageURL is base bucket link for the content images
	GoogleCloudStorageURL = "GOOGLE_CLOUD_STORAGE_URL"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// GetInviteLink generates a custom invite link for PRO or CONSUMER
func GetInviteLink(flavour feedlib.Flavour) (string, error) {
	proLink := serverutils.MustGetEnvVar(ProInviteLink)
	consumerLink := serverutils.MustGetEnvVar(ConsumerInviteLink)

	switch flavour {
	case feedlib.FlavourConsumer:
		return consumerLink, nil
	case feedlib.FlavourPro:
		return proLink, nil
	default:
		return "", fmt.Errorf("failed to get invite link for flavor: %v", flavour)
	}
}

// CreateInviteMessage creates a new invite message
func CreateInviteMessage(user *domain.User, inviteLink string, pin string, flavour feedlib.Flavour) string {
	consumerAppName := serverutils.MustGetEnvVar("CONSUMER_APP_NAME")
	proAppName := serverutils.MustGetEnvVar("PRO_APP_NAME")
	switch flavour {
	case feedlib.FlavourConsumer:
		message := fmt.Sprintf("You have been invited to %s. Download the app on %v. Your username is %v and single use pin is %v",
			consumerAppName, inviteLink, user.Username, pin)
		return message
	case feedlib.FlavourPro:
		message := fmt.Sprintf("You have been invited to %s. Download the app on %v. Your username is %v and single use pin is %v",
			proAppName, inviteLink, user.Username, pin)
		return message
	default:
		return ""
	}
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// EncryptSensitiveData encrypts sensitive data	for a user
func EncryptSensitiveData(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DecryptSensitiveData decrypts sensitive data for a user
func DecryptSensitiveData(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText, _ := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// GetPinExpiryDate returns the expiry date for the given pin
func GetPinExpiryDate() (*time.Time, error) {
	pinExpiryDays := serverutils.MustGetEnvVar("PIN_EXPIRY_DAYS")
	pinExpiryInt, err := strconv.Atoi(pinExpiryDays)
	if err != nil {
		return nil, exceptions.InternalErr(fmt.Errorf("failed to convert PIN expiry days to int: %v", err))
	}
	expiryDate := time.Now().AddDate(0, 0, pinExpiryInt)

	return &expiryDate, nil
}

// RestAPIResponseHelper returns custom standardised response for frontend response consistency
func RestAPIResponseHelper(key string, value interface{}) *dto.RestEndpointResponses {
	response := &dto.RestEndpointResponses{
		Data: map[string]interface{}{
			key: value,
		},
	}
	return response
}

// ReportErrorToSentry captures the exception thrown and registers an issue in sentry
func ReportErrorToSentry(err error) {
	defer sentry.Flush(2 * time.Millisecond)
	sentry.CaptureException(err)
}

// ComposeNotificationPayload creates the notification payload that will sent via pubsub
func ComposeNotificationPayload(userProfile *domain.User, notificationData dto.FCMNotificationMessage) *firebasetools.SendNotificationPayload {
	return &firebasetools.SendNotificationPayload{
		RegistrationTokens: userProfile.PushTokens,
		Notification: &firebasetools.FirebaseSimpleNotificationInput{
			Title: notificationData.Title,
			Body:  notificationData.Body,
		},
	}
}
