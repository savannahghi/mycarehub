package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *MyCareHubDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if err := facility.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("facility input validation failed: %s", err)
	}

	facilityObj := &gorm.Facility{
		Name:               facility.Name,
		Code:               facility.Code,
		Active:             facility.Active,
		County:             facility.County,
		Phone:              facility.Phone,
		Description:        facility.Description,
		FHIROrganisationID: facility.FHIROrganisationID,
	}

	facilitySession, err := d.create.GetOrCreateFacility(ctx, facilityObj)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// SaveTemporaryUserPin does the actual saving of the users PIN in the database
func (d *MyCareHubDb) SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	pinObj := &gorm.PINData{
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: pinData.ValidFrom,
		ValidTo:   pinData.ValidTo,
		IsValid:   pinData.IsValid,
		Flavour:   pinData.Flavour,
		Salt:      pinData.Salt,
	}

	_, err := d.create.SaveTemporaryUserPin(ctx, pinObj)
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		Flavour:   pinInput.Flavour,
		Salt:      pinInput.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}

// SaveOTP saves the otp to the database
func (d *MyCareHubDb) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	otpObject := &gorm.UserOTP{
		UserID:      otpInput.UserID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		PhoneNumber: otpInput.PhoneNumber,
		Flavour:     otpInput.Flavour,
		OTP:         otpInput.OTP,
	}

	err := d.create.SaveOTP(ctx, otpObject)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to save OTP")
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
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to save security question response data")
	}

	return nil
}

// CreateHealthDiaryEntry is used to add a health diary record to the database.
func (d *MyCareHubDb) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
	healthDiaryResponse := &gorm.ClientHealthDiaryEntry{
		Active:                healthDiaryInput.Active,
		Mood:                  healthDiaryInput.Mood,
		Note:                  healthDiaryInput.Note,
		EntryType:             healthDiaryInput.EntryType,
		ShareWithHealthWorker: healthDiaryInput.ShareWithHealthWorker,
		SharedAt:              healthDiaryInput.SharedAt,
		ClientID:              healthDiaryInput.ClientID,
	}

	err := d.create.CreateHealthDiaryEntry(ctx, healthDiaryResponse)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// CreateServiceRequest creates  a service request which will be handled by a staff user.
// This happens in a transaction because we do not want to
// create a health diary entry without a subsequent service request when the client's mood is "VERY_BAD"
func (d *MyCareHubDb) CreateServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	meta, err := json.Marshal(serviceRequestInput.Meta)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to marshal meta data: %v", err)
	}
	serviceRequest := &gorm.ClientServiceRequest{
		Active:      serviceRequestInput.Active,
		RequestType: serviceRequestInput.RequestType,
		Request:     serviceRequestInput.Request,
		Status:      serviceRequestInput.Status,
		ClientID:    serviceRequestInput.ClientID,
		FacilityID:  serviceRequestInput.FacilityID,
		Meta:        string(meta),
	}

	err = d.create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// CreateStaffServiceRequest creates a new service request for the specified staff
func (d *MyCareHubDb) CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
	meta, err := json.Marshal(serviceRequestInput.Meta)
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
	}

	err = d.create.CreateStaffServiceRequest(ctx, serviceRequest)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// CreateClientCaregiver creates a client's caregiver
func (d *MyCareHubDb) CreateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	caregiver := &gorm.Caregiver{
		FirstName:     caregiverInput.FirstName,
		LastName:      caregiverInput.LastName,
		PhoneNumber:   caregiverInput.PhoneNumber,
		CaregiverType: caregiverInput.CaregiverType,
	}

	err := d.create.CreateClientCaregiver(ctx, caregiverInput.ClientID, caregiver)
	if err != nil {
		return err
	}

	return nil
}

// CreateCommunity creates a channel in the database
func (d *MyCareHubDb) CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {

	var genderList pq.StringArray
	for _, g := range communityInput.Gender {
		genderList = append(genderList, string(*g))
	}

	var clientTypeList pq.StringArray
	for _, c := range communityInput.ClientType {
		clientTypeList = append(clientTypeList, string(*c))
	}

	input := &gorm.Community{
		Name:         communityInput.Name,
		Description:  communityInput.Description,
		Active:       true,
		MinimumAge:   communityInput.AgeRange.LowerBound,
		MaximumAge:   communityInput.AgeRange.UpperBound,
		Gender:       genderList,
		ClientTypes:  clientTypeList,
		InviteOnly:   communityInput.InviteOnly,
		Discoverable: true,
	}

	channel, err := d.create.CreateCommunity(ctx, input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var genders []enumutils.Gender
	for _, k := range channel.Gender {
		genders = append(genders, enumutils.Gender(k))
	}

	var clientTypes []enums.ClientType
	for _, k := range channel.ClientTypes {
		clientTypes = append(clientTypes, enums.ClientType(k))
	}

	return &domain.Community{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		AgeRange: &domain.AgeRange{
			LowerBound: channel.MinimumAge,
			UpperBound: channel.MaximumAge,
		},
		Gender:     genders,
		ClientType: clientTypes,
		InviteOnly: channel.InviteOnly,
	}, nil
}

// GetOrCreateNextOfKin creates a related person who is a next of kin
func (d *MyCareHubDb) GetOrCreateNextOfKin(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {

	pn := &gorm.RelatedPerson{
		FirstName:        person.Name,
		RelationshipType: "NEXT_OF_KIN",
	}

	return d.create.GetOrCreateNextOfKin(ctx, pn, clientID, contactID)
}

// GetOrCreateContact creates a contact
func (d *MyCareHubDb) GetOrCreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {

	ct := &gorm.Contact{
		Active:       true,
		ContactType:  contact.ContactType,
		ContactValue: contact.ContactValue,
		UserID:       contact.UserID,
		Flavour:      contact.Flavour,
		OptedIn:      contact.OptedIn,
	}

	c, err := d.create.GetOrCreateContact(ctx, ct)
	if err != nil {
		return nil, err
	}

	return &domain.Contact{
		ID:           c.ContactID,
		ContactType:  *c.ContactID,
		ContactValue: c.ContactValue,
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
		}
		screeningToolResponsesObj = append(screeningToolResponsesObj, stq)
	}
	err := d.create.AnswerScreeningToolQuestions(ctx, screeningToolResponsesObj)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}
	return nil
}

// CreateUser creates a new user
func (d *MyCareHubDb) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {

	u := &gorm.User{
		Active:      true,
		Username:    user.Username,
		Name:        user.Name,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		UserType:    user.UserType,
		Flavour:     user.Flavour,
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
		FacilityID:              client.FacilityID,
		ClientCounselled:        client.ClientCounselled,
		ClientTypes:             clientTypes,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
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
		TreatmentBuddy:          c.TreatmentBuddy,
		ClientCounselled:        c.ClientCounselled,
		OrganisationID:          c.OrganisationID,
		FacilityID:              c.FacilityID,
		CHVUserID:               c.CHVUserID,
	}, nil
}

// RegisterClient registers a client in the database
func (d *MyCareHubDb) RegisterClient(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
	contact := &gorm.Contact{
		ContactType:  payload.Phone.ContactType,
		ContactValue: payload.Phone.ContactValue,
		Active:       payload.Phone.Active,
		OptedIn:      payload.Phone.Active,
		UserID:       payload.Phone.UserID,
		Flavour:      payload.Phone.Flavour,
	}

	identifier := &gorm.Identifier{
		IdentifierType:      payload.ClientIdentifier.IdentifierType,
		IdentifierValue:     payload.ClientIdentifier.IdentifierValue,
		IdentifierUse:       payload.ClientIdentifier.IdentifierUse,
		Description:         payload.ClientIdentifier.Description,
		IsPrimaryIdentifier: payload.ClientIdentifier.IsPrimaryIdentifier,
		Active:              payload.ClientIdentifier.Active,
	}

	var pgClientTypes pq.StringArray
	for _, clientType := range payload.Client.ClientTypes {
		pgClientTypes = append(pgClientTypes, clientType.String())
	}
	clientProfile := &gorm.Client{
		UserID:                  &payload.Client.UserID,
		ClientTypes:             pgClientTypes,
		TreatmentEnrollmentDate: payload.Client.TreatmentEnrollmentDate,
		FacilityID:              payload.Client.FacilityID,
		ClientCounselled:        payload.Client.ClientCounselled,
		Active:                  payload.Client.Active,
	}

	err := d.create.RegisterClient(ctx, contact, identifier, clientProfile)
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
		UserID:                  *clientProfile.UserID,
		TreatmentEnrollmentDate: clientProfile.TreatmentEnrollmentDate,
		TreatmentBuddy:          clientProfile.TreatmentBuddy,
		ClientCounselled:        clientProfile.ClientCounselled,
		FacilityID:              clientProfile.FacilityID,
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
	}, nil
}

// SaveNotification saves a notification in the database
func (d *MyCareHubDb) SaveNotification(ctx context.Context, payload *domain.Notification) error {
	notification := &gorm.Notification{
		Active:     true,
		Title:      payload.Title,
		Body:       payload.Body,
		Type:       payload.Type.String(),
		Flavour:    payload.Flavour,
		IsRead:     false,
		UserID:     payload.UserID,
		FacilityID: payload.FacilityID,
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
		})
	}

	return d.create.CreateUserSurveys(ctx, userSurveys)
}

// CreateMetric saves a metric to the database
func (d *MyCareHubDb) CreateMetric(ctx context.Context, payload *domain.Metric) error {
	event, err := json.Marshal(payload.Event)
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
	}

	return d.create.SaveFeedback(ctx, feedback)
}

// RegisterStaff registers a new staff member into the portal
func (d *MyCareHubDb) RegisterStaff(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
	usr := &gorm.User{
		Username:    payload.UserProfile.Username,
		Name:        payload.UserProfile.Name,
		Gender:      payload.UserProfile.Gender,
		DateOfBirth: payload.UserProfile.DateOfBirth,
		UserType:    payload.UserProfile.UserType,
		Flavour:     payload.UserProfile.Flavour,
		Active:      payload.UserProfile.Active,
	}

	contact := &gorm.Contact{
		ContactType:  payload.Phone.ContactType,
		ContactValue: payload.Phone.ContactValue,
		Active:       payload.Phone.Active,
		OptedIn:      payload.Phone.Active,
		Flavour:      payload.Phone.Flavour,
	}

	identifier := &gorm.Identifier{
		IdentifierType:      payload.StaffIdentifier.IdentifierType,
		IdentifierValue:     payload.StaffIdentifier.IdentifierValue,
		IdentifierUse:       payload.StaffIdentifier.IdentifierUse,
		Description:         payload.StaffIdentifier.Description,
		IsPrimaryIdentifier: payload.StaffIdentifier.IsPrimaryIdentifier,
		Active:              payload.StaffIdentifier.Active,
	}

	staffProfile := &gorm.StaffProfile{
		Active:            payload.Staff.Active,
		StaffNumber:       payload.Staff.StaffNumber,
		DefaultFacilityID: payload.Staff.DefaultFacilityID,
	}

	staff, err := d.create.RegisterStaff(ctx, usr, contact, identifier, staffProfile)
	if err != nil {
		return nil, err
	}

	return &domain.StaffProfile{
		ID:                staff.ID,
		UserID:            staff.UserID,
		Active:            staff.Active,
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
	}, nil
}

// CreateScreeningTool maps the screening tool domain model to database model to create screening tools
func (d *MyCareHubDb) CreateScreeningTool(ctx context.Context, input *domain.ScreeningTool) error {
	questionnaire := &gorm.Questionnaire{
		Active:      input.Questionnaire.Active,
		Name:        input.Questionnaire.Name,
		Description: input.Questionnaire.Description,
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
		genders = append(genders, g.String())
	}
	screeningtool := &gorm.ScreeningTool{
		Active:          input.Active,
		QuestionnaireID: questionnaire.ID,
		Threshold:       input.Threshold,
		ClientTypes:     clientTypes,
		Genders:         genders,
		MinimumAge:      input.AgeRange.LowerBound,
		MaximumAge:      input.AgeRange.UpperBound,
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
		}
		err := d.create.CreateQuestion(ctx, question)
		if err != nil {
			return err
		}
		for _, c := range q.Choices {
			choice := &gorm.QuestionInputChoice{
				Active:     c.Active,
				QuestionID: question.ID,
				Choice:     c.Choice,
				Value:      c.Value,
				Score:      c.Score,
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
	}

	screeningToolQuestionResponses := []*gorm.ScreeningToolQuestionResponse{}
	for _, q := range input.QuestionResponses {
		screeningToolQuestionResponses = append(screeningToolQuestionResponses, &gorm.ScreeningToolQuestionResponse{
			Active:                  q.Active,
			ScreeningToolResponseID: screeningToolResponse.ID,
			QuestionID:              q.QuestionID,
			Response:                q.Response,
			Score:                   q.Score,
		})
	}

	return d.create.CreateScreeningToolResponse(ctx, screeningToolResponse, screeningToolQuestionResponses)

}
