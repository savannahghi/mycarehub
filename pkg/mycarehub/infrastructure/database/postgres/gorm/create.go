package gorm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility, identifier *FacilityIdentifier) (*Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinPayload *PINData) (bool, error)
	SavePin(ctx context.Context, pinData *PINData) (bool, error)
	SaveOTP(ctx context.Context, otpInput *UserOTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*SecurityQuestionResponse) error
	CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *ClientHealthDiaryEntry) (*ClientHealthDiaryEntry, error)
	CreateServiceRequest(ctx context.Context, serviceRequestInput *ClientServiceRequest) error
	CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *StaffServiceRequest) error
	CreateCommunity(ctx context.Context, community *Community) (*Community, error)
	GetOrCreateNextOfKin(ctx context.Context, person *RelatedPerson, clientID, contactID string) error
	GetOrCreateContact(ctx context.Context, contact *Contact) (*Contact, error)
	CreateAppointment(ctx context.Context, appointment *Appointment) error
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*ScreeningToolsResponse) error
	CreateUser(ctx context.Context, user *User) error
	CreateClient(ctx context.Context, client *Client, contactID, identifierID string) error
	CreateIdentifier(ctx context.Context, identifier *Identifier) error
	CreateNotification(ctx context.Context, notification *Notification) error
	CreateUserSurveys(ctx context.Context, userSurvey []*UserSurvey) error
	CreateMetric(ctx context.Context, metric *Metric) error
	RegisterStaff(ctx context.Context, user *User, contact *Contact, identifier *Identifier, staffProfile *StaffProfile) (*StaffProfile, error)
	SaveFeedback(ctx context.Context, feedback *Feedback) error
	RegisterClient(ctx context.Context, user *User, contact *Contact, identifier *Identifier, client *Client) (*Client, error)
	RegisterCaregiver(ctx context.Context, user *User, contact *Contact, caregiver *Caregiver) error
	CreateCaregiver(ctx context.Context, caregiver *Caregiver) error
	CreateQuestionnaire(ctx context.Context, input *Questionnaire) error
	CreateScreeningTool(ctx context.Context, input *ScreeningTool) error
	CreateQuestion(ctx context.Context, input *Question) error
	CreateQuestionChoice(ctx context.Context, input *QuestionInputChoice) error
	CreateScreeningToolResponse(ctx context.Context, screeningToolResponse *ScreeningToolResponse, screeningToolQuestionResponses []*ScreeningToolQuestionResponse) (*string, error)
	AddCaregiverToClient(ctx context.Context, clientCaregiver *CaregiverClient) error
	CreateProgram(ctx context.Context, program *Program) error
	CreateOrganisation(ctx context.Context, organization *Organisation) error
	AddFacilityToProgram(ctx context.Context, programID string, facilityIDs []string) error
}

// GetOrCreateFacility is used to get or create a facility
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility, identifier *FacilityIdentifier) (*Facility, error) {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed initialize database transaction %v", err)
	}

	if facility.FHIROrganisationID == "" {
		tx = tx.Omit("fhir_organization_id")
	}

	err := tx.FirstOrCreate(facility).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	identifier.FacilityID = *facility.FacilityID
	if err := tx.FirstOrCreate(identifier).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a facility identifier: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to create/update security question responses failed: %v", err)
	}

	return facility, nil
}

// SaveTemporaryUserPin is used to save a temporary user pin
func (db *PGInstance) SaveTemporaryUserPin(ctx context.Context, pinPayload *PINData) (bool, error) {
	if pinPayload == nil {
		return false, fmt.Errorf("pinPayload must be provided")
	}
	err := db.DB.WithContext(ctx).Create(pinPayload).Error
	if err != nil {
		return false, fmt.Errorf("failed to save a pin: %v", err)
	}
	return true, nil
}

// SavePin saves the pin to the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.WithContext(ctx).Create(pinData).Error
	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}

	return true, nil
}

// SaveFeedback saves the feedback to the database
func (db *PGInstance) SaveFeedback(ctx context.Context, feedback *Feedback) error {
	err := db.DB.WithContext(ctx).Create(feedback).Error
	if err != nil {
		return fmt.Errorf("failed to save feedback: %v", err)
	}
	return nil
}

// SaveOTP saves the generated otp to the database
func (db *PGInstance) SaveOTP(ctx context.Context, otpInput *UserOTP) error {
	//Invalidate other OTPs before saving the new OTP by setting valid to false
	if otpInput.PhoneNumber == "" || !otpInput.Flavour.IsValid() {
		return fmt.Errorf("phone number cannot be empty")
	}

	err := db.DB.WithContext(ctx).Model(&UserOTP{}).Where(&UserOTP{PhoneNumber: otpInput.PhoneNumber, Flavour: otpInput.Flavour}).
		Updates(map[string]interface{}{"is_valid": false}).Error
	if err != nil {
		return fmt.Errorf("failed to update OTP data: %v", err)
	}

	//Save the OTP by setting valid to true
	err = db.DB.WithContext(ctx).Create(otpInput).Error
	if err != nil {
		return fmt.Errorf("failed to save otp data")
	}
	return nil
}

// SaveSecurityQuestionResponse saves the security question response to the database if it does not exist,
// otherwise it updates the existing one
func (db *PGInstance) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*SecurityQuestionResponse) error {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed initialize database transaction %v", err)
	}
	for _, questionResponse := range securityQuestionResponse {
		SaveSecurityQuestionResponseUpdatePayload := &SecurityQuestionResponse{
			Response: questionResponse.Response,
		}
		err := tx.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: questionResponse.UserID, QuestionID: questionResponse.QuestionID}).First(&questionResponse).Error
		if err == nil {
			err := tx.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: questionResponse.UserID, QuestionID: questionResponse.QuestionID}).Updates(&SaveSecurityQuestionResponseUpdatePayload).Error
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update security question response data: %v", err)
			}
		} else if err == gorm.ErrRecordNotFound {
			err = tx.Create(&questionResponse).Error
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create security question response data: %v", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get security question response data: %v", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create/update security question responses failed: %v", err)
	}

	return nil
}

// CreateHealthDiaryEntry records the health diary entries from a client. This is necessary for engagement with clients
// on a day-by-day basis
func (db *PGInstance) CreateHealthDiaryEntry(ctx context.Context, healthDiary *ClientHealthDiaryEntry) (*ClientHealthDiaryEntry, error) {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(healthDiary).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	return healthDiary, nil
}

// CreateServiceRequest creates a service request entry into the database. This step is reached only IF the client is
// in a VERY_BAD mood. We get this mood from the mood scale provided by the front end.
// This operation is done within a transaction to prevent a situation where a health entry can be recorded
// but a service request is not successfully created.
func (db *PGInstance) CreateServiceRequest(
	ctx context.Context,
	serviceRequestInput *ClientServiceRequest,
) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(serviceRequestInput).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// CreateStaffServiceRequest creates a new staff service request
func (db *PGInstance) CreateStaffServiceRequest(ctx context.Context, serviceRequestInput *StaffServiceRequest) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(serviceRequestInput).Error
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// CreateCommunity creates a channel in the database
func (db *PGInstance) CreateCommunity(ctx context.Context, community *Community) (*Community, error) {
	err := db.DB.WithContext(ctx).Create(community).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create a community: %v", err)
	}
	return community, nil
}

// GetOrCreateNextOfKin get or creates a related person in the database
// The client ID and contact ID are used to link the created person with a client
// and the associated contact for the person
func (db *PGInstance) GetOrCreateNextOfKin(ctx context.Context, person *RelatedPerson, clientID, contactID string) error {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where(person).FirstOrCreate(person).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get or create related person: %v", err)
	}

	// link contact
	contact := RelatedPersonContacts{
		RelatedPersonID: &person.ID,
		ContactID:       &contactID,
	}
	err = tx.Where(contact).FirstOrCreate(&contact).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get or create related person contact: %v", err)
	}

	// link client
	client := ClientRelatedPerson{
		ClientID:        &clientID,
		RelatedPersonID: &person.ID,
	}
	err = tx.Where(client).FirstOrCreate(&client).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get or create related person client: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create related person failed: %v", err)
	}

	return nil
}

// GetOrCreateContact creates a person's contact in the database if they do not exist or gets them if they already exist
func (db *PGInstance) GetOrCreateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where(Contact{Value: contact.Value, Flavour: contact.Flavour}).FirstOrCreate(contact).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create contact: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to create contact failed: %v", err)
	}

	return contact, nil
}

// CreateAppointment creates an appointment in the database
func (db *PGInstance) CreateAppointment(ctx context.Context, appointment *Appointment) error {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(appointment).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create an appointment: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create an appointment failed: %v", err)
	}

	return nil
}

// AnswerScreeningToolQuestions answers the screening tool questions
func (db *PGInstance) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*ScreeningToolsResponse) error {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, response := range screeningToolResponses {
		err := tx.Create(response).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create screening tool response: %v", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create screening tool responses failed: %v", err)
	}

	return nil
}

// CreateUser creates a new user
func (db *PGInstance) CreateUser(ctx context.Context, user *User) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit create user transaction: %w", err)
	}

	return nil
}

// CreateClient creates a new client
func (db *PGInstance) CreateClient(ctx context.Context, client *Client, contactID, identifierID string) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(client).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create client: %w", err)
	}

	// link contact
	// contact := ClientContacts{
	// 	ClientID:  client.ID,
	// 	ContactID: &contactID,
	// }
	// err = tx.Where(contact).FirstOrCreate(&contact).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("failed to get or create client contact: %w", err)
	// }

	// link identifiers
	identifier := ClientIdentifiers{
		ClientID:     client.ID,
		IdentifierID: &identifierID,
	}
	err = tx.Where(identifier).FirstOrCreate(&identifier).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get or create client identifier: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit create client transaction: %w", err)
	}

	return nil
}

// RegisterClient registers a client with the system
func (db *PGInstance) RegisterClient(ctx context.Context, user *User, contact *Contact, identifier *Identifier, client *Client) (*Client, error) {
	tx := db.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(user).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create user: %v", err)
	}

	// create contact
	contact.UserID = user.UserID
	err = tx.Where(Contact{Value: contact.Value, Flavour: contact.Flavour}).FirstOrCreate(contact).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create contact: %v", err)
	}

	// create identifier
	err = tx.Where(identifier).FirstOrCreate(identifier).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create identifier: %v", err)
	}

	// create client
	client.UserID = user.UserID
	err = tx.Create(client).First(&client).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// // link contact
	// clientContact := ClientContacts{
	// 	ClientID:  client.ID,
	// 	ContactID: contact.ContactID,
	// }
	// err = tx.Where(clientContact).FirstOrCreate(&clientContact).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, fmt.Errorf("failed to get or create client contact: %v", err)
	// }

	// link identifiers
	clientIdentifier := ClientIdentifiers{
		ClientID:     client.ID,
		IdentifierID: &identifier.ID,
	}
	err = tx.Where(clientIdentifier).FirstOrCreate(&clientIdentifier).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create client identifier: %v", err)
	}

	// Append client facilities
	clientFacilities := ClientFacilities{
		ClientID:   client.ID,
		FacilityID: &client.FacilityID,
	}
	err = tx.Where(clientFacilities).Create(&clientFacilities).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get client facilities: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit register client transaction: %v", err)
	}

	return client, nil
}

// RegisterCaregiver registers a new caregiver
func (db *PGInstance) RegisterCaregiver(ctx context.Context, user *User, contact *Contact, caregiver *Caregiver) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	contact.UserID = user.UserID
	err = tx.Create(contact).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user: %w", err)
	}

	caregiver.UserID = *user.UserID
	err = tx.Create(caregiver).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create caregiver: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit create caregiver transaction: %w", err)
	}

	return nil
}

// CreateCaregiver creates a caregiver record linked to a user
func (db *PGInstance) CreateCaregiver(ctx context.Context, caregiver *Caregiver) error {
	err := db.DB.WithContext(ctx).Create(caregiver).Error
	if err != nil {
		return fmt.Errorf("failed to create caregiver: %w", err)
	}

	return nil
}

// CreateIdentifier creates a new identifier
func (db *PGInstance) CreateIdentifier(ctx context.Context, identifier *Identifier) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(identifier).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create identifier: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit create identifier transaction: %w", err)
	}

	return nil
}

// CreateNotification saves a notification to the database
func (db *PGInstance) CreateNotification(ctx context.Context, notification *Notification) error {
	tx := db.DB.WithContext(ctx).Begin()

	err := tx.Create(notification).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create notification: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit create notification transaction: %w", err)
	}

	return nil
}

// CreateUserSurveys saves a user survey details including the survey link
func (db *PGInstance) CreateUserSurveys(ctx context.Context, userSurveys []*UserSurvey) error {
	if len(userSurveys) == 0 {
		return nil
	}

	err := db.DB.WithContext(ctx).Create(userSurveys).Error
	if err != nil {
		return fmt.Errorf("failed to create user survey: %w", err)
	}

	return nil
}

// CreateMetric saves a metric to the database
func (db *PGInstance) CreateMetric(ctx context.Context, metric *Metric) error {
	err := db.DB.WithContext(ctx).Create(metric).Error
	if err != nil {
		return fmt.Errorf("failed to create metric: %w", err)
	}

	return nil
}

// RegisterStaff registers a staff member to the database
func (db *PGInstance) RegisterStaff(ctx context.Context, user *User, contact *Contact, identifier *Identifier, staffProfile *StaffProfile) (*StaffProfile, error) {
	tx := db.DB.WithContext(ctx).Begin()

	// create user
	err := tx.Create(user).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// create contact
	contact.UserID = user.UserID
	err = tx.Where(Contact{Value: contact.Value, Flavour: contact.Flavour}).FirstOrCreate(contact).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create contact: %v", err)
	}

	// create identifier
	err = tx.Where(identifier).FirstOrCreate(identifier).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create identifier: %v", err)
	}

	// create staff profile
	staffProfile.UserID = *user.UserID
	err = tx.Create(staffProfile).FirstOrCreate(&staffProfile).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create staff profile: %w", err)
	}

	// link contact
	// contactLink := StaffContacts{
	// 	StaffID:   staffProfile.ID,
	// 	ContactID: contact.ContactID,
	// }
	// err = tx.Where(contactLink).FirstOrCreate(&contactLink).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, fmt.Errorf("failed to get or create staff contact: %w", err)
	// }

	// link identifier
	identifierLink := StaffIdentifiers{
		StaffID:      staffProfile.ID,
		IdentifierID: &identifier.ID,
	}
	err = tx.Where(identifierLink).FirstOrCreate(&identifierLink).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get or create staff identifier: %w", err)
	}

	// Append staff facilities
	staffFacilities := StaffFacilities{
		StaffID:    staffProfile.ID,
		FacilityID: &staffProfile.DefaultFacilityID,
	}
	err = tx.Where(staffFacilities).Create(&staffFacilities).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get staff facilities: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit create staff profile transaction: %w", err)
	}

	return staffProfile, nil
}

// CreateQuestionnaire saves a questionnaire to the database
func (db *PGInstance) CreateQuestionnaire(ctx context.Context, input *Questionnaire) error {
	if err := db.DB.WithContext(ctx).Create(&input).Error; err != nil {
		return fmt.Errorf("failed to create questionnaire: %w", err)
	}
	return nil
}

// CreateScreeningTool saves a screening tool to the database
func (db *PGInstance) CreateScreeningTool(ctx context.Context, input *ScreeningTool) error {
	if err := db.DB.WithContext(ctx).Create(&input).Error; err != nil {
		return fmt.Errorf("failed to create screening tool: %w", err)
	}
	return nil
}

// CreateQuestion saves a question to the database
func (db *PGInstance) CreateQuestion(ctx context.Context, input *Question) error {
	if err := db.DB.WithContext(ctx).Create(&input).Error; err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}
	return nil
}

// CreateQuestionChoice saves a question choice to the database
func (db *PGInstance) CreateQuestionChoice(ctx context.Context, input *QuestionInputChoice) error {
	if err := db.DB.WithContext(ctx).Create(&input).Error; err != nil {
		return fmt.Errorf("failed to create question choice: %w", err)
	}
	return nil
}

// CreateScreeningToolResponse saves a screening tool response to the database and returns the response ID
func (db *PGInstance) CreateScreeningToolResponse(ctx context.Context, screeningToolResponse *ScreeningToolResponse, screeningToolQuestionResponses []*ScreeningToolQuestionResponse) (*string, error) {
	tx := db.DB.WithContext(ctx).Begin()

	if err := tx.Create(screeningToolResponse).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create screening tool response: %w", err)
	}

	for _, questionResponse := range screeningToolQuestionResponses {
		questionResponse.ScreeningToolResponseID = screeningToolResponse.ID
		if err := tx.Create(questionResponse).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create screening tool question response: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit create screening tool response transaction: %w", err)
	}

	return &screeningToolResponse.ID, nil
}

// AddCaregiverToClient adds a caregiver to a client
func (db *PGInstance) AddCaregiverToClient(ctx context.Context, clientCaregiver *CaregiverClient) error {
	if err := db.DB.WithContext(ctx).Create(&clientCaregiver).Error; err != nil {
		return fmt.Errorf("failed to create client caregiver: %w", err)
	}

	return nil
}

// CreateOrganisation is used to create an organization into the database
func (db *PGInstance) CreateOrganisation(ctx context.Context, organization *Organisation) error {
	if err := db.DB.WithContext(ctx).Create(&organization).Error; err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

// CreateProgram adds a new program record
func (db *PGInstance) CreateProgram(ctx context.Context, program *Program) error {
	if err := db.DB.WithContext(ctx).Create(&program).Error; err != nil {
		return fmt.Errorf("failed to create program: %w", err)
	}
	return nil
}

// AddFacilityToProgram is used to add a facility to a program
func (db *PGInstance) AddFacilityToProgram(ctx context.Context, programID string, facilityIDs []string) error {
	for _, facilityID := range facilityIDs {
		programFacility := ProgramFacility{
			ProgramID:  programID,
			FacilityID: facilityID,
		}

		if err := db.DB.WithContext(ctx).Where(programFacility).FirstOrCreate(&programFacility).Error; err != nil {
			return fmt.Errorf("failed to create program facility: %w", err)
		}
	}

	return nil
}
