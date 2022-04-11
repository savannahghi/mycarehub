package dto

import (
	"time"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

// RestEndpointResponses represents the rest endpoints response(s) output
type RestEndpointResponses struct {
	Data map[string]interface{} `json:"data"`
}

// ClientRegistrationOutput models the client registration API response
type ClientRegistrationOutput struct {
	ID                string           `json:"id"`
	Active            bool             `json:"active"`
	ClientType        enums.ClientType `json:"client_type"`
	EnrollmentDate    *time.Time       `json:"enrollment_date"`
	FHIRPatientID     string           `json:"fhir_patient_id"`
	EMRHealthRecordID string           `json:"emr_health_record_id"`
	TreatmentBuddy    string           `json:"treatment_buddy"`
	Counselled        bool             `json:"counselled"`
	Organisation      string           `json:"organisation"`
	UserID            string           `json:"user"`
	CurrentFacilityID string           `json:"current_facility"`
	CHV               string           `json:"chv"`
	Caregiver         string           `json:"caregiver"`
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

	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Value     string    `json:"value"`
	Date      time.Time `json:"date"`
}

// PatientTestOrderOutput contains details of an orderered test and the date
type PatientTestOrderOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`

	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Date      time.Time `json:"date"`
}

// PatientTestResultOutput contains results for a completed test
type PatientTestResultOutput struct {
	PatientID      string `json:"patientID"`
	OrganizationID string `json:"organizationID"`

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
	ID              string `json:"id"`
	Active          bool   `json:"active"`
	StaffNumber     string `json:"staff_number"`
	UserID          string `json:"user"`
	DefaultFacility string `json:"default_facility"`
}

// AppointmentServiceRequestsOutput is the response returned after querying the
// appointment service requests for a specific facility
type AppointmentServiceRequestsOutput struct {
	AppointmentServiceRequests []domain.AppointmentServiceRequests `json:"Results"`
}

// PatientCreationOutput is the payload sent to the clinical service for creation of a patient
type PatientCreationOutput struct {
	ID     string `json:"id"`
	UserID string `json:"user"`
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
