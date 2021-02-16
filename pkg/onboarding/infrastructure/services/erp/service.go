package erp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// ServiceERP represents logic required to communicate with ERP
type ServiceERP interface {
	FetchERPClient() *base.ServerClient
	CreateERPCustomer(
		method string,
		path string,
		payload map[string]interface{},
		customer base.Customer,
	) error
	CreateERPSupplier(
		method string,
		path string,
		payload map[string]interface{},
		supplier base.Supplier,
	) error
}

// ServiceERPImpl represents ERP usecases
type ServiceERPImpl struct {
	ERPClient *base.ServerClient
}

// NewERPService returns new instance of ServiceImpl
func NewERPService() ServiceERP {

	erpClient, err := base.NewERPClient()
	if err != nil {
		log.Panicf("unable to initialize ERP client for profile service: %s", err)
		os.Exit(1)
	}
	if !erpClient.IsInitialized() {
		log.Panicf("uninitialized ERP client")
		os.Exit(1)
	}

	return &ServiceERPImpl{ERPClient: erpClient}
}

// CreateERPCustomer makes a call to create erp supplier
func (e *ServiceERPImpl) CreateERPCustomer(
	method string,
	path string,
	payload map[string]interface{},
	customer base.Customer,
) error {

	content, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal to JSON: %v", err)
	}

	return base.ReadRequestToTarget(
		e.ERPClient,
		http.MethodPost,
		path,
		"",
		content,
		&customer,
	)
}

// CreateERPSupplier makes a call to create erp supplier
func (e *ServiceERPImpl) CreateERPSupplier(
	method string,
	path string,
	payload map[string]interface{},
	supplier base.Supplier,
) error {

	content, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal to JSON: %v", err)
	}

	return base.ReadRequestToTarget(
		e.ERPClient,
		http.MethodPost,
		path,
		"",
		content,
		&supplier,
	)
}

// FetchERPClient retrieves the erp client
func (e *ServiceERPImpl) FetchERPClient() *base.ServerClient {
	return e.ERPClient
}
