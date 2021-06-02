package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

// FakeServiceERP is an `ERP` service mock .
type FakeServiceERP struct {
	FetchERPClientFn    func() *base.ServerClient
	CreateERPCustomerFn func(
		ctx context.Context,
		customerPayload dto.CustomerPayload,
		UID string,
	) (*base.Customer, error)
	CreateERPSupplierFn func(
		ctx context.Context,
		supplierPayload dto.SupplierPayload,
		UID string,
	) (*base.Supplier, error)
}

// FetchERPClient ...
func (f *FakeServiceERP) FetchERPClient() *base.ServerClient {
	return f.FetchERPClientFn()
}

// CreateERPCustomer ...
func (f *FakeServiceERP) CreateERPCustomer(
	ctx context.Context,
	customerPayload dto.CustomerPayload,
	UID string,
) (*base.Customer, error) {
	return f.CreateERPCustomerFn(ctx, customerPayload, UID)
}

// CreateERPSupplier ...
func (f *FakeServiceERP) CreateERPSupplier(
	ctx context.Context,
	supplierPayload dto.SupplierPayload,
	UID string,
) (*base.Supplier, error) {
	return f.CreateERPSupplierFn(ctx, supplierPayload, UID)
}
