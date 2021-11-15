package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SecurityQuestionsUseCaseMock mocks the implementation of security question usecase methods.
type SecurityQuestionsUseCaseMock struct {
	MockGetSecurityQuestionsFn func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
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
				ResponseType: enums.NumberResponse,
			}
			return []*domain.SecurityQuestion{securityQuestion}, nil
		},
	}
}

//GetSecurityQuestions mocks the implementation of getting all the security questions.
func (sq *SecurityQuestionsUseCaseMock) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return sq.MockGetSecurityQuestionsFn(ctx, flavour)
}
