package auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/authutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/auth"
	authMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/auth/mock"
)

func TestSILAuthServiceImpl_AuthenticateWithSlade360(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: authenticate",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to authenticate",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuth := authMock.NewAuthClientMock()
			auth := auth.NewAuthService(fakeAuth)

			if tt.name == "Sad case: unable to authenticate" {
				fakeAuth.MockAuthenticateFn = func() (*authutils.OAUTHResponse, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := auth.AuthenticateWithSlade360(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILAuthServiceImpl.AuthenticateWithSlade360() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
