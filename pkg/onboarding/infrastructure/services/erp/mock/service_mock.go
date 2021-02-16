package mock

import (
	"gitlab.slade360emr.com/go/base"
)

// FakeServiceERP is an `ERP` service mock .
type FakeServiceERP struct {
	FetchERPClientFn    func() *base.ServerClient
	CreateERPCustomerFn func(
		method string,
		path string,
		payload map[string]interface{},
		customer base.Customer,
	) error
	CreateERPSupplierFn func(
		method string,
		path string,
		payload map[string]interface{},
		supplier base.Supplier,
	) error
}

// FetchERPClient ...
func (f *FakeServiceERP) FetchERPClient() *base.ServerClient {
	return f.FetchERPClientFn()
}

// CreateERPCustomer ...
func (f *FakeServiceERP) CreateERPCustomer(
	method string,
	path string,
	payload map[string]interface{},
	customer base.Customer,
) error {
	return f.CreateERPCustomerFn(method, path, payload, customer)
}

// CreateERPSupplier ...
func (f *FakeServiceERP) CreateERPSupplier(
	method string,
	path string,
	payload map[string]interface{},
	supplier base.Supplier,
) error {
	return f.CreateERPSupplierFn(method, path, payload, supplier)
}
