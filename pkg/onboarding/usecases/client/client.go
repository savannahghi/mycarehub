package client

import "github.com/savannahghi/onboarding-service/pkg/onboarding/domain"

// IRegisterClient ...
type IRegisterClient interface {
	// TODO: the input client profile must not have an ID set
	//		validate identifiers when creating
	//		if the enrollemnt date is not supplied, set it automatically
	//		default to the client profile being active right after creation
	//		create a patient on FHIR (HealthRecordID
	//		if identifers not supplied (e.g patient being created on app), set
	//			an internal identifier as the default. It should be updated later
	//			with the CCC number or other final identifier
	// TODO: ensure the user exists...supplied user ID
	// TODO: only register clients who've been counselled
	// TODO: consider: after successful registration, send invite link automatically
	RegisterClient(user domain.User, profile domain.ClientProfileRegistrationPayload) (*domain.ClientProfile, error)
}

// IAddClientIdentifier ...
type IAddClientIdentifier interface {
	// TODO idType is an enum
	// TODO use idType and settings to decide if it's a primary identifier or not
	AddIdentifier(clientID string, idType string, idValue string, isPrimary bool) (*domain.Identifier, error)
}

// IInactivateClient ...
type IInactivateClient interface {
	// TODO Consider making reasons an enum
	InactivateClient(clientID string, reason string, notes string) (bool, error)
}

// IReactivateClient ...
type IReactivateClient interface {
	ReactivateClient(clientID string, reason string, notes string) (bool, error)
}

// ITransferClient ...
type ITransferClient interface {
	// TODO: maintain log of past transfers, who did it etc
	TransferClient(
		clientID string,
		OriginFacilityID string,
		DestinationFacilityID string,
		Reason string, // TODO: consider making this an enum
		Notes string, // optional notes...e.g if the reason given is "Other"
	) (bool, error)
}

// IGetClientIdentifiers ...
type IGetClientIdentifiers interface {
	GetIdentifiers(clientID string, active bool) ([]*domain.Identifier, error)
}

// IInactivateClientIdentifier ...
type IInactivateClientIdentifier interface {
	InactivateIdentifier(clientID string, identifierID string) (bool, error)
}

// IAssignTreatmentSupporter ...
type IAssignTreatmentSupporter interface {
	AssignTreatmentSupporter(
		clientID string,
		treatmentSupporterID string,
		treatmentSupporterType string, // TODO: enum, start with CHV and Treatment buddy
	) (bool, error)
}

// IUnassignTreatmentSupporter ...
type IUnassignTreatmentSupporter interface {
	UnassignTreatmentSupporter(
		clientID string,
		treatmentSupporterID string,
		reason string, // TODO: ensure these are in an audit log
		notes string, // TODO: Optional
	) (bool, error)
}

// IAddRelatedPerson ...
type IAddRelatedPerson interface {
	// add next of kin
	AddRelatedPerson(
		clientID string,
		relatedPerson *domain.RelatedPerson,
	) (*domain.RelatedPerson, bool)
}

// ClientProfileUseCases ...
type ClientProfileUseCases interface {
	IAddClientIdentifier
	IGetClientIdentifiers
	IInactivateClientIdentifier
	IRegisterClient
	IInactivateClient
	IReactivateClient
	ITransferClient
	IAssignTreatmentSupporter
	IUnassignTreatmentSupporter
	IAddRelatedPerson
}
