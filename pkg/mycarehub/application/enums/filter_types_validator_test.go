package enums

import (
	"reflect"
	"testing"
)

func TestValidateFilterSortCategoriesss(t *testing.T) {
	type args struct {
		category FilterSortCategoryType
		filter   FilterSortDataType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			args: args{
				category: FilterSortCategoryTypeFacility,
				filter:   FilterSortDataTypeName,
			},
			wantErr: false,
		},
		{
			name: "invalid category type",
			args: args{
				category: FilterSortCategoryType("invalid"),
				filter:   FilterSortDataTypeName,
			},
			wantErr: true,
		},
		{
			name: "invalid filter sorts type",
			args: args{
				category: FilterSortCategoryTypeFacility,
				filter:   FilterSortDataType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid filter not in category",
			args: args{
				category: FilterSortCategoryTypeFacility,
				filter:   FilterSortDataTypeCreatedAt,
			},
			wantErr: true,
		},
		{
			name:    "empty params passed",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFilterSortCategories(tt.args.category, tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilterSortCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findSelectedCategoryFilters(t *testing.T) {

	filterSortDataTypeCategories := []FilterSortDataTypeCategory{
		{
			FilterSortCategory: FilterSortCategoryTypeFacility,
			FilterSort:         FacilityFilterDataTypes,
		},
	}
	want1 := FilterSortDataTypeCategory{
		FilterSortCategory: FilterSortCategoryTypeFacility,
		FilterSort:         FacilityFilterDataTypes,
	}
	type args struct {
		filterCategories []FilterSortDataTypeCategory
		categoryInput    FilterSortCategoryType
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *FilterSortDataTypeCategory
	}{
		{
			name: "valid type",
			args: args{
				filterCategories: filterSortDataTypeCategories,
				categoryInput:    FilterSortCategoryTypeFacility,
			},
			want:  true,
			want1: &want1,
		},
		{
			name: "invalid country type",
			args: args{
				filterCategories: filterSortDataTypeCategories,
				categoryInput:    FilterSortCategoryType("invalid"),
			},
			want:  false,
			want1: nil,
		},
		{
			name: "invalid country list type",
			args: args{
				filterCategories: []FilterSortDataTypeCategory{},
				categoryInput:    FilterSortCategoryTypeFacility,
			},
			want:  false,
			want1: nil,
		},
		{
			name:  "empty args",
			args:  args{},
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findSelectedCategoryFilters(tt.args.filterCategories, tt.args.categoryInput)
			if got != tt.want {
				t.Errorf("findSelectedCategoryFilters() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("findSelectedCategoryFilters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
