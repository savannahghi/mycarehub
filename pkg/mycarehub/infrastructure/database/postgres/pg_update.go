package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// ReactivateFacility changes the status of an active facility from false to true
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.ReactivateFacility(ctx, mflCode)
}

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.InactivateFacility(ctx, mflCode)
}

// AcceptTerms can be used to accept or review terms of service
func (d *MyCareHubDb) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}
	return d.update.AcceptTerms(ctx, userID, termsID)
}

// UpdateUserFailedLoginCount increments a user's failed login count in the event where they fail to
// log into the app e.g when an invalid pin is passed
func (d *MyCareHubDb) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserFailedLoginCount(ctx, userID, failedLoginAttempts)
}

// UpdateUserLastFailedLoginTime updates the failed login time for a user in case an error occurs while logging in
func (d *MyCareHubDb) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserLastFailedLoginTime(ctx, userID)
}

// UpdateUserNextAllowedLoginTime updates the user's next allowed login time. This field is used to check whether we can
// allow a user to log in immediately or wait for some time before retrying the login process.
func (d *MyCareHubDb) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserNextAllowedLoginTime(ctx, userID, nextAllowedLoginTime)
}

// UpdateUserLastSuccessfulLoginTime updates the user's last successful login time to the current time in case a user
// successfully logs into the app
func (d *MyCareHubDb) UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserLastSuccessfulLoginTime(ctx, userID)
}

// SetNickName is used to set the user's nickname
func (d *MyCareHubDb) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	if userID == nil || nickname == nil {
		return false, fmt.Errorf("userID or nickname cannot be empty ")
	}

	return d.update.SetNickName(ctx, userID, nickname)
}

// UpdateUserPinChangeRequiredStatus updates the user's pin change required from true to false. It'll be used to
// determine the onboarding journey for a user.
func (d *MyCareHubDb) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserPinChangeRequiredStatus(ctx, userID, flavour)
}

// InvalidatePIN invalidates a pin that is linked to the user profile.
// This is done by toggling the IsValid field to false
func (d *MyCareHubDb) InvalidatePIN(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.InvalidatePIN(ctx, userID)
}

// UpdateIsCorrectSecurityQuestionResponse updates the user's security question response
func (d *MyCareHubDb) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.UpdateIsCorrectSecurityQuestionResponse(ctx, userID, isCorrectSecurityQuestionResponse)
}

// ShareContent updates content share count
func (d *MyCareHubDb) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	if input.Validate() != nil {
		return false, fmt.Errorf("input cannot be empty")
	}
	return d.update.ShareContent(ctx, input)
}

//BookmarkContent updates the user's bookmark status for a content
func (d *MyCareHubDb) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	return d.update.BookmarkContent(ctx, userID, contentID)
}

// UnBookmarkContent removes the bookmark for a given user
func (d *MyCareHubDb) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	return d.update.UnBookmarkContent(ctx, userID, contentID)
}

// LikeContent updates the number of likes for a particular content
func (d *MyCareHubDb) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	return d.update.LikeContent(ctx, userID, contentID)
}

// UnlikeContent updates the number of likes for a particular content
func (d *MyCareHubDb) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	return d.update.UnlikeContent(ctx, userID, contentID)
}

// ViewContent gets a content item and updates the view count
func (d *MyCareHubDb) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return d.update.ViewContent(ctx, userID, contentID)
}
