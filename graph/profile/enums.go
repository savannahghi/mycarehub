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
