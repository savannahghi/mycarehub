package dto

import (
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
	validator "gopkg.in/go-playground/validator.v9"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name               string `json:"name" validate:"required,min=3,max=100"`
	Code               int    `json:"code" validate:"required"`
	Phone              string `json:"phone" validate:"required"`
	Active             bool   `json:"active"`
	County             string `json:"county" validate:"required"`
	Description        string `json:"description" validate:"required,min=3,max=256"`
	FHIROrganisationID string `json:"fhirOrganisationId"`
}

// Validate helps with validation of facility input fields
func (f *FacilityInput) Validate() error {
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
	PhoneNumber *string         `json:"phoneNumber" validate:"required"`
	PIN         *string         `json:"pin" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// Validate helps with validation of LoginInput fields
func (f *LoginInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// PINInput represents the Pin input data structure
type PINInput struct {
	UserID     *string         `json:"id" validate:"required"`
	PIN        *string         `json:"pin" validate:"required"`
	ConfirmPIN *string         `json:"confirmPin" validate:"required"`
	Flavour    feedlib.Flavour `json:"flavour" validate:"required"`
}

// VerifyOTPInput represents the verify OTP input data structure
type VerifyOTPInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// SendOTPInput represents the send OTP input data structure
type SendOTPInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// PhoneInput is used to define the inputs needed carrying out an activity that requires a phone number and flavour.
type PhoneInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
}

// SendRetryOTPPayload is used to define the inputs passed when calling the endpoint
// that resends an otp
type SendRetryOTPPayload struct {
	Phone   string          `json:"phoneNumber" validate:"required"`
	Flavour feedlib.Flavour `json:"flavour" validate:"required"`
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
	Response           string `json:"Response" validate:"required"`
}

// Validate helps with validation of SecurityQuestionResponseInput fields
func (f *SecurityQuestionResponseInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// VerifySecurityQuestionInput defines the field passed when verifying the set security questions
type VerifySecurityQuestionInput struct {
	QuestionID  string          `json:"questionID" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
	Response    string          `json:"response" validate:"required"`
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
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
	PhoneNumber string          `json:"phoneNumber"`
	Flavour     feedlib.Flavour `json:"flavour"`
}

// GetSecurityQuestionsInput defines the field passed when getting the security questions
type GetSecurityQuestionsInput struct {
	Flavour feedlib.Flavour `json:"flavour" validate:"required"`
}

// GetUserRespondedSecurityQuestionsInput defines the field passed when getting the security questions
type GetUserRespondedSecurityQuestionsInput struct {
	PhoneNumber string          `json:"phonenumber" validate:"required"`
	Flavour     feedlib.Flavour `json:"flavour" validate:"required"`
	OTP         string          `json:"otp" validate:"required"`
}

// Validate helps with validation of GetUserRespondedSecurityQuestionsInput fields
func (f *GetUserRespondedSecurityQuestionsInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// UserResetPinInput contains the fields requires when a user is resetting a pin
type UserResetPinInput struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
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
	UserID    string `json:"userID" validate:"required"`
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
	UserID            string
	FeedbackType      enums.FeedbackType
	SatisfactionLevel int
	ServiceName       string
	Feedback          string
	RequiresFollowUp  bool
	PhoneNumber       string
}

// FeedbackEmail defines the field to be parsed when sending feedback emails
type FeedbackEmail struct {
	User              string
	FeedbackType      enums.FeedbackType
	SatisfactionLevel int
	ServiceName       string
	Feedback          string
	PhoneNumber       string
}

// CaregiverInput defines the field passed when creating a caregiver
type CaregiverInput struct {
	ClientID      string              `json:"clientID"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	PhoneNumber   string              `json:"phoneNumber"`
	CaregiverType enums.CaregiverType `json:"caregiverType"`
}

// Validate helps with validation of CaregiverInput fields
func (f *CaregiverInput) Validate() error {
	v := validator.New()
	err := v.Struct(f)
	return err
}

// ClientRegistrationInput defines the fields passed as a payload to the client registration API
type ClientRegistrationInput struct {
	Facility       string             `json:"facility" validate:"required"`
	ClientTypes    []enums.ClientType `json:"client_types" validate:"required"`
	ClientName     string             `json:"name" validate:"required"`
	Gender         enumutils.Gender   `json:"gender" validate:"required"`
	DateOfBirth    scalarutils.Date   `json:"date_of_birth" validate:"required"`
	PhoneNumber    string             `json:"phone_number" validate:"required"`
	EnrollmentDate scalarutils.Date   `json:"enrollment_date" validate:"required"`
	CCCNumber      string             `json:"ccc_number" validate:"required"`
	Counselled     bool               `json:"counselled" validate:"required"`
	InviteClient   bool               `json:"inviteClient"`
}

// Validate helps with validation of ClientRegistrationInput fields
func (f *ClientRegistrationInput) Validate() error {
	v := validator.New()
	err := v.Struct(f)
	return err
}

// CommunityInput defines the payload to create a channel
type CommunityInput struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	AgeRange    *AgeRangeInput      `json:"ageRange"`
	Gender      []*enumutils.Gender `json:"gender"`
	ClientType  []*enums.ClientType `json:"clientType"`
	InviteOnly  bool                `json:"inviteOnly"`
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

//ServiceRequestPayload defines the payload from KenyaEMR used to fetch
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

// ScreeningToolQuestionResponseInput defines the field passed when answering screening tools questions
type ScreeningToolQuestionResponseInput struct {
	ClientID         string                              `json:"clientID" validate:"required"`
	QuestionID       string                              `json:"questionID" validate:"required"`
	Response         string                              `json:"response" validate:"required"`
	ToolType         enums.ScreeningToolType             `json:"toolType" validate:"required"`
	ResponseType     enums.ScreeningToolResponseType     `json:"responseType" validate:"required"`
	ResponseCategory enums.ScreeningToolResponseCategory `json:"responseCategory" validate:"required"`
	QuestionSequence int                                 `json:"questionSequence" validate:"required"`
}

// Validate helps with validation of ScreeningToolQuestionResponseInput fields
func (f *ScreeningToolQuestionResponseInput) Validate() error {
	v := validator.New()
	err := v.Struct(f)
	return err
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
	Facility    string           `json:"facility" validate:"required"`
	StaffName   string           `json:"name" validate:"required"`
	Gender      enumutils.Gender `json:"gender" validate:"required"`
	DateOfBirth scalarutils.Date `json:"date_of_birth" validate:"required"`
	PhoneNumber string           `json:"phone_number" validate:"required"`
	IDNumber    string           `json:"id_number" validate:"required"`
	StaffNumber string           `json:"staff_number" validate:"required"`
	StaffRoles  string           `json:"role"`
	InviteStaff bool             `json:"invite_staff"`
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

// PinResetServiceRequestPayload models the details passed to an API when a pin reset service request
// is being created
type PinResetServiceRequestPayload struct {
	CCCNumber   string          `json:"cccNumber"`
	PhoneNumber string          `json:"phoneNumber"`
	Flavour     feedlib.Flavour `json:"flavour"`
}

// ServiceRequestInput is a domain entity that represents a service request.
type ServiceRequestInput struct {
	Active       bool                   `json:"active"`
	RequestType  string                 `json:"requestType"`
	Status       string                 `json:"status"`
	Request      string                 `json:"request"`
	ClientID     string                 `json:"clientID"`
	StaffID      string                 `json:"staffID"`
	InProgressBy *string                `json:"inProgressBy"`
	ResolvedBy   *string                `json:"resolvedBy"`
	FacilityID   string                 `json:"facility_id"`
	ClientName   *string                `json:"client_name"`
	StaffName    string                 `json:"staff_name"`
	Flavour      feedlib.Flavour        `json:"flavour"`
	Meta         map[string]interface{} `json:"meta"`
}

// ClientFHIRPayload is the payload from clinical service with patient's fhir ID
type ClientFHIRPayload struct {
	ClientID string `json:"clientID,omitempty"`
	FHIRID   string `json:"fhirID,omitempty"`
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

// UpdateFacilityPayload is the payload for updating faacility(s) fhir organization ID
type UpdateFacilityPayload struct {
	FacilityID         string `json:"facilityID"`
	FHIROrganisationID string `json:"fhirOrganisationID"`
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
	UserID      string `json:"userID"`
	FormID      string `json:"formID"`
	ProjectID   int    `json:"projectID"`
	LinkID      int    `json:"linkID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Token       string `json:"token"`
}

// VerifySurveySubmissionInput represents the payload that is to be sent when a user has filled a survey.
type VerifySurveySubmissionInput struct {
	ProjectID   int
	FormID      string
	SubmitterID int
}
