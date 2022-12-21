package gorm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

// GCSBaseURL is the Google Cloud Storage URL
var GCSBaseURL = serverutils.MustGetEnvVar(helpers.GoogleCloudStorageURL)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByIdentifier(ctx context.Context, identifier *FacilityIdentifier, isActive bool) (*Facility, error)
	RetrieveFacilityIdentifierByFacilityID(ctx context.Context, facilityID *string) (*FacilityIdentifier, error)
	SearchFacility(ctx context.Context, searchParameter *string) ([]Facility, error)
	GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*Facility, error)
	GetOrganisation(ctx context.Context, id string) (*Organisation, error)
	ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	ListNotifications(ctx context.Context, params *Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*Notification, *domain.Pagination, error)
	ListSurveyRespondents(ctx context.Context, params map[string]interface{}, facilityID string, pagination *domain.Pagination) ([]*UserSurvey, *domain.Pagination, error)
	ListAvailableNotificationTypes(ctx context.Context, params *Notification) ([]enums.NotificationType, error)
	ListAppointments(ctx context.Context, params *Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*Appointment, *domain.Pagination, error)
	GetUserProfileByUsername(ctx context.Context, username string) (*User, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error)
	GetUserProfileByUserID(ctx context.Context, userID *string) (*User, error)
	GetCurrentTerms(ctx context.Context) (*TermsOfService, error)
	GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*SecurityQuestion, error)
	GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error)
	GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*SecurityQuestionResponse, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error)
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
	CheckStaffExists(ctx context.Context, userID string) (bool, error)
	CheckClientExists(ctx context.Context, userID string) (bool, error)
	CheckCaregiverExists(ctx context.Context, userID string) (bool, error)
	GetClientProfile(ctx context.Context, userID string, programID string) (*Client, error)
	GetCaregiverByUserID(ctx context.Context, userID string) (*Caregiver, error)
	GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*Client, error)
	GetStaffProfile(ctx context.Context, userID string, programID string) (*StaffProfile, error)
	CheckUserHasPin(ctx context.Context, userID string) (bool, error)
	GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*UserOTP, error)
	GetClientsPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error)
	GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*SecurityQuestionResponse, error)
	GetContactByUserID(ctx context.Context, userID *string, contactType string) (*Contact, error)
	FindContacts(ctx context.Context, contactType string, contactValue string) ([]*Contact, error)
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
	GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*ClientHealthDiaryQuote, error)
	GetClientHealthDiaryEntries(ctx context.Context, params map[string]interface{}) ([]*ClientHealthDiaryEntry, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*Client, error)
	GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*ClientServiceRequest, error)
	GetStaffServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*StaffServiceRequest, error)
	CheckUserRole(ctx context.Context, userID string, role string) (bool, error)
	CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error)
	GetUserRoles(ctx context.Context, userID string, organisationID string) ([]*AuthorityRole, error)
	GetUserPermissions(ctx context.Context, userID string, organisationID string) ([]*AuthorityPermission, error)
	CheckIfUsernameExists(ctx context.Context, username string) (bool, error)
	GetCommunityByID(ctx context.Context, communityID string) (*Community, error)
	CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error)
	CheckFacilityExistsByIdentifier(ctx context.Context, identifier *FacilityIdentifier) (bool, error)
	GetClientsInAFacility(ctx context.Context, facilityID string) ([]*Client, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*ClientHealthDiaryEntry, error)
	GetClientsByParams(ctx context.Context, query Client, lastSyncTime *time.Time) ([]*Client, error)
	GetClientCCCIdentifier(ctx context.Context, clientID string) (*Identifier, error)
	SearchClientProfile(ctx context.Context, searchParameter string) ([]*Client, error)
	SearchStaffProfile(ctx context.Context, searchParameter string) ([]*StaffProfile, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*ClientServiceRequest, error)
	GetScreeningToolQuestions(ctx context.Context, toolType string) ([]ScreeningToolQuestion, error)
	GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*ScreeningToolQuestion, error)
	CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error)
	GetAllRoles(ctx context.Context) ([]*AuthorityRole, error)
	GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*ClientHealthDiaryEntry, error)
	GetUserProfileByStaffID(ctx context.Context, staffID string) (*User, error)
	GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*ClientHealthDiaryEntry, error)
	GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*ClientServiceRequest, error)
	GetStaffProfileByStaffID(ctx context.Context, staffID string) (*StaffProfile, error)
	GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*ClientServiceRequest, error)
	GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*ClientServiceRequest, error)
	GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*ScreeningToolsResponse, error)
	CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error)
	GetAnsweredScreeningToolQuestions(ctx context.Context, facilityID string, toolType string) ([]*ScreeningToolsResponse, error)
	GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID, toolType string, active bool) ([]*ScreeningToolsResponse, error)
	GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*ClientServiceRequest, error)
	GetAppointment(ctx context.Context, params *Appointment) (*Appointment, error)
	GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*UserSurvey, error)
	CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error)
	GetFacilityStaffs(ctx context.Context, facilityID string) ([]*StaffProfile, error)
	GetNotification(ctx context.Context, notificationID string) (*Notification, error)
	GetClientsByFilterParams(ctx context.Context, facilityID string, filterParams *dto.ClientFilterParamsInput) ([]*Client, error)
	SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*ClientServiceRequest, error)
	SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*StaffServiceRequest, error)
	GetScreeningToolByID(ctx context.Context, toolID string) (*ScreeningTool, error)
	GetQuestionnaireByID(ctx context.Context, questionnaireID string) (*Questionnaire, error)
	GetQuestionsByQuestionnaireID(ctx context.Context, questionnaireID string) ([]*Question, error)
	GetQuestionInputChoicesByQuestionID(ctx context.Context, questionID string) ([]*QuestionInputChoice, error)
	GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*ScreeningTool, error)
	GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, pagination *domain.Pagination) ([]*ScreeningTool, *domain.Pagination, error)
	GetScreeningToolServiceRequestOfRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*ClientServiceRequest, *domain.Pagination, error)
	GetScreeningToolResponseByID(ctx context.Context, id string) (*ScreeningToolResponse, error)
	GetScreeningToolQuestionResponsesByResponseID(ctx context.Context, responseID string) ([]*ScreeningToolQuestionResponse, error)
	GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*UserSurvey, error)
	GetClientsSurveyServiceRequest(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*ClientServiceRequest, *domain.Pagination, error)
	GetNotificationsCount(ctx context.Context, notification Notification) (int, error)
	GetStaffFacilities(ctx context.Context, staffFacility StaffFacilities, pagination *domain.Pagination) ([]*StaffFacilities, *domain.Pagination, error)
	GetClientFacilities(ctx context.Context, clientFacility ClientFacilities, pagination *domain.Pagination) ([]*ClientFacilities, *domain.Pagination, error)
	GetClientsSurveyCount(ctx context.Context, userID string) (int, error)
	SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*Caregiver, error)
	GetCaregiverManagedClients(ctx context.Context, caregiverID string, pagination *domain.Pagination) ([]*Client, *domain.Pagination, error)
	GetCaregiversClient(ctx context.Context, caregiverClient CaregiverClient) ([]*CaregiverClient, error)
	GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*Caregiver, error)
	ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*CaregiverClient, *domain.Pagination, error)
	CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error)
	GetStaffUserPrograms(ctx context.Context, userID string) ([]*Program, error)
	GetClientUserPrograms(ctx context.Context, userID string) ([]*Program, error)
	CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error)
	ListOrganisations(ctx context.Context) ([]*Organisation, error)
	GetProgramFacilities(ctx context.Context, programID string) ([]*ProgramFacility, error)
	GetProgramByID(ctx context.Context, programID string) (*Program, error)
}

// GetFacilityStaffs returns a list of staff at a particular facility
func (db PGInstance) GetFacilityStaffs(ctx context.Context, facilityID string) ([]*StaffProfile, error) {
	var staffs []*StaffProfile
	var staffModel StaffProfile
	if err := db.DB.Scopes(OrganisationScope(ctx, staffModel.TableName())).Where(StaffProfile{Active: true, DefaultFacilityID: facilityID}).Find(&staffs).Error; err != nil {
		return nil, fmt.Errorf("error retrieving staffs: %v", err)
	}

	return staffs, nil
}

// CheckStaffExists checks if there is a staff profile that exists for a user
func (db PGInstance) CheckStaffExists(ctx context.Context, userID string) (bool, error) {
	var staff StaffProfile
	err := db.DB.Where(StaffProfile{UserID: userID, Active: true}).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if staff exists: %v", err)
	}
	return true, nil
}

// CheckClientExists checks if there is a client profile that exists for a user
func (db PGInstance) CheckClientExists(ctx context.Context, userID string) (bool, error) {
	var client Client
	err := db.DB.Where(Client{UserID: &userID, Active: true}).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if client exists: %v", err)
	}
	return true, nil
}

// CheckCaregiverExists checks if there is a caregiver profile that exists for a user
func (db PGInstance) CheckCaregiverExists(ctx context.Context, userID string) (bool, error) {
	var caregiver Caregiver
	err := db.DB.Where(Caregiver{UserID: userID, Active: true}).First(&caregiver).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if caregiver exists: %v", err)
	}
	return true, nil
}

// RetrieveFacilityIdentifierByFacilityID gets a facility identifier by facility id
func (db PGInstance) RetrieveFacilityIdentifierByFacilityID(ctx context.Context, facilityID *string) (*FacilityIdentifier, error) {
	facilityIdentifier := &FacilityIdentifier{}
	if err := db.DB.Where(&FacilityIdentifier{FacilityID: *facilityID}).First(facilityIdentifier).Error; err != nil {
		return nil, err
	}

	return facilityIdentifier, nil
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility id cannot be nil")
	}
	var facility Facility
	err := db.DB.Where(&Facility{FacilityID: id, Active: isActive}).First(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}

// GetOrganisation retrieves an organisation using the provided id
func (db *PGInstance) GetOrganisation(ctx context.Context, id string) (*Organisation, error) {
	var organisation Organisation

	err := db.DB.Where(&Organisation{ID: &id, Active: true}).First(&organisation).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get organisation by ID %s err: %w", id, err)
	}

	return &organisation, nil
}

// CheckIfPhoneNumberExists checks if phone exists in the database.
func (db *PGInstance) CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
	var contact Contact
	if phone == "" || !flavour.IsValid() {
		return false, fmt.Errorf("invalid flavour: %v", flavour)
	}
	err := db.DB.Scopes(OrganisationScope(ctx, contact.TableName())).Model(&Contact{}).Where(&Contact{Value: phone, OptedIn: isOptedIn}).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if contact exists %v: %v", phone, err)
	}
	return true, nil
}

// RetrieveFacilityByIdentifier fetches a single facility using MFL Code
func (db *PGInstance) RetrieveFacilityByIdentifier(ctx context.Context, identifier *FacilityIdentifier, isActive bool) (*Facility, error) {
	if identifier.Type == "" || identifier.Value == "" {
		return nil, fmt.Errorf("facility  identifier cannot be nil")
	}
	var facility Facility

	if err := db.DB.Joins("JOIN common_facility_identifier on common_facility.id = common_facility_identifier.facility_id").
		Where("common_facility_identifier.identifier_type = ? AND common_facility_identifier.identifier_value = ?", identifier.Type, identifier.Value).
		First(&facility).Error; err != nil {
		return nil, err
	}

	return &facility, nil
}

// SearchFacility fetches facilities by pattern matching against the facility name or mflcode
func (db *PGInstance) SearchFacility(ctx context.Context, searchParameter *string) ([]Facility, error) {
	var facility []Facility
	var facilityModel Facility

	err := db.DB.Scopes(OrganisationScope(ctx, facilityModel.TableName())).
		Where(db.DB.Where("common_facility.name ILIKE ?", "%"+*searchParameter+"%").
			Or("CAST(common_facility.mfl_code as text) ILIKE ?", "%"+*searchParameter+"%")).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: false}).Find(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query facilities %w", err)
	}
	return facility, nil
}

// GetFacilitiesWithoutFHIRID fetches all the healthcare facilities in the platform without FHIR Organisation ID
func (db *PGInstance) GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*Facility, error) {
	var facility []*Facility
	orgID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organisation id from context: %v", err)
	}

	err = db.DB.Raw(
		`SELECT * FROM common_facility WHERE fhir_organization_id IS NULL AND organisation_id = ? `, orgID).Scan(&facility).Error
	if err != nil {
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

	err := db.DB.Where(&SecurityQuestion{Flavour: flavour, Active: true}).Find(&securityQuestion).Error
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

	f := Facility{}
	tx := db.DB.Scopes(OrganisationScope(ctx, f.TableName())).Begin()
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
		paginate(facilities, &paginatedFacilities.Pagination, resultCount, db.DB.Scopes(OrganisationScope(ctx, f.TableName()))),
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
			Active:      f.Active,
			County:      f.Country,
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
	a := Appointment{}
	tx := db.DB.Scopes(OrganisationScope(ctx, a.TableName())).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, pageInfo, fmt.Errorf("failed to initialize filter appointments transaction %v", err)
	}

	transaction := tx.Where(params)

	if err := addFilters(transaction, filters); err != nil {
		return nil, pageInfo, fmt.Errorf("failed to add filters: %v", err)
	}

	transaction.Find(&appointments)

	resultCount = int64(len(appointments))

	if pagination != nil {
		transaction = tx.Scopes(
			paginate(appointments, pagination, resultCount, db.DB.Scopes(OrganisationScope(ctx, a.TableName()))),
		).Where(params)

		if err := addFilters(transaction, filters); err != nil {
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

// ListNotifications retrieves notifications using the provided parameters and filters
func (db *PGInstance) ListNotifications(ctx context.Context, params *Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*Notification, *domain.Pagination, error) {
	var count int64
	var notifications []*Notification

	n := Notification{}
	userNotificationsQuery := db.DB.Scopes(OrganisationScope(ctx, n.TableName())).Where(Notification{UserID: params.UserID, Flavour: params.Flavour, Active: params.Active})
	if err := addFilters(userNotificationsQuery, filters); err != nil {
		return nil, pagination, fmt.Errorf("failed to add filters to transaction: %v", err)
	}

	tx := db.DB.Scopes(OrganisationScope(ctx, n.TableName())).Model(&Notification{}).Or(userNotificationsQuery)

	// include facility notifications
	if params.FacilityID != nil {
		facilityNotificationsQuery := db.DB.Scopes(OrganisationScope(ctx, n.TableName())).Where(Notification{FacilityID: params.FacilityID, Flavour: params.Flavour, Active: params.Active})
		if err := addFilters(facilityNotificationsQuery, filters); err != nil {
			return nil, pagination, fmt.Errorf("failed to add filters to transaction: %v", err)
		}

		tx.Or(facilityNotificationsQuery)
	}

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, pagination, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&notifications).Error; err != nil {
		return nil, pagination, fmt.Errorf("failed to execute paginated query: %v", err)
	}

	return notifications, pagination, nil
}

// ListSurveyRespondents retrieves survey respondents using the provided parameters. It also paginates the results
func (db *PGInstance) ListSurveyRespondents(ctx context.Context, params map[string]interface{}, facilityID string, pagination *domain.Pagination) ([]*UserSurvey, *domain.Pagination, error) {
	var count int64
	var userSurveys []*UserSurvey
	u := UserSurvey{}

	tx := db.DB.Scopes(OrganisationScope(ctx, u.TableName())).
		Model(&UserSurvey{}).
		Joins("JOIN clients_client on clients_client.user_id = common_usersurveys.user_id").
		Where("clients_client.current_facility_id = ?", facilityID).
		Where(params)

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&userSurveys).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to execute paginated query: %v", err)
	}

	return userSurveys, pagination, nil
}

// ListAvailableNotificationTypes retrieves the distinct notification types available for a user
func (db *PGInstance) ListAvailableNotificationTypes(ctx context.Context, params *Notification) ([]enums.NotificationType, error) {
	var notificationTypes []enums.NotificationType

	n := Notification{}

	tx := db.DB.Scopes(OrganisationScope(ctx, n.TableName())).Model(&Notification{}).Or(Notification{UserID: params.UserID, Flavour: params.Flavour, Active: params.Active})

	// include facility notification types
	if params.FacilityID != nil {
		tx.Or(Notification{FacilityID: params.FacilityID, Flavour: params.Flavour, Active: params.Active})
	}

	if err := tx.Distinct("notification_type").Find(&notificationTypes).Error; err != nil {
		return notificationTypes, fmt.Errorf("failed to execute query: %v", err)
	}

	return notificationTypes, nil
}

// GetUserProfileByPhoneNumber retrieves a user profile using their phone number
func (db *PGInstance) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	var user User

	if err := db.DB.Joins("JOIN common_contact on users_user.id = common_contact.user_id").Where("common_contact.contact_value = ?", phoneNumber).
		Preload(clause.Associations).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by phonenumber %v: %v", phoneNumber, err)
	}
	return &user, nil
}

// GetUserProfileByUsername retrieves a user using their username
func (db *PGInstance) GetUserProfileByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	if err := db.DB.Preload(clause.Associations).Where(User{Username: username}).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %w", username, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user's pin using the user ID and Flavour
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, IsValid: true}).First(&pin).Error; err != nil {
		return nil, fmt.Errorf("failed to get pin: %v", err)
	}
	return &pin, nil
}

// GetCurrentTerms fetches the most recent terms of service depending on the flavour
func (db *PGInstance) GetCurrentTerms(ctx context.Context) (*TermsOfService, error) {
	var termsOfService TermsOfService
	validTo := time.Now()
	if err := db.DB.WithContext(ctx).Model(&TermsOfService{}).Where(db.DB.Where("valid_to > ?", validTo).Or("valid_to = ?", nil).Order("valid_to desc")).First(&termsOfService).Statement.Error; err != nil {
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
		return nil, fmt.Errorf("failed to get user by user ID %v: %v", userID, err)
	}
	return &user, nil
}

// GetSecurityQuestionByID fetches a security question using the security question ID
func (db *PGInstance) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*SecurityQuestion, error) {
	var securityQuestion SecurityQuestion
	if err := db.DB.Where(&SecurityQuestion{SecurityQuestionID: securityQuestionID}).First(&securityQuestion).Error; err != nil {
		return nil, fmt.Errorf("failed to get security question by ID %v: %w", securityQuestionID, err)
	}
	return &securityQuestion, nil
}

// GetNotification retrieve a notification using the provided ID
func (db *PGInstance) GetNotification(ctx context.Context, notificationID string) (*Notification, error) {
	var notification Notification
	if err := db.DB.Scopes(OrganisationScope(ctx, notification.TableName())).Where(&Notification{ID: notificationID}).First(&notification).Error; err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

// GetAnsweredScreeningToolQuestions returns the answered screening tool questions
func (db *PGInstance) GetAnsweredScreeningToolQuestions(ctx context.Context, facilityID string, toolType string) ([]*ScreeningToolsResponse, error) {
	var screeningToolResponses []*ScreeningToolsResponse

	orgID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organisation ID from context: %w", err)
	}

	err = db.DB.Raw(
		`
		SELECT * FROM screeningtools_screeningtoolsquestion
		JOIN screeningtools_screeningtoolsresponse
		ON screeningtools_screeningtoolsquestion.id = screeningtools_screeningtoolsresponse.question_id
		JOIN clients_client
		ON clients_client.id = screeningtools_screeningtoolsresponse.client_id
		JOIN clients_servicerequest
		ON clients_client.id = clients_servicerequest.client_id
		WHERE clients_servicerequest.status = 'PENDING'
		AND screeningtools_screeningtoolsresponse.active = ?
		AND tool_type = ? 
		AND clients_servicerequest.meta->>'question_type'  = ?
		AND clients_client.current_facility_id = ?
		AND screeningtools_screeningtoolsresponse.organisation_id = ?
		AND clients_client.organisation_id = ?
		`, true, toolType, toolType, facilityID, orgID, orgID).
		Scan(&screeningToolResponses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get answered screening tool questions: %v", err)
	}

	return screeningToolResponses, nil
}

// GetSecurityQuestionResponse returns the security question response
func (db *PGInstance) GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*SecurityQuestionResponse, error) {
	var questionResponse SecurityQuestionResponse
	if err := db.DB.Scopes(OrganisationScope(ctx, questionResponse.TableName())).Where(&SecurityQuestionResponse{QuestionID: questionID, UserID: userID}).First(&questionResponse).Error; err != nil {
		return nil, fmt.Errorf("failed to get the security question response by ID")
	}
	return &questionResponse, nil
}

// VerifyOTP checks from the database for the validity of the provided OTP
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify otp with %v: %v", payload.OTP, payload.Flavour)
	}

	return true, nil
}

// GetClientProfile returns the client profile based on the user ID provided
func (db *PGInstance) GetClientProfile(ctx context.Context, userID string, programID string) (*Client, error) {
	var client Client

	tx := db.DB.Where(&Client{UserID: &userID})

	if programID != "" {
		tx = tx.Where("program_id = ?", programID)
	}

	tx = tx.Preload(clause.Associations).First(&client)

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed to get client profile %w", err)
	}

	return &client, nil
}

// GetCaregiverByUserID returns the caregiver record of the provided user ID
func (db *PGInstance) GetCaregiverByUserID(ctx context.Context, userID string) (*Caregiver, error) {
	var caregiver *Caregiver

	if err := db.DB.Scopes(OrganisationScope(ctx, caregiver.TableName())).Where(Caregiver{UserID: userID}).First(&caregiver).Error; err != nil {
		return nil, fmt.Errorf("failed to get caregiver by user ID %v: %w", userID, err)
	}

	return caregiver, nil
}

// GetStaffProfile returns the staff profile
func (db *PGInstance) GetStaffProfile(ctx context.Context, userID string, programID string) (*StaffProfile, error) {
	var staff StaffProfile

	tx := db.DB.Where(&StaffProfile{UserID: userID})

	if programID != "" {
		tx = tx.Where("program_id = ?", programID)
	}

	tx = tx.Preload(clause.Associations).First(&staff)

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("unable to get staff profile %w", err)
	}

	return &staff, nil
}

// SearchStaffProfile searches retrieves staff profile(s) based on pattern matching against the username, staff number
// or the phonenumber.
func (db *PGInstance) SearchStaffProfile(ctx context.Context, searchParameter string) ([]*StaffProfile, error) {
	var staff []*StaffProfile

	if err := db.DB.Joins("JOIN users_user ON users_user.id = staff_staff.user_id").
		Joins("JOIN common_contact on users_user.id = common_contact.user_id").
		Where(
			db.DB.Where("staff_staff.staff_number ILIKE ? ", "%"+searchParameter+"%").
				Or("users_user.username ILIKE ? ", "%"+searchParameter+"%").
				Or("common_contact.contact_value ILIKE ?", "%"+searchParameter+"%"),
		).Where("users_user.active = ?", true).Find(&staff).Error; err != nil {
		return nil, fmt.Errorf("unable to get staff user %w", err)
	}

	return staff, nil
}

// SearchCaregiverUser searches and retrieves caregiver user(s) based on pattern matching against the username, phone number or the caregiver number
func (db *PGInstance) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*Caregiver, error) {
	var caregivers []*Caregiver
	var caregiverModel Caregiver

	if err := db.DB.Scopes(OrganisationScope(ctx, caregiverModel.TableName())).
		Joins("JOIN users_user ON users_user.id = caregivers_caregiver.user_id").
		Joins("JOIN common_contact on users_user.id = common_contact.user_id").
		Where(
			db.DB.Scopes(OrganisationScope(ctx, caregiverModel.TableName())).
				Where("caregivers_caregiver.caregiver_number ILIKE ? ", "%"+searchParameter+"%").
				Or("users_user.username ILIKE ? ", "%"+searchParameter+"%").
				Or("common_contact.contact_value ILIKE ?", "%"+searchParameter+"%"),
		).Where("users_user.active = ?", true).
		Find(&caregivers).Error; err != nil {
		return nil, fmt.Errorf("unable to get caregiver user %w", err)
	}
	return caregivers, nil
}

// CheckUserHasPin performs a look-up on the pins' table to check whether a user has a pin
func (db *PGInstance) CheckUserHasPin(ctx context.Context, userID string) (bool, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID}).Find(&pin).Error; err != nil {
		return false, err
	}
	return true, nil
}

// GetOTP fetches an OTP from the database
func (db *PGInstance) GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*UserOTP, error) {
	var userOTP UserOTP
	if err := db.DB.Where(&UserOTP{PhoneNumber: phoneNumber, Flavour: flavour}).First(&userOTP).Error; err != nil {
		return nil, fmt.Errorf("failed to get otp: %v", err)
	}
	return &userOTP, nil
}

// GetUserSecurityQuestionsResponses fetches the security question responses that the user has responded to
func (db *PGInstance) GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*SecurityQuestionResponse, error) {
	var securityQuestionResponses []*SecurityQuestionResponse
	var responses SecurityQuestionResponse

	if err := db.DB.Scopes(OrganisationScope(ctx, responses.TableName())).Where(&SecurityQuestionResponse{UserID: userID, Active: true}).Find(&securityQuestionResponses).Error; err != nil {
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
	if err := db.DB.Scopes(OrganisationScope(ctx, contact.TableName())).Where(&Contact{UserID: userID, Type: contactType}).First(&contact).Error; err != nil {
		return nil, fmt.Errorf("failed to get contact: %v", err)
	}
	return &contact, nil
}

// CanRecordHeathDiary checks whether a user can record a health diary
// if the last record is less than 24 hours ago, the user cannot record a new entry
// if the last record is more than 24 hours ago, the user can record a new entry
func (db *PGInstance) CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error) {
	var clientHealthDiaryEntry []*ClientHealthDiaryEntry

	err := db.DB.Where("client_id = ?", clientID).Order("created desc").Find(&clientHealthDiaryEntry).Error
	if err != nil {
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
func (db *PGInstance) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*ClientHealthDiaryQuote, error) {
	var healthDiaryQuote []*ClientHealthDiaryQuote
	var quote ClientHealthDiaryQuote

	err := db.DB.Scopes(OrganisationScope(ctx, quote.TableName())).Where("active = true").Limit(limit).Order("RANDOM()").Find(&healthDiaryQuote).Error
	if err != nil {
		return nil, err
	}
	return healthDiaryQuote, nil
}

// GetClientHealthDiaryEntries gets all health diary entries that belong to a specific client
func (db *PGInstance) GetClientHealthDiaryEntries(ctx context.Context, params map[string]interface{}) ([]*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry []*ClientHealthDiaryEntry
	var entry ClientHealthDiaryEntry

	err := db.DB.Scopes(OrganisationScope(ctx, entry.TableName())).Where(params).Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&healthDiaryEntry).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all client health diary entries: %v", err)
	}
	return healthDiaryEntry, nil
}

// GetServiceRequestsForKenyaEMR gets all the service requests to be used by the KenyaEMR.
func (db *PGInstance) GetServiceRequestsForKenyaEMR(ctx context.Context, facilityID string, lastSyncTime time.Time) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	var clientServiceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).
		Where(&ClientServiceRequest{FacilityID: facilityID}).Where("created > ?", lastSyncTime).
		Where(db.DB.Where(&ClientServiceRequest{RequestType: string(enums.ServiceRequestTypeScreeningToolsRedFlag)}).
			Or(&ClientServiceRequest{RequestType: string(enums.ServiceRequestTypeRedFlag)})).
		Find(&serviceRequests).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get service requests: %v", err)
	}
	return serviceRequests, nil
}

// GetStaffPendingServiceRequestsCount gets the number of staffs pending pin reset service requests
func (db *PGInstance) GetStaffPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	var staffServiceRequest []*StaffServiceRequest
	var serviceRequest StaffServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, serviceRequest.TableName())).Model(&StaffServiceRequest{}).Where(&StaffServiceRequest{DefaultFacilityID: &facilityID, RequestType: "STAFF_PIN_RESET", Status: "PENDING"}).Find(&staffServiceRequest).Error
	if err != nil {
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
	var clientServiceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{FacilityID: facilityID, Status: "PENDING"}).Find(&serviceRequests).Error
	if err != nil {
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
				RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag,
				Total:       0,
			},
			{
				RequestType: enums.ServiceRequestTypeSurveyRedFlag,
			},
			{
				RequestType: enums.ServiceRequestTypeHomePageHealthDiary,
			},
			{
				RequestType: enums.ServiceRequestTypeStaffPinReset,
			},
			{
				RequestType: enums.ServiceRequestTypeAppointments,
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
		if request.RequestType == enums.ServiceRequestTypeScreeningToolsRedFlag.String() {
			serviceRequestsCount.RequestsTypeCount[2].Total++
		}
		if request.RequestType == enums.ServiceRequestTypeSurveyRedFlag.String() {
			serviceRequestsCount.RequestsTypeCount[3].Total++
		}
		if request.RequestType == enums.ServiceRequestTypeHomePageHealthDiary.String() {
			serviceRequestsCount.RequestsTypeCount[4].Total++
		}
		if request.RequestType == enums.ServiceRequestTypeStaffPinReset.String() {
			serviceRequestsCount.RequestsTypeCount[5].Total++
		}
		if request.RequestType == enums.ServiceRequestTypeAppointments.String() {
			serviceRequestsCount.RequestsTypeCount[6].Total++
		}
	}

	return &serviceRequestsCount, nil
}

// GetClientProfileByClientID fetches a client from the database
func (db *PGInstance) GetClientProfileByClientID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	err := db.DB.Scopes(OrganisationScope(ctx, client.TableName())).Where(&Client{ID: &clientID}).Preload("User.Contacts").Preload(clause.Associations).First(&client).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	return &client, nil
}

// GetStaffProfileByStaffID fetches a staff from the database
func (db *PGInstance) GetStaffProfileByStaffID(ctx context.Context, staffID string) (*StaffProfile, error) {
	var staff StaffProfile

	err := db.DB.Scopes(OrganisationScope(ctx, staff.TableName())).Where(&StaffProfile{ID: &staffID}).Preload("UserProfile.Contacts").Preload(clause.Associations).First(&staff).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %v", err)
	}
	return &staff, nil
}

// GetServiceRequests fetches clients service requests from the database according to the type and or status passed
func (db *PGInstance) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	var clientServiceRequest ClientServiceRequest

	if requestType != nil && requestStatus == nil {
		err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).
			Where(&ClientServiceRequest{RequestType: *requestType, FacilityID: facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&serviceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}

	} else if requestType == nil && requestStatus != nil {
		err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).Where(&ClientServiceRequest{Status: *requestStatus, FacilityID: facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&serviceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}

	} else if requestType != nil && requestStatus != nil {
		err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).Where(&ClientServiceRequest{RequestType: *requestType, Status: *requestStatus, FacilityID: facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&serviceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}

	} else {
		err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).Where(&ClientServiceRequest{FacilityID: facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&serviceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get service requests: %v", err)
		}
	}

	return serviceRequests, nil
}

// GetStaffServiceRequests gets all the staff's service requests depending on the provided parameters
func (db *PGInstance) GetStaffServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string) ([]*StaffServiceRequest, error) {
	var staffServiceRequests []*StaffServiceRequest
	var staffServiceRequest StaffServiceRequest

	if requestType != nil && requestStatus != nil {
		err := db.DB.Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).Where(&StaffServiceRequest{RequestType: *requestType, Status: *requestStatus, DefaultFacilityID: &facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&staffServiceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else if requestType == nil && requestStatus != nil {
		err := db.DB.Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).Where(&StaffServiceRequest{Status: *requestStatus, DefaultFacilityID: &facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&staffServiceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else if requestType != nil && requestStatus == nil {
		err := db.DB.Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).Where(&StaffServiceRequest{RequestType: *requestType, DefaultFacilityID: &facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&staffServiceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	} else {
		err := db.DB.Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).Where(&StaffServiceRequest{DefaultFacilityID: &facilityID}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).
			Find(&staffServiceRequests).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get staff service requests: %v", err)
		}
	}

	return staffServiceRequests, nil
}

// CheckUserRole checks if a user has a specific role
func (db *PGInstance) CheckUserRole(ctx context.Context, userID string, role string) (bool, error) {
	var returneduserID string
	orgID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get organisation id from context: %v", err)
	}

	err = db.DB.Raw(
		`
		SELECT user_id 
		FROM authority_authorityrole_users 
		WHERE user_id = ? 
		AND authority_authorityrole_users.authorityrole_id = 
		(SELECT id FROM authority_authorityrole WHERE name = ? AND organisation_id = ?)
		`, userID, role, orgID).Scan(&returneduserID).Error

	if returneduserID == "" {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if user has role: %v", err)
	}

	return true, nil
}

// CheckUserPermission checks if a user has a specific permission
func (db *PGInstance) CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error) {
	var returneduserID string
	orgID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get organisation id from context: %v", err)
	}

	err = db.DB.Raw(
		`
		SELECT user_id 
		FROM authority_authorityrole_users 
		WHERE authorityrole_id =
		(
			SELECT 	authorityrole_id 
			FROM authority_authorityrole_permissions
			WHERE authoritypermission_id = (SELECT id FROM authority_authoritypermission WHERE name = ? AND organisation_id = ?)
		)
		And 
		user_id = ?
		`, permission, orgID, userID).Scan(&returneduserID).Error

	if returneduserID == "" {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if user has permission: %v", err)
	}

	return true, nil
}

// GetUserRoles fetches a user's roles from the database
func (db *PGInstance) GetUserRoles(ctx context.Context, userID string, organisationID string) ([]*AuthorityRole, error) {
	var roles []*AuthorityRole

	err := db.DB.Raw(
		`
		SELECT * 
		FROM authority_authorityrole_users 
		JOIN authority_authorityrole ON authority_authorityrole_users.authorityrole_id = authority_authorityrole.id
		WHERE user_id = ?
		AND authority_authorityrole.organisation_id = ?
		`, userID, organisationID,
	).Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}

	return roles, nil
}

// GetUserPermissions fetches a user's permissions from the database
func (db *PGInstance) GetUserPermissions(ctx context.Context, userID string, organisationID string) ([]*AuthorityPermission, error) {
	var permissions []*AuthorityPermission
	err := db.DB.Raw(
		`
		SELECT authority_authoritypermission.created, authority_authoritypermission.updated, authority_authoritypermission.id, authority_authoritypermission.name, authority_authoritypermission.organisation_id
		FROM authority_authorityrole_users 
		JOIN authority_authorityrole_permissions ON authority_authorityrole_users.authorityrole_id = authority_authorityrole_permissions.authorityrole_id
		JOIN authority_authoritypermission ON authority_authorityrole_permissions.authoritypermission_id = authority_authoritypermission.id
		WHERE user_id = ? AND authority_authoritypermission.organisation_id = ?
		`, userID, organisationID).Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %v", err)
	}

	return permissions, nil
}

// CheckIfUsernameExists checks to see whether the provided username exists
func (db *PGInstance) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	var user User
	err := db.DB.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if username exists: %v", err)
	}

	return true, nil
}

// GetCommunityByID fetches the community using its ID
func (db *PGInstance) GetCommunityByID(ctx context.Context, communityID string) (*Community, error) {
	var community *Community

	err := db.DB.Scopes(OrganisationScope(ctx, community.TableName())).Where(&Community{ID: communityID}).First(&community).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find community by ID %s", communityID)
	}

	return community, nil
}

// CheckIdentifierExists checks whether an identifier of a certain type and value exists
// Used to validate uniqueness and prevent duplicates
func (db *PGInstance) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	var identifier *Identifier

	err := db.DB.Scopes(OrganisationScope(ctx, identifier.TableName())).Where(&Identifier{IdentifierType: identifierType, IdentifierValue: identifierValue, Active: true}).First(&identifier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check identifier of type: %s and value: %s", identifierType, identifierValue)
	}

	return true, nil
}

// CheckFacilityExistsByIdentifier checks whether a facility exists using the mfl code.
// Used to validate existence of a facility
func (db *PGInstance) CheckFacilityExistsByIdentifier(ctx context.Context, identifier *FacilityIdentifier) (bool, error) {
	_, err := db.RetrieveFacilityByIdentifier(ctx, identifier, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check for facility: %s", err)
	}

	return true, nil
}

// GetClientsInAFacility returns all the clients registered within a specified facility
func (db *PGInstance) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*Client, error) {
	var clientProfiles []*Client
	var client Client
	if err := db.DB.Scopes(OrganisationScope(ctx, client.TableName())).Where(&Client{FacilityID: facilityID}).Preload(clause.Associations).Find(&clientProfiles).Error; err != nil {
		return nil, fmt.Errorf("failed to get all clients in the specified facility: %v", err)
	}
	return clientProfiles, nil
}

// GetRecentHealthDiaryEntries fetches the health diary entries that were added after the last time the entries were
// synced to KenyaEMR
func (db *PGInstance) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, clientID string) ([]*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry []*ClientHealthDiaryEntry
	var clientHealthDiaryEntry ClientHealthDiaryEntry

	err := db.DB.Scopes(OrganisationScope(ctx, clientHealthDiaryEntry.TableName())).Where(&ClientHealthDiaryEntry{ClientID: clientID, Active: true}).Where("? > ?", clause.Column{Name: "created"}, lastSyncTime).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&healthDiaryEntry).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all client health diary entries: %v", err)
	}
	return healthDiaryEntry, nil
}

// GetClientsByParams retrieves clients using the parameters provided
func (db *PGInstance) GetClientsByParams(ctx context.Context, params Client, lastSyncTime *time.Time) ([]*Client, error) {
	var clients []*Client
	var client Client

	// add active parameter
	params.Active = true

	query := db.DB.Scopes(OrganisationScope(ctx, client.TableName())).Where(&params)

	if lastSyncTime != nil {
		query.Where("? > ?", clause.Column{Name: "created"}, lastSyncTime)
	}

	err := query.Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find clients: %v", err)
	}

	return clients, nil
}

// GetClientCCCIdentifier retrieves a client's ccc identifier
func (db *PGInstance) GetClientCCCIdentifier(ctx context.Context, clientID string) (*Identifier, error) {
	var clientIdentifiers []*ClientIdentifiers

	err := db.DB.Where(&ClientIdentifiers{ClientID: &clientID}).Find(&clientIdentifiers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find client identifiers: %v", err)
	}

	if len(clientIdentifiers) == 0 {
		err := fmt.Errorf("client has no associated identifiers, clientID: %v", clientID)
		return nil, err
	}

	ids := []string{}
	for _, clientIdentifier := range clientIdentifiers {
		ids = append(ids, *clientIdentifier.IdentifierID)
	}

	var identifier Identifier
	err = db.DB.Scopes(OrganisationScope(ctx, identifier.TableName())).Where(ids).First(&identifier).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find client identifiers: %v", err)
	}

	return &identifier, nil
}

// GetScreeningToolQuestions fetches the screening tools questions
func (db *PGInstance) GetScreeningToolQuestions(ctx context.Context, toolType string) ([]ScreeningToolQuestion, error) {
	var screeningToolsQuestions []ScreeningToolQuestion
	var screeningToolQuestion ScreeningToolQuestion

	err := db.DB.Scopes(OrganisationScope(ctx, screeningToolQuestion.TableName())).Where(&ScreeningToolQuestion{ToolType: toolType}).
		Order("sequence asc").
		Find(&screeningToolsQuestions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get screening tools questions: %v", err)
	}
	return screeningToolsQuestions, nil
}

// GetScreeningToolQuestionByQuestionID fetches the screening tool question by question ID
func (db *PGInstance) GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*ScreeningToolQuestion, error) {
	var screeningToolQuestion ScreeningToolQuestion

	err := db.DB.Scopes(OrganisationScope(ctx, screeningToolQuestion.TableName())).Where(&ScreeningToolQuestion{ID: questionID}).First(&screeningToolQuestion).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get screening tool question: %v", err)
	}
	return &screeningToolQuestion, nil
}

// GetClientProfileByCCCNumber returns a client profile using the CCC number
func (db *PGInstance) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*Client, error) {
	var client Client
	var identifier Identifier

	if err := db.DB.Scopes(OrganisationScope(ctx, client.TableName())).
		Joins("JOIN clients_client_identifiers on clients_client.id = clients_client_identifiers.client_id").
		Joins("JOIN clients_identifier on clients_identifier.id = clients_client_identifiers.identifier_id").
		Scopes(OrganisationScope(ctx, identifier.TableName())).
		Where("clients_identifier.identifier_type = ? AND clients_identifier.identifier_value = ? ", "CCC", CCCNumber).
		Preload(clause.Associations).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// CheckIfClientHasUnresolvedServiceRequests checks whether a client has a pending or in progress service request of the type passed in
func (db *PGInstance) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	var unresolvedServiceRequests []*ClientServiceRequest
	var clientServiceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, clientServiceRequest.TableName())).Where(&ClientServiceRequest{ClientID: clientID, RequestType: serviceRequestType, Status: enums.ServiceRequestStatusPending.String()}).
		Or(&ClientServiceRequest{ClientID: clientID, RequestType: serviceRequestType, Status: enums.ServiceRequestStatusInProgress.String()}).
		Find(&unresolvedServiceRequests).Error
	if err != nil {
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
	var roleModel AuthorityRole

	err := db.DB.Scopes(OrganisationScope(ctx, roleModel.TableName())).Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %v", err)
	}
	return roles, nil
}

// SearchClientProfile is used to query for a client profile. It uses pattern matching against the ccc number, phonenumber or username
func (db *PGInstance) SearchClientProfile(ctx context.Context, searchParameter string) ([]*Client, error) {
	var client []*Client
	var identifier Identifier
	var clientModel Client

	if err := db.DB.Scopes(OrganisationScope(ctx, clientModel.TableName())).
		Joins("JOIN users_user on users_user.id = clients_client.user_id").
		Joins("JOIN clients_client_identifiers on clients_client.id = clients_client_identifiers.client_id").
		Joins("JOIN clients_identifier on clients_identifier.id = clients_client_identifiers.identifier_id").
		Joins("JOIN common_contact on users_user.id = common_contact.user_id").
		Where(db.DB.Scopes(OrganisationScope(ctx, identifier.TableName())).
			Where("clients_identifier.identifier_value ILIKE ? AND clients_identifier.identifier_type = ?", "%"+searchParameter+"%", "CCC").
			Or("users_user.username ILIKE ? ", "%"+searchParameter+"%").
			Or("common_contact.contact_value ILIKE ?", "%"+searchParameter+"%"),
		).Where("users_user.active = ?", true).Preload(clause.Associations).Find(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to get client profile: %w", err)
	}

	return client, nil
}

// GetUserProfileByStaffID returns a user profile using the staff ID
func (db *PGInstance) GetUserProfileByStaffID(ctx context.Context, staffID string) (*User, error) {
	var user User

	orgID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organisation id from context: %v", err)
	}

	if err := db.DB.Raw(`
	 SELECT * FROM users_user
	 WHERE id = (
		SELECT user_id FROM staff_staff
		WHERE id = ? AND organisation_id = ?
	)`, staffID, orgID).Scan(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user profile by staff ID: %v", err)
	}

	return &user, nil
}

// GetHealthDiaryEntryByID gets the health diary entry with the given ID
func (db *PGInstance) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry *ClientHealthDiaryEntry

	err := db.DB.Scopes(OrganisationScope(ctx, healthDiaryEntry.TableName())).Where(&ClientHealthDiaryEntry{ClientHealthDiaryEntryID: &healthDiaryEntryID}).Find(&healthDiaryEntry).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get health diary entry: %v", err)
	}

	return healthDiaryEntry, nil
}

// GetSharedHealthDiaryEntries gets the recently shared health diary entry shared by the client to a health care worker
// and returns the entry.
// The health care worker will only see the entry as long as they share the same facility with the health care worker
func (db *PGInstance) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*ClientHealthDiaryEntry, error) {
	var healthDiaryEntry []*ClientHealthDiaryEntry
	var clientHealthDiaryEntry ClientHealthDiaryEntry
	var clientFacility ClientFacility

	err := db.DB.Scopes(OrganisationScope(ctx, clientHealthDiaryEntry.TableName())).
		Joins("JOIN clients_clientfacility on clients_healthdiaryentry.client_id = clients_clientfacility.client_id").
		Scopes(OrganisationScope(ctx, clientFacility.TableName())).
		Where("clients_healthdiaryentry.share_with_health_worker = ? AND clients_healthdiaryentry.client_id = ? AND clients_clientfacility.facility_id = ? ", true, clientID, facilityID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "shared_at"}, Desc: true}).Find(&healthDiaryEntry).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get shared health diary entry: %v", err)
	}

	return healthDiaryEntry, nil
}

// GetServiceRequestByID returns a service request by ID
func (db *PGInstance) GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*ClientServiceRequest, error) {
	var serviceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, serviceRequest.TableName())).Where(&ClientServiceRequest{ID: &serviceRequestID}).First(&serviceRequest).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get service request by ID: %v", err)
	}
	return &serviceRequest, nil
}

// GetAppointmentServiceRequests returns all appointments service requests that have been updated since the last sync time
func (db *PGInstance) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, facilityID string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	var serviceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, serviceRequest.TableName())).Where("created > ?", lastSyncTime).
		Where(&ClientServiceRequest{
			RequestType: enums.ServiceRequestTypeAppointments.String(),
			Status:      enums.ServiceRequestStatusPending.String(),
			FacilityID:  facilityID,
		}).
		Find(&serviceRequests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments service requests by last sync time: %v", err)
	}
	return serviceRequests, nil
}

// GetAppointment returns an appointment by provided params
func (db *PGInstance) GetAppointment(ctx context.Context, params *Appointment) (*Appointment, error) {
	var appointment Appointment
	err := db.DB.Scopes(OrganisationScope(ctx, appointment.TableName())).Where(params).Order(clause.OrderByColumn{Column: clause.Column{Name: "updated"}, Desc: true}).First(&appointment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment by ID: %w", err)
	}
	return &appointment, nil
}

// GetClientServiceRequests returns all system generated service requests by status passed in param
func (db *PGInstance) GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*ClientServiceRequest, error) {
	var serviceRequests []*ClientServiceRequest
	var serviceRequest ClientServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, serviceRequest.TableName())).Where(&ClientServiceRequest{
		RequestType: requestType,
		Status:      status,
		ClientID:    clientID,
		FacilityID:  facilityID,
	}).Find(&serviceRequests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get client service requests by status: %v", err)
	}
	return serviceRequests, nil
}

// GetActiveScreeningToolResponses returns all active screening tool responses that are within 24 hours of previous response
func (db *PGInstance) GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*ScreeningToolsResponse, error) {
	var responses []*ScreeningToolsResponse
	var response ScreeningToolsResponse

	err := db.DB.Scopes(OrganisationScope(ctx, response.TableName())).Where(&ScreeningToolsResponse{
		ClientID: clientID,
		Active:   true,
	}).Where("created >  ?", time.Now().Add(time.Hour*-24)).Find(&responses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses for client: %v", err)
	}
	return responses, nil
}

// CheckAppointmentExistsByExternalID checks if an appointment with the external id exists
func (db *PGInstance) CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	var appointment Appointment

	err := db.DB.Scopes(OrganisationScope(ctx, appointment.TableName())).Where(&Appointment{ExternalID: externalID}).First(&appointment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get appointment by appointment UUID: %v", err)
	}

	return true, nil
}

// GetClientScreeningToolResponsesByToolType returns all screening tool responses for a client based on the tooltype
func (db *PGInstance) GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID, toolType string, active bool) ([]*ScreeningToolsResponse, error) {
	var responses []*ScreeningToolsResponse
	var (
		screeningToolQuestion         ScreeningToolQuestion
		screeningToolQuestionResponse ScreeningToolsResponse
	)
	err := db.DB.Joins(
		"JOIN screeningtools_screeningtoolsquestion ON screeningtools_screeningtoolsquestion.id = screeningtools_screeningtoolsresponse.question_id",
	).Where(`
	    screeningtools_screeningtoolsquestion.tool_type = ?
		AND screeningtools_screeningtoolsresponse.client_id = ?
		AND screeningtools_screeningtoolsresponse.active = ?
	`, toolType, clientID, active).
		Scopes(OrganisationScope(ctx, screeningToolQuestion.TableName())).
		Scopes(OrganisationScope(ctx, screeningToolQuestionResponse.TableName())).
		Find(&responses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get responses for client: %v", err)
	}
	return responses, nil
}

// GetUserSurveyForms retrieves all user survey forms
func (db *PGInstance) GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*UserSurvey, error) {
	var userSurveys []*UserSurvey
	var userSurvey UserSurvey

	err := db.DB.Scopes(OrganisationScope(ctx, userSurvey.TableName())).Where(params).Find(&userSurveys).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user surveys: %v", err)
	}

	return userSurveys, nil
}

// GetClientScreeningToolServiceRequestByToolType returns a screening tool of type service request by based on tool type
func (db *PGInstance) GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*ClientServiceRequest, error) {
	var serviceRequest ClientServiceRequest
	err := db.DB.Scopes(OrganisationScope(ctx, serviceRequest.TableName())).Where(`
		client_id = ?
		AND meta->>'question_type' = ?
		AND request_type = ?
		AND status = ?
	`, clientID,
		toolType,
		enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
		status,
	).First(&serviceRequest).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get client service request by question ID: %v", err)
	}
	return &serviceRequest, nil
}

// CheckIfStaffHasUnresolvedServiceRequests returns true if the staff has unresolved service requests
func (db *PGInstance) CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
	var unresolvedServiceRequests []*StaffServiceRequest
	var staffServiceRequest StaffServiceRequest

	err := db.DB.Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).Where(&StaffServiceRequest{StaffID: staffID, RequestType: serviceRequestType}).
		Not(&StaffServiceRequest{Status: enums.ServiceRequestStatusResolved.String()}).
		Find(&unresolvedServiceRequests).Error
	if err != nil {
		return false, fmt.Errorf("failed to check for unresolved service requests: %v", err)
	}

	if len(unresolvedServiceRequests) > 0 {
		return true, nil
	}

	return false, nil
}

// GetAvailableScreeningTools returns all the available screening tools following the set criteria
func (db *PGInstance) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*ScreeningTool, error) {
	var screeningTools []*ScreeningTool
	organizationID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	t := time.Now().Add(time.Hour * -24)
	err = db.DB.Raw(
		`
		SELECT 
			questionnaires_screeningtool.id,  questionnaires_screeningtool.active, 
			questionnaires_screeningtool.questionnaire_id, questionnaires_screeningtool.threshold, 
			questionnaires_screeningtool.min_age, questionnaires_screeningtool.max_age,
			questionnaires_screeningtool.client_types,  questionnaires_screeningtool.genders
		FROM questionnaires_screeningtool
		JOIN clients_client
		ON clients_client.client_types && questionnaires_screeningtool.client_types
		JOIN users_user
		ON clients_client.user_id = users_user.id
		WHERE clients_client.id = ?
		AND clients_client.organisation_id = ?
		AND questionnaires_screeningtool.organisation_id = ?
		AND clients_client.current_facility_id = ?
		AND users_user.gender =  ANY (questionnaires_screeningtool.genders)
		AND DATE_PART( 'year', AGE(CURRENT_DATE, users_user.date_of_birth))::int >=  questionnaires_screeningtool.min_age
		AND DATE_PART( 'year', AGE(CURRENT_DATE, users_user.date_of_birth))::int <=  questionnaires_screeningtool.max_age
		AND questionnaires_screeningtool.id NOT IN
		(
			SELECT questionnaires_screeningtoolresponse.screeningtool_id FROM clients_servicerequest
			JOIN questionnaires_screeningtoolresponse
			ON (questionnaires_screeningtoolresponse.id)::text=(clients_servicerequest.meta->>'response_id')::text
			WHERE  clients_servicerequest.client_id = ?
			AND questionnaires_screeningtoolresponse.organisation_id = ?
			AND clients_servicerequest.organisation_id = ?
			AND clients_servicerequest.request_type = ?
			AND clients_servicerequest.status = ?
			OR questionnaires_screeningtoolresponse.created > ?
		)
		`, clientID, organizationID, organizationID, facilityID, clientID, organizationID, organizationID, enums.ServiceRequestTypeScreeningToolsRedFlag.String(), enums.ServiceRequestStatusPending.String(), t).
		Scan(&screeningTools).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get service requests for client: %w", err)
	}
	return screeningTools, nil
}

// GetClientsByFilterParams returns clients based on the filter params
// The query is constructed dynamically  based on the filterparams passed; empty filters are allowed
// facility ID is required hence it will be the first query passed, then the rest are optional
// For filter params, each check will compound to the final query that is being performed on the DB
func (db *PGInstance) GetClientsByFilterParams(ctx context.Context, facilityID string, params *dto.ClientFilterParamsInput) ([]*Client, error) {
	var (
		clients      []*Client
		filterParams dto.ClientFilterParamsInput
	)

	tx := db.DB.Where(&Client{FacilityID: facilityID})

	if params != nil {
		tx = tx.Joins("JOIN users_user on users_user.id = clients_client.user_id")
	}

	err := mapstructure.Decode(params, &filterParams)
	if err != nil {
		return nil, fmt.Errorf("failed to decode filter params: %v", err)
	}

	if len(filterParams.ClientTypes) > 0 {
		clientTypesString := fmt.Sprintf("%s", filterParams.ClientTypes)
		clientTypesString = strings.ReplaceAll(clientTypesString, "[", "{")
		clientTypesString = strings.ReplaceAll(clientTypesString, "]", "}")
		clientTypesString = strings.ReplaceAll(clientTypesString, " ", ",")

		c := Client{}
		tx = tx.Where("clients_client.client_types && ?", clientTypesString).Scopes(OrganisationScope(ctx, c.TableName()))
	}

	if filterParams.AgeRange != nil {
		lowerBoundDate := time.Now().AddDate(-filterParams.AgeRange.LowerBound, 0, 0).Format("2006-01-02")
		upperBoundDate := time.Now().AddDate(-filterParams.AgeRange.UpperBound, 0, 0).Format("2006-01-02")

		tx = tx.Where("(? > users_user.date_of_birth  AND ? < users_user.date_of_birth)", lowerBoundDate, upperBoundDate)
	}

	if len(filterParams.Gender) > 0 {
		var (
			genderString string
			genders      = filterParams.Gender
		)

		for i, gender := range genders {
			genderString += fmt.Sprintf("'%s'", strings.ToUpper(gender.String()))

			if len(genders) > 1 && i < len(genders)-1 {
				genderString += ", "
			}
		}

		tx = tx.Where(fmt.Sprintf("users_user.gender IN (%s)", genderString))
	}

	err = tx.Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get clients by filter params: %w", err)
	}

	return clients, err
}

// SearchClientServiceRequests is used to query(search) for client service requests depending on the search parameter and the type of service request passed
func (db *PGInstance) SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*ClientServiceRequest, error) {
	var clientServiceRequests []*ClientServiceRequest

	if err := db.DB.Joins("JOIN clients_client on clients_servicerequest.client_id=clients_client.id").
		Joins("JOIN users_user on clients_client.user_id=users_user.id").
		Joins("JOIN common_contact on users_user.id=common_contact.user_id").
		Where(db.DB.Or("users_user.username ILIKE ? ", "%"+searchParameter+"%").
			Or("common_contact.contact_value ILIKE ?", "%"+searchParameter+"%").
			Or("users_user.name ILIKE ? ", "%"+searchParameter+"%")).
		Where("clients_servicerequest.status = ?", enums.ServiceRequestStatusPending.String()).
		Where("clients_servicerequest.request_type = ?", requestType).
		Where("clients_servicerequest.facility_id = ?", facilityID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).
		Preload(clause.Associations).Find(&clientServiceRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to get client service requests: %w", err)
	}

	return clientServiceRequests, nil
}

// SearchStaffServiceRequests is used to query(search) for staff's service requests depending on the search parameter and the type of service request
func (db *PGInstance) SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*StaffServiceRequest, error) {
	var staffServiceRequests []*StaffServiceRequest
	var (
		staff               StaffProfile
		staffServiceRequest StaffServiceRequest
	)

	if err := db.DB.Joins("JOIN staff_staff on staff_servicerequest.staff_id=staff_staff.id").
		Joins("JOIN users_user on staff_staff.user_id=users_user.id").
		Joins("JOIN common_contact on users_user.id=common_contact.user_id").
		Where(db.DB.Scopes(OrganisationScope(ctx, staff.TableName())).
			Or("users_user.username ILIKE ? ", "%"+searchParameter+"%").Or("common_contact.contact_value ILIKE ?", "%"+searchParameter+"%").
			Or("users_user.name ILIKE ? ", "%"+searchParameter+"%")).
		Where("staff_servicerequest.status = ? ", enums.ServiceRequestStatusPending.String()).
		Where("staff_servicerequest.request_type = ?", requestType).
		Where("staff_servicerequest.facility_id = ?", facilityID).
		Scopes(OrganisationScope(ctx, staffServiceRequest.TableName())).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).
		Preload(clause.Associations).Find(&staffServiceRequests).Error; err != nil {
		return nil, fmt.Errorf("failed to get staff service requests: %w", err)
	}

	return staffServiceRequests, nil
}

// GetScreeningToolByID is used to get a screening tool by its ID
func (db *PGInstance) GetScreeningToolByID(ctx context.Context, id string) (*ScreeningTool, error) {
	var screeningTool ScreeningTool
	if err := db.DB.Scopes(OrganisationScope(ctx, screeningTool.TableName())).Where(&ScreeningTool{ID: id}).First(&screeningTool).Error; err != nil {
		return nil, fmt.Errorf("failed to get screening tool: %w", err)
	}

	return &screeningTool, nil
}

// GetQuestionnaireByID is used to get a questionnaire by its ID
func (db *PGInstance) GetQuestionnaireByID(ctx context.Context, id string) (*Questionnaire, error) {
	var questionnaire Questionnaire
	if err := db.DB.Scopes(OrganisationScope(ctx, questionnaire.TableName())).Where(&Questionnaire{ID: id}).First(&questionnaire).Error; err != nil {
		return nil, fmt.Errorf("failed to get questionnaire: %w", err)
	}

	return &questionnaire, nil
}

// GetQuestionsByQuestionnaireID is used to get questions by questionnaire ID
func (db *PGInstance) GetQuestionsByQuestionnaireID(ctx context.Context, questionnaireID string) ([]*Question, error) {
	var questions []*Question
	var q Question

	if err := db.DB.Scopes(OrganisationScope(ctx, q.TableName())).Where(&Question{QuestionnaireID: questionnaireID}).Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
}

// GetQuestionInputChoicesByQuestionID is used to get question input choices by question ID
func (db *PGInstance) GetQuestionInputChoicesByQuestionID(ctx context.Context, questionID string) ([]*QuestionInputChoice, error) {
	var questionInputChoices []*QuestionInputChoice
	var choice QuestionInputChoice
	if err := db.DB.Scopes(OrganisationScope(ctx, choice.TableName())).Where(&QuestionInputChoice{QuestionID: questionID}).Find(&questionInputChoices).Error; err != nil {
		return nil, fmt.Errorf("failed to get question input choices: %w", err)
	}

	return questionInputChoices, nil
}

// GetFacilityRespondedScreeningTools is used to get facility's responded screening tools questions
// These are screening tools that have red flag service requests and have been resolved
func (db *PGInstance) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, pagination *domain.Pagination) ([]*ScreeningTool, *domain.Pagination, error) {
	var count int64
	var screeningTools []*ScreeningTool
	var (
		questionnaire         Questionnaire
		screeningToolResponse ScreeningToolResponse
		clientServiceRequest  ClientServiceRequest
	)

	tx := db.DB.Model(&ScreeningTool{}).Joins("JOIN questionnaires_questionnaire ON questionnaires_screeningtool.questionnaire_id = questionnaires_questionnaire.id").
		Joins("JOIN questionnaires_screeningtoolresponse ON questionnaires_screeningtoolresponse.screeningtool_id = questionnaires_screeningtool.id").
		Joins("JOIN clients_servicerequest ON (questionnaires_screeningtoolresponse.id)::text=(clients_servicerequest.meta->>'response_id')::text").
		Where("questionnaires_screeningtoolresponse.facility_id = ?", facilityID).
		Where("clients_servicerequest.status = ?", enums.ServiceRequestStatusPending.String()).
		Where("clients_servicerequest.request_type = ?", enums.ServiceRequestTypeScreeningToolsRedFlag.String()).
		Scopes(OrganisationScope(ctx, questionnaire.TableName())).
		Scopes(OrganisationScope(ctx, screeningToolResponse.TableName())).
		Scopes(OrganisationScope(ctx, clientServiceRequest.TableName()))

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to get screening tools count: %w", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "questionnaires_questionnaire.name"}, Desc: true}).Find(&screeningTools).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get screening tools: %w", err)
	}

	return screeningTools, pagination, nil
}

// GetScreeningToolServiceRequestOfRespondents is used to get screening tool service request by respondents
// the clients who have a pending screening tool service request for the given facility are returned
func (db *PGInstance) GetScreeningToolServiceRequestOfRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm string, pagination *domain.Pagination) ([]*ClientServiceRequest, *domain.Pagination, error) {
	var serviceRequests []*ClientServiceRequest
	var count int64
	var (
		clients               Client
		screeningToolResponse ScreeningToolResponse
		clientsServicerequest ClientServiceRequest
		contact               Contact
	)

	tx := db.DB.Model(&ClientServiceRequest{}).Joins("JOIN questionnaires_screeningtoolresponse ON questionnaires_screeningtoolresponse.id::TEXT = clients_servicerequest.meta ->> 'response_id'::TEXT").
		Joins("JOIN clients_client ON clients_client.id = questionnaires_screeningtoolresponse.client_id").
		Joins("JOIN users_user ON clients_client.user_id = users_user.id").
		Joins("JOIN common_contact ON common_contact.user_id = users_user.id").
		Where("clients_servicerequest.request_type = ?", enums.ServiceRequestTypeScreeningToolsRedFlag.String()).
		Where("clients_servicerequest.status = ?", enums.ServiceRequestStatusPending.String()).
		Where("questionnaires_screeningtoolresponse.facility_id = ?", facilityID).
		Where("questionnaires_screeningtoolresponse.screeningtool_id = ?", screeningToolID).
		Scopes(OrganisationScope(ctx, clients.TableName())).
		Scopes(OrganisationScope(ctx, screeningToolResponse.TableName())).
		Scopes(OrganisationScope(ctx, clientsServicerequest.TableName())).
		Scopes(OrganisationScope(ctx, contact.TableName())).
		Or("common_contact.contact_value ILIKE ?", "%"+searchTerm+"%").
		Or("users_user.name ILIKE ?", "%"+searchTerm+"%")

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, err
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Find(&serviceRequests).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get screening tool serviceRequests: %w", err)
	}
	return serviceRequests, pagination, nil

}

// GetScreeningToolResponseByID is used to get a screening tool response by its ID
func (db *PGInstance) GetScreeningToolResponseByID(ctx context.Context, id string) (*ScreeningToolResponse, error) {
	var screeningToolResponse ScreeningToolResponse
	if err := db.DB.Scopes(OrganisationScope(ctx, screeningToolResponse.TableName())).Where(&ScreeningToolResponse{ID: id}).First(&screeningToolResponse).Error; err != nil {
		return nil, fmt.Errorf("failed to get screening tool response: %w", err)
	}

	return &screeningToolResponse, nil
}

// GetScreeningToolQuestionResponsesByResponseID is used to get screening tool question responses by screening tool response ID
func (db *PGInstance) GetScreeningToolQuestionResponsesByResponseID(ctx context.Context, responseID string) ([]*ScreeningToolQuestionResponse, error) {
	var screeningToolQuestionResponses []*ScreeningToolQuestionResponse
	var screeningToolQuestionResponse ScreeningToolQuestionResponse

	if err := db.DB.Scopes(OrganisationScope(ctx, screeningToolQuestionResponse.TableName())).Where(&ScreeningToolQuestionResponse{ScreeningToolResponseID: responseID}).Find(&screeningToolQuestionResponses).Error; err != nil {
		return nil, fmt.Errorf("failed to get screening tool question responses: %w", err)
	}

	return screeningToolQuestionResponses, nil
}

// GetSurveysWithServiceRequests is used to retrieve surveys with service requests for a particular facility
func (db *PGInstance) GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*UserSurvey, error) {
	var surveys []*UserSurvey

	organizationID, err := utils.GetOrganisationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := db.DB.Raw(
		`
		SELECT * FROM common_usersurveys
		JOIN clients_servicerequest
		ON (common_usersurveys.project_id)::int=(clients_servicerequest.meta->>'projectID')::int
		AND (common_usersurveys.link_id)::int=(clients_servicerequest.meta->>'submitterID')::int
		AND (common_usersurveys.form_id)::text=(clients_servicerequest.meta->>'formID')::text
		WHERE clients_servicerequest.request_type= ? 
		AND clients_servicerequest.status= ? 
		AND clients_servicerequest.facility_id= ?
		AND common_usersurveys.organisation_id= ?
		AND clients_servicerequest.organisation_id= ?
		`, enums.ServiceRequestTypeSurveyRedFlag.String(), enums.ServiceRequestStatusPending, facilityID, organizationID, organizationID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "clients_servicerequest.created"}, Desc: true}).
		Scan(&surveys).Error; err != nil {
		return nil, fmt.Errorf("failed to get surveys with service requests: %w", err)
	}

	return surveys, nil
}

// GetClientsSurveyServiceRequest retrieves a list of clients with a surveys service request
func (db *PGInstance) GetClientsSurveyServiceRequest(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*ClientServiceRequest, *domain.Pagination, error) {
	var clientsServiceRequest []*ClientServiceRequest
	var count int64

	var (
		client                Client
		clientsServicerequest ClientServiceRequest
	)

	tx := db.DB.Model(&ClientServiceRequest{}).Joins("JOIN clients_client ON clients_client.id=clients_servicerequest.client_id").
		Where("(clients_servicerequest.meta->>'projectID')::int = ? AND (clients_servicerequest.meta->>'formID')::text = ? AND clients_servicerequest.request_type = ? AND clients_servicerequest.status = ? AND clients_client.current_facility_id = ?", projectID, formID, enums.ServiceRequestTypeSurveyRedFlag, enums.ServiceRequestStatusPending, facilityID).
		Scopes(OrganisationScope(ctx, client.TableName())).
		Scopes(OrganisationScope(ctx, clientsServicerequest.TableName()))

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&clientsServiceRequest).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to execute paginated query: %v", err)
	}

	return clientsServiceRequest, pagination, nil
}

// GetStaffFacilities gets facilities belonging to a given staff
func (db *PGInstance) GetStaffFacilities(ctx context.Context, staffFacility StaffFacilities, pagination *domain.Pagination) ([]*StaffFacilities, *domain.Pagination, error) {
	var staffFacilities []*StaffFacilities
	var count int64

	tx := db.DB.Model(&staffFacilities)

	if staffFacility.StaffID != nil {
		tx = tx.Where("staff_id = ?", staffFacility.StaffID)
	}
	if staffFacility.FacilityID != nil {
		tx = tx.Where("facility_id = ?", staffFacility.FacilityID)
	}

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, err
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Find(&staffFacilities).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get staff facilities: %w", err)
	}

	return staffFacilities, pagination, nil

}

// GetClientFacilities gets facilities belonging to a given client
func (db *PGInstance) GetClientFacilities(ctx context.Context, clientFacility ClientFacilities, pagination *domain.Pagination) ([]*ClientFacilities, *domain.Pagination, error) {
	var clientFacilities []*ClientFacilities
	var count int64

	tx := db.DB.Model(&clientFacilities)

	if clientFacility.ClientID != nil {
		tx = tx.Where("client_id = ?", clientFacility.ClientID)
	}
	if clientFacility.FacilityID != nil {
		tx = tx.Where("facility_id = ?", clientFacility.FacilityID)
	}

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Find(&clientFacilities).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get client facilities: %w", err)
	}

	return clientFacilities, pagination, nil

}

// GetCaregiverProfileByCaregiverID retrieves the caregivers profile based on the user ID provided
func (db *PGInstance) GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*Caregiver, error) {
	var caregiver *Caregiver
	if err := db.DB.Scopes(OrganisationScope(ctx, caregiver.TableName())).Where(&Caregiver{ID: caregiverID}).Preload("UserProfile.Contacts").Preload(clause.Associations).First(&caregiver).Error; err != nil {
		return nil, fmt.Errorf("failed to get caregiver by user ID %v: %v", caregiverID, err)
	}

	return caregiver, nil
}

// ListClientsCaregivers retrieves a list of clients caregivers
func (db *PGInstance) ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) ([]*CaregiverClient, *domain.Pagination, error) {
	var caregiverClients []*CaregiverClient
	var count int64
	var caregiverClient CaregiverClient

	tx := db.DB.Scopes(OrganisationScope(ctx, caregiverClient.TableName())).Model(&CaregiverClient{}).Where(&CaregiverClient{ClientID: clientID})

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "created"}, Desc: true}).Find(&caregiverClients).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to execute paginated query: %v", err)
	}

	return caregiverClients, pagination, nil
}

// GetNotificationsCount fetches the total number of user's notifications in a given facility
func (db *PGInstance) GetNotificationsCount(ctx context.Context, notification Notification) (int, error) {
	var count int64
	filters := []*firebasetools.FilterParam{
		{
			FieldName:           "is_read",
			FieldType:           enumutils.FieldTypeBoolean,
			ComparisonOperation: enumutils.OperationEqual,
			FieldValue:          notification.IsRead,
		},
	}

	facilityNotificationsQuery := db.DB.Scopes(OrganisationScope(ctx, notification.TableName())).Where(&Notification{Flavour: notification.Flavour, FacilityID: notification.FacilityID})
	if err := addFilters(facilityNotificationsQuery, filters); err != nil {
		return 0, fmt.Errorf("failed to get notifications count: %w", err)
	}

	tx := db.DB.Model(&Notification{}).Or(facilityNotificationsQuery)

	if notification.UserID != nil {
		userNotificationsQuery := db.DB.Scopes(OrganisationScope(ctx, notification.TableName())).Where(&Notification{UserID: notification.UserID})
		tx.Or(userNotificationsQuery)
	}

	if err := tx.Find(&notification).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get notifications count: %w", err)
	}

	return int(count), nil
}

// FindContacts retrieves all the contacts that match the given contact type and value.
// Contacts can be shared by users thus the same contact can have multiple records stored
func (db *PGInstance) FindContacts(ctx context.Context, contactType, contactValue string) ([]*Contact, error) {
	var contacts []*Contact

	if err := db.DB.Where(Contact{Type: contactType, Value: contactValue}).Find(&contacts).Error; err != nil {
		return nil, fmt.Errorf("failed to find contacts: %w", err)
	}

	return contacts, nil
}

// GetClientsSurveyCount retrieves the total number of clients survey
func (db *PGInstance) GetClientsSurveyCount(ctx context.Context, userID string) (int, error) {
	var count int64
	var survey UserSurvey

	if err := db.DB.Scopes(OrganisationScope(ctx, survey.TableName())).Where("common_usersurveys.user_id = ? AND has_submitted = ?", userID, false).
		Find(&survey).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to execute count query: %v", err)
	}

	return int(count), nil
}

// GetCaregiversClient queries the association table for an occurrence of a caregiver's  client
func (db *PGInstance) GetCaregiversClient(ctx context.Context, caregiverClient CaregiverClient) ([]*CaregiverClient, error) {
	caregiversClients := []*CaregiverClient{}

	tx := db.DB.Scopes(OrganisationScope(ctx, caregiverClient.TableName())).Model(&caregiverClient)

	if caregiverClient.CaregiverID != "" {
		tx = tx.Where("caregiver_id = ?", caregiverClient.CaregiverID)
	}
	if caregiverClient.ClientID != "" {
		tx = tx.Where("client_id = ?", caregiverClient.ClientID)
	}

	if err := tx.Find(&caregiversClients).Error; err != nil {
		return nil, fmt.Errorf("failed to get caregivers's client: %w", err)
	}

	return caregiversClients, nil

}

// GetCaregiverManagedClients lists clients who are managed by the caregivers
// The clients should have given their consent to be managed by the caregivers
func (db *PGInstance) GetCaregiverManagedClients(ctx context.Context, caregiverID string, pagination *domain.Pagination) ([]*Client, *domain.Pagination, error) {

	var clients []*Client
	var count int64

	tx := db.DB.Model(&clients)

	tx = tx.Joins("JOIN caregivers_caregiver_client ON clients_client.id = caregivers_caregiver_client.client_id").
		Where("caregivers_caregiver_client.caregiver_id = ?", caregiverID).Where("caregivers_caregiver_client.client_consent = ?", enums.ConsentStateAccepted)

	if pagination != nil {
		if err := tx.Count(&count).Error; err != nil {
			return nil, nil, fmt.Errorf("failed to execute count query: %v", err)
		}

		pagination.Count = count
		paginateQuery(tx, pagination)
	}

	if err := tx.Find(&clients).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to get caregivers clients: %w", err)
	}

	return clients, pagination, nil
}

// CheckOrganisationExists checks if the organisation exists
func (db *PGInstance) CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error) {
	var organisation *Organisation
	err := db.DB.Where(&Organisation{ID: &organisationID}).First(&organisation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckIfProgramNameExists checks if a program exists in the organization
// the program name should be unique for each program in a given organization
func (db *PGInstance) CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error) {
	var program Program
	err := db.DB.Where(&Program{Name: programName, OrganisationID: organisationID}).First(&program).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetStaffUserPrograms retrieves all programs associated with a staff user
// a user can have multiple staff profiles for multiple programs
func (db *PGInstance) GetStaffUserPrograms(ctx context.Context, userID string) ([]*Program, error) {
	staff := []StaffProfile{}
	err := db.DB.Where(StaffProfile{UserID: userID}).Find(&staff).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user %s programs err: %w", userID, err)
	}

	var programs []*Program
	for _, staffUser := range staff {
		var program Program
		err = db.DB.Where(&Program{ID: staffUser.ProgramID}).First(&program).Error
		if err != nil {
			return nil, fmt.Errorf("failed to find program %s err: %w", staffUser.ProgramID, err)
		}

		programs = append(programs, &program)
	}

	return programs, nil
}

// GetClientUserPrograms retrieves all programs associated with a client user
// a user can have multiple client profiles for multiple programs
func (db *PGInstance) GetClientUserPrograms(ctx context.Context, userID string) ([]*Program, error) {
	clients := []Client{}
	err := db.DB.Where(Client{UserID: &userID}).Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find user %s programs err: %w", userID, err)
	}

	var programs []*Program
	for _, client := range clients {
		var program Program
		err = db.DB.Where(&Program{ID: client.ProgramID}).First(&program).Error
		if err != nil {
			return nil, fmt.Errorf("failed to find program %s err: %w", client.ProgramID, err)
		}

		programs = append(programs, &program)
	}

	return programs, nil
}

// ListOrganisations retrieves all organisations
func (db *PGInstance) ListOrganisations(ctx context.Context) ([]*Organisation, error) {
	var organisations []*Organisation

	if err := db.DB.WithContext(ctx).Find(&organisations).Error; err != nil {
		return nil, fmt.Errorf("failed to list organisations: %w", err)
	}

	return organisations, nil
}

// GetProgramFacilities gets the facilities that belong the program
func (db *PGInstance) GetProgramFacilities(ctx context.Context, programID string) ([]*ProgramFacility, error) {
	var programFacilities []*ProgramFacility

	if err := db.DB.Where(&ProgramFacility{ProgramID: programID}).Find(&programFacilities).Error; err != nil {
		return nil, err
	}

	return programFacilities, nil
}

// GetProgramByID retrieves a program by id
func (db *PGInstance) GetProgramByID(ctx context.Context, programID string) (*Program, error) {
	var program Program

	if err := db.DB.WithContext(ctx).Where("id = ?", programID).First(&program).Error; err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	return &program, nil
}
