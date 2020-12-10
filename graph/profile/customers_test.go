package profile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestService_AddCustomer(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	assert.NotNil(t, ctx)

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
			name: "Happy case",
			args: args{
				ctx:  ctx,
				name: "Be.Well Test customer",
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
			customer, err := s.AddCustomer(tt.args.ctx, nil, tt.args.name)
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

func TestService_SuspendCustomer(t *testing.T) {
	service := NewService()
	ctx, token := createNewUser(context.Background(), t)
	_, err := service.AddCustomer(ctx, nil, "To Be Deleted")
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
			name: "sad case: suspend a nonexisting customer",
			args: args{
				ctx: context.Background(),
				uid: "some random uid",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy case: suspend an existing customer",
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
			got, err := s.SuspendCustomer(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SuspendCustomer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.SuspendCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}
