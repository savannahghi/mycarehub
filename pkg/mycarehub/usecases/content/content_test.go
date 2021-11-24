package content_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/testutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/serverutils"
)

func TestUseCasesContentImpl_LikeContent_Integration_test(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	termsID := gofakeit.Number(1, 10000)
	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	termsInput := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
		Active:    true,
	}

	err = pg.DB.Create(&termsInput).Error
	if err != nil {
		t.Errorf("failed to create terms: %v", err)
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
		OrganisationID:      serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
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
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
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
			name: "Happy case",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentLike.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: contentLike.ContentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Content.LikeContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.LikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}

	//Teardown
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
	if err = pg.DB.Where("id", termsInput.TermsID).Unscoped().Delete(&gorm.TermsOfService{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestUseCasesContentImpl_UnlikeContent_Integration_test(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	nickname := uuid.New().String()
	currentTime := time.Now()
	flavour := feedlib.FlavourConsumer
	pastTime := time.Now().AddDate(0, 0, -1)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	termsID := gofakeit.Number(1, 10000)
	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	termsInput := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
		Active:    true,
	}

	err = pg.DB.Create(&termsInput).Error
	if err != nil {
		t.Errorf("failed to create terms: %v", err)
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
		OrganisationID:      serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
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
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
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
			name: "Happy case",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: contentLike.ContentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: contentLike.ContentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				userID:    *userInput.UserID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Content.UnlikeContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.UnlikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.UnlikeContent() = %v, want %v", got, tt.want)
			}
		})
	}

	//Teardown
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
	if err = pg.DB.Where("id", termsInput.TermsID).Unscoped().Delete(&gorm.TermsOfService{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
