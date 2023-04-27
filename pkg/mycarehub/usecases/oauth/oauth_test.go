package oauth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth"
	"gorm.io/gorm"
)

func TestUseCasesOauthImpl_CreateOauthClient(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.OauthClientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create oauth client",
			args: args{
				ctx: context.Background(),
				input: dto.OauthClientInput{
					Name:   "Client One",
					Secret: gofakeit.Password(true, true, true, true, false, 10),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating client",
			args: args{
				ctx: context.Background(),
				input: dto.OauthClientInput{
					Name:   "Client One",
					Secret: gofakeit.Password(true, true, true, true, false, 10),
				},
			},
			wantErr: true,
		},
		{
			name: "happy case: create oauth client",
			args: args{
				ctx: context.Background(),
				input: dto.OauthClientInput{
					Name:   "Client One",
					Secret: gofakeit.Password(true, true, true, true, false, 10),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			u := oauth.NewUseCasesOauthImplementation(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case: error creating client" {
				fakeDB.MockCreateOauthClient = func(ctx context.Context, client *domain.OauthClient) error {
					return fmt.Errorf("failed to create client")
				}
			}

			got, err := u.CreateOauthClient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesOauthImpl.CreateOauthClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesOauthImpl.CreateOauthClient() got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestUseCasesOauthImpl_GenerateUserAuthTokens(t *testing.T) {

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
			name: "happy case: generate user tokens existing client",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "happy case: generate user tokens new client",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to get oauth client",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get user profile",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid oauth client",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			u := oauth.NewUseCasesOauthImplementation(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "happy case: generate user tokens new client" {
				fakeDB.MockGetOauthClient = func(ctx context.Context, id string) (*domain.OauthClient, error) {
					return nil, gorm.ErrRecordNotFound
				}
			}

			if tt.name == "sad case: fail to get oauth client" {
				fakeDB.MockGetOauthClient = func(ctx context.Context, id string) (*domain.OauthClient, error) {
					return nil, fmt.Errorf("database error")
				}
			}

			if tt.name == "sad case: fail to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, gorm.ErrRecordNotFound
				}
			}

			if tt.name == "sad case: invalid oauth client" {
				fakeDB.MockGetOauthClient = func(ctx context.Context, id string) (*domain.OauthClient, error) {
					return &domain.OauthClient{}, nil
				}
			}

			got, err := u.GenerateUserAuthTokens(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesOauthImpl.GenerateUserAuthTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesOauthImpl.GenerateUserAuthTokens() got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestUseCasesOauthImpl_RefreshAutToken(t *testing.T) {
	type args struct {
		ctx          context.Context
		refreshToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad case: invalid token",
			args: args{
				ctx:          context.Background(),
				refreshToken: gofakeit.BS(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed refresh token",
			args: args{
				ctx:          context.Background(),
				refreshToken: gofakeit.BS(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			u := oauth.NewUseCasesOauthImplementation(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case: failed refresh token" {
				fakeDB.MockGetOauthClient = func(ctx context.Context, id string) (*domain.OauthClient, error) {
					return nil, gorm.ErrRecordNotFound
				}
			}

			_, err := u.RefreshAutToken(tt.args.ctx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesOauthImpl.RefreshAutToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
