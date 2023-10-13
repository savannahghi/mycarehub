package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestDayOfWeek_String(t *testing.T) {
	tests := []struct {
		name string
		e    DayOfWeek
		want string
	}{
		{
			name: "MONDAY",
			e:    DayOfWeekMonday,
			want: "MONDAY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("DayOfWeek.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    DayOfWeek
		want bool
	}{
		{
			name: "valid type",
			e:    DayOfWeekFriday,
			want: true,
		},
		{
			name: "invalid type",
			e:    DayOfWeek("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("DayOfWeek.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayOfWeek_UnmarshalGQL(t *testing.T) {
	value := DayOfWeekSunday
	invalid := DayOfWeek("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *DayOfWeek
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "SUNDAY",
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
				t.Errorf("DayOfWeek.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayOfWeek_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     DayOfWeek
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     DayOfWeekSunday,
			b:     w,
			wantW: strconv.Quote("SUNDAY"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DayOfWeek.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
