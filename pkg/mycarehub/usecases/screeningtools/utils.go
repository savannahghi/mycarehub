package screeningtools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func validateResponses(ctx context.Context, questions []*domain.ScreeningToolQuestion, answersInput []*dto.ScreeningToolQuestionResponseInput) error {
	questionMap := make(map[string]*domain.ScreeningToolQuestion)

	for _, question := range questions {
		questionMap[question.ID] = question
	}

	for _, v := range answersInput {
		if err := questionMap[v.QuestionID].ValidateResponseQUestionType(v.Response, v.ResponseType); err != nil {
			return err
		}
		if err := questionMap[v.QuestionID].ValidateResponseQuestionCategory(v.Response, v.ResponseCategory); err != nil {
			return err
		}
	}

	return nil
}

func generateServiceRequest(ctx context.Context, clientProfile *domain.ClientProfile, answersInput []*dto.ScreeningToolQuestionResponseInput) (*dto.ServiceRequestInput, error) {
	var (
		score   int
		request string
	)
	toolType := answersInput[0].ToolType

	serviceRequest := &dto.ServiceRequestInput{
		Active:      true,
		RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
		Status:      enums.ServiceRequestStatusPending.String(),
		Request:     request,
		ClientID:    *clientProfile.ID,
		FacilityID:  clientProfile.FacilityID,
		ClientName:  &clientProfile.User.Name,
		Flavour:     feedlib.FlavourConsumer,
		Meta: map[string]interface{}{
			"question_id":   answersInput[0].QuestionID,
			"question_type": toolType,
			"score":         score,
		},
	}

	switch toolType {
	case enums.ScreeningToolTypeTB:
		score, err := calculateToolScore(ctx, answersInput)
		if err != nil {
			return nil, err
		}
		if score >= 3 {
			serviceRequest.Request = fmt.Sprintf("%s has a score of %v in  the TB Assessment Tool. Consider reaching out to them", clientProfile.User.Name, score)
			serviceRequest.Meta["score"] = score
			return serviceRequest, nil
		}
	case enums.ScreeningToolTypeGBV:
		score, err := calculateToolScore(ctx, answersInput)
		if err != nil {
			return nil, err
		}
		if score >= 1 {
			serviceRequest.Request = fmt.Sprintf("%s has a score of %v in  the Gender Based Violence Assessment Tool. Consider reaching out to them", clientProfile.User.Name, score)
			serviceRequest.Meta["score"] = score
			return serviceRequest, nil
		}
	case enums.ScreeningToolTypeAlcoholSubstanceAssessment:
		score, err := calculateToolScore(ctx, answersInput)
		if err != nil {
			return nil, err
		}
		if score >= 3 {
			serviceRequest.Request = fmt.Sprintf("%s has a score of %v in  the Alcohol and Substance Assessment Tool. Consider reaching out to them", clientProfile.User.Name, score)
			serviceRequest.Meta["score"] = score
			return serviceRequest, nil
		}
	case enums.ScreeningToolTypeCUI:
		questionSequence := answersInput[0].QuestionSequence
		// if yes to question 4, trigger a service request to schedule a visit with the healthcare worker
		if questionSequence == 3 && strings.ToLower(answersInput[0].Response) == "yes" {
			serviceRequest.Request = fmt.Sprintf("%s has answered yes to question 4 in the Contraceptive Use Assessment Tool. Consider reaching out to them", clientProfile.User.Name)
			return serviceRequest, nil
		}
	default:
		return nil, fmt.Errorf("unsupported tool type: %s", toolType)
	}
	return nil, nil
}

func calculateToolScore(ctx context.Context, answersInput []*dto.ScreeningToolQuestionResponseInput) (int, error) {
	var score int
	for _, v := range answersInput {

		s, err := strconv.Atoi(v.Response)
		if err != nil {
			return 0, err
		}
		score += s
	}
	return score, nil
}
