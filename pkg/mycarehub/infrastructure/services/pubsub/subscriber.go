package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

// ReceivePubSubPushMessages receives and processes a pubsub message
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	message, err := ps.baseExt.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	topicID, err := pubsubtools.GetPubSubTopic(message)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	switch topicID {
	case ps.AddPubSubNamespace(common.CreateGetstreamEventTopicName, MyCareHubServiceName):
		var data dto.GetStreamEvent
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		switch data.Type {
		case stream.EventMessageNew:
			channel, err := ps.GetStream.GetChannel(ctx, data.ChannelID)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

			var channelMetadata domain.CommunityMetadata
			err = mapstructure.Decode(channel.ExtraData, &channelMetadata)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

			for _, member := range data.Members {
				var metadata domain.MemberMetadata
				err := mapstructure.Decode(member.User.ExtraData, &metadata)
				if err != nil {
					serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
						Err:     err,
						Message: err.Error(),
					}, http.StatusBadRequest)
					return
				}

				if metadata.UserType == "STAFF" {
					staffProfile, err := ps.Query.GetStaffProfileByStaffID(ctx, member.User.ID)
					if err != nil {
						helpers.ReportErrorToSentry(err)
					}

					notificationData := &dto.FCMNotificationMessage{
						Title: channelMetadata.Name,
						Body:  fmt.Sprintf("%v: %v", data.User.Name, data.Message.Text),
					}

					payload := helpers.ComposeNotificationPayload(staffProfile.User, *notificationData)
					_, err = ps.FCM.SendNotification(ctx, payload)
					if err != nil {
						helpers.ReportErrorToSentry(err)
						log.Printf("failed to send notification: %v", err)
					}
				}
			}
		}
	}

	resp := map[string]string{"Status": "Success"}
	returnedResponse, err := json.Marshal(resp)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	_, _ = w.Write(returnedResponse)
}
