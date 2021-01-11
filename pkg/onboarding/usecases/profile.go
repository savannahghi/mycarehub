package usecases

import (
	"context"
	"fmt"
	"strings"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// ProfileUseCase represents all the profile business logi
type ProfileUseCase interface {
	UserProfile(ctx context.Context) (*base.UserProfile, error)
	GetProfileByID(ctx context.Context, id *string) (*base.UserProfile, error)
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
	UpdateBioData(ctx context.Context, data base.BioData) error
	GetUserProfileByUID(
		ctx context.Context,
		UID string,
	) (*base.UserProfile, error)

	// masks phone number.
	MaskPhoneNumbers(phones []string) []string
}

// ProfileUseCaseImpl represents usecase implementation object
type ProfileUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewProfileUseCase returns a new a onboarding usecase
func NewProfileUseCase(r repository.OnboardingRepository) ProfileUseCase {
	return &ProfileUseCaseImpl{r}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}
	return profile, nil
}

// GetProfileByID returns the profile identified by the indicated ID
func (p *ProfileUseCaseImpl) GetProfileByID(ctx context.Context, id *string) (*base.UserProfile, error) {
	profile, err := p.onboardingRepository.GetUserProfileByID(ctx, *id)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}
	return profile, nil
}

// UpdatePrimaryPhoneNumber updates the primary phone number of a specific user profile
// this should be called after a prior check of uniqueness is done
// this call if valid for both unauthenticated  rest and authenticated graphql. We use `useContext` to determine
// which mode to fetch the user profile
func (p *ProfileUseCaseImpl) UpdatePrimaryPhoneNumber(ctx context.Context, phone string, useContext bool) error {

	var profile *base.UserProfile

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}

	// fetch the user profile
	if useContext {
		uid, err := base.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}
		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, uid)
		if err != nil {
			return exceptions.ProfileNotFoundError(err)
		}

	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
		if err != nil {
			return exceptions.ProfileNotFoundError(err)
		}

	}

	previousPrimaryPhone := profile.PrimaryPhone
	previousSecondaryPhones := profile.SecondaryPhoneNumbers

	if err := p.UpdatePrimaryPhoneNumber(ctx, phone, useContext); err != nil {
		return err
	}

	// removes the new primary phone number from the list of selected primary phones and addes the previus primary phone number
	// into the list of secondary phone numbers
	newSecPhones := func(oldSecondaryPhones []string, oldPrimaryPhone string, newPrimaryPhone string) []string {
		n := []string{}
		for _, phone := range oldSecondaryPhones {
			if phone != newPrimaryPhone {
				n = append(n, phone)
			}
		}
		n = append(n, oldPrimaryPhone)

		return n
	}(previousSecondaryPhones, *previousPrimaryPhone, *phoneNumber)

	if err := p.UpdateSecondaryPhoneNumbers(ctx, newSecPhones); err != nil {
		return err
	}

	return nil
}

// UpdatePrimaryEmailAddress updates primary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
// this call is only valid via graphql api
func (p *ProfileUseCaseImpl) UpdatePrimaryEmailAddress(ctx context.Context, emailAddress string) error {

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfilePrimaryEmailAddress(ctx, p.onboardingRepository, emailAddress)
}

// UpdateSecondaryPhoneNumbers updates secondary phone numberss of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryPhoneNumbers(ctx context.Context, phoneNumbers []string) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfileSecondaryPhoneNumbers(ctx, p.onboardingRepository, phoneNumbers)
}

// UpdateSecondaryEmailAddresses updates secondary email address of a specific user profile
// this should be called after a prior check of uniqueness is done
func (p *ProfileUseCaseImpl) UpdateSecondaryEmailAddresses(ctx context.Context, emailAddresses []string) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfileSecondaryEmailAddresses(ctx, p.onboardingRepository, emailAddresses)
}

// UpdateVerifiedUIDS updates the profile's verified uids
func (p *ProfileUseCaseImpl) UpdateVerifiedUIDS(ctx context.Context, uids []string) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfileVerifiedUIDS(ctx, p.onboardingRepository, uids)
}

// UpdateVerifiedIdentifiers updates the profile's verified identifiers
func (p *ProfileUseCaseImpl) UpdateVerifiedIdentifiers(ctx context.Context, identifiers []base.VerifiedIdentifier) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfileVerifiedIdentifiers(ctx, p.onboardingRepository, identifiers)
}

// UpdateSuspended updates primary suspend attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdateSuspended(ctx context.Context, status bool, phone string, useContext bool) error {
	var profile *base.UserProfile

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	// fetch the user profile
	if useContext {
		uid, err := base.GetLoggedInUserUID(ctx)
		if err != nil {
			return exceptions.UserNotFoundError(err)
		}
		profile, err = p.onboardingRepository.GetUserProfileByUID(ctx, uid)
		if err != nil {
			return exceptions.ProfileNotFoundError(err)
		}
	} else {
		profile, err = p.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
		if err != nil {
			return exceptions.ProfileNotFoundError(err)
		}

	}
	return profile.UpdateProfileSuspended(ctx, p.onboardingRepository, status)
}

// UpdatePhotoUploadID updates photouploadid attribute of a specific user profile
func (p *ProfileUseCaseImpl) UpdatePhotoUploadID(ctx context.Context, uploadID string) error {

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	return profile.UpdateProfilePhotoUploadID(ctx, p.onboardingRepository, uploadID)

}

// UpdateCovers updates primary covers of a specific user profile
func (p *ProfileUseCaseImpl) UpdateCovers(ctx context.Context, covers []base.Cover) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}
	return profile.UpdateProfileCovers(ctx, p.onboardingRepository, covers)

}

// UpdatePushTokens updates primary push tokens of a specific user profile
func (p *ProfileUseCaseImpl) UpdatePushTokens(ctx context.Context, pushToken string, retire bool) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
	}

	if retire {
		// remove the supplied push token then update the profile
		previousTokens := profile.PushTokens
		newTokens := []string{}
		for _, token := range previousTokens {
			if token != pushToken {
				newTokens = append(newTokens, token)
			}
		}

		if err := profile.UpdateProfilePushTokens(ctx, p.onboardingRepository, newTokens); err != nil {
			return err
		}

	}

	return profile.UpdateProfilePushTokens(ctx, p.onboardingRepository, []string{pushToken})
}

// UpdateBioData updates primary biodata of a specific user profile
func (p *ProfileUseCaseImpl) UpdateBioData(ctx context.Context, data base.BioData) error {

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return exceptions.ProfileNotFoundError(err)
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
	return p.onboardingRepository.GetUserProfileByUID(ctx, UID)
}
