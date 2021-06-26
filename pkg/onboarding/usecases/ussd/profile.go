package ussd

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
)

//CreateUsddUserProfile creates and updates a user profile
// TODO FIXME ```getorcreateProfile```
func (u *Impl) CreateUsddUserProfile(ctx context.Context, phoneNumber string, PIN string, userProfile *dto.UserProfileInput) error {
	user, err := u.onboardingRepository.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
	if err != nil {
		return err
	}
	profile, err := u.onboardingRepository.CreateUserProfile(
		ctx,
		phoneNumber,
		user.UID,
	)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	_, err = u.pinUsecase.SetUserPIN(
		ctx,
		PIN,
		profile.ID,
	)

	if err != nil {
		return err
	}
	_, err = u.onboardingRepository.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	data := base.BioData{
		FirstName:   &userFirstName,
		LastName:    &userLastName,
		DateOfBirth: userProfile.DateOfBirth,
	}
	err = u.onboardingRepository.UpdateBioData(ctx, profile.ID, data)
	if err != nil {
		return err
	}
	return nil

}
