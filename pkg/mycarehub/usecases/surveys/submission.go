package surveys

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// getFormResponse composes the response to a survey using the form and submission data
func getFormResponse(ctx context.Context, form, submissionData map[string]interface{}) ([]*domain.SurveyResponse, error) {
	responses := []*domain.SurveyResponse{}

	formHTML, ok := form["html"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid form: expected a 'html' key")
	}

	formBody, ok := formHTML["body"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid form: expected a 'body' key")
	}

	// single choice question responses
	singleChoiceQuestions, ok := formBody["select1"].([]interface{})
	if ok {
		response, err := getSingleChoiceResponses(ctx, singleChoiceQuestions, submissionData)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		responses = append(responses, response...)
	}

	return responses, nil
}

// getSingleChoiceResponses composes the single choice responses for a form
func getSingleChoiceResponses(ctx context.Context, questions []interface{}, submissionData map[string]interface{}) ([]*domain.SurveyResponse, error) {
	responses := []*domain.SurveyResponse{}

	// holds the choice/selection of an individual in their submission
	// {question_id:choice}
	submissions := make(map[string]string)
	for key, value := range submissionData {
		v, ok := value.(string)
		if ok {
			submissionData[key] = v
		}
	}

	for _, node := range questions {
		questionNode, ok := node.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid single choice question: %v", node)
		}

		reference, ok := questionNode["-ref"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid question: expected a question 'ref' key for %v", node)
		}

		refSplit := strings.Split(reference, "/")
		questionID := refSplit[len(refSplit)-1]

		questionText, ok := questionNode["label"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid question: expected a question 'label' key for %v", node)
		}

		// holds the available choices for a question
		// {value:label}
		// value is the representation of a choice as stored in the db eg 1, 2
		// label is the human readable representation of a choice e.g yes, no
		choices := make(map[string]string)
		for _, item := range questionNode["item"].([]interface{}) {
			i := item.(map[string]interface{})

			value, ok := i["value"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid single choice: expected a 'value' key for %v", item)
			}

			label, ok := i["label"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid single choice: expected a 'label' key for %v", item)
			}

			choices[value] = label
		}

		response := &domain.SurveyResponse{
			Question:     questionText,
			QuestionType: "SINGLE_CHOICE",
			Answer:       choices[submissions[questionID]],
		}

		responses = append(responses, response)
	}

	return responses, nil
}
