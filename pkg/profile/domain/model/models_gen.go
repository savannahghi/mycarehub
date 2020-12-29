// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/profile/domain"
)

type Identification struct {
	IdentificationDocType           domain.IdentificationDocType `json:"identificationDocType"`
	IdentificationDocNumber         string                       `json:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                       `json:"identificationDocNumberUploadID"`
}

type IdentificationInput struct {
	IdentificationDocType           domain.IdentificationDocType `json:"identificationDocType"`
	IdentificationDocNumber         string                       `json:"identificationDocNumber"`
	IdentificationDocNumberUploadID string                       `json:"identificationDocNumberUploadID"`
}

type IndividualCoach struct {
	IdentificationDoc           *Identification `json:"identificationDoc"`
	Krapin                      string          `json:"KRAPIN"`
	KRAPINUploadID              string          `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string       `json:"supportingDocumentsUploadID"`
	PracticeLicenseID           string          `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string         `json:"practiceLicenseUploadID"`
}

type IndividualCoachInput struct {
	IdentificationDoc           *IdentificationInput `json:"identificationDoc"`
	Krapin                      string               `json:"KRAPIN"`
	KRAPINUploadID              string               `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string            `json:"supportingDocumentsUploadID"`
	PracticeLicenseID           string               `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string              `json:"practiceLicenseUploadID"`
}

type IndividualNutrition struct {
	IdentificationDoc           *Identification `json:"identificationDoc"`
	Krapin                      string          `json:"KRAPIN"`
	KRAPINUploadID              string          `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string       `json:"supportingDocumentsUploadID"`
	PracticeLicenseID           string          `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string         `json:"practiceLicenseUploadID"`
}

type IndividualNutritionInput struct {
	IdentificationDoc           *IdentificationInput `json:"identificationDoc"`
	Krapin                      string               `json:"KRAPIN"`
	KRAPINUploadID              string               `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string            `json:"supportingDocumentsUploadID"`
	PracticeLicenseID           string               `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string              `json:"practiceLicenseUploadID"`
}

type IndividualPharmaceutical struct {
	IdentificationDoc           *Identification `json:"identificationDoc"`
	Krapin                      string          `json:"KRAPIN"`
	KRAPINUploadID              string          `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string       `json:"supportingDocumentsUploadID"`
	RegistrationNumber          string          `json:"registrationNumber"`
	PracticeLicenseID           string          `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string         `json:"practiceLicenseUploadID"`
}

type IndividualPharmaceuticalInput struct {
	IdentificationDoc           *IdentificationInput `json:"identificationDoc"`
	Krapin                      string               `json:"KRAPIN"`
	KRAPINUploadID              string               `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string            `json:"supportingDocumentsUploadID"`
	RegistrationNumber          string               `json:"registrationNumber"`
	PracticeLicenseID           string               `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string              `json:"practiceLicenseUploadID"`
}

type IndividualPractitioner struct {
	IdentificationDoc           *Identification              `json:"identificationDoc"`
	Krapin                      string                       `json:"KRAPIN"`
	KRAPINUploadID              string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string                    `json:"supportingDocumentsUploadID"`
	RegistrationNumber          string                       `json:"registrationNumber"`
	PracticeLicenseID           string                       `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string                      `json:"practiceLicenseUploadID"`
	PracticeServices            []domain.PractitionerService `json:"practiceServices"`
	Cadre                       domain.PractitionerCadre     `json:"cadre"`
}

type IndividualPractitionerInput struct {
	IdentificationDoc           *IdentificationInput         `json:"identificationDoc"`
	Krapin                      string                       `json:"KRAPIN"`
	KRAPINUploadID              string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID []*string                    `json:"supportingDocumentsUploadID"`
	RegistrationNumber          string                       `json:"registrationNumber"`
	PracticeLicenseID           string                       `json:"practiceLicenseID"`
	PracticeLicenseUploadID     *string                      `json:"practiceLicenseUploadID"`
	PracticeServices            []domain.PractitionerService `json:"practiceServices"`
	Cadre                       domain.PractitionerCadre     `json:"cadre"`
}

type IndividualRider struct {
	IdentificationDoc              *Identification `json:"identificationDoc"`
	Krapin                         string          `json:"KRAPIN"`
	KRAPINUploadID                 string          `json:"KRAPINUploadID"`
	DrivingLicenseID               string          `json:"drivingLicenseID"`
	DrivingLicenseUploadID         *string         `json:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string          `json:"certificateGoodConductUploadID"`
	SupportingDocumentsUploadID    []*string       `json:"supportingDocumentsUploadID"`
}

type IndividualRiderInput struct {
	IdentificationDoc              *IdentificationInput `json:"identificationDoc"`
	Krapin                         string               `json:"KRAPIN"`
	KRAPINUploadID                 string               `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID    []*string            `json:"supportingDocumentsUploadID"`
	DrivingLicenseID               string               `json:"drivingLicenseID"`
	DrivingLicenseUploadID         *string              `json:"drivingLicenseUploadID"`
	CertificateGoodConductUploadID string               `json:"certificateGoodConductUploadID"`
}

type OrganizationCoach struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification       `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationCoachInput struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput  `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationNutrition struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification       `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationNutritionInput struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput  `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationPharmaceutical struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification       `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationPharmaceuticalInput struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput  `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	RegistrationNumber                 string                  `json:"registrationNumber"`
	PracticeLicenseID                  string                  `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                 `json:"practiceLicenseUploadID"`
}

type OrganizationPractitioner struct {
	OrganizationTypeName               domain.OrganizationType      `json:"organizationTypeName"`
	Krapin                             string                       `json:"KRAPIN"`
	KRAPINUploadID                     string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string                    `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                      `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                      `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification            `json:"directorIdentifications"`
	OrganizationCertificate            *string                      `json:"organizationCertificate"`
	RegistrationNumber                 string                       `json:"registrationNumber"`
	PracticeLicenseUploadID            string                       `json:"practiceLicenseUploadID"`
	PracticeServices                   []domain.PractitionerService `json:"practiceServices"`
	Cadre                              domain.PractitionerCadre     `json:"cadre"`
}

type OrganizationPractitionerInput struct {
	OrganizationTypeName               domain.OrganizationType      `json:"organizationTypeName"`
	Krapin                             string                       `json:"KRAPIN"`
	KRAPINUploadID                     string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string                    `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         string                       `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID string                       `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput       `json:"directorIdentifications"`
	OrganizationCertificate            *string                      `json:"organizationCertificate"`
	RegistrationNumber                 string                       `json:"registrationNumber"`
	PracticeLicenseID                  string                       `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                      `json:"practiceLicenseUploadID"`
	PracticeServices                   []domain.PractitionerService `json:"practiceServices"`
	Cadre                              domain.PractitionerCadre     `json:"cadre"`
}

type OrganizationProvider struct {
	OrganizationTypeName               domain.OrganizationType      `json:"organizationTypeName"`
	Krapin                             string                       `json:"KRAPIN"`
	KRAPINUploadID                     string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string                    `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                      `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                      `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification            `json:"directorIdentifications"`
	OrganizationCertificate            *string                      `json:"organizationCertificate"`
	RegistrationNumber                 string                       `json:"registrationNumber"`
	PracticeLicenseID                  string                       `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                      `json:"practiceLicenseUploadID"`
	PracticeServices                   []domain.PractitionerService `json:"practiceServices"`
	Cadre                              domain.PractitionerCadre     `json:"cadre"`
}

type OrganizationProviderInput struct {
	OrganizationTypeName               domain.OrganizationType      `json:"organizationTypeName"`
	Krapin                             string                       `json:"KRAPIN"`
	KRAPINUploadID                     string                       `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string                    `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                      `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                      `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput       `json:"directorIdentifications"`
	OrganizationCertificate            *string                      `json:"organizationCertificate"`
	RegistrationNumber                 string                       `json:"registrationNumber"`
	PracticeLicenseID                  string                       `json:"practiceLicenseID"`
	PracticeLicenseUploadID            *string                      `json:"practiceLicenseUploadID"`
	PracticeServices                   []domain.PractitionerService `json:"practiceServices"`
	Cadre                              domain.PractitionerCadre     `json:"cadre"`
}

type OrganizationRider struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*Identification       `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
}

type OrganizationRiderInput struct {
	OrganizationTypeName               domain.OrganizationType `json:"organizationTypeName"`
	Krapin                             string                  `json:"KRAPIN"`
	KRAPINUploadID                     string                  `json:"KRAPINUploadID"`
	SupportingDocumentsUploadID        []*string               `json:"supportingDocumentsUploadID"`
	CertificateOfIncorporation         *string                 `json:"certificateOfIncorporation"`
	CertificateOfInCorporationUploadID *string                 `json:"certificateOfInCorporationUploadID"`
	DirectorIdentifications            []*IdentificationInput  `json:"directorIdentifications"`
	OrganizationCertificate            *string                 `json:"organizationCertificate"`
}

type Practitioner struct {
	Profile                  *base.UserProfile          `json:"profile"`
	License                  string                     `json:"license"`
	Cadre                    domain.PractitionerCadre   `json:"cadre"`
	Specialty                base.PractitionerSpecialty `json:"specialty"`
	ProfessionalProfile      base.Markdown              `json:"professionalProfile"`
	AverageConsultationPrice float64                    `json:"averageConsultationPrice"`
	Services                 *domain.ServicesOffered    `json:"services"`
}

type PractitionerConnection struct {
	Edges    []*PractitionerEdge `json:"edges"`
	PageInfo *base.PageInfo      `json:"pageInfo"`
}

type PractitionerEdge struct {
	Cursor *string       `json:"cursor"`
	Node   *Practitioner `json:"node"`
}

type TesterWhitelist struct {
	Email string `json:"email"`
}
