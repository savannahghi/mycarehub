package domain

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// PractitionerCadre is a list of health worker cadres.
type PractitionerCadre string

// practitioner cadre constants
const (
	PractitionerCadreDoctor          PractitionerCadre = "DOCTOR"
	PractitionerCadreClinicalOfficer PractitionerCadre = "CLINICAL_OFFICER"
	PractitionerCadreNurse           PractitionerCadre = "NURSE"
)

// AllPractitionerCadre is the set of known valid practitioner cadres
var AllPractitionerCadre = []PractitionerCadre{
	PractitionerCadreDoctor,
	PractitionerCadreClinicalOfficer,
	PractitionerCadreNurse,
}

// IsValid returns true if a practitioner cadre is valid
func (e PractitionerCadre) IsValid() bool {
	switch e {
	case PractitionerCadreDoctor, PractitionerCadreClinicalOfficer, PractitionerCadreNurse:
		return true
	}
	return false
}

func (e PractitionerCadre) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a practitioner cadre
func (e *PractitionerCadre) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PractitionerCadre(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PractitionerCadre", str)
	}
	return nil
}

// MarshalGQL writes the practitioner cadre to the supplied writer
func (e PractitionerCadre) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// FivePointRating is used to implement
type FivePointRating string

// known ratings
const (
	FivePointRatingPoor           FivePointRating = "POOR"
	FivePointRatingUnsatisfactory FivePointRating = "UNSATISFACTORY"
	FivePointRatingAverage        FivePointRating = "AVERAGE"
	FivePointRatingSatisfactory   FivePointRating = "SATISFACTORY"
	FivePointRatingExcellent      FivePointRating = "EXCELLENT"
)

// AllFivePointRating is a list of all known ratings
var AllFivePointRating = []FivePointRating{
	FivePointRatingPoor,
	FivePointRatingUnsatisfactory,
	FivePointRatingAverage,
	FivePointRatingSatisfactory,
	FivePointRatingExcellent,
}

// IsValid returns true for valid ratings
func (e FivePointRating) IsValid() bool {
	switch e {
	case FivePointRatingPoor, FivePointRatingUnsatisfactory, FivePointRatingAverage, FivePointRatingSatisfactory, FivePointRatingExcellent:
		return true
	}
	return false
}

func (e FivePointRating) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a rating value
func (e *FivePointRating) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FivePointRating(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FivePointRating", str)
	}
	return nil
}

// MarshalGQL converts the rating into a valid JSON string
func (e FivePointRating) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// PractitionerService defines the various services practitioners offer
type PractitionerService string

// PractitionerServiceOutpatientServices is a constant of all known practitioner service
const (
	PractitionerServiceOutpatientServices PractitionerService = "OUTPATIENT_SERVICES"
	PractitionerServiceInpatientServices  PractitionerService = "INPATIENT_SERVICES"
	PractitionerServicePharmacy           PractitionerService = "PHARMACY"
	PractitionerServiceMaternity          PractitionerService = "MATERNITY"
	PractitionerServiceLabServices        PractitionerService = "LAB_SERVICES"
	PractitionerServiceOther              PractitionerService = "OTHER"
)

//AllPractitionerService is a list of all known practitioner service
var AllPractitionerService = []PractitionerService{
	PractitionerServiceOutpatientServices,
	PractitionerServiceInpatientServices,
	PractitionerServicePharmacy,
	PractitionerServiceMaternity,
	PractitionerServiceLabServices,
	PractitionerServiceOther,
}

// IsValid returns true for valid practitioner service
func (e PractitionerService) IsValid() bool {
	switch e {
	case PractitionerServiceOutpatientServices, PractitionerServiceInpatientServices, PractitionerServicePharmacy, PractitionerServiceMaternity, PractitionerServiceLabServices, PractitionerServiceOther:
		return true
	}
	return false
}

func (e PractitionerService) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a practitioner service value
func (e *PractitionerService) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PractitionerService(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PractitionerService", str)
	}
	return nil
}

// MarshalGQL converts the practitioner service into a valid JSON string
func (e PractitionerService) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// BeneficiaryRelationship defines the various relationships with beneficiaries
type BeneficiaryRelationship string

// BeneficiaryRelationshipSpouse is a constant of beneficiary spouse relationship
const (
	BeneficiaryRelationshipSpouse BeneficiaryRelationship = "SPOUSE"
	BeneficiaryRelationshipChild  BeneficiaryRelationship = "CHILD"
)

//AllBeneficiaryRelationship is a list of all known beneficiary relationships
var AllBeneficiaryRelationship = []BeneficiaryRelationship{
	BeneficiaryRelationshipSpouse,
	BeneficiaryRelationshipChild,
}

// IsValid returns true for valid beneficiary relationship
func (e BeneficiaryRelationship) IsValid() bool {
	switch e {
	case BeneficiaryRelationshipSpouse, BeneficiaryRelationshipChild:
		return true
	}
	return false
}

func (e BeneficiaryRelationship) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a beneficiary relationship value
func (e *BeneficiaryRelationship) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BeneficiaryRelationship(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BeneficiaryRelationship", str)
	}
	return nil
}

// MarshalGQL converts the beneficiary relationship into a valid JSON string
func (e BeneficiaryRelationship) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// OrganizationType defines the various OrganizationTypes
type OrganizationType string

// OrganizationTypeLimitedCompany is an example of a OrganizationType
const (
	OrganizationTypeLimitedCompany OrganizationType = "LIMITED_COMPANY"
	OrganizationTypeTrust          OrganizationType = "TRUST"
	OrganizationTypeUniversity     OrganizationType = "UNIVERSITY"
)

// AllOrganizationType contains a slice of all OrganizationType
var AllOrganizationType = []OrganizationType{
	OrganizationTypeLimitedCompany,
	OrganizationTypeTrust,
	OrganizationTypeUniversity,
}

// IsValid checks if the OrganizationType is valid
func (e OrganizationType) IsValid() bool {
	switch e {
	case OrganizationTypeLimitedCompany, OrganizationTypeTrust, OrganizationTypeUniversity:
		return true
	}
	return false
}

func (e OrganizationType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an OrganizationType value
func (e *OrganizationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrganizationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OrganizationType", str)
	}
	return nil
}

// MarshalGQL converts OrganizationType into a valid JSON string
func (e OrganizationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// KYCProcessStatus  status for processing KYC for suppliers
type KYCProcessStatus string

// Valid KYCProcessStatus
const (
	KYCProcessStatusApproved KYCProcessStatus = "APPROVED"
	KYCProcessStatusRejected KYCProcessStatus = "REJECTED"
	KYCProcessStatusPending  KYCProcessStatus = "PENDING"
)

// AllKYCProcessStatus ...
var AllKYCProcessStatus = []KYCProcessStatus{
	KYCProcessStatusApproved,
	KYCProcessStatusRejected,
	KYCProcessStatusPending,
}

// IsValid checks if the KYCProcessStatus is valid
func (e KYCProcessStatus) IsValid() bool {
	switch e {
	case KYCProcessStatusApproved, KYCProcessStatusRejected, KYCProcessStatusPending:
		return true
	}
	return false
}

func (e KYCProcessStatus) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an KYCProcessStatus value
func (e *KYCProcessStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = KYCProcessStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid KYCProcessStatus", str)
	}
	return nil
}

// MarshalGQL converts KYCProcessStatus into a valid JSON string
func (e KYCProcessStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// EmploymentType ...
type EmploymentType string

// EmploymentTypeEmployed ..
const (
	EmploymentTypeEmployed     EmploymentType = "EMPLOYED"
	EmploymentTypeSelfEmployed EmploymentType = "SELF_EMPLOYED"
)

// AllEmploymentType ..
var AllEmploymentType = []EmploymentType{
	EmploymentTypeEmployed,
	EmploymentTypeSelfEmployed,
}

// IsValid ..
func (e EmploymentType) IsValid() bool {
	switch e {
	case EmploymentTypeEmployed, EmploymentTypeSelfEmployed:
		return true
	}
	return false
}

func (e EmploymentType) String() string {
	return string(e)
}

// UnmarshalGQL ..
func (e *EmploymentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EmploymentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EmploymentType", str)
	}
	return nil
}

// MarshalGQL ..
func (e EmploymentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

//AgentType is the different kind of agent groups
type AgentType string

// Valid AgentTypes that can possibly be given to a user
const (
	//FreelanceAgent are agents that work at part time with savannah
	FreelanceAgent AgentType = "Independent Agent"

	//CompanyAgent are agents who are fully employed by savannah
	CompanyAgent AgentType = "SIL Agent"
)

// IsValid ..
func (e AgentType) IsValid() bool {
	switch e {
	case FreelanceAgent, CompanyAgent:
		return true
	}
	return false
}

func (e AgentType) String() string {
	return string(e)
}

// UnmarshalGQL ..
func (e *AgentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AgentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AgentType", str)
	}
	return nil
}

// MarshalGQL ..
func (e AgentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Channel can be SMS, Whatsapp, Email etc
type Channel string

//Possible valid channels
const (
	WhatsAppChannel Channel = "WhatsApp"
	SMSChannel      Channel = "SMS"
	EmailChannel    Channel = "Email"
)

//IsValid ...
func (c Channel) IsValid() bool {
	switch c {
	case WhatsAppChannel, SMSChannel, EmailChannel:
		return true
	}
	return false
}

//String ...
func (c Channel) String() string {
	return string(c)
}

//Int ...
func (c Channel) Int() int {
	switch c {
	case WhatsAppChannel:
		return 1
	case SMSChannel:
		return 2
	case EmailChannel:
		return 3
	}
	return 0
}

// UnmarshalGQL ..
func (c *Channel) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}
	*c = Channel(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid Channel", str)
	}
	return nil
}

// MarshalGQL ..
func (c Channel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}
