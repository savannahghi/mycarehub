package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
)

func TestAdminUseCaseImpl_RegisterAdmin(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	// admin
	UID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
	id := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
	fName := "Tobias"
	lName := "Rieper"
	dob := scalarutils.Date{
		Year:  1995,
		Month: 6,
		Day:   1,
	}
	admin := dto.RegisterAdminInput{
		FirstName:   fName,
		LastName:    lName,
		Gender:      enumutils.GenderMale,
		PhoneNumber: firebasetools.TestUserEmail,
		Email:       firebasetools.TestUserEmail,
		DateOfBirth: dob,
	}

	type args struct {
		ctx   context.Context
		input dto.RegisterAdminInput
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "sad: unable to get logged in user",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to check user permissions",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: user do not have required permissions",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to get user profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to normalize phonenumber",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create detailed user profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create customer profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create supplier profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create admin profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create communication settings",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to create temporary pin",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad: unable to notify user",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy: registered new admin",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want: &profileutils.UserProfile{
				ID: id,
				UserBioData: profileutils.BioData{
					FirstName: &fName,
					LastName:  &lName,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad: unable to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("error unable to get logged in user profile")
				}
			}

			if tt.name == "sad: unable to check user permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("error unable to check user permissions")
				}
			}

			if tt.name == "sad: user do not have required permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad: unable to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to get user profile")
				}
			}

			if tt.name == "sad: unable to normalize phonenumber" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("error unable to normalize phone number")
				}
			}

			if tt.name == "sad: unable to create detailed user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to create user profile")
				}
			}

			if tt.name == "sad: unable to create customer profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return nil, fmt.Errorf("error unable to create customer profile")
				}
			}

			if tt.name == "sad: unable to create supplier profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return nil, fmt.Errorf("error unable to create supplier profile")
				}
			}

			if tt.name == "sad: unable to create admin profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}
				fakeRepo.CreateAdminProfileFn = func(ctx context.Context, adminProfile domain.AdminProfile) error {
					return fmt.Errorf("error unable to create admin profile")
				}
			}

			if tt.name == "sad: unable to create communication settings" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}
				fakeRepo.CreateAdminProfileFn = func(ctx context.Context, adminProfile domain.AdminProfile) error {
					return nil
				}
				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp, allowTextSms, allowPush, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return nil, fmt.Errorf("error unable to create communication settings")
				}
			}

			if tt.name == "sad: unable to create temporary pin" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}
				fakeRepo.CreateAdminProfileFn = func(ctx context.Context, adminProfile domain.AdminProfile) error {
					return nil
				}
				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp, allowTextSms, allowPush, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "nil", fmt.Errorf("error, unable to generate user pin")
				}
			}

			if tt.name == "sad: unable to notify user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}
				fakeRepo.CreateAdminProfileFn = func(ctx context.Context, adminProfile domain.AdminProfile) error {
					return nil
				}
				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp, allowTextSms, allowPush, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{}, nil
				}
				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "123", nil
				}
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "pin", "sha"
				}
				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
				fakeEngagementSvs.SendSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return fmt.Errorf("error unable to notify user")
				}
			}

			if tt.name == "happy: registered new admin" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{UID: UID}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phoneNumber := interserviceclient.TestUserPhoneNumber
					return &phoneNumber, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: id, UserBioData: profileutils.BioData{FirstName: &fName, LastName: &lName}}, nil
				}
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return &profileutils.Customer{}, nil
				}
				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}
				fakeRepo.CreateAdminProfileFn = func(ctx context.Context, adminProfile domain.AdminProfile) error {
					return nil
				}
				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp, allowTextSms, allowPush, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{}, nil
				}
				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "123", nil
				}
				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "pin", "sha"
				}
				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}
				fakeEngagementSvs.SendSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return nil
				}
			}

			got, err := i.Admin.RegisterAdmin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminUseCaseImpl.RegisterAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminUseCaseImpl.RegisterAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminUseCaseImpl_FetchAdmins(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.Admin
		wantErr bool
	}{
		{
			name: "success:_non_empty_list_of_user_admins",
			args: args{
				ctx: ctx,
			},
			want: []*dto.Admin{
				{
					ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					PrimaryPhone: interserviceclient.TestUserPhoneNumber,
					ResendPIN:    true,
					Roles:        []dto.RoleOutput{},
				},
				{
					ID:           "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					PrimaryPhone: interserviceclient.TestUserPhoneNumber,
					ResendPIN:    true,
					Roles:        []dto.RoleOutput{},
				},
			},
			wantErr: false,
		},
		{
			name: "success:_empty_list_of_user_admins",
			args: args{
				ctx: ctx,
			},
			want:    []*dto.Admin{},
			wantErr: false,
		},
		{
			name: "fail:error_fetching_list_of_user_admins",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success:_non_empty_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					p := interserviceclient.TestUserPhoneNumber
					s := []*profileutils.UserProfile{
						{
							ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
							PrimaryPhone: &p,
							VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
							Role:         profileutils.RoleTypeEmployee,
						},
						{
							ID:           "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
							PrimaryPhone: &p,
							VerifiedUIDS: []string{"c9d62c7e-93e5-44a6-b503-6fc159c1782f"},
							Role:         profileutils.RoleTypeEmployee,
						},
					}
					return s, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return &domain.PIN{IsOTP: true}, nil
				}

				fakeRepo.GetRolesByIDsFn = func(ctx context.Context, roleIDs []string) (*[]profileutils.Role, error) {
					roles := []profileutils.Role{}
					return &roles, nil
				}
			}
			if tt.name == "success:_empty_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					return []*profileutils.UserProfile{}, nil
				}

				fakeRepo.GetPINByProfileIDFn = func(ctx context.Context, ProfileID string) (*domain.PIN, error) {
					return &domain.PIN{}, nil
				}
			}
			if tt.name == "fail:error_fetching_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("cannot fetch list of user profiles")
				}
			}
			got, err := i.Admin.FetchAdmins(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminUseCaseImpl.FetchAdmins() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminUseCaseImpl.FetchAdmins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminUseCaseImpl_ActivateAdmin(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx   context.Context
		input dto.ProfileSuspensionInput
	}

	input := args{
		ctx: ctx,
		input: dto.ProfileSuspensionInput{
			ID:     "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
			Reason: "",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "sad unable to get logged in user",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to check user permissions",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad user do not have required permissions",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to get user profile by id",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to activate admin account",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "happy activated admin account",
			args:    input,
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if tt.name == "sad unable to get logged in user" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return nil, fmt.Errorf("unable to get logged in user")
			}
		}
		if tt.name == "sad unable to check user permissions" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return &dto.UserInfo{}, nil
			}
			fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
				return false, fmt.Errorf("unable to check user permissions")
			}
		}
		if tt.name == "sad user do not have required permissions" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return &dto.UserInfo{}, nil
			}
			fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
				return false, nil
			}
		}

		if tt.name == "sad unable to get user profile by id" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return &dto.UserInfo{}, nil
			}
			fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
				return false, nil
			}
			fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
				return nil, fmt.Errorf("unable to get user profile by id")
			}
		}
		if tt.name == "sad unable to activate admin account" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return &dto.UserInfo{}, nil
			}
			fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
				return true, nil
			}
			fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
				return &profileutils.UserProfile{}, nil
			}

			fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
				return fmt.Errorf("unable to activate account")
			}
		}
		if tt.name == "happy activated admin account" {
			fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
				return &dto.UserInfo{}, nil
			}
			fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
				return true, nil
			}
			fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
				return &profileutils.UserProfile{}, nil
			}
			fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
				return nil
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Admin.ActivateAdmin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminUseCaseImpl.ActivateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AdminUseCaseImpl.ActivateAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminUseCaseImpl_DeactivateAdmin(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}
	type args struct {
		ctx   context.Context
		input dto.ProfileSuspensionInput
	}

	input := args{
		ctx: ctx,
		input: dto.ProfileSuspensionInput{
			ID:     "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
			Reason: "",
		},
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "sad unable to get logged in user",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to check user permissions",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad user do not have required permissions",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to get user profile by id",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to deactivate admin account",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "happy deactivated admin account",
			args:    input,
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad unable to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}
			if tt.name == "sad unable to check user permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, fmt.Errorf("unable to check user permissions")
				}
			}
			if tt.name == "sad user do not have required permissions" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad unable to get user profile by id" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile by id")
				}
			}
			if tt.name == "sad unable to deactivate admin account" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}

				fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
					return fmt.Errorf("unable to deactivate account")
				}
			}
			if tt.name == "happy deactivated admin account" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{}, nil
				}
				fakeRepo.CheckIfUserHasPermissionFn = func(ctx context.Context, UID string, requiredPermission profileutils.Permission) (bool, error) {
					return true, nil
				}
				fakeRepo.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeRepo.UpdateSuspendedFn = func(ctx context.Context, id string, status bool) error {
					return nil
				}
			}
			got, err := i.Admin.DeactivateAdmin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminUseCaseImpl.DeactivateAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AdminUseCaseImpl.DeactivateAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminUseCaseImpl_FindAdminByNameOrPhone(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	nameOrPhone := "Test"
	fName := "Test"
	lName := "User"

	type args struct {
		ctx         context.Context
		nameOrPhone *string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "sad: unable get user profiles",
			args:    args{ctx: ctx, nameOrPhone: &nameOrPhone},
			want:    0,
			wantErr: true,
		},
		{
			name:    "sad: did not get any user",
			args:    args{ctx: ctx, nameOrPhone: &nameOrPhone},
			want:    0,
			wantErr: false,
		},
		{
			name:    "happy: got user profiles",
			args:    args{ctx: ctx, nameOrPhone: &nameOrPhone},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad: unable get user profiles" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error unable to get user profile by phone")
				}
			}

			if tt.name == "sad: did not get any user" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					return []*profileutils.UserProfile{}, nil
				}
			}

			if tt.name == "happy: got user profiles" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					phone := interserviceclient.TestUserPhoneNumber
					profile1 := profileutils.UserProfile{
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
						},
						PrimaryPhone: &phone,
					}
					profile2 := profileutils.UserProfile{
						UserBioData: profileutils.BioData{
							FirstName: &lName,
							LastName:  &lName,
						},
						PrimaryPhone: &phone,
					}
					return []*profileutils.UserProfile{
						&profile1,
						&profile2,
					}, nil
				}
			}
			got, err := i.Admin.FindAdminByNameOrPhone(tt.args.ctx, tt.args.nameOrPhone)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminUseCaseImpl.FindAdminByNameOrPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("AdminUseCaseImpl.FindAdminByNameOrPhone() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
