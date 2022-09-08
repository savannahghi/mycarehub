package helpers

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
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

	name := "John Doe"
	inviteLink := "https://example.com"
	pin := "99033"
	consumerAppName := serverutils.MustGetEnvVar("CONSUMER_APP_NAME")
	proAppName := serverutils.MustGetEnvVar("PRO_APP_NAME")

	consumerMessage := fmt.Sprintf("You have been invited to %s. Download the app on %v. Your single use pin is %v",
		consumerAppName, inviteLink, pin)
	proMessage := fmt.Sprintf("You have been invited to %s. Download the app on %v. Your single use pin is %v",
		proAppName, inviteLink, pin)

	type args struct {
		user       *domain.User
		inviteLink string
		pin        string
		flavour    feedlib.Flavour
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Case - Create consumer invite message",
			args: args{
				user: &domain.User{
					Name: name,
				},
				inviteLink: inviteLink,
				pin:        pin,
				flavour:    feedlib.FlavourConsumer,
			},
			want: consumerMessage,
		},
		{
			name: "Happy Case - Create Pro invite message",
			args: args{
				user: &domain.User{
					Name: name,
				},
				inviteLink: inviteLink,
				pin:        pin,
				flavour:    feedlib.FlavourPro,
			},
			want: proMessage,
		},
		{
			name: "Sad Case - Fail to create message",
			args: args{
				user: &domain.User{
					Name: name,
				},
				inviteLink: inviteLink,
				pin:        pin,
				flavour:    feedlib.Flavour("invalid"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateInviteMessage(tt.args.user, tt.args.inviteLink, tt.args.pin, tt.args.flavour); got != tt.want {
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

func TestGetPinExpiryDate(t *testing.T) {
	tests := []struct {
		name    string
		want    *time.Time
		wantErr bool
	}{
		{
			name:    "default case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetPinExpiryDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPinExpiryDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRestAPIResponseHelper(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *dto.RestEndpointResponses
	}{
		{
			name: "Happy case",
			args: args{
				key:   "test",
				value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RestAPIResponseHelper(tt.args.key, tt.args.value)
			assert.NotNil(t, got)
		})
	}
}

func TestCaptureSentryError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case",
			args: args{
				err: fmt.Errorf("an error occurred"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportErrorToSentry(tt.args.err)
		})
	}
}
