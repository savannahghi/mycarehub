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
	"github.com/savannahghi/interserviceclient"
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
				ctx: ctx,
				healthDiaryInput: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientHealthDiaryEntryID,
					Active:                   true,
					Mood:                     "Very Cool",
					Note:                     "I'm happy",
					EntryType:                "HOME_PAGE_HEALTH_DIARY_ENTRY",
					ShareWithHealthWorker:    true,
					SharedAt:                 &currentTime,
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
				ctx: ctx,
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
				ctx: context.Background(),
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
				ctx: context.Background(),
				client: &gorm.Client{
					Active:                  true,
					UserID:                  &userIDtoAssignClient,
					FacilityID:              facilityID,
					ClientCounselled:        true,
					ClientTypes:             clientTypeList,
					TreatmentEnrollmentDate: &enrollment,
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
				identifier: &gorm.Identifier{
					Active:              true,
					IdentifierType:      "CCC",
					IdentifierValue:     "5678901234789",
					IdentifierUse:       "OFFICIAL",
					Description:         "CCC Number, Primary Identifier",
					IsPrimaryIdentifier: true,
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
				ctx: context.Background(),
				notification: &gorm.Notification{
					Active:     true,
					Title:      "New Teleconsult",
					Body:       "Teleconsult with Doctor Who at the Tardis",
					Type:       "TELECONSULT",
					Flavour:    feedlib.FlavourConsumer,
					IsRead:     false,
					UserID:     &userID,
					FacilityID: &facilityID,
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
				ctx: context.Background(),
				userSurveys: []*gorm.UserSurvey{
					{
						UserID:      userID,
						Title:       gofakeit.Name(),
						Description: gofakeit.Sentence(1),
						Link:        gofakeit.URL(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: empty slice of user surveys",
			args: args{
				ctx:         context.Background(),
				userSurveys: []*gorm.UserSurvey{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: create user survey, invalid user ID",
			args: args{
				ctx: context.Background(),
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
	ctx := context.Background()

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
				ctx:      ctx,
				feedback: feedback,
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to save feedback",
			args: args{
				ctx:      ctx,
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
	chvID := userIDToRegisterClient
	caregiverID := "28b20a42-cbb8-4553-aedb-c575602d04fc"

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
	}
	clientData := &gorm.Client{
		ID:                      &clientID,
		Active:                  true,
		ClientTypes:             []string{"PMTCT"},
		UserID:                  &userID,
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &FHIRPatientID,
		HealthRecordID:          &HealthPatientID,
		TreatmentBuddy:          uuid.New().String(),
		ClientCounselled:        true,
		OrganisationID:          orgID,
		FacilityID:              facilityID,
		CHVUserID:               &chvID,
		CaregiverID:             &caregiverID,
	}
	InvalidClientData := &gorm.Client{
		ID:                      &invalidID,
		Active:                  true,
		ClientTypes:             []string{"PMTCT"},
		UserID:                  &userID,
		TreatmentEnrollmentDate: &currentTime,
		FHIRPatientID:           &FHIRPatientID,
		HealthRecordID:          &HealthPatientID,
		TreatmentBuddy:          uuid.New().String(),
		ClientCounselled:        true,
		OrganisationID:          orgID,
		FacilityID:              facilityID,
		CHVUserID:               &chvID,
		CaregiverID:             &caregiverID,
	}
	type args struct {
		ctx        context.Context
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
				ctx:        context.Background(),
				contact:    contactData,
				identifier: identifierData,
				client:     clientData,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register client",
			args: args{
				ctx:        context.Background(),
				contact:    contactData,
				identifier: identifierData,
				client:     InvalidClientData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.RegisterClient(tt.args.ctx, tt.args.contact, tt.args.identifier, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_RegisterStaff(t *testing.T) {
	identifierID := "c40b09a8-44b1-409c-bc5b-7e7623fcd7d5"
	staffID := staffIDToRegister
	userStaff := userToRegisterStaff
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
	}
	staff := &gorm.StaffProfile{
		ID:                &staffID,
		UserID:            userStaff,
		Active:            true,
		StaffNumber:       "123445679890",
		DefaultFacilityID: facilityID,
		OrganisationID:    orgID,
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
				ctx:          context.Background(),
				contact:      contactData,
				identifier:   identifierData,
				staffProfile: staff,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register staff",
			args: args{
				ctx:          context.Background(),
				contact:      contactData,
				identifier:   identifierData,
				staffProfile: invalidStaff,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.RegisterStaff(tt.args.ctx, tt.args.contact, tt.args.identifier, tt.args.staffProfile); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
