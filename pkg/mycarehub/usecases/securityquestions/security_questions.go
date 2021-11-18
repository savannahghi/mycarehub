package securityquestions

import (
	"context"
	"fmt"

	"strings"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/serverutils"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// SensitiveContentPassphrase is the secret key used when encrypting and decrypting a security question response
var SensitiveContentPassphrase = serverutils.MustGetEnvVar("SENSITIVE_CONTENT_SECRET_KEY")

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

// IVerifySecurityQuestionResponses verifies the security questions
type IVerifySecurityQuestionResponses interface {
	VerifySecurityQuestionResponses(
		ctx context.Context,
		responses *[]dto.VerifySecurityQuestionInput,
	) (bool, error)
}

// UseCaseSecurityQuestion groups all the security questions method interfaces
type UseCaseSecurityQuestion interface {
	IGetSecurityQuestions
	IRecordSecurityQuestionResponses
	IVerifySecurityQuestionResponses
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

	for _, i := range input {
		err := i.Validate()
		if err != nil {
			return nil, exceptions.InputValidationErr(fmt.Errorf("security question response validation failed: %s", err))
		}

		securityQuestion, err := s.Query.GetSecurityQuestionByID(ctx, &i.SecurityQuestionID)
		if err != nil {
			return nil, exceptions.ItemNotFoundErr(fmt.Errorf("security question id %s does not exist", i.SecurityQuestionID))
		}

		err = securityQuestion.Validate(i.Response)
		if err != nil {
			return nil, exceptions.InputValidationErr(fmt.Errorf("security question response %s is invalid: %v", i.Response, err))
		}

		encryptedResponse, err := helpers.EncryptSensitiveData(i.Response, SensitiveContentPassphrase)
		if err != nil {
			return nil, exceptions.EncryptionErr(fmt.Errorf("failed to encrypt sensitive data response: %v", err))
		}

		securityQuestionResponsePayload := &dto.SecurityQuestionResponseInput{
			UserID:             i.UserID,
			SecurityQuestionID: i.SecurityQuestionID,
			Response:           encryptedResponse,
		}
		// save the response
		err = s.Create.SaveSecurityQuestionResponse(ctx, securityQuestionResponsePayload)
		if err != nil {
			return nil, exceptions.FailedToSaveItemErr(fmt.Errorf("failed to save security question response data %v", err))
		}

		recordSecurityQuestionResponses = append(recordSecurityQuestionResponses,
			&domain.RecordSecurityQuestionResponse{
				SecurityQuestionID: i.SecurityQuestionID,
				IsCorrect:          true,
			})

	}

	return recordSecurityQuestionResponses, nil
}

// VerifySecurityQuestionResponses verifies the security questions against the recorded responses.
func (s *UseCaseSecurityQuestionsImpl) VerifySecurityQuestionResponses(
	ctx context.Context,
	responses *[]dto.VerifySecurityQuestionInput,
) (bool, error) {
	for _, securityQuestionResponse := range *responses {
		questionResponse, err := s.Query.GetSecurityQuestionResponseByID(ctx, securityQuestionResponse.QuestionID)
		if err != nil {
			return false, fmt.Errorf("failed to fetch security question response")
		}

		decryptedResponse, err := helpers.DecryptSensitiveData(questionResponse.Response, SensitiveContentPassphrase)
		if err != nil {
			return false, fmt.Errorf("failed to decrypt the response")
		}

		if !strings.EqualFold(securityQuestionResponse.Response, decryptedResponse) {
			return false, fmt.Errorf("the security question response does not match")
		}
	}
	return true, nil
}
