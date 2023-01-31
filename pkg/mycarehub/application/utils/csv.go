package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ReadCSVFile reads the content of a csv file
func ReadCSVFile(path string) (*csv.Reader, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(file)

	return csvReader, nil
}

// ParseFacilitiesFromCSV maps the values of the csv file to the Facilities object
func ParseFacilitiesFromCSV(path string) ([]*domain.Facility, error) {
	csvReader, err := ReadCSVFile(path)
	if err != nil {
		return nil, err
	}

	var (
		count  int
		labels []string
	)
	facilities := []*domain.Facility{}
	facilityCSVOutput := dto.FacilityCSVOutput{}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		count++

		if count == 1 {
			labels = row
			err = facilityCSVOutput.ValidateLabels(labels)
			if err != nil {
				return nil, err
			}
		}

		if count > 1 {
			facility, err := facilityCSVOutput.ParseValues(labels, row)
			if err != nil {
				return nil, err
			}

			facilities = append(facilities, &domain.Facility{
				Name:        facility.Name,
				Phone:       facility.Contact,
				Active:      true,
				Country:     facility.Country,
				Description: fmt.Sprintf("%s %s owned by %s and regulated by %s", facility.Level, facility.FacilityType, facility.OwnerType, facility.RegulatoryBody),
				Identifier: domain.FacilityIdentifier{
					Active: true,
					Type:   facility.IdentifierType,
					Value:  facility.Code,
				},
			})
		}
	}
	return facilities, nil
}
