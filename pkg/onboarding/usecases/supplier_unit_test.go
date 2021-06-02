package usecases_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestProfileUseCaseImpl_FindSupplierByID(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_find_supplier_by_id",
			args: args{
				ctx: ctx,
				id:  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
			},
			wantErr: false,
		},
		{
			name: "invalid:_find_supplier_by_id_fails",
			args: args{
				ctx: ctx,
				id:  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_find_supplier_by_id" {
				fakeRepo.GetSupplierProfileByIDFn = func(ctx context.Context, id string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

			}

			if tt.name == "invalid:_find_supplier_by_id_fails" {
				fakeRepo.GetSupplierProfileByIDFn = func(ctx context.Context, id string) (*base.Supplier, error) {
					return nil, fmt.Errorf("unable to get supp;ier profile")
				}

			}

			sup, err := i.Supplier.FindSupplierByID(tt.args.ctx, tt.args.id)
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

				if sup.ID == "" {
					t.Errorf("empty ID returned %v", sup.ID)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_SendKYCEmail(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx          context.Context
		text         string
		emailaddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_send_kyc_mail",
			args: args{
				ctx:          ctx,
				text:         "Dear user this is a sample email",
				emailaddress: "kalulu@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "invalid:_send_mail_fails",
			args: args{
				ctx:          ctx,
				text:         "Dear user this is a sample email",
				emailaddress: "kalulu@gmail.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_send_kyc_mail" {
				fakeEngagementSvs.SendMailFn = func(email string, message string, subject string) error {
					return nil
				}

			}

			if tt.name == "invalid:_send_mail_fails" {
				fakeEngagementSvs.SendMailFn = func(email string, message string, subject string) error {
					return fmt.Errorf("unable to send mail")
				}

			}

			err := i.Supplier.SendKYCEmail(tt.args.ctx, tt.args.text, tt.args.emailaddress)
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

func TestProfileUseCaseImpl_PublishKYCFeedItem(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			name: "valid:_publish_kyc_item",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
			},
			wantErr: false,
		},
		{
			name: "invalid:_unexpected_status_code",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_publish_kyc_item",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_publish_kyc_item" {
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_publish_kyc_item" {
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return nil, fmt.Errorf("unable to publish kyc item")
				}
			}

			if tt.name == "invalid:_unexpected_status_code" {
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "",
						StatusCode: 400,
						Body:       respReader,
					}, nil
				}
			}

			err := i.Supplier.PublishKYCFeedItem(tt.args.ctx, tt.args.uids...)
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

func TestProfileUseCaseImpl_RetireKYCRequest(t *testing.T) {
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
			name: "valid:_retire_kyc_nudge",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_login",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:_remove_kyc_nudge_fails",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_supplier_profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_user_profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_retire_kyc_nudge" {

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "hj539af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.RemoveKYCProcessingRequestFn = func(ctx context.Context, supplierProfileID string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_login" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to log in")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get userprofile")
				}
			}

			if tt.name == "invalid:_unable_to_get_supplier_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("unable to get supplier profile")
				}
			}

			if tt.name == "invalid:_remove_kyc_nudge_fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "hj539af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.RemoveKYCProcessingRequestFn = func(ctx context.Context, supplierProfileID string) error {
					return fmt.Errorf("unable to retire nudge")
				}
			}

			err := i.Supplier.RetireKYCRequest(tt.args.ctx)
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

func TestProfileUseCaseImpl_ProcessKYCRequest(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	rejectionReason := "You can do better :("
	type args struct {
		ctx             context.Context
		id              string
		status          domain.KYCProcessStatus
		rejectionReason *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_approved_a_kyc_request",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusApproved,
			},
			wantErr: false,
		},
		{
			name: "valid:_rejected_a_kyc_request",
			args: args{
				ctx:             ctx,
				id:              uuid.New().String(),
				status:          domain.KYCProcessStatusRejected,
				rejectionReason: &rejectionReason,
			},
			wantErr: false,
		},
		{
			name: "invalid:_failed_to_get_process_kyc_request",
			args: args{
				ctx:             ctx,
				id:              uuid.New().String(),
				status:          domain.KYCProcessStatusRejected,
				rejectionReason: &rejectionReason,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_update_supplier_profile",
			args: args{
				ctx:             ctx,
				id:              uuid.New().String(),
				status:          domain.KYCProcessStatusRejected,
				rejectionReason: &rejectionReason,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_update_user_profile",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusApproved,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_get_supplier_profile_when_approved",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusApproved,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_update_supplier_profile_when_approved",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusApproved,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_send_email",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusRejected,
			},
			wantErr: true,
		},
		{
			name: "invalid:_failed_to_send_sms",
			args: args{
				ctx:    ctx,
				id:     uuid.New().String(),
				status: domain.KYCProcessStatusRejected,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_approved_a_kyc_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:          uuid.New().String(),
						Permissions: []base.PermissionType{base.PermissionTypeAdmin},
					}, nil
				}
				fakeRepo.CheckIfAdminFn = func(profile *base.UserProfile) bool {
					return true
				}

				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client,
				) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return "supplier.topic"
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(
					ctx context.Context,
					profileID string,
					data *base.Supplier,
				) error {
					return nil
				}

				fakeRepo.ActivateSupplierProfileFn = func(
					profileID string,
					supplier base.Supplier,
				) (*base.Supplier, error) {
					return &base.Supplier{}, nil
				}

				fakeEngagementSvs.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeEngagementSvs.SendSMSFn = func(
					phoneNumbers []string,
					message string,
				) error {
					return nil
				}
			}

			if tt.name == "valid:_rejected_a_kyc_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:          uuid.New().String(),
						Permissions: []base.PermissionType{base.PermissionTypeAdmin},
					}, nil
				}
				fakeRepo.CheckIfAdminFn = func(profile *base.UserProfile) bool {
					return true
				}

				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeEngagementSvs.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeEngagementSvs.SendSMSFn = func(
					phoneNumbers []string,
					message string,
				) error {
					return nil
				}
			}

			if tt.name == "invalid:_failed_to_get_process_kyc_request" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					return nil, fmt.Errorf("failed to get the request")
				}
			}

			if tt.name == "invalid:_failed_to_update_supplier_profile" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return fmt.Errorf("failed to update supplier profile")
				}
			}

			if tt.name == "invalid:_failed_to_update_user_profile" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "invalid:_failed_to_get_supplier_profile_when_approved" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client,
				) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakeEPRSvc.CreateERPSupplierFn = func(
					ctx context.Context,
					supplierPayload dto.SupplierPayload,
					UID string,
				) (*base.Supplier, error) {
					return &base.Supplier{}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier profile")
				}
			}

			if tt.name == "invalid:_failed_to_update_supplier_profile_when_approved" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client,
				) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakeEPRSvc.CreateERPSupplierFn = func(
					ctx context.Context,
					supplierPayload dto.SupplierPayload,
					UID string,
				) (*base.Supplier, error) {
					return &base.Supplier{}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(
					ctx context.Context,
					profileID string,
					data *base.Supplier,
				) error {
					return fmt.Errorf("failed to update supplier profile")
				}
			}

			if tt.name == "invalid:_failed_to_send_email" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagementSvs.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return fmt.Errorf("failed to send email")
				}
			}

			if tt.name == "invalid:_failed_to_send_sms" {
				fakeRepo.FetchKYCProcessingRequestByIDFn = func(
					ctx context.Context,
					id string,
				) (*domain.KYCRequest, error) {
					profileID := uuid.New().String()
					return &domain.KYCRequest{
						ID: uuid.New().String(),
						SupplierRecord: &base.Supplier{
							ProfileID: &profileID,
						},
					}, nil
				}

				fakeRepo.UpdateKYCProcessingRequestFn = func(
					ctx context.Context,
					sup *domain.KYCRequest,
				) error {
					return nil
				}

				fakeRepo.GetUserProfileByIDFn = func(
					ctx context.Context,
					id string,
					suspended bool,
				) (*base.UserProfile, error) {
					email := base.GenerateRandomEmail()
					phone := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  uuid.New().String(),
						PrimaryEmailAddress: &email,
						PrimaryPhone:        &phone,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(
					ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeEngagementSvs.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeEngagementSvs.SendSMSFn = func(
					phoneNumbers []string,
					message string,
				) error {
					return fmt.Errorf("failed to send sms")
				}
			}

			_, err := i.Supplier.ProcessKYCRequest(
				tt.args.ctx,
				tt.args.id,
				tt.args.status,
				tt.args.rejectionReason,
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

func TestSupplierUseCasesImpl_AddOrganizationPharmaceuticalKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "VALID-ID-TYPE",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationPharmaceutical{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "A0352HDAKCS",
		KRAPINUploadID:       "kra-pin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		CertificateOfIncorporation:         "certificate_of_incorporation",
		CertificateOfInCorporationUploadID: "certificate_of_incorporation_upload_id",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345",
				IdentificationDocNumberUploadID: "upload_id",
			},
		},
		RegistrationNumber:      "registration_number",
		PracticeLicenseID:       "license_id",
		PracticeLicenseUploadID: "license_upload_id",
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationPharmaceutical
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationPharmaceutical
		wantErr bool
	}{
		{
			name: "valid:_successfully_AddOrganizationPharmaceuticalKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:_fail_to_AddOrganizationPharmaceuticalKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: true,
		},
		{
			name: "invalid:_use_invalid_organization_name",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPharmaceutical{
					OrganizationTypeName: "invalid organization name",
				},
			},
			want:    &validInput,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_find_supplierByUID",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_successfully_AddOrganizationPharmaceuticalKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_fail_to_AddOrganizationPharmaceuticalKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to add Organization Pharmaceutical Kyc")
				}
			}

			if tt.name == "invalid:_use_invalid_organization_name" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid organization name")
				}
			}

			if tt.name == "invalid:_unable_to_find_supplierByUID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "-91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
							},
						},
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get the supplier profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			_, err := i.Supplier.AddOrganizationPharmaceuticalKyc(tt.args.ctx, tt.args.input)

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

func TestSupplierUseCasesImpl_SuspendSupplier(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	suspensionReason := `
	"This email is to inform you that as a result of your actions on April 12th, 2021, you have been issued a suspension for 1 week (7 days)"
	`

	type args struct {
		ctx              context.Context
		suspensionReason *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:successfully_suspend_supplier",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_suspend_supplier",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_user_profile",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier_profile",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_logged_in_user",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_suspend_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3-bvihjskdn", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af315a4e",
							},
						},
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ProfileID:    &profileID,
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeEngagementSvs.NotifySupplierOnSuspensionFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_suspend_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3-bvihjskdn", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af315a4e",
							},
						},
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						ID: "-91bd-42b3af315a4e",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("failed tp suspend supplier")
				}
			}

			if tt.name == "invalid:fail_to_get_user_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get a user profile")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3-bvihjskdn", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af315a4e",
							},
						},
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier profile")
				}
			}

			if tt.name == "invalid:fail_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			got, err := i.Supplier.SuspendSupplier(tt.args.ctx, tt.args.suspensionReason)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.SuspendSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SupplierUseCasesImpl.SuspendSupplier() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddOrganizationRiderKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "VALID-ID-TYPE1",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationRider{
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		CertificateOfIncorporation:         "some-incorp-certificate",
		CertificateOfInCorporationUploadID: "some-incorp-certificate-uploadID",
		KRAPIN:                             "some-someKRAPIN",
		KRAPINUploadID:                     "some-KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload-id",
			},
		},
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationRider
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationRider
		wantErr bool
	}{
		{
			name: "valid:successfully_AddOrganizationRiderKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_Add_OrganizationRiderKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_use_invalid_organization_name",
			args: args{
				ctx: ctx,
				input: domain.OrganizationRider{
					OrganizationTypeName: "invalid organization name",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_AddOrganizationRiderKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumer := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumer,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_Add_OrganizationRiderKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to add Organization rider Kyc")
				}
			}

			if tt.name == "invalid:_use_invalid_organization_name" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid organization name")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddOrganizationRiderKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationRiderKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationRiderKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddOrganizationPractitionerKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "91bd-42b3af315a5c-p4f39af7-5b64-4c2f",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationPractitioner{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "provider-random-kra-pin",
		KRAPINUploadID:       "provider-krapin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		CertificateOfIncorporation:         "provider-incorp-certificate",
		CertificateOfInCorporationUploadID: "provider-incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "provider-id-upload",
			},
		},
		RegistrationNumber:      "provider-reg-no",
		PracticeLicenseID:       "provider-practice-license-id",
		PracticeLicenseUploadID: "provider-practice-license-uploadid",
		PracticeServices:        domain.AllPractitionerService,
		Cadre:                   domain.PractitionerCadreDoctor,
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationPractitioner
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationPractitioner
		wantErr bool
	}{
		{
			name: "valid:successfully_AddOrganizationPractitionerKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:_use_invalid_organization_name",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPractitioner{
					OrganizationTypeName: "invalid organization name",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_AddOrganizationPractitionerKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "91bd-42b3af3-15a4e-f4f39af7-5b64-4c2f-", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-42b3af315a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_use_invalid_organization_name" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid organization name")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddOrganizationPractitionerKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationPractitionerKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("%v", got)
				t.Logf("%v", tt.want)
				t.Errorf("SupplierUseCasesImpl.AddOrganizationPractitionerKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddOrganizationProviderKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-7e2aead29f2c-8716-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationProvider{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "random-kra-pin",
		KRAPINUploadID:       "krapin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		CertificateOfIncorporation:         "incorp-certificate",
		CertificateOfInCorporationUploadID: "incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
		RegistrationNumber:      "regn-no",
		PracticeLicenseID:       "practice-license-id",
		PracticeLicenseUploadID: "practice-license-uploadid",
		PracticeServices:        domain.AllPractitionerService,
	}
	type args struct {
		ctx   context.Context
		input domain.OrganizationProvider
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationProvider
		wantErr bool
	}{
		{
			name: "valid:successfully_add_organizationProviderKYC",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:_use_invalid_organization_name",
			args: args{
				ctx: ctx,
				input: domain.OrganizationProvider{
					OrganizationTypeName: "invalid organization name",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_organizationProviderKYC" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-5cf354a2-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-42b3af3-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_use_invalid_organization_name" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid organization name")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddOrganizationProviderKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationProviderKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationProviderKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddOrganizationCoachKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-7e2aead29f2c-8716-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationCoach{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "coach-random-kra-pin",
		KRAPINUploadID:       "coach-krapin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		CertificateOfIncorporation:         "incorp-certificate",
		CertificateOfInCorporationUploadID: "incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
		RegistrationNumber:      "coach-reg-no",
		PracticeLicenseUploadID: "coach-practice-license-uploadid",
	}
	type args struct {
		ctx   context.Context
		input domain.OrganizationCoach
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationCoach
		wantErr bool
	}{
		{
			name: "valid:successfully_AddOrganizationCoachKyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:_use_invalid_organization_name",
			args: args{
				ctx: ctx,
				input: domain.OrganizationCoach{
					OrganizationTypeName: "invalid organization name",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_AddOrganizationCoachKyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_use_invalid_organization_name" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid organization name")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "400d-8716-1d3e-7e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af791bd-5b64-4c2f-15a4e",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddOrganizationCoachKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationCoachKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationCoachKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddOrganizationNutritionKyc(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-7e2aead29f2c-8716-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.OrganizationNutrition{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "nutrition-random-kra-pin",
		KRAPINUploadID:       "nutrition-krapin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		CertificateOfIncorporation:         "incorp-certificate",
		CertificateOfInCorporationUploadID: "incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           base.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
		RegistrationNumber:      "nutrition-reg-no",
		PracticeLicenseUploadID: "nutrition-practice-license-uploadid",
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationNutrition
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.OrganizationNutrition
		wantErr bool
	}{
		{
			name: "valid:successfully_add_organizationNutritionKYC",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_add_organizationNutritionKYC",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_organizationNutritionKYC" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeOrganisation
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_add_organizationNutritionKYC" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to add organization nutrition kyc")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}
			}
			got, err := i.Supplier.AddOrganizationNutritionKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationNutritionKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddOrganizationNutritionKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_RetireKYCRequest(t *testing.T) {
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
			name: "valid:retire_kyc_request",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_retire_kyc_request",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_logged_in_user",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_user_profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier_profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:retire_kyc_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.RemoveKYCProcessingRequestFn = func(ctx context.Context, supplierProfileID string) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_retire_kyc_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.RemoveKYCProcessingRequestFn = func(ctx context.Context, supplierProfileID string) error {
					return fmt.Errorf("failed to remove kyc processing request")
				}
			}

			if tt.name == "invalid:fail_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "invalid:fail_to_get_user_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get userprofile")
				}
			}

			if tt.name == "invalid:fail_to_get_supplier_profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier profile")
				}
			}

			err := i.Supplier.RetireKYCRequest(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.RetireKYCRequest() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSupplierUseCasesImpl_AddIndividualRiderKyc(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-8716-7e2aead29f2c-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.IndividualRider{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           base.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678910",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:                         "A034RND82",
		KRAPINUploadID:                 "kra-pin-upload-id",
		DrivingLicenseID:               "driving-license-id",
		DrivingLicenseUploadID:         "license-upload-id",
		CertificateGoodConductUploadID: "good-conduct-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualRider
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.IndividualRider
		wantErr bool
	}{
		{
			name: "valid:successfully_add_individual_rider_kyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:kyc_already_submitted",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_kyc_processing_request",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:user_invalid_identificationDocType",
			args: args{
				ctx: ctx,
				input: domain.IndividualRider{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           "invalidDoc",
						IdentificationDocNumber:         "12345678910",
						IdentificationDocNumberUploadID: "id-upload",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_individual_rider_kyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:kyc_already_submitted" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: true,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("kyc already submitted")
				}
			}

			if tt.name == "invalid:fail_to_save_kyc_processing_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage kyc processing request")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			if tt.name == "invalid:user_invalid_identificationDocType" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid doctype used")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddIndividualRiderKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddIndividualRiderKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddIndividualRiderKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddIndividualPractitionerKyc(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-8716-7e2aead29f2c-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.IndividualPractitioner{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           base.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678910",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "A034RND82",
		KRAPINUploadID: "kra-pin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		RegistrationNumber:      "123456",
		PracticeLicenseID:       "license-id",
		PracticeLicenseUploadID: "practice-license-uploadID",
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceLabServices,
		},
		Cadre: domain.PractitionerCadreClinicalOfficer,
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualPractitioner
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.IndividualPractitioner
		wantErr bool
	}{
		{
			name: "valid:successfully_add_individual_practitioner_kyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:kyc_already_submitted",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_kyc_processing_request",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_invalid_PracticeService",
			args: args{
				ctx: ctx,
				input: domain.IndividualPractitioner{
					PracticeServices: []domain.PractitionerService{
						"invalidPracticeService",
						"invalidPracticeService2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_individual_practitioner_kyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:kyc_already_submitted" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: true,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("kyc already submitted")
				}
			}

			if tt.name == "invalid:fail_to_save_kyc_processing_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage kyc processing request")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			if tt.name == "invalid:_invalid_PracticeService" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid PracticeService used")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddIndividualPractitionerKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddIndividualPractitionerKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddIndividualPractitionerKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddIndividualPharmaceuticalKyc(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-8716-7e2aead29f2c-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.IndividualPharmaceutical{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           base.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678910",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "A034RND82",
		KRAPINUploadID: "kra-pin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		RegistrationNumber:      "123456",
		PracticeLicenseID:       "license-id",
		PracticeLicenseUploadID: "practice-license-uploadID",
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualPharmaceutical
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.IndividualPharmaceutical
		wantErr bool
	}{
		{
			name: "valid:successfully_add_individual_pharmaceutical_kyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:kyc_already_submitted",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_kyc_processing_request",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_invalid_IdentificationDocType",
			args: args{
				ctx: ctx,
				input: domain.IndividualPharmaceutical{
					IdentificationDoc: domain.Identification{
						IdentificationDocType: "invalid DocType",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_individual_pharmaceutical_kyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:kyc_already_submitted" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: true,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("kyc already submitted")
				}
			}

			if tt.name == "invalid:fail_to_save_kyc_processing_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage kyc processing request")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			if tt.name == "invalid:_invalid_IdentificationDocType" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid IdentificationDocType used")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddIndividualPharmaceuticalKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddIndividualPharmaceuticalKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddIndividualPharmaceuticalKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddIndividualCoachKyc(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-8716-7e2aead29f2c-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.IndividualCoach{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           base.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678910",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "A034RND82",
		KRAPINUploadID: "kra-pin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		PracticeLicenseID:       "license-id",
		PracticeLicenseUploadID: "practice-license-uploadID",
		AccreditationID:         "ACR-12345678",
		AccreditationUploadID:   "ACR-UPLOAD-12345678",
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualCoach
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.IndividualCoach
		wantErr bool
	}{
		{
			name: "valid:successfully_add_individual_coach_kyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:kyc_already_submitted",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_kyc_processing_request",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_invalid_IdentificationDocType",
			args: args{
				ctx: ctx,
				input: domain.IndividualCoach{
					IdentificationDoc: domain.Identification{
						IdentificationDocType: "invalid DocType",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_add_individual_coach_kyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:kyc_already_submitted" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: true,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("kyc already submitted")
				}
			}

			if tt.name == "invalid:fail_to_save_kyc_processing_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage kyc processing request")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			if tt.name == "invalid:_invalid_IdentificationDocType" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("invalid IdentificationDocType used")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddIndividualCoachKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddIndividualCoachKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddIndividualCoachKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_AddIndividualNutritionKyc(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	validRespPayload := `{"IsPublished":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	admin1 := &base.UserProfile{
		ID: "8716-8716-7e2aead29f2c-7e2aead29f2c",
	}
	adminUsers := []*base.UserProfile{}
	adminUsers = append(adminUsers, admin1)

	validInput := domain.IndividualNutrition{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           base.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678910",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "A034RND82",
		KRAPINUploadID: "kra-pin-upload-id",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		PracticeLicenseID:       "license-id",
		PracticeLicenseUploadID: "practice-license-uploadID",
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualNutrition
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.IndividualNutrition
		wantErr bool
	}{
		{
			name: "valid:successfully_add_individual_nutrition_kyc",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want:    &validInput,
			wantErr: false,
		},
		{
			name: "invalid:kyc_already_submitted",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_save_kyc_processing_request",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_add_individual_nutrition_kyc" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						ID:           "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID:    &profileID,
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accountType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
						AccountType:  &accountType,
					}, nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					email := "test@example.com"
					firstName := "Makmende"
					primaryPhoneNumber := base.TestUserPhoneNumber
					return &base.UserProfile{
						ID:                  "400d-8716--91bd-42b3af315a4e",
						PrimaryPhone:        &primaryPhoneNumber,
						PrimaryEmailAddress: &email,
						UserBioData: base.BioData{
							FirstName: &firstName,
							LastName:  &firstName,
						},
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}
				fakeEngagementSvs.SendAlertToSupplierFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeEngagementSvs.NotifyAdminsFn = func(input dto.EmailNotificationPayload) error {
					return nil
				}
				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}
				fakeEngagementSvs.PublishKYCFeedItemFn = func(uid string, payload base.Item) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:kyc_already_submitted" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: true,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("kyc already submitted")
				}
			}

			if tt.name == "invalid:fail_to_save_kyc_processing_request" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeRepo.FetchAdminUsersFn = func(ctx context.Context) ([]*base.UserProfile, error) {
					return adminUsers, nil
				}

				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage kyc processing request")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			got, err := i.Supplier.AddIndividualNutritionKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddIndividualNutritionKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.AddIndividualNutritionKyc() = %v, want %v", got, tt.want)
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

func TestSupplierUseCasesImpl_CreateSupplierAccount(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		name        string
		partnerType base.PartnerType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:)",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeCoach,
			},
			wantErr: false,
		},
		{
			name: "sad:( can't get logged in user",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeRider,
			},
			wantErr: true,
		},
		{
			name: "sad:( currency not found",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeRider,
			},
			wantErr: true,
		},
		{
			name: "sad:( failed to publsih to PubSub",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeRider,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy:)" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}
			}

			if tt.name == "sad:( can't get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("fail to fetch default currency")
				}
			}

			if tt.name == "sad:( currency not found" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					return nil, fmt.Errorf("fail to fetch default currency")
				}
			}

			if tt.name == "sad:( failed to publsih to PubSub" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return fmt.Errorf("error")
				}
			}

			err := i.Supplier.CreateSupplierAccount(tt.args.ctx, tt.args.name, tt.args.partnerType)
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

func TestSupplierUseCasesImpl_SupplierSetDefaultLocation(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	testChargeMasterBranchID := "94294577-6b27-4091-9802-1ce0f2ce4153"

	cursor := "1234"
	edges := &dto.BranchEdge{
		Cursor: &cursor,
		Node: &domain.Branch{
			ID:                    testChargeMasterBranchID,
			Name:                  "BRANCH-NAME",
			OrganizationSladeCode: "PRO-1234",
			BranchSladeCode:       "1",
		},
	}

	newEdges := []*dto.BranchEdge{}
	newEdges = append(newEdges, edges)

	type args struct {
		ctx        context.Context
		locationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:set_default_location_with_a_valid_locationID",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			wantErr: false,
		},
		{
			name: "invalid:set_default_location_with_an_invalid_locationID",
			args: args{
				ctx:        ctx,
				locationID: "invalid-location-id",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_logged_in_user",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			wantErr: true,
		},
		{
			name: "invalid:supplier_not_found",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_user_profile_by_uid",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_find_branch",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:set_default_location_with_a_valid_locationID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "83d3479d-e902-4aab-a27d-6d5067454daf", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "83d3479d-e902-4aab-a27d-6d5067454daf", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "94294577-6b27-4091-9802-1ce0f2ce4153",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newEdges,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
					}, nil
				}

			}

			if tt.name == "invalid:set_default_location_with_an_invalid_locationID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newEdges,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("fail to get the location")
				}
			}

			if tt.name == "invalid:fail_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "invalid:supplier_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")

				}
			}

			if tt.name == "invalid:_unable_to_get_user_profile_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "FSO798-AD3", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			if tt.name == "invalid:_unable_to_find_branch" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return nil, fmt.Errorf("unable to find branch")
				}
			}

			_, err := i.Supplier.SupplierSetDefaultLocation(tt.args.ctx, tt.args.locationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.SupplierSetDefaultLocation() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSupplierUseCasesImpl_FetchSupplierAllowedLocations(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	testChargeMasterParentOrgId := "83d3479d-e902-4aab-a27d-6d5067454daf"
	testChargeMasterBranchID := "94294577-6b27-4091-9802-1ce0f2ce4153"

	sladeCode := "1"
	cursor := "4567"
	edges := &dto.BranchEdge{
		Cursor: &cursor,
		Node: &domain.Branch{
			ID:                    testChargeMasterParentOrgId,
			Name:                  "BRANCH-NAME",
			OrganizationSladeCode: "PRO-1234",
			BranchSladeCode:       sladeCode,
		},
	}
	newEdges := []*dto.BranchEdge{}
	newEdges = append(newEdges, edges)

	// The Node ID is different from the supplier Location ID
	// This helps to test all cases
	payload2 := &dto.BranchEdge{
		Cursor: &cursor,
		Node: &domain.Branch{
			ID:                    testChargeMasterBranchID,
			Name:                  "BRANCH-NAME",
			OrganizationSladeCode: "PRO-1234",
			BranchSladeCode:       sladeCode,
		},
	}
	newPayload := []*dto.BranchEdge{}
	newPayload = append(newPayload, payload2)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:supplier_allowed_location_found",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "valid:supplier_location_found",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "valid:nil_supplier_location",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_find_branch",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid:logged_in_user_not_found",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:supplier_allowed_location_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newEdges,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}
			}

			if tt.name == "valid:supplier_location_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newPayload,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}
			}

			if tt.name == "valid:nil_supplier_location" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				// Here we dont pass a location
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newPayload,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_find_branch" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID: "8716-7e2ae-5cf354a2-1d3e-ad29f2c-400d",
						Location: &base.Location{
							ID:              testChargeMasterParentOrgId,
							Name:            "BRANCH-NAME",
							BranchSladeCode: &sladeCode,
						},
					}, nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newEdges,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, fmt.Errorf("failed to find branch")
				}
			}

			if tt.name == "invalid:logged_in_user_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("user not found")
				}
			}

			_, err := i.Supplier.FetchSupplierAllowedLocations(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FetchSupplierAllowedLocations() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSupplierUseCasesImpl_SupplierEDILogin(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	sladeCode := "1"
	savannahOrgName := "Savannah Informatics"
	cursor := "8765"
	parent := "parent"
	edges := &dto.BusinessPartnerEdge{
		Cursor: &cursor,
		Node: &domain.BusinessPartner{
			ID:        "BUS1N3SS-P123-1D",
			Name:      savannahOrgName,
			SladeCode: sladeCode,
			Parent:    &parent,
		},
	}

	newEdges := []*dto.BusinessPartnerEdge{}
	newEdges = append(newEdges, edges)

	payload2 := &dto.BranchEdge{
		Cursor: &cursor,
		Node: &domain.Branch{
			ID:                    "BUS1N3SS-P123-1D",
			Name:                  savannahOrgName,
			OrganizationSladeCode: "123456",
			BranchSladeCode:       sladeCode,
		},
	}
	newPayload := []*dto.BranchEdge{}
	newPayload = append(newPayload, payload2)

	payload3 := &dto.BusinessPartnerEdge{
		Cursor: &cursor,
		Node: &domain.BusinessPartner{
			ID:        "BUS1N3SS-P123",
			Name:      "Random Org",
			SladeCode: "PRO-1234",
			Parent:    &parent,
		},
	}

	newPayload3 := []*dto.BusinessPartnerEdge{}
	newPayload3 = append(newPayload3, payload3)

	payload4 := &dto.BranchEdge{
		Cursor: &cursor,
		Node: &domain.Branch{
			ID:                    "BUS1N3SS-P123",
			Name:                  "Random Org",
			OrganizationSladeCode: "1234",
			BranchSladeCode:       "PRO-1234",
		},
	}
	newPayload4 := []*dto.BranchEdge{}
	newPayload4 = append(newPayload4, payload4)

	// This will help test the case where a parent is nil
	payload5 := &dto.BusinessPartnerEdge{
		Cursor: &cursor,
		Node: &domain.BusinessPartner{
			ID:        "BUS1N3SS-P123",
			Name:      "Random Org",
			SladeCode: "PRO-1234",
		},
	}

	newPayload5 := []*dto.BusinessPartnerEdge{}
	newPayload5 = append(newPayload5, payload5)

	type args struct {
		ctx       context.Context
		username  string
		password  string
		sladeCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:successful_login",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "1",
			},
			wantErr: false,
		},
		{
			name: "valid:successful_login_with_a_non-savannah_sladeCode",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "PRO-1234",
			},
			wantErr: false,
		},
		{
			name: "valid:nil_business_parent",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "PRO-1234",
			},
			wantErr: false,
		},
		{
			name: "valid:nil_business_parent_fail_to_updateSupplierProfile",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "PRO-1234",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_update_permissions",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "1",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_ediUserProfile_using_SavannahSladeCode",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "1",
			},
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_ediUserProfile_using_non-SavannahSladeCode",
			args: args{
				ctx:       ctx,
				username:  "bewell@slade360.co.ke",
				password:  "please change me",
				sladeCode: "1023",
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_get_logged_in_user",
			args: args{
				ctx:       ctx,
				username:  "userName",
				password:  "1234der5",
				sladeCode: "PRO-1234",
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_get_userProfileByUID",
			args: args{
				ctx:       ctx,
				username:  "userName",
				password:  "1234der5",
				sladeCode: "PRO-1234",
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_find_SupplierByUID",
			args: args{
				ctx:       ctx,
				username:  "userName",
				password:  "1234der5",
				sladeCode: "PRO-1234",
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_find_supplier_by_id",
			args: args{
				ctx:       ctx,
				username:  "userName",
				password:  "1234der5",
				sladeCode: "PRO-1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successful_login" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:        578278332,
						GUID:      uuid.New().String(),
						Email:     "juhakalulu@gmail.com",
						FirstName: "Juha",
						LastName:  "Kalulu",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}

				fakeChargeMasterSvc.FindProviderFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
					sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
					return &dto.BusinessPartnerConnection{
						Edges: newEdges,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(uid string, payload base.Nudge) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}

				fakeChargeMasterSvc.FetchProviderByIDFn = func(ctx context.Context, id string) (*domain.BusinessPartner, error) {
					return &domain.BusinessPartner{
						ID:        "BUS1N3SS-P123-1D",
						Name:      savannahOrgName,
						SladeCode: "1",
						Parent:    &parent,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newPayload,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

			}

			if tt.name == "invalid:unable_to_find_supplier_by_id" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("unable to get supplier profile by profile id")
				}
			}

			if tt.name == "valid:successful_login_with_a_non-savannah_sladeCode" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:              578278332,
						GUID:            uuid.New().String(),
						Email:           "juhakalulu@gmail.com",
						FirstName:       "Juha",
						LastName:        "Kalulu",
						BusinessPartner: "1234",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}

				fakeChargeMasterSvc.FindProviderFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
					sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
					return &dto.BusinessPartnerConnection{
						Edges: newPayload3,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(uid string, payload base.Nudge) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}

				fakeChargeMasterSvc.FetchProviderByIDFn = func(ctx context.Context, id string) (*domain.BusinessPartner, error) {
					parent := "parent"
					return &domain.BusinessPartner{
						ID:        "BUS1N3SS-P123",
						Name:      "Random Org",
						SladeCode: "PRO-1234",
						Parent:    &parent,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeChargeMasterSvc.FindBranchFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
					sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
					return &dto.BranchConnection{
						Edges: newPayload4,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

			}

			if tt.name == "valid:nil_business_parent" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:              578278332,
						GUID:            uuid.New().String(),
						Email:           "juhakalulu@gmail.com",
						FirstName:       "Juha",
						LastName:        "Kalulu",
						BusinessPartner: "1234",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}

				fakeChargeMasterSvc.FindProviderFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
					sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
					return &dto.BusinessPartnerConnection{
						Edges: newPayload5,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(uid string, payload base.Nudge) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}
			}

			if tt.name == "valid:nil_business_parent_fail_to_updateSupplierProfile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:              578278332,
						GUID:            uuid.New().String(),
						Email:           "juhakalulu@gmail.com",
						FirstName:       "Juha",
						LastName:        "Kalulu",
						BusinessPartner: "1234",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return nil
				}

				fakeChargeMasterSvc.FindProviderFn = func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
					sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
					return &dto.BusinessPartnerConnection{
						Edges: newPayload5,
						PageInfo: &base.PageInfo{
							HasNextPage: false,
						},
					}, nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(uid string, payload base.Nudge) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return fmt.Errorf("failed to update supplier profile")
				}
			}

			if tt.name == "invalid:fail_to_update_permissions" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:        578278332,
						GUID:      uuid.New().String(),
						Email:     "juhakalulu@gmail.com",
						FirstName: "Juha",
						LastName:  "Kalulu",
					}, nil
				}

				fakeRepo.UpdateSupplierProfileFn = func(ctx context.Context, profileID string, data *base.Supplier) error {
					return nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "f4f39af7-5b64-4c2f-91bd-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{ID: "12334"}, nil
				}

				fakeRepo.UpdatePermissionsFn = func(ctx context.Context, id string, perms []base.PermissionType) error {
					return fmt.Errorf("failed to update permissions")
				}
			}

			if tt.name == "invalid:fail_to_get_ediUserProfile_using_SavannahSladeCode" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return nil, fmt.Errorf("cannot get edi user profile")
				}
			}

			if tt.name == "invalid:fail_to_get_ediUserProfile_using_non-SavannahSladeCode" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "42b3af315a4e-f4f39af7-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					accType := base.AccountTypeIndividual
					return &base.Supplier{
						SupplierID:        "5cf354a2-8716-7e2ae-1d3e-ad29f2c-400d",
						ID:                uid,
						AccountType:       &accType,
						UnderOrganization: true,
					}, nil
				}

				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, fmt.Errorf("edi user profile not found")
				}

				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return nil, fmt.Errorf("cannot get edi user profile")
				}
			}

			if tt.name == "invalid:unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "invalid:unable_to_get_userProfileByUID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get user profile by UID")
				}
			}

			if tt.name == "unable_to_find_SupplierByUID" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd-42b3af315a4e",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "f4f39af7-91bd-42b3af-315a4e",
							},
						},
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("supplier not found")
				}
			}

			resp, err := i.Supplier.SupplierEDILogin(tt.args.ctx, tt.args.username, tt.args.password, tt.args.sladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.SupplierEDILogin() error = %v, wantErr %v", err, tt.wantErr)
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

				if resp == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestUnitSupplierUseCasesImpl_AddPartnerType(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	testRiderName := "Test Rider"
	rider := base.PartnerTypeRider

	type args struct {
		ctx         context.Context
		name        *string
		partnerType *base.PartnerType
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid:add_partner_type",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid : missing name arg",
			args: args{
				ctx: ctx,
			},
			want:        false,
			wantErr:     true,
			expectedErr: "expected `name` to be defined and `partnerType` to be valid",
		},
		{
			name: "invalid:unable_to_login",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:unable_to_get_user_profile_by_id",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:unable_to_add_partner_type",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:add_partner_type" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.AddPartnerTypeFn = func(ctx context.Context, profileID string, name *string, partnerType *base.PartnerType) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid:unable_to_login" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("unable to login")
				}
			}

			if tt.name == "invalid:unable_to_get_user_profile_by_id" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile by uid")
				}

			}

			if tt.name == "invalid:unable_to_add_partner_type" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspend bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.AddPartnerTypeFn = func(ctx context.Context, profileID string, name *string, partnerType *base.PartnerType) (bool, error) {
					return false, fmt.Errorf("unable to add partner type")
				}
			}

			got, err := i.Supplier.AddPartnerType(tt.args.ctx, tt.args.name, tt.args.partnerType)

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

				if got != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, got)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_FindSupplierByUID(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	profileID := "93ca42bb-5cfc-4499-b137-2df4d67b4a21"
	supplier := &base.Supplier{
		ProfileID: &profileID,
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Supplier
		wantErr bool
	}{
		{
			name: "valid:_find_supplier_by_uid",
			args: args{
				ctx: ctx,
			},
			want:    supplier,
			wantErr: false,
		},
		{
			name: "invalid:_find_supplier_by_uid",
			args: args{
				ctx: ctx,
			},
			want:    supplier,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_find_supplier_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return &base.Supplier{
						ID:        "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						ProfileID: &profileID,
					}, nil
				}

			}

			if tt.name == "invalid:_find_supplier_by_uid" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-87167-e2aead29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "5b64-4c2f-15a4e-f4f39af791bd-42b3af3",
					}, nil
				}
				fakeRepo.GetSupplierProfileByProfileIDFn = func(ctx context.Context, profileID string) (*base.Supplier, error) {
					return nil, fmt.Errorf("failed to get supplier")
				}

			}

			sup, err := i.Supplier.FindSupplierByUID(tt.args.ctx)
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

				if *sup.ProfileID == "" {
					t.Errorf("empty profileID returned %v", sup.ProfileID)
					return
				}
			}

		})
	}
}
func TestSupplierUseCase_StageKYCProcessingRequest(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	profileID := "93ca42bb-5cfc-4499-b137-2df4d67b4a21"
	supplier := &base.Supplier{
		ProfileID: &profileID,
	}
	type args struct {
		ctx context.Context
		sup *base.Supplier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_stage_KYC_processing",
			args: args{
				ctx: ctx,
				sup: supplier,
			},
			wantErr: false,
		},
		{
			name: "invalid:_stage_KYC_processing",
			args: args{
				ctx: ctx,
				sup: supplier,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_stage_KYC_processing" {
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return nil
				}

			}
			if tt.name == "invalid:_stage_KYC_processing" {
				fakeRepo.StageKYCProcessingRequestFn = func(ctx context.Context, data *domain.KYCRequest) error {
					return fmt.Errorf("failed to stage KYC processing request")
				}

			}
			err := i.Supplier.StageKYCProcessingRequest(tt.args.ctx, tt.args.sup)

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

func TestUnitSupplierUseCasesImpl_SetUpSupplier(t *testing.T) {
	ctx := context.Background()

	individualPartner := base.AccountTypeIndividual
	organizationPartner := base.AccountTypeOrganisation

	s, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		accountType base.AccountType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful individual supplier account setup",
			args: args{
				ctx:         ctx,
				accountType: individualPartner,
			},
			wantErr: false,
		},
		{
			name: "Successful organization supplier account setup",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: false,
		},
		{
			name: "invalid failed to get the logged in user",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: true,
		},
		{
			name: "invalid failed to get user profile",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: true,
		},
		{
			name: "invalid failed to add supplier account type",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: true,
		},
		{
			name: "invalid:_resolving_the_consumer_nudge_fails",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: false, // the error is logged
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Successful individual supplier account setup" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.AddSupplierAccountTypeFn = func(
					ctx context.Context,
					profileID string,
					accountType base.AccountType,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID:          uuid.New().String(),
						AccountType: &individualPartner,
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(
					ctx context.Context,
					profileID string,
					data *base.Supplier,
				) error {
					return nil
				}
				fakeEngagementSvs.PublishKYCNudgeFn = func(
					uid string,
					payload base.Nudge,
				) (*http.Response, error) {
					return &http.Response{StatusCode: 200}, nil
				}

				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}
			}

			if tt.name == "Successful organization supplier account setup" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.AddSupplierAccountTypeFn = func(
					ctx context.Context,
					profileID string,
					accountType base.AccountType,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID:           uuid.New().String(),
						AccountType:  &organizationPartner,
						SupplierName: "Juha Kalulu",
					}, nil
				}
				fakeRepo.UpdateSupplierProfileFn = func(
					ctx context.Context,
					profileID string,
					data *base.Supplier,
				) error {
					return nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(
					uid string,
					payload base.Nudge,
				) (*http.Response, error) {
					return &http.Response{StatusCode: 200}, nil
				}

				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}
			}

			if tt.name == "invalid failed to get the logged in user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "invalid failed to get user profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "invalid failed to add supplier account type" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.AddSupplierAccountTypeFn = func(
					ctx context.Context,
					profileID string,
					accountType base.AccountType,
				) (*base.Supplier, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "invalid:_resolving_the_consumer_nudge_fails" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(
					ctx context.Context,
					uid string,
					suspend bool,
				) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.GetSupplierProfileByProfileIDFn = func(
					ctx context.Context,
					profileID string,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID: uuid.New().String(),
					}, nil
				}

				fakeRepo.AddSupplierAccountTypeFn = func(
					ctx context.Context,
					profileID string,
					accountType base.AccountType,
				) (*base.Supplier, error) {
					return &base.Supplier{
						ID:          uuid.New().String(),
						AccountType: &individualPartner,
					}, nil
				}

				fakeEngagementSvs.PublishKYCNudgeFn = func(
					uid string,
					payload base.Nudge,
				) (*http.Response, error) {
					return &http.Response{StatusCode: 200}, nil
				}

				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					UID string,
					flavour base.Flavour,
					nudgeTitle string,
				) error {
					return fmt.Errorf("an error occurred")
				}
			}

			_, err := s.Supplier.SetUpSupplier(
				tt.args.ctx,
				tt.args.accountType,
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

func TestUnitSupplierUseCasesImplUnit_EDIUserLogin(t *testing.T) {
	s, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	username := "user"
	password := "pass"

	// var username1 string
	// var password1 string

	type args struct {
		username *string
		password *string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "valid:login_user_with_username_and_password",
			args: args{
				username: &username,
				password: &password,
			},
			wantErr: false,
		},
		{
			name: "invalid:unable_to_initialize_login_client",
			args: args{
				username: &username,
				password: &password,
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_fetch_user_profile",
			args: args{
				username: &username,
				password: &password,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeBaseExt.LoginClientFn = nil
			fakeBaseExt.FetchUserProfileFn = nil
			if tt.name == "valid:login_user_with_username_and_password" {
				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}
				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return &base.EDIUserProfile{
						ID:              578278332,
						GUID:            uuid.New().String(),
						Email:           "juhakalulu@gmail.com",
						FirstName:       "Juha",
						LastName:        "Kalulu",
						BusinessPartner: "1234",
					}, nil
				}

			}

			if tt.name == "invalid:unable_to_initialize_login_client" {
				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, fmt.Errorf("unable to login the client")
				}

			}

			if tt.name == "invalid:unable_to_fetch_user_profile" {
				fakeBaseExt.LoginClientFn = func(username string, password string) (base.Client, error) {
					return nil, nil
				}
				fakeBaseExt.FetchUserProfileFn = func(authClient base.Client) (*base.EDIUserProfile, error) {
					return nil, fmt.Errorf("unable to fetch user profile")
				}
			}

			profile, err := s.Supplier.EDIUserLogin(
				tt.args.username,
				tt.args.password,
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
				if profile.Email != "juhakalulu@gmail.com" {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestSupplierUseCasesImpl_CreateCustomerAccount(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		name        string
		partnerType base.PartnerType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:)",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeConsumer,
			},
			wantErr: false,
		},
		{
			name: "sad:( currency not found",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeConsumer,
			},
			wantErr: true,
		},
		{
			name: "sad:( can't get logged in user",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeConsumer,
			},
			wantErr: true,
		},
		{
			name: "sad:( failed to publsih to PubSub",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: base.PartnerTypeConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy:)" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return nil
				}
			}

			if tt.name == "sad:( can't get logged in user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("error")
				}
			}

			if tt.name == "sad:( currency not found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					return nil, fmt.Errorf("fail to fetch default currency")
				}
			}

			if tt.name == "sad:( failed to publsih to PubSub" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeEPRSvc.FetchERPClientFn = func() *base.ServerClient {
					return &base.ServerClient{}
				}

				fakeBaseExt.FetchDefaultCurrencyFn = func(c base.Client) (*base.FinancialYearAndCurrency, error) {
					id := uuid.New().String()
					return &base.FinancialYearAndCurrency{
						ID: &id,
					}, nil
				}

				fakePubSub.TopicIDsFn = func() []string {
					return []string{uuid.New().String()}
				}

				fakePubSub.AddPubSubNamespaceFn = func(topicName string) string {
					return uuid.New().String()
				}

				fakePubSub.PublishToPubsubFn = func(ctx context.Context, topicID string, payload []byte) error {
					return fmt.Errorf("error")
				}
			}

			err := i.Supplier.CreateCustomerAccount(tt.args.ctx, tt.args.name, tt.args.partnerType)
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
