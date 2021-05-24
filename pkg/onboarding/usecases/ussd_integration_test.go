package usecases_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestCreateUSSDData(t *testing.T) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	phone := "+254721026491"
	invalidPhone := ""
	validInput := dto.EndSessionDetails{
		SessionID:   ":12345678",
		PhoneNumber: &phone,
		Input:       "1",
	}

	invalidInput := dto.EndSessionDetails{
		SessionID:   "",
		PhoneNumber: &invalidPhone,
	}
	type args struct {
		ctx   context.Context
		input dto.EndSessionDetails
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:) successfully add USSD Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "sad:( unsuccessfully add USSD Details",
			args: args{
				ctx:   context.Background(),
				input: invalidInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.AITUSSD.CreateUSSDData(tt.args.ctx, &tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("USSDUseCaseImpl.AddUSSDDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
