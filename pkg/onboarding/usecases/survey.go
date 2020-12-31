package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// SurveyUseCases represents all the business logic involved in user post visit surveys.
type SurveyUseCases interface {
	RecordPostVisitSurvey(ctx context.Context, input domain.PostVisitSurveyInput) (bool, error)
}

// SurveyUseCasesImpl represents the usecase implementation object
type SurveyUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewSurveyUseCases initializes a new sign up usecase
func NewSurveyUseCases(r repository.OnboardingRepository) *SurveyUseCasesImpl {
	return &SurveyUseCasesImpl{r}
}

// RecordPostVisitSurvey records the survey input supplied by the user
func (rs *SurveyUseCasesImpl) RecordPostVisitSurvey(
	ctx context.Context,
	input domain.PostVisitSurveyInput,
) (bool, error) {
	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return false, &domain.CustomError{
			Err:     nil,
			Message: exceptions.LikelyToRecommendErrMsg,
			Code:    0, // TODO: Add a code for this error
		}
	}

	UID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, err
	}

	return rs.onboardingRepository.RecordPostVisitSurvey(ctx, input, UID)
}
