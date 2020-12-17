package profile

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

// SignUpMethod defines the various frontend sign up options
type SignUpMethod string

// SignUpMethodAnonymous is a constant of all known sign up methods
const (
	SignUpMethodAnonymous SignUpMethod = "anonymous"
	SignUpMethodApple     SignUpMethod = "apple"
	SignUpMethodFacebook  SignUpMethod = "facebook"
	SignUpMethodGoogle    SignUpMethod = "google"
	SignUpMethodPhone     SignUpMethod = "phone"
)

// AllSignUpMethod is a list of all known sign up methods
var AllSignUpMethod = []SignUpMethod{
	SignUpMethodAnonymous,
	SignUpMethodApple,
	SignUpMethodFacebook,
	SignUpMethodGoogle,
	SignUpMethodPhone,
}

// IsValid returns true for valid sign up method
func (e SignUpMethod) IsValid() bool {
	switch e {
	case SignUpMethodAnonymous, SignUpMethodApple, SignUpMethodFacebook, SignUpMethodGoogle, SignUpMethodPhone:
		return true
	}
	return false
}

func (e SignUpMethod) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a signup method value
func (e *SignUpMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SignUpMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SignUpMethod", str)
	}
	return nil
}

// MarshalGQL converts the sign up method into a valid JSON string
func (e SignUpMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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

// AccountType defines the various supplier account types
type AccountType string

// AccountTypeIndivdual is an example of a suppiler account type
const (
	AccountTypeIndividual   AccountType = "INDIVIDUAL"
	AccountTypeOrganisation AccountType = "ORGANISATION"
)

// AllAccountType is a slice that represents all the account types
var AllAccountType = []AccountType{
	AccountTypeIndividual,
	AccountTypeOrganisation,
}

// IsValid checks if the account type is valid
func (e AccountType) IsValid() bool {
	switch e {
	case AccountTypeIndividual, AccountTypeOrganisation:
		return true
	}
	return false
}

func (e AccountType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into a account type value
func (e *AccountType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccountType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccountType", str)
	}
	return nil
}

// MarshalGQL converts AccountType into a valid JSON string
func (e AccountType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IdentificationDocType defines the various supplier IdentificationDocTypes
type IdentificationDocType string

// IdentificationDocTypeNationalid is an example of a IdentificationDocType
const (
	IdentificationDocTypeNationalid IdentificationDocType = "NATIONALID"
	IdentificationDocTypePassport   IdentificationDocType = "PASSPORT"
	IdentificationDocTypeMilitary   IdentificationDocType = "MILITARY"
)

// AllIdentificationDocType contains a slice of all IdentificationDocTypes
var AllIdentificationDocType = []IdentificationDocType{
	IdentificationDocTypeNationalid,
	IdentificationDocTypePassport,
	IdentificationDocTypeMilitary,
}

// IsValid checks if the IdentificationDocType is valid
func (e IdentificationDocType) IsValid() bool {
	switch e {
	case IdentificationDocTypeNationalid, IdentificationDocTypePassport, IdentificationDocTypeMilitary:
		return true
	}
	return false
}

func (e IdentificationDocType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an IdentificationDocType value
func (e *IdentificationDocType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentificationDocType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentificationDocType", str)
	}
	return nil
}

// MarshalGQL converts IdentificationDocType into a valid JSON string
func (e IdentificationDocType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// PartnerType defines the different partner types we have in Be.Well
type PartnerType string

// PartnerTypeRider is an example of a partner type who is involved in delivery of goods
const (
	PartnerTypeRider          PartnerType = "RIDER"
	PartnerTypePractitioner   PartnerType = "PRACTITIONER"
	PartnerTypeProvider       PartnerType = "PROVIDER"
	PartnerTypePharmaceutical PartnerType = "PHARMACEUTICAL"
	PartnerTypeCoach          PartnerType = "COACH"
	PartnerTypeNutrition      PartnerType = "NUTRITION"
	PartnerTypeConsumer       PartnerType = "CONSUMER"
)

// AllPartnerType represents a list of the partner types we offer
var AllPartnerType = []PartnerType{
	PartnerTypeRider,
	PartnerTypePractitioner,
	PartnerTypeProvider,
	PartnerTypePharmaceutical,
	PartnerTypeCoach,
	PartnerTypeNutrition,
	PartnerTypeConsumer,
}

// IsValid checks if a partner type is valid or not
func (e PartnerType) IsValid() bool {
	switch e {
	case PartnerTypeRider, PartnerTypePractitioner, PartnerTypeProvider, PartnerTypePharmaceutical, PartnerTypeCoach, PartnerTypeNutrition, PartnerTypeConsumer:
		return true
	}
	return false
}

func (e PartnerType) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an correct partner type value
func (e *PartnerType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PartnerType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PartnerType", str)
	}
	return nil
}

// MarshalGQL converts partner type into a valid JSON string
func (e PartnerType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// OrganizationType defines the various OrganizationTypes
type OrganizationType string

// IdentificationDocTypeNationalid is an example of a IdentificationDocType
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

// IsValid checks if the IdentificationDocType is valid
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

// UnmarshalGQL converts the input, if valid, into an IdentificationDocType value
func (e *OrganizationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OrganizationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentificationDocType", str)
	}
	return nil
}

// MarshalGQL converts IdentificationDocType into a valid JSON string
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

// IsValid checks if the IdentificationDocType is valid
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

// UnmarshalGQL converts the input, if valid, into an IdentificationDocType value
func (e *KYCProcessStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = KYCProcessStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentificationDocType", str)
	}
	return nil
}

// MarshalGQL converts IdentificationDocType into a valid JSON string
func (e KYCProcessStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
