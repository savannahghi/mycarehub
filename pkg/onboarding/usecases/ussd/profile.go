package ussd

import (
	"context"

	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

//CreateUsddUserProfile creates and updates a user profile
func (u *Impl) CreateUsddUserProfile(ctx context.Context, phoneNumber string, PIN string, userProfile *dto.UserProfileInput) error {
	ctx, span := tracer.Start(ctx, "CreateUsddUserProfile")
	defer span.End()

	user, err := u.onboardingRepository.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	profile, err := u.onboardingRepository.CreateUserProfile(
		ctx,
		phoneNumber,
		user.UID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}
	_, err = u.pinUsecase.SetUserPIN(
		ctx,
		PIN,
		profile.ID,
	)

	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	_, err = u.onboardingRepository.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.InternalServerError(err)
	}

	data := profileutils.BioData{
		FirstName:   &userFirstName,
		LastName:    &userLastName,
		DateOfBirth: userProfile.DateOfBirth,
	}
	err = u.onboardingRepository.UpdateBioData(ctx, profile.ID, data)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	return nil

}
