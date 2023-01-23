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

var (
	cmsServiceBaseURL      = serverutils.MustGetEnvVar("CONTENT_SERVICE_BASE_URL")
	removeClientPath       = "client_remove"
	removeStaffPath        = "staff_remove"
	registerStaffPath      = "staff_registration"
	registerClientPath     = "/api/clients/"
	createProgramPath      = "api/programs/"
	createOrganisationPath = "api/organisations/"
)

// ReceivePubSubPushMessages receives and processes a pubsub message
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	message, err := ps.BaseExt.VerifyPubSubJWTAndDecodePayload(w, r)
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
		var data dto.PubsubCreateCMSClientPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		clientInput := &dto.PubsubCreateCMSClientPayload{
			ClientID:       data.ClientID,
			Name:           data.Name,
			Gender:         data.Gender,
			DateOfBirth:    data.DateOfBirth,
			OrganisationID: data.OrganisationID,
			ProgramID:      data.ProgramID,
		}

		registerClientAPIEndpoint := fmt.Sprintf("%s/%s", cmsServiceBaseURL, registerClientPath)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPost, registerClientAPIEndpoint, clientInput)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		if resp.StatusCode != http.StatusCreated {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.CreateCMSStaffTopicName, MyCareHubServiceName):
		var data dto.PubsubCreateCMSStaffPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
		}

		staffInput := &dto.PubsubCreateCMSStaffPayload{
			UserID:         data.UserID,
			Name:           data.Name,
			Gender:         data.Gender,
			UserType:       data.UserType,
			PhoneNumber:    data.PhoneNumber,
			Handle:         data.Handle,
			Flavour:        data.Flavour,
			DateOfBirth:    data.DateOfBirth,
			StaffNumber:    data.StaffNumber,
			StaffID:        data.StaffID,
			FacilityID:     data.FacilityID,
			FacilityName:   data.FacilityName,
			OrganisationID: data.OrganisationID,
		}

		registerStaffAPIEndpoint := fmt.Sprintf("%s/%s", cmsServiceBaseURL, registerStaffPath)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPost, registerStaffAPIEndpoint, staffInput)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.DeleteCMSClientTopicName, MyCareHubServiceName):
		var data dto.DeleteCMSUserPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		deleteClientAPIEndpoint := fmt.Sprintf("%s/%s/%s", cmsServiceBaseURL, removeClientPath, data.UserID)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodDelete, deleteClientAPIEndpoint, nil)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.DeleteCMSStaffTopicName, MyCareHubServiceName):
		var data dto.DeleteCMSUserPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		deleteStaffAPIEndpoint := fmt.Sprintf("%s/%s/%s", cmsServiceBaseURL, removeStaffPath, data.UserID)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodDelete, deleteStaffAPIEndpoint, nil)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.CreateCMSProgramTopicName, MyCareHubServiceName):
		var data dto.CreateCMSProgramPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		programPayload := &dto.CreateCMSProgramPayload{
			ProgramID:      data.ProgramID,
			Name:           data.Name,
			OrganisationID: data.OrganisationID,
		}

		createCMSProgramPath := fmt.Sprintf("%s/%s", cmsServiceBaseURL, createProgramPath)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPost, createCMSProgramPath, programPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.CreateCMSOrganisationTopicName, MyCareHubServiceName):
		var data dto.CreateCMSOrganisationPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		organisationPayload := &dto.CreateCMSOrganisationPayload{
			OrganisationID: data.OrganisationID,
			Name:           data.Name,
			Email:          data.Email,
			PhoneNumber:    data.PhoneNumber,
			Code:           data.Code,
		}

		createCMSOrganisationPath := fmt.Sprintf("%s/%s", cmsServiceBaseURL, createOrganisationPath)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPost, createCMSOrganisationPath, organisationPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			err := fmt.Errorf("invalid status code :%v", resp.StatusCode)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	default:
		err := fmt.Errorf("unknown topic ID: %v", topicID)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
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
