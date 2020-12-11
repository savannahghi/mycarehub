package profile

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_AddSupplier(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)

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
			supplier, err := s.AddSupplier(tt.args.ctx, nil, tt.args.name, tt.args.partnerType)
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
	supplier, err := service.AddSupplier(ctx, &token.UID, gofakeit.Name(), PartnerTypeProvider)
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

func TestService_AddSupplierKyc(t *testing.T) {
	srv := NewService()
	idDocType := IdentificationDocTypeNationalid
	idNo := "234567"
	license := "UL15KIAWAP!"
	cadre := PractitionerCadreDoctor
	profession := "Wa Mifupa"
	orgAccType := AccountTypeOrganisation
	type args struct {
		ctx   context.Context
		input SupplierKYCInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case: add account type",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				input: SupplierKYCInput{
					AccountType: orgAccType,
				},
			},
			wantErr: false,
		},
		{
			name: "good case: add identification docs",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				input: SupplierKYCInput{
					IdentificationDocType:   &idDocType,
					IdentificationDocNumber: &idNo,
				},
			},
			wantErr: false,
		},
		{
			name: "good case: add practice details",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				input: SupplierKYCInput{
					License:    &license,
					Cadre:      &cadre,
					Profession: &profession,
				},
			},
			wantErr: false,
		},
		{
			name: "bad case: nonexistent supplier",
			args: args{
				ctx: context.Background(),
				input: SupplierKYCInput{
					AccountType: AccountTypeIndividual,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			_, err := s.AddSupplierKyc(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddSupplierKyc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_SuspendSupplier(t *testing.T) {
	service := NewService()
	ctx, token := createNewUser(context.Background(), t)
	_, err := service.AddSupplier(ctx, nil, "To Be Deleted", PartnerTypeProvider)
	if err != nil {
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

	newIndividualSupplierInput := SupplierAccountInput{
		AccountType:       AccountTypeIndividual,
		UnderOrganization: false,
	}

	newOrganisationSupplierInput := SupplierAccountInput{
		AccountType:       AccountTypeOrganisation,
		UnderOrganization: false,
	}

	type args struct {
		ctx   context.Context
		input SupplierAccountInput
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
				input: newIndividualSupplierInput,
			},
			wantErr: false,
		},
		{
			name: "Successful basic organisation account set up",
			args: args{
				ctx:   ctx,
				input: newOrganisationSupplierInput,
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
