package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
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
	InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
	LikeContent(context context.Context, userID string, contentID int) (bool, error)
	UnlikeContent(context context.Context, userID string, contentID int) (bool, error)
	ViewContent(ctx context.Context, userID string, contentID int) (bool, error)
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
	UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error)
	AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error
	UpdateServiceRequests(ctx context.Context, payload []*ClientServiceRequest) (bool, error)
}

// LikeContent perfoms the actual database operation to update content like. The operation
// is carried out in a transaction.
func (db *PGInstance) LikeContent(context context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	contentLike := &ContentLike{
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
		return false, fmt.Errorf("failed to initialize like transaction")
	}

	var contentItem ContentItem
	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).First(&contentItem).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to get content item: %v", err)
	}

	var likeCount = contentItem.LikeCount
	err := tx.Model(&ContentLike{}).Create(&ContentLike{ContentID: contentID, UserID: userID, Active: true}).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return true, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to update content likes: %v", err)
	}

	err = tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).
		Updates(map[string]interface{}{
			"like_count": likeCount + 1,
		}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to update like count in content item table: %v", err)
	}

	err = tx.Where(contentLike).FirstOrCreate(contentLike).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to save or get like content: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to like count failed: %v", err)
	}

	return true, nil
}

// UnlikeContent perfoms the actual database operation to update content unlike. The operation
// is carried out in a transaction.
func (db *PGInstance) UnlikeContent(context context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to initialize like transaction")
	}

	var contentLike ContentLike
	if err := tx.Model(&ContentLike{}).Where(&ContentLike{UserID: userID, ContentID: contentID}).First(&contentLike).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return true, nil
		}
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to get content like for the specifies user: %v", err)
	}

	var contentItem ContentItem
	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).First(&contentItem).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to get content item: %v", err)
	}

	err := tx.Model(&ContentLike{}).Where(&ContentLike{UserID: userID, ContentID: contentID}).Delete(&ContentLike{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to delete content likes: %v", err)
	}

	err = tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).
		Updates(map[string]interface{}{
			"like_count": contentItem.LikeCount - 1,
		}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to update like count in content item table: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to like count failed: %v", err)
	}

	return true, nil
}

// ReactivateFacility perfoms the actual re-activation of the facility in the database
func (db *PGInstance) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: false}).
		Updates(&Facility{Active: true}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		return err
	}
	return nil
}

// UpdateUserLastFailedLoginTime updates the user's last failed login time
func (db *PGInstance) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	currentTime := time.Now()
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(&User{LastFailedLogin: &currentTime}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}
	return nil
}

// UpdateUserNextAllowedLoginTime updates the user's next allowed login time. This field is used to check whether we can
// allow a user to log in immediately or wait for some time before retrying the login process.
func (db *PGInstance) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID}).Updates(&User{NextAllowedLogin: &nextAllowedLoginTime}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to set nickname")
	}

	return true, nil
}

// SetInProgressBy updates the staff assigned to a service request
func (db *PGInstance) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	if requestID == "" || staffID == "" {
		return false, fmt.Errorf("requestID or staffID cannot be empty")
	}
	if err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: &requestID}).Updates(map[string]interface{}{
		"status":            enums.ServiceRequestStatusInProgress,
		"in_progress_by_id": staffID,
		"in_progress_at":    time.Now(),
	}).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update the assi")
	}
	return true, nil
}

// UpdateUserPinChangeRequiredStatus updates the user's pin change required from true to false. It'll be used to
// determine the onboarding journey for a user.
func (db *PGInstance) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if !flavour.IsValid() {
		return false, fmt.Errorf("invalid flavour provided")
	}
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID, Flavour: flavour}).Updates(map[string]interface{}{
		"pin_change_required":        false,
		"has_set_pin":                true,
		"has_set_security_questions": true,
		"is_phone_verified":          true,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	return true, nil
}

// InvalidatePIN toggles the valid field of a pin from true to false
func (db *PGInstance) InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")

	}
	err := db.DB.Model(&PINData{}).Where(&PINData{UserID: userID, IsValid: true, Flavour: flavour}).Select("active").Updates(PINData{IsValid: false}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	updatedShareCount := contentItem.ShareCount + 1

	err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentShare.ContentID}).Updates(map[string]interface{}{
		"share_count": updatedShareCount,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to update content share count: %v", err)
	}

	err = tx.Where(contentShare).FirstOrCreate(contentShare).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	var updatedBookmarkCount = contentItem.BookmarkCount
	if err := tx.Model(&ContentBookmark{}).Where(&ContentBookmark{UserID: userID, ContentID: contentID}).First(&contentBookmark).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		updatedBookmarkCount++
	}

	err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentID}).Updates(map[string]interface{}{
		"bookmark_count": updatedBookmarkCount,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to update content share count: %v", err)
	}

	err = tx.Where(contentBookmark).FirstOrCreate(contentBookmark).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
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
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to update content share count: %v", err)
		}

		err = tx.Delete(contentBookmark).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
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

// ViewContent gets a specific content and updates the count each time it is viewed
func (db *PGInstance) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	var (
		contentItem ContentItem
		contentView *ContentView
	)

	contentView = &ContentView{
		Active:    true,
		ContentID: contentID,
		UserID:    userID,
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

	if err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentView.ContentID}).First(&contentItem).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get content item: %v", err)
	}

	updatedViewCount := contentItem.ViewCount + 1

	err := tx.Model(&ContentItem{}).Where(&ContentItem{PagePtrID: contentView.ContentID}).Updates(map[string]interface{}{
		"view_count": updatedViewCount,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update content view count: %v", err)
	}

	err = tx.Where(contentView).FirstOrCreate(contentView).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to save view content count: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to view count/content failed: %v", err)
	}

	return true, nil
}

// UpdateClientCaregiver updates the caregiver for a client
func (db *PGInstance) UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	var (
		client Client
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failied initialize database transaction %v", err)
	}

	if err := tx.Model(&Client{}).Where(&Client{ID: &caregiverInput.ClientID}).First(&client).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get client: %v", err)
	}

	err := tx.Model(&Caregiver{}).Where(&Caregiver{CaregiverID: client.CaregiverID}).Updates(map[string]interface{}{
		"first_name":     caregiverInput.FirstName,
		"last_name":      caregiverInput.LastName,
		"phone_number":   caregiverInput.PhoneNumber,
		"caregiver_type": caregiverInput.CaregiverType,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update caregiver: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to update client caregiver failed: %v", err)
	}

	return nil
}

// ResolveServiceRequest resolves a service request for a given client
func (db *PGInstance) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
	var (
		serviceRequest ClientServiceRequest
	)

	currentTime := time.Now()

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	if err := tx.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: serviceRequestID}).First(&serviceRequest).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get service request: %v", err)
	}

	err := tx.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: serviceRequestID}).Updates(ClientServiceRequest{
		Status:       enums.ServiceRequestStatusResolved.String(),
		ResolvedByID: staffID,
		ResolvedAt:   &currentTime,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to update service request: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update service request failed: %v", err)
	}

	return true, nil
}

// AssignRoles assigns roles to a user
func (db *PGInstance) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	var (
		user User
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	err := tx.Model(&User{}).Where(&User{UserID: &userID}).First(&user).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	for _, role := range roles {
		var (
			roleID string
		)

		err := tx.Raw(`SELECT id FROM authority_authorityrole WHERE name = ?`, role.String()).Row().Scan(&roleID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to get authority role: %v", err)
		}

		err = tx.Model(&AuthorityRoleUser{}).Where(&AuthorityRoleUser{UserID: user.UserID, RoleID: &roleID}).FirstOrCreate(&AuthorityRoleUser{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to assign role: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update user roles failed: %v", err)
	}

	return true, nil
}

// RevokeRoles revokes roles from a user
func (db *PGInstance) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	var user User

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	err := tx.Model(&User{}).Where(&User{UserID: &userID}).First(&user).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	for _, role := range roles {
		var roleID string

		err := tx.Raw(`SELECT id FROM authority_authorityrole WHERE name = ?`, role.String()).Row().Scan(&roleID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to get authority role: %v", err)
		}

		err = tx.Model(&AuthorityRoleUser{}).Where(&AuthorityRoleUser{UserID: user.UserID, RoleID: &roleID}).Delete(&AuthorityRoleUser{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to revoke role: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update user roles failed: %v", err)
	}

	return true, nil
}

// InvalidateScreeningToolResponse invalidates a screening tool response
func (db *PGInstance) InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error {
	var (
		client               Client
		screeningToolRespose ScreeningToolsResponse
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failied initialize database transaction %v", err)
	}

	err := tx.Model(&Client{}).Where(&Client{ID: &clientID}).First(&client).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to get client: %v", err)
	}

	err = tx.Model(&screeningToolRespose).Where(
		&ScreeningToolsResponse{
			ClientID:   clientID,
			QuestionID: questionID,
		}).Updates(map[string]interface{}{"active": false}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to invalidate screening tool response: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("transaction commit to update screening tool response failed: %v", err)
	}
	return nil
}

// UpdateServiceRequests performs and update to the client service requests
func (db *PGInstance) UpdateServiceRequests(ctx context.Context, payload []*ClientServiceRequest) (bool, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to initialize database transaction: %v", err)
	}

	for _, k := range payload {
		err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: k.ID, RequestType: k.RequestType}).Updates(map[string]interface{}{
			"status":         k.Status,
			"in_progress_at": k.InProgressAt,
			"resolved_at":    k.ResolvedAt,
		}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("unable to update client's service request: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to commit transaction: %v", err)
	}

	return true, nil
}
