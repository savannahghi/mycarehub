package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

func TestProfileUseCaseImpl_UpdateVerifiedUIDS(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_profile_uids",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: false,
		},

		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_profile_uids" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.UpdateVerifiedUIDSFn = func(ctx context.Context, id string, uids []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateVerifiedUIDS(tt.args.ctx, tt.args.uids)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_UpdateSecondaryEmailAddresses(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx            context.Context
		emailAddresses []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_profile_secondary_email",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: false,
		},
		{
			name: "invalid:_update_profile_secondary_email", // no primary email
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_profile_secondary_email" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := base.TestUserEmail
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &email,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, uids []string) error {
					return nil
				}

				fakeRepo.CheckIfEmailExistsFn = func(ctx context.Context, email string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:_update_profile_secondary_email" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, uids []string) error {
					return fmt.Errorf("unable to update secondary email")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateSecondaryEmailAddresses(tt.args.ctx, tt.args.emailAddresses)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateUserName(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_name_succeeds",
			args: args{
				ctx:      ctx,
				userName: "kamau",
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:      ctx,
				userName: "mwas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_name_succeeds" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateUserNameFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged user")
				}
			}
			err := i.Onboarding.UpdateUserName(tt.args.ctx, tt.args.userName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateVerifiedIdentifiers(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		identifiers []base.VerifiedIdentifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_name_succeeds",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "a4f39af7-5b64-4c2f-91bd-42b3af315a5h",
						LoginProvider: "Facebook",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "j4f39af7-5b64-4c2f-91bd-42b3af225a5c",
						LoginProvider: "Phone",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "p4f39af7-5b64-4c2f-91bd-42b3af315a5c",
						LoginProvider: "Google",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_name_succeeds" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateVerifiedIdentifiersFn = func(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateVerifiedIdentifiers(tt.args.ctx, tt.args.identifiers)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePrimaryEmailAddress(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	primaryEmail := "me@gmail.com"

	type args struct {
		ctx          context.Context
		emailAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_email_succeeds",
			args: args{
				ctx:          ctx,
				emailAddress: primaryEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:          ctx,
				emailAddress: "kalulu@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_primary_email_address",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_secondary_email_address",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_update_email_succeeds" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_update_primary_email_address" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return fmt.Errorf("unable to update primary address")
				}
			}

			if tt.name == "invalid:_unable_to_update_secondary_email_address" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
						SecondaryEmailAddresses: []string{
							"", "lulu@gmail.com",
						},
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return fmt.Errorf("unable to update secondary email")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdatePrimaryEmailAddress(tt.args.ctx, tt.args.emailAddress)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_SetPrimaryEmailAddress(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	primaryEmail := "juha@gmail.com"

	type args struct {
		ctx          context.Context
		emailAddress string
		otp          string
		UID          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_set_primary_address_succeeds",
			args: args{
				ctx:          ctx,
				emailAddress: primaryEmail,
				otp:          "689552",
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_to_get_logged_in_uid",
			args: args{
				ctx:          ctx,
				emailAddress: "kichwa@gmail.com",
				otp:          "453852",
			},
			wantErr: true,
		},
		{
			name: "invalid:_verify_otp_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "kichwa@gmail.com",
				otp:          "453852",
			},
			wantErr: true,
		},
		{
			name: "invalid:_verify_otp_returns_false",
			args: args{
				ctx:          ctx,
				emailAddress: "kalu@gmail.com",
				otp:          "235789",
			},
			wantErr: true,
		},
		{
			name: "invalid:_update_primary_address_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "mwendwapole@gmail.com",
				otp:          "897523",
			},
			wantErr: true,
		},
		{
			name: "invalid:_resolving_the_consumer_nudge_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "mwendwapole@gmail.com",
				otp:          "897523",
			},
			wantErr: false,
		},
		{
			name: "invalid:_resolving_the_pro_nudge_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "mwendwapole@gmail.com",
				otp:          "897523",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_set_primary_address_succeeds" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}

				// Resolve the second nudge
				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}
			}

			if tt.name == "invalid:_failed_to_get_logged_in_uid" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error has occurred")
				}
			}

			if tt.name == "invalid:_verify_otp_fails" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify email otp")
				}
			}

			if tt.name == "invalid:_verify_otp_returns_false" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:_update_primary_address_fails" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return fmt.Errorf("unable to update primary email")
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get loggedin user")
				}
			}

			if tt.name == "invalid:_resolving_the_consumer_nudge_fails" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "invalid:_resolving_the_pro_nudge_fails" {
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeEngagementSvs.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}

				// Resolve the second nudge
				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return fmt.Errorf("an error occurred")
				}
			}

			err := i.Onboarding.SetPrimaryEmailAddress(
				tt.args.ctx,
				tt.args.emailAddress,
				tt.args.otp,
			)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_UpdatePermissions(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx   context.Context
		perms []base.PermissionType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: successfully updates permissions",
			args: args{
				ctx:   ctx,
				perms: []base.PermissionType{base.PermissionTypeSuperAdmin},
			},
			wantErr: false,
		},
		{
			name: "invalid: get logged in user uid fails",
			args: args{
				ctx:   ctx,
				perms: []base.PermissionType{base.PermissionTypeSuperAdmin},
			},
			wantErr: true,
		},
		{
			name: "invalid: get user profile by UID fails",
			args: args{
				ctx:   ctx,
				perms: []base.PermissionType{base.PermissionTypeSuperAdmin},
			},
			wantErr: true,
		},
		{
			name: "invalid: update permissions fails",
			args: args{
				ctx:   ctx,
				perms: []base.PermissionType{base.PermissionTypeSuperAdmin},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid: successfully updates permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}
				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}
			}

			if tt.name == "invalid: get logged in user uid fails" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("failed to get loggeg in user UID")
				}
			}

			if tt.name == "invalid: get user profile by UID fails" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile by UID")
				}
			}

			if tt.name == "invalid: update permissions fails" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}
				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return fmt.Errorf("unable to update permissions")
				}
			}

			err := i.Onboarding.UpdatePermissions(tt.args.ctx, tt.args.perms)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_AddRoleToUser(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx   context.Context
		phone string
		role  base.RoleType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: successfully updates role",
			args: args{
				ctx:   ctx,
				phone: "+254721123123",
				role:  base.RoleTypeEmployee,
			},
			wantErr: false,
		},
		{
			name: "invalid: get profile by primary phone number failed",
			args: args{
				ctx:   ctx,
				phone: "+254721123123",
				role:  base.RoleTypeEmployee,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: successfully updates role" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}

				fakeRepo.UpdateRoleFn = func(ctx context.Context, id string, role base.RoleType) error {
					return nil
				}
			}

			if tt.name == "invalid: get profile by primary phone number failed" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeBaseExt.GetUserProfileByPrimaryPhoneNumberFn = func(ctx context.Context, phone string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("UserProfile matching PhoneNumber not found")
				}
				fakeRepo.UpdateRoleFn = func(ctx context.Context, id string, role base.RoleType) error {
					return fmt.Errorf("User Roles not updated")
				}
			}

			err := i.Onboarding.AddRoleToUser(tt.args.ctx, tt.args.phone, tt.args.role)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_GetUserProfileAttributes(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx       context.Context
		UIDs      []string
		attribute string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid:_get_user_profile_emails",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.EmailsAttribute,
			},
			wantErr: false,
		},
		{
			name: "valid:_get_user_profile_phone_numbers",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.PhoneNumbersAttribute,
			},
			wantErr: false,
		},
		{
			name: "valid:_get_user_profile_fcm_tokens",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.FCMTokensAttribute,
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_get_user_profile_attribute",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: "not-an-attribute",
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_get_user_profile",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.FCMTokensAttribute,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_get_user_profile_emails" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					return &base.UserProfile{
						PrimaryEmailAddress: &email,
						SecondaryEmailAddresses: []string{
							base.GenerateRandomEmail(),
						},
					}, nil
				}
			}

			if tt.name == "valid:_get_user_profile_phone_numbers" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						PrimaryPhone:          &phone,
						SecondaryPhoneNumbers: []string{"+254700000000"},
					}, nil
				}
			}

			if tt.name == "valid:_get_user_profile_fcm_tokens" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						PushTokens: []string{uuid.New().String()},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						PrimaryEmailAddress: &email,
						SecondaryEmailAddresses: []string{
							base.GenerateRandomEmail(),
						},
						PrimaryPhone:          &phone,
						SecondaryPhoneNumbers: []string{"+254700000000"},
						PushTokens:            []string{uuid.New().String()},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
				}
			}

			attribute, err := i.Onboarding.GetUserProfileAttributes(
				tt.args.ctx,
				tt.args.UIDs,
				tt.args.attribute,
			)

			if tt.wantErr && attribute != nil {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr && attribute == nil {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_ConfirmedEmailAddresses(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx  context.Context
		UIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_get_confirmed_emails",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_get_user_profile",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_get_confirmed_emails" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					return &base.UserProfile{
						PrimaryEmailAddress: &email,
						SecondaryEmailAddresses: []string{
							base.GenerateRandomEmail(),
						},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
				}
			}

			confirmedEmails, err := i.Onboarding.ConfirmedEmailAddresses(
				tt.args.ctx,
				tt.args.UIDs,
			)
			if tt.wantErr && confirmedEmails != nil {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr && confirmedEmails == nil {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_ConfirmedPhoneNumbers(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx  context.Context
		UIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_get_confirmed_emails",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_get_user_profile",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_get_confirmed_emails" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						PrimaryPhone:          &phone,
						SecondaryPhoneNumbers: []string{"+254700000000"},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
				}
			}

			confirmedEmails, err := i.Onboarding.ConfirmedPhoneNumbers(
				tt.args.ctx,
				tt.args.UIDs,
			)
			if tt.wantErr && confirmedEmails != nil {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr && confirmedEmails == nil {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_validFCM(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx  context.Context
		UIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_valid_fcm_tokens",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_get_user_profile",
			args: args{
				ctx:  ctx,
				UIDs: []string{uuid.New().String()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_valid_fcm_tokens" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						PushTokens: []string{uuid.New().String()},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
				}
			}

			validFCM, err := i.Onboarding.ValidFCMTokens(
				tt.args.ctx,
				tt.args.UIDs,
			)
			if tt.wantErr && validFCM != nil {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr && validFCM == nil {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_ProfileAttributes(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx       context.Context
		UIDs      []string
		attribute string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid:_get_user_profile_emails",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.EmailsAttribute,
			},
			wantErr: false,
		},
		{
			name: "valid:_get_user_profile_phone_numbers",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.PhoneNumbersAttribute,
			},
			wantErr: false,
		},
		{
			name: "valid:_get_user_profile_fcm_tokens",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.FCMTokensAttribute,
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_get_user_profile_attribute",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: "not-an-attribute",
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_get_user_profile",
			args: args{
				ctx:       ctx,
				UIDs:      []string{uuid.New().String()},
				attribute: usecases.FCMTokensAttribute,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_get_user_profile_emails" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					return &base.UserProfile{
						PrimaryEmailAddress: &email,
						SecondaryEmailAddresses: []string{
							base.GenerateRandomEmail(),
						},
					}, nil
				}
			}

			if tt.name == "valid:_get_user_profile_phone_numbers" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						PrimaryPhone:          &phone,
						SecondaryPhoneNumbers: []string{"+254700000000"},
					}, nil
				}
			}

			if tt.name == "valid:_get_user_profile_fcm_tokens" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						PushTokens: []string{uuid.New().String()},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						PrimaryEmailAddress: &email,
						SecondaryEmailAddresses: []string{
							base.GenerateRandomEmail(),
						},
						PrimaryPhone:          &phone,
						SecondaryPhoneNumbers: []string{"+254700000000"},
						PushTokens:            []string{uuid.New().String()},
					}, nil
				}
			}

			if tt.name == "invalid:_failed_get_user_profile" {
				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError(fmt.Errorf("user profile not found"))
				}
			}

			attribute, err := i.Onboarding.ProfileAttributes(
				tt.args.ctx,
				tt.args.UIDs,
				tt.args.attribute,
			)

			if tt.wantErr && attribute != nil {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr && attribute == nil {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateSuspended(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx        context.Context
		status     bool
		phone      string
		useContext bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_suspend_with_use_context_false",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: false,
			},
			wantErr: false,
		},
		{
			name: "invalid:_suspend_with_use_context_false_update_user_fails",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: false,
			},
			wantErr: true,
		},
		{
			name: "valid:_suspend_with_use_context_true",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: true,
			},
			wantErr: false,
		},
		{
			name: "invalid:_suspend_with_use_context_true_get_user_profile_fails",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: true,
			},
			wantErr: true,
		},
		{
			name: "invalid:_suspend_with_use_context_true_get_loggedin_user_fails",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: true,
			},
			wantErr: true,
		},
		{
			name: "invalid:_normalize_msisdn_fails",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: false,
			},
			wantErr: true,
		},
		{
			name: "invalid:_suspend_use_context_false_get_user_profile_fails",
			args: args{
				ctx:        ctx,
				status:     true,
				phone:      "0721152896",
				useContext: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_suspend_with_use_context_false" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}

				fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
					return nil
				}
			}

			if tt.name == "invalid:_suspend_with_use_context_false_update_user_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "123",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0721521456", "0721856741",
						},
					}, nil
				}

				fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
					return fmt.Errorf("unable to update user profile")
				}
			}

			if tt.name == "valid:_suspend_with_use_context_true" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
					return nil
				}
			}

			if tt.name == "invalid:_suspend_with_use_context_true_get_loggedin_user_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get loggedin user")
				}

			}

			if tt.name == "invalid:_suspend_with_use_context_true_get_user_profile_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get userprofile")
				}

			}

			if tt.name == "invalid:_normalize_msisdn_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize phone")
				}
			}

			if tt.name == "invalid:_suspend_use_context_false_get_user_profile_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}

			err := i.Onboarding.UpdateSuspended(
				tt.args.ctx,
				tt.args.status,
				tt.args.phone,
				tt.args.useContext,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_UpdatePrimaryPhoneNumber(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	primaryPhone := "+254719789543"
	primaryPhone1 := "+254765739201"
	type args struct {
		ctx        context.Context
		phone      string
		useContext bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_primaryPhoneNumber_with_useContext_false",
			args: args{
				ctx:        ctx,
				phone:      primaryPhone,
				useContext: false,
			},
			wantErr: false,
		},

		{
			name: "valid:_update_primaryPhoneNumber_with_useContext_true",
			args: args{
				ctx:        ctx,
				phone:      primaryPhone1,
				useContext: true,
			},
			wantErr: false,
		},
		{
			name: "invalid:_missing_primaryPhoneNumber",
			args: args{
				ctx:        ctx,
				phone:      " ",
				useContext: true,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:        ctx,
				phone:      "+25463728192",
				useContext: true,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_userProfile_by_phonenumber",
			args: args{
				ctx:        ctx,
				phone:      "+254736291036",
				useContext: false,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:        ctx,
				phone:      "+25463728192",
				useContext: true,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_secondary_phonenumber",
			args: args{
				ctx:        ctx,
				phone:      "+25463728192",
				useContext: false,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_primary_phonenumber",
			args: args{
				ctx:        ctx,
				phone:      "+25463728192",
				useContext: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_update_primaryPhoneNumber_with_useContext_false" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254755889922"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "ABCDE",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}

				fakeRepo.UpdateSecondaryPhoneNumbersFn = func(ctx context.Context, id string, phoneNumbers []string) error {
					return nil
				}
			}

			if tt.name == "valid:_update_primaryPhoneNumber_with_useContext_true" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254789029156"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "f4f39af7--91bd-42b3af315a4e",
						PrimaryPhone: &primaryPhone1,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}

				fakeRepo.UpdateSecondaryPhoneNumbersFn = func(ctx context.Context, id string, phoneNumbers []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_missing_primaryPhoneNumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number provided")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254789029156"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "invalid:_unable_to_get_userProfile_by_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254799774466"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile by phonenumber")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_update_secondary_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254755889922"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "ABCDE",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}

				fakeRepo.UpdateSecondaryPhoneNumbersFn = func(ctx context.Context, id string, phoneNumbers []string) error {
					return fmt.Errorf("unable to update secondary phonenumber")
				}
			}

			if tt.name == "invalid:_unable_to_update_secondary_phonenumber" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254755889922"
					return &phone, nil
				}

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "ABCDE",
						PrimaryPhone: &phoneNumber,
						SecondaryPhoneNumbers: []string{
							"0765839203", "0789437282",
						},
					}, nil
				}

				fakeRepo.UpdatePrimaryPhoneNumberFn = func(ctx context.Context, id string, phoneNumber string) error {
					return fmt.Errorf("unable to update primary phonenumber")
				}

			}

			err := i.Onboarding.UpdatePrimaryPhoneNumber(tt.args.ctx, tt.args.phone, tt.args.useContext)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCase_UpdateBioData(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	dateOfBirth := base.Date{
		Day:   12,
		Year:  2000,
		Month: 2,
	}

	firstName := "Jatelo"
	lastName := "Omera"
	bioData := base.BioData{
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
	}

	var gender base.Gender = "female"
	updateGender := base.BioData{
		Gender: gender,
	}

	updateDOB := base.BioData{
		DateOfBirth: &dateOfBirth,
	}

	updateFirstName := base.BioData{
		FirstName: &firstName,
	}

	updateLastName := base.BioData{
		LastName: &lastName,
	}
	type args struct {
		ctx  context.Context
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: update primary biodata of a specific user profile",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: false,
		},
		{
			name: "valid: update primary biodata of a specific user profile - gender",
			args: args{
				ctx:  ctx,
				data: updateGender,
			},
			wantErr: false,
		},
		{
			name: "valid: update primary biodata of a specific user profile - DOB",
			args: args{
				ctx:  ctx,
				data: updateDOB,
			},
			wantErr: false,
		},
		{
			name: "valid: update primary biodata of a specific user profile - First Name",
			args: args{
				ctx:  ctx,
				data: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "valid: update primary biodata of a specific user profile - Last Name",
			args: args{
				ctx:  ctx,
				data: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "invalid: get logged in user uid fails",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: true,
		},
		{
			name: "invalid: get user profile by UID fails",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: true,
		},
		{
			name: "invalid: update primary biodata of a specific user profile",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: update primary biodata of a specific user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

			}
			if tt.name == "valid: update primary biodata of a specific user profile - gender" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

			}
			if tt.name == "valid: update primary biodata of a specific user profile - DOB" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

			}
			if tt.name == "valid: update primary biodata of a specific user profile - First Name" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

			}
			if tt.name == "valid: update primary biodata of a specific user profile - Last Name" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}

			}
			if tt.name == "invalid: get logged in user uid fails" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("failed to get loggeg in user UID")
				}
			}

			if tt.name == "invalid: get user profile by UID fails" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile by UID")
				}
			}
			if tt.name == "invalid: update primary biodata of a specific user profile" {

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return fmt.Errorf("failed update primary biodata of a user profile")
				}

			}

			err := i.Onboarding.UpdateBioData(tt.args.ctx, tt.args.data)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCase_CheckPhoneExists(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_check phone number exists",
			args: args{
				ctx:   ctx,
				phone: "+254711223344",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_check phone number exists - empty phone number provided",
			args: args{
				ctx:   ctx,
				phone: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_check phone number exists",
			args: args{
				ctx:   ctx,
				phone: "+254711223344",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_check phone number exists" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254711223344"
					return &phone, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "invalid:_check phone number exists - empty phone number provided" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("empty phone number provided")
				}
			}
			if tt.name == "invalid:_check phone number exists" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254711223344"
					return &phone, nil
				}
				fakeRepo.CheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("error checking if phone number exists")
				}
			}
			_, err := i.Onboarding.CheckPhoneExists(tt.args.ctx, tt.args.phone)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCase_CheckEmailExists(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validEmail := "me4@gmail.com"
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_check email exists",
			args: args{
				ctx:   ctx,
				email: validEmail,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_check email exists",
			args: args{
				ctx:   ctx,
				email: "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_check email exists" {
				fakeRepo.CheckIfEmailExistsFn = func(ctx context.Context, email string) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "invalid:_check email exists" {
				fakeRepo.CheckIfEmailExistsFn = func(ctx context.Context, email string) (bool, error) {
					return false, fmt.Errorf("failed to if email exists")
				}
			}
			_, err := i.Onboarding.CheckEmailExists(tt.args.ctx, tt.args.email)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePhotoUploadID(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx      context.Context
		uploadID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:successfully_updatePhotoUploadID",
			args: args{
				ctx:      ctx,
				uploadID: "some-upload-id",
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_update_photoUploadID",
			args: args{
				ctx:      ctx,
				uploadID: "some-upload-id",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_loggedInUser",
			args: args{
				ctx:      ctx,
				uploadID: "some-upload-id",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_userProfile",
			args: args{
				ctx:      ctx,
				uploadID: "some-upload-id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_updatePhotoUploadID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdatePhotoUploadIDFn = func(ctx context.Context, id string, uploadID string) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_get_loggedInUser" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "invalid:fail_to_get_userProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "invalid:fail_to_update_photoUploadID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdatePhotoUploadIDFn = func(ctx context.Context, id string, uploadID string) error {
					return fmt.Errorf("failed to update photo upload ID")
				}
			}
			err := i.Onboarding.UpdatePhotoUploadID(tt.args.ctx, tt.args.uploadID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdatePhotoUploadID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_AddAddress(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	addr := dto.UserAddressInput{
		Latitude:  1.2,
		Longitude: -34.001,
	}
	type args struct {
		ctx         context.Context
		input       dto.UserAddressInput
		addressType base.AddressType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:) add home address",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeHome,
			},
			wantErr: false,
		},
		{
			name: "happy:) add work address",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: false,
		},
		{
			name: "sad:( failed to get logged in user",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: true,
		},
		{
			name: "sad:( failed to get user profile",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: true,
		},
		{
			name: "sad:( failed to update user profile",
			args: args{
				ctx:         ctx,
				input:       addr,
				addressType: base.AddressTypeWork,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) add home address" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateAddressesFn = func(ctx context.Context, id string, address base.Address, addressType base.AddressType) error {
					return nil
				}
			}

			if tt.name == "happy:) add work address" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateAddressesFn = func(ctx context.Context, id string, address base.Address, addressType base.AddressType) error {
					return nil
				}
			}

			if tt.name == "sad:( failed to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			if tt.name == "sad:( failed to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			if tt.name == "sad:( failed to update user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateAddressesFn = func(ctx context.Context, id string, address base.Address, addressType base.AddressType) error {
					return fmt.Errorf("an error occurred")
				}
			}

			_, err := i.Onboarding.AddAddress(
				tt.args.ctx,
				tt.args.input,
				tt.args.addressType,
			)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_GetAddresses(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "happy:) get addresses",
			args:    args{ctx: ctx},
			wantErr: false,
		},
		{
			name:    "sad:( failed to get logged in user",
			args:    args{ctx: ctx},
			wantErr: true,
		},
		{
			name:    "sad:( failed to get user profile",
			args:    args{ctx: ctx},
			wantErr: true,
		},
		{
			name:    "sad:/ failed to get the home addresses",
			args:    args{ctx: ctx},
			wantErr: true,
		},
		{
			name:    "sad:/ failed to get the work addresses",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) get addresses" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						HomeAddress: &base.Address{
							Latitude:  "123",
							Longitude: "-1.2",
						},
						WorkAddress: &base.Address{
							Latitude:  "123",
							Longitude: "-1.2",
						},
					}, nil
				}
			}

			if tt.name == "sad:( failed to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			if tt.name == "sad:( failed to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			if tt.name == "sad:/ failed to get the home addresses" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:          uuid.New().String(),
						HomeAddress: &base.Address{},
						WorkAddress: &base.Address{
							Latitude:  "123",
							Longitude: "-1.2",
						},
					}, nil
				}
			}

			if tt.name == "sad:/ failed to get the work addresses" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						HomeAddress: &base.Address{
							Latitude:  "123",
							Longitude: "-1.2",
						},
						WorkAddress: &base.Address{},
					}, nil
				}
			}

			_, err := i.Onboarding.GetAddresses(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_GetUserCommunicationsSettings(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid: get comms settings",
			args:    args{ctx: ctx},
			wantErr: false,
		},
		{
			name:    "invalid: failed to get logged in user",
			args:    args{ctx: ctx},
			wantErr: true,
		},
		{
			name:    "invalid: failed to get user profile",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: get comms settings" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
						ID:            uuid.New().String(),
						AllowWhatsApp: true,
						AllowTextSMS:  true,
						AllowEmail:    true,
						AllowPush:     true,
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}
			}

			if tt.name == "invalid: failed to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			if tt.name == "invalid: failed to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			_, err := i.Onboarding.GetUserCommunicationsSettings(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetUserCommunicationsSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_SetUserCommunicationsSettings(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx           context.Context
		allowWhatsApp bool
		allowTextSms  bool
		allowPush     bool
		allowEmail    bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: set comms settings",
			args: args{
				ctx:           ctx,
				allowWhatsApp: true,
				allowTextSms:  true,
				allowPush:     true,
				allowEmail:    true,
			},
			wantErr: false,
		},
		{
			name: "invalid: failed to get logged in user",
			args: args{
				ctx:           ctx,
				allowWhatsApp: true,
				allowTextSms:  true,
				allowPush:     true,
				allowEmail:    true,
			},
			wantErr: true,
		},
		{
			name: "invalid: failed to get user profile",
			args: args{
				ctx:           ctx,
				allowWhatsApp: true,
				allowTextSms:  true,
				allowPush:     true,
				allowEmail:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: set comms settings" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string,
					allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
						ID:            uuid.New().String(),
						AllowWhatsApp: true,
						AllowTextSMS:  true,
						AllowEmail:    true,
						AllowPush:     true,
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}
			}

			if tt.name == "invalid: failed to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			if tt.name == "invalid: failed to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error ocurred")
				}
			}

			_, err := i.Onboarding.SetUserCommunicationsSettings(tt.args.ctx, &tt.args.allowWhatsApp,
				&tt.args.allowTextSms, &tt.args.allowEmail, &tt.args.allowPush)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.SetUserCommunicationsSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}
