package screeningtools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetScreeningToolsQuestion represents the interface to get screening tools questions
type IGetScreeningToolsQuestion interface {
	GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
	GetAvailableScreeningToolQuestions(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error)
	GetAvailableFacilityScreeningTools(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error)
}

// IAnswerScreeningToolQuestion represents the interface to answer screening tools questions
type IAnswerScreeningToolQuestion interface {
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error)
}

// IGetAssessmentResponses is used to get the screening tools assessment responses
type IGetAssessmentResponses interface {
	GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error)
}

// IGetScreeningToolServiceRequestResponses represents the interface to get screening tool responses
type IGetScreeningToolServiceRequestResponses interface {
	GetScreeningToolServiceRequestResponses(ctx context.Context, clientID string, toolType enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error)
}

// UseCasesScreeningTools represents the usecases for screening tools
type UseCasesScreeningTools interface {
	IGetScreeningToolsQuestion
	IAnswerScreeningToolQuestion
	IGetScreeningToolServiceRequestResponses
	IGetAssessmentResponses
}

// ServiceScreeningToolsImpl represents screening tools implementation object
type ServiceScreeningToolsImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
}

// NewUseCasesScreeningTools is the controller for the screening tools usecases
func NewUseCasesScreeningTools(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
) *ServiceScreeningToolsImpl {
	return &ServiceScreeningToolsImpl{
		Query:       query,
		Create:      create,
		Update:      update,
		ExternalExt: externalExt,
	}
}

// GetScreeningToolQuestions get all the screening tools questions
func (t *ServiceScreeningToolsImpl) GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error) {
	emptyString := ""
	if questionType == nil {
		questionType = &emptyString
	} else {
		ok := enums.ScreeningToolType(*questionType).IsValid()
		if !ok {
			return nil, fmt.Errorf("invalid question type: %s in input", *questionType)
		}
	}

	screeningToolsQuestions, err := t.Query.GetScreeningToolQuestions(ctx, *questionType)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get screening tools questions: %v", err))
	}

	return screeningToolsQuestions, nil
}

// AnswerScreeningToolQuestions answer screening tools questions
// a condition is a collection of values collected when answering a set of screening tool questions.
//
//	some of the data we can save in the condition includes:
//		1. response value
//		2. meta data of the question
//		3. the number of times the a given choice has been selected for the given set of questions (based on `ToolType`)
//
// Example:
//
//	given a set of screening tool questions:
//	[
//		{
//			"ID": "fe8f8f8f-8f8f-8f8f-8f8f-8f8f8f8f8f8f",
//			"Question": "In the past, has anyone made you feel threatened, fearful or in danger?",
//			"ToolType": "VIOLENCE_ASSESSMENT",
//			"ResponseChoices": map[string]interface{}{"1": "Yes", "2": "No"},
//			"ResponseType": "INTEGER",
//			"ResponseCategory": "SINGLE_CHOICE",
//			"Active": true,
//			"Sequence": 0
//			"Meta": map[string]interface{}{
//				"helper_text": "Emotional violence Assessment",
//				"violence_type": "EMOTIONAL",
//				"violence_code": "GBV-EV",
//			}
//		}
//	]
//	we can formulate a condition like:
//	we assume the user answered yes for this question
//	{
//		"VIOLENCE_ASSESSMENT_question_number_0": "yes", // response value
//		"VIOLENCE_ASSESSMENT_question_number_0_meta": {
//			"helper_text": "Emotional violence Assessment",
//			"violence_type": "EMOTIONAL",
//			"violence_code": "GBV-EV",
//		},// meta data of question number 0
//		"VIOLENCE_ASSESSMENT_yes_count": 1, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered yes
//		"VIOLENCE_ASSESSMENT_no_count": 0, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered no
//	}
func (t *ServiceScreeningToolsImpl) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
	condition := make(map[string]interface{})
	serviceRequests := make([]*domain.ServiceRequest, 0)
	toolTypeCategory := make(map[string]string)

	if len(screeningToolResponses) == 0 {
		return false, fmt.Errorf("no screening tool responses provided")
	}

	clientProfile, err := t.Query.GetClientProfileByClientID(ctx, screeningToolResponses[0].ClientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get client profile: %v", err))
	}

	for _, screeningToolResponse := range screeningToolResponses {
		err := screeningToolResponse.Validate()
		if err != nil {
			return false, fmt.Errorf("screening tool responses are empty: %v", err)
		}

		screeningToolQuestion, err := t.Query.GetScreeningToolQuestionByQuestionID(ctx, screeningToolResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to get screening tool question: %v", err)
		}

		err = screeningToolQuestion.ValidateResponseQuestionCategory(screeningToolResponse.Response, screeningToolQuestion.ResponseCategory)
		if err != nil {
			return false, fmt.Errorf("invalid response: %v", err)
		}

		err = screeningToolQuestion.ValidateResponseQUestionType(screeningToolResponse.Response, screeningToolQuestion.ResponseType)
		if err != nil {
			return false, fmt.Errorf("invalid response: %v", err)
		}
		err = t.Update.InvalidateScreeningToolResponse(ctx, screeningToolResponse.ClientID, screeningToolResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to invalidate previous screening tool response: %v", err)
		}

		condition = addCondition(screeningToolQuestion, screeningToolResponse.Response, condition)
		serviceRequest := createServiceRequest(screeningToolQuestion, screeningToolResponse.Response, condition)
		if serviceRequest != nil {
			score := strconv.Itoa(utils.InterfaceToInt(condition[screeningToolQuestion.ToolType.String()+"_score"]))
			serviceRequest.Active = true
			serviceRequest.Status = enums.ServiceRequestStatusPending.String()
			serviceRequest.ClientID = screeningToolResponse.ClientID
			serviceRequest.FacilityID = clientProfile.DefaultFacilityID
			serviceRequest.Meta = map[string]interface{}{
				"question_id":         screeningToolQuestion.ID,
				"question_type":       screeningToolQuestion.ToolType,
				"score":               score,
				"screening_tool_name": utils.InterfaceToString(condition[screeningToolQuestion.ToolType.String()+"_screening_tool_name"]),
			}
			// delete(toolTypeCategory, ScreeningToolTypeGBV)
			if _, ok := toolTypeCategory[screeningToolQuestion.ToolType.String()]; ok {
				continue
			}
			toolTypeCategory[string(screeningToolQuestion.ToolType)] = string(screeningToolQuestion.ToolType)
			serviceRequests = append(serviceRequests, serviceRequest)
		}
	}

	err = t.Create.AnswerScreeningToolQuestions(ctx, screeningToolResponses)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to answer screening tools questions: %v", err)
	}

	for s := range serviceRequests {
		serviceRequestInput := &dto.ServiceRequestInput{
			Active:       serviceRequests[s].Active,
			RequestType:  serviceRequests[s].RequestType,
			Status:       serviceRequests[s].Status,
			Request:      serviceRequests[s].Request,
			ClientID:     serviceRequests[s].ClientID,
			InProgressBy: serviceRequests[s].InProgressBy,
			ResolvedBy:   serviceRequests[s].ResolvedBy,
			FacilityID:   serviceRequests[s].FacilityID,
			ClientName:   serviceRequests[s].ClientName,
			Flavour:      feedlib.FlavourConsumer,
			Meta:         serviceRequests[s].Meta,
		}
		err = t.Create.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create service request: %v", err)
		}
	}
	return true, nil
}

// GetAvailableScreeningToolQuestions returns all screening tool questions that fit the following criteria:
// 1. A screening tool response for each client should be after 24 hours of the last response
// 2. A screening tool response that creates a service request should be resolved first before the next response
// 3. A user who is MALE should not answer a contraceptives question
func (t *ServiceScreeningToolsImpl) GetAvailableScreeningToolQuestions(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error) {
	availableScreeningTools := []*domain.AvailableScreeningTools{}

	validToolTypes := make(map[enums.ScreeningToolType]*domain.AvailableScreeningTools)

	_, err := t.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ProfileNotFoundErr(fmt.Errorf("failed to get client profile: %v", err))
	}

	for i := range enums.ScreeningToolQuestions {
		validToolTypes[enums.ScreeningToolQuestions[i]] = &domain.AvailableScreeningTools{
			ToolType: enums.ScreeningToolQuestions[i],
		}
	}

	loggedInUserID, err := t.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}
	userProfile, err := t.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ProfileNotFoundErr(fmt.Errorf("failed to get user profile: %v", err))
	}
	if userProfile.Gender == "MALE" {
		delete(validToolTypes, enums.ScreeningToolTypeCUI)

	}

	clientProfile, err := t.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ProfileNotFoundErr(fmt.Errorf("failed to get client profile: %v", err))
	}

	activeScreeningResponses, err := t.Query.GetActiveScreeningToolResponses(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get active screening tool responses: %v", err)
	}
	for _, activeScreeningResponse := range activeScreeningResponses {
		screeningtool, err := t.Query.GetScreeningToolQuestionByQuestionID(ctx, activeScreeningResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get screening tool question: %v", err)
		}
		delete(validToolTypes, screeningtool.ToolType)
	}

	pendingServiceRequests, err := t.Query.GetClientServiceRequests(
		ctx,
		enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
		enums.ServiceRequestStatusPending.String(),
		clientID,
		clientProfile.DefaultFacilityID,
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get pending service requests: %v", err)
	}
	for _, pendingServiceRequest := range pendingServiceRequests {
		delete(validToolTypes, enums.ScreeningToolType(utils.InterfaceToString(pendingServiceRequest.Meta["question_type"])))
	}

	for _, v := range validToolTypes {
		// backward compatibility for old screening tool service requests
		// the new service requests for screening tools do not save question_type
		if !v.ToolType.IsValid() {
			continue
		}
		availableScreeningTools = append(availableScreeningTools, v)
	}

	return availableScreeningTools, nil
}

// GetAvailableFacilityScreeningTools returns all screening tool questions that fit the following criteria:
// 1. Each tool type returned must have a response by the client and it's service request status must be pending
// 2. The client and the staff must belong to the same facility
func (t *ServiceScreeningToolsImpl) GetAvailableFacilityScreeningTools(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error) {
	availableScreeningTools := []*domain.AvailableScreeningTools{}
	validToolTypes := make(map[enums.ScreeningToolType]*domain.AvailableScreeningTools)

	toolType := enums.ServiceRequestTypeScreeningToolsRedFlag.String()
	status := enums.ServiceRequestStatusPending.String()
	pendingServiceRequests, err := t.Query.GetServiceRequests(
		ctx,
		&toolType,
		&status,
		facilityID,
		feedlib.FlavourConsumer,
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get pending service requests: %v", err)
	}
	for _, pendingServiceRequest := range pendingServiceRequests {
		validToolTypes[enums.ScreeningToolType(utils.InterfaceToString(pendingServiceRequest.Meta["question_type"]))] = &domain.AvailableScreeningTools{
			ToolType: enums.ScreeningToolType(utils.InterfaceToString(pendingServiceRequest.Meta["question_type"])),
		}
	}

	for _, v := range validToolTypes {
		// backward compatibility for old screening tool service requests
		// the new service requests for screening tools do not save question_type
		if !v.ToolType.IsValid() {
			continue
		}
		availableScreeningTools = append(availableScreeningTools, v)
	}
	return availableScreeningTools, nil
}

// GetAssessmentResponses returns the assessment responses for a given facility
func (t *ServiceScreeningToolsImpl) GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
	if facilityID == "" || !enums.ScreeningToolType(toolType).IsValid() {
		return nil, fmt.Errorf("invalid facility id or tool type")
	}
	return t.Query.GetAssessmentResponses(ctx, facilityID, toolType)
}

// GetScreeningToolServiceRequestResponses returns all screening tool responses for a client who has a service request of the specified tool type in param
func (t *ServiceScreeningToolsImpl) GetScreeningToolServiceRequestResponses(ctx context.Context, clientID string, toolType enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error) {
	screeningToolResponses := []*domain.ScreeningToolResponse{}
	ScreeningToolResponsePayload := &domain.ScreeningToolResponsePayload{}
	if clientID == "" {
		return nil, fmt.Errorf("client id is required")
	}
	ok := toolType.IsValid()
	if !ok {
		err := fmt.Errorf("invalid screening tool type")
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	client, err := t.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ProfileNotFoundErr(fmt.Errorf("failed to get client profile: %v", err))
	}
	clientResponses, err := t.Query.GetClientScreeningToolResponsesByToolType(ctx, clientID, toolType.String(), true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client screening tool responses: %v", err)
	}

	serviceRequest, err := t.Query.GetClientScreeningToolServiceRequestByToolType(ctx, clientID, string(toolType), enums.ServiceRequestStatusPending.String())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client screening tool service request: %v", err)
	}

	for _, clientResponse := range clientResponses {
		screeningToolQuestion, err := t.Query.GetScreeningToolQuestionByQuestionID(ctx, clientResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get screening tool question: %v", err)
		}

		screeningToolResponses = append(screeningToolResponses, &domain.ScreeningToolResponse{
			ToolIndex: screeningToolQuestion.Sequence,
			Tool:      screeningToolQuestion.Question,
			Response:  utils.InterfaceToString(screeningToolQuestion.ResponseChoices[clientResponse.Answer]),
		})
	}
	ScreeningToolResponsePayload.ScreeningToolResponses = screeningToolResponses
	ScreeningToolResponsePayload.ServiceRequestID = serviceRequest.ID
	ScreeningToolResponsePayload.ClientContact = client.User.Contacts.ContactValue
	return ScreeningToolResponsePayload, nil

}
