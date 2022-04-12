package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFilterSortCategoryType_String(t *testing.T) {
	tests := []struct {
		name string
		e    FilterSortCategoryType
		want string
	}{
		{
			name: "SortFacility",
			e:    FilterSortCategoryTypeSortFacility,
			want: "SortFacility",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("FilterSortCategoryType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSortCategoryType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    FilterSortCategoryType
		want bool
	}{
		{
			name: "valid type",
			e:    FilterSortCategoryTypeSortFacility,
			want: true,
		},
		{
			name: "invalid type",
			e:    FilterSortCategoryType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("FilterSortCategoryType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSortCategoryType_UnmarshalGQL(t *testing.T) {
	value := FilterSortCategoryTypeSortFacility
	invalid := FilterSortCategoryType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *FilterSortCategoryType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "SortFacility",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FilterSortCategoryType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterSortCategoryType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     FilterSortCategoryType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     FilterSortCategoryTypeSortFacility,
			b:     w,
			wantW: strconv.Quote("SortFacility"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FilterSortCategoryType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
