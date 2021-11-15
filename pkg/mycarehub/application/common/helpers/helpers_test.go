package helpers

import (
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func TestGetInviteLink(t *testing.T) {
	type args struct {
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "valid: flavour is valid, pro",
			args: args{
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "valid: flavour is valid, consumer",
			args: args{
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid: flavour is invalid",
			args: args{
				flavour: feedlib.Flavour("invalid"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetInviteLink(tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInviteLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreateInviteMessage(t *testing.T) {

	firstName := "Jon"
	lastName := "Doe"
	inviteLink := "https://example.com"
	pin := "99033"

	want := fmt.Sprintf("Dear %v %v, you have been invited to My Afya Hub. Download the app on %v. Your single use pin is %v",
		firstName, lastName, inviteLink, pin)
	type args struct {
		user       *domain.User
		inviteLink string
		pin        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default case",
			args: args{
				user: &domain.User{
					FirstName: firstName,
					LastName:  lastName,
				},
				inviteLink: inviteLink,
				pin:        pin,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateInviteMessage(tt.args.user, tt.args.inviteLink, tt.args.pin); got != tt.want {
				t.Errorf("CreateInviteMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateResetPinMessage(t *testing.T) {
	firstName := "Jon"
	lastName := "Doe"

	pin := "99033"

	want := fmt.Sprintf("Dear %v %v, your PIN for My Afya Hub has been reset successfully. Your single use pin is %v",
		firstName, lastName, pin)

	type args struct {
		user *domain.User
		pin  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default case",
			args: args{
				user: &domain.User{
					FirstName: firstName,
					LastName:  lastName,
				},
				pin: pin,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateResetPinMessage(tt.args.user, tt.args.pin); got != tt.want {
				t.Errorf("CreateResetPinMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
