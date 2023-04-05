package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestServiceRequestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestStatus
		want string
	}{
		{
			name: "PENDING",
			e:    ServiceRequestStatusPending,
			want: "PENDING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ServiceRequestStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestStatus
		want bool
	}{
		{
			name: "valid type",
			e:    ServiceRequestStatusInProgress,
			want: true,
		},
		{
			name: "invalid type",
			e:    ServiceRequestStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ServiceRequestStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestStatus_UnmarshalGQL(t *testing.T) {
	value := ServiceRequestStatusResolved
	invalid := ServiceRequestStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ServiceRequestStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "RESOLVED",
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
				t.Errorf("ServiceRequestStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceRequestStatus_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     ServiceRequestStatus
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     ServiceRequestStatusResolved,
			b:     w,
			wantW: strconv.Quote("RESOLVED"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ServiceRequestStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPINResetVerificationStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    PINResetVerificationStatus
		want bool
	}{
		{
			name: "Happy Case - Valid",
			e:    PINResetVerificationStatusApproved,
			want: true,
		},
		{
			name: "Sad Case - Invalid State",
			e:    PINResetVerificationStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PINResetVerificationStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPINResetVerificationStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    PINResetVerificationStatus
		want string
	}{
		{
			name: "approved",
			e:    PINResetVerificationStatusApproved,
			want: "APPROVED",
		},
		{
			name: "rejected",
			e:    PINResetVerificationStatusRejected,
			want: "REJECTED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PINResetVerificationStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPINResetVerificationStatus_UnmarshalGQL(t *testing.T) {
	valid := PINResetVerificationStatusApproved
	invalid := PINResetVerificationStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *PINResetVerificationStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "APPROVED",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not a real kyc process status",
			},
			wantErr: true,
		},
		{
			name: "non string",
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
				t.Errorf("PINResetVerificationStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
