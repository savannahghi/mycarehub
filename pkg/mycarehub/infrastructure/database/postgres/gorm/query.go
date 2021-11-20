package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
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
	GetCurrentTerms(ctx context.Context) (*TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error)
	GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	GetClientProfileByUserID(ctx context.Context, userID string) (*Client, error)
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	var facility Facility
	err := db.DB.Where(&Facility{FacilityID: id, Active: isActive}).First(&facility).Error
	if err != nil {
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
		return false, fmt.Errorf("failed to check if contact exists %v: %v", phone, err)
	}
	return true, nil
}

// RetrieveFacilityByMFLCode fetches a single facility using MFL Code
func (db *PGInstance) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*Facility, error) {
	var facility Facility
	if err := db.DB.Where(&Facility{Code: MFLCode, Active: isActive}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by MFL Code %v and status %v: %v", MFLCode, isActive, err)
	}
	return &facility, nil
}

// GetFacilities fetches all the healthcare facilities in the platform.
func (db *PGInstance) GetFacilities(ctx context.Context) ([]Facility, error) {
	var facility []Facility
	err := db.DB.Find(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query all facilities %v", err)
	}
	return facility, nil
}

// GetSecurityQuestions fetches all the security questions.
func (db *PGInstance) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error) {
	if flavour == "" {
		return nil, fmt.Errorf("flavour cannot be empty")
	}
	var securityQuestion []*SecurityQuestion
	err := db.DB.Where(&SecurityQuestion{Flavour: flavour}).Find(&securityQuestion).Error
	if err != nil {
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
			return nil, fmt.Errorf("failed to validate filter %v: %v", f.Value, err)
		}
		err = enums.ValidateFilterSortCategories(enums.FilterSortCategoryTypeFacility, f.DataType)
		if err != nil {
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
		return nil, fmt.Errorf("failed to get user by phonenumber %v: %v", phoneNumber, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user's pin using the user ID
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, IsValid: true}).First(&pin).Error; err != nil {
		return nil, fmt.Errorf("failed to get pin: %v", err)
	}
	return &pin, nil
}

// GetCurrentTerms fetches the most most recent terms of service
func (db *PGInstance) GetCurrentTerms(ctx context.Context) (*TermsOfService, error) {
	var termsOfService TermsOfService
	validTo := time.Now()
	if err := db.DB.Model(&TermsOfService{}).Where("valid_to > ?", validTo).Or("valid_to = ?", nil).Order("valid_to desc").First(&termsOfService).Error; err != nil {
		return nil, fmt.Errorf("failed to get the current terms : %v", err)
	}

	return &termsOfService, nil
}

// GetUserProfileByUserID fetches a user profile using the user ID
func (db *PGInstance) GetUserProfileByUserID(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := db.DB.Where(&User{UserID: &userID}).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by user ID %v: %v", userID, err)
	}
	return &user, nil
}

// GetSecurityQuestionByID fetches a security question using the security question ID
func (db *PGInstance) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error) {
	var securityQuestion SecurityQuestion
	if err := db.DB.Where(&SecurityQuestion{SecurityQuestionID: securityQuestionID}).First(&securityQuestion).Error; err != nil {
		return nil, fmt.Errorf("failed to get security question by ID %v: %v", securityQuestionID, err)
	}
	return &securityQuestion, nil
}

// GetSecurityQuestionResponseByID returns the security question response
func (db *PGInstance) GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*SecurityQuestionResponse, error) {
	var questionResponse SecurityQuestionResponse
	if err := db.DB.Where(&SecurityQuestionResponse{QuestionID: questionID}).First(&questionResponse).Error; err != nil {
		return nil, fmt.Errorf("failed to get the security question response by ID")
	}
	return &questionResponse, nil
}

//VerifyOTP checks from the database the validity of the provided OTP
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
		return false, fmt.Errorf("failed to verify otp with %v: %v", payload.OTP, payload.Flavour)
	}

	return true, nil
}

// GetClientProfileByUserID returns the client profile based on the user ID provided
func (db *PGInstance) GetClientProfileByUserID(ctx context.Context, userID string) (*Client, error) {
	var client Client
	if err := db.DB.Where(&Client{UserID: &userID}).Preload(clause.Associations).First(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to get client by user ID %v: %v", userID, err)
	}
	return &client, nil
}
