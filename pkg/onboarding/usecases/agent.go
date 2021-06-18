package usecases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	agentWelcomeMessage      = "You have been successfully registered as an agent. We look forward to working with you."
	agentWelcomeEmailSubject = "Successfully registered as an agent"
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

// RegisterAgent creates a new Agent in bewell
func (a *AgentUseCaseImpl) RegisterAgent(ctx context.Context, input dto.RegisterAgentInput) (*base.UserProfile, error) {

	_, err := a.baseExt.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

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

	agentProfile := base.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: base.BioData{
			FirstName: &input.FirstName,
			LastName:  &input.LastName,
			Gender:    input.Gender,
		},
		Role: base.RoleTypeAgent,
	}

	// create a user profile in bewell
	profile, err := a.repo.CreateDetailedUserProfile(ctx, input.PhoneNumber, agentProfile)
	if err != nil {
		return nil, err
	}

	_, err = a.repo.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	sup := base.Supplier{
		IsOrganizationVerified: true,
		SladeCode:              SavannahSladeCode,
		KYCSubmitted:           true,
		PartnerSetupComplete:   true,
		OrganizationName:       SavannahOrgName,
	}

	_, err = a.repo.CreateDetailedSupplierProfile(ctx, profile.ID, sup)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	// set the user default communications settings
	defaultCommunicationSetting := true
	_, err = a.repo.SetUserCommunicationsSettings(
		ctx,
		profile.ID,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
		&defaultCommunicationSetting,
	)
	if err != nil {
		return nil, err
	}

	if err := a.notifyNewAgent([]string{input.Email}, []string{input.PhoneNumber}); err != nil {
		return nil, fmt.Errorf("unable to send agent registration notifications: %w", err)
	}

	return profile, nil
}

func (a *AgentUseCaseImpl) notifyNewAgent(emails []string, phoneNumbers []string) error {
	t := template.Must(template.New("agentApprovalEmail").Parse(utils.AgentApprovalEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("error while generating agent approval email template: %s", err)
	}

	if err := a.engagement.SendSMS(phoneNumbers, agentWelcomeMessage); err != nil {
		return fmt.Errorf("unable to send agent registration message: %w", err)
	}

	text := buf.String()
	for _, email := range emails {
		if err := a.engagement.SendMail(email, text, agentWelcomeEmailSubject); err != nil {
			return fmt.Errorf("unable to send agent registration email: %w", err)
		}
	}

	return nil
}
