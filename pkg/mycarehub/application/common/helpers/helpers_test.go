package helpers

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
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

func TestEncryptSensitiveData(t *testing.T) {

	text := "texttoencrypt"
	bs := []byte(gofakeit.HipsterSentence(32))
	mysecret := string(bs[0:32])

	type args struct {
		text     string
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				text:     text,
				MySecret: mysecret,
			},
			wantErr: false,
		},
		{
			name: "invalid: short secret",
			args: args{
				text:     text,
				MySecret: "mysecret",
			},
			wantErr: true,
		},
		{
			name: "invalid: long secret",
			args: args{
				text:     text,
				MySecret: gofakeit.HipsterSentence(10),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptSensitiveData(tt.args.text, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptSensitiveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}

		})
	}
}

func TestDecryptSensitiveData(t *testing.T) {
	text := "texttoencrypt"
	bs := []byte(gofakeit.HipsterSentence(32))
	mysecret := string(bs[0:32])

	encryptedText, err := EncryptSensitiveData(text, mysecret)
	if err != nil {
		t.Errorf("failed to encrypt text: %v", err)
	}

	type args struct {
		text     string
		MySecret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				text:     encryptedText,
				MySecret: mysecret,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptSensitiveData(tt.args.text, tt.args.MySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptSensitiveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}
