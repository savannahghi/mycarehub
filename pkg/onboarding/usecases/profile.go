package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	// EmailsAttribute is an attribute that represents
	// a user profile's email addresses
	EmailsAttribute = "emails"

	// PhoneNumbersAttribute is an attribute that represents
	// a user profile's phone numbers
	PhoneNumbersAttribute = "phonenumbers"

	// FCMTokensAttribute is an attribute that represents
	// a user profile's FCM push tokens
	FCMTokensAttribute = "tokens"
)

// VerifyEmailNudgeTitle is the title defined in the `engagement service`
// for the `VerifyEmail` nudge
const VerifyEmailNudgeTitle = "Add Primary Email Address"

// ProfileUseCase represents all the profile business logic
type ProfileUseCase interface {
	// profile related
	UserProfile(ctx context.Context) (*base.UserProfile, error)
	GetProfileByID(ctx context.Context, id *string) (*base.UserProfile, error)
	UpdateUserName(ctx context.Context, userName string) error
	UpdatePrimaryPhoneNumber(ctx context.Context, phoneNumber string, useContext bool) error
	UpdatePrimaryEmailAddress(ctx context.Context, emailAddress string) error
	UpdateSecondaryPhoneNumbers(ctx context.Context, phoneNumbers []string) error
	UpdateSecondaryEmailAddresses(ctx context.Context, emailAddresses []string) error
	UpdateVerifiedIdentifiers(ctx context.Context, identifiers []base.VerifiedIdentifier) error
	UpdateVerifiedUIDS(ctx context.Context, uids []string) error
	UpdateSuspended(ctx context.Context, status bool, phoneNumber string, useContext bool) error
	UpdatePhotoUploadID(ctx context.Context, uploadID string) error
	UpdateCovers(ctx context.Context, covers []base.Cover) error
	UpdatePushTokens(ctx context.Context, pushToken string, retire bool) error
	UpdatePermissions(ctx context.Context, perms []base.PermissionType) error
	AddAdminPermsToUser(ctx context.Context, phone string) error
	RemoveAdminPermsToUser(ctx context.Context, phone string) error
	UpdateBioData(ctx context.Context, data base.BioData) error
	GetUserProfileByUID(
		ctx context.Context,
		UID string,
	) (*base.UserProfile, error)

	// masks phone number.
	MaskPhoneNumbers(phones []string) []string
	// called to set the primary phone number of a specific profile.
	// useContext is used to mark under which scenario the method is been called.
	SetPrimaryPhoneNumber(
		ctx context.Context,
		phoneNumber string,
		otp string,
		useContext bool,
	) error
	SetPrimaryEmailAddress(
		ctx context.Context,
		emailAddress string,
		otp string,
	) error

	// checks whether a phone number has been registered by another user. Checks both primary and
	// secondary phone numbers. If the the phone number is foreign, it returns false
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)

	// check whether a email has been registered by another user. Checks both primary and
	// secondary emails. If the the phone number is foreign, it returns false
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// called to remove specific secondary phone numbers from the user's profile.'
	RetireSecondaryPhoneNumbers(ctx context.Context, toRemovePhoneNumbers []string) (bool, error)

	// called to remove specific secondary email addresses from the user's profile.
	RetireSecondaryEmailAddress(ctx context.Context, toRemoveEmails []string) (bool, error)

	GetUserProfileAttributes(
		ctx context.Context,
		UIDs []string,
		attribute string,
	) (map[string][]string, error)

	ConfirmedEmailAddresses(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ConfirmedPhoneNumbers(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ValidFCMTokens(
		ctx context.Context,
		UIDs []string,
	) (map[string][]string, error)

	ProfileAttributes(
		ctx context.Context,
		UIDs []string,
		attribute string,
	) (map[string][]string, error)

	SetupAsExperimentParticipant(ctx context.Context, participate *bool) (bool, error)

	AddAddress(
		ctx context.Context,
		input resources.UserAddressInput,
		addressType base.AddressType,
	) (*base.Address, error)

	GetAddresses(ctx context.Context) (*domain.UserAddresses, error)

	GetUserCommunicationsSettings(ctx context.Context) (*base.UserCommunicationsSetting, error)

	SetUserCommunicationsSettings(
		ctx context.Context,
		allowWhatsApp *bool,
		allowTextSms *bool,
		allowPush *bool,
		allowEmail *bool,
	) (*base.UserCommunicationsSetting, error)
}

// ProfileUseCaseImpl represents usecase implementation object
type ProfileUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	baseExt              extension.BaseExtension
	engagement           engagement.ServiceEngagement
}

// NewProfileUseCase returns a new a onboarding usecase
func NewProfileUseCase(
	r repository.OnboardingRepository,
	ext extension.BaseExtension,
	eng engagement.ServiceEngagement,
) ProfileUseCase {
	return &ProfileUseCaseImpl{
		onboardingRepository: r,
		baseExt:              ext,
		engagement:           eng,
	}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	return profile, nil
}

// GetProfileByID returns the profile identified by the indicated ID
func (p *ProfileUseCaseImpl) GetProfileByID(
	ctx context.Context,
	id *string,
) (*base.UserProfile, error) {
	profile, err := p.onboardingRepository.GetUserProfileByID(ctx, *id, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	return profile, nil
}

// UpdateUserName updates the user username.
func (p *ProfileUseCaseImpl) UpdateUserName(ctx context.Context, userName string) error {
	profile, err := p.UserProfile(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return profile.UpdateProfileUserName(ctx, p.onboardingRepository, userName)
}

// UpdatePrimaryPhoneNumber updates the primary phone number of a specific user profile
// this should be called after a prior check of uniqueness is done
// We use `useContext` to determine
// which mode to fetch the user profile
func (p *ProfileUseCaseImpl) UpdatePrimaryPhoneNumber(
	ctx context.Context,
	phone string,
	useContext bool,
) error {

	var profile *base.UserProfile

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	// fetch the user profile
	if useContext {
		uid, err := p.baseExt.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}
		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

	}

	previousPrimaryPhone := profile.PrimaryPhone
	secondaryPhones := profile.SecondaryPhoneNumbers
	if err := profile.UpdateProfilePrimaryPhoneNumber(ctx, p.onboardingRepository, phone); err != nil {
		return err
	}

	// check if number to be set as primary exists in the list of secondary phones
	index, exists := utils.FindItem(secondaryPhones, *phoneNumber)
	if exists {
		// remove the phoneNumber from the Secondary Phones slice
		secondaryPhones = append(secondaryPhones[:index], secondaryPhones[index+1:]...)
	}

	secondaryPhones = append(secondaryPhones, *previousPrimaryPhone)

	if len(secondaryPhones) >= 1 {
		if err := profile.UpdateProfileSecondaryPhoneNumbers(ctx, p.onboardingRepository, secondaryPhones); err != nil {
			return err
		}
	}

	return nil
}

// UpdatePrimaryEmailAddress updates primary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdatePrimaryEmailAddress(
	ctx context.Context,
	emailAddress string,
) error {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err := profile.UpdateProfilePrimaryEmailAddress(ctx, p.onboardingRepository, emailAddress); err != nil {
		return err
	}

	previousPrimaryEmail := profile.PrimaryEmailAddress
	secondaryEmails := profile.SecondaryEmailAddresses

	if profile.PrimaryEmailAddress != nil {
		// Check if the email to be set as primary exists in the list
		// of secondary emails.
		index, exists := utils.FindItem(secondaryEmails, emailAddress)
		if exists {
			secondaryEmails = append(secondaryEmails[:index], secondaryEmails[index+1:]...)
		}

		secondaryEmails = append(secondaryEmails, *previousPrimaryEmail)

		if len(secondaryEmails) >= 1 {
			if err := profile.UpdateProfileSecondaryEmailAddresses(ctx, p.onboardingRepository, secondaryEmails); err != nil {
				return err
			}
		}
	}

	return nil

}

// UpdateSecondaryPhoneNumbers updates secondary phone numbers of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryPhoneNumbers(
	ctx context.Context,
	phoneNumbers []string,
) error {
	uniquePhones := []string{}
	// assert that the phone numbers are unique
	for _, phone := range phoneNumbers {
		exist, err := p.CheckPhoneExists(ctx, phone)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

		if !exist {
			uniquePhones = append(uniquePhones, phone)
		}
	}

	if len(uniquePhones) >= 1 {
		uid, err := p.baseExt.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}

		profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

		return profile.UpdateProfileSecondaryPhoneNumbers(ctx, p.onboardingRepository, phoneNumbers)
	}

	// throw an error indicating the phone number(s) is/are already in the use
	return exceptions.CheckPhoneNumberExistError()
}

// UpdateSecondaryEmailAddresses updates secondary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryEmailAddresses(
	ctx context.Context,
	emailAddresses []string,
) error {
	uniqueEmails := []string{}
	for _, email := range emailAddresses {
		exist, err := p.CheckEmailExists(ctx, email)
		if err != nil {
			return err
		}

		if !exist {
			uniqueEmails = append(uniqueEmails, email)
		}
	}

	if len(uniqueEmails) >= 1 {
		uid, err := p.baseExt.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}

		profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
		if err != nil {
			return err
		}

		if profile.PrimaryEmailAddress != nil {
			return profile.UpdateProfileSecondaryEmailAddresses(
				ctx,
				p.onboardingRepository,
				uniqueEmails,
			)
		}

		// internal error. primary email addresses must be present before addong secondary email
		// addresses.
		return exceptions.InternalServerError(
			fmt.Errorf(
				"primary email addresses must be present before adding secondary email addresses",
			),
		)

	}

	// throw an error indicating the email(s) is/are already in the use
	return exceptions.CheckEmailExistError()
}

// UpdateVerifiedUIDS updates the profile's verified uids
func (p *ProfileUseCaseImpl) UpdateVerifiedUIDS(ctx context.Context, uids []string) error {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return profile.UpdateProfileVerifiedUIDS(ctx, p.onboardingRepository, uids)
}

// UpdateVerifiedIdentifiers updates the profile's verified identifiers
func (p *ProfileUseCaseImpl) UpdateVerifiedIdentifiers(
	ctx context.Context,
	identifiers []base.VerifiedIdentifier,
) error {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	return profile.UpdateProfileVerifiedIdentifiers(ctx, p.onboardingRepository, identifiers)
}

// UpdateSuspended updates primary suspend attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdateSuspended(
	ctx context.Context,
	status bool,
	phone string,
	useContext bool,
) error {
	var profile *base.UserProfile

	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	// fetch the user profile
	if useContext {
		uid, err := p.baseExt.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}
		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}
	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber, false)
		if err != nil {
			// this is a wrapped error. No need to wrap it again
			return err
		}

	}
	return profile.UpdateProfileSuspended(ctx, p.onboardingRepository, status)
}

// UpdatePhotoUploadID updates photouploadid attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdatePhotoUploadID(ctx context.Context, uploadID string) error {

	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	return profile.UpdateProfilePhotoUploadID(ctx, p.onboardingRepository, uploadID)

}

// UpdateCovers updates primary covers of a specific user profile
func (p *ProfileUseCaseImpl) UpdateCovers(ctx context.Context, covers []base.Cover) error {
	if len(covers) == 0 {
		return fmt.Errorf("no covers to update found")
	}

	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		return err
	}

	return profile.UpdateProfileCovers(
		ctx,
		p.onboardingRepository,
		utils.AddHashToCovers(covers),
	)
}

// UpdatePushTokens updates primary push tokens of a specific user profile.
func (p *ProfileUseCaseImpl) UpdatePushTokens(
	ctx context.Context,
	pushToken string,
	retire bool,
) error {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}

	if retire {
		// check if supplied push token exist in the list of pushtokens
		index, exist := utils.FindItem(profile.PushTokens, pushToken)
		if exist {
			// remove it from the list of push tokens
			profile.PushTokens = append(profile.PushTokens[:index], profile.PushTokens[index+1:]...)
		}

		return profile.UpdateProfilePushTokens(ctx, p.onboardingRepository, profile.PushTokens)

	}
	newToken := []string{}
	newTokens := append(newToken, pushToken)

	return profile.UpdateProfilePushTokens(ctx, p.onboardingRepository, newTokens)
}

// UpdatePermissions updates the profiles permissions
func (p *ProfileUseCaseImpl) UpdatePermissions(
	ctx context.Context,
	perms []base.PermissionType,
) error {
	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return profile.UpdateProfilePermissions(ctx, p.onboardingRepository, perms)
}

// AddAdminPermsToUser updates the profiles permissions
func (p *ProfileUseCaseImpl) AddAdminPermsToUser(ctx context.Context, phone string) error {
	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		return err
	}
	perms := base.DefaultSuperAdminPermissions
	return profile.UpdateProfilePermissions(ctx, p.onboardingRepository, perms)
}

// RemoveAdminPermsToUser updates the profiles permissions by removing the admin permissions
// This also flips back userProfile field IsAdmin to false
func (p *ProfileUseCaseImpl) RemoveAdminPermsToUser(ctx context.Context, phone string) error {
	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		return err
	}
	permissions := profile.Permissions
	if len(permissions) >= 1 {
		permissions = nil
	}
	return profile.UpdateProfilePermissions(ctx, p.onboardingRepository, permissions)
}

// UpdateBioData updates primary biodata of a specific user profile
func (p *ProfileUseCaseImpl) UpdateBioData(ctx context.Context, data base.BioData) error {

	uid, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	return profile.UpdateProfileBioData(ctx, p.onboardingRepository, data)
}

// MaskPhoneNumbers masks phone number. the masked phone numbers will be in the form +254700***123
func (p *ProfileUseCaseImpl) MaskPhoneNumbers(phones []string) []string {
	masked := make([]string, 0, len(phones))
	for _, num := range phones {
		var b strings.Builder
		max := len(num)
		for i, p := range num {
			if i+1 == max-3 || i+1 == max-4 || i+1 == max-5 {
				fmt.Fprintf(&b, "*")
			} else {
				fmt.Fprint(&b, string(p))
			}
		}
		masked = append(masked, b.String())
	}
	return masked
}

// GetUserProfileByUID retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) GetUserProfileByUID(
	ctx context.Context,
	UID string,
) (*base.UserProfile, error) {
	return p.onboardingRepository.GetUserProfileByUID(ctx, UID, false)
}

// SetPrimaryPhoneNumber set the primary phone number of the user after verifying the otp code
func (p *ProfileUseCaseImpl) SetPrimaryPhoneNumber(
	ctx context.Context,
	phoneNumber string,
	otp string,
	useContext bool,
) error {
	// verify otp code
	verified, err := p.engagement.VerifyOTP(
		ctx,
		phoneNumber,
		otp,
	)
	if err != nil {
		return exceptions.VerifyOTPError(err)
	}

	if !verified {
		return exceptions.VerifyOTPError(nil)
	}

	// now set the primary phone number
	if err := p.UpdatePrimaryPhoneNumber(ctx, phoneNumber, useContext); err != nil {
		return err
	}

	return nil
}

// SetPrimaryEmailAddress set the primary email address of the user after verifying the otp code
func (p *ProfileUseCaseImpl) SetPrimaryEmailAddress(
	ctx context.Context,
	emailAddress string,
	otp string,
) error {
	UID, err := p.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	// verify otp code
	verified, err := p.engagement.VerifyEmailOTP(
		ctx,
		emailAddress,
		otp,
	)
	if err != nil {
		return exceptions.VerifyOTPError(err)
	}
	if !verified {
		return exceptions.VerifyOTPError(nil)
	}
	if err := p.UpdatePrimaryEmailAddress(ctx, emailAddress); err != nil {
		return err
	}

	// The `VerifyEmail` nudge is by default created for both flavours, `PRO`
	// and `CONSUMER`, thus if a user adds and verifies their `Primary Email`
	// we need to `Resolve` the nudge for this user in both flavours
	// Resolve the nudge in `CONSUMER`
	go func() {
		cons := func() error {
			return p.engagement.ResolveDefaultNudgeByTitle(
				UID,
				base.FlavourConsumer,
				VerifyEmailNudgeTitle,
			)
		}
		if err := backoff.Retry(
			cons,
			backoff.NewExponentialBackOff(),
		); err != nil {
			logrus.Error(err)
		}

		pro := func() error {
			return p.engagement.ResolveDefaultNudgeByTitle(
				UID,
				base.FlavourPro,
				VerifyEmailNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			backoff.NewExponentialBackOff(),
		); err != nil {
			logrus.Error(err)
		}
	}()

	return nil
}

// CheckPhoneExists checks whether a phone number has been registered by another user.
// Checks both primary and secondary phone numbers.
func (p *ProfileUseCaseImpl) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	phoneNumber, err := p.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}
	exists, err := p.onboardingRepository.CheckIfPhoneNumberExists(ctx, *phoneNumber)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckEmailExists checks whether a email has been registered by another user.
// Checks both primary and secondary emails.
func (p *ProfileUseCaseImpl) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := p.onboardingRepository.CheckIfEmailExists(ctx, email)
	if err != nil {
		return false, exceptions.CheckEmailExistError()
	}
	return exists, nil
}

// RetireSecondaryPhoneNumbers overwrites an existing secondary phone number,
// if any, with the provided phone number.
func (p *ProfileUseCaseImpl) RetireSecondaryPhoneNumbers(
	ctx context.Context,
	phoneNumbers []string,
) (bool, error) {
	profile, err := p.UserProfile(ctx)
	if err != nil {
		return false, err
	}

	secondaryPhoneNumbers := profile.SecondaryPhoneNumbers
	if len(secondaryPhoneNumbers) == 0 {
		return false, exceptions.SecondaryResourceHardResetError()
	}

	for _, phoneNo := range phoneNumbers {
		// Check if the passed number exists in the list of secondary numbers
		index, exists := utils.FindItem(secondaryPhoneNumbers, phoneNo)
		if exists {
			// remove the passed number from the list of secondary numbers
			secondaryPhoneNumbers = append(secondaryPhoneNumbers[:index], secondaryPhoneNumbers[index+1:]...)
		} else {
			return false, exceptions.RecordDoesNotExistError(fmt.Errorf("record does not exist"))
		}
	}

	if err := p.onboardingRepository.HardResetSecondaryPhoneNumbers(
		ctx,
		profile,
		secondaryPhoneNumbers,
	); err != nil {
		return false, err
	}

	return true, nil
}

// RetireSecondaryEmailAddress removes specific secondary email addresses from the user's profile.
func (p *ProfileUseCaseImpl) RetireSecondaryEmailAddress(
	ctx context.Context,
	emailAddresses []string,
) (bool, error) {
	profile, err := p.UserProfile(ctx)
	if err != nil {
		return false, err
	}

	secondaryEmails := profile.SecondaryEmailAddresses
	if len(secondaryEmails) == 0 {
		return false, exceptions.SecondaryResourceHardResetError()
	}

	for _, email := range emailAddresses {
		// Check if the passed email exists in the list of secondary emails
		index, exists := utils.FindItem(secondaryEmails, email)
		if exists {
			// remove the passed email from the list of secondary emails
			secondaryEmails = append(secondaryEmails[:index], secondaryEmails[index+1:]...)
		} else {
			return false, exceptions.RecordDoesNotExistError(fmt.Errorf("record does not exist"))
		}
	}

	if err := p.onboardingRepository.HardResetSecondaryEmailAddress(
		ctx,
		profile,
		secondaryEmails,
	); err != nil {
		return false, err
	}

	return true, nil
}

// GetUserProfileAttributes takes a slice of UIDs and for each UID,
// it fetches the user profiles confirmed emails, phone numbers and
// FCM push tokens
func (p *ProfileUseCaseImpl) GetUserProfileAttributes(
	ctx context.Context,
	UIDs []string,
	attribute string,
) (map[string][]string, error) {
	output := make(map[string][]string)
	values := []string{}

	for _, UID := range UIDs {
		profile, err := p.onboardingRepository.GetUserProfileByUID(
			ctx,
			UID,
			false,
		)
		if err != nil {
			return output, err
		}

		switch attribute {
		case EmailsAttribute:
			primaryEmail := profile.PrimaryEmailAddress

			if primaryEmail == nil {
				//if not found just show an empty list
				output[UID] = []string{}
				continue
			}

			output[UID] = append(values, *primaryEmail)

			secondaryEmails := profile.SecondaryEmailAddresses
			if len(secondaryEmails) != 0 {
				for _, secondaryEmail := range secondaryEmails {
					output[UID] = append(
						values,
						*primaryEmail,
						secondaryEmail,
					)
				}
			}

		case PhoneNumbersAttribute:
			output[UID] = append(values, *profile.PrimaryPhone)
			secondaryPhones := profile.SecondaryPhoneNumbers
			if len(secondaryPhones) != 0 {
				for _, secondaryPhone := range secondaryPhones {
					output[UID] = append(
						values,
						*profile.PrimaryPhone,
						secondaryPhone,
					)
				}
			}

		case FCMTokensAttribute:
			if len(profile.PushTokens) == 0 {
				// We do not expect there to be a user profile wtihout
				// FCM push tokens, but we can't be too sure
				// if not found just show an empty list
				output[UID] = []string{}
				continue
			}
			output[UID] = append(values, profile.PushTokens...)

		default:
			err := fmt.Errorf("failed to retrieve user profile attribute %s",
				attribute,
			)
			return nil, exceptions.RetrieveRecordError(err)
		}
	}

	return output, nil
}

// ConfirmedEmailAddresses returns verified email addresses for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ConfirmedEmailAddresses(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		EmailsAttribute,
	)
}

// ConfirmedPhoneNumbers returns verified phone numbers for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ConfirmedPhoneNumbers(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		PhoneNumbersAttribute,
	)
}

// ValidFCMTokens returns valid FCM push tokens for
// each of the UID in the slice of UIDs provided
func (p *ProfileUseCaseImpl) ValidFCMTokens(
	ctx context.Context,
	UIDs []string,
) (map[string][]string, error) {
	return p.GetUserProfileAttributes(
		ctx,
		UIDs,
		FCMTokensAttribute,
	)
}

// ProfileAttributes retrieves the user profiles confirmed emails,
// phone numbers and FCM push tokens
func (p *ProfileUseCaseImpl) ProfileAttributes(
	ctx context.Context,
	UIDs []string,
	attribute string,
) (map[string][]string, error) {
	switch attribute {
	case EmailsAttribute:
		return p.ConfirmedEmailAddresses(
			ctx,
			UIDs,
		)

	case PhoneNumbersAttribute:
		return p.ConfirmedPhoneNumbers(
			ctx,
			UIDs,
		)

	case FCMTokensAttribute:
		return p.ValidFCMTokens(
			ctx,
			UIDs,
		)

	default:
		err := fmt.Errorf("failed to retrieve user profile attribute %s",
			attribute,
		)
		return nil, exceptions.RetrieveRecordError(err)
	}
}

// SetupAsExperimentParticipant sets up the logged-in user as an experiment participant.
// An experiment participant will be able to see unstable or otherwise flaged-feature in the UI of
// the app
func (p *ProfileUseCaseImpl) SetupAsExperimentParticipant(
	ctx context.Context,
	participate *bool,
) (bool, error) {
	// fetch the user profile
	pr, err := p.UserProfile(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}

	if *participate {
		// add the user to the list of experiment participants
		return p.onboardingRepository.AddUserAsExperimentParticipant(ctx, pr)
	}

	// remove the user to the list of experiment participants
	return p.onboardingRepository.RemoveUserAsExperimentParticipant(ctx, pr)
}

// AddAddress adds a user's home or work address to thir user's profile
func (p *ProfileUseCaseImpl) AddAddress(
	ctx context.Context,
	input resources.UserAddressInput,
	addressType base.AddressType,
) (*base.Address, error) {
	var address *base.Address
	profile, err := p.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	address = &base.Address{
		// Longitude and latitude coordinates are stored with
		// 15 decimal digits right of the decimal points
		Latitude:         fmt.Sprintf("%.15f", input.Latitude),
		Longitude:        fmt.Sprintf("%.15f", input.Longitude),
		Locality:         input.Locality,
		Name:             input.Name,
		PlaceID:          input.PlaceID,
		FormattedAddress: input.FormattedAddress,
	}
	err = p.onboardingRepository.UpdateAddresses(
		ctx,
		profile.ID,
		*address,
		addressType,
	)
	if err != nil {
		return address, err
	}

	return address, nil
}

// GetAddresses returns a user's home and work addresses
func (p *ProfileUseCaseImpl) GetAddresses(
	ctx context.Context,
) (*domain.UserAddresses, error) {
	profile, err := p.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	var thinHomeAddress domain.ThinAddress
	if profile.HomeAddress != nil {
		homeLatitude, err := strconv.ParseFloat(
			profile.HomeAddress.Latitude,
			64,
		)
		if err != nil {
			return nil, err
		}

		homeLongitude, err := strconv.ParseFloat(
			profile.HomeAddress.Longitude,
			64,
		)
		if err != nil {
			return nil, err
		}

		thinHomeAddress = domain.ThinAddress{
			Latitude:  homeLatitude,
			Longitude: homeLongitude,
		}
	}

	var thinWorkAddress domain.ThinAddress
	if profile.WorkAddress != nil {
		workLatitude, err := strconv.ParseFloat(
			profile.WorkAddress.Latitude,
			64,
		)
		if err != nil {
			return nil, err
		}

		workLongitude, err := strconv.ParseFloat(
			profile.WorkAddress.Longitude,
			64,
		)
		if err != nil {
			return nil, err
		}

		thinWorkAddress = domain.ThinAddress{
			Latitude:  workLatitude,
			Longitude: workLongitude,
		}
	}

	return &domain.UserAddresses{
		HomeAddress: thinHomeAddress,
		WorkAddress: thinWorkAddress,
	}, nil
}

// GetUserCommunicationsSettings  retrives the logged in user communications settings.
func (p *ProfileUseCaseImpl) GetUserCommunicationsSettings(
	ctx context.Context,
) (*base.UserCommunicationsSetting, error) {
	pr, err := p.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return p.onboardingRepository.GetUserCommunicationsSettings(ctx, pr.ID)
}

// SetUserCommunicationsSettings sets the user communication settings
func (p *ProfileUseCaseImpl) SetUserCommunicationsSettings(
	ctx context.Context,
	allowWhatsApp *bool,
	allowTextSms *bool,
	allowPush *bool,
	allowEmail *bool,
) (*base.UserCommunicationsSetting, error) {
	pr, err := p.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return p.onboardingRepository.SetUserCommunicationsSettings(
		ctx,
		pr.ID,
		allowWhatsApp,
		allowTextSms,
		allowPush,
		allowEmail,
	)

}
