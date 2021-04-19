package domain

// SupportingDocument used to add more documents when
type SupportingDocument struct {
	SupportingDocumentTitle       string `json:"supportingDocumentTitle" mapstructure:"supportingDocumentTitle"`
	SupportingDocumentDescription string `json:"supportingDocumentDescription" mapstructure:"supportingDocumentDescription"`
	SupportingDocumentUpload      string `json:"supportingDocumentUpload" mapstructure:"supportingDocumentUpload"`
}

// Identification identify model
type Identification struct {
	IdentificationDocType           IdentificationDocType `json:"identificationDocType" mapstructure:"identificationDocType"`
	IdentificationDocNumber         string                `json:"identificationDocNumber" mapstructure:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                `json:"identificationDocNumberUploadID" mapstructure:"identificationDocNumberUploadID"`
}

// IndividualRider represents the KYC information required for an Individual Rider
type IndividualRider struct {
	IdentificationDoc              Identification       `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                         string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                 string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	DrivingLicenseID               string               `json:"drivingLicenseID" mapstructure:"drivingLicenseID"`
	DrivingLicenseUploadID         string               `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string               `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocuments            []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
}

// IndividualPractitioner represents the KYC information required for an Individual Rider
type IndividualPractitioner struct {
	IdentificationDoc       Identification        `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                  string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID          string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments     []SupportingDocument  `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	RegistrationNumber      string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID       string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices        []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
	Cadre                   PractitionerCadre     `json:"cadre" mapstructure:"cadre"`
}

// IndividualPharmaceutical represents the KYC information required for an Individual Pharmaceutical
type IndividualPharmaceutical struct {
	IdentificationDoc       Identification       `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                  string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID          string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments     []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	RegistrationNumber      string               `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID       string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// IndividualCoach represents the KYC information required for an Individual Coach
type IndividualCoach struct {
	IdentificationDoc       Identification       `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                  string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID          string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments     []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	PracticeLicenseID       string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	AccreditationID         string               `json:"accreditationID" mapstructure:"accreditationID"`
	AccreditationUploadID   string               `json:"accreditationUploadID" mapstructure:"accreditationUploadID"`
}

// IndividualNutrition represents the KYC information required for an Individual Nutrition
type IndividualNutrition struct {
	IdentificationDoc       Identification       `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                  string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID          string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments     []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	PracticeLicenseID       string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// OrganizationRider represents the KYC information required for an Organization Rider
type OrganizationRider struct {
	OrganizationTypeName               OrganizationType     `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	CertificateOfIncorporation         string               `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string               `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification     `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string               `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	KRAPIN                             string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
}

// OrganizationPractitioner represents the KYC information required for an Organization Practitioner
type OrganizationPractitioner struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument  `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification      `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices                   []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
	Cadre                              PractitionerCadre     `json:"cadre" mapstructure:"cadre"`
}

// OrganizationProvider represents the KYC information required for an Organization Provider
type OrganizationProvider struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument  `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification      `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices                   []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
}

// OrganizationNutrition represents the KYC information required for an Organization Nutrition
type OrganizationNutrition struct {
	OrganizationTypeName               OrganizationType     `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	CertificateOfIncorporation         string               `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string               `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification     `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string               `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string               `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// OrganizationCoach represents the KYC information required for an Organization Coach
type OrganizationCoach struct {
	OrganizationTypeName               OrganizationType     `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	CertificateOfIncorporation         string               `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string               `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification     `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string               `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string               `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// OrganizationPharmaceutical represents the KYC information required for an Organization Pharmaceutical
type OrganizationPharmaceutical struct {
	OrganizationTypeName               OrganizationType     `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string               `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string               `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocuments                []SupportingDocument `json:"supportingDocuments" mapstructure:"supportingDocuments"`
	CertificateOfIncorporation         string               `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string               `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification     `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string               `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string               `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string               `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string               `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}
