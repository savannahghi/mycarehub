package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth/storage"
)

func TestStorage_CreateRefreshTokenSession(t *testing.T) {

	type args struct {
		ctx       context.Context
		signature string
		request   fosite.Requester
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create access token",
			args: args{
				ctx:       context.Background(),
				signature: gofakeit.Username(),
				request: &fosite.Request{
					Session: &domain.Session{
						ID: gofakeit.UUID(),
					},
					Client: &domain.OauthClient{
						ID: gofakeit.UUID(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to create session",
			args: args{
				ctx:       context.Background(),
				signature: gofakeit.Username(),
				request: &fosite.Request{
					Session: &domain.Session{
						ID: gofakeit.UUID(),
					},
					Client: &domain.OauthClient{
						ID: gofakeit.UUID(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to create token",
			args: args{
				ctx:       context.Background(),
				signature: gofakeit.Username(),
				request: &fosite.Request{
					Session: &domain.Session{
						ID: gofakeit.UUID(),
					},
					Client: &domain.OauthClient{
						ID: gofakeit.UUID(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			s := storage.NewFositeStorage(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case: fail to create session" {
				fakeDB.MockCreateOrUpdateSessionFn = func(ctx context.Context, session *domain.Session) error {
					return fmt.Errorf("failed to create session")
				}
			}

			if tt.name == "sad case: fail to create token" {
				fakeDB.MockCreateRefreshTokenFn = func(ctx context.Context, token *domain.RefreshToken) error {
					return fmt.Errorf("failed to create token")
				}
			}
			if err := s.CreateRefreshTokenSession(tt.args.ctx, tt.args.signature, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateRefreshTokenSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_GetRefreshTokenSession(t *testing.T) {
	type args struct {
		ctx       context.Context
		signature string
		session   fosite.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get token",
			args: args{
				ctx:       context.Background(),
				signature: "signed",
				session:   nil,
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get token",
			args: args{
				ctx:       context.Background(),
				signature: "signed",
				session:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			s := storage.NewFositeStorage(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case: failed to get token" {
				fakeDB.MockGetRefreshTokenFn = func(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
					return nil, fmt.Errorf("failed to get token")
				}
			}

			got, err := s.GetRefreshTokenSession(tt.args.ctx, tt.args.signature, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetRefreshTokenSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Storage.GetRefreshTokenSession() got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteRefreshTokenSession(t *testing.T) {
	type args struct {
		ctx       context.Context
		signature string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete token",
			args: args{
				ctx:       context.Background(),
				signature: "signature",
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to delete token",
			args: args{
				ctx:       context.Background(),
				signature: "signature",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			s := storage.NewFositeStorage(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case: fail to delete token" {
				fakeDB.MockDeleteRefreshTokenFn = func(ctx context.Context, signature string) error {
					return fmt.Errorf("failed to remove token")
				}
			}

			if err := s.DeleteRefreshTokenSession(tt.args.ctx, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteRefreshTokenSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_RevokeRefreshToken(t *testing.T) {
	type args struct {
		ctx       context.Context
		requestID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: revoke token",
			args: args{
				ctx:       context.Background(),
				requestID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to get token",
			args: args{
				ctx:       context.Background(),
				requestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to update token",
			args: args{
				ctx:       context.Background(),
				requestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			s := storage.NewFositeStorage(fakeDB, fakeDB, fakeDB, fakeDB)

			if tt.name == "sad case: fail to get token" {
				fakeDB.MockGetRefreshTokenFn = func(ctx context.Context, token domain.RefreshToken) (*domain.RefreshToken, error) {
					return nil, fmt.Errorf("failed to get token")
				}
			}

			if tt.name == "sad case: fail to update token" {
				fakeDB.MockUpdateRefreshTokenFn = func(ctx context.Context, code *domain.RefreshToken, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update token")
				}
			}

			if err := s.RevokeRefreshToken(tt.args.ctx, tt.args.requestID); (err != nil) != tt.wantErr {
				t.Errorf("Storage.RevokeRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
