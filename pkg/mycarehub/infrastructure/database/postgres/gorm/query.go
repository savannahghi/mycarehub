package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	ListAppointments(ctx context.Context, params *Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*Appointment, *domain.Pagination, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*User, error)
	GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*PINData, error)
	GetUserProfileByUserID(ctx context.Context, userID *string) (*User, error)
	GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*TermsOfService, error)
	CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error)
	GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*Client, error)
	GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*Client, error)
	GetStaffProfileByUserID(ctx context.Context, userID string) (*StaffProfile, error)
	CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*UserOTP, error)
	GetClientsPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
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
	GetClientProfileByClientID(ctx context.Context, clientID string) (*Client, error)
	GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*ClientServiceRequest, error)
	GetStaffServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*StaffServiceRequest, error)
	CheckUserRole(ctx context.Context, userID string, role string) (bool, error)
	CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error)
	GetUserRoles(ctx context.Context, userID string) ([]*AuthorityRole, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*AuthorityPermission, error)
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)
	GetCommunityByID(ctx context.Context, communityID string) (*Community, error)
	CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error)
	GetClientsInAFacility(ctx context.Context, facilityID string) ([]*Client, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*ClientHealthDiaryEntry, error)
	GetClientsByParams(ctx context.Context, query Client, lastSyncTime *time.Time) ([]*Client, error)
	GetClientCCCIdentifier(ctx context.Context, clientID string) (*Identifier, error)
	SearchClientProfilesByCCCNumber(ctx context.Context, CCCNumber string) ([]*Client, error)
	SearchStaffProfileByStaffNumber(ctx context.Context, staffNumber string) ([]*StaffProfile, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*ClientServiceRequest, error)
	GetScreeningToolQuestions(ctx context.Context, toolType string) ([]ScreeningToolQuestion, error)
	GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*ScreeningToolQuestion, error)
	CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	GetAllRoles(ctx context.Context) ([]*AuthorityRole, error)
	GetUserProfileByStaffID(ctx context.Context, staffID string) (*User, error)
	GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*ClientHealthDiaryEntry, error)
	GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*ClientServiceRequest, error)
	GetStaffProfileByStaffID(ctx context.Context, staffID string) (*StaffProfile, error)
	GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time) ([]*ClientServiceRequest, error)
	GetAppointmentByID(ctx context.Context, appointmentID string) (*Appointment, error)
	GetAppointmentByAppointmentUUID(ctx context.Context, appointmentUUID string) (*Appointment, error)
	GetClientAppointmentByID(ctx context.Context, appointmentID, clientID string) (*Appointment, error)
	GetClientServiceRequests(ctx context.Context, requestType, status, clientID string) ([]*ClientServiceRequest, error)
	GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*ScreeningToolsResponse, error)
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

// GetFacilitiesWithoutFHIRID fetches all the healthcare facilities in the platform without FHIR Organisation ID
func (db *PGInstance) GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*Facility, error) {
	var facility []*Facility
	err := db.DB.Raw(
		`SELECT * FROM common_facility WHERE fhir_organization_id IS NULL`).Scan(&facility).Error
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

// ListAppointments Retrieves appointments using the provided parameters and filters
func (db *PGInstance) ListAppointments(ctx context.Context, params *Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*Appointment, *domain.Pagination, error) {
	var appointments []*Appointment
	pageInfo := &domain.Pagination{} // TODO: fix pagination implementation
	// this will keep track of the results for pagination
	// Count query is unreliable for this since it is returning the count for all rows instead of results
	var resultCount int64

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, pageInfo, fmt.Errorf("failed to initialize filter appointments transaction %v", err)
	}

	transaction := tx.Where(params)

	transaction, err := addFilters(transaction, filters)
	if err != nil {
		return nil, pageInfo, fmt.Errorf("failed to add filters: %v", err)
	}

	transaction.Find(&appointments)

	resultCount = int64(len(appointments))

	if pagination != nil {
		transaction = tx.Scopes(
			paginate(appointments, pagination, resultCount, db.DB),
		).Where(params)

		transaction, err := addFilters(transaction, filters)
		if err != nil {
			return nil, pageInfo, fmt.Errorf("failed to add filters: %v", err)
		}

		transaction.Find(&appointments)

		pageInfo = pagination
	}

	if err := tx.Commit().Error; err != nil {
		return nil, pageInfo, fmt.Errorf("failed to commit transaction list facilities transaction%v", err)
	}

	return appointments, pageInfo, nil
}

// GetUserProfileByPhoneNumber retrieves a user profile using their phonenumber
func (db *PGInstance) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*User, error) {
	var user User
	if err := db.DB.Joins("JOIN common_contact on users_user.id = common_contact.user_id").Where("common_contact.contact_value = ? AND common_contact.flavour = ?", phoneNumber, flavour).
		Preload(clause.Associations).First(&user).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user by phonenumber %v: %v", phoneNumber, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user's pin using the user ID and Flavour
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*PINData, error) {
	if !flavour.IsValid() {
		return nil, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, IsValid: true, Flavour: flavour}).First(&pin).Error; err != nil {
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
func (db *PGInstance) GetUserProfileByUserID(ctx context.Context, userID *string) (*User, error) {
	if userID == nil {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	var user User
	if err := db.DB.Where(&User{UserID: userID}).Preload(clause.Associations).First(&user).Error; err != nil {
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

// GetSecurityQuestionResponse returns the security question response
func (db *PGInstance) GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*SecurityQuestionResponse, error) {
	var questionResponse SecurityQuestionResponse
	if err := db.DB.Where(&SecurityQuestionResponse{QuestionID: questionID, UserID: userID}).First(&questionResponse).Error; err != nil {
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
		return false, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
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

// SearchStaffProfileByStaffNumber searches for the staff profile(s) of a given staff.
// It may also return other staffs whose staff number may match at a given time.
func (db *PGInstance) SearchStaffProfileByStaffNumber(ctx context.Context, staffNumber string) ([]*StaffProfile, error) {
	var staff []*StaffProfile

	if err := db.DB.Where("staff_number ~~* ? ", "%"+staffNumber+"%").Preload(clause.Associations).Find(&staff).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to ge staff %v", err)
	}

	return staff, nil
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

// GetServiceRequestsForKenyaEMR gets all the service requests to be used by the KenyaEMR.
func (db *PGInstance) GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest

	err := db.DB.Where(&ClientServiceRequest{FacilityID: facilityID}).Where("created > ?", lastSyncTime).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&serviceRequests).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get service requests: %v", err)
	}
	return serviceRequests, nil
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

// GetStaffPendingServiceRequestsCount gets the number of staffs pending pin reser service requests
func (db *PGInstance) GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	var staffServiceRequest []*StaffServiceRequest

	err := db.DB.Model(&StaffServiceRequest{}).Where(&StaffServiceRequest{DefaultFacilityID: &facilityID, RequestType: "STAFF_PIN_RESET", Status: "PENDING"}).Find(&staffServiceRequest).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	staffServiceRequestCount := &domain.ServiceRequestsCount{
		Total: len(staffServiceRequest),
		RequestsTypeCount: []*domain.RequestTypeCount{
			{
				RequestType: enums.ServiceRequestTypeStaffPinReset,
				Total:       0,
			},
		},
	}

	for _, req := range staffServiceRequest {
		if req.RequestType == enums.ServiceRequestTypeStaffPinReset.String() {
			staffServiceRequestCount.RequestsTypeCount[0].Total++
		}
	}

	return staffServiceRequestCount, nil
}

// GetClientsPendingServiceRequestsCount gets the number of clients service requests
func (db *PGInstance) GetClientsPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	var serviceRequests []*ClientServiceRequest

	err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{FacilityID: facilityID, Status: "PENDING"}).Find(&serviceRequests).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client's service requests:: %v", err)
	}

	serviceRequestsCount := domain.ServiceRequestsCount{
		Total: len(serviceRequests),
		RequestsTypeCount: []*domain.RequestTypeCount{
			{
				RequestType: enums.ServiceRequestTypeRedFlag,
				Total:       0,
			},
			{
				RequestType: enums.ServiceRequestTypePinReset,
				Total:       0,
			},
			{
				RequestType: enums.ServiceRequestTypeProfileUpdate,
				Total:       0,
			},
		},
	}

	for _, request := range serviceRequests {
		if request.RequestType == enums.ServiceRequestTypeRedFlag.String() {
			serviceRequestsCount.RequestsTypeCount[0].Total++
		}
		if request.RequestType == enums.ServiceRequestTypePinReset.String() {
			serviceRequestsCount.RequestsTypeCount[1].Total++
		}
		if request.RequestType == enums.ServiceRequestTypeProfileUpdate.String() {
			serviceRequestsCount.RequestsTypeCount[2].Total++
		}
	}

	return &serviceRequestsCount, nil
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

// GetClientProfileByClientID fetches a client from the database
func (db *PGInstance) GetClientProfileByClientID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	err := db.DB.Where(&Client{ID: &clientID}).Preload(clause.Associations).First(&client).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	return &client, nil
}

// GetStaffProfileByStaffID fetches a staff from the database
func (db *PGInstance) GetStaffProfileByStaffID(ctx context.Context, staffID string) (*StaffProfile, error) {
	var staff StaffProfile
	err := db.DB.Where(&StaffProfile{ID: &staffID}).Preload(clause.Associations).First(&staff).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %v", err)
	}
	return &staff, nil
}

// GetServiceRequests fetches clients service requests from the database according to the type and or status passed
func (db *PGInstance) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	if requestType != nil && requestStatus == nil {
		err := db.DB.Where(&ClientServiceRequest{RequestType: *requestType, FacilityID: facilityID}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else if requestType == nil && requestStatus != nil {
		err := db.DB.Where(&ClientServiceRequest{Status: *requestStatus, FacilityID: facilityID}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else if requestType != nil && requestStatus != nil {
		err := db.DB.Where(&ClientServiceRequest{RequestType: *requestType, Status: *requestStatus, FacilityID: facilityID}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	} else {
		err := db.DB.Where(&ClientServiceRequest{FacilityID: facilityID}).Find(&serviceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	}

	return serviceRequests, nil
}

// GetStaffServiceRequests gets all the staff's service requests depending on the provided parameters
func (db *PGInstance) GetStaffServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*StaffServiceRequest, error) {
	var staffServiceRequests []*StaffServiceRequest
	if requestType != nil && requestStatus != nil {
		err := db.DB.Where(&StaffServiceRequest{RequestType: *requestType, Status: *requestStatus, DefaultFacilityID: &facilityID}).Find(&staffServiceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else if requestType == nil && requestStatus != nil {
		err := db.DB.Where(&StaffServiceRequest{Status: *requestStatus, DefaultFacilityID: &facilityID}).Find(&staffServiceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else if requestType != nil && requestStatus == nil {
		err := db.DB.Where(&StaffServiceRequest{RequestType: *requestType, DefaultFacilityID: &facilityID}).Find(&staffServiceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else {
		err := db.DB.Where(&StaffServiceRequest{DefaultFacilityID: &facilityID}).Find(&staffServiceRequests).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	}

	return staffServiceRequests, nil
}

// CheckUserRole checks if a user has a specific role
func (db *PGInstance) CheckUserRole(ctx context.Context, userID string, role string) (bool, error) {
	var returneduserID string
	err := db.DB.Raw(
		`
		SELECT user_id 
		FROM authority_authorityrole_users 
		WHERE user_id = ? 
		AND authority_authorityrole_users.authorityrole_id = 
		(SELECT id FROM authority_authorityrole WHERE name = ?)
		`, userID, role,
	).Scan(&returneduserID).Error

	if returneduserID == "" {
		return false, nil
	}

	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check if user has role: %v", err)
	}

	return true, nil
}

// CheckUserPermission checks if a user has a specific permission
func (db *PGInstance) CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error) {
	var returneduserID string
	err := db.DB.Raw(
		`
		SELECT user_id 
		FROM authority_authorityrole_users 
		WHERE authorityrole_id =
		(
			SELECT 	authorityrole_id 
			FROM authority_authorityrole_permissions
			WHERE authoritypermission_id = (SELECT id FROM authority_authoritypermission WHERE name = ?)
		)
		And 
		user_id = ?
		`, permission, userID,
	).Scan(&returneduserID).Error

	if returneduserID == "" {
		return false, nil
	}

	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check if user has permission: %v", err)
	}

	return true, nil
}

// GetUserRoles fetches a user's roles from the database
func (db *PGInstance) GetUserRoles(ctx context.Context, userID string) ([]*AuthorityRole, error) {
	var roles []*AuthorityRole
	err := db.DB.Raw(
		`
		SELECT * 
		FROM authority_authorityrole_users 
		JOIN authority_authorityrole ON authority_authorityrole_users.authorityrole_id = authority_authorityrole.id
		WHERE user_id = ?
		`, userID,
	).Find(&roles).Error

	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	return roles, nil
}

// GetUserPermissions fetches a user's permissions from the database
func (db *PGInstance) GetUserPermissions(ctx context.Context, userID string) ([]*AuthorityPermission, error) {
	var permissions []*AuthorityPermission
	err := db.DB.Raw(
		`
		SELECT * 
		FROM authority_authorityrole_users 
		JOIN authority_authorityrole_permissions ON authority_authorityrole_users.authorityrole_id = authority_authorityrole_permissions.authorityrole_id
		JOIN authority_authoritypermission ON authority_authorityrole_permissions.authoritypermission_id = authority_authoritypermission.id
		WHERE user_id = ?
		`, userID,
	).Find(&permissions).Error

	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user permissions: %v", err)
	}

	return permissions, nil
}

// CheckIfUsernameExists checks to see whether the provided username exists
func (db *PGInstance) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	var user User
	err := db.DB.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if username exists: %v", err)
	}

	return true, nil
}

// GetCommunityByID fetches the community using its ID
func (db *PGInstance) GetCommunityByID(ctx context.Context, communityID string) (*Community, error) {
	var community *Community

	err := db.DB.Where(&Community{ID: communityID}).First(&community).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to find community by ID %s", communityID)
	}

	return community, nil
}

// CheckIdentifierExists checks whether an identifier of a certain type and value exists
// Used to validate uniqueness and prevent duplicates
func (db *PGInstance) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	var identifier *Identifier

	err := db.DB.Where(&Identifier{IdentifierType: identifierType, IdentifierValue: identifierValue, Active: true}).First(&identifier).Error
	if err != nil {
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check identifier of type: %s and value: %s", identifierType, identifierValue)
	}

	return true, nil
}

// CheckFacilityExistsByMFLCode checks whether a facility exists using the mfl code.
// Used to validate existence of a facility
func (db *PGInstance) CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error) {
	_, err := db.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return false, nil
		}
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check for facility: %s", err)
	}

	return true, nil
}

// GetClientsInAFacility returns all the clients registered within a specified facility
func (db *PGInstance) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*Client, error) {
	var clientProfiles []*Client
	if err := db.DB.Where(&Client{FacilityID: facilityID}).Preload(clause.Associations).Find(&clientProfiles).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get all clients in the specified facility: %v", err)
	}
	return clientProfiles, nil
}

// GetRecentHealthDiaryEntries fetches the health diary entries that were added after the last time the entries were
// synced to KenyaEMR
func (db *PGInstance) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry []*ClientHealthDiaryEntry
	err := db.DB.Where(&ClientHealthDiaryEntry{ClientID: clientID, Active: true}).Where("? > ?", clause.Column{Name: "created"}, lastSyncTime).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&healthDiaryEntry).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get all client health diary entries: %v", err)
	}
	return healthDiaryEntry, nil
}

// GetClientsByParams retrieves clients using the parameters provided
func (db *PGInstance) GetClientsByParams(ctx context.Context, params Client, lastSyncTime *time.Time) ([]*Client, error) {
	var clients []*Client

	// add active parameter
	params.Active = true

	query := db.DB.Where(&params)

	if lastSyncTime != nil {
		query.Where("? > ?", clause.Column{Name: "created"}, lastSyncTime)
	}

	err := query.Find(&clients).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to find clients: %v", err)
	}

	return clients, nil
}

// GetClientCCCIdentifier retrieves a client's ccc identifier
func (db *PGInstance) GetClientCCCIdentifier(ctx context.Context, clientID string) (*Identifier, error) {
	var clientIdentifiers []*ClientIdentifiers

	err := db.DB.Where(&ClientIdentifiers{ClientID: &clientID}).Find(&clientIdentifiers).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to find client identifiers: %v", err)
	}

	if len(clientIdentifiers) == 0 {
		err := fmt.Errorf("client has no associated identifiers, clientID: %v", clientID)
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	ids := []string{}
	for _, clientIdentifier := range clientIdentifiers {
		ids = append(ids, *clientIdentifier.IdentifierID)
	}

	var identifier Identifier
	err = db.DB.Where(ids).Where("identifier_type = ?", "CCC").Where("active = ?", true).First(&identifier).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to find client identifiers: %v", err)
	}

	return &identifier, nil
}

// GetScreeningToolQuestions fetches the screening tools questions
func (db *PGInstance) GetScreeningToolQuestions(ctx context.Context, toolType string) ([]ScreeningToolQuestion, error) {
	var screeningToolsQuestions []ScreeningToolQuestion
	err := db.DB.Where(&ScreeningToolQuestion{ToolType: toolType}).Find(&screeningToolsQuestions).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tools questions: %v", err)
	}
	return screeningToolsQuestions, nil
}

// GetScreeningToolQuestionByQuestionID fetches the screening tool question by question ID
func (db *PGInstance) GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*ScreeningToolQuestion, error) {
	var screeningToolQuestion ScreeningToolQuestion
	err := db.DB.Where(&ScreeningToolQuestion{ID: questionID}).First(&screeningToolQuestion).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tool question: %v", err)
	}
	return &screeningToolQuestion, nil
}

// GetClientProfileByCCCNumber returns a client profile using the CCC number
func (db *PGInstance) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*Client, error) {
	var client Client
	if err := db.DB.Joins("JOIN clients_client_identifiers on clients_client.id = clients_client_identifiers.client_id").
		Joins("JOIN clients_identifier on clients_identifier.id = clients_client_identifiers.identifier_id").
		Where("clients_identifier.identifier_type = ? AND clients_identifier.identifier_value = ? ", "CCC", CCCNumber).
		Preload(clause.Associations).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// CheckIfClientHasUnresolvedServiceRequests checks whether a client has a pending or in progress service request of the type passed in
func (db *PGInstance) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	var unresolvedServiceRequests []*ClientServiceRequest
	err := db.DB.Where(&ClientServiceRequest{ClientID: clientID, RequestType: serviceRequestType, Status: enums.ServiceRequestStatusPending.String()}).
		Or(&ClientServiceRequest{ClientID: clientID, RequestType: serviceRequestType, Status: enums.ServiceRequestStatusInProgress.String()}).
		Find(&unresolvedServiceRequests).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check for unresolved service requests: %v", err)
	}

	if len(unresolvedServiceRequests) > 0 {
		return true, nil
	}

	return false, nil
}

// GetAllRoles returns all roles
func (db *PGInstance) GetAllRoles(ctx context.Context) ([]*AuthorityRole, error) {
	var roles []*AuthorityRole
	err := db.DB.Find(&roles).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get all roles: %v", err)
	}
	return roles, nil
}

// SearchClientProfilesByCCCNumber is used to search for client profiles
// It returns clients profiles whose parts of the CCC number matches
func (db *PGInstance) SearchClientProfilesByCCCNumber(ctx context.Context, CCCNumber string) ([]*Client, error) {
	var client []*Client
	if err := db.DB.Joins("JOIN clients_client_identifiers on clients_client.id = clients_client_identifiers.client_id").
		Joins("JOIN clients_identifier on clients_identifier.id = clients_client_identifiers.identifier_id").
		Where("clients_identifier.identifier_type = ? AND clients_identifier.identifier_value ~~* ? ", "CCC", "%"+CCCNumber+"%").
		Preload(clause.Associations).Find(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to get client profile by CCC number: %v", err)
	}
	return client, nil
}

// GetUserProfileByStaffID returns a user profile using the staff ID
func (db *PGInstance) GetUserProfileByStaffID(ctx context.Context, staffID string) (*User, error) {
	var user User
	if err := db.DB.Raw(`
	 SELECT * FROM users_user
	 WHERE id = (
		SELECT user_id FROM staff_staff
		WHERE id = ?
	)`, staffID).Scan(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user profile by staff ID: %v", err)
	}

	return &user, nil
}

// GetHealthDiaryEntryByID gets the health diary entry with the given ID
func (db *PGInstance) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry *ClientHealthDiaryEntry

	err := db.DB.Where(&ClientHealthDiaryEntry{ClientHealthDiaryEntryID: &healthDiaryEntryID}).Find(&healthDiaryEntry).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get health diary entry: %v", err)
	}

	return healthDiaryEntry, nil
}

// GetServiceRequestByID returns a service request by ID
func (db *PGInstance) GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*ClientServiceRequest, error) {
	var serviceRequest ClientServiceRequest
	err := db.DB.Where(&ClientServiceRequest{ID: &serviceRequestID}).First(&serviceRequest).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get service request by ID: %v", err)
	}
	return &serviceRequest, nil
}

// GetAppointmentServiceRequests returns all appointments service requests that have been updated since the last sync time
func (db *PGInstance) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	err := db.DB.Where("created > ?", lastSyncTime).
		Where(&ClientServiceRequest{
			RequestType: enums.ServiceRequestTypeAppointments.String(),
			Status:      enums.ServiceRequestStatusPending.String(),
		}).
		Find(&serviceRequests).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get appointments service requests by last sync time: %v", err)
	}
	return serviceRequests, nil
}

// GetAppointmentByID returns an appointment by ID
func (db *PGInstance) GetAppointmentByID(ctx context.Context, appointmentID string) (*Appointment, error) {
	var appointment Appointment
	err := db.DB.Where(&Appointment{ID: appointmentID}).First(&appointment).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get appointment by ID: %v", err)
	}
	return &appointment, nil
}

// GetClientAppointmentByID returns a client appointment by ID
func (db *PGInstance) GetClientAppointmentByID(ctx context.Context, appointmentID, clientID string) (*Appointment, error) {
	var appointment Appointment
	err := db.DB.Where(&Appointment{ID: appointmentID, ClientID: clientID}).First(&appointment).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client appointment by ID: %v", err)
	}
	return &appointment, nil
}

// GetAppointmentByAppointmentUUID returns an appointment by appointment UUID
func (db *PGInstance) GetAppointmentByAppointmentUUID(ctx context.Context, appointmentUUID string) (*Appointment, error) {
	var appointment Appointment
	err := db.DB.Where(&Appointment{AppointmentUUID: appointmentUUID}).First(&appointment).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get appointment by appointment UUID: %v", err)
	}
	return &appointment, nil
}

// GetClientServiceRequests returns all system generated service requests by status passed in param
func (db *PGInstance) GetClientServiceRequests(ctx context.Context, requestType, status, clientID string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	err := db.DB.Where(&ClientServiceRequest{
		RequestType: requestType,
		Status:      status,
		ClientID:    clientID,
	}).Find(&serviceRequests).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get client service requests by status: %v", err)
	}
	return serviceRequests, nil
}

// GetActiveScreeningToolResponses returns all active screening tool responses that are within 24 hours of previous response
func (db *PGInstance) GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*ScreeningToolsResponse, error) {
	var responses []*ScreeningToolsResponse
	err := db.DB.Where(&ScreeningToolsResponse{
		ClientID: clientID,
		Active:   true,
	}).Where("created >  ?", time.Now().Add(time.Hour*-24)).Find(&responses).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get responses for client: %v", err)
	}
	return responses, nil
}
