package dto

import (
	"encoding/json"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
	validator "gopkg.in/go-playground/validator.v9"
)

// RestEndpointResponses represents the rest endpoints response(s) output
type RestEndpointResponses struct {
	Data map[string]interface{} `json:"data"`
}

// ClientRegistrationOutput models the client registration API response
type ClientRegistrationOutput struct {
	ID                string             `json:"id"`
	Active            bool               `json:"active"`
	ClientTypes       []enums.ClientType `json:"client_types"`
	EnrollmentDate    *time.Time         `json:"enrollment_date"`
	FHIRPatientID     string             `json:"fhir_patient_id"`
	EMRHealthRecordID string             `json:"emr_health_record_id"`
	TreatmentBuddy    string             `json:"treatment_buddy"`
	Counselled        bool               `json:"counselled"`
	Organisation      string             `json:"organisation"`
	UserID            string             `json:"user"`
	CurrentFacilityID string             `json:"current_facility"`
	CHV               string             `json:"chv"`
	Caregiver         string             `json:"caregiver"`
}

// FacilityAppointmentsResponse is the response sent after creating/updating an appointment
type FacilityAppointmentsResponse struct {
	MFLCode      string                `json:"MFLCODE"`
	Appointments []AppointmentResponse `json:"appointments"`
}

// AppointmentResponse is the response representing an appointment
type AppointmentResponse struct {
	CCCNumber         string           `json:"ccc_number"`
	ExternalID        string           `json:"appointment_id"`
	AppointmentDate   scalarutils.Date `json:"appointment_date"`
	AppointmentReason string           `json:"appointment_reason"`
}

// HealthDiaryEntriesResponse is the response returned after querying the
// health diary entries for a specific facility
type HealthDiaryEntriesResponse struct {
	MFLCode            int                              `json:"MFLCODE"`
	HealthDiaryEntries []*domain.ClientHealthDiaryEntry `json:"healthDiaries"`
}

// RedFlagServiceRequestResponse models the response returned when fetching for
// red flags service requests
type RedFlagServiceRequestResponse struct {
	RedFlagServiceRequests []*domain.ServiceRequest `json:"redFlags"`
}

// PatientSyncResponse is the response to a patient sync poll
// the patients payload contains ccc numbers
type PatientSyncResponse struct {
	MFLCode int `json:"MFLCODE"`
	// Patients is a slice of CCC numbers
	Patients []string `json:"patients"`
}

// PatientAllergyOutput contains allergy details for a client/patient
type PatientAllergyOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`

	Name      string          `json:"name"`
	ConceptID *string         `json:"conceptID"`
	Date      time.Time       `json:"date"`
	Reaction  AllergyReaction `json:"reaction"`
	Severity  AllergySeverity `json:"severity"`
}

// AllergyReaction ...
type AllergyReaction struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// AllergySeverity ...
type AllergySeverity struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// PatientVitalSignOutput contains vital signs collected for a particular client/patient
type PatientVitalSignOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`

	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Value     string    `json:"value"`
	Date      time.Time `json:"date"`
}

// PatientTestOrderOutput contains details of an orderered test and the date
type PatientTestOrderOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`

	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Date      time.Time `json:"date"`
}

// PatientTestResultOutput contains results for a completed test
type PatientTestResultOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`

	Name      string     `json:"name"`
	ConceptID *string    `json:"conceptId"`
	Date      time.Time  `json:"date"`
	Result    TestResult `json:"result"`
}

// TestResult ...
type TestResult struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptId"`
}

// PatientMedicationOutput contains details for medication that a patient/client is prescribed or using
type PatientMedicationOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`

	Name      string          `json:"medication"`
	ConceptID *string         `json:"conceptId"`
	Date      time.Time       `json:"date"`
	Value     string          `json:"value"`
	Drug      *MedicationDrug `json:"drug"`
}

// MedicationDrug ...
type MedicationDrug struct {
	ConceptID *string `json:"conceptId"`
}

// StaffRegistrationOutput models the staff registration api response
type StaffRegistrationOutput struct {
	ID                  string       `json:"id"`
	Active              bool         `json:"active"`
	StaffNumber         string       `json:"staff_number"`
	UserID              string       `json:"user"`
	DefaultFacility     string       `json:"default_facility"`
	UserProfile         *domain.User `json:"user_profile"`
	IsOrganisationAdmin bool         `json:"isOrganisationAdmin"`
}

// AppointmentServiceRequestsOutput is the response returned after querying the
// appointment service requests for a specific facility
type AppointmentServiceRequestsOutput struct {
	AppointmentServiceRequests []domain.AppointmentServiceRequests `json:"Results"`
}

// PatientCreationOutput is the payload sent to the clinical service for creation of a patient
type PatientCreationOutput struct {
	UserID      string           `json:"userID"`
	ClientID    string           `json:"clientID"`
	Name        string           `json:"name"`
	DateOfBirth *time.Time       `json:"dateOfBirth"`
	Gender      enumutils.Gender `json:"gender"`
	Active      bool             `json:"active"`
	PhoneNumber string           `json:"phoneNumber"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `json:"facilityID"`
}

// FCMNotificationMessage models the payload passed when composing a notification payload
//
// The title is what will appear as the notification's title message on the phone's notification tray
// Most of the notifications will be `BLIND` meaning that the body will be empty
type FCMNotificationMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SavedNotification is used to serialize and save successful FCM notifications.
//
// It's the basis for a primitive "inbox" - a mechanism by which an app can
// request it's messages in bulk.
type SavedNotification struct {
	ID                string                                         `json:"id,omitempty"`
	RegistrationToken string                                         `json:"registrationToken,omitempty"`
	MessageID         string                                         `json:"messageID,omitempty"`
	Timestamp         time.Time                                      `json:"timestamp,omitempty"`
	Data              map[string]interface{}                         `json:"data,omitempty"`
	Notification      *firebasetools.FirebaseSimpleNotificationInput `json:"notification,omitempty"`
	AndroidConfig     *firebasetools.FirebaseAndroidConfigInput      `json:"androidConfig,omitempty"`
	WebpushConfig     *firebasetools.FirebaseWebpushConfigInput      `json:"webpushConfig,omitempty"`
	APNSConfig        *firebasetools.FirebaseAPNSConfigInput         `json:"apnsConfig,omitempty"`
}

// SurveyForm is contains the information about a survey form
type SurveyForm struct {
	ProjectID    int       `json:"projectId"`
	XMLFormID    string    `json:"xmlFormId"`
	State        string    `json:"state"`
	EnketoID     string    `json:"enketoId"`
	EnketoOnceID string    `json:"enketoOnceId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	KeyID        string    `json:"keyId"`
	Version      string    `json:"version"`
	Hash         string    `json:"hash"`
	SHA          string    `json:"sha"`
	SHA256       string    `json:"sha256"`
	DraftToken   string    `json:"draftToken"`
	PublishedAt  time.Time `json:"publishedAt"`
	Name         string    `json:"name"`
}

// SurveyPublicLink is contains the information about a survey public link
type SurveyPublicLink struct {
	ID          int        `json:"id"`
	DisplayName string     `json:"displayName"`
	Once        bool       `json:"once"`
	Token       string     `json:"token"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	CSRF        string     `json:"csrf"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

// SurveysWithServiceRequest is the domain model used to list surveys with a service request
type SurveysWithServiceRequest struct {
	Title     string `json:"name"`
	ProjectID int    `json:"projectID"`
	LinkID    int    `json:"linkID"`
	FormID    string `json:"formID"`
}

// PubsubCreateCMSClientPayload contains the user's data model to be used to publish the user who is to be created in the content CMS service.
type PubsubCreateCMSClientPayload struct {
	ClientID       string           `json:"client_id"`
	Name           string           `json:"name"`
	Gender         string           `json:"gender"`
	DateOfBirth    scalarutils.Date `json:"date_of_birth"`
	OrganisationID string           `json:"organisation_id"`
	ProgramID      string           `json:"program_id"`
}

// PubsubCreateCMSStaffPayload is the payload passed when creating a staff user on the CMS service
type PubsubCreateCMSStaffPayload struct {
	UserID         string           `json:"user_id"`
	Name           string           `json:"name"`
	Gender         enumutils.Gender `json:"gender"`
	UserType       enums.UsersType  `json:"user_type"`
	PhoneNumber    string           `json:"phone_number"`
	Handle         string           `json:"handle"`
	Flavour        feedlib.Flavour  `json:"flavour"`
	DateOfBirth    scalarutils.Date `json:"date_of_birth"`
	StaffNumber    string           `json:"staff_number"`
	StaffID        string           `json:"staff_id"`
	FacilityID     string           `json:"facility_id"`
	FacilityName   string           `json:"facility_name"`
	OrganisationID string           `json:"organisation_id"`
}

// FacilityOutputPage returns a paginated list of users facility
type FacilityOutputPage struct {
	Pagination *domain.Pagination
	Facilities []*domain.Facility
}

// ManagedClientOutputPage returns a paginated list of managed client profiles
type ManagedClientOutputPage struct {
	Pagination     *domain.Pagination      `json:"pagination"`
	ManagedClients []*domain.ManagedClient `json:"managedClients"`
}

// CaregiverProfileOutputPage returns a paginated list of users caregiver profile
type CaregiverProfileOutputPage struct {
	Pagination *domain.Pagination         `json:"pagination"`
	Caregivers []*domain.CaregiverProfile `json:"caregivers"`
}

// Organisation represents output for a tenant/organisation
type Organisation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// OrganisationsOutput is the output for fetching organisations associated with a contact
type OrganisationsOutput struct {
	Count         int            `json:"count"`
	Organisations []Organisation `json:"organisations"`
}

// ProgramOutput represents the output for fetching programs
type ProgramOutput struct {
	Count    int               `json:"count"`
	Programs []*domain.Program `json:"programs"`
}

// CreateCMSProgramPayload is the payload passed when creating a program on the CMS service using PubSub
type CreateCMSProgramPayload struct {
	ProgramID      string `json:"program_id"`
	Name           string `json:"name"`
	OrganisationID string `json:"organisation_id"`
}

// CreateCMSOrganisationPayload is the payload passed when creating an organisation on the CMS service using PubSub
type CreateCMSOrganisationPayload struct {
	OrganisationID string `json:"organisation_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Code           int    `json:"code"`
}

// CreateCMSFacilityPayload is the payload passed when creating a facility on the CMS service using PubSub
type CreateCMSFacilityPayload struct {
	FacilityID string `json:"id"`
	Name       string `json:"name"`
}

// CMSLinkFacilityToProgramPayload is the payload passed when adding a facility to a program on the CMS service using PubSub
type CMSLinkFacilityToProgramPayload struct {
	FacilityID []string `json:"facilities"`
	ProgramID  string   `json:"program_id,omitempty"`
}

// ClinicalTenantPayload is the dataclass to create a clinical tenant
type ClinicalTenantPayload struct {
	Name        string                     `json:"name,omitempty"`
	PhoneNumber string                     `json:"phoneNumber,omitempty"`
	Identifiers []ClinicalTenantIdentifier `json:"identifiers,omitempty"`
}

// ClinicalTenantIdentifier models the clinical's tenant identification model
type ClinicalTenantIdentifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type OrganisationOutputPage struct {
	Pagination    *domain.Pagination     `json:"pagination"`
	Organisations []*domain.Organisation `json:"organisations"`
}

// Output expresses a type constraint satisfied by the output structs
type Output interface {
	OrganisationOutput | ProgramJsonOutput
}

// ParseValues is a generic function that takes in any concrete type, parses the values
// with the provided slice of bytes and validates the parsed generic is not empty
func ParseValues[T Output](concreteType T, values []byte) (*T, error) {
	output := new(T)

	err := json.Unmarshal(values, &output)
	if err != nil {
		return nil, err
	}

	v := validator.New()
	err = v.Struct(output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// OrganisationOutput is a struct that stores the output of organisation json values
type OrganisationOutput struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	EmailAddress    string `json:"emailAddress" validate:"required"`
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	PostalAddress   string `json:"postalAddress" validate:"required"`
	PhysicalAddress string `json:"physicalAddress" validate:"required"`
	DefaultCountry  string `json:"defaultCountry" validate:"required"`
}

// ParseValues transforms and validates the json organisation to type OrganisationInput
func (o *OrganisationOutput) ParseValues(values []byte) (*OrganisationInput, error) {
	organisation, err := ParseValues(*o, values)
	if err != nil {
		return nil, err
	}

	return &OrganisationInput{
		Name:            organisation.Name,
		Description:     organisation.Description,
		EmailAddress:    organisation.EmailAddress,
		PhoneNumber:     organisation.PhoneNumber,
		PostalAddress:   organisation.PostalAddress,
		PhysicalAddress: organisation.PhysicalAddress,
		DefaultCountry:  organisation.DefaultCountry,
	}, nil
}

// ProgramJsonOutput is a struct that stores the output of program json values
type ProgramJsonOutput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// ParseValues transforms and validates the json program to type ProgramInput
func (p *ProgramJsonOutput) ParseValues(values []byte) (*ProgramInput, error) {
	program, err := ParseValues(*p, values)
	if err != nil {
		return nil, err
	}

	return &ProgramInput{
		Name:        program.Name,
		Description: program.Description,
	}, nil
}

// MatrixUserRegistrationOutput is used to show the response after a user has been registered in Matrix
type MatrixUserRegistrationOutput struct {
	Name        string `json:"name"`
	Admin       bool   `json:"admin"`
	DisplayName string `json:"displayname"`
}

// MatrixUserRegistrationPayload is the payload passed when registering a Matrix user via pubsub
type MatrixUserRegistrationPayload struct {
	Auth             *domain.MatrixAuth             `json:"auth"`
	RegistrationData *domain.MatrixUserRegistration `json:"registrationData"`
}

// MatrixNotifyOutput returns a list of the rejected push keys when a notification is sent from Matrix
type MatrixNotifyOutput struct {
	Rejected []string `json:"rejected,omitempty"`
}

// UpdatePatientFHIRID represents the data structure used for updating the FHIR ID in a user's profile.
type UpdatePatientFHIRID struct {
	FhirID   string `json:"fhirID"`
	ClientID string `json:"clientID"`
}

// UpdateProgramFHIRID represents the data structure used for updating the fhir id field in a program
type UpdateProgramFHIRID struct {
	ProgramID    string `json:"programID"`
	FHIRTenantID string `json:"fhirTenantID"`
}

// UpdateFacilityFHIRID represents the data structure used for updating the fhir id field in a facility
type UpdateFacilityFHIRID struct {
	FacilityID string `json:"facilityID"`
	FhirID     string `json:"fhirOrganisationID"`
}
