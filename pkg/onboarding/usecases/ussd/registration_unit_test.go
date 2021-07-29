package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/scalarutils"

	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

func TestImpl_HandleUserRegistration_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
	if err != nil {
		t.Errorf("unable to initialize service %v", err)
		return
	}

	phoneNumber := "+254700200210"
	DateOfBirth := scalarutils.Date{
		Day:   0,
		Month: 0,
		Year:  0,
	}
	PIN := "1234"
	FirstName := gofakeit.FirstName()
	LastName := gofakeit.LastName()

	SessionID := uuid.New().String()
	Level := 2
	Text := ""

	sessionDet := &dto.SessionDetails{
		SessionID:   SessionID,
		PhoneNumber: &phoneNumber,
		Level:       Level,
		Text:        Text,
	}

	validUSSDLeadDetails := &domain.USSDLeadDetails{
		ID:          uuid.New().String(),
		Level:       InitialState,
		PhoneNumber: phoneNumber,
		SessionID:   SessionID,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: DateOfBirth,
		PIN:         PIN,
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
			name: "Happy_case:optout",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: RegOptOutInput,
			},
			want: "CON We have successfully opted you\r\n" +
				"out of marketing messages\r\n" +
				"0. Go back home",
		},

		{
			name: "Sad_case:optout",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "12",
			},
			want: "CON Invalid choice. Try again.\r\n" +
				"1. Register\r\n" +
				"2. Opt Out\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy_case:optout" {
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}

				fakeCrm.OptOutFn = func(ctx context.Context, phoneNumber string) (*hubspotDomain.CRMContact, error) {
					return &hubspotDomain.CRMContact{ContactID: uuid.NewString()}, nil
				}

				fakeRepo.SaveUSSDEventFn = func(ctx context.Context, input *dto.USSDEvent) (*dto.USSDEvent, error) {
					return &dto.USSDEvent{
						SessionID: uuid.NewString(),
					}, nil
				}
			}

			if tt.name == "Sad_case:optout" {
				fakeCrm.OptOutFn = func(ctx context.Context, phoneNumber string) (*hubspotDomain.CRMContact, error) {
					return nil, fmt.Errorf("an error occurred %w", err)
				}
			}

			updatedSession, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occurred %v", err)
				return
			}

			if gotresp := u.AITUSSD.HandleUserRegistration(tt.args.ctx, updatedSession, tt.args.userResponse); gotresp != tt.want {
				t.Errorf("Impl.HandleUserRegistration() = %v, want %v", gotresp, tt.want)
			}

		})
	}
}
