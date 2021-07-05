package extension

import (
	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
)

// CRMExtension define crm extension's interface
type CRMExtension interface {
	CreateContact(contact domain.CRMContact) (*domain.CRMContact, error)
	UpdateContact(
		phone string,
		properties domain.ContactProperties,
	) (*domain.CRMContact, error)
}

// CRMExtensionImpl repreents crm expension implementation
type CRMExtensionImpl struct {
	Hubspot hubspot.ServiceHubSpotInterface
}

// NewCRMExtension initializes a new hubspot crm extension
func NewCRMExtension(hubspot hubspot.ServiceHubSpotInterface) CRMExtension {
	return &CRMExtensionImpl{Hubspot: hubspot}
}

// CreateContact extends commontool library's create crm contact logic
func (c *CRMExtensionImpl) CreateContact(
	contact domain.CRMContact,
) (*domain.CRMContact, error) {
	return c.Hubspot.CreateContact(contact)
}

// UpdateContact extends commontool library's update crm contact logic
func (c *CRMExtensionImpl) UpdateContact(
	phone string,
	properties domain.ContactProperties,
) (*domain.CRMContact, error) {
	return c.Hubspot.UpdateContact(phone, properties)
}
