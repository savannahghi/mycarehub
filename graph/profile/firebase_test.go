package profile_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func TestService_RetrieveUserProfileFirebaseDocSnapshotByMSISDN(t *testing.T) {

	s := profile.NewService()

	ctx := base.GetAuthenticatedContext(t)

	_, err := s.CreateUserByPhone(ctx, base.TestUserPhoneNumberWithPin)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
		return
	}

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := s.SetUserPIN(ctx, base.TestUserPhoneNumberWithPin, base.TestUserPin)
	if !set {
		t.Errorf("can't set a test pin")
	}
	if err != nil {
		t.Errorf("can't set a test pin: %v", err)
		return
	}

	type args struct {
		ctx    context.Context
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good Case",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumberWithPin,
			},
			wantErr: false,
		},
		{
			name: "Bad Case - Empty Msisdn Supplied",
			args: args{
				ctx:    ctx,
				msisdn: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.RetrieveUserProfileFirebaseDocSnapshotByMSISDN(tt.args.ctx, tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RetrieveUserProfileFirebaseDocSnapshotByMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Error("got nil firebase document snapshot")
				return
			}
		})
	}
}

func TestService_RetrievePINFirebaseDocSnapshotByMSISDN(t *testing.T) {
	s := profile.NewService()

	ctx := base.GetAuthenticatedContext(t)

	_, err := s.CreateUserByPhone(ctx, base.TestUserPhoneNumberWithPin)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
		return
	}

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	set, err := s.SetUserPIN(ctx, base.TestUserPhoneNumberWithPin, base.TestUserPin)
	if !set {
		t.Errorf("can't set a test pin")
	}
	if err != nil {
		t.Errorf("can't set a test pin: %v", err)
		return
	}

	type args struct {
		ctx    context.Context
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good Case",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumberWithPin,
			},
			wantErr: false,
		},
		{
			name: "Bad Case - Empty Msisdn Supplied",
			args: args{
				ctx:    ctx,
				msisdn: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(tt.args.ctx, tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RetrievePINFirebaseDocSnapshotByMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Error("got nil firebase document snapshot")
				return
			}
		})
	}
}
