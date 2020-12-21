package profile

import (
	"context"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_AddPartnerType(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)

	type args struct {
		ctx         context.Context
		name        string
		partnerType PartnerType
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid: add PartnerTypeRider ",
			args: args{
				ctx:         ctx,
				name:        "Test Rider",
				partnerType: PartnerTypeRider,
			},
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePractitioner ",
			args: args{
				ctx:         ctx,
				name:        "Test Rider",
				partnerType: PartnerTypePractitioner,
			},
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeProvider ",
			args: args{
				ctx:         ctx,
				name:        "Test Provider",
				partnerType: PartnerTypeProvider,
			},
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePharmaceutical",
			args: args{
				ctx:         ctx,
				name:        "Test Pharmaceutical",
				partnerType: PartnerTypePharmaceutical,
			},
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeCoach",
			args: args{
				ctx:         ctx,
				name:        "Test coach",
				partnerType: PartnerTypeCoach,
			},
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeNutrition",
			args: args{
				ctx:         ctx,
				name:        "Test nutrition",
				partnerType: PartnerTypeNutrition,
			},
			wantErr: false,
		},

		{
			name: "invalid: add PartnerTypeConsumer",
			args: args{
				ctx:         ctx,
				name:        "Test consumer",
				partnerType: PartnerTypeConsumer,
			},
			wantErr:     true,
			expectedErr: "invalid `partnerType`. cannot use CONSUMER in this context",
		},

		{
			name: "invalid : invalid context",
			args: args{
				ctx:         context.Background(),
				name:        "Test Rider",
				partnerType: PartnerTypeRider,
			},
			wantErr:     true,
			expectedErr: `unable to get the logged in user: auth token not found in context: unable to get auth token from context with key "UID" `,
		},
		{
			name: "invalid : missing name arg",
			args: args{
				ctx: ctx,
			},
			wantErr:     true,
			expectedErr: "expected `name` to be defined and `partnerType` to be valid",
		},
		{
			name: "invalid : unknownn partner type",
			args: args{
				ctx:         ctx,
				name:        "Test Partner",
				partnerType: "not a valid partner type",
			},
			wantErr:     true,
			expectedErr: "expected `name` to be defined and `partnerType` to be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			resp, err := s.AddPartnerType(tt.args.ctx, &tt.args.name, &tt.args.partnerType)
			if tt.wantErr {
				assert.Equal(t, false, resp)
				assert.NotNil(t, err)
				assert.Contains(t, tt.expectedErr, err.Error())
			}
			if err == nil {
				assert.Nil(t, err)
				assert.Equal(t, true, resp)
			}
		})
	}

}

func TestService_AddSupplier(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	type args struct {
		ctx         context.Context
		name        string
		partnerType PartnerType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successfully create a supplier",
			args: args{
				ctx:         ctx,
				name:        "Be.Well Test Supplier",
				partnerType: PartnerTypeProvider,
			},
			wantErr: false,
		},
		{
			name: "sad case:add supplier without basic partner details",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: add supplier with the wrong partner type",
			args: args{
				ctx:         context.Background(),
				partnerType: "not a valid partner type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			supplier, err := s.AddSupplier(tt.args.ctx, tt.args.name, tt.args.partnerType)
			if err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, supplier)
			}
			if err != nil {
				assert.Nil(t, supplier)
				assert.NotNil(t, err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_FindSupplier(t *testing.T) {
	service := NewService()
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	if token == nil {
		t.Errorf("nil token")
		return
	}
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	supplier, err := service.AddSupplier(ctx, name, PartnerTypeProvider)
	if err != nil {
		t.Errorf("can't add supplier: %v", err)
		return
	}
	if supplier == nil {
		t.Errorf("nil supplier after adding a supplier")
		return
	}

	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : authenticated context",
			args: args{
				ctx: ctx,
				uid: token.UID,
			},
			wantErr: false,
		},
		{
			name: "invalid : unautheniticated context",
			args: args{
				ctx: context.Background(),
				uid: "not a uid",
			},
			wantErr:     true,
			expectedErr: "a user with the UID not a uid does not have a supplier's account",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			supplier, err := s.FindSupplier(tt.args.ctx, tt.args.uid)

			if tt.wantErr && (err == nil) {
				t.Errorf("expected an error. got %v", err)
			}

			if tt.wantErr {
				assert.Nil(t, supplier)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}

			if !tt.wantErr {
				assert.NotNil(t, supplier)
				assert.Nil(t, err)
			}

		})
	}
}

func TestService_SuspendSupplier(t *testing.T) {
	service := NewService()
	ctx, token := createNewUser(context.Background(), t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, errS := service.AddSupplier(ctx, name, PartnerTypeProvider)
	if errS != nil {
		t.Errorf("can't create a supplier")
		return
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad case: suspend a nonexisting supplier",
			args: args{
				ctx: context.Background(),
				uid: "some random uid",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy case: suspend an existing supplier",
			args: args{
				ctx: ctx,
				uid: token.UID,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.SuspendSupplier(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SuspendSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.SuspendSupplier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SetUpSupplier(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	partnerName := gofakeit.Name()
	partnerType := PartnerTypePractitioner
	service.AddPartnerType(ctx, &partnerName, &partnerType)

	type args struct {
		ctx   context.Context
		input AccountType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful basic individual account set up",
			args: args{
				ctx:   ctx,
				input: "INDIVIDUAL",
			},
			wantErr: false,
		},
		{
			name: "Successful basic organisation account set up",
			args: args{
				ctx:   ctx,
				input: "ORGANISATION",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			_, err := s.SetUpSupplier(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SetUpSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_PublishKYCNudge(t *testing.T) {
	service := NewService()

	_, token := base.GetAuthenticatedContextAndToken(t)

	type args struct {
		uid     string
		partner PartnerType
		account AccountType
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedErr    string
		invalidService bool
	}{
		{
			name: "valid : Individual Rider KYC Nudge",
			args: args{
				uid:     token.UID,
				partner: PartnerTypeRider,
				account: AccountTypeIndividual,
			},
			wantErr:        false,
			invalidService: false,
		},
		{
			name: "valid : Organization Practitioner KYC Nudge",
			args: args{
				uid:     token.UID,
				partner: PartnerTypePractitioner,
				account: AccountTypeOrganisation,
			},
			wantErr:        false,
			invalidService: false,
		},

		{
			name: "invalid : unknown partner type",
			args: args{
				uid:     token.UID,
				partner: "alien partner",
				account: AccountTypeOrganisation,
			},
			wantErr:        true,
			invalidService: false,
			expectedErr:    "expected `partner` to be defined and to be valid",
		},
		{
			name: "invalid : consumer partner",
			args: args{
				uid:     token.UID,
				partner: PartnerTypeConsumer,
				account: AccountTypeOrganisation,
			},
			wantErr:        true,
			invalidService: false,
			expectedErr:    "invalid `partner`. cannot use CONSUMER in this context",
		},
		{
			name: "invalid : unknown account type",
			args: args{
				uid:     token.UID,
				partner: PartnerTypePractitioner,
				account: "alien account",
			},
			wantErr:        true,
			invalidService: false,
			expectedErr:    "provided `account` is not valid",
		},

		{
			name: "invalid : wrong engagement service",
			args: args{
				uid:     token.UID,
				partner: PartnerTypePractitioner,
				account: AccountTypeOrganisation,
			},
			wantErr:        true,
			invalidService: true,
			expectedErr:    "unable to publish kyc nudge. unexpected status code  404",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.invalidService {
				s := service
				err := s.PublishKYCNudge(tt.args.uid, &tt.args.partner, &tt.args.account)
				if !tt.wantErr {
					assert.Nil(t, err)
				}

				if tt.wantErr {
					assert.NotNil(t, err)
					assert.Contains(t, tt.expectedErr, err.Error())
				}
				return
			}

			is := service
			is.engagement = service.Otp
			err := is.PublishKYCNudge(tt.args.uid, &tt.args.partner, &tt.args.account)
			assert.NotNil(t, err)
			assert.Contains(t, tt.expectedErr, err.Error())

		})
	}
}

func TestService_AddIndividualRiderKyc(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = service.SetUpSupplier(ctx, AccountTypeIndividual)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	riderInput := IndividualRider{
		IdentificationDoc: Identification{
			IdentificationDocType:           IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678",
			IdentificationDocNumberUploadID: "12345678",
		},
		KRAPIN:                         "12345678",
		KRAPINUploadID:                 "12345678",
		DrivingLicenseUploadID:         "12345678",
		CertificateGoodConductUploadID: "12345678",
	}
	riderKYC := &IndividualRider{
		IdentificationDoc: Identification{
			IdentificationDocType:           IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "12345678",
			IdentificationDocNumberUploadID: "12345678",
		},
		KRAPIN:                         "12345678",
		KRAPINUploadID:                 "12345678",
		DrivingLicenseUploadID:         "12345678",
		CertificateGoodConductUploadID: "12345678",
	}

	type args struct {
		ctx   context.Context
		input IndividualRider
	}
	tests := []struct {
		name    string
		args    args
		want    *IndividualRider
		wantErr bool
	}{
		{
			name: "Successful Add individual rider KYC",
			args: args{
				ctx:   ctx,
				input: riderInput,
			},
			wantErr: false,
			want:    riderKYC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddIndividualRiderKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddIndividualRiderKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AddIndividualRiderKyc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_AddOrganizationRiderKyc(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx      context.Context
		input    OrganizationRider
		resource OrganizationRider
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: OrganizationRider{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
					OrganizationCertificate:            "ORG-CERT",
					KRAPIN:                             "KRA-PIN-12345678",
					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
				},
				resource: OrganizationRider{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
					OrganizationCertificate:            "ORG-CERT",
					KRAPIN:                             "KRA-PIN-12345678",
					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
				},
			},
			wantErr: false,
		},

		{
			name: "invalid : organization type name",
			args: args{
				ctx: ctx,
				input: OrganizationRider{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					CertificateOfIncorporation:         "CERT-OF-CORP-ID",
					CertificateOfInCorporationUploadID: "CERT-OF-CORP-UPLOAD-ID",
					OrganizationCertificate:            "ORG-CERT",
					KRAPIN:                             "KRA-PIN-12345678",
					KRAPINUploadID:                     "KRA-PIN-UPLOAD-ID12345678",
					SupportingDocumentsUploadID:        []string{"SUPPORTING-UPLOAD-ID12345678"},
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddOrganizationRiderKyc(tt.args.ctx, tt.args.input)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				if !reflect.DeepEqual(*got, tt.args.resource) {
					t.Errorf("Service.AddOrganizationRiderKyc() = %v, want %v", got, tt.args.resource)
				}
				return
			}
			assert.NotNil(t, err)
			assert.Nil(t, got)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestService_AddIndividualPractitionerKyc(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx      context.Context
		input    IndividualPractitioner
		resource IndividualPractitioner
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: IndividualPractitioner{
					IdentificationDoc: Identification{
						IdentificationDocType:           IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
				resource: IndividualPractitioner{
					IdentificationDoc: Identification{
						IdentificationDocType:           IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr: false,
		},

		{
			name: "invalid : practice services",
			args: args{
				ctx: ctx,
				input: IndividualPractitioner{
					IdentificationDoc: Identification{
						IdentificationDocType:           IdentificationDocTypeNationalid,
						IdentificationDocNumber:         "12345678",
						IdentificationDocNumberUploadID: "12345678",
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{"SUPPORTING"},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddIndividualPractitionerKyc(tt.args.ctx, tt.args.input)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				if !reflect.DeepEqual(*got, tt.args.resource) {
					t.Errorf("Service.AddIndividualPractitionerKyc() = %v, want %v", got, tt.args.resource)
				}
				return
			}
			assert.NotNil(t, err)
			assert.Nil(t, got)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestService_AddOrganizationPractitionerKyc(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx      context.Context
		input    OrganizationPractitioner
		resource OrganizationPractitioner
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: OrganizationPractitioner{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
				resource: OrganizationPractitioner{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr: false,
		},

		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: OrganizationPractitioner{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{"SUPPORTING"},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},

		{
			name: "invalid : practice services",
			args: args{
				ctx: ctx,
				input: OrganizationPractitioner{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{"SUPPORTING"},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddOrganizationPractitionerKyc(tt.args.ctx, tt.args.input)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				if !reflect.DeepEqual(*got, tt.args.resource) {
					t.Errorf("Service.AddOrganizationPractitionerKyc() = %v, want %v", got, tt.args.resource)
				}
				return
			}
			assert.NotNil(t, err)
			assert.Nil(t, got)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestService_AddOrganizationProviderKyc(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	name := gofakeit.Name()
	partnerRider := PartnerTypeRider
	_, err := service.AddPartnerType(ctx, &name, &partnerRider)
	if err != nil {
		t.Errorf("can't create a supplier")
		return
	}

	_, err = service.SetUpSupplier(ctx, AccountTypeOrganisation)
	if err != nil {
		t.Errorf("can't set up a supplier")
		return
	}

	type args struct {
		ctx      context.Context
		input    OrganizationProvider
		resource OrganizationProvider
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid : should pass",
			args: args{
				ctx: ctx,
				input: OrganizationProvider{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
				resource: OrganizationProvider{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServiceOther},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr: false,
		},

		{
			name: "invalid : organization type name ",
			args: args{
				ctx: ctx,
				input: OrganizationProvider{
					OrganizationTypeName: "AWESOME ORG",
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{"SUPPORTING"},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `OrganizationTypeName` provided : AWESOME ORG",
		},

		{
			name: "invalid : practice services",
			args: args{
				ctx: ctx,
				input: OrganizationProvider{
					OrganizationTypeName: OrganizationTypeLimitedCompany,
					DirectorIdentifications: []Identification{
						{
							IdentificationDocType:           IdentificationDocTypeNationalid,
							IdentificationDocNumber:         "12345678",
							IdentificationDocNumberUploadID: "12345678",
						},
					},
					KRAPIN:             "KRA-12345678",
					KRAPINUploadID:     "KRA-UPLOAD-12345678",
					RegistrationNumber: "REG-12345",
					PracticeLicenseID:  "PRAC-12345",
					PracticeServices:   []PractitionerService{"SUPPORTING"},
					Cadre:              PractitionerCadreDoctor,
				},
			},
			wantErr:     true,
			expectedErr: "invalid `PracticeService` provided : SUPPORTING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddOrganizationProviderKyc(tt.args.ctx, tt.args.input)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				if !reflect.DeepEqual(*got, tt.args.resource) {
					t.Errorf("Service.AddOrganizationProviderKyc() = %v, want %v", got, tt.args.resource)
				}
				return
			}
			assert.NotNil(t, err)
			assert.Nil(t, got)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestService_SendKYCEmail(t *testing.T) {
	type args struct {
		ctx          context.Context
		text         string
		emailaddress string
	}

	approvalEmail := generateProcessKYCApprovalEmailTemplate()
	rejectionEmail := generateProcessKYCRejectionEmailTemplate()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good case: approved",
			args: args{
				ctx:          context.Background(),
				text:         approvalEmail,
				emailaddress: base.GenerateRandomEmail(),
			},
			wantErr: false,
		},
		{
			name: "Good case: rejected",
			args: args{
				ctx:          context.Background(),
				text:         rejectionEmail,
				emailaddress: base.GenerateRandomEmail(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.SendKYCEmail(tt.args.ctx, tt.args.text, tt.args.emailaddress); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendKYCEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_PublishKYCFeedItem(t *testing.T) {
	service := NewService()
	ctx, token := base.GetAuthenticatedContextAndToken(t)

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
			name: "Send feed item to user",
			args: args{
				ctx:  ctx,
				uids: []string{token.UID},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			if err := s.PublishKYCFeedItem(tt.args.ctx, tt.args.uids...); (err != nil) != tt.wantErr {
				t.Errorf("Service.PublishKYCFeedItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
