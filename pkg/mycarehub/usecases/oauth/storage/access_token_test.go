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

func TestStorage_CreateAccessTokenSession(t *testing.T) {

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
			name: "happy case: create token",
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
				fakeDB.MockCreateAccessTokenFn = func(ctx context.Context, token *domain.AccessToken) error {
					return fmt.Errorf("failed to create token")
				}
			}

			if err := s.CreateAccessTokenSession(tt.args.ctx, tt.args.signature, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateAccessTokenSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_GetAccessTokenSession(t *testing.T) {
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
				fakeDB.MockGetAccessTokenFn = func(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error) {
					return nil, fmt.Errorf("failed to get token")
				}
			}

			got, err := s.GetAccessTokenSession(tt.args.ctx, tt.args.signature, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetAccessTokenSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Storage.GetAccessTokenSession() got = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteAccessTokenSession(t *testing.T) {
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
				fakeDB.MockDeleteAccessTokenFn = func(ctx context.Context, signature string) error {
					return fmt.Errorf("failed to remove token")
				}
			}

			if err := s.DeleteAccessTokenSession(tt.args.ctx, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteAccessTokenSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_RevokeAccessToken(t *testing.T) {
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
				fakeDB.MockGetAccessTokenFn = func(ctx context.Context, token domain.AccessToken) (*domain.AccessToken, error) {
					return nil, fmt.Errorf("failed to get token")
				}
			}

			if tt.name == "sad case: fail to update token" {
				fakeDB.MockUpdateAccessTokenFn = func(ctx context.Context, code *domain.AccessToken, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update token")
				}
			}

			if err := s.RevokeAccessToken(tt.args.ctx, tt.args.requestID); (err != nil) != tt.wantErr {
				t.Errorf("Storage.RevokeAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
