package utils

import (
	"context"
	"testing"
)

func TestGetLOggedInUserID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "Sad case: failed to get logged in user id",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLoggedInUserID(tt.args.ctx); got != tt.want {
				t.Errorf("GetLoggedInUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
