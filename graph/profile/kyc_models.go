package profile

// Identification identify model
type Identification struct {
	IdentificationDocType           IdentificationDocType `json:"identificationDocType" mapstructure:"identificationDocType"`
	IdentificationDocNumber         string                `json:"identificationDocNumber" mapstructure:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                `json:"identificationDocNumberUploadID" mapstructure:"identificationDocNumberUploadID"`
}

// IndividualRider holds the KYC for an individual rider
type IndividualRider struct {
	IdentificationDoc              Identification `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                         string         `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                 string         `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	DrivingLicenseID               string         `json:"drivingLicenseID" mapstructure:"drivingLicenseID"`
	DrivingLicenseUploadID         string         `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string         `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// IndividualPractitioner ...
type IndividualPractitioner struct {
	IdentificationDoc           Identification        `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	RegistrationNumber          string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID           string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices            []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
	Cadre                       PractitionerCadre     `json:"cadre" mapstructure:"cadre"`
}

// IndividualPharmaceutical ...
type IndividualPharmaceutical struct {
	IdentificationDoc           Identification `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string         `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string         `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	RegistrationNumber          string         `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID           string         `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string         `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// IndividualCoach ...
type IndividualCoach struct {
	IdentificationDoc           Identification `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string         `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string         `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	PracticeLicenseID           string         `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string         `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// IndividualNutrition ...
type IndividualNutrition struct {
	IdentificationDoc           Identification `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string         `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string         `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	PracticeLicenseID           string         `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string         `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// OrganizationPractitioner ...
type OrganizationPractitioner struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
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

// OrganizationProvider ...
type OrganizationProvider struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
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

// OrganizationNutrition ...
type OrganizationNutrition struct {
	OrganizationTypeName               OrganizationType `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string           `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string           `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string         `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string           `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string           `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string           `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string           `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string           `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string           `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// OrganizationCoach ...
type OrganizationCoach struct {
	OrganizationTypeName               OrganizationType `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string           `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string           `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string         `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string           `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string           `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string           `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string           `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string           `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string           `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}

// OrganizationPharmaceutical ...
type OrganizationPharmaceutical struct {
	OrganizationTypeName               OrganizationType `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string           `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string           `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string         `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string           `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string           `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []Identification `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string           `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string           `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string           `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string           `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}
