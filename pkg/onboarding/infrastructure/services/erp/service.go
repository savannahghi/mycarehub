package erp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	supplierAPIPath = "/api/business_partners/suppliers/"
	customerAPIPath = "/api/business_partners/customers/"
)

// ServiceERP represents logic required to communicate with ERP
type ServiceERP interface {
	FetchERPClient() *base.ServerClient
	CreateERPCustomer(
		ctx context.Context,
		customerPayload dto.CustomerPayload,
		UID string,
	) (*base.Customer, error)
	CreateERPSupplier(
		ctx context.Context,
		supplierPayload dto.SupplierPayload,
		UID string,
	) (*base.Supplier, error)
}

// ServiceERPImpl represents ERP usecases
type ServiceERPImpl struct {
	ERPClient *base.ServerClient
	repo      repository.OnboardingRepository
}

// NewERPService returns new instance of ServiceImpl
func NewERPService(r repository.OnboardingRepository) ServiceERP {

	erpClient, err := base.NewERPClient()
	if err != nil {
		log.Panicf("unable to initialize ERP client for profile service: %s", err)
		os.Exit(1)
	}
	if !erpClient.IsInitialized() {
		log.Panicf("uninitialized ERP client")
		os.Exit(1)
	}

	return &ServiceERPImpl{
		ERPClient: erpClient,
		repo:      r,
	}
}

// CreateERPCustomer makes a call to create erp supplier
func (e *ServiceERPImpl) CreateERPCustomer(
	ctx context.Context,
	customerPayload dto.CustomerPayload,
	UID string,
) (*base.Customer, error) {
	profile, err := e.repo.GetUserProfileByUID(
		ctx,
		UID,
		false,
	)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(customerPayload)
	if err != nil {
		return nil, err
	}

	data, err := base.ReadWriteRequestToTarget(
		e.ERPClient,
		http.MethodPost,
		customerAPIPath,
		"",
		payload,
		&base.Customer{},
	)
	if err != nil {
		return nil, err
	}

	customer := data.(*base.Customer)
	customer.Active = true

	return e.repo.UpdateCustomerProfile(
		ctx,
		profile.ID,
		*customer,
	)
}

// CreateERPSupplier makes a call to create erp supplier
func (e *ServiceERPImpl) CreateERPSupplier(
	ctx context.Context,
	supplierPayload dto.SupplierPayload,
	UID string,
) (*base.Supplier, error) {
	profile, err := e.repo.GetUserProfileByUID(ctx, UID, false)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(supplierPayload)
	if err != nil {
		return nil, err
	}

	data, err := base.ReadWriteRequestToTarget(
		e.ERPClient,
		http.MethodPost,
		supplierAPIPath,
		"",
		payload,
		&base.Supplier{},
	)
	if err != nil {
		return nil, err
	}

	supplier := data.(*base.Supplier)
	supplier.Active = true

	return e.repo.ActivateSupplierProfile(
		profile.ID,
		*supplier,
	)
}

// FetchERPClient retrieves the erp client
func (e *ServiceERPImpl) FetchERPClient() *base.ServerClient {
	return e.ERPClient
}
