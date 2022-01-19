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
	userID                             = "6ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userID2                            = "6ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	userIDtoAddCaregiver               = "8ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	userIDToInvalidate                 = "5ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToAcceptTerms                = "4ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToIncreaseFailedLoginCount   = "7ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDtoUpdateLastFailedLoginTime  = "8ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDToUpdateNextAllowedLoginTime = "9ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	userIDUpdatePinRequireChangeStatus = "5ecbbc80-24b8-421a-9f1a-e14e12678ee0"
	userIDToSavePin                    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef0"
	treatmentBuddyID                   = "5ecbbc80-24c8-421a-9f1a-e14e12678ee1"
	treatmentBuddyID2                  = "5ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	fhirPatientID                      = "5ecbbc80-24c8-421a-9f1a-e14e12678ee2"
	fhirPatientID2                     = "f933fd4b-1e3c-4ecd-9d7a-82b2790c0543"
	testEmrHealthRecordID              = "5ecbbc80-24c8-421a-9f1a-e14e12678ee3"
	testEmrHealthRecordID2             = "5ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	testChvID                          = "5ecbbc80-24c8-421a-9f1a-e14e12678ee4"
	testChvID2                         = "5ecbbc80-24c8-421a-9f1a-e14e12678ef4"
	userNickname                       = "test user"
	clientID                           = "26b20a42-cbb8-4553-aedb-c539602d04fc"
	clientID2                          = "00a6a0cd-42ac-417b-97d9-e939a1232de1"
	ClientToAddCaregiver               = "00a6a0cd-42ac-417b-97d9-e939a1232de2"
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

	// contact variables
	// contactID = "bdc22436-e314-43f2-bb39-ba1ab332f9b0"
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
