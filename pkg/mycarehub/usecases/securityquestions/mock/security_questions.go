package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SecurityQuestionsUseCaseMock mocks the implementation of security question usecase methods.
type SecurityQuestionsUseCaseMock struct {
	MockGetSecurityQuestionsFn         func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
	MockGetSecurityQuestionByIDFn      func(ctx context.Context, id string, flavour feedlib.Flavour) (*domain.SecurityQuestion, error)
	MockSaveSecurityQuestionResponseFn func(ctx context.Context, securityQuestionResponse *dto.SecurityQuestionResponseInput) error
}

// NewSecurityQuestionsUseCaseMock creates and itializes security question mocks
func NewSecurityQuestionsUseCaseMock() *SecurityQuestionsUseCaseMock {
	return &SecurityQuestionsUseCaseMock{
		MockGetSecurityQuestionsFn: func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
			securityQuestion := &domain.SecurityQuestion{
				QuestionStem: "test",
				Description:  "test",
				Flavour:      feedlib.FlavourConsumer,
				Active:       true,
				ResponseType: enums.SecurityQuestionResponseTypeNumber,
			}
			return []*domain.SecurityQuestion{securityQuestion}, nil
		},
		MockGetSecurityQuestionByIDFn: func(ctx context.Context, id string, flavour feedlib.Flavour) (*domain.SecurityQuestion, error) {
			securityQuestion := &domain.SecurityQuestion{
				QuestionStem: "test",
				Description:  "test",
				Flavour:      feedlib.FlavourConsumer,
				Active:       true,
				ResponseType: enums.SecurityQuestionResponseTypeNumber,
			}
			return securityQuestion, nil
		},
		MockSaveSecurityQuestionResponseFn: func(ctx context.Context, securityQuestionResponse *dto.SecurityQuestionResponseInput) error {
			return nil
		},
	}
}

//GetSecurityQuestions mocks the implementation of getting all the security questions.
func (sq *SecurityQuestionsUseCaseMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return sq.MockGetSecurityQuestionsFn(ctx, flavour)
}

// GetSecurityQuestionByID mocks the implementation of getting a security question by ID.
func (sq *SecurityQuestionsUseCaseMock) GetSecurityQuestionByID(ctx context.Context, id string, flavour feedlib.Flavour) (*domain.SecurityQuestion, error) {
	return sq.MockGetSecurityQuestionByIDFn(ctx, id, flavour)
}

// SaveSecurityQuestionResponse saves the security question response.
func (sq *SecurityQuestionsUseCaseMock) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse *dto.SecurityQuestionResponseInput) error {
	return sq.MockSaveSecurityQuestionResponseFn(ctx, securityQuestionResponse)
}
