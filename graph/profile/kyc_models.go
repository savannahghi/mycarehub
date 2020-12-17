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
	KRAPIN                         string         `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                 string         `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	DrivingLicenseUploadID         string         `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string         `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []string       `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// IndividualRiderInput is used to record the KYC for an individual rider
type IndividualRiderInput struct {
	IdentificationDoc              IdentificationInput `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                         string              `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                 string              `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	DrivingLicenseUploadID         string              `json:"drivingLicenseUploadID" mapstructure:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string              `json:"certificateGoodConductUploadID" mapstructure:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []string            `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
}

// IndividualPractitionerInput ...
type IndividualPractitionerInput struct {
	IdentificationDoc           IdentificationInput   `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	RegistrationNumber          string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID           string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices            []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
	Cadre                       PractitionerCadre     `json:"cadre" mapstructure:"cadre"`
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

// IndividualPharmaceuticalInput ...
type IndividualPharmaceuticalInput struct {
	IdentificationDoc           IdentificationInput `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string              `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string              `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string            `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	RegistrationNumber          string              `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID           string              `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string              `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// IndividualCoachInput ...
type IndividualCoachInput struct {
	IdentificationDoc           IdentificationInput `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string              `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string              `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string            `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	PracticeLicenseID           string              `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string              `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// IndividualNutritionInput ...
type IndividualNutritionInput struct {
	IdentificationDoc           IdentificationInput `json:"identificationDoc" mapstructure:"identificationDoc"`
	KRAPIN                      string              `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID              string              `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID []string            `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	PracticeLicenseID           string              `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID     string              `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// OrganizationPractitionerInput ...
type OrganizationPractitionerInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
	PracticeServices                   []PractitionerService `json:"practiceServices" mapstructure:"practiceServices"`
	Cadre                              PractitionerCadre     `json:"cadre" mapstructure:"cadre"`
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

// OrganizationProviderInput ...
type OrganizationProviderInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
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

// OrganizationNutritionInput ...
type OrganizationNutritionInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// OrganizationCoachInput ...
type OrganizationCoachInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
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

// OrganizationPharmaceuticalInput ...
type OrganizationPharmaceuticalInput struct {
	OrganizationTypeName               OrganizationType      `json:"organizationTypeName" mapstructure:"organizationTypeName"`
	KRAPIN                             string                `json:"KRAPIN" mapstructure:"KRAPIN"`
	KRAPINUploadID                     string                `json:"KRAPINUploadID" mapstructure:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []string              `json:"supportingDocumentsUploadID" mapstructure:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                `json:"certificateOfIncorporation" mapstructure:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                `json:"certificateOfInCorporationUploadID" mapstructure:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []IdentificationInput `json:"directorIdentifications" mapstructure:"directorIdentifications"`
	OrganizationCertificate            string                `json:"organizationCertificate" mapstructure:"organizationCertificate"`
	RegistrationNumber                 string                `json:"registrationNumber" mapstructure:"registrationNumber"`
	PracticeLicenseID                  string                `json:"practiceLicenseID" mapstructure:"practiceLicenseID"`
	PracticeLicenseUploadID            string                `json:"practiceLicenseUploadID" mapstructure:"practiceLicenseUploadID"`
}
