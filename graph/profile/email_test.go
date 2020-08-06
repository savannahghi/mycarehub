package profile

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func Test_isTester(t *testing.T) {
	validTesterEmail := gofakeit.Email()
	srv := NewService()
	ctx := context.Background()
	added, err := srv.AddTester(ctx, validTesterEmail)
	assert.Nil(t, err)
	assert.True(t, added)

	type args struct {
		emails []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "exists",
			args: args{
				emails: []string{validTesterEmail},
			},
			want: true,
		},
		{
			name: "Apple special case",
			args: args{
				emails: []string{"jobs@apple.com"},
			},
			want: true,
		},
		{
			name: "does not exist",
			args: args{
				emails: []string{gofakeit.Email()},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTester(ctx, tt.args.emails); got != tt.want {
				t.Errorf("isTester() = %v, want %v", got, tt.want)
			}
		})
	}
}
