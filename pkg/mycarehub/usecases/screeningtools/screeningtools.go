package screeningtools

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetScreeningToolsQuestion represents the interface to get screening tools questions
type IGetScreeningToolsQuestion interface {
	GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
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
	Query  infrastructure.Query
	Create infrastructure.Create
	Update infrastructure.Update
}

// NewUseCasesScreeningTools is the controller for the screening tools usecases
func NewUseCasesScreeningTools(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
) *ServiceScreeningToolsImpl {
	return &ServiceScreeningToolsImpl{
		Query:  query,
		Create: create,
		Update: update,
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
			serviceRequest.ClientID = screeningToolResponse.ClientID
			serviceRequest.FacilityID = clientProfile.FacilityID
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
