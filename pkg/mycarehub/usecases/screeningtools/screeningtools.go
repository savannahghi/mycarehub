package screeningtools

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetScreeningToolsQuestion represents the interface to get screening tools questions
type IGetScreeningToolsQuestion interface {
	GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
	GetAvailableScreeningToolQuestions(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error)
}

// IAnswerScreeningToolQuestion represents the interface to answer screening tools questions
type IAnswerScreeningToolQuestion interface {
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error)
}

// UseCasesScreeningTools represents the usecases for screening tools
type UseCasesScreeningTools interface {
	IGetScreeningToolsQuestion
	IAnswerScreeningToolQuestion
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
// 		some of the data we can save in the condition includes:
// 			1. response value
// 			2. meta data of the question
// 			3. the number of times the a given choice has been selected for the given set of questions (based on `ToolType`)
// Example:
// 		given a set of screening tool questions:
// 		[
// 			{
// 				"ID": "fe8f8f8f-8f8f-8f8f-8f8f-8f8f8f8f8f8f",
// 				"Question": "In the past, has anyone made you feel threatened, fearful or in danger?",
// 				"ToolType": "VIOLENCE_ASSESSMENT",
// 				"ResponseChoices": map[string]interface{}{"1": "Yes", "2": "No"},
// 				"ResponseType": "INTEGER",
// 				"ResponseCategory": "SINGLE_CHOICE",
// 				"Active": true,
// 				"Sequence": 0
// 				"Meta": map[string]interface{}{
// 					"helper_text": "Emotional violence Assessment",
// 					"violence_type": "EMOTIONAL",
// 					"violence_code": "GBV-EV",
// 				}
// 			}
// 		]
// 		we can formulate a condition like:
// 		we assume the user answered yes for this question
// 		{
// 			"VIOLENCE_ASSESSMENT_question_number_0": "yes", // response value
// 			"VIOLENCE_ASSESSMENT_question_number_0_meta": {
// 				"helper_text": "Emotional violence Assessment",
// 				"violence_type": "EMOTIONAL",
// 				"violence_code": "GBV-EV",
// 			},// meta data of question number 0
// 			"VIOLENCE_ASSESSMENT_yes_count": 1, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered yes
// 			"VIOLENCE_ASSESSMENT_no_count": 0, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered no
// 		}
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
			serviceRequest.Active = true
			serviceRequest.Status = enums.ServiceRequestStatusPending.String()
			serviceRequest.ClientID = screeningToolResponse.ClientID
			serviceRequest.FacilityID = clientProfile.FacilityID
			serviceRequest.Meta = map[string]interface{}{
				"question_id":   screeningToolQuestion.ID,
				"question_type": screeningToolQuestion.ToolType,
			}
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
		enums.ServiceRequestTypeScreeningTools.String(),
		enums.ServiceRequestStatusPending.String(),
		clientID,
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get pending service requests: %v", err)
	}
	for _, pendingServiceRequest := range pendingServiceRequests {
		delete(validToolTypes, enums.ScreeningToolType(interfaceToString(pendingServiceRequest.Meta["question_type"])))
	}

	for _, v := range validToolTypes {
		availableScreeningTools = append(availableScreeningTools, v)
	}

	return availableScreeningTools, nil
}
