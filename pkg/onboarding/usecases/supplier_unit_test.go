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
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
)

const (
	testSladeCode               = "BRA-PRO-4190-4"
	testEDIPortalUsername       = "avenue-4190@healthcloud.co.ke"
	testEDIPortalPassword       = "test provider"
	testChargeMasterParentOrgId = "83d3479d-e902-4aab-a27d-6d5067454daf"
	testChargeMasterBranchID    = "94294577-6b27-4091-9802-1ce0f2ce4153"
)

func TestSupplierUseCasesImpl_AddPartnerType(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	testRiderName := "Test Rider"
	rider := base.PartnerTypeRider
	testPractitionerName := "Test Practitioner"
	practitioner := base.PartnerTypePractitioner
	testProviderName := "Test Provider"
	provider := base.PartnerTypeProvider
	testPharmaceuticalName := "Test Pharmaceutical"
	pharmaceutical := base.PartnerTypePharmaceutical
	testCoachName := "Test Coach"
	coach := base.PartnerTypeCoach
	testNutritionName := "Test Nutrition"
	nutrition := base.PartnerTypeNutrition
	testConsumerName := "Test Consumer"
	consumer := base.PartnerTypeConsumer

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
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
			name: "valid: add PartnerTypeRider ",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePractitioner ",
			args: args{
				ctx:         ctx,
				name:        &testPractitionerName,
				partnerType: &practitioner,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeProvider ",
			args: args{
				ctx:         ctx,
				name:        &testProviderName,
				partnerType: &provider,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePharmaceutical",
			args: args{
				ctx:         ctx,
				name:        &testPharmaceuticalName,
				partnerType: &pharmaceutical,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeCoach",
			args: args{
				ctx:         ctx,
				name:        &testCoachName,
				partnerType: &coach,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeNutrition",
			args: args{
				ctx:         ctx,
				name:        &testNutritionName,
				partnerType: &nutrition,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "invalid: add PartnerTypeConsumer",
			args: args{
				ctx:         ctx,
				name:        &testConsumerName,
				partnerType: &consumer,
			},
			want:        false,
			wantErr:     true,
			expectedErr: "invalid `partnerType`. cannot use CONSUMER in this context",
		},

		{
			name: "invalid : invalid context",
			args: args{
				ctx:         context.Background(),
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:        false,
			wantErr:     true,
			expectedErr: `unable to get the logged in user: auth token not found in context: unable to get auth token from context with key "UID" `,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier := s
			got, err := supplier.Supplier.AddPartnerType(tt.args.ctx, tt.args.name, tt.args.partnerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddPartnerType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SupplierUseCasesImpl.AddPartnerType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUpSupplier(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	individualPartner := base.AccountTypeIndividual
	organizationPartner := base.AccountTypeOrganisation

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
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
			name: "SadCase - Invalid supplier setup",
			args: args{
				ctx:         ctx,
				accountType: "non existent type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier := s
			_, err := supplier.Supplier.SetUpSupplier(tt.args.ctx, tt.args.accountType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUpSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestSupplierUseCasesImpl_EDIUserLogin(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	validUsername := "avenue-4190@healthcloud.co.ke"
	validPassword := "test provider"

	invalidUsername := "username"
	invalidPassword := "password"

	emptyUsername := ""
	emptyPassword := ""
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
			name: "Happy Case: valid credentials",
			args: args{
				username: &validUsername,
				password: &validPassword,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Wrong userame and password",
			args: args{
				username: &invalidUsername,
				password: &invalidPassword,
			},
			wantErr: true,
		},
		{
			name: "sad case: empty username and password",
			args: args{
				username: &emptyUsername,
				password: &emptyPassword,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ediLogin := s
			_, err := ediLogin.Supplier.EDIUserLogin(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.EDIUserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSupplierUseCasesImpl_CoreEDIUserLogin(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: valid credentials",
			args: args{
				username: "bewell@slade360.co.ke",
				password: "please change me",
			},
			wantErr: true, // TODO: switch to true when https://accounts-core.release.slade360.co.ke/
			// comes back live
		},
		{
			name: "Sad Case: Wrong userame and password",
			args: args{
				username: "invalid Username",
				password: "invalid Password",
			},
			wantErr: true,
		},
		{
			name: "sad case: empty username and password",
			args: args{
				username: "",
				password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coreEdiLogin := s
			_, err := coreEdiLogin.Supplier.CoreEDIUserLogin(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.CoreEDIUserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func clean(newCtx context.Context, testPhoneNumber string, t *testing.T, service *interactor.Interactor) {
	err := service.Signup.RemoveUserByPhoneNumber(newCtx, testPhoneNumber)
	if err != nil {
		t.Errorf("failed to clean data after test error: %v", err)
		return
	}
}

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
				fakeMailgunSvc.SendMailFn = func(email string, message string, subject string) error {
					return nil
				}

			}

			if tt.name == "invalid:_send_mail_fails" {
				fakeMailgunSvc.SendMailFn = func(email string, message string, subject string) error {
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

				fakeMailgunSvc.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeMessagingSvc.SendSMSFn = func(
					phoneNumbers []string,
					message string,
				) error {
					return nil
				}
			}

			if tt.name == "valid:_rejected_a_kyc_request" {
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

				fakeMailgunSvc.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeMessagingSvc.SendSMSFn = func(
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

				fakeMailgunSvc.SendMailFn = func(
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

				fakeMailgunSvc.SendMailFn = func(
					email string,
					message string,
					subject string,
				) error {
					return nil
				}

				fakeMessagingSvc.SendSMSFn = func(
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		SupportingDocumentsUploadID: []string{
			"supporting_docs_upload_id",
			"random_upload_id",
		},
		CertificateOfIncorporation:         "certificate_of_incorporation",
		CertificateOfInCorporationUploadID: "certificate_of_incorporation_upload_id",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345",
				IdentificationDocNumberUploadID: "upload_id",
			},
		},
		OrganizationCertificate: "organization_certificate",
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
					return nil
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		want    bool
		wantErr bool
	}{
		{
			name: "valid:successfully_suspend_supplier",
			args: args{
				ctx: ctx,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:fail_to_suspend_supplier",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_user_profile",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_supplier_profile",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:fail_to_get_logged_in_user",
			args: args{
				ctx: ctx,
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

			got, err := i.Supplier.SuspendSupplier(tt.args.ctx)
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		OrganizationCertificate:            "some-org-cert",
		KRAPIN:                             "some-someKRAPIN",
		KRAPINUploadID:                     "some-KRAPINUploadID",
		SupportingDocumentsUploadID:        []string{"SupportingDocumentsUploadID", "Support"},
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
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
					return nil
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "provider-random-kra-pin",
		KRAPINUploadID:                     "provider-krapin-upload-id",
		SupportingDocumentsUploadID:        []string{"uploadid-1", "uploadid-2"},
		CertificateOfIncorporation:         "provider-incorp-certificate",
		CertificateOfInCorporationUploadID: "provider-incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "provider-id-upload",
			},
		},
		OrganizationCertificate: "provider-organization-cert",
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
					return &base.Supplier{
						ID:        "42b3af315a4e-5b64-4c2f-91bd",
						ProfileID: &profileID,
					}, nil
				}

				fakeRepo.GetSupplierProfileByUIDFn = func(ctx context.Context, uid string) (*base.Supplier, error) {
					return &base.Supplier{
						SupplierID:   "8716-7e2ae-ad29f2c-400d",
						KYCSubmitted: false,
					}, nil
				}

				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c-42b3af315a4e", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "400d-8716--91bd",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: "42b3af315a4e-91bd-42b3af-315a4e",
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
					return nil
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "random-kra-pin",
		KRAPINUploadID:                     "krapin-upload-id",
		SupportingDocumentsUploadID:        []string{"uploadid-1", "uploadid-2"},
		CertificateOfIncorporation:         "incorp-certificate",
		CertificateOfInCorporationUploadID: "incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
		OrganizationCertificate: "organization-cert",
		RegistrationNumber:      "regn-no",
		PracticeLicenseID:       "practice-license-id",
		PracticeLicenseUploadID: "practice-license-uploadid",
		Cadre:                   domain.PractitionerCadreDoctor,
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
					return nil
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

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
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "coach-random-kra-pin",
		KRAPINUploadID:                     "coach-krapin-upload-id",
		SupportingDocumentsUploadID:        []string{"uploadid-1", "uploadid-2"},
		CertificateOfIncorporation:         "incorp-certificate",
		CertificateOfInCorporationUploadID: "incorp-certificate-uploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678910",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
		OrganizationCertificate: "coach-organization-cert",
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
					return nil
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
