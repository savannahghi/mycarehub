package gorm_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := gofakeit.Name()
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()

	facility := &gorm.Facility{
		Name:               name,
		Code:               code,
		Active:             true,
		County:             county,
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
	}

	invalidFacility := &gorm.Facility{
		Name:        name,
		Code:        -458789,
		Active:      true,
		County:      county,
		Description: description,
	}

	type args struct {
		ctx      context.Context
		facility *gorm.Facility
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get or create facility",
			args: args{
				ctx:      ctx,
				facility: facility,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail tp get or create facility",
			args: args{
				ctx:      ctx,
				facility: nil,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get or create facility",
			args: args{
				ctx: ctx,
				facility: &gorm.Facility{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      gofakeit.HipsterSentence(50),
					Description: description,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create an invalid facility",
			args: args{
				ctx:      ctx,
				facility: invalidFacility,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
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
	// teardown
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}
	if err = pg.DB.Where("id", facility.FacilityID).Unscoped().Delete(&gorm.Facility{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_SaveTemporaryUserPin(t *testing.T) {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	ctx := context.Background()

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
				ctx:        ctx,
				pinPayload: pinPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing payload",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "invalid: invalid payload",
			args: args{
				ctx: ctx,
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
				ctx:        ctx,
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
	ctx := context.Background()

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
				ctx:     ctx,
				pinData: pinPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing user id",
			args: args{
				ctx: ctx,
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
				ctx: ctx,
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
				ctx: ctx,
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

	ctx := context.Background()

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
				ctx: ctx,
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
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	generatedAt := time.Now()
	validUntil := time.Now().AddDate(0, 0, 2)

	ext := extension.NewExternalMethodsImpl()

	otp, err := ext.GenerateOTP(ctx)
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

	newOTP, err := ext.GenerateOTP(ctx)
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
				ctx:      ctx,
				otpInput: gormOTPInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				otpInput: invalidgormOTPInput1,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:      ctx,
				otpInput: invalidgormOTPInput2,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavour and phone",
			args: args{
				ctx:      ctx,
				otpInput: invalidgormOTPInput3,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx: ctx,
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
				ctx: ctx,
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
	ctx := context.Background()
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
				ctx:                 ctx,
				serviceRequestInput: serviceRequestInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid meta data",
			args: args{
				ctx:                 ctx,
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

func TestPGInstance_CreateClientCaregiver(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx             context.Context
		clientID        string
		clientCaregiver *gorm.Caregiver
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      ctx,
				clientID: ClientToAddCaregiver,
				clientCaregiver: &gorm.Caregiver{
					FirstName:     gofakeit.Name(),
					LastName:      gofakeit.Name(),
					PhoneNumber:   testPhone,
					CaregiverType: enums.CaregiverTypeFather,
					Active:        true,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx:      ctx,
				clientID: ClientToAddCaregiver,
				clientCaregiver: &gorm.Caregiver{
					FirstName:     gofakeit.Name(),
					LastName:      gofakeit.Name(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverType(gofakeit.HipsterSentence(20)),
					Active:        true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.CreateClientCaregiver(tt.args.ctx, tt.args.clientID, tt.args.clientCaregiver); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()

	clientHealthDiaryEntryID := uuid.New().String()

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
				ctx: ctx,
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientHealthDiaryEntryID,
					Active:                   true,
					Mood:                     "Very Cool",
					Note:                     "I'm happy",
					EntryType:                "HOME_PAGE_HEALTH_DIARY_ENTRY",
					ShareWithHealthWorker:    true,
					SharedAt:                 time.Now(),
					ClientID:                 clientID,
					OrganisationID:           uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: invalid input",
			args: args{
				ctx: ctx,
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					Active:                true,
					Mood:                  gofakeit.HipsterSentence(20),
					Note:                  "test",
					EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY",
					ShareWithHealthWorker: false,
					SharedAt:              time.Now(),
					ClientID:              clientID,
					OrganisationID:        orgID,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - no client ID",
			args: args{
				ctx: ctx,
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientHealthDiaryEntryID,
					Active:                   true,
					Mood:                     "Very Cool",
					Note:                     "I'm happy",
					EntryType:                "Test",
					ShareWithHealthWorker:    true,
					SharedAt:                 time.Now(),
					OrganisationID:           uuid.New().String(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - no health diary input",
			args: args{
				ctx:              ctx,
				healthDiaryInput: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.CreateHealthDiaryEntry(tt.args.ctx, tt.args.healthDiaryInput); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CreateHealthDiaryEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_CreateCommunity(t *testing.T) {
	ctx := context.Background()

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
				ctx: ctx,
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
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
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
				ctx: context.Background(),
				person: &gorm.RelatedPerson{
					Active:           true,
					FirstName:        gofakeit.Name(),
					LastName:         gofakeit.Name(),
					Gender:           "MALE",
					RelationshipType: "Next of Kin",
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
				ctx: context.Background(),
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
				ctx: context.Background(),
				screeningToolResponses: []*gorm.ScreeningToolsResponse{
					{
						QuestionID: screeningToolsQuestionID,
						ClientID:   clientID,
						Response:   "0",
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
				ctx: context.Background(),
				appointment: &gorm.Appointment{
					Active:                    true,
					ExternalID:                strconv.Itoa(gofakeit.Number(0, 1000)),
					ClientID:                  clientID,
					FacilityID:                facilityID,
					Reason:                    "Dental",
					Date:                      time.Now().Add(time.Duration(100)),
					HasRescheduledAppointment: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create an appointment",
			args: args{
				ctx: context.Background(),
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
	ctx := context.Background()

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
				ctx: ctx,
				serviceRequestInput: &gorm.StaffServiceRequest{
					ID:                &staffServiceRequestID,
					Active:            true,
					RequestType:       gofakeit.BeerName(),
					Request:           gofakeit.BeerName(),
					Status:            gofakeit.BeerName(),
					ResolvedAt:        &currentTime,
					StaffID:           staffID,
					OrganisationID:    orgID,
					DefaultFacilityID: &facilityID,
					Meta:              meta,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
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
				ctx: ctx,
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
