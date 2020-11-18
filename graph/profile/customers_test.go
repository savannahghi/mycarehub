package profile

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_AddCustomer(t *testing.T) {
	service := NewService()
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)

	fireBaseClient, clientErr := base.GetFirebaseAuthClient(ctx)
	assert.Nil(t, clientErr)
	assert.NotNil(t, fireBaseClient)

	user, userErr := fireBaseClient.GetUser(ctx, authToken.UID)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	params := (&auth.UserToUpdate{}).
		DisplayName("Be.Well Test User")
	u, err := fireBaseClient.UpdateUser(ctx, authToken.UID, params)
	assert.Nil(t, err)
	assert.NotNil(t, u)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			customer, err := s.AddCustomer(tt.args.ctx, nil)
			if err != nil {
				assert.Nil(t, customer)
			}
			if err == nil {
				assert.NotNil(t, customer)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_SaveCustomerToFireStore(t *testing.T) {
	service := NewService()
	type args struct {
		customer Customer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save customer happy case",
			args: args{
				customer: Customer{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			if err := s.SaveCustomerToFireStore(tt.args.customer); (err != nil) != tt.wantErr {
				t.Errorf("Service.SaveCustomerToFireStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_AddCustomerKYC(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	type args struct {
		ctx   context.Context
		input CustomerKYCInput
	}
	tests := []struct {
		name    string
		args    args
		want    *CustomerKYC
		wantErr bool
	}{
		{
			name: "add customer kyc happy case",
			args: args{
				ctx: ctx,
				input: CustomerKYCInput{
					KRAPin:     "a valid kra pin",
					Occupation: "hustler",
					IDNumber:   "totally an id number",
					Address:    "we still use this",
					City:       "Nairobi",
				},
			},
			wantErr: false,
		},
		{
			name: "add customer kyc sad case",
			args: args{
				ctx:   context.Background(),
				input: CustomerKYCInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			customerKYC, err := s.AddCustomerKYC(tt.args.ctx, tt.args.input)
			if err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, customerKYC)
			}
			if err != nil {
				assert.Nil(t, customerKYC)
				assert.NotNil(t, err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddCustomerKYC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_UpdateCustomer(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	type args struct {
		ctx   context.Context
		input CustomerKYCInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update customer happy case",
			args: args{
				ctx: ctx,
				input: CustomerKYCInput{
					Occupation: "changed to employee",
				},
			},
			wantErr: false,
		},
		{
			name: "update customer sad case",
			args: args{
				ctx:   context.Background(),
				input: CustomerKYCInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			customer, err := s.UpdateCustomer(tt.args.ctx, tt.args.input)
			if err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, customer)
			}
			if err != nil {
				assert.Nil(t, customer)
				assert.NotNil(t, err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_FindCustomer(t *testing.T) {
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
			name: "valid : authenicated context",
			args: args{
				ctx: ctx,
				uid: token.UID,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				uid: "not a uid",
			},
			wantErr:     true,
			expectedErr: "unable to get Firebase user with UID not a uid: cannot find user from uid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			customer, err := s.FindCustomer(tt.args.ctx, tt.args.uid)
			if tt.wantErr && (err == nil) {
				t.Errorf("expected an error. got %v", err)
			}

			if tt.wantErr {
				assert.Nil(t, customer)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}

			if !tt.wantErr {
				assert.NotNil(t, customer)
				assert.Nil(t, err)
			}

		})
	}
}

func TestFindCustomerByUID(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)
	service := NewService()
	profile, err := service.UserProfile(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	findCustomer := FindCustomerByUIDHandler(ctx, service)

	uid := &BusinessPartnerUID{UID: &profile.UID}
	goodUIDJSONBytes, err := json.Marshal(uid)
	assert.Nil(t, err)
	assert.NotNil(t, goodUIDJSONBytes)
	goodCustomerRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	goodCustomerRequest.Body = ioutil.NopCloser(bytes.NewReader(goodUIDJSONBytes))

	emptyUID := &BusinessPartnerUID{}
	badUIDJSONBytes, err := json.Marshal(emptyUID)
	assert.Nil(t, err)
	assert.NotNil(t, badUIDJSONBytes)
	badCustomerRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	badCustomerRequest.Body = ioutil.NopCloser(bytes.NewReader(badUIDJSONBytes))

	badUID := "this uid does not exist"
	nonExistentUID := &BusinessPartnerUID{UID: &badUID}
	nonExistentUIDJSONBytes, err := json.Marshal(nonExistentUID)
	assert.Nil(t, err)
	assert.NotNil(t, nonExistentUIDJSONBytes)
	nonExistentCustomerRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	nonExistentCustomerRequest.Body = ioutil.NopCloser(bytes.NewReader(nonExistentUIDJSONBytes))

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
			name: "valid : find customer",
			args: args{
				w: httptest.NewRecorder(),
				r: goodCustomerRequest,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "invalid : missing user uid",
			args: args{
				w: httptest.NewRecorder(),
				r: badCustomerRequest,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "not found customer request",
			args: args{
				w: httptest.NewRecorder(),
				r: nonExistentCustomerRequest,
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findCustomer(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}
