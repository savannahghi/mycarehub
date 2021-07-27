package crm

import (
	"context"

	hubspotDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

// ServiceCrm represents commontools crm lib usecases extension
type ServiceCrm interface {
	OptOut(ctx context.Context, phoneNumber string) (*hubspotDomain.CRMContact, error)
	CreateHubSpotContact(ctx context.Context, contact *hubspotDomain.CRMContact) (*hubspotDomain.CRMContact, error)
	UpdateHubSpotContact(ctx context.Context, contact *hubspotDomain.CRMContact) (*hubspotDomain.CRMContact, error)
	GetContactByPhone(ctx context.Context, phone string) (*hubspotDomain.CRMContact, error)
}

// Hubspot interacts with `HubSpot` CRM usecases
type Hubspot struct {
	hubSpotUsecases hubspotUsecases.HubSpotUsecases
}

// NewCrmService inits a new crm instance
func NewCrmService(hubSpotUsecases hubspotUsecases.HubSpotUsecases) *Hubspot {
	return &Hubspot{
		hubSpotUsecases: hubSpotUsecases,
	}
}

// OptOut marks a user as opted out of our marketing sms on both our firestore snd hubspot
func (h *Hubspot) OptOut(ctx context.Context, phoneNumber string) (*hubspotDomain.CRMContact, error) {
	return h.hubSpotUsecases.OptOut(ctx, phoneNumber)
}

// CreateHubSpotContact creates a hubspot contact on both our crm and firestore
func (h *Hubspot) CreateHubSpotContact(ctx context.Context, contact *hubspotDomain.CRMContact) (*hubspotDomain.CRMContact, error) {
	return h.hubSpotUsecases.CreateHubSpotContact(ctx, contact)
}

// UpdateHubSpotContact updates a hubspot contact on both our crm and firestore
func (h *Hubspot) UpdateHubSpotContact(ctx context.Context, contact *hubspotDomain.CRMContact) (*hubspotDomain.CRMContact, error) {
	return h.hubSpotUsecases.UpdateHubSpotContact(ctx, contact)
}

// GetContactByPhone gets a hubspot contact on both our crm and firestore
func (h *Hubspot) GetContactByPhone(ctx context.Context, phone string) (*hubspotDomain.CRMContact, error) {
	return h.hubSpotUsecases.GetContactByPhone(ctx, phone)
}
