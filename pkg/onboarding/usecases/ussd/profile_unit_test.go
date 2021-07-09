package ussd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_CreateUsddUserProfile_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	date := &base.Date{
		Year:  2000,
		Month: 10,
		Day:   01,
	}
	gender := base.GenderMale
	firstname := gofakeit.FirstName()
	lastname := gofakeit.LastName()
	phone := "+254700100200"
	pin := "4321"

	userP := &dto.UserProfileInput{
		DateOfBirth: date,
		Gender:      &gender,
		FirstName:   &firstname,
		LastName:    &lastname,
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
		PIN         string
		userProfile *dto.UserProfileInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         pin,
				userProfile: userP,
			},
			wantErr: false,
		},

		{
			name: "Sad case:_Non_digit_PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         "pin",
				userProfile: userP,
			},
			wantErr: false,
		},

		{
			name: "Sad case:_Longer_than_4_digit_PIN",
			args: args{
				ctx:         ctx,
				phoneNumber: phone,
				PIN:         "12345",
				userProfile: userP,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &dto.CreatedUserResponse{}, nil
				}

				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return &base.Customer{}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "pin", "pin"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case:_Non_digit_PIN" {
				err := utils.ValidatePIN(tt.args.PIN)
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("pin should be a valid number"))
					return
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return nil, err
				}
			}

			if tt.name == "Sad case:_Longer_than_4_digit_PIN" {
				err := utils.ValidatePIN(tt.args.PIN)
				if err != nil {
					exceptions.ValidatePINLengthError(fmt.Errorf("PIN should be of 4 digits"))
					return
				}

				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return nil, err
				}
			}

			if err := u.AITUSSD.CreateUsddUserProfile(tt.args.ctx, tt.args.phoneNumber, tt.args.PIN, tt.args.userProfile); (err != nil) != tt.wantErr {
				t.Errorf("Impl.CreateUsddUserProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
