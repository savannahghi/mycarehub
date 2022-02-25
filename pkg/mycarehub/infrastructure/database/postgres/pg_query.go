package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/serverutils"
)

// CheckWhetherUserHasLikedContent performs a operation to check whether user has liked the content
func (d *MyCareHubDb) CheckWhetherUserHasLikedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID < 1 {
		return false, fmt.Errorf("invalid userID or contentID")
	}
	return d.query.CheckWhetherUserHasLikedContent(ctx, userID, contentID)
}

//GetFacilities returns a slice of healthcare facilities in the platform.
func (d *MyCareHubDb) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.GetFacilities(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}
	for _, m := range facilities {
		singleFacility := domain.Facility{
			ID:          m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Phone:       m.Phone,
			Active:      m.Active,
			County:      m.County,
			Description: m.Description,
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
	if MFLCode == 0 {
		return nil, fmt.Errorf("facility ID should be defined")
	}
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

// GetSecurityQuestionResponseByID returns the security question response from the database
func (d *MyCareHubDb) GetSecurityQuestionResponseByID(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error) {
	if questionID == "" {
		return nil, fmt.Errorf("security question ID must be defined")
	}

	response, err := d.query.GetSecurityQuestionResponseByID(ctx, questionID)
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

//VerifyOTP performs the checking of OTP's existence for the specified user.
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

	response, err := d.query.GetClientProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	user := createMapUser(&response.UserProfile)
	return &domain.ClientProfile{
		ID:                      response.ID,
		User:                    user,
		Active:                  response.Active,
		ClientType:              response.ClientType,
		TreatmentEnrollmentDate: response.TreatmentEnrollmentDate,
		FHIRPatientID:           response.FHIRPatientID,
		HealthRecordID:          response.HealthRecordID,
		TreatmentBuddy:          response.TreatmentBuddy,
		ClientCounselled:        response.ClientCounselled,
		OrganisationID:          response.OrganisationID,
		FacilityID:              response.FacilityID,
		CHVUserID:               response.CHVUserID,
	}, nil
}

//GetStaffProfileByUserID fetches the staff's profile using the user's ID and returns the staff's profile in the login response.
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

	var facility []domain.Facility
	domainFacility := &domain.Facility{
		ID:          staffDefaultFacility.FacilityID,
		Name:        staffDefaultFacility.Name,
		Code:        staffDefaultFacility.Code,
		Phone:       staffDefaultFacility.Phone,
		Active:      staffDefaultFacility.Active,
		County:      staffDefaultFacility.County,
		Description: staffDefaultFacility.Description,
	}
	facility = append(facility, *domainFacility)

	user := createMapUser(&staff.UserProfile)
	return &domain.StaffProfile{
		ID:                staff.ID,
		User:              user,
		UserID:            staff.UserID,
		Active:            staff.Active,
		StaffNumber:       staff.StaffNumber,
		Facilities:        facility,
		DefaultFacilityID: staff.DefaultFacilityID,
	}, nil
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

//ListContentCategories retrieves the list of all content categories
func (d *MyCareHubDb) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	var contentItemCategory []*domain.ContentItemCategory

	allContentCategories, err := d.query.ListContentCategories(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	if len(allContentCategories) == 0 {
		return contentItemCategory, nil
	}

	for _, contentCategories := range allContentCategories {
		iconURL := fmt.Sprintf(serverutils.MustGetEnvVar(helpers.GoogleCloudStorageURL) + contentCategories.IconURL)

		contentCategoryItem := &domain.ContentItemCategory{
			ID:      contentCategories.ID,
			Name:    contentCategories.Name,
			IconURL: iconURL,
		}

		contentItemCategory = append(contentItemCategory, contentCategoryItem)
	}

	return contentItemCategory, nil
}

// GetUserBookmarkedContent is used to retrieve a user's bookmarked/pinned content
func (d *MyCareHubDb) GetUserBookmarkedContent(ctx context.Context, userID string) ([]*domain.ContentItem, error) {
	var domainContent []*domain.ContentItem
	bookmarkedContent, err := d.query.GetUserBookmarkedContent(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch user's bookmarked content: %v", err)
	}

	if len(bookmarkedContent) == 0 {
		return []*domain.ContentItem{}, nil
	}

	for _, content := range bookmarkedContent {
		contentItem := &domain.ContentItem{
			ID:                  content.PagePtrID,
			LikeCount:           content.LikeCount,
			BookmarkCount:       content.BookmarkCount,
			Body:                content.Body,
			ShareCount:          content.ShareCount,
			ViewCount:           content.ViewCount,
			Intro:               content.Intro,
			ItemType:            content.ItemType,
			TimeEstimateSeconds: content.TimeEstimateSeconds,
			Date:                content.Date.Format("2006-01-02"),
		}
		domainContent = append(domainContent, contentItem)
	}

	return domainContent, nil
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
func (d *MyCareHubDb) GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
	clientHealthDiaryQuote, err := d.query.GetClientHealthDiaryQuote(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to fetch client health diary quote: %v", err)
	}
	return &domain.ClientHealthDiaryQuote{
		Author: clientHealthDiaryQuote.Author,
		Quote:  clientHealthDiaryQuote.Quote,
	}, nil
}

// CheckIfUserBookmarkedContent is used to check if the user has bookmarked the content
func (d *MyCareHubDb) CheckIfUserBookmarkedContent(ctx context.Context, userID string, contentID int) (bool, error) {
	bookmarked, err := d.query.CheckIfUserBookmarkedContent(ctx, userID, contentID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to check if user bookmarked content: %v", err)
	}
	return bookmarked, nil
}

// GetPendingServiceRequestsCount gets the total number of service requests
func (d *MyCareHubDb) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
	return d.query.GetPendingServiceRequestsCount(ctx, facilityID)
}

// GetClientHealthDiaryEntries queries the database to return a clients all health diary records
func (d *MyCareHubDb) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	var healthDiaryEntries []*domain.ClientHealthDiaryEntry
	clientHealthDiaryEntry, err := d.query.GetClientHealthDiaryEntries(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, healthdiary := range clientHealthDiaryEntry {
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			Active:                healthdiary.Active,
			Mood:                  healthdiary.Mood,
			Note:                  healthdiary.Note,
			EntryType:             healthdiary.EntryType,
			ShareWithHealthWorker: healthdiary.ShareWithHealthWorker,
			SharedAt:              healthdiary.SharedAt,
			ClientID:              healthdiary.ClientID,
			CreatedAt:             healthdiary.CreatedAt,
		}
		healthDiaryEntries = append(healthDiaryEntries, healthDiaryEntry)
	}

	return healthDiaryEntries, nil
}

// GetFAQContent retrieves the FAQ content for the specified flavour
// an optional limit can be passed to the function to limit the number of FAQs returned
func (d *MyCareHubDb) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
	var faq []*domain.FAQ
	faqs, err := d.query.GetFAQContent(ctx, flavour, limit)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, faqItem := range faqs {
		faqItem := &domain.FAQ{
			ID:          faqItem.FAQID,
			Active:      faqItem.Active,
			Title:       faqItem.Title,
			Description: faqItem.Description,
			Body:        faqItem.Body,
			Flavour:     faqItem.Flavour,
		}
		faq = append(faq, faqItem)
	}

	return faq, nil
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
	user := createMapUser(&response.UserProfile)
	return &domain.ClientProfile{
		ID:                      response.ID,
		User:                    user,
		Active:                  response.Active,
		ClientType:              response.ClientType,
		TreatmentEnrollmentDate: response.TreatmentEnrollmentDate,
		FHIRPatientID:           response.FHIRPatientID,
		HealthRecordID:          response.HealthRecordID,
		TreatmentBuddy:          response.TreatmentBuddy,
		ClientCounselled:        response.ClientCounselled,
		OrganisationID:          response.OrganisationID,
		FacilityID:              response.FacilityID,
		CHVUserID:               response.CHVUserID,
		CaregiverID:             response.CaregiverID,
	}, nil
}

// GetServiceRequests retrieves the service requests by the type passed in the parameters
func (d *MyCareHubDb) GetServiceRequests(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
	var serviceRequests []*domain.ServiceRequest
	clientServiceRequests, err := d.query.GetServiceRequests(ctx, requestType, requestStatus, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, serviceRequest := range clientServiceRequests {
		clientProfile, err := d.query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		userProfile, err := d.query.GetUserProfileByUserID(ctx, clientProfile.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return nil, err
		}

		serviceRequest := &domain.ServiceRequest{
			ID:            *serviceRequest.ID,
			RequestType:   serviceRequest.RequestType,
			Request:       serviceRequest.Request,
			Status:        serviceRequest.Status,
			ClientID:      serviceRequest.ClientID,
			InProgressAt:  serviceRequest.InProgressAt,
			InProgressBy:  serviceRequest.InProgressByID,
			ResolvedAt:    serviceRequest.ResolvedAt,
			ResolvedBy:    serviceRequest.ResolvedByID,
			FacilityID:    facilityID,
			ClientName:    &userProfile.Name,
			ClientContact: &userProfile.Contacts.ContactValue,
		}
		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return serviceRequests, err
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
			RoleID: *role.AuthorityRoleID,
			Name:   enums.UserRoleType(role.Name),
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
