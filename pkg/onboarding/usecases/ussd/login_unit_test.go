package ussd_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestUSSDImpl_HandleLogin_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
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

		{
			name: "Sad case : bad login PIN longer than 4 digits",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "12345",
			},
			want: "CON The PIN you entered is not correct\r\n" +
				"Please try again (enter 00 if you\r\n" +
				"forgot your PIN)",
		},

		{
			name: "Sad case : invalid input",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "hello",
			},
			want: "CON The PIN you entered is not correct\r\n" +
				"Please try again (enter 00 if you\r\n" +
				"forgot your PIN)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case : empty input" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}

				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}
			}

			if tt.name == "Sad case : bad login PIN" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return nil, err
				}

				fakePinExt.ComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad case : bad login PIN longer than 4 digits" {
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return nil, err
				}

				fakePinExt.ComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad case : invalid input" {

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return nil, err
				}

				fakePinExt.ComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
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
