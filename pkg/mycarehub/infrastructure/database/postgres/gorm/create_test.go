package gorm_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
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

	facility := &gorm.Facility{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
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

	serviceRequestInput := &gorm.ClientServiceRequest{
		Active:         false,
		RequestType:    "HealthDiary",
		Request:        gofakeit.Sentence(5),
		Status:         "PENDING",
		InProgressAt:   time.Now(),
		ResolvedAt:     time.Now(),
		ClientID:       clientID,
		OrganisationID: orgID,
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
					EntryType:                "Test",
					ShareWithHealthWorker:    true,
					SharedAt:                 time.Now(),
					ClientID:                 clientID,
					OrganisationID:           uuid.New().String(),
				},
			},
			wantErr: false,
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
