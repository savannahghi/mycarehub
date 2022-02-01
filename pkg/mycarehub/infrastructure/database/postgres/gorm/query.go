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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*User, error)
	GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*TermsOfService, error)
	CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error)
	GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*Client, error)
	GetStaffProfileByUserID(ctx context.Context, userID string) (*StaffProfile, error)
	CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*UserOTP, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*Contact, error)
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
	GetUserBookmarkedContent(ctx context.Context, userID string) ([]*ContentItem, error)
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
	GetClientHealthDiaryQuote(ctx context.Context) (*ClientHealthDiaryQuote, error)
	CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error)
	GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*ClientHealthDiaryEntry, error)
	GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*FAQ, error)
	GetClientCaregiver(ctx context.Context, caregiverID string) (*Caregiver, error)
	GetClientByClientID(ctx context.Context, clientID string) (*Client, error)
	SearchUser(ctx context.Context, CCCNumber string) (*Client, error)
	GetServiceRequests(ctx context.Context, requestType *string, requestStatus *string) ([]*ClientServiceRequest, error)
}

// CheckWhetherUserHasLikedContent performs a operation to check whether user has liked the content
func (db *PGInstance) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	var contentItemLike ContentLike
	if err := db.DB.Where(&ContentLike{UserID: userID, ContentID: contentID}).First(&contentItemLike).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred: %v", err)
	}

	return true, nil
}

//ListContentCategories performs the actual database query to get the list of content categories
func (db *PGInstance) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	var contentItemCategories []*ContentItemCategory

	err := db.DB.Find(&contentItemCategories).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to query all content categories %v", err)
	}

	var domainContentItemCategory []*domain.ContentItemCategory
	for _, contentCategory := range contentItemCategories {
		var wagtailImage *WagtailImages
		err := db.DB.Model(&WagtailImages{}).Where(&WagtailImages{ID: contentCategory.IconID}).Find(&wagtailImage).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to fetch wagtail images %v", err)
		}
		contentItemCategory := &domain.ContentItemCategory{
			ID:      contentCategory.ID,
			Name:    contentCategory.Name,
			IconURL: wagtailImage.File,
		}
		domainContentItemCategory = append(domainContentItemCategory, contentItemCategory)
	}

	return domainContentItemCategory, nil
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility id cannot be nil")
	}
	var facility Facility
	err := db.DB.Where(&Facility{FacilityID: id, Active: isActive}).First(&facility).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}

// CheckIfPhoneNumberExists checks if phone exists in the database.
func (db *PGInstance) CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
	var contact Contact
	if phone == "" || !flavour.IsValid() {
		return false, fmt.Errorf("invalid flavour: %v", flavour)
	}
	err := db.DB.Model(&Contact{}).Where(&Contact{ContactValue: phone, OptedIn: isOptedIn, Flavour: flavour}).First(&contact).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check if contact exists %v: %v", phone, err)
	}
	return true, nil
}

// RetrieveFacilityByMFLCode fetches a single facility using MFL Code
func (db *PGInstance) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*Facility, error) {
	if MFLCode == 0 {
		return nil, fmt.Errorf("facility mfl code cannot be nil")
	}
	var facility Facility
	if err := db.DB.Where(&Facility{Code: MFLCode, Active: isActive}).First(&facility).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get facility by MFL Code %v and status %v: %v", MFLCode, isActive, err)
	}
	return &facility, nil
}

// GetFacilities fetches all the healthcare facilities in the platform.
func (db *PGInstance) GetFacilities(ctx context.Context) ([]Facility, error) {
	var facility []Facility
	err := db.DB.Find(&facility).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to query all facilities %v", err)
	}
	return facility, nil
}

// GetSecurityQuestions fetches all the security questions.
func (db *PGInstance) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error) {
	if flavour == "" || !flavour.IsValid() {
		return nil, fmt.Errorf("bad flavor specified: %v", flavour)
	}
	var securityQuestion []*SecurityQuestion
	err := db.DB.Where(&SecurityQuestion{Flavour: flavour}).Find(&securityQuestion).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to query all security questions %v", err)
	}
	return securityQuestion, nil
}

// ListFacilities lists all facilities, the results returned are
// from search, and provided filters. they are also paginated
func (db *PGInstance) ListFacilities(
	ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
	var facilities []Facility
	// this will keep track of the results for pagination
	// Count query is unreliable for this since it is returning the count for all rows instead of results
	var resultCount int64

	facilitiesOutput := []domain.Facility{}

	for _, f := range filter {
		err := f.Validate()
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to validate filter %v: %v", f.Value, err)
		}
		err = enums.ValidateFilterSortCategories(enums.FilterSortCategoryTypeFacility, f.DataType)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("filter param %v is not available in facilities: %v", f.Value, err)
		}
	}

	paginatedFacilities := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:        pagination.Pagination.Limit,
			CurrentPage:  pagination.Pagination.CurrentPage,
			Count:        pagination.Pagination.Count,
			TotalPages:   pagination.Pagination.TotalPages,
			NextPage:     pagination.Pagination.NextPage,
			PreviousPage: pagination.Pagination.PreviousPage,
			Sort:         pagination.Pagination.Sort,
		},
		Facilities: pagination.Facilities,
	}

	mappedFilterParams := filterParamsToMap(filter)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed to initialize filter facilities transaction %v", err)
	}

	tx.Where(
		"name ~* ? OR county ~* ? OR description ~* ?",
		*searchTerm, *searchTerm, *searchTerm,
	).Where(mappedFilterParams).Find(&facilities).Find(&facilities)

	resultCount = int64(len(facilities))

	tx.Scopes(
		paginate(facilities, &paginatedFacilities.Pagination, resultCount, db.DB),
	).Where(
		"name ~* ?  OR county ~* ? OR description ~* ?",
		*searchTerm, *searchTerm, *searchTerm,
	).Where(mappedFilterParams).Find(&facilities)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction list facilities transaction%v", err)
	}

	for _, f := range facilities {
		facility := domain.Facility{
			ID:          f.FacilityID,
			Name:        f.Name,
			Code:        f.Code,
			Active:      f.Active,
			County:      f.County,
			Description: f.Description,
		}
		facilitiesOutput = append(facilitiesOutput, facility)
	}

	pagination.Pagination.Count = paginatedFacilities.Pagination.Count
	pagination.Pagination.TotalPages = paginatedFacilities.Pagination.TotalPages
	pagination.Pagination.Limit = paginatedFacilities.Pagination.Limit
	pagination.Facilities = facilitiesOutput
	pagination.Pagination.NextPage = paginatedFacilities.Pagination.NextPage

	pagination.Pagination.PreviousPage = paginatedFacilities.Pagination.PreviousPage

	return pagination, nil
}

// GetUserProfileByPhoneNumber retrieves a user profile using their phonenumber
func (db *PGInstance) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	var user User
	if err := db.DB.Joins("JOIN common_contact on users_user.id = common_contact.user_id").Where("common_contact.contact_value = ?", phoneNumber).Preload(clause.Associations).First(&user).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user by phonenumber %v: %v", phoneNumber, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user's pin using the user ID
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, IsValid: true}).First(&pin).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get pin: %v", err)
	}
	return &pin, nil
}

// GetCurrentTerms fetches the most most recent terms of service depending on the flavour
func (db *PGInstance) GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*TermsOfService, error) {
	var termsOfService TermsOfService
	validTo := time.Now()
	if err := db.DB.Model(&TermsOfService{}).Where(db.DB.Where(&TermsOfService{Flavour: flavour}).Where("valid_to > ?", validTo).Or("valid_to = ?", nil).Order("valid_to desc")).First(&termsOfService).Statement.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get the current terms : %v", err)
	}

	return &termsOfService, nil
}

// GetUserProfileByUserID fetches a user profile using the user ID
func (db *PGInstance) GetUserProfileByUserID(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	var user User
	if err := db.DB.Where(&User{UserID: &userID}).First(&user).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user by user ID %v: %v", userID, err)
	}
	return &user, nil
}

// GetSecurityQuestionByID fetches a security question using the security question ID
func (db *PGInstance) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error) {
	var securityQuestion SecurityQuestion
	if err := db.DB.Where(&SecurityQuestion{SecurityQuestionID: securityQuestionID}).First(&securityQuestion).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get security question by ID %v: %v", securityQuestionID, err)
	}
	return &securityQuestion, nil
}

// GetSecurityQuestionResponseByID returns the security question response
func (db *PGInstance) GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*SecurityQuestionResponse, error) {
	var questionResponse SecurityQuestionResponse
	if err := db.DB.Where(&SecurityQuestionResponse{QuestionID: questionID}).First(&questionResponse).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get the security question response by ID")
	}
	return &questionResponse, nil
}

//VerifyOTP checks from the database for the validity of the provided OTP
func (db *PGInstance) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	var userOTP UserOTP
	if payload.PhoneNumber == "" || payload.OTP == "" {
		return false, fmt.Errorf("user ID or phone number or OTP cannot be empty")
	}
	if !payload.Flavour.IsValid() {
		return false, exceptions.InvalidFlavourDefinedError()
	}

	err := db.DB.Model(&UserOTP{}).Where(&UserOTP{PhoneNumber: payload.PhoneNumber, Valid: true, OTP: payload.OTP, Flavour: payload.Flavour}).First(&userOTP).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to verify otp with %v: %v", payload.OTP, payload.Flavour)
	}

	return true, nil
}

// GetClientProfileByUserID returns the client profile based on the user ID provided
func (db *PGInstance) GetClientProfileByUserID(ctx context.Context, userID string) (*Client, error) {
	var client Client
	if err := db.DB.Where(&Client{UserID: &userID}).Preload(clause.Associations).First(&client).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client by user ID %v: %v", userID, err)
	}
	return &client, nil
}

// GetStaffProfileByUserID returns the staff profile
func (db *PGInstance) GetStaffProfileByUserID(ctx context.Context, userID string) (*StaffProfile, error) {
	var staff StaffProfile

	if err := db.DB.Where(&StaffProfile{UserID: userID}).Preload(clause.Associations).First(&staff).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get staff by the provided user id %v", userID)
	}

	return &staff, nil
}

// CheckUserHasPin performs a look up on the pins table to check whether a user has a pin
func (db *PGInstance) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if !flavour.IsValid() {
		return false, fmt.Errorf("invalid flavour defined")
	}
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, Flavour: flavour}).Find(&pin).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	return true, nil
}

// GetOTP fetches an OTP from the database
func (db *PGInstance) GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*UserOTP, error) {
	var userOTP UserOTP
	if err := db.DB.Where(&UserOTP{PhoneNumber: phoneNumber, Flavour: flavour}).First(&userOTP).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get otp: %v", err)
	}
	return &userOTP, nil
}

// GetUserSecurityQuestionsResponses fetches the security question responses that the user has responded to
func (db *PGInstance) GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*SecurityQuestionResponse, error) {
	var securityQuestionResponses []*SecurityQuestionResponse
	if err := db.DB.Where(&SecurityQuestionResponse{UserID: userID, Active: true}).Find(&securityQuestionResponses).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get security questions: %v", err)
	}
	return securityQuestionResponses, nil
}

// GetContactByUserID fetches a user's contact using the user ID
func (db *PGInstance) GetContactByUserID(ctx context.Context, userID *string, contactType string) (*Contact, error) {
	var contact Contact

	if userID == nil {
		return nil, fmt.Errorf("user ID is required")
	}
	if contactType == "" {
		return nil, fmt.Errorf("contact type is required")
	}

	if contactType != "PHONE" && contactType != "EMAIL" {
		return nil, fmt.Errorf("contact type must be PHONE or EMAIL")
	}
	if err := db.DB.Where(&Contact{UserID: userID, ContactType: contactType}).First(&contact).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get contact: %v", err)
	}
	return &contact, nil
}

// GetUserBookmarkedContent retrieves a user's pinned content from the database
func (db *PGInstance) GetUserBookmarkedContent(ctx context.Context, userID string) ([]*ContentItem, error) {
	var contentItem []*ContentItem
	err := db.DB.Joins("JOIN content_contentbookmark ON content_contentitem.page_ptr_id = content_contentbookmark.content_item_id").
		Where("content_contentbookmark.user_id = ?", userID).Preload(clause.Associations).Find(&contentItem).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	return contentItem, nil
}

// CanRecordHeathDiary checks whether a user can record a health diary
// if the last record is less than 24 hours ago, the user cannot record a new entry
// if the last record is more than 24 hours ago, the user can record a new entry
func (db *PGInstance) CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error) {
	var clientHealthDiaryEntry []*ClientHealthDiaryEntry
	err := db.DB.Where("client_id = ?", clientID).Order("created desc").Find(&clientHealthDiaryEntry).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get client health diary: %v", err)
	}
	if len(clientHealthDiaryEntry) > 0 {
		if time.Since(clientHealthDiaryEntry[0].CreatedAt) < time.Hour*24 {
			return false, nil
		}
	}

	return true, nil
}

// GetClientHealthDiaryQuote fetches a client's health diary quote.
// it should be a random quote from the health diary
func (db *PGInstance) GetClientHealthDiaryQuote(ctx context.Context) (*ClientHealthDiaryQuote, error) {
	var healthDiaryQuote ClientHealthDiaryQuote
	err := db.DB.Where("active = true").Order("RANDOM()").First(&healthDiaryQuote).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	return &healthDiaryQuote, nil
}

// CheckIfUserBookmarkedContent fetches a user's pinned content from the database
func (db *PGInstance) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	var contentBookmark ContentBookmark
	err := db.DB.Where(&ContentBookmark{ContentID: contentID, UserID: userID}).First(&contentBookmark).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get content bookmark: %v", err)
	}
	return true, nil
}

// GetClientHealthDiaryEntries gets all health diary entries that belong to a specific client
func (db *PGInstance) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry []*ClientHealthDiaryEntry
	err := db.DB.Where(&ClientHealthDiaryEntry{ClientID: clientID, Active: true}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&healthDiaryEntry).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get all client health diary entries: %v", err)
	}
	return healthDiaryEntry, nil
}

// GetFAQContent fetches the FAQ content from the database
// when the limit is not provided, it defaults to 10
func (db *PGInstance) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*FAQ, error) {
	var faq []*FAQ
	err := db.DB.Where(&FAQ{Flavour: flavour, Active: true}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Limit(*limit).Find(&faq).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get FAQ content: %v", err)
	}
	return faq, nil

}

// GetClientCaregiver fetches a client's caregiver from the database
func (db *PGInstance) GetClientCaregiver(ctx context.Context, caregiverID string) (*Caregiver, error) {
	var (
		caregiver Caregiver
	)

	err := db.DB.Where(&Caregiver{CaregiverID: &caregiverID}).First(&caregiver).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get caregiver: %v", err)
	}
	return &caregiver, nil
}

// GetClientByClientID fetches a client from the database
func (db *PGInstance) GetClientByClientID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	err := db.DB.Where(&Client{ID: &clientID}).First(&client).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	return &client, nil
}

// GetServiceRequests fetches clients service requests from the database according to the type and or status passed
func (db *PGInstance) GetServiceRequests(ctx context.Context, requestType *string, requestStatus *string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	if requestType != nil && requestStatus == nil {
		err := db.DB.Where(&ClientServiceRequest{RequestType: *requestType}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else if requestType == nil && requestStatus != nil {
		err := db.DB.Where(&ClientServiceRequest{Status: *requestStatus}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else if requestType != nil && requestStatus != nil {
		err := db.DB.Where(&ClientServiceRequest{RequestType: *requestType, Status: *requestStatus}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else {
		err := db.DB.Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	}

	return serviceRequests, nil
}

// SearchUser searches for a patient using their CCC Number
func (db *PGInstance) SearchUser(ctx context.Context, CCCNumber string) (*Client, error) {
	var user *User

	err := db.DB.Joins("JOIN common_contact on users_user.id = common_contact.user_id "+
		"JOIN clients_client on users_user.id = clients_client.user_id "+
		"JOIN clients_client_identifiers "+
		"ON clients_client.id = clients_client_identifiers.client_id "+
		"JOIN clients_identifier "+
		"ON clients_client_identifiers.identifier_id = clients_identifier.id").
		Where("clients_identifier.identifier_value = ?", CCCNumber).Preload(clause.Associations).
		First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("an error occurred: %v", err)
	}
	logrus.Print("THE USER PROFILE IS: ", user)
	logrus.Print("THE CLIENT USER PROFILE FN IS: ", user.FirstName)
	return &Client{
		UserProfile: *user,
	}, nil
}
