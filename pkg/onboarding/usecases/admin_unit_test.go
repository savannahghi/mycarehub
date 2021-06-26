package usecases_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
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
	dob := base.Date{
		Year:  1995,
		Month: 6,
		Day:   1,
	}
	admin := dto.RegisterAdminInput{
		FirstName:   fName,
		LastName:    lName,
		Gender:      base.GenderMale,
		PhoneNumber: base.TestUserEmail,
		Email:       base.TestUserEmail,
		DateOfBirth: dob,
	}

	type args struct {
		ctx   context.Context
		input dto.RegisterAdminInput
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid:register_new_admin",
			args: args{
				ctx:   ctx,
				input: admin,
			},
			want: &base.UserProfile{
				ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
				UserBioData: base.BioData{
					FirstName:   &fName,
					LastName:    &lName,
					Gender:      base.GenderMale,
					DateOfBirth: &dob,
				},
				Role: base.RoleTypeEmployee,
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Permissions: base.DefaultAdminPermissions,
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultSuperAdminPermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					return &base.Supplier{}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}

				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
					return &base.UserCommunicationsSetting{
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Supplier{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.SetUserCommunicationsSettingsFn = func(ctx context.Context, profileID string, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"
					return &base.Customer{
						ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
						ProfileID: &prID,
					}, nil
				}

				fakeRepo.CreateDetailedSupplierProfileFn = func(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName:   &fName,
							LastName:    &lName,
							Gender:      base.GenderMale,
							DateOfBirth: &dob,
						},
						Role: base.RoleTypeEmployee,
					}, nil
				}

				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
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
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
						UserBioData: base.BioData{
							FirstName: &fName,
							LastName:  &lName,
							Gender:    base.GenderMale,
						},
						Permissions: base.DefaultEmployeePermissions,
					}, nil
				}
				fakeRepo.CreateDetailedUserProfileFn = func(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error) {
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
					PrimaryPhone:        base.TestUserPhoneNumber,
					PrimaryEmailAddress: base.TestUserEmail,
				},
				{
					ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					PrimaryPhone:        base.TestUserPhoneNumber,
					PrimaryEmailAddress: base.TestUserEmail,
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
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role base.RoleType) ([]*base.UserProfile, error) {
					p := base.TestUserPhoneNumber
					e := base.TestUserEmail
					s := []*base.UserProfile{
						{
							ID:                  "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
							PrimaryPhone:        &p,
							PrimaryEmailAddress: &e,
							VerifiedUIDS:        []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
							Role:                base.RoleTypeEmployee,
						},
						{
							ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
							PrimaryPhone:        &p,
							PrimaryEmailAddress: &e,
							VerifiedUIDS:        []string{"c9d62c7e-93e5-44a6-b503-6fc159c1782f"},
							Role:                base.RoleTypeEmployee,
						},
					}
					return s, nil
				}
			}
			if tt.name == "success:_empty_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role base.RoleType) ([]*base.UserProfile, error) {
					return []*base.UserProfile{}, nil
				}
			}
			if tt.name == "fail:error_fetching_list_of_user_admins" {
				fakeRepo.ListUserProfilesFn = func(ctx context.Context, role base.RoleType) ([]*base.UserProfile, error) {
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
