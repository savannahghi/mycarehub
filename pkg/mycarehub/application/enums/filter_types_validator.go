package enums

import "fmt"

// FilterDataTypeCategory defines a collection of filter data type categories
type FilterDataTypeCategory struct {
	FilterCategory FilterCategoryType
	Filters        []FilterDataType
}

// Create a struct with a combined list of counties of countries
var filterDataTypeCategories = []FilterDataTypeCategory{
	{
		FilterCategory: FilterCategoryTypeFacility,
		Filters:        FacilityFilterDataTypes,
	},
	// Other Filter categories
}

// ValidateFilterCategories validates the filter belongs to the specified category
//  (translates to having more filter types for each table)
func ValidateFilterCategories(category FilterCategoryType, filter FilterDataType) error {
	// ensure we are working with a correct filter category
	if !category.IsValid() {
		return fmt.Errorf("failed to validate category: %s", category)
	}
	// Validate the filter being passed
	if !filter.IsValid() {
		return fmt.Errorf("failed to validate filter: %s", filter)
	}

	ok, filters := findSelectedCategoryFilters(filterDataTypeCategories, category)
	if !ok {
		return fmt.Errorf("failed to find selected category filters: %s", filter)
	}

	err := findFilter(filters.Filters, filter)
	if err != nil {
		return fmt.Errorf("failed to find filter: %s", err)
	}
	return nil
}

// finds the selected filter category te ensure it's part of the enum, then return the respective filters for that category
func findSelectedCategoryFilters(filterCategories []FilterDataTypeCategory, categoryInput FilterCategoryType) (bool, *FilterDataTypeCategory) {
	for i, filterCategory := range filterCategories {
		if filterCategory.FilterCategory == categoryInput {
			return true, &filterDataTypeCategories[i]
		}
	}
	return false, nil
}

// checks whether the filter provided is present in the list of the selected category of filters
func findFilter(filters []FilterDataType, filterInput FilterDataType) error {
	for _, filter := range filters {
		if filter == filterInput {
			return nil
		}
	}
	return fmt.Errorf("failed to find filter: %s", filterInput)
}
