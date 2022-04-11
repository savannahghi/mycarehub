package fcm

import (
	"fmt"
	"strings"
)

// ValidateFCMData checks that the supplied FCM data does not use re
func ValidateFCMData(data map[string]string) error {
	if data != nil {
		fcmReservedWords := []string{"from", "notification", "message_type"}
		for _, reservedWord := range fcmReservedWords {
			_, present := data[reservedWord]
			if present {
				return fmt.Errorf("invalid use of FCM reserved word: %s", reservedWord)
			}
		}
		fcmReservedPrefixes := []string{"gcm", "google"}
		for _, reservedPrefix := range fcmReservedPrefixes {
			for k := range data {
				if strings.HasPrefix(k, reservedPrefix) {
					return fmt.Errorf("illegal FCM prefix: %s", reservedPrefix)
				}
			}
		}
	}
	return nil
}
