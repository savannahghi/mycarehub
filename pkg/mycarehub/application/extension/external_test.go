package extension

import (
	"context"
	"testing"

	"github.com/segmentio/ksuid"
)

func TestExternal_CreateFirebaseCustomToken(t *testing.T) {
	ctx := context.Background()
	uid := ksuid.New().String()
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx: ctx,
				uid: uid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewExternalMethodsImpl()
			got, err := s.CreateFirebaseCustomToken(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.CreateFirebaseCustomToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected to get a response but got: %v", got)
				return
			}
		})
	}
}

func TestExternal_SendFeedback(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx             context.Context
		subject         string
		feedbackMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:             ctx,
				subject:         "test",
				feedbackMessage: "test",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewExternalMethodsImpl()
			got, err := s.SendFeedback(tt.args.ctx, tt.args.subject, tt.args.feedbackMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("External.SendFeedback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("External.SendFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}
