package common

const (
	// OrganizationID is the default organization ID that's added to all models on the django side
	OrganizationID = "DEFAULT_ORG_ID"

	// ProgramID is the default program ID that's added to all models that require program ID
	ProgramID = "DEFAULT_PROGRAM_ID"

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

	// CreateCMSClientTopicName is the topic where cms user who is to be created, is published to.
	CreateCMSClientTopicName = "cms.client.create"

	// DeleteCMSClientTopicName is the topic where cms user who is to be deleted, is published to.
	DeleteCMSClientTopicName = "cms.client.delete"

	// DeleteCMSStaffTopicName is the topic where cms staff who is to be deleted, is published to.
	DeleteCMSStaffTopicName = "cms.staff.delete"

	// CreateCMSStaffTopicName is the topic where cms user(staff in this case) is to be created, is published to.
	CreateCMSStaffTopicName = "cms.staff.create"

	// CreateCMSProgramTopicName is the topic where the program, which has been created in myCareHub, is published to.
	CreateCMSProgramTopicName = "cms.program.create"

	// CreateCMSOrganisationTopicName is the topic where the organisation, which has been created in myCareHub, is published to.
	CreateCMSOrganisationTopicName = "cms.organisation.create"

	// CreateCMSFacilityTopicName is the topic where the facility, which has been created in myCareHub, is published to.
	CreateCMSFacilityTopicName = "cms.facility.create"

	// CreateCMSProgramFacilityTopicName is the topic where when a facility is added to a program in myCareHub, is published to.
	CreateCMSProgramFacilityTopicName = "cms.programfacility.link"

	// TenantTopicName is the topic where program is registered in clinical service as a tenant
	TenantTopicName = "mycarehub.tenant.create"
)
