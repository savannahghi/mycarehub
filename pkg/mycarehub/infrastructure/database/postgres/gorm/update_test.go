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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
)

func createInactiveTestFacility() *gorm.Facility {
	ID := uuid.New().String()
	name := ksuid.New().String()
	code := rand.Intn(1000000)
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      false,
		County:      county,
		Description: description,
	}

	return facility
}

func TestPGInstance_InactivateFacility(t *testing.T) {

	ctx := context.Background()

	testFacility := createTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &facility.Code,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ReactivateFacility(t *testing.T) {

	ctx := context.Background()

	testFacility := createInactiveTestFacility()

	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &facility.Code,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_SetNickname(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	nickname := uuid.New().String()
	invalidUserID := ksuid.New().String()
	invalidNickname := gofakeit.HipsterSentence(50)
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	user, err := testingDB.GetUserProfileByUserID(ctx, *userInput.UserID)
	if err != nil {
		t.Errorf("failed to get user profile with the provided user ID: %v", err)
	}

	type args struct {
		ctx      context.Context
		userID   *string
		nickname *string
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
				ctx:      ctx,
				userID:   user.UserID,
				nickname: &nickname,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				userID:   &invalidUserID,
				nickname: &nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:      ctx,
				userID:   nil,
				nickname: &nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no nickname",
			args: args{
				ctx:      ctx,
				userID:   userInput.UserID,
				nickname: &invalidNickname,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SetNickName(tt.args.ctx, tt.args.userID, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SetNickName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SetNickName() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_InvalidatePIN(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	nickname := uuid.New().String()
	currentTime := time.Now()
	futureTime := time.Now().AddDate(0, 0, 1)
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)
	newExtension := extension.NewExternalMethodsImpl()
	salt, encryptedPin := newExtension.EncryptPIN("0000", nil)

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	pinInput := &gorm.PINData{
		UserID:    *userInput.UserID,
		HashedPIN: encryptedPin,
		ValidFrom: time.Now(),
		ValidTo:   futureTime,
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
		Salt:      salt,
	}

	err = pg.DB.Create(pinInput).Error
	if err != nil {
		t.Errorf("failed to create pin: %v", err)
	}

	type args struct {
		ctx    context.Context
		userID string
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
				ctx:    ctx,
				userID: *userInput.UserID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: no user id provided",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.InvalidatePIN(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InvalidatePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InvalidatePIN() = %v, want %v", got, tt.want)
			}
		})
	}
	// TearDown
	if err := pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.PINData{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UpdateIsCorrectSecurityQuestionResponse(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	// Setup test user
	userInput := &gorm.User{
		Username:         uuid.New().String(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.FirstName(),
		UserType:         enums.ClientUser,
		Gender:           enumutils.GenderMale,
		Flavour:          feedlib.FlavourConsumer,
		AcceptedTermsID:  &termsID,
		TermsAccepted:    true,
		IsSuspended:      true,
		OrganisationID:   orgID,
		NextAllowedLogin: &pastTime,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	sequence := 1

	securityQuestionInput := &gorm.SecurityQuestion{
		QuestionStem:   gofakeit.Sentence(3),
		Description:    gofakeit.Sentence(3),
		ResponseType:   enums.SecurityQuestionResponseTypeNumber,
		Flavour:        feedlib.FlavourConsumer,
		Sequence:       &sequence,
		OrganisationID: orgID,
		Active:         true,
	}

	err = pg.DB.Create(securityQuestionInput).Error
	if err != nil {
		t.Errorf("Create securityQuestion failed: %v", err)
	}

	securityQuestionResponseInput := &gorm.SecurityQuestionResponse{
		UserID:         *userInput.UserID,
		QuestionID:     *securityQuestionInput.SecurityQuestionID,
		Response:       "23",
		Timestamp:      time.Now(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(securityQuestionResponseInput).Error
	if err != nil {
		t.Errorf("Create securityQuestionResponse failed: %v", err)
	}

	type args struct {
		ctx                               context.Context
		userID                            string
		isCorrectSecurityQuestionResponse bool
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
				ctx:                               ctx,
				userID:                            *userInput.UserID,
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid user id",
			args: args{
				ctx:                               ctx,
				userID:                            uuid.New().String(),
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateIsCorrectSecurityQuestionResponse(tt.args.ctx, tt.args.userID, tt.args.isCorrectSecurityQuestionResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() = %v, want %v", got, tt.want)
			}
		})
	}
	// TearDown
	if err = pg.DB.Where("id", securityQuestionResponseInput.ResponseID).Unscoped().Delete(&gorm.SecurityQuestionResponse{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err := pg.DB.Where("id", securityQuestionInput.SecurityQuestionID).Unscoped().Delete(&gorm.SecurityQuestion{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_AcceptTerms(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	// Setup test user
	userInput := &gorm.User{
		Username:         uuid.New().String(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.FirstName(),
		UserType:         enums.ClientUser,
		Gender:           enumutils.GenderMale,
		Flavour:          feedlib.FlavourConsumer,
		AcceptedTermsID:  &termsID,
		TermsAccepted:    true,
		IsSuspended:      true,
		OrganisationID:   orgID,
		NextAllowedLogin: &pastTime,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
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
				userID:  userInput.UserID,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing args",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
		if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
			t.Errorf("failed to delete record = %v", err)
		}
	}
	// Teardown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UpdateUserFailedLoginCount(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}
	// Setup test user
	userInput := &gorm.User{
		Username:         uuid.New().String(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.FirstName(),
		UserType:         enums.ClientUser,
		Gender:           enumutils.GenderMale,
		Flavour:          feedlib.FlavourConsumer,
		AcceptedTermsID:  &termsID,
		TermsAccepted:    true,
		IsSuspended:      true,
		OrganisationID:   orgID,
		NextAllowedLogin: &pastTime,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx                 context.Context
		userID              string
		failedLoginAttempts int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:                 ctx,
				userID:              *userInput.UserID,
				failedLoginAttempts: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.UpdateUserFailedLoginCount(tt.args.ctx, tt.args.userID, tt.args.failedLoginAttempts); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserFailedLoginCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Teardown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UpdateUserLastFailedLoginTime(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}
	// Setup test user
	userInput := &gorm.User{
		Username:         uuid.New().String(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.FirstName(),
		UserType:         enums.ClientUser,
		Gender:           enumutils.GenderMale,
		Flavour:          feedlib.FlavourConsumer,
		AcceptedTermsID:  &termsID,
		TermsAccepted:    true,
		IsSuspended:      true,
		OrganisationID:   orgID,
		NextAllowedLogin: &pastTime,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:    ctx,
				userID: *userInput.UserID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserLastFailedLoginTime(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserLastFailedLoginTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Teardown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UpdateUserNextAllowedLoginTime(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}
	// Setup test user
	userInput := &gorm.User{
		Username:         uuid.New().String(),
		FirstName:        gofakeit.FirstName(),
		LastName:         gofakeit.LastName(),
		MiddleName:       gofakeit.FirstName(),
		UserType:         enums.ClientUser,
		Gender:           enumutils.GenderMale,
		Flavour:          feedlib.FlavourConsumer,
		AcceptedTermsID:  &termsID,
		TermsAccepted:    true,
		IsSuspended:      true,
		OrganisationID:   orgID,
		NextAllowedLogin: &pastTime,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx                  context.Context
		userID               string
		nextAllowedLoginTime time.Time
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:                  ctx,
				userID:               *userInput.UserID,
				nextAllowedLoginTime: time.Now().Add(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserNextAllowedLoginTime(tt.args.ctx, tt.args.userID, tt.args.nextAllowedLoginTime); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserNextAllowedLoginTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Teardown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
func TestPGInstance_ShareContent(t *testing.T) {
	ctx := context.Background()
	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		// WagtailCorePageID:     gofakeit.Number(200, 500),
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1, // default to 1 => wagtailcore page
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1, // default to 1 => en
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           100,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	type args struct {
		ctx   context.Context
		input dto.ShareContentInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					UserID:    *userInput.UserID,
					ContentID: contentItemInput.PagePtrID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing input",
			args: args{
				ctx:   ctx,
				input: dto.ShareContentInput{},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ShareContent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ShareContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ShareContent() = %v, want %v", got, tt.want)
			}
		})
	}
	//TearDown
	if err = pg.DB.Where("content_item_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentShare{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_BookmarkContent(t *testing.T) {
	ctx := context.Background()
	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		// WagtailCorePageID:     gofakeit.Number(200, 500),
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1, // default to 1 => wagtailcore page
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1, // default to 1 => en
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           100,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentItemInput.PagePtrID,
			},
			wantErr: false,
			want:    true,
		},
		{
			// Ensures there is idepotency
			name: "bookmark already exists",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentItemInput.PagePtrID,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invald: missing parama",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.BookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.BookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.BookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("content_item_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentBookmark{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UnBookmarkContent(t *testing.T) {
	ctx := context.Background()
	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		// WagtailCorePageID:     gofakeit.Number(200, 500),
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1, // default to 1 => wagtailcore page
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1, // default to 1 => en
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           100,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentItemInput.PagePtrID,
			},
			wantErr: false,
			want:    true,
		},
		{
			// Ensures there is idepotency
			name: "bookmark already exists",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentItemInput.PagePtrID,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invald: missing parama",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UnBookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UnBookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UnBookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}
	//TearDown
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UpdateUserPinChangeRequiredStatus(t *testing.T) {
	ctx := context.Background()
	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)
	futureTime := time.Now().AddDate(0, 0, 2)
	newExtension := extension.NewExternalMethodsImpl()
	salt, encryptedPin := newExtension.EncryptPIN("0000", nil)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	pinInput := &gorm.PINData{
		UserID:    *userInput.UserID,
		HashedPIN: encryptedPin,
		ValidFrom: time.Now(),
		ValidTo:   futureTime,
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
		Salt:      salt,
	}

	err = pg.DB.Create(pinInput).Error
	if err != nil {
		t.Errorf("failed to create pin: %v", err)
	}

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name string

		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				userID:  *userInput.UserID,
				flavour: flavour,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateUserPinChangeRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() = %v, want %v", got, tt.want)
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

func TestPGInstance_LikeContent(t *testing.T) {
	ctx := context.Background()

	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1,
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1,
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           10,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	contentID := uuid.New().String()

	contentLike := &gorm.ContentLike{
		Base:           gorm.Base{},
		ContentLikeID:  contentID,
		Active:         true,
		ContentID:      contentItemInput.PagePtrID,
		UserID:         *userInput.UserID,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}
	err = pg.DB.Create(contentLike).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		context   context.Context
		userID    string
		contentID int
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
				context:   ctx,
				userID:    *userInput.UserID,
				contentID: contentLike.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: contentLike.ContentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				context:   ctx,
				userID:    *userInput.UserID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.LikeContent(tt.args.context, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.LikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("content_item_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentLike{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_UnlikeContent(t *testing.T) {
	ctx := context.Background()

	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1,
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1,
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           10,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	contentID := uuid.New().String()

	contentLike := &gorm.ContentLike{
		Base:           gorm.Base{},
		ContentLikeID:  contentID,
		Active:         true,
		ContentID:      contentItemInput.PagePtrID,
		UserID:         *userInput.UserID,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}
	err = pg.DB.Create(contentLike).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		context   context.Context
		userID    string
		contentID int
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
				context:   ctx,
				userID:    *userInput.UserID,
				contentID: contentLike.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: contentLike.ContentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				context:   ctx,
				userID:    *userInput.UserID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UnlikeContent(tt.args.context, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.LikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("content_item_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentLike{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_ViewContent(t *testing.T) {
	ctx := context.Background()
	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	// Setup test user
	userInput := &gorm.User{
		Username:            gofakeit.BeerHop(),
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
		NextAllowedLogin:    &pastTime,
		TermsAccepted:       true,
		AcceptedTermsID:     &termsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         true,
		IsStaff:             true,
		Email:               "",
		DateJoined:          "",
		Name:                nickname,
		IsApproved:          true,
		ApprovalNotified:    true,
		Handle:              "",
	}

	err = pg.DB.Create(userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	wagtailCorePageInput := &gorm.WagtailCorePage{
		// WagtailCorePageID:     gofakeit.Number(200, 500),
		Path:                  "/home/123",
		Depth:                 0,
		Numchild:              0,
		Title:                 "test title",
		Slug:                  "test-title",
		Live:                  true,
		HasUnpublishedChanges: false,
		URLPath:               "https://example.com",
		SEOTitle:              "test title",
		ShowInMenus:           false,
		SearchDescription:     "description",
		Expired:               false,
		ContentTypeID:         1, // default to 1 => wagtailcore page
		Locked:                false,
		DraftTitle:            "default title",
		TranslationKey:        uuid.New().String(),
		LocaleID:              1, // default to 1 => en
	}

	err = pg.DB.Create(wagtailCorePageInput).Error
	if err != nil {
		t.Errorf("failed to create wagtail content page: %v", err)
	}

	contentAuthorInput := &gorm.ContentAuthor{
		Active:         true,
		Name:           gofakeit.Name(),
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contentAuthorInput).Error
	if err != nil {
		t.Errorf("failed to create content author: %v", err)
	}

	contentItemInput := &gorm.ContentItem{
		PagePtrID:           wagtailCorePageInput.WagtailCorePageID,
		Date:                time.Now(),
		Intro:               gofakeit.Name(),
		ItemType:            "text",
		TimeEstimateSeconds: 3000,
		Body:                `gofakeit.HipsterParagraph(30, 10, 20, ",")`,
		LikeCount:           10,
		BookmarkCount:       40,
		ShareCount:          0,
		ViewCount:           100,
		AuthorID:            *contentAuthorInput.ContentAuthorID,
	}

	err = pg.DB.Create(contentItemInput).Error
	if err != nil {
		t.Errorf("failed to create content: %v", err)
	}

	type args struct {
		ctx       context.Context
		UserID    string
		ContentID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:       ctx,
				UserID:    *userInput.UserID,
				ContentID: contentItemInput.PagePtrID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing input",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ViewContent(tt.args.ctx, tt.args.UserID, tt.args.ContentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got {
				t.Errorf("expected %v, but got: %v", tt.want, got)
				return
			}
		})
	}
	//TearDown
	if err = pg.DB.Where("content_item_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentView{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("page_ptr_id", contentItemInput.PagePtrID).Unscoped().Delete(&gorm.ContentItem{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", contentAuthorInput.ContentAuthorID).Unscoped().Delete(&gorm.ContentAuthor{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", wagtailCorePageInput.WagtailCorePageID).Unscoped().Delete(&gorm.WagtailCorePage{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestPGInstance_ForgetMe(t *testing.T) {
	ctx := context.Background()

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	testFacility := createTestFacility()
	facility, err := testingDB.GetOrCreateFacility(ctx, testFacility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

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
		Flavour:             feedlib.FlavourConsumer,
		Avatar:              "",
		IsSuspended:         false,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               gofakeit.Email(),
		DateJoined:          "",
		Name:                gofakeit.BeerAlcohol(),
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              gofakeit.BeerName(),
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
		return
	}

	time := time.Now()
	client := &gorm.Client{
		Active:                  true,
		ClientType:              "OVC",
		UserProfile:             gorm.User{},
		TreatmentEnrollmentDate: &time,
		FHIRPatientID:           uuid.New().String(),
		HealthRecordID:          uuid.New().String(),
		TreatmentBuddy:          gofakeit.BeerName(),
		ClientCounselled:        true,
		OrganisationID:          orgID,
		FacilityID:              *facility.FacilityID,
		CHVUserID:               *userInput.UserID,
		UserID:                  userInput.UserID,
	}

	err = pg.DB.Create(client).Error
	if err != nil {
		t.Errorf("failed to create client: %v", err)
	}

	contact := &gorm.Contact{
		ContactType:    "PHONE",
		ContactValue:   gofakeit.PhoneFormatted(),
		Active:         true,
		OptedIn:        true,
		UserID:         userInput.UserID,
		Flavour:        userInput.Flavour,
		OrganisationID: orgID,
	}

	err = pg.DB.Create(contact).Error
	if err != nil {
		t.Errorf("failed to create contact: %v", err)
	}

	invalidUserID := uuid.New().String()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
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
				userID:  *userInput.UserID,
				flavour: userInput.Flavour,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  invalidUserID,
				flavour: userInput.Flavour,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     ctx,
				userID:  invalidUserID,
				flavour: "bad-flavour",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ForgetMe(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ForgetMe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("expected %v, but got: %v", tt.want, got)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("mfl_code", facility.Code).Unscoped().Delete(&gorm.Facility{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.Contact{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("user_id", userInput.UserID).Unscoped().Delete(&gorm.Client{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
