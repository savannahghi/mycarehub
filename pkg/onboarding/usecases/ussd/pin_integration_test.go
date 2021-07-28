package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/scalarutils"
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
	// UserPINResetState represents workflows required to reset a forgotten user PIN
	UserPINResetState = 10
)

func TestImpl_HandleChangePIN_IntegrationTest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	phoneNumber := "+254750100200"
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
		DateOfBirth: scalarutils.Date{},
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
		t.Errorf("an error occurred %v", err)
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
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "Happy case :_Change PIN_selected",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: ChangePINInput,
			},
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "Sad case :bad choice selected",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "90",
			},
			want: "CON Invalid choice. Please try again.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case : empty input" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Happy case :_Change PIN_selected" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ChangePINEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Sad case :bad choice selected" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, HomeMenuState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			session, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occurred %v", err)
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

	phoneNumber := "+254733333333"
	PIN := "1234"
	FirstName := gofakeit.LastName()
	LastName := gofakeit.LastName()
	SessionID := uuid.New().String()
	Level := 15

	ussdDet := &domain.USSDLeadDetails{
		ID:          gofakeit.UUID(),
		Level:       Level,
		PhoneNumber: phoneNumber,
		SessionID:   SessionID,
		FirstName:   FirstName,
		LastName:    LastName,
		DateOfBirth: scalarutils.Date{},
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
		t.Errorf("an error occurred %v", err)
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
				userResponse: "00",
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
				userResponse: "14032000",
			},
			want: "CON Date of birth entered does not match the date of birth on record. Please enter your valid date of birth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case : Reset PIN_new_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, ForgetPINResetState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Happy case : Reset Confirm new_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Sad case : Reset Confirm invalid_PIN" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetEnterNewPINState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Sad case : Reset Confirm PIN do not match" {
				err = u.AITUSSD.UpdateSessionLevel(ctx, PINResetProcessState, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

			if tt.name == "Sad case : Forgot PIN verify date" {

				user, err := u.AITUSSD.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

				userProfile, err := u.AITUSSD.CreateUserProfile(ctx, user.PhoneNumber, user.UID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

				_, err = u.AITUSSD.SetUserPIN(ctx, ussdDet.PIN, userProfile.ID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

				_, err = u.AITUSSD.GetUserProfileByPrimaryPhoneNumber(ctx, *userProfile.PrimaryPhone, false)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
				err = u.AITUSSD.UpdateSessionLevel(ctx, ForgotPINVerifyDate, sessionDetails.SessionID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

			}

			session, err := u.AITUSSD.GetOrCreateSessionState(ctx, sessionDet)
			if err != nil {
				t.Errorf("an error occurred %v", err)
				return
			}

			if got := u.AITUSSD.HandlePINReset(tt.args.ctx, session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandlePINReset() = %v, want %v", got, tt.want)
			}

		})
	}
	err = u.AITUSSD.RemoveUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		t.Errorf("removing user: an error occurred %v", err)
		return
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
			wantErr: false,
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
			if tt.name == "Happy case: Reset PIN successfully" {
				err := utils.ValidatePINLength("12345")
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}

				userResponse, err := u.AITUSSD.GetOrCreatePhoneNumberUser(ctx, phone)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

				userProfile, err := u.AITUSSD.CreateUserProfile(ctx, userResponse.PhoneNumber, userResponse.UID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

				_, err = u.AITUSSD.GetUserProfileByPrimaryPhoneNumber(ctx, *userProfile.PrimaryPhone, false)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}

			}

			if tt.name == "Sad case: Unable to Reset PIN" {
				err := utils.ValidatePINLength("12345")
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
			}

			if tt.name == "Sad case: empty PIN" {
				err := utils.ValidatePINLength("12345")
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
			}

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
			want:    true,
			wantErr: false,
		},

		{
			name: "Sad case: Unable to Change user PIN",
			args: args{
				ctx:   ctx,
				phone: interserviceclient.TestUserPhoneNumber,
				pin:   invalidPIN,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case: Change user PIN successfully" {
				err := utils.ValidatePINLength(invalidPIN)
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
				userResp, err := u.AITUSSD.GetOrCreatePhoneNumberUser(ctx, phoneNumber)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
				userprofile, err := u.AITUSSD.CreateUserProfile(ctx, phoneNumber, userResp.UID)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
				_, err = u.AITUSSD.GetUserProfileByPrimaryPhoneNumber(ctx, *userprofile.PrimaryPhone, false)
				if err != nil {
					t.Errorf("an error occurred %v", err)
					return
				}
			}

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
