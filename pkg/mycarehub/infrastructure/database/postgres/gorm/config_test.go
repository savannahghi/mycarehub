package gorm_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/authorization"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

var (
	fixtures  *testfixtures.Loader
	testingDB *gorm.PGInstance

	orgID                  = os.Getenv("DEFAULT_ORG_ID")
	orgID2                 = "3766b8ca-8cfa-43d5-a334-83507130de1a"
	orgIDToAddToProgram    = "a25a69ef-027d-4f57-8ea5-b2e43d9c1d34"
	organisationIDToDelete = "1c396506-607c-42d1-8abc-425b1e00d029"

	termsID         = 50005
	proTermsID      = 50006
	consumerTermsID = 50007
	db              *sql.DB

	testPhone   = gofakeit.Phone()
	testPhone2  = interserviceclient.TestUserPhoneNumber
	testFlavour = feedlib.FlavourConsumer
	futureTime  = time.Now().Add(time.Hour * 24 * 365 * 10)
	testOTP     = "1234"
	// user variables
	userID                             = "6ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userID2                            = "6ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	userIDtoAddCaregiver               = "8ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	userIDtoAssignClient               = "4181df12-ca96-4f28-b78b-8e8ad88b25df"
	userIDtoAssignStaff                = "6ecccc80-24c8-421a-9f1a-e14e13678ef0"
	userIDToInvalidate                 = "5ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToAcceptTerms                = "4ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDUpdatePinRequireChangeStatus = "5ecbbc80-24b8-421a-9f1a-e14e12678ee0"
	userIDToSavePin                    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	testUserWithCaregiver              = "e1e90ea3-fc06-442e-a1ec-251a031c0ca7"
	testUserWithoutCaregiver           = "723b64b3-e4d6-4416-98b2-18798279e457"
	testUserHasNotGivenConsent         = "839f9a85-bbe6-48e7-a730-42d56a39b532"
	staffUserToAddAsClient             = "f186100a-2b6c-4656-9bbd-960492f6bfb4"
	clientUserToAddAsClient            = "4aa35fa8-a720-4c6f-9510-86fe4b4addbd"

	treatmentBuddyID       = "5ecbbc80-24c8-421a-9f1a-e14e12678ee1"
	treatmentBuddyID2      = "5ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	fhirPatientID          = "5ecbbc80-24c8-421a-9f1a-e14e12678ee2"
	fhirPatientID2         = "f933fd4b-1e3c-4ecd-9d7a-82b2790c0543"
	testEmrHealthRecordID  = "5ecbbc80-24c8-421a-9f1a-e14e12678ee3"
	testEmrHealthRecordID2 = "5ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	testChvID              = "5ecbbc80-24c8-421a-9f1a-e14e12678ee4"
	testChvID2             = "5ecbbc80-24c8-421a-9f1a-e14e12678ef4"

	clientID                        = "26b20a42-cbb8-4553-aedb-c539602d04fc"
	clientID2                       = "00a6a0cd-42ac-417b-97d9-e939a1232de1"
	testClientWithCaregiver         = "f3265be7-54cd-4df9-a078-66bcb31e4dcc"
	testClientWithoutCaregiver      = "13bc475c-6fa8-40a1-ae20-2c9d137ca6e4"
	testClientHasNotGivenConsent    = "5f279d05-0df4-431d-8f70-6f7c76feb425"
	testClientToAddToAnotherProgram = "01bd8f8d-a1f6-45cf-973d-afb9bde23d87"
	clientWithRolesID               = "79b0aae0-1c42-4b2b-8920-12f7c05dddd9"

	contactID    = "bdc22436-e314-43f2-bb39-ba1ab332f9b0"
	identifierID = "bcbdaf68-3d36-4365-b575-4182d6749af5"

	// Facility variables
	facilityID                                = "4181df12-ca96-4f28-b78b-8e8ad88b25df"
	facilityIdentifierID                      = "b432032a-6957-11ed-a1eb-0242ac120002"
	facilityToAddToUserProfile                = "5181df12-ca96-4f28-b78b-8e8ad88b25de"
	facilityIdentifierToAddToUserProfile      = "dac51586-6957-11ed-a1eb-0242ac120002"
	facilityToRemoveFromUserProfile           = "bdc22436-e314-43f2-bb39-ba1ab332f9b0"
	facilityIdentifierToRemoveFromUserProfile = "2ec1f62c-6958-11ed-a1eb-0242ac120002"
	facilityToAddExistingStaff                = "7fb061a6-827e-462f-8a7e-0144643468c4"

	mflIdentifier                  = "324459"
	inactiveFacilityIdentifier     = "229900"
	facilityIdentifierToInactivate = "223900"
	mflIdentifierType              = enums.FacilityIdentifierTypeMFLCode.String()
	// Pin variables
	salt, encryptedPin string
	// Securityquestions variables
	securityQuestionID  = "26b20a42-cbb8-4553-aedb-c539602d04fb"
	securityQuestionID2 = "fada0b8a-4f3c-4df2-82be-35b82753f66c"
	securityQuestionID3 = "bdc22436-e314-43f2-bb39-ba1ab332f9b6"
	securityQuestionID4 = "e7f0e561-40fc-46db-84c2-18c6f26db40e"

	securityQuestionResponseID  = "6da66afc-58d4-11ec-bf63-0242ac130002"
	securityQuestionResponseID2 = "312d63a4-58d5-11ec-bf63-0242ac130002"
	securityQuestionResponseID3 = "f4cf3ffa-8d4e-45fa-ad19-c5cac7701e61"
	securityQuestionResponseID4 = "7225e76b-7780-46a9-a217-8e858789a869"

	// Caregiver
	testCaregiverID      = "26b20a42-cbb8-4553-aedb-c593602d04fc"
	testClientCaregiver1 = "28b20a42-cbb8-4553-aedb-c575602d04fc"
	testCaregiverOrg2ID  = "4e4ef3d2-eb26-407a-82c3-31243dc923cd"
	caregiverWithRolesID = "484831c5-9b63-4580-9aef-4bffb4bdd230"

	//Terms
	termsText = "Test terms"

	// Staff
	staffNumber = "test-Staff-101"
	staffID     = "8ecbbc80-24c8-421a-9f1a-e14e12678ef1"

	clientsServiceRequestID   = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	staffServiceRequestID     = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	clientsHealthDiaryEntryID = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	// Service Request
	serviceRequestID               = "8ecbbc80-24c8-421a-9f1a-e14e12678ef2"
	clientServiceRequestIDToUpdate = "fffbb75c-9138-47e8-a75b-d7ee5df5e9a0"

	// Authority
	canInviteUserPermissionID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	canEditOwnRolePermissionID   = "29672457-d081-48e0-a007-8f49cedb5c6f"
	canManageContentPermissionID = "1b2ecba8-010b-46f8-8976-58dad7812189"
	canCreateContentPermissionID = "a991f301-319b-4311-82cf-277551b71b4e"

	systemAdminRoleID      = "2063dd58-4550-4340-a003-6dcf51d3ee10"
	systemAdminRole        = authorization.DefaultRoleAdmin.String()
	defaultClientRoleID    = "043f12aa-6f51-434f-8e96-35020206f161"
	defaultClientRole      = authorization.DefaultRoleClient.String()
	defaultCaregiverRoleID = "6337eda5-9520-44a6-a4f2-81c32da8dbf2"
	defaultCaregiverRole   = authorization.DefaultRoleCaregiver.String()

	communityID         = "043f12aa-6f51-434f-8e96-35030306f161"
	communityIDToDelete = "043f12aa-6f51-434f-8e96-35030306f162"

	// Appointments
	appointmentID         = "2fc2b603-05ef-40f1-987a-3259eab87aef"
	externalAppointmentID = "5"

	// screeningtools
	screeningToolsQuestionID = "8ecbbc80-24c8-421a-9f1a-e14e12678ef4"
	screeningToolsResponseID = "8ecbbc80-24c8-421a-9f1a-e14e12678ef5"

	clientUnresolvedRequestID     = "8ecbbc80-24c8-421a-9f1a-e14e12678ef6"
	clientUserUnresolvedRequestID = "6ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	pendingServiceRequestID       = "8ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	inProgressServiceRequestID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef8"
	userFailedSecurityCountID     = "07ee2012-18c7-4cc7-8fd8-27249afb091d"
	resolvedServiceRequestID      = "8ecbbc80-24c8-421a-9f1a-e14e12678ef9"
	screeningToolServiceRequestID = "8ecbbc80-24c8-421a-9f1a-e14e12678efa"
	staffUnresolvedRequestID      = "8ecbbc80-24c8-421a-9f1a-e14e12678efb"
	staffUserUnresolvedRequestID  = "8ecbbc80-24c8-421a-9f1a-e14e12678efc"
	userWithRolesID               = "8ecbbc80-24c8-421a-9f1a-e14e12678efd"
	staffWithRolesID              = "8ecbbc80-24c8-421a-9f1a-e14e12678efe"

	userIDToDelete            = "6ecbbc80-24c8-421a-9f7a-e14e12678ef0"
	userIDToRegisterClient    = "6ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	userToRegisterStaff       = "6ecbbc90-24c8-431a-9f7a-e14e12678ef1"
	userToGetNotifications    = "6ecbbc90-24c8-431a-9f7a-e14e12678ef2"
	staffUserIDToDelete       = "6ecbbc80-24c8-421a-9f7a-e14e21678ef0"
	testStaffContact          = "teststaff@staff.com"
	testFlavourPRO            = feedlib.FlavourPro
	fhirPatientID3            = "f933fd4b-1e3c-4ecd-9d7a-82b2790c0544"
	clientID3                 = "11a6a0cd-42ac-714b-97d9-e939a1232de2"
	identifierIDToDelete      = "bcbdaf68-3d36-4365-b575-4392d6749af6"
	staffIdentifierIDToDelete = "bcbdaf89-3d36-4365-b575-4392d6749af7"
	randomIdentifierValue     = "test-identifier-value"
	contactIDToDelete         = "bdc36422-e314-43f2-bb39-ba1ab332f9b0"
	contactIDToRegisterStaff  = "bdc36422-e314-43f2-bb39-ba1ab332f9b1"
	staffContactIDToDelete    = "bdc36422-e314-43f2-bb39-ba1ab332f9c2"
	staffIDToDelete           = "8ecbbc80-24c8-124a-9f1a-e14e12678ef2"
	staffIDToRegister         = "8ecbbc70-24c8-154a-9f1a-e14e13678ef3"
	staffToAddAsClient        = "15b28e4d-4dca-4b80-aed7-0113ab0a20de"

	notificationID = "bf33ba36-30bc-487e-9a7b-bcb54da0bdfe"
	userSurveyID   = "4181df12-ca96-4f28-b78b-8e8ad88b25df"
	feedbackID     = "7281df12-ca96-4f28-b78b-8e8ad88b52df"

	// Questionnaires
	questionnaireID                 = "8ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	screeningToolID                 = "8ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	questionID                      = "8ecbbc80-24c8-421a-9f7a-e14e12678ef4"
	firstChoiceID                   = "8ecbbc80-24c8-421a-9f7a-e14e12678ef0"
	secondChoiceID                  = "8ecbbc80-24c8-421a-9f7a-e14e12678ef1"
	screeningToolQuestionResponseID = "8ecbbc80-24c8-421a-9f7a-e14e12678ef5"
	screeningToolServiceRequestID2  = "8ecbbc80-24c8-421a-9f7a-e14e12678ef6"

	// surveys
	projectID = 1
	formID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef4"

	testCaregiverNumber = "CG0001"

	programID  = "6ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	programID2 = "887dd3ef-3184-4114-86d7-aeafe809f861"

	programName = "test program"
	roomID      = "!vctkCBSzQoVghyPKau:prohealth360.org"
)

// addRequiredContext sets the organisation, program and the user context
func addRequiredContext(ctx context.Context, t *testing.T) context.Context {
	userToken := firebasetools.GetAuthToken(ctx, t)
	userToken.UID = userID
	ctx = context.WithValue(ctx, firebasetools.AuthTokenContextKey, userToken)
	return ctx
}

func TestMain(m *testing.M) {
	isLocalDB := utils.CheckIfCurrentDBIsLocal()
	if !isLocalDB {
		fmt.Println("Cannot run tests. The current database is not a local database.")
		os.Exit(1)
	}

	log.Println("setting up test database")
	var err error

	testingDB, err = gorm.NewPGInstance()
	if err != nil {
		fmt.Println("failed to initialize db:", err)
		os.Exit(1)
	}
	db, err = testingDB.DB.DB()
	if err != nil {
		fmt.Println("failed to initialize db:", err)
		os.Exit(1)
	}

	// setup test variables
	salt, encryptedPin = utils.EncryptPIN("0000", nil)

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateData(template.FuncMap{
			"salt":        salt,
			"hash":        encryptedPin,
			"valid_to":    time.Now().Add(500).String(),
			"test_phone":  "\"" + testPhone + "\"",
			"test_phone2": "\"" + testPhone2 + "\"",

			"test_user_id":                                   userID,
			"user_with_roles_id":                             userWithRolesID,
			"test_user_id2":                                  userID2,
			"test_user_with_caregiver":                       testUserWithCaregiver,
			"test_user_without_caregiver":                    testUserWithoutCaregiver,
			"test_user_has_not_given_consent":                testUserHasNotGivenConsent,
			"staff_user_id":                                  userIDtoAssignStaff,
			"test_staff_user_to_add_as_client":               staffUserToAddAsClient,
			"existing_user_client_to_add_to_another_program": clientUserToAddAsClient,

			"test_flavour":         testFlavour,
			"test_organisation_id": orgID,
			"future_time":          futureTime.String(),
			"test_otp":             "\"" + testOTP + "\"",

			"treatment_buddy_id":  treatmentBuddyID,
			"treatment_buddy_id2": treatmentBuddyID2,

			"test_fhir_patient_id":       fhirPatientID,
			"test_fhir_patient_id2":      fhirPatientID2,
			"test_emr_health_record_id":  testEmrHealthRecordID,
			"test_emr_health_record_id2": testEmrHealthRecordID2,

			"test_facility_id":                                facilityID,
			"facility_identifier_id":                          facilityIdentifierID,
			"facility_to_add_to_user_profile":                 facilityToAddToUserProfile,
			"facility_identifier_to_add_to_user_profile":      facilityIdentifierToAddToUserProfile,
			"facility_to_remove_from_user_profile":            facilityToRemoveFromUserProfile,
			"facility_identifier_to_remove_from_user_profile": facilityIdentifierToRemoveFromUserProfile,
			"facility_to_add_existing_staff":                  facilityToAddExistingStaff,

			"mfl_identifier_value":           mflIdentifier,
			"inactivate_facility_identifier": inactiveFacilityIdentifier,
			"active_facility_identifier":     facilityIdentifierToInactivate,
			"mfl_identifier_type":            mflIdentifierType,
			"test_chv_id":                    testChvID,
			"test_chv_id2":                   testChvID2,
			"test_password":                  gofakeit.Password(false, false, true, true, false, 10),
			"test_terms_id":                  termsID,
			"pro_terms_id":                   proTermsID,
			"consumer_terms_id":              consumerTermsID,
			"test_terms_text":                termsText,
			"security_question_id":           securityQuestionID,
			"security_question_id2":          securityQuestionID2,
			"security_question_id3":          securityQuestionID3,
			"security_question_id4":          securityQuestionID4,

			"security_question_response_id":  securityQuestionResponseID,
			"security_question_response_id2": securityQuestionResponseID2,
			"security_question_response_id3": securityQuestionResponseID3,
			"security_question_response_id4": securityQuestionResponseID4,
			"user_id_to_add_caregiver":       userIDtoAddCaregiver,

			"test_caregiver_id":       testCaregiverID,
			"test_caregiver_org_2_id": testCaregiverOrg2ID,
			"caregiver_with_roles_id": caregiverWithRolesID,

			"staff_number":                        staffNumber,
			"clients_service_request_id":          clientsServiceRequestID,
			"client_service_request_id_to_update": clientServiceRequestIDToUpdate,
			"staff_service_request_id":            staffServiceRequestID,
			"clients_healthdiaryentry_id":         clientsHealthDiaryEntryID,
			"staff_default_facility":              facilityID,
			"staff_id":                            staffID,
			"staff_with_roles_id":                 staffWithRolesID,
			"test_client_caregiver_one_id":        testClientCaregiver1,

			"test_service_request_id": serviceRequestID,
			"test_client_id":          clientID,

			"can_invite_user_permission":   canInviteUserPermissionID,
			"can_resolve_service_request":  canEditOwnRolePermissionID,
			"can_create_screeningtool":     canManageContentPermissionID,
			"can_send_client_survey_links": canCreateContentPermissionID,

			"system_admin_role_id":      systemAdminRoleID,
			"system_admin_role":         systemAdminRole,
			"default_client_role_id":    defaultClientRoleID,
			"default_client_role":       defaultClientRole,
			"default_caregiver_role_id": defaultCaregiverRoleID,
			"default_caregiver_role":    defaultCaregiverRole,

			"community_id":           communityID,
			"community_id_to_delete": communityIDToDelete,

			"screenintoolsquestion_id": screeningToolsQuestionID,
			"screenintoolsresponse_id": screeningToolsResponseID,

			"user_survey_id": userSurveyID,
			"appointment_id": appointmentID,

			"client_user_unresolved_request_id":         clientUserUnresolvedRequestID,
			"test_client_id_with_unresolved_request":    clientUnresolvedRequestID,
			"test_client_with_caregiver":                testClientWithCaregiver,
			"test_client_without_caregiver":             testClientWithoutCaregiver,
			"test_client_has_not_given_consent":         testClientHasNotGivenConsent,
			"existing_client_to_add_to_another_program": testClientToAddToAnotherProgram,

			"pending_service_request_id":        pendingServiceRequestID,
			"in_progress_service_request_id":    inProgressServiceRequestID,
			"user_failed_security_count_id":     userFailedSecurityCountID,
			"resolved_service_request_id":       resolvedServiceRequestID,
			"screening_tool_service_request_id": screeningToolServiceRequestID,
			"staff_unresolved_request_id":       staffUnresolvedRequestID,
			"staff_user_unresolved_request_id":  staffUserUnresolvedRequestID,
			"staff_to_add_as_client":            staffToAddAsClient,

			"test_client_id_to_delete": clientID3,
			"client_with_roles_id":     clientWithRolesID,

			"contact_id_to_delete":            contactIDToDelete,
			"contact_id_to_register_staff":    contactIDToRegisterStaff,
			"staff_contact_id_to_delete":      staffContactIDToDelete,
			"staff_id_to_delete":              staffIDToDelete,
			"staff_id_to_register":            staffIDToRegister,
			"test_ransdom_identifier_value":   randomIdentifierValue,
			"test_staff_identifier_to_delete": staffIdentifierIDToDelete,
			"test_fhir_patient_id3":           fhirPatientID3,
			"test_identifier_to_delete":       identifierIDToDelete,
			"test_flavour_pro":                testFlavourPRO,
			"test_staff_user_id_to_delete":    staffUserIDToDelete,
			"test_user_id_to_delete":          userIDToDelete,
			"test_user_id_to_register_client": userIDToRegisterClient,
			"test_user_id_to_register_staff":  userToRegisterStaff,
			"test_user_to_get_notifications":  userToGetNotifications,
			"test_staff_contact":              testStaffContact,
			"test_feedback_id":                feedbackID,

			"test_questionnaire_id":                                questionnaireID,
			"test_screeningtool_id":                                screeningToolID,
			"test_question_id":                                     questionID,
			"test_first_choice_id":                                 firstChoiceID,
			"test_second_choice_id":                                secondChoiceID,
			"test_questionnaires_screeningtoolquestionresponse_id": screeningToolQuestionResponseID,
			"screening_tool_service_request_id2":                   screeningToolServiceRequestID2,
			"test_caregiver_number":                                testCaregiverNumber,

			"test_project_id": projectID,
			"test_form_id":    formID,

			"test_program_id":  programID,
			"test_program_id2": programID2,

			"org_id_to_add_to_program": orgIDToAddToProgram,
			"program_name":             programName,

			"org_id_to_delete": organisationIDToDelete,
			"org_id_2":         orgID2,
			"test_room_id":     roomID,
		}),
		// this is the directory containing the YAML files.
		// The file name should be the same as the table name
		// order of inserting values matter to avoid foreign key constraint errors
		testfixtures.Paths(
			"../../../../../../fixtures/common_organisation.yml",
			"../../../../../../fixtures/users_user.yml",
			"../../../../../../fixtures/users_termsofservice.yml",
			"../../../../../../fixtures/common_securityquestion.yml",
			"../../../../../../fixtures/common_securityquestionresponse.yml",
			"../../../../../../fixtures/common_contact.yml",
			"../../../../../../fixtures/common_notification.yml",
			"../../../../../../fixtures/users_userotp.yml",
			"../../../../../../fixtures/common_facility.yml",
			"../../../../../../fixtures/common_facility_identifier.yml",
			"../../../../../../fixtures/users_userpin.yml",
			"../../../../../../fixtures/clients_client.yml",
			"../../../../../../fixtures/clients_client_facilities.yml",
			"../../../../../../fixtures/staff_staff.yml",
			"../../../../../../fixtures/staff_staff_identifiers.yml",
			"../../../../../../fixtures/clients_servicerequest.yml",
			"../../../../../../fixtures/staff_staff_facilities.yml",
			"../../../../../../fixtures/staff_servicerequest.yml",
			"../../../../../../fixtures/authority_authoritypermission.yml",
			"../../../../../../fixtures/authority_authorityrole.yml",
			"../../../../../../fixtures/authority_authorityrole_permissions.yml",
			"../../../../../../fixtures/authority_authorityrole_staff.yml",
			"../../../../../../fixtures/authority_authorityrole_clients.yml",
			"../../../../../../fixtures/authority_authorityrole_caregivers.yml",
			"../../../../../../fixtures/communities_community.yml",
			"../../../../../../fixtures/common_identifiers.yml",
			"../../../../../../fixtures/clients_client_identifiers.yml",
			"../../../../../../fixtures/appointments_appointment.yml",
			"../../../../../../fixtures/clients_healthdiaryentry.yml",
			"../../../../../../fixtures/common_usersurveys.yml",
			"../../../../../../fixtures/common_feedback.yml",
			"../../../../../../fixtures/clients_healthdiaryquote.yml",
			"../../../../../../fixtures/questionnaires_questionnaire.yml",
			"../../../../../../fixtures/questionnaires_screeningtool.yml",
			"../../../../../../fixtures/questionnaires_question.yml",
			"../../../../../../fixtures/questionnaires_questioninputchoice.yml",
			"../../../../../../fixtures/questionnaires_screeningtoolresponse.yml",
			"../../../../../../fixtures/questionnaires_screeningtoolquestionresponse.yml",
			"../../../../../../fixtures/caregivers_caregiver.yml",
			"../../../../../../fixtures/caregivers_caregiver_client.yml",
			"../../../../../../fixtures/common_program.yml",
			"../../../../../../fixtures/common_program_facility.yml",
		),
		// uncomment when running tests locally, if your db is not a test db
		// Ensure the testing db in the ci is named `test`
		// !!Warning!!: this can corrupt data, do not turn on or run tests while in non-test db
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		fmt.Println("failed to create fixtures:", err)
		os.Exit(1)

	}

	err = prepareTestDatabase()
	if err != nil {
		fmt.Println("failed to prepare test database:", err)
		os.Exit(1)
	}

	log.Printf("Running tests ...")
	os.Exit(m.Run())
}

func prepareTestDatabase() error {
	if err := fixtures.Load(); err != nil {
		return err
	}
	return nil
}
