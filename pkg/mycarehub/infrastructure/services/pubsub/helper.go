package pubsubmessaging

import (
	"context"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

// ProcessGetStreamEvent is used to process the getstream events that are sent to our webhook endpoints.
// When the data is published to the `getstream.events` topic, a subscriber will receive the messages and
// process it. The intention is to send a FCM message to a staff user.
func (ps ServicePubSubMessaging) ProcessGetStreamEvent(
	ctx context.Context,
	w http.ResponseWriter,
	data *dto.GetStreamEvent,
	notification *dto.FCMNotificationMessage,
) error {
	channel, err := ps.GetStream.GetChannel(ctx, data.ChannelID)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}

	var channelMetadata domain.CommunityMetadata
	err = mapstructure.Decode(channel.ExtraData, &channelMetadata)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}

	for _, member := range data.Members {
		var metadata domain.MemberMetadata
		err := mapstructure.Decode(member.User.ExtraData, &metadata)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return err
		}

		if metadata.UserType == enums.StaffUser.String() {
			staffProfile, err := ps.Query.GetStaffProfileByStaffID(ctx, member.User.ID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
			}

			notification.Title = channelMetadata.Name

			if staffProfile != nil {
				payload := helpers.ComposeNotificationPayload(staffProfile.User, *notification)
				_, err = ps.FCM.SendNotification(ctx, payload)
				if err != nil {
					helpers.ReportErrorToSentry(err)
					log.Printf("failed to send notification: %v", err)
				}
			}
		}
	}
	return nil
}
