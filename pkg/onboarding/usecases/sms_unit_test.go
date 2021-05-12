package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestSMSImpl_CreateSMSData(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize the onboarding interactor: %v", err)
		return
	}

	validLinkId := uuid.New().String()
	text := "Test Covers"
	to := "3601"
	id := "60119"
	from := "+254705385894"
	date := "2021-05-17T13:20:04.490Z"

	validData := &resources.AfricasTalkingMessage{
		LinkID: validLinkId,
		Text:   text,
		To:     to,
		ID:     id,
		Date:   date,
		From:   from,
	}

	invalidData := &resources.AfricasTalkingMessage{
		LinkID: " ",
		Text:   text,
		To:     to,
		ID:     id,
		Date:   date,
		From:   from,
	}

	type args struct {
		ctx   context.Context
		input *resources.AfricasTalkingMessage
	}
	tests := []struct {
		name    string
		args    args
		want    *resources.AfricasTalkingMessage
		wantErr bool
	}{
		{
			name: "Happy:) successfully persist sms message data",
			args: args{
				ctx:   ctx,
				input: validData,
			},
			wantErr: false,
		},
		{
			name: "Sad:( fail to persist sms message data",
			args: args{
				ctx:   ctx,
				input: invalidData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy:) successfully persist sms message data" {
				fakeRepo.PersistIncomingSMSDataFn = func(ctx context.Context, input *resources.AfricasTalkingMessage) error {
					return nil
				}
			}

			if tt.name == "Sad:( fail to persist sms message data" {
				fakeRepo.PersistIncomingSMSDataFn = func(ctx context.Context, input *resources.AfricasTalkingMessage) error {
					return fmt.Errorf("unable to persist sms message data")
				}
			}

			err := i.SMS.CreateSMSData(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SMSImpl.CreateSMSData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected. got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected. got %v", err)
					return
				}
			}
		})
	}
}
