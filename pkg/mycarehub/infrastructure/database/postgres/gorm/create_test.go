package gorm_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	ID := ksuid.New().String()
	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := gofakeit.Name()
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
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
}

func TestPGInstance_SaveTemporaryUserPin(t *testing.T) {

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	ctx := context.Background()

	flavor := feedlib.FlavourConsumer
	userID := uuid.New().String()

	newExtension := extension.NewExternalMethodsImpl()

	tempPin, err := newExtension.GenerateTempPIN(ctx)
	if err != nil {
		t.Errorf("failed to generate temporary pin: %v", err)
	}
	salt, encryptedTempPin := newExtension.EncryptPIN(tempPin, nil)

	userInput := &gorm.User{
		UserID:        &userID,
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		MiddleName:    gofakeit.FirstName(),
		UserType:      enums.ClientUser,
		Gender:        enumutils.GenderMale,
		TermsAccepted: true,
		Flavour:       flavor,
		IsSuspended:   false,
	}
	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	pinPayload := &gorm.PINData{
		UserID:    *userInput.UserID,
		HashedPIN: encryptedTempPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   flavor,
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

	if err = pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.PINData{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}

	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
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

	newExtension := extension.NewExternalMethodsImpl()
	salt, encryptedPin := newExtension.EncryptPIN("0000", nil)
	flavour := feedlib.FlavourConsumer
	userID := uuid.New().String()

	// Setup test user
	userInput := &gorm.User{
		UserID:          &userID,
		FirstName:       gofakeit.FirstName(),
		LastName:        gofakeit.LastName(),
		MiddleName:      gofakeit.FirstName(),
		UserType:        enums.ClientUser,
		Gender:          enumutils.GenderMale,
		Flavour:         flavour,
		AcceptedTermsID: &termsID,
		TermsAccepted:   true,
		IsSuspended:     false,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	pinPayload := &gorm.PINData{
		UserID:    *userInput.UserID,
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
	if err := pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.PINData{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_SaveSecurityQuestionResponse(t *testing.T) {

	ctx := context.Background()

	userID := uuid.New().String()

	flavour := feedlib.FlavourConsumer

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		UserID:          &userID,
		FirstName:       gofakeit.FirstName(),
		LastName:        gofakeit.LastName(),
		MiddleName:      gofakeit.FirstName(),
		UserType:        enums.ClientUser,
		Gender:          enumutils.GenderMale,
		Flavour:         flavour,
		AcceptedTermsID: &termsID,
		TermsAccepted:   true,
		IsSuspended:     false,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	sequence := 1

	securityQuestionInput := &gorm.SecurityQuestion{
		QuestionStem: gofakeit.Sentence(3),
		Description:  gofakeit.Sentence(3),
		ResponseType: enums.SecurityQuestionResponseTypeNumber,
		Flavour:      flavour,
		Sequence:     &sequence,
	}

	err = pg.DB.Create(securityQuestionInput).Error
	if err != nil {
		t.Errorf("Create securityQuestion failed: %v", err)
	}

	type args struct {
		ctx                      context.Context
		securityQuestionResponse *gorm.SecurityQuestionResponse
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
				securityQuestionResponse: &gorm.SecurityQuestionResponse{
					QuestionID: *securityQuestionInput.SecurityQuestionID,
					UserID:     *userInput.UserID,
					Response:   "20",
				},
			},
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
	if err := pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.SecurityQuestionResponse{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err := pg.DB.Where("id", securityQuestionInput.SecurityQuestionID).Unscoped().Delete(&gorm.SecurityQuestion{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
