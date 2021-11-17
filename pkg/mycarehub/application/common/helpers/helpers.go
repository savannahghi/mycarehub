package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

const (
	// ProInviteLink will store the pro app invite link at the settings
	ProInviteLink = "PRO_INVITE_LINK"

	// ConsumerInviteLink will store the consumer invite link at the settings
	ConsumerInviteLink = "CONSUMER_INVITE_LINK"
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
func CreateInviteMessage(user *domain.User, inviteLink string, pin string) string {
	message := fmt.Sprintf("Dear %v %v, you have been invited to My Afya Hub. Download the app on %v. Your single use pin is %v",
		user.FirstName, user.LastName, inviteLink, pin)
	return message
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

func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// DecryptSensitiveData decrypts sensitive data for a user
func DecryptSensitiveData(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
