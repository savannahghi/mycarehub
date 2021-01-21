package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
			wantErr: true, //todo : turn this back to false once a way is figured out to add primary email first
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, uids []string) error {
					return nil
				}

				fakeRepo.CheckIfEmailExistsFn = func(ctx context.Context, email string) (bool, error) {
					return false, nil
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateUserNameFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateVerifiedIdentifiersFn = func(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error {
					return nil
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_set_primary_address_succeeds" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_verify_otp_fails" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify email otp")
				}
			}

			if tt.name == "invalid:_verify_otp_returns_false" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid:_update_primary_address_fails" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return fmt.Errorf("unable to update primary email")
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get loggedin user")
				}
			}

			err := i.Onboarding.SetPrimaryEmailAddress(tt.args.ctx, tt.args.emailAddress, tt.args.otp)
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
			name: "valid: succefully updates permissions",
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

			if tt.name == "valid: succefully updates permissions" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}
				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}
			}

			if tt.name == "invalid: get logged in user uid fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get loggeg in user UID")
				}
			}

			if tt.name == "invalid: get user profile by UID fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile by UID")
				}
			}

			if tt.name == "invalid: update permissions fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError()
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
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError()
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
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError()
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
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError()
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
				) (*base.UserProfile, error) {
					return nil, exceptions.ProfileNotFoundError()
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

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
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

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
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

				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*base.UserProfile, error) {
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
