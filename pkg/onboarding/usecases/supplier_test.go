package usecases_test

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
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

func TestSuspendSupplier(t *testing.T) {
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
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid case - Suspend an existing supplier",
			args: args{
				ctx: ctx,
			},
			want:    true,
			wantErr: false,
		}, {
			name: "invalid case - Suspend a nonexistent supplier",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := s
			_, err := service.Supplier.SuspendSupplier(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuspendSupplier() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSupplierUseCasesImpl_AddOrganizationProviderKyc(t *testing.T) {
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

	name := "Makmende"
	partnerProvider := base.PartnerTypeProvider
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerProvider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationProvider
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.OrganizationProvider
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.OrganizationProvider{
					OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []domain.PractitionerService{domain.PractitionerServiceOutpatientServices, domain.PractitionerServiceInpatientServices, domain.PractitionerServiceOther},
					Cadre:              domain.PractitionerCadreDoctor,
				},
			},
			want: &domain.OrganizationProvider{
				OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
				DirectorIdentifications: []domain.Identification{
					{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
				},
				KRAPIN:             "KRA-12345678",
				KRAPINUploadID:     "KRA-UPLOAD-12345678",
				RegistrationNumber: "REG-12345",
				PracticeLicenseID:  "PRAC-12345",
				PracticeServices:   []domain.PractitionerService{domain.PractitionerServiceOutpatientServices, domain.PractitionerServiceInpatientServices, domain.PractitionerServiceOther},
				Cadre:              domain.PractitionerCadreDoctor,
			},
			wantErr: false,
		},
		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: domain.OrganizationProvider{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []domain.PractitionerService{"SUPPORTING"},
					Cadre:              domain.PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},
		{
			name: "invalid : practice services",
			args: args{
				ctx: ctx,
				input: domain.OrganizationProvider{
					OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []domain.PractitionerService{"SUPPORTING"},
					Cadre:              domain.PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddOrganizationProviderKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationProviderKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationProviderKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationProviderKyc() = %v, want %v", got, tt.want)
				}
				return
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

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	name := "Makmende"
	partnerPharmaceutical := base.PartnerTypePharmaceutical
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPharmaceutical)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx   context.Context
		input domain.OrganizationPharmaceutical
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.OrganizationPharmaceutical
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPharmaceutical{
					OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:                             "KRA-12345678",
					KRAPINUploadID:                     "KRA-UPLOAD-12345678",
					RegistrationNumber:                 "REG-12345",
					PracticeLicenseID:                  "PRAC-12345",
					PracticeLicenseUploadID:            "PRAC-UPLOAD-12345",
					CertificateOfIncorporation:         "CERT-12345678",
					CertificateOfInCorporationUploadID: "CERT-UPLOAD-12345",
				},
			},
			want: &domain.OrganizationPharmaceutical{
				OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
				DirectorIdentifications: []domain.Identification{
					{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
				},
				KRAPIN:                             "KRA-12345678",
				KRAPINUploadID:                     "KRA-UPLOAD-12345678",
				RegistrationNumber:                 "REG-12345",
				PracticeLicenseID:                  "PRAC-12345",
				PracticeLicenseUploadID:            "PRAC-UPLOAD-12345",
				CertificateOfIncorporation:         "CERT-12345678",
				CertificateOfInCorporationUploadID: "CERT-UPLOAD-12345",
			},
			wantErr: false,
		},
		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPharmaceutical{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:                             "KRA-12345678",
					KRAPINUploadID:                     "KRA-UPLOAD-12345678",
					RegistrationNumber:                 "REG-12345",
					PracticeLicenseID:                  "PRAC-12345",
					PracticeLicenseUploadID:            "PRAC-UPLOAD-12345",
					CertificateOfIncorporation:         "CERT-12345678",
					CertificateOfInCorporationUploadID: "CERT-UPLOAD-12345",
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddOrganizationPharmaceuticalKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPharmaceuticalKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPharmaceuticalKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPharmaceuticalKyc() = %v, want %v", got, tt.want)
				}
				return
			}

		})
	}
}

func TestSupplierUseCasesImpl_AddIndividualPharmaceuticalKyc(t *testing.T) {
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

	name := "Makmende"
	partnerPharmaceutical := base.PartnerTypePharmaceutical
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPharmaceutical)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeIndividual)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualPharmaceutical
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.IndividualPharmaceutical
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.IndividualPharmaceutical{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                  "KRA-12345678",
					KRAPINUploadID:          "KRA-UPLOAD-12345678",
					RegistrationNumber:      "REG-12345",
					PracticeLicenseID:       "PRAC-12345",
					PracticeLicenseUploadID: "PRAC-UPLOAD-12345",
				},
			},
			want: &domain.IndividualPharmaceutical{
				IdentificationDoc: domain.Identification{
					IdentificationDocType:           domain.IdentificationDocTypeNationalid,
					IdentificationDocNumber:         "12345678",
					IdentificationDocNumberUploadID: "12345678",
				},
				KRAPIN:                  "KRA-12345678",
				KRAPINUploadID:          "KRA-UPLOAD-12345678",
				RegistrationNumber:      "REG-12345",
				PracticeLicenseID:       "PRAC-12345",
				PracticeLicenseUploadID: "PRAC-UPLOAD-12345",
			},
			wantErr: false,
		},
		{
			name: "invalid : unauthenticated context",
			args: args{
				ctx: context.Background(),
				input: domain.IndividualPharmaceutical{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                  "KRA-12345678",
					KRAPINUploadID:          "KRA-UPLOAD-12345678",
					RegistrationNumber:      "REG-12345",
					PracticeLicenseID:       "PRAC-12345",
					PracticeLicenseUploadID: "PRAC-UPLOAD-12345",
				},
			},
			wantErr:     true,
			expectedErr: exceptions.SupplierNotFoundError(fmt.Errorf("unauthenticated context")).Error(),
		},
		{
			name: "invalid : wrong identification document type",
			args: args{
				ctx: ctx,
				input: domain.IndividualPharmaceutical{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           "SCHOOL ID",
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                  "KRA-12345678",
					KRAPINUploadID:          "KRA-UPLOAD-12345678",
					RegistrationNumber:      "REG-12345",
					PracticeLicenseID:       "PRAC-12345",
					PracticeLicenseUploadID: "PRAC-UPLOAD-12345",
				},
			},
			wantErr:     true,
			expectedErr: "invalid `IdentificationDocType` provided : SCHOOL ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddIndividualPharmaceuticalKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPharmaceuticalKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPharmaceuticalKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPharmaceuticalKyc() = %v, want %v", got, tt.want)
				}
				return
			}

		})
	}
}

func TestSupplierUseCasesImpl_AddIndividualCoachKyc(t *testing.T) {
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

	name := "Makmende"
	partnerCoach := base.PartnerTypeCoach
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerCoach)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeIndividual)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualCoach
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.IndividualCoach
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.IndividualCoach{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                      "KRA-12345678",
					KRAPINUploadID:              "KRA-UPLOAD-12345678",
					PracticeLicenseID:           "PRAC-12345",
					PracticeLicenseUploadID:     "PRAC-UPLOAD-12345",
					SupportingDocumentsUploadID: []string{"SUPP-UPLOAD-ID-1234", "SUPP-UPLOAD-ID-1234"},
				},
			},
			want: &domain.IndividualCoach{
				IdentificationDoc: domain.Identification{
					IdentificationDocType:           domain.IdentificationDocTypeNationalid,
					IdentificationDocNumber:         "12345678",
					IdentificationDocNumberUploadID: "12345678",
				},
				KRAPIN:                      "KRA-12345678",
				KRAPINUploadID:              "KRA-UPLOAD-12345678",
				PracticeLicenseID:           "PRAC-12345",
				PracticeLicenseUploadID:     "PRAC-UPLOAD-12345",
				SupportingDocumentsUploadID: []string{"SUPP-UPLOAD-ID-1234", "SUPP-UPLOAD-ID-1234"},
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				input: domain.IndividualCoach{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                      "KRA-12345678",
					KRAPINUploadID:              "KRA-UPLOAD-12345678",
					PracticeLicenseID:           "PRAC-12345",
					PracticeLicenseUploadID:     "PRAC-UPLOAD-12345",
					SupportingDocumentsUploadID: []string{"SUPP-UPLOAD-ID-1234", "SUPP-UPLOAD-ID-1234"},
				},
			},
			wantErr:     true,
			expectedErr: exceptions.SupplierNotFoundError(fmt.Errorf("unauthenticated context")).Error(),
		},
		{
			name: "invalid: wrong identification document type",
			args: args{
				ctx: ctx,
				input: domain.IndividualCoach{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           "SCHOOL ID",
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:                      "KRA-12345678",
					KRAPINUploadID:              "KRA-UPLOAD-12345678",
					PracticeLicenseID:           "PRAC-12345",
					PracticeLicenseUploadID:     "PRAC-UPLOAD-12345",
					SupportingDocumentsUploadID: []string{"SUPP-UPLOAD-ID-1234", "SUPP-UPLOAD-ID-1234"},
				},
			},
			wantErr:     true,
			expectedErr: "invalid `IdentificationDocType` provided : SCHOOL ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddIndividualCoachKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddIndividualCoachKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddIndividualCoachKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddIndividualCoachKyc() = %v, want %v", got, tt.want)
				}
				return
			}

		})
	}
}

func TestAddIndividualRiderKYC(t *testing.T) {
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

	name := "Jatelo"
	partnerRider := base.PartnerTypeRider
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeIndividual)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx   context.Context
		input domain.IndividualRider
	}

	tests := []struct {
		name        string
		args        args
		want        *domain.IndividualRider
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Happy Case - Successfully add individual rider KYC",
			args: args{
				ctx: ctx,
				input: domain.IndividualRider{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "23456789",
					},
					KRAPIN:                         "A0123456",
					KRAPINUploadID:                 "34567890",
					DrivingLicenseID:               "12345678",
					CertificateGoodConductUploadID: "34567890",
				},
			},
			want: &domain.IndividualRider{
				IdentificationDoc: domain.Identification{
					IdentificationDocType:           domain.IdentificationDocTypeNationalid,
					IdentificationDocNumber:         "12345678",
					IdentificationDocNumberUploadID: "23456789",
				},
				KRAPIN:                         "A0123456",
				KRAPINUploadID:                 "34567890",
				DrivingLicenseID:               "12345678",
				CertificateGoodConductUploadID: "34567890",
			},
			wantErr: false,
		}, {
			name: "Sad Case - Attempt adding rider KYC with invalid details",
			args: args{
				ctx: ctx,
				input: domain.IndividualRider{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           "RANDOM STRING",
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "23456789",
					},
					KRAPIN:                         "A0123456",
					KRAPINUploadID:                 "34567890",
					DrivingLicenseID:               "12345678",
					CertificateGoodConductUploadID: "34567890",
				},
			},
			wantErr:     true,
			expectedErr: exceptions.WrongEnumTypeError("RANDOM STRING").Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := s
			got, err := service.Supplier.AddIndividualRiderKyc(tt.args.ctx, tt.args.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddIndividualRiderKYC() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("AddIndividualRiderKYC() error = %v, expectedErr %v", err, tt.expectedErr)
					return
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AddIndividualRiderKYC() = %v, want %v", got, tt.want)
				}
				return
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
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	name := "Makmende"
	partnerRider := base.PartnerTypeRider
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}
	type args struct {
		ctx   context.Context
		input domain.OrganizationRider
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.OrganizationRider
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.OrganizationRider{
					OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:                             "KRA-12345678",
					KRAPINUploadID:                     "KRA-UPLOAD-12345678",
					CertificateOfIncorporation:         "CERT-12345",
					CertificateOfInCorporationUploadID: "CERT-UPLOAD-1234",
					OrganizationCertificate:            "ORG-12345",
				},
			},
			want: &domain.OrganizationRider{
				OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
				DirectorIdentifications: []domain.Identification{
					{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
				},
				KRAPIN:                             "KRA-12345678",
				KRAPINUploadID:                     "KRA-UPLOAD-12345678",
				CertificateOfIncorporation:         "CERT-12345",
				CertificateOfInCorporationUploadID: "CERT-UPLOAD-1234",
				OrganizationCertificate:            "ORG-12345",
			},
			wantErr: false,
		},
		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: domain.OrganizationRider{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []domain.Identification{
						{
							IdentificationDocType:           domain.IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:                             "KRA-12345678",
					KRAPINUploadID:                     "KRA-UPLOAD-12345678",
					CertificateOfIncorporation:         "CERT-12345",
					CertificateOfInCorporationUploadID: "CERT-UPLOAD-1234",
					OrganizationCertificate:            "ORG-12345",
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddOrganizationRiderKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationRiderKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationRiderKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationRiderKyc() = %v, want %v", got, tt.want)
				}
				return
			}

		})
	}
}

func TestSupplierUseCasesImpl_FetchKYCProcessingRequests(t *testing.T) {
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

	name := "Makmende"
	partnerRider := base.PartnerTypeRider

	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	riderKYC := domain.OrganizationRider{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678",
				IdentificationDocNumberUploadID: "12345678",
			},
		},
		KRAPIN:                             "KRA-12345678",
		KRAPINUploadID:                     "KRA-UPLOAD-12345678",
		CertificateOfIncorporation:         "CERT-12345",
		CertificateOfInCorporationUploadID: "CERT-UPLOAD-1234",
		OrganizationCertificate:            "ORG-12345",
	}
	_, err = s.Supplier.AddOrganizationRiderKyc(ctx, riderKYC)
	if err != nil {
		t.Errorf("can't create KYC for a rider's organisation")
		return
	}

	// get supplier after KYC is added above
	supplier, err := s.Supplier.FindSupplierByUID(ctx)
	if err != nil {
		t.Errorf("cannot get supplier")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		want     []*domain.KYCRequest
		supplier *base.Supplier
		wantErr  bool
	}{
		{
			name: "successful fetch single KYC request",
			args: args{
				ctx: ctx,
			},
			supplier: supplier,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.FetchKYCProcessingRequests(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FetchKYCProcessingRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, request := range got {
					if request.ID == "" {
						t.Errorf("KYC request should have ID")
						return
					}
					if request.Processed != false {
						t.Errorf("SupplierUseCasesImpl.FetchKYCProcessingRequests() = %v, want %v", request.Processed, false)
						return
					}
					if request.Status != domain.KYCProcessStatusPending {
						t.Errorf("SupplierUseCasesImpl.FetchKYCProcessingRequests() = %v, want %v", request.Status, domain.KYCProcessStatusPending)
						return
					}
				}
			}
		})
	}
}

func TestAddOrganizationCoachKyc(t *testing.T) {

	/*
	 * Run tests
	 */
	test1 := domain.OrganizationCoach{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "SOMEKRAPIN",
	}
	test2 := domain.OrganizationCoach{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "someKraPin",
	}
	test2Want := domain.OrganizationCoach{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		KRAPIN:               "someKraPin",
	}
	tests := []struct {
		name    string
		coach   domain.OrganizationCoach
		want    domain.OrganizationCoach
		wantErr bool
	}{
		{
			name:    "valid case",
			coach:   test1,
			want:    test1,
			wantErr: false,
		},
		{
			name:    "valid case2",
			coach:   test2,
			want:    test2Want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*
			 * create a supplier account.
			 */
			ctx := context.Background()
			service, err := InitializeTestService(ctx)
			if err != nil {
				t.Errorf("failed to create service")
				return
			}

			seed := rand.NewSource(time.Now().UnixNano())
			unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
			testPhoneNumber := "+25475698" + unique
			testPhoneNumberPin := "7463"

			// first create a user account
			newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			coachResponse, err := service.Supplier.AddOrganizationCoachKyc(newCtx, tt.coach)
			if tt.wantErr && err == nil {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("error was expected but got no error")
				return
			}

			if err != nil && !tt.wantErr {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("error was not expected but got: %v", err)
				return
			}

			// check data matches expected

			if coachResponse.OrganizationTypeName != tt.want.OrganizationTypeName {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("wanted: %v, got: %v", tt.want.OrganizationTypeName, coachResponse.OrganizationTypeName)
				return
			}

			if coachResponse.KRAPIN != tt.want.KRAPIN {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("wanted: %v, got: %v", tt.want.KRAPIN, coachResponse.KRAPIN)
				return
			}

			if coachResponse.KRAPINUploadID != tt.want.KRAPINUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, coachResponse.KRAPINUploadID)
				return
			}

			for index, document := range coachResponse.SupportingDocumentsUploadID {
				if document != tt.want.SupportingDocumentsUploadID[index] {
					clean(newCtx, testPhoneNumber, t, service)
					t.Errorf("wanted: %v, got: %v", tt.want.SupportingDocumentsUploadID[index], document)
					return
				}
			}

			if coachResponse.CertificateOfIncorporation != tt.want.CertificateOfIncorporation {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("wanted: %v, got: %v", tt.want.CertificateOfIncorporation, coachResponse.CertificateOfIncorporation)
				return
			}

			if coachResponse.CertificateOfInCorporationUploadID != tt.want.CertificateOfInCorporationUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("wanted: %v, got: %v", tt.want.CertificateOfInCorporationUploadID, coachResponse.CertificateOfInCorporationUploadID)
				return
			}

			/*
			 * delete the created account and its data
			 */
			clean(newCtx, testPhoneNumber, t, service)
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

func createUserTestAccount(ctx context.Context, service *interactor.Interactor,
	testPhoneNumber string, testPhoneNumberPin string, flavour base.Flavour, t *testing.T) (context.Context, error) {
	// try do clean up first
	_ = service.Signup.RemoveUserByPhoneNumber(ctx, testPhoneNumber)
	otp, err := generateTestOTP(t, testPhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a test OTP: %v", err)
	}

	response, err := service.Signup.CreateUserByPhone(
		ctx,
		&resources.SignUpInput{
			PhoneNumber: &testPhoneNumber,
			PIN:         &testPhoneNumberPin,
			Flavour:     flavour,
			OTP:         &otp.OTP,
		},
	)
	if err != nil {
		t.Errorf("failed to create a user error: %v", err)
		return nil, err
	}

	// get context of the created account.
	idTokens, err := base.AuthenticateCustomFirebaseToken(*response.Auth.CustomToken)
	if err != nil {
		t.Errorf("%v", err)
		return nil, err
	}
	authToken, err := base.ValidateBearerToken(ctx, idTokens.IDToken)
	if err != nil {
		t.Errorf("%v", err)
		return nil, err
	}
	newCtx := context.WithValue(ctx, base.AuthTokenContextKey, authToken)
	return newCtx, nil
}

func TestAddIndividualNutritionKYC(t *testing.T) {
	test1ID := uuid.New().String()
	test1NutritionInput := domain.IndividualNutrition{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeMilitary,
			IdentificationDocNumber:         "1111111111",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
	}

	test2NutritionInput := domain.IndividualNutrition{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeMilitary,
			IdentificationDocNumber:         "1111111111",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
	}

	test2NutritionOutPut := domain.IndividualNutrition{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeMilitary,
			IdentificationDocNumber:         "000000",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
	}

	tests := []struct {
		name      string
		nutrition domain.IndividualNutrition
		want      domain.IndividualNutrition
		wantErr   bool
	}{
		{
			name:      "valid case",
			nutrition: test1NutritionInput,
			want:      test1NutritionInput,
			wantErr:   false,
		},
		{
			name:      "invalid case: IdentificationDocNumber different from input",
			nutrition: test2NutritionInput,
			want:      test2NutritionOutPut,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*
			 * create a supplier account.
			 */
			ctx := context.Background()
			service, err := InitializeTestService(ctx)
			if err != nil {
				t.Errorf("failed to create service")
				return
			}

			seed := rand.NewSource(time.Now().UnixNano())
			unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
			testPhoneNumber := "+25475698" + unique
			testPhoneNumberPin := "7463"

			newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			response, err := service.Supplier.AddIndividualNutritionKyc(newCtx, tt.nutrition)
			if err != nil {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("failed to add individual nutritionkyc got error: %v", err)
				return
			}

			// check the data returned is the expected
			if response.IdentificationDoc != tt.want.IdentificationDoc {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.IdentificationDoc, response.IdentificationDoc)
				}
				return
			}

			if response.KRAPIN != tt.want.KRAPIN {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPIN, response.KRAPIN)
				}
				return
			}

			if response.KRAPINUploadID != tt.want.KRAPINUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, response.KRAPINUploadID)
				}
				return
			}

			if len(response.SupportingDocumentsUploadID) != len(tt.want.SupportingDocumentsUploadID) {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", len(response.SupportingDocumentsUploadID), len(tt.want.SupportingDocumentsUploadID))
				}
				return
			}

			if response.PracticeLicenseID != tt.want.PracticeLicenseID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseID, response.PracticeLicenseID)
				}
				return
			}

			if response.PracticeLicenseUploadID != tt.want.PracticeLicenseUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseUploadID, response.PracticeLicenseUploadID)
				}
				return
			}

			// do clean up
			clean(newCtx, testPhoneNumber, t, service)
		})
	}
}

func TestSupplierUseCasesImpl_FindSupplierByUID(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
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
			name: "valid: supplier found",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context provided",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.FindSupplierByUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FindSupplierByUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantErr {
				t.Errorf("nil supplier returned")
				return
			}
		})
	}
}

func TestAddOrganizationNutritionKyc(t *testing.T) {
	ctx := context.Background()
	service, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("failed to create service")
		return
	}
	/*
	 * create a supplier account.
	 */

	seed := rand.NewSource(time.Now().UnixNano())
	unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
	testPhoneNumber := "+25475698" + unique
	testPhoneNumberPin := "7463"

	newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	test1Input := domain.OrganizationNutrition{
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "someKRAPIN",
		KRAPINUploadID:                     "KRAPINUploadID",
		SupportingDocumentsUploadID:        []string{"SupportingDocumentsUploadID", "Support"},
		CertificateOfIncorporation:         "CertificateOfIncorporation",
		CertificateOfInCorporationUploadID: "CertificateOfInCorporationUploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeMilitary,
				IdentificationDocNumber:         "IdentificationDocNumber",
				IdentificationDocNumberUploadID: "IdentificationDocNumberUploadID",
			},
		},
		OrganizationCertificate: "OrganizationCertificate",
		RegistrationNumber:      "RegistrationNumber",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
	}
	test2Input := test1Input
	test2OutPut := test2Input
	test2OutPut.CertificateOfInCorporationUploadID = " some "
	tests := []struct {
		name    string
		input   domain.OrganizationNutrition
		want    domain.OrganizationNutrition
		wantErr bool
	}{
		{
			name:  "valid case",
			input: test1Input,
			want:  test1Input,
		},
		{
			name:  "invalid case",
			input: test2Input,
			want:  test2Input,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// run the function being tested
			response, err := service.Supplier.AddOrganizationNutritionKyc(newCtx, tt.input)
			if err != nil {
				t.Errorf("failed to add organization nutrition kyc, returned error: %v", err)
				return
			}

			// validate response

			if response.OrganizationTypeName != tt.want.OrganizationTypeName && tt.wantErr {

				t.Errorf("wanted: %v, got: %v", tt.want.OrganizationTypeName, response.OrganizationTypeName)
				return
			}

			if response.KRAPIN != tt.want.KRAPIN && tt.wantErr {

				t.Errorf("wanted: %v, got: %v", tt.want.KRAPIN, response.KRAPIN)
				return
			}

			if response.KRAPINUploadID != tt.want.KRAPINUploadID && tt.wantErr {

				t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, response.KRAPINUploadID)
				return
			}

			if len(response.SupportingDocumentsUploadID) != len(tt.want.SupportingDocumentsUploadID) {

				t.Errorf("wanted: %v, got: %v", len(tt.want.SupportingDocumentsUploadID), len(response.SupportingDocumentsUploadID))
				return
			}

			if response.CertificateOfIncorporation != tt.want.CertificateOfIncorporation {

				t.Errorf("wanted: %v, got: %v", tt.want.CertificateOfIncorporation, response.CertificateOfIncorporation)
				return
			}

			if response.CertificateOfInCorporationUploadID != tt.want.CertificateOfInCorporationUploadID {

				t.Errorf("wanted: %v, got: %v", tt.want.CertificateOfInCorporationUploadID, response.CertificateOfInCorporationUploadID)
				return
			}

			if len(response.DirectorIdentifications) != len(tt.want.DirectorIdentifications) {

				t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, response.KRAPINUploadID)
				return
			}

			if response.OrganizationCertificate != tt.want.OrganizationCertificate {

				t.Errorf("wanted: %v, got: %v", tt.want.OrganizationCertificate, response.OrganizationCertificate)
				return
			}

			if response.RegistrationNumber != tt.want.RegistrationNumber {

				t.Errorf("wanted: %v, got: %v", tt.want.RegistrationNumber, response.RegistrationNumber)
				return
			}

			if response.PracticeLicenseID != tt.want.PracticeLicenseID {

				t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseID, response.PracticeLicenseID)
				return
			}

			if response.PracticeLicenseUploadID != tt.want.PracticeLicenseUploadID {

				t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseUploadID, response.PracticeLicenseUploadID)

				return
			}

		})
	}
	// clean up
	clean(newCtx, testPhoneNumber, t, service)
}

func TestSupplierUseCasesImpl_AddOrganizationPractitionerKyc(t *testing.T) {
	test1ID := uuid.New().String()
	test1PractitionerKYC := domain.OrganizationPractitioner{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678",
				IdentificationDocNumberUploadID: test1ID,
			},
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	test2PractitionerKYC := domain.OrganizationPractitioner{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "12345678",
				IdentificationDocNumberUploadID: test1ID,
			},
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	test2PractitionerOutPut := domain.OrganizationPractitioner{
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "0000000",
				IdentificationDocNumberUploadID: test1ID,
			},
		},
		KRAPIN:                      test1ID,
		KRAPINUploadID:              test1ID,
		SupportingDocumentsUploadID: []string{test1ID, strings.ToUpper(test1ID)},
		PracticeLicenseID:           test1ID,
		PracticeLicenseUploadID:     test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	tests := []struct {
		name         string
		practitioner domain.OrganizationPractitioner
		want         domain.OrganizationPractitioner
		wantErr      bool
	}{
		{
			name:         "valid case",
			practitioner: test1PractitionerKYC,
			want:         test1PractitionerKYC,
			wantErr:      false,
		},
		{
			name:         "invalid case: IdentificationDocNumber different from input",
			practitioner: test2PractitionerKYC,
			want:         test2PractitionerOutPut,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*
			 * create a supplier account.
			 */
			ctx := context.Background()
			service, err := InitializeTestService(ctx)
			if err != nil {
				t.Errorf("failed to create service")
				return
			}

			seed := rand.NewSource(time.Now().UnixNano())
			unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
			testPhoneNumber := "+25475698" + unique
			testPhoneNumberPin := "7463"

			newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			response, err := service.Supplier.AddOrganizationPractitionerKyc(newCtx, tt.practitioner)
			if err != nil {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("failed to add organizational practitionerkyc got error: %v", err)
				return
			}

			// check the data returned is the expected
			if len(response.DirectorIdentifications) != len(tt.want.DirectorIdentifications) {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", len(tt.want.DirectorIdentifications), len(response.DirectorIdentifications))
				}
				return
			}

			if response.KRAPIN != tt.want.KRAPIN {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPIN, response.KRAPIN)
				}
				return
			}

			if response.KRAPINUploadID != tt.want.KRAPINUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, response.KRAPINUploadID)
				}
				return
			}

			if len(response.SupportingDocumentsUploadID) != len(tt.want.SupportingDocumentsUploadID) {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", len(response.SupportingDocumentsUploadID), len(tt.want.SupportingDocumentsUploadID))
				}
				return
			}

			if response.PracticeLicenseID != tt.want.PracticeLicenseID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseID, response.PracticeLicenseID)
				}
				return
			}

			if response.PracticeLicenseUploadID != tt.want.PracticeLicenseUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseUploadID, response.PracticeLicenseUploadID)
				}
				return
			}

			// do clean up
			clean(newCtx, testPhoneNumber, t, service)
		})
	}
}

func TestSupplierUseCasesImpl_AddIndividualPractitionerKyc(t *testing.T) {
	test1ID := uuid.New().String()
	test1PractitionerInput := domain.IndividualPractitioner{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:             "KRA-12345678",
		KRAPINUploadID:     test1ID,
		RegistrationNumber: "REG-12345",
		PracticeLicenseID:  test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	test2PractitionerInput := domain.IndividualPractitioner{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "111111",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:             "KRA-12345678",
		KRAPINUploadID:     test1ID,
		RegistrationNumber: "REG-12345",
		PracticeLicenseID:  test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	test2PractitionerOutput := domain.IndividualPractitioner{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "000000",
			IdentificationDocNumberUploadID: test1ID,
		},
		KRAPIN:             "KRA-12345678",
		KRAPINUploadID:     test1ID,
		RegistrationNumber: "REG-12345",
		PracticeLicenseID:  test1ID,
		PracticeServices: []domain.PractitionerService{
			domain.PractitionerServiceOutpatientServices,
			domain.PractitionerServiceInpatientServices,
			domain.PractitionerServiceOther,
		},
		Cadre: domain.PractitionerCadreDoctor,
	}

	tests := []struct {
		name         string
		practitioner domain.IndividualPractitioner
		want         domain.IndividualPractitioner
		wantErr      bool
	}{
		{
			name:         "valid case",
			practitioner: test1PractitionerInput,
			want:         test1PractitionerInput,
			wantErr:      false,
		},
		{
			name:         "invalid case: IdentificationDocNumber different from input",
			practitioner: test2PractitionerInput,
			want:         test2PractitionerOutput,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*
			 * create a supplier account.
			 */
			ctx := context.Background()
			service, err := InitializeTestService(ctx)
			if err != nil {
				t.Errorf("failed to create service")
				return
			}

			seed := rand.NewSource(time.Now().UnixNano())
			unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
			testPhoneNumber := "+25475698" + unique
			testPhoneNumberPin := "7463"

			newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			response, err := service.Supplier.AddIndividualPractitionerKyc(newCtx, tt.practitioner)
			if err != nil {
				clean(newCtx, testPhoneNumber, t, service)
				t.Errorf("failed to add individual practitionerkyc got error: %v", err)
				return
			}

			// check the data returned is the expected
			if response.IdentificationDoc != tt.want.IdentificationDoc {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.IdentificationDoc, response.IdentificationDoc)
				}
				return
			}

			if response.KRAPIN != tt.want.KRAPIN {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPIN, response.KRAPIN)
				}
				return
			}

			if response.KRAPINUploadID != tt.want.KRAPINUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.KRAPINUploadID, response.KRAPINUploadID)
				}
				return
			}

			if len(response.SupportingDocumentsUploadID) != len(tt.want.SupportingDocumentsUploadID) {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", len(response.SupportingDocumentsUploadID), len(tt.want.SupportingDocumentsUploadID))
				}
				return
			}

			if response.PracticeLicenseID != tt.want.PracticeLicenseID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseID, response.PracticeLicenseID)
				}
				return
			}

			if response.PracticeLicenseUploadID != tt.want.PracticeLicenseUploadID {
				clean(newCtx, testPhoneNumber, t, service)
				if !tt.wantErr {
					t.Errorf("wanted: %v, got: %v", tt.want.PracticeLicenseUploadID, response.PracticeLicenseUploadID)
				}
				return
			}

			// do clean up
			clean(newCtx, testPhoneNumber, t, service)
		})
	}
}

func TestSupplierUseCasesImpl_FetchSupplierAllowedLocations(t *testing.T) {
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

	name := "Makmende"
	partnerRider := base.PartnerTypeRider

	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	_, err = s.Supplier.SupplierEDILogin(ctx, testEDIPortalUsername, testEDIPortalPassword, testSladeCode)
	if err != nil {
		t.Errorf("can't perform supplier edi login: %v", err)
		return
	}

	cmParentOrgId := testChargeMasterParentOrgId
	filter := []*resources.BranchFilterInput{
		{
			ParentOrganizationID: &cmParentOrgId,
		},
	}

	brs, err := s.ChargeMaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		t.Errorf("can't find branch")
		return
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.BranchConnection
		wantErr bool
	}{
		{
			name: "valid: supplier allowed locations found",
			args: args{
				ctx: ctx,
			},
			want:    brs,
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context provided",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.FetchSupplierAllowedLocations(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FetchSupplierAllowedLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SupplierUseCasesImpl.FetchSupplierAllowedLocations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupplierSetDefaultLocation(t *testing.T) {
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

	name := "Makmende"
	partnerPractitioner := base.PartnerTypePractitioner
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, base.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	_, err = s.Supplier.SupplierEDILogin(ctx, testEDIPortalUsername, testEDIPortalPassword, testSladeCode)
	if err != nil {
		t.Errorf("unable to login user")
		return
	}

	cmParentOrgId := testChargeMasterParentOrgId
	filter := []*resources.BranchFilterInput{
		{
			ParentOrganizationID: &cmParentOrgId,
		},
	}

	_, err = s.ChargeMaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		t.Errorf("can't find branch")
		return
	}

	type args struct {
		ctx        context.Context
		locationID string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Sad case - Set default location with an invalid locationID",
			args: args{
				ctx:        ctx,
				locationID: "fdd",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Set default location with an empty locationID",
			args: args{
				ctx:        ctx,
				locationID: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy case - Set default location with a valid locationID",
			args: args{
				ctx:        ctx,
				locationID: testChargeMasterBranchID,
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := s
			got, err := service.Supplier.SupplierSetDefaultLocation(tt.args.ctx, tt.args.locationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierSetDefaultLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SupplierUseCasesImpl.SupplierSetDefaultLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessKYCRequest(t *testing.T) {
	ctx := context.Background()
	service, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("failed to create service")
		return
	}
	supplier := service.Supplier

	/*
	 * create a supplier account.
	 */

	seed := rand.NewSource(time.Now().UnixNano())
	unique := fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9)) + fmt.Sprint(rand.New(seed).Intn(9))
	testPhoneNumber := "+25475698" + unique
	testPhoneNumberPin := "7463"

	newCtx, err := createUserTestAccount(ctx, service, testPhoneNumber, testPhoneNumberPin, base.FlavourPro, t)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	partnerName := "jubileeisnotinsurance"
	partnerType := base.PartnerTypeNutrition

	_, err = supplier.AddPartnerType(newCtx, &partnerName, &partnerType)
	if err != nil {
		t.Errorf("failed to add partner type, error %v", err)
		return
	}
	_, err = supplier.SetUpSupplier(newCtx, base.AccountTypeIndividual)

	if err != nil {
		t.Errorf("failed to add partner type, error %v", err)
		return
	}

	test1Input := domain.IndividualNutrition{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
	}

	_, err = supplier.AddIndividualNutritionKyc(newCtx, test1Input)
	if err != nil {
		t.Errorf("failed to add organization nutrition kyc, returned error: %v", err)
		clean(newCtx, testPhoneNumber, t, service)
		return
	}

	tests := []struct {
		name string
	}{
		{
			name: "valid case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kycrequests, err := supplier.FetchKYCProcessingRequests(newCtx)
			if err != nil {
				t.Errorf("failed to fetch kyc requests")
				clean(newCtx, testPhoneNumber, t, service)
				return
			}
			firstKYC := kycrequests[0]

			/* validate data */
			if firstKYC == nil {
				t.Errorf("nil kyc returned")
				clean(newCtx, testPhoneNumber, t, service)
				return
			}

			reason := "some reason"
			response, err := supplier.ProcessKYCRequest(newCtx, firstKYC.ID, domain.KYCProcessStatusApproved, &reason)
			if err != nil {
				t.Errorf("failed to process kyc requests: %v", err)
				clean(newCtx, testPhoneNumber, t, service)
				return
			}

			if !response {
				t.Errorf("%v", err)
				clean(newCtx, testPhoneNumber, t, service)
				return
			}
		})
	}
	clean(newCtx, testPhoneNumber, t, service)
}
