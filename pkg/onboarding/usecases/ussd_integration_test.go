package usecases_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestGenerateUSSD(t *testing.T) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	phone := "+254721026491"
	invalidPhone := ""
	validInput := dto.SessionDetails{
		SessionID:   ":12345678",
		PhoneNumber: &phone,
	}

	invalidInput := dto.SessionDetails{
		SessionID:   "",
		PhoneNumber: &invalidPhone,
	}
	type args struct {
		ctx   context.Context
		input dto.SessionDetails
	}
	tests := []struct {
		name string
		args args
		resp string
	}{
		{
			name: "happy:) successfully add USSD Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			resp: "CON Welcome to Be.Well \n1. Register",
		},
		{
			name: "sad:( unsuccessfully add USSD Details",
			args: args{
				ctx:   context.Background(),
				input: invalidInput,
			},
			resp: "2: unable to normalize the msisdn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := s.AITUSSD.GenerateUSSD(tt.args.ctx, &tt.args.input)
			if resp != tt.resp {
				t.Errorf("USSDUseCaseImpl.AddUSSDDetails() resp = %v, want %v", resp, tt.resp)
				return
			}
		})
	}
}
