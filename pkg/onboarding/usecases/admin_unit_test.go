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
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	// admin 47
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
			name: "valid:register_new_admin",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want: &profileutils.UserProfile{
				ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
				UserBioData: profileutils.BioData{
					FirstName:   &fName,
					LastName:    &lName,
					Gender:      enumutils.GenderMale,
					DateOfBirth: &dob,
				},
				Role: profileutils.RoleTypeEmployee,
			},
			wantErr: false,
		},
		{
			name: "invalid:cannot_create_user_profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_create_customer_profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_create_supplier_profile",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_set_communication_settings",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_notify_new_admin_sms",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_notify_new_admin_email",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:get_logged_in_user_error",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:get_profile_by_uid_error",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:invalid_logged_in_user_role",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:normalizing_phonenumber_failed",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_set_admin_temporary_pin",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "invalid:normalizing_phonenumber_failed" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("cannot normalize the mobile number")
				}
			}

			if tt.name == "invalid:get_logged_in_user_error" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("cannot get logged in user")
				}
			}

			if tt.name == "invalid:get_profile_by_uid_error" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user bu UID")
				}
			}

			if tt.name == "invalid:invalid_logged_in_user_role" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Permissions: profileutils.DefaultAdminPermissions,
					}, fmt.Errorf("user do not have required permissions")
				}
			}

			if tt.name == "valid:register_new_admin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultSuperAdminPermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return &profileutils.Supplier{}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:        "4711a5e4-a211-4e2b-b40b-b1160049b984",
						ProfileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}

				fakeEngagementSvs.SendSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return nil
				}

				fakeEngagementSvs.SendMailFn = func(ctx context.Context, email string, message string, subject string) error {
					return nil
				}
			}

			if tt.name == "invalid:cannot_notify_new_admin_sms" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:        "4711a5e4-a211-4e2b-b40b-b1160049b984",
						ProfileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}

				fakeEngagementSvs.SendSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return fmt.Errorf("cannot send notification sms")
				}
			}

			if tt.name == "invalid:cannot_set_admin_temporary_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:        "4711a5e4-a211-4e2b-b40b-b1160049b984",
						ProfileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "1234", fmt.Errorf("cannot generate temporary PIN")
				}

			}

			if tt.name == "invalid:cannot_set_admin_temporary_pin" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:        "4711a5e4-a211-4e2b-b40b-b1160049b984",
						ProfileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "1234", fmt.Errorf("cannot generate temporary PIN")
				}

			}

			if tt.name == "invalid:cannot_notify_new_admin_email" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return &profileutils.UserCommunicationsSetting{
						ID:        "4711a5e4-a211-4e2b-b40b-b1160049b984",
						ProfileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}, nil
				}

				fakePinExt.GenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakePinExt.EncryptPINFn = func(rawPwd string, options *extension.Options) (string, string) {
					return "salt", "passw"
				}

				fakeRepo.SavePINFn = func(ctx context.Context, pin *domain.PIN) (bool, error) {
					return true, nil
				}

				fakeEngagementSvs.SendSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return nil
				}

				fakeEngagementSvs.SendMailFn = func(ctx context.Context, email string, message string, subject string) error {
					return fmt.Errorf("cannot send notification email")
				}
			}

			if tt.name == "invalid:cannot_set_communication_settings" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
					return nil, fmt.Errorf("")
				}
			}

			if tt.name == "invalid:cannot_create_supplier_profile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &profileutils.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier profileutils.Supplier) (*profileutils.Supplier, error) {
					return nil, fmt.Errorf("cannot create supplier profile")
				}
			}

			if tt.name == "invalid:cannot_create_customer_profile" {

				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      enumutils.GenderMale,
							DateOfBirth: &dob,
						},
						Role: profileutils.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*profileutils.Customer, error) {
					return nil, fmt.Errorf("cannot create customer profile")
				}
			}

			if tt.name == "invalid:cannot_create_user_profile" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254777886622"
					return &phone, nil
				}

				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: profileutils.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    enumutils.GenderMale,
						},
						Permissions: profileutils.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile profileutils.UserProfile) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("cannot create user profile")
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

	i, err := InitializeFakeOnboaridingInteractor()
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
					ID:                  "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					PrimaryPhone:        interserviceclient.TestUserPhoneNumber,
					PrimaryEmailAddress: firebasetools.TestUserEmail,
				},
				{
					ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					PrimaryPhone:        interserviceclient.TestUserPhoneNumber,
					PrimaryEmailAddress: firebasetools.TestUserEmail,
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
					e := firebasetools.TestUserEmail
					s := []*profileutils.UserProfile{
						{
							ID:                  "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
							PrimaryPhone:        &p,
							PrimaryEmailAddress: &e,
							VerifiedUIDS:        []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
							Role:                profileutils.RoleTypeEmployee,
						},
						{
							ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
							PrimaryPhone:        &p,
							PrimaryEmailAddress: &e,
							VerifiedUIDS:        []string{"c9d62c7e-93e5-44a6-b503-6fc159c1782f"},
							Role:                profileutils.RoleTypeEmployee,
						},
					}
					return s, nil
				}
			}
			if tt.name == "success:_empty_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role profileutils.RoleType) ([]*profileutils.UserProfile, error) {
					return []*profileutils.UserProfile{}, nil
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
