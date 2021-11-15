package helpers

import (
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

// CreateResetPinMessage creates reset pin message
func CreateResetPinMessage(user *domain.User, pin string) string {
	message := fmt.Sprintf("Dear %v %v, your PIN for My Afya Hub has been reset successfully. Your single use pin is %v",
		user.FirstName, user.LastName, pin)
	return message
}
