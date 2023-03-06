package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// SaveTemporaryUserPin does the actual saving of the users PIN in the database
func (d *MyCareHubDb) SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	pinObj := &gorm.PINData{
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: pinData.ValidFrom,
		ValidTo:   pinData.ValidTo,
		IsValid:   pinData.IsValid,
		Salt:      pinData.Salt,
	}

	_, err := d.create.SaveTemporaryUserPin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}

// SavePin gets the pin details from the user and saves it in the database
func (d *MyCareHubDb) SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error) {

	pinObj := &gorm.PINData{
		UserID:    pinInput.UserID,
		HashedPIN: pinInput.HashedPIN,
		ValidFrom: pinInput.ValidFrom,
		ValidTo:   pinInput.ValidTo,
		IsValid:   pinInput.IsValid,
		Salt:      pinInput.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}

// SaveOTP saves the otp to the database
func (d *MyCareHubDb) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	//Invalidate other OTPs before saving the new OTP by setting valid to false
	if otpInput.PhoneNumber == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	if !otpInput.Flavour.IsValid() {
		return fmt.Errorf("flavour %v is invalid", otpInput.Flavour)
	}

	otpObject := &gorm.UserOTP{
		UserID:      otpInput.UserID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         otpInput.OTP,
		Flavour:     otpInput.Flavour,
	}

	err := d.create.SaveOTP(ctx, otpObject)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	return nil
}

// SaveSecurityQuestionResponse saves the security question response to the database
func (d *MyCareHubDb) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error {
	var securityQuestionResponseObj []*gorm.SecurityQuestionResponse
	for _, sqr := range securityQuestionResponse {
		response := &gorm.SecurityQuestionResponse{
			UserID:     sqr.UserID,
			QuestionID: sqr.SecurityQuestionID,
			Active:     true,
			Response:   sqr.Response,
		}
		securityQuestionResponseObj = append(securityQuestionResponseObj, response)
	}

	err := d.create.SaveSecurityQuestionResponse(ctx, securityQuestionResponseObj)
	if err != nil {
		return fmt.Errorf("failed to save security question response data")
	}

	return nil
}

// CreateHealthDiaryEntry is used to add a health diary record to the database.
func (d *MyCareHubDb) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) (*domain.ClientHealthDiaryEntry, error) {
	healthDiaryResponse := &gorm.ClientHealthDiaryEntry{
		Active:                healthDiaryInput.Active,
		Mood:                  healthDiaryInput.Mood,
		Note:                  healthDiaryInput.Note,
		EntryType:             healthDiaryInput.EntryType,
		ShareWithHealthWorker: healthDiaryInput.ShareWithHealthWorker,
		SharedAt:              healthDiaryInput.SharedAt,
		ProgramID:             healthDiaryInput.ProgramID,
		ClientID:              healthDiaryInput.ClientID,
		OrganisationID:        healthDiaryInput.OrganisationID,
	}

	healthDiaryEntry, err := d.create.CreateHealthDiaryEntry(ctx, healthDiaryResponse)
	if err != nil {
		return nil, err
	}

	return &domain.ClientHealthDiaryEntry{
		ID:                    healthDiaryEntry.ClientHealthDiaryEntryID,
		Active:                healthDiaryEntry.Active,
		Mood:                  healthDiaryEntry.Mood,
		Note:                  healthDiaryEntry.Mood,
		EntryType:             healthDiaryEntry.EntryType,
		ShareWithHealthWorker: healthDiaryEntry.ShareWithHealthWorker,
		SharedAt:              healthDiaryEntry.SharedAt,
		ClientID:              healthDiaryEntry.ClientID,
		CreatedAt:             healthDiaryEntry.CreatedAt,
		ProgramID:             healthDiaryEntry.ProgramID,
		OrganisationID:        healthDiaryEntry.OrganisationID,
	}, nil
}

// CreateServiceRequest creates  a service request which will be handled by a staff user.
// This happens in a transaction because we do not want to
// create a health diary entry without a subsequent service request when the client's mood is "VERY_BAD"
func (d *MyCareHubDb) CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	meta, err := json.Marshal(serviceRequestInput.Meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta data: %v", err)
	}
	serviceRequest := &gorm.ClientServiceRequest{
		Active:         serviceRequestInput.Active,
		RequestType:    serviceRequestInput.RequestType,
		Request:        serviceRequestInput.Request,
		Status:         serviceRequestInput.Status,
		ClientID:       serviceRequestInput.ClientID,
		FacilityID:     serviceRequestInput.FacilityID,
		ProgramID:      serviceRequestInput.ProgramID,
		Meta:           string(meta),
		OrganisationID: serviceRequestInput.OrganisationID,
	}

	err = d.create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		return err
	}

	return nil
}

// CreateStaffServiceRequest creates a new service request for the specified staff
func (d *MyCareHubDb) CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	meta, err := json.Marshal(serviceRequestInput.Meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta data: %v", err)
	}
	serviceRequest := &gorm.StaffServiceRequest{
		Active:            serviceRequestInput.Active,
		RequestType:       serviceRequestInput.RequestType,
		Request:           serviceRequestInput.Request,
		Status:            serviceRequestInput.Status,
		StaffID:           serviceRequestInput.StaffID,
		DefaultFacilityID: &serviceRequestInput.FacilityID,
		Meta:              string(meta),
		ProgramID:         serviceRequestInput.ProgramID,
	}

	err = d.create.CreateStaffServiceRequest(ctx, serviceRequest)
	if err != nil {
		return err
	}

	return nil
}

// CreateCommunity creates a channel in the database
func (d *MyCareHubDb) CreateCommunity(ctx context.Context, community *domain.Community) (*domain.Community, error) {

	var genderList pq.StringArray
	for _, g := range community.Gender {
		genderList = append(genderList, string(g))
	}

	var clientTypeList pq.StringArray
	for _, c := range community.ClientType {
		clientTypeList = append(clientTypeList, string(c))
	}

	input := &gorm.Community{
		RoomID:         community.RoomID,
		Name:           community.Name,
		Description:    community.Description,
		Active:         true,
		MinimumAge:     community.AgeRange.LowerBound,
		MaximumAge:     community.AgeRange.UpperBound,
		Gender:         genderList,
		ClientTypes:    clientTypeList,
		ProgramID:      community.ProgramID,
		OrganisationID: community.OrganisationID,
	}

	record, err := d.create.CreateCommunity(ctx, input)
	if err != nil {
		return nil, err
	}

	var genders []enumutils.Gender
	for _, k := range record.Gender {
		genders = append(genders, enumutils.Gender(k))
	}

	var clientTypes []enums.ClientType
	for _, k := range record.ClientTypes {
		clientTypes = append(clientTypes, enums.ClientType(k))
	}

	return &domain.Community{
		ID:          record.ID,
		Name:        record.Name,
		RoomID:      record.RoomID,
		Description: record.Description,
		AgeRange: &domain.AgeRange{
			LowerBound: record.MinimumAge,
			UpperBound: record.MaximumAge,
		},
		Gender:         genders,
		ClientType:     clientTypes,
		OrganisationID: record.OrganisationID,
		ProgramID:      record.ProgramID,
	}, nil
}

// GetOrCreateNextOfKin creates a related person who is a next of kin
func (d *MyCareHubDb) GetOrCreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {

	pn := &gorm.RelatedPerson{
		FirstName:        person.Name,
		RelationshipType: "NEXT_OF_KIN",
		ProgramID:        person.ProgramID,
	}

	return d.create.GetOrCreateNextOfKin(ctx, pn, clientID, contactID)
}

// GetOrCreateContact creates a contact
func (d *MyCareHubDb) GetOrCreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {

	ct := &gorm.Contact{
		Active:         true,
		Type:           contact.ContactType,
		Value:          contact.ContactValue,
		UserID:         contact.UserID,
		OptedIn:        contact.OptedIn,
		OrganisationID: contact.OrganisationID,
	}

	c, err := d.create.GetOrCreateContact(ctx, ct)
	if err != nil {
		return nil, err
	}

	return &domain.Contact{
		ID:           &c.ID,
		ContactType:  c.Type,
		ContactValue: c.Value,
		Active:       c.Active,
		OptedIn:      c.OptedIn,
	}, nil
}

// CreateAppointment creates a new appointment
func (d *MyCareHubDb) CreateAppointment(ctx context.Context, appointment domain.Appointment) error {

	date := appointment.Date.AsTime()
	ap := &gorm.Appointment{
		Active:     true,
		ExternalID: appointment.ExternalID,
		ClientID:   appointment.ClientID,
		FacilityID: appointment.FacilityID,
		Reason:     appointment.Reason,
		Provider:   appointment.Provider,
		Date:       date,
		ProgramID:  appointment.ProgramID,
	}

	return d.create.CreateAppointment(ctx, ap)
}

// AnswerScreeningToolQuestions creates a screening tool answers
func (d *MyCareHubDb) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error {

	var screeningToolResponsesObj []*gorm.ScreeningToolsResponse
	for _, st := range screeningToolResponses {
		stq := &gorm.ScreeningToolsResponse{
			ClientID:   st.ClientID,
			QuestionID: st.QuestionID,
			Response:   st.Response,
			Active:     true,
			ProgramID:  st.ProgramID,
		}
		screeningToolResponsesObj = append(screeningToolResponsesObj, stq)
	}
	err := d.create.AnswerScreeningToolQuestions(ctx, screeningToolResponsesObj)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new user
func (d *MyCareHubDb) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {

	u := &gorm.User{
		Active:           true,
		Username:         user.Username,
		Name:             user.Name,
		Gender:           user.Gender,
		DateOfBirth:      user.DateOfBirth,
		CurrentProgramID: user.CurrentProgramID,
	}

	err := d.create.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return createMapUser(u), nil
}

// CreateClient creates a new client
func (d *MyCareHubDb) CreateClient(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
	var clientTypes pq.StringArray
	for _, c := range client.ClientTypes {
		clientTypes = append(clientTypes, c.String())
	}
	c := &gorm.Client{
		Active:                  true,
		UserID:                  &client.UserID,
		FacilityID:              *client.DefaultFacility.ID,
		ClientCounselled:        client.ClientCounselled,
		ClientTypes:             clientTypes,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
		ProgramID:               client.ProgramID,
	}

	err := d.create.CreateClient(ctx, c, contactID, identifierID)
	if err != nil {
		return nil, err
	}

	user := createMapUser(&c.User)

	var clientList []enums.ClientType
	for _, k := range c.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	return &domain.ClientProfile{
		ID:                      c.ID,
		User:                    user,
		Active:                  c.Active,
		ClientTypes:             clientList,
		TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
		FHIRPatientID:           c.FHIRPatientID,
		HealthRecordID:          c.HealthRecordID,
		ClientCounselled:        c.ClientCounselled,
		OrganisationID:          c.OrganisationID,
		DefaultFacility: &domain.Facility{
			ID: &c.FacilityID,
		},
		ProgramID: c.ProgramID,
	}, nil
}

// RegisterClient registers a client in the database
func (d *MyCareHubDb) RegisterClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
	usr := &gorm.User{
		Username:              payload.UserProfile.Username,
		Name:                  payload.UserProfile.Name,
		Gender:                payload.UserProfile.Gender,
		DateOfBirth:           payload.UserProfile.DateOfBirth,
		Active:                payload.UserProfile.Active,
		CurrentProgramID:      payload.UserProfile.CurrentProgramID,
		CurrentOrganisationID: payload.UserProfile.CurrentOrganizationID,
	}

	contact := &gorm.Contact{
		Type:           payload.Phone.ContactType,
		Value:          payload.Phone.ContactValue,
		Active:         payload.Phone.Active,
		OptedIn:        payload.Phone.Active,
		OrganisationID: payload.Phone.OrganisationID,
	}

	identifier := &gorm.Identifier{
		IdentifierType:      payload.ClientIdentifier.IdentifierType,
		IdentifierValue:     payload.ClientIdentifier.IdentifierValue,
		IdentifierUse:       payload.ClientIdentifier.IdentifierUse,
		Description:         payload.ClientIdentifier.Description,
		IsPrimaryIdentifier: payload.ClientIdentifier.IsPrimaryIdentifier,
		Active:              payload.ClientIdentifier.Active,
		ProgramID:           payload.ClientIdentifier.ProgramID,
		OrganisationID:      payload.ClientIdentifier.OrganisationID,
	}

	var pgClientTypes pq.StringArray
	for _, clientType := range payload.Client.ClientTypes {
		pgClientTypes = append(pgClientTypes, clientType.String())
	}
	clientProfile := &gorm.Client{
		ClientTypes:             pgClientTypes,
		TreatmentEnrollmentDate: payload.Client.TreatmentEnrollmentDate,
		FacilityID:              *payload.Client.DefaultFacility.ID,
		ClientCounselled:        payload.Client.ClientCounselled,
		Active:                  payload.Client.Active,
		ProgramID:               payload.Client.ProgramID,
		OrganisationID:          payload.Client.OrganisationID,
	}

	client, err := d.create.RegisterClient(ctx, usr, contact, identifier, clientProfile)
	if err != nil {
		return nil, err
	}

	var clientTypes []enums.ClientType
	for _, k := range clientProfile.ClientTypes {
		clientTypes = append(clientTypes, enums.ClientType(k))
	}

	return &domain.ClientProfile{
		ID:                      clientProfile.ID,
		Active:                  clientProfile.Active,
		ClientTypes:             clientTypes,
		TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
		UserID:                  *client.UserID,
		ClientCounselled:        clientProfile.ClientCounselled,
		DefaultFacility: &domain.Facility{
			ID: &clientProfile.FacilityID,
		},
		User:           createMapUser(usr),
		OrganisationID: clientProfile.OrganisationID,
		ProgramID:      clientProfile.ProgramID,
	}, nil
}

// RegisterExistingUserAsClient registers an existing user as a client and returns the client profile
func (d *MyCareHubDb) RegisterExistingUserAsClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
	identifier := &gorm.Identifier{
		Active:              payload.ClientIdentifier.Active,
		IdentifierType:      payload.ClientIdentifier.IdentifierType,
		IdentifierValue:     payload.ClientIdentifier.IdentifierValue,
		IdentifierUse:       payload.ClientIdentifier.IdentifierUse,
		Description:         payload.ClientIdentifier.Description,
		IsPrimaryIdentifier: payload.ClientIdentifier.IsPrimaryIdentifier,
		OrganisationID:      payload.ClientIdentifier.OrganisationID,
		ProgramID:           payload.ClientIdentifier.ProgramID,
	}

	var pgClientTypes pq.StringArray
	for _, clientType := range payload.Client.ClientTypes {
		pgClientTypes = append(pgClientTypes, clientType.String())
	}

	clientProfile := &gorm.Client{
		UserID:                  &payload.Client.UserID,
		ClientTypes:             pgClientTypes,
		TreatmentEnrollmentDate: payload.Client.TreatmentEnrollmentDate,
		FacilityID:              *payload.Client.DefaultFacility.ID,
		ClientCounselled:        payload.Client.ClientCounselled,
		Active:                  payload.Client.Active,
		ProgramID:               payload.Client.ProgramID,
		OrganisationID:          payload.Client.OrganisationID,
	}

	client, err := d.create.RegisterExistingUserAsClient(ctx, identifier, clientProfile)
	if err != nil {
		return nil, err
	}

	var clientTypes []enums.ClientType
	for _, k := range clientProfile.ClientTypes {
		clientTypes = append(clientTypes, enums.ClientType(k))
	}

	return &domain.ClientProfile{
		ID:                      clientProfile.ID,
		Active:                  clientProfile.Active,
		ClientTypes:             clientTypes,
		TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
		UserID:                  *client.UserID,
		ClientCounselled:        clientProfile.ClientCounselled,
		DefaultFacility: &domain.Facility{
			ID: &clientProfile.FacilityID,
		},
		OrganisationID: clientProfile.OrganisationID,
		ProgramID:      clientProfile.ProgramID,
	}, nil
}

// RegisterCaregiver registers a new caregiver on the platform
func (d *MyCareHubDb) RegisterCaregiver(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
	user := &gorm.User{
		Username:              input.User.Username,
		Name:                  input.User.Name,
		Gender:                input.User.Gender,
		DateOfBirth:           input.User.DateOfBirth,
		Active:                input.User.Active,
		CurrentProgramID:      input.User.CurrentProgramID,
		CurrentOrganisationID: input.User.CurrentOrganizationID,
	}

	contact := &gorm.Contact{
		Type:           input.Contact.ContactType,
		Value:          input.Contact.ContactValue,
		Active:         input.Contact.Active,
		OptedIn:        input.Contact.Active,
		OrganisationID: input.Contact.OrganisationID,
	}

	caregiver := &gorm.Caregiver{
		Active:          input.Caregiver.Active,
		CaregiverNumber: input.Caregiver.CaregiverNumber,
		OrganisationID:  input.Caregiver.OrganisationID,
	}

	err := d.create.RegisterCaregiver(ctx, user, contact, caregiver)
	if err != nil {
		return nil, err
	}

	profile := domain.CaregiverProfile{
		ID: caregiver.ID,
		User: domain.User{
			ID:               user.UserID,
			Username:         user.Username,
			Name:             user.Name,
			Gender:           user.Gender,
			Active:           user.Active,
			CurrentProgramID: user.CurrentProgramID,
		},
		CaregiverNumber: caregiver.CaregiverNumber,
	}

	return &profile, nil
}

// RegisterExistingUserAsCaregiver registers an existing user as a caregiver
func (d *MyCareHubDb) RegisterExistingUserAsCaregiver(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
	caregiver := &gorm.Caregiver{
		UserID:          input.Caregiver.UserID,
		Active:          input.Caregiver.Active,
		CaregiverNumber: input.Caregiver.CaregiverNumber,
		OrganisationID:  input.Caregiver.OrganisationID,
	}

	caregiver, err := d.create.RegisterExistingUserAsCaregiver(ctx, caregiver)
	if err != nil {
		return nil, err
	}

	user, err := d.query.GetUserProfileByUserID(ctx, &input.Caregiver.UserID)
	if err != nil {
		return nil, err
	}

	profile := domain.CaregiverProfile{
		ID:              caregiver.ID,
		User:            *createMapUser(user),
		CaregiverNumber: caregiver.CaregiverNumber,
	}

	return &profile, nil
}

// CreateCaregiver creates a caregiver record using the provided input
func (d *MyCareHubDb) CreateCaregiver(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error) {
	cgv := &gorm.Caregiver{
		Active:          caregiver.Active,
		CaregiverNumber: caregiver.CaregiverNumber,
		UserID:          caregiver.UserID,
	}

	err := d.create.CreateCaregiver(ctx, cgv)
	if err != nil {
		return nil, err
	}

	return &domain.Caregiver{
		ID:              cgv.ID,
		UserID:          cgv.UserID,
		CaregiverNumber: cgv.CaregiverNumber,
		Active:          cgv.Active,
	}, nil
}

// CreateIdentifier creates a new identifier
func (d *MyCareHubDb) CreateIdentifier(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
	i := &gorm.Identifier{
		Active:              true,
		IdentifierType:      identifier.IdentifierType,
		IdentifierValue:     identifier.IdentifierValue,
		IdentifierUse:       identifier.IdentifierUse,
		Description:         identifier.Description,
		IsPrimaryIdentifier: identifier.IsPrimaryIdentifier,
		ProgramID:           identifier.ProgramID,
	}

	err := d.create.CreateIdentifier(ctx, i)
	if err != nil {
		return nil, err
	}

	return &domain.Identifier{
		ID:                  i.ID,
		IdentifierType:      i.IdentifierType,
		IdentifierValue:     i.IdentifierValue,
		IdentifierUse:       i.IdentifierUse,
		Description:         i.Description,
		ValidFrom:           i.ValidFrom,
		ValidTo:             i.ValidTo,
		IsPrimaryIdentifier: i.IsPrimaryIdentifier,
		ProgramID:           i.ProgramID,
	}, nil
}

// SaveNotification saves a notification in the database
func (d *MyCareHubDb) SaveNotification(ctx context.Context, payload *domain.Notification) error {
	notification := &gorm.Notification{
		Active:     true,
		Title:      payload.Title,
		Body:       payload.Body,
		Type:       payload.Type.String(),
		IsRead:     false,
		UserID:     payload.UserID,
		FacilityID: payload.FacilityID,
		ProgramID:  payload.ProgramID,
	}
	return d.create.CreateNotification(ctx, notification)
}

// CreateUserSurveys creates a new user survey
func (d *MyCareHubDb) CreateUserSurveys(ctx context.Context, surveys []*dto.UserSurveyInput) error {
	var userSurveys []*gorm.UserSurvey

	for _, survey := range surveys {
		userSurveys = append(userSurveys, &gorm.UserSurvey{
			Active:      true,
			Link:        survey.Link,
			Title:       survey.Title,
			Description: survey.Description,
			UserID:      survey.UserID,
			FormID:      survey.FormID,
			ProjectID:   survey.ProjectID,
			LinkID:      survey.LinkID,
			Token:       survey.Token,
			ProgramID:   survey.ProgramID,
		})
	}

	return d.create.CreateUserSurveys(ctx, userSurveys)
}

// CreateMetric saves a metric to the database
func (d *MyCareHubDb) CreateMetric(ctx context.Context, payload *domain.Metric) error {
	event, err := json.Marshal(payload.Event)
	if err != nil {
		return fmt.Errorf("failed to marshal meta data: %v", err)
	}

	metric := &gorm.Metric{
		Active:    true,
		UserID:    payload.UserID,
		Timestamp: payload.Timestamp,
		Type:      payload.Type,
		Payload:   string(event),
	}

	return d.create.CreateMetric(ctx, metric)
}

// SaveFeedback saves a feedback to the database
func (d *MyCareHubDb) SaveFeedback(ctx context.Context, payload *domain.FeedbackResponse) error {
	feedback := &gorm.Feedback{
		UserID:            payload.UserID,
		FeedbackType:      payload.FeedbackType.String(),
		SatisfactionLevel: payload.SatisfactionLevel,
		ServiceName:       payload.ServiceName,
		Feedback:          payload.Feedback,
		RequiresFollowUp:  payload.RequiresFollowUp,
		PhoneNumber:       payload.PhoneNumber,
		ProgramID:         payload.ProgramID,
		OrganisationID:    payload.OrganisationID,
	}

	return d.create.SaveFeedback(ctx, feedback)
}

// RegisterStaff registers a new staff member into the portal
func (d *MyCareHubDb) RegisterStaff(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
	user := &gorm.User{
		Username:              payload.UserProfile.Username,
		Name:                  payload.UserProfile.Name,
		Gender:                payload.UserProfile.Gender,
		DateOfBirth:           payload.UserProfile.DateOfBirth,
		Active:                payload.UserProfile.Active,
		CurrentProgramID:      payload.UserProfile.CurrentProgramID,
		CurrentOrganisationID: payload.UserProfile.CurrentOrganizationID,
	}

	contact := &gorm.Contact{
		Type:           payload.Phone.ContactType,
		Value:          payload.Phone.ContactValue,
		Active:         payload.Phone.Active,
		OptedIn:        payload.Phone.Active,
		OrganisationID: payload.Phone.OrganisationID,
	}

	identifier := &gorm.Identifier{
		IdentifierType:      payload.StaffIdentifier.IdentifierType,
		IdentifierValue:     payload.StaffIdentifier.IdentifierValue,
		IdentifierUse:       payload.StaffIdentifier.IdentifierUse,
		Description:         payload.StaffIdentifier.Description,
		IsPrimaryIdentifier: payload.StaffIdentifier.IsPrimaryIdentifier,
		Active:              payload.StaffIdentifier.Active,
		ProgramID:           payload.StaffIdentifier.ProgramID,
		OrganisationID:      payload.StaffIdentifier.OrganisationID,
	}

	staffProfile := &gorm.StaffProfile{
		Active:            payload.Staff.Active,
		StaffNumber:       payload.Staff.StaffNumber,
		DefaultFacilityID: *payload.Staff.DefaultFacility.ID,
		ProgramID:         payload.Staff.ProgramID,
		OrganisationID:    payload.Staff.OrganisationID,
	}

	staff, err := d.create.RegisterStaff(ctx, user, contact, identifier, staffProfile)
	if err != nil {
		return nil, err
	}

	return &domain.StaffProfile{
		ID:          staff.ID,
		UserID:      staff.UserID,
		Active:      staff.Active,
		StaffNumber: staff.StaffNumber,
		DefaultFacility: &domain.Facility{
			ID: &staff.DefaultFacilityID,
		},
		User:           createMapUser(user),
		OrganisationID: staff.OrganisationID,
		ProgramID:      staff.ProgramID,
	}, nil
}

// RegisterExistingUserAsStaff is used to create a staff profile of an already existing user in a certain program
func (d *MyCareHubDb) RegisterExistingUserAsStaff(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
	staffProfile := &gorm.StaffProfile{
		Active:            payload.Staff.Active,
		StaffNumber:       payload.Staff.StaffNumber,
		DefaultFacilityID: *payload.Staff.DefaultFacility.ID,
		ProgramID:         payload.Staff.ProgramID,
		OrganisationID:    payload.Staff.OrganisationID,
		UserID:            payload.Staff.UserID,
	}

	identifier := &gorm.Identifier{
		Active:              payload.StaffIdentifier.Active,
		IdentifierType:      payload.StaffIdentifier.IdentifierType,
		IdentifierValue:     payload.StaffIdentifier.IdentifierValue,
		IdentifierUse:       payload.StaffIdentifier.IdentifierUse,
		Description:         payload.StaffIdentifier.Description,
		IsPrimaryIdentifier: payload.StaffIdentifier.Active,
		OrganisationID:      payload.StaffIdentifier.OrganisationID,
		ProgramID:           payload.StaffIdentifier.ProgramID,
	}

	staff, err := d.create.RegisterExistingUserAsStaff(ctx, identifier, staffProfile)
	if err != nil {
		return nil, err
	}

	return &domain.StaffProfile{
		ID:          staff.ID,
		UserID:      staff.UserID,
		Active:      staff.Active,
		StaffNumber: staff.StaffNumber,
		DefaultFacility: &domain.Facility{
			ID: &staff.DefaultFacilityID,
		},
		OrganisationID: staff.OrganisationID,
		ProgramID:      staff.ProgramID,
	}, nil
}

// CreateScreeningTool maps the screening tool domain model to database model to create screening tools
func (d *MyCareHubDb) CreateScreeningTool(ctx context.Context, input *domain.ScreeningTool) error {
	questionnaire := &gorm.Questionnaire{
		Active:         input.Questionnaire.Active,
		Name:           input.Questionnaire.Name,
		Description:    input.Questionnaire.Description,
		ProgramID:      input.ProgramID,
		OrganisationID: input.OrganisationID,
	}

	err := d.create.CreateQuestionnaire(ctx, questionnaire)
	if err != nil {
		return err
	}

	clientTypes := pq.StringArray{}
	for _, t := range input.ClientTypes {
		clientTypes = append(clientTypes, t.String())
	}
	genders := pq.StringArray{}
	for _, g := range input.Genders {
		genders = append(genders, strings.ToUpper(g.String()))
	}
	screeningtool := &gorm.ScreeningTool{
		Active:          input.Active,
		QuestionnaireID: questionnaire.ID,
		Threshold:       input.Threshold,
		ClientTypes:     clientTypes,
		Genders:         genders,
		MinimumAge:      input.AgeRange.LowerBound,
		MaximumAge:      input.AgeRange.UpperBound,
		ProgramID:       input.ProgramID,
		OrganisationID:  input.OrganisationID,
	}

	err = d.create.CreateScreeningTool(ctx, screeningtool)
	if err != nil {
		return err
	}

	for _, q := range input.Questionnaire.Questions {
		question := &gorm.Question{
			Active:            q.Active,
			QuestionnaireID:   questionnaire.ID,
			Text:              q.Text,
			QuestionType:      q.QuestionType.String(),
			ResponseValueType: q.ResponseValueType.String(),
			SelectMultiple:    q.SelectMultiple,
			Required:          q.Required,
			Sequence:          q.Sequence,
			ProgramID:         q.ProgramID,
			OrganisationID:    q.OrganisationID,
		}
		err := d.create.CreateQuestion(ctx, question)
		if err != nil {
			return err
		}
		for _, c := range q.Choices {
			choice := &gorm.QuestionInputChoice{
				Active:         c.Active,
				QuestionID:     question.ID,
				Choice:         c.Choice,
				Value:          c.Value,
				Score:          c.Score,
				ProgramID:      c.ProgramID,
				OrganisationID: c.OrganisationID,
			}
			err := d.create.CreateQuestionChoice(ctx, choice)
			if err != nil {
				return err
			}
		}

	}

	return nil

}

// CreateScreeningToolResponse saves a screening tool response to the database
func (d *MyCareHubDb) CreateScreeningToolResponse(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error) {
	screeningToolResponse := &gorm.ScreeningToolResponse{
		Active:          input.Active,
		ScreeningToolID: input.ScreeningToolID,
		FacilityID:      input.FacilityID,
		ClientID:        input.ClientID,
		AggregateScore:  input.AggregateScore,
		ProgramID:       input.ProgramID,
		OrganisationID:  input.OrganisationID,
	}

	screeningToolQuestionResponses := []*gorm.ScreeningToolQuestionResponse{}
	for _, q := range input.QuestionResponses {
		screeningToolQuestionResponses = append(screeningToolQuestionResponses, &gorm.ScreeningToolQuestionResponse{
			Active:                  q.Active,
			ScreeningToolResponseID: screeningToolResponse.ID,
			QuestionID:              q.QuestionID,
			Response:                q.Response,
			Score:                   q.Score,
			ProgramID:               q.ProgramID,
			OrganisationID:          q.OrganisationID,
			FacilityID:              q.FacilityID,
		})
	}

	return d.create.CreateScreeningToolResponse(ctx, screeningToolResponse, screeningToolQuestionResponses)

}

// AddCaregiverToClient is used to assign a caregiver to a client
func (d *MyCareHubDb) AddCaregiverToClient(ctx context.Context, clientCaregiver *domain.CaregiverClient) error {
	caregiverClient := &gorm.CaregiverClient{
		CaregiverID:      clientCaregiver.CaregiverID,
		ClientID:         clientCaregiver.ClientID,
		RelationshipType: clientCaregiver.RelationshipType,
		Active:           true,
		AssignedBy:       clientCaregiver.AssignedBy,
		ProgramID:        clientCaregiver.ProgramID,
		OrganisationID:   clientCaregiver.OrganisationID,
	}

	return d.create.AddCaregiverToClient(ctx, caregiverClient)
}

// CreateOrganisation is used to create a new organisation in the database
func (d *MyCareHubDb) CreateOrganisation(ctx context.Context, organisation *domain.Organisation, programs []*domain.Program) (*domain.Organisation, error) {
	organisationObj := &gorm.Organisation{
		Active:          organisation.Active,
		Code:            organisation.Code,
		Name:            organisation.Name,
		Description:     organisation.Description,
		EmailAddress:    organisation.EmailAddress,
		PhoneNumber:     organisation.PhoneNumber,
		PostalAddress:   organisation.PostalAddress,
		PhysicalAddress: organisation.PhysicalAddress,
		DefaultCountry:  organisation.DefaultCountry,
	}

	org, err := d.create.CreateOrganisation(ctx, organisationObj)
	if err != nil {
		return nil, err
	}

	progs := []*domain.Program{}
	for _, program := range programs {
		prog, err := d.CreateProgram(ctx, &dto.ProgramInput{
			Name:           program.Name,
			Description:    program.Description,
			OrganisationID: *org.ID,
		})
		if err != nil {
			return nil, err
		}
		progs = append(progs, prog)
	}

	return &domain.Organisation{
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
		Programs:        progs,
	}, nil
}

// CreateProgram enables the creation of a new program
func (d *MyCareHubDb) CreateProgram(ctx context.Context, input *dto.ProgramInput) (*domain.Program, error) {
	programInput := &gorm.Program{
		Active:         true,
		Name:           input.Name,
		Description:    input.Description,
		OrganisationID: input.OrganisationID,
	}

	program, err := d.create.CreateProgram(ctx, programInput)
	if err != nil {
		return nil, err
	}

	return &domain.Program{
		ID:          program.ID,
		Active:      program.Active,
		Name:        program.Name,
		Description: program.Description,
		Organisation: domain.Organisation{
			ID: program.OrganisationID,
		},
	}, nil
}

// AddFacilityToProgram is used to add a facility to a program which the currently logged in staff member belongs to.
func (d *MyCareHubDb) AddFacilityToProgram(ctx context.Context, programID string, facilityIDs []string) ([]*domain.Facility, error) {
	err := d.create.AddFacilityToProgram(ctx, programID, facilityIDs)
	if err != nil {
		return nil, err
	}

	var facilities []*domain.Facility
	records, err := d.query.GetProgramFacilities(ctx, programID)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		facilities = append(facilities, &domain.Facility{
			ID: &record.FacilityID,
		})
	}

	return facilities, nil
}

// CreateFacilities inserts multiple facility records in the database
func (d *MyCareHubDb) CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
	facilitiesObj := []*gorm.Facility{}
	for _, facility := range facilities {
		facilitiesObj = append(facilitiesObj, &gorm.Facility{
			Name:        facility.Name,
			Active:      facility.Active,
			Country:     facility.Country,
			Phone:       facility.Phone,
			Description: facility.Description,
			Identifier: gorm.FacilityIdentifier{
				Active: facility.Identifier.Active,
				Type:   string(facility.Identifier.Type),
				Value:  facility.Identifier.Value,
			},
		})
	}

	output, err := d.create.CreateFacilities(ctx, facilitiesObj)
	if err != nil {
		return nil, err
	}

	result := []*domain.Facility{}
	for _, facility := range output {
		result = append(result, d.mapFacilityObjectToDomain(facility, &facility.Identifier))
	}

	return result, nil
}
