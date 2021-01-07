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

// import (
// 	"context"
// 	"testing"

// 	"github.com/brianvoe/gofakeit"
// 	"github.com/stretchr/testify/assert"
// 	"gitlab.slade360emr.com/go/base"
// )

// func TestService_AddSupplier(t *testing.T) {
// 	service := NewService()
// 	ctx := base.GetAuthenticatedContext(t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	type args struct {
// 		ctx         context.Context
// 		name        string
// 		partnerType PartnerType
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy case: successfully create a supplier",
// 			args: args{
// 				ctx:         ctx,
// 				name:        "Be.Well Test Supplier",
// 				partnerType: PartnerTypeProvider,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "sad case:add supplier without basic partner details",
// 			args: args{
// 				ctx: context.Background(),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "sad case: add supplier with the wrong partner type",
// 			args: args{
// 				ctx:         context.Background(),
// 				partnerType: "not a valid partner type",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			supplier, err := s.AddSupplier(tt.args.ctx, tt.args.name, tt.args.partnerType)
// 			if err == nil {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, supplier)
// 			}
// 			if err != nil {
// 				assert.Nil(t, supplier)
// 				assert.NotNil(t, err)
// 			}
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.AddSupplier() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestService_FindSupplier(t *testing.T) {
// 	service := NewService()
// 	ctx, token := base.GetAuthenticatedContextAndToken(t)
// 	if token == nil {
// 		t.Errorf("nil token")
// 		return
// 	}
// 	if ctx == nil {
// 		t.Errorf("nil context")
// 		return
// 	}

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	supplier, err := service.AddSupplier(ctx, name, PartnerTypeProvider)
// 	if err != nil {
// 		t.Errorf("can't add supplier: %v", err)
// 		return
// 	}
// 	if supplier == nil {
// 		t.Errorf("nil supplier after adding a supplier")
// 		return
// 	}

// 	type args struct {
// 		ctx context.Context
// 		uid string
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantErr     bool
// 		expectedErr string
// 	}{
// 		{
// 			name: "valid : authenticated context",
// 			args: args{
// 				ctx: ctx,
// 				uid: token.UID,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "invalid : unautheniticated context",
// 			args: args{
// 				ctx: context.Background(),
// 				uid: "not a uid",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			supplier, err := s.FindSupplier(tt.args.ctx, tt.args.uid)

// 			if tt.wantErr && (err == nil) {
// 				t.Errorf("expected an error. got %v", err)
// 			}

// 			if tt.wantErr {
// 				assert.Nil(t, supplier)
// 			}

// 			if !tt.wantErr {
// 				assert.NotNil(t, supplier)
// 				assert.Nil(t, err)
// 			}

// 		})
// 	}
// }

// func TestService_SuspendSupplier(t *testing.T) {
// 	service := NewService()
// 	ctx, token := createNewUser(context.Background(), t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	_, errS := service.AddSupplier(ctx, name, PartnerTypeProvider)
// 	if errS != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}
// 	type args struct {
// 		ctx context.Context
// 		uid string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "sad case: suspend a nonexisting supplier",
// 			args: args{
// 				ctx: context.Background(),
// 				uid: "some random uid",
// 			},
// 			want:    false,
// 			wantErr: true,
// 		},
// 		{
// 			name: "Happy case: suspend an existing supplier",
// 			args: args{
// 				ctx: ctx,
// 				uid: token.UID,
// 			},
// 			want:    true,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			got, err := s.SuspendSupplier(tt.args.ctx, tt.args.uid)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.SuspendSupplier() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("Service.SuspendSupplier() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestService_SetUpSupplier(t *testing.T) {
// 	service := NewService()
// 	ctx, _ := base.GetAuthenticatedContextAndToken(t)

// 	partnerName := gofakeit.Name()
// 	partnerType := PartnerTypePractitioner
// 	service.AddPartnerType(ctx, &partnerName, &partnerType)

// 	type args struct {
// 		ctx   context.Context
// 		input AccountType
// 	}

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successful basic individual account set up",
// 			args: args{
// 				ctx:   ctx,
// 				input: "INDIVIDUAL",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Successful basic organisation account set up",
// 			args: args{
// 				ctx:   ctx,
// 				input: "ORGANISATION",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			_, err := s.SetUpSupplier(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.SetUpSupplier() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestService_PublishKYCNudge(t *testing.T) {
// 	service := NewService()

// 	_, token := base.GetAuthenticatedContextAndToken(t)

// 	type args struct {
// 		uid     string
// 		partner PartnerType
// 		account AccountType
// 	}

// 	tests := []struct {
// 		name           string
// 		args           args
// 		wantErr        bool
// 		expectedErr    string
// 		invalidService bool
// 	}{
// 		{
// 			name: "valid : Individual Rider KYC Nudge",
// 			args: args{
// 				uid:     token.UID,
// 				partner: PartnerTypeRider,
// 				account: AccountTypeIndividual,
// 			},
// 			wantErr:        false,
// 			invalidService: false,
// 		},
// 		{
// 			name: "valid : Organization Practitioner KYC Nudge",
// 			args: args{
// 				uid:     token.UID,
// 				partner: PartnerTypePractitioner,
// 				account: AccountTypeOrganisation,
// 			},
// 			wantErr:        false,
// 			invalidService: false,
// 		},

// 		{
// 			name: "invalid : unknown partner type",
// 			args: args{
// 				uid:     token.UID,
// 				partner: "alien partner",
// 				account: AccountTypeOrganisation,
// 			},
// 			wantErr:        true,
// 			invalidService: false,
// 			expectedErr:    "expected `partner` to be defined and to be valid",
// 		},
// 		{
// 			name: "invalid : consumer partner",
// 			args: args{
// 				uid:     token.UID,
// 				partner: PartnerTypeConsumer,
// 				account: AccountTypeOrganisation,
// 			},
// 			wantErr:        true,
// 			invalidService: false,
// 			expectedErr:    "invalid `partner`. cannot use CONSUMER in this context",
// 		},
// 		{
// 			name: "invalid : unknown account type",
// 			args: args{
// 				uid:     token.UID,
// 				partner: PartnerTypePractitioner,
// 				account: "alien account",
// 			},
// 			wantErr:        true,
// 			invalidService: false,
// 			expectedErr:    "provided `account` is not valid",
// 		},

// 		{
// 			name: "invalid : wrong engagement service",
// 			args: args{
// 				uid:     token.UID,
// 				partner: PartnerTypePractitioner,
// 				account: AccountTypeOrganisation,
// 			},
// 			wantErr:        true,
// 			invalidService: true,
// 			expectedErr:    "unable to publish kyc nudge. unexpected status code  404",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if !tt.invalidService {
// 				s := service
// 				err := s.PublishKYCNudge(tt.args.uid, &tt.args.partner, &tt.args.account)
// 				if !tt.wantErr {
// 					assert.Nil(t, err)
// 				}

// 				if tt.wantErr {
// 					assert.NotNil(t, err)
// 					assert.Contains(t, tt.expectedErr, err.Error())
// 				}
// 				return
// 			}

// 			is := service
// 			is.engagement = service.Otp
// 			err := is.PublishKYCNudge(tt.args.uid, &tt.args.partner, &tt.args.account)
// 			assert.NotNil(t, err)
// 			assert.Contains(t, tt.expectedErr, err.Error())

// 		})
// 	}
// }

// func TestService_AddIndividualRiderKyc(t *testing.T) {
// 	service := NewService()
// 	ctx, _ := base.GetAuthenticatedContextAndToken(t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	_, err = service.SetUpSupplier(ctx, AccountTypeIndividual)
// 	if err != nil {
// 		t.Errorf("can't set up a supplier")
// 		return
// 	}

// 	riderInput := IndividualRider{
// 		IdentificationDoc: Identification{
// 			IdentificationDocType:           IdentificationDocTypeNationalid,
// 			IdentificationDocNumber:         "12345678",
// 			IdentificationDocNumberUploadID: "12345678",
// 		},
// 		KRAPIN:                         "12345678",
// 		KRAPINUploadID:                 "12345678",
// 		DrivingLicenseUploadID:         "12345678",
// 		CertificateGoodConductUploadID: "12345678",
// 	}
// 	riderKYC := &IndividualRider{
// 		IdentificationDoc: Identification{
// 			IdentificationDocType:           IdentificationDocTypeNationalid,
// 			IdentificationDocNumber:         "12345678",
// 			IdentificationDocNumberUploadID: "12345678",
// 		},
// 		KRAPIN:                         "12345678",
// 		KRAPINUploadID:                 "12345678",
// 		DrivingLicenseUploadID:         "12345678",
// 		CertificateGoodConductUploadID: "12345678",
// 	}

// 	type args struct {
// 		ctx   context.Context
// 		input IndividualRider
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *IndividualRider
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successful Add individual rider KYC",
// 			args: args{
// 				ctx:   ctx,
// 				input: riderInput,
// 			},
// 			wantErr: false,
// 			want:    riderKYC,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			got, err := s.AddIndividualRiderKyc(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Service.AddIndividualRiderKyc() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Service.AddIndividualRiderKyc() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestService_AddOrganizationRiderKyc(t *testing.T) {
// 	service := NewService()
// 	ctx, _ := base.GetAuthenticatedContextAndToken(t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
// 	if err != nil {
// 		t.Errorf("can't set up a supplier")
// 		return
// 	}

// 	type args struct {
// 		ctx      context.Context
// 		input    OrganizationRider
// 		resource OrganizationRider
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantErr     bool
// 		expectedErr string
// 	}{
// 		{
// 			name: "valid : should pass",
// 			args: args{
// 				ctx: ctx,
// 				input: OrganizationRider{
// 					OrganizationTypeName: OrganizationTypeLimitedCompany,
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
// 					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
// 					OrganizationCertificate:            "ORG-CERT",
// 					KRAPIN:                             "KRA-PIN-12345678",
// 					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
// 					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
// 				},
// 				resource: OrganizationRider{
// 					OrganizationTypeName: OrganizationTypeLimitedCompany,
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
// 					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
// 					OrganizationCertificate:            "ORG-CERT",
// 					KRAPIN:                             "KRA-PIN-12345678",
// 					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
// 					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "invalid : organization type name",
// 			args: args{
// 				ctx: ctx,
// 				input: OrganizationRider{
// 					OrganizationTypeName: "AWESOME ORG",
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
// 					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
// 					OrganizationCertificate:            "ORG-CERT",
// 					KRAPIN:                             "KRA-PIN-12345678",
// 					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
// 					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
// 				},
// 			},
// 			wantErr:     true,
// 			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			got, err := s.AddOrganizationRiderKyc(tt.args.ctx, tt.args.input)
// 			if !tt.wantErr {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, got)
// 				if !reflect.DeepEqual(*got, tt.args.resource) {
// 					t.Errorf("Service.AddOrganizationRiderKyc() = %v, want %v", got, tt.args.resource)
// 				}
// 				return
// 			}
// 			assert.NotNil(t, err)
// 			assert.Nil(t, got)
// 			assert.Contains(t, err.Error(), tt.expectedErr)
// 		})
// 	}
// }

// func TestService_AddIndividualPractitionerKyc(t *testing.T) {
// 	service := NewService()
// 	ctx, _ := base.GetAuthenticatedContextAndToken(t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
// 	if err != nil {
// 		t.Errorf("can't set up a supplier")
// 		return
// 	}

// 	type args struct {
// 		ctx      context.Context
// 		input    IndividualPractitioner
// 		resource IndividualPractitioner
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantErr     bool
// 		expectedErr string
// 	}{
// 		{
// 			name: "valid : should pass",
// 			args: args{
// 				ctx: ctx,
// 				input: IndividualPractitioner{
// 					IdentificationDoc: Identification{
// 						IdentificationDocType:           IdentificationDocTypeNationalid,
// 						IdentificationDocNumber:         "12345678",
// 						IdentificationDocNumberUploadID: "12345678",
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 				resource: IndividualPractitioner{
// 					IdentificationDoc: Identification{
// 						IdentificationDocType:           IdentificationDocTypeNationalid,
// 						IdentificationDocNumber:         "12345678",
// 						IdentificationDocNumberUploadID: "12345678",
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "invalid : practice services",
// 			args: args{
// 				ctx: ctx,
// 				input: IndividualPractitioner{
// 					IdentificationDoc: Identification{
// 						IdentificationDocType:           IdentificationDocTypeNationalid,
// 						IdentificationDocNumber:         "12345678",
// 						IdentificationDocNumberUploadID: "12345678",
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{"SUPPORTING"},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 			},
// 			wantErr:     true,
// 			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			got, err := s.AddIndividualPractitionerKyc(tt.args.ctx, tt.args.input)
// 			if !tt.wantErr {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, got)
// 				if !reflect.DeepEqual(*got, tt.args.resource) {
// 					t.Errorf("Service.AddIndividualPractitionerKyc() = %v, want %v", got, tt.args.resource)
// 				}
// 				return
// 			}
// 			assert.NotNil(t, err)
// 			assert.Nil(t, got)
// 			assert.Contains(t, err.Error(), tt.expectedErr)
// 		})
// 	}
// }

// func TestService_AddOrganizationPractitionerKyc(t *testing.T) {
// 	service := NewService()
// 	ctx, _ := base.GetAuthenticatedContextAndToken(t)

// 	name := gofakeit.Name()
// 	partnerRider := PartnerTypeRider
// 	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
// 	if err != nil {
// 		t.Errorf("can't create a supplier")
// 		return
// 	}

// 	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
// 	if err != nil {
// 		t.Errorf("can't set up a supplier")
// 		return
// 	}

// 	type args struct {
// 		ctx      context.Context
// 		input    OrganizationPractitioner
// 		resource OrganizationPractitioner
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantErr     bool
// 		expectedErr string
// 	}{
// 		{
// 			name: "valid : should pass",
// 			args: args{
// 				ctx: ctx,
// 				input: OrganizationPractitioner{
// 					OrganizationTypeName: OrganizationTypeLimitedCompany,
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 				resource: OrganizationPractitioner{
// 					OrganizationTypeName: OrganizationTypeLimitedCompany,
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "invalid : organization type name ",
// 			args: args{
// 				ctx: ctx,
// 				input: OrganizationPractitioner{
// 					OrganizationTypeName: "AWESOME ORG",
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{"SUPPORTING"},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 			},
// 			wantErr:     true,
// 			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
// 		},

// 		{
// 			name: "invalid : practice services",
// 			args: args{
// 				ctx: ctx,
// 				input: OrganizationPractitioner{
// 					OrganizationTypeName: OrganizationTypeLimitedCompany,
// 					DirectorIdentifications: []Identification{
// 						{
// 							IdentificationDocType:           IdentificationDocTypeNationalid,
// 							IdentificationDocNumber:         "12345678",
// 							IdentificationDocNumberUploadID: "12345678",
// 						},
// 					},
// 					KRAPIN:             "KRA-12345678",
// 					KRAPINUploadID:     "KRA-UPLOAD-12345678",
// 					RegistrationNumber: "REG-12345",
// 					PracticeLicenseID:  "PRAC-12345",
// 					PracticeServices:   []PractitionerService{"SUPPORTING"},
// 					Cadre:              PractitionerCadreDoctor,
// 				},
// 			},
// 			wantErr:     true,
// 			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			got, err := s.AddOrganizationPractitionerKyc(tt.args.ctx, tt.args.input)
// 			if !tt.wantErr {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, got)
// 				if !reflect.DeepEqual(*got, tt.args.resource) {
// 					t.Errorf("Service.AddOrganizationPractitionerKyc() = %v, want %v", got, tt.args.resource)
// 				}
// 				return
// 			}
// 			assert.NotNil(t, err)
// 			assert.Nil(t, got)
// 			assert.Contains(t, err.Error(), tt.expectedErr)
// 		})
// 	}
// }

// func TestService_SendKYCEmail(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		text         string
// 		emailaddress string
// 	}

// 	approvalEmail := generateProcessKYCApprovalEmailTemplate()
// 	rejectionEmail := generateProcessKYCRejectionEmailTemplate()

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Good case: approved",
// 			args: args{
// 				ctx:          context.Background(),
// 				text:         approvalEmail,
// 				emailaddress: base.GenerateRandomEmail(),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Good case: rejected",
// 			args: args{
// 				ctx:          context.Background(),
// 				text:         rejectionEmail,
// 				emailaddress: base.GenerateRandomEmail(),
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewService()
// 			if err := s.SendKYCEmail(tt.args.ctx, tt.args.text, tt.args.emailaddress); (err != nil) != tt.wantErr {
// 				t.Errorf("Service.SendKYCEmail() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestService_PublishKYCFeedItem(t *testing.T) {
// 	service := NewService()
// 	ctx, token := base.GetAuthenticatedContextAndToken(t)

// 	type args struct {
// 		ctx  context.Context
// 		uids []string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Send feed item to user",
// 			args: args{
// 				ctx:  ctx,
// 				uids: []string{token.UID},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service
// 			if err := s.PublishKYCFeedItem(tt.args.ctx, tt.args.uids...); (err != nil) != tt.wantErr {
// 				t.Errorf("Service.PublishKYCFeedItem() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestSupplierUseCasesImpl_AddPartnerType(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	testRiderName := "Test Rider"
	rider := domain.PartnerTypeRider
	testPractitionerName := "Test Practitioner"
	practitioner := domain.PartnerTypePractitioner
	testProviderName := "Test Provider"
	provider := domain.PartnerTypeProvider
	testPharmaceuticalName := "Test Pharmaceutical"
	pharmaceutical := domain.PartnerTypePharmaceutical
	testCoachName := "Test Coach"
	coach := domain.PartnerTypeCoach
	testNutritionName := "Test Nutrition"
	nutrition := domain.PartnerTypeNutrition
	testConsumerName := "Test Consumer"
	consumer := domain.PartnerTypeConsumer

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx         context.Context
		name        *string
		partnerType *domain.PartnerType
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

	individualPartner := domain.AccountTypeIndividual
	organizationPartner := domain.AccountTypeOrganisation

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx         context.Context
		accountType domain.AccountType
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
	partnerProvider := domain.PartnerTypeProvider
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerProvider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeOrganisation)
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
	partnerPharmaceutical := domain.PartnerTypePharmaceutical
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPharmaceutical)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeOrganisation)
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
	partnerPharmaceutical := domain.PartnerTypePharmaceutical
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPharmaceutical)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeIndividual)
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
			expectedErr: "unable to get the logged in user supplier profile: 10: failed to get a user",
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
	partnerCoach := domain.PartnerTypeCoach
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerCoach)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeIndividual)
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
			expectedErr: "unable to get the logged in user supplier profile: 10: failed to get a user",
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
	partnerRider := domain.PartnerTypeRider
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeIndividual)
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
			expectedErr: exceptions.WrongEnumTypeError("RANDOM STRING", nil).Error(),
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

func TestSupplierUseCasesImpl_AddIndividualPractitionerKyc(t *testing.T) {
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
	partnerPractitioner := domain.PartnerTypePractitioner
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeIndividual)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}
	type args struct {
		ctx   context.Context
		input domain.IndividualPractitioner
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.IndividualPractitioner
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.IndividualPractitioner{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices: []domain.PractitionerService{
						domain.PractitionerServiceOutpatientServices,
						domain.PractitionerServiceInpatientServices,
						domain.PractitionerServiceOther,
					},
					Cadre: domain.PractitionerCadreDoctor,
				},
			},
			want: &domain.IndividualPractitioner{
				IdentificationDoc: domain.Identification{
					IdentificationDocType:           domain.IdentificationDocTypeNationalid,
					IdentificationDocNumber:         "12345678",
					IdentificationDocNumberUploadID: "12345678",
				},
				KRAPIN:             "KRA-12345678",
				KRAPINUploadID:     "KRA-UPLOAD-12345678",
				RegistrationNumber: "REG-12345",
				PracticeLicenseID:  "PRAC-12345",
				PracticeServices: []domain.PractitionerService{
					domain.PractitionerServiceOutpatientServices,
					domain.PractitionerServiceInpatientServices,
					domain.PractitionerServiceOther,
				},
				Cadre: domain.PractitionerCadreDoctor,
			},
			wantErr: false,
		},
		{
			name: "invalid : practice services",
			args: args{
				ctx: ctx,
				input: domain.IndividualPractitioner{
					IdentificationDoc: domain.Identification{
						IdentificationDocType:           domain.IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
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
			got, err := s.Supplier.AddIndividualPractitionerKyc(tt.args.ctx, tt.args.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPractitionerKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPractitionerKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddIndividualPractitionerKyc() = %v, want %v", got, tt.want)
				}
				return
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
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	name := "Makmende"
	partnerPractitioner := domain.PartnerTypePractitioner
	_, err = s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = s.Supplier.SetUpSupplier(ctx, domain.AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}
	type args struct {
		ctx   context.Context
		input domain.OrganizationPractitioner
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.OrganizationPractitioner
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPractitioner{
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
					PracticeServices: []domain.PractitionerService{
						domain.PractitionerServiceOutpatientServices,
						domain.PractitionerServiceInpatientServices,
						domain.PractitionerServiceOther,
					},
					Cadre: domain.PractitionerCadreDoctor,
				},
			},
			want: &domain.OrganizationPractitioner{
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
				PracticeServices: []domain.PractitionerService{
					domain.PractitionerServiceOutpatientServices,
					domain.PractitionerServiceInpatientServices,
					domain.PractitionerServiceOther,
				},
				Cadre: domain.PractitionerCadreDoctor,
			},
			wantErr: false,
		},

		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: domain.OrganizationPractitioner{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.AddOrganizationPractitionerKyc(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPractitionerKyc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPractitionerKyc() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SupplierUseCasesImpl.AddOrganizationPractitionerKyc() = %v, want %v", got, tt.want)
				}
				return
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
		fmt.Print()
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
