package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// AgentUseCase represent the business logic required for management of agents
type AgentUseCase interface {
	RegisterAgent(ctx context.Context, input dto.RegisterAgentInput) (*base.UserProfile, error)
}

// AgentUseCaseImpl  represents usecase implementation object
type AgentUseCaseImpl struct {
	repo    repository.OnboardingRepository
	profile ProfileUseCase

	engagement engagement.ServiceEngagement
	messaging  messaging.ServiceMessaging
	baseExt    extension.BaseExtension
}

// NewAgentUseCases returns a new a onboarding usecase
func NewAgentUseCases(
	r repository.OnboardingRepository,
	p ProfileUseCase,
	eng engagement.ServiceEngagement,
	messaging messaging.ServiceMessaging,
	ext extension.BaseExtension,
) AgentUseCase {

	return &AgentUseCaseImpl{
		repo:       r,
		profile:    p,
		engagement: eng,
		messaging:  messaging,
		baseExt:    ext,
	}
}

// RegisterAgent ...
func (a *AgentUseCaseImpl) RegisterAgent(ctx context.Context, input dto.RegisterAgentInput) (*base.UserProfile, error) {
	// Check logged in user has permissions/role of employee
	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	// Get Logged In user profile
	usp, err := a.repo.GetUserProfileByUID(ctx, p.UID, false)
	if err != nil {
		return nil, err
	}

	if usp.Role != base.RoleTypeEmployee {
		return nil, exceptions.RoleNotValid(fmt.Errorf("error: logged in user does not have `EMPLOYEE` role"))
	}

	/*


		Create user profile and Assign agent role to user profile
		save the user profile to db --> create agent

		Create agent specific supplier profile

		Notify agent: email + text message

	*/
	return nil, nil
}
