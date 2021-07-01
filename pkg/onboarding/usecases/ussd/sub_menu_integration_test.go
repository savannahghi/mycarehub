package ussd_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_HandleHomeMenu(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	type args struct {
		ctx          context.Context
		level        int
		session      *domain.USSDLeadDetails
		userResponse string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.AITUSSD.HandleHomeMenu(tt.args.ctx, tt.args.level, tt.args.session, tt.args.userResponse); got != tt.want {
				t.Errorf("Impl.HandleHomeMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}
