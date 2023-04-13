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
	userIDSameProgramWithClientID      = "650b7958-12fd-4fa6-9309-ec11618263ae"
	userIDOauthUser                    = "d15c3bb1-bc52-44cc-875e-bf7f4d921dee"

	treatmentBuddyID       = "5ecbbc80-24c8-421a-9f1a-e14e12678ee1"
	treatmentBuddyID2      = "5ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	fhirPatientID          = "5ecbbc80-24c8-421a-9f1a-e14e12678ee2"
	fhirPatientID2         = "f933fd4b-1e3c-4ecd-9d7a-82b2790c0543"
	testEmrHealthRecordID  = "5ecbbc80-24c8-421a-9f1a-e14e12678ee3"
	testEmrHealthRecordID2 = "5ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	testChvID              = "5ecbbc80-24c8-421a-9f1a-e14e12678ee4"
	testChvID2             = "5ecbbc80-24c8-421a-9f1a-e14e12678ef4"

	clientID                         = "26b20a42-cbb8-4553-aedb-c539602d04fc"
	clientID2                        = "00a6a0cd-42ac-417b-97d9-e939a1232de1"
	clientDifferentUserSameProgramID = "b65572ac-d676-4de7-9d9a-031c87b7d2fc"
	clientSameUserDifferentProgramID = "c65cb23c-5c59-40e7-882b-7414af4ca648"
	testClientWithCaregiver          = "f3265be7-54cd-4df9-a078-66bcb31e4dcc"
	testClientWithoutCaregiver       = "13bc475c-6fa8-40a1-ae20-2c9d137ca6e4"
	testClientHasNotGivenConsent     = "5f279d05-0df4-431d-8f70-6f7c76feb425"
	testClientToAddToAnotherProgram  = "01bd8f8d-a1f6-45cf-973d-afb9bde23d87"
	clientWithRolesID                = "79b0aae0-1c42-4b2b-8920-12f7c05dddd9"
	testClientToAssignToCaregiver    = "4a9552c7-ddbb-423c-89b6-626099087b37"

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

	clientUnresolvedRequestID     = "8ecbbc80-24c8-421a-9f1a-e14e12678ef6"
	clientUserUnresolvedRequestID = "6ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	pendingServiceRequestID       = "8ecbbc80-24c8-421a-9f1a-e14e12678ef7"
	inProgressServiceRequestID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef8"
	userFailedSecurityCountID     = "07ee2012-18c7-4cc7-8fd8-27249afb091d"
	resolvedServiceRequestID      = "8ecbbc80-24c8-421a-9f1a-e14e12678ef9"
	staffUnresolvedRequestID      = "8ecbbc80-24c8-421a-9f1a-e14e12678efb"
	staffUserUnresolvedRequestID  = "8ecbbc80-24c8-421a-9f1a-e14e12678efc"
	userWithRolesID               = "8ecbbc80-24c8-421a-9f1a-e14e12678efd"
	staffWithRolesID              = "8ecbbc80-24c8-421a-9f1a-e14e12678efe"

	userIDToDelete              = "6ecbbc80-24c8-421a-9f7a-e14e12678ef0"
	userIDToRegisterClient      = "6ecbbc80-24c8-421a-9f1a-e14e12678ef1"
	userToRegisterStaff         = "6ecbbc90-24c8-431a-9f7a-e14e12678ef1"
	userToGetNotifications      = "6ecbbc90-24c8-431a-9f7a-e14e12678ef2"
	testUserToAssignToCaregiver = "411189bd-4615-4a92-9a0c-f1ca3a3fe1e8"

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
	questionnaireID                                                   = "8ecbbc80-24c8-421a-9f1a-e14e12678ef3"
	questionnaireHasResponseWithin24HoursID                           = "4a639098-0504-43fc-9ca7-0457402ddc42"
	questionnaireHasNoPendingServiceRequestAndResponseWithin24HoursID = "ca24d34d-84e9-426a-8d4a-74635b6d337a"
	questionnaireHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	questionnaireHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	questionnaireHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	questionnaireHasGenderMismatchID                                  = "81e85ba7-430b-4f96-9a0b-d411f2a30258"
	questionnaireHasAgeMismatchID                                     = "966ef828-fd0e-4312-bbb6-8dbb80221f77"
	questionnaireHasClientTypeMismatchID                              = "0bde2dde-31e4-4904-a8f0-35d586d1c841"
	questionnaireSameUserDifferentProgramID                           = "3b036f75-4799-4d38-a493-383e2f437321"
	questionnaireDifferentUserSameProgramID                           = "466890ef-fce9-4cb7-88ab-8cad4fa2c077"

	screeningToolID                                                   = "8ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	screeningToolHasResponseWithin24HoursID                           = "4a639098-0504-43fc-9ca7-0457402ddc42"
	screeningToolHasNoPendingServiceRequestAndResponseWithin24HoursID = "ca24d34d-84e9-426a-8d4a-74635b6d337a"
	screeningToolHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	screeningToolHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	screeningToolHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	screeningToolHasGenderMismatchID                                  = "81e85ba7-430b-4f96-9a0b-d411f2a30258"
	screeningToolHasAgeMismatchID                                     = "966ef828-fd0e-4312-bbb6-8dbb80221f77"
	screeningToolHasClientTypeMismatchID                              = "0bde2dde-31e4-4904-a8f0-35d586d1c841"
	screeningToolSameUserDifferentProgramID                           = "3b036f75-4799-4d38-a493-383e2f437321"
	screeningToolDifferentUserSameProgramID                           = "466890ef-fce9-4cb7-88ab-8cad4fa2c077"

	questionID                                                   = "8ecbbc80-24c8-421a-9f7a-e14e12678ef4"
	questionHasResponseWithin24HoursID                           = "4a639098-0504-43fc-9ca7-0457402ddc42"
	questionHasNoPendingServiceRequestAndResponseWithin24HoursID = "ca24d34d-84e9-426a-8d4a-74635b6d337a"
	questionHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	questionHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	questionHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	questionHasGenderMismatchID                                  = "81e85ba7-430b-4f96-9a0b-d411f2a30258"
	questionHasAgeMismatchID                                     = "966ef828-fd0e-4312-bbb6-8dbb80221f77"
	questionHasClientTypeMismatchID                              = "0bde2dde-31e4-4904-a8f0-35d586d1c841"
	questionSameUserDifferentProgramID                           = "3b036f75-4799-4d38-a493-383e2f437321"
	questionDifferentUserSameProgramID                           = "466890ef-fce9-4cb7-88ab-8cad4fa2c077"

	firstChoiceID                                                   = "8ecbbc80-24c8-421a-9f7a-e14e12678ef0"
	firstChoiceHasResponseWithin24HoursID                           = "4a639098-0504-43fc-9ca7-0457402ddc42"
	firstChoiceHasNoPendingServiceRequestAndResponseWithin24HoursID = "ca24d34d-84e9-426a-8d4a-74635b6d337a"
	firstChoiceHasNoPendingServiceRequestAndResponseAfter24HoursID  = "871d89bb-b249-4838-8938-5f66afdaceff"
	firstChoiceHasPendingServiceRequestAndResponseWithin24HoursID   = "fd7141e8-9213-43e4-bbf8-8c29e3d404b6"
	firstChoiceHasPendingServiceRequestAndResponseAfter24HoursID    = "7a2f18cd-1efc-471f-a2c1-628e4ef4314d"
	firstChoiceHasGenderMismatchID                                  = "81e85ba7-430b-4f96-9a0b-d411f2a30258"
	firstChoiceHasAgeMismatchID                                     = "966ef828-fd0e-4312-bbb6-8dbb80221f77"
	firstChoiceHasClientTypeMismatchID                              = "0bde2dde-31e4-4904-a8f0-35d586d1c841"
	firstChoiceSameUserDifferentProgramID                           = "3b036f75-4799-4d38-a493-383e2f437321"
	firstChoiceDifferentUserSameProgramID                           = "466890ef-fce9-4cb7-88ab-8cad4fa2c077"

	secondChoiceID                                                   = "8ecbbc80-24c8-421a-9f7a-e14e12678ef1"
	secondChoiceHasResponseWithin24HoursID                           = "5a639098-0504-43fc-9ca7-0457402ddc42"
	secondChoiceHasNoPendingServiceRequestAndResponseWithin24HoursID = "da24d34d-84e9-426a-8d4a-74635b6d337a"
	secondChoiceHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	secondChoiceHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	secondChoiceHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	secondChoiceHasGenderMismatchID                                  = "71e85ba7-430b-4f96-9a0b-d411f2a30258"
	secondChoiceHasAgeMismatchID                                     = "866ef828-fd0e-4312-bbb6-8dbb80221f77"
	secondChoiceHasClientTypeMismatchID                              = "4bde2dde-31e4-4904-a8f0-35d586d1c841"
	secondChoiceSameUserDifferentProgramID                           = "6b036f75-4799-4d38-a493-383e2f437321"
	secondChoiceDifferentUserSameProgramID                           = "566890ef-fce9-4cb7-88ab-8cad4fa2c077"

	screeningToolsResponseID                                                   = "8ecbbc80-24c8-421a-9f1a-e14e12678ef5"
	screeningToolsResponseHasResponseWithin24HoursID                           = "5a639098-0504-43fc-9ca7-0457402ddc42"
	screeningToolsResponseHasNoPendingServiceRequestAndResponseWithin24HoursID = "da24d34d-84e9-426a-8d4a-74635b6d337a"
	screeningToolsResponseHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	screeningToolsResponseHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	screeningToolsResponseHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	screeningToolsResponseHasGenderMismatchID                                  = "71e85ba7-430b-4f96-9a0b-d411f2a30258"
	screeningToolsResponseHasAgeMismatchID                                     = "866ef828-fd0e-4312-bbb6-8dbb80221f77"
	screeningToolsResponseHasClientTypeMismatchID                              = "4bde2dde-31e4-4904-a8f0-35d586d1c841"
	screeningToolsResponseSameUserDifferentProgramID                           = "6b036f75-4799-4d38-a493-383e2f437321"
	screeningToolsResponseDifferentUserSameProgramID                           = "566890ef-fce9-4cb7-88ab-8cad4fa2c077"

	screeningToolQuestionResponseID                                                   = "8ecbbc80-24c8-421a-9f7a-e14e12678ef5"
	screeningToolQuestionResponseHasResponseWithin24HoursID                           = "5a639098-0504-43fc-9ca7-0457402ddc42"
	screeningToolQuestionResponseHasNoPendingServiceRequestAndResponseWithin24HoursID = "da24d34d-84e9-426a-8d4a-74635b6d337a"
	screeningToolQuestionResponseHasNoPendingServiceRequestAndResponseAfter24HoursID  = "94fa4b2d-53b5-4065-9932-6d4a802db3e5"
	screeningToolQuestionResponseHasPendingServiceRequestAndResponseWithin24HoursID   = "dea431ec-2bdd-4e8b-9a43-f30486391dec"
	screeningToolQuestionResponseHasPendingServiceRequestAndResponseAfter24HoursID    = "33f9f045-3b78-4441-8cc1-347f181313b3"
	screeningToolQuestionResponseHasGenderMismatchID                                  = "71e85ba7-430b-4f96-9a0b-d411f2a30258"
	screeningToolQuestionResponseHasAgeMismatchID                                     = "866ef828-fd0e-4312-bbb6-8dbb80221f77"
	screeningToolQuestionResponseHasClientTypeMismatchID                              = "4bde2dde-31e4-4904-a8f0-35d586d1c841"
	screeningToolQuestionResponseSameUserDifferentProgramID                           = "6b036f75-4799-4d38-a493-383e2f437321"
	screeningToolQuestionResponseDifferentUserSameProgramID                           = "566890ef-fce9-4cb7-88ab-8cad4fa2c077"

	serviceRequestIDHasNoPendingServiceRequestAndResponseWithin24HoursID = "8ecbbc80-24c8-421a-9f7a-e14e12678ef6"
	serviceRequestIDHasNoPendingServiceRequestAndResponseAfter24HoursID  = "6cc3099d-d467-4e11-bfdf-01d1941ce28a"
	serviceRequestIDHasPendingServiceRequestAndResponseWithin24HoursID   = "45d205bd-c05f-4634-876d-d5b56b64e97e"
	serviceRequestIDHasPendingServiceRequestAndResponseAfter24HoursID    = "8ff43e6d-61bc-49fa-812c-2f67cd146761"

	// surveys
	projectID = 1
	formID    = "8ecbbc80-24c8-421a-9f1a-e14e12678ef4"

	testCaregiverNumber = "CG0001"

	programID  = "6ecbbc80-24c8-421a-9f1a-e14e12678ee0"
	programID2 = "887dd3ef-3184-4114-86d7-aeafe809f861"

	programName = "test program"
	roomID      = "!vctkCBSzQoVghyPKau:prohealth360.org"

	oauthClientOneID = "548394d2-1992-40eb-b82e-ac56f08e779c"

	oauthSessionOneID = "1631203b-0182-4d4d-9a6c-4b270759427d"
	oauthSessionTwoID = "2c3a5a48-b638-4e21-9460-297af43331f7"

	oauthAuthorizationCode = "e455b001-faa4-42ec-835f-16dec96d68d9"
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
			"test_user_to_assign_to_caregiver":               testUserToAssignToCaregiver,
			"test_user_has_not_given_consent":                testUserHasNotGivenConsent,
			"staff_user_id":                                  userIDtoAssignStaff,
			"test_staff_user_to_add_as_client":               staffUserToAddAsClient,
			"existing_user_client_to_add_to_another_program": clientUserToAddAsClient,
			"test_user_id_different_user_same_program":       userIDSameProgramWithClientID,
			"test_oauth_user_id":                             userIDOauthUser,

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
			"test_client_to_assign_to_caregiver":  testClientToAssignToCaregiver,

			"test_service_request_id": serviceRequestID,

			"test_client_id": clientID,
			"test_client_id_same_user_different_program": clientSameUserDifferentProgramID,
			"test_client_id_different_user_same_program": clientDifferentUserSameProgramID,

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

			"user_survey_id": userSurveyID,
			"appointment_id": appointmentID,

			"client_user_unresolved_request_id":         clientUserUnresolvedRequestID,
			"test_client_id_with_unresolved_request":    clientUnresolvedRequestID,
			"test_client_with_caregiver":                testClientWithCaregiver,
			"test_client_without_caregiver":             testClientWithoutCaregiver,
			"test_client_has_not_given_consent":         testClientHasNotGivenConsent,
			"existing_client_to_add_to_another_program": testClientToAddToAnotherProgram,

			"pending_service_request_id":       pendingServiceRequestID,
			"in_progress_service_request_id":   inProgressServiceRequestID,
			"user_failed_security_count_id":    userFailedSecurityCountID,
			"resolved_service_request_id":      resolvedServiceRequestID,
			"staff_unresolved_request_id":      staffUnresolvedRequestID,
			"staff_user_unresolved_request_id": staffUserUnresolvedRequestID,
			"staff_to_add_as_client":           staffToAddAsClient,

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

			"test_questionnaire_id":                                                             questionnaireID,
			"test_questionnaire_id_has_response_within_24_hours":                                questionnaireHasResponseWithin24HoursID,
			"test_questionnaire_id_has_no_pending_service_request_and_response_within_24_hours": questionnaireHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_questionnaire_id_has_no_pending_service_request_and_response_after_24_hours":  questionnaireHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_questionnaire_id_has_pending_service_request_and_response_within_24_hours":    questionnaireHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_questionnaire_id_has_pending_service_request_and_response_after_24_hours":     questionnaireHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_questionnaire_id_client_gender_mismatch":                                      questionnaireHasGenderMismatchID,
			"test_questionnaire_id_age_mismatch":                                                questionnaireHasAgeMismatchID,
			"test_questionnaire_id_client_types_mismatch":                                       questionnaireHasClientTypeMismatchID,
			"test_questionnaire_id_same_user_different_program":                                 questionnaireSameUserDifferentProgramID,
			"test_questionnaire_id_different_user_same_program":                                 questionnaireDifferentUserSameProgramID,

			"test_screeningtool_id":                                                             screeningToolID,
			"test_screeningtool_id_has_response_within_24_hours":                                screeningToolHasResponseWithin24HoursID,
			"test_screeningtool_id_has_no_pending_service_request_and_response_within_24_hours": screeningToolHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_screeningtool_id_has_no_pending_service_request_and_response_after_24_hours":  screeningToolHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_screeningtool_id_has_pending_service_request_and_response_within_24_hours":    screeningToolHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_screeningtool_id_has_pending_service_request_and_response_after_24_hours":     screeningToolHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_screeningtool_id_client_gender_mismatch":                                      screeningToolHasGenderMismatchID,
			"test_screeningtool_id_age_mismatch":                                                screeningToolHasAgeMismatchID,
			"test_screeningtool_id_client_types_mismatch":                                       screeningToolHasClientTypeMismatchID,
			"test_screeningtool_id_same_user_different_program":                                 screeningToolSameUserDifferentProgramID,
			"test_screeningtool_id_different_user_same_program":                                 screeningToolDifferentUserSameProgramID,

			"test_question_id": questionID,
			"test_question_id_has_response_within_24_hours":                                questionHasResponseWithin24HoursID,
			"test_question_id_has_no_pending_service_request_and_response_within_24_hours": questionHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_question_id_has_no_pending_service_request_and_response_after_24_hours":  questionHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_question_id_has_pending_service_request_and_response_within_24_hours":    questionHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_question_id_has_pending_service_request_and_response_after_24_hours":     questionHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_question_id_client_gender_mismatch":                                      questionHasGenderMismatchID,
			"test_question_id_age_mismatch":                                                questionHasAgeMismatchID,
			"test_question_id_client_types_mismatch":                                       questionHasClientTypeMismatchID,
			"test_question_id_same_user_different_program":                                 questionSameUserDifferentProgramID,
			"test_question_id_different_user_same_program":                                 questionDifferentUserSameProgramID,

			"test_first_choice_id":                                                             firstChoiceID,
			"test_first_choice_id_has_response_within_24_hours":                                firstChoiceHasResponseWithin24HoursID,
			"test_first_choice_id_has_no_pending_service_request_and_response_within_24_hours": firstChoiceHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_first_choice_id_has_no_pending_service_request_and_response_after_24_hours":  firstChoiceHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_first_choice_id_has_pending_service_request_and_response_within_24_hours":    firstChoiceHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_first_choice_id_has_pending_service_request_and_response_after_24_hours":     firstChoiceHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_first_choice_id_client_gender_mismatch":                                      firstChoiceHasGenderMismatchID,
			"test_first_choice_id_age_mismatch":                                                firstChoiceHasAgeMismatchID,
			"test_first_choice_id_client_types_mismatch":                                       firstChoiceHasClientTypeMismatchID,
			"test_first_choice_id_same_user_different_program":                                 firstChoiceSameUserDifferentProgramID,
			"test_first_choice_id_different_user_same_program":                                 firstChoiceDifferentUserSameProgramID,

			"test_second_choice_id":                                                             secondChoiceID,
			"test_second_choice_id_has_response_within_24_hours":                                secondChoiceHasResponseWithin24HoursID,
			"test_second_choice_id_has_no_pending_service_request_and_response_within_24_hours": secondChoiceHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_second_choice_id_has_no_pending_service_request_and_response_after_24_hours":  secondChoiceHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_second_choice_id_has_pending_service_request_and_response_within_24_hours":    secondChoiceHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_second_choice_id_has_pending_service_request_and_response_after_24_hours":     secondChoiceHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_second_choice_id_client_gender_mismatch":                                      secondChoiceHasGenderMismatchID,
			"test_second_choice_id_age_mismatch":                                                secondChoiceHasAgeMismatchID,
			"test_second_choice_id_client_types_mismatch":                                       secondChoiceHasClientTypeMismatchID,
			"test_second_choice_id_same_user_different_program":                                 secondChoiceSameUserDifferentProgramID,
			"test_second_choice_id_different_user_same_program":                                 secondChoiceDifferentUserSameProgramID,

			"screenintoolsresponse_id":                                                             screeningToolsResponseID,
			"screenintoolsresponse_id_has_response_within_24_hours":                                screeningToolsResponseHasResponseWithin24HoursID,
			"screenintoolsresponse_id_has_no_pending_service_request_and_response_within_24_hours": screeningToolsResponseHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"screenintoolsresponse_id_has_no_pending_service_request_and_response_after_24_hours":  screeningToolsResponseHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"screenintoolsresponse_id_has_pending_service_request_and_response_within_24_hours":    screeningToolsResponseHasPendingServiceRequestAndResponseWithin24HoursID,
			"screenintoolsresponse_id_has_pending_service_request_and_response_after_24_hours":     screeningToolsResponseHasPendingServiceRequestAndResponseAfter24HoursID,
			"screenintoolsresponse_id_client_gender_mismatch":                                      screeningToolsResponseHasGenderMismatchID,
			"screenintoolsresponse_id_age_mismatch":                                                screeningToolsResponseHasAgeMismatchID,
			"screenintoolsresponse_id_client_types_mismatch":                                       screeningToolsResponseHasClientTypeMismatchID,
			"screenintoolsresponse_id_same_user_different_program":                                 screeningToolsResponseSameUserDifferentProgramID,
			"screenintoolsresponse_id_different_user_same_program":                                 screeningToolsResponseDifferentUserSameProgramID,

			"test_screening_tool_question_response_id":                                                             screeningToolQuestionResponseID,
			"test_screening_tool_question_response_id_has_response_within_24_hours":                                screeningToolQuestionResponseHasResponseWithin24HoursID,
			"test_screening_tool_question_response_id_has_no_pending_service_request_and_response_within_24_hours": screeningToolQuestionResponseHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"test_screening_tool_question_response_id_has_no_pending_service_request_and_response_after_24_hours":  screeningToolQuestionResponseHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"test_screening_tool_question_response_id_has_pending_service_request_and_response_within_24_hours":    screeningToolQuestionResponseHasPendingServiceRequestAndResponseWithin24HoursID,
			"test_screening_tool_question_response_id_has_pending_service_request_and_response_after_24_hours":     screeningToolQuestionResponseHasPendingServiceRequestAndResponseAfter24HoursID,
			"test_screening_tool_question_response_id_client_gender_mismatch":                                      screeningToolQuestionResponseHasGenderMismatchID,
			"test_screening_tool_question_response_id_age_mismatch":                                                screeningToolQuestionResponseHasAgeMismatchID,
			"test_screening_tool_question_response_id_client_types_mismatch":                                       screeningToolQuestionResponseHasClientTypeMismatchID,
			"test_screening_tool_question_response_id_same_user_different_program":                                 screeningToolQuestionResponseSameUserDifferentProgramID,
			"test_screening_tool_question_response_id_different_user_same_program":                                 screeningToolQuestionResponseDifferentUserSameProgramID,

			"service_request_id_has_no_pending_service_request_and_response_within_24_hours": serviceRequestIDHasNoPendingServiceRequestAndResponseWithin24HoursID,
			"service_request_id_has_no_pending_service_request_and_response_after_24_hours":  serviceRequestIDHasNoPendingServiceRequestAndResponseAfter24HoursID,
			"service_request_id_has_pending_service_request_and_response_within_24_hours":    serviceRequestIDHasPendingServiceRequestAndResponseWithin24HoursID,
			"service_request_id_has_pending_service_request_and_response_after_24_hours":     serviceRequestIDHasPendingServiceRequestAndResponseAfter24HoursID,

			"test_caregiver_number": testCaregiverNumber,

			"test_project_id": projectID,
			"test_form_id":    formID,

			"test_program_id":  programID,
			"test_program_id2": programID2,

			"org_id_to_add_to_program": orgIDToAddToProgram,
			"program_name":             programName,

			"org_id_to_delete":      organisationIDToDelete,
			"test_organisation_id2": orgID2,
			"test_room_id":          roomID,

			"test_oauth_client_one": oauthClientOneID,

			"test_oauth_session_one": oauthSessionOneID,
			"test_oauth_session_two": oauthSessionTwoID,

			"test_oauth_auth_code_one": oauthAuthorizationCode,
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
			"../../../../../../fixtures/oauth_client.yml",
			"../../../../../../fixtures/oauth_session.yml",
			"../../../../../../fixtures/oauth_authorization_code.yml",
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
