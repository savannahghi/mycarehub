package erp

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// Service represents logic required to communicate with ERP
type Service interface {
	FetchERPClient() *base.ServerClient
	CreateERPSupplier(method string, path string, payload map[string]interface{}, partner domain.PartnerType) error
}

// ServiceImpl represents ERP usecases
type ServiceImpl struct {
	ERPClient *base.ServerClient
}

// NewERPService returns new instance of ServiceImpl
func NewERPService(r repository.OnboardingRepository) Service {

	erpClient, err := base.NewERPClient()
	if err != nil {
		log.Panicf("unable to initialize ERP client for profile service: %s", err)
		os.Exit(1)
	}
	if !erpClient.IsInitialized() {
		log.Panicf("uninitialized ERP client")
		os.Exit(1)
	}

	return &ServiceImpl{ERPClient: erpClient}
}

// CreateERPSupplier makes a call to create erp supplier
func (e *ServiceImpl) CreateERPSupplier(method string, path string, payload map[string]interface{}, partner domain.PartnerType) error {

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}

	if err := base.ReadRequestToTarget(e.ERPClient, "POST", path, "", content, &domain.Supplier{PartnerType: partner}); err != nil {
		return fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	return nil
}

// FetchERPClient retrieves the erp client
func (e *ServiceImpl) FetchERPClient() *base.ServerClient {
	return e.ERPClient
}
