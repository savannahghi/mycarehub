package profile

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_GetConfirmedEmailAddresses(t *testing.T) {
	service := NewService()
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	assert.NotNil(t, ctx)
	assert.NotNil(t, token)

	profile, err := service.UserProfile(ctx)
	assert.NotNil(t, profile)
	assert.Nil(t, err)

	phoneNumberCtx := base.GetPhoneNumberAuthenticatedContext(t)

	phoneNumberProfile, err := service.UserProfile(phoneNumberCtx)
	assert.Nil(t, err)
	assert.NotNil(t, phoneNumberProfile)

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    map[string][]string
	}{
		{
			name: "Existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{token.UID},
			},
			wantErr: false,
			want: map[string][]string{
				token.UID: profile.Emails,
			},
		},
		{
			name: "Non existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{"not a uid"},
			},
			wantErr: false,
			want: map[string][]string{
				"not a uid": {},
			},
		},
		{
			name: "Slice of uids case",
			args: args{
				ctx:  ctx,
				uids: []string{phoneNumberProfile.UID, token.UID},
			},
			wantErr: false,
			want: map[string][]string{
				token.UID:              profile.Emails,
				phoneNumberProfile.UID: {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			confirmedEmails, err := s.GetConfirmedEmailAddresses(tt.args.ctx, tt.args.uids)
			log.Printf("Confirmed Emails....%v", confirmedEmails)
			if err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, confirmedEmails)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetConfirmedEmailAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(confirmedEmails, tt.want) {
				t.Errorf("Service.GetConfirmedEmailAddresses() = %v, want %v", confirmedEmails, tt.want)
			}
		})
	}
}

func TestService_GetConfirmedPhoneNumbers(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	assert.NotNil(t, ctx)

	profile, err := service.UserProfile(ctx)
	assert.NotNil(t, profile)
	assert.Nil(t, err)

	phoneNumberCtx := base.GetPhoneNumberAuthenticatedContext(t)

	phoneNumberProfile, err := service.UserProfile(phoneNumberCtx)
	assert.Nil(t, err)
	assert.NotNil(t, phoneNumberProfile)

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "Existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{profile.UID},
			},
			wantErr: false,
			want: map[string][]string{
				profile.UID: profile.Msisdns,
			},
		},
		{
			name: "Non existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{"not a uid"},
			},
			wantErr: false,
			want: map[string][]string{
				"not a uid": {},
			},
		},
		{
			name: "Slice of uids case",
			args: args{
				ctx:  ctx,
				uids: []string{phoneNumberProfile.UID, profile.UID},
			},
			wantErr: false,
			want: map[string][]string{
				profile.UID:            profile.Msisdns,
				phoneNumberProfile.UID: {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.GetConfirmedPhoneNumbers(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetConfirmedPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetConfirmedPhoneNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetValidFCMTokens(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	assert.NotNil(t, ctx)

	profile, err := service.UserProfile(ctx)
	assert.NotNil(t, profile)
	assert.Nil(t, err)

	phoneNumberCtx := base.GetPhoneNumberAuthenticatedContext(t)

	phoneNumberProfile, err := service.UserProfile(phoneNumberCtx)
	assert.Nil(t, err)
	assert.NotNil(t, phoneNumberProfile)

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "Existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{profile.UID},
			},
			wantErr: false,
			want: map[string][]string{
				profile.UID: profile.PushTokens,
			},
		},
		{
			name: "Non existing uid case",
			args: args{
				ctx:  ctx,
				uids: []string{"not a uid"},
			},
			wantErr: false,
			want: map[string][]string{
				"not a uid": {},
			},
		},
		{
			name: "Slice of uids case",
			args: args{
				ctx:  ctx,
				uids: []string{phoneNumberProfile.UID, profile.UID},
			},
			wantErr: false,
			want: map[string][]string{
				profile.UID:            profile.PushTokens,
				phoneNumberProfile.UID: {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.GetValidFCMTokens(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetValidFCMTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetValidFCMTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
