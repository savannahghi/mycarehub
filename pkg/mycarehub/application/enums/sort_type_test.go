package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestSortDataType_String(t *testing.T) {
	tests := []struct {
		name string
		e    SortDataType
		want string
	}{
		{
			name: "asc",
			e:    SortDataTypeAsc,
			want: "asc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SortDataType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDataType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    SortDataType
		want bool
	}{
		{
			name: "valid type",
			e:    SortDataTypeAsc,
			want: true,
		},
		{
			name: "invalid type",
			e:    SortDataType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("SortDataType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortDataType_UnmarshalGQL(t *testing.T) {
	value := SortDataTypeAsc
	invalid := SortDataType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *SortDataType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "asc",
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
				t.Errorf("SortDataType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSortDataType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     SortDataType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     SortDataTypeAsc,
			b:     w,
			wantW: strconv.Quote("asc"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SortDataType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
