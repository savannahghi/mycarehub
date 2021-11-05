package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFilterSortDataType_String(t *testing.T) {
	tests := []struct {
		name string
		e    FilterSortDataType
		want string
	}{
		{
			name: "created_at",
			e:    FilterSortDataTypeCreatedAt,
			want: "created_at",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("FilterSortDataType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSortDataType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    FilterSortDataType
		want bool
	}{
		{
			name: "valid type",
			e:    FilterSortDataTypeCreatedAt,
			want: true,
		},
		{
			name: "invalid type",
			e:    FilterSortDataType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("FilterSortDataType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSortDataType_UnmarshalGQL(t *testing.T) {
	pmtc := FilterSortDataTypeCreatedAt
	invalid := FilterSortDataType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *FilterSortDataType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &pmtc,
			args: args{
				v: "created_at",
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
				t.Errorf("FilterSortDataType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterSortDataType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     FilterSortDataType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     FilterSortDataTypeCreatedAt,
			b:     w,
			wantW: strconv.Quote("created_at"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FilterSortDataType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
