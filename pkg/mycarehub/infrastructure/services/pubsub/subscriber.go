package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (
	flaggedMessageNotificationBody = "A message from %v has been flagged"
	addMemberNotificationBody      = "A new member has been added to the community"
	removeMemberNotificationBody   = "%v has been removed from the community"
	bannedUserNotificationBody     = "%v has been banned from the community"
	unbanUserNotificationBody      = "%v has been unbanned and can rejoin the community"
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
			notificationData := &dto.FCMNotificationMessage{
				Body: fmt.Sprintf("%v: %v", data.User.Name, data.Message.Text),
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

		case "message.flagged":
			notificationData := &dto.FCMNotificationMessage{
				Body: fmt.Sprintf(flaggedMessageNotificationBody, data.User.Name),
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

		case stream.EventMemberAdded:
			notificationData := &dto.FCMNotificationMessage{
				Body: addMemberNotificationBody,
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

		case stream.EventMemberRemoved:
			notificationData := &dto.FCMNotificationMessage{
				Body: fmt.Sprintf(removeMemberNotificationBody, data.User.Name),
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

		case "user.banned":
			notificationData := &dto.FCMNotificationMessage{
				Body: fmt.Sprintf(bannedUserNotificationBody, data.User.Name),
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}

		case "user.unbanned":
			notificationData := &dto.FCMNotificationMessage{
				Body: fmt.Sprintf(unbanUserNotificationBody, data.User.Name),
			}
			err := ps.ProcessGetStreamEvent(ctx, w, &data, notificationData)
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
		}

	case ps.AddPubSubNamespace(common.CreateCMSClientTopicName, MyCareHubServiceName):
		var data dto.CMSClientOutput
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		clientInput := &dto.PubSubCMSClientInput{
			UserID:         data.UserID,
			Name:           data.Name,
			Gender:         data.Gender,
			UserType:       data.UserType,
			PhoneNumber:    data.PhoneNumber,
			Handle:         data.Handle,
			Flavour:        data.Flavour,
			DateOfBirth:    data.DateOfBirth,
			ClientID:       data.ClientID,
			ClientTypes:    data.ClientTypes,
			EnrollmentDate: data.EnrollmentDate,
			FacilityID:     data.FacilityID,
			FacilityName:   data.FacilityName,
			OrganisationID: data.OrganisationID,
		}

		err = ps.CMS.RegisterClient(ctx, clientInput)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
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
