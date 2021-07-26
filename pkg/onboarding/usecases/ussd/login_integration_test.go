package ussd_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
)

const (
	// EmptyInput is used to load a default menu when user has not supplied any input
	EmptyInput = ""
	// ForgotPINInput indicates the user has forgotten their PIN and would like to reset it
	ForgotPINInput = "00"
	//ForgotPINVerifyDate indicates the state when a use wants to reset PIN
	ForgotPINVerifyDate = 15
)

func TestImpl_HandleLogin_Integration_Test(t *testing.T) {
	ctx := context.Background()
	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service %v", err)
		return
	}

	phoneNumber := "+254700100200"
	dateOfBirth := "12122000"
	PIN := "1234"
	FirstName := gofakeit.LastName()
	LastName := gofakeit.LastName()
	SessionID := uuid.New().String()
	Level := 0

	ussdDet := &domain.USSDLeadDetails{
		ID:          gofakeit.UUID(),
		Level:       Level,
		PhoneNumber: phoneNumber,
		SessionID:   SessionID,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: dateOfBirth,
		PIN:         PIN,
	}

	sessionDet := &dto.SessionDetails{
		SessionID:   SessionID,
		PhoneNumber: &phoneNumber,
		Level:       Level,
		Text:        "",
	}

	sessionDetails, err := u.AITUSSD.AddAITSessionDetails(ctx, sessionDet)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	type args struct {
		ctx          context.Context
		session      *domain.USSDLeadDetails
		userResponse string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy case : empty input",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: EmptyInput,
			},
			want: "CON Welcome to Be.Well.Please enter\r\n" +
				"your PIN to continue(enter 00 if\r\n" +
				"you forgot your PIN)\r\n",
		},

		{
			name: "Happy case : forgot PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: ForgotPINInput,
			},
			want: "CON Please enter your date of birth in\r\n" +
				"DDMMYYYY format e.g 14031996 for\r\n" +
				"14th March 1996\r\n" +
				"to be able to reset PIN\r\n",
		},

		{
			name: "Sad case : forgot PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "01",
			},
			want: "END Something went wrong. Please try again.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case : empty input" {

				err = u.AITUSSD.UpdateSessionLevel(ctx, HomeMenuState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Happy case : forgot PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ForgotPINVerifyDate, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case : forgot PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, HomeMenuState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			session, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occured %v", err)
				return
			}
			login := u
			got := login.AITUSSD.HandleLogin(tt.args.ctx, session, tt.args.userResponse)
			if got != tt.want {
				t.Errorf("Impl.HandleLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
