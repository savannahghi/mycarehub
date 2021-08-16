package usecases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
)

const (
	agentWelcomeMessage      = "We look forward to working with you."
	agentWelcomeEmailSubject = "Successfully registered as an agent"
)

// AgentUseCase represent the business logic required for management of agents
type AgentUseCase interface {
	RegisterAgent(
		ctx context.Context,
		input dto.RegisterAgentInput,
	) (*profileutils.UserProfile, error)
	ActivateAgent(ctx context.Context, input dto.ProfileSuspensionInput) (bool, error)
	DeactivateAgent(ctx context.Context, input dto.ProfileSuspensionInput) (bool, error)
	FetchAgents(ctx context.Context) ([]*dto.Agent, error)
	FindAgentbyPhone(ctx context.Context, phoneNumber *string) (*dto.Agent, error)
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

func (a *AgentUseCaseImpl) checkPreconditions() {
	if a.repo == nil {
		log.Panicf("nil repository in agent usecase implementation")
	}

	if a.engagement == nil {
		log.Panicf("nil engagement service in agent usecase implementation")
	}

	if a.baseExt == nil {
		log.Panicf("nil base extension in agent usecase implementation")
	}

	if a.pin == nil {
		log.Panicf("nil pin usecase in agent usecase implementation")
	}
}

// RegisterAgent creates a new Agent in bewell
func (a *AgentUseCaseImpl) RegisterAgent(
	ctx context.Context,
	input dto.RegisterAgentInput,
) (*profileutils.UserProfile, error) {
	a.checkPreconditions()
	ctx, span := tracer.Start(ctx, "RegisterAgent")
	defer span.End()

	msisdn, err := a.baseExt.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	// Check logged in user has permissions/role of employee
	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Get Logged In user profile
	usp, err := a.repo.GetUserProfileByUID(ctx, p.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)

	agentProfile := profileutils.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: profileutils.BioData{
			FirstName:   &input.FirstName,
			LastName:    &input.LastName,
			Gender:      input.Gender,
			DateOfBirth: &input.DateOfBirth,
		},
		Role:        profileutils.RoleTypeAgent,
		Permissions: profileutils.RoleTypeAgent.Permissions(),
		Roles:       input.RoleIDs,
		CreatedByID: &usp.ID,
		Created:     &timestamp,
	}

	// create a user profile in bewell
	profile, err := a.repo.CreateDetailedUserProfile(ctx, *msisdn, agentProfile)
	if err != nil {
		// wrapped error
		utils.RecordSpanError(span, err)
		return nil, err
	}

	_, err = a.repo.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.InternalServerError(err)
	}

	sup := profileutils.Supplier{
		IsOrganizationVerified: true,
		SladeCode:              SavannahSladeCode,
		KYCSubmitted:           true,
		PartnerSetupComplete:   true,
		OrganizationName:       SavannahOrgName,
	}

	_, err = a.repo.CreateDetailedSupplierProfile(ctx, profile.ID, sup)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return nil, err
	}

	otp, err := a.pin.SetUserTempPIN(ctx, profile.ID)
	if err != nil {
		// wrapped error
		utils.RecordSpanError(span, err)
		return nil, err
	}

	if err := a.notifyNewAgent(ctx, input.Email, input.PhoneNumber, *profile.UserBioData.FirstName, otp); err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to send agent registration notifications: %w", err)
	}

	return profile, nil
}

func (a *AgentUseCaseImpl) notifyNewAgent(
	ctx context.Context,
	email, phoneNumber, firstName, tempPIN string,
) error {
	type pin struct {
		Name string
		Pin  string
	}

	message := fmt.Sprintf(domain.WelcomeMessage, firstName, tempPIN)
	message += " " + agentWelcomeMessage

	if err := a.engagement.SendSMS(ctx, []string{phoneNumber}, message); err != nil {
		return fmt.Errorf("unable to send agent registration message: %w", err)
	}

	if email != "" {
		t := template.Must(template.New("agentApprovalEmail").Parse(utils.AgentApprovalEmail))

		buf := new(bytes.Buffer)

		err := t.Execute(buf, pin{firstName, tempPIN})
		if err != nil {
			log.Fatalf("error while generating agent approval email template: %s", err)
		}

		text := buf.String()

		if err := a.engagement.SendMail(ctx, email, text, agentWelcomeEmailSubject); err != nil {
			return fmt.Errorf("unable to send agent registration email: %w", err)
		}

	}

	return nil
}

// ActivateAgent activates/unsuspend the agent profile
func (a *AgentUseCaseImpl) ActivateAgent(
	ctx context.Context,
	input dto.ProfileSuspensionInput,
) (bool, error) {
	a.checkPreconditions()
	ctx, span := tracer.Start(ctx, "ActivateAgent")
	defer span.End()

	agent, err := a.repo.GetUserProfileByID(ctx, input.ID, true)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	err = a.repo.UpdateSuspended(ctx, agent.ID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// DeactivateAgent deactivates/suspends the agent profile
func (a *AgentUseCaseImpl) DeactivateAgent(
	ctx context.Context,
	input dto.ProfileSuspensionInput,
) (bool, error) {
	a.checkPreconditions()
	ctx, span := tracer.Start(ctx, "DeactivateAgent")
	defer span.End()

	// Get agent profile using phoneNumber
	agent, err := a.repo.GetUserProfileByID(ctx, input.ID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	if agent.Role != profileutils.RoleTypeAgent {
		return false, exceptions.InternalServerError(fmt.Errorf("this user is not an agent"))
	}

	err = a.repo.UpdateSuspended(ctx, agent.ID, true)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// FetchAgents fetches registered agents
func (a *AgentUseCaseImpl) FetchAgents(ctx context.Context) ([]*dto.Agent, error) {
	a.checkPreconditions()
	ctx, span := tracer.Start(ctx, "FetchAgents")
	defer span.End()

	profiles, err := a.repo.ListUserProfiles(ctx, profileutils.RoleTypeAgent)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	agents := []*dto.Agent{}

	for _, profile := range profiles {
		// Retrieve the agent PIN
		pin, err := a.repo.GetPINByProfileID(ctx, profile.ID)
		if err != nil {
			utils.RecordSpanError(span, err)
			// the error is wrapped already. No need to wrap it again
			return nil, err
		}

		agent := &dto.Agent{
			ID:                      profile.ID,
			PhotoUploadID:           profile.PhotoUploadID,
			UserBioData:             profile.UserBioData,
			PrimaryPhone:            *profile.PrimaryPhone,
			PrimaryEmailAddress:     profile.PrimaryEmailAddress,
			SecondaryPhoneNumbers:   profile.SecondaryPhoneNumbers,
			SecondaryEmailAddresses: profile.SecondaryEmailAddresses,
			TermsAccepted:           profile.TermsAccepted,
			Suspended:               profile.Suspended,
			ResendPIN:               pin.IsOTP,
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

// FindAgentbyPhone is used to find an agent using their phone number
func (a *AgentUseCaseImpl) FindAgentbyPhone(
	ctx context.Context,
	phoneNumber *string,
) (*dto.Agent, error) {
	a.checkPreconditions()
	ctx, span := tracer.Start(ctx, "FindAgentbyPhone")
	defer span.End()

	phoneNumber, err := a.baseExt.NormalizeMSISDN(*phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	profile, err := a.repo.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.AgentNotFoundError(err)
	}

	agent := dto.Agent{
		ID:                      profile.ID,
		PhotoUploadID:           profile.PhotoUploadID,
		UserBioData:             profile.UserBioData,
		PrimaryPhone:            *profile.PrimaryPhone,
		PrimaryEmailAddress:     profile.PrimaryEmailAddress,
		SecondaryPhoneNumbers:   profile.SecondaryPhoneNumbers,
		SecondaryEmailAddresses: profile.SecondaryEmailAddresses,
		TermsAccepted:           profile.TermsAccepted,
		Suspended:               profile.Suspended,
	}

	return &agent, nil
}
