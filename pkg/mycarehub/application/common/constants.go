package common

const (
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

	// CreateOrganizationTopicName is the topic for publishing an organization
	CreateOrganizationTopicName = "organization.create"

	// CreateGetstreamEventTopicName is the topic where getstream events are published to
	CreateGetstreamEventTopicName = "getstream.event"

	// CreateCMSClientTopicName is the topic where cms user who is to be created, is published to.
	CreateCMSClientTopicName = "cms.client.create"

	// DeleteCMSClientTopicName is the topic where cms user who is to be deleted, is published to.
	DeleteCMSClientTopicName = "cms.client.delete"

	// DeleteCMSStaffTopicName is the topic where cms staff who is to be deleted, is published to.
	DeleteCMSStaffTopicName = "cms.staff.delete"

	// CreateCMSStaffTopicName is the topic where cms user(staff in this case) is to be created, is published to.
	CreateCMSStaffTopicName = "cms.staff.create"
)
