package securityquestions

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetSecurityQuestions gets the security questions
type IGetSecurityQuestions interface {
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error)
}

// UseCaseSecurityQuestion groups all the security questions method interfaces
type UseCaseSecurityQuestion interface {
	IGetSecurityQuestions
}

// UseCaseSecurityQuestionsImpl represents security question implementation object
type UseCaseSecurityQuestionsImpl struct {
	Query infrastructure.Query
}

// NewSecurityQuestionsUsecase returns a new security question instance
func NewSecurityQuestionsUsecase(query infrastructure.Query) *UseCaseSecurityQuestionsImpl {
	return &UseCaseSecurityQuestionsImpl{
		Query: query,
	}
}

// GetSecurityQuestions gets all the security questions
func (s *UseCaseSecurityQuestionsImpl) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	return s.Query.GetSecurityQuestions(ctx, flavour)
}
