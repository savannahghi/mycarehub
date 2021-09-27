package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
)

func TestLoginUseCasesImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()

	flavour := feedlib.FlavourConsumer
	phone := "+254720000000"
	l, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	_, err = CreateOrLoginTestUserByPhone(t)
	if err != nil {
		t.Errorf("failed to create or login test user: %v", err)
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid login",
			args: args{
				ctx:     ctx,
				phone:   phone,
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "sad case: wrong pin number supplied",
			args: args{
				ctx:     ctx,
				phone:   phone,
				PIN:     "4567",
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: user profile without a primary phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2547900900", // not a primary phone number
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2541234",
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect flavour",
			args: args{
				ctx:     ctx,
				phone:   phone,
				PIN:     interserviceclient.TestUserPin,
				flavour: "not-a-correct-flavour",
			},
			// TODO: Return this to true
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := l.LoginByPhone(
				tt.args.ctx,
				tt.args.phone,
				tt.args.PIN,
				tt.args.flavour,
			)
			fmt.Printf("214: THE ERROR IS: %v\n", err)
			if tt.wantErr && authResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_ResumeWIthPin(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	fmt.Printf("THE ERROR IS: %v\n\n", err)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	_, err = CreateOrLoginTestUserByPhone(t)
	if err != nil {
		t.Errorf("failed to create or login test user by phone: %v",
			err,
		)
		return
	}

	l, err := InitializeTestService(ctx)
	fmt.Printf("244: THE ERROR IS: %v\n", err)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
		pin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "valid:_login_with_pin",
			args: args{
				ctx: ctx,
				pin: "1234",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid:_unable_to_get_profile",
			args: args{
				ctx: context.Background(),
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_userprofile_returns_nil",
			args: args{
				ctx: context.Background(),
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_unable_to_get_pin_by_profile_id",
			args: args{
				ctx: context.Background(),
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_pin_data_returns_nil",
			args: args{
				ctx: context.Background(),
				pin: "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid:_pin_mismatch",
			args: args{
				ctx: ctx,
				pin: "1284",
			},
			// if the pins don't match, return false and dont throw an error.
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			isLogin, err := l.ResumeWithPin(
				tt.args.ctx,
				tt.args.pin,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if tt.want != isLogin {
					t.Errorf("expected %v got %v", tt.want, isLogin)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_RefreshToken(t *testing.T) {
	ctx := context.Background()

	l, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_refreshToken",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			wantErr: false,
		},
		// {
		// name: "invalid:invalid_refreshtoken",
		// args: args{
		// ctx: ctx,
		// token: "",
		// },
		// wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := l.RefreshToken(tt.args.ctx, tt.args.token)
			fmt.Printf("284: GOT: %v\n", got)
			fmt.Printf("285: THE ERROR IS : %v\n", err)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ProfileUseCaseImpl.RefreshToken() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_LoginAsAnonymous(t *testing.T) {
	ctx := context.Background()

	l, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.AuthCredentialResponse
		wantErr bool
	}{
		{
			name: "Default Case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := l.LoginAsAnonymous(tt.args.ctx)
			fmt.Printf("GOT IS: %v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ProfileUseCaseImpl.LoginAsAnonymous() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()

	l, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	_, err = CreateOrLoginTestUserByPhone(t)
	if err != nil {
		t.Errorf("unable to create or login test user")
		return
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserResponse
		wantErr bool
	}{
		{
			name: "valid:successfully_login_by_phone",
			args: args{
				ctx:     ctx,
				phone:   "+254720000000",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_normalize_phone",
			args: args{
				ctx:     ctx,
				phone:   "+21",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getUserProfile",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getPin",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_generateAuthCredentials",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_getCustomerOrSupplierProfile",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_comparePin",
			args: args{
				ctx:     ctx,
				phone:   "+254761829103",
				PIN:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := l.LoginByPhone(
				tt.args.ctx,
				tt.args.phone,
				tt.args.PIN,
				tt.args.flavour,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ProfileUseCaseImpl.LoginByPhone() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil user response returned")
					return
				}
			}
		})
	}

}
