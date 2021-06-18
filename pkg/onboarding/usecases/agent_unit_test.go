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

func TestAgentUseCaseImpl_RegisterAgent(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	// agent 47
	fName := "Tobias"
	lName := "Rieper"
	dob := base.Date{
		Year:  1995,
		Month: 6,
		Day:   1,
	}
	agent := dto.RegisterAgentInput{
		FirstName:   fName,
		LastName:    lName,
		Gender:      base.GenderMale,
		PhoneNumber: base.TestUserEmail,
		Email:       base.TestUserEmail,
		DateOfBirth: dob,
	}

	type args struct {
		ctx   context.Context
		input dto.RegisterAgentInput
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid:register_new_agent",
			args: args{
				ctx:   ctx,
				input: agent,
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
				Role: base.RoleTypeAgent,
			},
			wantErr: false,
		},
		{
			name: "invalid:cannot_create_user_profile",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_create_customer_profile",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_create_supplier_profile",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_set_communication_settings",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_notify_new_agent_sms",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_notify_new_agent_email",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:get_logged_in_user_error",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:get_profile_by_uid_error",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:invalid_logged_in_user_role",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:normalizing_phonenumber_failed",
			args: args{
				ctx:   ctx,
				input: agent,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:cannot_set_agent_temporary_pin",
			args: args{
				ctx:   ctx,
				input: agent,
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
						Permissions: base.DefaultAgentPermissions,
					}, fmt.Errorf("user do not have required permissions")
				}
			}

			if tt.name == "valid:register_new_agent" {
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
						Role: base.RoleTypeAgent,
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

				fakeEngagementSvs.SendSMSFn = func(phoneNumbers []string, message string) error {
					return nil
				}

				fakeEngagementSvs.SendMailFn = func(email string, message string, subject string) error {
					return nil
				}
			}

			if tt.name == "invalid:cannot_notify_new_agent_sms" {
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
						Role: base.RoleTypeAgent,
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

				fakeEngagementSvs.SendSMSFn = func(phoneNumbers []string, message string) error {
					return fmt.Errorf("cannot send notification sms")
				}

			}

			if tt.name == "invalid:cannot_set_agent_temporary_pin" {
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
						Role: base.RoleTypeAgent,
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

			if tt.name == "invalid:cannot_notify_new_agent_email" {
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
						Role: base.RoleTypeAgent,
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

				fakeEngagementSvs.SendSMSFn = func(phoneNumbers []string, message string) error {
					return nil
				}

				fakeEngagementSvs.SendMailFn = func(email string, message string, subject string) error {
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
						Role: base.RoleTypeAgent,
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
						Role: base.RoleTypeAgent,
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
						Role: base.RoleTypeAgent,
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

			got, err := i.Agent.RegisterAgent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AgentUseCaseImpl.RegisterAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AgentUseCaseImpl.RegisterAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
