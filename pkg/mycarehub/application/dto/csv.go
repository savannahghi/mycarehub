package dto

import (
	"encoding/json"
	"fmt"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	validator "gopkg.in/go-playground/validator.v9"
)

// FacilityCSVOutput is a struct that stores the output of facility csv values
type FacilityCSVOutput struct {
	Code            string                       `json:"code" validate:"required"`
	IdentifierType  enums.FacilityIdentifierType `json:"identifierType" validate:"required"`
	Name            string                       `json:"name" validate:"required"`
	Level           string                       `json:"level" validate:"required"`
	FacilityType    string                       `json:"facilityType" validate:"required"`
	OwnerType       string                       `json:"ownerType" validate:"required"`
	RegulatoryBody  string                       `json:"regulatoryBody" validate:"required"`
	Country         string                       `json:"country" validate:"required"`
	County          string                       `json:"county" validate:"required"`
	OperationStatus string                       `json:"operationStatus" validate:"required"`
	Contact         string                       `json:"contact" validate:"required"`
}

// ValidateLabels ensures the labels of the facility csv are valid. the json tag matches the respective label value
func (f *FacilityCSVOutput) ValidateLabels(labels []string) error {
	if len(labels) != 11 {
		return fmt.Errorf("invalid facility csv: invalid label length: expected 11, got %v", len(labels))
	}
	var labelsObj map[string]interface{}

	bs, err := json.Marshal(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &labelsObj)
	if err != nil {
		return err
	}

	for _, label := range labels {
		if _, ok := labelsObj[label]; !ok {
			return fmt.Errorf("invalid facility csv: invalid label: %v", label)
		}
	}

	return nil
}

// ParseValues transforms the actual facility values from row 2 of the csv
func (f *FacilityCSVOutput) ParseValues(labels []string, values []string) (*FacilityCSVOutput, error) {
	if len(values) != 11 {
		return nil, fmt.Errorf("invalid facility csv: invalid values length: expected 11, got %v", len(values))
	}

	_, err := converterandformatter.NormalizeMSISDN(values[10])
	if err != nil {
		return nil, err
	}

	if ok := enums.FacilityIdentifierType(values[1]).IsValid(); !ok {
		return nil, fmt.Errorf("invalid facility identifier type: %v", values[1])
	}

	f = &FacilityCSVOutput{
		Code:            values[0],
		IdentifierType:  enums.FacilityIdentifierType(values[1]),
		Name:            values[2],
		Level:           values[3],
		FacilityType:    values[4],
		OwnerType:       values[5],
		RegulatoryBody:  values[6],
		Country:         values[7],
		County:          values[8],
		OperationStatus: values[9],
		Contact:         values[10],
	}

	v := validator.New()
	err = v.Struct(f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
