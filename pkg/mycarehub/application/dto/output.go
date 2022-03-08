package dto

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
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

// HealthDiaryEntriesResponse is the response returned after querying the
// health diary entries for a specific facility
type HealthDiaryEntriesResponse struct {
	MFLCode            int                              `json:"MFLCODE"`
	HealthDiaryEntries []*domain.ClientHealthDiaryEntry `json:"healthDiaries"`
}
