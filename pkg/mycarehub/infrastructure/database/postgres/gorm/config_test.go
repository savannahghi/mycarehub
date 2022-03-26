package gorm_test

import (
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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

var (
	fixtures  *testfixtures.Loader
	testingDB *gorm.PGInstance
	orgID     = os.Getenv("DEFAULT_ORG_ID")
	termsID   = 50005
	db        *sql.DB

	newExtension extension.ExternalMethodsExtension

	testPhone   = gofakeit.Phone()
	testFlavour = feedlib.FlavourConsumer
	futureTime  = time.Now().Add(time.Hour * 24 * 365 * 10)
	testOTP     = "1234"
	// user variables
	userID                                     = "6ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userID2                                    = "6ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	userIDtoAddCaregiver                       = "8ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	userIDtoAssignStaff                        = "6ecccc80-24c8-421a-9f1a-e14e13678ef0"
	userIDToInvalidate                         = "5ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToAcceptTerms                        = "4ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToIncreaseFailedLoginCount           = "7ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDtoUpdateLastFailedLoginTime          = "8ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToUpdateUserProfileAfterLoginSuccess = "9ecbbc81-24c8-421a-9f1a-e14e12678ee1"
	userIDToUpdateNextAllowedLoginTime         = "9ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDUpdatePinRequireChangeStatus         = "5ecbbc80-24b8-421a-9f1a-e14e12678ee0"
	userIDToSavePin                            = "8ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	treatmentBuddyID                           = "5ecbbc80-24c8-421a-9f1a-e14e12678ee1"
	treatmentBuddyID2                          = "5ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	fhirPatientID                              = "5ecbbc80-24c8-421a-9f1a-e14e12678ee2"
	fhirPatientID2                             = "f933fd4b-1e3c-4ecd-9d7a-82b2790c0543"
	testEmrHealthRecordID                      = "5ecbbc80-24c8-421a-9f1a-e14e12678ee3"
	testEmrHealthRecordID2                     = "5ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	testChvID                                  = "5ecbbc80-24c8-421a-9f1a-e14e12678ee4"
	testChvID2                                 = "5ecbbc80-24c8-421a-9f1a-e14e12678ef4"
	userNickname                               = "test user"
	clientID                                   = "26b20a42-cbb8-4553-aedb-c539602d04fc"
	clientID2                                  = "00a6a0cd-42ac-417b-97d9-e939a1232de1"
	contactID                                  = "bdc22436-e314-43f2-bb39-ba1ab332f9b0"
	ClientToAddCaregiver                       = "00a6a0cd-42ac-417b-97d9-e939a1232de2"
	// Facility variables
	facilityID          = "4181df12-ca96-4f28-b78b-8e8ad88b25df"
	mflCode             = 324459
	inactiveMflCode     = 229900
	mflCodeToInactivate = 223900
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

	// Content
	contentID  = 10000
	contentID2 = 20000
	authorID   = "4181df12-ca96-4f28-b78b-8e8ad88b25df"
	authorID2  = "4181df12-ca96-4f28-b78b-8e8ad88b25de"

	// Caregiver
	testCaregiverID = "26b20a42-cbb8-4553-aedb-c539602d04fc"

	//Terms
	termsText = "Test terms"

	// Staff
	staffNumber = "test-Staff-101"
	staffID     = "8ecbbc80-24c8-421a-9f1a-e14e12678ef1"

	clientsServiceRequestID   = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	staffServiceRequestID     = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	clientsHealthDiaryEntryID = "8ecbbc10-24c8-421a-9f1a-e17f12678ef1"
	// Service Request
	serviceRequestID = "8ecbbc80-24c8-421a-9f1a-e14e12678ef2"

	// Authority
	canInviteUserPermissionID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	canEditOwnRolePermissionID   = "29672457-d081-48e0-a007-8f49cedb5c6f"
	canManageContentPermissionID = "1b2ecba8-010b-46f8-8976-58dad7812189"
	canCreateContentPermissionID = "a991f301-319b-4311-82cf-277551b71b4e"

	systemAdminRoleID       = "2063dd58-4550-4340-a003-6dcf51d3ee10"
	contentManagementRoleID = "043f12aa-6f51-434f-8e96-35020206f161"
	systemAdminRole         = enums.UserRoleTypeSystemAdministrator.String()
	contentManagementRole   = enums.UserRoleTypeContentManagement.String()

	communityID = "043f12aa-6f51-434f-8e96-35030306f161"

	// Appointments
	appointmentID   = "2fc2b603-05ef-40f1-987a-3259eab87aef"
	appointmentUUID = "d0ba38f2-5a9c-4969-8eb4-beea0a4ff9a5"

	// screeningtools
	screeningToolsQuestionID = "8ecbbc80-24c8-421a-9f1a-e14e12678ef4"
	screeningToolsResponseID = "8ecbbc80-24c8-421a-9f1a-e14e12678ef5"

	clientUnresolvedRequestID     = "8ecbbc80-24c8-421a-9f1a-e14e12678ef6"
	clientUserUnresolvedRequestID = "6ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	pendingServiceRequestID       = "8ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	inProgressServiceRequestID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef8"
	userFailedSecurityCountID     = "07ee2012-18c7-4cc7-8fd8-27249afb091d"
)

func TestMain(m *testing.M) {
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
	newExtension = extension.NewExternalMethodsImpl()
	salt, encryptedPin = newExtension.EncryptPIN("0000", nil)

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Template(),
		testfixtures.TemplateData(template.FuncMap{
			"salt":                       salt,
			"hash":                       encryptedPin,
			"valid_to":                   time.Now().Add(500).String(),
			"test_phone":                 "\"" + testPhone + "\"",
			"test_user_id":               userID,
			"test_user_id2":              userID2,
			"staff_user_id":              userIDtoAssignStaff,
			"test_flavour":               testFlavour,
			"test_organisation_id":       orgID,
			"future_time":                futureTime.String(),
			"test_otp":                   "\"" + testOTP + "\"",
			"treatment_buddy_id":         treatmentBuddyID,
			"treatment_buddy_id2":        treatmentBuddyID2,
			"test_fhir_patient_id":       fhirPatientID,
			"test_fhir_patient_id2":      fhirPatientID2,
			"test_emr_health_record_id":  testEmrHealthRecordID,
			"test_emr_health_record_id2": testEmrHealthRecordID2,
			"test_facility_id":           facilityID,
			"test_chv_id":                testChvID,
			"test_chv_id2":               testChvID2,
			"test_password":              gofakeit.Password(false, false, true, true, false, 10),
			"test_terms_id":              termsID,
			"test_terms_text":            termsText,
			"content_id":                 contentID,
			"content_id2":                contentID2,
			"author_id":                  authorID,
			"author_id2":                 authorID2,
			"security_question_id":       securityQuestionID,
			"security_question_id2":      securityQuestionID2,
			"security_question_id3":      securityQuestionID3,
			"security_question_id4":      securityQuestionID4,

			"security_question_response_id":  securityQuestionResponseID,
			"security_question_response_id2": securityQuestionResponseID2,
			"security_question_response_id3": securityQuestionResponseID3,
			"security_question_response_id4": securityQuestionResponseID4,
			"user_id_to_add_caregiver":       userIDtoAddCaregiver,
			"test_caregiver_id":              testCaregiverID,
			"staff_number":                   staffNumber,
			"clients_service_request_id":     clientsServiceRequestID,
			"staff_service_request_id":       staffServiceRequestID,
			"clients_healthdiaryentry_id":    clientsHealthDiaryEntryID,
			"staff_default_facility":         facilityID,
			"staff_id":                       staffID,

			"test_service_request_id": serviceRequestID,
			"test_client_id":          clientID,

			"can_invite_user_permission":    canInviteUserPermissionID,
			"can_edit_own_role_permission":  canEditOwnRolePermissionID,
			"can_manage_content_permission": canManageContentPermissionID,
			"can_create_content_permission": canCreateContentPermissionID,

			"system_admin_role_id":       systemAdminRoleID,
			"content_management_role_id": contentManagementRoleID,
			"system_admin_role":          systemAdminRole,
			"content_management_role":    contentManagementRole,

			"community_id": communityID,

			"screenintoolsquestion_id": screeningToolsQuestionID,
			"screenintoolsresponse_id": screeningToolsResponseID,

			"client_user_unresolved_request_id":      clientUserUnresolvedRequestID,
			"test_client_id_with_unresolved_request": clientUnresolvedRequestID,
			"pending_service_request_id":             pendingServiceRequestID,
			"in_progress_service_request_id":         inProgressServiceRequestID,
			"user_failed_security_count_id":          userFailedSecurityCountID,
		}),
		// this is the directory containing the YAML files.
		// The file name should be the same as the table name
		// order of inserting values matter to avoid foreign key constraint errors
		testfixtures.Paths(
			"../../../../../../fixtures/common_organisation.yml",
			"../../../../../../fixtures/users_termsofservice.yml",
			"../../../../../../fixtures/clients_securityquestion.yml",
			"../../../../../../fixtures/content_author.yml",
			"../../../../../../fixtures/wagtailcore_page.yml",
			"../../../../../../fixtures/content_contentitem.yml",
			"../../../../../../fixtures/users_user.yml",
			"../../../../../../fixtures/clients_securityquestionresponse.yml",
			"../../../../../../fixtures/common_contact.yml",
			"../../../../../../fixtures/users_userotp.yml",
			"../../../../../../fixtures/common_facility.yml",
			"../../../../../../fixtures/users_userpin.yml",
			"../../../../../../fixtures/clients_caregiver.yml",
			"../../../../../../fixtures/clients_client.yml",
			"../../../../../../fixtures/staff_staff.yml",
			"../../../../../../fixtures/clients_servicerequest.yml",
			"../../../../../../fixtures/staff_staff_facilities.yml",
			"../../../../../../fixtures/authority_authoritypermission.yml",
			"../../../../../../fixtures/authority_authorityrole.yml",
			"../../../../../../fixtures/authority_authorityrole_permissions.yml",
			"../../../../../../fixtures/authority_authorityrole_users.yml",
			"../../../../../../fixtures/communities_community.yml",
			"../../../../../../fixtures/clients_identifier.yml",
			"../../../../../../fixtures/clients_client_identifiers.yml",
			"../../../../../../fixtures/appointments_appointment.yml",
			"../../../../../../fixtures/screeningtools_screeningtoolsquestion.yml",
			"../../../../../../fixtures/screeningtools_screeningtoolsresponse.yml",
			"../../../../../../fixtures/clients_healthdiaryentry.yml",
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
