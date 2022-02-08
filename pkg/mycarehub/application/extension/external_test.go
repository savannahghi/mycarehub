package extension

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	openSourceDto "github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	ext     External
	fakeExt mock.FakeExtensionImpl
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")

	log.Printf("Running tests ...")
	code := m.Run()

	os.Exit(code)
}

func TestExternal_CreateFirebaseCustomToken(t *testing.T) {
	ctx := context.Background()
	uid := ksuid.New().String()
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx: ctx,
				uid: uid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "default case" {
				fakeExt.MockAuthenticateCustomFirebaseTokenFn = func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
					return &firebasetools.FirebaseUserTokens{}, nil
				}
			}
			got, err := ext.CreateFirebaseCustomToken(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.CreateFirebaseCustomToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestExternal_SendFeedback(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx             context.Context
		subject         string
		feedbackMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
		panic   bool
	}{
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			panic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				fcSendFeedback := func() { _, _ = ext.SendFeedback(tt.args.ctx, tt.args.subject, tt.args.feedbackMessage) }
				assert.Panics(t, fcSendFeedback)
				return
			}
			got, err := ext.SendFeedback(tt.args.ctx, tt.args.subject, tt.args.feedbackMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.SendFeedback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("External.SendFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_AuthenticateCustomFirebaseToken(t *testing.T) {
	type args struct {
		customAuthToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *firebasetools.FirebaseUserTokens
		wantErr bool
	}{
		{
			name:    "invalid: missing token",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ext.AuthenticateCustomFirebaseToken(tt.args.customAuthToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.AuthenticateCustomFirebaseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestExternal_ComparePIN(t *testing.T) {
	type args struct {
		rawPwd     string
		salt       string
		encodedPwd string
		options    *extension.Options
	}
	tests := []struct {
		name   string
		args   args
		want   bool
		panics bool
	}{
		{
			name:   "invalid: missing params",
			args:   args{},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				fcComparePIN := func() { _ = ext.ComparePIN(tt.args.rawPwd, tt.args.salt, tt.args.encodedPwd, tt.args.options) }
				assert.Panics(t, fcComparePIN)
				return
			}
			if got := ext.ComparePIN(tt.args.rawPwd, tt.args.salt, tt.args.encodedPwd, tt.args.options); got != tt.want {
				t.Errorf("External.ComparePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_GetLoggedInUserUID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid: user not in context",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ext.GetLoggedInUserUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.GetLoggedInUserUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("External.GetLoggedInUserUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_SendSMS(t *testing.T) {

	type args struct {
		ctx          context.Context
		phoneNumbers string
		message      string
		from         enumutils.SenderID
	}
	tests := []struct {
		name    string
		args    args
		want    *openSourceDto.SendMessageResponse
		wantErr bool
		panics  bool
	}{
		{
			name:   "invalid: missing params",
			args:   args{},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				fcSendSMS := func() { _, _ = ext.SendSMS(tt.args.ctx, tt.args.phoneNumbers, tt.args.message, tt.args.from) }
				assert.Panics(t, fcSendSMS)
				return
			}
			got, err := ext.SendSMS(tt.args.ctx, tt.args.phoneNumbers, tt.args.message, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("External.SendSMS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_GenerateAndSendOTP(t *testing.T) {
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		panics  bool
	}{
		{
			name:    "invalid: missing params",
			args:    args{},
			wantErr: true,
			panics:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				fcGenerateAndSendOTP := func() { _, _ = ext.GenerateAndSendOTP(tt.args.ctx, tt.args.phoneNumber) }
				assert.Panics(t, fcGenerateAndSendOTP)
				return
			}
			got, err := ext.GenerateAndSendOTP(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.GenerateAndSendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("External.GenerateAndSendOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_GenerateRetryOTP(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload *dto.SendRetryOTPPayload
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		panics  bool
	}{
		{
			name:    "invalid: missing params",
			args:    args{},
			wantErr: true,
			panics:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				fcGenerateRetryOTP := func() { _, _ = ext.GenerateRetryOTP(tt.args.ctx, tt.args.payload) }
				assert.Panics(t, fcGenerateRetryOTP)
				return
			}
			got, err := ext.GenerateRetryOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.GenerateRetryOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("External.GenerateRetryOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternal_SendInviteSMS(t *testing.T) {
	type args struct {
		ctx         context.Context
		phoneNumber string
		message     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		panics  bool
	}{
		{
			name:    "invalid: missing params",
			args:    args{},
			wantErr: true,
			panics:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				fcSendInviteSMS := func() { _ = ext.SendInviteSMS(tt.args.ctx, tt.args.phoneNumber, tt.args.message) }
				assert.Panics(t, fcSendInviteSMS)
				return
			}
			if err := ext.SendInviteSMS(tt.args.ctx, tt.args.phoneNumber, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("External.SendInviteSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
