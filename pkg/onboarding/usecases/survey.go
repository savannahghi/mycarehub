package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// SurveyUseCases represents all the business logic involved in user post visit surveys.
type SurveyUseCases interface {
	RecordPostVisitSurvey(ctx context.Context, input dto.PostVisitSurveyInput) (bool, error)
}

// SurveyUseCasesImpl represents the usecase implementation object
type SurveyUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	baseExt              extension.BaseExtension
}

// NewSurveyUseCases initializes a new sign up usecase
func NewSurveyUseCases(
	r repository.OnboardingRepository,
	ext extension.BaseExtension,
) *SurveyUseCasesImpl {
	return &SurveyUseCasesImpl{r, ext}
}

// RecordPostVisitSurvey records the survey input supplied by the user
func (rs *SurveyUseCasesImpl) RecordPostVisitSurvey(
	ctx context.Context,
	input dto.PostVisitSurveyInput,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "RecordPostVisitSurvey")
	defer span.End()

	if input.LikelyToRecommend < 0 || input.LikelyToRecommend > 10 {
		return false, exceptions.LikelyToRecommendError(
			fmt.Errorf(exceptions.LikelyToRecommendErrMsg),
		)
	}

	UID, err := rs.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.UserNotFoundError(err)
	}

	if err := rs.onboardingRepository.RecordPostVisitSurvey(ctx, input, UID); err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(fmt.Errorf(exceptions.InternalServerErrorMsg))
	}

	return true, nil
}
