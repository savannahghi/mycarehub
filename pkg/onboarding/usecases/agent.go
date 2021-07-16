package usecases

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/savannahghi/pubsubtools"
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
	ActivateAgent(ctx context.Context, agentID string) (bool, error)
	DeactivateAgent(ctx context.Context, agentID string) (bool, error)
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

// RegisterAgent creates a new Agent in bewell
func (a *AgentUseCaseImpl) RegisterAgent(
	ctx context.Context,
	input dto.RegisterAgentInput,
) (*base.UserProfile, error) {
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

	if !usp.HasPermission(base.PermissionTypeRegisterAgent) {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to create agent"),
		)
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)
	agentProfile := base.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: base.BioData{
			FirstName:   &input.FirstName,
			LastName:    &input.LastName,
			Gender:      input.Gender,
			DateOfBirth: &input.DateOfBirth,
		},
		Role:        base.RoleTypeAgent,
		Permissions: base.RoleTypeAgent.Permissions(),
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

	sup := base.Supplier{
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

	if err := a.notifyNewAgent(ctx, []string{input.Email}, []string{input.PhoneNumber}, *profile.UserBioData.FirstName, otp); err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to send agent registration notifications: %w", err)
	}

	return profile, nil
}

func (a *AgentUseCaseImpl) notifyNewAgent(
	ctx context.Context,
	emails []string,
	phoneNumbers []string,
	name, otp string,
) error {
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

	message := fmt.Sprintf(
		"%sPlease use this One Time PIN: %s to log onto Bewell with your phone number. You will be prompted to change the PIN on login.",
		agentWelcomeMessage,
		otp,
	)
	if err := a.engagement.SendSMS(ctx, phoneNumbers, message); err != nil {
		return fmt.Errorf("unable to send agent registration message: %w", err)
	}

	text := buf.String()
	for _, email := range emails {
		if err := a.engagement.SendMail(ctx, email, text, agentWelcomeEmailSubject); err != nil {
			return fmt.Errorf("unable to send agent registration email: %w", err)
		}
	}

	return nil
}

// ActivateAgent activates/unsuspends the agent profile
func (a *AgentUseCaseImpl) ActivateAgent(ctx context.Context, agentID string) (bool, error) {
	ctx, span := tracer.Start(ctx, "ActivateAgent")
	defer span.End()

	// Check logged in user has permissions/role of employee
	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	usp, err := a.repo.GetUserProfileByUID(ctx, p.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	if !usp.HasPermission(base.PermissionTypeUnsuspendAgent) {
		return false, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to activate agent"),
		)
	}

	agent, err := a.repo.GetUserProfileByID(ctx, agentID, true)
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

// DeactivateAgent deacivates/suspends the agent profile
func (a *AgentUseCaseImpl) DeactivateAgent(ctx context.Context, agentID string) (bool, error) {
	ctx, span := tracer.Start(ctx, "DeactivateAgent")
	defer span.End()

	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	usp, err := a.repo.GetUserProfileByUID(ctx, p.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	if !usp.HasPermission(base.PermissionTypeSuspendAgent) {
		return false, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to suspend agent"),
		)
	}

	// Get agent profile using phoneNumber
	agent, err := a.repo.GetUserProfileByID(ctx, agentID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	if agent.Role != base.RoleTypeAgent {
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
	ctx, span := tracer.Start(ctx, "FetchAgents")
	defer span.End()

	profiles, err := a.repo.ListUserProfiles(ctx, base.RoleTypeAgent)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	agents := []*dto.Agent{}

	for _, profile := range profiles {
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
