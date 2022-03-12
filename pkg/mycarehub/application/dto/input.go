package dto

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
	"gopkg.in/go-playground/validator.v9"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Code        int    `json:"code" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Active      bool   `json:"active"`
	County      string `json:"county" validate:"required"`
	Description string `json:"description" validate:"required,min=3,max=256"`
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
	UserID           string
	Message          string
	RequiresFollowUp bool
}

// FeedbackEmail defines the field to be parsed when sending feedback emails
type FeedbackEmail struct {
	User             string
	Message          string
	RequiresFollowUp string
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
	Facility       string           `json:"facility"`
	ClientType     enums.ClientType `json:"client_type"`
	ClientName     string           `json:"name"`
	Gender         enumutils.Gender `json:"gender"`
	DateOfBirth    scalarutils.Date `json:"date_of_birth"`
	PhoneNumber    string           `json:"phone_number"`
	EnrollmentDate scalarutils.Date `json:"enrollment_date"`
	CCCNumber      string           `json:"ccc_number"`
	Counselled     bool             `json:"counselled"`
	InviteClient   bool             `json:"inviteClient"`
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
	MFLCode            int              `json:"MFLCODE"`
	CCCNumber          string           `json:"cccNumber"`
	Name               string           `json:"name"`
	DateOfBirth        scalarutils.Date `json:"dateOfBirth"`
	ClientType         string           `json:"clientType"`
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

// ScreeningToolQuestionResponseInput defines the field passed when answering screening tools questions
type ScreeningToolQuestionResponseInput struct {
	ClientID   string `json:"clientID" validate:"required"`
	QuestionID string `json:"questionID" validate:"required"`
	Response   string `json:"response" validate:"required"`
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
	Facility    string           `json:"facility"`
	StaffName   string           `json:"name"`
	Gender      enumutils.Gender `json:"gender"`
	DateOfBirth scalarutils.Date `json:"date_of_birth"`
	PhoneNumber string           `json:"phone_number"`
	IDNumber    int              `json:"id_number"`
	StaffNumber string           `json:"staff_number"`
	StaffRoles  string           `json:"staff_roles"`
	InviteStaff bool             `json:"invite_staff"`
}
