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
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
)

const (
	adminWelcomeMessage      = "You have been successfully registered as an admin. We look forward to working with you."
	adminWelcomeEmailSubject = "Successfully registered as an admin"
)

// AdminUseCase represent the business logic required for management of admins
type AdminUseCase interface {
	RegisterAdmin(ctx context.Context, input dto.RegisterAdminInput) (*profileutils.UserProfile, error)
	FetchAdmins(ctx context.Context) ([]*dto.Admin, error)
}

// AdminUseCaseImpl  represents usecase implementation object
type AdminUseCaseImpl struct {
	repo       repository.OnboardingRepository
	engagement engagement.ServiceEngagement
	baseExt    extension.BaseExtension
	pin        UserPINUseCases
}

// NewAdminUseCases returns a new a onboarding usecase
func NewAdminUseCases(
	r repository.OnboardingRepository,
	eng engagement.ServiceEngagement,
	ext extension.BaseExtension,
	pin UserPINUseCases,
) AdminUseCase {

	return &AdminUseCaseImpl{
		repo:       r,
		engagement: eng,
		baseExt:    ext,
		pin:        pin,
	}
}

// RegisterAdmin creates a new Admin in bewell
func (a *AdminUseCaseImpl) RegisterAdmin(ctx context.Context, input dto.RegisterAdminInput) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "RegisterAdmin")
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

	if !usp.HasPermission(profileutils.PermissionTypeCreateAdmin) {
		return nil, exceptions.RoleNotValid(fmt.Errorf("error: logged in user does not have permissions to create admin"))
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)
	adminProfile := profileutils.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: profileutils.BioData{
			FirstName:   &input.FirstName,
			LastName:    &input.LastName,
			Gender:      input.Gender,
			DateOfBirth: &input.DateOfBirth,
		},
		Role:        profileutils.RoleTypeEmployee,
		Permissions: profileutils.RoleTypeEmployee.Permissions(),
		CreatedByID: &usp.ID,
		Created:     &timestamp,
	}

	// create a user profile in bewell
	profile, err := a.repo.CreateDetailedUserProfile(ctx, *msisdn, adminProfile)
	if err != nil {
		utils.RecordSpanError(span, err)
		// wrapped error
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
		utils.RecordSpanError(span, err)
		// wrapped error
		return nil, err
	}

	otp, err := a.pin.SetUserTempPIN(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		// wrapped error
		return nil, err
	}

	if err := a.notifyNewAdmin(ctx, []string{input.Email}, []string{input.PhoneNumber}, *profile.UserBioData.FirstName, otp); err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to send admin registration notifications: %w", err)
	}

	return profile, nil
}

func (a *AdminUseCaseImpl) notifyNewAdmin(ctx context.Context, emails []string, phoneNumbers []string, name, otp string) error {
	type pin struct {
		Name string
		Pin  string
	}

	t := template.Must(template.New("adminApprovalEmail").Parse(utils.AdminApprovalEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, pin{name, otp})
	if err != nil {
		log.Fatalf("error while generating admin approval email template: %s", err)
	}

	message := fmt.Sprintf("%sPlease use this One Time PIN: %s to log onto Bewell with your phone number. For enquiries call us on 0790360360", adminWelcomeMessage, otp)
	if err := a.engagement.SendSMS(ctx, phoneNumbers, message); err != nil {
		return fmt.Errorf("unable to send admin registration message: %w", err)
	}

	text := buf.String()
	for _, email := range emails {
		if err := a.engagement.SendMail(ctx, email, text, adminWelcomeEmailSubject); err != nil {
			return fmt.Errorf("unable to send admin registration email: %w", err)
		}
	}

	return nil
}

// FetchAdmins fetches registered admins
func (a *AdminUseCaseImpl) FetchAdmins(ctx context.Context) ([]*dto.Admin, error) {
	ctx, span := tracer.Start(ctx, "FetchAdmins")
	defer span.End()

	profiles, err := a.repo.ListUserProfiles(ctx, profileutils.RoleTypeEmployee)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	admins := []*dto.Admin{}

	for _, profile := range profiles {
		admin := &dto.Admin{
			ID:                      profile.ID,
			PhotoUploadID:           profile.PhotoUploadID,
			UserBioData:             profile.UserBioData,
			PrimaryPhone:            *profile.PrimaryPhone,
			PrimaryEmailAddress:     *profile.PrimaryEmailAddress,
			SecondaryPhoneNumbers:   profile.SecondaryPhoneNumbers,
			SecondaryEmailAddresses: profile.SecondaryEmailAddresses,
			TermsAccepted:           profile.TermsAccepted,
			Suspended:               profile.Suspended,
		}

		admins = append(admins, admin)
	}

	return admins, nil
}
