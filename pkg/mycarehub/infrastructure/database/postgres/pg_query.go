package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ory/fosite"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// ListFacilities returns a slice of healthcare facilities in the platform.
func (d *MyCareHubDb) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	filtersOutput := []*domain.FiltersParam{}
	for _, f := range filterInput {
		filter := &domain.FiltersParam{
			Name:     string(f.DataType),
			DataType: f.DataType,
			Value:    f.Value,
		}
		filtersOutput = append(filtersOutput, filter)
	}

	facilities, page, err := d.query.ListFacilities(ctx, searchTerm, filtersOutput, paginationsInput)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	facilitiesOutput := []*domain.Facility{}
	for _, f := range facilities {
		facility, err := d.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}
		facilitiesOutput = append(facilitiesOutput, facility)
	}
	return facilitiesOutput, page, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	identifierSession, err := d.query.RetrieveFacilityIdentifierByFacilityID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve facility identifier: %w", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession, identifierSession), nil
}

// GetOrganisation retrieves an organisation using the provided id
func (d *MyCareHubDb) GetOrganisation(ctx context.Context, id string) (*domain.Organisation, error) {
	record, err := d.query.GetOrganisation(ctx, id)
	if err != nil {
		return nil, err
	}

	programs, _, err := d.query.ListPrograms(ctx, record.ID, nil)
	if err != nil {
		return nil, err
	}

	var mappedPrograms []*domain.Program
	for _, program := range programs {
		mappedPrograms = append(mappedPrograms, &domain.Program{
			ID:                 program.ID,
			Active:             program.Active,
			Name:               program.Name,
			FHIROrganisationID: program.FHIROrganisationID,
			Description:        program.Description,
		})
	}

	return &domain.Organisation{
		ID:              *record.ID,
		Active:          record.Active,
		Code:            record.Code,
		Name:            record.Name,
		Description:     record.Description,
		EmailAddress:    record.EmailAddress,
		PhoneNumber:     record.PhoneNumber,
		PostalAddress:   record.PostalAddress,
		PhysicalAddress: record.PhysicalAddress,
		DefaultCountry:  record.DefaultCountry,
		Programs:        mappedPrograms,
	}, nil
}

// RetrieveFacilityByIdentifier gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	identifierObj := &gorm.FacilityIdentifier{
		Type:  identifier.Type.String(),
		Value: identifier.Value,
	}
	facilitySession, err := d.query.RetrieveFacilityByIdentifier(ctx, identifierObj, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve facility by identifier: %s", err)
	}

	identifierSession, err := d.query.RetrieveFacilityIdentifierByFacilityID(ctx, facilitySession.FacilityID)
	if err != nil {
		return nil, fmt.Errorf("failed retrieve facility identifier: %w", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession, identifierSession), nil
}

// ListProgramFacilities gets facilities that are filtered from search and filter,
// the results are also paginated
func (d *MyCareHubDb) ListProgramFacilities(ctx context.Context, programID, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	filtersOutput := []*domain.FiltersParam{}
	for _, f := range filterInput {
		filter := &domain.FiltersParam{
			Name:     string(f.DataType),
			DataType: f.DataType,
			Value:    f.Value,
		}
		filtersOutput = append(filtersOutput, filter)
	}

	facilities, page, err := d.query.ListProgramFacilities(ctx, programID, searchTerm, filtersOutput, paginationsInput)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	facilitiesOutput := []*domain.Facility{}
	for _, f := range facilities {
		facility, err := d.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}
		facilitiesOutput = append(facilitiesOutput, facility)
	}
	return facilitiesOutput, page, nil
}

// GetUserProfileByPhoneNumber fetches and returns a userprofile using their phonenumber
func (d *MyCareHubDb) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number should be provided")
	}

	user, err := d.query.GetUserProfileByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile by phonenumber: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetUserProfileByUsername retrieves a user using their username
func (d *MyCareHubDb) GetUserProfileByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := d.query.GetUserProfileByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile by username: %w", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetUserPINByUserID fetches a user pin by the user ID
func (d *MyCareHubDb) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id cannot be empty")
	}
	pinData, err := d.query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve user PIN data: %s", err)
	}

	return &domain.UserPIN{
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: pinData.ValidFrom,
		ValidTo:   pinData.ValidTo,
		IsValid:   pinData.IsValid,
		Salt:      pinData.Salt,
	}, nil
}

// GetCurrentTerms fetches the current terms service
func (d *MyCareHubDb) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	terms, err := d.query.GetCurrentTerms(ctx)
	if err != nil {
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
		return nil, fmt.Errorf("failed to get user profile by user ID: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetCaregiverByUserID returns the caregiver record of the provided user ID
func (d *MyCareHubDb) GetCaregiverByUserID(ctx context.Context, userID string) (*domain.Caregiver, error) {
	cgv, err := d.query.GetCaregiverByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	caregiver := &domain.Caregiver{
		ID:              cgv.ID,
		UserID:          cgv.UserID,
		CaregiverNumber: cgv.CaregiverNumber,
		Active:          cgv.Active,
	}

	return caregiver, nil
}

// GetSecurityQuestions fetches all the security questions
func (d *MyCareHubDb) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	var securityQuestion []*domain.SecurityQuestion

	allSecurityQuestions, err := d.query.GetSecurityQuestions(ctx, flavour)
	if err != nil {
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
		return false, fmt.Errorf("failed check whether phone exists: %s", err)
	}

	return exists, nil
}

// VerifyOTP performs the checking of OTP's existence for the specified user.
func (d *MyCareHubDb) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return d.query.VerifyOTP(ctx, payload)
}

// GetClientProfile fetched a client profile using the supplied user ID. This will be used to return the client
// details as part of the login response
func (d *MyCareHubDb) GetClientProfile(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID must be defined")
	}

	client, err := d.query.GetClientProfile(ctx, userID, programID)
	if err != nil {
		return nil, err
	}

	var clientList []enums.ClientType
	for _, k := range client.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	facility, err := d.RetrieveFacility(ctx, &client.FacilityID, true)
	if err != nil {
		return nil, err
	}

	identifiers, err := d.GetClientIdentifiers(ctx, *client.ID)
	if err != nil {
		return nil, err
	}

	facilities, _, err := d.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: client.ID, ProgramID: programID}, nil)
	if err != nil {
		log.Printf("failed to get client facilities: %v", err)
	}

	user, err := d.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.ClientProfile{
		ID:                      client.ID,
		User:                    user,
		Active:                  client.Active,
		ClientTypes:             clientList,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
		FHIRPatientID:           client.FHIRPatientID,
		HealthRecordID:          client.HealthRecordID,
		ClientCounselled:        client.ClientCounselled,
		OrganisationID:          client.OrganisationID,
		ProgramID:               client.ProgramID,
		DefaultFacility:         facility,
		Facilities:              facilities,
		Identifiers:             identifiers,
	}, nil
}

// GetStaffProfile fetches the staff's profile using the user's ID and returns the staff's profile in the login response.
func (d *MyCareHubDb) GetStaffProfile(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("staff's user ID must be defined")
	}

	staff, err := d.query.GetStaffProfile(ctx, userID, programID)
	if err != nil {
		return nil, err
	}
	facilities, _, err := d.GetStaffFacilities(ctx, dto.StaffFacilityInput{StaffID: staff.ID, ProgramID: programID}, nil)
	if err != nil {
		log.Printf("unable to get staff facilities: %v", err)
	}

	facility, err := d.RetrieveFacility(ctx, &staff.DefaultFacilityID, true)
	if err != nil {
		return nil, err
	}
	user, err := d.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	nationalIDIdentifier := enums.UserIdentifierTypeNationalID.String()

	identifiersObj, err := d.query.GetStaffIdentifiers(ctx, *staff.ID, &nationalIDIdentifier)
	if err != nil {
		return nil, err
	}

	var identifiers []*domain.Identifier

	for _, identifier := range identifiersObj {
		identifiers = append(identifiers, &domain.Identifier{
			ID:                  identifier.ID,
			Type:                enums.UserIdentifierType(identifier.Type),
			Value:               identifier.Value,
			Use:                 identifier.Use,
			Description:         identifier.Description,
			ValidFrom:           identifier.ValidFrom,
			ValidTo:             identifier.ValidTo,
			IsPrimaryIdentifier: identifier.IsPrimaryIdentifier,
			Active:              identifier.Active,
			ProgramID:           identifier.ProgramID,
			OrganisationID:      identifier.OrganisationID,
		})
	}

	return &domain.StaffProfile{
		ID:                  staff.ID,
		User:                user,
		UserID:              staff.UserID,
		Active:              staff.Active,
		StaffNumber:         staff.StaffNumber,
		Facilities:          facilities,
		ProgramID:           staff.ProgramID,
		DefaultFacility:     facility,
		IsOrganisationAdmin: staff.IsOrganisationAdmin,
		Identifiers:         identifiers,
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
			return nil, err
		}
		user := createMapUser(userProfile)

		facility, err := d.RetrieveFacility(ctx, &s.DefaultFacilityID, true)
		if err != nil {
			return nil, err
		}

		staffProfile := &domain.StaffProfile{
			ID:                  s.ID,
			User:                user,
			UserID:              s.UserID,
			Active:              s.Active,
			StaffNumber:         s.StaffNumber,
			DefaultFacility:     facility,
			IsOrganisationAdmin: s.IsOrganisationAdmin,
		}

		staffProfiles = append(staffProfiles, staffProfile)
	}

	return staffProfiles, nil
}

// SearchStaffProfile searches for the staff profile(s) based on the passed parameter. It might be
// a username, phonenumber or staff number. It uses pattern matching and returns all values that match
// the parameter passed
func (d *MyCareHubDb) SearchStaffProfile(ctx context.Context, searchParameter string, programID *string) ([]*domain.StaffProfile, error) {
	var staffProfiles []*domain.StaffProfile

	staffs, err := d.query.SearchStaffProfile(ctx, searchParameter, programID)
	if err != nil {
		return nil, err
	}

	for _, s := range staffs {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &s.UserID)
		if err != nil {
			return nil, err
		}
		user := createMapUser(userProfile)

		facility, err := d.RetrieveFacility(ctx, &s.DefaultFacilityID, true)
		if err != nil {
			return nil, err
		}

		staffProfile := &domain.StaffProfile{
			ID:                  s.ID,
			User:                user,
			UserID:              s.UserID,
			Active:              s.Active,
			StaffNumber:         s.StaffNumber,
			DefaultFacility:     facility,
			IsOrganisationAdmin: s.IsOrganisationAdmin,
		}

		staffProfiles = append(staffProfiles, staffProfile)
	}

	return staffProfiles, nil
}

// SearchCaregiverUser searches for the caregiver user(s) based on the passed parameter.
// Search parameter can be username, phonenumber or caregiver number.
// the results are scoped to the program of the healthcare worker
func (d *MyCareHubDb) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	caregiverProfiles := []*domain.CaregiverProfile{}

	caregivers, err := d.query.SearchCaregiverUser(ctx, searchParameter)
	if err != nil {
		return nil, err
	}

	for _, caregiver := range caregivers {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &caregiver.UserID)
		if err != nil {
			return nil, err
		}
		user := createMapUser(userProfile)

		clientProfile, err := d.query.GetClientProfile(ctx, *userProfile.UserID, "")
		if err != nil {
			// Do not lock the search if no client profile is found since we are only using the response to know if the caregiver is a client
			log.Printf("unable to get client profile: %v", err)
		}

		var isClient bool
		if clientProfile != nil {
			isClient = true
		}

		caregiverProfile := &domain.CaregiverProfile{
			ID:              caregiver.ID,
			User:            *user,
			CaregiverNumber: caregiver.CaregiverNumber,
			IsClient:        isClient,
			Consent:         domain.ConsentStatus{},
			CurrentClient:   caregiver.CurrentClient,
			CurrentFacility: caregiver.CurrentFacility,
		}

		caregiverProfiles = append(caregiverProfiles, caregiverProfile)
	}

	return caregiverProfiles, nil
}

// SearchPlatformCaregivers searches for the caregiver user(s) based on the passed parameter.
// Search parameter can be username, phonenumber or caregiver number.
// the results are scoped to the whole platform
func (d *MyCareHubDb) SearchPlatformCaregivers(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	caregiverProfiles := []*domain.CaregiverProfile{}

	caregivers, err := d.query.SearchPlatformCaregivers(ctx, searchParameter)
	if err != nil {
		return nil, err
	}

	for _, caregiver := range caregivers {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &caregiver.UserID)
		if err != nil {
			return nil, err
		}
		user := createMapUser(userProfile)

		clientProfile, err := d.query.GetClientProfile(ctx, *userProfile.UserID, "")
		if err != nil {
			// Do not lock the search if no client profile is found since we are only using the response to know if the caregiver is a client
			log.Printf("unable to get client profile: %v", err)
		}

		var isClient bool
		if clientProfile != nil {
			isClient = true
		}

		caregiverProfile := &domain.CaregiverProfile{
			ID:              caregiver.ID,
			User:            *user,
			CaregiverNumber: caregiver.CaregiverNumber,
			IsClient:        isClient,
			Consent:         domain.ConsentStatus{},
			CurrentClient:   caregiver.CurrentClient,
			CurrentFacility: caregiver.CurrentFacility,
		}

		caregiverProfiles = append(caregiverProfiles, caregiverProfile)
	}

	return caregiverProfiles, nil
}

// CheckUserHasPin performs a look up on the pins table to check whether a user has a pin
func (d *MyCareHubDb) CheckUserHasPin(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("user ID must be defined")
	}

	exists, err := d.query.CheckUserHasPin(ctx, userID)
	if err != nil {
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
		return nil, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("invalid flavour defined"))
	}

	otp, err := d.query.GetOTP(ctx, phoneNumber, flavour)
	if err != nil {
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
func (d *MyCareHubDb) GetUserSecurityQuestionsResponses(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID should be provided")
	}

	securityQuestionResponses, err := d.query.GetUserSecurityQuestionsResponses(ctx, userID, flavour)
	if err != nil {
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
		ID:           &contact.ID,
		ContactType:  contact.Type,
		ContactValue: contact.Value,
		Active:       contact.Active,
		OptedIn:      contact.OptedIn,
	}, nil
}

// CanRecordHeathDiary is used to check if the user can record their health diary
func (d *MyCareHubDb) CanRecordHeathDiary(ctx context.Context, userID string) (bool, error) {
	canRecord, err := d.query.CanRecordHeathDiary(ctx, userID)
	if err != nil {
		return false, err
	}

	return canRecord, nil
}

// GetClientHealthDiaryQuote fetches the health diary quote for the specified user
func (d *MyCareHubDb) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	var clientHealthDiaryQuotes []*domain.ClientHealthDiaryQuote
	clientHealthDiaryQuote, err := d.query.GetClientHealthDiaryQuote(ctx, limit)
	if err != nil {
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
func (d *MyCareHubDb) GetPendingServiceRequestsCount(ctx context.Context, facilityID string, programID string) (*domain.ServiceRequestsCountResponse, error) {
	if facilityID == "" {
		return nil, fmt.Errorf("facility ID cannot be empty")
	}
	clientsPendingServiceRequestsCount, err := d.query.GetClientsPendingServiceRequestsCount(ctx, facilityID, &programID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clients pending service requests count: %v", err)
	}

	staffPendingServiceRequestsCount, err := d.query.GetStaffPendingServiceRequestsCount(ctx, facilityID, programID)
	if err != nil {
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
		return nil, err
	}

	//Get user profile information using the client ID
	clientProfile, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
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
			PhoneNumber:           clientProfile.User.Contacts.Value,
			ClientName:            clientProfile.User.Name,
			CaregiverID:           healthdiary.CaregiverID,
		}
		healthDiaryEntries = append(healthDiaryEntries, healthDiaryEntry)
	}

	return healthDiaryEntries, nil
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

	facility, err := d.RetrieveFacility(ctx, &response.FacilityID, true)
	if err != nil {
		return nil, err
	}

	identifiers, err := d.GetClientIdentifiers(ctx, clientID)
	if err != nil {
		return nil, err
	}

	program, err := d.GetProgramByID(ctx, response.ProgramID)
	if err != nil {
		return nil, err
	}

	organisation, err := d.GetOrganisation(ctx, response.OrganisationID)
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
		ClientCounselled:        response.ClientCounselled,
		OrganisationID:          response.OrganisationID,
		ProgramID:               response.ProgramID,
		DefaultFacility:         facility,
		UserID:                  *response.UserID,
		Identifiers:             identifiers,
		Program:                 program,
		Organisation:            organisation,
	}, nil

}

// GetServiceRequests retrieves the service requests by the type passed in the parameters
func (d *MyCareHubDb) GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, programID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		clientServiceRequests, err := d.query.GetServiceRequests(ctx, requestType, requestStatus, facilityID, programID)
		if err != nil {
			return nil, err
		}

		serviceRequests, err := d.ReturnClientsServiceRequests(ctx, clientServiceRequests)
		if err != nil {
			return nil, err
		}

		return serviceRequests, nil

	case feedlib.FlavourPro:
		staffServiceRequests, err := d.query.GetStaffServiceRequests(ctx, requestType, requestStatus, facilityID)
		if err != nil {
			return nil, err
		}

		serviceRequests, err := d.ReturnStaffServiceRequests(ctx, staffServiceRequests)
		if err != nil {
			return nil, err
		}
		return serviceRequests, nil

	default:
		return nil, fmt.Errorf("invalid flavour %v defined: ", flavour)
	}
}

// ReturnClientsServiceRequests returns all the clients service requests
func (d *MyCareHubDb) ReturnClientsServiceRequests(ctx context.Context, clientServiceRequests []*gorm.ClientServiceRequest) ([]*domain.ServiceRequest, error) {
	var serviceRequests []*domain.ServiceRequest

	for _, serviceRequest := range clientServiceRequests {
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
		if err != nil {
			return nil, err
		}

		var caregiverName, caregiverContact, caregiverID string

		if serviceRequest.CaregiverID != nil {
			caregiverProfile, err := d.GetCaregiverProfileByCaregiverID(ctx, *serviceRequest.CaregiverID)
			if err != nil {
				return nil, err
			}
			caregiverID = caregiverProfile.ID
			caregiverName = caregiverProfile.User.Name
			caregiverContact = caregiverProfile.User.Contacts.ContactValue
		}
		var meta map[string]interface{}
		if serviceRequest.Meta != "" {
			meta, err = utils.ConvertJSONStringToMap(serviceRequest.Meta)
			if err != nil {
				return nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}
		var resolvedByName string
		if serviceRequest.ResolvedByID != nil {
			resolvedBy, err := d.query.GetUserProfileByStaffID(ctx, *serviceRequest.ResolvedByID)
			if err != nil {
				return nil, err
			}
			resolvedByName = resolvedBy.Name
		}

		clientServiceRequest := &domain.ServiceRequest{
			ID:               *serviceRequest.ID,
			RequestType:      serviceRequest.RequestType,
			Request:          serviceRequest.Request,
			Status:           serviceRequest.Status,
			ClientID:         serviceRequest.ClientID,
			CreatedAt:        serviceRequest.Base.CreatedAt,
			InProgressAt:     serviceRequest.InProgressAt,
			InProgressBy:     serviceRequest.InProgressByID,
			ResolvedAt:       serviceRequest.ResolvedAt,
			ResolvedBy:       serviceRequest.ResolvedByID,
			ResolvedByName:   &resolvedByName,
			FacilityID:       serviceRequest.FacilityID,
			ClientName:       &clientProfile.User.Name,
			ClientContact:    &clientProfile.User.Contacts.Value,
			Meta:             meta,
			CaregiverID:      caregiverID,
			CaregiverName:    &caregiverName,
			CaregiverContact: &caregiverContact,
		}
		if serviceRequest.CaregiverID != nil {
			clientServiceRequest.CaregiverID = *serviceRequest.CaregiverID
		}

		serviceRequests = append(serviceRequests, clientServiceRequest)
	}
	return serviceRequests, nil
}

// ReturnStaffServiceRequests returns a response of all the staffs service requests
func (d *MyCareHubDb) ReturnStaffServiceRequests(ctx context.Context, staffServiceRequests []*gorm.StaffServiceRequest) ([]*domain.ServiceRequest, error) {
	var serviceRequests []*domain.ServiceRequest

	for _, serviceReq := range staffServiceRequests {
		staffProfile, err := d.query.GetStaffProfileByStaffID(ctx, serviceReq.StaffID)
		if err != nil {
			return nil, err
		}
		var meta map[string]interface{}
		if serviceReq.Meta != "" {
			meta, err = utils.ConvertJSONStringToMap(serviceReq.Meta)
			if err != nil {
				return nil, fmt.Errorf("error converting meta json string to map: %v", err)
			}
		}

		var resolvedByName string
		if serviceReq.ResolvedByID != nil {
			resolvedBy, err := d.query.GetUserProfileByStaffID(ctx, *serviceReq.ResolvedByID)
			if err != nil {
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
			StaffContact:   &staffProfile.UserProfile.Contacts.Value,
			Meta:           meta,
		}

		serviceRequests = append(serviceRequests, serviceRequest)

	}
	return serviceRequests, nil
}

// CheckStaffExists checks if there is a staff profile that exists for a user
func (d *MyCareHubDb) CheckStaffExists(ctx context.Context, userID string) (bool, error) {
	return d.query.CheckStaffExists(ctx, userID)
}

// CheckClientExists checks if there is a client profile that exists for a user
func (d *MyCareHubDb) CheckClientExists(ctx context.Context, userID string) (bool, error) {
	return d.query.CheckClientExists(ctx, userID)
}

// CheckCaregiverExists checks if there is a caregiver profile that exists for a user
func (d *MyCareHubDb) CheckCaregiverExists(ctx context.Context, userID string) (bool, error) {
	return d.query.CheckCaregiverExists(ctx, userID)
}

// CheckIfUsernameExists checks whether the provided username exists
func (d *MyCareHubDb) CheckIfUsernameExists(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username must be defined")
	}

	ok, err := d.query.CheckIfUsernameExists(ctx, username)
	if err != nil {
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
		return nil, err
	}

	return &domain.Community{
		ID:          community.ID,
		Name:        community.Name,
		Description: community.Description,
	}, nil
}

// CheckIdentifierExists checks whether an identifier of a certain type and value exists
// Used to validate uniqueness and prevent duplicates
func (d *MyCareHubDb) CheckIdentifierExists(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
	return d.query.CheckIdentifierExists(ctx, identifierType.String(), identifierValue)
}

// CheckFacilityExistsByIdentifier checks whether a facility exists using the mfl code.
// Used to validate existence of a facility
func (d *MyCareHubDb) CheckFacilityExistsByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	identifierObj := &gorm.FacilityIdentifier{
		Type:  identifier.Type.String(),
		Value: identifier.Value,
	}
	return d.query.CheckFacilityExistsByIdentifier(ctx, identifierObj)
}

// GetClientsInAFacility fetches all the clients that belong to a specific facility
func (d *MyCareHubDb) GetClientsInAFacility(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
	clientProfiles, err := d.query.GetClientsInAFacility(ctx, facilityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clients that belong to a facility: %v", err)
	}
	var clients []*domain.ClientProfile
	for _, client := range clientProfiles {
		var clientList []enums.ClientType
		for _, k := range client.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}
		facility, err := d.RetrieveFacility(ctx, &client.FacilityID, true)
		if err != nil {
			return nil, err
		}
		user := createMapUser(&client.User)
		clients = append(clients, &domain.ClientProfile{
			ID:                      client.ID,
			User:                    user,
			Active:                  client.Active,
			ClientTypes:             clientList,
			TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
			FHIRPatientID:           client.FHIRPatientID,
			HealthRecordID:          client.HealthRecordID,
			ClientCounselled:        client.ClientCounselled,
			OrganisationID:          client.OrganisationID,
			DefaultFacility:         facility,
			UserID:                  *client.UserID,
		})
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
		return nil, err
	}

	clientIdentifier, err := d.GetClientIdentifiers(ctx, *client.ID)
	if err != nil {
		log.Printf("failed to get client CCC identifier: %v", err)
		// This should not be blocking. In an event where an identifier value is not found, is should not
		// fail and return
	}

	var identifierValue string
	for _, identifier := range clientIdentifier {
		if identifier.Type == enums.UserIdentifierTypeCCC {
			identifierValue = identifier.Value
		}
	}

	contact, err := d.query.GetContactByUserID(ctx, &client.UserID, "PHONE")
	if err != nil {
		log.Printf("failed to get contact for user: %v", err)
		// This should not be blocking. In an event where an identifier value is not found, is should not
		// fail and return
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
			CCCNumber:             identifierValue,
			PhoneNumber:           contact.Value,
			ClientName:            client.User.Name,
			CaregiverID:           healthdiary.CaregiverID,
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
		facility, err := d.RetrieveFacility(ctx, &c.FacilityID, true)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, &domain.ClientProfile{
			ID:                      c.ID,
			Active:                  c.Active,
			ClientTypes:             clientList,
			UserID:                  *c.UserID,
			TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
			FHIRPatientID:           c.FHIRPatientID,
			HealthRecordID:          c.HealthRecordID,
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			DefaultFacility:         facility,
		})
	}

	return profiles, nil
}

// GetClientIdentifiers retrieves a client's ccc identifier record
func (d *MyCareHubDb) GetClientIdentifiers(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
	identifiersObj, err := d.query.GetClientIdentifiers(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var identifiers []*domain.Identifier
	for _, identifier := range identifiersObj {
		identifiers = append(identifiers, &domain.Identifier{
			ID:                  identifier.ID,
			Type:                enums.UserIdentifierType(identifier.Type),
			Value:               identifier.Value,
			Use:                 identifier.Use,
			Description:         identifier.Description,
			ValidFrom:           identifier.ValidFrom,
			ValidTo:             identifier.ValidTo,
			IsPrimaryIdentifier: identifier.IsPrimaryIdentifier,
			Active:              identifier.Active,
			ProgramID:           identifier.ProgramID,
			OrganisationID:      identifier.OrganisationID,
		})
	}
	return identifiers, nil
}

// GetHealthDiaryEntryByID gets the health diary entry with the given ID
func (d *MyCareHubDb) GetHealthDiaryEntryByID(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
	healthDiaryEntry, err := d.query.GetHealthDiaryEntryByID(ctx, healthDiaryEntryID)
	if err != nil {
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
		CaregiverID:           healthDiaryEntry.CaregiverID,
	}, nil
}

// GetServiceRequestsForKenyaEMR retrieves from the database all service requests belonging to a specific facility
func (d *MyCareHubDb) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {

	mflIdentifier := &gorm.FacilityIdentifier{
		Type:  enums.FacilityIdentifierTypeMFLCode.String(),
		Value: strconv.Itoa(payload.MFLCode),
	}

	facility, err := d.query.RetrieveFacilityByIdentifier(ctx, mflIdentifier, true)
	if err != nil {
		return nil, err
	}

	serviceRequests := []*domain.ServiceRequest{}
	allServiceRequests, err := d.query.GetServiceRequestsForKenyaEMR(ctx, *facility.FacilityID, *payload.LastSyncTime)
	if err != nil {
		return nil, err
	}
	for _, serviceReq := range allServiceRequests {
		var (
			screeningToolName  string
			screeningToolScore string
		)
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, serviceReq.ClientID)
		if err != nil {
			return nil, err
		}

		clientIdentifier, err := d.GetClientIdentifiers(ctx, *clientProfile.ID)
		if err != nil {
			// This should not be blocking. In an event where an identifier value is not found, is should not
			// fail and return
			continue
		}

		var identifierValue string
		for _, identifier := range clientIdentifier {
			if identifier.Type == enums.UserIdentifierTypeCCC {
				identifierValue = identifier.Value
			}
		}

		if serviceReq.Meta == "" {
			serviceReq.Meta = "{}"
		}

		meta, err := utils.ConvertJSONStringToMap(serviceReq.Meta)
		if err != nil {
			return nil, err
		}
		if serviceReq.RequestType == string(enums.ServiceRequestTypeScreeningToolsRedFlag) {
			screeningToolName = utils.InterfaceToString(meta["screening_tool_name"])
			score := utils.InterfaceToFloat64(meta["score"])
			screeningToolScore = strconv.FormatFloat(score, 'f', 2, 64)
		}

		userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
		if err != nil {
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
			ClientContact:      &userProfile.Contacts.Value,
			CCCNumber:          &identifierValue,
			ScreeningToolName:  screeningToolName,
			ScreeningToolScore: screeningToolScore,
		}
		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return serviceRequests, nil
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
func (d *MyCareHubDb) ListSurveyRespondents(ctx context.Context, params *domain.UserSurvey, facilityID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error) {
	userSurveyParams := &gorm.UserSurvey{
		HasSubmitted: params.HasSubmitted,
		FormID:       params.FormID,
		ProjectID:    params.ProjectID,
		ProgramID:    params.ProgramID,
	}

	respondents, pageInfo, err := d.query.ListSurveyRespondents(ctx, userSurveyParams, facilityID, pagination)
	if err != nil {
		return nil, nil, err
	}

	mapped := []*domain.SurveyRespondent{}
	for _, a := range respondents {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, &a.UserID)
		if err != nil {
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
			CaregiverID: a.CaregiverID,
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
		ProgramID:  params.ProgramID,
	}

	notificationTypes, err := d.query.ListAvailableNotificationTypes(ctx, parameters)
	if err != nil {
		return nil, err
	}

	return notificationTypes, nil
}

// GetProgramClientProfileByIdentifier fetches a client using their CCC number
func (d *MyCareHubDb) GetProgramClientProfileByIdentifier(ctx context.Context, programID, identifierType, value string) (*domain.ClientProfile, error) {
	clientProfile, err := d.query.GetProgramClientProfileByIdentifier(ctx, programID, identifierType, value)
	if err != nil {
		return nil, err
	}

	userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
	if err != nil {
		return nil, err
	}

	identifiers, err := d.GetClientIdentifiers(ctx, *clientProfile.ID)
	if err != nil {
		return nil, err
	}

	var clientList []enums.ClientType
	for _, k := range clientProfile.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	facility, err := d.RetrieveFacility(ctx, &clientProfile.FacilityID, true)
	if err != nil {
		return nil, err
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
		ClientCounselled:        clientProfile.ClientCounselled,
		OrganisationID:          clientProfile.OrganisationID,
		ProgramID:               clientProfile.ProgramID,
		DefaultFacility:         facility,
		Identifiers:             identifiers,
	}, nil
}

// GetClientProfilesByIdentifier fetches all client profiles that match the given identifier
func (d *MyCareHubDb) GetClientProfilesByIdentifier(ctx context.Context, identifierType, value string) ([]*domain.ClientProfile, error) {
	clientProfilesObject, err := d.query.GetClientProfilesByIdentifier(ctx, identifierType, value)
	if err != nil {
		return nil, err
	}
	var clientProfiles []*domain.ClientProfile

	for _, clientProfile := range clientProfilesObject {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
		if err != nil {
			return nil, err
		}

		identifiers, err := d.GetClientIdentifiers(ctx, *clientProfile.ID)
		if err != nil {
			return nil, err
		}

		var clientList []enums.ClientType
		for _, k := range clientProfile.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}

		facility, err := d.RetrieveFacility(ctx, &clientProfile.FacilityID, true)
		if err != nil {
			return nil, err
		}

		user := createMapUser(userProfile)
		clientProfiles = append(clientProfiles, &domain.ClientProfile{
			ID:                      clientProfile.ID,
			User:                    user,
			Active:                  clientProfile.Active,
			ClientTypes:             clientList,
			UserID:                  *clientProfile.UserID,
			TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
			FHIRPatientID:           clientProfile.FHIRPatientID,
			HealthRecordID:          clientProfile.HealthRecordID,
			ClientCounselled:        clientProfile.ClientCounselled,
			OrganisationID:          clientProfile.OrganisationID,
			ProgramID:               clientProfile.ProgramID,
			DefaultFacility:         facility,
			Identifiers:             identifiers,
		},
		)
	}

	return clientProfiles, nil
}

// SearchClientProfile searches for client profiles with the specified CCC number, phonenumber or username
// It returns a list of profiles that match the passed parameter
func (d *MyCareHubDb) SearchClientProfile(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
	clientProfile, err := d.query.SearchClientProfile(ctx, searchParameter)
	if err != nil {
		return nil, err
	}

	var clients []*domain.ClientProfile

	for _, c := range clientProfile {
		userProfile, err := d.query.GetUserProfileByUserID(ctx, c.UserID)
		if err != nil {
			return nil, err
		}
		user := createMapUser(userProfile)

		facility, err := d.RetrieveFacility(ctx, &c.FacilityID, true)
		if err != nil {
			return nil, err
		}

		identifiers, err := d.GetClientIdentifiers(ctx, *c.ID)
		if err != nil {
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
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			DefaultFacility:         facility,
			Identifiers:             identifiers,
		}

		clients = append(clients, client)
	}

	return clients, nil
}

// CheckIfClientHasUnresolvedServiceRequests checks if a client has an unresolved service request
func (d *MyCareHubDb) CheckIfClientHasUnresolvedServiceRequests(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
	return d.query.CheckIfClientHasUnresolvedServiceRequests(ctx, clientID, serviceRequestType)
}

// GetUserProfileByStaffID fetches a user profile using their staff ID
func (d *MyCareHubDb) GetUserProfileByStaffID(ctx context.Context, staffID string) (*domain.User, error) {
	userProfile, err := d.query.GetUserProfileByStaffID(ctx, staffID)
	if err != nil {
		return nil, err
	}

	user := createMapUser(userProfile)
	return user, nil
}

// GetClientServiceRequestByID fetches a service request by id
func (d *MyCareHubDb) GetClientServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	serviceRequest, err := d.query.GetClientServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		return nil, err
	}

	metadata, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
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
		return nil, err
	}
	user := createMapUser(&staffProfile.UserProfile)

	facility, err := d.RetrieveFacility(ctx, &staffProfile.DefaultFacilityID, true)
	if err != nil {
		return nil, err
	}

	return &domain.StaffProfile{
		ID:                  staffProfile.ID,
		User:                user,
		UserID:              staffProfile.UserID,
		Active:              staffProfile.Active,
		StaffNumber:         staffProfile.StaffNumber,
		DefaultFacility:     facility,
		ProgramID:           staffProfile.ProgramID,
		IsOrganisationAdmin: staffProfile.IsOrganisationAdmin,
	}, nil
}

// GetAppointmentServiceRequests fetches all service requests of request type appointment given the last sync time
func (d *MyCareHubDb) GetAppointmentServiceRequests(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error) {
	mflIdentifier := &gorm.FacilityIdentifier{
		Type:  enums.FacilityIdentifierTypeMFLCode.String(),
		Value: mflCode,
	}

	facility, err := d.query.RetrieveFacilityByIdentifier(ctx, mflIdentifier, true)
	if err != nil {
		return nil, err
	}

	serviceRequests, err := d.query.GetAppointmentServiceRequests(ctx, lastSyncTime, *facility.FacilityID)
	if err != nil {
		return nil, err
	}

	appointmentServiceRequests := []domain.AppointmentServiceRequests{}
	for _, request := range serviceRequests {
		metaMap, err := utils.ConvertJSONStringToMap(request.Meta)
		if err != nil {
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
			return nil, err
		}

		var rescheduleTime time.Time
		valueRescheduleTime, exists := metaMap["rescheduleTime"]
		if !exists {
			continue
		}
		rescheduleTime, err = time.Parse(time.RFC3339, valueRescheduleTime.(string))
		if err != nil {
			return nil, err
		}

		suggestedDate, err := utils.ConvertTimeToScalarDate(rescheduleTime)
		if err != nil {
			return nil, err
		}

		var inProgressByName string
		if request.InProgressByID != nil {
			inProgressBy, err := d.GetUserProfileByStaffID(ctx, *request.InProgressByID)
			if err != nil {
				return nil, err
			}
			inProgressByName = inProgressBy.Name
		}

		var resolvedByName string
		if request.ResolvedByID != nil {
			resolvedBy, err := d.GetUserProfileByStaffID(ctx, *request.ResolvedByID)
			if err != nil {
				return nil, err
			}
			resolvedByName = resolvedBy.Name
		}

		clientProfile, err := d.query.GetClientProfileByClientID(ctx, request.ClientID)
		if err != nil {
			return nil, err
		}

		identifiers, err := d.GetClientIdentifiers(ctx, request.ClientID)
		if err != nil {
			continue
		}
		var identifierValue string
		for _, identifier := range identifiers {
			if identifier.Type == enums.UserIdentifierTypeCCC {
				identifierValue = identifier.Value
			}
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
			ClientContact: &clientProfile.User.Contacts.Value,
			CCCNumber:     identifierValue,
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
		return nil, err
	}

	for _, f := range results {
		facility, err := d.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, err
		}
		facilities = append(facilities, facility)
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
		return nil, err
	}

	date := appointment.Date
	appointmentDate, err := scalarutils.NewDate(date.Day(), int(date.Month()), date.Year())
	if err != nil {
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
		return nil, err
	}

	var serviceRequestList []*domain.ServiceRequest
	for _, r := range serviceRequests {
		meta, err := utils.ConvertJSONStringToMap(r.Meta)
		if err != nil {
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

// CheckAppointmentExistsByExternalID checks if an appointment with the external id exists
func (d *MyCareHubDb) CheckAppointmentExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	return d.query.CheckAppointmentExistsByExternalID(ctx, externalID)
}

// GetUserSurveyForms retrives all user survey forms
func (d *MyCareHubDb) GetUserSurveyForms(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
	var userSurveys []*domain.UserSurvey

	surveys, err := d.query.GetUserSurveyForms(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, s := range surveys {
		userSurveys = append(userSurveys, &domain.UserSurvey{
			ID:             s.ID,
			Active:         s.Active,
			Created:        s.CreatedAt,
			Link:           s.Link,
			Title:          s.Title,
			Description:    s.Description,
			HasSubmitted:   s.HasSubmitted,
			UserID:         s.UserID,
			Token:          s.Token,
			ProjectID:      s.ProjectID,
			FormID:         s.FormID,
			LinkID:         s.LinkID,
			ProgramID:      s.ProgramID,
			OrganisationID: s.OrganisationID,
		})
	}

	return userSurveys, nil
}

// GetSharedHealthDiaryEntries fetches the most recent shared health diary entry
func (d *MyCareHubDb) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	clientProfile, err := d.query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var healthDiaryEntries []*domain.ClientHealthDiaryEntry
	entries, err := d.query.GetSharedHealthDiaryEntries(ctx, clientID, facilityID)
	if err != nil {
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
			PhoneNumber:           clientProfile.User.Contacts.Value,
			ClientName:            clientProfile.User.Name,
			CaregiverID:           k.CaregiverID,
		})
	}

	return healthDiaryEntries, nil
}

// GetClientScreeningToolServiceRequestByToolType fetches a screening tool service request by tooltype, client ID and status
func (d *MyCareHubDb) GetClientScreeningToolServiceRequestByToolType(ctx context.Context, clientID, toolType, status string) (*domain.ServiceRequest, error) {
	serviceRequest, err := d.query.GetClientScreeningToolServiceRequestByToolType(ctx, clientID, toolType, status)
	if err != nil {
		return nil, err
	}
	meta, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
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
			return nil, err
		}
		facility, err := d.RetrieveFacility(ctx, &c.FacilityID, true)
		if err != nil {
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
			ClientCounselled:        c.ClientCounselled,
			OrganisationID:          c.OrganisationID,
			DefaultFacility:         facility,
		})
	}

	return clientList, nil
}

// SearchClientServiceRequests is used to query(search) for client service requests depending on the search parameter
func (d *MyCareHubDb) SearchClientServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	serviceRequests, err := d.query.SearchClientServiceRequests(ctx, searchParameter, requestType, facilityID)
	if err != nil {
		return nil, err
	}

	return d.ReturnClientsServiceRequests(ctx, serviceRequests)
}

// SearchStaffServiceRequests is used to query(search) for staff's service requests depending on the search parameter
func (d *MyCareHubDb) SearchStaffServiceRequests(ctx context.Context, searchParameter string, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	serviceRequests, err := d.query.SearchStaffServiceRequests(ctx, searchParameter, requestType, facilityID)
	if err != nil {
		return nil, err
	}

	return d.ReturnStaffServiceRequests(ctx, serviceRequests)
}

// GetScreeningToolByID fetches a screening tool by ID including the whole questions payload
func (d *MyCareHubDb) GetScreeningToolByID(ctx context.Context, toolID string) (*domain.ScreeningTool, error) {
	tool, err := d.query.GetScreeningToolByID(ctx, toolID)
	if err != nil {
		return nil, err
	}

	questionnaire, err := d.query.GetQuestionnaireByID(ctx, tool.QuestionnaireID)
	if err != nil {
		return nil, err
	}

	questionsPayload, err := d.query.GetQuestionsByQuestionnaireID(ctx, questionnaire.ID)
	if err != nil {
		return nil, err
	}

	questions := []domain.Question{}

	for _, q := range questionsPayload {
		choices := []domain.QuestionInputChoice{}
		choicesPayload, err := d.query.GetQuestionInputChoicesByQuestionID(ctx, q.ID)
		if err != nil {
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
func (d *MyCareHubDb) GetAvailableScreeningTools(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error) {
	clientTypes := []string{}
	for _, clientType := range screeningTool.ClientTypes {
		clientTypes = append(clientTypes, clientType.String())
	}

	genders := []string{}
	for _, gender := range screeningTool.Genders {
		genders = append(genders, gender.String())
	}

	screeningToolObj := gorm.ScreeningTool{
		ClientTypes: clientTypes,
		Genders:     genders,
		MinimumAge:  screeningTool.AgeRange.LowerBound,
		MaximumAge:  screeningTool.AgeRange.UpperBound,
		ProgramID:   screeningTool.ProgramID,
	}

	screeningTools, err := d.query.GetAvailableScreeningTools(ctx, clientID, screeningToolObj, screeningToolIDs)
	if err != nil {
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

// GetScreeningToolResponsesWithin24Hours gets the user screening response that are within 24 hours
func (d *MyCareHubDb) GetScreeningToolResponsesWithin24Hours(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
	screeningToolResponsesList, err := d.query.GetScreeningToolResponsesWithin24Hours(ctx, clientID, programID)
	if err != nil {
		return nil, err
	}
	var screeningToolResponses []*domain.QuestionnaireScreeningToolResponse

	for _, screeningToolResponse := range screeningToolResponsesList {
		screeningToolResponses = append(screeningToolResponses, &domain.QuestionnaireScreeningToolResponse{
			ID:              screeningToolResponse.ID,
			Active:          screeningToolResponse.Active,
			ScreeningToolID: screeningToolResponse.ScreeningToolID,
			FacilityID:      screeningToolResponse.FacilityID,
			ClientID:        screeningToolResponse.ClientID,
			DateOfResponse:  screeningToolResponse.CreatedAt,
			AggregateScore:  screeningToolResponse.AggregateScore,
			ProgramID:       screeningToolResponse.ProgramID,
			OrganisationID:  screeningToolResponse.OrganisationID,
		})
	}
	return screeningToolResponses, nil
}

// GetScreeningToolResponsesWithPendingServiceRequests gets the user screening response that have pending service requests
func (d *MyCareHubDb) GetScreeningToolResponsesWithPendingServiceRequests(ctx context.Context, clientID, programID string) ([]*domain.QuestionnaireScreeningToolResponse, error) {
	screeningToolResponsesList, err := d.query.GetScreeningToolResponsesWithPendingServiceRequests(ctx, clientID, programID)
	if err != nil {
		return nil, err
	}
	var screeningToolResponses []*domain.QuestionnaireScreeningToolResponse

	for _, screeningToolResponse := range screeningToolResponsesList {
		screeningToolResponses = append(screeningToolResponses, &domain.QuestionnaireScreeningToolResponse{
			ID:              screeningToolResponse.ID,
			Active:          screeningToolResponse.Active,
			ScreeningToolID: screeningToolResponse.ScreeningToolID,
			FacilityID:      screeningToolResponse.FacilityID,
			ClientID:        screeningToolResponse.ClientID,
			DateOfResponse:  screeningToolResponse.CreatedAt,
			AggregateScore:  screeningToolResponse.AggregateScore,
			ProgramID:       screeningToolResponse.ProgramID,
			OrganisationID:  screeningToolResponse.OrganisationID,
		})
	}
	return screeningToolResponses, nil
}

// GetFacilityRespondedScreeningTools fetches responded screening tools for a given facility
func (d *MyCareHubDb) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error) {
	screeningTools, pageInfo, err := d.query.GetFacilityRespondedScreeningTools(ctx, facilityID, programID, pagination)
	if err != nil {
		return nil, nil, err
	}

	var screeningToolList []*domain.ScreeningTool
	for _, s := range screeningTools {
		questionnaire, err := d.query.GetQuestionnaireByID(ctx, s.QuestionnaireID)
		if err != nil {
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
func (d *MyCareHubDb) GetScreeningToolRespondents(ctx context.Context, facilityID, programID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error) {

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	serviceRequests, pageInfo, err := d.query.GetScreeningToolServiceRequestOfRespondents(ctx, facilityID, programID, screeningToolID, searchTerm, page)
	if err != nil {
		return nil, nil, err
	}

	var respondents []*domain.ScreeningToolRespondent

	for _, s := range serviceRequests {
		meta, err := utils.ConvertJSONStringToMap(s.Meta)
		if err != nil {
			return nil, nil, err
		}
		responseID := meta["response_id"].(string)
		response, err := d.query.GetScreeningToolResponseByID(ctx, responseID)
		if err != nil {
			return nil, nil, err
		}
		client, err := d.query.GetClientProfileByClientID(ctx, s.ClientID)
		if err != nil {
			return nil, nil, err
		}
		respondent := &domain.ScreeningToolRespondent{
			ClientID:                s.ClientID,
			ScreeningToolResponseID: response.ID,
			ServiceRequestID:        *s.ID,
			ServiceRequest:          s.Request,
			Name:                    client.User.Name,
			PhoneNumber:             client.User.Contacts.Value,
		}

		respondents = append(respondents, respondent)
	}

	return respondents, pageInfo, nil
}

// GetScreeningToolResponseByID fetches a screening tool response by ID
func (d *MyCareHubDb) GetScreeningToolResponseByID(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	response, err := d.query.GetScreeningToolResponseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	screeningToolResponses, err := d.query.GetScreeningToolQuestionResponsesByResponseID(ctx, response.ID)
	if err != nil {
		return nil, err
	}
	screeningTool, err := d.GetScreeningToolByID(ctx, response.ScreeningToolID)
	if err != nil {
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
		CaregiverID:       response.CaregiverID,
	}, nil
}

// GetSurveysWithServiceRequests fetches all the surveys with a service request for a given facility
func (d *MyCareHubDb) GetSurveysWithServiceRequests(ctx context.Context, facilityID, programID string) ([]*dto.SurveysWithServiceRequest, error) {
	surveys, err := d.query.GetSurveysWithServiceRequests(ctx, facilityID, programID)
	if err != nil {
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
			PhoneNumber:      clientProfile.User.Contacts.Value,
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
		facility, err := d.RetrieveFacility(ctx, f.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}

		notification := &gorm.Notification{
			FacilityID: facility.ID,
			Flavour:    feedlib.FlavourPro,
			ProgramID:  input.ProgramID,
		}

		notificationCount, err := d.query.GetNotificationsCount(ctx, *notification)
		if err != nil {
			return nil, nil, err
		}

		staffPendingServiceRequest, err := d.query.GetClientsPendingServiceRequestsCount(ctx, *f.FacilityID, &input.ProgramID)
		if err != nil {
			return nil, nil, err
		}

		identifier, err := d.query.RetrieveFacilityIdentifierByFacilityID(ctx, facility.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed retrieve facility identifier: %w", err)
		}

		facilities = append(facilities, &domain.Facility{
			ID:                 facility.ID,
			Name:               facility.Name,
			Phone:              facility.Phone,
			Active:             facility.Active,
			Country:            facility.Country,
			Description:        facility.Description,
			FHIROrganisationID: facility.FHIROrganisationID,
			Identifier: domain.FacilityIdentifier{
				ID:     identifier.ID,
				Active: identifier.Active,
				Type:   enums.FacilityIdentifierType(identifier.Type),
				Value:  identifier.Value,
			},
			WorkStationDetails: domain.WorkStationDetails{
				Notifications:   notificationCount,
				ServiceRequests: staffPendingServiceRequest.Total,
			},
		})
	}

	return facilities, pageInfo, nil

}

// FindContacts retrieves all the contacts that match the given contact type and value.
// Contacts can be shared by users thus the same contact can have multiple records stored
func (d *MyCareHubDb) FindContacts(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
	records, err := d.query.FindContacts(ctx, contactType, contactValue)
	if err != nil {
		return nil, err
	}

	var contacts []*domain.Contact
	for _, record := range records {
		contact := domain.Contact{
			ID:             &record.ID,
			ContactType:    record.Type,
			ContactValue:   record.Value,
			Active:         record.Active,
			OptedIn:        record.OptedIn,
			UserID:         record.UserID,
			OrganisationID: record.OrganisationID,
		}

		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

// GetStaffUserPrograms retrieves all programs associated with a staff user
func (d *MyCareHubDb) GetStaffUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	records, err := d.query.GetStaffUserPrograms(ctx, userID)
	if err != nil {
		return nil, err
	}

	programs := []*domain.Program{}
	for _, record := range records {
		organisation, err := d.query.GetOrganisation(ctx, record.OrganisationID)
		if err != nil {
			return nil, err
		}
		program := domain.Program{
			ID:                 record.ID,
			Active:             record.Active,
			Name:               record.Name,
			Description:        record.Description,
			FHIROrganisationID: record.FHIROrganisationID,
			Organisation: domain.Organisation{
				ID:              *organisation.ID,
				Active:          organisation.Active,
				Code:            organisation.Code,
				Name:            organisation.Name,
				Description:     organisation.Description,
				EmailAddress:    organisation.EmailAddress,
				PhoneNumber:     organisation.PhoneNumber,
				PostalAddress:   organisation.PostalAddress,
				PhysicalAddress: organisation.PhysicalAddress,
				DefaultCountry:  organisation.DefaultCountry,
			},
		}

		programs = append(programs, &program)
	}

	return programs, nil
}

// GetClientUserPrograms retrieves all programs associated with a client user
func (d *MyCareHubDb) GetClientUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	records, err := d.query.GetClientUserPrograms(ctx, userID)
	if err != nil {
		return nil, err
	}

	programs := []*domain.Program{}
	for _, record := range records {
		organisation, err := d.query.GetOrganisation(ctx, record.OrganisationID)
		if err != nil {
			return nil, err
		}
		program := domain.Program{
			ID:                 record.ID,
			Active:             record.Active,
			Name:               record.Name,
			Description:        record.Description,
			FHIROrganisationID: record.FHIROrganisationID,
			Organisation: domain.Organisation{
				ID:              *organisation.ID,
				Active:          organisation.Active,
				Code:            organisation.Code,
				Name:            organisation.Name,
				Description:     organisation.Description,
				EmailAddress:    organisation.EmailAddress,
				PhoneNumber:     organisation.PhoneNumber,
				PostalAddress:   organisation.PostalAddress,
				PhysicalAddress: organisation.PhysicalAddress,
				DefaultCountry:  organisation.DefaultCountry,
			},
		}

		programs = append(programs, &program)
	}

	return programs, nil
}

// GetClientFacilities gets a list of client facilities
func (d *MyCareHubDb) GetClientFacilities(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
	clientProfile, err := d.query.GetClientProfileByClientID(ctx, *input.ClientID)
	if err != nil {
		return nil, nil, err
	}

	facilities := []*domain.Facility{}

	clientFacility := gorm.ClientFacilities{
		ClientID:   clientProfile.ID,
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
			ProgramID:  input.ProgramID,
		}

		notificationCount, err := d.query.GetNotificationsCount(ctx, *notification)
		if err != nil {
			return nil, nil, err
		}

		surveyCount, err := d.query.GetClientsSurveyCount(ctx, *clientProfile.UserID)
		if err != nil {
			return nil, nil, err
		}

		identifier, err := d.query.RetrieveFacilityIdentifierByFacilityID(ctx, facility.FacilityID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed retrieve facility identifier: %w", err)
		}

		facilities = append(facilities, &domain.Facility{
			ID:                 facility.FacilityID,
			Name:               facility.Name,
			Phone:              facility.Phone,
			Active:             facility.Active,
			Country:            facility.Country,
			Description:        facility.Description,
			FHIROrganisationID: facility.FHIROrganisationID,
			Identifier: domain.FacilityIdentifier{
				ID:     identifier.ID,
				Active: identifier.Active,
				Type:   enums.FacilityIdentifierType(identifier.Type),
				Value:  identifier.Value,
			},
			WorkStationDetails: domain.WorkStationDetails{
				Notifications: notificationCount,
				Surveys:       surveyCount,
			},
		})
	}

	return facilities, pageInfo, nil
}

// GetCaregiverManagedClients lists clients who are managed by the caregivers
// The clients should have given their consent to be managed by the caregivers
func (d *MyCareHubDb) GetCaregiverManagedClients(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error) {
	managedClients := []*domain.ManagedClient{}
	caregiverClients, pageInfo, err := d.query.GetCaregiverManagedClients(ctx, userID, pagination)
	if err != nil {
		return nil, nil, err
	}

	for _, caregiverClient := range caregiverClients {
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, caregiverClient.ClientID)
		if err != nil {
			return nil, nil, err
		}

		clientFacilityInput := dto.ClientFacilityInput{
			ClientID: clientProfile.ID,
		}

		clientFacilities, _, err := d.GetClientFacilities(ctx, clientFacilityInput, nil)
		if err != nil {
			return nil, nil, err
		}

		notification := &gorm.Notification{
			UserID:  clientProfile.UserID,
			Flavour: feedlib.FlavourConsumer,
		}

		notificationCount, err := d.query.GetNotificationsCount(ctx, *notification)
		if err != nil {
			return nil, nil, err
		}

		surveyCount, err := d.query.GetClientsSurveyCount(ctx, *clientProfile.UserID)
		if err != nil {
			return nil, nil, err
		}
		domainUser := createMapUser(&clientProfile.User)
		facility, err := d.RetrieveFacility(ctx, &clientProfile.FacilityID, true)
		if err != nil {
			return nil, nil, err
		}
		managedClient := &domain.ManagedClient{
			ClientProfile: &domain.ClientProfile{
				ID:              clientProfile.ID,
				User:            domainUser,
				DefaultFacility: facility,
				Facilities:      clientFacilities,
			},
			CaregiverConsent: caregiverClient.CaregiverConsent,
			ClientConsent:    caregiverClient.ClientConsent,
			WorkStationDetails: domain.WorkStationDetails{
				Notifications: notificationCount,
				Surveys:       surveyCount,
			},
		}
		managedClients = append(managedClients, managedClient)

	}

	return managedClients, pageInfo, nil

}

// ListClientsCaregivers retrieves a list of clients caregivers
func (d *MyCareHubDb) ListClientsCaregivers(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error) {
	caregivers := []*domain.CaregiverProfile{}

	clientCaregivers, pageInfo, err := d.query.ListClientsCaregivers(ctx, clientID, pagination)
	if err != nil {
		return nil, nil, err
	}

	caregiversClient := &domain.ClientCaregivers{}
	for _, clientCaregiver := range clientCaregivers {
		caregiver, err := d.query.GetCaregiverProfileByCaregiverID(ctx, clientCaregiver.CaregiverID)
		if err != nil {
			return nil, nil, err
		}

		user := createMapUser(&caregiver.UserProfile)

		caregivers = append(caregivers, &domain.CaregiverProfile{
			ID:              caregiver.ID,
			User:            *user,
			CaregiverNumber: caregiver.CaregiverNumber,
			Consent: domain.ConsentStatus{
				ConsentStatus: clientCaregiver.CaregiverConsent,
			},
			CurrentClient:   caregiver.CurrentClient,
			CurrentFacility: caregiver.CurrentFacility,
		})

		caregiversClient = &domain.ClientCaregivers{
			Caregivers: caregivers,
		}
	}

	return caregiversClient, pageInfo, nil
}

// CheckOrganisationExists check whether an organisation exists
func (d *MyCareHubDb) CheckOrganisationExists(ctx context.Context, organisationID string) (bool, error) {
	return d.query.CheckOrganisationExists(ctx, organisationID)
}

// CheckIfProgramNameExists checks if a program exists in the organization
// the program name should be unique for each program in a given organization
func (d *MyCareHubDb) CheckIfProgramNameExists(ctx context.Context, organisationID string, programName string) (bool, error) {
	return d.query.CheckIfProgramNameExists(ctx, organisationID, programName)
}

// ListOrganisations lists all organisations
func (d *MyCareHubDb) ListOrganisations(ctx context.Context, pagination *domain.Pagination) ([]*domain.Organisation, *domain.Pagination, error) {
	organisationObj, paginationInfo, err := d.query.ListOrganisations(ctx, pagination)
	if err != nil {
		return nil, nil, err
	}

	organisations := []*domain.Organisation{}
	for _, organisation := range organisationObj {
		programs, _, err := d.ListPrograms(ctx, organisation.ID, nil)
		if err != nil {
			return nil, nil, err
		}
		organisations = append(organisations, &domain.Organisation{
			ID:              *organisation.ID,
			Active:          organisation.Active,
			Code:            organisation.Code,
			Name:            organisation.Name,
			Description:     organisation.Description,
			EmailAddress:    organisation.EmailAddress,
			PhoneNumber:     organisation.PhoneNumber,
			PostalAddress:   organisation.PostalAddress,
			PhysicalAddress: organisation.PhysicalAddress,
			DefaultCountry:  organisation.DefaultCountry,
			Programs:        programs,
		})
	}

	return organisations, paginationInfo, nil
}

// GetProgramFacilities gets the facilities that belong the program
func (d *MyCareHubDb) GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error) {
	facilities := []*domain.Facility{}

	programFacilities, err := d.query.GetProgramFacilities(ctx, programID)
	if err != nil {
		return nil, err
	}

	for _, programFacility := range programFacilities {
		facility, err := d.RetrieveFacility(ctx, &programFacility.FacilityID, true)
		if err != nil {
			return nil, fmt.Errorf("failed retrieve facility by id: %w", err)
		}

		facilities = append(facilities, facility)
	}

	return facilities, nil

}

// GetProgramByID retrieves a program by its ID
func (d *MyCareHubDb) GetProgramByID(ctx context.Context, programID string) (*domain.Program, error) {
	program, err := d.query.GetProgramByID(ctx, programID)
	if err != nil {
		return nil, err
	}

	organisation, err := d.query.GetOrganisation(ctx, program.OrganisationID)
	if err != nil {
		return nil, err
	}

	programFacilities, err := d.query.GetProgramFacilities(ctx, programID)
	if err != nil {
		return nil, err
	}

	var facilities []*domain.Facility
	for _, programFacility := range programFacilities {
		facility, err := d.RetrieveFacility(ctx, &programFacility.FacilityID, true)
		if err != nil {
			return nil, fmt.Errorf("failed retrieve facility by id: %w", err)
		}

		facilities = append(facilities, facility)
	}

	return &domain.Program{
		ID:                 programID,
		Active:             program.Active,
		Name:               program.Name,
		Description:        program.Description,
		FHIROrganisationID: program.FHIROrganisationID,
		Organisation: domain.Organisation{
			ID:          program.OrganisationID,
			Name:        organisation.Name,
			Description: organisation.Description,
		},
		Facilities: facilities,
	}, nil
}

// ListPrograms gets a list of programs
func (d *MyCareHubDb) ListPrograms(ctx context.Context, organisationID *string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error) {
	programsObj, pageInfo, err := d.query.ListPrograms(ctx, organisationID, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get programs: %v", err)
	}

	programs := []*domain.Program{}
	for _, program := range programsObj {
		organisation, err := d.query.GetOrganisation(ctx, program.OrganisationID)
		if err != nil {
			return nil, nil, err
		}
		programs = append(programs, &domain.Program{
			ID:                 program.ID,
			Active:             program.Active,
			Name:               program.Name,
			Description:        program.Description,
			FHIROrganisationID: program.FHIROrganisationID,
			Organisation: domain.Organisation{
				ID:              *organisation.ID,
				Active:          organisation.Active,
				Code:            organisation.Code,
				Name:            organisation.Name,
				Description:     organisation.Description,
				EmailAddress:    organisation.EmailAddress,
				PhoneNumber:     organisation.PhoneNumber,
				PostalAddress:   organisation.PostalAddress,
				PhysicalAddress: organisation.PhysicalAddress,
				DefaultCountry:  organisation.DefaultCountry,
			},
		})
	}

	return programs, pageInfo, nil
}

// CheckIfSuperUserExists checks if there is a platform superuser
func (d *MyCareHubDb) CheckIfSuperUserExists(ctx context.Context) (bool, error) {
	return d.query.CheckIfSuperUserExists(ctx)
}

// GetCaregiverProfileByUserID gets the caregiver profile by user ID and organisation ID.
func (d *MyCareHubDb) GetCaregiverProfileByUserID(ctx context.Context, userID string, organisationID string) (*domain.CaregiverProfile, error) {
	caregiver, err := d.query.GetCaregiverProfileByUserID(ctx, userID, organisationID)
	if err != nil {
		return nil, err
	}

	user, err := d.query.GetUserProfileByUserID(ctx, &caregiver.UserID)
	if err != nil {
		return nil, err
	}
	userProfile := createMapUser(user)

	isClient, err := d.query.CheckClientExists(ctx, *userProfile.ID)
	if err != nil {
		return nil, err
	}

	return &domain.CaregiverProfile{
		ID:              caregiver.ID,
		UserID:          userID,
		User:            *userProfile,
		CaregiverNumber: caregiver.CaregiverNumber,
		IsClient:        isClient,
		CurrentClient:   caregiver.CurrentClient,
		CurrentFacility: caregiver.CurrentFacility,
	}, nil
}

// GetCaregiversClient gets the caregivers clients details
func (d *MyCareHubDb) GetCaregiversClient(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
	caregiversClientInput := gorm.CaregiverClient{
		CaregiverID: caregiverClient.CaregiverID,
		ClientID:    caregiverClient.ClientID,
	}

	caregiverClientProfile, err := d.query.GetCaregiversClient(ctx, caregiversClientInput)
	if err != nil {
		return nil, err
	}

	caregiverClients := []*domain.CaregiverClient{}

	for _, client := range caregiverClientProfile {
		caregiverClients = append(caregiverClients, &domain.CaregiverClient{
			CaregiverID:        client.CaregiverID,
			ClientID:           client.ClientID,
			Active:             client.Active,
			RelationshipType:   client.RelationshipType,
			CaregiverConsent:   client.ClientConsent,
			CaregiverConsentAt: client.CaregiverConsentAt,
			ClientConsent:      client.ClientConsent,
			ClientConsentAt:    client.ClientConsentAt,
			OrganisationID:     client.OrganisationID,
			AssignedBy:         client.AssignedBy,
			ProgramID:          client.ProgramID,
		})
	}

	return caregiverClients, nil
}

// GetCaregiverProfileByCaregiverID retrieves the caregivers profile based on the caregiver ID provided
func (d *MyCareHubDb) GetCaregiverProfileByCaregiverID(ctx context.Context, caregiverID string) (*domain.CaregiverProfile, error) {
	caregiver, err := d.query.GetCaregiverProfileByCaregiverID(ctx, caregiverID)
	if err != nil {
		return nil, err
	}

	user := createMapUser(&caregiver.UserProfile)

	isClient, err := d.query.CheckClientExists(ctx, caregiver.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.CaregiverProfile{
		ID:              caregiver.ID,
		UserID:          caregiver.UserID,
		User:            *user,
		CaregiverNumber: caregiver.CaregiverNumber,
		IsClient:        isClient,
		CurrentClient:   caregiver.CurrentClient,
		CurrentFacility: caregiver.CurrentFacility,
	}, nil
}

// SearchOrganisation searches for organisations based on the search parameter provided
func (d *MyCareHubDb) SearchOrganisation(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
	organisations, err := d.query.SearchOrganisation(ctx, searchParameter)
	if err != nil {
		return nil, err
	}

	orgs := []*domain.Organisation{}
	for _, org := range organisations {
		programs, _, err := d.ListPrograms(ctx, org.ID, nil)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, &domain.Organisation{
			ID:              *org.ID,
			Active:          org.Active,
			Code:            org.Code,
			Name:            org.Name,
			Description:     org.Description,
			EmailAddress:    org.EmailAddress,
			PhoneNumber:     org.PhoneNumber,
			PostalAddress:   org.PostalAddress,
			PhysicalAddress: org.PhysicalAddress,
			DefaultCountry:  org.DefaultCountry,
			Programs:        programs,
		})
	}

	return orgs, nil
}

// SearchPrograms searches for programs based on the search parameter provided and from the provided organisation
func (d *MyCareHubDb) SearchPrograms(ctx context.Context, searchParameter string, organisationID string, pagination *domain.Pagination) ([]*domain.Program, *domain.Pagination, error) {
	programs, pageInfo, err := d.query.SearchPrograms(ctx, searchParameter, organisationID, pagination)
	if err != nil {
		return nil, nil, err
	}

	programList := []*domain.Program{}

	for _, program := range programs {
		organisation, err := d.GetOrganisation(ctx, program.OrganisationID)
		if err != nil {
			return nil, nil, err
		}

		programList = append(programList, &domain.Program{
			ID:                 program.ID,
			Active:             program.Active,
			Name:               program.Name,
			Description:        program.Description,
			FHIROrganisationID: program.FHIROrganisationID,
			Organisation:       *organisation,
		})
	}

	return programList, pageInfo, nil
}

// ListCommunities  is used to list Matrix communities(rooms)
func (d *MyCareHubDb) ListCommunities(ctx context.Context, programID string, organisationID string) ([]*domain.Community, error) {
	records, err := d.query.ListCommunities(ctx, programID, organisationID)
	if err != nil {
		return nil, err
	}

	var communities []*domain.Community
	for _, record := range records {
		clientTypes := []enums.ClientType{}
		for _, k := range record.ClientTypes {
			clientTypes = append(clientTypes, enums.ClientType(k))
		}

		genders := []enumutils.Gender{}
		for _, k := range record.Gender {
			genders = append(genders, enumutils.Gender(k))
		}

		communities = append(communities, &domain.Community{
			ID:          record.ID,
			RoomID:      record.RoomID,
			Name:        record.Name,
			Description: record.Description,
			AgeRange: &domain.AgeRange{
				LowerBound: record.MinimumAge,
				UpperBound: record.MaximumAge,
			},
			Gender:         genders,
			ClientType:     clientTypes,
			OrganisationID: record.OrganisationID,
			ProgramID:      record.ProgramID,
		})
	}

	return communities, nil
}

// CheckPhoneExists is used to check if the phone number exists
func (d *MyCareHubDb) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	return d.query.CheckPhoneExists(ctx, phone)
}

// GetStaffServiceRequestByID gets the specified staff service request by ID
func (d *MyCareHubDb) GetStaffServiceRequestByID(ctx context.Context, serviceRequestID string) (*domain.ServiceRequest, error) {
	serviceRequest, err := d.query.GetStaffServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		return nil, err
	}
	metadata, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
		return nil, err
	}

	return &domain.ServiceRequest{
		ID:          *serviceRequest.ID,
		RequestType: serviceRequest.RequestType,
		Request:     serviceRequest.Request,
		Status:      serviceRequest.Status,
		Active:      serviceRequest.Active,
		StaffID:     serviceRequest.StaffID,
		CreatedAt:   serviceRequest.CreatedAt,
		Meta:        metadata,
	}, nil
}

// GetClientJWT retrieves a JWT by unique JTI
func (d *MyCareHubDb) GetClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
	result, err := d.query.GetClientJWT(ctx, jti)
	if err != nil {
		return nil, err
	}

	jwt := &domain.OauthClientJWT{
		ID:        result.ID,
		Active:    result.Active,
		JTI:       result.JTI,
		ExpiresAt: result.ExpiresAt,
	}

	return jwt, nil
}

// GetOauthClient retrieves a client by ID
func (d *MyCareHubDb) GetOauthClient(ctx context.Context, id string) (*domain.OauthClient, error) {
	result, err := d.query.GetOauthClient(ctx, id)
	if err != nil {
		return nil, err
	}

	client := &domain.OauthClient{
		ID:                      result.ID,
		Name:                    result.Name,
		Active:                  result.Active,
		Secret:                  result.Secret,
		RotatedSecrets:          result.RotatedSecrets,
		Public:                  result.Public,
		RedirectURIs:            result.RedirectURIs,
		Scopes:                  result.Scopes,
		Audience:                result.Audience,
		Grants:                  result.Grants,
		ResponseTypes:           result.ResponseTypes,
		TokenEndpointAuthMethod: result.TokenEndpointAuthMethod,
	}

	return client, nil
}

// GetValidClientJWT retrieves a JWT that is still valid i.e not expired
func (d *MyCareHubDb) GetValidClientJWT(ctx context.Context, jti string) (*domain.OauthClientJWT, error) {
	result, err := d.query.GetValidClientJWT(ctx, jti)
	if err != nil {
		return nil, err
	}

	jwt := &domain.OauthClientJWT{
		ID:        result.ID,
		Active:    result.Active,
		JTI:       result.JTI,
		ExpiresAt: result.ExpiresAt,
	}

	return jwt, nil
}

// GetAuthorizationCode retrieves an authorization code using the code
func (d *MyCareHubDb) GetAuthorizationCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	result, err := d.query.GetAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	}

	var form map[string][]string
	err = result.Form.AssignTo(&form)
	if err != nil {
		return nil, err
	}

	var sessionExtra map[string]interface{}
	err = result.Session.Extra.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	var sessionExpiresAt map[fosite.TokenType]time.Time
	err = result.Session.ExpiresAt.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	session := domain.Session{
		ID:        result.Session.ID,
		ClientID:  result.Session.ClientID,
		Username:  result.Session.Username,
		Subject:   result.Session.Subject,
		ExpiresAt: sessionExpiresAt,
		Extra:     sessionExtra,
		UserID:    result.Session.UserID,
	}

	client := result.Client

	authCode := &domain.AuthorizationCode{
		ID:                result.ID,
		Active:            result.Active,
		Code:              result.Code,
		RequestedAt:       result.RequestedAt,
		RequestedScopes:   result.RequestedScopes,
		GrantedScopes:     result.GrantedScopes,
		Form:              form,
		RequestedAudience: result.RequestedAudience,
		GrantedAudience:   result.GrantedAudience,
		SessionID:         result.SessionID,
		Session:           session,
		ClientID:          result.ClientID,
		Client: domain.OauthClient{
			ID:                      client.ID,
			Name:                    client.Name,
			Active:                  client.Active,
			Secret:                  client.Secret,
			RotatedSecrets:          client.RotatedSecrets,
			Public:                  client.Public,
			RedirectURIs:            client.RedirectURIs,
			Scopes:                  client.Scopes,
			Audience:                client.Audience,
			Grants:                  client.Grants,
			ResponseTypes:           client.ResponseTypes,
			TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		},
	}

	return authCode, nil
}

// GetAccessToken retrieves an access token using the signature
func (d *MyCareHubDb) GetAccessToken(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error) {
	params := gorm.AccessToken{
		ID:        token.ID,
		Signature: token.Signature,
	}

	result, err := d.query.GetAccessToken(ctx, params)
	if err != nil {
		return nil, err
	}

	var form map[string][]string
	err = result.Form.AssignTo(&form)
	if err != nil {
		return nil, err
	}

	var sessionExtra map[string]interface{}
	err = result.Session.Extra.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	var sessionExpiresAt map[fosite.TokenType]time.Time
	err = result.Session.ExpiresAt.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	session := domain.Session{
		ID:        result.Session.ID,
		ClientID:  result.Session.ClientID,
		Username:  result.Session.Username,
		Subject:   result.Session.Subject,
		ExpiresAt: sessionExpiresAt,
		Extra:     sessionExtra,
		UserID:    result.Session.UserID,
	}

	client := result.Client

	accessToken := &domain.AccessToken{
		ID:                result.ID,
		Active:            result.Active,
		Signature:         result.Signature,
		RequestedAt:       result.RequestedAt,
		RequestedScopes:   result.RequestedScopes,
		GrantedScopes:     result.GrantedScopes,
		Form:              form,
		RequestedAudience: result.RequestedAudience,
		GrantedAudience:   result.GrantedAudience,
		SessionID:         result.SessionID,
		Session:           session,
		ClientID:          result.ClientID,
		Client: domain.OauthClient{
			ID:                      client.ID,
			Name:                    client.Name,
			Active:                  client.Active,
			Secret:                  client.Secret,
			RotatedSecrets:          client.RotatedSecrets,
			Public:                  client.Public,
			RedirectURIs:            client.RedirectURIs,
			Scopes:                  client.Scopes,
			Audience:                client.Audience,
			Grants:                  client.Grants,
			ResponseTypes:           client.ResponseTypes,
			TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		},
	}

	return accessToken, nil
}

// GetRefreshToken retrieves a refresh token using the signature
func (d *MyCareHubDb) GetRefreshToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
	params := gorm.RefreshToken{
		ID:        token.ID,
		Signature: token.Signature,
	}

	result, err := d.query.GetRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}

	var form map[string][]string
	err = result.Form.AssignTo(&form)
	if err != nil {
		return nil, err
	}

	var sessionExtra map[string]interface{}
	err = result.Session.Extra.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	var sessionExpiresAt map[fosite.TokenType]time.Time
	err = result.Session.ExpiresAt.AssignTo(&sessionExtra)
	if err != nil {
		return nil, err
	}

	session := domain.Session{
		ID:        result.Session.ID,
		ClientID:  result.Session.ClientID,
		Username:  result.Session.Username,
		Subject:   result.Session.Subject,
		ExpiresAt: sessionExpiresAt,
		Extra:     sessionExtra,
		UserID:    result.Session.UserID,
	}

	client := result.Client

	refreshToken := &domain.RefreshToken{
		ID:                result.ID,
		Active:            result.Active,
		Signature:         result.Signature,
		RequestedAt:       result.RequestedAt,
		RequestedScopes:   result.RequestedScopes,
		GrantedScopes:     result.GrantedScopes,
		Form:              form,
		RequestedAudience: result.RequestedAudience,
		GrantedAudience:   result.GrantedAudience,
		SessionID:         result.SessionID,
		Session:           session,
		ClientID:          result.ClientID,
		Client: domain.OauthClient{
			ID:                      client.ID,
			Name:                    client.Name,
			Active:                  client.Active,
			Secret:                  client.Secret,
			RotatedSecrets:          client.RotatedSecrets,
			Public:                  client.Public,
			RedirectURIs:            client.RedirectURIs,
			Scopes:                  client.Scopes,
			Audience:                client.Audience,
			Grants:                  client.Grants,
			ResponseTypes:           client.ResponseTypes,
			TokenEndpointAuthMethod: client.TokenEndpointAuthMethod,
		},
	}

	return refreshToken, nil
}

// CheckIfClientHasPendingSurveyServiceRequest returns true if client has a pending survey service request
func (d *MyCareHubDb) CheckIfClientHasPendingSurveyServiceRequest(ctx context.Context, clientID string, projectID int, formID string) (bool, error) {
	return d.query.CheckIfClientHasPendingSurveyServiceRequest(ctx, clientID, projectID, formID)
}

// GetUserProfileByPushToken is used to fetch user's profile using their device token. Device token is unique for every user
func (d *MyCareHubDb) GetUserProfileByPushToken(ctx context.Context, pushToken string) (*domain.User, error) {
	user, err := d.query.GetUserProfileByPushToken(ctx, pushToken)
	if err != nil {
		return nil, err
	}

	return d.mapProfileObjectToDomain(user), nil
}

// CheckStaffExistsInProgram checks if a staff user is registered in a program
func (d *MyCareHubDb) CheckStaffExistsInProgram(ctx context.Context, userID, programID string) (bool, error) {
	return d.query.CheckStaffExistsInProgram(ctx, userID, programID)
}

// CheckIfFacilityExistsInProgram checks if a facility is associated with a program
func (d *MyCareHubDb) CheckIfFacilityExistsInProgram(ctx context.Context, programID, facilityID string) (bool, error) {
	return d.query.CheckIfFacilityExistsInProgram(ctx, programID, facilityID)
}

// CheckIfClientExistsInProgram checks if a client exists in a program
func (d *MyCareHubDb) CheckIfClientExistsInProgram(ctx context.Context, userID, programID string) (bool, error) {
	return d.query.CheckIfClientExistsInProgram(ctx, userID, programID)
}

// GetUserClientProfiles gets all client profiles for a given user
func (d *MyCareHubDb) GetUserClientProfiles(ctx context.Context, userID string) ([]*domain.ClientProfile, error) {
	clientProfilseObj, err := d.query.GetUserClientProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var clientProfiles []*domain.ClientProfile

	for _, clientProfile := range clientProfilseObj {
		var clientList []enums.ClientType
		for _, k := range clientProfile.ClientTypes {
			clientList = append(clientList, enums.ClientType(k))
		}

		facility, err := d.RetrieveFacility(ctx, &clientProfile.FacilityID, true)
		if err != nil {
			return nil, err
		}

		identifiers, err := d.GetClientIdentifiers(ctx, *clientProfile.ID)
		if err != nil {
			return nil, err
		}

		facilities, _, err := d.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: clientProfile.ID}, nil)
		if err != nil {
			log.Printf("failed to get client facilities: %v", err)
		}

		user, err := d.GetUserProfileByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		clientProfiles = append(clientProfiles, &domain.ClientProfile{
			ID:                      clientProfile.ID,
			User:                    user,
			Active:                  clientProfile.Active,
			ClientTypes:             clientList,
			TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
			FHIRPatientID:           clientProfile.FHIRPatientID,
			HealthRecordID:          clientProfile.HealthRecordID,
			ClientCounselled:        clientProfile.ClientCounselled,
			OrganisationID:          clientProfile.OrganisationID,
			ProgramID:               clientProfile.ProgramID,
			DefaultFacility:         facility,
			Facilities:              facilities,
			Identifiers:             identifiers,
		})
	}

	return clientProfiles, nil
}

// GetUserStaffProfiles gets all staff profiles for a given user
func (d *MyCareHubDb) GetUserStaffProfiles(ctx context.Context, userID string) ([]*domain.StaffProfile, error) {
	staffProfilesObj, err := d.query.GetUserStaffProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var staffProfiles []*domain.StaffProfile

	for _, staffProfile := range staffProfilesObj {
		facilities, _, err := d.GetStaffFacilities(ctx, dto.StaffFacilityInput{StaffID: staffProfile.ID}, nil)
		if err != nil {
			log.Printf("unable to get staff facilities: %v", err)
		}

		facility, err := d.RetrieveFacility(ctx, &staffProfile.DefaultFacilityID, true)
		if err != nil {
			return nil, err
		}

		user, err := d.GetUserProfileByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		nationalIDIdentifier := enums.UserIdentifierTypeNationalID.String()

		identifiersObj, err := d.query.GetStaffIdentifiers(ctx, *staffProfile.ID, &nationalIDIdentifier)
		if err != nil {
			return nil, err
		}

		var identifiers []*domain.Identifier

		for _, identifier := range identifiersObj {
			identifiers = append(identifiers, &domain.Identifier{
				ID:                  identifier.ID,
				Type:                enums.UserIdentifierType(identifier.Type),
				Value:               identifier.Value,
				Use:                 identifier.Use,
				Description:         identifier.Description,
				ValidFrom:           identifier.ValidFrom,
				ValidTo:             identifier.ValidTo,
				IsPrimaryIdentifier: identifier.IsPrimaryIdentifier,
				Active:              identifier.Active,
				ProgramID:           identifier.ProgramID,
				OrganisationID:      identifier.OrganisationID,
			})
		}

		staffProfiles = append(staffProfiles, &domain.StaffProfile{
			ID:                  staffProfile.ID,
			User:                user,
			UserID:              staffProfile.UserID,
			Active:              staffProfile.Active,
			StaffNumber:         staffProfile.StaffNumber,
			Facilities:          facilities,
			ProgramID:           staffProfile.ProgramID,
			DefaultFacility:     facility,
			IsOrganisationAdmin: staffProfile.IsOrganisationAdmin,
			Identifiers:         identifiers,
		})
	}
	return staffProfiles, nil
}
