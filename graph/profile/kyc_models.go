package profile

// Identification identify model
type Identification struct {
	IdentificationDocType           IdentificationDocType `json:"identificationDocType" mapstructure:"identificationDocType"`
	IdentificationDocNumber         string                `json:"identificationDocNumber" mapstructure:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                `json:"identificationDocNumberUploadID" mapstructure:"identificationDocNumberUploadID"`
}

// IdentificationInput ...
type IdentificationInput struct {
	IdentificationDocType           IdentificationDocType `json:"identificationDocType" mapstructure:"identificationDocType"`
	IdentificationDocNumber         string                `json:"identificationDocNumber" mapstructure:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                `json:"identificationDocNumberUploadID" mapstructure:"identificationDocNumberUploadID"`
}

// IndividualRider holds the KYC for an individual rider
type IndividualRider struct {
	IdentificationDoc              Identification `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                         string         `json:"KRAPIN" mapstructure:"kraPIN"`
	KRAPINUploadID                 string         `json:"KRAPINUploadID" mapstructure:"kraPINUploadID"`
	DrivingLicenseUploadID         string         `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string         `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// IndividualRiderInput is used to record the KYC for an individual rider
type IndividualRiderInput struct {
	IdentificationDoc              IdentificationInput `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                         string              `json:"KRAPIN" mapstructure:"kraPIN"`
	KRAPINUploadID                 string              `json:"KRAPINUploadID" mapstructure:"kraPINUploadID"`
	DrivingLicenseUploadID         string              `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string              `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []string            `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// OrganizationRider ...
type OrganizationRider struct {
	OrganizationTypeName               OrganizationType `json:"identificationDoc" mapstructure:"identificationDoc"`
	CertificateOfIncorporation         string           `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string           `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string           `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	KRAPIN                             string           `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string           `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string         `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// OrganizationRiderInput ...
type OrganizationRiderInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}
