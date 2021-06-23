package usecases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
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
	repo       repository.OnboardingRepository
	engagement engagement.ServiceEngagement
	baseExt    extension.BaseExtension
	pin        UserPINUseCases
}

// NewAgentUseCases returns a new a onboarding usecase
func NewAgentUseCases(
	r repository.OnboardingRepository,
	eng engagement.ServiceEngagement,
	ext extension.BaseExtension,
	pin UserPINUseCases,
) AgentUseCase {

	return &AgentUseCaseImpl{
		repo:       r,
		engagement: eng,
		baseExt:    ext,
		pin:        pin,
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

	timestamp := time.Now().In(base.TimeLocation)
	agentProfile := base.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: base.BioData{
			FirstName: &input.FirstName,
			LastName:  &input.LastName,
			Gender:    input.Gender,
		},
		Role:        base.RoleTypeAgent,
		Permissions: base.RoleTypeAgent.Permissions(),
		CreatedByID: &usp.ID,
		Created:     &timestamp,
	}

	// create a user profile in bewell
	profile, err := a.repo.CreateDetailedUserProfile(ctx, input.PhoneNumber, agentProfile)
	if err != nil {
		// wrapped error
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
		// wrapped error
		return nil, err
	}

	otp, err := a.pin.SetUserTempPIN(ctx, profile.ID)
	if err != nil {
		// wrapped error
		return nil, err
	}

	if err := a.notifyNewAgent([]string{input.Email}, []string{input.PhoneNumber}, *profile.UserBioData.FirstName, otp); err != nil {
		return nil, fmt.Errorf("unable to send agent registration notifications: %w", err)
	}

	return profile, nil
}

func (a *AgentUseCaseImpl) notifyNewAgent(emails []string, phoneNumbers []string, name, otp string) error {
	type pin struct {
		Name string
		Pin  string
	}

	t := template.Must(template.New("agentApprovalEmail").Parse(utils.AgentApprovalEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, pin{name, otp})
	if err != nil {
		log.Fatalf("error while generating agent approval email template: %s", err)
	}

	message := fmt.Sprintf("%sPlease use this One Time PIN: %s to log onto Bewell with your phone number", agentWelcomeMessage, otp)
	if err := a.engagement.SendSMS(phoneNumbers, message); err != nil {
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
