package ussd_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_HandleLogin_Unittest(t *testing.T) {
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
			name: "Happy case : good login PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "4321",
			},
			want: "CON Welcome to Be.Well\r\n" +
				"1. Opt out from marketing messages\r\n" +
				"2. Change PIN",
		},

		{
			name: "Sad case : bad login PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "1",
			},
			want: "CON The PIN you entered is not correct\r\n" +
				"Please try again (enter 00 if you\r\n" +
				"forgot your PIN)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case : good login PIN" {
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{
						Level:       Level,
						PhoneNumber: phoneNumber,
						SessionID:   sessionDetails.SessionID,
						FirstName:   FirstName,
						LastName:    LastName,
						DateOfBirth: dateOfBirth,
						PIN:         PIN,
					}, nil
				}
			}

			if tt.name == "Happy case : forgot PIN" {
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{
						Level:       Level,
						PhoneNumber: phoneNumber,
						SessionID:   sessionDetails.SessionID,
						FirstName:   FirstName,
						LastName:    LastName,
						DateOfBirth: dateOfBirth,
						PIN:         "00",
					}, nil
				}
			}

			session, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occured %v", err)
				return
			}

			if got := u.AITUSSD.HandleLogin(tt.args.ctx, session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandleLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
