package dto

import (
	"time"

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
	CCCNumber       string                  `json:"ccc_number"`
	AppointmentUUID string                  `json:"appointment_uuid"`
	AppointmentType string                  `json:"appointment_type"`
	Status          enums.AppointmentStatus `json:"status"`
	AppointmentDate scalarutils.Date        `json:"appointment_date"`
	TimeSlot        string                  `json:"time_slot"`
}

// HealthDiaryEntriesResponse is the response returned after querying the
// health diary entries for a specific facility
type HealthDiaryEntriesResponse struct {
	MFLCode            int                              `json:"MFLCODE"`
	HealthDiaryEntries []*domain.ClientHealthDiaryEntry `json:"healthDiaries"`
}

// PatientSyncResponse is the response to a patient sync poll
// the patients payload contains ccc numbers
type PatientSyncResponse struct {
	MFLCode int `json:"MFLCODE"`
	// Patients is a slice of CCC numbers
	Patients []string `json:"patients"`
}

// StaffRegistrationOutput models the staff registration api response
type StaffRegistrationOutput struct {
	ID              string `json:"id"`
	Active          bool   `json:"active"`
	StaffNumber     string `json:"staff_number"`
	UserID          string `json:"user"`
	DefaultFacility string `json:"default_facility"`
}
