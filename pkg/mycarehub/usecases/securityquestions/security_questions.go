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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/serverutils"
)

// SensitiveContentPassphrase is the secret key used when encrypting and decrypting a security question response
var SensitiveContentPassphrase = serverutils.MustGetEnvVar("SENSITIVE_CONTENT_SECRET_KEY")

// IGetSecurityQuestions creates the security questions
type ICreateSecurityQuestions interface {
	CreateSecurityQuestions(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error)
}

// IGetSecurityQuestions gets the security questions
type IGetSecurityQuestions interface {
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
}

// IGetUserRespondedSecurityQuestions gets the security questions that a user had set earlier during onboarding
type IGetUserRespondedSecurityQuestions interface {
	GetUserRespondedSecurityQuestions(ctx context.Context, input dto.GetUserRespondedSecurityQuestionsInput) ([]*domain.SecurityQuestion, error)
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
		responses *dto.VerifySecurityQuestionsPayload,
	) (bool, error)
}

// UseCaseSecurityQuestion groups all the security questions method interfaces
type UseCaseSecurityQuestion interface {
	IGetSecurityQuestions
	IRecordSecurityQuestionResponses
	IVerifySecurityQuestionResponses
	IGetUserRespondedSecurityQuestions
	ICreateSecurityQuestions
}

// UseCaseSecurityQuestionsImpl represents security question implementation object
type UseCaseSecurityQuestionsImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
	ExternalExt extension.ExternalMethodsExtension
}

// NewSecurityQuestionsUsecase returns a new security question instance
func NewSecurityQuestionsUsecase(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
) *UseCaseSecurityQuestionsImpl {
	return &UseCaseSecurityQuestionsImpl{
		Query:       query,
		Create:      create,
		Update:      update,
		ExternalExt: externalExt,
	}
}

// GetSecurityQuestions gets all the security questions
func (s *UseCaseSecurityQuestionsImpl) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return s.Query.GetSecurityQuestions(ctx, flavour)
}

// RecordSecurityQuestionResponses records the security question responses during user onboarding
func (s *UseCaseSecurityQuestionsImpl) RecordSecurityQuestionResponses(ctx context.Context, input []*dto.SecurityQuestionResponseInput) ([]*domain.RecordSecurityQuestionResponse, error) {
	var recordSecurityQuestionResponses []*domain.RecordSecurityQuestionResponse
	var securityQuestionResponseInput []*dto.SecurityQuestionResponseInput

	for _, i := range input {
		err := i.Validate()
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.InputValidationErr(fmt.Errorf("security question response validation failed: %s", err))
		}

		securityQuestion, err := s.Query.GetSecurityQuestionByID(ctx, &i.SecurityQuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.ItemNotFoundErr(fmt.Errorf("security question id %s does not exist", i.SecurityQuestionID))
		}

		err = securityQuestion.Validate(i.Response)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.InputValidationErr(fmt.Errorf("security question response %s is invalid: %v", i.Response, err))
		}

		encryptedResponse, err := helpers.EncryptSensitiveData(i.Response, SensitiveContentPassphrase)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.EncryptionErr(fmt.Errorf("failed to encrypt sensitive data response: %v", err))
		}

		securityQuestionResponsePayload := &dto.SecurityQuestionResponseInput{
			UserID:             i.UserID,
			SecurityQuestionID: i.SecurityQuestionID,
			Response:           encryptedResponse,
		}

		securityQuestionResponseInput = append(securityQuestionResponseInput, securityQuestionResponsePayload)

		recordSecurityQuestionResponses = append(recordSecurityQuestionResponses,
			&domain.RecordSecurityQuestionResponse{
				SecurityQuestionID: i.SecurityQuestionID,
				IsCorrect:          true,
			})
	}

	// save the response
	err := s.Create.SaveSecurityQuestionResponse(ctx, securityQuestionResponseInput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.FailedToSaveItemErr(fmt.Errorf("failed to save security question response data %v", err))
	}

	return recordSecurityQuestionResponses, nil
}

// VerifySecurityQuestionResponses verifies the security questions against the recorded responses.
func (s *UseCaseSecurityQuestionsImpl) VerifySecurityQuestionResponses(
	ctx context.Context,
	responses *dto.VerifySecurityQuestionsPayload,
) (bool, error) {
	failCountInstance := make(map[string]int)
	if len(responses.SecurityQuestionsInput) == 0 {
		helpers.ReportErrorToSentry(fmt.Errorf("no responses provided"))
		return false, exceptions.EmptyInputErr(fmt.Errorf("no responses provided"))
	}
	userProfile, err := s.Query.GetUserProfileByUsername(ctx, responses.SecurityQuestionsInput[0].Username)
	if err != nil {
		return false, exceptions.ProfileNotFoundErr(err)
	}

	failCount := userProfile.FailedSecurityCount
	if failCount >= 3 {
		err := fmt.Errorf("failed: security questions answering attempts exceeded 3 attempts")
		helpers.ReportErrorToSentry(err)
		return false, exceptions.FailedSecurityCountExceededErr(err)
	}
	for _, securityQuestionResponse := range responses.SecurityQuestionsInput {
		questionResponse, err := s.Query.GetSecurityQuestionResponse(ctx, securityQuestionResponse.QuestionID, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.SecurityQuestionNotFoundErr(fmt.Errorf("security question does not exist"))
		}

		decryptedResponse, err := helpers.DecryptSensitiveData(questionResponse.Response, SensitiveContentPassphrase)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.InternalErr(fmt.Errorf("failed to decrypt sensitive data response: %v", err))
		}

		if !strings.EqualFold(securityQuestionResponse.Response, decryptedResponse) {

			_, err := s.Update.UpdateIsCorrectSecurityQuestionResponse(ctx, *userProfile.ID, false)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, exceptions.InternalErr(fmt.Errorf("failed to update security question response: %v", err))
			}

			failCount++
			failCountInstance[securityQuestionResponse.Username] = failCount

			err = s.Update.UpdateFailedSecurityQuestionsAnsweringAttempts(ctx, *userProfile.ID, failCount)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, exceptions.InternalErr(fmt.Errorf("failed to update security question response fail count %v", err))
			}

			helpers.ReportErrorToSentry(err)
			return false, exceptions.SecurityQuestionResponseMismatchErr(fmt.Errorf("the security question response does not match: %d attempts left", 3-failCount))
		}
	}

	if failCountInstance[responses.SecurityQuestionsInput[0].Username] <= 3 {
		err := s.Update.UpdateFailedSecurityQuestionsAnsweringAttempts(ctx, *userProfile.ID, 0)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.InternalErr(fmt.Errorf("failed to reset security question response fail count %v", err))
		}
	}
	return true, nil
}

// GetUserRespondedSecurityQuestions gets the security questions that the user had responded to during onboarding
// 3 random question will be drawn when the user is resetting their pin
func (s *UseCaseSecurityQuestionsImpl) GetUserRespondedSecurityQuestions(ctx context.Context, input dto.GetUserRespondedSecurityQuestionsInput) ([]*domain.SecurityQuestion, error) {
	// ensure the phone/flavour is verified
	if err := input.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.EmptyInputErr(fmt.Errorf("empty value passed in input: %v", err))
	}

	if !input.Flavour.IsValid() {
		return nil, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := s.Query.GetUserProfileByUsername(ctx, input.Username)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(fmt.Errorf("failed to get a user profile by phonenumber: %v", err))
	}

	// ensure the otp for the phone is valid
	ok, err := s.Query.VerifyOTP(ctx, &dto.VerifyOTPInput{
		PhoneNumber: userProfile.Contacts.ContactValue,
		OTP:         input.OTP,
		Flavour:     input.Flavour,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to verify otp: %v", err))
	}

	if !ok {
		return nil, exceptions.InternalErr(fmt.Errorf("failed to verify otp: %v", err))
	}

	// ensure the questions are associated with the user who set the responses
	securityQuestionResponses, err := s.Query.GetUserSecurityQuestionsResponses(ctx, *userProfile.ID, input.Flavour.String())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get security questions: %v", err))
	}

	if len(securityQuestionResponses) < 3 {
		return nil, fmt.Errorf("failed to get security questions, user must have answered at least 3")
	}

	shuffledQuestions, err := utils.ShuffleSecurityQuestionResponses(securityQuestionResponses)
	if err != nil {
		return nil, err
	}

	randomTwoSecurityQuestionresponses := shuffledQuestions[:2]
	securityQuestions := []*domain.SecurityQuestion{}

	// return random 2 security questions
	for _, i := range randomTwoSecurityQuestionresponses {
		securityQuestion, err := s.Query.GetSecurityQuestionByID(ctx, &i.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get security question: %v", err))
		}
		securityQuestions = append(securityQuestions, securityQuestion)
	}

	return securityQuestions, nil
}

// CreateSecurityQuestions enables the creation of security questions in bulk. It is mainly used by the CLI
func (s *UseCaseSecurityQuestionsImpl) CreateSecurityQuestions(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error) {
	return s.Create.CreateSecurityQuestions(ctx, securityQuestions)
}
