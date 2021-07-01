package ussd_test

import (
	"context"
	"reflect"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_UpdateSessionLevel(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	type args struct {
		ctx       context.Context
		level     int
		sessionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := u.AITUSSD.UpdateSessionLevel(tt.args.ctx, tt.args.level, tt.args.sessionID); (err != nil) != tt.wantErr {
				t.Errorf("Impl.UpdateSessionLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_AddAITSessionDetails(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	type args struct {
		ctx   context.Context
		input *dto.SessionDetails
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.AITUSSD.AddAITSessionDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.AddAITSessionDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.AddAITSessionDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_GetOrCreateSessionState(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	type args struct {
		ctx     context.Context
		payload *dto.SessionDetails
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.AITUSSD.GetOrCreateSessionState(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetOrCreateSessionState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.GetOrCreateSessionState() = %v, want %v", got, tt.want)
			}
		})
	}
}
