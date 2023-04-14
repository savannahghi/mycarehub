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
