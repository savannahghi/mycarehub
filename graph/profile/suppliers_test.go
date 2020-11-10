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
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add supplier happy case",
			args: args{
				ctx: ctx,
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
			supplier, err := s.AddSupplier(tt.args.ctx)
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
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
				uid: token.UID,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx: context.Background(),
				uid: "not a uid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			supplier, err := s.FindSupplier(tt.args.ctx, tt.args.uid)
			if supplier == nil {
				assert.Nil(t, err)
				assert.Nil(t, supplier)
			}
			if supplier != nil {
				assert.Nil(t, err)
				assert.NotNil(t, supplier)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFindSupplierByUID(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	service := NewService()
	profile, err := service.UserProfile(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	findSupplier := FindSupplierByUIDHandler(ctx, service)

	uid := &BusinessPartnerUID{
		UID: profile.UID,
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

	nonExistentUID := &BusinessPartnerUID{UID: "this uid does not exist"}
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
			name: "happy find Supplier request",
			args: args{
				w: httptest.NewRecorder(),
				r: goodSupplierRequest,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "sad find supplier request",
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
