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

// PractitionerSpecialty is a list of recognised health worker specialties.
//
// See: https://medicalboard.co.ke/resources_page/gazetted-specialties/
type PractitionerSpecialty string

// list of known practitioner specialties
const (
	PractitionerSpecialtyUnspecified                     PractitionerSpecialty = "UNSPECIFIED"
	PractitionerSpecialtyAnaesthesia                     PractitionerSpecialty = "ANAESTHESIA"
	PractitionerSpecialtyCardiothoracicSurgery           PractitionerSpecialty = "CARDIOTHORACIC_SURGERY"
	PractitionerSpecialtyClinicalMedicalGenetics         PractitionerSpecialty = "CLINICAL_MEDICAL_GENETICS"
	PractitionerSpecialtyClincicalPathology              PractitionerSpecialty = "CLINCICAL_PATHOLOGY"
	PractitionerSpecialtyGeneralPathology                PractitionerSpecialty = "GENERAL_PATHOLOGY"
	PractitionerSpecialtyAnatomicPathology               PractitionerSpecialty = "ANATOMIC_PATHOLOGY"
	PractitionerSpecialtyClinicalOncology                PractitionerSpecialty = "CLINICAL_ONCOLOGY"
	PractitionerSpecialtyDermatology                     PractitionerSpecialty = "DERMATOLOGY"
	PractitionerSpecialtyEarNoseAndThroat                PractitionerSpecialty = "EAR_NOSE_AND_THROAT"
	PractitionerSpecialtyEmergencyMedicine               PractitionerSpecialty = "EMERGENCY_MEDICINE"
	PractitionerSpecialtyFamilyMedicine                  PractitionerSpecialty = "FAMILY_MEDICINE"
	PractitionerSpecialtyGeneralSurgery                  PractitionerSpecialty = "GENERAL_SURGERY"
	PractitionerSpecialtyGeriatrics                      PractitionerSpecialty = "GERIATRICS"
	PractitionerSpecialtyImmunology                      PractitionerSpecialty = "IMMUNOLOGY"
	PractitionerSpecialtyInfectiousDisease               PractitionerSpecialty = "INFECTIOUS_DISEASE"
	PractitionerSpecialtyInternalMedicine                PractitionerSpecialty = "INTERNAL_MEDICINE"
	PractitionerSpecialtyMicrobiology                    PractitionerSpecialty = "MICROBIOLOGY"
	PractitionerSpecialtyNeurosurgery                    PractitionerSpecialty = "NEUROSURGERY"
	PractitionerSpecialtyObstetricsAndGynaecology        PractitionerSpecialty = "OBSTETRICS_AND_GYNAECOLOGY"
	PractitionerSpecialtyOccupationalMedicine            PractitionerSpecialty = "OCCUPATIONAL_MEDICINE"
	PractitionerSpecialtyOpthalmology                    PractitionerSpecialty = "OPTHALMOLOGY"
	PractitionerSpecialtyOrthopaedicSurgery              PractitionerSpecialty = "ORTHOPAEDIC_SURGERY"
	PractitionerSpecialtyOncology                        PractitionerSpecialty = "ONCOLOGY"
	PractitionerSpecialtyOncologyRadiotherapy            PractitionerSpecialty = "ONCOLOGY_RADIOTHERAPY"
	PractitionerSpecialtyPaediatricsAndChildHealth       PractitionerSpecialty = "PAEDIATRICS_AND_CHILD_HEALTH"
	PractitionerSpecialtyPalliativeMedicine              PractitionerSpecialty = "PALLIATIVE_MEDICINE"
	PractitionerSpecialtyPlasticAndReconstructiveSurgery PractitionerSpecialty = "PLASTIC_AND_RECONSTRUCTIVE_SURGERY"
	PractitionerSpecialtyPsychiatry                      PractitionerSpecialty = "PSYCHIATRY"
	PractitionerSpecialtyPublicHealth                    PractitionerSpecialty = "PUBLIC_HEALTH"
	PractitionerSpecialtyRadiology                       PractitionerSpecialty = "RADIOLOGY"
	PractitionerSpecialtyUrology                         PractitionerSpecialty = "UROLOGY"
)

// AllPractitionerSpecialty is the set of known practitioner specialties
var AllPractitionerSpecialty = []PractitionerSpecialty{
	PractitionerSpecialtyUnspecified,
	PractitionerSpecialtyAnaesthesia,
	PractitionerSpecialtyCardiothoracicSurgery,
	PractitionerSpecialtyClinicalMedicalGenetics,
	PractitionerSpecialtyClincicalPathology,
	PractitionerSpecialtyGeneralPathology,
	PractitionerSpecialtyAnatomicPathology,
	PractitionerSpecialtyClinicalOncology,
	PractitionerSpecialtyDermatology,
	PractitionerSpecialtyEarNoseAndThroat,
	PractitionerSpecialtyEmergencyMedicine,
	PractitionerSpecialtyFamilyMedicine,
	PractitionerSpecialtyGeneralSurgery,
	PractitionerSpecialtyGeriatrics,
	PractitionerSpecialtyImmunology,
	PractitionerSpecialtyInfectiousDisease,
	PractitionerSpecialtyInternalMedicine,
	PractitionerSpecialtyMicrobiology,
	PractitionerSpecialtyNeurosurgery,
	PractitionerSpecialtyObstetricsAndGynaecology,
	PractitionerSpecialtyOccupationalMedicine,
	PractitionerSpecialtyOpthalmology,
	PractitionerSpecialtyOrthopaedicSurgery,
	PractitionerSpecialtyOncology,
	PractitionerSpecialtyOncologyRadiotherapy,
	PractitionerSpecialtyPaediatricsAndChildHealth,
	PractitionerSpecialtyPalliativeMedicine,
	PractitionerSpecialtyPlasticAndReconstructiveSurgery,
	PractitionerSpecialtyPsychiatry,
	PractitionerSpecialtyPublicHealth,
	PractitionerSpecialtyRadiology,
	PractitionerSpecialtyUrology,
}

// IsValid returns True if the practitioner specialty is valid
func (e PractitionerSpecialty) IsValid() bool {
	switch e {
	case PractitionerSpecialtyUnspecified, PractitionerSpecialtyAnaesthesia, PractitionerSpecialtyCardiothoracicSurgery, PractitionerSpecialtyClinicalMedicalGenetics, PractitionerSpecialtyClincicalPathology, PractitionerSpecialtyGeneralPathology, PractitionerSpecialtyAnatomicPathology, PractitionerSpecialtyClinicalOncology, PractitionerSpecialtyDermatology, PractitionerSpecialtyEarNoseAndThroat, PractitionerSpecialtyEmergencyMedicine, PractitionerSpecialtyFamilyMedicine, PractitionerSpecialtyGeneralSurgery, PractitionerSpecialtyGeriatrics, PractitionerSpecialtyImmunology, PractitionerSpecialtyInfectiousDisease, PractitionerSpecialtyInternalMedicine, PractitionerSpecialtyMicrobiology, PractitionerSpecialtyNeurosurgery, PractitionerSpecialtyObstetricsAndGynaecology, PractitionerSpecialtyOccupationalMedicine, PractitionerSpecialtyOpthalmology, PractitionerSpecialtyOrthopaedicSurgery, PractitionerSpecialtyOncology, PractitionerSpecialtyOncologyRadiotherapy, PractitionerSpecialtyPaediatricsAndChildHealth, PractitionerSpecialtyPalliativeMedicine, PractitionerSpecialtyPlasticAndReconstructiveSurgery, PractitionerSpecialtyPsychiatry, PractitionerSpecialtyPublicHealth, PractitionerSpecialtyRadiology, PractitionerSpecialtyUrology:
		return true
	}
	return false
}

func (e PractitionerSpecialty) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a practitioner specialty
func (e *PractitionerSpecialty) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PractitionerSpecialty(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PractitionerSpecialty", str)
	}
	return nil
}

// MarshalGQL writes the practitioner specialty to the supplied writer
func (e PractitionerSpecialty) MarshalGQL(w io.Writer) {
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

// IsValid returs true for valid ratings
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

// MarshalGQL coverts the rating into a valid JSON string
func (e FivePointRating) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}
