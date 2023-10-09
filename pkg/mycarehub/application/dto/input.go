package dto

import (
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
	validator "gopkg.in/go-playground/validator.v9"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name               string                  `json:"name" validate:"required,min=3,max=100"`
	Phone              string                  `json:"phone" validate:"required"`
	Active             bool                    `json:"active"`
	Country            string                  `json:"country" validate:"required"`
	County             string                  `json:"county" validate:"required"`
	Address            string                  `json:"address" validate:"required"`
	Description        string                  `json:"description" validate:"required,min=3,max=256"`
	FHIROrganisationID string                  `json:"fhirOrganisationID"`
	Identifier         FacilityIdentifierInput `json:"identifier" validate:"required"`
	Coordinates        CoordinatesInput        `json:"coordinates" validate:"required"`
}

// Validate helps with validation of facility input fields
func (f *FacilityInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// FacilityIdentifierInput is the identifier of the facility
type FacilityIdentifierInput struct {
	Type       enums.FacilityIdentifierType `json:"type" validate:"required"`
	Value      string                       `json:"value" validate:"required"`
	FacilityID string                       `json:"facilityID"`
}

// Validate helps with validation of facility identifier input fields
func (f *FacilityIdentifierInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// PaginationsInput contains fields required for pagination
type PaginationsInput struct {
	Limit       int        `json:"limit"`
	CurrentPage int        `json:"currentPage" validate:"required"`
	Sort        SortsInput `json:"sort"`
}

// Validate helps with validation of PaginationsInput fields
func (f *PaginationsInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// FiltersInput contains fields required for filtering
type FiltersInput struct {
	DataType enums.FilterSortDataType `json:"dataType" validate:"required"`
	Value    string                   `json:"value" validate:"required"` // TODO: Clear spec on validation e.g dates must be ISO 8601. This is the actual data being filtered
}

// Validate helps with validation of FiltersInput fields
func (f *FiltersInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// SortsInput includes the fields required for sorting the different types of fields
type SortsInput struct {
	Direction enums.SortDataType       `json:"direction"`
	Field     enums.FilterSortDataType `json:"field"`
}

// LoginInput represents the Login input data structure
type LoginInput struct {
	Username string          `json:"username" validate:"required"`
	PIN      string          `json:"pin" validate:"required"`
	Flavour  feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of LoginInput fields
func (f *LoginInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// PINInput represents the Pin input data structure
type PINInput struct {
	UserID     *string         `json:"userID" validate:"required"`
	PIN        *string         `json:"pin" validate:"required"`
	ConfirmPIN *string         `json:"confirmPIN" validate:"required"`
	Flavour    feedlib.Flavour `json:"flavour" validate:"required"`
}

// VerifyOTPInput represents the verify OTP input data structure
type VerifyOTPInput struct {
	PhoneNumber string          `json:"phoneNumber"`
	Username    string          `json:"username" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of VerifyOTPInput fields
func (f *VerifyOTPInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	if !f.Flavour.IsValid() {
		err = fmt.Errorf("invalid flavour provided: %v", f.Flavour)
	}

	if !converterandformatter.IsMSISDNValid(f.PhoneNumber) {
		err = fmt.Errorf("invalid phone provided: %v", f.PhoneNumber)
	}

	return err
}

// SendOTPInput represents the send OTP input data structure
type SendOTPInput struct {
	Username string          `json:"username" validate:"required"`
	Flavour  feedlib.Flavour `json:"flavour" validate:"required"`
}

// BasicUserInput is used to define the inputs needed carrying out an activity that requires either a username, phone number and flavour.
type BasicUserInput struct {
	Username string          `json:"username"`
	Flavour  feedlib.Flavour `json:"flavour" validate:"required"`
}

// SendRetryOTPPayload is used to define the inputs passed when calling the endpoint
// that resends an otp
type SendRetryOTPPayload struct {
	Username string          `json:"username" validate:"required"`
	Flavour  feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of PINInput fields
func (f *PINInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// SecurityQuestionResponseInput represents the SecurityQuestionResponse input data structure
type SecurityQuestionResponseInput struct {
	UserID             string `json:"userID" validate:"required"`
	SecurityQuestionID string `json:"securityQuestionID" validate:"required"`
	Response           string `json:"response" validate:"required"`
}

// Validate helps with validation of SecurityQuestionResponseInput fields
func (f *SecurityQuestionResponseInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// VerifySecurityQuestionInput defines the field passed when verifying the set security questions
type VerifySecurityQuestionInput struct {
	QuestionID string          `json:"questionID" validate:"required"`
	Flavour    feedlib.Flavour `json:"flavour" validate:"required"`
	Response   string          `json:"response" validate:"required"`
	Username   string          `json:"username" validate:"required"`
}

// VerifySecurityQuestionsPayload holds a list of security question inputs.
type VerifySecurityQuestionsPayload struct {
	SecurityQuestionsInput []*VerifySecurityQuestionInput `json:"verifySecurityQuestionsInput"`
}

// Validate checks to validate whether the field inputs for verifying a security question
// are filled
func (f *VerifySecurityQuestionInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// VerifyPhoneInput carries the OTP data used to send OTP messages to a particular phone number
type VerifyPhoneInput struct {
	Username string          `json:"username"`
	Flavour  feedlib.Flavour `json:"flavour"`
}

// GetSecurityQuestionsInput defines the field passed when getting the security questions
type GetSecurityQuestionsInput struct {
	Flavour feedlib.Flavour `json:"flavour" validate:"required"`
}

// GetUserRespondedSecurityQuestionsInput defines the field passed when getting the security questions
type GetUserRespondedSecurityQuestionsInput struct {
	Username string          `json:"username" validate:"required"`
	Flavour  feedlib.Flavour `json:"flavour" validate:"required"`
	OTP      string          `json:"otp" validate:"required"`
}

// Validate helps with validation of GetUserRespondedSecurityQuestionsInput fields
func (f *GetUserRespondedSecurityQuestionsInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// UserResetPinInput contains the fields requires when a user is resetting a pin
type UserResetPinInput struct {
	PhoneNumber string          `json:"phoneNumber"`
	Username    string          `json:"username" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
	PIN         string          `json:"pin" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
}

// Validate checks to validate whether the field inputs for verifying user pin
func (f *UserResetPinInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// ShareContentInput defines the field passed when sharing content
type ShareContentInput struct {
	ClientID  string `json:"clientID" validate:"required"`
	ContentID int    `json:"contentID" validate:"required"`
	Channel   string `json:"channel" validate:"required"`
}

// Validate helps with validation of ShareContentInput fields
func (f *ShareContentInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// RefreshTokenPayload is used when calling the REST API to
// exchange a Refresh Token for new ID Token
type RefreshTokenPayload struct {
	UserID *string `json:"userID"`
}

// FeedbackResponseInput defines the field passed when sending feedback
type FeedbackResponseInput struct {
	UserID            string             `json:"userID" validate:"required"`
	FeedbackType      enums.FeedbackType `json:"feedbackType" validate:"required"`
	SatisfactionLevel int                `json:"satisfactionLevel" validate:"required"`
	ServiceName       string             `json:"serviceName" validate:"required"`
	Feedback          string             `json:"feedback" validate:"required"`
	RequiresFollowUp  bool               `json:"requiresFollowUp" validate:"required"`
}

// FeedbackEmail defines the field to be parsed when sending feedback emails
type FeedbackEmail struct {
	User              string             `json:"user"`
	FeedbackType      enums.FeedbackType `json:"feedbackType"`
	SatisfactionLevel int                `json:"satisfactionLevel"`
	ServiceName       string             `json:"serviceName"`
	Feedback          string             `json:"feedback"`
	PhoneNumber       string             `json:"phoneNumber"`
	ProgramID         string             `json:"programID"`
}

// CaregiverInput defines the field passed when creating a caregiver
type CaregiverInput struct {
	Username        string                 `json:"username"`
	Name            string                 `json:"name"`
	Gender          enumutils.Gender       `json:"gender"`
	DateOfBirth     scalarutils.Date       `json:"dateOfBirth"`
	PhoneNumber     string                 `json:"phoneNumber"`
	CaregiverNumber string                 `json:"caregiverNumber"`
	SendInvite      bool                   `json:"sendInvite"`
	AssignedClients []ClientCaregiverInput `json:"assignedClients"`
}

// Validate helps with validation of CaregiverInput fields
func (f *CaregiverInput) Validate() error {
	v := validator.New()
	err := v.Struct(f)
	return err
}

// ClientRegistrationInput defines the fields passed as a payload to the client registration API
type ClientRegistrationInput struct {
	Username       string             `json:"username" validate:"required"`
	Facility       string             `json:"facility" validate:"required"`
	ClientTypes    []enums.ClientType `json:"clientTypes" validate:"required"`
	ClientName     string             `json:"clientName" validate:"required"`
	Gender         enumutils.Gender   `json:"gender" validate:"required"`
	DateOfBirth    scalarutils.Date   `json:"dateOfBirth" validate:"required"`
	PhoneNumber    string             `json:"phoneNumber" validate:"required"`
	EnrollmentDate scalarutils.Date   `json:"enrollmentDate" validate:"required"`
	CCCNumber      string             `json:"cccNumber" validate:"required"`
	Counselled     bool               `json:"counselled" validate:"required"`
	InviteClient   bool               `json:"inviteClient"`
	ProgramID      string             `json:"programID"`
}

// Validate helps with validation of ClientRegistrationInput fields
func (f *ClientRegistrationInput) Validate() error {
	v := validator.New()
	err := v.Struct(f)
	return err
}

// ExistingUserClientInput defines the fields passed as a payload to create a client profile of an already existing user
type ExistingUserClientInput struct {
	UserID         string             `json:"userID" validate:"required"`
	ProgramID      string             `json:"programID" validate:"required"`
	FacilityID     string             `json:"facilityID" validate:"required"`
	CCCNumber      *string            `json:"cccNumber" validate:"required"`
	ClientTypes    []enums.ClientType `json:"clientTypes" validate:"required"`
	EnrollmentDate scalarutils.Date   `json:"enrollmentDate" validate:"required"`
	Counselled     bool               `json:"counselled" validate:"required"`
	InviteClient   bool               `json:"inviteClient"`
}

// Validate helps with validation of ExistingUserClientRegistrationInput fields
func (e *ExistingUserClientInput) Validate() error {
	v := validator.New()
	err := v.Struct(e)
	return err
}

// CommunityInput defines the payload to create a channel
type CommunityInput struct {
	Name           string              `json:"name"`
	Topic          string              `json:"topic"`
	AgeRange       *AgeRangeInput      `json:"ageRange"`
	Gender         []*enumutils.Gender `json:"gender"`
	Visibility     enums.Visibility    `json:"visibility"`
	Preset         enums.Preset        `json:"preset"`
	ClientType     []*enums.ClientType `json:"clientType"`
	OrganisationID string              `json:"organisationID"`
	ProgramID      string              `json:"programID"`
	FacilityID     string              `json:"facilityID"`
}

// AgeRangeInput defines the channel users age input
type AgeRangeInput struct {
	LowerBound int `json:"lowerBound"`
	UpperBound int `json:"upperBound"`
}

// NextOfKinPayload defines the payload from KenyaEMR
// used for client registration
type NextOfKinPayload struct {
	Name         string `json:"name"`
	Contact      string `json:"contact"`
	Relationship string `json:"relationship"`
	ProgramID    string `json:"programID"`
}

// PatientRegistrationPayload defines the payload from KenyaEMR
// used for client registration
type PatientRegistrationPayload struct {
	MFLCode            string           `json:"MFLCODE"`
	CCCNumber          string           `json:"cccNumber"`
	Name               string           `json:"name"`
	DateOfBirth        scalarutils.Date `json:"dateOfBirth"`
	ClientType         enums.ClientType `json:"clientType"`
	PhoneNumber        string           `json:"phoneNumber"`
	EnrollmentDate     scalarutils.Date `json:"enrollmentDate"`
	BirthDateEstimated bool             `json:"birthDateEstimated"`
	Gender             string           `json:"gender"`
	Counselled         bool             `json:"counselled"`
	NextOfKin          NextOfKinPayload `json:"nextOfKin"`
	ProgramID          string           `json:"programID"`
}

// FetchHealthDiaryEntries models the payload that is passed when
// fetching the health diary entries that were recorded by patients assigned to
// the matching facility
type FetchHealthDiaryEntries struct {
	MFLCode      int        `json:"MFLCODE"`
	LastSyncTime *time.Time `json:"lastSyncTime"`
}

// PatientsPayload is the payload for registering patients
type PatientsPayload struct {
	Patients []*PatientRegistrationPayload `json:"patients"`
}

// PatientSyncPayload is the payload for polling newly created patients/clients
// since the last polling/sync time
type PatientSyncPayload struct {
	MFLCode  int        `json:"MFLCODE"`
	SyncTime *time.Time `json:"lastSyncTime"`
}

// ServiceRequestPayload defines the payload from KenyaEMR used to fetch
// service requests.
type ServiceRequestPayload struct {
	MFLCode      int        `json:"MFLCODE"`
	LastSyncTime *time.Time `json:"lastSyncTime"`
}

// FacilityAppointmentsPayload is the payload sent for creating/updating an appointment
type FacilityAppointmentsPayload struct {
	MFLCode      string               `json:"MFLCODE"`
	Appointments []AppointmentPayload `json:"appointments"`
}

// AppointmentPayload is the payload representing an appointment
type AppointmentPayload struct {
	CCCNumber         string           `json:"ccc_number"`
	ExternalID        string           `json:"appointment_id"`
	AppointmentDate   scalarutils.Date `json:"appointment_date"`
	AppointmentReason string           `json:"appointment_reason"`
}

// UpdateServiceRequestsPayload defined a list of service requests to synchronize MyCareHub with.
type UpdateServiceRequestsPayload struct {
	ServiceRequests []UpdateServiceRequestPayload `json:"serviceRequests" validate:"required"`
}

// UpdateServiceRequestPayload defines the payload that is used to synchronize KenyaEMR service requests to MyCareHub.
type UpdateServiceRequestPayload struct {
	ID           string    `json:"id" validate:"required"`
	RequestType  string    `json:"request_type" validate:"required"`
	Status       string    `json:"status" validate:"required"`
	InProgressAt time.Time `json:"in_progress_at"`
	InProgressBy string    `json:"in_progress_by"`
	ResolvedAt   time.Time `json:"resolved_at"`
	ResolvedBy   string    `json:"resolved_by"`
}

// StaffRegistrationInput is a model that represents the inputs passed when registering a staff user
type StaffRegistrationInput struct {
	Username            string           `json:"username" validate:"required"`
	Facility            string           `json:"facility" validate:"required"`
	StaffName           string           `json:"staffName" validate:"required"`
	Gender              enumutils.Gender `json:"gender" validate:"required"`
	DateOfBirth         scalarutils.Date `json:"dateOfBirth" validate:"required"`
	PhoneNumber         string           `json:"phoneNumber" validate:"required"`
	IDNumber            string           `json:"idNumber" validate:"required"`
	StaffNumber         string           `json:"staffNumber" validate:"required"`
	StaffRoles          string           `json:"staffRoles"`
	InviteStaff         bool             `json:"inviteStaff"`
	ProgramID           string           `json:"programID"`
	OrganisationID      string           `json:"organisationID"`
	IsOrganisationAdmin bool             `json:"isOrganisationAdmin"`
}

// Validate helps with validation of StaffRegistrationInput fields
func (s StaffRegistrationInput) Validate() error {
	var err error

	// try converting the ID number to an int
	_, err = strconv.Atoi(s.IDNumber)
	if err != nil {
		return fmt.Errorf("ID number must be an integer")
	}
	v := validator.New()

	err = v.Struct(s)

	return err
}

// ExistingUserStaffInput is a model that represents the inputs passed when registering an existing staff user to a program
type ExistingUserStaffInput struct {
	UserID      string  `json:"userID"`
	ProgramID   string  `json:"programID" validate:"required"`
	FacilityID  string  `json:"facilityID" validate:"required"`
	IDNumber    *string `json:"idNumber"`
	StaffNumber string  `json:"staffNumber" validate:"required"`
	StaffRoles  string  `json:"staffRoles"`
	InviteStaff bool    `json:"inviteStaff"`
}

// Validate helps with validation of StaffRegistrationInput fields
func (e ExistingUserStaffInput) Validate() error {
	var err error
	v := validator.New()

	err = v.Struct(e)

	return err
}

// PinResetServiceRequestPayload models the details passed to an API when a pin reset service request
// is being created
type PinResetServiceRequestPayload struct {
	CCCNumber string          `json:"cccNumber"`
	Username  string          `json:"username"`
	Flavour   feedlib.Flavour `json:"flavour"`
}

// ServiceRequestInput is a domain entity that represents a service request.
type ServiceRequestInput struct {
	Active         bool                   `json:"active"`
	RequestType    string                 `json:"requestType"`
	Status         string                 `json:"status"`
	Request        string                 `json:"request"`
	ClientID       string                 `json:"clientID"`
	StaffID        string                 `json:"staffID"`
	InProgressBy   *string                `json:"inProgressBy"`
	ResolvedBy     *string                `json:"resolvedBy"`
	FacilityID     string                 `json:"facilityID"`
	ClientName     *string                `json:"clientName"`
	StaffName      string                 `json:"staffName"`
	Flavour        feedlib.Flavour        `json:"flavour"`
	Meta           map[string]interface{} `json:"meta"`
	ProgramID      string                 `json:"programID"`
	OrganisationID string                 `json:"organisationID"`
	CaregiverID    *string                `json:"caregiverID"`
}

// AllergyPayload contains allergy details for a client/patient
type AllergyPayload struct {
	Name              string    `json:"allergy"`
	AllergyConceptID  *string   `json:"allergyConceptId"`
	Reaction          string    `json:"reaction"`
	ReactionConceptID *string   `json:"reactionConceptId"`
	Severity          string    `json:"severity"`
	SeverityConceptID *string   `json:"severityConceptId"`
	Date              time.Time `json:"allergyDateTime"`
}

// VitalSignPayload contains vital signs collected for a particular client/patient
type VitalSignPayload struct {
	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Value     string    `json:"value"`
	Date      time.Time `json:"obsDateTime"`
}

// TestOrderPayload contains details of an orderered test and the date
type TestOrderPayload struct {
	Name      string    `json:"orderedTestName"`
	ConceptID *string   `json:"conceptId"`
	Date      time.Time `json:"orderDateTime"`
}

// TestResultPayload contains results for a completed test
type TestResultPayload struct {
	Name            string    `json:"test"`
	TestConceptID   *string   `json:"testConceptId"`
	Date            time.Time `json:"testDateTime"`
	Result          string    `json:"result"`
	ResultConceptID *string   `json:"resultConceptId"`
}

// MedicationPayload contains details for medication that a patient/client is prescribed or using
type MedicationPayload struct {
	Name                string    `json:"medication"`
	MedicationConceptID *string   `json:"medicationConceptId"`
	Date                time.Time `json:"medicationDateTime"`
	Value               string    `json:"value"`
	DrugConceptID       *string   `json:"drugConceptId"`
}

// PatientRecordPayload contains all the available records for a patient that is available from KenyaEMR
// for syncing an updating on myCareHub
type PatientRecordPayload struct {
	CCCNumber   string               `json:"ccc_number"`
	MFLCode     int                  `json:"MFLCODE"`
	Allergies   []*AllergyPayload    `json:"allergies"`
	VitalSigns  []*VitalSignPayload  `json:"vitalSigns"`
	TestOrders  []*TestOrderPayload  `json:"testOrders"`
	TestResults []*TestResultPayload `json:"testResults"`
	Medications []*MedicationPayload `json:"medications"`
}

// PatientsRecordsPayload is the payload sent from a Kenya EMR instance containing records
// of all newly created/updated patients/clients since the last sync
type PatientsRecordsPayload struct {
	MFLCode string                 `json:"MFLCODE"`
	Records []PatientRecordPayload `json:"records"`
}

// AppointmentServiceRequestInput models the payload that is passed when
// fetching appointment service requests
type AppointmentServiceRequestInput struct {
	MFLCode      int        `json:"MFLCODE"`
	LastSyncTime *time.Time `json:"lastSyncTime"`
}

// SurveyLinkInput is the payload for creating a survey public access link
type SurveyLinkInput struct {
	ProjectID   int    `json:"projectID"`
	FormID      string `json:"formID"`
	DisplayName string `json:"displayName"`
	OnceOnly    bool   `json:"onceOnly"`
}

// ClientFilterParamsInput is the payload for filtering clients
type ClientFilterParamsInput struct {
	ClientTypes []enums.ClientType `json:"clientTypes"`
	AgeRange    *AgeRangeInput     `json:"ageRange"`
	Gender      []enumutils.Gender `json:"gender"`
}

// UserSurveyInput represents a user survey input data structure
type UserSurveyInput struct {
	UserID         string `json:"userID"`
	FormID         string `json:"formID"`
	ProjectID      int    `json:"projectID"`
	LinkID         int    `json:"linkID"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Link           string `json:"link"`
	Token          string `json:"token"`
	ProgramID      string `json:"programID"`
	OrganisationID string `json:"organisationID"`
}

// VerifySurveySubmissionInput represents the payload that is to be sent when a user has filled a survey.
type VerifySurveySubmissionInput struct {
	ProjectID   int     `json:"projectID"`
	FormID      string  `json:"formID"`
	SubmitterID int     `json:"submitterID"`
	CaregiverID *string `json:"caregiverID"`
}

// QuestionnaireInput represents the payload that is to be used when creating a questionnaire.
type QuestionnaireInput struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Questions   []*QuestionInput `json:"questions"`
}

// Validate helps with validation of the QuestionnaireInput
func (q QuestionnaireInput) Validate() error {
	v := validator.New()
	err := v.Struct(q)
	for _, question := range q.Questions {
		if err := question.Validate(); err != nil {
			return err
		}
		switch question.QuestionType {
		case enums.QuestionTypeOpenEnded:
			if question.SelectMultiple {
				return fmt.Errorf("select multiple is not supported for open ended questions")
			}
			if len(question.Choices) > 0 {
				return fmt.Errorf("choices are not supported for open ended questions")
			}
		case enums.QuestionTypeCloseEnded:
			if len(question.Choices) < 2 {
				return fmt.Errorf("at least two choices are required for close ended questions")
			}
		}
	}
	return err
}

// ScreeningToolInput represents the payload that is to be used when creating a questionnaire
type ScreeningToolInput struct {
	Questionnaire QuestionnaireInput `json:"questionnaire"`
	Threshold     int                `json:"threshold"`
	ClientTypes   []enums.ClientType `json:"clientTypes"`
	Genders       []enumutils.Gender `json:"genders"`
	AgeRange      AgeRangeInput      `json:"ageRange"`
	ProgramID     string             `json:"programID"`
}

// QuestionInput represents the input for a Question for a given screening tool in a questionnaire
type QuestionInput struct {
	Text              string                          `json:"text" validate:"required"`
	QuestionType      enums.QuestionType              `json:"questionType" validate:"required"`
	ResponseValueType enums.QuestionResponseValueType `json:"responseValueType" validate:"required"`
	Required          bool                            `json:"required" validate:"required"`
	SelectMultiple    bool                            `json:"selectMultiple"`
	Sequence          int                             `json:"sequence" validate:"required"`
	Choices           []QuestionInputChoiceInput      `json:"choices"`
	ProgramID         string                          `json:"programID"`
}

// Validate helps with validation of a question input
func (s QuestionInput) Validate() error {
	v := validator.New()
	err := v.Struct(s)

	// validate response value type against the the choice provided
	for _, c := range s.Choices {
		switch s.ResponseValueType {
		case enums.QuestionResponseValueTypeNumber:
			_, err := strconv.Atoi(c.Value)
			if err != nil {
				return fmt.Errorf("choice value must be a number")
			}
		case enums.QuestionResponseValueTypeBoolean:
			if _, err := strconv.ParseBool(c.Value); err != nil {
				return fmt.Errorf("response value must be a boolean")
			}
		}
	}
	return err
}

// QuestionInputChoiceInput represents choices for a given question
type QuestionInputChoiceInput struct {
	Choice    *string `json:"choice" validate:"required"`
	Value     string  `json:"value" validate:"required"`
	Score     int     `json:"score"`
	ProgramID string  `json:"programID"`
}

// QuestionnaireScreeningToolResponseInput represents the payload that is to be used when creating a questionnaire screening tool response.
type QuestionnaireScreeningToolResponseInput struct {
	ScreeningToolID   string                                             `json:"screeningToolID" validate:"required"`
	ClientID          string                                             `json:"clientID" validate:"required"`
	QuestionResponses []*QuestionnaireScreeningToolQuestionResponseInput `json:"questionResponses" validate:"required"`
	ProgramID         string                                             `json:"programID"`
	CaregiverID       *string                                            `json:"caregiverID"`
}

// Validate helps with validation of a QuestionnaireScreeningToolResponseInput
func (s QuestionnaireScreeningToolResponseInput) Validate() error {
	v := validator.New()
	err := v.Struct(s)
	for _, qr := range s.QuestionResponses {
		if err := qr.Validate(); err != nil {
			return err
		}
	}
	return err
}

// QuestionnaireScreeningToolQuestionResponseInput represents the payload that is to be used when creating a questionnaire screening tool question response.
type QuestionnaireScreeningToolQuestionResponseInput struct {
	QuestionID string `json:"questionID"  validate:"required"`
	Response   string `json:"response"`
	ProgramID  string `json:"programID"`
}

// Validate helps with validation of a question response input
func (s QuestionnaireScreeningToolQuestionResponseInput) Validate() error {
	v := validator.New()
	err := v.Struct(s)
	return err
}

// SurveyResponseInput is the input for getting a survey response
type SurveyResponseInput struct {
	ProjectID   int    `json:" projectID"`
	FormID      string `json:"formID"`
	SubmitterID int    `json:"submitterID"`
	ProgramID   string `json:"programID"`
}

// StaffFacilityInput is the input for getting a staff facility from the through table
type StaffFacilityInput struct {
	StaffID    *string `json:"staffID"`
	FacilityID *string `json:"facilityID"`
	ProgramID  string  `json:"programID"`
}

// ClientFacilityInput is the input for getting a client facility from the through table
type ClientFacilityInput struct {
	ClientID   *string `json:"clientID"`
	FacilityID *string `json:"facilityID"`
	ProgramID  string  `json:"programID"`
}

// ClientCaregiverInput is the input for used to assign a caregiver to a client
type ClientCaregiverInput struct {
	ClientID      string              `json:"clientID"`
	CaregiverID   string              `json:"caregiverID"`
	CaregiverType enums.CaregiverType `json:"caregiverType"`
	Consent       enums.ConsentState  `json:"consentState"`
}

// OrganisationInput is the input for creating an organisation
type OrganisationInput struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	EmailAddress    string `json:"emailAddress"`
	PhoneNumber     string `json:"phoneNumber"`
	PostalAddress   string `json:"postalAddress"`
	PhysicalAddress string `json:"physicalAddress"`
	DefaultCountry  string `json:"defaultCountry"`
}

// ProgramInput defines the program input structure
type ProgramInput struct {
	Name           string   `json:"name"  validate:"required"`
	Description    string   `json:"description"  validate:"required"`
	OrganisationID string   `json:"organisationID" validate:"required"`
	Facilities     []string `json:"facilities"`
}

// Validate helps with validation of a question response input
func (s ProgramInput) Validate() error {
	v := validator.New()
	err := v.Struct(s)
	return err
}

// OauthClientInput is the input for creating an oauth client
type OauthClientInput struct {
	Name          string   `json:"name"`
	Secret        string   `json:"secret"`
	RedirectURIs  []string `json:"redirectURIs"`
	ResponseTypes []string `json:"responseTypes"`
	Grants        []string `json:"grants"`
}

// MatrixNotifyInput is the input for receiving Matrix's notification event data
type MatrixNotifyInput struct {
	Notification Notification `json:"notification,omitempty"`
}

// Notification contains the notification data
type Notification struct {
	Content           EventContent `json:"content,omitempty"`
	Counts            Counts       `json:"counts,omitempty"`
	Devices           []Devices    `json:"devices,omitempty"`
	EventID           string       `json:"event_id,omitempty"`
	Prio              string       `json:"prio,omitempty"`
	RoomAlias         string       `json:"room_alias,omitempty"`
	RoomID            string       `json:"room_id,omitempty"`
	RoomName          string       `json:"room_name,omitempty"`
	Sender            string       `json:"sender,omitempty"`
	SenderDisplayName string       `json:"sender_display_name,omitempty"`
	Type              string       `json:"type,omitempty"`
}

// EventContent is the events content
type EventContent struct {
	Body    string `json:"body,omitempty"`
	Msgtype string `json:"msgtype,omitempty"`
}

// Counts  dictionary of the current number of unacknowledged communications for the recipient user. Counts whose value is zero should be omitted.
type Counts struct {
	MissedCalls int `json:"missed_calls,omitempty"`
	Unread      int `json:"unread,omitempty"`
}

// Data is a dictionary of additional pusher-specific data. For ‘http’ pushers, this is the data dictionary passed in at pusher creation minus the url key.
type Data struct {
	URL    string `json:"url,omitempty"`
	Format string `json:"format,omitempty"`
}

// Tweaks are a dictionary of customizations made to the way this notification is to be presented.
type Tweaks struct {
	Sound string `json:"sound,omitempty"`
}

// Devices is the device to which the notification should be sent to.
type Devices struct {
	AppID            string `json:"app_id,omitempty"`
	Data             Data   `json:"data,omitempty"`
	Pushkey          string `json:"pushkey,omitempty"`
	PushkeyTimeStamp int    `json:"pushkey_ts,omitempty"`
	Tweaks           Tweaks `json:"tweaks,omitempty"`
}

// CoordinatesInput is used to get the coordinates of a facility
type CoordinatesInput struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
