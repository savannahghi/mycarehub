package content_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
)

func TestUseCasesContentImpl_LikeContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		clientID  string
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
				clientID:  uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				clientID:  "",
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx:       ctx,
				clientID:  "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 1,
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.LikeContent(tt.args.ctx, tt.args.clientID, tt.args.contentID)
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
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ContentID: gofakeit.Number(1, 100),
					ClientID:  uuid.New().String(),
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
		{
			name: "sad case: make request error",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ClientID:  gofakeit.UUID(),
					ContentID: gofakeit.Number(1, 100),
					Channel:   "SMS",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ClientID:  gofakeit.UUID(),
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.ShareContent(tt.args.ctx, tt.args.input)
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
		clientID  string
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
				clientID:  uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:       ctx,
				clientID:  "",
				contentID: 1,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				ctx:       ctx,
				clientID:  "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 1,
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					result := struct {
						Count   int `json:"count"`
						Results []struct {
							ID string `json:"id"`
						} `json:"results"`
					}{
						Count: 1,
						Results: []struct {
							ID string "json:\"id\""
						}{
							{
								ID: gofakeit.UUID(),
							},
						},
					}

					payload, err := json.Marshal(result)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					result := struct {
						Count   int `json:"count"`
						Results []struct {
							ID string `json:"id"`
						} `json:"results"`
					}{
						Count: 1,
						Results: []struct {
							ID string "json:\"id\""
						}{
							{
								ID: gofakeit.UUID(),
							},
						},
					}

					payload, err := json.Marshal(result)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.UnlikeContent(tt.args.ctx, tt.args.clientID, tt.args.contentID)
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
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user bookmarked content",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get user bookmarked content, no content",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
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
			name: "sad case: make request error",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

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
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			got, err := c.GetUserBookmarkedContent(tt.args.ctx, tt.args.clientID)
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
	clientID := gofakeit.UUID()
	type args struct {
		ctx           context.Context
		categoryIDs   []int
		limit         string
		categoryNames []string
		clientID      *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get content",
			args: args{
				ctx:           ctx,
				limit:         "10",
				categoryIDs:   []int{categoryID},
				categoryNames: []string{"Chemotherapy"},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get content as caregiver",
			args: args{
				ctx:           ctx,
				limit:         "10",
				categoryIDs:   []int{categoryID},
				clientID:      &clientID,
				categoryNames: []string{"Chemotherapy"},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to get logged in user",
			args: args{
				ctx:         ctx,
				limit:       "10",
				categoryIDs: []int{categoryID},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get user profile of the logged in user",
			args: args{
				ctx:         ctx,
				limit:       "10",
				categoryIDs: []int{categoryID},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get client profile of the logged in user",
			args: args{
				ctx:         ctx,
				limit:       "10",
				categoryIDs: []int{categoryID},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to make http request",
			args: args{
				ctx:         ctx,
				limit:       "10",
				categoryIDs: []int{categoryID},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get client",
			args: args{
				ctx:         ctx,
				limit:       "10",
				categoryIDs: []int{categoryID},
				clientID:    &clientID,
			},
			wantErr: true,
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
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}
			if tt.name == "Sad Case - unable to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}
			if tt.name == "Sad Case - unable to get user profile of the logged in user" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}
			if tt.name == "Sad Case - unable to get client profile of the logged in user" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}
			if tt.name == "Sad Case - unable to make http request" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make http request")
				}
			}

			if tt.name == "Sad Case - unable to get client" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client")
				}
			}

			got, err := c.GetContent(tt.args.ctx, tt.args.categoryIDs, tt.args.categoryNames, tt.args.limit, tt.args.clientID)
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

func TestUseCasesContentImpl_GetContentItemByID(t *testing.T) {
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

					cntnt := domain.ContentItem{
						ID: 10,
					}

					payload, err := json.Marshal(cntnt)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.GetContentItemByID(tt.args.ctx, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.GetContentItemByID() error = %v, wantErr %v", err, tt.wantErr)
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
		clientID  string
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
				clientID:  uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.ViewContent(tt.args.ctx, tt.args.clientID, tt.args.contentID)
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
		clientID  string
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
				clientID:  uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
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
				ctx:      ctx,
				clientID: uuid.New().String(),
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy Case - Successfully bookmark content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.BookmarkContent(tt.args.ctx, tt.args.clientID, tt.args.contentID)
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
			name: "happy case: has liked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: has not liked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "sad case: make request error",
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
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: gofakeit.Number(1, 1001),
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 1,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "happy case: has liked content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 1,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "happy case: has not liked content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 0,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
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
			name: "happy case: has bookmarked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: has not bookmarked content",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "sad case: make request error",
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
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				contentID: 12,
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 1,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "happy case: has bookmarked content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 1,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "happy case: has not bookmarked content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					count := struct {
						Count int `json:"count"`
					}{
						Count: 0,
					}

					payload, err := json.Marshal(count)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
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
		clientID  string
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
				clientID:  uuid.New().String(),
				contentID: 1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
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
				ctx:      ctx,
				clientID: uuid.New().String(),
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
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:       ctx,
				clientID:  uuid.New().String(),
				contentID: 1,
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

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					result := struct {
						Count   int `json:"count"`
						Results []struct {
							ID string `json:"id"`
						} `json:"results"`
					}{
						Count: 1,
						Results: []struct {
							ID string "json:\"id\""
						}{
							{
								ID: gofakeit.UUID(),
							},
						},
					}

					payload, err := json.Marshal(result)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					result := struct {
						Count   int `json:"count"`
						Results []struct {
							ID string `json:"id"`
						} `json:"results"`
					}{
						Count: 1,
						Results: []struct {
							ID string "json:\"id\""
						}{
							{
								ID: gofakeit.UUID(),
							},
						},
					}

					payload, err := json.Marshal(result)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			got, err := c.UnBookmarkContent(tt.args.ctx, tt.args.clientID, tt.args.contentID)
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
					ClientID:  uuid.New().String(),
					ContentID: 20,
					Channel:   "test",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: user id not provided",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ContentID: 20,
					Channel:   "test",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: content id not provided",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ClientID: uuid.New().String(),
					Channel:  "test",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ClientID:  uuid.New().String(),
					ContentID: 20,
					Channel:   "test",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					ClientID:  uuid.New().String(),
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
			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Happy case" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					payload, err := json.Marshal([]byte{})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
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

func TestUseCasesContentImpl_GetFAQs(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get consumer faqs",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Happy case: get pro faqs",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "sad case: make request error",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "sad case: make request error" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			_, err := c.GetFAQs(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.GetFAQs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesContentImpl_ListContentCategories(t *testing.T) {

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
			name: "happy case: fetch content categories",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get logged in user",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get user profile of the logged in user",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to make request",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "sad case: unable to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}
			if tt.name == "sad case: unable to get user profile of the logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}
			if tt.name == "sad case: unable to make request" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make request")
				}
			}

			_, err := c.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesContentImpl_FetchContent(t *testing.T) {
	type args struct {
		ctx   context.Context
		limit string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: fetch content",
			args: args{
				limit: "15",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch content",
			args: args{
				limit: "15",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			c := content.NewUseCasesContentImplementation(fakeDB, fakeDB, fakeExt)

			if tt.name == "Sad case: unable to fetch content" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make request")
				}
			}

			_, err := c.FetchContent(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesContentImpl.FetchContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
