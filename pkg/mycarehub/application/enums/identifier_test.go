package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestIdentifierType_String(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierType
		want string
	}{
		{
			name: "CCC",
			e:    IdentifierTypeCCC,
			want: "CCC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IdentifierType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierType
		want bool
	}{
		{
			name: "valid type",
			e:    IdentifierTypeCCC,
			want: true,
		},
		{
			name: "invalid type",
			e:    IdentifierType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IdentifierType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierType_UnmarshalGQL(t *testing.T) {
	value := IdentifierTypeCCC
	invalid := IdentifierType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IdentifierType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CCC",
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
				t.Errorf("IdentifierType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdentifierType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     IdentifierType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     IdentifierTypeCCC,
			b:     w,
			wantW: strconv.Quote("CCC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IdentifierType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestIdentifierUse_String(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierUse
		want string
	}{
		{
			name: "OFFICIAL",
			e:    IdentifierUseOfficial,
			want: "OFFICIAL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("IdentifierUse.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierUse_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    IdentifierUse
		want bool
	}{
		{
			name: "valid type",
			e:    IdentifierUseOfficial,
			want: true,
		},
		{
			name: "invalid type",
			e:    IdentifierUse("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IdentifierUse.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifierUse_UnmarshalGQL(t *testing.T) {
	value := IdentifierUseOfficial
	invalid := IdentifierUse("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *IdentifierUse
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "OFFICIAL",
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
				t.Errorf("IdentifierUse.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdentifierUse_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     IdentifierUse
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     IdentifierUseOfficial,
			b:     w,
			wantW: strconv.Quote("OFFICIAL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IdentifierUse.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
