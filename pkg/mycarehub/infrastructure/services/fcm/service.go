package fcm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// ServiceFCMImpl provides methods for sending fcm notifications
type ServiceFCMImpl struct {
	fcmClient *messaging.Client
}

func initializeFCMClient() (*messaging.Client, error) {
	ctx := context.Background()
	fc := &firebasetools.FirebaseClient{}
	app, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app: %s", err)
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Messaging client: %v", err)
	}

	return fcmClient, nil
}

// ServiceFCM defines all interactions with the FCM service
type ServiceFCM interface {
	SendNotification(
		ctx context.Context,
		payload *firebasetools.SendNotificationPayload,
	) (bool, error)
}

// NewService initializes a service to interact with Firebase Cloud Messaging
func NewService() ServiceFCM {
	fcmClient, err := initializeFCMClient()
	if err != nil {
		log.Panicf("error getting Messaging client: %v\n", err)
	}

	return &ServiceFCMImpl{
		fcmClient: fcmClient,
	}
}

// SendNotification sends a data message to the specified registration tokens.
//
// It returns:
//
//  - a list of registration tokens for which message sending failed
//  - an error, if no message sending occurred
//
// Notification messages can also be accompanied by custom `data`.
//
// For data messages, the following keys should be avoided:
//
//  - reserved words: "from", "notification" and "message_type"
//  - any word starting with "gcm" or "google"
//
// Messages that are time sensitive (e.g video calls) should be sent with
// `HIGH_PRIORITY`. Their time to live should also be limited (or the expiry)
// set on iOS. For Android, there is a `TTL` key in `messaging.AndroidConfig`.
// For iOS, the `apns-expiration` header should be set to a specific timestamp
// e.g `"apns-expiration":"1604750400"`. For web, there's a `TTL` header that
// is also a number of seconds e.g. `"TTL":"4500"`.
//
// For Android, priority is set via the `messaging.AndroidConfig` `priority`
// key to either "normal" or "high". It should be set to "high" only for urgent
// notification e.g video call notifications. For web, it is set via the
// `Urgency` header e.g "Urgency": "high". For iOS, the "apns-priority" header
// is used, with "5" for normal/low and "10" to mean urgent/high.
//
// The callers of this method should implement retries and exponential backoff,
// if necessary.
func (s ServiceFCMImpl) SendNotification(
	ctx context.Context,
	payload *firebasetools.SendNotificationPayload,
) (bool, error) {
	if payload.RegistrationTokens == nil {
		return false, fmt.Errorf("can't send FCM notifications to nil registration tokens")
	}
	message := &messaging.MulticastMessage{Tokens: payload.RegistrationTokens}

	if payload.Data != nil {
		err := ValidateFCMData(payload.Data)
		if err != nil {
			return false, err
		}
		message.Data = payload.Data
	}
	if payload.Notification != nil {
		message.Notification = &messaging.Notification{
			Title: payload.Notification.Title,
			Body:  payload.Notification.Body,
		}
		if payload.Notification.ImageURL != nil {
			message.Notification.ImageURL = *payload.Notification.ImageURL
		}
	}
	if payload.Android != nil {
		message.Android = &messaging.AndroidConfig{
			Priority: payload.Android.Priority,
			Data:     converterandformatter.ConvertInterfaceMap(payload.Android.Data),
		}
		if payload.Android.CollapseKey != nil {
			message.Android.CollapseKey = *payload.Android.CollapseKey
		}
		if payload.Android.RestrictedPackageName != nil {
			message.Android.RestrictedPackageName = *payload.Android.RestrictedPackageName
		}
	}
	if payload.Web != nil {
		message.Webpush = &messaging.WebpushConfig{
			Headers: converterandformatter.ConvertInterfaceMap(payload.Web.Headers),
			Data:    converterandformatter.ConvertInterfaceMap(payload.Web.Data),
		}
	}
	if payload.Ios != nil {
		message.APNS = &messaging.APNSConfig{
			Headers: converterandformatter.ConvertInterfaceMap(payload.Web.Headers),
		}
	}

	batchResp, err := s.fcmClient.SendMulticast(ctx, message)
	if err != nil {
		return false, fmt.Errorf("unable to send FCM messages: %w", err)
	}

	var errorMessages []string
	for idx, resp := range batchResp.Responses {
		if !resp.Success {
			// The order of responses corresponds to the order of the registration tokens.
			msg := fmt.Sprintf(
				"fcm: failed to send message to %s: %v",
				payload.RegistrationTokens[idx],
				resp.Error,
			)
			errorMessages = append(errorMessages, msg)
		}
		if payload.Notification != nil {
			savedNotification := dto.SavedNotification{
				ID:                uuid.New().String(),
				RegistrationToken: payload.RegistrationTokens[idx],
				MessageID:         resp.MessageID,
				Timestamp:         time.Now(),
			}
			if payload.Notification != nil {
				savedNotification.Notification = &firebasetools.FirebaseSimpleNotificationInput{
					Title:    payload.Notification.Title,
					Body:     payload.Notification.Body,
					ImageURL: payload.Notification.ImageURL,
				}
			}
			if payload.Data != nil {
				savedNotification.Data = converterandformatter.ConvertStringMap(payload.Data)
			}
			if payload.Android != nil {
				savedNotification.AndroidConfig = &firebasetools.FirebaseAndroidConfigInput{
					CollapseKey: payload.Android.CollapseKey,
					Priority:    payload.Android.Priority,
					Data:        payload.Android.Data,
				}
			}
			if payload.Web != nil {
				savedNotification.WebpushConfig = &firebasetools.FirebaseWebpushConfigInput{
					Headers: payload.Web.Headers,
					Data:    payload.Web.Data,
				}
			}
			if payload.Ios != nil {
				savedNotification.APNSConfig = &firebasetools.FirebaseAPNSConfigInput{
					Headers: payload.Ios.Headers,
				}
			}
			// TODO - Save the notification
		}
	}
	if len(errorMessages) > 0 {
		return false, fmt.Errorf(strings.Join(errorMessages, "; "))
	}
	return true, nil
}
