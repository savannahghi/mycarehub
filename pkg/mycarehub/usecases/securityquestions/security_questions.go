package securityquestions

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetSecurityQuestions gets the security questions
type IGetSecurityQuestions interface {
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
}

// IRecordSecurityQuestionResponses contains method to  record security question responses
type IRecordSecurityQuestionResponses interface {

	// TODO: Validate the responses...all fields in the struct are required
	// get userID from responses
	// infer flavour and question from responses
	// TODO Save responses for the user...for use in future verification
	// TODO Wire in metrics
	RecordSecurityQuestionResponses(ctx context.Context, input []*dto.SecurityQuestionResponseInput) ([]*domain.RecordSecurityQuestionResponse, error)
}

// UseCaseSecurityQuestion groups all the security questions method interfaces
type UseCaseSecurityQuestion interface {
	IGetSecurityQuestions
	IRecordSecurityQuestionResponses
}

// UseCaseSecurityQuestionsImpl represents security question implementation object
type UseCaseSecurityQuestionsImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	ExternalExt extension.ExternalMethodsExtension
}

// NewSecurityQuestionsUsecase returns a new security question instance
func NewSecurityQuestionsUsecase(
	query infrastructure.Query,
	create infrastructure.Create,
	externalExt extension.ExternalMethodsExtension,
) *UseCaseSecurityQuestionsImpl {
	return &UseCaseSecurityQuestionsImpl{
		Query:       query,
		Create:      create,
		ExternalExt: externalExt,
	}
}

// GetSecurityQuestions gets all the security questions
func (s *UseCaseSecurityQuestionsImpl) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return s.Query.GetSecurityQuestions(ctx, flavour)
}

// RecordSecurityQuestionResponses records the security question responses
func (s *UseCaseSecurityQuestionsImpl) RecordSecurityQuestionResponses(ctx context.Context, input []*dto.SecurityQuestionResponseInput) ([]*domain.RecordSecurityQuestionResponse, error) {
	var recordSecurityQuestionResponses []*domain.RecordSecurityQuestionResponse

	var sensitiveContentPassphrase = "the operating system for health."
	for _, i := range input {
		err := i.Validate()
		if err != nil {
			return nil, fmt.Errorf("security question response validation failed: %s", err.Error())
		}

		securityQuestion, err := s.Query.GetSecurityQuestionByID(ctx, &i.SecurityQuestionID)
		if err != nil {
			return nil, fmt.Errorf("security question id %s does not exist", i.SecurityQuestionID)
		}

		err = securityQuestion.Validate(i.Response)
		if err != nil {
			return nil, fmt.Errorf("response %s is invalid: %v", i.Response, err)
		}

		encryptedResponse, err := helpers.EncryptSensitiveData(i.Response, sensitiveContentPassphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt response: %v", err)
		}

		securityQuestionResponsePayload := &dto.SecurityQuestionResponseInput{
			UserID:             i.UserID,
			SecurityQuestionID: i.SecurityQuestionID,
			Response:           encryptedResponse,
		}
		// save the response
		err = s.Create.SaveSecurityQuestionResponse(ctx, securityQuestionResponsePayload)
		if err != nil {
			return nil, fmt.Errorf("failed to save security question response data")
		}

		recordSecurityQuestionResponses = append(recordSecurityQuestionResponses,
			&domain.RecordSecurityQuestionResponse{
				SecurityQuestionID: i.SecurityQuestionID,
				IsCorrect:          true,
			})

	}

	return recordSecurityQuestionResponses, nil
}
