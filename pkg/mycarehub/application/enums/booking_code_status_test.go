package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestBookingCodeStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    BookingCodeStatus
		want string
	}{
		{
			name: "VERIFIED",
			e:    Verified,
			want: "VERIFIED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("BookingCodeStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingCodeStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    BookingCodeStatus
		want bool
	}{
		{
			name: "valid type",
			e:    Verified,
			want: true,
		},
		{
			name: "invalid type",
			e:    BookingCodeStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("BookingCodeStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingCodeStatus_UnmarshalGQL(t *testing.T) {
	value := Verified
	invalid := BookingCodeStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *BookingCodeStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "VERIFIED",
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
				t.Errorf("BookingCodeStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingCodeStatus_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     BookingCodeStatus
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     Verified,
			b:     w,
			wantW: strconv.Quote("VERIFIED"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("BookingCodeStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestBookingStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    BookingStatus
		want string
	}{
		{
			name: "PENDING",
			e:    Pending,
			want: "PENDING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("BookingStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    BookingStatus
		want bool
	}{
		{
			name: "valid type",
			e:    Pending,
			want: true,
		},
		{
			name: "invalid type",
			e:    BookingStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("BookingStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookingStatus_UnmarshalGQL(t *testing.T) {
	value := Pending
	invalid := BookingStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *BookingStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "PENDING",
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
				t.Errorf("BookingStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingStatus_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     BookingStatus
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     Pending,
			b:     w,
			wantW: strconv.Quote("PENDING"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("BookingStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
