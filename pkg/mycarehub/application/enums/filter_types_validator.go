package enums

import "fmt"

// FilterSortDataTypeCategory defines a collection of filter data type categories
type FilterSortDataTypeCategory struct {
	FilterSortCategory FilterSortCategoryType
	FilterSort         []FilterSortDataType
}

// Create a struct with a combined list of counties of sort and filter
var filterSortDataTypeCategories = []FilterSortDataTypeCategory{
	{
		FilterSortCategory: FilterSortCategoryTypeFacility,
		FilterSort:         FacilityFilterDataTypes,
	},
	{
		FilterSortCategory: FilterSortCategoryTypeSortFacility,
		FilterSort:         FacilitySortDataTypes,
	},
	// Other Filter/Sort categories
}

// ValidateFilterSortCategories validates the filter belongs to the specified category
//  (translates to having more filter types for each table)
func ValidateFilterSortCategories(category FilterSortCategoryType, filter FilterSortDataType) error {
	// ensure we are working with a correct filter category
	if !category.IsValid() {
		return fmt.Errorf("failed to validate category: %s", category)
	}
	// Validate the filter being passed
	if !filter.IsValid() {
		return fmt.Errorf("failed to validate filter: %s", filter)
	}

	ok, filters := findSelectedCategoryFilters(filterSortDataTypeCategories, category)
	if !ok {
		return fmt.Errorf("failed to find selected category filters: %s", filter)
	}

	err := findFilter(filters.FilterSort, filter)
	if err != nil {
		return fmt.Errorf("failed to find filter: %s", err)
	}
	return nil
}

// finds the selected filter category te ensure it's part of the enum, then return the respective filters for that category
func findSelectedCategoryFilters(filterCategories []FilterSortDataTypeCategory, categoryInput FilterSortCategoryType) (bool, *FilterSortDataTypeCategory) {
	for i, filterCategory := range filterCategories {
		if filterCategory.FilterSortCategory == categoryInput {
			return true, &filterSortDataTypeCategories[i]
		}
	}
	return false, nil
}

// checks whether the filter provided is present in the list of the selected category of filters
func findFilter(filters []FilterSortDataType, filterInput FilterSortDataType) error {
	for _, filter := range filters {
		if filter == filterInput {
			return nil
		}
	}
	return fmt.Errorf("failed to find filter: %s", filterInput)
}
