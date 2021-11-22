package otp_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/testutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
)

var (
	termsID   = 50005
	orgID     = uuid.New().String()
	testingDB *gorm.PGInstance
)

func TestMain(m *testing.M) {
	log.Println("setting up test database")
	var err error
	testingDB, err = gorm.NewPGInstance()
	if err != nil {
		os.Exit(1)
	}
	// add organization
	createOrganization()

	//create terms
	createTermsOfService()

	log.Printf("Running tests ...")
	os.Exit(m.Run())

	// teardown
	// remove organization
	log.Println("tearing down test database")

	testingDB.DB.Unscoped().Delete(gorm.Organisation{OrganisationID: &orgID})
	testingDB.DB.Unscoped().Delete(gorm.TermsOfService{TermsID: &termsID})
}

func createOrganization() {
	organisation := &gorm.Organisation{
		OrganisationID:   &orgID,
		Active:           true,
		OrgCode:          gofakeit.Name(),
		Code:             gofakeit.Number(100, 344),
		OrganisationName: gofakeit.Name(),
		EmailAddress:     gofakeit.Email(),
		PhoneNumber:      gofakeit.Phone(),
		PostalAddress:    gofakeit.Address().Address,
		PhysicalAddress:  gofakeit.Address().City,
		DefaultCountry:   "KEN",
	}

	testingDB.DB.Create(organisation)
}

func createTermsOfService() {
	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	terms := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
	}

	testingDB.DB.Create(terms)
}

func TestUseCaseOTPImpl_VerifyPhoneNumber_Integration(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	inexistentNo := "+254700000520"
	invalidPhone := ksuid.New().String()
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contactID := uuid.New().String()
	contact := &gorm.Contact{
		ContactID:      &contactID,
		ContactType:    "SMS",
		ContactValue:   "+254710000100",
		Active:         true,
		OptedIn:        true,
		UserID:         userInput.UserID,
		Flavour:        userInput.Flavour,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}

	err = pg.DB.Create(&contact).Error
	if err != nil {
		t.Errorf("failed to create contact: %v", err)
	}

	type args struct {
		ctx     context.Context
		phone   string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.OtpResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				phone:   contact.ContactValue,
				flavour: contact.Flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   contact.ContactValue,
				flavour: "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - inexistent number and good flavour",
			args: args{
				ctx:     ctx,
				phone:   inexistentNo,
				flavour: contact.Flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - inexistent number and bad flavour",
			args: args{
				ctx:     ctx,
				phone:   inexistentNo,
				flavour: "gofakeit.HipsterSentence(100)",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and valid flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: contact.Flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and empty flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.OTP.VerifyPhoneNumber(tt.args.ctx, tt.args.phone, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", contact.ContactID).Unscoped().Delete(&gorm.Contact{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestUseCaseOTPImpl_VerifyOTP_integration_test(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	userID := uuid.New().String()
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		UserID:              &userID,
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	otpID := gofakeit.Number(1, 10000)
	generatedAt := time.Now()
	validUntil := time.Now().AddDate(0, 0, 2)

	otp, err := i.OTP.GenerateOTP(ctx)
	if err != nil {
		t.Errorf("unable to generate OTP")
	}

	otpInput := &gorm.UserOTP{
		OTPID:       otpID,
		UserID:      *userInput.UserID,
		Valid:       true,
		GeneratedAt: generatedAt,
		ValidUntil:  validUntil,
		Channel:     "SMS",
		Flavour:     userInput.Flavour,
		PhoneNumber: "+254710000111",
		OTP:         otp,
	}

	err = pg.DB.Create(&otpInput).Error
	if err != nil {
		t.Errorf("failed to create otp: %v", err)
	}

	validOTPPayload := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         otpInput.OTP,
		Flavour:     flavour,
	}
	invalidOTPPayload2 := &dto.VerifyOTPInput{
		PhoneNumber: "",
		OTP:         otpInput.OTP,
		Flavour:     flavour,
	}
	invalidOTPPayload3 := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         "",
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		PhoneNumber: otpInput.PhoneNumber,
		OTP:         otpInput.OTP,
		Flavour:     "flavour",
	}
	invalidOTPPayload5 := &dto.VerifyOTPInput{
		PhoneNumber: "otpInput.PhoneNumber",
		OTP:         "otpInput.OTP",
		Flavour:     "flavour",
	}
	invalidOTPPayload6 := &dto.VerifyOTPInput{
		PhoneNumber: gofakeit.HipsterParagraph(1, 10, 100, ""),
		OTP:         gofakeit.HipsterParagraph(1, 10, 100, ""),
		Flavour:     "gofakeit.HipsterParagraph(300, 10, 100)",
	}

	type args struct {
		ctx     context.Context
		payload *dto.VerifyOTPInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				payload: validOTPPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload2,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no otp",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload3,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload4,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload5,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - very bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload6,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.OTP.VerifyOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.VerifyOTPInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseOTPImpl.VerifyOTPInput() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", otpInput.OTPID).Unscoped().Delete(&gorm.UserOTP{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestUseCaseOTPImpl_GenerateAndSendOTP_Integration_test(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contact := &gorm.Contact{
		ContactType:    "SMS",
		ContactValue:   "+254711111111",
		Active:         true,
		OptedIn:        true,
		UserID:         userInput.UserID,
		Flavour:        userInput.Flavour,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}

	err = pg.DB.Create(&contact).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: contact.ContactValue,
				flavour:     userInput.Flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				flavour:     userInput.Flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				flavour:     feedlib.Flavour("invalid-flavour"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := i.OTP.GenerateAndSendOTP(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.GenerateAndSendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", contact.ContactID).Unscoped().Delete(&gorm.Contact{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
