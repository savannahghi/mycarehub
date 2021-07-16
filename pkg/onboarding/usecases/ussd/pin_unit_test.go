package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_HandleChangePIN_Unittest(t *testing.T) {
	ctx := context.Background()
	u, err := InitializeFakeUSSDTestService()
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
			name: "change pin",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: EmptyInput,
			},
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "go back home",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: GoBackHomeInput,
			},
			want: "CON Welcome to Be.Well\r\n" +
				"1. Opt out from marketing messages\r\n" +
				"2. Change PIN",
		},

		{
			name: "change pin option 2 selected",
			args: args{
				ctx:          ctx,
				session:      ussdDet,
				userResponse: "2",
			},
			want: "CON Enter your old PIN to continue\r\n" +
				"0. Go back home\r\n",
		},

		{
			name: "Sad case :invalid input",
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
			if tt.name == "change pin" {
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}
			}
			if tt.name == "go back home" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return &domain.PIN{}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return true
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}
			}
			if tt.name == "change pin option 2 selected" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return &domain.PIN{}, nil
				}
				fakePinExt.ComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return true
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}
			}

			if tt.name == "Sad case :invalid input" {
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("an error occured: %v", err)
				}
			}

			if got := u.AITUSSD.HandleChangePIN(tt.args.ctx, tt.args.session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandleChangePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_SetUSSDUserPin_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
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
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "pin", "pin"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: Unable to Reset PIN" {
				fakeRepo.SetUSSDUserPinFn = func(ctx context.Context, phoneNumber, PIN string) error {
					return err
				}
			}

			if tt.name == "Sad case: empty PIN" {
				fakeRepo.SetUSSDUserPinFn = func(ctx context.Context, phoneNumber, PIN string) error {
					return err
				}
			}

			if tt.name == "Sad case: invalid PIN" {
				err := utils.ValidatePINLength("12345")
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}
				fakeRepo.SetUSSDUserPinFn = func(ctx context.Context, phoneNumber, PIN string) error {
					return err
				}
			}

			if err := u.AITUSSD.SetUSSDUserPin(tt.args.ctx, tt.args.phoneNumber, tt.args.PIN); (err != nil) != tt.wantErr {
				t.Errorf("Impl.SetUSSDUserPin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_ChangeUSSDUserPIN_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
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
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "pin", "pin"
				}

				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: Unable to Change user PIN" {
				err := utils.ValidatePINLength(invalidPIN)
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))

					return
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.UpdatePINFn = func(ctx context.Context, id string, pin *domain.PIN) (bool, error) {
					return false, err
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
