package usecases_test

import (
	"context"
	"log"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestUserPinUseCaseImpl_SetUserPIN(t *testing.T) {
	ctx := context.Background()
	flavour := base.FlavourConsumer
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	u, err := createTestUserByPhone(
		ctx,
		flavour,
		base.TestUserPhoneNumber,
		base.TestUserPin,
	)
	if err != nil {
		t.Errorf("failed to create test phone user: %v", err)
		return
	}
	if !u {
		t.Errorf("failed to create a test user")
		return
	}

	type args struct {
		ctx   context.Context
		pin   string
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: valid pin setup - valid payload",
			args: args{
				ctx:   ctx,
				pin:   "1234",
				phone: base.TestUserPhoneNumber,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid payload",
			args: args{
				ctx:   ctx,
				pin:   "",
				phone: base.TestUserPhoneNumber,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid payload - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			authResponse, err := pin.UserPIN.SetUserPIN(tt.args.ctx, tt.args.pin, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.SetUserPIN() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && authResponse != false {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == false {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestUserPinUseCaseImpl_ChangeUserPIN(t *testing.T) {
	ctx := context.Background()
	flavour := base.FlavourConsumer
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	u, err := createTestUserByPhone(
		ctx,
		flavour,
		base.TestUserPhoneNumber,
		base.TestUserPin,
	)
	if err != nil {
		t.Errorf("failed to create test phone user: %v", err)
		return
	}
	if !u {
		t.Errorf("failed to create a test user")
		return
	}

	type args struct {
		ctx   context.Context
		phone string
		pin   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: valid pin setup - valid payload",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				pin:   "12356",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid payload- empty payload",
			args: args{
				ctx:   ctx,
				pin:   "",
				phone: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid payload - empty phone number",
			args: args{
				ctx:   ctx,
				phone: "",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:   ctx,
				phone: "+2541234",
				pin:   base.TestUserPin,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pin := s
			authResponse, err := pin.UserPIN.ChangeUserPIN(tt.args.ctx, tt.args.phone, tt.args.pin)
			log.Println("Error gani hii:", authResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPinUseCaseImpl.ChangeUserPIN() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
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
