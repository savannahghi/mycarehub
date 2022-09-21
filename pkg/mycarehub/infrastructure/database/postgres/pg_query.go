package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/scalarutils"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// SearchFacility returns a slice of healthcare facilities in the platform.
func (d *MyCareHubDb) SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.SearchFacility(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get facilities: %w", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}

	for _, m := range facilities {
		singleFacility := domain.Facility{
			ID:                 m.FacilityID,
			Name:               m.Name,
			Code:               m.Code,
			Phone:              m.Phone,
			Active:             m.Active,
			County:             m.County,
			Description:        m.Description,
			FHIROrganisationID: m.FHIROrganisationID,
		}

		facility = append(facility, &singleFacility)
	}

	return facility, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// RetrieveFacilityByMFLCode gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
	facilitySession, err := d.query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed query and retrieve facility by MFLCode: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// ListFacilities gets facilities that are filtered from search and filter,
// the results are also paginated
func (d *MyCareHubDb) ListFacilities(
	ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	// if user did not provide current page, throw an error
	if err := paginationsInput.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("pagination input validation failed: %v", err)
	}

	sortOutput := &domain.SortParam{
		Field:     paginationsInput.Sort.Field,
		Direction: paginationsInput.Sort.Direction,
	}
	paginationOutput := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       paginationsInput.Limit,
			CurrentPage: paginationsInput.CurrentPage,
			Sort:        sortOutput,
		},
	}
	filtersOutput := []*domain.FiltersParam{}
	for _, f := range filterInput {
		filter := &domain.FiltersParam{
			Name:     string(f.DataType),
			DataType: f.DataType,
			Value:    f.Value,
		}
		filtersOutput = append(filtersOutput, filter)
	}

	facilities, err := d.query.ListFacilities(ctx, searchTerm, filtersOutput, &paginationOutput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}
	return facilities, nil
}

// GetUserProfileByPhoneNumber fetches and returns a userprofile using their phonenumber
func (d *MyCareHubDb) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number should be provided")
	}

	user, err := d.query.GetUserProfileByPhoneNumber(ctx, phoneNumber, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user profile by phonenumber: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetUserPINByUserID fetches a user pin by the user ID
func (d *MyCareHubDb) GetUserPINByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id cannot be empty")
	}
	pinData, err := d.query.GetUserPINByUserID(ctx, userID, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed query and retrieve user PIN data: %s", err)
	}

	return &domain.UserPIN{
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: pinData.ValidFrom,
		ValidTo:   pinData.ValidTo,
		Flavour:   pinData.Flavour,
		IsValid:   pinData.IsValid,
		Salt:      pinData.Salt,
	}, nil
}

// GetCurrentTerms fetches the current terms service
func (d *MyCareHubDb) GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error) {
	terms, err := d.query.GetCurrentTerms(ctx, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get current terms of service: %v", err)
	}

	return &domain.TermsOfService{
		TermsID: *terms.TermsID,
		Text:    terms.Text,
	}, nil
}

// GetUserProfileByUserID fetches and returns a userprofile using their user ID
func (d *MyCareHubDb) GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID should be provided")
	}

	user, err := d.query.GetUserProfileByUserID(ctx, &userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user profile by user ID: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetSecurityQuestions fetches all the security questions
func (d *MyCareHubDb) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	var securityQuestion []*domain.SecurityQuestion

	allSecurityQuestions, err := d.query.GetSecurityQuestions(ctx, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get security questions: %v", err)
	}

	if len(allSecurityQuestions) == 0 {
		return securityQuestion, nil
	}

	for _, sq := range allSecurityQuestions {
		singleSecurityQuestion := &domain.SecurityQuestion{
			SecurityQuestionID: *sq.SecurityQuestionID,
			QuestionStem:       sq.QuestionStem,
			Description:        sq.Description,
			Flavour:            sq.Flavour,
			Active:             sq.Active,
			ResponseType:       sq.ResponseType,
		}

		securityQuestion = append(securityQuestion, singleSecurityQuestion)
	}

	return securityQuestion, nil
}

// GetSecurityQuestionByID fetches a security question by ID
func (d *MyCareHubDb) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
	securityQuestion, err := d.query.GetSecurityQuestionByID(ctx, securityQuestionID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get security question by ID: %v", err)
	}

	return &domain.SecurityQuestion{
		SecurityQuestionID: *securityQuestion.SecurityQuestionID,
		QuestionStem:       securityQuestion.QuestionStem,
		Description:        securityQuestion.Description,
		Flavour:            securityQuestion.Flavour,
		Active:             securityQuestion.Active,
		ResponseType:       securityQuestion.ResponseType,
	}, nil
}

// GetSecurityQuestionResponse returns the security question response from the database
func (d *MyCareHubDb) GetSecurityQuestionResponse(ctx context.Context, questionID string, userID string) (*domain.SecurityQuestionResponse, error) {
	if questionID == "" {
		return nil, fmt.Errorf("security question ID must be defined")
	}

	response, err := d.query.GetSecurityQuestionResponse(ctx, questionID, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.SecurityQuestionResponse{
		ResponseID: response.ResponseID,
		QuestionID: response.QuestionID,
		UserID:     response.UserID,
		Active:     response.Active,
		Response:   response.Response,
	}, nil
}

// CheckIfPhoneNumberExists checks if phone exists in the database
func (d *MyCareHubDb) CheckIfPhoneNumberExists(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
	if phone == "" {
		return false, fmt.Errorf("phone should be defined")
	}
	exists, err := d.query.CheckIfPhoneNumberExists(ctx, phone, isOptedIn, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed check whether phone exists: %s", err)
	}

	return exists, nil
}

// VerifyOTP performs the checking of OTP's existence for the specified user.
func (d *MyCareHubDb) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	if payload.PhoneNumber == "" || payload.OTP == "" {
		return false, fmt.Errorf("user ID or phone number or OTP cannot be empty")
	}
	if !payload.Flavour.IsValid() {
		return false, exceptions.InvalidFlavourDefinedError()
	}

	return d.query.VerifyOTP(ctx, payload)
}

// GetClientProfileByUserID fetched a client profile using the supplied user ID. This will be used to return the client
// details as part of the login response
func (d *MyCareHubDb) GetClientProfileByUserID(ctx context.Context, userID string) (*domain.ClientProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID must be defined")
	}

	client, err := d.query.GetClientProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var clientList []enums.ClientType
	for _, k := range client.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	facility, err := d.query.RetrieveFacility(ctx, &client.FacilityID, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	facilities, _, err := d.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: client.ID}, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
	}

	user := createMapUser(&client.User)
	return &domain.ClientProfile{
		ID:                      client.ID,
		User:                    user,
		Active:                  client.Active,
		ClientTypes:             clientList,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
		FHIRPatientID:           client.FHIRPatientID,
		HealthRecordID:          client.HealthRecordID,
		TreatmentBuddy:          client.TreatmentBuddy,
		ClientCounselled:        client.ClientCounselled,
		OrganisationID:          client.OrganisationID,
		FacilityID:              client.FacilityID,
		FacilityName:            facility.Name,
		Facilities:              facilities,
		CHVUserID:               client.CHVUserID,
	}, nil
}

// GetStaffProfileByUserID fetches the staff's profile using the user's ID and returns the staff's profile in the login response.
func (d *MyCareHubDb) GetStaffProfileByUserID(ctx context.Context, userID string) (*domain.StaffProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("staff's user ID must be defined")
	}

	staff, err := d.query.GetStaffProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get staff profile: %v", err)
	}

	staffDefaultFacility, err := d.query.RetrieveFacility(ctx, &staff.DefaultFacilityID, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get the staff facility: %v", err)
	}

	facilities, _, err := d.GetStaffFacilities(ctx, dto.StaffFacilityInput{StaffID: staff.ID}, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
	}

	user := createMapUser(&staff.UserProfile)
	return &domain.StaffProfile{
		ID:                  staff.ID,
		User:                user,
		UserID:              staff.UserID,
		Active:              staff.Active,
		StaffNumber:         staff.StaffNumber,
		Facilities:          facilities,
		DefaultFacilityID:   staff.DefaultFacilityID,
		DefaultFacilityName: staffDefaultFacility.Name,
	}, nil
}

// GetFacilityStaffs returns a list of staff at a particular facility
func (d *MyCareHubDb) GetFacilityStaffs(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error) {
	staffs, err := d.query.GetFacilityStaffs(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	staffProfiles := []*domain.StaffProfile{}
	for _, s := range staffs {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &s.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		user := createMapUser(userProfile)

		staffProfile := &domain.StaffProfile{
			ID:                s.ID,
			User:              user,
			UserID:            s.UserID,
			Active:            s.Active,
			StaffNumber:       s.StaffNumber,
			DefaultFacilityID: s.DefaultFacilityID,
		}

		staffProfiles = append(staffProfiles, staffProfile)
	}

	return staffProfiles, nil
}

// SearchStaffProfile searches for the staff profile(s) based on the passed parameter. It might be
// a username, phonenumber or staff number. It uses pattern matching and returns all values that match
// the parameter passed
func (d *MyCareHubDb) SearchStaffProfile(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	var staffProfiles []*domain.StaffProfile

	staffs, err := d.query.SearchStaffProfile(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, s := range staffs {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &s.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		user := createMapUser(userProfile)

		staffProfile := &domain.StaffProfile{
			ID:                s.ID,
			User:              user,
			UserID:            s.UserID,
			Active:            s.Active,
			StaffNumber:       s.StaffNumber,
			DefaultFacilityID: s.DefaultFacilityID,
		}

		staffProfiles = append(staffProfiles, staffProfile)
	}

	return staffProfiles, nil
}

// SearchCaregiverUser searches for the caregiver user(s) based on the passed parameter.
//
// Search parameter can be username, phonenumber or caregiver number.
func (d *MyCareHubDb) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	var caregiverProfiles []*domain.CaregiverProfile

	caregivers, err := d.query.SearchCaregiverUser(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, caregiver := range caregivers {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &caregiver.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		user := createMapUser(userProfile)

		caregiverProfile := &domain.CaregiverProfile{
			ID:              caregiver.ID,
			User:            *user,
			CaregiverNumber: caregiver.CaregiverNumber,
		}

		caregiverProfiles = append(caregiverProfiles, caregiverProfile)
	}

	return caregiverProfiles, nil
}

// CheckUserHasPin performs a look up on the pins table to check whether a user has a pin
func (d *MyCareHubDb) CheckUserHasPin(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("user ID must be defined")
	}

	exists, err := d.query.CheckUserHasPin(ctx, userID, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return exists, nil
}

// GetOTP fetches the OTP for the specified user.
func (d *MyCareHubDb) GetOTP(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTP, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number should be provided")
	}
	if !flavour.IsValid() {
		return nil, exceptions.InvalidFlavourDefinedError()
	}

	otp, err := d.query.GetOTP(ctx, phoneNumber, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get OTP: %v", err)
	}

	return &domain.OTP{
		UserID:      otp.UserID,
		OTP:         otp.OTP,
		GeneratedAt: otp.GeneratedAt,
		ValidUntil:  otp.ValidUntil,
		Flavour:     otp.Flavour,
		Valid:       otp.Valid,
	}, nil
}

// GetUserSecurityQuestionsResponses fetches all the security questions that the user has responded to
func (d *MyCareHubDb) GetUserSecurityQuestionsResponses(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID should be provided")
	}

	securityQuestionResponses, err := d.query.GetUserSecurityQuestionsResponses(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get security questions: %v", err)
	}

	if len(securityQuestionResponses) == 0 {
		return []*domain.SecurityQuestionResponse{}, nil
	}

	var securityQuestionResponse []*domain.SecurityQuestionResponse

	for _, sqr := range securityQuestionResponses {
		singleSecurityQuestionResponse := &domain.SecurityQuestionResponse{
			ResponseID: sqr.ResponseID,
			QuestionID: sqr.QuestionID,
			UserID:     sqr.UserID,
			Active:     sqr.Active,
			Response:   sqr.Response,
			IsCorrect:  sqr.IsCorrect,
		}

		securityQuestionResponse = append(securityQuestionResponse, singleSecurityQuestionResponse)
	}

	return securityQuestionResponse, nil
}

// GetContactByUserID fetches and returns a contact using their user ID
func (d *MyCareHubDb) GetContactByUserID(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
	if userID == nil {
		return nil, fmt.Errorf("user ID should be provided")
	}

	if contactType == "" {
		return nil, fmt.Errorf("contact type is required")
	}

	if contactType != "PHONE" && contactType != "EMAIL" {
		return nil, fmt.Errorf("contact type must be PHONE or EMAIL")
	}

	contact, err := d.query.GetContactByUserID(ctx, userID, contactType)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact by user ID: %v", err)
	}

	return &domain.Contact{
		ID:           contact.ContactID,
		ContactType:  contact.ContactType,
		ContactValue: contact.ContactValue,
		Active:       contact.Active,
		OptedIn:      contact.OptedIn,
	}, nil
}

// CanRecordHeathDiary is used to check if the user can record their health diary
func (d *MyCareHubDb) CanRecordHeathDiary(ctx context.Context, userID string) (bool, error) {
	canRecord, err := d.query.CanRecordHeathDiary(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return canRecord, nil
}

// GetClientHealthDiaryQuote fetches the health diary quote for the specified user
func (d *MyCareHubDb) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	var clientHealthDiaryQuotes []*domain.ClientHealthDiaryQuote
	clientHealthDiaryQuote, err := d.query.GetClientHealthDiaryQuote(ctx, limit)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch client health diary quote: %v", err)
	}
	for _, quote := range clientHealthDiaryQuote {
		clientHealthDiaryQuotes = append(clientHealthDiaryQuotes, &domain.ClientHealthDiaryQuote{
			Author: quote.Author,
			Quote:  quote.Quote,
		})
	}

	return clientHealthDiaryQuotes, nil
}

// GetPendingServiceRequestsCount gets the total number of service requests
func (d *MyCareHubDb) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
	if facilityID == "" {
		return nil, fmt.Errorf("facility ID cannot be empty")
	}
	clientsPendingServiceRequestsCount, err := d.query.GetClientsPendingServiceRequestsCount(ctx, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch clients pending service requests count: %v", err)
	}

	staffPendingServiceRequestsCount, err := d.query.GetStaffPendingServiceRequestsCount(ctx, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch staff pending service requests count: %v", err)
	}

	return &domain.ServiceRequestsCountResponse{
		ClientsServiceRequestCount: clientsPendingServiceRequestsCount,
		StaffServiceRequestCount:   staffPendingServiceRequestsCount,
	}, nil

}

// GetClientHealthDiaryEntries queries the database to return a clients all health diary records
func (d *MyCareHubDb) GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
	var healthDiaryEntries []*domain.ClientHealthDiaryEntry

	queryParams := map[string]interface{}{
		"client_id": clientID,
	}
	if moodType != nil {
		queryParams["mood"] = moodType.String()
	}
	if shared != nil {
		queryParams["share_with_health_worker"] = shared
	}

	clientHealthDiaryEntry, err := d.query.GetClientHealthDiaryEntries(ctx, queryParams)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	//Get user profile information using the client ID
	clientProfile, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, healthdiary := range clientHealthDiaryEntry {
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			ID:                    healthdiary.ClientHealthDiaryEntryID,
			Active:                healthdiary.Active,
			Mood:                  healthdiary.Mood,
			Note:                  healthdiary.Note,
			EntryType:             healthdiary.EntryType,
			ShareWithHealthWorker: healthdiary.ShareWithHealthWorker,
			SharedAt:              healthdiary.SharedAt,
			ClientID:              healthdiary.ClientID,
			CreatedAt:             healthdiary.CreatedAt,
			PhoneNumber:           clientProfile.User.Contacts.ContactValue,
			ClientName:            clientProfile.User.Name,
		}
		healthDiaryEntries = append(healthDiaryEntries, healthDiaryEntry)
	}

	return healthDiaryEntries, nil
}

// GetClientCaregiver retrieves the caregiver for the specified client
func (d *MyCareHubDb) GetClientCaregiver(ctx context.Context, caregiverID string) (*domain.Caregiver, error) {
	caregiver, err := d.query.GetClientCaregiver(ctx, caregiverID)
	if err != nil {
		return nil, err
	}

	return &domain.Caregiver{
		ID:            *caregiver.CaregiverID,
		FirstName:     caregiver.FirstName,
		LastName:      caregiver.LastName,
		PhoneNumber:   caregiver.PhoneNumber,
		CaregiverType: caregiver.CaregiverType,
	}, nil
}

// GetClientProfileByClientID retrieves the client for the specified clientID
func (d *MyCareHubDb) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	response, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	var clientList []enums.ClientType
	for _, k := range response.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}
	userProfile, err := d.query.GetUserProfileByUserID(ctx, response.UserID)
	if err != nil {
		return nil, err
	}
	user := createMapUser(userProfile)
	return &domain.ClientProfile{
		ID:                      response.ID,
		User:                    user,
		Active:                  response.Active,
		ClientTypes:             clientList,
		TreatmentEnrollmentDate: response.TreatmentEnrollmentDate,
		FHIRPatientID:           response.FHIRPatientID,
		HealthRecordID:          response.HealthRecordID,
		TreatmentBuddy:          response.TreatmentBuddy,
		ClientCounselled:        response.ClientCounselled,
		OrganisationID:          response.OrganisationID,
		FacilityID:              response.FacilityID,
		CHVUserID:               response.CHVUserID,
		CaregiverID:             response.CaregiverID,
		UserID:                  *response.UserID,
	}, nil

}

// GetServiceRequests retrieves the service requests by the type passed in the parameters
func (d *MyCareHubDb) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		clientServiceRequests, err := d.query.GetServiceRequests(ctx, requestType, requestStatus, facilityID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		serviceRequests, err := d.ReturnClientsServiceRequests(ctx, clientServiceRequests)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		return serviceRequests, nil

	case feedlib.FlavourPro:
		if facilityID == "" {
			return nil, fmt.Errorf("facility ID is required")
		}
		staffServiceRequests, err := d.query.GetStaffServiceRequests(ctx, requestType, requestStatus, facilityID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		serviceRequests, err := d.ReturnStaffServiceRequests(ctx, staffServiceRequests)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		return serviceRequests, nil

	default:
		return nil, fmt.Errorf("invalid flavour %v defined: ", flavour)
	}
}

// ReturnClientsServiceRequests returns all the clients service requests
func (d *MyCareHubDb) ReturnClientsServiceRequests(ctx context.Context, clientServiceRequests []*gorm.ClientServiceRequest) ([]*domain.ServiceRequest, error) {
	var (
		serviceRequests []*domain.ServiceRequest
		meta            map[string]interface{}
		resolvedByName  string
	)

	for _, serviceRequest := range clientServiceRequests {
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		if serviceRequest.Meta != "" {
			meta, err = utils.ConvertJSONStringToMap(serviceRequest.Meta)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}

		if serviceRequest.ResolvedByID != nil {
			resolvedBy, err := d.query.GetUserProfileByStaffID(ctx, *serviceRequest.ResolvedByID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}
			resolvedByName = resolvedBy.Name
		}

		serviceRequest := &domain.ServiceRequest{
			ID:             *serviceRequest.ID,
			RequestType:    serviceRequest.RequestType,
			Request:        serviceRequest.Request,
			Status:         serviceRequest.Status,
			ClientID:       serviceRequest.ClientID,
			CreatedAt:      serviceRequest.Base.CreatedAt,
			InProgressAt:   serviceRequest.InProgressAt,
			InProgressBy:   serviceRequest.InProgressByID,
			ResolvedAt:     serviceRequest.ResolvedAt,
			ResolvedBy:     serviceRequest.ResolvedByID,
			ResolvedByName: &resolvedByName,
			FacilityID:     serviceRequest.FacilityID,
			ClientName:     &clientProfile.User.Name,
			ClientContact:  &clientProfile.User.Contacts.ContactValue,
			Meta:           meta,
		}
		serviceRequests = append(serviceRequests, serviceRequest)
	}
	return serviceRequests, nil
}

// ReturnStaffServiceRequests returns a response of all the staffs service requests
func (d *MyCareHubDb) ReturnStaffServiceRequests(ctx context.Context, staffServiceRequests []*gorm.StaffServiceRequest) ([]*domain.ServiceRequest, error) {
	var (
		serviceRequests []*domain.ServiceRequest
		meta            map[string]interface{}
		resolvedByName  string
	)

	for _, serviceReq := range staffServiceRequests {
		staffProfile, err := d.query.GetStaffProfileByStaffID(ctx, serviceReq.StaffID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		if serviceReq.Meta != "" {
			meta, err = utils.ConvertJSONStringToMap(serviceReq.Meta)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}

		if serviceReq.ResolvedByID != nil {
			resolvedBy, err := d.query.GetUserProfileByStaffID(ctx, *serviceReq.ResolvedByID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}
			resolvedByName = resolvedBy.Name
		}

		serviceRequest := &domain.ServiceRequest{
			ID:             *serviceReq.ID,
			RequestType:    serviceReq.RequestType,
			Request:        serviceReq.Request,
			Status:         serviceReq.Status,
			StaffID:        serviceReq.StaffID,
			CreatedAt:      serviceReq.CreatedAt,
			ResolvedAt:     serviceReq.ResolvedAt,
			ResolvedBy:     serviceReq.ResolvedByID,
			ResolvedByName: &resolvedByName,
			FacilityID:     staffProfile.DefaultFacilityID,
			StaffName:      &staffProfile.UserProfile.Name,
			StaffContact:   &staffProfile.UserProfile.Contacts.ContactValue,
			Meta:           meta,
		}

		serviceRequests = append(serviceRequests, serviceRequest)

	}
	return serviceRequests, nil
}

// CheckUserRole check if a user has a role
func (d *MyCareHubDb) CheckUserRole(ctx context.Context, userID string, role string) (bool, error) {
	return d.query.CheckUserRole(ctx, userID, role)
}

// CheckUserPermission check if a user has a permission
func (d *MyCareHubDb) CheckUserPermission(ctx context.Context, userID string, permission string) (bool, error) {
	return d.query.CheckUserPermission(ctx, userID, permission)
}

// GetUserRoles retrieves the roles for the specified user
func (d *MyCareHubDb) GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
	var roles []*domain.AuthorityRole
	rolesList, err := d.query.GetUserRoles(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, role := range rolesList {
		role := &domain.AuthorityRole{
			AuthorityRoleID: *role.AuthorityRoleID,
			Name:            enums.UserRoleType(role.Name),
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetUserPermissions retrieves the permissions for the specified user
func (d *MyCareHubDb) GetUserPermissions(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
	var permissions []*domain.AuthorityPermission
	permissionsList, err := d.query.GetUserPermissions(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, permission := range permissionsList {
		permission := &domain.AuthorityPermission{
			PermissionID: *permission.AuthorityPermissionID,
			Name:         enums.PermissionType(permission.Name),
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// CheckIfUsernameExists checks whether the provided username exists
func (d *MyCareHubDb) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username must be defined")
	}

	ok, err := d.query.CheckIfUsernameExists(ctx, username)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return ok, nil
}

// GetCommunityByID fetches the community by ID
func (d *MyCareHubDb) GetCommunityByID(ctx context.Context, communityID string) (*domain.Community, error) {
	if communityID == "" {
		return nil, fmt.Errorf("communityID cannot be empty")
	}
	community, err := d.query.GetCommunityByID(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.Community{
		ID:          community.ID,
		Name:        community.Name,
		Description: community.Description,
		InviteOnly:  community.InviteOnly,
	}, nil
}

// CheckIdentifierExists checks whether an identifier of a certain type and value exists
// Used to validate uniqueness and prevent duplicates
func (d *MyCareHubDb) CheckIdentifierExists(ctx context.Context, identifierType string, identifierValue string) (bool, error) {
	return d.query.CheckIdentifierExists(ctx, identifierType, identifierValue)
}

// CheckFacilityExistsByMFLCode checks whether a facility exists using the mfl code.
// Used to validate existence of a facility
func (d *MyCareHubDb) CheckFacilityExistsByMFLCode(ctx context.Context, MFLCode int) (bool, error) {
	return d.query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
}

// GetClientsInAFacility fetches all the clients that belong to a specific facility
func (d *MyCareHubDb) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
	clientProfiles, err := d.query.GetClientsInAFacility(ctx, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch clients that belong to a facility: %v", err)
	}
	var clients []*domain.ClientProfile

	for _, cli := range clientProfiles {
		var clientList []enums.ClientType
		for _, k := range cli.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}
		user := createMapUser(&cli.User)
		client := &domain.ClientProfile{
			ID:                      cli.ID,
			User:                    user,
			Active:                  cli.Active,
			ClientTypes:             clientList,
			TreatmentEnrollmentDate: cli.TreatmentEnrollmentDate,
			FHIRPatientID:           cli.FHIRPatientID,
			HealthRecordID:          cli.HealthRecordID,
			TreatmentBuddy:          cli.TreatmentBuddy,
			ClientCounselled:        cli.ClientCounselled,
			OrganisationID:          cli.OrganisationID,
			FacilityID:              cli.FacilityID,
			CHVUserID:               cli.CHVUserID,
			CaregiverID:             cli.CaregiverID,
			UserID:                  *cli.UserID,
		}
		clients = append(clients, client)
	}
	return clients, nil
}

// GetRecentHealthDiaryEntries queries the database for health diary entries that were
// recorded after the last time the entries were synced to KenyaEMR.
func (d *MyCareHubDb) GetRecentHealthDiaryEntries(
	ctx context.Context,
	lastSyncTime time.Time,
	client *domain.ClientProfile,
) ([]*domain.ClientHealthDiaryEntry, error) {
	var healthDiaryEntries []*domain.ClientHealthDiaryEntry
	clientHealthDiaryEntry, err := d.query.GetRecentHealthDiaryEntries(ctx, lastSyncTime, *client.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	clientIdentifier, err := d.GetClientCCCIdentifier(ctx, *client.ID)
	if err != nil {
		// This should not be blocking. In an event where an identifier value is not found, is should not
		// fail and return
		helpers.ReportErrorToSentry(err)
	}

	contact, err := d.query.GetContactByUserID(ctx, &client.UserID, "PHONE")
	if err != nil {
		// This should not be blocking. In an event where an identifier value is not found, is should not
		// fail and return
		helpers.ReportErrorToSentry(err)
	}

	for _, healthdiary := range clientHealthDiaryEntry {
		if clientIdentifier == nil || contact == nil {
			continue
		}
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			ID:                    healthdiary.ClientHealthDiaryEntryID,
			Active:                healthdiary.Active,
			Mood:                  healthdiary.Mood,
			Note:                  healthdiary.Note,
			EntryType:             healthdiary.EntryType,
			ShareWithHealthWorker: healthdiary.ShareWithHealthWorker,
			SharedAt:              healthdiary.SharedAt,
			ClientID:              healthdiary.ClientID,
			CreatedAt:             healthdiary.CreatedAt,
			CCCNumber:             clientIdentifier.IdentifierValue,
			PhoneNumber:           contact.ContactValue,
			ClientName:            client.User.Name,
		}
		healthDiaryEntries = append(healthDiaryEntries, healthDiaryEntry)
	}

	return healthDiaryEntries, nil
}

// GetClientsByParams retrieves client profiles matching the provided parameters
func (d *MyCareHubDb) GetClientsByParams(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
	clients, err := d.query.GetClientsByParams(ctx, params, lastSyncTime)
	if err != nil {
		return nil, err
	}

	profiles := []*domain.ClientProfile{}
	for _, c := range clients {
		var clientList []enums.ClientType
		for _, k := range c.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}
		profiles = append(profiles, &domain.ClientProfile{
			ID:                      c.ID,
			Active:                  c.Active,
			ClientTypes:             clientList,
			UserID:                  *c.UserID,
			TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
			FHIRPatientID:           c.FHIRPatientID,
			HealthRecordID:          c.HealthRecordID,
			TreatmentBuddy:          c.TreatmentBuddy,
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			FacilityID:              c.FacilityID,
			CHVUserID:               c.CHVUserID,
			CaregiverID:             c.CaregiverID,
		})
	}

	return profiles, nil
}

// GetClientCCCIdentifier retrieves a client's ccc identifier record
func (d *MyCareHubDb) GetClientCCCIdentifier(ctx context.Context, clientID string) (*domain.Identifier, error) {
	identifier, err := d.query.GetClientCCCIdentifier(ctx, clientID)
	if err != nil {
		return nil, err
	}

	id := domain.Identifier{
		ID:                  identifier.ID,
		IdentifierType:      identifier.IdentifierType,
		IdentifierValue:     identifier.IdentifierValue,
		IdentifierUse:       identifier.IdentifierUse,
		Description:         identifier.Description,
		ValidFrom:           identifier.ValidFrom,
		ValidTo:             identifier.ValidTo,
		IsPrimaryIdentifier: identifier.IsPrimaryIdentifier,
	}

	return &id, nil
}

// GetHealthDiaryEntryByID gets the health diary entry with the given ID
func (d *MyCareHubDb) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
	healthDiaryEntry, err := d.query.GetHealthDiaryEntryByID(ctx, healthDiaryEntryID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.ClientHealthDiaryEntry{
		ID:                    healthDiaryEntry.ClientHealthDiaryEntryID,
		Active:                healthDiaryEntry.Active,
		Mood:                  healthDiaryEntry.Mood,
		Note:                  healthDiaryEntry.Note,
		EntryType:             healthDiaryEntry.EntryType,
		ShareWithHealthWorker: healthDiaryEntry.ShareWithHealthWorker,
		SharedAt:              healthDiaryEntry.SharedAt,
		ClientID:              healthDiaryEntry.ClientID,
		CreatedAt:             healthDiaryEntry.CreatedAt,
	}, nil
}

// GetServiceRequestsForKenyaEMR retrieves from the database all service requests belonging to a specific facility
func (d *MyCareHubDb) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {

	facility, err := d.query.RetrieveFacilityByMFLCode(ctx, payload.MFLCode, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	serviceRequests := []*domain.ServiceRequest{}
	allServiceRequests, err := d.query.GetServiceRequestsForKenyaEMR(ctx, *facility.FacilityID, *payload.LastSyncTime)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	for _, serviceReq := range allServiceRequests {
		var (
			screeningToolName  string
			screeningToolScore string
		)
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, serviceReq.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		clientIdentifier, err := d.GetClientCCCIdentifier(ctx, *clientProfile.ID)
		if err != nil {
			// This should not be blocking. In an event where an identifier value is not found, is should not
			// fail and return
			helpers.ReportErrorToSentry(err)
			continue
		}

		if serviceReq.Meta == "" {
			serviceReq.Meta = "{}"
		}

		meta, err := utils.ConvertJSONStringToMap(serviceReq.Meta)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		if serviceReq.RequestType == string(enums.ServiceRequestTypeScreeningToolsRedFlag) {
			screeningToolName = utils.InterfaceToString(meta["screening_tool_name"])
			screeningToolScore = utils.InterfaceToString(meta["score"])
		}

		userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		serviceRequest := &domain.ServiceRequest{
			ID:                 *serviceReq.ID,
			RequestType:        serviceReq.RequestType,
			Request:            serviceReq.Request,
			Status:             serviceReq.Status,
			ClientID:           serviceReq.ClientID,
			InProgressAt:       serviceReq.InProgressAt,
			InProgressBy:       serviceReq.InProgressByID,
			ResolvedAt:         serviceReq.ResolvedAt,
			ResolvedBy:         serviceReq.ResolvedByID,
			FacilityID:         serviceReq.FacilityID,
			ClientName:         &userProfile.Name,
			ClientContact:      &userProfile.Contacts.ContactValue,
			CCCNumber:          &clientIdentifier.IdentifierValue,
			ScreeningToolName:  screeningToolName,
			ScreeningToolScore: screeningToolScore,
		}
		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return serviceRequests, nil
}

// GetAssessmentResponses retrieves from the database all violence assessment responses belonging
func (d *MyCareHubDb) GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
	var responses []*domain.ScreeningToolAssessmentResponse
	answeredQuestions, err := d.query.GetAnsweredScreeningToolQuestions(ctx, facilityID, toolType)
	responsesMap := make(map[string]bool)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, answeredQuestion := range answeredQuestions {
		if _, ok := responsesMap[answeredQuestion.ClientID]; ok {
			continue
		}
		responsesMap[answeredQuestion.ClientID] = true

		clientProfile, err := d.query.GetClientProfileByClientID(ctx, answeredQuestion.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		responses = append(responses, &domain.ScreeningToolAssessmentResponse{
			ClientName:   userProfile.Name,
			DateAnswered: answeredQuestion.Base.CreatedAt,
			ClientID:     *clientProfile.ID,
		})
	}

	return responses, nil
}

// GetScreeningToolQuestions fetches the screening tools questions
func (d *MyCareHubDb) GetScreeningToolQuestions(ctx context.Context, questionType string) ([]*domain.ScreeningToolQuestion, error) {
	var screeningToolQuestions []*domain.ScreeningToolQuestion
	screeningToolQuestionsList, err := d.query.GetScreeningToolQuestions(ctx, questionType)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, screeningToolQuestion := range screeningToolQuestionsList {
		choices, err := utils.ConvertJSONStringToMap(screeningToolQuestion.ResponseChoices)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("error converting response choices json string to map: %v", err)
		}

		var meta map[string]interface{}

		if screeningToolQuestion.Meta != "" {
			meta, err = utils.ConvertJSONStringToMap(screeningToolQuestion.Meta)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}
		screeningToolQuestion := &domain.ScreeningToolQuestion{
			ID:               screeningToolQuestion.ID,
			Question:         screeningToolQuestion.Question,
			ToolType:         enums.ScreeningToolType(screeningToolQuestion.ToolType),
			ResponseChoices:  choices,
			ResponseType:     enums.ScreeningToolResponseType(screeningToolQuestion.ResponseType),
			ResponseCategory: enums.ScreeningToolResponseCategory(screeningToolQuestion.ResponseCategory),
			Sequence:         screeningToolQuestion.Sequence,
			Meta:             meta,
			Active:           screeningToolQuestion.Active,
		}
		screeningToolQuestions = append(screeningToolQuestions, screeningToolQuestion)
	}

	return screeningToolQuestions, nil
}

// ListAppointments lists appointments at a facility
func (d *MyCareHubDb) ListAppointments(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {

	parameters := &gorm.Appointment{
		Active:   true,
		ClientID: params.ClientID,
		Reason:   params.Reason,
		Provider: params.Provider,
	}

	appointments, pageInfo, err := d.query.ListAppointments(ctx, parameters, filters, pagination)
	if err != nil {
		return nil, nil, err
	}

	mapped := []*domain.Appointment{}
	for _, a := range appointments {
		m := &domain.Appointment{
			ID:         a.ID,
			ExternalID: a.ExternalID,
			Reason:     a.Reason,
			Provider:   a.Provider,
			Date: scalarutils.Date{
				Year:  a.Date.Year(),
				Month: int(a.Date.Month()),
				Day:   a.Date.Day(),
			},
			HasRescheduledAppointment: a.HasRescheduledAppointment,
		}

		mapped = append(mapped, m)
	}

	return mapped, pageInfo, nil
}

// ListNotifications lists notifications based on the provided parameters
func (d *MyCareHubDb) ListNotifications(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error) {

	parameters := &gorm.Notification{
		Active:     true,
		UserID:     params.UserID,
		FacilityID: params.FacilityID,
		Flavour:    params.Flavour,
	}

	notifications, pageInfo, err := d.query.ListNotifications(ctx, parameters, filters, pagination)
	if err != nil {
		return nil, nil, err
	}

	mapped := []*domain.Notification{}
	for _, a := range notifications {
		notificationType := enums.NotificationType(a.Type)
		if !notificationType.IsValid() {
			continue
		}

		m := &domain.Notification{
			ID:        a.ID,
			Title:     a.Title,
			Body:      a.Body,
			Type:      notificationType,
			IsRead:    a.IsRead,
			CreatedAt: a.CreatedAt,
		}

		mapped = append(mapped, m)
	}

	return mapped, pageInfo, nil
}

// ListSurveyRespondents lists survey respondents based on the provided parameters
func (d *MyCareHubDb) ListSurveyRespondents(ctx context.Context, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error) {
	params := map[string]interface{}{
		"project_id":    projectID,
		"has_submitted": true,
		"form_id":       formID,
	}

	respondents, pageInfo, err := d.query.ListSurveyRespondents(ctx, params, pagination)
	if err != nil {
		return nil, nil, err
	}

	mapped := []*domain.SurveyRespondent{}
	for _, a := range respondents {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &a.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, nil, err
		}

		var submittedAt *time.Time
		if a.SubmittedAt == nil {
			submittedAt = &a.Base.UpdatedAt
		} else {
			submittedAt = a.SubmittedAt
		}

		m := &domain.SurveyRespondent{
			ID:          a.ID,
			Name:        userProfile.Name,
			SubmittedAt: *submittedAt,
			ProjectID:   a.ProjectID,
			SubmitterID: a.LinkID,
			FormID:      a.FormID,
		}

		mapped = append(mapped, m)
	}

	return mapped, pageInfo, nil
}

// ListAvailableNotificationTypes retrieves the distinct notification types available for a user
func (d *MyCareHubDb) ListAvailableNotificationTypes(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error) {
	parameters := &gorm.Notification{
		Active:     true,
		UserID:     params.UserID,
		FacilityID: params.FacilityID,
		Flavour:    params.Flavour,
	}

	notificationTypes, err := d.query.ListAvailableNotificationTypes(ctx, parameters)
	if err != nil {
		return nil, err
	}

	return notificationTypes, nil
}

// GetScreeningToolQuestionByQuestionID fetches a screening tool question by question id
func (d *MyCareHubDb) GetScreeningToolQuestionByQuestionID(ctx context.Context, questionID string) (*domain.ScreeningToolQuestion, error) {
	screeningToolQuestion, err := d.query.GetScreeningToolQuestionByQuestionID(ctx, questionID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	choices, err := utils.ConvertJSONStringToMap(screeningToolQuestion.ResponseChoices)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("error converting response choices json string to map: %v", err)
	}

	var meta map[string]interface{}

	if screeningToolQuestion.Meta != "" {
		meta, err = utils.ConvertJSONStringToMap(screeningToolQuestion.Meta)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, fmt.Errorf("error converting meta json string to map: %v", err)
		}
	}
	screeningToolQuestionObj := &domain.ScreeningToolQuestion{
		ID:               screeningToolQuestion.ID,
		Question:         screeningToolQuestion.Question,
		ToolType:         enums.ScreeningToolType(screeningToolQuestion.ToolType),
		ResponseChoices:  choices,
		ResponseType:     enums.ScreeningToolResponseType(screeningToolQuestion.ResponseType),
		ResponseCategory: enums.ScreeningToolResponseCategory(screeningToolQuestion.ResponseCategory),
		Sequence:         screeningToolQuestion.Sequence,
		Meta:             meta,
		Active:           screeningToolQuestion.Active,
	}

	return screeningToolQuestionObj, nil
}

// GetClientProfileByCCCNumber fetches a client using their CCC number
func (d *MyCareHubDb) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
	clientProfile, err := d.query.GetClientProfileByCCCNumber(ctx, CCCNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	cccIdentifier, err := d.query.GetClientCCCIdentifier(ctx, *clientProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var clientList []enums.ClientType
	for _, k := range clientProfile.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	user := createMapUser(userProfile)
	return &domain.ClientProfile{
		ID:                      clientProfile.ID,
		User:                    user,
		Active:                  clientProfile.Active,
		ClientTypes:             clientList,
		UserID:                  *clientProfile.UserID,
		TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
		FHIRPatientID:           clientProfile.FHIRPatientID,
		HealthRecordID:          clientProfile.HealthRecordID,
		TreatmentBuddy:          clientProfile.TreatmentBuddy,
		ClientCounselled:        clientProfile.ClientCounselled,
		OrganisationID:          clientProfile.OrganisationID,
		FacilityID:              clientProfile.FacilityID,
		CHVUserID:               clientProfile.CHVUserID,
		CCCNumber:               cccIdentifier.IdentifierValue,
	}, nil
}

// SearchClientProfile searches for client profiles with the specified CCC number, phonenumber or username
// It returns a list of profiles that match the passed parameter
func (d *MyCareHubDb) SearchClientProfile(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
	clientProfile, err := d.query.SearchClientProfile(ctx, searchParameter)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var clients []*domain.ClientProfile

	for _, c := range clientProfile {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, c.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		user := createMapUser(userProfile)

		identifier, err := d.query.GetClientCCCIdentifier(ctx, *c.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		var clientList []enums.ClientType
		for _, k := range c.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}

		client := &domain.ClientProfile{
			ID:                      c.ID,
			User:                    user,
			Active:                  c.Active,
			ClientTypes:             clientList,
			UserID:                  *c.UserID,
			TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
			FHIRPatientID:           c.FHIRPatientID,
			HealthRecordID:          c.HealthRecordID,
			TreatmentBuddy:          c.TreatmentBuddy,
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			FacilityID:              c.FacilityID,
			CHVUserID:               c.CHVUserID,
			CCCNumber:               identifier.IdentifierValue,
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// CheckIfClientHasUnresolvedServiceRequests checks if a client has an unresolved service request
func (d *MyCareHubDb) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return d.query.CheckIfClientHasUnresolvedServiceRequests(ctx, clientID, serviceRequestType)
}

// GetAllRoles fetches all roles
func (d *MyCareHubDb) GetAllRoles(ctx context.Context) ([]*domain.AuthorityRole, error) {
	roles, err := d.query.GetAllRoles(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	mapped := []*domain.AuthorityRole{}
	for _, r := range roles {
		m := &domain.AuthorityRole{
			AuthorityRoleID: *r.AuthorityRoleID,
			Name:            enums.UserRoleType(r.Name),
			Active:          r.Active,
		}

		mapped = append(mapped, m)
	}

	return mapped, nil
}

// GetUserProfileByStaffID fetches a user profile using their staff ID
func (d *MyCareHubDb) GetUserProfileByStaffID(ctx context.Context, staffID string) (*domain.User, error) {
	userProfile, err := d.query.GetUserProfileByStaffID(ctx, staffID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	user := createMapUser(userProfile)
	return user, nil
}

// GetServiceRequestByID fetches a service request by id
func (d *MyCareHubDb) GetServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	serviceRequest, err := d.query.GetServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	metadata, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.ServiceRequest{
		ID:           *serviceRequest.ID,
		RequestType:  serviceRequest.RequestType,
		Request:      serviceRequest.Request,
		Status:       serviceRequest.Status,
		ClientID:     serviceRequest.ClientID,
		CreatedAt:    serviceRequest.CreatedAt,
		InProgressAt: serviceRequest.InProgressAt,
		InProgressBy: serviceRequest.InProgressByID,
		ResolvedAt:   serviceRequest.ResolvedAt,
		ResolvedBy:   serviceRequest.ResolvedByID,
		FacilityID:   serviceRequest.FacilityID,
		Meta:         metadata,
	}, nil
}

// GetStaffProfileByStaffID is used to retrieve staff profile using their staff ID
func (d *MyCareHubDb) GetStaffProfileByStaffID(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
	staffProfile, err := d.query.GetStaffProfileByStaffID(ctx, staffID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	user := createMapUser(&staffProfile.UserProfile)

	return &domain.StaffProfile{
		ID:                staffProfile.ID,
		User:              user,
		UserID:            staffProfile.UserID,
		Active:            staffProfile.Active,
		StaffNumber:       staffProfile.StaffNumber,
		DefaultFacilityID: staffProfile.DefaultFacilityID,
	}, nil
}

// GetAppointmentServiceRequests fetches all service requests of request type appointment given the last sync time
func (d *MyCareHubDb) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error) {
	MFLCode, err := strconv.Atoi(mflCode)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	facility, err := d.query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	serviceRequests, err := d.query.GetAppointmentServiceRequests(ctx, lastSyncTime, *facility.FacilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	appointmentServiceRequests := []domain.AppointmentServiceRequests{}
	for _, request := range serviceRequests {
		metaMap, err := utils.ConvertJSONStringToMap(request.Meta)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		var appointmentID string
		valueID, exists := metaMap["appointmentID"]
		if !exists {
			continue
		}
		appointmentID = valueID.(string)

		param := gorm.Appointment{ID: appointmentID}
		appointment, err := d.query.GetAppointment(ctx, &param)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		var rescheduleTime time.Time
		valueRescheduleTime, exists := metaMap["rescheduleTime"]
		if !exists {
			continue
		}
		rescheduleTime, err = time.Parse(time.RFC3339, valueRescheduleTime.(string))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		suggestedDate, err := utils.ConvertTimeToScalarDate(rescheduleTime)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		var inProgressByName string
		if request.InProgressByID != nil {
			inProgressBy, err := d.GetUserProfileByStaffID(ctx, *request.InProgressByID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}
			inProgressByName = inProgressBy.Name
		}

		var resolvedByName string
		if request.ResolvedByID != nil {
			resolvedBy, err := d.GetUserProfileByStaffID(ctx, *request.ResolvedByID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}
			resolvedByName = resolvedBy.Name
		}

		clientProfile, err := d.query.GetClientProfileByClientID(ctx, request.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		identifier, err := d.GetClientCCCIdentifier(ctx, request.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			continue
		}

		m := domain.AppointmentServiceRequests{
			ID:         *request.ID,
			ExternalID: appointment.ExternalID,
			Reason:     appointment.Reason,
			Date:       suggestedDate,

			Status:        request.Status,
			InProgressAt:  request.InProgressAt,
			InProgressBy:  &inProgressByName,
			ResolvedAt:    request.ResolvedAt,
			ResolvedBy:    &resolvedByName,
			ClientName:    &clientProfile.User.Name,
			ClientContact: &clientProfile.User.Contacts.ContactValue,
			CCCNumber:     identifier.IdentifierValue,
			MFLCODE:       strconv.Itoa(facility.Code),
		}

		appointmentServiceRequests = append(appointmentServiceRequests, m)
	}

	return appointmentServiceRequests, nil
}

// GetFacilitiesWithoutFHIRID fetches all facilities without FHIR Organisation ID
func (d *MyCareHubDb) GetFacilitiesWithoutFHIRID(ctx context.Context) ([]*domain.Facility, error) {
	var facilities []*domain.Facility
	results, err := d.query.GetFacilitiesWithoutFHIRID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, f := range results {
		facilities = append(facilities, &domain.Facility{
			ID:                 f.FacilityID,
			Name:               f.Name,
			Code:               f.Code,
			Phone:              f.Phone,
			Active:             f.Active,
			County:             f.County,
			Description:        f.Description,
			FHIROrganisationID: f.FHIROrganisationID,
		})
	}

	return facilities, nil
}

// GetAppointment fetches an appointment given the provided parameters
func (d *MyCareHubDb) GetAppointment(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
	parameters := &gorm.Appointment{
		ID:         params.ID,
		Active:     true,
		ClientID:   params.ClientID,
		Reason:     params.Reason,
		Provider:   params.Provider,
		ExternalID: params.ExternalID,
		FacilityID: params.FacilityID,
	}

	appointment, err := d.query.GetAppointment(ctx, parameters)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	date := appointment.Date
	appointmentDate, err := scalarutils.NewDate(date.Day(), int(date.Month()), date.Year())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	ap := &domain.Appointment{
		ID:         appointment.ID,
		ExternalID: appointment.ExternalID,
		Date:       *appointmentDate,
		Reason:     appointment.Reason,
		ClientID:   appointment.ClientID,
		FacilityID: appointment.FacilityID,
		Provider:   appointment.Provider,
	}

	return ap, nil
}

// GetClientServiceRequests fetches all client service requests generated by the system given the status
func (d *MyCareHubDb) GetClientServiceRequests(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error) {
	serviceRequests, err := d.query.GetClientServiceRequests(ctx, requestType, status, clientID, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var serviceRequestList []*domain.ServiceRequest
	for _, r := range serviceRequests {
		meta, err := utils.ConvertJSONStringToMap(r.Meta)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		serviceRequestList = append(serviceRequestList,
			&domain.ServiceRequest{
				ID:          *r.ID,
				RequestType: r.RequestType,
				Request:     r.Request,
				Status:      r.Status,
				Active:      r.Active,
				ClientID:    r.ClientID,
				CreatedAt:   r.CreatedAt,
				Meta:        meta,
			},
		)
	}

	return serviceRequestList, nil
}

// GetActiveScreeningToolResponses fetches all active screening tool responses
func (d *MyCareHubDb) GetActiveScreeningToolResponses(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error) {
	responses, err := d.query.GetActiveScreeningToolResponses(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var responseList []*domain.ScreeningToolQuestionResponse
	for _, r := range responses {
		responseList = append(responseList,
			&domain.ScreeningToolQuestionResponse{
				ID:         r.ID,
				QuestionID: r.QuestionID,
				ClientID:   r.ClientID,
				Answer:     r.Response,
				Active:     r.Active,
			},
		)
	}

	return responseList, nil
}

// CheckAppointmentExistsByExternalID checks if an appointment with the external id exists
func (d *MyCareHubDb) CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	return d.query.CheckAppointmentExistsByExternalID(ctx, externalID)
}

// GetUserSurveyForms retrives all user survey forms
func (d *MyCareHubDb) GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
	var userSurveys []*domain.UserSurvey

	surveys, err := d.query.GetUserSurveyForms(ctx, params)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, s := range surveys {
		userSurveys = append(userSurveys, &domain.UserSurvey{
			ID:           s.ID,
			Active:       s.Active,
			Created:      s.CreatedAt,
			Link:         s.Link,
			Title:        s.Title,
			Description:  s.Description,
			HasSubmitted: s.HasSubmitted,
			UserID:       s.UserID,
			Token:        s.Token,
			ProjectID:    s.ProjectID,
			FormID:       s.FormID,
			LinkID:       s.LinkID,
		})
	}

	return userSurveys, nil
}

// GetSharedHealthDiaryEntries fetches the most recent shared health diary entry
func (d *MyCareHubDb) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	clientProfile, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var healthDiaryEntries []*domain.ClientHealthDiaryEntry
	entries, err := d.query.GetSharedHealthDiaryEntries(ctx, clientID, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, k := range entries {
		healthDiaryEntries = append(healthDiaryEntries, &domain.ClientHealthDiaryEntry{
			ID:                    k.ClientHealthDiaryEntryID,
			Active:                k.Active,
			Mood:                  k.Mood,
			Note:                  k.Note,
			EntryType:             k.EntryType,
			ShareWithHealthWorker: k.ShareWithHealthWorker,
			SharedAt:              k.SharedAt,
			ClientID:              k.ClientID,
			CreatedAt:             k.CreatedAt,
			PhoneNumber:           clientProfile.User.Contacts.ContactValue,
			ClientName:            clientProfile.User.Name,
		})
	}

	return healthDiaryEntries, nil
}

// GetClientScreeningToolResponsesByToolType fetches all screening tool responses
func (d *MyCareHubDb) GetClientScreeningToolResponsesByToolType(ctx context.Context, clientID, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error) {
	responses, err := d.query.GetClientScreeningToolResponsesByToolType(ctx, clientID, toolType, active)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var responseList []*domain.ScreeningToolQuestionResponse
	for _, r := range responses {
		responseList = append(responseList,
			&domain.ScreeningToolQuestionResponse{
				ID:         r.ID,
				QuestionID: r.QuestionID,
				ClientID:   r.ClientID,
				Answer:     r.Response,
				Active:     r.Active,
			},
		)
	}

	return responseList, nil
}

// GetClientScreeningToolServiceRequestByToolType fetches a screening tool service request by tooltype, client ID and status
func (d *MyCareHubDb) GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error) {
	serviceRequest, err := d.query.GetClientScreeningToolServiceRequestByToolType(ctx, clientID, toolType, status)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	meta, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &domain.ServiceRequest{
		ID:          *serviceRequest.ID,
		RequestType: serviceRequest.RequestType,
		Request:     serviceRequest.Request,
		Status:      serviceRequest.Status,
		Active:      serviceRequest.Active,
		ClientID:    serviceRequest.ClientID,
		CreatedAt:   serviceRequest.CreatedAt,
		Meta:        meta,
	}, nil

}

// CheckIfStaffHasUnresolvedServiceRequests checks if a staff has unresolved service requests
func (d *MyCareHubDb) CheckIfStaffHasUnresolvedServiceRequests(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
	return d.query.CheckIfStaffHasUnresolvedServiceRequests(ctx, staffID, serviceRequestType)
}

// GetNotification retrieve a notification using the provided ID
func (d *MyCareHubDb) GetNotification(ctx context.Context, notificationID string) (*domain.Notification, error) {
	n, err := d.query.GetNotification(ctx, notificationID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	notification := &domain.Notification{
		ID:        n.ID,
		Title:     n.Title,
		Body:      n.Body,
		Type:      enums.NotificationType(n.Type),
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}

	return notification, nil
}

// GetClientsByFilterParams fetches clients by filter params
func (d *MyCareHubDb) GetClientsByFilterParams(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
	clients, err := d.query.GetClientsByFilterParams(ctx, *facilityID, filterParams)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var clientList []*domain.ClientProfile
	for _, c := range clients {
		var clientTypes []enums.ClientType
		for _, k := range c.ClientTypes {
			clientTypes = append(clientTypes, enums.ClientType(k))
		}
		user, err := d.query.GetUserProfileByUserID(ctx, c.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		domainUser := createMapUser(user)
		clientList = append(clientList, &domain.ClientProfile{
			ID:                      c.ID,
			Active:                  c.Active,
			ClientTypes:             clientTypes,
			UserID:                  *c.UserID,
			User:                    domainUser,
			TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
			FHIRPatientID:           c.FHIRPatientID,
			HealthRecordID:          c.HealthRecordID,
			TreatmentBuddy:          c.TreatmentBuddy,
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			FacilityID:              c.FacilityID,
			CHVUserID:               c.CHVUserID,
			CaregiverID:             c.CaregiverID,
		})
	}

	return clientList, nil
}

// SearchClientServiceRequests is used to query(search) for client service requests depending on the search parameter
func (d *MyCareHubDb) SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	serviceRequests, err := d.query.SearchClientServiceRequests(ctx, searchParameter, requestType, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return d.ReturnClientsServiceRequests(ctx, serviceRequests)
}

// SearchStaffServiceRequests is used to query(search) for staff's service requests depending on the search parameter
func (d *MyCareHubDb) SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	serviceRequests, err := d.query.SearchStaffServiceRequests(ctx, searchParameter, requestType, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return d.ReturnStaffServiceRequests(ctx, serviceRequests)
}

// GetScreeningToolByID fetches a screening tool by ID including the whole questions payload
func (d *MyCareHubDb) GetScreeningToolByID(ctx context.Context, toolID string) (*domain.ScreeningTool, error) {
	tool, err := d.query.GetScreeningToolByID(ctx, toolID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	questionnaire, err := d.query.GetQuestionnaireByID(ctx, tool.QuestionnaireID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	questionsPayload, err := d.query.GetQuestionsByQuestionnaireID(ctx, questionnaire.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	questions := []domain.Question{}

	for _, q := range questionsPayload {
		choices := []domain.QuestionInputChoice{}
		choicesPayload, err := d.query.GetQuestionInputChoicesByQuestionID(ctx, q.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}
		for _, c := range choicesPayload {
			choices = append(choices, domain.QuestionInputChoice{
				ID:         c.ID,
				Active:     c.Active,
				QuestionID: c.QuestionID,
				Choice:     c.Choice,
				Value:      c.Value,
				Score:      c.Score,
			})
		}

		questions = append(questions, domain.Question{
			ID:                q.ID,
			Active:            q.Active,
			QuestionnaireID:   q.QuestionnaireID,
			Text:              q.Text,
			QuestionType:      enums.QuestionType(q.QuestionType),
			ResponseValueType: enums.QuestionResponseValueType(q.ResponseValueType),
			Required:          q.Required,
			SelectMultiple:    q.SelectMultiple,
			Sequence:          q.Sequence,
			Choices:           choices,
		})
	}

	clientTypes := []enums.ClientType{}
	for _, k := range tool.ClientTypes {
		clientTypes = append(clientTypes, enums.ClientType(k))
	}

	genders := []enumutils.Gender{}
	for _, k := range tool.Genders {
		genders = append(genders, enumutils.Gender(k))
	}

	return &domain.ScreeningTool{
		ID:              tool.ID,
		Active:          tool.Active,
		QuestionnaireID: tool.QuestionnaireID,
		Threshold:       tool.Threshold,
		ClientTypes:     clientTypes,
		Genders:         genders,
		AgeRange: domain.AgeRange{
			LowerBound: tool.MinimumAge,
			UpperBound: tool.MaximumAge,
		},
		Questionnaire: domain.Questionnaire{
			ID:          questionnaire.ID,
			Active:      questionnaire.Active,
			Name:        questionnaire.Name,
			Description: questionnaire.Description,
			Questions:   questions,
		},
	}, nil
}

// GetAvailableScreeningTools fetches available screening tools for a client based on set criteria settings
func (d *MyCareHubDb) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
	screeningTools, err := d.query.GetAvailableScreeningTools(ctx, clientID, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var screeningToolList []*domain.ScreeningTool
	for _, s := range screeningTools {
		var clientTypes []enums.ClientType
		for _, k := range s.ClientTypes {
			clientTypes = append(clientTypes, enums.ClientType(k))
		}
		var genders []enumutils.Gender
		for _, g := range s.Genders {
			genders = append(genders, enumutils.Gender(g))
		}

		questionnaire, err := d.query.GetQuestionnaireByID(ctx, s.QuestionnaireID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		screeningToolList = append(screeningToolList, &domain.ScreeningTool{
			ID:              s.ID,
			Active:          s.Active,
			QuestionnaireID: s.QuestionnaireID,
			Threshold:       s.Threshold,
			ClientTypes:     clientTypes,
			Genders:         genders,
			AgeRange: domain.AgeRange{
				LowerBound: s.MinimumAge,
				UpperBound: s.MaximumAge,
			},
			Questionnaire: domain.Questionnaire{
				ID:          questionnaire.ID,
				Active:      questionnaire.Active,
				Name:        questionnaire.Name,
				Description: questionnaire.Description,
			},
		})
	}
	return screeningToolList, nil
}

// GetFacilityRespondedScreeningTools fetches responded screening tools for a given facility
func (d *MyCareHubDb) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error) {
	screeningTools, pageInfo, err := d.query.GetFacilityRespondedScreeningTools(ctx, facilityID, pagination)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, nil, err
	}

	var screeningToolList []*domain.ScreeningTool
	for _, s := range screeningTools {
		questionnaire, err := d.query.GetQuestionnaireByID(ctx, s.QuestionnaireID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, nil, err
		}

		screeningToolList = append(screeningToolList, &domain.ScreeningTool{
			ID:              s.ID,
			Active:          s.Active,
			QuestionnaireID: s.QuestionnaireID,
			Questionnaire: domain.Questionnaire{
				ID:          questionnaire.ID,
				Active:      questionnaire.Active,
				Name:        questionnaire.Name,
				Description: questionnaire.Description,
			},
		})
	}

	return screeningToolList, pageInfo, nil
}

// GetScreeningToolRespondents fetches the respondents for a screening tool
func (d *MyCareHubDb) GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error) {

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	serviceRequests, pageInfo, err := d.query.GetScreeningToolServiceRequestOfRespondents(ctx, facilityID, screeningToolID, searchTerm, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, nil, err
	}

	var respondents []*domain.ScreeningToolRespondent

	for _, s := range serviceRequests {
		meta, err := utils.ConvertJSONStringToMap(s.Meta)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, nil, err
		}
		responseID := meta["response_id"].(string)
		response, err := d.query.GetScreeningToolResponseByID(ctx, responseID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, nil, err
		}
		client, err := d.query.GetClientProfileByClientID(ctx, s.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, nil, err
		}
		respondent := &domain.ScreeningToolRespondent{
			ClientID:                s.ClientID,
			ScreeningToolResponseID: response.ID,
			ServiceRequestID:        *s.ID,
			ServiceRequest:          s.Request,
			Name:                    client.User.Name,
			PhoneNumber:             client.User.Contacts.ContactValue,
		}

		respondents = append(respondents, respondent)
	}

	return respondents, pageInfo, nil
}

// GetScreeningToolResponseByID fetches a screening tool response by ID
func (d *MyCareHubDb) GetScreeningToolResponseByID(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	response, err := d.query.GetScreeningToolResponseByID(ctx, id)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	screeningToolResponses, err := d.query.GetScreeningToolQuestionResponsesByResponseID(ctx, response.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	screeningTool, err := d.GetScreeningToolByID(ctx, response.ScreeningToolID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	questionResponsesPayload := []*domain.QuestionnaireScreeningToolQuestionResponse{}
	for _, s := range screeningToolResponses {
		question := screeningTool.GetQuestion(s.QuestionID)
		questionResponsesPayload = append(questionResponsesPayload, &domain.QuestionnaireScreeningToolQuestionResponse{
			ID:                      s.ID,
			Active:                  s.Active,
			ScreeningToolResponseID: id,
			QuestionID:              s.QuestionID,
			QuestionType:            question.QuestionType,
			SelectMultiple:          question.SelectMultiple,
			ResponseValueType:       question.ResponseValueType,
			Sequence:                question.Sequence,
			QuestionText:            question.Text,
			Response:                s.Response,
			NormalizedResponse:      screeningTool.GetNormalizedResponse(s.QuestionID, s.Response),
			Score:                   s.Score,
		})
	}
	return &domain.QuestionnaireScreeningToolResponse{
		ID:                response.ID,
		Active:            response.Active,
		ScreeningToolID:   response.ScreeningToolID,
		FacilityID:        response.FacilityID,
		ClientID:          response.ClientID,
		DateOfResponse:    response.CreatedAt,
		AggregateScore:    response.AggregateScore,
		QuestionResponses: questionResponsesPayload,
	}, nil
}

// GetSurveysWithServiceRequests fetches all the surveys with a service request for a given facility
func (d *MyCareHubDb) GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
	surveys, err := d.query.GetSurveysWithServiceRequests(ctx, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var surveysList []*dto.SurveysWithServiceRequest
	for _, survey := range surveys {
		surveysList = append(surveysList, &dto.SurveysWithServiceRequest{
			Title:     survey.Title,
			ProjectID: survey.ProjectID,
			LinkID:    survey.LinkID,
			FormID:    survey.FormID,
		})
	}

	// If we have similar title names from surveyList, only show one from the list
	var uniqueSurveyList []*dto.SurveysWithServiceRequest
	surveyMap := make(map[string]string)
	for _, surveyList := range surveysList {
		// If we have not seen this title before, add it to the unique list
		if _, ok := surveyMap[surveyList.Title]; !ok {
			uniqueSurveyList = append(uniqueSurveyList, surveyList)
			surveyMap[surveyList.Title] = surveyList.Title
		}
	}

	return uniqueSurveyList, nil
}

// GetSurveyServiceRequestUser returns a list of users who have a survey service request
func (d *MyCareHubDb) GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyServiceRequestUser, *domain.Pagination, error) {

	serviceReq, pageInfo, err := d.query.GetClientsSurveyServiceRequest(ctx, facilityID, projectID, formID, pagination)
	if err != nil {
		return nil, nil, err
	}

	mapped := []*domain.SurveyServiceRequestUser{}
	for _, s := range serviceReq {
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, s.ClientID)
		if err != nil {
			return nil, nil, err
		}

		var metaMap map[string]interface{}
		if s.Meta != "" {
			metaMap, err = utils.ConvertJSONStringToMap(s.Meta)
			if err != nil {
				return nil, nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}

		formID, ok := metaMap["formID"].(string)
		if !ok {
			return nil, nil, fmt.Errorf("error converting meta json %T to string", metaMap["formID"])
		}

		var projectID int
		project, ok := metaMap["projectID"].(float64)
		if !ok {
			return nil, nil, fmt.Errorf("error converting meta json %T to float64", metaMap["projectID"])
		}
		projectID = int(project)

		var submitterID int
		submitter, ok := metaMap["submitterID"].(float64)
		if !ok {
			return nil, nil, fmt.Errorf("error converting meta json %T to float64", metaMap["submitterID"])
		}
		submitterID = int(submitter)

		surveyName, ok := metaMap["surveyName"].(string)
		if !ok {
			return nil, nil, fmt.Errorf("error converting meta json %T to string", metaMap["surveyName"])
		}

		m := &domain.SurveyServiceRequestUser{
			Name:             clientProfile.User.Name,
			FormID:           formID,
			ProjectID:        projectID,
			SubmitterID:      submitterID,
			SurveyName:       surveyName,
			ServiceRequestID: *s.ID,
			PhoneNumber:      clientProfile.User.Contacts.ContactValue,
		}

		mapped = append(mapped, m)
	}

	return mapped, pageInfo, nil
}

// GetStaffFacilities gets a list of staff facilities
func (d *MyCareHubDb) GetStaffFacilities(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	facilities := []*domain.Facility{}

	staffFacility := gorm.StaffFacilities{
		StaffID:    input.StaffID,
		FacilityID: input.FacilityID,
	}

	staffFacilities, pageInfo, err := d.query.GetStaffFacilities(ctx, staffFacility, pagination)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range staffFacilities {
		facility, err := d.query.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}

		notification := &gorm.Notification{
			FacilityID: facility.FacilityID,
			Flavour:    feedlib.FlavourPro,
		}

		notificationCount, err := d.query.GetNotificationsCount(ctx, *notification)
		if err != nil {
			return nil, nil, err
		}

		staffPendingServiceRequest, err := d.query.GetClientsPendingServiceRequestsCount(ctx, *f.FacilityID)
		if err != nil {
			return nil, nil, err
		}

		facilities = append(facilities, &domain.Facility{
			ID:                 facility.FacilityID,
			Name:               facility.Name,
			Code:               facility.Code,
			Phone:              facility.Phone,
			Active:             facility.Active,
			County:             facility.County,
			Description:        facility.Description,
			FHIROrganisationID: facility.FHIROrganisationID,
			WorkStationDetails: domain.WorkStationDetails{
				Notifications:   notificationCount,
				ServiceRequests: staffPendingServiceRequest.Total,
			},
		})
	}

	return facilities, pageInfo, nil

}

// GetClientFacilities gets a list of client facilities
func (d *MyCareHubDb) GetClientFacilities(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	userProfile, err := d.query.GetClientProfileByClientID(ctx, *input.ClientID)
	if err != nil {
		return nil, nil, err
	}

	facilities := []*domain.Facility{}

	clientFacility := gorm.ClientFacilities{
		ClientID:   input.ClientID,
		FacilityID: input.FacilityID,
	}

	clientFacilities, pageInfo, err := d.query.GetClientFacilities(ctx, clientFacility, pagination)
	if err != nil {
		return nil, nil, err
	}

	for _, f := range clientFacilities {
		facility, err := d.query.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}

		notification := &gorm.Notification{
			FacilityID: facility.FacilityID,
			Flavour:    feedlib.FlavourConsumer,
		}

		notificationCount, err := d.query.GetNotificationsCount(ctx, *notification)
		if err != nil {
			return nil, nil, err
		}

		surveyCount, err := d.query.GetClientsSurveyCount(ctx, *userProfile.UserID)
		if err != nil {
			return nil, nil, err
		}

		facilities = append(facilities, &domain.Facility{
			ID:                 facility.FacilityID,
			Name:               facility.Name,
			Code:               facility.Code,
			Phone:              facility.Phone,
			Active:             facility.Active,
			County:             facility.County,
			Description:        facility.Description,
			FHIROrganisationID: facility.FHIROrganisationID,
			WorkStationDetails: domain.WorkStationDetails{
				Notifications: notificationCount,
				Surveys:       surveyCount,
			},
		})
	}

	return facilities, pageInfo, nil
}
