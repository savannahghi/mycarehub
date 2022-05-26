package content_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	helpers_mock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content/mock"
)

func TestUsecaseContentImpl_ListContentCategories(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ContentItemCategory
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewContentUsecaseMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case" {
				fakeDB.MockListContentCategoriesFn = func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseContentImpl.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected content to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected content not to be nil for %v", tt.name)
				return
			}
		})
	}
}
func TestUseCasesContentImpl_LikeContent(t *testing.T) {
	ctx := context.Background()

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
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
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
			_ = mock.NewContentUsecaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case" {
				fakeDB.MockLikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockLikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no contentID" {
				fakeDB.MockLikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and contentID" {
				fakeDB.MockLikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.LikeContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.LikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseContentImpl_ShareContent(t *testing.T) {

	ctx := context.Background()

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
			name: "Happy Case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ContentID: gofakeit.Number(1, 100),
					UserID:    uuid.New().String(),
					Channel:   "SMS",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - no userID",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ContentID: gofakeit.Number(1, 100),
					Channel:   "SMS",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad Case - no userID" {
				fakeDB.MockShareContentFn = func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.Update.ShareContent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseContentImpl.ShareContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseContentImpl.ShareContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_UnlikeContent(t *testing.T) {
	ctx := context.Background()

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
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
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
			_ = mock.NewContentUsecaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case" {
				fakeDB.MockUnlikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockUnlikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no contentID" {
				fakeDB.MockUnlikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and contentID" {
				fakeDB.MockUnlikeContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.UnlikeContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.UnlikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.UnlikeContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_GetUserBookmarkedContent(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully get user bookmarked content",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get user bookmarked content, no content",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing user ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get content",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - user not found",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeContent := mock.NewContentUsecaseMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)
			fakeHelpers := helpers_mock.NewFakeHelper()

			if tt.name == "Happy Case - Successfully get user bookmarked content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {

					cntnt := domain.Content{
						Items: []domain.ContentItem{
							{
								ID: 10,
							},
						},
					}

					payload, err := json.Marshal(cntnt)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeContent.MockGetUserBookmarkedContentFn = func(ctx context.Context, userID string) (*domain.Content, error) {
					fakeHelpers.MockFakeReportErrorToSentryFn = func(err error) {}
					return nil, fmt.Errorf("user ID is required")
				}
			}

			if tt.name == "Sad Case - Fail to get content" {
				fakeDB.MockGetUserBookmarkedContentFn = func(ctx context.Context, userID string) ([]*domain.ContentItem, error) {
					return nil, fmt.Errorf("failed to get bookmarked content")
				}
			}
			if tt.name == "Happy Case - Successfully get user bookmarked content, no content" {
				fakeDB.MockGetUserBookmarkedContentFn = func(ctx context.Context, userID string) ([]*domain.ContentItem, error) {
					return []*domain.ContentItem{}, nil
				}
			}

			if tt.name == "Sad Case - user not found" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("user not found")
				}
			}
			got, err := c.GetUserBookmarkedContent(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.GetUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesContentImpl_GetContent(t *testing.T) {
	ctx := context.Background()
	categoryID := 1
	type args struct {
		ctx        context.Context
		categoryID *int
		limit      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get content",
			args: args{
				ctx:        ctx,
				limit:      "10",
				categoryID: &categoryID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Happy Case - Successfully get content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {

					cntnt := domain.Content{
						Items: []domain.ContentItem{
							{
								ID: 10,
							},
						},
					}

					payload, err := json.Marshal(cntnt)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.GetContent(tt.args.ctx, tt.args.categoryID, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesContentImpl_GetContentByContentItemID(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get content",
			args: args{
				ctx:       ctx,
				contentID: int(uuid.New()[8]),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to make request",
			args: args{
				ctx:       ctx,
				contentID: int(uuid.New()[8]),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad Case - Failed to make request" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make request")
				}
			}

			if tt.name == "Happy Case - Successfully get content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {

					cntnt := domain.Content{
						Items: []domain.ContentItem{
							{
								ID: 10,
							},
						},
					}

					payload, err := json.Marshal(cntnt)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.GetContentByContentItemID(tt.args.ctx, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.GetContentByContentItemID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesContentImpl_ViewContent(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy case - Successfully update view count",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to update view count",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Sad Case - Missing user ID",
			args:    args{},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad Case - Fail to update view count" {
				fakeDB.MockViewContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("failed to update view count")
				}
			}

			if tt.name == "Sad Case - Missing user ID" {
				fakeDB.MockViewContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("failed to update view count")
				}
			}

			got, err := c.ViewContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.ViewContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.ViewContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_BookmarkContent(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully bookmark content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to bookmark content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad Case - Fail to bookmark content" {
				fakeDB.MockBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("failed to bookmark content")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no contentID" {
				fakeDB.MockBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and contentID" {
				fakeDB.MockBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.BookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.BookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.BookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_CheckWhetherUserHasLikedContent(t *testing.T) {
	ctx := context.Background()

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
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty userID",
			args: args{
				ctx:       ctx,
				userID:    "",
				contentID: gofakeit.Number(1, 1001),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid contentID",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: -5,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case" {
				fakeDB.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty userID" {
				fakeDB.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty contentID" {
				fakeDB.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - invalid contentID" {
				fakeDB.MockCheckWhetherUserHasLikedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.CheckWhetherUserHasLikedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.CheckWhetherUserHasLikedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.CheckWhetherUserHasLikedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_CheckIfUserBookmarkedContent(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully check if user bookmarked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully check if user bookmarked content, but user has not bookmarked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to check if user bookmarked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Missing user ID",
			args: args{
				ctx:       ctx,
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Missing content ID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Happy Case - Successfully check if user bookmarked content, but user has not bookmarked content" {
				fakeDB.MockCheckIfUserBookmarkedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "Sad Case - Fail to check if user bookmarked content" {
				fakeDB.MockCheckIfUserBookmarkedContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("failed to check if user bookmarked content")
				}
			}

			got, err := c.CheckIfUserBookmarkedContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.CheckIfUserBookmarkedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.CheckIfUserBookmarkedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_UnBookmarkContent(t *testing.T) {
	ctx := context.Background()

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
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case" {
				fakeDB.MockUnBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockUnBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no contentID" {
				fakeDB.MockUnBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and contentID" {
				fakeDB.MockUnBookmarkContentFn = func(ctx context.Context, userID string, contentID int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.UnBookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.UnBookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.UnBookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesContentImpl_ShareContent(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakeExt := extensionMock.NewFakeExtension()
	c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

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
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					UserID:    uuid.New().String(),
					ContentID: 20,
					Channel:   "test",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					UserID:    uuid.New().String(),
					ContentID: 20,
					Channel:   "test",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockShareContentFn = func(ctx context.Context, input dto.ShareContentInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := c.ShareContent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.ShareContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesContentImpl.ShareContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
