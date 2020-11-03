package profile

import (
	"context"
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
			customer, err := s.AddCustomer(tt.args.ctx)
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
