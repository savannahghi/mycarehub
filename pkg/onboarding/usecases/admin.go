package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	adminWelcomeMessage = " You can now access the administrator panel"
)

// AdminUseCase represent the business logic required for management of admins
type AdminUseCase interface {
	RegisterAdmin(
		ctx context.Context,
		input dto.RegisterAdminInput,
	) (*profileutils.UserProfile, error)
	FetchAdmins(ctx context.Context) ([]*dto.Admin, error)
	ActivateAdmin(ctx context.Context, input dto.ProfileSuspensionInput) (bool, error)
	DeactivateAdmin(ctx context.Context, input dto.ProfileSuspensionInput) (bool, error)
	FindAdminByNameOrPhone(ctx context.Context, nameOrPhone *string) ([]*dto.Admin, error)
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
func (a *AdminUseCaseImpl) RegisterAdmin(
	ctx context.Context,
	input dto.RegisterAdminInput,
) (*profileutils.UserProfile, error) {
	ctx, span := tracer.Start(ctx, "RegisterAdmin")
	defer span.End()

	// Check logged in user has permissions to register admin
	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	allowed, err := a.repo.CheckIfUserHasPermission(ctx, p.UID, profileutils.CanCreateEmployee)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	if !allowed {
		err = fmt.Errorf("error, user do not have required permissions")
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Get Logged In user profile
	usp, err := a.repo.GetUserProfileByUID(ctx, p.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	phoneNumber, err := a.baseExt.NormalizeMSISDN(input.PhoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)
	userProfile := profileutils.UserProfile{
		PrimaryEmailAddress: &input.Email,
		UserBioData: profileutils.BioData{
			FirstName:   &input.FirstName,
			LastName:    &input.LastName,
			Gender:      input.Gender,
			DateOfBirth: &input.DateOfBirth,
		},
		Role:        profileutils.RoleTypeEmployee,
		Permissions: profileutils.RoleTypeEmployee.Permissions(),
		Roles:       input.RoleIDs,
		CreatedByID: &usp.ID,
		Created:     &timestamp,
	}

	// create a user profile in bewell
	profile, err := a.repo.CreateDetailedUserProfile(ctx, *phoneNumber, userProfile)
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

	adminProfile := domain.AdminProfile{
		ID:             uuid.New().String(),
		ProfileID:      profile.ID,
		OrganizationID: SavannahSladeCode,
	}

	err = a.repo.CreateAdminProfile(ctx, adminProfile)

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

	message := fmt.Sprintf(domain.WelcomeMessage, input.FirstName, otp)
	message += adminWelcomeMessage

	if err := a.engagement.SendSMS(ctx, []string{*phoneNumber}, message); err != nil {
		return nil, fmt.Errorf("unable to send admin registration message: %w", err)
	}

	return profile, nil
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
		// Retrieve the admin PIN
		pin, err := a.repo.GetPINByProfileID(ctx, profile.ID)
		if err != nil {
			utils.RecordSpanError(span, err)
			// the error is wrapped already. No need to wrap it again
			return nil, err
		}

		roles, err := a.repo.GetRolesByIDs(ctx, profile.Roles)
		if err != nil {
			utils.RecordSpanError(span, err)
			// the error is wrapped already. No need to wrap it again
			return nil, err
		}

		// role output
		ro := []dto.RoleOutput{}
		for _, role := range *roles {
			perms, err := role.Permissions(ctx)
			if err != nil {
				utils.RecordSpanError(span, err)
				return nil, err
			}
			o := dto.RoleOutput{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Active:      role.Active,
				Scopes:      role.Scopes,
				Permissions: perms,
			}

			ro = append(ro, o)
		}

		admin := &dto.Admin{
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
			Roles:                   ro,
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

// ActivateAdmin activates/unsuspend the admin profile
func (a *AdminUseCaseImpl) ActivateAdmin(
	ctx context.Context,
	input dto.ProfileSuspensionInput,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "ActivateAdmin")
	defer span.End()

	// Check logged in user has permissions/role of employee
	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	allowed, err := a.repo.CheckIfUserHasPermission(ctx, p.UID, profileutils.CanRemoveEmployee)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	if !allowed {
		return false, fmt.Errorf("error, user do not have the permissions to activate an employee")
	}

	admin, err := a.repo.GetUserProfileByID(ctx, input.ID, true)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	err = a.repo.UpdateSuspended(ctx, admin.ID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// DeactivateAdmin deactivates/suspends the admin profile
func (a *AdminUseCaseImpl) DeactivateAdmin(
	ctx context.Context,
	input dto.ProfileSuspensionInput,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "DeactivateAdmin")
	defer span.End()

	p, err := a.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	allowed, err := a.repo.CheckIfUserHasPermission(ctx, p.UID, profileutils.CanRemoveEmployee)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	if !allowed {
		return false, fmt.Errorf(
			"error, user do not have the permissions to deactivate an employee",
		)
	}

	// Get admin profile using phoneNumber
	admin, err := a.repo.GetUserProfileByID(ctx, input.ID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.InternalServerError(err)
	}

	err = a.repo.UpdateSuspended(ctx, admin.ID, true)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return true, nil
}

// FindAdminByNameOrPhone is used to find an Admin using their phone number
func (a *AdminUseCaseImpl) FindAdminByNameOrPhone(
	ctx context.Context,
	nameOrPhone *string,
) ([]*dto.Admin, error) {
	ctx, span := tracer.Start(ctx, "FindAdminByNameOrPhone")
	defer span.End()

	profiles, err := a.repo.ListUserProfiles(ctx, profileutils.RoleTypeEmployee)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.UserNotFoundError(err)
	}

	admins := []*dto.Admin{}

	for _, profile := range profiles {

		fullName := strings.ToLower(fmt.Sprintf("%v %v", *profile.UserBioData.FirstName, *profile.UserBioData.LastName))
		phoneNumber := profile.PrimaryPhone

		if strings.Contains(*phoneNumber, *nameOrPhone) || strings.Contains(fullName, strings.ToLower(*nameOrPhone)) {
			admin := dto.Admin{
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
			admins = append(admins, &admin)
		}
	}

	return admins, nil
}
