package gorm_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_GetOrCreateFacility(t *testing.T) {

	name := ksuid.New().String()
	county := gofakeit.Name()
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()
	identifierValue := strconv.Itoa(gofakeit.Number(3000, 100000))

	facility := &gorm.Facility{
		Name:               name,
		Active:             true,
		Country:            county,
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
	}

	identifier := &gorm.FacilityIdentifier{
		Type:  enums.FacilityIdentifierTypeMFLCode.String(),
		Value: identifierValue,
	}

	type args struct {
		ctx        context.Context
		facility   *gorm.Facility
		identifier *gorm.FacilityIdentifier
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get or create facility",
			args: args{
				ctx:        addOrganizationContext(context.Background()),
				facility:   facility,
				identifier: identifier,
			},
			wantErr: false,
		},
		// {
		// 	name: "Sad Case - Fail to get or create facility",
		// 	args: args{
		// 		ctx: addOrganizationContext(context.Background()),
		// 		facility: &gorm.Facility{
		// 			Name:        name,
		// 			Active:      true,
		// 			Country:     gofakeit.HipsterSentence(50),
		// 			Description: description,
		// 		},
		// 		identifier: identifier,
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetOrCreateFacility(tt.args.ctx, tt.args.facility, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
	// // teardown
	// pg, err := gorm.NewPGInstance()
	// if err != nil {
	// 	t.Errorf("pgInstance.Teardown() = %v", err)
	// }
	// if err = pg.DB.Where("id", facility.FacilityID).Unscoped().Delete(&gorm.Facility{}).Error; err != nil {
	// 	t.Errorf("failed to delete record = %v", err)
	// }
}

func TestPGInstance_SaveTemporaryUserPin(t *testing.T) {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer

	pinPayload := &gorm.PINData{
		UserID:    userIDToSavePin,
		HashedPIN: encryptedPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   flavour,
		Salt:      salt,
	}

	invalidPinPayload := &gorm.PINData{
		HashedPIN: encryptedPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   flavour,
		Salt:      salt,
	}

	type args struct {
		ctx        context.Context
		pinPayload *gorm.PINData
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:        addOrganizationContext(context.Background()),
				pinPayload: pinPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing payload",
			args: args{
				ctx: addOrganizationContext(context.Background()),
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "invalid: invalid payload",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				pinPayload: &gorm.PINData{
					UserID:    userIDToSavePin,
					HashedPIN: encryptedPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   feedlib.Flavour(gofakeit.HipsterSentence(30)),
					Salt:      salt,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: no userID",
			args: args{
				ctx:        addOrganizationContext(context.Background()),
				pinPayload: invalidPinPayload,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SaveTemporaryUserPin(tt.args.ctx, tt.args.pinPayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SaveTemporaryUserPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SaveTemporaryUserPin() = %v, want %v", got, tt.want)
			}
		})
	}

	// Teardown
	if err = pg.DB.Where("user_id", userIDToSavePin).Unscoped().Delete(&gorm.PINData{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}

}

func TestPGInstance_SavePin(t *testing.T) {

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	longString := gofakeit.Sentence(300)

	pinPayload := &gorm.PINData{
		UserID:    userIDToSavePin,
		HashedPIN: encryptedPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
		Salt:      salt,
	}

	type args struct {
		ctx     context.Context
		pinData *gorm.PINData
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     addOrganizationContext(context.Background()),
				pinData: pinPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing user id",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				pinData: &gorm.PINData{
					HashedPIN: encryptedPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   feedlib.FlavourConsumer,
					Salt:      salt,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: user does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				pinData: &gorm.PINData{
					UserID:    ksuid.New().String(),
					HashedPIN: encryptedPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   feedlib.FlavourConsumer,
					Salt:      salt,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid user id",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				pinData: &gorm.PINData{
					UserID:    longString,
					HashedPIN: encryptedPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   feedlib.FlavourConsumer,
					Salt:      salt,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SavePin(tt.args.ctx, tt.args.pinData)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SavePin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SavePin() = %v, want %v", got, tt.want)
			}
		})
	}

	// Teardown
	if err := pg.DB.Where("user_id", userIDToSavePin).Unscoped().Delete(&gorm.PINData{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_SaveSecurityQuestionResponse(t *testing.T) {

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	type args struct {
		ctx                      context.Context
		securityQuestionResponse []*gorm.SecurityQuestionResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				securityQuestionResponse: []*gorm.SecurityQuestionResponse{
					{
						QuestionID: securityQuestionID,
						UserID:     userID,
						Response:   "1917",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.SaveSecurityQuestionResponse(tt.args.ctx, tt.args.securityQuestionResponse); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SaveSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Teardown
	if err := pg.DB.Where("user_id", userID).Unscoped().Delete(&gorm.SecurityQuestionResponse{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_SaveOTP(t *testing.T) {

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	generatedAt := time.Now()
	validUntil := time.Now().AddDate(0, 0, 2)

	otp, err := utils.GenerateOTP()
	if err != nil {
		t.Errorf("unable to generate OTP")
	}

	otpInput := &gorm.UserOTP{
		OTPID:       gofakeit.Number(100, 120020),
		UserID:      userID,
		Valid:       true,
		GeneratedAt: generatedAt,
		ValidUntil:  validUntil,
		Channel:     "SMS",
		Flavour:     feedlib.FlavourConsumer,
		PhoneNumber: "+254710000111",
		OTP:         otp,
	}

	err = pg.DB.Create(&otpInput).Error
	if err != nil {
		t.Errorf("failed to create otp: %v", err)
	}

	newOTP, err := utils.GenerateOTP()
	if err != nil {
		t.Errorf("unable to generate OTP")
	}

	gormOTPInput := &gorm.UserOTP{
		OTPID:       gofakeit.Number(100, 200200),
		UserID:      userID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		Flavour:     otpInput.Flavour,
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         newOTP,
	}

	invalidgormOTPInput1 := &gorm.UserOTP{
		UserID:      userID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		Flavour:     otpInput.Flavour,
		PhoneNumber: "",
		OTP:         newOTP,
	}

	invalidgormOTPInput2 := &gorm.UserOTP{
		UserID:      userID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		Flavour:     feedlib.Flavour("Invalid-flavour"),
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         newOTP,
	}

	invalidgormOTPInput3 := &gorm.UserOTP{
		UserID:      userID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		Flavour:     "invalid",
		PhoneNumber: "",
		OTP:         newOTP,
	}

	type args struct {
		ctx      context.Context
		otpInput *gorm.UserOTP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				otpInput: gormOTPInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				otpInput: invalidgormOTPInput1,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				otpInput: invalidgormOTPInput2,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavour and phone",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				otpInput: invalidgormOTPInput3,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				otpInput: &gorm.UserOTP{
					UserID:      userID,
					Valid:       otpInput.Valid,
					GeneratedAt: otpInput.GeneratedAt,
					ValidUntil:  otpInput.ValidUntil,
					Channel:     otpInput.Channel,
					Flavour:     feedlib.Flavour(gofakeit.HipsterSentence(30)),
					PhoneNumber: otpInput.PhoneNumber,
					OTP:         newOTP,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				otpInput: &gorm.UserOTP{
					UserID:      userID,
					Valid:       otpInput.Valid,
					GeneratedAt: otpInput.GeneratedAt,
					ValidUntil:  otpInput.ValidUntil,
					Channel:     otpInput.Channel,
					Flavour:     otpInput.Flavour,
					PhoneNumber: otpInput.PhoneNumber,
					OTP:         "12345678910",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.SaveOTP(tt.args.ctx, tt.args.otpInput); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SaveOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Teardown
	if err = pg.DB.Where("id", otpInput.OTPID).Unscoped().Delete(&gorm.UserOTP{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", gormOTPInput.OTPID).Unscoped().Delete(&gorm.UserOTP{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_CreateServiceRequest(t *testing.T) {

	testTime := time.Now()
	meta := `{"test":"test"}`
	serviceRequestInput := &gorm.ClientServiceRequest{
		Active:         false,
		RequestType:    "HealthDiary",
		Request:        gofakeit.Sentence(5),
		Status:         "PENDING",
		InProgressAt:   &testTime,
		ResolvedAt:     &testTime,
		ClientID:       clientID,
		OrganisationID: orgID,
		FacilityID:     facilityID,
		Meta:           meta,
		ProgramID:      programID,
	}
	InvalidServiceRequestInput := &gorm.ClientServiceRequest{
		Active:         false,
		RequestType:    "HealthDiary",
		Request:        gofakeit.Sentence(5),
		Status:         "PENDING",
		InProgressAt:   &testTime,
		ResolvedAt:     &testTime,
		ClientID:       clientID,
		OrganisationID: orgID,
		FacilityID:     facilityID,
		Meta:           "",
	}
	type args struct {
		ctx                 context.Context
		serviceRequestInput *gorm.ClientServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                 addOrganizationContext(context.Background()),
				serviceRequestInput: serviceRequestInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid meta data",
			args: args{
				ctx:                 addOrganizationContext(context.Background()),
				serviceRequestInput: InvalidServiceRequestInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.CreateServiceRequest(tt.args.ctx, tt.args.serviceRequestInput); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateHealthDiaryEntry(t *testing.T) {

	clientHealthDiaryEntryID := uuid.New().String()
	currentTime := time.Now()

	type args struct {
		ctx              context.Context
		healthDiaryInput *gorm.ClientHealthDiaryEntry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientHealthDiaryEntryID,
					Active:                   true,
					Mood:                     "Very Cool",
					Note:                     "I'm happy",
					EntryType:                "HOME_PAGE_HEALTH_DIARY_ENTRY",
					ShareWithHealthWorker:    true,
					SharedAt:                 &currentTime,
					ClientID:                 clientID,
					ProgramID:                programID,
					OrganisationID:           uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					Active:                true,
					Mood:                  gofakeit.HipsterSentence(20),
					Note:                  "test",
					EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY",
					ShareWithHealthWorker: false,
					SharedAt:              &currentTime,
					ClientID:              clientID,
					OrganisationID:        orgID,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - no client ID",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientHealthDiaryEntryID,
					Active:                   true,
					Mood:                     "Very Cool",
					Note:                     "I'm happy",
					EntryType:                "Test",
					ShareWithHealthWorker:    true,
					SharedAt:                 &currentTime,
					OrganisationID:           uuid.New().String(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - no health diary input",
			args: args{
				ctx:              addOrganizationContext(context.Background()),
				healthDiaryInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if _, err := testingDB.CreateHealthDiaryEntry(tt.args.ctx, tt.args.healthDiaryInput); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateHealthDiaryEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateCommunity(t *testing.T) {

	var genderList pq.StringArray
	for _, g := range enumutils.AllGender {
		genderList = append(genderList, string(g))
	}

	var clientTypeList pq.StringArray
	for _, c := range enums.AllClientType {
		clientTypeList = append(clientTypeList, string(c))
	}

	type args struct {
		ctx       context.Context
		community *gorm.Community
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				community: &gorm.Community{
					Name:           "test",
					Description:    "test",
					Active:         true,
					MinimumAge:     19,
					MaximumAge:     30,
					Gender:         genderList,
					ClientTypes:    clientTypeList,
					InviteOnly:     true,
					Discoverable:   true,
					ProgramID:      programID,
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				community: &gorm.Community{
					Name:           "test",
					Description:    "test",
					Active:         true,
					MinimumAge:     19,
					MaximumAge:     30,
					InviteOnly:     true,
					Discoverable:   true,
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.CreateCommunity(tt.args.ctx, tt.args.community)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_GetOrCreateNextOfKin(t *testing.T) {
	type args struct {
		ctx       context.Context
		person    *gorm.RelatedPerson
		clientID  string
		contactID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create related person",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				person: &gorm.RelatedPerson{
					Active:           true,
					FirstName:        gofakeit.Name(),
					LastName:         gofakeit.Name(),
					Gender:           "MALE",
					RelationshipType: "Next of Kin",
					ProgramID:        programID,
				},
				clientID:  clientID,
				contactID: contactID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.GetOrCreateNextOfKin(tt.args.ctx, tt.args.person, tt.args.clientID, tt.args.contactID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOrCreateNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_GetOrCreateContact(t *testing.T) {
	type args struct {
		ctx     context.Context
		contact *gorm.Contact
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create a contacts",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				contact: &gorm.Contact{
					Active:       true,
					ContactType:  "Phone",
					ContactValue: gofakeit.Phone(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetOrCreateContact(tt.args.ctx, tt.args.contact)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOrCreateContact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected contact to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil && got.ContactID == nil {
				t.Errorf("expected contact not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_AnswerScreeningToolQuestions(t *testing.T) {
	type args struct {
		ctx                    context.Context
		screeningToolResponses []*gorm.ScreeningToolsResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create screening tool responses",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				screeningToolResponses: []*gorm.ScreeningToolsResponse{
					{
						QuestionID: screeningToolsQuestionID,
						ClientID:   clientID,
						Response:   "0",
						ProgramID:  programID,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.AnswerScreeningToolQuestions(tt.args.ctx, tt.args.screeningToolResponses); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AnswerScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateAppointment(t *testing.T) {

	type args struct {
		ctx         context.Context
		appointment *gorm.Appointment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create an appointment",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				appointment: &gorm.Appointment{
					Active:                    true,
					ExternalID:                strconv.Itoa(gofakeit.Number(0, 1000)),
					ClientID:                  clientID,
					FacilityID:                facilityID,
					Reason:                    "Dental",
					Date:                      time.Now().Add(time.Duration(100)),
					HasRescheduledAppointment: true,
					ProgramID:                 programID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create an appointment",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				appointment: &gorm.Appointment{
					Active:                    true,
					ExternalID:                strconv.Itoa(gofakeit.Number(0, 1000)),
					ClientID:                  clientID,
					FacilityID:                "facilityID",
					Reason:                    "Dental",
					Date:                      time.Now().Add(time.Duration(100)),
					HasRescheduledAppointment: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateAppointment(tt.args.ctx, tt.args.appointment); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateStaffServiceRequest(t *testing.T) {

	ID := uuid.New().String()
	currentTime := time.Now()
	meta := `{"test":"test"}`

	type args struct {
		ctx                 context.Context
		serviceRequestInput *gorm.StaffServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				serviceRequestInput: &gorm.StaffServiceRequest{
					ID:                &staffServiceRequestID,
					Active:            true,
					RequestType:       gofakeit.BeerName(),
					Request:           gofakeit.BeerName(),
					Status:            gofakeit.BeerName(),
					ResolvedAt:        &currentTime,
					StaffID:           staffID,
					OrganisationID:    orgID,
					ProgramID:         programID,
					DefaultFacilityID: &facilityID,
					Meta:              meta,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				serviceRequestInput: &gorm.StaffServiceRequest{
					ID:             &ID,
					Active:         true,
					RequestType:    gofakeit.BeerName(),
					Request:        gofakeit.BeerName(),
					Status:         gofakeit.BeerName(),
					ResolvedAt:     &currentTime,
					StaffID:        "staffID",
					OrganisationID: orgID,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid metadata",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				serviceRequestInput: &gorm.StaffServiceRequest{
					ID:             &ID,
					Active:         true,
					RequestType:    gofakeit.BeerName(),
					Request:        gofakeit.BeerName(),
					Status:         gofakeit.BeerName(),
					ResolvedAt:     &currentTime,
					StaffID:        "staffID",
					OrganisationID: orgID,
					Meta:           "meta",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateStaffServiceRequest(tt.args.ctx, tt.args.serviceRequestInput); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateStaffServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateUser(t *testing.T) {
	date := gofakeit.Date()

	type args struct {
		ctx  context.Context
		user *gorm.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a new user",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				user: &gorm.User{
					Active:      true,
					Username:    gofakeit.Username(),
					Name:        gofakeit.Name(),
					Gender:      enumutils.GenderMale,
					DateOfBirth: &date,
					UserType:    enums.ClientUser,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.args.user.UserID == nil {
				t.Errorf("expected user ID not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_CreateClient(t *testing.T) {
	enrollment := time.Now()
	var clientTypeList pq.StringArray
	for _, c := range enums.AllClientType {
		clientTypeList = append(clientTypeList, string(c))
	}

	type args struct {
		ctx          context.Context
		client       *gorm.Client
		contactID    string
		identifierID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create client",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				client: &gorm.Client{
					Active:                  true,
					UserID:                  &userIDtoAssignClient,
					FacilityID:              facilityID,
					ClientCounselled:        true,
					ClientTypes:             clientTypeList,
					TreatmentEnrollmentDate: &enrollment,
					ProgramID:               programID,
				},
				contactID:    contactID,
				identifierID: identifierID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateClient(tt.args.ctx, tt.args.client, tt.args.contactID, tt.args.identifierID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateClient() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.args.client.ID == nil {
				t.Errorf("expected client ID not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_CreateMetric(t *testing.T) {
	inv := "invalid-id"

	type args struct {
		ctx    context.Context
		metric *gorm.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a metric",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				metric: &gorm.Metric{
					Active:    true,
					UserID:    &userID,
					Timestamp: time.Now(),
					Type:      enums.MetricTypeContent,
					Payload:   `{"contentID":"","duration":32}`,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid metric data",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				metric: &gorm.Metric{
					Active:    true,
					UserID:    &inv,
					Timestamp: time.Now(),
					Type:      enums.MetricTypeContent,
					Payload:   `{"contentID":"","duration":32}`,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateMetric(tt.args.ctx, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateIdentifier(t *testing.T) {
	type args struct {
		ctx        context.Context
		identifier *gorm.Identifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create identifier",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				identifier: &gorm.Identifier{
					Active:              true,
					IdentifierType:      "CCC",
					IdentifierValue:     "5678901234789",
					IdentifierUse:       "OFFICIAL",
					Description:         "CCC Number, Primary Identifier",
					IsPrimaryIdentifier: true,
					ProgramID:           programID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateIdentifier(tt.args.ctx, tt.args.identifier); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.args.identifier.ID == "" {
				t.Errorf("expected identifier ID not to be empty for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_CreateNotification(t *testing.T) {
	type args struct {
		ctx          context.Context
		notification *gorm.Notification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create appointment",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				notification: &gorm.Notification{
					Active:     true,
					Title:      "New Teleconsult",
					Body:       "Teleconsult with Doctor Who at the Tardis",
					Type:       "TELECONSULT",
					Flavour:    feedlib.FlavourConsumer,
					IsRead:     false,
					UserID:     &userID,
					FacilityID: &facilityID,
					ProgramID:  programID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateNotification(tt.args.ctx, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateUserSurvey(t *testing.T) {
	type args struct {
		ctx         context.Context
		userSurveys []*gorm.UserSurvey
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create user survey",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				userSurveys: []*gorm.UserSurvey{
					{
						UserID:      userID,
						Title:       gofakeit.Name(),
						Description: gofakeit.Sentence(1),
						Link:        gofakeit.URL(),
						ProgramID:   programID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: empty slice of user surveys",
			args: args{
				ctx:         addOrganizationContext(context.Background()),
				userSurveys: []*gorm.UserSurvey{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: create user survey, invalid user ID",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				userSurveys: []*gorm.UserSurvey{
					{
						UserID:      "userID",
						Title:       gofakeit.Name(),
						Description: gofakeit.Sentence(1),
						Link:        gofakeit.URL(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.CreateUserSurveys(tt.args.ctx, tt.args.userSurveys); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateUserSurveys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_SaveFeedback(t *testing.T) {

	feedback := &gorm.Feedback{
		ID:                feedbackID,
		UserID:            userID,
		FeedbackType:      "GENERAL_TYPE",
		SatisfactionLevel: 5,
		ServiceName:       "TEST",
		Feedback:          "I am a test feedback",
		RequiresFollowUp:  true,
		PhoneNumber:       interserviceclient.TestUserPhoneNumber,
		OrganisationID:    orgID,
		ProgramID:         programID,
	}

	invalidFeedback := &gorm.Feedback{
		ID:                "invalidFeedbackID",
		UserID:            userID,
		FeedbackType:      "GENERAL_TYPE",
		SatisfactionLevel: 5,
		ServiceName:       "TEST",
		Feedback:          "I am a test feedback",
		RequiresFollowUp:  true,
		PhoneNumber:       interserviceclient.TestUserPhoneNumber,
		OrganisationID:    orgID,
		ProgramID:         programID,
	}

	type args struct {
		ctx      context.Context
		feedback *gorm.Feedback
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: save feedback",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				feedback: feedback,
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to save feedback",
			args: args{
				ctx:      addOrganizationContext(context.Background()),
				feedback: invalidFeedback,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.SaveFeedback(tt.args.ctx, tt.args.feedback); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SaveFeedback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_RegisterClient(t *testing.T) {
	contactID := "bdc22436-e314-43f2-bb39-ba1ab332f9b0"
	identifierID := "bcbdaf68-3d36-4365-b575-4182d6759ad9"
	clientID := "26b30a42-cbb8-4773-aedb-c539602d04fc"
	userID := userIDToRegisterClient
	currentTime := time.Now()
	FHIRPatientID := "26b30a43-cbb8-4773-aedb-c539602d04fc"
	HealthPatientID := "29b30a42-cbb8-4553-aedb-c539602d04fc"

	invalidID := "invalidID"
	contactData := &gorm.Contact{
		ContactID:      &contactID,
		ContactType:    "PHONE",
		ContactValue:   testPhone,
		Active:         true,
		OptedIn:        true,
		UserID:         &userID,
		Flavour:        testFlavour,
		OrganisationID: orgID,
	}
	identifierData := &gorm.Identifier{
		ID:                  identifierID,
		OrganisationID:      orgID,
		Active:              true,
		IdentifierType:      "CCC",
		IdentifierValue:     "123456789",
		IdentifierUse:       "OFFICIAL",
		Description:         "A CCC Number",
		ValidFrom:           time.Now(),
		ValidTo:             time.Now(),
		IsPrimaryIdentifier: true,
		ProgramID:           programID,
	}
	clientData := &gorm.Client{
		ID:                      &clientID,
		Active:                  true,
		ClientTypes:             []string{"PMTCT"},
		UserID:                  &userID,
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &FHIRPatientID,
		HealthRecordID:          &HealthPatientID,
		ClientCounselled:        true,
		OrganisationID:          orgID,
		FacilityID:              facilityID,
		ProgramID:               programID,
	}
	InvalidClientData := &gorm.Client{
		ID:                      &invalidID,
		Active:                  true,
		ClientTypes:             []string{"PMTCT"},
		UserID:                  &userID,
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &FHIRPatientID,
		HealthRecordID:          &HealthPatientID,
		ClientCounselled:        true,
		OrganisationID:          orgID,
		FacilityID:              facilityID,
	}
	userProfile := &gorm.User{
		UserID:                 &userIDToRegisterClient,
		Username:               gofakeit.Name(),
		UserType:               enums.HealthcareWorkerUser,
		Gender:                 enumutils.GenderMale,
		Active:                 true,
		Contacts:               gorm.Contact{},
		PushTokens:             []string{},
		LastSuccessfulLogin:    &currentTime,
		LastFailedLogin:        &currentTime,
		FailedLoginCount:       3,
		NextAllowedLogin:       &currentTime,
		TermsAccepted:          true,
		AcceptedTermsID:        &termsID,
		Flavour:                feedlib.FlavourPro,
		Avatar:                 "test",
		IsSuspended:            true,
		PinChangeRequired:      true,
		HasSetPin:              true,
		HasSetSecurityQuestion: true,
		IsPhoneVerified:        true,
		OrganisationID:         uuid.New().String(),
		IsSuperuser:            true,
		Name:                   gofakeit.BeerBlg(),
		DateOfBirth:            &currentTime,
	}
	type args struct {
		ctx        context.Context
		user       *gorm.User
		contact    *gorm.Contact
		identifier *gorm.Identifier
		client     *gorm.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: register client",
			args: args{
				ctx:        addOrganizationContext(context.Background()),
				user:       userProfile,
				contact:    contactData,
				identifier: identifierData,
				client:     clientData,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register client",
			args: args{
				ctx:        addOrganizationContext(context.Background()),
				contact:    contactData,
				identifier: identifierData,
				client:     InvalidClientData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.RegisterClient(tt.args.ctx, tt.args.user, tt.args.contact, tt.args.identifier, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestPGInstance_RegisterStaff(t *testing.T) {
	identifierID := "c40b09a8-44b1-409c-bc5b-7e7623fcd7d5"
	staffID := staffIDToRegister
	userStaff := userToRegisterStaff
	currentTime := time.Now()

	userProfile := &gorm.User{
		UserID:                 &userToRegisterStaff,
		Username:               gofakeit.Name(),
		UserType:               enums.HealthcareWorkerUser,
		Gender:                 enumutils.GenderMale,
		Active:                 true,
		Contacts:               gorm.Contact{},
		PushTokens:             []string{},
		LastSuccessfulLogin:    &currentTime,
		LastFailedLogin:        &currentTime,
		FailedLoginCount:       3,
		NextAllowedLogin:       &currentTime,
		TermsAccepted:          true,
		AcceptedTermsID:        &termsID,
		Flavour:                feedlib.FlavourPro,
		Avatar:                 "test",
		IsSuspended:            true,
		PinChangeRequired:      true,
		HasSetPin:              true,
		HasSetSecurityQuestion: true,
		IsPhoneVerified:        true,
		OrganisationID:         uuid.New().String(),
		IsSuperuser:            true,
		Name:                   gofakeit.BeerBlg(),
		DateOfBirth:            &currentTime,
	}

	contactData := &gorm.Contact{
		ContactID:      &contactIDToRegisterStaff,
		ContactType:    "PHONE",
		ContactValue:   "+123445679890",
		Active:         true,
		OptedIn:        true,
		UserID:         &userToRegisterStaff,
		Flavour:        feedlib.FlavourPro,
		OrganisationID: orgID,
	}
	identifierData := &gorm.Identifier{
		ID:                  identifierID,
		Active:              true,
		IdentifierType:      "NATIONAL_ID",
		IdentifierValue:     "95454545",
		IdentifierUse:       "OFFICIAL",
		Description:         "A national ID number",
		ValidFrom:           time.Now(),
		ValidTo:             time.Now(),
		IsPrimaryIdentifier: true,
		OrganisationID:      orgID,
		ProgramID:           programID,
	}
	staff := &gorm.StaffProfile{
		ID:                &staffID,
		UserID:            userStaff,
		Active:            true,
		StaffNumber:       "123445679890",
		DefaultFacilityID: facilityID,
		OrganisationID:    orgID,
		ProgramID:         programID,
	}

	invalidStaff := &gorm.StaffProfile{
		ID:                &staffID,
		UserID:            "userStaff",
		Active:            true,
		StaffNumber:       "123445679890",
		DefaultFacilityID: facilityID,
		OrganisationID:    orgID,
	}
	type args struct {
		ctx          context.Context
		usr          *gorm.User
		contact      *gorm.Contact
		identifier   *gorm.Identifier
		staffProfile *gorm.StaffProfile
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: register staff",
			args: args{
				ctx:          addOrganizationContext(context.Background()),
				usr:          userProfile,
				contact:      contactData,
				identifier:   identifierData,
				staffProfile: staff,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register staff",
			args: args{
				ctx:          addOrganizationContext(context.Background()),
				contact:      contactData,
				identifier:   identifierData,
				staffProfile: invalidStaff,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.RegisterStaff(tt.args.ctx, tt.args.usr, tt.args.contact, tt.args.identifier, tt.args.staffProfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateQuestionnaire(t *testing.T) {
	name := gofakeit.BeerIbu()
	type args struct {
		ctx   context.Context
		input *gorm.Questionnaire
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create questionnaire",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.Questionnaire{
					ID:          uuid.NewString(),
					Active:      true,
					Name:        name,
					Description: gofakeit.Sentence(1),
					ProgramID:   programID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: create questionnaire, duplicate name",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.Questionnaire{
					ID:          uuid.NewString(),
					Active:      true,
					Name:        name,
					Description: gofakeit.Sentence(1),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateQuestionnaire(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateScreeningTool(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *gorm.ScreeningTool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create screening tool",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.ScreeningTool{
					ID:              uuid.NewString(),
					Active:          true,
					QuestionnaireID: questionnaireID,
					Threshold:       3,
					ClientTypes:     []string{string(enums.ClientTypeOtz)},
					Genders:         []string{enumutils.GenderFemale.String()},
					MinimumAge:      14,
					MaximumAge:      25,
					ProgramID:       programID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: create screening tool, questionnaire does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.ScreeningTool{
					ID:              uuid.NewString(),
					Active:          true,
					QuestionnaireID: uuid.NewString(),
					Threshold:       3,
					ClientTypes:     []string{string(enums.ClientTypeOtz)},
					Genders:         []string{enumutils.GenderFemale.String()},
					MinimumAge:      14,
					MaximumAge:      25,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateScreeningTool(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateScreeningTool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateQuestion(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *gorm.Question
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create question",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.Question{
					ID:                uuid.NewString(),
					Active:            true,
					QuestionnaireID:   questionnaireID,
					Text:              gofakeit.Sentence(1),
					QuestionType:      string(enums.QuestionTypeCloseEnded),
					ResponseValueType: string(enums.QuestionResponseValueTypeNumber),
					SelectMultiple:    false,
					Required:          true,
					Sequence:          1,
					ProgramID:         programID,
				},
			},
		},
		{
			name: "Sad case: create question, questionnaire does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.Question{
					ID:                uuid.NewString(),
					Active:            true,
					QuestionnaireID:   uuid.NewString(),
					Text:              gofakeit.Sentence(1),
					QuestionType:      string(enums.QuestionTypeCloseEnded),
					ResponseValueType: string(enums.QuestionResponseValueTypeNumber),
					SelectMultiple:    false,
					Required:          true,
					Sequence:          1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateQuestion(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateQuestion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateQuestionChoice(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *gorm.QuestionInputChoice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create question choice",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.QuestionInputChoice{
					ID:         uuid.NewString(),
					Active:     true,
					QuestionID: questionID,
					Choice:     gofakeit.Sentence(1),
					Value:      "1",
					Score:      1,
					ProgramID:  programID,
				},
			},
		},
		{
			name: "Sad case: create question choice, question does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				input: &gorm.QuestionInputChoice{
					ID:         uuid.NewString(),
					Active:     true,
					QuestionID: uuid.NewString(),
					Choice:     gofakeit.Sentence(1),
					Value:      "1",
					Score:      1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateQuestionChoice(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateQuestionChoice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateScreeningToolResponse(t *testing.T) {
	screeningToolsResponseID := uuid.NewString()

	type args struct {
		ctx                            context.Context
		screeningToolResponse          *gorm.ScreeningToolResponse
		screeningToolQuestionResponses []*gorm.ScreeningToolQuestionResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create screening tool response",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				screeningToolResponse: &gorm.ScreeningToolResponse{
					ID:              screeningToolsResponseID,
					Active:          true,
					ScreeningToolID: screeningToolID,
					FacilityID:      facilityID,
					ClientID:        clientID,
					AggregateScore:  1,
					ProgramID:       programID,
				},
				screeningToolQuestionResponses: []*gorm.ScreeningToolQuestionResponse{
					{
						ID:                      uuid.NewString(),
						Active:                  true,
						ScreeningToolResponseID: screeningToolsResponseID,
						QuestionID:              questionID,
						Response:                "0",
						Score:                   1,
						ProgramID:               programID,
						FacilityID:              facilityID,
					},
				},
			},
		},
		{
			name: "Sad case: create screening tool response, screening tool does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				screeningToolResponse: &gorm.ScreeningToolResponse{
					ID:              screeningToolsResponseID,
					Active:          true,
					ScreeningToolID: uuid.NewString(),
					FacilityID:      facilityID,
					ClientID:        clientID,
					AggregateScore:  1,
				},
				screeningToolQuestionResponses: []*gorm.ScreeningToolQuestionResponse{
					{
						ID:                      uuid.NewString(),
						Active:                  true,
						ScreeningToolResponseID: screeningToolsResponseID,
						QuestionID:              questionID,
						Response:                "0",
						Score:                   1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: screening tool response, question does not exist",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				screeningToolResponse: &gorm.ScreeningToolResponse{
					ID:              screeningToolsResponseID,
					Active:          true,
					ScreeningToolID: screeningToolID,
					FacilityID:      facilityID,
					ClientID:        clientID,
					AggregateScore:  1,
				},
				screeningToolQuestionResponses: []*gorm.ScreeningToolQuestionResponse{
					{
						ID:                      uuid.NewString(),
						Active:                  true,
						ScreeningToolResponseID: screeningToolsResponseID,
						QuestionID:              uuid.NewString(),
						Response:                "0",
						Score:                   1,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.CreateScreeningToolResponse(tt.args.ctx, tt.args.screeningToolResponse, tt.args.screeningToolQuestionResponses)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestPGInstance_RegisterCaregiver(t *testing.T) {
	dob := time.Now()
	type args struct {
		ctx       context.Context
		user      *gorm.User
		contact   *gorm.Contact
		caregiver *gorm.Caregiver
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register caregiver",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				user: &gorm.User{
					Username:    gofakeit.Username(),
					Name:        gofakeit.Name(),
					Gender:      enumutils.GenderMale,
					DateOfBirth: &dob,
					UserType:    enums.CaregiverUser,
					Flavour:     feedlib.FlavourConsumer,
					Active:      true,
				},
				contact: &gorm.Contact{
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
					Active:       true,
					OptedIn:      false,
					Flavour:      feedlib.FlavourConsumer,
				},
				caregiver: &gorm.Caregiver{
					CaregiverNumber: gofakeit.SSN(),
					Active:          true,
					ProgramID:       programID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.RegisterCaregiver(tt.args.ctx, tt.args.user, tt.args.contact, tt.args.caregiver); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_AddCaregiverToClient(t *testing.T) {

	now := time.Now()

	type args struct {
		ctx             context.Context
		clientCaregiver *gorm.CaregiverClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add new caregiver to client",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				clientCaregiver: &gorm.CaregiverClient{
					CaregiverID:        "8ecbbc80-24c8-421a-9f1a-e14e12678ef2",
					ClientID:           clientID2,
					Active:             true,
					RelationshipType:   enums.CaregiverTypeFather,
					CaregiverConsent:   enums.ConsentStateAccepted,
					CaregiverConsentAt: &now,
					ClientConsent:      enums.ConsentStateAccepted,
					ClientConsentAt:    &now,
					OrganisationID:     orgID,
					AssignedBy:         staffID,
					ProgramID:          programID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to add new caregiver to client",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				clientCaregiver: &gorm.CaregiverClient{
					CaregiverID:        testCaregiverID,
					ClientID:           "clientID",
					Active:             true,
					RelationshipType:   enums.CaregiverTypeFather,
					CaregiverConsent:   enums.ConsentStateAccepted,
					CaregiverConsentAt: &now,
					ClientConsent:      enums.ConsentStateAccepted,
					ClientConsentAt:    &now,
					OrganisationID:     orgID,
					AssignedBy:         staffID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.AddCaregiverToClient(tt.args.ctx, tt.args.clientCaregiver); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AddCaregiverToClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateCaregiver(t *testing.T) {

	type args struct {
		ctx       context.Context
		caregiver *gorm.Caregiver
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a caregiver",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				caregiver: &gorm.Caregiver{
					Active:          true,
					CaregiverNumber: gofakeit.SSN(),
					UserID:          userIDtoAssignClient,
					ProgramID:       programID,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid user id",
			args: args{
				ctx: addOrganizationContext(context.Background()),
				caregiver: &gorm.Caregiver{
					Active:          true,
					CaregiverNumber: gofakeit.SSN(),
					UserID:          "invalid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateCaregiver(tt.args.ctx, tt.args.caregiver); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateOrganisation(t *testing.T) {
	invalidUUID := "1"
	type args struct {
		ctx          context.Context
		organization *gorm.Organisation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create an organisation",
			args: args{
				ctx: context.Background(),
				organization: &gorm.Organisation{
					ID:               &orgID,
					Active:           true,
					OrganisationCode: gofakeit.SSN(),
					Name:             gofakeit.SSN(),
					Description:      gofakeit.Sentence(10),
					EmailAddress:     gofakeit.Email(),
					PhoneNumber:      gofakeit.Phone(),
					PostalAddress:    gofakeit.Address().Address,
					PhysicalAddress:  gofakeit.Address().Address,
					DefaultCountry:   gofakeit.Country(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create an organisation with invalid org code",
			args: args{
				ctx: context.Background(),
				organization: &gorm.Organisation{
					ID:               &invalidUUID,
					Active:           true,
					OrganisationCode: uuid.New().String(),
					Name:             "test",
					Description:      gofakeit.Sentence(10),
					EmailAddress:     gofakeit.Email(),
					PhoneNumber:      gofakeit.Phone(),
					PostalAddress:    gofakeit.HipsterParagraph(3, 200, 200, " "),
					PhysicalAddress:  gofakeit.HipsterParagraph(3, 200, 200, " "),
					DefaultCountry:   gofakeit.Country(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateOrganisation(tt.args.ctx, tt.args.organization); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
