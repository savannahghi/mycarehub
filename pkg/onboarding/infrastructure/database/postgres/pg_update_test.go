package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	gormmock "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm/mock"
)

func TestOnboardingDb_UpdateUserLastSuccessfulLogin(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx           context.Context
		userID        string
		lastLoginTime time.Time
		flavour       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:           ctx,
				userID:        "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				lastLoginTime: time.Now(),
				flavour:       "CONSUMER",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:           ctx,
				userID:        "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				lastLoginTime: time.Now(),
				flavour:       "invalid-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormmock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.UpdateUserLastSuccessfulLoginFn = func(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if err := d.UpdateUserLastSuccessfulLogin(tt.args.ctx, tt.args.userID, tt.args.lastLoginTime, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.UpdateUserLastSuccessfulLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOnboardingDb_UpdateUserLastFailedLogin(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                 context.Context
		userID              string
		lastFailedLoginTime time.Time
		flavour             string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                 ctx,
				userID:              "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				lastFailedLoginTime: time.Now(),
				flavour:             "CONSUMER",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:                 ctx,
				userID:              "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				lastFailedLoginTime: time.Now(),
				flavour:             "Invalid -flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormmock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.UpdateUserLastFailedLoginFn = func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if err := d.UpdateUserLastFailedLogin(tt.args.ctx, tt.args.userID, tt.args.lastFailedLoginTime, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.UpdateUserLastFailedLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOnboardingDb_UpdateUserFailedLoginCount(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx              context.Context
		userID           string
		failedLoginCount string
		flavour          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				userID:           "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				failedLoginCount: "0",
				flavour:          "CONSUMER",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:              ctx,
				userID:           "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				failedLoginCount: "0",
				flavour:          "Invalid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormmock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.UpdateUserLastFailedLoginFn = func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if err := d.UpdateUserFailedLoginCount(tt.args.ctx, tt.args.userID, tt.args.failedLoginCount, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.UpdateUserFailedLoginCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOnboardingDb_UpdateUserNextAllowedLogin(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                  context.Context
		userID               string
		nextAllowedLoginTime time.Time
		flavour              string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                  ctx,
				userID:               "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				nextAllowedLoginTime: time.Now(),
				flavour:              "CONSUMER",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:                  ctx,
				userID:               "1zixbASMwkk3QTnSDmH0EDHZ6H8",
				nextAllowedLoginTime: time.Now(),
				flavour:              "Invalid-CONSUMER",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormmock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.UpdateUserNextAllowedLoginFn = func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if err := d.UpdateUserNextAllowedLogin(tt.args.ctx, tt.args.userID, tt.args.nextAllowedLoginTime, tt.args.flavour); (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.UpdateUserNextAllowedLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
