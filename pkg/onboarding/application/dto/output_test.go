package dto_test

import (
	"reflect"
	"testing"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

func TestNewOKResp(t *testing.T) {
	type args struct {
		rawResponse interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      *dto.OKResp
		wantError bool
	}{
		{
			name: "Happy",
			args: args{
				rawResponse: "hello",
			},
			want: &dto.OKResp{
				Status:   "OK",
				Response: "hello",
			},
			wantError: false,
		},
		{
			name: "Sad case",
			args: args{
				rawResponse: "hello",
			},
			want: &dto.OKResp{
				Status:   "badstatus",
				Response: "hello",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dto.NewOKResp(tt.args.rawResponse); !reflect.DeepEqual(got, tt.want) && !tt.wantError {
				t.Errorf("NewOKResp() = %v, want %v", got, tt.want)
			}
		})
	}
}
