package ussd_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// InitialState ...
	InitialState = 0
	// GetFirstNameState ...
	GetFirstNameState = 1
	// GetLastNameState ...
	GetLastNameState = 2
	// GetDOBState ...
	GetDOBState = 3
	// GetPINState ...
	GetPINState = 4
	// SaveRecordState ...
	SaveRecordState = 5
	// RegisterInput ...
	RegisterInput = "1"
	//RegOptOutInput ...
	RegOptOutInput = "2"
	//RegChangePINInput ...
	RegChangePINInput = "2"
	//HomeMenuState represents inner submenu once user is logged in
	HomeMenuState = 5
)

func TestImpl_HandleUserRegistration(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service %v", err)
		return
	}

	phoneNumber := "+254700100200"
	dateOfBirth := "12122000"
	PIN := "1234"

	ContactType := "phone"
	ContactValue := phoneNumber
	FirstName := gofakeit.FirstName()
	LastName := gofakeit.LastName()
	DateOfBirth := base.Date{
		Day:   0,
		Month: 0,
		Year:  0,
	}
	IsSync := false
	TimeSync := time.Now()
	OptOut := "NO"
	WantCover := false
	ContactChannel := "USSD"
	IsRegistered := false

	contactLeadPayload := &dto.ContactLeadInput{
		ContactType:    ContactType,
		ContactValue:   ContactValue,
		FirstName:      FirstName,
		LastName:       LastName,
		DateOfBirth:    DateOfBirth,
		IsSync:         IsSync,
		TimeSync:       &TimeSync,
		OptOut:         CRMDomain.GeneralOptionType(OptOut),
		WantCover:      WantCover,
		ContactChannel: ContactChannel,
		IsRegistered:   IsRegistered,
	}

	SessionID := uuid.New().String()
	Level := 0
	Text := ""

	sessionDet := &dto.SessionDetails{
		SessionID:   SessionID,
		PhoneNumber: &phoneNumber,
		Level:       Level,
		Text:        Text,
	}

	//Staging
	err = u.AITUSSD.StageCRMPayload(ctx, contactLeadPayload)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	sessionDetails, err := u.AITUSSD.AddAITSessionDetails(ctx, sessionDet)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	err = u.Onboarding.SetOptOut(ctx, "STOP", sessionDetails.PhoneNumber)
	if err != nil {
		t.Errorf("an error occured %v", err)
		return
	}

	validUSSDLeadDetails := &domain.USSDLeadDetails{
		ID:          uuid.New().String(),
		Level:       InitialState,
		PhoneNumber: phoneNumber,
		SessionID:   SessionID,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: dateOfBirth,
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
			name: "Happy_case:_firstname",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: gofakeit.FirstName(),
			},
			want: "CON Please enter your lastname(e.g.\r\n" +
				"Doe)\r\n",
		},

		{
			name: "Happy_case:_lastname",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: gofakeit.LastName(),
			},
			want: "CON Please enter your date of birth in\r\n" +
				"DDMMYYYY format e.g 14031996 for\r\n" +
				"14th March 1996\r\n",
		},

		{
			name: "Happy_case:_dob",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "12122000",
			},
			want: "CON Please enter a 4 digit PIN to secure your account",
		},

		{
			name: "Happy_case:_pin_one",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "1234",
			},
			want: "CON Please enter a 4 digit PIN again to confirm",
		},

		//Bad cases

		{
			name: "Sad_case:_invalid_firstname",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "1234",
			},
			want: "CON Invalid name. Please enter a valid name (e.g John)",
		},

		{
			name: "Sad_case:_invalid_lastname",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "1234",
			},
			want: "CON Invalid name. Please enter a valid name (e.g Doe)",
		},

		{
			name: "Sad_case:_invalid_dob",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "hello",
			},
			want: "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996",
		},

		{
			name: "Sad case:_invalid_pin_one",
			args: args{
				ctx:          ctx,
				session:      validUSSDLeadDetails,
				userResponse: "hello",
			},
			want: "CON The PIN you entered in not correct please enter a 4 digit PIN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("test name is", tt.name)
			if tt.name == "Happy_case:_firstname" {
				//Get firstname state
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetFirstNameState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}

			}

			fmt.Println("test name is", tt.name)
			if tt.name == "Happy_case:_lastname" {
				//Get last state
				fmt.Println("updating level ")
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetLastNameState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}

			}

			if tt.name == "Happy_case:_dob" {
				//Get dob state
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetDOBState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Happy_case:_pin_one" {
				//Get pin state
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad_case:_invalid_firstname" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetFirstNameState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}

			}

			if tt.name == "Sad_case:_invalid_lastname" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetLastNameState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}

			}

			if tt.name == "Sad_case:_invalid_dob" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetDOBState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case:_invalid_pin_one" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, GetPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			updatedSession, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occured %v", err)
				return
			}

			if gotresp := u.AITUSSD.HandleUserRegistration(tt.args.ctx, updatedSession, tt.args.userResponse); gotresp != tt.want {
				t.Errorf("Impl.HandleUserRegistration() = %v, want %v", gotresp, tt.want)
			}

		})
	}
}
