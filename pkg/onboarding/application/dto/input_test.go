package dto

import (
	"net/url"
	"reflect"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestBusinessPartnerFilterInput_ToURLValues(t *testing.T) {
	var (
		search    = "data"
		name      = "somename"
		sladeCode = "slasde"
	)

	correctFilters := BusinessPartnerFilterInput{
		Search:    &search,
		Name:      &name,
		SladeCode: &sladeCode,
	}

	failingFilters := BusinessPartnerFilterInput{
		Search: &search,
	}

	expectedFilter := url.Values{
		"search":     []string{search},
		"name":       []string{name},
		"slade_code": []string{sladeCode},
	}

	tests := []struct {
		name       string
		filter     BusinessPartnerFilterInput
		wantValues url.Values
		wantError  bool
	}{
		{
			name:       "success url values transformation",
			filter:     correctFilters,
			wantValues: expectedFilter,
			wantError:  false,
		},
		{
			name:       "bad filter data",
			filter:     failingFilters,
			wantValues: expectedFilter,
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.wantValues, tt.filter.ToURLValues()) && !tt.wantError {
				t.Errorf("BusinessPartnerFilterInput.ToURLValues() = %v, want %v", tt.filter.ToURLValues(), tt.wantValues)
			}
		})
	}
}

func TestBusinessPartnerSortInput_ToURLValues(t *testing.T) {

	var name base.SortOrder = base.SortOrderAsc
	var sladeCode base.SortOrder = base.SortOrderAsc
	var sladeCode2 base.SortOrder = base.SortOrderDesc

	correctFilters := BusinessPartnerSortInput{
		Name:      &name,
		SladeCode: &sladeCode,
	}

	failingFilters := BusinessPartnerSortInput{
		Name:      &name,
		SladeCode: &sladeCode2,
	}

	expectedSortValue := url.Values{
		"order_by":   []string{"name"},
		"slade_code": []string{"number"},
	}

	tests := []struct {
		name       string
		filter     BusinessPartnerSortInput
		wantValues url.Values
		wantError  bool
	}{
		{
			name:       "success passing sort filter",
			filter:     correctFilters,
			wantValues: expectedSortValue,
			wantError:  false,
		},
		{
			name:       "bad filters",
			filter:     failingFilters,
			wantValues: expectedSortValue,
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.wantValues, tt.filter.ToURLValues()) && !tt.wantError {
				t.Errorf("BusinessPartnerSortInput.ToURLValues() = %v, want %v", tt.filter.ToURLValues(), tt.wantValues)
			}
		})
	}
}

func TestBranchSortInput_ToURLValues(t *testing.T) {
	var name base.SortOrder = base.SortOrderAsc
	var sladeCode base.SortOrder = base.SortOrderDesc
	var sladeCode2 base.SortOrder = base.SortOrderAsc

	correctFilters := BranchSortInput{
		Name:      &name,
		SladeCode: &sladeCode,
	}

	failingFilters := BranchSortInput{
		Name:      &name,
		SladeCode: &sladeCode2,
	}

	expectedFilters := url.Values{
		"order_by":   []string{"name"},
		"slade_code": []string{"-number"},
	}

	tests := []struct {
		name       string
		filter     BranchSortInput
		wantValues url.Values
		wantError  bool
	}{
		{
			name:       "success building filters",
			filter:     correctFilters,
			wantValues: expectedFilters,
			wantError:  false,
		},
		{
			name:       "bad filters",
			filter:     failingFilters,
			wantValues: expectedFilters,
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &BranchSortInput{
				Name:      tt.filter.Name,
				SladeCode: tt.filter.SladeCode,
			}
			if gotValues := i.ToURLValues(); !reflect.DeepEqual(gotValues, tt.wantValues) && !tt.wantError {
				t.Errorf("BranchSortInput.ToURLValues() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestBranchFilterInput_ToURLValues(t *testing.T) {

	var (
		search               = "somesearch"
		sladeCode            = "sladecode"
		parentOrganizationID = "parentorg"
	)

	correctFilters := BranchFilterInput{
		Search:               &search,
		SladeCode:            &sladeCode,
		ParentOrganizationID: &parentOrganizationID,
	}

	failingFilter := BranchFilterInput{
		Search: &search,
	}

	expectedFilters := url.Values{
		"search":     []string{search},
		"slade_code": []string{sladeCode},
		"parent":     []string{parentOrganizationID},
	}

	tests := []struct {
		name       string
		filter     BranchFilterInput
		wantValues url.Values
		wantError  bool
	}{
		{
			name:       "success transforming url values ",
			filter:     correctFilters,
			wantValues: expectedFilters,
			wantError:  false,
		},
		{
			name:       "bad filters",
			filter:     failingFilter,
			wantValues: expectedFilters,
			wantError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &BranchFilterInput{
				Search:               tt.filter.Search,
				SladeCode:            tt.filter.SladeCode,
				ParentOrganizationID: tt.filter.ParentOrganizationID,
			}
			if got := i.ToURLValues(); !reflect.DeepEqual(got, tt.wantValues) && !tt.wantError {
				t.Errorf("BranchFilterInput.ToURLValues() = %v, want %v", got, tt.wantValues)
			}
		})
	}
}
