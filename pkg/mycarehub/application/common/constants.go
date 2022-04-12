package common

const (
	// FacilityTableName holds the table name for the facilities. This is a table that stores
	// the facilities data
	FacilityTableName = "common_facility"

	// OrganizationID is the default organization ID that's added to all models on the django side
	OrganizationID = "DEFAULT_ORG_ID"

	// CreatePatientTopic is the topic ID for patient creation in clinical service
	CreatePatientTopic = "patient.create"

	// CreateVitalsTopicName is the topic for publishing a patient's vital signs
	CreateVitalsTopicName = "vitals.create"

	// CreateAllergyTopicName is the topic for publishing a patient's allergy
	CreateAllergyTopicName = "allergy.create"

	// CreateMedicationTopicName is the topic for publishing a patient's medication
	CreateMedicationTopicName = "medication.create"

	// CreateTestResultTopicName is the topic for publishing a patient's test results
	CreateTestResultTopicName = "test.result.create"

	// CreateTestOrderTopicName is the topic for publishing a patient's test order
	CreateTestOrderTopicName = "test.order.create"

	// CreateOrganizationTopicName is the topic for publishing a organization
	CreateOrganizationTopicName = "organization.create"
)
