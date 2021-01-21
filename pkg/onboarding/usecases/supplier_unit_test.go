package usecases_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"gitlab.slade360emr.com/go/base"
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
