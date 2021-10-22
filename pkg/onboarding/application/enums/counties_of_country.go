package enums

import "fmt"

// CountiesOfCountry defines the counties of a country
type CountiesOfCountry struct {
	Country  CountryType
	Counties []CountyType
}

// Create a struct with a combined list of counties of countries
var countiesOfCountries = []CountiesOfCountry{
	{
		Country:  CountryTypeKenya,
		Counties: KenyanCounties,
	},
	// Other CountiesOfCountries
}

// ValidateCountiesOfCountries validates the county passed to a country is valid
func ValidateCountiesOfCountries(country CountryType, county CountyType) error {
	// ensure we are working with a correct country type to begin with
	if !country.IsValid() {
		return fmt.Errorf("failed to validate country: %s", country)
	}
	// Validate the county too
	if !county.IsValid() {
		return fmt.Errorf("failed to validate county: %s", county)
	}

	ok, counties := findSelectedCountryCounties(countiesOfCountries, country)
	if !ok {
		return fmt.Errorf("failed to find selected country's counties: %s", county)
	}

	err := findCounty(counties.Counties, county)
	if err != nil {
		return fmt.Errorf("failed to find county: %s", err)
	}
	return nil
}

// finds the selected country te ensure it's part of the enum, then return the respective counties it has
func findSelectedCountryCounties(countriesCounties []CountiesOfCountry, countryInput CountryType) (bool, *CountiesOfCountry) {
	for i, countryCounty := range countriesCounties {
		if countryCounty.Country == countryInput {
			return true, &countiesOfCountries[i]
		}
	}
	return false, nil
}

// checks whether the county provided is present in the list of the selected country's counties
func findCounty(counties []CountyType, countyInput CountyType) error {
	for _, county := range counties {
		if county == countyInput {
			return nil
		}
	}
	return fmt.Errorf("failed to find county: %s", countyInput)
}
