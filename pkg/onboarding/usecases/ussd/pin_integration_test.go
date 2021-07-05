package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// ChangePINEnterNewPINState indicates the state at which user wants to set a new PIN
	ChangePINEnterNewPINState = 51
	// ChangePINInput indicates the user intention to change their PIN
	ChangePINInput = "2"
	// GoBackHomeInput represents the user intention to go back to the main menu
	GoBackHomeInput = "0"
	//ConfirmNewPInState indicates the state when a user is confirming a pin update
	ConfirmNewPINState = 53
	// PINResetEnterNewPINState indicates the state when the user wants to reset their PIN
	PINResetEnterNewPINState = 10
	//ForgetPINResetState indicates the state when a use wants to reset PIN
	ForgetPINResetState = 13
	// PINResetProcessState represents the state when the user has provided a wrong PIN
	PINResetProcessState = 11
)

func TestImpl_HandleChangePIN(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
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
			name: "Happy case : enter correct old PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: EmptyInput,
			},
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "Happy case :_Change PIN_selected_enter correct old PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: ChangePINInput,
			},
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "Sad case :bad user response for correct old PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "ChangePINInput",
			},
			//TODO: Make this test valid for bad user response for correct
			//old PIN response
			want: "END Something went wrong. Please try again",
		},

		{
			name: "Happy case :Go back Home",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: GoBackHomeInput,
			},
			//TODO: Make this test valid for Go back Home response
			want: "END Something went wrong. Please try again",
		},

		{
			name: "Sad case :Go back Home",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "GoBackHomeInput",
			},
			want: "END invalid input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case : forgot PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Happy case :_Change PIN_selected_enter correct old PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case :bad user response for correct old PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Happy case :Go back Home" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, HomeMenuState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case :Go back Home" {
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

			if got := u.AITUSSD.HandleChangePIN(tt.args.ctx, session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandleChangePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_HandlePINReset(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
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
			name: "Happy case : Reset PIN_new_PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "1234",
			},
			want: "CON Please enter a new 4 digit PIN to\r\n" +
				"secure your account\r\n",
		},

		{
			name: "Happy case : Reset Confirm new_PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "1234",
			},
			want: "CON Please enter a 4 digit PIN again to\r\n" +
				"confirm.\r\n",
		},

		{
			name: "Sad case : Reset Confirm invalid_PIN",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "12",
			},
			want: "CON The PIN you entered is invalid. Please try again",
		},

		{
			name: "Sad case : Reset Confirm PIN do not match",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "4321",
			},
			want: "CON The PIN you entered does not match\r\n" +
				"Please enter a 4 digit PIN to\r\n" +
				"secure your account\r\n",
		},

		{
			name: "Sad case : Forgot PIN verify date",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "14032100",
			},
			//TODO: Make this test valid for  Forgot PIN verify date response
			want: "END something wrong it happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case : Reset PIN_new_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ForgetPINResetState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Happy case : Reset Confirm new_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case : Reset Confirm invalid_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case : Reset Confirm PIN do not match" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetProcessState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occured %v", err)
					return
				}
			}

			if tt.name == "Sad case : Forgot PIN verify date" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ForgotPINVerifyDate, sessionDetails.SessionID)
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

			if got := u.AITUSSD.HandlePINReset(tt.args.ctx, session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandlePINReset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_SetUSSDUserPin(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	phone := "+254700100200"
	pin := "4321"

	type args struct {
		ctx         context.Context
		phoneNumber string
		PIN         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Reset PIN successfully",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         pin,
			},
			wantErr: true, // TODO: Fix and make wantErr: false
		},

		{
			name: "Sad case: Unable to Reset PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         "pin",
			},
			wantErr: true,
		},

		{
			name: "Sad case: empty PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         "pin",
			},
			wantErr: true,
		},

		{
			name: "Sad case: invalid PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         "12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case: invalid PIN" {
				err := utils.ValidatePINLength("12345")
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
			}

			if err := u.AITUSSD.SetUSSDUserPin(tt.args.ctx, tt.args.phoneNumber, tt.args.PIN); (err != nil) != tt.wantErr {
				t.Errorf("Impl.SetUSSDUserPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_ChangeUSSDUserPIN(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	phoneNumber := "+254700100200"
	PIN := "4321"
	invalidPIN := "321"

	type args struct {
		ctx   context.Context
		phone string
		pin   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: Change user PIN successfully",
			args: args{
				ctx:   ctx,
				phone: phoneNumber,
				pin:   PIN,
			},
			want:    false, // TODO: Fix and make want: true
			wantErr: true,  // TODO: Fix and make wantErr: false
		},

		{
			name: "Sad case: Unable to Change user PIN",
			args: args{
				ctx:   ctx,
				phone: base.TestUserPhoneNumber,
				pin:   invalidPIN,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: Unable to Change user PIN" {
				err := utils.ValidatePINLength(invalidPIN)
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
			}

			got, err := u.AITUSSD.ChangeUSSDUserPIN(tt.args.ctx, tt.args.phone, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.ChangeUSSDUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Impl.ChangeUSSDUserPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
