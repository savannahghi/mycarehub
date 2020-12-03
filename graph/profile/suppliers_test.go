package profile

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_AddSupplier(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add supplier happy case",
			args: args{
				ctx:  ctx,
				name: "Be.Well Test Supplier",
			},
			wantErr: false,
		},
		{
			name: "add supplier sad case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			supplier, err := s.AddSupplier(tt.args.ctx, nil, tt.args.name)
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
	assert.NotNil(t, ctx)
	assert.NotNil(t, token)

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
			expectedErr: "unable to get Firebase user with UID not a uid: cannot find user from uid: \"not a uid\"",
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

func TestFindSupplierByUID(t *testing.T) {
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	service := NewService()
	findSupplier := FindSupplierByUIDHandler(ctx, service)

	uid := &BusinessPartnerUID{
		UID: &token.UID,
	}
	goodUIDJSONBytes, err := json.Marshal(uid)
	assert.Nil(t, err)
	assert.NotNil(t, goodUIDJSONBytes)
	goodSupplierRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	goodSupplierRequest.Body = ioutil.NopCloser(bytes.NewReader(goodUIDJSONBytes))

	emptyUID := &BusinessPartnerUID{}
	badUIDJSONBytes, err := json.Marshal(emptyUID)
	assert.Nil(t, err)
	assert.NotNil(t, badUIDJSONBytes)
	badSupplierRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	badSupplierRequest.Body = ioutil.NopCloser(bytes.NewReader(badUIDJSONBytes))

	badUID := "this uid does not exist"
	nonExistentUID := &BusinessPartnerUID{UID: &badUID}
	nonExistentUIDJSONBytes, err := json.Marshal(nonExistentUID)
	assert.Nil(t, err)
	assert.NotNil(t, nonExistentUIDJSONBytes)
	nonExistentSupplierRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	nonExistentSupplierRequest.Body = ioutil.NopCloser(bytes.NewReader(nonExistentUIDJSONBytes))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "valid : find Supplier",
			args: args{
				w: httptest.NewRecorder(),
				r: goodSupplierRequest,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid : missing user uid",
			args: args{
				w: httptest.NewRecorder(),
				r: badSupplierRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "not found customer request",
			args: args{
				w: httptest.NewRecorder(),
				r: nonExistentSupplierRequest,
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findSupplier(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, rec.Code, tt.wantStatusCode)
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
	_, err := service.AddSupplier(ctx, nil, "To Be Deleted")
	if err != nil {
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
