package domain_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
)

func TestKYCProcessStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.KYCProcessStatus
		want string
	}{
		{
			name: "approved",
			e:    domain.KYCProcessStatusApproved,
			want: "APPROVED",
		},
		{
			name: "rejected",
			e:    domain.KYCProcessStatusRejected,
			want: "REJECTED",
		},
		{
			name: "pending",
			e:    domain.KYCProcessStatusPending,
			want: "PENDING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("KYCProcessStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKYCProcessStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.KYCProcessStatus
		want bool
	}{
		{
			name: "valid approved",
			e:    domain.KYCProcessStatusApproved,
			want: true,
		},
		{
			name: "invalid kyc process status",
			e:    domain.KYCProcessStatus("this is not a real kyc process status"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("KYCProcessStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGender_UnmarshalGQL(t *testing.T) {
	valid := domain.KYCProcessStatusApproved
	invalid := domain.KYCProcessStatus("this is not a real kyc process status")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.KYCProcessStatus
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
				t.Errorf("KYCProcessStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGender_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.KYCProcessStatus
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.KYCProcessStatusPending,
			wantW: strconv.Quote("PENDING"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("KYCProcessStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFivePointRating_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.FivePointRating
		want string
	}{
		{
			name: "poor",
			e:    domain.FivePointRatingPoor,
			want: "POOR",
		},
		{
			name: "unsatisfactory",
			e:    domain.FivePointRatingUnsatisfactory,
			want: "UNSATISFACTORY",
		},
		{
			name: "average",
			e:    domain.FivePointRatingAverage,
			want: "AVERAGE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("FivePointRating.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFivePointRating_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.FivePointRating
		want bool
	}{
		{
			name: "valid",
			e:    domain.FivePointRatingPoor,
			want: true,
		},
		{
			name: "invalid",
			e:    domain.FivePointRating("this is not a real rating"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("FivePointRating.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFivePointRating_UnmarshalGQL(t *testing.T) {
	valid := domain.FivePointRatingPoor
	invalid := domain.FivePointRating("this is not a real rating")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.FivePointRating
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "SATISFACTORY",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not a real five points rating",
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
				t.Errorf("FivePointRating.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFivePointRating_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.FivePointRating
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.FivePointRatingAverage,
			wantW: strconv.Quote("AVERAGE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FivePointRating.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestBeneficiaryRelationship_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.BeneficiaryRelationship
		want string
	}{
		{
			name: "spouse",
			e:    domain.BeneficiaryRelationshipSpouse,
			want: "SPOUSE",
		},
		{
			name: "child",
			e:    domain.BeneficiaryRelationshipChild,
			want: "CHILD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeneficiaryRelationship_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.BeneficiaryRelationship
		want bool
	}{
		{
			name: "valid",
			e:    domain.BeneficiaryRelationshipSpouse,
			want: true,
		},
		{
			name: "invalid",
			e:    domain.BeneficiaryRelationship("this is not real"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeneficiaryRelationship_UnmarshalGQL(t *testing.T) {
	valid := domain.BeneficiaryRelationshipSpouse
	invalid := domain.BeneficiaryRelationship("this is not real")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.BeneficiaryRelationship
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "SPOUSE",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not real",
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
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBeneficiaryRelationship_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.BeneficiaryRelationship
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.BeneficiaryRelationshipSpouse,
			wantW: strconv.Quote("SPOUSE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestOrganizationType_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.OrganizationType
		want string
	}{
		{
			name: "LIMITED_COMPANY",
			e:    domain.OrganizationTypeLimitedCompany,
			want: "LIMITED_COMPANY",
		},
		{
			name: "TRUST",
			e:    domain.OrganizationTypeTrust,
			want: "TRUST",
		},
		{
			name: "UNIVERSITY",
			e:    domain.OrganizationTypeUniversity,
			want: "UNIVERSITY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganizationType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.OrganizationType
		want bool
	}{
		{
			name: "valid",
			e:    domain.OrganizationTypeLimitedCompany,
			want: true,
		},
		{
			name: "invalid",
			e:    domain.OrganizationType("this is not real"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrganizationType_UnmarshalGQL(t *testing.T) {
	valid := domain.OrganizationTypeUniversity
	invalid := domain.OrganizationType("this is not real")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.OrganizationType
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "UNIVERSITY",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not real",
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
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrganizationType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.OrganizationType
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.OrganizationTypeTrust,
			wantW: strconv.Quote("TRUST"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPractitionerService_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.PractitionerService
		want string
	}{
		{
			name: "INPATIENT_SERVICES",
			e:    domain.PractitionerServiceInpatientServices,
			want: "INPATIENT_SERVICES",
		},
		{
			name: "PHARMACY",
			e:    domain.PractitionerServicePharmacy,
			want: "PHARMACY",
		},
		{
			name: "MATERNITY",
			e:    domain.PractitionerServiceMaternity,
			want: "MATERNITY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerService_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.PractitionerService
		want bool
	}{
		{
			name: "valid",
			e:    domain.PractitionerServiceOutpatientServices,
			want: true,
		},
		{
			name: "invalid",
			e:    domain.PractitionerService("this is not real"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerService_UnmarshalGQL(t *testing.T) {
	valid := domain.PractitionerServiceOutpatientServices
	invalid := domain.PractitionerService("this is not real")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.PractitionerService
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "OUTPATIENT_SERVICES",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not real",
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
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPractitionerService_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.PractitionerService
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.PractitionerServiceOutpatientServices,
			wantW: strconv.Quote("OUTPATIENT_SERVICES"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestAgentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    domain.AgentType
		want bool
	}{
		{
			name: "valid",
			e:    domain.CompanyAgent,
			want: true,
		},
		{
			name: "invalid",
			e:    domain.AgentType("FreelanceAgent"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AgentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentType_String(t *testing.T) {
	tests := []struct {
		name string
		e    domain.AgentType
		want string
	}{
		{
			name: "valid",
			e:    domain.CompanyAgent,
			want: "SIL Agent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AgentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentType_UnmarshalGQL(t *testing.T) {
	valid := domain.CompanyAgent
	invalid := domain.AgentType("CompanyAgent")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *domain.AgentType
		args    args
		wantErr bool
	}{
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid string",
			e:    &invalid,
			args: args{
				v: "CompanyAgent",
			},
			wantErr: true,
		},
		{
			name: "valid type",
			e:    &valid,
			args: args{
				v: "SIL Agent",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AgentType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgentType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     domain.AgentType
		wantW string
	}{
		{
			name:  "valid",
			e:     domain.CompanyAgent,
			wantW: strconv.Quote("SIL Agent"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AgentType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestChannel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    domain.Channel
		want bool
	}{
		{
			name: "valid",
			c:    domain.WhatsAppChannel,
			want: true,
		},
		{
			name: "invalid",
			c:    domain.Channel("this is not real"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("Channel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		c    domain.Channel
		want string
	}{
		{
			name: "whatsapp",
			c:    domain.WhatsAppChannel,
			want: "WhatsApp",
		},
		{
			name: "sms",
			c:    domain.SMSChannel,
			want: "SMS",
		},
		{
			name: "email",
			c:    domain.EmailChannel,
			want: "Email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_Int(t *testing.T) {
	tests := []struct {
		name string
		c    domain.Channel
		want int
	}{
		{
			name: "whatsapp",
			c:    domain.WhatsAppChannel,
			want: 1,
		},
		{
			name: "sms",
			c:    domain.SMSChannel,
			want: 2,
		},
		{
			name: "email",
			c:    domain.EmailChannel,
			want: 3,
		},
		{
			name: "test",
			c:    domain.Channel("test"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Int(); got != tt.want {
				t.Errorf("Channel.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_UnmarshalGQL(t *testing.T) {
	valid := domain.WhatsAppChannel
	invalid := domain.Channel("test")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *domain.Channel
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			c:    &valid,
			args: args{
				v: "WhatsApp",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			c:    &invalid,
			args: args{
				v: "test",
			},
			wantErr: true,
		},
		{
			name: "non string",
			c:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Channel.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannel_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     domain.Channel
		wantW string
	}{
		{
			name:  "valid",
			c:     domain.WhatsAppChannel,
			wantW: strconv.Quote("WhatsApp"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Channel.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
