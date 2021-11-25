package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// Update represents all `update` operations to the database
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error
	UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error
	UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error
	UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
	UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// ReactivateFacility perfoms the actual re-activation of the facility in the database
func (db *PGInstance) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: false}).
		Updates(&Facility{Active: true}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// InactivateFacility perfoms the actual inactivation of the facility in the database
func (db *PGInstance) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: true}).
		Updates(&Facility{Active: false}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// AcceptTerms perfoms the actual modification of the users data for the terms accepted as well as the id of the terms accepted
func (db *PGInstance) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}

	if err := db.DB.Model(&User{}).Where(&User{UserID: userID}).
		Updates(&User{TermsAccepted: true, AcceptedTermsID: termsID}).Error; err != nil {
		return false, fmt.Errorf("an error occurred while updating the user: %v", err)
	}

	return true, nil
}

// UpdateUserFailedLoginCount updates the user's failed login count field in an event where a user fails to
// log into the app
func (db *PGInstance) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(map[string]interface{}{
		"failed_login_count": failedLoginAttempts,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserLastFailedLoginTime updates the user's last failed login time
func (db *PGInstance) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	currentTime := time.Now()
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(&User{LastFailedLogin: &currentTime}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserNextAllowedLoginTime updates the user's next allowed login time. This field is used to check whether we can
// allow a user to log in immediately or wait for some time before retrying the login process.
func (db *PGInstance) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(&User{NextAllowedLogin: &nextAllowedLoginTime}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserLastSuccessfulLoginTime updates the `lastSuccessfulLogin` field in the event where a user
// successfully logs into the app
func (db *PGInstance) UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error {
	currentTime := time.Now()
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(&User{LastSuccessfulLogin: &currentTime}).Error
	if err != nil {
		return err
	}
	return nil
}

// SetNickName is used to set the user's nickname in the database
func (db *PGInstance) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	if userID == nil || nickname == nil {
		return false, fmt.Errorf("userID or nickname cannot be nil")
	}
	err := db.DB.Model(&User{}).Where(&User{UserID: userID}).Updates(&User{Username: *nickname}).Error
	if err != nil {
		return false, fmt.Errorf("failed to set nickname")
	}

	return true, nil
}

// UpdateUserPinChangeRequiredStatus updates the user's pin change required from true to false. It'll be used to
// determine the onboarding journey for a user.
func (db *PGInstance) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID, Flavour: flavour}).Updates(map[string]interface{}{
		"pin_change_required":        false,
		"has_set_pin":                true,
		"has_set_security_questions": true,
		"is_phone_verified":          true,
	}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// InvalidatePIN toggles the valid field of a pin from true to false
func (db *PGInstance) InvalidatePIN(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")

	}
	err := db.DB.Model(&PINData{}).Where(&PINData{UserID: userID, IsValid: true}).Select("active").Updates(PINData{IsValid: false}).Error
	if err != nil {
		return false, fmt.Errorf("an error occurred while invalidating the pin: %v", err)
	}
	return true, nil
}

// UpdateIsCorrectSecurityQuestionResponse updates the is_correct_security_question_response field in the database
func (db *PGInstance) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")

	}
	err := db.DB.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: userID}).Updates(map[string]interface{}{
		"is_correct": isCorrectSecurityQuestionResponse,
	}).Error
	if err != nil {
		return false, fmt.Errorf("an error occurred while updating the is correct security question response: %v", err)
	}
	return true, nil
}

// ShareContent  updates the user shared content count
func (db *PGInstance) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	if input.ContentID == 0 || input.UserID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}

	var (
		contentItem  ContentItem
		contentShare *ContentShare
	)

	contentShare = &ContentShare{
		Active:    true,
		ContentID: input.ContentID,
		UserID:    input.UserID,
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentShare.ContentID}).First(&contentItem).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	updatedShareCount := contentItem.ShareCount + 1

	err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentShare.ContentID}).Updates(map[string]interface{}{
		"share_count": updatedShareCount,
	}).Error
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update content share count: %v", err)
	}

	err = tx.Where(contentShare).FirstOrCreate(contentShare).Error
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to save or get share content: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to share count failed: %v", err)
	}

	return true, nil
}

//BookmarkContent enable a user to set a bookmark on a content
func (db *PGInstance) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	var (
		contentItem     ContentItem
		contentBookmark *ContentBookmark
	)

	contentBookmark = &ContentBookmark{
		UserID:    userID,
		ContentID: contentID,
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).First(&contentItem).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	var updatedBookmarkCount = contentItem.BookmarkCount
	if err := tx.Model(&ContentBookmark{}).Where(&ContentBookmark{UserID: userID, ContentID: contentID}).First(&contentBookmark).Error; err != nil {
		updatedBookmarkCount++
	}

	err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).Updates(map[string]interface{}{
		"bookmark_count": updatedBookmarkCount,
	}).Error
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update content share count: %v", err)
	}

	err = tx.Where(contentBookmark).FirstOrCreate(contentBookmark).Error
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to save or get share content: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to share count failed: %v", err)
	}

	return true, nil
}

// UnBookmarkContent removes a bookmark from a content
func (db *PGInstance) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	var (
		contentItem     ContentItem
		contentBookmark *ContentBookmark
	)

	contentBookmark = &ContentBookmark{
		UserID:    userID,
		ContentID: contentID,
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).First(&contentItem).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	var updatedBookmarkCount = contentItem.BookmarkCount

	if err := tx.Model(&ContentBookmark{}).Where(&ContentBookmark{UserID: userID, ContentID: contentID}).First(&contentBookmark).Error; err == nil {
		updatedBookmarkCount--

		err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).Updates(map[string]interface{}{
			"bookmark_count": updatedBookmarkCount,
		}).Error
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to update content share count: %v", err)
		}

		err = tx.Delete(contentBookmark).Error
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to delete content bookmark: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to share count failed: %v", err)
	}
	return true, nil
}
