package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

var (
	cmsServiceBaseURL = serverutils.MustGetEnvVar("CONTENT_SERVICE_BASE_URL")
)

var (
	clientsPath       = "api/clients/"
	programsPath      = "api/programs/"
	organisationsPath = "api/organisations/"
	facilitiesPath    = "api/facilities/"
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

		registerClientAPIEndpoint := fmt.Sprintf("%s/%s", cmsServiceBaseURL, clientsPath)
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

		deleteClientAPIEndpoint := fmt.Sprintf("%s/%s%s/", cmsServiceBaseURL, clientsPath, data.UserID)
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

		createCMSProgramPath := fmt.Sprintf("%s/%s", cmsServiceBaseURL, programsPath)
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

		createCMSOrganisationPath := fmt.Sprintf("%s/%s", cmsServiceBaseURL, organisationsPath)
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

	case ps.AddPubSubNamespace(common.CreateCMSFacilityTopicName, MyCareHubServiceName):
		var data dto.CreateCMSFacilityPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		facilityPayload := &dto.CreateCMSFacilityPayload{
			FacilityID: data.FacilityID,
			Name:       data.Name,
		}

		facilitiesURL := fmt.Sprintf("%s/%s", cmsServiceBaseURL, facilitiesPath)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPost, facilitiesURL, facilityPayload)
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

	case ps.AddPubSubNamespace(common.CreateCMSProgramFacilityTopicName, MyCareHubServiceName):
		var data dto.CMSLinkFacilityToProgramPayload
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		programFacilityPayload := &dto.CMSLinkFacilityToProgramPayload{
			FacilityID: data.FacilityID,
		}

		addFacilityToProgramPath := fmt.Sprintf("%s/%s%s/", cmsServiceBaseURL, programsPath, data.ProgramID)
		resp, err := ps.BaseExt.MakeRequest(ctx, http.MethodPatch, addFacilityToProgramPath, programFacilityPayload)
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
