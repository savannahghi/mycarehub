package domain

import (
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestPagination_GetOffset(t *testing.T) {
	nextPage := 2
	previousPage := 1
	sort := SortParam{
		Field:     enums.FilterSortDataTypeActive,
		Direction: enums.SortDataTypeAsc,
	}
	type fields struct {
		Limit        int
		CurrentPage  int
		Count        int64
		TotalPages   int
		NextPage     *int
		PreviousPage *int
		Sort         *SortParam
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "default case",
			fields: fields{
				Limit:        2,
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:        tt.fields.Limit,
				CurrentPage:  tt.fields.CurrentPage,
				Count:        tt.fields.Count,
				TotalPages:   tt.fields.TotalPages,
				NextPage:     tt.fields.NextPage,
				PreviousPage: tt.fields.PreviousPage,
				Sort:         tt.fields.Sort,
			}
			if got := p.GetOffset(); got != tt.want {
				t.Errorf("Pagination.GetOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_GetLimit(t *testing.T) {
	nextPage := 2
	previousPage := 1
	sort := SortParam{
		Field:     enums.FilterSortDataTypeActive,
		Direction: enums.SortDataTypeAsc,
	}
	type fields struct {
		Limit        int
		CurrentPage  int
		Count        int64
		TotalPages   int
		NextPage     *int
		PreviousPage *int
		Sort         *SortParam
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "default case",
			fields: fields{
				Limit:        2,
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: 2,
		},
		{
			name: "happy case: no limit passed",
			fields: fields{
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:        tt.fields.Limit,
				CurrentPage:  tt.fields.CurrentPage,
				Count:        tt.fields.Count,
				TotalPages:   tt.fields.TotalPages,
				NextPage:     tt.fields.NextPage,
				PreviousPage: tt.fields.PreviousPage,
				Sort:         tt.fields.Sort,
			}
			if got := p.GetLimit(); got != tt.want {
				t.Errorf("Pagination.GetLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_GetPage(t *testing.T) {
	nextPage := 2
	previousPage := 1
	sort := SortParam{
		Field:     enums.FilterSortDataTypeActive,
		Direction: enums.SortDataTypeAsc,
	}
	type fields struct {
		Limit        int
		CurrentPage  int
		Count        int64
		TotalPages   int
		NextPage     *int
		PreviousPage *int
		Sort         *SortParam
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "default case",
			fields: fields{
				Limit:        2,
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: 1,
		},
		{
			name: "happy case: no current page passed",
			fields: fields{
				Limit:        2,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:        tt.fields.Limit,
				CurrentPage:  tt.fields.CurrentPage,
				Count:        tt.fields.Count,
				TotalPages:   tt.fields.TotalPages,
				NextPage:     tt.fields.NextPage,
				PreviousPage: tt.fields.PreviousPage,
				Sort:         tt.fields.Sort,
			}
			if got := p.GetPage(); got != tt.want {
				t.Errorf("Pagination.GetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_GetSort(t *testing.T) {
	nextPage := 2
	previousPage := 1
	sort := SortParam{
		Field:     enums.FilterSortDataTypeActive,
		Direction: enums.SortDataTypeAsc,
	}
	type fields struct {
		Limit        int
		CurrentPage  int
		Count        int64
		TotalPages   int
		NextPage     *int
		PreviousPage *int
		Sort         *SortParam
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default case",
			fields: fields{
				Limit:        2,
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort:         &sort,
			},
			want: "active asc",
		},
		{
			name: "happy case: no sort params passed",
			fields: fields{
				Limit:        2,
				CurrentPage:  1,
				Count:        4,
				TotalPages:   2,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				Sort: &SortParam{
					Field:     "",
					Direction: "",
				},
			},
			want: "updated desc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pagination{
				Limit:        tt.fields.Limit,
				CurrentPage:  tt.fields.CurrentPage,
				Count:        tt.fields.Count,
				TotalPages:   tt.fields.TotalPages,
				NextPage:     tt.fields.NextPage,
				PreviousPage: tt.fields.PreviousPage,
				Sort:         tt.fields.Sort,
			}
			if got := p.GetSort(); got != tt.want {
				t.Errorf("Pagination.GetSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
