package screeningtools

import (
	"strconv"
	"strings"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func addCondition(question *domain.ScreeningToolQuestion, response string, condition map[string]interface{}) map[string]interface{} {
	var (
		responseValue string
	)

	// saves the response value for a given question in the condition
	responseKey := question.ToolType.String() + "_question_number_" + strconv.Itoa(question.Sequence)

	condition[responseKey] = response

	// saves extra meta data for a given question
	condition[responseKey+"_question_meta"] = question.Meta

	switch question.ResponseCategory {
	case enums.ScreeningToolResponseCategorySingleChoice:
		responseValue = question.ResponseChoices[response].(string)
		responseTrimmedValue := strings.TrimSpace(strings.ToLower(responseValue))
		// saves the count of responses for a given question
		// this is for  `ResponseCategory` `SINGLE_CHOICE`
		// a user is expected to select only one response from `ResponseChoices`
		// for each set of question answered for a given `ToolType`
		// the count is incremented for each response give,
		// eg. we record how many yes responses were given for a given question
		// of tool type `VIOLENCE_ASSESSMENT`
		//
		// the key will be `VIOLENCE_ASSESSMENT_yes_count`
		responseCountKey := question.ToolType.String() + "_" + responseTrimmedValue + "_count"
		condition[responseCountKey] = utils.InterfaceToInt(condition[responseCountKey]) + 1
	}
	/*
				The condition would resemble this:
				{
						"VIOLENCE_ASSESSMENT_question_number_0": "yes", // response value
						"VIOLENCE_ASSESSMENT_question_number_0_question_meta": {
							"helper_text": "Emotional violence Assessment",
							"violence_type": "EMOTIONAL",
							"violence_code": "GBV-EV",
						},// meta data of question number 0
						"VIOLENCE_ASSESSMENT_yes_count": 1, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered yes
		             	"VIOLENCE_ASSESSMENT_no_count": 0, // number of times question of tool type VIOLENCE_ASSESSMENT has been answered no
					}
	*/
	return condition
}

func createServiceRequest(question *domain.ScreeningToolQuestion, response string, condition map[string]interface{}) *domain.ServiceRequest {
	serviceRequestTemplate := serviceRequestTemplate(question, response, condition)

	yesCount := utils.InterfaceToInt(condition[question.ToolType.String()+"_"+"yes"+"_count"])
	noCount := utils.InterfaceToInt(condition[question.ToolType.String()+"_"+"no"+"_count"])

	switch question.ToolType {
	case enums.ScreeningToolTypeTB:
		if yesCount >= 3 {
			condition[question.ToolType.String()+"_score"] = yesCount
			condition[question.ToolType.String()+"_screening_tool_name"] = "TB"
			return &domain.ServiceRequest{
				RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
				Request:     serviceRequestTemplate,
			}
		}
		if noCount == 4 {
			// 	//TODO:  TPT and repeat screening on subsequent visits
			return nil
		}
	case enums.ScreeningToolTypeGBV:
		if yesCount >= 1 {
			condition[question.ToolType.String()+"_score"] = yesCount
			condition[question.ToolType.String()+"_screening_tool_name"] = "GBV"
			return &domain.ServiceRequest{
				RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
				Request:     serviceRequestTemplate,
			}
		}
	case enums.ScreeningToolTypeCUI:
		toolTypeResponse := enums.ScreeningToolTypeCUI.String() + "_question_number_" + "3"
		// TODO:  send a notification to HCW in the client’s facility for discussion with the clinician during the next visit
		if condition[toolTypeResponse] == "yes" {
			return nil
		}
	case enums.ScreeningToolTypeAlcoholSubstanceAssessment:
		if yesCount >= 3 {
			condition[question.ToolType.String()+"_score"] = yesCount
			condition[question.ToolType.String()+"_screening_tool_name"] = "Alcohol and Substance"
			return &domain.ServiceRequest{
				RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
				Request:     serviceRequestTemplate,
			}
		}
	}
	return nil
}

func serviceRequestTemplate(question *domain.ScreeningToolQuestion, response string, condition map[string]interface{}) string {
	var (
		requestTemplate string
		gbvMetaTemplate string
	)

	questionMeta := getQuestionMeta(question, condition)

	callClientString := "Consider calling Client for further discussion."
	repeatScreeningString := "Repeat screening for client on subsequent visits."

	yesCount := utils.InterfaceToInt(condition[question.ToolType.String()+"_"+"yes"+"_count"])
	noCount := utils.InterfaceToInt(condition[question.ToolType.String()+"_"+"no"+"_count"])

	switch question.ToolType {
	case enums.ScreeningToolTypeTB:
		if yesCount >= 3 {
			requestTemplate = "TB assessment: greater than or equal to 3 yes responses indicates positive TB cases. " + callClientString
		}
		if noCount == 4 {
			requestTemplate = "TB assessment: all no responses indicates no TB cases. " + repeatScreeningString
		}
	case enums.ScreeningToolTypeGBV:
		gbvMetaTemplate = "The GBV code is " + utils.InterfaceToString(questionMeta["violence_code"]) + ". "
		if yesCount >= 1 {
			requestTemplate = "Violence assessment: greater than or equal to 1 yes responses indicates positive Violence cases. " + gbvMetaTemplate + callClientString
		}
		if noCount == 4 {
			requestTemplate = "Violence assessment: all no responses indicates no GBV cases. " + repeatScreeningString
		}
	case enums.ScreeningToolTypeCUI:
		// the sequence is zero in
		toolTypeResponse := enums.ScreeningToolTypeCUI.String() + "_question_number_" + "3"
		if condition[toolTypeResponse] == "yes" {
			requestTemplate = "Contraceptive assessment: yes response to question number 4. " + callClientString
		}
	case enums.ScreeningToolTypeAlcoholSubstanceAssessment:
		if yesCount >= 3 {
			requestTemplate = "Alcohol/Substance Assessment: greater than or equal to 3 yes responses indicates positive alcohol/substance cases. " + callClientString
		}
		if noCount == 4 {
			requestTemplate = "Alcohol/Substance Assessment: all no responses indicates no alcohol/substance cases. " + repeatScreeningString
		}
	}
	return requestTemplate
}

func getQuestionMeta(question *domain.ScreeningToolQuestion, condition map[string]interface{}) map[string]interface{} {
	questionMetaKey := question.ToolType.String() + "_question_number_" + strconv.Itoa(question.Sequence) + "_question_meta"

	if condition[questionMetaKey] != nil {
		return condition[questionMetaKey].(map[string]interface{})
	}
	return nil
}
